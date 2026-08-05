package vcs

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	PlaceholderDescription      = "(no description set)"
	PendingRemainderDescription = "gx: pending remainder"
	// RevisionTrailerFormat is the gx identity trailer line template.
	RevisionTrailerFormat = "gx: https://gx.run/r/%s"
)

var (
	ErrNoRecordedAdds     = errors.New("no commit recorded for this stack")
	ErrEmptyCommitMessage = errors.New("commit message is required")
)

func IsPlaceholderDescription(desc string) bool {
	trimmed := strings.TrimSpace(desc)
	if trimmed == "" {
		return true
	}
	return strings.EqualFold(trimmed, PlaceholderDescription)
}

func ValidateCommitMessage(message string) error {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return fmt.Errorf("%w; use -m \"describe this revision\"", ErrEmptyCommitMessage)
	}
	if IsPlaceholderDescription(trimmed) {
		return fmt.Errorf("commit message cannot be %q", PlaceholderDescription)
	}
	return nil
}

func RevisionTrailerLine(revisionID string) string {
	revisionID = strings.TrimSpace(revisionID)
	if revisionID == "" {
		return ""
	}
	return fmt.Sprintf(RevisionTrailerFormat, revisionID)
}

// StampRevisionTrailer returns message carrying exactly one gx identity
// trailer for revisionID. gx trailers from a prior revision identity are
// replaced, never accumulated.
func StampRevisionTrailer(message, revisionID string) string {
	revisionID = strings.TrimSpace(revisionID)
	if revisionID == "" {
		return message
	}
	trailer := RevisionTrailerLine(revisionID)
	lines := strings.Split(message, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if isRevisionTrailerLine(line) {
			continue
		}
		kept = append(kept, line)
	}
	body := strings.TrimRight(strings.Join(kept, "\n"), "\n")
	if body == "" {
		return trailer
	}
	return body + "\n\n" + trailer
}

func isRevisionTrailerLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), revisionTrailerPrefix())
}

// revisionTrailerPresent reports whether message already carries the canonical
// gx trailer for revisionID and no stale trailers, i.e. re-stamping would be a
// no-op.
func revisionTrailerPresent(message, revisionID string) bool {
	return StampRevisionTrailer(message, revisionID) == strings.TrimRight(message, "\n")
}

func revisionTrailerPrefix() string {
	return fmt.Sprintf(RevisionTrailerFormat, "")
}

func FinalizeCommitMessage(message, revisionID string) (string, error) {
	if err := ValidateCommitMessage(message); err != nil {
		return "", err
	}
	return StampRevisionTrailer(message, revisionID), nil
}

func validateRecordedChangeDescription(description string) error {
	if err := ValidateCommitMessage(description); err != nil {
		return fmt.Errorf("recorded revision has invalid description: %w", err)
	}
	return nil
}

func (s *Service) requireRecordedAddsForPublish(ctx context.Context, stack StackInfo) error {
	if stack.ID == 0 {
		return fmt.Errorf(
			"%w: stage changes with `git add`, then run `git commit -m \"describe this revision\"` before `git push`",
			ErrNoRecordedAdds,
		)
	}
	store, err := openStore(ctx)
	if err != nil {
		return err
	}
	defer store.Close()

	changes, err := store.ListChangesByStackID(ctx, stack.ID)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		return fmt.Errorf(
			"%w: stage changes with `git add`, then run `git commit -m \"...\"` to record at least one change before `git push`",
			ErrNoRecordedAdds,
		)
	}
	for _, change := range changes {
		if !IsPlaceholderDescription(change.Description) {
			return nil
		}
	}
	return fmt.Errorf(
		"%w: recorded changes have no description; stage changes with `git add`, then run `git commit -m \"...\"`",
		ErrNoRecordedAdds,
	)
}
