package vcs

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var revisionTrailerLinePattern = regexp.MustCompile(
	`(?m)^` + regexp.QuoteMeta(strings.TrimSuffix(RevisionTrailerFormat, "%s")) + `(\S+)\s*$`,
)

// ParseRevisionIDsFromMessage returns lgtm revision IDs embedded in a commit message.
func ParseRevisionIDsFromMessage(message string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, match := range revisionTrailerLinePattern.FindAllStringSubmatch(message, -1) {
		if len(match) < 2 {
			continue
		}
		id := strings.TrimSpace(match[1])
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// RevisionTrailerValue formats the trailer value for a revision ID.
func RevisionTrailerValue(revisionID string) string {
	return RevisionTrailerLine(revisionID)
}

// RevisionIDsInGitRange lists unique lgtm revision IDs from commit messages in a ref range.
func RevisionIDsInGitRange(ctx context.Context, repoRoot, refRange string) ([]string, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	refRange = strings.TrimSpace(refRange)
	if repoRoot == "" {
		return nil, fmt.Errorf("repo root required")
	}
	if refRange == "" {
		return nil, nil
	}
	args := []string{"-C", repoRoot, "log", "--format=%B", refRange}
	if !strings.Contains(refRange, "..") {
		// A bare rev means a new-branch push with no remote base; without a
		// bound the walk collects every lgtm revision in repo history. Only
		// commits that are not already on a remote are being pushed.
		args = append(args, "--not", "--remotes")
	}
	cmd := exec.CommandContext(ctx, "git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git log %s: %w\n%s", refRange, err, strings.TrimSpace(string(out)))
	}
	seen := map[string]struct{}{}
	var ids []string
	for _, message := range strings.Split(string(out), "\n\n") {
		for _, id := range ParseRevisionIDsFromMessage(message) {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids, nil
}
