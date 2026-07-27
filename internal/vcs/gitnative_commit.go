package vcs

import (
	"context"
	"fmt"
	"strings"
)

// PrepareCommitMessageHook ensures exactly one GX revision trailer is present.
func PrepareCommitMessageHook(message string) (string, error) {
	existing := ParseRevisionIDsFromMessage(message)
	if len(existing) > 0 {
		return StampRevisionTrailer(message, existing[len(existing)-1]), nil
	}
	revisionID, err := GenerateRevisionID()
	if err != nil {
		return "", err
	}
	return StampRevisionTrailer(message, revisionID), nil
}

// RunPostCommitHook records the latest git commit as a GX revision.
func (s *Service) RunPostCommitHook(ctx context.Context, repoRoot string) error {
	repo, err := s.ResolveGXRepoAtPath(ctx, repoRoot)
	if err != nil {
		return err
	}
	return withRepoIdentityLock(repo.GitCommonDir, func() error {
		pending, err := readPendingCommitContext(repo.GitDir)
		if err != nil {
			pending = PendingCommitContext{
				WorktreeRoot: repo.RootPath,
				GitCommonDir: repo.GitCommonDir,
			}
		}
		if pending.WorktreeRoot == "" {
			pending.WorktreeRoot = repo.RootPath
		}
		headOID, err := s.runTrimmed(ctx, repo.RootPath, "git", "rev-parse", "HEAD")
		if err != nil {
			return err
		}
		_, recordErr := s.RecordGitCommit(ctx, repo, strings.TrimSpace(headOID), pending)
		clearPendingCommitContext(repo.GitDir)
		return recordErr
	})
}

// RunPostRewriteHook updates rewritten commit OIDs from git post-rewrite stdin.
func (s *Service) RunPostRewriteHook(ctx context.Context, repoRoot string, mappings []CommitOIDMapping) error {
	repo, err := s.ResolveGXRepoAtPath(ctx, repoRoot)
	if err != nil {
		return err
	}
	for _, mapping := range mappings {
		if err := s.RewriteGitCommitOIDs(ctx, repo, mapping.OldOID, mapping.NewOID); err != nil {
			return err
		}
	}
	return nil
}

// CommitOIDMapping pairs old and new commit object IDs from post-rewrite.
type CommitOIDMapping struct {
	OldOID string
	NewOID string
}

// ParsePostRewriteMappings parses git post-rewrite stdin lines.
func ParsePostRewriteMappings(input string) []CommitOIDMapping {
	var out []CommitOIDMapping
	for _, line := range strings.Split(input, "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			continue
		}
		out = append(out, CommitOIDMapping{OldOID: fields[0], NewOID: fields[1]})
	}
	return out
}

func (s *Service) commitResultFromDB(ctx context.Context, repo RepoInfo, commitOID string) (CommitResult, error) {
	message, err := s.runTrimmed(ctx, repo.RootPath, "git", "log", "-1", "--format=%B", commitOID)
	if err != nil {
		return CommitResult{}, err
	}
	revisionIDs := ParseRevisionIDsFromMessage(message)
	if len(revisionIDs) == 0 {
		return CommitResult{}, fmt.Errorf("commit has no GX revision trailer")
	}
	revisionID := revisionIDs[len(revisionIDs)-1]
	store, err := openStore(ctx)
	if err != nil {
		return CommitResult{}, err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByIdentity(ctx, repo.GitCommonDir, repo.RootPath)
	if err != nil || repoRow == nil {
		return CommitResult{}, fmt.Errorf("repo not initialized")
	}
	changeRow, err := store.FindChangeByJJChangeID(ctx, repoRow.ID, revisionID)
	if err != nil || changeRow == nil {
		return CommitResult{}, fmt.Errorf("revision not recorded")
	}
	change := ChangeInfo{
		ChangeID:       changeRow.JJChangeID,
		CommitID:       changeRow.CurrentCommitID,
		Description:    changeRow.Description,
		ParentChangeID: changeRow.ParentChangeID,
	}
	var stack *StackInfo
	if branch, branchErr := s.currentStagedCommitBranch(ctx, repo.RootPath); branchErr == nil {
		if existing, err := store.FindStackByBookmark(ctx, repoRow.ID, branch); err == nil && existing != nil {
			info := stackInfoFromStorage(*existing)
			stack = &info
		}
	}
	return CommitResult{
		Repo:             repo,
		Change:           change,
		Stack:            stack,
		OperationID:      commitOID,
		ProvenanceStatus: "recorded",
	}, nil
}
