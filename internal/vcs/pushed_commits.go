package vcs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const (
	// maxPushedCommits bounds how many commits one publish bundle describes.
	// A push of an entire imported history must not turn into a bundle with
	// thousands of patches; the newest commits are the ones under review, so
	// when the range is larger than this the oldest entries are dropped.
	maxPushedCommits = 200
	// maxPushedCommitPatchBytes bounds one commit's inline diff. Oversized
	// patches are omitted (files and message stay) rather than truncated, so a
	// partial diff can never be mistaken for the whole change downstream.
	maxPushedCommitPatchBytes = 512 << 10
)

// PushedCommitsInGitRange describes every commit in a pushed ref range,
// oldest first: SHA, full message, GX trailer revision ID (empty when the
// commit has none), changed files, and the unified diff.
func PushedCommitsInGitRange(ctx context.Context, repoRoot, refRange string) ([]PushedCommit, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	refRange = strings.TrimSpace(refRange)
	if repoRoot == "" {
		return nil, fmt.Errorf("repo root required")
	}
	if refRange == "" {
		return nil, nil
	}
	args := []string{"-C", repoRoot, "rev-list", "--reverse", refRange}
	if !strings.Contains(refRange, "..") {
		// A bare rev means a new-branch push with no remote base; without a
		// bound the walk collects the entire repo history. Mirror the trailer
		// scan: only commits that are not already on a remote are being pushed.
		args = append(args, "--not", "--remotes")
	}
	out, err := exec.CommandContext(ctx, "git", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git rev-list %s: %w\n%s", refRange, err, strings.TrimSpace(string(out)))
	}
	shas := splitLines(strings.TrimSpace(string(out)))
	if len(shas) > maxPushedCommits {
		shas = shas[len(shas)-maxPushedCommits:]
	}
	commits := make([]PushedCommit, 0, len(shas))
	for _, sha := range shas {
		sha = strings.TrimSpace(sha)
		if sha == "" {
			continue
		}
		commit, err := describePushedCommit(ctx, repoRoot, sha)
		if err != nil {
			return nil, err
		}
		commits = append(commits, commit)
	}
	return commits, nil
}

func describePushedCommit(ctx context.Context, repoRoot, sha string) (PushedCommit, error) {
	message, err := gitCombined(ctx, repoRoot, "log", "-1", "--format=%B", sha)
	if err != nil {
		return PushedCommit{}, fmt.Errorf("read commit message %s: %w", shortID(sha, 7), err)
	}
	message = strings.TrimSpace(message)
	revisionID := ""
	if ids := ParseRevisionIDsFromMessage(message); len(ids) > 0 {
		revisionID = ids[len(ids)-1]
	}
	filesOut, err := gitCombined(ctx, repoRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", sha)
	if err != nil {
		return PushedCommit{}, fmt.Errorf("list commit files %s: %w", shortID(sha, 7), err)
	}
	patch, err := gitCombined(ctx, repoRoot, "show", "--format=", "--patch", sha)
	if err != nil {
		return PushedCommit{}, fmt.Errorf("read commit patch %s: %w", shortID(sha, 7), err)
	}
	if len(patch) > maxPushedCommitPatchBytes {
		patch = ""
	}
	return PushedCommit{
		CommitID:   sha,
		RevisionID: revisionID,
		Message:    message,
		Files:      splitLines(strings.TrimSpace(filesOut)),
		Patch:      patch,
	}, nil
}

func gitCombined(ctx context.Context, repoRoot string, args ...string) (string, error) {
	fullArgs := append([]string{"-C", repoRoot}, args...)
	out, err := exec.CommandContext(ctx, "git", fullArgs...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
