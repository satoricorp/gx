package vcs

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) SyncCloudBookmarkTip(ctx context.Context, repo RepoInfo, branchName string) error {
	branchName = strings.TrimSpace(branchName)
	if branchName == "" {
		return fmt.Errorf("bookmark name is required")
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "git", "fetch"); err != nil {
		return err
	}
	remoteRef := branchName + "@origin"
	if err := s.setBookmarkTargetAtRev(ctx, repo.RootPath, branchName, remoteRef, true); err != nil {
		return fmt.Errorf("set local bookmark %s: %w", branchName, err)
	}
	return nil
}
