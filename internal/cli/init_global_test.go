package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// isolateGlobalGitConfigForCLI points HOME, TOTALITY_HOME and git's global/system
// config at throwaway files so `tx init --global` can be executed in tests
// without ever touching the developer's real ~/.gitconfig or ~/.totality.
func isolateGlobalGitConfigForCLI(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	configPath := filepath.Join(home, "gitconfig")
	if err := os.WriteFile(configPath, []byte("[tx]\n\tisolationprobe = temp\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_HOME", filepath.Join(home, "tx"))
	t.Setenv("GIT_CONFIG_GLOBAL", configPath)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("NO_COLOR", "1")
	out, err := exec.Command("git", "config", "--global", "--get", "tx.isolationprobe").Output()
	if err != nil || strings.TrimSpace(string(out)) != "temp" {
		t.Skipf("git does not honor GIT_CONFIG_GLOBAL, refusing to touch the real global config: %v", err)
	}
}

func TestInitGlobalInstallsMachineWideHooks(t *testing.T) {
	isolateGlobalGitConfigForCLI(t)

	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"init", "--global"})
	if err := root.Execute(); err != nil {
		t.Fatalf("tx init --global error = %v\n%s", err, out.String())
	}

	hooksDir := filepath.Join(os.Getenv("TOTALITY_HOME"), "hooks")
	for _, name := range []string{"prepare-commit-msg", "post-commit", "post-rewrite", "pre-push", "pre-commit"} {
		data, err := os.ReadFile(filepath.Join(hooksDir, name))
		if err != nil {
			t.Fatalf("missing global hook %s: %v", name, err)
		}
		if !strings.Contains(string(data), "tx_resolve_local_hook "+name) {
			t.Fatalf("global %s hook does not chain to the repo hook:\n%s", name, data)
		}
	}

	configured, err := exec.Command("git", "config", "--global", "--get", "core.hooksPath").Output()
	if err != nil {
		t.Fatalf("core.hooksPath not set: %v", err)
	}
	if got := strings.TrimSpace(string(configured)); got != hooksDir {
		t.Fatalf("core.hooksPath = %q, want %q", got, hooksDir)
	}

	text := out.String()
	for _, want := range []string{
		"Global hooks",
		hooksDir,
		"git config tx.enabled false",
		"git config --global --unset core.hooksPath",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("tx init --global output missing %q:\n%s", want, text)
		}
	}
}

func TestInitGlobalRefusesForeignHooksPath(t *testing.T) {
	isolateGlobalGitConfigForCLI(t)
	foreign := t.TempDir()
	if out, err := exec.Command("git", "config", "--global", "core.hooksPath", foreign).CombinedOutput(); err != nil {
		t.Fatalf("seed core.hooksPath: %v\n%s", err, out)
	}

	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"init", "--global"})
	err := root.Execute()
	if err == nil {
		t.Fatalf("tx init --global error = nil, want refusal\n%s", out.String())
	}
	if !strings.Contains(err.Error(), "not managed by Totality") || !strings.Contains(err.Error(), foreign) {
		t.Fatalf("error = %q, want a refusal naming %s", err, foreign)
	}
	if got, _ := exec.Command("git", "config", "--global", "--get", "core.hooksPath").Output(); strings.TrimSpace(string(got)) != foreign {
		t.Fatalf("core.hooksPath = %q, want it left at %q", strings.TrimSpace(string(got)), foreign)
	}
}

func TestInitGlobalFlagIsOptIn(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("TOTALITY_HOME", t.TempDir())
	root := NewRoot(context.Background())
	init, _, err := root.Find([]string{"init"})
	if err != nil {
		t.Fatal(err)
	}
	flag := init.Flags().Lookup("global")
	if flag == nil {
		t.Fatal("tx init has no --global flag")
	}
	if flag.DefValue != "false" {
		t.Fatalf("--global default = %q, want false", flag.DefValue)
	}
	for _, name := range []string{"name", "email", "yes"} {
		if init.Flags().Lookup(name) == nil {
			t.Fatalf("tx init lost its --%s flag", name)
		}
	}
}

func TestCaptureGlobalHooksDoctorValue(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	for _, tc := range []struct {
		name  string
		state globalHookDoctorJSON
		want  string
	}{
		{name: "not opted in", state: globalHookDoctorJSON{Dir: "/home/dev/.totality/hooks"}, want: ""},
		{
			name:  "enabled",
			state: globalHookDoctorJSON{Dir: "/home/dev/.totality/hooks", Enabled: true, Installed: true},
			want:  "ok: /home/dev/.totality/hooks",
		},
		{
			name:  "overridden",
			state: globalHookDoctorJSON{Dir: "/home/dev/.totality/hooks", Installed: true, Conflict: "/home/dev/.husky", NeedsRepair: true},
			want:  "warn: core.hooksPath points at /home/dev/.husky; run `tx init --global` to restore",
		},
		{
			name:  "incomplete",
			state: globalHookDoctorJSON{Dir: "/home/dev/.totality/hooks", Enabled: true, NeedsRepair: true},
			want:  "warn: incomplete; run `tx init --global`",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := captureGlobalHooksDoctorValue(tc.state); got != tc.want {
				t.Fatalf("captureGlobalHooksDoctorValue() = %q, want %q", got, tc.want)
			}
		})
	}
}
