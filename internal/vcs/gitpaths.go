package vcs

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// GitPaths describes the active worktree and shared git metadata directory.
type GitPaths struct {
	WorktreeRoot string
	CommonDir    string
	GitDir       string
	HooksDir     string
}

// ResolveGitPaths resolves git paths for startPath using git itself.
func (s *Service) ResolveGitPaths(ctx context.Context, startPath string) (GitPaths, error) {
	worktree, err := s.runTrimmed(ctx, startPath, "git", "rev-parse", "--show-toplevel")
	if err != nil {
		return GitPaths{}, fmt.Errorf("not a git repository: %w", err)
	}
	worktree = strings.TrimSpace(worktree)
	commonDir, err := s.runTrimmed(ctx, worktree, "git", "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return GitPaths{}, fmt.Errorf("read git common dir: %w", err)
	}
	commonDir = NormalizeGitCommonDir(worktree, commonDir)
	gitDir, err := s.runTrimmed(ctx, worktree, "git", "rev-parse", "--path-format=absolute", "--git-dir")
	if err != nil {
		return GitPaths{}, fmt.Errorf("read git dir: %w", err)
	}
	gitDir = NormalizeGitCommonDir(worktree, gitDir)
	hooksDir, err := s.runTrimmed(ctx, worktree, "git", "rev-parse", "--path-format=absolute", "--git-path", "hooks")
	if err != nil {
		return GitPaths{}, fmt.Errorf("read git hooks dir: %w", err)
	}
	hooksDir = filepath.Clean(strings.TrimSpace(hooksDir))
	return GitPaths{
		WorktreeRoot: worktree,
		CommonDir:    commonDir,
		GitDir:       gitDir,
		HooksDir:     hooksDir,
	}, nil
}

// NormalizeGitCommonDir returns an absolute common directory path.
func NormalizeGitCommonDir(worktreeRoot, raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !filepath.IsAbs(raw) {
		raw = filepath.Join(worktreeRoot, raw)
	}
	return filepath.Clean(raw)
}

// RepoIdentityKey is the stable repository identity shared by linked worktrees.
func RepoIdentityKey(commonDir string) string {
	return NormalizeGitCommonDir("", commonDir)
}

func (s *Service) ResolveGxRepoAtPath(ctx context.Context, startPath string) (RepoInfo, error) {
	paths, err := s.ResolveGitPaths(ctx, startPath)
	if err != nil {
		return RepoInfo{}, err
	}
	info, err := s.resolveGitInfo(ctx, paths.WorktreeRoot)
	if err != nil {
		return RepoInfo{}, err
	}
	info.RootPath = paths.WorktreeRoot
	info.GitCommonDir = paths.CommonDir
	info.GitDir = paths.GitDir
	if info.Backend == "" {
		info.Backend = "git"
	}
	info = s.withStoredRepoConfigByIdentity(ctx, info)
	return info, nil
}

func (s *Service) withStoredRepoConfigByIdentity(ctx context.Context, info RepoInfo) RepoInfo {
	store, err := openStore(ctx)
	if err != nil {
		return info
	}
	defer store.Close()
	repo, err := store.FindRepoByIdentity(ctx, info.GitCommonDir, info.RootPath)
	if err != nil || repo == nil {
		return info
	}
	info.AuthoringBase = repo.AuthoringBase
	if repo.Backend != "" {
		info.Backend = repo.Backend
	}
	return info
}
