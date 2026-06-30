package authoring

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/provenance"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/structural"
)

var hunkHeaderPattern = regexp.MustCompile(`@@ -([0-9]+)(?:,([0-9]+))? \+([0-9]+)(?:,([0-9]+))? @@`)

type demuxGroup struct {
	Files      []string
	Hunks      []HunkRange
	UseHunks   bool
	Symbol     string
	SymbolKind string
	Confidence float64
}

func (e *Engine) ProposeDemux(ctx context.Context, opts ProposeDemuxOptions) (DemuxProposal, error) {
	repo, err := e.vcs.ResolveJJRepo(ctx)
	if err != nil {
		return DemuxProposal{}, err
	}
	workRev, err := e.vcs.WorkRev(ctx, repo.RootPath)
	if err != nil {
		return DemuxProposal{}, err
	}
	current, err := e.vcs.CurrentChange(ctx, repo.RootPath, workRev)
	if err != nil {
		return DemuxProposal{}, err
	}
	files := cleanFiles(current.Files)
	files, err = filterFilesByFilesets(files, opts.Filesets)
	if err != nil {
		return DemuxProposal{}, err
	}
	files, err = excludeFilesByFilesets(files, opts.ExcludeFilesets)
	if err != nil {
		return DemuxProposal{}, err
	}
	if len(files) == 0 {
		return DemuxProposal{}, fmt.Errorf("no working-copy changes to demux")
	}
	diff, err := e.vcs.DiffGit(ctx, repo.RootPath, workRev)
	if err != nil {
		return DemuxProposal{}, err
	}
	hunks := parseGitHunks(diff)
	hunks = filterHunksByFiles(hunks, files)

	anchor, err := e.wipePendingDemuxProposals(ctx, repo, current.ChangeID)
	if err != nil {
		return DemuxProposal{}, err
	}

	sessionIDs, provenanceStatus, err := e.sessionIDsForProposal(ctx, repo.RootPath)
	if err != nil {
		return DemuxProposal{}, err
	}

	facts := structural.Analyze(repo.RootPath, files)
	changedSymbols := structural.EnclosingSymbols(facts, structuralHunkInputs(hunks))
	hunks = annotateHunksWithSymbols(hunks, changedSymbols)
	hunksByFile := hunksByFile(hunks)
	fileGroups := orderGroupsByStructuralDependencies(groupFiles(files), facts)
	groups := expandGroupsByChangedSymbols(fileGroups, hunksByFile)
	revisions := make([]RevisionProposal, 0, len(groups))
	for index, group := range groups {
		revisions = append(revisions, RevisionProposal{
			ID:               fmt.Sprintf("u%d", index+1),
			Intent:           proposalIntentForGroup(opts.Intent, group),
			Files:            group.Files,
			UseHunks:         group.UseHunks,
			HunkIDs:          hunkIDs(group.Hunks),
			Hunks:            group.Hunks,
			ProvenanceStatus: provenanceStatus,
			SessionIDs:       sessionIDs,
			Confidence:       group.Confidence,
		})
	}

	warnings := []string{
		"deterministic proposal uses file, test-pair, and changed-symbol heuristics; MCP/LLM callers may still regroup returned hunk ids",
		"intermediate stack buildability is not guaranteed; verify ordered revisions before apply",
	}
	if len(facts.Edges) == 0 {
		warnings = append(warnings, "no structural dependency edges found; proposal order falls back to file heuristics")
	} else {
		warnings = append(warnings, "lightweight structural dependency edges were used for initial revision ordering")
	}
	if provenanceStatus == "absent" {
		warnings = append(warnings, "no explicit or repo-local pending session found; proposed revisions will not have exact session provenance")
	} else if provenanceStatus == "repo_local" {
		warnings = append(warnings, "repo-local pending session provenance was inferred from captured sessions for this repository")
	}

	proposal := DemuxProposal{
		RepoRoot:         repo.RootPath,
		ProposedChangeID: current.ChangeID,
		ProposedCommitID: current.CommitID,
		Status:           ProposalPending,
		PlanInstructions: demuxPlanInstructions(),
		Hunks:            hunks,
		StructuralFacts:  structuralFacts(facts),
		StructuralDeps:   structuralDependencies(facts),
		ChangedSymbols:   changedSymbolsForProposal(changedSymbols),
		Revisions:        revisions,
		Warnings:         warnings,
	}
	proposal, _ = shapeDemuxProposalForReview(proposal)
	proposal, err = e.planDemuxRoutes(ctx, proposal)
	if err != nil {
		return DemuxProposal{}, err
	}
	proposal, _ = shapeDemuxProposalForReview(proposal)
	proposal.FeasibilityWarnings = feasibilityWarningsForProposal(proposal)
	if anchor.ID != "" {
		proposal.ID = anchor.ID
		proposal.CreatedAt = anchor.CreatedAt
	}
	proposal = normalizeConventionalDemuxStackRoutes(proposal)
	proposal, _ = shapeDemuxProposalForReview(proposal)
	proposal.FeasibilityWarnings = feasibilityWarningsForProposal(proposal)
	return e.SaveDemuxProposal(ctx, proposal)
}

type pendingDemuxAnchor struct {
	ID        string
	CreatedAt int64
}

func (e *Engine) wipePendingDemuxProposals(ctx context.Context, repo RepoInfo, proposedChangeID string) (pendingDemuxAnchor, error) {
	store, err := openStore(ctx)
	if err != nil {
		return pendingDemuxAnchor{}, err
	}
	defer store.Close()
	now := time.Now().UnixMilli()
	repoID, err := upsertProposalRepo(ctx, store, repo, now)
	if err != nil {
		return pendingDemuxAnchor{}, err
	}
	anchor := pendingDemuxAnchor{}
	if row, err := store.FindLatestDemuxProposal(ctx, repoID, string(ProposalPending)); err == nil && row != nil {
		if strings.TrimSpace(row.BaseChangeID) == strings.TrimSpace(proposedChangeID) {
			anchor.ID = row.ID
			anchor.CreatedAt = row.CreatedAt
		}
	}
	if err := store.DeletePendingDemuxProposalsForRepo(ctx, repoID); err != nil {
		return pendingDemuxAnchor{}, err
	}
	return anchor, nil
}

func (e *Engine) DemuxChanges(ctx context.Context, opts ProposeDemuxOptions) (DemuxPlanPacket, error) {
	return e.demuxPipeline().proposeChanges(ctx, opts)
}

func (e *Engine) ShowDemuxProposal(ctx context.Context, selector string) (DemuxPlanPacket, error) {
	return e.demuxPipeline().showProposal(ctx, selector)
}

func (e *Engine) ShowDemuxRevision(ctx context.Context, proposalID, revisionID string) (DemuxRevisionView, error) {
	revisionID = strings.TrimSpace(revisionID)
	if revisionID == "" {
		return DemuxRevisionView{}, fmt.Errorf("revision id is required")
	}
	var (
		proposal DemuxProposal
		err      error
	)
	if strings.TrimSpace(proposalID) == "" {
		proposal, err = e.LoadLatestPendingDemuxProposal(ctx)
	} else {
		proposal, err = e.LoadDemuxProposalBySelector(ctx, proposalID)
	}
	if err != nil {
		return DemuxRevisionView{}, err
	}
	proposal, err = normalizeDemuxProposal(proposal)
	if err != nil {
		return DemuxRevisionView{}, err
	}
	for _, revision := range proposal.Revisions {
		if revision.ID != revisionID {
			continue
		}
		files := cleanFiles(revision.Files)
		hunkIDs := revision.HunkIDs
		hunks := revision.Hunks
		if len(hunks) == 0 && len(hunkIDs) > 0 {
			for _, id := range hunkIDs {
				if hunk := findRevisionHunkByID(proposal.Hunks, id); hunk.ID != "" {
					hunks = append(hunks, hunk)
				}
			}
		}
		if len(hunks) == 0 {
			hunks = hunksForFiles(hunksByFile(proposal.Hunks), files)
		}
		return DemuxRevisionView{
			ProposalID:          proposal.ID,
			ProposalStatus:      proposal.Status,
			ProposedChangeID:    proposal.ProposedChangeID,
			ProposedCommitID:    proposal.ProposedCommitID,
			Revision:            revision,
			Hunks:               hunks,
			FeasibilityWarnings: relevantFeasibilityWarnings(proposal.FeasibilityWarnings, revision.ID),
			StructuralFacts:     relevantStructuralFacts(proposal.StructuralFacts, files),
			StructuralDeps:      relevantStructuralDependencies(proposal.StructuralDeps, files),
			ChangedSymbols:      relevantChangedSymbols(proposal.ChangedSymbols, files, hunkIDs),
		}, nil
	}
	return DemuxRevisionView{}, fmt.Errorf("demux proposal %s has no revision %q", proposal.ID, revisionID)
}

func demuxWorkflowState(review ReviewDemuxResult, proposal DemuxProposal) DemuxWorkflowState {
	if len(review.Errors) > 0 || len(blockingFeasibilityWarnings(proposal.FeasibilityWarnings)) > 0 {
		return DemuxWorkflowRepairRequired
	}
	return DemuxWorkflowReadyToApply
}

func (e *Engine) sessionIDsForProposal(ctx context.Context, repoRoot string) ([]string, string, error) {
	store, err := openStore(ctx)
	if err != nil {
		return nil, "", err
	}
	defer store.Close()
	attachment, err := provenance.Resolve(ctx, store, repoRoot)
	if err != nil {
		return nil, "", err
	}
	return attachment.SessionIDs, attachment.Status, nil
}

func (e *Engine) ApplyDemuxProposal(ctx context.Context, proposalID string) (ApplyDemuxResult, error) {
	return e.ApplyDemuxProposalWithOptions(ctx, proposalID, ApplyDemuxOptions{})
}

