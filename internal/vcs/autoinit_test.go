package vcs

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureReadyRepoSkipsOutsideRepository(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	root := t.TempDir()
	svc := NewServiceWithRunner(&fakeRunner{})

	result, err := svc.EnsureReadyRepo(context.Background(), root)
	if err != nil {
		t.Fatalf("EnsureReadyRepo() error = %v", err)
	}
	if result.Prepared {
		t.Fatal("EnsureReadyRepo() prepared work outside a repository")
	}
}

func TestEnsureReadyRepoNoOpWhenInitializedAndBaseExists(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey("/repo", "git", "rev-parse", "--show-toplevel"):                              {"/repo\n"},
			runnerKey("/repo", "jj", "root"):                                                        {"/repo\n"},
			runnerKey("/repo", "git", "branch", "--show-current"):                                   {"main\n"},
			runnerKey("/repo", "jj", "log", "-r", "main", "--limit", "1", "--no-graph", "-T", "commit_id"): {"abc123\n"},
		},
	}
	svc := NewServiceWithRunner(runner)
	store, err := openStore(context.Background())
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	defer store.Close()
	if err := store.RecordInitializedRepo(context.Background(), "/repo", 1); err != nil {
		t.Fatalf("RecordInitializedRepo() error = %v", err)
	}

	result, err := svc.EnsureReadyRepo(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("EnsureReadyRepo() error = %v", err)
	}
	if result.Prepared {
		t.Fatal("EnsureReadyRepo() prepared work for an already-ready repo")
	}
	for _, call := range runner.calls {
		if strings.Contains(call, "jj|git init") {
			t.Fatalf("EnsureReadyRepo() unexpectedly initialized repo: %v", runner.calls)
		}
	}
}

func TestEnsureReadyRepoBootstrapsRegisteredRepoWithoutBase(t *testing.T) {
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj executable not found")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found")
	}

	root := t.TempDir()
	runGit(t, root, "init", "-b", "main")
	runGit(t, root, "config", "user.name", "Joe Example")
	runGit(t, root, "config", "user.email", "joe@example.com")
	runJJ(t, root, "git", "init", ".")

	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("HOME", t.TempDir())
	store, err := openStore(context.Background())
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	if err := store.RecordInitializedRepo(context.Background(), root, 1); err != nil {
		t.Fatalf("RecordInitializedRepo() error = %v", err)
	}
	store.Close()

	svc := NewServiceWithRunner(ExecRunner{})
	result, err := svc.EnsureReadyRepo(context.Background(), root)
	if err != nil {
		t.Fatalf("EnsureReadyRepo() error = %v", err)
	}
	if !result.Prepared {
		t.Fatal("EnsureReadyRepo() did not prepare registered repo missing main")
	}
	exists, err := svc.RevisionExists(context.Background(), root, "main")
	if err != nil {
		t.Fatalf("RevisionExists(main) error = %v", err)
	}
	if !exists {
		t.Fatal("RevisionExists(main) = false after EnsureReadyRepo")
	}
	if _, err := os.Stat(filepath.Join(root, bootstrapReadmeName)); err != nil {
		t.Fatalf("README.md missing after EnsureReadyRepo: %v", err)
	}
}
