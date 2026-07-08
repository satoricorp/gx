package vcs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
)

// ErrSharedStackBookmark means two or more GX stacks point at the same JJ change.
var ErrSharedStackBookmark = errors.New("gx stacks must not share the same jj change")

type jjBookmarkTarget struct {
	Name     string
	ChangeID string
}

const (
	legacyGXStackBookmarkPrefix = "gx/"
	legacyGXDraftBookmarkPrefix = "gx/draft/"
)

func isGXStackBookmark(name string) bool {
	name = strings.TrimSpace(name)
	if legacyGXInternalCheckoutRef(name) {
		return false
	}
	return isConventionalStackBookmark(name) || strings.HasPrefix(name, legacyGXStackBookmarkPrefix)
}

func legacyGXInternalCheckoutRef(name string) bool {
	name = strings.TrimSpace(name)
	return name == "gx/base" || name == "gx/edit" ||
		strings.HasPrefix(name, "gx/base/") ||
		strings.HasPrefix(name, "gx/edit/")
}

func baseCheckoutRef(baseRef string) string {
	baseRef = cleanRefName(baseRef)
	if legacyStackBookmarkName(baseRef) {
		return stackBookmarkName(baseRef, "")
	}
	if baseRef == "" {
		return "main"
	}
	return baseRef
}

func gxAuthoringCheckoutRef(baseRef string) string {
	return baseCheckoutRef(baseRef)
}

func gxAuthoringBaseFromCheckoutRef(name string) (string, bool) {
	return "", false
}

func stackBookmarkName(name, headChangeID string) string {
	name = cleanRefName(name)
	name = strings.TrimPrefix(name, legacyGXDraftBookmarkPrefix)
	name = strings.TrimPrefix(name, legacyGXStackBookmarkPrefix)
	if kind, rest, ok := splitConventionalStackBookmark(name); ok {
		slug := bookmarkSlug(rest)
		if slug == "" {
			slug = fallbackStackSlug(headChangeID)
		}
		return kind + "/" + slug
	}
	slug := bookmarkSlug(name)
	if slug == "" {
		slug = fallbackStackSlug(headChangeID)
	}
	return conventionalStackKindForName(name) + "/" + slug
}

func stackNameFromBookmark(name string) string {
	name = cleanRefName(name)
	name = strings.TrimPrefix(name, legacyGXDraftBookmarkPrefix)
	name = strings.TrimPrefix(name, legacyGXStackBookmarkPrefix)
	if _, rest, ok := splitConventionalStackBookmark(name); ok {
		return rest
	}
	return name
}

func legacyStackBookmarkName(name string) bool {
	name = cleanRefName(name)
	return strings.HasPrefix(name, legacyGXDraftBookmarkPrefix) ||
		strings.HasPrefix(name, legacyGXStackBookmarkPrefix)
}

func cleanRefName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.TrimPrefix(name, "refs/heads/")
	name = strings.TrimPrefix(name, "origin/")
	return name
}

func splitConventionalStackBookmark(name string) (string, string, bool) {
	name = cleanRefName(name)
	kind, rest, ok := strings.Cut(name, "/")
	if !ok || strings.TrimSpace(rest) == "" || !isConventionalStackKind(kind) {
		return "", "", false
	}
	return kind, rest, true
}

func isConventionalStackBookmark(name string) bool {
	_, _, ok := splitConventionalStackBookmark(name)
	return ok
}

func isConventionalStackKind(kind string) bool {
	switch strings.TrimSpace(kind) {
	case "feature", "bug", "docs", "test", "chore":
		return true
	default:
		return false
	}
}

func conventionalStackKindForName(name string) string {
	slug := bookmarkSlug(name)
	switch {
	case strings.HasPrefix(slug, "bug-") ||
		strings.HasPrefix(slug, "fix-") ||
		strings.HasPrefix(slug, "repair-") ||
		strings.Contains(slug, "regression") ||
		strings.Contains(slug, "failure") ||
		strings.Contains(slug, "panic") ||
		strings.Contains(slug, "crash"):
		return "bug"
	case strings.HasPrefix(slug, "docs-") ||
		strings.HasPrefix(slug, "doc-") ||
		strings.HasPrefix(slug, "readme-") ||
		strings.Contains(slug, "documentation"):
		return "docs"
	case strings.HasPrefix(slug, "test-") ||
		strings.HasPrefix(slug, "tests-") ||
		strings.Contains(slug, "-test-") ||
		strings.Contains(slug, "-tests-"):
		return "test"
	case strings.HasPrefix(slug, "chore-") ||
		strings.HasPrefix(slug, "deps-") ||
		strings.HasPrefix(slug, "build-") ||
		strings.HasPrefix(slug, "release-"):
		return "chore"
	default:
		return "feature"
	}
}

