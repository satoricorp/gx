package vcs

import (
	"context"
	"slices"
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
	if err := s.reconcileStoredStacksRemoteState(ctx, store, repo, repoID, stored); err != nil {
		return nil, err
	}
	reloaded, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return nil, err
	}
	stored = reloaded
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
	targets, err := s.gitBranchTargets(ctx, repo.RootPath)
	if err != nil {
		return nil, err
	}
	bookmarkTargets := make(map[string]string, len(targets))
	for _, target := range targets {
		bookmarkTargets[target.Name] = target.CommitID
	}
	mergedByStack := s.stackMergeStates(ctx, repo.RootPath, repo.defaultBaseBranch(), infos, bookmarkTargets)

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
		if !stackVisibleInGx(info) {
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

func (s *Service) stackMergeStates(ctx context.Context, repoRoot, fallbackBaseRef string, stacks []StackInfo, bookmarkTargets map[string]string) map[int64]bool {
	merged := make(map[int64]bool, len(stacks))
	for _, stack := range stacks {
		merged[stack.ID] = s.stackMergedIntoBase(ctx, repoRoot, fallbackBaseRef, stack, bookmarkTargets)
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

func (s *Service) stackMergeBaseSelectors(ctx context.Context, repoRoot, baseRef, remoteName, defaultBaseBranch string) []string {
	baseRef = strings.TrimSpace(baseRef)
	if baseRef == "" {
		return nil
	}
	selectors := make([]string, 0, 3)
	if s.stackBaseRefExists(ctx, repoRoot, baseRef, remoteName) {
		selectors = append(selectors, baseRef)
		if !strings.Contains(baseRef, "@") {
			remoteName = strings.TrimSpace(remoteName)
			if remoteName == "" {
				remoteName = "origin"
			}
			remoteBase := baseRef + "@" + remoteName
			if remoteBase != baseRef && s.revExists(ctx, repoRoot, remoteBase) {
				selectors = append(selectors, remoteBase)
			}
		}
	}
	defaultBaseBranch = strings.TrimSpace(defaultBaseBranch)
	if defaultBaseBranch != "" && defaultBaseBranch != baseRef && s.revExists(ctx, repoRoot, defaultBaseBranch) {
		if len(selectors) == 0 || !slices.Contains(selectors, defaultBaseBranch) {
			selectors = append(selectors, defaultBaseBranch)
		}
	}
	return selectors
}

func (s *Service) revExists(ctx context.Context, repoRoot, rev string) bool {
	rev = strings.TrimSpace(rev)
	if rev == "" {
		return false
	}
	_, err := s.runStdoutTrimmed(ctx, repoRoot, "git", "rev-parse", "--verify", rev+"^{commit}")
	return err == nil
}
