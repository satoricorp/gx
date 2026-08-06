package vcs

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrDetachedHEAD = errors.New("detached HEAD")

func isUnbornHEADError(err error) bool {
	if err == nil {
		return false
	}
	value := strings.ToLower(err.Error())
	return strings.Contains(value, "unknown revision") ||
		strings.Contains(value, "bad revision") ||
		strings.Contains(value, "ambiguous argument 'head'") ||
		strings.Contains(value, "your current branch does not have any commits yet")
}

func (s *Service) currentStagedCommitBranch(ctx context.Context, repoRoot string) (string, error) {
	branch, err := s.runTrimmed(ctx, repoRoot, "git", "branch", "--show-current")
	if err != nil {
		return "", err
	}
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return "", fmt.Errorf("%w; check out a branch before committing", ErrDetachedHEAD)
	}
	return cleanRefName(branch), nil
}

// resolveStagedCommitStack maps a recorded commit onto a stack. The stack is
// the checked-out branch: gx records what git already did and never moves HEAD
// or mints branches of its own.
func (s *Service) resolveStagedCommitStack(ctx context.Context, repo RepoInfo, branch, headCommit string) (StackInfo, error) {
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
	onBase := branch == baseCheckoutRef(baseRef) || branch == baseRef
	bookmark := branch
	store, err := openStore(ctx)
	if err != nil {
		return StackInfo{}, err
	}
	defer store.Close()
	repoID, err := upsertRepo(ctx, store, repo)
	if err != nil {
		return StackInfo{}, err
	}
	if existing, err := store.FindStackByBookmark(ctx, repoID, bookmark); err != nil {
		return StackInfo{}, err
	} else if existing != nil {
		return stackInfoFromStorage(*existing), nil
	}
	baseCommit := s.stackBaseCommitID(ctx, repo.RootPath, baseRef)
	if onBase && strings.TrimSpace(headCommit) != "" {
		baseCommit = strings.TrimSpace(headCommit)
	}
	return StackInfo{
		Name:         stackNameFromBookmark(bookmark),
		BookmarkName: bookmark,
		BaseRef:      baseRef,
		BaseCommitID: baseCommit,
		Status:       "draft",
	}, nil
}

func (s *Service) refExists(ctx context.Context, repoRoot, ref string) bool {
	_, err := s.runTrimmed(ctx, repoRoot, "git", "rev-parse", "--verify", "refs/heads/"+cleanRefName(ref))
	return err == nil
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
