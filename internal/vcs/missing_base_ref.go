package vcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
)

const (
	MissingStackBaseFixAction = "rebase_onto_default"
)

// MissingStackBaseRef describes a stack whose stored parent base ref no longer exists.
type MissingStackBaseRef struct {
	StackID        int64  `json:"stack_id"`
	Name           string `json:"name"`
	BookmarkName   string `json:"bookmark_name"`
	MissingBaseRef string `json:"missing_base_ref"`
	DefaultBaseRef string `json:"default_base_ref"`
}

// MissingStackBaseRefStatus is attached to status output when parent base refs are gone.
type MissingStackBaseRefStatus struct {
	Issues                 []MissingStackBaseRef `json:"issues,omitempty"`
	NeedsRebaseOntoDefault bool                  `json:"needs_rebase_onto_default,omitempty"`
	FixAction              string                `json:"fix_action,omitempty"`
	FixPrompt              string                `json:"fix_prompt,omitempty"`
}

// RebaseOntoDefaultResult reports stacks rebased after a missing-base repair.
type RebaseOntoDefaultResult struct {
	RepoRoot string                `json:"repo_root"`
	Fixed    []MissingStackBaseRef `json:"fixed,omitempty"`
	Actions  []string              `json:"actions,omitempty"`
}

// ErrMissingStackBaseRefs blocks status until missing parent base refs are repaired or declined.
type ErrMissingStackBaseRefs struct {
	Status MissingStackBaseRefStatus
}

func (e *ErrMissingStackBaseRefs) Error() string {
	if len(e.Status.Issues) == 0 {
		return "stack parent base ref is missing"
	}
	if len(e.Status.Issues) == 1 {
		issue := e.Status.Issues[0]
		return fmt.Sprintf(
			"stack %s stores missing parent base ref %q (parent branch was likely merged); rebase onto %s to continue",
			firstNonEmpty(issue.BookmarkName, issue.Name),
			issue.MissingBaseRef,
			issue.DefaultBaseRef,
		)
	}
	return fmt.Sprintf("%d stacks store missing parent base refs; rebase onto the default branch to continue", len(e.Status.Issues))
}

func missingStackBaseFixPrompt(issues []MissingStackBaseRef) string {
	if len(issues) == 0 {
		return ""
	}
	if len(issues) == 1 {
		issue := issues[0]
		return fmt.Sprintf(
			"Rebase stack %s onto %s because parent base ref %q is missing (parent branch was likely merged)? [y/N]",
			firstNonEmpty(issue.BookmarkName, issue.Name),
			issue.DefaultBaseRef,
			issue.MissingBaseRef,
		)
	}
	defaultBase := issues[0].DefaultBaseRef
	return fmt.Sprintf(
		"Rebase %d stacks onto %s because their stored parent base refs are missing (parent branches were likely merged)? [y/N]",
		len(issues),
		defaultBase,
	)
}

func (s *Service) DetectMissingStackBaseRefs(ctx context.Context) (MissingStackBaseRefStatus, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return MissingStackBaseRefStatus{}, err
	}
	store, err := openStore(ctx)
	if err != nil {
		return MissingStackBaseRefStatus{}, err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return MissingStackBaseRefStatus{}, err
	}
	if repoRow == nil {
		return MissingStackBaseRefStatus{}, nil
	}
	issues, err := s.missingStackBaseRefsForRepo(ctx, store, repo, repoRow.ID)
	if err != nil {
		return MissingStackBaseRefStatus{}, err
	}
	return missingStackBaseRefStatusFromIssues(issues), nil
}

func missingStackBaseRefStatusFromIssues(issues []MissingStackBaseRef) MissingStackBaseRefStatus {
	if len(issues) == 0 {
		return MissingStackBaseRefStatus{}
	}
	return MissingStackBaseRefStatus{
		Issues:                 issues,
		NeedsRebaseOntoDefault: true,
		FixAction:              MissingStackBaseFixAction,
		FixPrompt:              missingStackBaseFixPrompt(issues),
	}
}

func (s *Service) missingStackBaseRefsForRepo(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64) ([]MissingStackBaseRef, error) {
	stored, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}
	defaultBase := repo.defaultBaseBranch()
	issues := make([]MissingStackBaseRef, 0)
	for _, stack := range stored {
		if IsTerminalStackStatus(stack.Status) {
			continue
		}
		baseRef := strings.TrimSpace(stack.BaseRef)
		if baseRef == "" || baseRef == defaultBase {
			continue
		}
		remoteName := ""
		if stack.RemoteName != nil {
			remoteName = strings.TrimSpace(*stack.RemoteName)
		}
		if s.stackBaseRefExists(ctx, repo.RootPath, baseRef, remoteName) {
			continue
		}
		bookmark := strings.TrimSpace(stack.BookmarkName)
		if bookmark == "" {
			continue
		}
		exists, err := s.RevisionExists(ctx, repo.RootPath, bookmark)
		if err != nil {
			return nil, err
		}
		if !exists {
			continue
		}
		issues = append(issues, MissingStackBaseRef{
			StackID:        stack.ID,
			Name:           stack.Name,
			BookmarkName:   bookmark,
			MissingBaseRef: baseRef,
			DefaultBaseRef: defaultBase,
		})
	}
	return issues, nil
}

