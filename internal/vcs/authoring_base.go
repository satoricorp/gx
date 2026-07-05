package vcs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	bootstrapInitCommitMessage = "init commit"
	bootstrapReadmeName        = "README.md"
)

func (s *Service) ensureAuthoringBaseRevision(ctx context.Context, repo RepoInfo, baseRef string) error {
	baseRef = strings.TrimSpace(baseRef)
	if baseRef == "" {
		baseRef = repo.defaultBaseBranch()
	}
	exists, err := s.RevisionExists(ctx, repo.RootPath, baseRef)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return withRepoLock(repo.RootPath, func() error {
		return s.bootstrapAuthoringBaseUnlocked(ctx, repo, baseRef)
	})
}

func (s *Service) resolveBootstrapBaseRef(ctx context.Context, repo RepoInfo, baseRef string) string {
	baseRef = strings.TrimPrefix(strings.TrimSpace(baseRef), "origin/")
	if baseRef != "" {
		if exists, err := s.RevisionExists(ctx, repo.RootPath, baseRef); err == nil && exists {
			return baseRef
		}
	}
	if branch := s.currentGitCheckoutRef(ctx, repo.RootPath); branch != "" {
		return branch
	}
	if repo.DefaultBranch != nil && strings.TrimSpace(*repo.DefaultBranch) != "" {
		return strings.TrimSpace(*repo.DefaultBranch)
	}
	return repo.defaultBaseBranch()
}

func (s *Service) bootstrapAuthoringBaseUnlocked(ctx context.Context, repo RepoInfo, baseRef string) error {
	baseRef = s.resolveBootstrapBaseRef(ctx, repo, baseRef)
	exists, err := s.RevisionExists(ctx, repo.RootPath, baseRef)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if _, _, _, err := s.ensureIdentity(ctx, repo.RootPath, InitOptions{}); err != nil {
		return err
	}
	if err := s.ensureGXInternalIgnored(ctx, repo.RootPath); err != nil {
		return err
	}
	if err := s.ensureBootstrapReadme(repo.RootPath); err != nil {
		return err
	}
	readmePath := filepath.Join(repo.RootPath, bootstrapReadmeName)
	if _, err := os.Stat(readmePath); err != nil {
		return fmt.Errorf("bootstrap readme missing at %s: %w", readmePath, err)
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "git", "add", bootstrapReadmeName); err != nil {
		return err
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "git", "commit", "-m", bootstrapInitCommitMessage); err != nil {
		return err
	}
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "git", "import"); err != nil {
		return err
	}
	exists, err = s.RevisionExists(ctx, repo.RootPath, baseRef)
	if err != nil {
		return err
	}
	if !exists {
		headRev := "main"
		if head, err := s.runStdoutTrimmed(ctx, repo.RootPath, "git", "rev-parse", "--abbrev-ref", "HEAD"); err == nil && strings.TrimSpace(head) != "" {
			headRev = strings.TrimSpace(head)
		}
		if err := s.setBookmarkTargetAtRev(ctx, repo.RootPath, baseRef, headRev); err != nil {
			return err
		}
	}
	store, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return err
	}
	return store.SetRepoAuthoringBase(ctx, repoID, baseRef, time.Now().UnixMilli())
}

func (s *Service) ensureBootstrapReadme(repoRoot string) error {
	path := filepath.Join(repoRoot, bootstrapReadmeName)
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	name := filepath.Base(repoRoot)
	if name == "" || name == "." || name == string(os.PathSeparator) {
		name = "repo"
	}
	content := fmt.Sprintf("# %s\n", name)
	return os.WriteFile(path, []byte(content), 0o644)
}
