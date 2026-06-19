package vcs

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/storage"
)

const (
	jjStackLineTmpl    = `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`
	jjBookmarkListTmpl = `name ++ "|" ++ normal_target.change_id() ++ "\n"`
)

var bookmarkSlugPattern = regexp.MustCompile(`[^a-z0-9]+`)

type jjStackEntry struct {
	ChangeID       string
	CommitID       string
	Description    string
	ParentChangeID *string
}

func bookmarkSlug(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if idx := strings.IndexAny(value, "\n\r"); idx >= 0 {
		value = value[:idx]
	}
	value = strings.ToLower(value)
	value = bookmarkSlugPattern.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if len(value) > 48 {
		value = strings.Trim(value[:48], "-")
	}
	return value
}

func jjStackRevset(container string) string {
	return jjStackRevsetForHead(container, "@")
}

func jjStackRevsetForHead(container, head string) string {
	container = strings.TrimSpace(container)
	head = strings.TrimSpace(head)
	if head == "" {
		head = "@"
	}
	base := "mutable() & ~empty() & ~hidden() & ancestors(" + quoteJJRev(head) + ")"
	if container == "" {
		return base
	}
	return base + " & ~ancestors(" + quoteJJRev(container) + ")"
}

func quoteJJRev(name string) string {
	if strings.ContainsAny(name, " \t\n\"()&|") {
		return `"` + strings.ReplaceAll(name, `"`, `\"`) + `"`
	}
	return name
}

func (s *Service) jjStackEntries(ctx context.Context, repoRoot, container string) ([]jjStackEntry, error) {
	return s.jjStackEntriesForHead(ctx, repoRoot, container, "@")
}

func (s *Service) jjStackEntriesForHead(ctx context.Context, repoRoot, container, head string) ([]jjStackEntry, error) {
	revset := jjStackRevsetForHead(container, head)
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", revset, "--reversed", "--no-graph", "-T", jjStackLineTmpl)
	if err != nil {
		return nil, err
	}
	lines := splitLines(out)
	entries := make([]jjStackEntry, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) < 3 {
			continue
		}
		entry := jjStackEntry{
			ChangeID:    parts[0],
			CommitID:    parts[1],
			Description: parts[2],
		}
		if len(parts) == 4 && strings.TrimSpace(parts[3]) != "" {
			first := strings.SplitN(parts[3], ",", 2)[0]
			entry.ParentChangeID = &first
		}
		if entry.ChangeID == "" || entry.CommitID == "" {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *Service) isChangeHidden(ctx context.Context, repoRoot, changeID string) (bool, error) {
	value, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", changeID, "-n", "1", "--no-graph", "-T", "hidden")
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no matching commits") {
			return true, nil
		}
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(lastNonEmptyLine(value)), "true"), nil
}

func (s *Service) changeExistsInJJ(ctx context.Context, repoRoot, changeID string) bool {
	_, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "log", "-r", changeID, "-n", "1", "--no-graph", "-T", "change_id")
	return err == nil
}

func (s *Service) jjContainerBookmark(ctx context.Context, repoRoot string) (string, error) {
	workChangeID, err := s.currentWorkChangeID(ctx, repoRoot)
	if err != nil || workChangeID == "" {
		return "", err
	}
	out, err := s.runStdoutTrimmed(ctx, repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl)
	if err != nil {
		return "", err
	}
	var fallback string
	for _, line := range splitLines(out) {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		target := strings.TrimSpace(parts[1])
		if name == "" || strings.HasPrefix(name, "gx/") {
			continue
		}
		if target == workChangeID {
			return name, nil
		}
		if fallback == "" {
			fallback = name
		}
	}
	return fallback, nil
}

func (s *Service) reconcileRepoChanges(ctx context.Context, repo RepoInfo, container string) error {
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
		return nil
	}

	entries, err := s.jjStackEntries(ctx, repo.RootPath, container)
	if err != nil {
		return err
	}
	activeIDs := make(map[string]jjStackEntry, len(entries))
	for _, entry := range entries {
		activeIDs[entry.ChangeID] = entry
	}

	changes, err := store.ListChangesByRepoID(ctx, repoRow.ID)
	if err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	for _, change := range changes {
		entry, onStack := activeIDs[change.JJChangeID]
		switch {
		case onStack:
			if change.Status == "abandoned" || change.CurrentCommitID != entry.CommitID || change.Description != entry.Description {
				parent := entry.ParentChangeID
				if _, err := store.UpsertChange(ctx, storage.Change{
					RepoID:          repoRow.ID,
					JJChangeID:      entry.ChangeID,
					CurrentCommitID: entry.CommitID,
					Description:     firstNonEmpty(entry.Description, change.Description, "(no description set)"),
					ParentChangeID:  parent,
					Status:          "draft",
					FirstSeenAt:     change.FirstSeenAt,
					UpdatedAt:       now,
				}); err != nil {
					return err
				}
			}
		case !s.changeExistsInJJ(ctx, repo.RootPath, change.JJChangeID):
			if change.Status != "abandoned" {
				if _, err := store.UpsertChange(ctx, storage.Change{
					RepoID:          repoRow.ID,
					JJChangeID:      change.JJChangeID,
					CurrentCommitID: change.CurrentCommitID,
					Description:     change.Description,
					ParentChangeID:  change.ParentChangeID,
					Status:          "abandoned",
					FirstSeenAt:     change.FirstSeenAt,
					UpdatedAt:       now,
				}); err != nil {
					return err
				}
			}
		default:
			hidden, err := s.isChangeHidden(ctx, repo.RootPath, change.JJChangeID)
			if err != nil {
				return err
			}
			if hidden && change.Status != "abandoned" {
				if _, err := store.UpsertChange(ctx, storage.Change{
					RepoID:          repoRow.ID,
					JJChangeID:      change.JJChangeID,
					CurrentCommitID: change.CurrentCommitID,
					Description:     change.Description,
					ParentChangeID:  change.ParentChangeID,
					Status:          "abandoned",
					FirstSeenAt:     change.FirstSeenAt,
					UpdatedAt:       now,
				}); err != nil {
					return err
				}
			}
		}
	}

	for _, entry := range entries {
		if _, err := upsertChange(ctx, store, repoRow.ID, ChangeInfo{
			ChangeID:       entry.ChangeID,
			CommitID:       entry.CommitID,
			Description:    entry.Description,
			ParentChangeID: entry.ParentChangeID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func stackRevisionsFromEntries(entries []jjStackEntry, activeChangeID string, publishedThrough int64, changeByJJ map[string]storage.Change) ([]RevisionSummary, int) {
	revisions := make([]RevisionSummary, 0, len(entries))
	publishedCount := 0
	for _, entry := range entries {
		description := firstNonEmpty(entry.Description, "(no description set)")
		if IsPlaceholderDescription(description) {
			continue
		}
		row, ok := changeByJJ[entry.ChangeID]
		status := "draft"
		if ok {
			status = row.Status
		}
		if status == "abandoned" {
			continue
		}
		published := ok && ChangePublishedThrough(row, publishedThrough)
		if published {
			publishedCount++
		}
		revisions = append(revisions, RevisionSummary{
			Index:       len(revisions) + 1,
			ChangeID:    entry.ChangeID,
			CommitID:    entry.CommitID,
			Description: description,
			Status:      status,
			Active:      activeChangeID != "" && entry.ChangeID == activeChangeID,
			Published:   published,
		})
	}
	return revisions, publishedCount
}
