package vcs

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) SyncCloudBookmarkTip(ctx context.Context, repo RepoInfo, branchName string) error {
	branchName = strings.TrimSpace(branchName)
	if branchName == "" {
		return fmt.Errorf("branch name is required")
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "git", "fetch", "origin", branchName); err != nil {
		return err
	}
	remoteRef := "refs/remotes/origin/" + branchName
	if err := s.setBookmarkTargetAtRev(ctx, repo.RootPath, branchName, remoteRef, true); err != nil {
		return fmt.Errorf("set local branch %s: %w", branchName, err)
	}
	return nil
}
