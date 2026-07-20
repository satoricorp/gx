package vcs

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// CommitStagedViaGit records staged changes by delegating to native git commit.
func (s *Service) CommitStagedViaGit(ctx context.Context, opts StagedRevisionOptions) (CommitResult, error) {
	if err := ValidateCommitMessage(opts.Message); err != nil {
		return CommitResult{}, err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return CommitResult{}, err
	}
	repo, err := s.ResolveGXRepoAtPath(ctx, cwd)
	if err != nil {
		return CommitResult{}, fmt.Errorf("gx commit must run inside a git repository: %w", err)
	}
	var result CommitResult
	err = func() error {
		if err := s.preflightStagedIndex(ctx, repo.RootPath); err != nil {
			return err
		}
		stagedFiles, err := s.stagedFileNames(ctx, repo.RootPath)
		if err != nil {
			return err
		}
		if len(stagedFiles) == 0 {
			return codedError(ExitCodeNoStagedChanges, ErrNoStagedChanges)
		}
		branch, err := s.currentStagedCommitBranch(ctx, repo.RootPath)
		if err != nil {
			return err
		}
		requestedBranch := cleanRefName(opts.Branch)
		if requestedBranch != "" {
			if s.isProtectedRef(ctx, repo.RootPath, requestedBranch) {
				return fmt.Errorf("branch %s is protected; choose a non-base branch name", requestedBranch)
			}
			if s.refExists(ctx, repo.RootPath, requestedBranch) {
				return fmt.Errorf("branch %s already exists; git switch %s and run gx commit without --branch", requestedBranch, requestedBranch)
			}
			if _, err := s.runner.Run(ctx, repo.RootPath, "git", "checkout", "-b", requestedBranch); err != nil {
				return fmt.Errorf("create branch %s: %w", requestedBranch, err)
			}
			branch = requestedBranch
		}
		pending := PendingCommitContext{
			WorktreeRoot:        repo.RootPath,
			GitCommonDir:        repo.GitCommonDir,
			Branch:              branch,
			PreferredSessionIDs: opts.PreferredSessionIDs,
			SessionContexts:     opts.SessionContexts,
			SelfReport:          opts.SelfReport,
		}
		if err := writePendingCommitContext(repo.GitDir, pending); err != nil {
			return fmt.Errorf("write commit context: %w", err)
		}
		defer clearPendingCommitContext(repo.GitDir)
		stamped, err := PrepareCommitMessageHook(opts.Message)
		if err != nil {
			return err
		}
		if err := s.runner.RunStream(ctx, repo.RootPath, "git", "commit", "-m", stamped); err != nil {
			return fmt.Errorf("git commit: %w", err)
		}
		headOID, err := s.runTrimmed(ctx, repo.RootPath, "git", "rev-parse", "HEAD")
		if err != nil {
			return fmt.Errorf("read HEAD: %w", err)
		}
		result, err = s.loadCommitResultFromHEAD(ctx, repo, strings.TrimSpace(headOID), pending)
		if err != nil {
			return err
		}
		result.CreatedBranch = requestedBranch != ""
		return s.assertStagedCommitPostcondition(ctx, repo.RootPath)
	}()
	return result, err
}

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

// CommitContextFromOptions builds pending hook context from staged commit options.
func CommitContextFromOptions(repo RepoInfo, branch string, opts StagedRevisionOptions) PendingCommitContext {
	return PendingCommitContext{
		WorktreeRoot:        repo.RootPath,
		GitCommonDir:        repo.GitCommonDir,
		Branch:              branch,
		RequestedBranch:     cleanRefName(opts.Branch),
		PreferredSessionIDs: opts.PreferredSessionIDs,
		SessionContexts:     opts.SessionContexts,
		SelfReport:          opts.SelfReport,
	}
}

func (s *Service) loadCommitResultFromHEAD(ctx context.Context, repo RepoInfo, commitOID string, pending PendingCommitContext) (CommitResult, error) {
	if result, err := s.commitResultFromDB(ctx, repo, commitOID); err == nil {
		result.SelfReport = pending.SelfReport
		return result, nil
	}
	return s.RecordGitCommit(ctx, repo, commitOID, pending)
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
