package vcs

import (
	"errors"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestStackEditRevisionUsesBookmarkForPublishedStack(t *testing.T) {
	head := "rsnwoxsqtmmmrlnkqvyqrqwolywoqozo"
	stack := storage.Stack{
		BookmarkName: "feature/remove-logo-shimmer-animation-from-root-intro-pr",
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
		BookmarkName: "feature/login",
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
		"gx/base",
		"gx/base/worktree-abc123",
		"gx/edit",
		"gx/edit/worktree-abc123",
	} {
		if isGXStackBookmark(name) {
			t.Fatalf("isGXStackBookmark(%q) = true, want false", name)
		}
	}
}

func TestStackBookmarkNameUsesConventionalPrefixes(t *testing.T) {
	for _, tc := range []struct {
		name string
		want string
	}{
		{name: "add login flow", want: "feature/add-login-flow"},
		{name: "fix compose panic", want: "bug/fix-compose-panic"},
		{name: "docs readme", want: "docs/docs-readme"},
		{name: "test publish flow", want: "test/test-publish-flow"},
		{name: "chore release notes", want: "chore/chore-release-notes"},
		{name: "gx/legacy-name", want: "feature/legacy-name"},
		{name: "gx/draft/legacy-name", want: "feature/legacy-name"},
		{name: "feature/keep-prefix", want: "feature/keep-prefix"},
	} {
		if got := stackBookmarkName(tc.name, "abcdef123456"); got != tc.want {
			t.Fatalf("stackBookmarkName(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}
