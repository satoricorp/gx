package vcs

import (
	"context"
	"fmt"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestStackPublicationStatePublishesThroughHeadInOrder(t *testing.T) {
	state := NewStackPublicationState(2, nil)
	changes := []storage.Change{
		{ID: 1, CurrentCommitID: "commit-one"},
		{ID: 2, CurrentCommitID: "commit-two"},
		{ID: 3, CurrentCommitID: "commit-three"},
	}
	got := []bool{
		state.RevisionPublished(changes[0]),
		state.RevisionPublished(changes[1]),
		state.RevisionPublished(changes[2]),
	}
	want := []bool{true, true, false}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("published[%d] = %t, want %t", i, got[i], want[i])
		}
	}
}

func TestStackPublicationStateUsesExactPushedCommit(t *testing.T) {
	state := NewStackPublicationState(-1, map[int64]string{2: "commit-two"})
	if state.RevisionPublished(storage.Change{ID: 2, CurrentCommitID: "commit-two"}) != true {
		t.Fatal("expected exact pushed commit to be published")
	}
	if state.RevisionPublished(storage.Change{ID: 2, CurrentCommitID: "commit-new"}) != false {
		t.Fatal("expected changed commit to be unpublished")
	}
}

func TestDeriveStackStatusPreservesTerminalStatus(t *testing.T) {
	if got := DeriveStackStatus("merged", 2, 2, false); got != "merged" {
		t.Fatalf("DeriveStackStatus(merged) = %q", got)
	}
	if got := DeriveStackStatus("draft", 2, 2, false); got != "published" {
		t.Fatalf("DeriveStackStatus(draft, all published) = %q", got)
	}
	if got := DeriveStackStatus("published", 2, 1, false); got != "draft" {
		t.Fatalf("DeriveStackStatus(published, partial) = %q", got)
	}
	if got := DeriveStackStatus("draft", 2, 0, true); got != "merged" {
		t.Fatalf("DeriveStackStatus(merged into base) = %q", got)
	}
	if got := DeriveStackStatus("published", 2, 2, true); got != "merged" {
		t.Fatalf("DeriveStackStatus(published, merged into base) = %q", got)
	}
}

func TestStackMergedIntoBaseUsesLiveBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "(feature/work) & ancestors(main)", "-n", "1", "--no-graph", "-T", "change_id"): {
				"workchange\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergedIntoBase(context.Background(), repoRoot, "main", StackInfo{
		BookmarkName: "feature/work",
		BaseRef:      "main",
		Status:       "draft",
	}, map[string]string{"feature/work": "workchange"})
	if !merged {
		t.Fatal("expected live bookmark ancestor of main to be merged")
	}
}

func TestStackMergedIntoBaseIgnoresStoredHeadWhenBookmarkExists(t *testing.T) {
	repoRoot := t.TempDir()
	headCommit := "abc123"
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "(feature/work) & ancestors(main)", "-n", "1", "--no-graph", "-T", "change_id"): {
				"",
			},
		},
		errors: map[string][]error{
			runnerKey(repoRoot, "jj", "log", "-r", "main@origin", "-n", "1", "--no-graph", "-T", "change_id"): {
				fmt.Errorf("revision not found"),
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergedIntoBase(context.Background(), repoRoot, "main", StackInfo{
		BookmarkName: "feature/work",
		BaseRef:      "main",
		HeadCommitID: &headCommit,
		Status:       "draft",
	}, map[string]string{"feature/work": "workchange"})
	if merged {
		t.Fatal("stored head should not classify a stack as merged while its live bookmark is not merged")
	}
}

func TestStackMergedIntoBaseUsesRemoteTrackingBase(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "main@origin", "-n", "1", "--no-graph", "-T", "change_id"): {
				"basechange\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "(feature/work) & ancestors(main)", "-n", "1", "--no-graph", "-T", "change_id"): {
				"",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "(feature/work) & ancestors(main@origin)", "-n", "1", "--no-graph", "-T", "change_id"): {
				"workchange\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergedIntoBase(context.Background(), repoRoot, "main", StackInfo{
		BookmarkName: "feature/work",
		BaseRef:      "main",
		Status:       "published",
	}, map[string]string{"feature/work": "workchange"})
	if !merged {
		t.Fatal("expected stack merged into remote-tracking base to be merged")
	}
}

func TestStackMergedIntoBaseFallsBackToStoredHeadWithoutBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	headCommit := "abc123"
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "(abc123) & ancestors(main)", "-n", "1", "--no-graph", "-T", "change_id"): {
				"workchange\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergedIntoBase(context.Background(), repoRoot, "main", StackInfo{
		BookmarkName: "feature/work",
		BaseRef:      "main",
		HeadCommitID: &headCommit,
		Status:       "draft",
	}, nil)
	if !merged {
		t.Fatal("expected stored head ancestor of main to be merged when bookmark is gone")
	}
}

func TestStackMergeStatesUsesDefaultBranchWhenStackBaseMissing(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "main@origin", "-n", "1", "--no-graph", "-T", "change_id"): {
				"basechange\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "(feature/work) & ancestors(feature/authoring) | (feature/work) & ancestors(main) | (feature/work) & ancestors(main@origin)", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "\n"`): {
				"work-change|work-commit\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergeStates(context.Background(), repoRoot, "main", []StackInfo{
		{ID: 1, BookmarkName: "feature/work", BaseRef: "feature/authoring", Status: "published"},
	}, map[string]string{
		"feature/work": "work-change",
	})
	if !merged[1] {
		t.Fatal("expected stack merged into default branch when stack base ref is missing")
	}
}

func TestStackMergeStatesBatchesBookmarkChecks(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "main@origin", "-n", "1", "--no-graph", "-T", "change_id"): {
				"basechange\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "(feature/one) & ancestors(main) | (feature/one) & ancestors(main@origin) | (feature/two) & ancestors(main) | (feature/two) & ancestors(main@origin)", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "\n"`): {
				"one-change|one-commit\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergeStates(context.Background(), repoRoot, "main", []StackInfo{
		{ID: 1, BookmarkName: "feature/one", BaseRef: "main", Status: "draft"},
		{ID: 2, BookmarkName: "feature/two", BaseRef: "main", Status: "draft"},
	}, map[string]string{
		"feature/one": "one-change",
		"feature/two": "two-change",
	})

	if !merged[1] {
		t.Fatal("expected feature/one to be merged")
	}
	if merged[2] {
		t.Fatal("expected feature/two to remain unmerged")
	}
	if len(runner.calls) != 3 {
		t.Fatalf("merge checks made %d runner calls, want 3: %v", len(runner.calls), runner.calls)
	}
}
