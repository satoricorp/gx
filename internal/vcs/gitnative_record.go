package vcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
)

// RecordGitCommit records metadata for a git commit that already exists.
func (s *Service) RecordGitCommit(ctx context.Context, repo RepoInfo, commitOID string, pending PendingCommitContext) (CommitResult, error) {
	commitOID = strings.TrimSpace(commitOID)
	if commitOID == "" {
		return CommitResult{}, fmt.Errorf("commit oid required")
	}
	message, err := s.runTrimmed(ctx, repo.RootPath, "git", "log", "-1", "--format=%B", commitOID)
	if err != nil {
		return CommitResult{}, fmt.Errorf("read commit message: %w", err)
	}
	revisionIDs := ParseRevisionIDsFromMessage(message)
	if len(revisionIDs) == 0 {
		return CommitResult{}, fmt.Errorf("commit %s has no gx revision trailer", shortID(commitOID, 7))
	}
	revisionID := revisionIDs[len(revisionIDs)-1]
	if existing, lookupErr := s.commitResultFromDB(ctx, repo, commitOID); lookupErr == nil &&
		existing.Change.ChangeID == revisionID && existing.Change.CommitID == commitOID {
		return existing, nil
	}
	filesOut, err := s.runTrimmed(ctx, repo.RootPath, "git", "diff-tree", "--no-commit-id", "--name-only", "-r", commitOID)
	if err != nil {
		return CommitResult{}, fmt.Errorf("list commit files: %w", err)
	}
	parentOID, _ := s.runTrimmed(ctx, repo.RootPath, "git", "rev-parse", commitOID+"^")
	parentOID = strings.TrimSpace(parentOID)
	var parentChangeID *string
	if parentOID != "" {
		parentMessage, parentErr := s.runTrimmed(ctx, repo.RootPath, "git", "log", "-1", "--format=%B", parentOID)
		if parentErr == nil {
			if ids := ParseRevisionIDsFromMessage(parentMessage); len(ids) > 0 {
				parentChangeID = &ids[len(ids)-1]
			}
		}
	}
	subject := strings.TrimSpace(strings.Split(message, "\n")[0])
	change := ChangeInfo{
		ChangeID:       revisionID,
		CommitID:       commitOID,
		Description:    subject,
		ParentChangeID: parentChangeID,
		Files:          splitLines(filesOut),
	}
	branch := strings.TrimSpace(pending.Branch)
	if branch == "" {
		branch, _ = s.currentStagedCommitBranch(ctx, repo.RootPath)
	}
	stack, err := s.resolveStagedCommitStack(ctx, repo, branch, parentOID)
	if err != nil {
		return CommitResult{}, err
	}
	// Sessions are deliberately not matched here. A `sessions` row exists
	// because a transcript was observed, never because a commit happened to
	// match one; the links are written at push time from the hunk links the
	// capture pipeline already produces (AttachSessionsFromHunkLinks).
	result := CommitResult{
		Repo:        repo,
		Change:      change,
		Stack:       &stack,
		OperationID: commitOID,
	}
	if err := recordCommit(ctx, result); err != nil {
		return result, err
	}
	return result, nil
}

// RewriteGitCommitOIDs updates stored commit OIDs after history rewrite.
func (s *Service) RewriteGitCommitOIDs(ctx context.Context, repo RepoInfo, oldOID, newOID string) error {
	oldOID = strings.TrimSpace(oldOID)
	newOID = strings.TrimSpace(newOID)
	if oldOID == "" || newOID == "" || oldOID == newOID {
		return nil
	}
	return withRepoIdentityLock(repo.GitCommonDir, func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		repoRow, err := store.FindRepoByIdentity(ctx, repo.GitCommonDir, repo.RootPath)
		if err != nil {
			return err
		}
		if repoRow == nil {
			return nil
		}
		change, err := store.FindChangeByCommitID(ctx, repoRow.ID, oldOID)
		if err != nil {
			return err
		}
		if change == nil {
			return s.recoverGitCommitByOID(ctx, store, repo, *repoRow, newOID)
		}
		now := time.Now().UnixMilli()
		if err := store.UpdateChangeCommitID(ctx, change.ID, newOID, now); err != nil {
			return err
		}
		filesOut, _ := s.runTrimmed(ctx, repo.RootPath, "git", "diff-tree", "--no-commit-id", "--name-only", "-r", newOID)
		filesJSON, err := json.Marshal(splitLines(filesOut))
		if err != nil {
			return err
		}
		return store.WriteChangeRevision(ctx, storage.ChangeRevision{
			ChangeID:      change.ID,
			JJCommitID:    newOID,
			JJOperationID: newOID,
			ChangedFiles:  string(filesJSON),
			CreatedAt:     now,
		})
	})
}

