package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	calls      [][]string
	kickstarts int
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	call := append([]string{name}, args...)
	f.calls = append(f.calls, call)
	if name != "launchctl" {
		return nil, nil
	}
	switch args[0] {
	case "kickstart":
		f.kickstarts++
		if f.kickstarts == 1 {
			return []byte(`Could not find service "dev.gx.capture" in domain for user gui: 501`), errors.New("exit status 113")
		}
		return nil, nil
	case "bootout", "bootstrap", "enable", "setenv":
		return nil, nil
	default:
		return nil, nil
	}
}

func TestStartInstallsWhenPlistMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("GX_HOME", filepath.Join(home, ".gx"))

	runner := &successRunner{}
	mgr := Manager{Runner: runner}
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	agentPath := filepath.Join(home, "Library", "LaunchAgents", Label+".plist")
	if _, err := os.Stat(agentPath); err != nil {
		t.Fatalf("Start() did not write launch agent: %v", err)
	}
}

type successRunner struct{}

func (successRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	return nil, nil
}

func TestStartBootstrapsWhenServiceUnloaded(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	agentPath := filepath.Join(home, "Library", "LaunchAgents", Label+".plist")
	if err := os.MkdirAll(filepath.Dir(agentPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(agentPath, []byte(LaunchAgentPlist("/opt/homebrew/bin/gx")), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(agentPath) })

	runner := &fakeRunner{}
	mgr := Manager{Runner: runner}
	if err := mgr.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if runner.kickstarts != 2 {
		t.Fatalf("kickstart calls = %d, want 2", runner.kickstarts)
	}

	joined := make([]string, len(runner.calls))
	for i, call := range runner.calls {
		joined[i] = strings.Join(call, " ")
	}
	for _, want := range []string{
		"launchctl bootout " + guiTarget() + "/" + Label,
		"launchctl bootstrap " + guiTarget() + " " + agentPath,
		"launchctl enable " + guiTarget() + "/" + Label,
	} {
		if !containsCall(joined, want) {
			t.Fatalf("missing launchctl call %q in %#v", want, runner.calls)
		}
	}
}

func containsCall(calls []string, want string) bool {
	for _, call := range calls {
		if call == want {
			return true
		}
	}
	return false
}

func TestStopIsIdempotentWhenServiceMissing(t *testing.T) {
	runner := &bootoutMissingRunner{}
	mgr := Manager{Runner: runner}
	if err := mgr.Stop(context.Background()); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
}

type bootoutMissingRunner struct{}

func (bootoutMissingRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	if name == "launchctl" && len(args) > 0 && args[0] == "bootout" {
		return []byte("Boot-out failed: 3: No such process"), errors.New("exit status 3")
	}
	return nil, nil
}

func TestLaunchAgentPlistUsesAmbientDaemon(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	plist := LaunchAgentPlist("/opt/homebrew/bin/gx")
	for _, want := range []string{
		"<string>dev.gx.capture</string>",
		"<string>/opt/homebrew/bin/gx</string>",
		"<string>__gx-daemon</string>",
		"<string>--ambient</string>",
		"<key>KeepAlive</key>",
	} {
		if !strings.Contains(plist, want) {
			t.Fatalf("plist missing %q:\n%s", want, plist)
		}
	}
}

func TestShellBlockUsesStableProxyURLs(t *testing.T) {
	t.Setenv("GX_PROXY_ADDR", "127.0.0.1:49999")
	block := ShellBlock()
	for _, want := range []string{
		"export ANTHROPIC_BASE_URL='http://127.0.0.1:49999'",
		"export OPENAI_BASE_URL='http://127.0.0.1:49999/v1'",
	} {
		if !strings.Contains(block, want) {
			t.Fatalf("shell block missing %q:\n%s", want, block)
		}
	}
}

func TestControlAddressDefaultsAndOverride(t *testing.T) {
	t.Setenv("GX_DAEMON_ADDR", "")
	if got := ControlAddress(); got != DefaultControl {
		t.Fatalf("ControlAddress() = %q, want %q", got, DefaultControl)
	}
	t.Setenv("GX_DAEMON_ADDR", "127.0.0.1:49998")
	if got := ControlAddress(); got != "127.0.0.1:49998" {
		t.Fatalf("ControlAddress() override = %q", got)
	}
	if got := ControlURL(); got != "http://127.0.0.1:49998" {
		t.Fatalf("ControlURL() = %q", got)
	}
}

func TestUpsertBlockIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".zshrc")
	if err := os.WriteFile(path, []byte("export PATH=/bin\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := upsertBlock(path, ShellBlock()); err != nil {
		t.Fatal(err)
	}
	if err := upsertBlock(path, ShellBlock()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(string(data), shellBlockOpen); got != 1 {
		t.Fatalf("gx shell block count = %d, want 1:\n%s", got, string(data))
	}
	if !strings.Contains(string(data), "export PATH=/bin") {
		t.Fatalf("preserved shell content missing:\n%s", string(data))
	}
}
