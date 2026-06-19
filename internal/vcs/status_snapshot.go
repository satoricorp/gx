package vcs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

// RevisionSnapshot is one ordered revision in a bookmark stack.
type RevisionSnapshot struct {
	Index       int
	ChangeID    string
	CommitID    string
	Description string
	ShortID     string
	Active      bool
	Published   bool
	Working     bool
	SyncNote    string
}

// BookmarkSnapshot summarizes one gx bookmark/stack.
type BookmarkSnapshot struct {
	Stack         StackInfo
	Current       bool
	ChangeCount   int
	ApprovedCount int
	FileCount     int
	Units         []UnitSummary
	Revisions     []RevisionSnapshot
}

// StatusSnapshot is the structured input for gx status/switch rendering.
type StatusSnapshot struct {
	RepoLabel string
	RepoRoot  string
	BaseRef   string
	Bookmarks []BookmarkSnapshot
}

func (s *Service) StatusSnapshot(ctx context.Context) (StatusSnapshot, error) {
	stack, err := s.Stack(ctx)
	if err != nil {
		return StatusSnapshot{}, err
	}

	store, err := openStore(ctx)
	if err != nil {
		return StatusSnapshot{}, err
	}
	defer store.Close()

	repoRow, err := store.FindRepoByRoot(ctx, stack.Repo.RootPath)
	if err != nil {
		return StatusSnapshot{}, err
	}

	publishedThrough := int64(-1)
	if repoRow != nil {
		latestPush, pushErr := store.LatestPushByRepoID(ctx, repoRow.ID)
		if pushErr != nil {
			return StatusSnapshot{}, pushErr
		}
		if latestPush != nil && latestPush.CurrentChangeID != nil {
			publishedThrough = *latestPush.CurrentChangeID
		}
	}

	activeFiles := 0
	if stack.Repo.RootPath != "" {
		if current, currentErr := s.CurrentChange(ctx, stack.Repo.RootPath, "@"); currentErr == nil {
			activeFiles = len(current.Files)
		}
	}

	baseRef := stack.Repo.defaultBaseBranch()
	if stack.Stack != nil && strings.TrimSpace(stack.Stack.BaseRef) != "" {
		baseRef = stack.Stack.BaseRef
	}

	snapshot := StatusSnapshot{
		RepoLabel: repoLabel(stack.Repo),
		RepoRoot:  stack.Repo.RootPath,
		BaseRef:   baseRef,
	}

	if len(stack.Stacks) == 0 {
		if stack.Stack != nil {
			snapshot.Bookmarks = []BookmarkSnapshot{
				buildBookmarkSnapshot(*stack.Stack, true, stack.Units, activeFiles, publishedThrough, nil),
			}
		}
		return snapshot, nil
	}

	bookmarks := make([]BookmarkSnapshot, 0, len(stack.Stacks))
	for index, entry := range stack.Stacks {
		entry.Alias = stackAlias(index)
		current := stack.Stack != nil && entry.BookmarkName == stack.Stack.BookmarkName
		var units []UnitSummary
		fileCount := 0
		if current {
			units = stack.Units
			fileCount = activeFiles
		}
		var storedChanges []storage.Change
		if repoRow != nil && entry.ID != 0 {
			rows, listErr := store.ListChangesByStackID(ctx, entry.ID)
			if listErr != nil {
				return StatusSnapshot{}, listErr
			}
			storedChanges = rows
		}
		bookmarks = append(bookmarks, buildBookmarkSnapshot(entry, current, units, fileCount, publishedThrough, storedChanges))
	}
	snapshot.Bookmarks = bookmarks
	return snapshot, nil
}

