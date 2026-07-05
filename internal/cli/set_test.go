package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/inference"
)

func TestSetKeyStoresInteractiveInferenceKey(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetIn(strings.NewReader("1\nanthropic-key\n"))
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"set", "key"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v\n%s", err, out.String())
	}

	creds, ok, err := inference.LoadStored()
	if err != nil || !ok {
		t.Fatalf("LoadStored() = ok %v err %v, want stored key", ok, err)
	}
	if creds.Provider != inference.ProviderAnthropic || creds.APIKey != "anthropic-key" {
		t.Fatalf("stored creds = %#v, want anthropic key", creds)
	}
	if strings.Contains(out.String(), "anthropic-key") {
		t.Fatalf("set key output leaked key:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "Anthropic") || !strings.Contains(out.String(), "key saved") {
		t.Fatalf("set key output missing confirmation:\n%s", out.String())
	}
}

func TestSetKeyRejectsEmptyKey(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetIn(strings.NewReader("2\n\n"))
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"set", "key"})

	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "api key is required") {
		t.Fatalf("root.Execute() error = %v, want api key required", err)
	}
}

func TestRootHelpShowsInferenceKeyStatus(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	const apiKey = "apikey_openai-secret"
	if err := inference.Save(inference.Credentials{Provider: "openai", APIKey: apiKey}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "Using openai: apikey_...") {
		t.Fatalf("root help missing OpenAI key status:\n%s", text)
	}
	if strings.Contains(text, apiKey) {
		t.Fatalf("root help leaked OpenAI key:\n%s", text)
	}
}

func TestRootHelpShowsAnthropicInferenceKeyStatus(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	const apiKey = "apikey_anthropic-secret"
	if err := inference.Save(inference.Credentials{Provider: "anthropic", APIKey: apiKey}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "Using anthropic: apikey_...") {
		t.Fatalf("root help missing Anthropic key status:\n%s", text)
	}
	if strings.Contains(text, apiKey) {
		t.Fatalf("root help leaked Anthropic key:\n%s", text)
	}
}

func TestMaskedAPIKeyNeverShowsFullKey(t *testing.T) {
	for _, tc := range []struct {
		key  string
		want string
	}{
		{key: "apikey_secret", want: "apikey_..."},
		{key: "abcdefg", want: "abcdef..."},
		{key: "a", want: "..."},
		{key: "", want: "..."},
	} {
		if got := maskedAPIKey(tc.key); got != tc.want {
			t.Fatalf("maskedAPIKey(%q) = %q, want %q", tc.key, got, tc.want)
		}
	}
}

func TestRootHelpShowsMissingInferenceKeyStatus(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")
	root := NewRoot(context.Background())
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("root.Execute() error = %v", err)
	}

	text := out.String()
	if !strings.Contains(text, "API Key required. Run `gx set key` to set an API Key.") {
		t.Fatalf("root help missing API key required status:\n%s", text)
	}
}
