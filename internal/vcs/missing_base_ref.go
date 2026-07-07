package vcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
)

const MissingStackBaseRepairCommand = "gx doctor"

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
	RepairCommand          string                `json:"repair_command,omitempty"`
}

// MissingStackBaseRefFailure reports one stack that could not be rebased.
type MissingStackBaseRefFailure struct {
	Issue MissingStackBaseRef `json:"issue"`
	Error string              `json:"error"`
}

// RebaseOntoDefaultResult reports stacks rebased after a missing-base repair.
type RebaseOntoDefaultResult struct {
	RepoRoot string                       `json:"repo_root"`
	Fixed    []MissingStackBaseRef        `json:"fixed,omitempty"`
	Failed   []MissingStackBaseRefFailure `json:"failed,omitempty"`
	Actions  []string                     `json:"actions,omitempty"`
}

// ErrMissingStackBaseRefs blocks status until missing parent base refs are repaired or declined.
type ErrMissingStackBaseRefs struct {
	Status MissingStackBaseRefStatus
}

func (e *ErrMissingStackBaseRefs) Error() string {
	if len(e.Status.Issues) == 0 {
		return fmt.Sprintf("stack parent base ref is missing; run %s to repair", MissingStackBaseRepairCommand)
	}
	if len(e.Status.Issues) == 1 {
		issue := e.Status.Issues[0]
		return fmt.Sprintf(
			"missing parent base ref %q for stack %s (parent branch was likely merged); run %s to repair",
			issue.MissingBaseRef,
			firstNonEmpty(issue.BookmarkName, issue.Name),
			MissingStackBaseRepairCommand,
		)
	}
	return fmt.Sprintf("%d stacks store missing parent base refs; run %s to repair", len(e.Status.Issues), MissingStackBaseRepairCommand)
}

func (s *Service) ensureStackBaseRefsForPublish(ctx context.Context) error {
	status, err := s.DetectMissingStackBaseRefs(ctx)
	if err != nil {
		return err
	}
	if len(status.Issues) == 0 {
		return nil
	}
	return &ErrMissingStackBaseRefs{Status: status}
}

