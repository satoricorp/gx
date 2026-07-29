package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestLegacyLaunchAgentPathUsesHomeLaunchAgents(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := legacyLaunchAgentPath()
	if err != nil {
		t.Fatalf("legacyLaunchAgentPath() error = %v", err)
	}
	want := filepath.Join(home, "Library", "LaunchAgents", "dev.totality.capture.plist")
	if got != want {
		t.Fatalf("legacyLaunchAgentPath() = %q, want %q", got, want)
	}
}

func TestLegacyLaunchAgentTargetUsesGUIDomain(t *testing.T) {
	want := "gui/" + strconv.Itoa(os.Getuid()) + "/dev.totality.capture"
	if got := legacyLaunchAgentTarget(); got != want {
		t.Fatalf("legacyLaunchAgentTarget() = %q, want %q", got, want)
	}
}

func TestLegacyProxyBaseURLs(t *testing.T) {
	t.Setenv("TOTALITY_PROXY_ADDR", "")
	anthropicURL, openaiURL := legacyProxyBaseURLs()
	if anthropicURL != "http://127.0.0.1:43123" || openaiURL != "http://127.0.0.1:43123/v1" {
		t.Fatalf("legacyProxyBaseURLs() = %q, %q", anthropicURL, openaiURL)
	}

	t.Setenv("TOTALITY_PROXY_ADDR", "127.0.0.1:50000")
	anthropicURL, openaiURL = legacyProxyBaseURLs()
	if anthropicURL != "http://127.0.0.1:50000" || openaiURL != "http://127.0.0.1:50000/v1" {
		t.Fatalf("legacyProxyBaseURLs() with TOTALITY_PROXY_ADDR = %q, %q", anthropicURL, openaiURL)
	}
}

func TestStripLegacyShellBlock(t *testing.T) {
	block := strings.Join([]string{
		legacyShellBlockOpen,
		"export ANTHROPIC_BASE_URL='http://127.0.0.1:43123'",
		"export OPENAI_BASE_URL='http://127.0.0.1:43123/v1'",
		legacyShellBlockClose,
		"",
	}, "\n")

	stripped, removed := stripLegacyShellBlock("# mine\nalias ll='ls -la'\n" + block)
	if !removed {
		t.Fatal("stripLegacyShellBlock() removed = false, want true")
	}
	if stripped != "# mine\nalias ll='ls -la'\n" {
		t.Fatalf("stripLegacyShellBlock() = %q", stripped)
	}

	unchanged, removed := stripLegacyShellBlock("# mine\nalias ll='ls -la'\n")
	if removed || unchanged != "# mine\nalias ll='ls -la'\n" {
		t.Fatalf("stripLegacyShellBlock() without markers = %q, removed=%v", unchanged, removed)
	}

	// A dangling open marker without a close marker must be left untouched.
	dangling := "# mine\n" + legacyShellBlockOpen + "\nexport X=1\n"
	unchanged, removed = stripLegacyShellBlock(dangling)
	if removed || unchanged != dangling {
		t.Fatalf("stripLegacyShellBlock() with dangling marker = %q, removed=%v", unchanged, removed)
	}
}

type fakeLegacyRunner struct {
	calls  []string
	getenv map[string]string
}

func (f *fakeLegacyRunner) run(_ context.Context, name string, args ...string) ([]byte, error) {
	call := name + " " + strings.Join(args, " ")
	f.calls = append(f.calls, call)
	if name == "launchctl" && len(args) == 2 && args[0] == "getenv" {
		value, ok := f.getenv[args[1]]
		if !ok {
			return nil, fmt.Errorf("unknown env %s", args[1])
		}
		return []byte(value + "\n"), nil
	}
	return nil, nil
}

