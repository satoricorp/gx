package vcs

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureAuthoringBaseRevisionNoOpWhenBaseExists(t *testing.T) {
	cwd := t.TempDir()
	repo := RepoInfo{RootPath: cwd}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "log", "-r", "main", "--limit", "1", "--no-graph", "-T", "commit_id"): {"abc123\n"},
		},
	}
	svc := NewServiceWithRunner(runner)

	if err := svc.ensureAuthoringBaseRevision(context.Background(), repo, "main"); err != nil {
		t.Fatalf("ensureAuthoringBaseRevision() error = %v, want nil", err)
	}
	assertRunnerNotCalled(t, runner.calls, runnerKey(cwd, "jj", "commit", "-m", bootstrapInitCommitMessage, bootstrapReadmeName))
}

func TestEnsureBootstrapReadmeCreatesDefault(t *testing.T) {
	root := t.TempDir()
	svc := NewServiceWithRunner(&fakeRunner{})

	if err := svc.ensureBootstrapReadme(root); err != nil {
		t.Fatalf("ensureBootstrapReadme() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, bootstrapReadmeName))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	baseName := filepath.Base(root)
	if got, want := string(data), "# "+baseName+"\n"; got != want {
		t.Fatalf("README content = %q, want %q", got, want)
	}
}

func TestEnsureBootstrapReadmePreservesExisting(t *testing.T) {
	root := t.TempDir()
	readmePath := filepath.Join(root, bootstrapReadmeName)
	if err := os.WriteFile(readmePath, []byte("# custom\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	svc := NewServiceWithRunner(&fakeRunner{})

	if err := svc.ensureBootstrapReadme(root); err != nil {
		t.Fatalf("ensureBootstrapReadme() error = %v", err)
	}
	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "# custom\n" {
		t.Fatalf("README content = %q, want preserved custom content", string(data))
	}
}

func TestEnsureAuthoringBaseRevisionBootstrapsWithUntrackedWork(t *testing.T) {
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
	if err := os.MkdirAll(filepath.Join(root, "internal"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "foo.go"), []byte("package foo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runJJ(t, root, "git", "init", "--colocate", ".")

	svc := NewServiceWithRunner(ExecRunner{})
	main := "main"
	repo := RepoInfo{RootPath: root, BranchName: &main}
	if err := svc.ensureAuthoringBaseRevision(context.Background(), repo, "main"); err != nil {
		t.Fatalf("ensureAuthoringBaseRevision() error = %v", err)
	}
	readmePath := filepath.Join(root, bootstrapReadmeName)
	data, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	baseName := filepath.Base(root)
	if got, want := string(data), "# "+baseName+"\n"; got != want {
		t.Fatalf("README content = %q, want %q", got, want)
	}
	gitLogCmd := exec.Command("git", "log", "-1", "--format=%s")
	gitLogCmd.Dir = root
	out, err := gitLogCmd.CombinedOutput()
	if err != nil {
		allCmd := exec.Command("git", "log", "--oneline", "--all")
		allCmd.Dir = root
		allOut, _ := allCmd.CombinedOutput()
		t.Fatalf("git log after bootstrap: %v\n%s\nall:\n%s", err, out, allOut)
	}
	if strings.TrimSpace(string(out)) != bootstrapInitCommitMessage {
		t.Fatalf("initial commit message = %q, want %q", strings.TrimSpace(string(out)), bootstrapInitCommitMessage)
	}
	statCmd := exec.Command("jj", "log", "-r", "main", "--stat")
	statCmd.Dir = root
	statOut, err := statCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("jj log --stat: %v\n%s", err, statOut)
	}
	if !strings.Contains(string(statOut), "README.md") {
		t.Fatalf("init commit missing README.md:\n%s", statOut)
	}
}

func TestEnsureAuthoringCheckoutBootstrapsUnbornBaseWithWork(t *testing.T) {
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
	if err := os.MkdirAll(filepath.Join(root, "internal"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "foo.go"), []byte("package foo\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runJJ(t, root, "git", "init", "--colocate", ".")

	svc := NewServiceWithRunner(ExecRunner{})
	repo := RepoInfo{RootPath: root, BranchName: strPtr("main")}
	ref, mode, err := svc.ensureAuthoringCheckout(context.Background(), repo, "gx generate")
	if err != nil {
		t.Fatalf("ensureAuthoringCheckout() error = %v", err)
	}
	if ref != "main" || mode != "base" {
		t.Fatalf("ensureAuthoringCheckout() = (%q, %q), want (main, base)", ref, mode)
	}
	exists, err := svc.RevisionExists(context.Background(), root, "main")
	if err != nil {
		t.Fatalf("RevisionExists() error = %v", err)
	}
	if !exists {
		t.Fatalf("main bookmark missing after ensureAuthoringCheckout")
	}
}

func TestEnsureAuthoringCheckoutBootstrapsUnbornBase(t *testing.T) {
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
	runJJ(t, root, "git", "init", "--colocate", ".")

	svc := NewServiceWithRunner(ExecRunner{})
	repo := RepoInfo{RootPath: root, BranchName: strPtr("main")}
	ref, mode, err := svc.ensureAuthoringCheckout(context.Background(), repo, "gx generate")
	if err != nil {
		t.Fatalf("ensureAuthoringCheckout() error = %v", err)
	}
	if ref != "main" || mode != "base" {
		t.Fatalf("ensureAuthoringCheckout() = (%q, %q), want (main, base)", ref, mode)
	}
}


func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func runJJ(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("jj", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("jj %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func strPtr(value string) *string {
	return &value
}
