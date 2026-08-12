package codereview

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOpenAICompleteSendsChatShapeAndParsesReply drives one leg call end to
// end against a fake server: bearer auth, the prefix stripped from the model,
// reasoning headroom on top of the caller's output budget, and the reply's
// text and finish_reason mapped into the shared completion shape.
func TestOpenAICompleteSendsChatShapeAndParsesReply(t *testing.T) {
	var got openAIChatRequest
	var auth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("request body did not decode: %v", err)
		}
		w.Write([]byte(`{"choices":[{"message":{"content":"{\"ok\":true}"},"finish_reason":"stop"}]}`))
	}))
	defer server.Close()

	transport := &openAITransport{apiKey: "sk-test", url: server.URL, client: server.Client()}
	completion, err := transport.complete(context.Background(), "openai:gpt-5.2", "be brief", "review this", 512)
	if err != nil {
		t.Fatalf("complete() error = %v", err)
	}
	if auth != "Bearer sk-test" {
		t.Errorf("Authorization = %q, want the bearer key", auth)
	}
	if got.Model != "gpt-5.2" {
		t.Errorf("model = %q, want the openai: prefix stripped", got.Model)
	}
	if got.MaxCompletionTokens != 512+openAIReasoningHeadroomTokens {
		t.Errorf("max_completion_tokens = %d, want the output budget plus reasoning headroom", got.MaxCompletionTokens)
	}
	if len(got.Messages) != 2 || got.Messages[0].Role != "system" || got.Messages[1].Role != "user" {
		t.Fatalf("messages = %+v, want system then user", got.Messages)
	}
	if completion.Text != `{"ok":true}` || completion.truncated() {
		t.Errorf("completion = %+v, want the reply text, not truncated", completion)
	}
}

// TestOpenAITruncationReadsAsMaxTokens keeps the truncation signal uniform
// across wires: finish_reason "length" must look exactly like Bedrock's
// stop_reason "max_tokens" at the parse site.
func TestOpenAITruncationReadsAsMaxTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"choices":[{"message":{"content":"{\"partial"},"finish_reason":"length"}]}`))
	}))
	defer server.Close()
	transport := &openAITransport{apiKey: "sk-test", url: server.URL, client: server.Client()}
	completion, err := transport.complete(context.Background(), "openai:gpt-5.2", "", "input", 64)
	if err != nil {
		t.Fatalf("complete() error = %v", err)
	}
	if !completion.truncated() {
		t.Errorf("truncated() = false for finish_reason length (StopReason %q)", completion.StopReason)
	}
}

// TestOpenAILegWithoutKeyIsUnavailableWithTheFix pins the failure mode: an
// openai: leg with no OPENAI_API_KEY must read as an unavailable reviewer
// naming the variable, never as a clean review with one silent leg.
func TestOpenAILegWithoutKeyIsUnavailableWithTheFix(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	_, err := newOpenAITransport()
	if err == nil {
		t.Fatal("newOpenAITransport() succeeded without a key")
	}
	if !strings.Contains(err.Error(), "OPENAI_API_KEY") {
		t.Errorf("error %q does not name OPENAI_API_KEY", err)
	}
}

func TestModelUsesOpenAIRouting(t *testing.T) {
	if !modelUsesOpenAI("openai:gpt-5.2") {
		t.Error("openai:gpt-5.2 should route to the OpenAI transport")
	}
	for _, m := range []string{"us.anthropic.claude-opus-4-6-v1", "zai.glm-5", "gpt-5.2"} {
		if modelUsesOpenAI(m) {
			t.Errorf("%q should NOT route to the OpenAI transport", m)
		}
	}
}
