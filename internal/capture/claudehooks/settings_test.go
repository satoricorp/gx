package claudehooks_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/capture/claudehooks"
)

func TestMergeSettingsCreatesHooks(t *testing.T) {
	repo := t.TempDir()
	gxPath := filepath.Join(repo, "bin", "gx")
	if err := os.MkdirAll(filepath.Dir(gxPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(gxPath, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := claudehooks.MergeSettings(repo, gxPath); err != nil {
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
	gxPath := filepath.Join(repo, "gx")
	if err := os.WriteFile(gxPath, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := claudehooks.MergeSettings(repo, gxPath); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(string(data), "echo keep-me", "gx capture transcript") {
		t.Fatalf("merged settings missing entries: %s", data)
	}
}

func TestMergeSettingsIdempotent(t *testing.T) {
	repo := t.TempDir()
	gxPath := filepath.Join(repo, "gx")
	if err := os.WriteFile(gxPath, []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := claudehooks.MergeSettings(repo, gxPath); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(filepath.Join(repo, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := claudehooks.MergeSettings(repo, gxPath); err != nil {
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
