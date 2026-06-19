package vcs

import (
	"errors"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestStackEditRevisionUsesBookmarkForPublishedStack(t *testing.T) {
	head := "rsnwoxsqtmmmrlnkqvyqrqwolywoqozo"
	stack := storage.Stack{
		BookmarkName: "gx/remove-logo-shimmer-animation-from-root-intro-pr",
		HeadChangeID: &head,
		Status:       "published",
	}
	if got := stackEditRevision(stack); got != stack.BookmarkName {
		t.Fatalf("stackEditRevision() = %q, want bookmark %q", got, stack.BookmarkName)
	}
}

func TestStackEditRevisionUsesHeadForDraftStack(t *testing.T) {
	head := "chg123"
	commit := "commit123"
	stack := storage.Stack{
		BookmarkName: "gx/login",
		HeadChangeID: &head,
		HeadCommitID: &commit,
		Status:       "draft",
	}
	if got := stackEditRevision(stack); got != commit {
		t.Fatalf("stackEditRevision() = %q, want commit %q", got, commit)
	}
}

func TestStackEditRevisionFallsBackToHeadChangeForLegacyDraftStack(t *testing.T) {
	head := "chg123"
	stack := storage.Stack{
		BookmarkName: "gx/login",
		HeadChangeID: &head,
		Status:       "draft",
	}
	if got := stackEditRevision(stack); got != head {
		t.Fatalf("stackEditRevision() = %q, want change %q", got, head)
	}
}

func TestIsJJImmutableError(t *testing.T) {
	if isJJImmutableError(nil) {
		t.Fatal("expected false for nil error")
	}
	if !isJJImmutableError(errors.New("Commit abc is immutable")) {
		t.Fatal("expected true for immutable jj error")
	}
}

func TestInternalCheckoutRefsAreNotStackBookmarks(t *testing.T) {
	for _, name := range []string{
		gxInternalBaseRef,
		gxInternalBaseRef + "/worktree-abc123",
		gxInternalEditRef,
		gxInternalEditRef + "/worktree-abc123",
	} {
		if isGXStackBookmark(name) {
			t.Fatalf("isGXStackBookmark(%q) = true, want false", name)
		}
	}
}
