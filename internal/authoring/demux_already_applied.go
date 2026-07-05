package authoring

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/vcs"
)

type demuxRevisionApplyOutcome struct {
	Result  CheckpointResult
	Skipped bool
	Warning string
}

func (o demuxRevisionApplyOutcome) checkpoint() CheckpointResult {
	return o.Result
}

func demuxSkippedCheckpoint(revision RevisionProposal, reason string) demuxRevisionApplyOutcome {
	return demuxRevisionApplyOutcome{
		Skipped: true,
		Warning: reason,
		Result: CheckpointResult{
			Output: "skipped: " + reason,
		},
	}
}

func (e *Engine) excludeDemuxAlreadyAppliedContent(ctx context.Context, proposal DemuxProposal) (DemuxProposal, error) {
	if len(proposal.Revisions) == 0 {
		return proposal, nil
	}
	router, err := e.demuxRouter(ctx, proposal.RepoRoot)
	if err != nil {
		return proposal, err
	}
	sourceRev := strings.TrimSpace(proposal.ProposedCommitID)
	if sourceRev == "" {
		return proposal, nil
	}

	filtered := make([]RevisionProposal, 0, len(proposal.Revisions))
	dedupUnavailableWarned := map[string]struct{}{}
	for _, revision := range proposal.Revisions {
		target := routeTargetStack(revision)
		stack := router.lookup(target)
		if stack == nil || strings.TrimSpace(stack.BookmarkName) == "" {
			filtered = append(filtered, revision)
			continue
		}
		updated, skipped, warning, err := e.filterRevisionAlreadyOnStack(ctx, proposal.RepoRoot, stack.BookmarkName, sourceRev, revision)
		if err != nil {
			var unavailable demuxDedupStackUnavailableError
			if errors.As(err, &unavailable) {
				if _, ok := dedupUnavailableWarned[unavailable.StackBookmark]; !ok {
					dedupUnavailableWarned[unavailable.StackBookmark] = struct{}{}
					proposal.Warnings = append(proposal.Warnings, demuxDedupStackSkipWarning(unavailable.StackBookmark))
				}
				filtered = append(filtered, revision)
				continue
			}
			return proposal, err
		}
		if skipped {
			if warning != "" {
				proposal.Warnings = append(proposal.Warnings, warning)
			}
			continue
		}
		if warning != "" {
			proposal.Warnings = append(proposal.Warnings, warning)
		}
		filtered = append(filtered, updated)
	}
	proposal.Revisions = filtered
	proposal = pruneExcludedDemuxHunks(proposal)
	return proposal, nil
}

func pruneExcludedDemuxHunks(proposal DemuxProposal) DemuxProposal {
	if len(proposal.Hunks) == 0 {
		return proposal
	}
	assignedHunkIDs := map[string]struct{}{}
	assignedFiles := map[string]struct{}{}
	for _, revision := range proposal.Revisions {
		if revision.UseHunks {
			for _, id := range revision.HunkIDs {
				if id = strings.TrimSpace(id); id != "" {
					assignedHunkIDs[id] = struct{}{}
				}
			}
			for _, hunk := range revision.Hunks {
				if id := strings.TrimSpace(hunk.ID); id != "" {
					assignedHunkIDs[id] = struct{}{}
				}
			}
		}
		for _, file := range revision.Files {
			if file = strings.TrimSpace(file); file != "" {
				assignedFiles[file] = struct{}{}
			}
		}
	}
	remaining := make([]HunkRange, 0, len(proposal.Hunks))
	for _, hunk := range proposal.Hunks {
		id := strings.TrimSpace(hunk.ID)
		file := strings.TrimSpace(hunk.File)
		if _, ok := assignedHunkIDs[id]; ok {
			remaining = append(remaining, hunk)
			continue
		}
		if _, ok := assignedFiles[file]; ok {
			remaining = append(remaining, hunk)
			continue
		}
	}
	proposal.Hunks = remaining
	return proposal
}

