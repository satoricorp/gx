package vcs

import (
	"context"
	"fmt"
	"strings"
)

func (s *Service) ensureBranchMutationAllowed(ctx context.Context, repoRoot, name, targetCommit string) error {
	name = cleanRefName(name)
	targetCommit = strings.TrimSpace(targetCommit)
	if name == "" || !s.isProtectedRef(ctx, repoRoot, name) {
		return nil
	}
	current, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "--verify", "refs/heads/"+name+"^{commit}")
	if err == nil && strings.TrimSpace(current) == targetCommit {
		return nil
	}
	// Advancing the checked-out branch to one of its own descendants is
	// git-commit semantics, not a base move; gx commit on the default branch
	// depends on it. Rewinds and moves of branches that are not checked out
	// stay protected.
	if err == nil && s.isCheckedOutBranch(ctx, repoRoot, name) &&
		s.isAncestorCommit(ctx, repoRoot, strings.TrimSpace(current), targetCommit) {
		return nil
	}
	return fmt.Errorf("ref %s is protected; gx will not move the authoring base or default branch", name)
}

func (s *Service) isCheckedOutBranch(ctx context.Context, repoRoot, name string) bool {
	current, err := s.runTrimmed(ctx, repoRoot, "git", "branch", "--show-current")
	if err != nil {
		return false
	}
	return cleanRefName(strings.TrimSpace(current)) == name
}

func (s *Service) isAncestorCommit(ctx context.Context, repoRoot, ancestor, descendant string) bool {
	ancestor = strings.TrimSpace(ancestor)
	descendant = strings.TrimSpace(descendant)
	if ancestor == "" || descendant == "" {
		return false
	}
	_, err := s.runner.Run(ctx, repoRoot, "git", "merge-base", "--is-ancestor", ancestor, descendant)
	return err == nil
}

func (s *Service) ensureBookmarkMutationAllowed(ctx context.Context, repoRoot, name, targetRev string) error {
	name = cleanRefName(name)
	if name == "" || !s.isProtectedRef(ctx, repoRoot, name) {
		return nil
	}
	targetCommit, err := s.commitIDForRev(ctx, repoRoot, targetRev)
	if err != nil {
		return err
	}
	current, currentErr := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "--verify", "refs/heads/"+name+"^{commit}")
	if currentErr == nil && strings.TrimSpace(current) == strings.TrimSpace(targetCommit) {
		return nil
	}
	return fmt.Errorf("ref %s is protected; gx will not move the authoring base or default branch", name)
}

func (s *Service) rejectProtectedStackBookmark(ctx context.Context, repoRoot, name string) error {
	name = cleanRefName(name)
	if name == "" || !s.isProtectedRef(ctx, repoRoot, name) {
		return nil
	}
	return fmt.Errorf("stack bookmark %s is protected; choose a non-base branch name", name)
}

func (s *Service) isProtectedRef(ctx context.Context, repoRoot, name string) bool {
	name = cleanRefName(name)
	if name == "" {
		return false
	}
	protected := map[string]struct{}{
		"main":   {},
		"master": {},
		"trunk":  {},
	}
	if store, err := openStore(ctx); err == nil {
		if repo, repoErr := store.FindRepoByRoot(ctx, repoRoot); repoErr == nil && repo != nil {
			if repo.DefaultBranch != nil {
				if ref := cleanRefName(*repo.DefaultBranch); ref != "" {
					protected[ref] = struct{}{}
				}
			}
			if repo.AuthoringBase != nil {
				if ref := cleanRefName(*repo.AuthoringBase); ref != "" {
					protected[ref] = struct{}{}
				}
			}
		}
		_ = store.Close()
	}
	if head, err := s.runTrimmed(ctx, repoRoot, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"); err == nil {
		head = strings.TrimSpace(head)
		if idx := strings.LastIndex(head, "/"); idx >= 0 && idx+1 < len(head) {
			head = head[idx+1:]
		}
		if ref := cleanRefName(head); ref != "" {
			protected[ref] = struct{}{}
		}
	}
	_, ok := protected[name]
	return ok
}
