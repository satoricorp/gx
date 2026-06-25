package vcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
)

type stackReadModel struct {
	repo                  RepoInfo
	currentStack          StackInfo
	currentStackFound     bool
	stacks                []StackInfo
	currentRevisions      []RevisionSummary
	currentPublishedCount int
}

func (s *Service) loadStackReadModel(ctx context.Context, repo RepoInfo) (stackReadModel, error) {
	model := stackReadModel{repo: repo}
	store, err := openStore(ctx)
	if err != nil {
		return stackReadModel{}, err
	}
	defer store.Close()

	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return stackReadModel{}, err
	}
	if repoRow == nil {
		return model, nil
	}
	if err := s.normalizeLegacyStackBookmarks(ctx, store, repo.RootPath, repoRow.ID); err != nil {
		return stackReadModel{}, err
	}
	if err := s.repairMissingStackRowsFromBookmarks(ctx, store, repo, repoRow.ID); err != nil {
		return stackReadModel{}, err
	}

	current, found, err := s.resolveCurrentStackForReadWithStore(ctx, store, repo, repoRow.ID)
	if err != nil {
		return stackReadModel{}, err
	}
	model.currentStack = current
	model.currentStackFound = found

	stackContainer := repo.defaultBaseBranch()
	if found && strings.TrimSpace(current.BaseRef) != "" {
		stackContainer = current.BaseRef
	}
	entries, err := s.jjStackEntries(ctx, repo.RootPath, stackContainer)
	if err != nil {
		return stackReadModel{}, err
	}
	changes, err := store.ListChangesByRepoID(ctx, repoRow.ID)
	if err != nil {
		return stackReadModel{}, err
	}
	changeByJJ := make(map[string]storage.Change, len(changes))
	for _, change := range changes {
		changeByJJ[change.JJChangeID] = change
	}
	activeChangeID, err := s.currentWorkChangeID(ctx, repo.RootPath)
	if err != nil {
		activeChangeID = ""
	}
	publishedThrough := int64(-1)
	if latestPush, err := store.LatestPushByRepoID(ctx, repoRow.ID); err != nil {
		return stackReadModel{}, err
	} else if latestPush != nil && latestPush.CurrentChangeID != nil {
		publishedThrough = *latestPush.CurrentChangeID
	}
	model.currentRevisions, model.currentPublishedCount = stackRevisionsFromEntries(entries, activeChangeID, publishedThrough, changeByJJ)

	stacks, err := s.loadStoredStackInfos(ctx, store, repo, repoRow.ID)
	if err != nil {
		return stackReadModel{}, err
	}
	model.stacks = stacks
	if model.currentStackFound {
		if hydrated, ok := matchingHydratedStack(model.currentStack, stacks); ok {
			model.currentStack = hydrated
		}
	}
	return model, nil
}

func (s *Service) loadStoredStackInfos(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64) ([]StackInfo, error) {
	stored, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}
	return s.hydrateStoredStacks(ctx, store, repo, repoID, stored)
}

