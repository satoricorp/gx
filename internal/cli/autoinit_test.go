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
	gxHome := t.TempDir()
	writeTestGXConfig(t, gxHome, "Joe Example", "joe@example.com")
	t.Setenv("GX_HOME", gxHome)
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
		t.Fatalf("gx doctor error = %v\nstdout:\n%s\nstderr:\n%s", err, out.String(), errOut.String())
	}
	combined := out.String() + errOut.String()
	if strings.Contains(combined, "Revision `main` doesn't exist") {
		t.Fatalf("gx doctor leaked jj main error:\n%s", combined)
	}
	if strings.Contains(combined, "jj log -r mutable()") {
		t.Fatalf("gx doctor leaked raw jj command:\n%s", combined)
	}
	if !strings.Contains(errOut.String(), "Initializing gx for this repository") {
		t.Fatalf("gx doctor missing auto-init notice on stderr:\nstdout:\n%s\nstderr:\n%s", out.String(), errOut.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".jj")); !os.IsNotExist(err) {
		t.Fatalf("gx doctor unexpectedly created .jj: %v", err)
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
	gxHome := t.TempDir()
	writeTestGXConfig(t, gxHome, "Joe Example", "joe@example.com")
	t.Setenv("GX_HOME", gxHome)
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
		t.Fatalf("gx doctor should remain usable: %v\n%s", err, errOut.String())
	}
	if !strings.Contains(errOut.String(), "warning: GX lifecycle hooks not installed:") {
		t.Fatalf("missing lifecycle hook warning:\n%s", errOut.String())
	}
}

func TestDoctorHandlesRegisteredRepoWithoutMain(t *testing.T) {
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Joe Example")
	runGitTest(t, root, "config", "user.email", "joe@example.com")
	writeTestFile(t, root, "README.md", "# repo\n")

	gxHome := t.TempDir()
	writeTestGXConfig(t, gxHome, "Joe Example", "joe@example.com")
	t.Setenv("GX_HOME", gxHome)
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
		t.Fatalf("gx doctor error = %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "doesn't exist") {
		t.Fatalf("gx doctor reported a missing base revision:\n%s", out.String())
	}
}

func TestShouldSkipAutoInitForSetupCommands(t *testing.T) {
	root := NewRoot(context.Background())
	cases := map[string]bool{
		"gx init":                    true,
		"gx version":                 true,
		"gx auth status":             true,
		"gx set inference-key dummy": true,
		"gx review":                  true,
		"gx doctor":                  false,
	}
	for path, want := range cases {
		cmd, _, err := root.Find(strings.Fields(strings.TrimPrefix(path, "gx ")))
		if err != nil {
			t.Fatalf("Find(%q) error = %v", path, err)
		}
		if got := shouldSkipAutoInit(cmd); got != want {
			t.Fatalf("shouldSkipAutoInit(%q) = %v, want %v", path, got, want)
		}
	}
}
