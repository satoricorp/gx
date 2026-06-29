package inference

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveWritesInferenceJSON0600(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	if err := Save(Credentials{Provider: "anthropic", APIKey: "anthropic-key"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	path := filepath.Join(home, "inference.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat inference.json: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("inference.json mode = %o, want 0600", info.Mode().Perm())
	}
	creds, ok, err := LoadStored()
	if err != nil || !ok {
		t.Fatalf("LoadStored() = ok %v err %v, want stored credentials", ok, err)
	}
	if creds.Provider != ProviderAnthropic || creds.APIKey != "anthropic-key" {
		t.Fatalf("LoadStored() = %#v, want anthropic key", creds)
	}
}

func TestResolveStoredKeyWinsOverEnv(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("ANTHROPIC_API_KEY", "env-anthropic")
	t.Setenv("OPENAI_API_KEY", "env-openai")
	if err := Save(Credentials{Provider: "openai", APIKey: "stored-openai"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	creds, ok := Resolve()
	if !ok {
		t.Fatal("Resolve() ok = false, want true")
	}
	if creds.Provider != ProviderOpenAI || creds.APIKey != "stored-openai" {
		t.Fatalf("Resolve() = %#v, want stored openai", creds)
	}
}

func TestResolveEnvFallbackOrder(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("ANTHROPIC_API_KEY", "env-anthropic")
	t.Setenv("OPENAI_API_KEY", "env-openai")

	creds, ok := Resolve()
	if !ok {
		t.Fatal("Resolve() ok = false, want true")
	}
	if creds.Provider != ProviderAnthropic || creds.APIKey != "env-anthropic" {
		t.Fatalf("Resolve() = %#v, want anthropic env first", creds)
	}
}

func TestApplyToEnvironmentKeepsOnlyActiveProvider(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("ANTHROPIC_API_KEY", "env-anthropic")
	t.Setenv("OPENAI_API_KEY", "env-openai")
	if err := Save(Credentials{Provider: "anthropic", APIKey: "stored-anthropic"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	ApplyToEnvironment()

	if got := os.Getenv("ANTHROPIC_API_KEY"); got != "stored-anthropic" {
		t.Fatalf("ANTHROPIC_API_KEY = %q, want stored-anthropic", got)
	}
	if got := os.Getenv("OPENAI_API_KEY"); got != "" {
		t.Fatalf("OPENAI_API_KEY = %q, want unset", got)
	}
}
