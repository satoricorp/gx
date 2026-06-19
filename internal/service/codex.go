package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const codexProviderID = "gx-openai"

type CodexStatus struct {
	ConfigPath string
	Found      bool
	Correct    bool
	Value      string
}

func CheckCodexConfig() CodexStatus {
	path := codexConfigPath()
	status := CodexStatus{ConfigPath: path}
	data, err := os.ReadFile(path)
	if err != nil {
		return status
	}
	status.Found = true
	text := string(data)
	status.Value = findTOMLTableStringValue(text, "model_providers."+codexProviderID, "base_url")
	if status.Value == "" {
		status.Value = findTOMLRootStringValue(text, "openai_base_url")
	}
	status.Correct = findTOMLRootStringValue(text, "model_provider") == codexProviderID && status.Value == OpenAIBaseURL()
	return status
}

func RepairCodexConfig(ctx context.Context, runner Runner) (CodexStatus, error) {
	status := CheckCodexConfig()
	if status.Correct {
		return status, nil
	}
	if runner == nil {
		runner = ExecRunner{}
	}
	if out, err := runner.Run(ctx, "codex", "config", "set", "openai_base_url", OpenAIBaseURL()); err == nil {
		status = CheckCodexConfig()
		if status.Correct {
			return status, nil
		}
		_ = out
	}
	if err := repairCodexConfigFile(status.ConfigPath); err != nil {
		return status, err
	}
	return CheckCodexConfig(), nil
}

func repairCodexConfigFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create codex config dir: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read codex config: %w", err)
	}
	if len(data) > 0 {
		backup := fmt.Sprintf("%s.gx-backup-%d", path, time.Now().Unix())
		if err := os.WriteFile(backup, data, 0o600); err != nil {
			return fmt.Errorf("write codex config backup: %w", err)
		}
	}
	text := repairCodexConfigText(string(data))
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		return fmt.Errorf("write codex config: %w", err)
	}
	return nil
}

func codexConfigPath() string {
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".codex", "config.toml")
	}
	return filepath.Join(".codex", "config.toml")
}

func findTOMLStringValue(text, key string) string {
	if strings.Contains(key, ".") {
		parts := strings.Split(key, ".")
		if len(parts) > 1 {
			table := strings.Join(parts[:len(parts)-1], ".")
			return findTOMLTableStringValue(text, table, parts[len(parts)-1])
		}
	}
	return findTOMLRootStringValue(text, key)
}

func findTOMLRootStringValue(text, key string) string {
	prefix := key + " "
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			return ""
		}
		if !strings.HasPrefix(trimmed, prefix) && !strings.HasPrefix(trimmed, key+"=") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	}
	return ""
}

func findTOMLTableStringValue(text, table, key string) string {
	inTable := false
	prefix := key + " "
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inTable = trimmed == "["+table+"]"
			continue
		}
		if !inTable {
			continue
		}
		if !strings.HasPrefix(trimmed, prefix) && !strings.HasPrefix(trimmed, key+"=") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) != 2 {
			continue
		}
		return strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	}
	return ""
}

func setTOMLStringValue(text, key, value string) string {
	if strings.Contains(key, ".") {
		parts := strings.Split(key, ".")
		if len(parts) > 1 {
			table := strings.Join(parts[:len(parts)-1], ".")
			return setTOMLTableStringValue(text, table, parts[len(parts)-1], value)
		}
	}
	return setTOMLRootStringValue(text, key, value)
}

func setTOMLRootStringValue(text, key, value string) string {
	lines := strings.Split(text, "\n")
	replacement := fmt.Sprintf(`%s = "%s"`, key, strings.ReplaceAll(value, `"`, `\"`))
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			break
		}
		if strings.HasPrefix(trimmed, key+" ") || strings.HasPrefix(trimmed, key+"=") {
			lines[i] = replacement
			return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
		}
	}
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "[") {
			insertAt := i
			if insertAt > 0 && strings.TrimSpace(lines[insertAt-1]) == "" {
				insertAt--
			}
			inserted := append([]string{}, lines[:insertAt]...)
			inserted = append(inserted, replacement)
			inserted = append(inserted, "")
			inserted = append(inserted, lines[i:]...)
			return strings.TrimRight(strings.Join(inserted, "\n"), "\n") + "\n"
		}
	}
	if strings.TrimSpace(text) != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return text + replacement + "\n"
}

