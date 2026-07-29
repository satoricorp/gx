package vcs

import (
	"context"
	"strings"
	"time"

	"github.com/satoricorp/totality/internal/storage"
)

func repoRemoteName(repo RepoInfo) string {
	if repo.DefaultRemote != nil && strings.TrimSpace(*repo.DefaultRemote) != "" {
		return strings.TrimSpace(*repo.DefaultRemote)
	}
	return "origin"
}

func (s *Service) reconcileStoredStacksRemoteState(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64, stacks []storage.Stack) error {
	if len(stacks) == 0 {
		return nil
	}
	remoteName := repoRemoteName(repo)
	defaultBranch := repo.defaultBaseBranch()
	now := time.Now().UnixMilli()

	for _, stack := range stacks {
		if IsTerminalStackStatus(stack.Status) {
			continue
		}
		info := stackInfoFromStorage(stack)
		publishRef := publishRefForStack(info)
		if publishRef == "" {
			continue
		}
		remoteHead, err := s.remoteBranchHead(ctx, repo.RootPath, remoteName, publishRef)
		if err != nil || remoteHead == "" {
			continue
		}

		changes, err := store.ListChangesByStackID(ctx, stack.ID)
		if err != nil {
			return err
		}

		remoteRef := "refs/heads/" + publishRef
		updated := false
		if stack.RemoteRef == nil || strings.TrimSpace(*stack.RemoteRef) != remoteRef {
			stack.RemoteRef = &remoteRef
			updated = true
		}
		if stack.RemoteName == nil || strings.TrimSpace(*stack.RemoteName) != remoteName {
			stack.RemoteName = &remoteName
			updated = true
		}

		publishedCount := 0
		for _, change := range changes {
			commitID := strings.TrimSpace(change.CurrentCommitID)
			if commitID == "" {
				continue
			}
			onRemote, err := s.commitReachableFrom(ctx, repo.RootPath, commitID, remoteHead)
			if err != nil {
				return err
			}
			if !onRemote {
				continue
			}
			publishedCount++
			if err := store.UpsertChangeBookmark(ctx, storage.ChangeBookmark{
				ChangeID:           change.ID,
				BookmarkName:       publishRef,
				RemoteName:         &remoteName,
				RemoteRef:          &remoteRef,
				LastPushedCommitID: commitID,
				CreatedAt:          now,
				UpdatedAt:          now,
			}); err != nil {
				return err
			}
		}

		localHead := strings.TrimSpace(derefString(stack.HeadCommitID))
		nextStatus := strings.TrimSpace(stack.Status)
		if localHead != "" && localHead == remoteHead && publishedCount == len(changes) && len(changes) > 0 {
			if s.commitMergedIntoBranch(ctx, repo.RootPath, remoteHead, defaultBranch) {
				nextStatus = "merged"
			} else {
				nextStatus = "published"
			}
		} else if publishedCount > 0 && publishedCount < len(changes) {
			nextStatus = "draft"
		} else if publishedCount == len(changes) && len(changes) > 0 && localHead == remoteHead {
			nextStatus = "published"
		}

		if nextStatus != "" && nextStatus != strings.TrimSpace(stack.Status) {
			stack.Status = nextStatus
			updated = true
		}
		if updated {
			stack.UpdatedAt = now
			if _, err := store.UpsertStack(ctx, stack); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Service) commitReachableFrom(ctx context.Context, repoRoot, commit, head string) (bool, error) {
	commit = strings.TrimSpace(commit)
	head = strings.TrimSpace(head)
	if commit == "" || head == "" {
		return false, nil
	}
	if commit == head {
		return true, nil
	}
	return s.gitCommitIsAncestor(ctx, repoRoot, commit, head)
}

func (s *Service) commitMergedIntoBranch(ctx context.Context, repoRoot, commit, branch string) bool {
	commit = strings.TrimSpace(commit)
	branch = strings.TrimSpace(branch)
	if commit == "" || branch == "" {
		return false
	}
	if merged, _ := s.gitCommitIsAncestor(ctx, repoRoot, commit, branch); merged {
		return true
	}
	originBranch := branch
	if !strings.HasPrefix(originBranch, "origin/") {
		originBranch = "origin/" + strings.TrimPrefix(originBranch, "origin/")
	}
	merged, _ := s.gitCommitIsAncestor(ctx, repoRoot, commit, originBranch)
	return merged
}