func writeLegacyPlist(t *testing.T, home string) string {
	t.Helper()
	dir := filepath.Join(home, "Library", "LaunchAgents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "dev.totality.capture.plist")
	if err := os.WriteFile(path, []byte("<plist/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCleanupLegacyAmbientCaptureRemovesArtifacts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "")

	plistPath := writeLegacyPlist(t, home)
	zshrc := filepath.Join(home, ".zshrc")
	rcBody := strings.Join([]string{
		"# mine",
		legacyShellBlockOpen,
		"export ANTHROPIC_BASE_URL='http://127.0.0.1:43123'",
		legacyShellBlockClose,
		"",
	}, "\n")
	if err := os.WriteFile(zshrc, []byte(rcBody), 0o644); err != nil {
		t.Fatal(err)
	}

	runner := &fakeLegacyRunner{getenv: map[string]string{
		"ANTHROPIC_BASE_URL": "http://127.0.0.1:43123",
		"OPENAI_BASE_URL":    "http://127.0.0.1:43123/v1",
	}}
	result := cleanupLegacyAmbientCapture(context.Background(), runner.run)

	if !result.CleanedLaunchAgent || !result.CleanedShellBlocks || !result.CleanedEnv {
		t.Fatalf("cleanup result = %+v, want all cleaned", result)
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("cleanup warnings = %v, want none", result.Warnings)
	}
	if _, err := os.Stat(plistPath); !os.IsNotExist(err) {
		t.Fatalf("plist still present after cleanup: %v", err)
	}
	data, err := os.ReadFile(zshrc)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "tl ambient capture") {
		t.Fatalf(".zshrc still contains managed block:\n%s", data)
	}
	if !strings.Contains(string(data), "# mine") {
		t.Fatalf(".zshrc lost user content:\n%s", data)
	}

	wantCalls := []string{
		"launchctl bootout " + legacyLaunchAgentTarget(),
		"launchctl getenv ANTHROPIC_BASE_URL",
		"launchctl unsetenv ANTHROPIC_BASE_URL",
		"launchctl getenv OPENAI_BASE_URL",
		"launchctl unsetenv OPENAI_BASE_URL",
	}
	if strings.Join(runner.calls, "; ") != strings.Join(wantCalls, "; ") {
		t.Fatalf("launchctl calls = %v, want %v", runner.calls, wantCalls)
	}
}

func TestCleanupLegacyAmbientCaptureNoopWhenAbsent(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	runner := &fakeLegacyRunner{}
	result := cleanupLegacyAmbientCapture(context.Background(), runner.run)

	if result.cleanedAny() {
		t.Fatalf("cleanup result = %+v, want untouched", result)
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("cleanup warnings = %v, want none", result.Warnings)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("launchctl calls = %v, want none", runner.calls)
	}
}

func TestCleanupLegacyAmbientCaptureKeepsForeignEnvValues(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "")
	writeLegacyPlist(t, home)

	runner := &fakeLegacyRunner{getenv: map[string]string{
		"ANTHROPIC_BASE_URL": "https://my-own-gateway.example.com",
		"OPENAI_BASE_URL":    "http://127.0.0.1:43123/v1",
	}}
	result := cleanupLegacyAmbientCapture(context.Background(), runner.run)

	if !result.CleanedLaunchAgent {
		t.Fatalf("cleanup result = %+v, want LaunchAgent cleaned", result)
	}
	for _, call := range runner.calls {
		if call == "launchctl unsetenv ANTHROPIC_BASE_URL" {
			t.Fatalf("cleanup unset a user-owned env value: %v", runner.calls)
		}
	}
	found := false
	for _, call := range runner.calls {
		if call == "launchctl unsetenv OPENAI_BASE_URL" {
			found = true
		}
	}
	if !found {
		t.Fatalf("cleanup left the legacy proxy env value in place: %v", runner.calls)
	}
}

func writeLegacyCodexConfig(t *testing.T, home, body string) string {
	t.Helper()
	dir := filepath.Join(home, ".codex")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func legacyCodexBackups(t *testing.T, path string) []string {
	t.Helper()
	matches, err := filepath.Glob(path + ".totality-backup-*")
	if err != nil {
		t.Fatal(err)
	}
	return matches
}

func TestRevertLegacyCodexConfigFileFullRevert(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "")

	// Exactly what the retired service writer produced when repairing a
	// pre-existing config: the revert must round-trip back to that original.
	fixture := strings.Join([]string{
		`# personal settings`,
		`model = "gpt-5"`,
		`model_provider = "totality-openai"`,
		``,
		`[projects."/repo"]`,
		`trust_level = "trusted"`,
		``,
		`[model_providers.totality-openai]`,
		`name = "Totality OpenAI Proxy"`,
		`base_url = "http://127.0.0.1:43123/v1"`,
		`wire_api = "responses"`,
		`requires_openai_auth = true`,
		``,
	}, "\n")
	want := strings.Join([]string{
		`# personal settings`,
		`model = "gpt-5"`,
		``,
		`[projects."/repo"]`,
		`trust_level = "trusted"`,
		``,
	}, "\n")
	path := writeLegacyCodexConfig(t, home, fixture)

	reverted, err := revertLegacyCodexConfigFile(path)
	if err != nil {
		t.Fatalf("revertLegacyCodexConfigFile() error = %v", err)
	}
	if !reverted {
		t.Fatal("revertLegacyCodexConfigFile() reverted = false, want true")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("reverted config = %q, want %q", data, want)
	}
	backups := legacyCodexBackups(t, path)
	if len(backups) != 1 {
		t.Fatalf("backups = %v, want exactly one", backups)
	}
	backupData, err := os.ReadFile(backups[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(backupData) != fixture {
		t.Fatalf("backup content = %q, want original fixture", backupData)
	}
}

func TestRevertLegacyCodexConfigFileHonorsTotalityProxyAddr(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "127.0.0.1:50000")

	fixture := strings.Join([]string{
		`model_provider = "totality-openai"`,
		``,
		`[model_providers.totality-openai]`,
		`name = "Totality OpenAI Proxy"`,
		`base_url = "http://127.0.0.1:50000/v1"`,
		`wire_api = "responses"`,
		`requires_openai_auth = true`,
		``,
	}, "\n")
	path := writeLegacyCodexConfig(t, home, fixture)

	reverted, err := revertLegacyCodexConfigFile(path)
	if err != nil {
		t.Fatalf("revertLegacyCodexConfigFile() error = %v", err)
	}
	if !reverted {
		t.Fatal("revertLegacyCodexConfigFile() reverted = false, want true")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "totality-openai") {
		t.Fatalf("reverted config still mentions totality-openai:\n%s", data)
	}
	// The default port stays recognized even while the override is set.
	if !legacyCodexBaseURLPointsAtProxy("http://127.0.0.1:43123/v1") {
		t.Fatal("default proxy address no longer recognized with TOTALITY_PROXY_ADDR set")
	}
}

func TestRevertLegacyCodexConfigFileLeavesForeignConfigUntouched(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "")

	fixture := strings.Join([]string{
		`model_provider = "totality-openai"`,
		``,
		`[model_providers.totality-openai]`,
		`name = "my own gateway"`,
		`base_url = "https://llm.example.com/v1"`,
		``,
	}, "\n")
	path := writeLegacyCodexConfig(t, home, fixture)

	reverted, err := revertLegacyCodexConfigFile(path)
	if err != nil {
		t.Fatalf("revertLegacyCodexConfigFile() error = %v", err)
	}
	if reverted {
		t.Fatal("revertLegacyCodexConfigFile() reverted = true, want untouched")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != fixture {
		t.Fatalf("config changed: %q, want %q", data, fixture)
	}
	if backups := legacyCodexBackups(t, path); len(backups) != 0 {
		t.Fatalf("backups = %v, want none", backups)
	}
}

func TestRevertLegacyCodexConfigFileMissingFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".codex", "config.toml")
	reverted, err := revertLegacyCodexConfigFile(path)
	if err != nil {
		t.Fatalf("revertLegacyCodexConfigFile() error = %v", err)
	}
	if reverted {
		t.Fatal("revertLegacyCodexConfigFile() reverted = true, want no-op")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("config file was created: %v", err)
	}
}

