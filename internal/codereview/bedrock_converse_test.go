package codereview

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/cloud"
)

// capturingRoundTripper records the request it saw and returns one canned
// response, so a test can assert what actually went on the wire.
type capturingRoundTripper struct {
	request *http.Request
	body    []byte
	respond *http.Response
}

func (c *capturingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	c.request = req
	if req.Body != nil {
		c.body, _ = io.ReadAll(req.Body)
	}
	return c.respond, nil
}

func TestBedrockActionForModelRoutesByVendor(t *testing.T) {
	cases := []struct {
		model string
		want  string
	}{
		{"us.anthropic.claude-opus-4-6-v1", "invoke"},
		{"us.anthropic.claude-haiku-4-5-20251001-v1:0", "invoke"},
		{"arn:aws:bedrock:us-west-2::foundation-model/anthropic.claude-opus-4-6-v1", "invoke"},
		{"us.meta.llama4-maverick-17b-instruct-v1:0", "converse"},
		{"us.amazon.nova-2-lite-v1:0", "converse"},
		{"us.writer.palmyra-x5-v1:0", "converse"},
	}
	for _, tc := range cases {
		if got := bedrockActionForModel(tc.model); got != tc.want {
			t.Errorf("bedrockActionForModel(%q) = %q, want %q", tc.model, got, tc.want)
		}
	}
}

// TestConverseCompleteSendsConverseShapeAndParsesReply drives one non-Anthropic
// call end to end: the wire path must be the single-encoded /converse endpoint,
// the body must be the Converse shape (no anthropic_version), and the reply's
// text blocks must come back joined with the stop reason preserved.
func TestConverseCompleteSendsConverseShapeAndParsesReply(t *testing.T) {
	rt := &capturingRoundTripper{respond: cannedResponse(200, `{
		"output": {"message": {"role": "assistant", "content": [
			{"text": "first"},
			{"reasoningContent": {"reasoningText": {"text": "scratch work"}}},
			{"text": " second"}
		]}},
		"stopReason": "end_turn"
	}`)}
	transport := retryTestTransport(rt)
	const model = "us.meta.llama4-maverick-17b-instruct-v1:0"

	completion, err := transport.complete(context.Background(), model, "be brief", "review this", 512)
	if err != nil {
		t.Fatalf("complete() error = %v", err)
	}
	if completion.Text != "first second" {
		t.Errorf("Text = %q, want the text blocks joined without the reasoning block", completion.Text)
	}
	if completion.StopReason != "end_turn" {
		t.Errorf("StopReason = %q, want end_turn", completion.StopReason)
	}
	if completion.truncated() {
		t.Error("truncated() = true for an end_turn reply")
	}

	if got, want := rt.request.URL.EscapedPath(), "/model/us.meta.llama4-maverick-17b-instruct-v1%3A0/converse"; got != want {
		t.Errorf("wire path = %q, want %q", got, want)
	}
	var sent map[string]any
	if err := json.Unmarshal(rt.body, &sent); err != nil {
		t.Fatalf("request body is not JSON: %v", err)
	}
	if _, hasAnthropic := sent["anthropic_version"]; hasAnthropic {
		t.Error("Converse body carries anthropic_version; that is the InvokeModel dialect")
	}
	system, ok := sent["system"].([]any)
	if !ok || len(system) != 1 {
		t.Fatalf("system = %v, want a one-element block list", sent["system"])
	}
	if block := system[0].(map[string]any); block["text"] != "be brief" {
		t.Errorf("system text = %v, want the system prompt", block["text"])
	}
	inference, ok := sent["inferenceConfig"].(map[string]any)
	if !ok || inference["maxTokens"] != float64(512) {
		t.Errorf("inferenceConfig = %v, want maxTokens 512", sent["inferenceConfig"])
	}
}

// TestConverseTruncationSurfacesAsMaxTokens keeps the truncation signal alive
// across the second wire: a Converse reply cut off at the output cap must look
// exactly like an InvokeModel one to the parse site.
func TestConverseTruncationSurfacesAsMaxTokens(t *testing.T) {
	rt := &capturingRoundTripper{respond: cannedResponse(200, `{
		"output": {"message": {"content": [{"text": "{\"partial"}]}},
		"stopReason": "max_tokens"
	}`)}
	transport := retryTestTransport(rt)
	completion, err := transport.complete(context.Background(), "us.amazon.nova-2-lite-v1:0", "", "input", 64)
	if err != nil {
		t.Fatalf("complete() error = %v", err)
	}
	if !completion.truncated() {
		t.Errorf("truncated() = false for stopReason max_tokens (StopReason %q)", completion.StopReason)
	}
}

// TestCloudTransportRefusesConverseModels pins the honest failure: gx Cloud
// only normalizes the Anthropic shape, so a non-Anthropic leg through Cloud
// must fail here, naming the fix, not server-side with a misleading error.
// The cloud transport used to refuse every non-Anthropic model before the
// network, saying gx Cloud "does not speak Converse yet". It does — the fight
// route is a Converse passthrough — and whether a model is allowed is the
// server's decision, answered as model_not_allowed with the list. So the
// transport now forwards a Converse-family model like any other: with a bare
// client there is no server to reach, and the failure is a transport error
// from the HTTP layer, not the old pre-network refusal.
func TestCloudTransportForwardsConverseModelsToTheServer(t *testing.T) {
	transport := &cloudBedrockTransport{client: &cloud.Client{}, url: "https://example.invalid"}
	_, err := transport.complete(
		context.Background(), "zai.glm-5", "s", "i", 32)
	if err == nil {
		t.Fatal("a bare client cannot succeed; expected a transport error")
	}
	if strings.Contains(err.Error(), "does not speak") || strings.Contains(err.Error(), bedrockDirectEnvVar) {
		t.Errorf("cloud transport still refuses Converse-family models before the network: %q", err)
	}
	if !strings.Contains(err.Error(), "gx Cloud Bedrock call for zai.glm-5 failed") {
		t.Errorf("error should come from the forwarded call, got %q", err)
	}
}