func (s *Service) stackBaseRefExists(ctx context.Context, repoRoot, baseRef, remoteName string) bool {
	baseRef = strings.TrimSpace(baseRef)
	if baseRef == "" {
		return false
	}
	if s.revExists(ctx, repoRoot, baseRef) {
		return true
	}
	if strings.Contains(baseRef, "@") {
		return false
	}
	remoteName = strings.TrimSpace(remoteName)
	if remoteName == "" {
		remoteName = "origin"
	}
	return s.revExists(ctx, repoRoot, baseRef+"@"+remoteName)
}

// RebaseMissingStackBaseRefs rebases affected stacks onto the default branch and updates stored base refs.
func (s *Service) RebaseMissingStackBaseRefs(ctx context.Context, issues []MissingStackBaseRef) (RebaseOntoDefaultResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return RebaseOntoDefaultResult{}, err
	}
	result := RebaseOntoDefaultResult{RepoRoot: repo.RootPath}
	if len(issues) == 0 {
		return result, nil
	}
	err = withRepoLock(repo.RootPath, func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
		if err != nil {
			return err
		}
		if repoRow == nil {
			return fmt.Errorf("repo not registered in gx storage")
		}
		now := time.Now().UnixMilli()
		for _, issue := range issues {
			if err := s.rebaseStackOntoDefaultUnlocked(ctx, store, repo, repoRow.ID, issue, now); err != nil {
				return err
			}
			result.Fixed = append(result.Fixed, issue)
			result.Actions = append(result.Actions, fmt.Sprintf(
				"rebased stack %s onto %s and updated stored base ref from %s",
				firstNonEmpty(issue.BookmarkName, issue.Name),
				issue.DefaultBaseRef,
				issue.MissingBaseRef,
			))
		}
		return nil
	})
	return result, err
}

func (s *Service) rebaseStackOntoDefaultUnlocked(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64, issue MissingStackBaseRef, now int64) error {
	bookmark := strings.TrimSpace(issue.BookmarkName)
	defaultBase := firstNonEmpty(strings.TrimSpace(issue.DefaultBaseRef), repo.defaultBaseBranch())
	if bookmark == "" {
		return fmt.Errorf("missing bookmark for stack %q", issue.Name)
	}
	if !s.revExists(ctx, repo.RootPath, defaultBase) {
		return fmt.Errorf("default base ref %q does not exist", defaultBase)
	}
	sourceRevset := stackRebaseSourceRevset(bookmark, issue.MissingBaseRef, defaultBase)
	if _, err := s.runner.Run(ctx, repo.RootPath, "jj", "rebase", "-s", sourceRevset, "-d", defaultBase); err != nil {
		return fmt.Errorf("rebase stack %s onto %s: %w", bookmark, defaultBase, err)
	}
	stack, err := store.FindStackByBookmark(ctx, repoID, bookmark)
	if err != nil {
		return err
	}
	if stack == nil {
		return fmt.Errorf("stack %q not found in gx storage", bookmark)
	}
	stack.BaseRef = defaultBase
	stack.BaseCommitID = s.stackBaseCommitID(ctx, repo.RootPath, defaultBase)
	if change, changeErr := s.CurrentChange(ctx, repo.RootPath, bookmark); changeErr == nil {
		headChangeID := change.ChangeID
		headCommitID := change.CommitID
		stack.HeadChangeID = &headChangeID
		stack.HeadCommitID = &headCommitID
	}
	stack.UpdatedAt = now
	if _, err := store.UpsertStack(ctx, *stack); err != nil {
		return err
	}
	if err := store.RenameStackBaseRef(ctx, repoID, issue.MissingBaseRef, defaultBase, now); err != nil {
		return err
	}
	return s.reconcileRepoChanges(ctx, repo, defaultBase)
}

func stackRebaseSourceRevset(bookmark, missingBaseRef, defaultBase string) string {
	bookmark = strings.TrimSpace(bookmark)
	missingBaseRef = strings.TrimSpace(missingBaseRef)
	defaultBase = strings.TrimSpace(defaultBase)
	revset := fmt.Sprintf("ancestors(%s) & mutable() & ~empty() & ~hidden()", quoteJJRev(bookmark))
	if missingBaseRef != "" && missingBaseRef != defaultBase {
		revset += " & ~ancestors(" + quoteJJRev(missingBaseRef) + ")"
	}
	if defaultBase != "" {
		revset += " & ~ancestors(" + quoteJJRev(defaultBase) + ")"
	}
	return revset
}

func attachMissingStackBaseRefStatus(summary *StackSummary, status MissingStackBaseRefStatus) {
	if summary == nil || len(status.Issues) == 0 {
		return
	}
	summary.MissingBaseRefs = status.Issues
	summary.NeedsRebaseOntoDefault = status.NeedsRebaseOntoDefault
	summary.FixAction = status.FixAction
	summary.FixPrompt = status.FixPrompt
}
