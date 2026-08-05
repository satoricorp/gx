package repobind_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/satoricorp/lgtm/internal/capture/repobind"
)

// A session run in a linked worktree belongs to the repository, not to the
// checkout. Agent sessions in this project run in .claude/worktrees/<name>, so
// this is the ordinary case rather than an exotic one.
func TestForDirectoryBindsWorktreesToOneRepository(t *testing.T) {
	main := initRepo(t, "git@github.com:satoricorp/lgtm.git")
	worktree := filepath.Join(main, ".claude", "worktrees", "feature")
	runGit(t, main, "worktree", "add", "-q", "-b", "feature", worktree)

	r := repobind.NewResolver()
	fromMain := r.ForDirectory(context.Background(), main)
	fromWorktree := r.ForDirectory(context.Background(), worktree)

	if !fromMain.Bound() || !fromWorktree.Bound() {
		t.Fatalf("expected both to bind: main=%+v worktree=%+v", fromMain, fromWorktree)
	}
	if fromMain.Origin != fromWorktree.Origin {
		t.Fatalf("worktree bound to a different repository: %q vs %q", fromMain.Origin, fromWorktree.Origin)
	}
	if fromMain.Origin != "github.com/satoricorp/lgtm" {
		t.Fatalf("Origin = %q", fromMain.Origin)
	}
	// The checkouts differ; the repository does not.
	if fromMain.RepoRoot == fromWorktree.RepoRoot {
		t.Fatal("expected distinct checkout roots")
	}
	if fromMain.GitCommonDir != fromWorktree.GitCommonDir {
		t.Fatalf("worktrees disagree on the repository: %q vs %q", fromMain.GitCommonDir, fromWorktree.GitCommonDir)
	}
}

// Subdirectories bind to the same repository as the root, since a session's cwd
// is often somewhere inside the tree.
func TestForDirectoryBindsSubdirectories(t *testing.T) {
	root := initRepo(t, "https://github.com/satoricorp/console.git")
	sub := filepath.Join(root, "internal", "cli")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	r := repobind.NewResolver()
	if got := r.ForDirectory(context.Background(), sub); got.Origin != "github.com/satoricorp/console" {
		t.Fatalf("Origin = %q, want github.com/satoricorp/console", got.Origin)
	}
}

// The cases that must stay unbound. Each of these is a directory whose sessions
// have to remain local: unbound is the safe answer, and nothing may fall back
// to guessing identity from the path.
func TestForDirectoryLeavesUnidentifiableDirectoriesUnbound(t *testing.T) {
	r := repobind.NewResolver()
	ctx := context.Background()

	t.Run("not a git repository", func(t *testing.T) {
		if got := r.ForDirectory(ctx, t.TempDir()); got.Bound() {
			t.Fatalf("plain directory bound: %+v", got)
		}
	})

	t.Run("git repository with no remote", func(t *testing.T) {
		// A personal side project with no origin. It must not bind.
		local := initRepo(t, "")
		got := r.ForDirectory(ctx, local)
		if got.Bound() {
			t.Fatalf("remote-less repository bound: %+v", got)
		}
		if got.RepoRoot == "" {
			t.Fatal("expected the repo root to still resolve for diagnostics")
		}
	})

	t.Run("empty and missing paths", func(t *testing.T) {
		for _, dir := range []string{"", "   ", filepath.Join(t.TempDir(), "does-not-exist")} {
			if got := r.ForDirectory(ctx, dir); got.Bound() {
				t.Fatalf("ForDirectory(%q) bound: %+v", dir, got)
			}
		}
	})
}

// Two repositories must never bind to the same identity.
func TestForDirectoryKeepsRepositoriesApart(t *testing.T) {
	tx := initRepo(t, "git@github.com:satoricorp/lgtm.git")
	yeet := initRepo(t, "git@github.com:joe/yeet.git")

	r := repobind.NewResolver()
	ctx := context.Background()
	a := r.ForDirectory(ctx, tx)
	b := r.ForDirectory(ctx, yeet)
	if a.Origin == b.Origin {
		t.Fatalf("distinct repositories share an origin: %q", a.Origin)
	}
}

func initRepo(t *testing.T, origin string) string {
	t.Helper()
	dir := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(dir); err == nil {
		dir = resolved
	}
	runGit(t, dir, "init", "-q", "-b", "main")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
	if origin != "" {
		runGit(t, dir, "remote", "add", "origin", origin)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-q", "-m", "initial")
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
