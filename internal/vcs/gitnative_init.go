package vcs

import (
	"context"
)

func (s *Service) isRepoInitializedByIdentity(ctx context.Context, gitCommonDir, worktreeRoot string) (bool, error) {
	store, err := openStore(ctx)
	if err != nil {
		return false, err
	}
	defer store.Close()
	repo, err := store.FindRepoByIdentity(ctx, gitCommonDir, worktreeRoot)
	if err != nil {
		return false, err
	}
	if repo == nil {
		return false, nil
	}
	return store.IsInitializedRepo(ctx, repo.RootPath)
}
