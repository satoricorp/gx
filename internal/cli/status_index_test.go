package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatusPreservesGitIndexAcrossModes(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found")
	}
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj not found")
	}
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Test")
	runGitTest(t, root, "config", "user.email", "t@e.com")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGitTest(t, root, "add", "README.md")
	runGitTest(t, root, "commit", "-m", "init")
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGitTest(t, root, "add", "a.txt")
	before := statusGitCachedDiff(t, root)
	if len(before) == 0 {
		t.Fatal("expected staged content before gx status")
	}

	jjCmd := exec.Command("jj", "git", "init", ".")
	jjCmd.Dir = root
	if out, err := jjCmd.CombinedOutput(); err != nil {
		t.Fatalf("jj git init error = %v\n%s", err, out)
	}

	gxHome := t.TempDir()
	writeTestGXConfig(t, gxHome, "Test", "t@e.com")
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "dumb")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Chdir(root)

	modes := []struct {
		name string
		args []string
	}{
		{name: "default", args: []string{"status"}},
		{name: "json", args: []string{"status", "--json"}},
		{name: "agent", args: []string{"status", "--agent"}},
	}
	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			cmd := NewRoot(context.Background())
			cmd.SetOut(discardWriter{})
			cmd.SetErr(discardWriter{})
			cmd.SetArgs(mode.args)
			if err := cmd.Execute(); err != nil {
				if strings.Contains(err.Error(), "Missing parent base") {
					t.Skip("missing stack base in test repo")
				}
				t.Fatalf("status mode %q error = %v", mode.name, err)
			}
			after := statusGitCachedDiff(t, root)
			if string(after) != string(before) {
				t.Fatalf("git index changed after gx status (%s)\nbefore=%q\nafter=%q", mode.name, before, after)
			}
		})
	}
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

func statusGitCachedDiff(t *testing.T, root string) []byte {
	t.Helper()
	cmd := exec.Command("git", "diff", "--cached", "--binary")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git diff --cached: %v\n%s", err, out)
	}
	return out
}
