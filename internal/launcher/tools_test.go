package launcher

import (
	"reflect"
	"testing"
)

func TestCodexCommandArgs(t *testing.T) {
	tool := toolFor("codex")
	got := tool.CommandArgs(43123, []string{"codex", "exec", "hello"})
	want := []string{
		"-c", `model_provider="gx-openai"`,
		"-c", `model_providers.gx-openai.name="GX OpenAI Proxy"`,
		"-c", `model_providers.gx-openai.base_url="http://127.0.0.1:43123/v1"`,
		"-c", `model_providers.gx-openai.wire_api="responses"`,
		"-c", `model_providers.gx-openai.requires_openai_auth=true`,
		"-c", `transport="responses_http"`,
		"exec", "hello",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected codex args\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestClaudeCommandArgsUnchanged(t *testing.T) {
	tool := toolFor("claude")
	args := []string{"claude", "-p", "hi"}
	got := tool.CommandArgs(43123, args)
	if !reflect.DeepEqual(got, args) {
		t.Fatalf("expected args to stay unchanged, got %#v", got)
	}
}
