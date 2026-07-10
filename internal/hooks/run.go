package hooks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/background"
	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/orchestrator"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

// PushOptions configures a hook-triggered capture run.
type PushOptions struct {
	RepoRoot string
	Remote   string
	RefRange string
	Base     string
	Head     string
	LocalRef string
	HeadSHA  string
	HomeDir  string
	Tools    []string
	SkipCapture bool
}

// PushOutcome summarizes a non-blocking pre-push hook run.
type PushOutcome struct {
	Result           orchestrator.Result
	RevisionIDs      []string
	ShareableExtract int
	ShareableSession int
	Publication      publication.Result
}

// RunPush executes capture staging, shareable marking, and background upload kickoff.
// It always returns a nil error so git pre-push hooks never block pushes on network.
func RunPush(ctx context.Context, opts PushOptions) (PushOutcome, error) {
	outcome := PushOutcome{}
	homeDir := opts.HomeDir
	if homeDir == "" {
		homeDir, _ = os.UserHomeDir()
	}
	if capture.IsPaused(homeDir) {
		return outcome, nil
	}

	repoRoot := opts.RepoRoot
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return outcome, nil
		}
		repoRoot = cwd
	}
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return outcome, nil
	}

	base, head := parseRefRange(opts.RefRange, opts.Base, opts.Head)
	tools := opts.Tools
	if len(tools) == 0 {
		tools = []string{capture.ToolClaude, capture.ToolCodex, capture.ToolCursor}
	}

	if !opts.SkipCapture {
		result, err := orchestrator.Run(ctx, orchestrator.RunOptions{
			RepoRoot:  repoRoot,
			Base:      base,
			Head:      head,
			HomeDir:   opts.HomeDir,
			Tools:     tools,
			StageOnly: true,
		})
		if err == nil {
			outcome.Result = result
		}
	}

	refRange := strings.TrimSpace(opts.RefRange)
	if refRange == "" {
		refRange = strings.TrimSpace(outcome.Result.RefRange)
	}
	if refRange == "" && base != "" && head != "" {
		refRange = base + ".." + head
	}
	revisionIDs, err := vcs.RevisionIDsInGitRange(ctx, repoRoot, refRange)
	if err != nil {
		return outcome, nil
	}
	outcome.RevisionIDs = revisionIDs

	att, _ := gitCaptureAttestation(ctx, repoRoot)
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		return outcome, nil
	}
	if n, err := stager.MarkExtractShareable(ctx, revisionIDs, att); err == nil {
		outcome.ShareableExtract = n
	}
	if n, err := stager.MarkSessionShareable(ctx, revisionIDs, att); err == nil {
		outcome.ShareableSession = n
	}

	headSHA := strings.TrimSpace(opts.HeadSHA)
	if headSHA == "" && head != "" {
		headSHA, _ = gitRevParse(ctx, repoRoot, head)
	}
	pub, err := EnqueueAdoptedPublication(ctx, AdoptPushOptions{
		RepoRoot:    repoRoot,
		Remote:      opts.Remote,
		LocalRef:    opts.LocalRef,
		HeadSHA:     headSHA,
		RevisionIDs: revisionIDs,
	})
	if err == nil {
		outcome.Publication = pub
	}

	if backgroundWorkersEnabled() {
		_ = background.StartDetachedGX("capture", "sync", "--quiet")
		if pub.Queued {
			_ = background.StartDetachedGX("__gx-upload-outbox", "--quiet")
		}
	}
	return outcome, nil
}

func gitCaptureAttestation(ctx context.Context, repoRoot string) (storage.CaptureAttestation, error) {
	name, err := gitConfigValue(ctx, repoRoot, "user.name")
	if err != nil {
		return storage.CaptureAttestation{}, err
	}
	email, err := gitConfigValue(ctx, repoRoot, "user.email")
	if err != nil {
		return storage.CaptureAttestation{}, err
	}
	return storage.CaptureAttestation{Name: name, Email: email}, nil
}

func gitConfigValue(ctx context.Context, repoRoot, key string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "config", key)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git config %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func parseRefRange(refRange, base, head string) (string, string) {
	refRange = strings.TrimSpace(refRange)
	if base != "" || head != "" {
		if head == "" {
			head = "HEAD"
		}
		return base, head
	}
	if refRange == "" {
		return "", "HEAD"
	}
	if strings.Contains(refRange, "..") {
		parts := strings.SplitN(refRange, "..", 2)
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	// Single SHA from a new pushed ref. Let the orchestrator choose the best
	// available base instead of assuming "<sha>~20" exists.
	return "", refRange
}

// ParsePrePushLine parses one stdin line from git pre-push.
func ParsePrePushLine(line string) (localRef, localSHA, remoteRef, remoteSHA string, err error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) != 4 {
		return "", "", "", "", fmt.Errorf("expected 4 fields, got %d", len(fields))
	}
	return fields[0], fields[1], fields[2], fields[3], nil
}

// RefRangeForPush computes the capture ref range for one push update.
func RefRangeForPush(localSHA, remoteSHA string) string {
	zero := "0000000000000000000000000000000000000000"
	if localSHA == zero {
		return ""
	}
	if remoteSHA == zero {
		return localSHA
	}
	return remoteSHA + ".." + localSHA
}

func backgroundWorkersEnabled() bool {
	return os.Getenv("GX_DISABLE_BACKGROUND_WORKERS") == ""
}