func (e *Engine) ApplyDemuxProposalWithOptions(ctx context.Context, proposalID string, opts ApplyDemuxOptions) (result ApplyDemuxResult, err error) {
	proposal, err := e.LoadDemuxProposal(ctx, proposalID)
	if err != nil {
		return ApplyDemuxResult{}, err
	}
	if proposal.Status != ProposalPending {
		return ApplyDemuxResult{}, fmt.Errorf("demux proposal %s is %s, not pending", proposal.ID, proposal.Status)
	}
	if len(proposal.Revisions) == 0 {
		return ApplyDemuxResult{}, fmt.Errorf("demux proposal %s has no revisions", proposal.ID)
	}
	workRev, current, err := e.currentDemuxApplyChange(ctx, proposal.RepoRoot)
	if err != nil {
		return ApplyDemuxResult{}, err
	}
	if current.ChangeID != proposal.ProposedChangeID || current.CommitID != proposal.ProposedCommitID {
		return ApplyDemuxResult{}, fmt.Errorf(
			"demux proposal %s is stale: current revision is %s/%s, proposed revision is %s/%s",
			proposal.ID,
			current.ChangeID,
			shortID(current.CommitID, 12),
			proposal.ProposedChangeID,
			shortID(proposal.ProposedCommitID, 12),
		)
	}
	proposal, err = normalizeDemuxProposal(proposal)
	if err != nil {
		return ApplyDemuxResult{}, err
	}
	proposal.FeasibilityWarnings = feasibilityWarningsForProposal(proposal)
	routeReview, err := e.reviewDemuxRoutes(ctx, proposal)
	if err != nil {
		return ApplyDemuxResult{}, err
	}
	if len(routeReview.Errors) > 0 {
		return ApplyDemuxResult{}, fmt.Errorf("demux proposal %s has invalid route(s): %s", proposal.ID, strings.Join(routeReview.Errors, "; "))
	}
	proposal.FeasibilityWarnings = append(proposal.FeasibilityWarnings, routeReview.Warnings...)
	if !opts.AllowWarnings {
		if warnings := blockingFeasibilityWarnings(proposal.FeasibilityWarnings); len(warnings) > 0 {
			return ApplyDemuxResult{}, fmt.Errorf("compose proposal %s has %d blocking feasibility warning(s); run `gx compose review %s --raw` to inspect them, `gx compose fix %s` to ask AI to repair them, or rerun apply with --allow-warnings after accepting the risk", proposal.ID, len(warnings), proposal.ID, proposal.ID)
		}
	}
	if _, err := e.SaveDemuxProposal(ctx, proposal); err != nil {
		return ApplyDemuxResult{}, err
	}
	source, err := e.demuxApplySourceLocation(ctx, current)
	if err != nil {
		return ApplyDemuxResult{}, err
	}
	if err := e.editDemuxApplySourceIfNeeded(ctx, proposal.RepoRoot, workRev, source); err != nil {
		return ApplyDemuxResult{}, err
	}
	if opts.ReturnToDefaultBranch {
		if err := source.returnToDefaultBranch(ctx, e, proposal.RepoRoot); err != nil {
			return ApplyDemuxResult{}, err
		}
	}
	defer func() {
		restoreSource := source
		if err != nil {
			restoreSource.PreserveVisibleCurrent = true
			if editErr := e.vcs.EditWorkingCopyRevision(ctx, proposal.RepoRoot, restoreSource.changeRef()); editErr != nil {
				err = fmt.Errorf("%w; additionally failed to restore source revision %s: %v", err, shortID(restoreSource.ChangeID, 12), editErr)
			}
		}
		if restoreErr := restoreSource.restoreAfterSuccessfulApply(ctx, e, proposal.RepoRoot); restoreErr != nil {
			if err != nil {
				err = fmt.Errorf("%w; additionally failed to return visible Git checkout to %s: %v", err, restoreSource.visibleCheckoutLabel(), restoreErr)
				return
			}
			err = fmt.Errorf("return visible Git checkout to %s: %w", restoreSource.visibleCheckoutLabel(), restoreErr)
		}
	}()

	results := make([]CheckpointResult, 0, len(proposal.Revisions))
	for index, revision := range proposal.Revisions {
		if strings.TrimSpace(revision.Intent) == "" {
			return ApplyDemuxResult{}, fmt.Errorf("revision %s has empty intent", revision.ID)
		}
		if len(revision.Files) == 0 {
			return ApplyDemuxResult{}, fmt.Errorf("revision %s has no files", revision.ID)
		}
		if hasNonCurrentDemuxRoute(revision) {
			result, err := e.applyRoutedDemuxRevision(ctx, proposal, revision)
			if err != nil {
				return ApplyDemuxResult{}, fmt.Errorf("apply revision %s: %w", revision.ID, err)
			}
			results = append(results, result)
			if err := e.writeDemuxEvidence(ctx, proposal, revision, result); err != nil {
				return ApplyDemuxResult{}, fmt.Errorf("write demux evidence for revision %s: %w", revision.ID, err)
			}
			continue
		}
		result, err := e.applyCurrentDemuxRevision(ctx, revision)
		if err != nil {
			return ApplyDemuxResult{}, fmt.Errorf("apply revision %s: %w", revision.ID, err)
		}
		results = append(results, result)
		if err := e.writeDemuxEvidence(ctx, proposal, revision, result); err != nil {
			return ApplyDemuxResult{}, fmt.Errorf("write demux evidence for revision %s: %w", revision.ID, err)
		}
		if index < len(proposal.Revisions)-1 {
			if err := e.prepareNextDemuxRevision(ctx, proposal.RepoRoot, result.Change.ChangeID); err != nil {
				return ApplyDemuxResult{}, fmt.Errorf("prepare next revision after %s: %w", revision.ID, err)
			}
		}
	}
	if err := e.MarkDemuxProposalStatus(ctx, proposal.ID, ProposalApplied); err != nil {
		return ApplyDemuxResult{}, err
	}
	remainingChanges := demuxProposalLeavesSourceFiles(current.Files, proposal)
	if remainingChanges && source.Stack == nil && strings.TrimSpace(source.GitCheckoutRef) != "" {
		if err := e.vcs.RestorePathsFromRevision(ctx, proposal.RepoRoot, source.GitCheckoutRef, demuxProposalCoveredFiles(proposal)); err != nil {
			return ApplyDemuxResult{}, fmt.Errorf("remove accepted compose files from remaining worktree: %w", err)
		}
		remainingChanges, err = e.demuxHasRemainingChanges(ctx, proposal.RepoRoot)
		if err != nil {
			return ApplyDemuxResult{}, fmt.Errorf("inspect remaining compose changes after cleanup: %w", err)
		}
	}
	source.PreserveVisibleCurrent = remainingChanges
	proposal.Status = ProposalApplied
	result = ApplyDemuxResult{
		Proposal:         proposal,
		Revisions:        results,
		RemainingChanges: remainingChanges,
		AcceptedSubset:   remainingChanges,
	}
	if remainingChanges {
		result.NextAction = "gx compose"
	}
	return result, nil
}

func (e *Engine) currentDemuxApplyChange(ctx context.Context, repoRoot string) (string, ChangeInfo, error) {
	workRev, err := e.vcs.WorkRev(ctx, repoRoot)
	if err != nil {
		return "", ChangeInfo{}, err
	}
	current, err := e.vcs.CurrentChange(ctx, repoRoot, workRev)
	if err != nil {
		return "", ChangeInfo{}, err
	}
	return workRev, current, nil
}

func (e *Engine) editDemuxApplySourceIfNeeded(ctx context.Context, repoRoot, workRev string, source demuxSourceLocation) error {
	if strings.TrimSpace(workRev) == "@" {
		return nil
	}
	if err := e.vcs.EditWorkingCopyRevision(ctx, repoRoot, source.changeRef()); err != nil {
		return fmt.Errorf("edit compose source revision %s: %w", shortID(source.ChangeID, 12), err)
	}
	return nil
}

func (e *Engine) applyCurrentDemuxRevision(ctx context.Context, revision RevisionProposal) (CheckpointResult, error) {
	preferred := revision.SessionIDs
	if revision.UseHunks && len(revision.Hunks) > 0 {
		patchFile, cleanup, err := writeRevisionPatch(revision.Hunks)
		if err != nil {
			return CheckpointResult{}, err
		}
		defer cleanup()
		return e.Checkpoint(ctx, CheckpointOptions{
			Intent:              revision.Intent,
			Hunk:                true,
			PatchFile:           patchFile,
			PreferredSessionIDs: preferred,
		})
	}
	return e.Checkpoint(ctx, CheckpointOptions{
		Intent:              revision.Intent,
		Filesets:            revision.Files,
		PreferredSessionIDs: preferred,
	})
}

func (e *Engine) applyRoutedDemuxRevision(ctx context.Context, proposal DemuxProposal, revision RevisionProposal) (CheckpointResult, error) {
	target := routeTargetStack(revision)
	router, err := e.demuxRouter(ctx, proposal.RepoRoot)
	if err != nil {
		return CheckpointResult{}, err
	}
	targetStack := router.lookup(target)
	targetBookmark := ""
	if targetStack != nil {
		targetBookmark = targetStack.BookmarkName
	}
	sourceChange, err := e.vcs.CurrentChange(ctx, proposal.RepoRoot, "@")
	if err != nil {
		return CheckpointResult{}, err
	}
	source, err := e.demuxSourceLocation(ctx, sourceChange)
	if err != nil {
		return CheckpointResult{}, err
	}
	if source.matches(target) {
		return e.applyCurrentDemuxRevision(ctx, revision)
	}
	files := cleanFiles(revision.Files)
	if len(files) == 0 {
		return CheckpointResult{}, fmt.Errorf("route requires files for revision %s", revision.ID)
	}
	sourceRev := firstNonEmpty(proposal.ProposedCommitID, sourceChange.ChangeID)
	sourceIgnored := demuxPreflightIgnoredPaths(proposal.RepoRoot)

	restoreSource := func() error {
		if err := source.restore(ctx, e, proposal.RepoRoot); err != nil {
			return err
		}
		if err := e.vcs.RestorePathsFromRevision(ctx, proposal.RepoRoot, sourceChange.ChangeID, files); err != nil {
			return fmt.Errorf("restore routed files to source %s: %w", source.label(), err)
		}
		return nil
	}
	restoreAfterError := func(cause error) error {
		if restoreErr := restoreSource(); restoreErr != nil {
			return fmt.Errorf("%w; additionally failed to restore source %s: %v", cause, source.label(), restoreErr)
		}
		return cause
	}

	if targetStack != nil {
		if _, err := e.Switch(ctx, target); err != nil {
			return CheckpointResult{}, fmt.Errorf("switch to routed target %q: %w", target, err)
		}
		if err := e.vcs.NewRevisionChild(ctx, proposal.RepoRoot); err != nil {
			return CheckpointResult{}, restoreAfterError(fmt.Errorf("prepare empty child on routed target %q: %w", target, err))
		}
	} else {
		baseRef := firstNonEmpty(routeBaseStack(revision), "base")
		resolvedBase, ok, err := e.BaseSwitchTarget(ctx, baseRef)
		if err != nil {
			return CheckpointResult{}, fmt.Errorf("resolve routed base %q: %w", baseRef, err)
		}
		if !ok {
			if _, err := e.Switch(ctx, baseRef); err != nil {
				return CheckpointResult{}, fmt.Errorf("switch to routed base %q: %w", baseRef, err)
			}
			resolvedBase = baseRef
		}
		if err := e.vcs.NewRevisionFrom(ctx, proposal.RepoRoot, resolvedBase); err != nil {
			return CheckpointResult{}, fmt.Errorf("prepare empty child on routed base %q: %w", baseRef, err)
		}
	}
	targetBefore, err := e.vcs.CurrentChange(ctx, proposal.RepoRoot, "@")
	if err != nil {
		return CheckpointResult{}, restoreAfterError(err)
	}
	targetFiles := demuxTargetBlockingFiles(targetBefore.Files, files, sourceIgnored)
	if len(targetFiles) > 0 {
		return CheckpointResult{}, restoreAfterError(fmt.Errorf("target stack %q has existing working-copy changes: %s", target, strings.Join(targetFiles, ", ")))
	}
	if err := e.vcs.RestorePathsFromRevision(ctx, proposal.RepoRoot, sourceRev, files); err != nil {
		return CheckpointResult{}, restoreAfterError(fmt.Errorf("copy routed files to target stack %q: %w", target, err))
	}

	var result CheckpointResult
	if targetStack != nil {
		result, err = e.vcs.RecordCurrentRevisionInStack(ctx, targetBookmark, revision.Intent, revision.SessionIDs)
	} else {
		result, err = e.vcs.RecordCurrentRevisionInNewStack(ctx, newStackRouteName(revision), newStackRouteBookmark(revision), routeBaseStack(revision), revision.Intent, revision.SessionIDs)
	}
	if err != nil {
		return CheckpointResult{}, restoreAfterError(err)
	}
	if err := source.restore(ctx, e, proposal.RepoRoot); err != nil {
		return CheckpointResult{}, fmt.Errorf("return to source %s: %w", source.label(), err)
	}
	if err := e.vcs.RestorePathsFromRevision(ctx, proposal.RepoRoot, "@-", files); err != nil {
		return CheckpointResult{}, fmt.Errorf("remove routed files from source %s: %w", source.label(), err)
	}
	return result, nil
}