func (s *Service) hydrateStoredStacks(ctx context.Context, store *storage.Store, repo RepoInfo, repoID int64, stored []storage.Stack) ([]StackInfo, error) {
	if len(stored) == 0 {
		return []StackInfo{}, nil
	}
	stackIDs := make([]int64, 0, len(stored))
	publishRefs := make([]string, 0, len(stored))
	infos := make([]StackInfo, 0, len(stored))
	for index, stack := range stored {
		info := stackInfoFromStorage(stack)
		info.Alias = stackAlias(index)
		infos = append(infos, info)
		stackIDs = append(stackIDs, stack.ID)
		publishRefs = append(publishRefs, publishRefForStack(info))
	}
	changesByStack, err := store.ListChangesByStackIDs(ctx, stackIDs)
	if err != nil {
		return nil, err
	}
	bookmarksByName, err := store.ListChangeBookmarksByNames(ctx, publishRefs)
	if err != nil {
		return nil, err
	}
	latestPushByBranch, err := store.LatestPushesByBranchNames(ctx, repoID, publishRefs)
	if err != nil {
		return nil, err
	}
	targets, err := s.jjBookmarkTargets(ctx, repo.RootPath)
	if err != nil {
		return nil, err
	}
	bookmarkTargets := make(map[string]string, len(targets))
	for _, target := range targets {
		bookmarkTargets[target.Name] = target.ChangeID
	}
	mergedByStack := s.stackMergeStates(ctx, repo.RootPath, infos, bookmarkTargets)

	out := make([]StackInfo, 0, len(infos))
	for _, info := range infos {
		publishRef := publishRefForStack(info)
		revisions := s.storedRevisionsWithDiff(ctx, repo.RootPath, changesByStack[info.ID], stackPublicationState(bookmarksByName[publishRef], latestPushByBranch[publishRef]))
		info.Revisions = revisions
		info.RevisionCount = len(revisions)
		info.PublishedCount = 0
		for _, revision := range revisions {
			if revision.Published {
				info.PublishedCount++
			}
		}
		info.Status = DeriveStackStatus(info.Status, info.RevisionCount, info.PublishedCount, mergedByStack[info.ID])
		if info.Status == "merged" {
			if err := store.PrunePublishedStack(ctx, repoID, info.ID, publishRef, time.Now().UnixMilli()); err != nil {
				return nil, err
			}
		}
		if !stackVisibleInGX(info) {
			continue
		}
		out = append(out, info)
	}
	return out, nil
}

func stackPublicationState(bookmarks []storage.ChangeBookmark, latestPush storage.Push) StackPublicationState {
	lastPushedByChangeID := make(map[int64]string, len(bookmarks))
	for _, bookmark := range bookmarks {
		lastPushedByChangeID[bookmark.ChangeID] = bookmark.LastPushedCommitID
	}
	publishedThrough := int64(-1)
	if len(lastPushedByChangeID) == 0 && latestPush.CurrentChangeID != nil {
		publishedThrough = *latestPush.CurrentChangeID
	}
	return NewStackPublicationState(publishedThrough, lastPushedByChangeID)
}

func (s *Service) storedRevisionsWithDiff(ctx context.Context, repoRoot string, changes []storage.Change, publishedState StackPublicationState) []RevisionSummary {
	revisions := make([]RevisionSummary, 0, len(changes))
	for _, change := range changes {
		if IsPlaceholderDescription(change.Description) || change.Status == "abandoned" {
			continue
		}
		if !s.storedChangeHasDiff(ctx, repoRoot, change) {
			continue
		}
		revisions = append(revisions, RevisionSummary{
			Index:       len(revisions) + 1,
			ChangeID:    change.JJChangeID,
			CommitID:    change.CurrentCommitID,
			Description: change.Description,
			Status:      change.Status,
			Published:   publishedState.RevisionPublished(change),
		})
	}
	return revisions
}

func (s *Service) storedChangeHasDiff(ctx context.Context, repoRoot string, change storage.Change) bool {
	target := strings.TrimSpace(change.CurrentCommitID)
	if target == "" {
		target = strings.TrimSpace(change.JJChangeID)
	}
	if strings.TrimSpace(repoRoot) == "" || target == "" {
		return target != ""
	}
	empty, err := s.isRevisionEmpty(ctx, repoRoot, target)
	if err != nil {
		return true
	}
	return !empty
}

func matchingHydratedStack(current StackInfo, stacks []StackInfo) (StackInfo, bool) {
	for _, stack := range stacks {
		if current.ID != 0 && stack.ID == current.ID {
			return stack, true
		}
		if strings.TrimSpace(current.BookmarkName) != "" && stack.BookmarkName == current.BookmarkName {
			return stack, true
		}
	}
	return StackInfo{}, false
}

type stackMergeCandidate struct {
	stackID    int64
	selector   string
	changeID   string
	commitID   string
	baseRef    string
	remoteName string
}

