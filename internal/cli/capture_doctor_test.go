package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/hooks"
	"github.com/satoricorp/gx/internal/storage"
)

func TestCaptureHookStatusNotApplicableOutsideGitRepo(t *testing.T) {
	installed, applicable, resolvedRoot := captureHookStatus(context.Background(), t.TempDir())
	if installed {
		t.Fatal("expected hook not installed outside a git repo")
	}
	if applicable {
		t.Fatal("expected hook check not applicable outside a git repo")
	}
	if resolvedRoot != "" {
		t.Fatalf("resolvedRoot = %q, want empty", resolvedRoot)
	}
}

func TestCaptureHookStatusResolvesGitRoot(t *testing.T) {
	repo := initDoctorGitRepo(t)
	wantRepo, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	subdir := filepath.Join(repo, "nested", "dir")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}

	installed, applicable, resolvedRoot := captureHookStatus(context.Background(), subdir)
	if installed {
		t.Fatal("expected hook not installed before gx init")
	}
	if !applicable {
		t.Fatal("expected hook check applicable inside a git repo")
	}
	if resolvedRoot != wantRepo {
		t.Fatalf("resolvedRoot = %q, want %q", resolvedRoot, wantRepo)
	}

	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, GXPath: "/usr/local/bin/gx"}); err != nil {
		t.Fatal(err)
	}
	installed, applicable, resolvedRoot = captureHookStatus(context.Background(), subdir)
	if !installed || !applicable || resolvedRoot != wantRepo {
		t.Fatalf("captureHookStatus() = installed %v applicable %v root %q, want true true %q", installed, applicable, resolvedRoot, wantRepo)
	}
}

func TestCaptureRegisteredRepoHooksReportsInstalledAndMissingHooks(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	installedRepo := initDoctorGitRepo(t)
	missingRepo := initDoctorGitRepo(t)
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: installedRepo, GXPath: "/usr/local/bin/gx"}); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, ctx)
	defer store.Close()
	for _, repoRoot := range []string{installedRepo, missingRepo} {
		if err := store.RecordInitializedRepo(ctx, repoRoot, 1); err != nil {
			t.Fatalf("RecordInitializedRepo(%s) error = %v", repoRoot, err)
		}
	}

	got := captureRegisteredRepoHooks(ctx)
	if len(got) != 2 {
		t.Fatalf("captureRegisteredRepoHooks() returned %d repos, want 2: %#v", len(got), got)
	}
	var installed, missing int
	for _, repo := range got {
		if !repo.GitReachable {
			t.Fatalf("repo should be reachable: %#v", repo)
		}
		if repo.HookInstalled {
			installed++
		} else {
			missing++
		}
	}
	if installed != 1 || missing != 1 {
		t.Fatalf("registered hooks installed=%d missing=%d, want 1 and 1: %#v", installed, missing, got)
	}
}

func initDoctorGitRepo(t *testing.T) string {
	t.Helper()
	requireGit(t)
	repo := t.TempDir()
	runGit(t, repo, "init")
	return repo
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func openTestStore(t *testing.T, ctx context.Context) *storage.Store {
	t.Helper()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		_ = db.Close()
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	return store
}
