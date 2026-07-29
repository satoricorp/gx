package vcs

import (
	"context"
	"os"
)

// EnsureReadyResult reports auto-init/bootstrap work performed for a repository.
type EnsureReadyResult struct {
	InitResult
	Prepared bool
}

func (s *Service) repoRootFromPath(ctx context.Context, startPath string) (string, bool, error) {
	if gitRepo, err := s.ResolveGitRepoAtPath(ctx, startPath); err == nil {
		return gitRepo.RootPath, true, nil
	}
	return "", false, nil
}

func (s *Service) isRepoInitialized(ctx context.Context, repoRoot string) (bool, error) {
	repo, err := s.ResolveTotalityRepoAtPath(ctx, repoRoot)
	if err != nil {
		store, storeErr := openStore(ctx)
		if storeErr != nil {
			return false, storeErr
		}
		defer store.Close()
		return store.IsInitializedRepo(ctx, repoRoot)
	}
	return s.isRepoInitializedByIdentity(ctx, repo.GitCommonDir, repo.RootPath)
}

// EnsureReadyRepo initializes tl for the repository at startPath when needed.
func (s *Service) EnsureReadyRepo(ctx context.Context, startPath string) (EnsureReadyResult, error) {
	repoRoot, inRepo, err := s.repoRootFromPath(ctx, startPath)
	if err != nil {
		return EnsureReadyResult{}, err
	}
	if !inRepo {
		return EnsureReadyResult{}, nil
	}

	repo, err := s.ResolveTotalityRepoAtPath(ctx, repoRoot)
	if err != nil {
		return EnsureReadyResult{}, err
	}
	initialized, err := s.isRepoInitializedByIdentity(ctx, repo.GitCommonDir, repo.RootPath)
	if err != nil {
		return EnsureReadyResult{}, err
	}

	if initialized {
		return EnsureReadyResult{InitResult: InitResult{Repo: repo}}, nil
	}

	result := EnsureReadyResult{Prepared: true}
	initResult, err := s.InitAtPath(ctx, repoRoot, InitOptions{Interactive: false})
	if err != nil {
		return EnsureReadyResult{}, err
	}
	result.InitResult = initResult
	return result, nil
}

// RepoRootFromWorkingDirectory resolves the git repository root for cwd.
func (s *Service) RepoRootFromWorkingDirectory(ctx context.Context) (string, bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false, err
	}
	return s.repoRootFromPath(ctx, cwd)
}