func setTOMLTableStringValue(text, table, key, value string) string {
	lines := strings.Split(text, "\n")
	header := "[" + table + "]"
	replacement := fmt.Sprintf(`%s = "%s"`, key, strings.ReplaceAll(value, `"`, `\"`))
	tableStart := -1
	tableEnd := len(lines)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if tableStart >= 0 {
				tableEnd = i
				break
			}
			if trimmed == header {
				tableStart = i
			}
		}
	}
	if tableStart < 0 {
		text = strings.TrimRight(text, "\n")
		if strings.TrimSpace(text) != "" {
			text += "\n\n"
		}
		return text + header + "\n" + replacement + "\n"
	}
	for i := tableStart + 1; i < tableEnd; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, key+" ") || strings.HasPrefix(trimmed, key+"=") {
			lines[i] = replacement
			return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
		}
	}
	if tableEnd == len(lines) {
		for tableEnd > tableStart+1 && strings.TrimSpace(lines[tableEnd-1]) == "" {
			tableEnd--
		}
	}
	inserted := append([]string{}, lines[:tableEnd]...)
	inserted = append(inserted, replacement)
	inserted = append(inserted, lines[tableEnd:]...)
	return strings.TrimRight(strings.Join(inserted, "\n"), "\n") + "\n"
}

func setTOMLBoolValue(text, key string, value bool) string {
	if strings.Contains(key, ".") {
		parts := strings.Split(key, ".")
		if len(parts) > 1 {
			table := strings.Join(parts[:len(parts)-1], ".")
			return setTOMLTableLiteralValue(text, table, parts[len(parts)-1], fmt.Sprintf("%t", value))
		}
	}
	return setTOMLRootLiteralValue(text, key, fmt.Sprintf("%t", value))
}

func setTOMLRootLiteralValue(text, key, value string) string {
	lines := strings.Split(text, "\n")
	replacement := key + " = " + value
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			break
		}
		if strings.HasPrefix(trimmed, key+" ") || strings.HasPrefix(trimmed, key+"=") {
			lines[i] = replacement
			return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
		}
	}
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "[") {
			insertAt := i
			if insertAt > 0 && strings.TrimSpace(lines[insertAt-1]) == "" {
				insertAt--
			}
			inserted := append([]string{}, lines[:insertAt]...)
			inserted = append(inserted, replacement)
			inserted = append(inserted, "")
			inserted = append(inserted, lines[i:]...)
			return strings.TrimRight(strings.Join(inserted, "\n"), "\n") + "\n"
		}
	}
	if strings.TrimSpace(text) != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return text + replacement + "\n"
}

func setTOMLTableLiteralValue(text, table, key, value string) string {
	lines := strings.Split(text, "\n")
	header := "[" + table + "]"
	replacement := key + " = " + value
	tableStart := -1
	tableEnd := len(lines)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if tableStart >= 0 {
				tableEnd = i
				break
			}
			if trimmed == header {
				tableStart = i
			}
		}
	}
	if tableStart < 0 {
		text = strings.TrimRight(text, "\n")
		if strings.TrimSpace(text) != "" {
			text += "\n\n"
		}
		return text + header + "\n" + replacement + "\n"
	}
	for i := tableStart + 1; i < tableEnd; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, key+" ") || strings.HasPrefix(trimmed, key+"=") {
			lines[i] = replacement
			return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
		}
	}
	if tableEnd == len(lines) {
		for tableEnd > tableStart+1 && strings.TrimSpace(lines[tableEnd-1]) == "" {
			tableEnd--
		}
	}
	inserted := append([]string{}, lines[:tableEnd]...)
	inserted = append(inserted, replacement)
	inserted = append(inserted, lines[tableEnd:]...)
	return strings.TrimRight(strings.Join(inserted, "\n"), "\n") + "\n"
}

func repairCodexConfigText(text string) string {
	text = setTOMLStringValue(text, "model_provider", codexProviderID)
	text = setTOMLStringValue(text, "model_providers."+codexProviderID+".name", "GX OpenAI Proxy")
	text = setTOMLStringValue(text, "model_providers."+codexProviderID+".base_url", OpenAIBaseURL())
	text = setTOMLStringValue(text, "model_providers."+codexProviderID+".wire_api", "responses")
	text = setTOMLBoolValue(text, "model_providers."+codexProviderID+".requires_openai_auth", true)
	return text
}
