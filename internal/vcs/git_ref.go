package vcs

import (
	"context"
	"strings"
)

// setGitBranchRef points refs/heads/branchName at commitID.
// Uses update-ref so checked-out worktree branches can be moved (git branch -f cannot).
func (s *Service) setGitBranchRef(ctx context.Context, repoRoot, branchName, commitID string) error {
	branchName = strings.TrimSpace(branchName)
	commitID = strings.TrimSpace(commitID)
	if branchName == "" || commitID == "" {
		return nil
	}
	_, err := s.runner.Run(ctx, repoRoot, "git", "update-ref", "refs/heads/"+branchName, commitID)
	return err
}
