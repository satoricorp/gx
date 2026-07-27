package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// `gx review` has to be usable as a CI gate and on a checkout the reviewer
// does not own, so it must not auto-initialize anything: no git hooks, no GX
// home, no repo state. This test reviews a repo that has never run `gx init`
// and asserts the repo and the machine come out untouched.
func TestReviewNeverInitializesTheRepoOrTheMachine(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	// Point GX_HOME and HOME at paths that do not exist yet: anything that
	// opens the store or writes config has to create them, which makes the
	// mutation visible instead of silently landing in an existing directory.
	sandbox := t.TempDir()
	gxHome := filepath.Join(sandbox, "gx-home")
	fakeHome := filepath.Join(sandbox, "home")
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("HOME", fakeHome)

	hooksDir := filepath.Join(root, ".git", "hooks")
	before := hookDirEntries(t, hooksDir)

	out, err := runReviewCommand(t, "--base", "main")
	if err != nil {
		t.Fatalf("gx review error = %v\n%s", err, out)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatalf("gx review produced no report in an uninitialized repo")
	}

	if after := hookDirEntries(t, hooksDir); !equalStrings(before, after) {
		t.Fatalf("gx review changed .git/hooks:\nbefore: %v\nafter:  %v", before, after)
	}
	for _, name := range before {
		data, readErr := os.ReadFile(filepath.Join(hooksDir, name))
		if readErr != nil {
			t.Fatalf("read hook %s: %v", name, readErr)
		}
		if strings.Contains(string(data), "gx ") || strings.Contains(string(data), "GX_") {
			t.Fatalf("gx review left a GX marker in .git/hooks/%s:\n%s", name, data)
		}
	}
	for _, path := range []string{gxHome, fakeHome, filepath.Join(root, ".gx")} {
		if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
			t.Fatalf("gx review created %s (stat error = %v), want it untouched", path, statErr)
		}
	}
}

func hookDirEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("read %s: %v", dir, err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
