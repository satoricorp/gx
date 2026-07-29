package repopath_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/satoricorp/totality/internal/capture/repopath"
)

// An agent working in a linked worktree edits <worktree>/internal/cli/x.go,
// and the commit that lands contains internal/cli/x.go. Relativizing against
// the checkout that ran the push must produce the committed path, not the path
// through the worktree.
//
// The nesting matters: with the worktree inside the main checkout, the naive
// answer (".claude/worktrees/feature/internal/cli/x.go") is a clean relative
// path with no "..", so nothing downstream can tell it is wrong. It just never
// matches a hunk, and the session disappears.
func TestRelResolvesEditsMadeInANestedWorktree(t *testing.T) {
	repopath.ResetCacheForTest()
	main, worktree := repoWithNestedWorktree(t)

	edited := filepath.Join(worktree, "internal", "cli", "x.go")
	if got := repopath.Rel(edited, main); got != "internal/cli/x.go" {
		t.Fatalf("Rel(worktree edit) = %q, want %q", got, "internal/cli/x.go")
	}

	// The main checkout still behaves as it always did.
	if got := repopath.Rel(filepath.Join(main, "internal", "cli", "x.go"), main); got != "internal/cli/x.go" {
		t.Fatalf("Rel(main edit) = %q, want %q", got, "internal/cli/x.go")
	}
}

// A worktree outside the repository is the same repository too, and capture
// runs from whichever checkout pushed.
func TestRelResolvesEditsMadeInAnExternalWorktree(t *testing.T) {
	repopath.ResetCacheForTest()
	main := initRepo(t)
	outside := filepath.Join(t.TempDir(), "elsewhere")
	runGit(t, main, "worktree", "add", "-b", "external", outside)

	edited := filepath.Join(outside, "internal", "cli", "x.go")
	if got := repopath.Rel(edited, main); got != "internal/cli/x.go" {
		t.Fatalf("Rel(external worktree edit) = %q, want %q", got, "internal/cli/x.go")
	}
}

// Accepting every checkout must not turn into accepting everything: a session
// that edited an unrelated repository stays unattributable.
func TestRelLeavesPathsOutsideTheRepositoryAlone(t *testing.T) {
	repopath.ResetCacheForTest()
	main, _ := repoWithNestedWorktree(t)
	other := filepath.Join(t.TempDir(), "other-repo", "internal", "cli", "x.go")

	got := repopath.Rel(other, main)
	if got == "internal/cli/x.go" {
		t.Fatalf("Rel folded an unrelated repo's path into this repo: %q", got)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("Rel(outside) = %q, want the path left absolute", got)
	}
}

func repoWithNestedWorktree(t *testing.T) (main, worktree string) {
	t.Helper()
	main = initRepo(t)
	worktree = filepath.Join(main, ".claude", "worktrees", "feature")
	runGit(t, main, "worktree", "add", "-b", "feature", worktree)
	return main, worktree
}

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// macOS hands out /var symlinks for temp dirs; resolve so the path the
	// test compares against is the one git reports.
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}
	runGit(t, dir, "init", "-b", "main")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-m", "initial")
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