func TestRevertLegacyCodexConfigFileBlockOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "")

	fixture := strings.Join([]string{
		`[model_providers.totality-openai]`,
		`name = "Totality OpenAI Proxy"`,
		`base_url = "http://127.0.0.1:43123/v1"`,
		``,
		`[projects."/repo"]`,
		`trust_level = "trusted"`,
		``,
	}, "\n")
	want := strings.Join([]string{
		`[projects."/repo"]`,
		`trust_level = "trusted"`,
		``,
	}, "\n")
	path := writeLegacyCodexConfig(t, home, fixture)

	reverted, err := revertLegacyCodexConfigFile(path)
	if err != nil {
		t.Fatalf("revertLegacyCodexConfigFile() error = %v", err)
	}
	if !reverted {
		t.Fatal("revertLegacyCodexConfigFile() reverted = false, want true")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != want {
		t.Fatalf("reverted config = %q, want %q", data, want)
	}
}

func TestRevertLegacyCodexConfigFileAssignmentOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "")

	fixture := strings.Join([]string{
		`model = "gpt-5"`,
		`model_provider = "totality-openai"`,
		``,
	}, "\n")
	path := writeLegacyCodexConfig(t, home, fixture)

	reverted, err := revertLegacyCodexConfigFile(path)
	if err != nil {
		t.Fatalf("revertLegacyCodexConfigFile() error = %v", err)
	}
	if !reverted {
		t.Fatal("revertLegacyCodexConfigFile() reverted = false, want true")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "model = \"gpt-5\"\n" {
		t.Fatalf("reverted config = %q", data)
	}
}

