package gxtest_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
)

// TestNewRepoIsAMainCheckout pins the identity a main checkout has: its git dir
// IS its common dir. Repositories in this shape are the ones every existing
// test builds.
func TestNewRepoIsAMainCheckout(t *testing.T) {
	repo := gxtest.NewWorld(t).NewRepo(t)
	if repo.GitDir != repo.GitCommonDir {
		t.Fatalf("main checkout git_dir = %q, common_dir = %q, want equal", repo.GitDir, repo.GitCommonDir)
	}
	if want := filepath.Join(repo.Root, ".git"); repo.GitCommonDir != want {
		t.Fatalf("common_dir = %q, want %q", repo.GitCommonDir, want)
	}
	if repo.IsLinkedWorktree() {
		t.Fatal("main checkout reported as a linked worktree")
	}
}

// TestAddWorktreeSharesTheMainCommonDir pins the state the session bug lived
// in: a worktree with its own toplevel and its own git dir, sharing the main
// repository's common dir. Repository identity in gx is that shared common dir,
// so a push from here writes rows the pushing directory's own path can never
// find.
func TestAddWorktreeSharesTheMainCommonDir(t *testing.T) {
	world := gxtest.NewWorld(t)
	main := world.NewRepo(t)
	worktree := main.AddWorktree(t, filepath.Join(t.TempDir(), "linked"), "feature")

	if worktree.Root == main.Root {
		t.Fatalf("worktree root = %q, want a different toplevel from %q", worktree.Root, main.Root)
	}
	if worktree.GitCommonDir != main.GitCommonDir {
		t.Fatalf("worktree common_dir = %q, want the main repo's %q", worktree.GitCommonDir, main.GitCommonDir)
	}
	if worktree.GitDir == worktree.GitCommonDir {
		t.Fatalf("worktree git_dir = common_dir = %q; a linked worktree has its own git dir", worktree.GitDir)
	}
	if !worktree.IsLinkedWorktree() {
		t.Fatal("linked worktree not reported as one")
	}
}

// TestCommitStampsATrailer covers the property every session assertion built on
// this harness depends on: without the GX trailer no `changes` row is created,
// so sessions have nothing to attach to and the assertions pass vacuously.
func TestCommitStampsATrailer(t *testing.T) {
	world := gxtest.NewWorld(t)
	repo := world.NewRepo(t)
	commit := repo.Commit(t, map[string]string{"alpha.go": "package alpha\n"}, "add alpha")

	if commit.SHA == "" || commit.RevisionID == "" {
		t.Fatalf("Commit() = %+v, want both a sha and a revision id", commit)
	}
	message := gxtest.GitOutput(t, repo.Root, "log", "-1", "--pretty=%B")
	if want := gxtest.RevisionTrailerLine(commit.RevisionID); !strings.Contains(message, want) {
		t.Fatalf("commit message = %q, want it to carry %q", message, want)
	}
}