func (e *Engine) filterRevisionAlreadyOnStack(ctx context.Context, repoRoot, stackBookmark, sourceRev string, revision RevisionProposal) (RevisionProposal, bool, string, error) {
	files := cleanFiles(revision.Files)
	if revision.UseHunks && len(revision.Hunks) > 0 {
		remainingHunks := make([]HunkRange, 0, len(revision.Hunks))
		remainingIDs := make([]string, 0, len(revision.HunkIDs))
		skippedHunks := 0
		for index, hunk := range revision.Hunks {
			file := strings.TrimSpace(hunk.File)
			if file == "" {
				continue
			}
			already, err := e.demuxFileContentMatchesRevision(ctx, repoRoot, stackBookmark, sourceRev, file)
			if err != nil {
				var unavailable demuxDedupStackUnavailableError
				if errors.As(err, &unavailable) {
					return revision, false, "", unavailable
				}
				return revision, false, "", err
			}
			if already {
				skippedHunks++
				continue
			}
			remainingHunks = append(remainingHunks, hunk)
			if index < len(revision.HunkIDs) {
				remainingIDs = append(remainingIDs, revision.HunkIDs[index])
			}
		}
		if len(remainingHunks) == 0 {
			return revision, true, fmt.Sprintf(
				"revision %s skipped at plan time: all %d hunk(s) already present on stack %q",
				revision.ID, skippedHunks, stackBookmark,
			), nil
		}
		if skippedHunks > 0 {
			revision.Hunks = remainingHunks
			revision.HunkIDs = remainingIDs
			revision.UseHunks = true
			revision.Files = hunkFiles(remainingHunks)
			return revision, false, fmt.Sprintf(
				"revision %s: skipped %d hunk(s) already present on stack %q",
				revision.ID, skippedHunks, stackBookmark,
			), nil
		}
		return revision, false, "", nil
	}

	remainingFiles := make([]string, 0, len(files))
	skippedFiles := 0
	for _, file := range files {
		already, err := e.demuxFileContentMatchesRevision(ctx, repoRoot, stackBookmark, sourceRev, file)
		if err != nil {
			var unavailable demuxDedupStackUnavailableError
			if errors.As(err, &unavailable) {
				return revision, false, "", unavailable
			}
			return revision, false, "", err
		}
		if already {
			skippedFiles++
			continue
		}
		remainingFiles = append(remainingFiles, file)
	}
	if len(remainingFiles) == 0 {
		return revision, true, fmt.Sprintf(
			"revision %s skipped at plan time: all file(s) already present on stack %q",
			revision.ID, stackBookmark,
		), nil
	}
	if skippedFiles > 0 {
		revision.Files = remainingFiles
		return revision, false, fmt.Sprintf(
			"revision %s: skipped %d file(s) already present on stack %q",
			revision.ID, skippedFiles, stackBookmark,
		), nil
	}
	return revision, false, "", nil
}

func (e *Engine) demuxFileContentMatchesRevision(ctx context.Context, repoRoot, leftRev, rightRev, file string) (bool, error) {
	leftRev = strings.TrimSpace(leftRev)
	rightRev = strings.TrimSpace(rightRev)
	file = strings.TrimSpace(file)
	if leftRev == "" || rightRev == "" || file == "" {
		return false, nil
	}
	diff, err := e.vcs.DiffRevisionRange(ctx, repoRoot, leftRev, rightRev, []string{file})
	if err != nil {
		if demuxDedupRevisionUnavailable(err) {
			return false, demuxDedupStackUnavailableError{StackBookmark: leftRev}
		}
		return false, err
	}
	return strings.TrimSpace(diff) == "", nil
}

type demuxDedupStackUnavailableError struct {
	StackBookmark string
}

func (e demuxDedupStackUnavailableError) Error() string {
	return fmt.Sprintf("dedup stack %q bookmark not found", e.StackBookmark)
}

