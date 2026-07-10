package claudehooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const hookMarker = "gx capture transcript"

// MergeSettings registers GX transcript capture hooks in .claude/settings.json
// without replacing existing user hook entries.
func MergeSettings(repoRoot, gxPath string) error {
	if strings.TrimSpace(repoRoot) == "" {
		return fmt.Errorf("repo root required")
	}
	if strings.TrimSpace(gxPath) == "" {
		var err error
		gxPath, err = os.Executable()
		if err != nil {
			return err
		}
	}
	settingsPath := filepath.Join(repoRoot, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return err
	}

	settings := map[string]json.RawMessage{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("parse %s: %w", settingsPath, err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	hooksRaw, ok := settings["hooks"]
	var hooks map[string]json.RawMessage
	if ok {
		if err := json.Unmarshal(hooksRaw, &hooks); err != nil {
			return fmt.Errorf("parse hooks: %w", err)
		}
	}
	if hooks == nil {
		hooks = map[string]json.RawMessage{}
	}

	command := shellQuote(gxPath) + " capture transcript"
	for _, event := range []string{"Stop", "SessionEnd"} {
		merged, err := mergeHookEvent(hooks[event], command)
		if err != nil {
			return fmt.Errorf("merge %s hook: %w", event, err)
		}
		mergedBytes, err := json.Marshal(merged)
		if err != nil {
			return err
		}
		hooks[event] = mergedBytes
	}

	hooksBytes, err := json.Marshal(hooks)
	if err != nil {
		return err
	}
	settings["hooks"] = hooksBytes

	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(settingsPath, out, 0o644)
}

func mergeHookEvent(existing json.RawMessage, command string) ([]map[string]any, error) {
	var groups []map[string]any
	if len(existing) > 0 {
		if err := json.Unmarshal(existing, &groups); err != nil {
			return nil, err
		}
	}
	for _, group := range groups {
		if hookGroupContainsMarker(group) {
			return groups, nil
		}
	}
	entry := map[string]any{
		"hooks": []map[string]any{
			{
				"type":    "command",
				"command": command,
				"async":   true,
			},
		},
	}
	return append(groups, entry), nil
}

func hookGroupContainsMarker(group map[string]any) bool {
	rawHooks, ok := group["hooks"]
	if !ok {
		return false
	}
	data, err := json.Marshal(rawHooks)
	if err != nil {
		return false
	}
	var hooks []map[string]any
	if err := json.Unmarshal(data, &hooks); err != nil {
		return false
	}
	for _, hook := range hooks {
		cmd, _ := hook["command"].(string)
		if strings.Contains(cmd, hookMarker) {
			return true
		}
	}
	return false
}

func shellQuote(value string) string {
	if value == "" {
		return `""`
	}
	if !strings.ContainsAny(value, " \t\n\"'$\\`") {
		return value
	}
	return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
}