func (s *Service) recoverGitCommitByOID(ctx context.Context, store *storage.Store, repo RepoInfo, repoRow storage.Repo, commitOID string) error {
	message, err := s.runTrimmed(ctx, repo.RootPath, "git", "log", "-1", "--format=%B", commitOID)
	if err != nil {
		return nil
	}
	revisionIDs := ParseRevisionIDsFromMessage(message)
	if len(revisionIDs) == 0 {
		return nil
	}
	existing, err := store.FindChangeByJJChangeID(ctx, repoRow.ID, revisionIDs[len(revisionIDs)-1])
	if err != nil {
		return err
	}
	if existing != nil {
		if strings.TrimSpace(existing.CurrentCommitID) == strings.TrimSpace(commitOID) {
			return nil
		}
		return store.UpdateChangeCommitID(ctx, existing.ID, commitOID, time.Now().UnixMilli())
	}
	_, err = s.RecordGitCommit(ctx, repo, commitOID, PendingCommitContext{
		WorktreeRoot: repo.RootPath,
		GitCommonDir: repo.GitCommonDir,
	})
	return err
}

// RecoverMissingRevisions scans commit messages and records missing DB rows.
func (s *Service) RecoverMissingRevisions(ctx context.Context, repo RepoInfo, revisionIDs []string) (int, error) {
	if len(revisionIDs) == 0 {
		return 0, nil
	}
	recovered := 0
	var recoveryErrors []error
	store, err := openStore(ctx)
	if err != nil {
		return 0, err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByIdentity(ctx, repo.GitCommonDir, repo.RootPath)
	if err != nil {
		return 0, err
	}
	if repoRow == nil {
		return 0, nil
	}
	for _, revisionID := range revisionIDs {
		revisionID = strings.TrimSpace(revisionID)
		if revisionID == "" {
			continue
		}
		existing, err := store.FindChangeByJJChangeID(ctx, repoRow.ID, revisionID)
		if err != nil {
			return recovered, err
		}
		commitOID, err := s.findCommitByRevisionID(ctx, repo.RootPath, revisionID)
		if err != nil {
			recoveryErrors = append(recoveryErrors, fmt.Errorf("find revision %s: %w", revisionID, err))
			continue
		}
		if commitOID == "" {
			recoveryErrors = append(recoveryErrors, fmt.Errorf("find revision %s: no matching commit", revisionID))
			continue
		}
		if existing != nil {
			if strings.TrimSpace(existing.CurrentCommitID) != strings.TrimSpace(commitOID) {
				if err := store.UpdateChangeCommitID(ctx, existing.ID, commitOID, time.Now().UnixMilli()); err != nil {
					recoveryErrors = append(recoveryErrors, fmt.Errorf("recover revision %s commit oid: %w", revisionID, err))
					continue
				}
				recovered++
			}
			continue
		}
		if _, err := s.RecordGitCommit(ctx, repo, commitOID, PendingCommitContext{
			WorktreeRoot: repo.RootPath,
			GitCommonDir: repo.GitCommonDir,
		}); err != nil {
			recoveryErrors = append(recoveryErrors, fmt.Errorf("recover revision %s: %w", revisionID, err))
			continue
		}
		recovered++
	}
	return recovered, errors.Join(recoveryErrors...)
}

func (s *Service) findCommitByRevisionID(ctx context.Context, repoRoot, revisionID string) (string, error) {
	needle := RevisionTrailerLine(revisionID)
	out, err := s.runTrimmed(ctx, repoRoot, "git", "log", "--all", "--format=%H%x00%B", "--grep", needle, "-n", "1")
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(out, "\x00", 2)
	if len(parts) == 0 {
		return "", nil
	}
	return strings.TrimSpace(parts[0]), nil
}
