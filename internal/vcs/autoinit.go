package vcs

import (
	"context"
	"fmt"
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
	if jjRepo, err := s.ResolveJJRepoAtPath(ctx, startPath); err == nil {
		return jjRepo.RootPath, true, nil
	}
	return "", false, nil
}

func (s *Service) isRepoInitialized(ctx context.Context, repoRoot string) (bool, error) {
	store, err := openStore(ctx)
	if err != nil {
		return false, err
	}
	defer store.Close()
	return store.IsInitializedRepo(ctx, repoRoot)
}

// EnsureReadyRepo initializes gx for the repository at startPath when needed and
// bootstraps the authoring base revision so stack commands can resolve main.
func (s *Service) EnsureReadyRepo(ctx context.Context, startPath string) (EnsureReadyResult, error) {
	repoRoot, inRepo, err := s.repoRootFromPath(ctx, startPath)
	if err != nil {
		return EnsureReadyResult{}, err
	}
	if !inRepo {
		return EnsureReadyResult{}, nil
	}

	initialized, err := s.isRepoInitialized(ctx, repoRoot)
	if err != nil {
		return EnsureReadyResult{}, err
	}

	var repo RepoInfo
	needsBase := false
	jjRepo, jjErr := s.ResolveJJRepoAtPath(ctx, repoRoot)
	switch {
	case jjErr != nil:
		needsBase = true
	default:
		repo = jjRepo
		exists, revErr := s.RevisionExists(ctx, repoRoot, repo.defaultBaseBranch())
		if revErr != nil {
			return EnsureReadyResult{}, revErr
		}
		needsBase = !exists
	}

	if initialized && !needsBase {
		return EnsureReadyResult{InitResult: InitResult{Repo: repo}}, nil
	}

	result := EnsureReadyResult{Prepared: true}
	initResult, err := s.InitAtPath(ctx, repoRoot, InitOptions{Interactive: false})
	if err != nil {
		return EnsureReadyResult{}, fmt.Errorf("initialize gx: %w", err)
	}
	result.InitResult = initResult

	if needsBase {
		target := initResult.Repo
		if target.RootPath == "" {
			target = repo
		}
		if err := s.ensureAuthoringBaseRevision(ctx, target, target.defaultBaseBranch()); err != nil {
			return result, fmt.Errorf("prepare gx authoring base: %w", err)
		}
	}

	return result, nil
}

// RepoRootFromWorkingDirectory resolves the git or jj repository root for cwd.
func (s *Service) RepoRootFromWorkingDirectory(ctx context.Context) (string, bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false, err
	}
	return s.repoRootFromPath(ctx, cwd)
}
