package hooks

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	cmds := [][]string{
		{"git", "init", "-b", "main"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", "README.md")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, out)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v: %s", err, out)
	}
}

func TestRunPushRespectsPauseFlag(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()
	initGitRepo(t, repo)

	gxHome := filepath.Join(home, ".gx")
	t.Setenv("GX_HOME", gxHome)
	if err := os.MkdirAll(gxHome, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gxHome, "pause-capture"), []byte("1\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	outcome, err := RunPush(context.Background(), PushOptions{
		RepoRoot: repo,
		HomeDir:  home,
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if outcome.Result.RefRange != "" || outcome.Result.EligibleHunks != 0 {
		t.Fatalf("RunPush() = %+v, want empty result when paused", outcome)
	}
}

func TestParseRefRangeSingleSHALeavesBaseForOrchestrator(t *testing.T) {
	base, head := parseRefRange("abc123", "", "")
	if base != "" || head != "abc123" {
		t.Fatalf("parseRefRange() = (%q, %q), want empty base and abc123 head", base, head)
	}
}
