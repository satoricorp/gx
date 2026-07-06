package vcs

import (
	"context"
	"errors"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestValidateCommitMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		wantErr error
	}{
		{name: "ok", message: "add login flow"},
		{name: "whitespace", message: "   ", wantErr: ErrEmptyCommitMessage},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCommitMessage(tc.message)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ValidateCommitMessage() error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateCommitMessage() error = %v", err)
			}
		})
	}

	if err := ValidateCommitMessage(PlaceholderDescription); err == nil {
		t.Fatal("ValidateCommitMessage() error = nil, want placeholder rejection")
	}
}

func TestValidateRecordedChangeDescription(t *testing.T) {
	t.Parallel()

	if err := validateRecordedChangeDescription("add login flow"); err != nil {
		t.Fatalf("validateRecordedChangeDescription() error = %v", err)
	}
	if err := validateRecordedChangeDescription(""); err == nil {
		t.Fatal("validateRecordedChangeDescription() error = nil, want empty rejection")
	}
}

func TestRequireRecordedAddsForPublish(t *testing.T) {
	ctx := context.Background()
	t.Setenv("GX_HOME", t.TempDir())

	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}

	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:  t.TempDir(),
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}

	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "login",
		BookmarkName: "feature/login",
		BaseRef:      "main",
		BaseCommitID: "base",
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	svc := &Service{}
	if err := svc.requireRecordedAddsForPublish(ctx, StackInfo{ID: stackID}); !errors.Is(err, ErrNoRecordedAdds) {
		t.Fatalf("requireRecordedAddsForPublish() error = %v, want %v", err, ErrNoRecordedAdds)
	}

	db, err = storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err = storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}

	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "chg1",
		CurrentCommitID: "commit1",
		Description:     PlaceholderDescription,
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	if err := svc.requireRecordedAddsForPublish(ctx, StackInfo{ID: stackID}); !errors.Is(err, ErrNoRecordedAdds) {
		t.Fatalf("requireRecordedAddsForPublish() placeholder error = %v, want %v", err, ErrNoRecordedAdds)
	}

	db, err = storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err = storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	if _, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "chg1",
		CurrentCommitID: "commit1",
		Description:     "add login flow",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       2,
	}); err != nil {
		t.Fatalf("UpsertChange() update error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	if err := svc.requireRecordedAddsForPublish(ctx, StackInfo{ID: stackID}); err != nil {
		t.Fatalf("requireRecordedAddsForPublish() error = %v, want nil", err)
	}
}
