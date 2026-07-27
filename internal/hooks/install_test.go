package hooks_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/hooks"
)

func TestInstallPrePushHook(t *testing.T) {
	repo := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = repo
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	if err := hooks.Install(hooks.InstallOptions{
		RepoRoot: repo,
		GXPath:   "/usr/local/bin/gx",
	}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repo, ".git", "hooks", "pre-push"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !hooks.IsInstalled(repo) {
		t.Fatal("expected hook installed")
	}
	if content == "" {
		t.Fatal("empty hook script")
	}
}

// TestPerRepoHooksResolveGXAtRunTime pins the run-time gx resolution of the
// per-repo scripts: the pinned path when it is still an executable file, else
// gx from PATH, else a silent exit 0 so a missing binary never breaks a
// commit or push.
func TestPerRepoHooksResolveGXAtRunTime(t *testing.T) {
	repo := t.TempDir()
	runGit(t, repo, "init")
	work := t.TempDir()
	pinnedLog := filepath.Join(work, "pinned.log")
	pathLog := filepath.Join(work, "path.log")
	pinnedGX := filepath.Join(work, "pinned", "gx")
	writeExecutable(t, pinnedGX, "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \""+pinnedLog+"\"\n")
	pathDir := filepath.Join(work, "pathbin")
	writeExecutable(t, filepath.Join(pathDir, "gx"), "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \""+pathLog+"\"\n")
	// The pre-push loop shells out to git; nothing else on PATH is needed.
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(gitPath, filepath.Join(pathDir, "git")); err != nil {
		t.Fatal(err)
	}
	stdin := "refs/heads/main aaaa111 refs/heads/main bbbb222\n"

	install := func(gxPath string) string {
		t.Helper()
		hookPath := filepath.Join(repo, ".git", "hooks", "pre-push")
		if err := os.RemoveAll(hookPath); err != nil {
			t.Fatal(err)
		}
		if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, GXPath: gxPath}); err != nil {
			t.Fatal(err)
		}
		return hookPath
	}
	runPrePush := func(hookPath, pathEnv string) (string, error) {
		t.Helper()
		cmd := exec.Command(hookPath, "origin", "git@example.com:acme/app.git")
		cmd.Dir = repo
		cmd.Stdin = strings.NewReader(stdin)
		cmd.Env = append(os.Environ(), "PATH="+pathEnv)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	t.Run("pinned path wins when executable", func(t *testing.T) {
		out, err := runPrePush(install(pinnedGX), pathDir)
		if err != nil {
			t.Fatalf("pre-push: %v\n%s", err, out)
		}
		if got := readFile(t, pinnedLog); !strings.Contains(got, "capture push") {
			t.Fatalf("pinned gx not invoked:\n%s", got)
		}
		if _, err := os.Stat(pathLog); err == nil {
			t.Fatal("PATH gx invoked despite a valid pinned path")
		}
	})

	t.Run("falls back to PATH when pinned path is gone", func(t *testing.T) {
		hookPath := install(filepath.Join(work, "missing", "gx"))
		out, err := runPrePush(hookPath, pathDir)
		if err != nil {
			t.Fatalf("pre-push: %v\n%s", err, out)
		}
		got := readFile(t, pathLog)
		for _, want := range []string{"capture push", "--ref-range bbbb222..aaaa111", "--local-ref refs/heads/main"} {
			if !strings.Contains(got, want) {
				t.Fatalf("PATH gx args missing %q:\n%s", want, got)
			}
		}
	})

	t.Run("exits zero silently when no gx exists", func(t *testing.T) {
		hookPath := install(filepath.Join(work, "missing", "gx"))
		emptyDir := filepath.Join(work, "empty")
		if err := os.MkdirAll(emptyDir, 0o755); err != nil {
			t.Fatal(err)
		}
		out, err := runPrePush(hookPath, emptyDir)
		if err != nil {
			t.Fatalf("pre-push without gx must exit 0, got %v\n%s", err, out)
		}
		if strings.TrimSpace(out) != "" {
			t.Fatalf("pre-push without gx must be silent, got:\n%s", out)
		}
	})

	t.Run("prepare-commit-msg exits zero without gx", func(t *testing.T) {
		install(filepath.Join(work, "missing", "gx"))
		hookPath := filepath.Join(repo, ".git", "hooks", "prepare-commit-msg")
		messagePath := filepath.Join(repo, ".git", "COMMIT_EDITMSG")
		if err := os.WriteFile(messagePath, []byte("subject\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		emptyDir := filepath.Join(work, "empty")
		if err := os.MkdirAll(emptyDir, 0o755); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command(hookPath, messagePath)
		cmd.Dir = repo
		cmd.Env = append(os.Environ(), "PATH="+emptyDir)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("prepare-commit-msg without gx must exit 0 (a failing hook aborts the commit), got %v\n%s", err, out)
		}
		if strings.TrimSpace(string(out)) != "" {
			t.Fatalf("prepare-commit-msg without gx must be silent, got:\n%s", out)
		}
	})
}

func TestRefRangeForPush(t *testing.T) {
	zero := "0000000000000000000000000000000000000000"
	if got := hooks.RefRangeForPush("abc", zero); got != "abc" {
		t.Fatalf("new branch range = %q", got)
	}
	if got := hooks.RefRangeForPush("def", "abc"); got != "abc..def" {
		t.Fatalf("update range = %q", got)
	}
}

func TestParsePrePushLine(t *testing.T) {
	localRef, localSHA, remoteRef, remoteSHA, err := hooks.ParsePrePushLine(
		"refs/heads/main abc refs/heads/main def",
	)
	if err != nil {
		t.Fatal(err)
	}
	if localRef != "refs/heads/main" || localSHA != "abc" || remoteRef != "refs/heads/main" || remoteSHA != "def" {
		t.Fatalf("unexpected parse: %q %q %q %q", localRef, localSHA, remoteRef, remoteSHA)
	}
}

func TestHuskySnippet(t *testing.T) {
	snippet := hooks.HuskyPrePushSnippet("/bin/gx")
	if snippet == "" || snippet[:6] != "remote" {
		t.Fatalf("unexpected snippet: %q", snippet)
	}
}