func (s *Service) stackMergeStates(ctx context.Context, repoRoot string, stacks []StackInfo, bookmarkTargets map[string]string) map[int64]bool {
	merged := make(map[int64]bool, len(stacks))
	candidates := stackMergeCandidates(stacks, bookmarkTargets)
	if len(candidates) == 0 {
		return merged
	}
	revsets := make([]string, 0, len(candidates))
	baseSelectors := map[string][]string{}
	for _, candidate := range candidates {
		cacheKey := candidate.baseRef + "\x00" + candidate.remoteName
		selectors, ok := baseSelectors[cacheKey]
		if !ok {
			selectors = s.stackMergeBaseSelectors(ctx, repoRoot, candidate.baseRef, candidate.remoteName)
			baseSelectors[cacheKey] = selectors
		}
		for _, baseSelector := range selectors {
			revsets = append(revsets, fmt.Sprintf("(%s) & ancestors(%s)", quoteJJRev(candidate.selector), quoteJJRev(baseSelector)))
		}
	}
	if len(revsets) == 0 {
		return merged
	}
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", strings.Join(revsets, " | "), "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "\n"`)
	if err != nil {
		for _, stack := range stacks {
			merged[stack.ID] = s.stackMergedIntoBase(ctx, repoRoot, stack, bookmarkTargets)
		}
		return merged
	}
	seenChanges := map[string]struct{}{}
	seenCommits := map[string]struct{}{}
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) > 0 && strings.TrimSpace(parts[0]) != "" {
			seenChanges[strings.TrimSpace(parts[0])] = struct{}{}
		}
		if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
			seenCommits[strings.TrimSpace(parts[1])] = struct{}{}
		}
	}
	for _, candidate := range candidates {
		if candidate.changeID != "" {
			if _, ok := seenChanges[candidate.changeID]; ok {
				merged[candidate.stackID] = true
			}
		}
		if candidate.commitID != "" {
			if _, ok := seenCommits[candidate.commitID]; ok {
				merged[candidate.stackID] = true
			}
		}
	}
	return merged
}

func stackMergeCandidates(stacks []StackInfo, bookmarkTargets map[string]string) []stackMergeCandidate {
	candidates := make([]stackMergeCandidate, 0, len(stacks))
	for _, stack := range stacks {
		baseRef := strings.TrimSpace(stack.BaseRef)
		if baseRef == "" || IsTerminalStackStatus(stack.Status) {
			continue
		}
		remoteName := ""
		if stack.RemoteName != nil {
			remoteName = strings.TrimSpace(*stack.RemoteName)
		}
		if bookmark := strings.TrimSpace(stack.BookmarkName); bookmark != "" {
			if changeID, exists := bookmarkTargets[bookmark]; exists {
				candidates = append(candidates, stackMergeCandidate{
					stackID:    stack.ID,
					selector:   bookmark,
					changeID:   strings.TrimSpace(changeID),
					baseRef:    baseRef,
					remoteName: remoteName,
				})
				continue
			}
		}
		if stack.HeadCommitID != nil {
			if head := strings.TrimSpace(*stack.HeadCommitID); head != "" {
				candidates = append(candidates, stackMergeCandidate{
					stackID:    stack.ID,
					selector:   head,
					commitID:   head,
					baseRef:    baseRef,
					remoteName: remoteName,
				})
			}
		}
		if stack.HeadChangeID != nil {
			if head := strings.TrimSpace(*stack.HeadChangeID); head != "" {
				candidates = append(candidates, stackMergeCandidate{
					stackID:    stack.ID,
					selector:   head,
					changeID:   head,
					baseRef:    baseRef,
					remoteName: remoteName,
				})
			}
		}
	}
	return candidates
}

func (s *Service) stackMergeBaseSelectors(ctx context.Context, repoRoot, baseRef, remoteName string) []string {
	baseRef = strings.TrimSpace(baseRef)
	if baseRef == "" {
		return nil
	}
	selectors := []string{baseRef}
	if strings.Contains(baseRef, "@") {
		return selectors
	}
	remoteName = strings.TrimSpace(remoteName)
	if remoteName == "" {
		remoteName = "origin"
	}
	remoteBase := baseRef + "@" + remoteName
	if s.revExists(ctx, repoRoot, remoteBase) {
		selectors = append(selectors, remoteBase)
	}
	return selectors
}

func (s *Service) revExists(ctx context.Context, repoRoot, rev string) bool {
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return false
	}
	_, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", rev, "-n", "1", "--no-graph", "-T", "change_id")
	return err == nil
}
