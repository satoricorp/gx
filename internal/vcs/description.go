package vcs

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	PlaceholderDescription       = "(no description set)"
	PendingRemainderDescription  = "gx: pending remainder"
)

var (
	ErrNoRecordedAdds     = errors.New("no gx add recorded for this stack")
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

func validateRecordedChangeDescription(description string) error {
	if err := ValidateCommitMessage(description); err != nil {
		return fmt.Errorf("recorded revision has invalid description: %w", err)
	}
	return nil
}

func (s *Service) requireRecordedAddsForPublish(ctx context.Context, stack StackInfo) error {
	if stack.ID == 0 {
		return fmt.Errorf(
			"%w: run `gx add -m \"describe this revision\"` before `gx publish`",
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
			"%w: run `gx add -m \"...\"` to record at least one change before `gx publish`",
			ErrNoRecordedAdds,
		)
	}
	for _, change := range changes {
		if !IsPlaceholderDescription(change.Description) {
			return nil
		}
	}
	return fmt.Errorf(
		"%w: recorded changes have no description; run `gx add -m \"...\"`",
		ErrNoRecordedAdds,
	)
}