func isDemuxDedupStackUnavailable(err error) bool {
	var unavailable demuxDedupStackUnavailableError
	return errors.As(err, &unavailable)
}

func demuxDedupRevisionUnavailable(err error) bool {
	return vcs.JJRevisionUnavailable(err)
}

func demuxDedupStackSkipWarning(stackBookmark string) string {
	return fmt.Sprintf(
		"dedup: stack %s bookmark not found; skipped overlap check (consider gx cleanup)",
		strings.TrimSpace(stackBookmark),
	)
}

func (e *Engine) demuxRevisionAlreadyOnTargetStack(ctx context.Context, repoRoot, stackBookmark, sourceRev string, revision RevisionProposal) (bool, error) {
	files := cleanFiles(revisionFilesForApply(revision))
	if len(files) == 0 {
		return false, nil
	}
	for _, file := range files {
		matches, err := e.demuxFileContentMatchesRevision(ctx, repoRoot, stackBookmark, sourceRev, file)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

func revisionFilesForApply(revision RevisionProposal) []string {
	if revision.UseHunks && len(revision.Hunks) > 0 {
		return hunkFiles(revision.Hunks)
	}
	return revision.Files
}

func hunkFiles(hunks []HunkRange) []string {
	files := make([]string, 0, len(hunks))
	seen := map[string]struct{}{}
	for _, hunk := range hunks {
		file := strings.TrimSpace(hunk.File)
		if file == "" {
			continue
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		files = append(files, file)
	}
	return cleanFiles(files)
}

func (e *Engine) validateDemuxCheckpointResult(ctx context.Context, repoRoot string, revision RevisionProposal, result CheckpointResult, stackBookmark string, sourceRev string) error {
	changeRef := firstNonEmpty(result.Change.CommitID, result.Change.ChangeID)
	if changeRef == "" {
		return fmt.Errorf("revision %s produced no recorded change", revision.ID)
	}
	empty, err := e.vcs.RevisionEmpty(ctx, repoRoot, changeRef)
	if err != nil {
		return err
	}
	if !empty {
		return nil
	}
	if stackBookmark != "" && sourceRev != "" {
		already, alreadyErr := e.demuxRevisionAlreadyOnTargetStack(ctx, repoRoot, stackBookmark, sourceRev, revision)
		if alreadyErr == nil && already {
			return errDemuxRevisionAlreadyApplied
		}
	}
	return fmt.Errorf("revision %s produced an empty commit", revision.ID)
}

var errDemuxRevisionAlreadyApplied = fmt.Errorf("demux revision already applied")

func isDemuxRevisionAlreadyApplied(err error) bool {
	return err == errDemuxRevisionAlreadyApplied
}

func (e *Engine) restoreDemuxRevisionAttempt(ctx context.Context, repoRoot, opID string) error {
	opID = strings.TrimSpace(opID)
	if opID == "" {
		return nil
	}
	return e.vcs.RestoreOperation(ctx, repoRoot, opID)
}

func demuxStackBookmarkForRevision(ctx context.Context, e *Engine, repoRoot string, revision RevisionProposal) (string, error) {
	target := routeTargetStack(revision)
	if target == "" {
		status, err := e.Status(ctx)
		if err != nil {
			return "", err
		}
		if status.Stack != nil {
			return status.Stack.BookmarkName, nil
		}
		return "", nil
	}
	router, err := e.demuxRouter(ctx, repoRoot)
	if err != nil {
		return "", err
	}
	stack := router.lookup(target)
	if stack == nil {
		return "", nil
	}
	return stack.BookmarkName, nil
}

func appendDemuxApplySkipWarning(proposal *DemuxProposal, warning string) {
	warning = strings.TrimSpace(warning)
	if warning == "" {
		return
	}
	proposal.Warnings = append(proposal.Warnings, warning)
}

func demuxOutcomeFromCheckpoint(result CheckpointResult) demuxRevisionApplyOutcome {
	return demuxRevisionApplyOutcome{Result: result}
}
