package hooks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/background"
	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/capture/orchestrator"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

// SuppressAdoptedPublicationEnv tells the pre-push hook that gx push owns
// this push and will enqueue the publication itself.
const SuppressAdoptedPublicationEnv = "GX_SUPPRESS_ADOPTED_PUBLICATION"

// runCapture is a seam for tests that need one specific orchestrator outcome —
// above all a partial failure, which returns a populated Result *and* an error
// and is otherwise hard to produce on demand.
var runCapture = orchestrator.Run

// markCaptureShareable and revisionIDsInGitRange are seams for tests. Both
// failure modes are real — SQLite busy against the detached sync writer, a git
// invocation that fails mid-push — and both used to abort the rest of the hook,
// which is precisely what needs covering; neither is reachable on demand
// without racing SQLite or corrupting a repo mid-run.
var markCaptureShareable = markShareable

var revisionIDsInGitRange = vcs.RevisionIDsInGitRange

// PushOptions configures a hook-triggered capture run.
type PushOptions struct {
	RepoRoot    string
	Remote      string
	RefRange    string
	Base        string
	Head        string
	LocalRef    string
	HeadSHA     string
	HomeDir     string
	Tools       []string
	SkipCapture bool
}

// PushOutcome summarizes a non-blocking pre-push hook run.
//
// SkipReason is set whenever the hook deliberately did no work. Without it a
// paused capture, a repo opted out with `git config gx.enabled false`, and a
// failure to resolve the working directory all returned the same zero outcome,
// so all three rendered as the identical `capture staged extract= sessions=0`
// line — the same line a genuine error prints. An empty extract id has to mean
// one thing.
//
// UploadFailures reports rows whose last upload attempt failed. Uploads run in
// a detached process with both streams sent to /dev/null, so the push that
// starts one can never report its result; the push after it reports what the
// previous one left behind.
//
// AttachError reports a failure to link this push's sessions to their change
// rows. It is never returned: a push must not fail because SQLite was busy.
// But it must not be silent either — an empty sessions[] in the published
// artifact is exactly the symptom that hid the commit-time matcher's
// self-referential gate for 174 bundles.
type PushOutcome struct {
	Result           orchestrator.Result
	RevisionIDs      []string
	SkipReason       string
	CaptureError     string
	RecoveryError    string
	ShareableError   string
	AttachError      string
	PublicationError string
	ShareableExtract int
	ShareableSession int
	AttachedSessions int
	UploadFailures   storage.CaptureUploadFailures
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
		outcome.SkipReason = "capture is paused (GX_CAPTURE_PAUSED or ~/.gx/pause-capture)"
		return outcome, nil
	}

	repoRoot := opts.RepoRoot
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			// Not an opt-out — something is wrong. It is named as a failure so
			// it cannot hide behind a reason a user would read as intentional.
			outcome.SkipReason = fmt.Sprintf("resolve working directory: %v", err)
			outcome.CaptureError = outcome.SkipReason
			return outcome, nil
		}
		repoRoot = cwd
	}
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		outcome.SkipReason = fmt.Sprintf("resolve repo path %q: %v", repoRoot, err)
		outcome.CaptureError = outcome.SkipReason
		return outcome, nil
	}
	// Per-repo opt-out: `git config gx.enabled false` excludes one repository
	// from GX, including from a machine-wide `gx init --global` install.
	if !EnabledForRepo(ctx, repoRoot) {
		outcome.SkipReason = "this repository opted out (git config gx.enabled false)"
		return outcome, nil
	}

	base, head := parseRefRange(opts.RefRange, opts.Base, opts.Head)
	tools := opts.Tools
	if len(tools) == 0 {
		tools = []string{capture.ToolClaude, capture.ToolCodex, capture.ToolCursor}
	}

	if !opts.SkipCapture {
		result, err := runCapture(ctx, orchestrator.RunOptions{
			RepoRoot:  repoRoot,
			Base:      base,
			Head:      head,
			HomeDir:   opts.HomeDir,
			Tools:     tools,
			StageOnly: true,
		})
		// The orchestrator returns a populated Result alongside its error: a run
		// that staged the extract and then failed on one session has produced
		// real, uploadable rows. Keeping the Result on both branches is what
		// lets the steps below still mark those rows shareable and put them in
		// the publish bundle. Dropping it made a partial failure look like a
		// total no-op, which is the "silent staging" symptom.
		outcome.Result = result
		if err != nil {
			// The hook must never block a push, but a failed capture run must
			// not read as "captured nothing" either — that hid a staging bug
			// behind a plausible sessions=0 for a long time.
			outcome.CaptureError = err.Error()
		}
	}

	refRange := strings.TrimSpace(opts.RefRange)
	if refRange == "" {
		refRange = strings.TrimSpace(outcome.Result.RefRange)
	}
	if refRange == "" && base != "" && head != "" {
		refRange = base + ".." + head
	}
	revisionIDs, err := revisionIDsInGitRange(ctx, repoRoot, refRange)
	if err != nil {
		// Revisions are a backlog convenience, not a prerequisite: the rows
		// this run staged are addressed by id. Returning here also skipped the
		// publication and the upload kickoff, so a trailer-scan hiccup left the
		// PR with no GX artifact at all and nothing said so. Carry on with no
		// revisions and report the scan failure.
		outcome.RecoveryError = fmt.Sprintf("scan GX revision trailers: %v", err)
		revisionIDs = nil
	}
	outcome.RevisionIDs = revisionIDs

	if len(revisionIDs) > 0 {
		if repo, repoErr := vcs.NewService().ResolveGXRepoAtPath(ctx, repoRoot); repoErr == nil {
			if _, recoverErr := vcs.NewService().RecoverMissingRevisions(ctx, repo, revisionIDs); recoverErr != nil {
				outcome.RecoveryError = recoverErr.Error()
			}
		} else {
			outcome.RecoveryError = repoErr.Error()
		}
	}

	// Placement is load-bearing. RecoverMissingRevisions above guarantees the
	// `changes` rows this links to exist, and EnqueueAdoptedPublication below
	// runs reviewbundle.BuildPush synchronously in this same process — so the
	// artifact this very push writes already carries the sessions. Running it
	// before the SuppressAdoptedPublication early return covers `gx push`,
	// which builds its own artifact from the same rows.
	if attached, err := attachPushSessions(ctx, repoRoot, outcome.Result.HunkLinks); err != nil {
		outcome.AttachError = err.Error()
	} else {
		outcome.AttachedSessions = attached
	}

	if err := markCaptureShareable(ctx, repoRoot, &outcome, revisionIDs); err != nil {
		// SQLite is WAL with a 5s busy timeout and a detached sync from the
		// previous push writes the same tables, so this genuinely happens.
		// Aborting here skipped the publication and the upload kickoff too and
		// rendered as `shareable=0/0` with no warning: the exact silent
		// staging this marking exists to eliminate.
		outcome.ShareableError = err.Error()
	}

	// Read after marking and before kicking off the next sync, so what is
	// reported is what the previous detached sync actually left behind.
	if failures, err := storage.CaptureUploadFailureCounts(ctx); err == nil {
		outcome.UploadFailures = failures
	}

	if os.Getenv(SuppressAdoptedPublicationEnv) != "" {
		// gx push drives this git push and enqueues its own artifact with
		// the PR URL attached; a hook publication here would race it with
		// a PR-less artifact for the same head.
		if backgroundWorkersEnabled() {
			_ = background.StartDetachedGX("capture", "sync", "--quiet")
		}
		return outcome, nil
	}

	headSHA := strings.TrimSpace(opts.HeadSHA)
	if headSHA == "" && head != "" {
		headSHA, _ = gitRevParse(ctx, repoRoot, head)
	}
	pub, err := EnqueueAdoptedPublication(ctx, AdoptPushOptions{
		RepoRoot: repoRoot,
		Remote:   opts.Remote,
		LocalRef: opts.LocalRef,
		HeadSHA:  headSHA,
		RefRange: refRange,
	})
	if err == nil {
		outcome.Publication = pub
	} else {
		// Without this the PR simply never grows a GX summary and nothing
		// anywhere says why.
		outcome.PublicationError = err.Error()
	}

	if backgroundWorkersEnabled() {
		_ = background.StartDetachedGX("capture", "sync", "--quiet")
		if shouldStartOutboxWorker(pub.Queued) {
			_ = background.StartDetachedGX("__gx-upload-outbox", "--quiet")
		}
	}
	return outcome, nil
}

