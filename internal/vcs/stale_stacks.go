package vcs

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// StaleStack describes a stored stack whose jj bookmark no longer resolves.
type StaleStack struct {
	StackID      int64  `json:"stack_id"`
	Name         string `json:"name"`
	BookmarkName string `json:"bookmark_name"`
	Status       string `json:"status"`
}

// StaleStackCleanupResult reports stale stacks and any repair actions taken.
type StaleStackCleanupResult struct {
	RepoRoot string       `json:"repo_root"`
	Stale    []StaleStack `json:"stale"`
	Actions  []string     `json:"actions"`
}

// ListStaleStacks returns non-terminal stack rows whose bookmarks are missing from jj.
func (s *Service) ListStaleStacks(ctx context.Context) (StaleStackCleanupResult, error) {
	return s.staleStacks(ctx, false)
}

// CleanupStaleStacks marks stale stack rows as closed so routing and dedup skip them.
func (s *Service) CleanupStaleStacks(ctx context.Context) (StaleStackCleanupResult, error) {
	return s.staleStacks(ctx, true)
}

func (s *Service) staleStacks(ctx context.Context, repair bool) (StaleStackCleanupResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return StaleStackCleanupResult{}, err
	}
	result := StaleStackCleanupResult{RepoRoot: repo.RootPath}
	store, err := openStore(ctx)
	if err != nil {
		return result, err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return result, err
	}
	if repoRow == nil {
		return result, nil
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoRow.ID)
	if err != nil {
		return result, err
	}
	now := time.Now().UnixMilli()
	for _, stack := range stacks {
		if IsTerminalStackStatus(stack.Status) {
			continue
		}
		bookmark := strings.TrimSpace(stack.BookmarkName)
		if bookmark == "" {
			continue
		}
		exists, err := s.RevisionExists(ctx, repo.RootPath, bookmark)
		if err != nil {
			return result, err
		}
		if exists {
			continue
		}
		result.Stale = append(result.Stale, StaleStack{
			StackID:      stack.ID,
			Name:         stack.Name,
			BookmarkName: bookmark,
			Status:       stack.Status,
		})
		if repair {
			if err := store.MarkStackStatus(ctx, stack.ID, "closed", now); err != nil {
				return result, err
			}
			result.Actions = append(result.Actions, fmt.Sprintf(
				"marked stale stack %s as closed (bookmark not found in jj)", bookmark,
			))
		}
	}
	return result, nil
}
