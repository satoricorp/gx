package vcs

import (
	"regexp"
	"strings"

	"github.com/satoricorp/lgtm/internal/storage"
)

var bookmarkSlugPattern = regexp.MustCompile(`[^a-z0-9]+`)

type stackEntry struct {
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

func stackRevisionsFromEntries(entries []stackEntry, activeChangeID string, publishedThrough int64, changeByRevision map[string]storage.Change) ([]RevisionSummary, int) {
	revisions := make([]RevisionSummary, 0, len(entries))
	publishedCount := 0
	for _, entry := range entries {
		description := firstNonEmpty(entry.Description, "(no description set)")
		if IsPlaceholderDescription(description) {
			continue
		}
		row, ok := changeByRevision[entry.ChangeID]
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