func buildBookmarkSnapshot(body StackInfo, current bool, units []UnitSummary, fileCount int, publishedThrough int64, stored []storage.Change) BookmarkSnapshot {
	changeCount := len(units)
	approved := 0
	for _, unit := range units {
		if unit.Published {
			approved++
		}
	}
	if changeCount == 0 && len(stored) > 0 {
		changeCount = len(stored)
		approved = CountPublishedChanges(stored, publishedThrough)
	}

	revisions := revisionSnapshots(body, units)
	if changeCount == 0 {
		changeCount = len(revisions)
	}

	return BookmarkSnapshot{
		Stack:         body,
		Current:       current,
		ChangeCount:   changeCount,
		ApprovedCount: approved,
		FileCount:     fileCount,
		Units:         units,
		Revisions:     revisions,
	}
}

func revisionSnapshots(body StackInfo, units []UnitSummary) []RevisionSnapshot {
	revisions := make([]RevisionSnapshot, 0, len(units))
	for _, unit := range units {
		note := revisionSyncNote(body, unit.Published)
		if unit.Active {
			note = firstNonEmpty(note, "working change")
		}
		revisions = append(revisions, RevisionSnapshot{
			Index:       unit.Index,
			ChangeID:    unit.ChangeID,
			CommitID:    unit.CommitID,
			Description: unit.Description,
			ShortID:     shortID(firstNonEmpty(unit.CommitID, unit.ChangeID), 7),
			Active:      unit.Active,
			Published:   unit.Published,
			Working:     unit.Active,
			SyncNote:    note,
		})
	}
	return revisions
}

func revisionSyncNote(body StackInfo, published bool) string {
	parts := []string{"local"}
	if published || hasRemoteSync(body) {
		parts = append(parts, "origin")
	}
	if len(parts) == 1 {
		return parts[0] + " only"
	}
	return strings.Join(parts, ",")
}

func hasRemoteSync(body StackInfo) bool {
	if body.RemoteRef != nil && strings.TrimSpace(*body.RemoteRef) != "" {
		return true
	}
	return body.Status == "published"
}

func bookmarkSyncLabel(body StackInfo) string {
	parts := []string{"local"}
	if hasRemoteSync(body) {
		parts = append(parts, "origin")
	}
	return "sync " + strings.Join(parts, ",")
}

func bookmarkSyncIcons(body StackInfo) string {
	icons := "⌂"
	if hasRemoteSync(body) {
		icons += "⇡"
	}
	return icons
}

func bookmarkMetaLine(body StackInfo, changeCount, approved int) string {
	parts := []string{
		"alias " + body.Alias,
		"base " + firstNonEmpty(body.BaseRef, "main"),
		fmtChanges(changeCount),
		fmtApproved(approved, changeCount),
		bookmarkSyncLabel(body),
	}
	if body.Status == "published" {
		parts = append(parts, "published")
	}
	return strings.Join(parts, " · ")
}

func fmtChanges(n int) string {
	if n == 1 {
		return "1 change"
	}
	return fmt.Sprintf("%d changes", n)
}

func fmtApproved(approved, total int) string {
	if total == 0 {
		return "0/0 approved"
	}
	return fmt.Sprintf("%d/%d approved", approved, total)
}

func repoLabel(repo RepoInfo) string {
	if repo.RemoteURL != nil {
		if full := githubRepoFullName(*repo.RemoteURL); full != "" {
			return full
		}
	}
	root := strings.TrimSpace(repo.RootPath)
	if root == "" {
		return "."
	}
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		if rel, relErr := filepath.Rel(home, root); relErr == nil && !strings.HasPrefix(rel, "..") {
			return "~/" + rel
		}
	}
	return root
}

func githubRepoFullName(remoteURL string) string {
	remoteURL = strings.TrimSpace(remoteURL)
	if idx := strings.Index(remoteURL, "github.com"); idx >= 0 {
		suffix := remoteURL[idx+len("github.com"):]
		suffix = strings.TrimPrefix(suffix, ":")
		suffix = strings.TrimPrefix(suffix, "/")
		suffix = strings.TrimSuffix(suffix, ".git")
		parts := strings.Split(suffix, "/")
		if len(parts) >= 2 && parts[0] != "" && parts[1] != "" {
			return parts[0] + "/" + parts[1]
		}
	}
	return ""
}