func TestLegacyCodexConfigPointsAtRetiredProxy(t *testing.T) {
	t.Setenv("TOTALITY_PROXY_ADDR", "")
	cases := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "block with default proxy",
			text: "[model_providers.totality-openai]\nbase_url = \"http://127.0.0.1:43123/v1\"\n",
			want: true,
		},
		{
			name: "block pointing elsewhere",
			text: "[model_providers.totality-openai]\nbase_url = \"https://llm.example.com/v1\"\n",
			want: false,
		},
		{
			name: "block without base_url",
			text: "[model_providers.totality-openai]\nname = \"Totality OpenAI Proxy\"\n",
			want: false,
		},
		{
			name: "assignment only",
			text: "model_provider = \"totality-openai\"\n",
			want: true,
		},
		{
			name: "foreign assignment only",
			text: "model_provider = \"openai\"\n",
			want: false,
		},
	}
	for _, tc := range cases {
		if got := legacyCodexConfigPointsAtRetiredProxy(tc.text); got != tc.want {
			t.Errorf("%s: legacyCodexConfigPointsAtRetiredProxy() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestCleanupLegacyAmbientCaptureRevertsCodexConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("TOTALITY_PROXY_ADDR", "")

	// Only the codex config is present: no plist, no shell blocks. The revert
	// must run on the config's own evidence without any launchctl calls.
	path := writeLegacyCodexConfig(t, home, strings.Join([]string{
		`model_provider = "totality-openai"`,
		``,
		`[model_providers.totality-openai]`,
		`name = "Totality OpenAI Proxy"`,
		`base_url = "http://127.0.0.1:43123/v1"`,
		`wire_api = "responses"`,
		`requires_openai_auth = true`,
		``,
	}, "\n"))

	runner := &fakeLegacyRunner{}
	result := cleanupLegacyAmbientCapture(context.Background(), runner.run)

	if !result.CleanedCodexConfig {
		t.Fatalf("cleanup result = %+v, want CleanedCodexConfig", result)
	}
	if result.CleanedLaunchAgent || result.CleanedShellBlocks || result.CleanedEnv {
		t.Fatalf("cleanup result = %+v, want only codex config cleaned", result)
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("cleanup warnings = %v, want none", result.Warnings)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("launchctl calls = %v, want none", runner.calls)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "totality-openai") {
		t.Fatalf("codex config still mentions totality-openai:\n%s", data)
	}
}

func TestCleanupLegacyAmbientCaptureSkipsLaunchctlWithoutArtifacts(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	// A shell profile without the managed block and no plist: cleanup must not
	// shell out to launchctl at all.
	if err := os.WriteFile(filepath.Join(home, ".zshrc"), []byte("alias ll='ls -la'\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	runner := &fakeLegacyRunner{}
	result := cleanupLegacyAmbientCapture(context.Background(), runner.run)

	if result.cleanedAny() || len(runner.calls) != 0 {
		t.Fatalf("cleanup result = %+v calls = %v, want full no-op", result, runner.calls)
	}
}