func demuxTargetBlockingFiles(targetFiles []string, revisionFiles []string, sourceIgnored map[string]struct{}) []string {
	if len(targetFiles) == 0 {
		return nil
	}
	protected := map[string]struct{}{}
	for _, file := range revisionFiles {
		file = filepath.Clean(strings.TrimSpace(file))
		if file != "" && file != "." {
			protected[file] = struct{}{}
		}
	}
	blocking := make([]string, 0, len(targetFiles))
	for _, file := range targetFiles {
		file = filepath.Clean(strings.TrimSpace(file))
		if file == "" || file == "." {
			continue
		}
		if isProtectedPath(protected, file) {
			blocking = append(blocking, file)
			continue
		}
		if demuxPreflightIgnoredPath(file, sourceIgnored) {
			continue
		}
		blocking = append(blocking, file)
	}
	return blocking
}

type demuxSourceLocation struct {
	ChangeID               string
	CommitID               string
	GitCheckoutRef         string
	Stack                  *StackInfo
	ForceBase              bool
	PreserveVisibleCurrent bool
}

func (e *Engine) demuxSourceLocation(ctx context.Context, change ChangeInfo) (demuxSourceLocation, error) {
	status, err := e.Status(ctx)
	if err != nil {
		return demuxSourceLocation{}, err
	}
	source := demuxSourceLocation{ChangeID: change.ChangeID, CommitID: change.CommitID}
	if status.Stack != nil && strings.TrimSpace(status.Stack.BookmarkName) != "" {
		stack := *status.Stack
		source.Stack = &stack
	}
	return source, nil
}

func (e *Engine) demuxApplySourceLocation(ctx context.Context, change ChangeInfo) (demuxSourceLocation, error) {
	status, err := e.Status(ctx)
	if err != nil {
		return demuxSourceLocation{}, err
	}
	base, err := e.Base(ctx)
	if err != nil {
		return demuxSourceLocation{}, err
	}
	if base.OnBase {
		checkoutRef := firstNonEmpty(base.CurrentRef, base.BaseRef)
		if strings.TrimSpace(base.BaseRef) != "" && checkoutRef != strings.TrimSpace(base.BaseRef) {
			checkoutRef = strings.TrimSpace(base.BaseRef)
		}
		return demuxSourceLocation{
			ChangeID:       change.ChangeID,
			CommitID:       change.CommitID,
			GitCheckoutRef: checkoutRef,
		}, nil
	}
	if status.Stack == nil {
		return demuxSourceLocation{
			ChangeID:       change.ChangeID,
			CommitID:       change.CommitID,
			GitCheckoutRef: base.BaseRef,
			ForceBase:      true,
		}, nil
	}
	return demuxSourceLocation{
		ChangeID:       change.ChangeID,
		CommitID:       change.CommitID,
		GitCheckoutRef: firstNonEmpty(base.CurrentRef, base.BaseRef),
		ForceBase:      true,
	}, nil
}

func (s demuxSourceLocation) matches(target string) bool {
	if s.Stack == nil {
		return false
	}
	target = strings.TrimSpace(target)
	return target == s.Stack.BookmarkName || target == s.Stack.Name || target == s.Stack.Alias
}

func (s demuxSourceLocation) restore(ctx context.Context, e *Engine, repoRoot string) error {
	if err := s.restoreCheckout(ctx, e, repoRoot); err != nil {
		return err
	}
	if s.ForceBase || s.Stack == nil || strings.TrimSpace(s.Stack.BookmarkName) == "" {
		return e.vcs.EditRevision(ctx, repoRoot, s.revisionRef())
	}
	return nil
}

func (s demuxSourceLocation) restoreCheckout(ctx context.Context, e *Engine, repoRoot string) error {
	if s.ForceBase {
		_, err := e.Switch(ctx, "base")
		return err
	}
	if !s.ForceBase && s.Stack != nil && strings.TrimSpace(s.Stack.BookmarkName) != "" {
		if _, err := e.Switch(ctx, s.Stack.BookmarkName); err != nil {
			return err
		}
		if err := e.vcs.EditWorkingCopyRevision(ctx, repoRoot, s.revisionRef()); err != nil {
			return err
		}
		return e.vcs.ForceGitCheckoutPreservingWorktree(ctx, repoRoot, s.Stack.BookmarkName)
	}
	return e.vcs.EditRevision(ctx, repoRoot, s.revisionRef())
}

func (s demuxSourceLocation) revisionRef() string {
	return firstNonEmpty(s.CommitID, s.ChangeID)
}

func (s demuxSourceLocation) changeRef() string {
	return firstNonEmpty(s.ChangeID, s.CommitID)
}

func (s demuxSourceLocation) restoreVisibleGitCheckout(ctx context.Context, e *Engine, repoRoot string) error {
	if strings.TrimSpace(s.GitCheckoutRef) == "" {
		return nil
	}
	if s.PreserveVisibleCurrent {
		return e.vcs.ForceGitCheckoutPreservingWorktree(ctx, repoRoot, s.GitCheckoutRef)
	}
	if s.Stack == nil && !s.PreserveVisibleCurrent {
		current, err := e.vcs.CurrentChange(ctx, repoRoot, "@")
		if err != nil {
			return err
		}
		if len(current.Files) == 0 {
			if err := e.vcs.NewRevisionFrom(ctx, repoRoot, s.GitCheckoutRef); err != nil {
				return err
			}
		}
		return e.vcs.ForceGitCheckoutPreservingWorktree(ctx, repoRoot, s.GitCheckoutRef)
	}
	return e.vcs.ForceGitCheckout(ctx, repoRoot, s.GitCheckoutRef)
}

func (s *demuxSourceLocation) returnToDefaultBranch(ctx context.Context, e *Engine, repoRoot string) error {
	base, err := e.Base(ctx)
	if err != nil {
		return err
	}
	defaultBranch := strings.TrimSpace(base.DefaultBranch)
	if defaultBranch == "" {
		defaultBranch = "main"
	}
	s.GitCheckoutRef = defaultBranch
	return nil
}

func (s demuxSourceLocation) restoreAfterSuccessfulApply(ctx context.Context, e *Engine, repoRoot string) error {
	if s.Stack == nil && strings.TrimSpace(s.GitCheckoutRef) != "" {
		return s.restoreVisibleGitCheckout(ctx, e, repoRoot)
	}
	if err := s.restoreCheckout(ctx, e, repoRoot); err != nil {
		return err
	}
	return s.restoreVisibleGitCheckout(ctx, e, repoRoot)
}

func (s demuxSourceLocation) label() string {
	if s.Stack != nil && strings.TrimSpace(s.Stack.BookmarkName) != "" {
		return "stack " + s.Stack.BookmarkName
	}
	return "revision " + shortID(s.ChangeID, 12)
}

func (s demuxSourceLocation) checkoutLabel() string {
	if s.Stack != nil && strings.TrimSpace(s.Stack.BookmarkName) != "" {
		return "stack " + s.Stack.BookmarkName
	}
	return "base"
}

func (s demuxSourceLocation) visibleCheckoutLabel() string {
	return firstNonEmpty(s.GitCheckoutRef, s.checkoutLabel())
}

func (e *Engine) demuxHasRemainingChanges(ctx context.Context, repoRoot string) (bool, error) {
	current, err := e.vcs.CurrentChange(ctx, repoRoot, "@")
	if err != nil {
		return false, err
	}
	return len(current.Files) > 0, nil
}

func demuxProposalLeavesSourceFiles(sourceFiles []string, proposal DemuxProposal) bool {
	covered := demuxProposalCoveredFileSet(proposal)
	for _, file := range sourceFiles {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}
		if _, ok := covered[file]; !ok {
			return true
		}
	}
	return false
}

func demuxProposalCoveredFiles(proposal DemuxProposal) []string {
	covered := demuxProposalCoveredFileSet(proposal)
	files := make([]string, 0, len(covered))
	for file := range covered {
		files = append(files, file)
	}
	sort.Strings(files)
	return files
}

func demuxProposalCoveredFileSet(proposal DemuxProposal) map[string]struct{} {
	covered := map[string]struct{}{}
	for _, revision := range proposal.Revisions {
		for _, file := range revision.Files {
			file = strings.TrimSpace(file)
			if file != "" {
				covered[file] = struct{}{}
			}
		}
		for _, hunk := range revision.Hunks {
			file := strings.TrimSpace(hunk.File)
			if file != "" {
				covered[file] = struct{}{}
			}
		}
		for _, hunkID := range revision.HunkIDs {
			for _, hunk := range proposal.Hunks {
				if hunk.ID != hunkID {
					continue
				}
				file := strings.TrimSpace(hunk.File)
				if file != "" {
					covered[file] = struct{}{}
				}
			}
		}
	}
	return covered
}

