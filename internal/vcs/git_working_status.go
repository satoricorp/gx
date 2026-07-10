package vcs

import (
	"context"
	"fmt"
	"strings"
)

const HintPush = "gx push"

// GitWorkingStatus is pure-git working tree state from git status --porcelain.
type GitWorkingStatus struct {
	Staged    []string `json:"staged,omitempty"`
	Unstaged  []string `json:"unstaged,omitempty"`
	Untracked []string `json:"untracked,omitempty"`
}

func (s GitWorkingStatus) StagedCount() int   { return len(s.Staged) }
func (s GitWorkingStatus) UnstagedCount() int { return len(s.Unstaged) }
func (s GitWorkingStatus) UntrackedCount() int {
	return len(s.Untracked)
}
func (s GitWorkingStatus) Dirty() bool {
	return s.StagedCount()+s.UnstagedCount()+s.UntrackedCount() > 0
}

func (s *Service) GitWorkingStatus(ctx context.Context, repoRoot string) (GitWorkingStatus, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return GitWorkingStatus{}, fmt.Errorf("git repo root is required")
	}
	out, err := s.runTrimmed(ctx, repoRoot, "git", "status", "--porcelain")
	if err != nil {
		return GitWorkingStatus{}, fmt.Errorf("git status --porcelain: %w", err)
	}
	return parseGitWorkingStatus(out), nil
}

func parseGitWorkingStatus(out string) GitWorkingStatus {
	var status GitWorkingStatus
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if len(line) < 2 {
			continue
		}
		x, y := line[0], line[1]
		path := strings.TrimSpace(line[2:])
		if path == "" && x != '?' {
			continue
		}
		if x == '?' && y == '?' {
			status.Untracked = append(status.Untracked, porcelainPath(path))
			continue
		}
		path = porcelainPath(path)
		if x != ' ' && x != '?' {
			status.Staged = append(status.Staged, path)
		}
		if y != ' ' && y != '?' {
			status.Unstaged = append(status.Unstaged, path)
		}
	}
	return status
}

func porcelainPath(raw string) string {
	raw = strings.TrimSpace(raw)
	if i := strings.Index(raw, " -> "); i >= 0 {
		return strings.TrimSpace(raw[i+4:])
	}
	return raw
}
