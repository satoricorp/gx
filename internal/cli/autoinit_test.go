package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorAutoInitializesUnbornGitRepo(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Joe Example")
	runGitTest(t, root, "config", "user.email", "joe@example.com")
	writeTestFile(t, root, "README.md", "# repo\n")
	t.Chdir(root)
	totalityHome := t.TempDir()
	writeTestTotalityConfig(t, totalityHome, "Joe Example", "joe@example.com")
	t.Setenv("TOTALITY_HOME", totalityHome)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "dumb")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	cmd := NewRoot(context.Background())
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs([]string{"doctor"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("tl doctor error = %v\nstdout:\n%s\nstderr:\n%s", err, out.String(), errOut.String())
	}
	combined := out.String() + errOut.String()
	if strings.Contains(combined, "Revision `main` doesn't exist") {
		t.Fatalf("tl doctor leaked jj main error:\n%s", combined)
	}
	if strings.Contains(combined, "jj log -r mutable()") {
		t.Fatalf("tl doctor leaked raw jj command:\n%s", combined)
	}
	if !strings.Contains(errOut.String(), "Initializing tl for this repository") {
		t.Fatalf("tl doctor missing auto-init notice on stderr:\nstdout:\n%s\nstderr:\n%s", out.String(), errOut.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".jj")); !os.IsNotExist(err) {
		t.Fatalf("tl doctor unexpectedly created .jj: %v", err)
	}
}

func TestDoctorWarnsWhenLifecycleHooksCannotBeInstalled(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Joe Example")
	runGitTest(t, root, "config", "user.email", "joe@example.com")
	hooksDir := filepath.Join(root, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooksDir, "pre-push"), []byte("#!/bin/sh\necho foreign\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	totalityHome := t.TempDir()
	writeTestTotalityConfig(t, totalityHome, "Joe Example", "joe@example.com")
	t.Setenv("TOTALITY_HOME", totalityHome)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "dumb")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	cmd := NewRoot(context.Background())
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs([]string{"doctor"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("tl doctor should remain usable: %v\n%s", err, errOut.String())
	}
	if !strings.Contains(errOut.String(), "warning: Totality lifecycle hooks not installed:") {
		t.Fatalf("missing lifecycle hook warning:\n%s", errOut.String())
	}
}

func TestDoctorHandlesRegisteredRepoWithoutMain(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Joe Example")
	runGitTest(t, root, "config", "user.email", "joe@example.com")
	writeTestFile(t, root, "README.md", "# repo\n")

	totalityHome := t.TempDir()
	writeTestTotalityConfig(t, totalityHome, "Joe Example", "joe@example.com")
	t.Setenv("TOTALITY_HOME", totalityHome)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TERM", "dumb")
	t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	store := openTestStore(t, context.Background())
	if err := store.RecordInitializedRepo(context.Background(), root, 1); err != nil {
		t.Fatalf("RecordInitializedRepo() error = %v", err)
	}
	store.Close()

	t.Chdir(root)
	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"doctor"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("tl doctor error = %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "doesn't exist") {
		t.Fatalf("tl doctor reported a missing base revision:\n%s", out.String())
	}
}

func TestShouldSkipAutoInitForSetupCommands(t *testing.T) {
	root := NewRoot(context.Background())
	cases := map[string]bool{
		"tl init":        true,
		"tl version":     true,
		"tl auth status": true,
		"tl review":      true,
		"tl doctor":      false,
	}
	for path, want := range cases {
		cmd, _, err := root.Find(strings.Fields(strings.TrimPrefix(path, "tl ")))
		if err != nil {
			t.Fatalf("Find(%q) error = %v", path, err)
		}
		if got := shouldSkipAutoInit(cmd); got != want {
			t.Fatalf("shouldSkipAutoInit(%q) = %v, want %v", path, got, want)
		}
	}
}
