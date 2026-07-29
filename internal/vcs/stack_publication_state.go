package vcs

import (
	"strings"

	"github.com/satoricorp/totality/internal/storage"
)

type StackPublicationState struct {
	PublishedThrough     int64
	LastPushedByChangeID map[int64]string
	pastPublishedHead    bool
}

func NewStackPublicationState(publishedThrough int64, lastPushedByChangeID map[int64]string) StackPublicationState {
	if lastPushedByChangeID == nil {
		lastPushedByChangeID = map[int64]string{}
	}
	return StackPublicationState{
		PublishedThrough:     publishedThrough,
		LastPushedByChangeID: lastPushedByChangeID,
	}
}

func (s *StackPublicationState) RevisionPublished(change storage.Change) bool {
	if s == nil {
		return false
	}
	if s.LastPushedByChangeID[change.ID] == change.CurrentCommitID {
		return true
	}
	if s.PublishedThrough == -1 || s.pastPublishedHead {
		return false
	}
	published := change.ID <= s.PublishedThrough
	if change.ID == s.PublishedThrough {
		s.pastPublishedHead = true
	}
	return published
}

func ChangePublishedThrough(change storage.Change, publishedThrough int64) bool {
	return publishedThrough != -1 && change.ID <= publishedThrough
}

func CountPublishedChanges(changes []storage.Change, publishedThrough int64) int {
	count := 0
	for _, change := range changes {
		if ChangePublishedThrough(change, publishedThrough) {
			count++
		}
	}
	return count
}

func DeriveStackStatus(storedStatus string, revisionCount, publishedCount int, mergedIntoBase bool) string {
	status := strings.TrimSpace(storedStatus)
	if IsTerminalStackStatus(status) {
		return status
	}
	if mergedIntoBase {
		return "merged"
	}
	if revisionCount == 0 {
		return status
	}
	if publishedCount == revisionCount {
		return "published"
	}
	return "draft"
}

func IsTerminalStackStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "merged", "closed":
		return true
	default:
		return false
	}
}
