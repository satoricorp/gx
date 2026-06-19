package hooks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/orchestrator"
)

// PushOptions configures a hook-triggered capture run.
type PushOptions struct {
	RepoRoot string
	Remote   string
	RefRange string
	Base     string
	Head     string
	HomeDir  string
	Tools    []string
}

// RunPush executes the capture pipeline for a git push ref range.
func RunPush(ctx context.Context, opts PushOptions) (orchestrator.Result, error) {
	homeDir := opts.HomeDir
	if homeDir == "" {
		homeDir, _ = os.UserHomeDir()
	}
	if capture.IsPaused(homeDir) {
		return orchestrator.Result{}, nil
	}

	repoRoot := opts.RepoRoot
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return orchestrator.Result{}, err
		}
		repoRoot = cwd
	}
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return orchestrator.Result{}, err
	}

	base, head := parseRefRange(opts.RefRange, opts.Base, opts.Head)
	tools := opts.Tools
	if len(tools) == 0 {
		tools = []string{capture.ToolClaude, capture.ToolCodex, capture.ToolCursor}
	}

	return orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot:  repoRoot,
		Base:      base,
		Head:      head,
		HomeDir:   opts.HomeDir,
		Tools:     tools,
		StageOnly: !auth.HasUploadCredentials(),
	})
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