func routedRevisionHunks(proposal DemuxProposal, revision RevisionProposal) []HunkRange {
	if revision.UseHunks {
		return append([]HunkRange(nil), revision.Hunks...)
	}
	return hunksForFiles(hunksByFile(proposal.Hunks), revision.Files)
}

func blockingFeasibilityWarnings(warnings []FeasibilityWarning) []FeasibilityWarning {
	var out []FeasibilityWarning
	for _, warning := range warnings {
		if warning.Severity == "warning" {
			out = append(out, warning)
		}
	}
	return out
}

type demuxEvidencePayload struct {
	ProposalID          string                 `json:"proposal_id"`
	Revision            RevisionProposal       `json:"revision"`
	Files               []string               `json:"files"`
	HunkIDs             []string               `json:"hunk_ids,omitempty"`
	Hunks               []HunkRange            `json:"hunks,omitempty"`
	FeasibilityWarnings []FeasibilityWarning   `json:"feasibility_warnings,omitempty"`
	StructuralFacts     []StructuralFact       `json:"structural_facts,omitempty"`
	StructuralDeps      []StructuralDependency `json:"structural_dependencies,omitempty"`
	ChangedSymbols      []ChangedSymbol        `json:"changed_symbols,omitempty"`
}

func (e *Engine) writeDemuxEvidence(ctx context.Context, proposal DemuxProposal, revision RevisionProposal, result CheckpointResult) error {
	store, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer store.Close()

	repo, err := e.repoForProposal(ctx, proposal.RepoRoot)
	if err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	repoID, err := upsertProposalRepo(ctx, store, repo, now)
	if err != nil {
		return err
	}
	change, err := store.FindChangeByJJChangeID(ctx, repoID, result.Change.ChangeID)
	if err != nil {
		return err
	}
	var changeID int64
	if change != nil {
		changeID = change.ID
	} else {
		changeID, err = store.UpsertChange(ctx, storage.Change{
			RepoID:          repoID,
			JJChangeID:      result.Change.ChangeID,
			CurrentCommitID: result.Change.CommitID,
			Description:     result.Change.Description,
			ParentChangeID:  result.Change.ParentChangeID,
			Status:          "draft",
			FirstSeenAt:     now,
			UpdatedAt:       now,
		})
		if err != nil {
			return err
		}
	}

	files := cleanFiles(revision.Files)
	hunkIDs := append([]string(nil), revision.HunkIDs...)
	payload := demuxEvidencePayload{
		ProposalID:          proposal.ID,
		Revision:            revision,
		Files:               files,
		HunkIDs:             hunkIDs,
		Hunks:               revision.Hunks,
		FeasibilityWarnings: relevantFeasibilityWarnings(proposal.FeasibilityWarnings, revision.ID),
		StructuralFacts:     relevantStructuralFacts(proposal.StructuralFacts, files),
		StructuralDeps:      relevantStructuralDependencies(proposal.StructuralDeps, files),
		ChangedSymbols:      relevantChangedSymbols(proposal.ChangedSymbols, files, hunkIDs),
	}
	filesJSON, err := json.Marshal(files)
	if err != nil {
		return err
	}
	hunkIDsJSON, err := json.Marshal(hunkIDs)
	if err != nil {
		return err
	}
	evidenceJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return store.WriteChangeDemuxEvidence(ctx, storage.ChangeDemuxEvidence{
		ChangeID:           changeID,
		DemuxProposalID:    proposal.ID,
		RevisionProposalID: revision.ID,
		Intent:             revision.Intent,
		FilesJSON:          string(filesJSON),
		HunkIDsJSON:        string(hunkIDsJSON),
		UseHunks:           revision.UseHunks,
		Confidence:         revision.Confidence,
		ProvenanceStatus:   revision.ProvenanceStatus,
		EvidenceJSON:       string(evidenceJSON),
		CreatedAt:          now,
	})
}

func relevantFeasibilityWarnings(warnings []FeasibilityWarning, revisionID string) []FeasibilityWarning {
	var out []FeasibilityWarning
	for _, warning := range warnings {
		if warning.RevisionID == "" || warning.RevisionID == revisionID {
			out = append(out, warning)
		}
	}
	return out
}

func relevantStructuralFacts(facts []StructuralFact, files []string) []StructuralFact {
	fileSet := stringSet(files)
	var out []StructuralFact
	for _, fact := range facts {
		if _, ok := fileSet[fact.File]; ok {
			out = append(out, fact)
		}
	}
	return out
}

func relevantStructuralDependencies(deps []StructuralDependency, files []string) []StructuralDependency {
	fileSet := stringSet(files)
	var out []StructuralDependency
	for _, dep := range deps {
		_, from := fileSet[dep.FromFile]
		_, to := fileSet[dep.ToFile]
		if from || to {
			out = append(out, dep)
		}
	}
	return out
}

func relevantChangedSymbols(symbols []ChangedSymbol, files []string, hunkIDs []string) []ChangedSymbol {
	fileSet := stringSet(files)
	hunkSet := stringSet(hunkIDs)
	var out []ChangedSymbol
	for _, symbol := range symbols {
		_, fileMatch := fileSet[symbol.File]
		_, hunkMatch := hunkSet[symbol.HunkID]
		if fileMatch || hunkMatch {
			out = append(out, symbol)
		}
	}
	return out
}

func stringSet(values []string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out[value] = struct{}{}
		}
	}
	return out
}

func (e *Engine) prepareNextDemuxRevision(ctx context.Context, repoRoot, recordedChangeID string) error {
	current, err := e.vcs.CurrentChange(ctx, repoRoot, "@")
	if err != nil {
		return err
	}
	if len(current.Files) > 0 {
		return nil
	}
	parent, err := e.vcs.CurrentChange(ctx, repoRoot, "@-")
	if err != nil {
		return err
	}
	if parent.ChangeID == recordedChangeID || len(parent.Files) == 0 {
		return nil
	}
	return e.vcs.EditRevision(ctx, repoRoot, firstNonEmpty(parent.CommitID, parent.ChangeID))
}

func (e *Engine) ApplyDemuxPlan(ctx context.Context, proposal DemuxProposal) (ApplyDemuxResult, error) {
	return e.ApplyDemuxPlanWithOptions(ctx, proposal, ApplyDemuxOptions{})
}

func (e *Engine) ApplyDemuxPlanWithOptions(ctx context.Context, proposal DemuxProposal, opts ApplyDemuxOptions) (ApplyDemuxResult, error) {
	if strings.TrimSpace(proposal.ID) != "" {
		base, err := e.LoadDemuxProposal(ctx, proposal.ID)
		if err == nil {
			proposal = mergeDemuxPlan(base, proposal)
		} else if !strings.Contains(err.Error(), "unknown demux proposal") {
			return ApplyDemuxResult{}, err
		}
	}
	saved, err := e.SaveDemuxProposal(ctx, proposal)
	if err != nil {
		return ApplyDemuxResult{}, err
	}
	return e.ApplyDemuxProposalWithOptions(ctx, saved.ID, opts)
}

func (e *Engine) ReviewDemuxPlan(ctx context.Context, proposal DemuxProposal) (ReviewDemuxResult, error) {
	if strings.TrimSpace(proposal.ID) != "" {
		base, err := e.LoadDemuxProposal(ctx, proposal.ID)
		if err != nil {
			return ReviewDemuxResult{}, err
		}
		proposal = mergeDemuxPlan(base, proposal)
	}
	normalized, err := normalizeDemuxProposal(proposal)
	if err != nil {
		return ReviewDemuxResult{
			Valid:       false,
			Proposal:    proposal,
			Errors:      []string{err.Error()},
			RepairHints: repairHintsForReviewError(proposal, err.Error()),
		}, nil
	}
	normalized.FeasibilityWarnings = feasibilityWarningsForProposal(normalized)
	routeReview, err := e.reviewDemuxRoutes(ctx, normalized)
	if err != nil {
		return ReviewDemuxResult{}, err
	}
	normalized.FeasibilityWarnings = append(normalized.FeasibilityWarnings, routeReview.Warnings...)
	if len(routeReview.Errors) > 0 {
		return ReviewDemuxResult{
			Valid:       false,
			Proposal:    normalized,
			Errors:      routeReview.Errors,
			RepairHints: routeReview.Hints,
		}, nil
	}
	hints := append(repairHintsForWarnings(normalized.FeasibilityWarnings), routeReview.Hints...)
	return ReviewDemuxResult{
		Valid:       true,
		Proposal:    normalized,
		RepairHints: hints,
	}, nil
}

func mergeDemuxPlan(base, plan DemuxProposal) DemuxProposal {
	if strings.TrimSpace(plan.RepoRoot) == "" {
		plan.RepoRoot = base.RepoRoot
	}
	if strings.TrimSpace(plan.ProposedChangeID) == "" {
		plan.ProposedChangeID = base.ProposedChangeID
	}
	if strings.TrimSpace(plan.ProposedCommitID) == "" {
		plan.ProposedCommitID = base.ProposedCommitID
	}
	if plan.Status == "" {
		plan.Status = base.Status
	}
	if plan.CreatedAt == 0 {
		plan.CreatedAt = base.CreatedAt
	}
	if len(plan.PlanInstructions) == 0 {
		plan.PlanInstructions = base.PlanInstructions
	}
	if len(plan.Warnings) == 0 {
		plan.Warnings = base.Warnings
	}
	if len(plan.StructuralFacts) == 0 {
		plan.StructuralFacts = base.StructuralFacts
	}
	if len(plan.StructuralDeps) == 0 {
		plan.StructuralDeps = base.StructuralDeps
	}
	if len(plan.ChangedSymbols) == 0 {
		plan.ChangedSymbols = base.ChangedSymbols
	}
	if len(plan.Revisions) == 0 {
		plan.Revisions = base.Revisions
	}
	if len(plan.Hunks) == 0 {
		plan.Hunks = base.Hunks
		if len(plan.Revisions) > 0 && len(base.Revisions) > 0 {
			plan.Hunks = baseHunksForDemuxRevisions(base.Hunks, plan.Revisions)
		}
	}
	return plan
}

func baseHunksForDemuxRevisions(hunks []HunkRange, revisions []RevisionProposal) []HunkRange {
	if len(hunks) == 0 || len(revisions) == 0 {
		return nil
	}
	selectedHunks := map[string]struct{}{}
	selectedFiles := map[string]struct{}{}
	for _, revision := range revisions {
		for _, id := range revision.HunkIDs {
			if id = strings.TrimSpace(id); id != "" {
				selectedHunks[id] = struct{}{}
			}
		}
		for _, hunk := range revision.Hunks {
			if id := strings.TrimSpace(hunk.ID); id != "" {
				selectedHunks[id] = struct{}{}
			}
		}
		if !revision.UseHunks {
			for _, file := range revision.Files {
				if file = strings.TrimSpace(file); file != "" {
					selectedFiles[file] = struct{}{}
				}
			}
		}
	}
	out := make([]HunkRange, 0, len(hunks))
	for _, hunk := range hunks {
		if _, ok := selectedHunks[hunk.ID]; ok {
			out = append(out, hunk)
			continue
		}
		if _, ok := selectedFiles[hunk.File]; ok {
			out = append(out, hunk)
		}
	}
	return out
}

