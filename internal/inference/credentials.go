package inference

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satoricorp/totality/internal/storage"
)

const fileName = "inference.json"

const (
	ProviderAnthropic = "anthropic"
	ProviderOpenAI    = "openai"
)

type Credentials struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
	Source   string `json:"-"`
}

func Path() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fileName), nil
}

func Save(creds Credentials) error {
	provider, err := NormalizeProvider(creds.Provider)
	if err != nil {
		return err
	}
	key := strings.TrimSpace(creds.APIKey)
	if key == "" {
		return fmt.Errorf("api key is required")
	}
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create tx home dir: %w", err)
	}
	file := Credentials{Provider: provider, APIKey: key}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", fileName, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", fileName, err)
	}
	return nil
}

func LoadStored() (Credentials, bool, error) {
	path, err := Path()
	if err != nil {
		return Credentials{}, false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Credentials{}, false, nil
		}
		return Credentials{}, false, fmt.Errorf("read %s: %w", fileName, err)
	}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return Credentials{}, false, fmt.Errorf("parse %s: %w", fileName, err)
	}
	provider, err := NormalizeProvider(creds.Provider)
	if err != nil {
		return Credentials{}, false, err
	}
	key := strings.TrimSpace(creds.APIKey)
	if key == "" {
		return Credentials{}, false, nil
	}
	return Credentials{Provider: provider, APIKey: key, Source: fileName}, true, nil
}

func Resolve() (Credentials, bool) {
	if creds, ok, err := LoadStored(); err == nil && ok {
		return creds, true
	}
	if key := strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY")); key != "" {
		return Credentials{Provider: ProviderAnthropic, APIKey: key, Source: "ANTHROPIC_API_KEY"}, true
	}
	if key := strings.TrimSpace(os.Getenv("OPENAI_API_KEY")); key != "" {
		return Credentials{Provider: ProviderOpenAI, APIKey: key, Source: "OPENAI_API_KEY"}, true
	}
	return Credentials{}, false
}

func EnvVars() map[string]string {
	creds, ok := Resolve()
	if !ok {
		return nil
	}
	switch creds.Provider {
	case ProviderAnthropic:
		return map[string]string{"ANTHROPIC_API_KEY": creds.APIKey}
	case ProviderOpenAI:
		return map[string]string{"OPENAI_API_KEY": creds.APIKey}
	default:
		return nil
	}
}

func ApplyToEnvironment() {
	creds, ok := Resolve()
	if !ok {
		return
	}
	switch creds.Provider {
	case ProviderAnthropic:
		_ = os.Setenv("ANTHROPIC_API_KEY", creds.APIKey)
		_ = os.Unsetenv("OPENAI_API_KEY")
	case ProviderOpenAI:
		_ = os.Setenv("OPENAI_API_KEY", creds.APIKey)
		_ = os.Unsetenv("ANTHROPIC_API_KEY")
	}
}

func NormalizeProvider(provider string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case ProviderAnthropic, "claude":
		return ProviderAnthropic, nil
	case ProviderOpenAI:
		return ProviderOpenAI, nil
	default:
		return "", fmt.Errorf("provider must be anthropic or openai")
	}
}

func ProviderDisplay(provider string) string {
	switch provider {
	case ProviderAnthropic:
		return "Anthropic"
	case ProviderOpenAI:
		return "OpenAI"
	default:
		return strings.TrimSpace(provider)
	}
}