// attachPushSessions links the sessions this push's hunk links name to the
// change rows they edited, so the artifact built moments later carries them.
func attachPushSessions(ctx context.Context, repoRoot string, links []matcher.HunkLink) (int, error) {
	if len(links) == 0 {
		return 0, nil
	}
	service := vcs.NewService()
	repo, err := service.ResolveGXRepoAtPath(ctx, repoRoot)
	if err != nil {
		return 0, err
	}
	return service.AttachSessionsFromHunkLinks(ctx, repo, links)
}

// markShareable attests the capture rows this push may upload.
//
// The staged row ids are the load-bearing gate: they name the exact rows this
// run wrote, with no dependence on resolving a revision. Revision matching runs
// afterwards for the backlog — rows staged by an earlier run whose revisions
// only now reach the remote — and still works, since a repo with the GX
// prepare-commit-msg hook stamps a trailer on plain `git commit`. Rows the id
// pass already marked are excluded from the second pass's count.
//
// Every marking failure is collected and returned. Discarding them made a
// failed marking pass indistinguishable from a successful one that had nothing
// to mark.
func markShareable(ctx context.Context, repoRoot string, outcome *PushOutcome, revisionIDs []string) error {
	att, _ := gitCaptureAttestation(ctx, repoRoot)
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		return err
	}
	var problems []error
	count := func(n int, err error) int {
		if err != nil {
			problems = append(problems, err)
			return 0
		}
		return n
	}
	outcome.ShareableExtract += count(stager.MarkExtractsShareableByID(ctx, outcome.Result.StagedExtractIDs, att))
	outcome.ShareableSession += count(stager.MarkSessionsShareableByID(ctx, outcome.Result.StagedSessionIDs, att))
	// Rows left behind by the older row-key scheme describe the very same
	// transcripts this run just staged, so they carry the same attestation
	// rather than sitting unshareable for the rest of time.
	outcome.ShareableSession += count(stager.MarkSessionsShareableForSourcesOf(ctx, outcome.Result.StagedSessionIDs, att))
	outcome.ShareableExtract += count(stager.MarkExtractShareable(ctx, revisionIDs, att))
	outcome.ShareableSession += count(stager.MarkSessionShareable(ctx, revisionIDs, att))
	return errors.Join(problems...)
}

// shouldStartOutboxWorker decides whether this push spawns the publish-outbox
// worker. It is not gated on the current push having enqueued something: with
// `gx sync` retired, the pre-push hook is the retry path for items that failed
// or were left pending by an earlier push, so any backlog also spawns the
// worker.
func shouldStartOutboxWorker(queuedThisPush bool) bool {
	if queuedThisPush {
		return true
	}
	status, err := publication.QueuedUploadStatus()
	if err != nil {
		return false
	}
	return status.Pending > 0 || status.Failed > 0
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
