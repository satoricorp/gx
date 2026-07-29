package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Read-only commands run through ensureAutoInitializedRepo, which guards the
// user's staged Git index with PreservingGitIndex. Inspecting repo state must
// never disturb what the user has staged.
func TestReadOnlyCommandsPreserveGitIndex(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found")
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
	before := gitCachedDiff(t, root)
	if len(before) == 0 {
		t.Fatal("expected staged content before running tl")
	}

	totalityHome := t.TempDir()
	writeTestTotalityConfig(t, totalityHome, "Test", "t@e.com")
	t.Setenv("TOTALITY_HOME", totalityHome)
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
		{name: "doctor", args: []string{"doctor"}},
		{name: "doctor json", args: []string{"doctor", "--json"}},
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
				t.Fatalf("tl %s error = %v", strings.Join(mode.args, " "), err)
			}
			after := gitCachedDiff(t, root)
			if string(after) != string(before) {
				t.Fatalf("git index changed after tl %s\nbefore=%q\nafter=%q", strings.Join(mode.args, " "), before, after)
			}
		})
	}
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

func gitCachedDiff(t *testing.T, root string) []byte {
	t.Helper()
	cmd := exec.Command("git", "diff", "--cached", "--binary")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git diff --cached: %v\n%s", err, out)
	}
	return out
}