func fallbackStackSlug(headChangeID string) string {
	if id := shortID(strings.TrimSpace(headChangeID), 8); id != "" {
		return "change-" + id
	}
	return "draft"
}

func isJJImmutableError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "immutable")
}

func stackEditRevision(stack storage.Stack) string {
	bookmark := strings.TrimSpace(stack.BookmarkName)
	if stack.HeadCommitID != nil {
		if head := strings.TrimSpace(*stack.HeadCommitID); head != "" && stack.Status != "published" {
			return head
		}
	}
	if stack.HeadChangeID != nil {
		if head := strings.TrimSpace(*stack.HeadChangeID); head != "" && stack.Status != "published" {
			return head
		}
	}
	return bookmark
}

// editStackWorkingCopy checks out a stack's JJ line of work. Published tips are
// immutable after push, so gx forks a new mutable change with `jj new`.
// workingCopyOnBookmarkLine reports whether @ descends from the bookmark (including @ = bookmark tip).
func (s *Service) workingCopyOnBookmarkLine(ctx context.Context, repoRoot, bookmarkName string) (bool, error) {
	bookmarkName = strings.TrimSpace(bookmarkName)
	if bookmarkName == "" {
		return false, nil
	}
	changeID, err := s.currentWorkChangeID(ctx, repoRoot)
	if err != nil || strings.TrimSpace(changeID) == "" {
		return false, err
	}
	revset := quoteJJRev(bookmarkName) + "..@"
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", revset, "--no-graph", "-T", "change_id ++ \"\\n\"")
	if err != nil {
		return false, err
	}
	for _, line := range splitLines(out) {
		if strings.TrimSpace(line) == changeID {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) editStackWorkingCopy(ctx context.Context, repoRoot string, stack storage.Stack) (string, error) {
	bookmark := strings.TrimSpace(stack.BookmarkName)
	if bookmark == "" {
		return "", fmt.Errorf("stack is empty")
	}
	rev := stackEditRevision(stack)
	output, err := s.runner.Run(ctx, repoRoot, "jj", "edit", rev)
	if err == nil {
		return output, nil
	}
	if !isJJImmutableError(err) {
		return "", err
	}
	return s.runner.Run(ctx, repoRoot, "jj", "new", bookmark)
}

func (s *Service) normalizeLegacyStackBookmarks(ctx context.Context, store *storage.Store, repoRoot string, repoID int64) error {
	if repoID == 0 {
		return nil
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	if len(stacks) == 0 {
		return nil
	}
	targets, err := s.jjBookmarkTargets(ctx, repoRoot)
	if err != nil {
		return err
	}
	targetByName := make(map[string]string, len(targets))
	for _, target := range targets {
		targetByName[target.Name] = target.ChangeID
	}
	for _, stack := range stacks {
		if legacyStackBookmarkName(stack.BaseRef) {
			newBaseRef := stackBookmarkName(stack.BaseRef, firstNonEmpty(derefString(stack.HeadChangeID), derefString(stack.HeadCommitID)))
			if err := store.RenameStackBaseRef(ctx, repoID, stack.BaseRef, newBaseRef, time.Now().UnixMilli()); err != nil {
				return err
			}
		}
		oldName := strings.TrimSpace(stack.BookmarkName)
		if legacyStackBookmarkName(stack.Name) && !legacyStackBookmarkName(oldName) {
			if err := store.RenameStack(ctx, repoID, oldName, stackNameFromBookmark(stack.Name), time.Now().UnixMilli()); err != nil {
				return err
			}
		}
		if !legacyStackBookmarkName(oldName) {
			continue
		}
		normalizedName := stackNameFromBookmark(oldName)
		newName := stackBookmarkName(oldName, firstNonEmpty(derefString(stack.HeadChangeID), derefString(stack.HeadCommitID)))
		if newName == "" {
			continue
		}
		oldTarget, oldExists := targetByName[oldName]
		stackTarget := oldTarget
		if stack.HeadChangeID != nil && strings.TrimSpace(*stack.HeadChangeID) != "" {
			stackTarget = strings.TrimSpace(*stack.HeadChangeID)
		}
		if newTarget, exists := targetByName[newName]; exists {
			if stackTarget == "" || newTarget != stackTarget {
				if stackTarget == "" {
					continue
				}
				newName = newName + "-" + shortID(stackTarget, 8)
				if _, exists := targetByName[newName]; exists {
					continue
				}
			} else if oldExists {
				if _, err := s.runner.Run(ctx, repoRoot, "jj", "bookmark", "forget", oldName); err != nil {
					return fmt.Errorf("forget legacy stack %q: %w", oldName, err)
				}
				delete(targetByName, oldName)
				if err := store.RenameStackBookmark(ctx, repoID, oldName, newName, normalizedName, time.Now().UnixMilli()); err != nil {
					return err
				}
				continue
			}
		}
		if oldExists {
			if _, err := s.runner.Run(ctx, repoRoot, "jj", "bookmark", "rename", oldName, newName); err != nil {
				return fmt.Errorf("rename legacy stack %q to %q: %w", oldName, newName, err)
			}
			targetByName[newName] = oldTarget
			delete(targetByName, oldName)
		}
		if err := store.RenameStackBookmark(ctx, repoID, oldName, newName, normalizedName, time.Now().UnixMilli()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) jjBookmarkTargets(ctx context.Context, repoRoot string) ([]jjBookmarkTarget, error) {
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl)
	if err != nil {
		return nil, err
	}
	targets := make([]jjBookmarkTarget, 0)
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		changeID := strings.TrimSpace(parts[1])
		if name == "" || changeID == "" {
			continue
		}
		targets = append(targets, jjBookmarkTarget{Name: name, ChangeID: changeID})
	}
	return targets, nil
}

func (s *Service) jjBookmarkExists(ctx context.Context, repoRoot, bookmarkName string) (bool, error) {
	bookmarkName = strings.TrimSpace(bookmarkName)
	if bookmarkName == "" {
		return false, nil
	}
	targets, err := s.jjBookmarkTargets(ctx, repoRoot)
	if err != nil {
		return false, err
	}
	for _, target := range targets {
		if target.Name == bookmarkName {
			return true, nil
		}
	}
	return false, nil
}

func knownStackBookmarks(stacks []storage.Stack) map[string]struct{} {
	known := make(map[string]struct{}, len(stacks))
	for _, stack := range stacks {
		if stack.BookmarkName != "" {
			known[stack.BookmarkName] = struct{}{}
		}
	}
	return known
}

func (s *Service) validateStackBookmarkExclusivity(ctx context.Context, store *storage.Store, repoRoot string, repoID int64) error {
	if repoID == 0 {
		return nil
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	known := knownStackBookmarks(stacks)
	if len(known) == 0 {
		return nil
	}
	targets, err := s.jjBookmarkTargets(ctx, repoRoot)
	if err != nil {
		return err
	}
	byChange := make(map[string][]string)
	for _, target := range targets {
		if _, ok := known[target.Name]; !ok {
			continue
		}
		byChange[target.ChangeID] = append(byChange[target.ChangeID], target.Name)
	}
	for changeID, names := range byChange {
		if len(names) <= 1 {
			continue
		}
		return fmt.Errorf("%w: change %s is tracked by bookmarks %s",
			ErrSharedStackBookmark, shortID(changeID, 12), strings.Join(names, ", "))
	}
	return nil
}

func (s *Service) validateStackBookmarks(ctx context.Context, repoRoot string) error {
	store, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByRoot(ctx, repoRoot)
	if err != nil {
		return err
	}
	if repoRow == nil {
		return nil
	}
	return s.validateStackBookmarkExclusivity(ctx, store, repoRoot, repoRow.ID)
}

func (s *Service) relocateKnownStackBookmarksFromChange(ctx context.Context, store *storage.Store, repoID int64, repoRoot, keepBookmark, targetRev string) error {
	if repoID == 0 {
		return nil
	}
	changeID, err := s.changeIDForRev(ctx, repoRoot, targetRev)
	if err != nil {
		return err
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	stackByBookmark := make(map[string]storage.Stack, len(stacks))
	known := knownStackBookmarks(stacks)
	for _, stack := range stacks {
		stackByBookmark[stack.BookmarkName] = stack
	}
	targets, err := s.jjBookmarkTargets(ctx, repoRoot)
	if err != nil {
		return err
	}
	keepBookmark = strings.TrimSpace(keepBookmark)
	for _, target := range targets {
		if target.ChangeID != changeID {
			continue
		}
		if target.Name == keepBookmark {
			continue
		}
		if _, ok := known[target.Name]; !ok {
			continue
		}
		stack := stackByBookmark[target.Name]
		distinctRev, err := s.distinctChangeForStack(ctx, repoRoot, stack, changeID, known)
		if err != nil {
			return err
		}
		if _, err := s.runJJGitBacked(ctx, repoRoot, "bookmark", "set", target.Name, "-r", distinctRev, "--allow-backwards"); err != nil {
			return fmt.Errorf("move bookmark %q off shared change: %w", target.Name, err)
		}
	}
	return nil
}

func (s *Service) assertBookmarkTargetAvailable(ctx context.Context, store *storage.Store, repoID int64, repoRoot, bookmarkName, targetRev string) error {
	if repoID == 0 {
		return nil
	}
	changeID, err := s.changeIDForRev(ctx, repoRoot, targetRev)
	if err != nil {
		return err
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	known := knownStackBookmarks(stacks)
	targets, err := s.jjBookmarkTargets(ctx, repoRoot)
	if err != nil {
		return err
	}
	bookmarkName = strings.TrimSpace(bookmarkName)
	for _, target := range targets {
		if target.ChangeID != changeID {
			continue
		}
		if target.Name == bookmarkName {
			continue
		}
		if _, ok := known[target.Name]; ok {
			return fmt.Errorf("%w: cannot set %q on change %s while %q is already there",
				ErrSharedStackBookmark, bookmarkName, shortID(changeID, 12), target.Name)
		}
	}
	return nil
}

func (s *Service) gitBranchDetached(ctx context.Context, repoRoot string) bool {
	branch, err := s.gitValue(ctx, repoRoot, "branch", "--show-current")
	return err != nil || strings.TrimSpace(branch) == ""
}

// ensureGitAttachedToStack points git HEAD at the stack's JJ bookmark when detached or mismatched.
// GX uses JJ for commits; git branch attachment keeps colocated tools consistent.
func (s *Service) ensureGitAttachedToStack(ctx context.Context, repo RepoInfo, stack StackInfo) (RepoInfo, error) {
	if stack.BookmarkName == "" {
		return repo, nil
	}
	if repo.BranchName != nil && strings.TrimSpace(*repo.BranchName) == stack.BookmarkName && !s.gitBranchDetached(ctx, repo.RootPath) {
		return repo, nil
	}
	store, err := openStore(ctx)
	if err != nil {
		return repo, err
	}
	defer store.Close()
	repoRow, err := store.FindRepoByRoot(ctx, repo.RootPath)
	if err != nil {
		return repo, err
	}
	if repoRow != nil {
		if err := s.repairSharedStackBookmarks(ctx, store, repo.RootPath, repoRow.ID, stack.BookmarkName); err != nil {
			return repo, err
		}
	}
	commitID, err := s.commitIDForRev(ctx, repo.RootPath, stack.BookmarkName)
	if err != nil {
		return repo, err
	}
	if err := s.ensureBranchMutationAllowed(ctx, repo.RootPath, stack.BookmarkName, commitID); err != nil {
		return repo, err
	}
	if err := s.attachGitBranch(ctx, repo.RootPath, stack.BookmarkName, commitID); err != nil {
		return repo, err
	}
	return s.ResolveJJRepoAtPath(ctx, repo.RootPath)
}

func (s *Service) repairSharedStackBookmarks(ctx context.Context, store *storage.Store, repoRoot string, repoID int64, keepBookmark string) error {
	if err := s.validateStackBookmarkExclusivity(ctx, store, repoRoot, repoID); err == nil {
		return nil
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	stackByBookmark := make(map[string]storage.Stack, len(stacks))
	for _, stack := range stacks {
		stackByBookmark[stack.BookmarkName] = stack
	}
	targets, err := s.jjBookmarkTargets(ctx, repoRoot)
	if err != nil {
		return err
	}
	known := knownStackBookmarks(stacks)
	byChange := make(map[string][]string)
	for _, target := range targets {
		if _, ok := known[target.Name]; !ok {
			continue
		}
		byChange[target.ChangeID] = append(byChange[target.ChangeID], target.Name)
	}
	for changeID, names := range byChange {
		if len(names) <= 1 {
			continue
		}
		for _, name := range names {
			if name == keepBookmark {
				continue
			}
			stack, ok := stackByBookmark[name]
			if !ok {
				continue
			}
			targetRev, err := s.distinctChangeForStack(ctx, repoRoot, stack, changeID, known)
			if err != nil {
				return err
			}
			if _, err := s.runJJGitBacked(ctx, repoRoot, "bookmark", "set", name, "-r", targetRev, "--allow-backwards"); err != nil {
				return fmt.Errorf("move bookmark %q off shared change: %w", name, err)
			}
		}
	}
	return s.validateStackBookmarkExclusivity(ctx, store, repoRoot, repoID)
}

func (s *Service) distinctChangeForStack(ctx context.Context, repoRoot string, stack storage.Stack, sharedChangeID string, known map[string]struct{}) (string, error) {
	if stack.HeadChangeID != nil {
		head := strings.TrimSpace(*stack.HeadChangeID)
		if head != "" && head != sharedChangeID {
			return head, nil
		}
	}
	rev := firstNonEmpty(stack.BookmarkName, sharedChangeID)
	for range 64 {
		parent, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", rev, "--no-graph", "-T", `parents.first().change_id() ++ "|" ++ parents.first().commit_id()`)
		if err != nil {
			return "", fmt.Errorf("find distinct change for stack %q: %w", stack.BookmarkName, err)
		}
		parts := strings.SplitN(strings.TrimSpace(lastNonEmptyLine(parent)), "|", 2)
		parentChange := ""
		parentCommit := ""
		if len(parts) > 0 {
			parentChange = strings.TrimSpace(parts[0])
		}
		if len(parts) > 1 {
			parentCommit = strings.TrimSpace(parts[1])
		}
		if parentChange == "" || parentChange == rev {
			return "", fmt.Errorf("no distinct change found for stack %q below %s", stack.BookmarkName, shortID(sharedChangeID, 12))
		}
		targets, err := s.jjBookmarkTargets(ctx, repoRoot)
		if err != nil {
			return "", err
		}
		shared := false
		for _, target := range targets {
			if target.ChangeID != parentChange {
				continue
			}
			if _, ok := known[target.Name]; !ok {
				continue
			}
			if target.Name != stack.BookmarkName {
				shared = true
				break
			}
		}
		if !shared {
			return firstNonEmpty(parentCommit, parentChange), nil
		}
		rev = firstNonEmpty(parentCommit, parentChange)
	}
	return "", fmt.Errorf("no distinct change found for stack %q below %s", stack.BookmarkName, shortID(sharedChangeID, 12))
}

func (s *Service) ensureForkedForNewStack(ctx context.Context, store *storage.Store, repoID int64, repoRoot string) error {
	if repoID == 0 {
		return nil
	}
	changeID, err := s.currentWorkChangeID(ctx, repoRoot)
	if err != nil || strings.TrimSpace(changeID) == "" {
		return err
	}
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		return err
	}
	known := knownStackBookmarks(stacks)
	targets, err := s.jjBookmarkTargets(ctx, repoRoot)
	if err != nil {
		return err
	}
	for _, target := range targets {
		if target.ChangeID != changeID {
			continue
		}
		if _, ok := known[target.Name]; ok {
			if _, err := s.runner.Run(ctx, repoRoot, "jj", "new", "-m", "fork for new gx stack"); err != nil {
				return fmt.Errorf("fork jj change for new gx stack: %w", err)
			}
			return nil
		}
	}
	return nil
}
