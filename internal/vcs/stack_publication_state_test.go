package vcs

import (
	"context"
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
}

func TestStackMergedIntoBaseUsesLiveBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "(gx/work) & ancestors(main)", "-n", "1", "--no-graph", "-T", "change_id"): {
				"workchange\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergedIntoBase(context.Background(), repoRoot, StackInfo{
		BookmarkName: "gx/work",
		BaseRef:      "main",
		Status:       "draft",
	}, map[string]string{"gx/work": "workchange"})
	if !merged {
		t.Fatal("expected live bookmark ancestor of main to be merged")
	}
}

func TestStackMergedIntoBaseIgnoresStoredHeadWhenBookmarkExists(t *testing.T) {
	repoRoot := t.TempDir()
	headCommit := "abc123"
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "(gx/work) & ancestors(main)", "-n", "1", "--no-graph", "-T", "change_id"): {
				"",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	merged := svc.stackMergedIntoBase(context.Background(), repoRoot, StackInfo{
		BookmarkName: "gx/work",
		BaseRef:      "main",
		HeadCommitID: &headCommit,
		Status:       "draft",
	}, map[string]string{"gx/work": "workchange"})
	if merged {
		t.Fatal("stored head should not classify a stack as merged while its live bookmark is not merged")
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

	merged := svc.stackMergedIntoBase(context.Background(), repoRoot, StackInfo{
		BookmarkName: "gx/work",
		BaseRef:      "main",
		HeadCommitID: &headCommit,
		Status:       "draft",
	}, nil)
	if !merged {
		t.Fatal("expected stored head ancestor of main to be merged when bookmark is gone")
	}
}
