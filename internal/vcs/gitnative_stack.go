package vcs

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

func (s *Service) loadGitStackReadModel(ctx context.Context, repo RepoInfo) (stackReadModel, error) {
	model := stackReadModel{repo: repo}
	store, err := openStore(ctx)
	if err != nil {
		return stackReadModel{}, err
	}
	defer store.Close()

	repoRow, err := store.FindRepoByIdentity(ctx, repo.GitCommonDir, repo.RootPath)
	if err != nil {
		return stackReadModel{}, err
	}
	if repoRow == nil {
		return model, nil
	}
	var current StackInfo
	found := false
	if repo.BranchName != nil {
		branch := strings.TrimSpace(*repo.BranchName)
		if stored, findErr := store.FindStackByBookmark(ctx, repoRow.ID, branch); findErr != nil {
			return stackReadModel{}, findErr
		} else if stored != nil {
			current = stackInfoFromStorage(*stored)
			found = true
		}
	}
	model.currentStack = current
	model.currentStackFound = found

	entries, err := s.gitStackEntries(ctx, repo)
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
	publishedThrough := int64(-1)
	if latestPush, err := store.LatestPushByRepoID(ctx, repoRow.ID); err != nil {
		return stackReadModel{}, err
	} else if latestPush != nil && latestPush.CurrentChangeID != nil {
		publishedThrough = *latestPush.CurrentChangeID
	}
	activeRevisionID := ""
	if len(entries) > 0 {
		activeRevisionID = entries[len(entries)-1].ChangeID
	}
	model.currentRevisions, model.currentPublishedCount = stackRevisionsFromEntries(entries, activeRevisionID, publishedThrough, changeByJJ)

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

func (s *Service) gitStackEntries(ctx context.Context, repo RepoInfo) ([]stackEntry, error) {
	branch, err := s.runTrimmed(ctx, repo.RootPath, "git", "branch", "--show-current")
	if err != nil {
		return nil, err
	}
	branch = cleanRefName(strings.TrimSpace(branch))
	if branch == "" {
		return nil, fmt.Errorf("%w; check out a branch before running gx", ErrDetachedHEAD)
	}
	baseRef := s.publicStackBaseRef(ctx, repo, s.defaultStackBaseRef(repo))
	rangeSpec := fmt.Sprintf("%s..HEAD", baseCheckoutRef(baseRef))
	out, err := s.runTrimmed(ctx, repo.RootPath, "git", "log", "--reverse", "--format=%H%x00%B", rangeSpec)
	if err != nil {
		if isUnbornHEADError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("git log %s: %w", rangeSpec, err)
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}
	records := strings.Split(out, "\n")
	var entries []stackEntry
	index := 0
	for _, record := range records {
		parts := strings.SplitN(record, "\x00", 2)
		if len(parts) != 2 {
			continue
		}
		commitID := strings.TrimSpace(parts[0])
		message := parts[1]
		revisionIDs := ParseRevisionIDsFromMessage(message)
		if len(revisionIDs) == 0 {
			continue
		}
		index++
		subject := strings.TrimSpace(strings.Split(message, "\n")[0])
		entries = append(entries, stackEntry{
			ChangeID:    revisionIDs[len(revisionIDs)-1],
			CommitID:    commitID,
			Description: subject,
		})
	}
	return entries, nil
}

func withRepoIdentityLock(commonDir string, fn func() error) error {
	commonDir = strings.TrimSpace(commonDir)
	if commonDir == "" {
		return fn()
	}
	lockPath := filepath.Join(commonDir, "gx", "repo.lock")
	return withLockFile(lockPath, fn)
}
