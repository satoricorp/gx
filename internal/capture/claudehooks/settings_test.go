package claudehooks_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/lgtm/internal/capture/claudehooks"
)

func TestMergeSettingsCreatesHooks(t *testing.T) {
	repo := t.TempDir()
	lgtmPath := filepath.Join(repo, "bin", "lgtm")
	if err := os.MkdirAll(filepath.Dir(lgtmPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lgtmPath, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := claudehooks.MergeSettings(repo, lgtmPath); err != nil {
		t.Fatal(err)
	}
	settingsPath := filepath.Join(repo, ".claude", "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	var settings map[string]any
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatal(err)
	}
	hooks, ok := settings["hooks"].(map[string]any)
	if !ok {
		t.Fatalf("hooks = %T", settings["hooks"])
	}
	for _, event := range []string{"Stop", "SessionEnd"} {
		groups, ok := hooks[event].([]any)
		if !ok || len(groups) == 0 {
			t.Fatalf("%s hooks = %#v", event, hooks[event])
		}
	}
}

func TestMergeSettingsPreservesExistingHooks(t *testing.T) {
	repo := t.TempDir()
	settingsPath := filepath.Join(repo, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	existing := `{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          {"type": "command", "command": "echo keep-me"}
        ]
      }
    ]
  }
}`
	if err := os.WriteFile(settingsPath, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	lgtmPath := filepath.Join(repo, "lgtm")
	if err := os.WriteFile(lgtmPath, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := claudehooks.MergeSettings(repo, lgtmPath); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(string(data), "echo keep-me", "lgtm capture transcript") {
		t.Fatalf("merged settings missing entries: %s", data)
	}
}

func TestMergeSettingsIdempotent(t *testing.T) {
	repo := t.TempDir()
	lgtmPath := filepath.Join(repo, "lgtm")
	if err := os.WriteFile(lgtmPath, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := claudehooks.MergeSettings(repo, lgtmPath); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(filepath.Join(repo, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := claudehooks.MergeSettings(repo, lgtmPath); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(filepath.Join(repo, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Fatalf("second merge changed settings:\nfirst=%s\nsecond=%s", first, second)
	}
}

func containsAll(body string, parts ...string) bool {
	for _, part := range parts {
		if !contains(body, part) {
			return false
		}
	}
	return true
}

func contains(body, part string) bool {
	return len(part) == 0 || (len(body) >= len(part) && indexOf(body, part) >= 0)
}

func indexOf(body, part string) int {
	for i := 0; i+len(part) <= len(body); i++ {
		if body[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}

// The hook command lands in a repo-shared, conventionally committed file, so
// it must not embed this machine's home directory. Regression: it used to
// write the absolute os.Executable() path (/Users/<name>/…), which pointed
// every teammate's hooks at a binary that does not exist on their machine.
func TestMergeSettingsSubstitutesHomeWithVariable(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	repo := t.TempDir()
	lgtmPath := filepath.Join(home, ".local", "bin", "lgtm")

	if err := claudehooks.MergeSettings(repo, lgtmPath); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repo, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	settings := string(data)
	if !strings.Contains(settings, `"\"$HOME/.local/bin/lgtm\" capture transcript"`) {
		t.Fatalf("settings do not carry the $HOME-relative command:\n%s", settings)
	}
	if strings.Contains(settings, home) {
		t.Fatalf("settings leak the literal home directory %q:\n%s", home, settings)
	}
}

// A binary outside the home directory has nothing portable to substitute.
func TestMergeSettingsKeepsNonHomePathsAbsolute(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	repo := t.TempDir()

	if err := claudehooks.MergeSettings(repo, "/opt/lgtm/bin/lgtm"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repo, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "/opt/lgtm/bin/lgtm capture transcript") {
		t.Fatalf("settings = %s, want the absolute non-home path kept", data)
	}
}
