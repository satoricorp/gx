package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/capture"
)

// HunkOptions configures commit hunk extraction.
type HunkOptions struct {
	RepoRoot string
	Base     string
	Head     string
}

// RevListCount returns commits reachable from ref.
func RevListCount(repoRoot, ref string) (int, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		ref = "HEAD"
	}
	if _, err := runGit(repoRoot, "rev-parse", "--verify", ref+"^{commit}"); err != nil {
		return 0, fmt.Errorf("invalid ref %q: %w", ref, err)
	}
	out, err := runGit(repoRoot, "rev-list", "--count", ref)
	if err != nil {
		return 0, err
	}
	var count int
	if _, err := fmt.Sscanf(strings.TrimSpace(out), "%d", &count); err != nil {
		return 0, err
	}
	return count, nil
}

// CommitCount returns the number of commits in base..head.
func CommitCount(opts HunkOptions) (int, error) {
	base, head, err := resolveRefs(opts)
	if err != nil {
		return 0, err
	}
	out, err := runGit(opts.RepoRoot, "rev-list", "--count", base+".."+head)
	if err != nil {
		return 0, err
	}
	var count int
	_, scanErr := fmt.Sscanf(strings.TrimSpace(out), "%d", &count)
	if scanErr != nil {
		return 0, scanErr
	}
	return count, nil
}

// CommitRange returns first/last commit timestamps in base..head (unix ms).
func CommitRange(opts HunkOptions) (first, last int64, err error) {
	base, head, err := resolveRefs(opts)
	if err != nil {
		return 0, 0, err
	}
	firstOut, err := runGit(opts.RepoRoot, "log", base+".."+head, "--reverse", "--format=%ct", "-1")
	if err != nil {
		return 0, 0, err
	}
	lastOut, err := runGit(opts.RepoRoot, "log", base+".."+head, "--format=%ct", "-1")
	if err != nil {
		return 0, 0, err
	}
	var firstSec, lastSec int64
	if _, err := fmt.Sscanf(strings.TrimSpace(firstOut), "%d", &firstSec); err != nil {
		return 0, 0, err
	}
	if _, err := fmt.Sscanf(strings.TrimSpace(lastOut), "%d", &lastSec); err != nil {
		return 0, 0, err
	}
	return firstSec * 1000, lastSec * 1000, nil
}

// CollectHunks extracts added-line hunks from git commits in base..head.
func CollectHunks(opts HunkOptions) ([]capture.CommitHunk, error) {
	base, head, err := resolveRefs(opts)
	if err != nil {
		return nil, err
	}
	shas, err := runGit(opts.RepoRoot, "rev-list", "--reverse", base+".."+head)
	if err != nil {
		return nil, err
	}
	shaList := splitLines(shas)
	var hunks []capture.CommitHunk
	for _, sha := range shaList {
		if sha == "" {
			continue
		}
		tsOut, err := runGit(opts.RepoRoot, "show", "-s", "--format=%ct", sha)
		if err != nil {
			return nil, err
		}
		var tsSec int64
		if _, err := fmt.Sscanf(strings.TrimSpace(tsOut), "%d", &tsSec); err != nil {
			return nil, err
		}
		diffOut, err := runGit(opts.RepoRoot, "show", "--pretty=format:", "--unified=0", sha)
		if err != nil {
			return nil, err
		}
		parsed := parseUnifiedDiff(sha, tsSec*1000, diffOut)
		hunks = append(hunks, parsed...)
	}
	return hunks, nil
}

func resolveRefs(opts HunkOptions) (base, head string, err error) {
	base = strings.TrimSpace(opts.Base)
	head = strings.TrimSpace(opts.Head)
	if head == "" {
		head = "HEAD"
	}
	if base == "" {
		base = "main"
	}
	if _, err := runGit(opts.RepoRoot, "rev-parse", "--verify", base+"^{commit}"); err != nil {
		return "", "", fmt.Errorf("invalid base ref %q: %w", base, err)
	}
	if _, err := runGit(opts.RepoRoot, "rev-parse", "--verify", head+"^{commit}"); err != nil {
		return "", "", fmt.Errorf("invalid head ref %q: %w", head, err)
	}
	return base, head, nil
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func parseUnifiedDiff(sha string, commitTime int64, diff string) []capture.CommitHunk {
	lines := strings.Split(diff, "\n")
	var hunks []capture.CommitHunk
	var currentFile string
	var added []string
	flush := func() {
		if currentFile == "" || len(added) == 0 {
			added = nil
			return
		}
		hunks = append(hunks, capture.CommitHunk{
			CommitSHA:  sha,
			FilePath:   filepath.ToSlash(currentFile),
			AddedLines: append([]string(nil), added...),
			CommitTime: commitTime,
		})
		added = nil
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "+++ b/") {
			flush()
			currentFile = strings.TrimPrefix(line, "+++ b/")
			continue
		}
		if strings.HasPrefix(line, "+++ ") && !strings.HasPrefix(line, "+++ /dev/null") {
			flush()
			currentFile = strings.TrimSpace(strings.TrimPrefix(line, "+++ "))
			if strings.HasPrefix(currentFile, "b/") {
				currentFile = strings.TrimPrefix(currentFile, "b/")
			}
			continue
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			added = append(added, strings.TrimPrefix(line, "+"))
		}
	}
	flush()
	return hunks
}