// RepairMissingStackBaseRefs rebases stacks with missing parent base refs onto the default branch.
func (s *Service) RepairMissingStackBaseRefs(ctx context.Context) (RebaseOntoDefaultResult, error) {
	status, err := s.DetectMissingStackBaseRefs(ctx)
	if err != nil {
		return RebaseOntoDefaultResult{}, err
	}
	if len(status.Issues) == 0 {
		return RebaseOntoDefaultResult{}, nil
	}
	return s.RebaseMissingStackBaseRefs(ctx, status.Issues)
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
		RepairCommand:          MissingStackBaseRepairCommand,
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

// RepairMissingStackBaseRefs detects and rebases stacks with missing parent base refs.
func (s *Service) RepairMissingStackBaseRefs(ctx context.Context) (RebaseOntoDefaultResult, error) {
	status, err := s.DetectMissingStackBaseRefs(ctx)
	if err != nil {
		return RebaseOntoDefaultResult{}, err
	}
	if len(status.Issues) == 0 {
		repo, repoErr := s.configuredJJRepo(ctx)
		if repoErr != nil {
			return RebaseOntoDefaultResult{}, repoErr
		}
		result := RebaseOntoDefaultResult{RepoRoot: repo.RootPath}
		return result, s.finishMissingStackBaseRepair(ctx, &result)
	}
	return s.RebaseMissingStackBaseRefs(ctx, status.Issues)
}

// RebaseMissingStackBaseRefs rebases affected stacks onto the default branch and updates stored base refs.
func (s *Service) RebaseMissingStackBaseRefs(ctx context.Context, issues []MissingStackBaseRef) (RebaseOntoDefaultResult, error) {
	repo, err := s.configuredJJRepo(ctx)
	if err != nil {
		return RebaseOntoDefaultResult{}, err
	}
	result := RebaseOntoDefaultResult{RepoRoot: repo.RootPath}
	if len(issues) == 0 {
		return result, s.finishMissingStackBaseRepair(ctx, &result)
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
		missingBases := map[string]struct{}{}
		for _, issue := range issues {
			missingBases[strings.TrimSpace(issue.MissingBaseRef)] = struct{}{}
		}
		for _, issue := range issues {
			if err := s.rebaseStackOntoDefaultUnlocked(ctx, store, repo, repoRow.ID, issue, now); err != nil {
				result.Failed = append(result.Failed, MissingStackBaseRefFailure{
					Issue: issue,
					Error: err.Error(),
				})
				continue
			}
			result.Fixed = append(result.Fixed, issue)
			result.Actions = append(result.Actions, fmt.Sprintf(
				"rebased stack %s onto %s and updated stored base ref from %s",
				firstNonEmpty(issue.BookmarkName, issue.Name),
				issue.DefaultBaseRef,
				issue.MissingBaseRef,
			))
		}
		if repoRow.AuthoringBase != nil {
			authoringBase := strings.TrimSpace(*repoRow.AuthoringBase)
			if authoringBase != "" {
				if _, missing := missingBases[authoringBase]; missing {
					defaultBase := firstNonEmpty(strings.TrimSpace(issues[0].DefaultBaseRef), repo.defaultBaseBranch())
					if authoringBase != defaultBase {
						if err := store.SetRepoAuthoringBase(ctx, repoRow.ID, defaultBase, now); err != nil {
							return err
						}
						result.Actions = append(result.Actions, fmt.Sprintf(
							"changed GX authoring base from %s to %s",
							authoringBase,
							defaultBase,
						))
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, s.finishMissingStackBaseRepair(ctx, &result)
}

func (s *Service) finishMissingStackBaseRepair(ctx context.Context, result *RebaseOntoDefaultResult) error {
	return withBusyRetry(ctx, "repair repo authoring base", func() error {
		store, err := openStore(ctx)
		if err != nil {
			return err
		}
		defer store.Close()
		repoRow, err := store.FindRepoByRoot(ctx, result.RepoRoot)
		if err != nil {
			return err
		}
		if repoRow == nil {
			return nil
		}
		repo, err := s.configuredJJRepo(ctx)
		if err != nil {
			return err
		}
		previousBase := ""
		if repoRow.AuthoringBase != nil {
			previousBase = strings.TrimSpace(*repoRow.AuthoringBase)
		}
		if err := s.repairRepoAuthoringBaseIfMissing(ctx, store, repo, repoRow.ID, time.Now().UnixMilli()); err != nil {
			return err
		}
		if previousBase != "" && previousBase != repo.defaultBaseBranch() {
			result.Actions = append(result.Actions, fmt.Sprintf("updated repo authoring base from %s to %s", previousBase, repo.defaultBaseBranch()))
		}
		return nil
	})
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
	remoteName := ""
	if stack, stackErr := store.FindStackByBookmark(ctx, repoID, bookmark); stackErr == nil && stack != nil && stack.RemoteName != nil {
		remoteName = strings.TrimSpace(*stack.RemoteName)
	}
	missingBaseExists := s.stackBaseRefExists(ctx, repo.RootPath, issue.MissingBaseRef, remoteName)
	sourceRevset := stackRebaseSourceRevset(bookmark, issue.MissingBaseRef, defaultBase, missingBaseExists)
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
	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return err
	}
	if repoRow != nil && repoRow.AuthoringBase != nil && strings.TrimSpace(*repoRow.AuthoringBase) == strings.TrimSpace(issue.MissingBaseRef) {
		if err := store.SetRepoAuthoringBase(ctx, repoID, defaultBase, now); err != nil {
			return err
		}
	}
	return s.reconcileRepoChanges(ctx, repo, defaultBase)
}

func (s *Service) repairRepoAuthoringBaseIfMissing(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64, now int64) error {
	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return err
	}
	if repoRow == nil || repoRow.AuthoringBase == nil {
		return nil
	}
	baseRef := strings.TrimSpace(*repoRow.AuthoringBase)
	if baseRef == "" {
		return nil
	}
	defaultBase := repo.defaultBaseBranch()
	if baseRef == defaultBase || s.stackBaseRefExists(ctx, repo.RootPath, baseRef, repoRemoteName(repo)) {
		return nil
	}
	if err := store.SetRepoAuthoringBase(ctx, repoID, defaultBase, now); err != nil {
		return err
	}
	return store.RenameStackBaseRef(ctx, repoID, baseRef, defaultBase, now)
}

func stackRebaseSourceRevset(bookmark, missingBaseRef, defaultBase string, missingBaseExists bool) string {
	bookmark = strings.TrimSpace(bookmark)
	missingBaseRef = strings.TrimSpace(missingBaseRef)
	defaultBase = strings.TrimSpace(defaultBase)
	revset := fmt.Sprintf("ancestors(%s) & mutable() & ~empty() & ~hidden()", quoteJJRev(bookmark))
	if missingBaseExists && missingBaseRef != "" && missingBaseRef != defaultBase {
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
	summary.RepairCommand = status.RepairCommand
}