var (
	unassignedHunkErrorPattern = regexp.MustCompile(`^hunk ([^ ]+) in (.+) is not assigned to any revision$`)
	duplicateHunkErrorPattern  = regexp.MustCompile(`^hunk ([^ ]+) is assigned to both revision ([^ ]+) and revision ([^ ]+)$`)
	mixedCoverageErrorPattern  = regexp.MustCompile(`^hunk ([^ ]+) in (.+) is covered by hunk revision ([^ ]+) and whole-file revision ([^ ]+)$`)
)

func repairHintsForReviewError(proposal DemuxProposal, message string) []RepairHint {
	switch {
	case unassignedHunkErrorPattern.MatchString(message):
		match := unassignedHunkErrorPattern.FindStringSubmatch(message)
		return []RepairHint{{
			Kind:       "unassigned_hunk",
			HunkID:     match[1],
			File:       match[2],
			Candidates: revisionIDs(proposal.Revisions),
			Suggestion: fmt.Sprintf("assign hunk %s to an existing revision or create a new revision for %s", match[1], match[2]),
		}}
	case duplicateHunkErrorPattern.MatchString(message):
		match := duplicateHunkErrorPattern.FindStringSubmatch(message)
		return []RepairHint{{
			Kind:       "duplicate_hunk_assignment",
			HunkID:     match[1],
			Candidates: []string{match[2], match[3]},
			Suggestion: fmt.Sprintf("keep hunk %s in exactly one revision", match[1]),
		}}
	case mixedCoverageErrorPattern.MatchString(message):
		match := mixedCoverageErrorPattern.FindStringSubmatch(message)
		return []RepairHint{{
			Kind:       "mixed_hunk_and_whole_file",
			HunkID:     match[1],
			File:       match[2],
			Candidates: []string{match[3], match[4]},
			Suggestion: fmt.Sprintf("choose hunk-level or whole-file coverage for %s, not both", match[2]),
		}}
	default:
		return nil
	}
}

func repairHintsForWarnings(warnings []FeasibilityWarning) []RepairHint {
	var hints []RepairHint
	for _, warning := range warnings {
		switch warning.Source {
		case "inferred_dependency":
			dependsOn := strings.TrimSpace(warning.DependsOn)
			symbol := strings.TrimSpace(warning.Symbol)
			if dependsOn == "" {
				dependsOn, symbol = parseInferredDependencyWarning(warning.Message)
			}
			if dependsOn == "" {
				continue
			}
			hints = append(hints, RepairHint{
				Kind:       "inferred_dependency",
				RevisionID: warning.RevisionID,
				DependsOn:  dependsOn,
				Symbol:     symbol,
				Suggestion: fmt.Sprintf("add %s to %s.depends_on if this structural dependency is intentional", dependsOn, warning.RevisionID),
			})
		case "structural_dependency":
			dependsOn := strings.TrimSpace(warning.DependsOn)
			if dependsOn == "" {
				continue
			}
			hints = append(hints, RepairHint{
				Kind:       "reorder_dependency",
				RevisionID: warning.RevisionID,
				DependsOn:  dependsOn,
				Symbol:     warning.Symbol,
				Suggestion: fmt.Sprintf("move revision %s before %s, or merge the dependent hunks into one revision", dependsOn, warning.RevisionID),
			})
		default:
			continue
		}
	}
	return hints
}

func parseInferredDependencyWarning(message string) (dependsOn string, symbol string) {
	fields := strings.Fields(message)
	for index, field := range fields {
		if field == "on" && index+1 < len(fields) {
			dependsOn = fields[index+1]
		}
		if field == "via" && index+1 < len(fields) {
			symbol = strings.TrimRight(fields[index+1], ";")
		}
	}
	return dependsOn, symbol
}

