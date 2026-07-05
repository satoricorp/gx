package vcs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestCleanupStaleStacksMarksMissingBookmarkStacksClosed(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	healthyID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "healthy",
		BookmarkName: "feature/healthy",
		BaseRef:      "main",
		BaseCommitID: "base",
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack(healthy) error = %v", err)
	}
	staleID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "stale",
		BookmarkName: "feature/env-example",
		BaseRef:      "main",
		BaseCommitID: "base",
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack(stale) error = %v", err)
	}
	mergedID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "merged",
		BookmarkName: "feature/merged",
		BaseRef:      "main",
		BaseCommitID: "base",
		Status:       "merged",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack(merged) error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	revisionQuery := func(rev string) string {
		return runnerKey(cwd, "jj", "log", "-r", rev, "--limit", "1", "--no-graph", "-T", "commit_id")
	}
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):     {cwd + "\n", cwd + "\n"},
			revisionQuery("feature/healthy"): {"abc123\n", "abc123\n"},
		},
		errors: map[string][]error{
			revisionQuery("feature/env-example"): {
				fmt.Errorf(`Error: Revision "feature/env-example" doesn't exist`),
				fmt.Errorf(`Error: Revision "feature/env-example" doesn't exist`),
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	listed, err := svc.ListStaleStacks(ctx)
	if err != nil {
		t.Fatalf("ListStaleStacks() error = %v", err)
	}
	if len(listed.Stale) != 1 || listed.Stale[0].BookmarkName != "feature/env-example" {
		t.Fatalf("ListStaleStacks() = %#v, want feature/env-example stale", listed)
	}
	if len(listed.Actions) != 0 {
		t.Fatalf("ListStaleStacks() took actions: %#v", listed.Actions)
	}

	result, err := svc.CleanupStaleStacks(ctx)
	if err != nil {
		t.Fatalf("CleanupStaleStacks() error = %v", err)
	}
	if len(result.Stale) != 1 || result.Stale[0].BookmarkName != "feature/env-example" {
		t.Fatalf("CleanupStaleStacks() = %#v, want feature/env-example stale", result)
	}
	if len(result.Actions) != 1 {
		t.Fatalf("CleanupStaleStacks() actions = %#v, want one repair action", result.Actions)
	}

	db, err = storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() reopen error = %v", err)
	}
	store, err = storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() reopen error = %v", err)
	}
	defer store.Close()
	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		t.Fatalf("ListStacksByRepoID() error = %v", err)
	}
	statuses := map[int64]string{}
	for _, stack := range stacks {
		statuses[stack.ID] = stack.Status
	}
	if statuses[staleID] != "closed" {
		t.Fatalf("stale stack status = %q, want closed", statuses[staleID])
	}
	if statuses[healthyID] != "draft" {
		t.Fatalf("healthy stack status = %q, want draft", statuses[healthyID])
	}
	if statuses[mergedID] != "merged" {
		t.Fatalf("merged stack status = %q, want merged", statuses[mergedID])
	}
}