func revisionIDs(revisions []RevisionProposal) []string {
	ids := make([]string, 0, len(revisions))
	for _, revision := range revisions {
		id := strings.TrimSpace(revision.ID)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func StripDependsOnOutsideRevisions(revisions []RevisionProposal) []RevisionProposal {
	if len(revisions) == 0 {
		return revisions
	}
	ids := map[string]struct{}{}
	for _, revision := range revisions {
		if id := strings.TrimSpace(revision.ID); id != "" {
			ids[id] = struct{}{}
		}
	}
	out := make([]RevisionProposal, len(revisions))
	for index, revision := range revisions {
		var kept []string
		for _, dep := range revision.DependsOn {
			dep = strings.TrimSpace(dep)
			if dep == "" {
				continue
			}
			if _, ok := ids[dep]; ok {
				kept = append(kept, dep)
			}
		}
		revision.DependsOn = kept
		out[index] = revision
	}
	return out
}

func normalizeDemuxProposal(proposal DemuxProposal) (DemuxProposal, error) {
	hunksByID := map[string]HunkRange{}
	for _, hunk := range proposal.Hunks {
		if strings.TrimSpace(hunk.ID) == "" {
			continue
		}
		hunksByID[hunk.ID] = hunk
	}
	usedHunks := map[string]string{}
	wholeFileOwners := map[string]string{}
	revisionIDs := map[string]struct{}{}
	for index := range proposal.Revisions {
		revision := &proposal.Revisions[index]
		revision.ID = strings.TrimSpace(revision.ID)
		if revision.ID == "" {
			return DemuxProposal{}, fmt.Errorf("revision at position %d has empty id", index+1)
		}
		if _, ok := revisionIDs[revision.ID]; ok {
			return DemuxProposal{}, fmt.Errorf("duplicate revision id %s", revision.ID)
		}
		revisionIDs[revision.ID] = struct{}{}
		if revision.UseHunks && len(revision.HunkIDs) == 0 && len(revision.Hunks) > 0 {
			revision.HunkIDs = hunkIDs(revision.Hunks)
		}
		if revision.UseHunks && len(revision.HunkIDs) > 0 {
			resolved := make([]HunkRange, 0, len(revision.HunkIDs))
			for _, id := range revision.HunkIDs {
				id = strings.TrimSpace(id)
				if id == "" {
					continue
				}
				if owner, ok := usedHunks[id]; ok {
					return DemuxProposal{}, fmt.Errorf("hunk %s is assigned to both revision %s and revision %s", id, owner, revision.ID)
				}
				hunk, ok := hunksByID[id]
				if !ok {
					hunk = findRevisionHunkByID(revision.Hunks, id)
				}
				if hunk.ID == "" || hunk.Patch == "" {
					return DemuxProposal{}, fmt.Errorf("revision %s references unknown hunk %s", revision.ID, id)
				}
				resolved = append(resolved, hunk)
				usedHunks[id] = revision.ID
			}
			revision.Hunks = resolved
			if len(revision.Files) == 0 {
				revision.Files = filesForHunks(resolved)
			}
		}
		if revision.UseHunks && len(revision.Hunks) == 0 {
			return DemuxProposal{}, fmt.Errorf("revision %s has use_hunks=true but no hunk_ids or hunks", revision.ID)
		}
		if !revision.UseHunks {
			for _, file := range revision.Files {
				file = strings.TrimSpace(file)
				if file == "" {
					continue
				}
				if owner, ok := wholeFileOwners[file]; ok {
					return DemuxProposal{}, fmt.Errorf("file %s is assigned to both revision %s and revision %s", file, owner, revision.ID)
				}
				wholeFileOwners[file] = revision.ID
			}
		}
		for _, dep := range revision.DependsOn {
			dep = strings.TrimSpace(dep)
			if dep == "" {
				continue
			}
			if _, ok := revisionIDs[dep]; !ok {
				return DemuxProposal{}, fmt.Errorf("revision %s depends on unknown or later revision %s", revision.ID, dep)
			}
		}
	}
	if err := reviewDemuxHunkCoverage(proposal.Hunks, usedHunks, wholeFileOwners); err != nil {
		return DemuxProposal{}, err
	}
	return proposal, nil
}

func reviewDemuxHunkCoverage(hunks []HunkRange, usedHunks map[string]string, wholeFileOwners map[string]string) error {
	if len(hunks) == 0 {
		return nil
	}
	for _, hunk := range hunks {
		if strings.TrimSpace(hunk.ID) == "" {
			continue
		}
		hunkOwner, hasHunkOwner := usedHunks[hunk.ID]
		fileOwner, hasFileOwner := wholeFileOwners[hunk.File]
		switch {
		case hasHunkOwner && hasFileOwner:
			return fmt.Errorf("hunk %s in %s is covered by hunk revision %s and whole-file revision %s", hunk.ID, hunk.File, hunkOwner, fileOwner)
		case hasHunkOwner || hasFileOwner:
			continue
		default:
			return fmt.Errorf("hunk %s in %s is not assigned to any revision", hunk.ID, hunk.File)
		}
	}
	return nil
}

func feasibilityWarnings(revisions []RevisionProposal, hunks []HunkRange, facts structural.Facts) []FeasibilityWarning {
	var warnings []FeasibilityWarning
	warnings = append(warnings, unmappedHunkWarnings(revisions, hunks)...)
	warnings = append(warnings, structuralOrderWarnings(revisions, facts)...)
	warnings = append(warnings, testSeparationWarnings(revisions)...)
	return warnings
}

func feasibilityWarningsForProposal(proposal DemuxProposal) []FeasibilityWarning {
	return feasibilityWarnings(proposal.Revisions, proposal.Hunks, structuralFactsForProposal(proposal))
}

func structuralFactsForProposal(proposal DemuxProposal) structural.Facts {
	facts := structural.Facts{
		Files: make([]structural.FileFact, 0, len(proposal.StructuralFacts)),
		Edges: make([]structural.DependencyEdge, 0, len(proposal.StructuralDeps)),
	}
	for _, fact := range proposal.StructuralFacts {
		facts.Files = append(facts.Files, structural.FileFact{
			File:              fact.File,
			Language:          fact.Language,
			DefinedSymbols:    append([]string(nil), fact.DefinedSymbols...),
			ReferencedSymbols: append([]string(nil), fact.ReferencedSymbols...),
			Symbols:           structuralSymbolsForProposal(fact.File, fact.Symbols),
		})
	}
	for _, dep := range proposal.StructuralDeps {
		facts.Edges = append(facts.Edges, structural.DependencyEdge{
			FromFile: dep.FromFile,
			ToFile:   dep.ToFile,
			Symbol:   dep.Symbol,
		})
	}
	return facts
}

func structuralSymbolsForProposal(file string, symbols []StructuralSymbol) []structural.Symbol {
	out := make([]structural.Symbol, 0, len(symbols))
	for _, symbol := range symbols {
		out = append(out, structural.Symbol{
			Name:      symbol.Name,
			Kind:      symbol.Kind,
			File:      file,
			StartLine: symbol.StartLine,
			EndLine:   symbol.EndLine,
		})
	}
	return out
}

func unmappedHunkWarnings(revisions []RevisionProposal, hunks []HunkRange) []FeasibilityWarning {
	if len(hunks) == 0 {
		return nil
	}
	mapped := map[string]struct{}{}
	for _, hunk := range hunks {
		if strings.TrimSpace(hunk.Symbol) != "" {
			mapped[hunk.ID] = struct{}{}
		}
	}
	if len(mapped) == len(hunks) {
		return nil
	}
	revisionByHunk := map[string]string{}
	for _, revision := range revisions {
		for _, id := range revision.HunkIDs {
			revisionByHunk[id] = revision.ID
		}
	}
	var warnings []FeasibilityWarning
	for _, hunk := range hunks {
		if _, ok := mapped[hunk.ID]; ok {
			continue
		}
		warnings = append(warnings, FeasibilityWarning{
			RevisionID: revisionByHunk[hunk.ID],
			Severity:   "info",
			Source:     "changed_symbol",
			Message:    fmt.Sprintf("hunk %s in %s is not mapped to an enclosing symbol; demux kept a coarser grouping", hunk.ID, hunk.File),
		})
	}
	return warnings
}

func structuralOrderWarnings(revisions []RevisionProposal, facts structural.Facts) []FeasibilityWarning {
	if len(revisions) < 2 || len(facts.Edges) == 0 {
		return nil
	}
	index := structuralRevisionIndex(revisions)
	var warnings []FeasibilityWarning
	seen := map[string]struct{}{}
	for _, edge := range facts.Edges {
		fromIndex, hasFrom := index.sourceRevision(edge.FromFile, edge.Symbol)
		toIndex, hasTo := index.targetRevision(edge.ToFile, edge.Symbol)
		if !hasFrom || !hasTo || fromIndex == toIndex || toIndex < fromIndex {
			if hasFrom && hasTo && fromIndex != toIndex && toIndex < fromIndex {
				if warning, ok := inferredDependencyWarning(revisions, fromIndex, toIndex, edge); ok {
					key := warning.Source + "|" + revisions[fromIndex].ID + "|" + revisions[toIndex].ID + "|" + edge.Symbol
					if _, seenWarning := seen[key]; !seenWarning {
						seen[key] = struct{}{}
						warnings = append(warnings, warning)
					}
				}
			}
			continue
		}
		key := revisions[fromIndex].ID + "|" + revisions[toIndex].ID + "|" + edge.Symbol
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		warnings = append(warnings, FeasibilityWarning{
			RevisionID: revisions[fromIndex].ID,
			Severity:   "info",
			Source:     "structural_dependency",
			DependsOn:  revisions[toIndex].ID,
			Symbol:     edge.Symbol,
			FromFile:   edge.FromFile,
			ToFile:     edge.ToFile,
			Message: fmt.Sprintf(
				"%s references %s from %s, but %s is proposed after it",
				edge.FromFile,
				edge.Symbol,
				edge.ToFile,
				revisions[toIndex].ID,
			),
		})
	}
	return warnings
}

func inferredDependencyWarning(revisions []RevisionProposal, fromIndex, toIndex int, edge structural.DependencyEdge) (FeasibilityWarning, bool) {
	source := revisions[fromIndex]
	target := revisions[toIndex]
	if revisionDependsOn(source, target.ID) {
		return FeasibilityWarning{}, false
	}
	return FeasibilityWarning{
		RevisionID: source.ID,
		Severity:   "info",
		Source:     "inferred_dependency",
		DependsOn:  target.ID,
		Symbol:     edge.Symbol,
		FromFile:   edge.FromFile,
		ToFile:     edge.ToFile,
		Message: fmt.Sprintf(
			"%s structurally depends on %s via %s; add depends_on: [%q] to make the stack dependency explicit",
			source.ID,
			target.ID,
			edge.Symbol,
			target.ID,
		),
	}, true
}

func revisionDependsOn(revision RevisionProposal, dependencyID string) bool {
	dependencyID = strings.TrimSpace(dependencyID)
	if dependencyID == "" {
		return false
	}
	for _, current := range revision.DependsOn {
		if strings.TrimSpace(current) == dependencyID {
			return true
		}
	}
	return false
}

type structuralRevisionLocator struct {
	revisions          []RevisionProposal
	fileOwners         map[string]int
	sourceSymbolOwners map[string]int
	targetSymbolOwners map[string]int
}

func structuralRevisionIndex(revisions []RevisionProposal) structuralRevisionLocator {
	index := structuralRevisionLocator{
		revisions:          revisions,
		fileOwners:         map[string]int{},
		sourceSymbolOwners: map[string]int{},
		targetSymbolOwners: map[string]int{},
	}
	for revisionIndex, revision := range revisions {
		for _, file := range revision.Files {
			file = strings.TrimSpace(file)
			if file == "" {
				continue
			}
			if _, ok := index.fileOwners[file]; !ok {
				index.fileOwners[file] = revisionIndex
			}
		}
		for _, hunk := range revision.Hunks {
			file := strings.TrimSpace(hunk.File)
			symbol := strings.TrimSpace(hunk.Symbol)
			if file == "" || symbol == "" {
				continue
			}
			key := structuralSymbolKey(file, symbol)
			if _, ok := index.targetSymbolOwners[key]; !ok {
				index.targetSymbolOwners[key] = revisionIndex
			}
		}
	}
	for revisionIndex, revision := range revisions {
		for _, hunk := range revision.Hunks {
			file := strings.TrimSpace(hunk.File)
			if file == "" {
				continue
			}
			for _, other := range revisions {
				for _, targetHunk := range other.Hunks {
					symbol := strings.TrimSpace(targetHunk.Symbol)
					if symbol == "" || !patchMentionsSymbol(hunk.Patch, symbol) {
						continue
					}
					key := structuralSymbolKey(file, symbol)
					if _, ok := index.sourceSymbolOwners[key]; !ok {
						index.sourceSymbolOwners[key] = revisionIndex
					}
				}
			}
		}
	}
	return index
}

func (index structuralRevisionLocator) sourceRevision(file, symbol string) (int, bool) {
	if revisionIndex, ok := index.sourceSymbolOwners[structuralSymbolKey(file, symbol)]; ok {
		return revisionIndex, true
	}
	return index.fileRevision(file)
}

func (index structuralRevisionLocator) targetRevision(file, symbol string) (int, bool) {
	if revisionIndex, ok := index.targetSymbolOwners[structuralSymbolKey(file, symbol)]; ok {
		return revisionIndex, true
	}
	return index.fileRevision(file)
}

func (index structuralRevisionLocator) fileRevision(file string) (int, bool) {
	revisionIndex, ok := index.fileOwners[strings.TrimSpace(file)]
	return revisionIndex, ok
}

func structuralSymbolKey(file, symbol string) string {
	return strings.TrimSpace(file) + "\x00" + strings.TrimSpace(symbol)
}

func patchMentionsSymbol(patch, symbol string) bool {
	symbol = regexp.QuoteMeta(strings.TrimSpace(symbol))
	if symbol == "" || patch == "" {
		return false
	}
	pattern := regexp.MustCompile(`(^|[^A-Za-z0-9_$])` + symbol + `([^A-Za-z0-9_$]|$)`)
	return pattern.FindStringIndex(patch) != nil
}

func testSeparationWarnings(revisions []RevisionProposal) []FeasibilityWarning {
	fileToRevision := map[string]string{}
	for _, revision := range revisions {
		for _, file := range revision.Files {
			fileToRevision[file] = revision.ID
		}
	}
	var warnings []FeasibilityWarning
	for file, revisionID := range fileToRevision {
		counterpart := testCounterpartSource(file)
		if counterpart == "" {
			continue
		}
		sourceRevisionID, ok := fileToRevision[counterpart]
		if !ok || sourceRevisionID == revisionID {
			continue
		}
		warnings = append(warnings, FeasibilityWarning{
			RevisionID: revisionID,
			Severity:   "info",
			Source:     "test_pairing",
			Message:    fmt.Sprintf("%s is separated from likely source counterpart %s", file, counterpart),
		})
	}
	return warnings
}

func testCounterpartSource(file string) string {
	dir := filepath.Dir(file)
	base := filepath.Base(file)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	switch {
	case strings.HasSuffix(name, "_test"):
		return filepath.ToSlash(filepath.Join(dir, strings.TrimSuffix(name, "_test")+ext))
	case strings.HasSuffix(name, ".test"):
		return filepath.ToSlash(filepath.Join(dir, strings.TrimSuffix(name, ".test")+ext))
	case strings.HasSuffix(name, ".spec"):
		return filepath.ToSlash(filepath.Join(dir, strings.TrimSuffix(name, ".spec")+ext))
	default:
		return ""
	}
}

func findRevisionHunkByID(hunks []HunkRange, id string) HunkRange {
	for _, hunk := range hunks {
		if hunk.ID == id {
			return hunk
		}
	}
	return HunkRange{}
}

func cleanFiles(files []string) []string {
	out := make([]string, 0, len(files))
	seen := map[string]struct{}{}
	for _, file := range files {
		file = normalizeFilesetPath(file)
		if file == "" {
			continue
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		out = append(out, file)
	}
	sort.Strings(out)
	return out
}

func filterFilesByFilesets(files []string, filesets []string) ([]string, error) {
	filters := cleanFilesets(filesets)
	if len(filters) == 0 {
		return files, nil
	}
	out := make([]string, 0, len(files))
	for _, file := range files {
		for _, filter := range filters {
			if fileMatchesFileset(file, filter) {
				out = append(out, file)
				break
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no working-copy changes matched demux filesets: %s", strings.Join(filters, ", "))
	}
	return out, nil
}

func excludeFilesByFilesets(files []string, filesets []string) ([]string, error) {
	filters := cleanFilesets(filesets)
	if len(filters) == 0 {
		return files, nil
	}
	out := make([]string, 0, len(files))
	for _, file := range files {
		excluded := false
		for _, filter := range filters {
			if fileMatchesFileset(file, filter) {
				excluded = true
				break
			}
		}
		if !excluded {
			out = append(out, file)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("all working-copy changes were excluded from demux by filesets: %s", strings.Join(filters, ", "))
	}
	return out, nil
}

func cleanFilesets(filesets []string) []string {
	out := make([]string, 0, len(filesets))
	seen := map[string]struct{}{}
	for _, fileset := range filesets {
		fileset = normalizeFilesetPath(fileset)
		if fileset == "" || fileset == "." {
			continue
		}
		if _, ok := seen[fileset]; ok {
			continue
		}
		seen[fileset] = struct{}{}
		out = append(out, fileset)
	}
	sort.Strings(out)
	return out
}

func normalizeFilesetPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	path = filepath.ToSlash(filepath.Clean(path))
	path = strings.TrimPrefix(path, "./")
	if path == "." {
		return "."
	}
	return strings.Trim(path, "/")
}

func fileMatchesFileset(file, fileset string) bool {
	file = normalizeFilesetPath(file)
	fileset = normalizeFilesetPath(fileset)
	if file == "" || fileset == "" || fileset == "." {
		return false
	}
	return file == fileset || strings.HasPrefix(file, fileset+"/")
}

func filterHunksByFiles(hunks []HunkRange, files []string) []HunkRange {
	if len(hunks) == 0 {
		return nil
	}
	fileSet := stringSet(files)
	out := make([]HunkRange, 0, len(hunks))
	for _, hunk := range hunks {
		if _, ok := fileSet[normalizeFilesetPath(hunk.File)]; ok {
			out = append(out, hunk)
		}
	}
	return out
}

func groupFiles(files []string) [][]string {
	remaining := append([]string(nil), files...)
	sort.SliceStable(remaining, func(i, j int) bool {
		return fileOrderRank(remaining[i]) < fileOrderRank(remaining[j])
	})
	if looksLikeAppScaffold(remaining) {
		return groupAppScaffoldFiles(remaining)
	}

	groups := [][]string{}
	used := map[string]struct{}{}
	for _, file := range remaining {
		if _, ok := used[file]; ok {
			continue
		}
		group := []string{file}
		used[file] = struct{}{}
		for _, candidate := range remaining {
			if _, ok := used[candidate]; ok {
				continue
			}
			if counterpartKey(candidate) != "" && counterpartKey(candidate) == counterpartKey(file) {
				group = append(group, candidate)
				used[candidate] = struct{}{}
			}
		}
		sort.Strings(group)
		groups = append(groups, group)
	}
	return groups
}

func looksLikeAppScaffold(files []string) bool {
	if len(files) < 6 {
		return false
	}
	fileSet := stringSet(files)
	if _, ok := fileSet["package.json"]; !ok {
		return false
	}
	for _, marker := range []string{
		"index.html",
		"src/main.tsx",
		"src/main.ts",
		"src/main.jsx",
		"src/main.js",
		"src/App.tsx",
		"src/App.ts",
		"src/App.jsx",
		"src/App.js",
		"src-tauri/tauri.conf.json",
	} {
		if _, ok := fileSet[marker]; ok {
			return true
		}
	}
	return false
}

func groupAppScaffoldFiles(files []string) [][]string {
	type scaffoldGroup struct {
		key   string
		rank  int
		files []string
	}
	groupsByKey := map[string]*scaffoldGroup{}
	order := []string{}
	var passthrough []string
	for _, file := range files {
		key, rank := appScaffoldGroupKey(file)
		if key == "" {
			passthrough = append(passthrough, file)
			continue
		}
		group := groupsByKey[key]
		if group == nil {
			group = &scaffoldGroup{key: key, rank: rank}
			groupsByKey[key] = group
			order = append(order, key)
		}
		group.files = append(group.files, file)
	}
	sort.SliceStable(order, func(i, j int) bool {
		left := groupsByKey[order[i]]
		right := groupsByKey[order[j]]
		if left.rank != right.rank {
			return left.rank < right.rank
		}
		return left.key < right.key
	})
	groups := make([][]string, 0, len(order)+len(passthrough))
	for _, key := range order {
		group := groupsByKey[key]
		files := cleanFiles(group.files)
		if len(files) > 0 {
			groups = append(groups, files)
		}
	}
	if len(passthrough) > 0 {
		groups = append(groups, groupFiles(passthrough)...)
	}
	return groups
}

func appScaffoldGroupKey(file string) (string, int) {
	file = normalizeFilesetPath(file)
	switch {
	case file == "":
		return "", 0
	case file == "README.md" || strings.HasSuffix(file, ".md"):
		return "scaffold:docs", 90
	case strings.HasPrefix(file, "src-tauri/"):
		return "scaffold:tauri", 20
	case file == "index.html" || strings.HasPrefix(file, "src/"):
		return "scaffold:frontend", 30
	case file == ".gitignore" ||
		file == ".vscode/extensions.json" ||
		file == "package.json" ||
		file == "package-lock.json" ||
		strings.HasPrefix(file, "tsconfig") ||
		strings.HasPrefix(file, "vite.config.") ||
		strings.HasSuffix(file, ".config.ts") ||
		strings.HasSuffix(file, ".config.js"):
		return "scaffold:tooling", 10
	default:
		return "", 0
	}
}

func expandGroupsByChangedSymbols(fileGroups [][]string, hunksByFile map[string][]HunkRange) []demuxGroup {
	groups := make([]demuxGroup, 0, len(fileGroups))
	for _, files := range fileGroups {
		if len(files) == 1 {
			if symbolGroups := symbolGroupsForFile(files[0], hunksByFile[files[0]]); len(symbolGroups) > 0 {
				groups = append(groups, symbolGroups...)
				continue
			}
		}
		hunks := hunksForFiles(hunksByFile, files)
		groups = append(groups, demuxGroup{
			Files:      files,
			Hunks:      hunks,
			Confidence: confidenceForFileGroup(files),
		})
	}
	return groups
}

func symbolGroupsForFile(file string, hunks []HunkRange) []demuxGroup {
	if len(hunks) < 2 {
		return nil
	}
	bySymbol := map[string][]HunkRange{}
	symbolOrder := map[string]int{}
	symbolKind := map[string]string{}
	for index, hunk := range hunks {
		if strings.TrimSpace(hunk.Symbol) == "" {
			return nil
		}
		key := hunk.SymbolKind + ":" + hunk.Symbol
		bySymbol[key] = append(bySymbol[key], hunk)
		symbolKind[key] = hunk.SymbolKind
		if _, ok := symbolOrder[key]; !ok {
			symbolOrder[key] = index
		}
	}
	if len(bySymbol) < 2 {
		return nil
	}
	keys := make([]string, 0, len(bySymbol))
	for key := range bySymbol {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return symbolOrder[keys[i]] < symbolOrder[keys[j]]
	})
	groups := make([]demuxGroup, 0, len(keys))
	for _, key := range keys {
		_, symbol, _ := strings.Cut(key, ":")
		groups = append(groups, demuxGroup{
			Files:      []string{file},
			Hunks:      bySymbol[key],
			UseHunks:   true,
			Symbol:     symbol,
			SymbolKind: symbolKind[key],
			Confidence: 0.72,
		})
	}
	return groups
}

func orderGroupsByStructuralDependencies(groups [][]string, facts structural.Facts) [][]string {
	if len(groups) < 2 || len(facts.Edges) == 0 {
		return groups
	}
	fileGroup := map[string]int{}
	for index, group := range groups {
		for _, file := range group {
			fileGroup[file] = index
		}
	}
	dependsOn := map[int]map[int]struct{}{}
	dependents := map[int]map[int]struct{}{}
	for _, edge := range facts.Edges {
		fromGroup, hasFrom := fileGroup[edge.FromFile]
		toGroup, hasTo := fileGroup[edge.ToFile]
		if !hasFrom || !hasTo || fromGroup == toGroup {
			continue
		}
		if dependsOn[fromGroup] == nil {
			dependsOn[fromGroup] = map[int]struct{}{}
		}
		if dependents[toGroup] == nil {
			dependents[toGroup] = map[int]struct{}{}
		}
		dependsOn[fromGroup][toGroup] = struct{}{}
		dependents[toGroup][fromGroup] = struct{}{}
	}
	if len(dependsOn) == 0 {
		return groups
	}

	ready := make([]int, 0, len(groups))
	for index := range groups {
		if len(dependsOn[index]) == 0 {
			ready = append(ready, index)
		}
	}
	sortGroupIndexes(ready, groups)

	var ordered [][]string
	used := map[int]struct{}{}
	for len(ready) > 0 {
		index := ready[0]
		ready = ready[1:]
		if _, ok := used[index]; ok {
			continue
		}
		used[index] = struct{}{}
		ordered = append(ordered, groups[index])
		for dependent := range dependents[index] {
			delete(dependsOn[dependent], index)
			if len(dependsOn[dependent]) == 0 {
				ready = append(ready, dependent)
			}
		}
		sortGroupIndexes(ready, groups)
	}
	if len(ordered) != len(groups) {
		return groups
	}
	return ordered
}

func sortGroupIndexes(indexes []int, groups [][]string) {
	sort.SliceStable(indexes, func(i, j int) bool {
		left := groups[indexes[i]][0]
		right := groups[indexes[j]][0]
		leftRank := fileOrderRank(left)
		rightRank := fileOrderRank(right)
		if leftRank != rightRank {
			return leftRank < rightRank
		}
		return left < right
	})
}

func structuralFacts(facts structural.Facts) []StructuralFact {
	out := make([]StructuralFact, 0, len(facts.Files))
	for _, fact := range facts.Files {
		out = append(out, StructuralFact{
			File:              fact.File,
			Language:          fact.Language,
			DefinedSymbols:    fact.DefinedSymbols,
			ReferencedSymbols: fact.ReferencedSymbols,
			Symbols:           structuralSymbols(fact.Symbols),
		})
	}
	return out
}

func structuralSymbols(symbols []structural.Symbol) []StructuralSymbol {
	out := make([]StructuralSymbol, 0, len(symbols))
	for _, symbol := range symbols {
		out = append(out, StructuralSymbol{
			Name:      symbol.Name,
			Kind:      symbol.Kind,
			StartLine: symbol.StartLine,
			EndLine:   symbol.EndLine,
		})
	}
	return out
}

func structuralDependencies(facts structural.Facts) []StructuralDependency {
	out := make([]StructuralDependency, 0, len(facts.Edges))
	for _, edge := range facts.Edges {
		out = append(out, StructuralDependency{
			FromFile: edge.FromFile,
			ToFile:   edge.ToFile,
			Symbol:   edge.Symbol,
		})
	}
	return out
}

func changedSymbolsForProposal(symbols []structural.HunkSymbol) []ChangedSymbol {
	out := make([]ChangedSymbol, 0, len(symbols))
	for _, symbol := range symbols {
		out = append(out, ChangedSymbol{
			HunkID:    symbol.HunkID,
			File:      symbol.File,
			Symbol:    symbol.Symbol,
			Kind:      symbol.Kind,
			StartLine: symbol.StartLine,
			EndLine:   symbol.EndLine,
		})
	}
	return out
}

func structuralHunkInputs(hunks []HunkRange) []structural.HunkInput {
	out := make([]structural.HunkInput, 0, len(hunks))
	for _, hunk := range hunks {
		out = append(out, structural.HunkInput{
			ID:       hunk.ID,
			File:     hunk.File,
			NewStart: hunk.NewStart,
			NewLines: hunk.NewLines,
		})
	}
	return out
}

func annotateHunksWithSymbols(hunks []HunkRange, symbols []structural.HunkSymbol) []HunkRange {
	if len(hunks) == 0 || len(symbols) == 0 {
		return hunks
	}
	byHunk := map[string]structural.HunkSymbol{}
	for _, symbol := range symbols {
		byHunk[symbol.HunkID] = symbol
	}
	out := append([]HunkRange(nil), hunks...)
	for index := range out {
		symbol, ok := byHunk[out[index].ID]
		if !ok {
			continue
		}
		out[index].Symbol = symbol.Symbol
		out[index].SymbolKind = symbol.Kind
	}
	return out
}

func fileOrderRank(file string) int {
	lower := strings.ToLower(file)
	switch {
	case strings.Contains(lower, "schema"), strings.Contains(lower, "migration"), strings.Contains(lower, "migrations/"):
		return 10
	case strings.Contains(lower, "type"), strings.Contains(lower, "model"):
		return 20
	case strings.HasSuffix(lower, "_test.go"), strings.Contains(lower, ".test."), strings.Contains(lower, ".spec."):
		return 80
	case strings.HasSuffix(lower, ".md"):
		return 90
	default:
		return 50
	}
}

func counterpartKey(file string) string {
	dir := filepath.Dir(file)
	base := filepath.Base(file)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	switch name {
	case "package-lock", "npm-shrinkwrap", "pnpm-lock", "yarn", "bun":
		name = "package"
	}
	for _, suffix := range []string{"_test", ".test", ".spec"} {
		name = strings.TrimSuffix(name, suffix)
	}
	if name == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Join(dir, name))
}

func proposalIntent(prefix string, files []string) string {
	prefix = strings.TrimSpace(prefix)
	label := files[0]
	if len(files) > 1 {
		label = counterpartKey(files[0])
		if label == "" {
			label = files[0]
		}
	}
	label = strings.Trim(strings.ReplaceAll(label, "_", " "), "/")
	if prefix == "" {
		return "update " + label
	}
	return prefix + ": " + label
}

func proposalIntentForGroup(prefix string, group demuxGroup) string {
	if strings.TrimSpace(group.Symbol) == "" {
		return proposalIntent(prefix, group.Files)
	}
	label := group.Symbol
	if strings.TrimSpace(group.SymbolKind) != "" {
		label = group.SymbolKind + " " + group.Symbol
	}
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return "update " + label
	}
	return prefix + ": " + label
}

func confidenceForFileGroup(files []string) float64 {
	if len(files) == 1 {
		return 0.55
	}
	return 0.65
}

func parseGitHunks(diff string) []HunkRange {
	type filePatch struct {
		file   string
		header []string
	}
	var current filePatch
	var hunkLines []string
	var hunkHeader string
	var hunks []HunkRange
	hunkIndex := 1

	flush := func() {
		if current.file == "" || hunkHeader == "" || len(hunkLines) == 0 {
			return
		}
		oldStart, oldLines, newStart, newLines := parseHunkHeader(hunkHeader)
		patchLines := append([]string{}, current.header...)
		patchLines = append(patchLines, hunkLines...)
		hunks = append(hunks, HunkRange{
			ID:       fmt.Sprintf("h%d", hunkIndex),
			File:     current.file,
			Header:   hunkHeader,
			OldStart: oldStart,
			OldLines: oldLines,
			NewStart: newStart,
			NewLines: newLines,
			Patch:    strings.Join(patchLines, ""),
		})
		hunkIndex++
		hunkLines = nil
		hunkHeader = ""
	}

	for _, line := range strings.SplitAfter(diff, "\n") {
		trimmed := strings.TrimRight(line, "\n")
		switch {
		case strings.HasPrefix(trimmed, "diff --git "):
			flush()
			current = filePatch{header: []string{line}}
		case current.header != nil && hunkHeader == "" && strings.HasPrefix(trimmed, "+++ "):
			current.header = append(current.header, line)
			current.file = parseDiffFile(trimmed)
		case current.header != nil && hunkHeader == "" && !strings.HasPrefix(trimmed, "@@ "):
			current.header = append(current.header, line)
		case strings.HasPrefix(trimmed, "@@ "):
			flush()
			hunkHeader = trimmed
			hunkLines = []string{line}
		case hunkHeader != "":
			hunkLines = append(hunkLines, line)
		}
	}
	flush()
	return hunks
}

func parseDiffFile(line string) string {
	value := strings.TrimSpace(strings.TrimPrefix(line, "+++"))
	value = strings.TrimSpace(value)
	if value == "/dev/null" {
		return ""
	}
	return strings.TrimPrefix(value, "b/")
}

func parseHunkHeader(header string) (int, int, int, int) {
	matches := hunkHeaderPattern.FindStringSubmatch(header)
	if len(matches) == 0 {
		return 0, 0, 0, 0
	}
	return atoiDefault(matches[1], 0), atoiDefault(matches[2], 1), atoiDefault(matches[3], 0), atoiDefault(matches[4], 1)
}

func atoiDefault(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func hunksByFile(hunks []HunkRange) map[string][]HunkRange {
	out := map[string][]HunkRange{}
	for _, hunk := range hunks {
		out[hunk.File] = append(out[hunk.File], hunk)
	}
	return out
}

func hunksForFiles(hunks map[string][]HunkRange, files []string) []HunkRange {
	var out []HunkRange
	for _, file := range files {
		out = append(out, hunks[file]...)
	}
	return out
}

func hunkIDs(hunks []HunkRange) []string {
	ids := make([]string, 0, len(hunks))
	for _, hunk := range hunks {
		if strings.TrimSpace(hunk.ID) == "" {
			continue
		}
		ids = append(ids, hunk.ID)
	}
	return ids
}

func filesForHunks(hunks []HunkRange) []string {
	seen := map[string]struct{}{}
	var files []string
	for _, hunk := range hunks {
		if strings.TrimSpace(hunk.File) == "" {
			continue
		}
		if _, ok := seen[hunk.File]; ok {
			continue
		}
		seen[hunk.File] = struct{}{}
		files = append(files, hunk.File)
	}
	sort.Strings(files)
	return files
}

func demuxPlanInstructions() []string {
	return demuxPlanningContract()
}

func demuxWorkflowInstructions() []string {
	return []string{
		"Review the returned proposal and review result.",
		"If state is ready_to_apply, call gx_apply_revision_plan with review.proposal unless you can improve the revision grouping.",
		"If state is repair_recommended, revise review.proposal from repair_hints or intentionally accept the hints before apply.",
		"If state is repair_required, revise review.proposal from errors, repair_hints, and warning-severity feasibility_warnings.",
		"If review returns errors or repair_hints, revise the proposal and review again.",
		"Call gx_apply_revision_plan only after review is valid and warning-severity feasibility issues are fixed or intentionally accepted.",
	}
}

func demuxPlanningContract() []string {
	return []string{
		"Return a compose proposal with the same id, repo_root, proposed_change_id, and proposed_commit_id.",
		"Create ordered revisions with one logical intent each.",
		"For hunk-level revisions, set use_hunks=true and hunk_ids to ids from top-level hunks; do not copy patch payloads.",
		"Assign every top-level hunk exactly once, either by hunk_id or by one whole-file revision for that hunk's file.",
		"Do not mix hunk-level and whole-file revisions for the same file.",
		"Prefer grouping hunks that share the same changed symbol unless they represent separate intents.",
		"Use structural_dependencies as ordering hints: to_file should usually appear before from_file.",
		"Use target_stack for an existing stack when the revision clearly belongs on that stack.",
		"To create a new stack, set target_stack and include base_stack when the base should not be the repo default.",
		"Leave target_stack empty to keep a revision on the current stack.",
		"When setting a route, include route_source and route_reason so route review can distinguish heuristic, ai, and user decisions.",
		"Keep session_ids and provenance_status from the source proposal unless the caller supplies explicit provenance.",
	}
}

func demuxRevisionShape() map[string]any {
	return map[string]any{
		"id":                "short stable id, for example r1",
		"intent":            "one-line logical intent",
		"files":             []string{"optional; GX derives this from hunk_ids when omitted"},
		"use_hunks":         true,
		"hunk_ids":          []string{"h1"},
		"depends_on":        []string{"optional earlier revision id"},
		"target_stack":      "optional existing or new stack, for example feature/name, bug/name, docs/name, test/name, or chore/name",
		"base_stack":        "optional base stack for a new routed stack",
		"route_source":      "heuristic | ai | user",
		"route_reason":      "short explanation for the route",
		"route_confidence":  0.8,
		"provenance_status": "explicit | repo_local | absent",
		"session_ids":       []string{"optional captured GX session ids"},
		"confidence":        0.8,
	}
}

func demuxToolArgumentShape(proposal DemuxProposal) map[string]any {
	return map[string]any{
		"proposal": map[string]any{
			"id":                 proposal.ID,
			"repo_root":          proposal.RepoRoot,
			"proposed_change_id": proposal.ProposedChangeID,
			"proposed_commit_id": proposal.ProposedCommitID,
			"status":             proposal.Status,
			"revisions":          []RevisionProposal{},
		},
	}
}

func writeRevisionPatch(hunks []HunkRange) (string, func(), error) {
	var patch strings.Builder
	for _, hunk := range hunks {
		patch.WriteString(hunk.Patch)
		if !strings.HasSuffix(hunk.Patch, "\n") {
			patch.WriteString("\n")
		}
	}
	file, err := os.CreateTemp("", "gx-revision-*.patch")
	if err != nil {
		return "", func() {}, fmt.Errorf("create revision patch: %w", err)
	}
	if _, err := file.WriteString(patch.String()); err != nil {
		_ = file.Close()
		_ = os.Remove(file.Name())
		return "", func() {}, fmt.Errorf("write revision patch: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(file.Name())
		return "", func() {}, fmt.Errorf("close revision patch: %w", err)
	}
	return file.Name(), func() { _ = os.Remove(file.Name()) }, nil
}

func shortID(value string, n int) string {
	if len(value) <= n {
		return value
	}
	return value[:n]
}
