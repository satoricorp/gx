package codereview

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"
)

// anthropicTestServer answers /v1/messages with the given replies in order,
// recording each request body.
func anthropicTestServer(t *testing.T, replies ...func(w http.ResponseWriter)) (*httptest.Server, *[]map[string]any, *atomic.Int32) {
	t.Helper()
	var bodies []map[string]any
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(calls.Add(1)) - 1
		if !strings.HasSuffix(r.URL.Path, "/v1/messages") {
			t.Errorf("request path = %s, want /v1/messages", r.URL.Path)
		}
		if got := r.Header.Get("X-Api-Key"); got != "sk-ant-test" {
			t.Errorf("x-api-key = %q, want the configured key", got)
		}
		if r.Header.Get("Anthropic-Version") == "" {
			t.Error("request carries no anthropic-version header")
		}
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("request body is not JSON: %v", err)
		}
		bodies = append(bodies, body)
		if n >= len(replies) {
			t.Errorf("unexpected request %d", n+1)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		replies[n](w)
	}))
	t.Cleanup(srv.Close)
	return srv, &bodies, &calls
}

func anthropicReply(status int, body string) func(w http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		// Keeps the SDK's retry backoff at a millisecond in tests.
		w.Header().Set("Retry-After-Ms", "1")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}
}

func anthropicMessage(stopReason, text string) string {
	reply, _ := json.Marshal(map[string]any{
		"id": "msg_test", "type": "message", "role": "assistant", "model": "claude-sonnet-4-6",
		"content":     []map[string]any{{"type": "text", "text": text}},
		"stop_reason": stopReason,
		"usage":       map[string]any{"input_tokens": 10, "output_tokens": 5},
	})
	return string(reply)
}

func testAnthropicTransport(srv *httptest.Server) *anthropicTransport {
	return newAnthropicTransport(anthropicCredential{apiKey: "sk-ant-test"}, option.WithBaseURL(srv.URL))
}

// The panel's models are Bedrock profiles everywhere upstream, so what reaches
// the wire has to be the Anthropic API's own name for the same model, and the
// prompt has to arrive as the system prompt plus one user turn, as it does on
// Bedrock.
func TestAnthropicTransportSendsTheReviewAsOneMessagesCall(t *testing.T) {
	srv, bodies, _ := anthropicTestServer(t, anthropicReply(http.StatusOK, anthropicMessage("end_turn", `{"findings":[]}`)))

	completion, err := testAnthropicTransport(srv).complete(context.Background(), "us.anthropic.claude-sonnet-4-6", "review rules", "the diff", 6000)
	if err != nil {
		t.Fatalf("complete() error = %v", err)
	}
	if completion.Text != `{"findings":[]}` || completion.StopReason != "end_turn" {
		t.Fatalf("completion = %+v, want the reply text and its stop reason", completion)
	}
	body := (*bodies)[0]
	if body["model"] != "claude-sonnet-4-6" {
		t.Fatalf("model = %v, want the Anthropic API ID claude-sonnet-4-6", body["model"])
	}
	if body["max_tokens"] != float64(6000) {
		t.Fatalf("max_tokens = %v, want 6000", body["max_tokens"])
	}
	if system, _ := json.Marshal(body["system"]); !strings.Contains(string(system), "review rules") {
		t.Fatalf("system = %s, want the review prompt", system)
	}
	messages, _ := body["messages"].([]any)
	if len(messages) != 1 || !strings.Contains(mustJSON(messages[0]), "the diff") {
		t.Fatalf("messages = %v, want one user turn carrying the input", body["messages"])
	}
	if _, ok := body["anthropic_version"]; ok {
		t.Fatal("body carries Bedrock's anthropic_version; the API takes it as a header")
	}
}

// A reply cut off at the cap has to read as truncated, as it does on every
// other wire, or the parse site blames the model's formatting.
func TestAnthropicTransportReportsATruncatedReply(t *testing.T) {
	srv, _, _ := anthropicTestServer(t, anthropicReply(http.StatusOK, anthropicMessage("max_tokens", `{"findings":[`)))

	completion, err := testAnthropicTransport(srv).complete(context.Background(), "us.anthropic.claude-haiku-4-5-20251001-v1:0", "", "x", 100)
	if err != nil {
		t.Fatalf("complete() error = %v", err)
	}
	if !completion.truncated() {
		t.Fatalf("completion = %+v, want it reported as truncated", completion)
	}
}

// A declined review must fail the leg, never pass for a review with nothing in it.
func TestAnthropicTransportFailsALegTheModelDeclined(t *testing.T) {
	srv, _, _ := anthropicTestServer(t, anthropicReply(http.StatusOK, anthropicMessage("refusal", "")))

	if _, err := testAnthropicTransport(srv).complete(context.Background(), "claude-sonnet-4-6", "", "x", 100); err == nil || !strings.Contains(err.Error(), "refusal") {
		t.Fatalf("complete() error = %v, want the refusal reported", err)
	}
}

// Throttles are retried; a request the API rejects is not.
func TestAnthropicTransportRetriesThrottlesButNotRejections(t *testing.T) {
	srv, _, calls := anthropicTestServer(t,
		anthropicReply(http.StatusTooManyRequests, `{"type":"error","error":{"type":"rate_limit_error","message":"slow down"}}`),
		anthropicReply(http.StatusOK, anthropicMessage("end_turn", "ok")),
	)
	if _, err := testAnthropicTransport(srv).complete(context.Background(), "claude-sonnet-4-6", "", "x", 100); err != nil {
		t.Fatalf("complete() after one 429 = %v, want the retry to succeed", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want the throttle retried once", calls.Load())
	}

	srv, _, calls = anthropicTestServer(t,
		anthropicReply(http.StatusBadRequest, `{"type":"error","error":{"type":"invalid_request_error","message":"bad model"}}`),
	)
	_, err := testAnthropicTransport(srv).complete(context.Background(), "claude-sonnet-4-6", "", "x", 100)
	if err == nil || !strings.Contains(err.Error(), "400") || !strings.Contains(err.Error(), "bad model") {
		t.Fatalf("complete() on a 400 = %v, want the status and the API's reason", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want a rejected request sent once", calls.Load())
	}
}

// Every model gx ships as a default has to map onto an Anthropic API ID, and a
// model the Anthropic API cannot serve has to say so before any call is made.
func TestAnthropicAPIModelIDMapsTheShippedModels(t *testing.T) {
	for _, tc := range []struct{ model, want string }{
		{defaultBedrockReviewModelA, "claude-haiku-4-5-20251001"},
		{defaultBedrockReviewModelB, "claude-sonnet-4-6"},
		{defaultBedrockJudgeModel, "claude-haiku-4-5-20251001"},
		{defaultRuleSortModel, "claude-haiku-4-5-20251001"},
		{normalizeBedrockModelID("claude-opus-5-5"), "claude-opus-5-5"},
		{"anthropic.claude-opus-4-6-v1", "claude-opus-4-6"},
		{"arn:aws:bedrock:us-west-2:123456789012:inference-profile/us.anthropic.claude-sonnet-4-6", "claude-sonnet-4-6"},
	} {
		model, want := tc.model, tc.want
		got, err := anthropicAPIModelID(model)
		if err != nil || got != want {
			t.Errorf("anthropicAPIModelID(%q) = %q, %v; want %q", model, got, err, want)
		}
	}
	for _, model := range []string{"zai.glm-5", "nvidia.nemotron-super-3-120b", "arn:aws:bedrock:us-west-2:123456789012:application-inference-profile/abc123"} {
		if _, err := anthropicAPIModelID(model); err == nil || !strings.Contains(err.Error(), "not an Anthropic model") {
			t.Errorf("anthropicAPIModelID(%q) error = %v, want it refused by name", model, err)
		}
	}
}

// localReviewEnv is a default build: no gx Cloud and no request for AWS.
func localReviewEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_CLOUD_URL", "")
	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
}

func TestTransportPlanUsesTheAnthropicKeyWhenNoCloudIsConfigured(t *testing.T) {
	localReviewEnv(t)
	if _, err := resolveBedrockTransportPlan(); err == nil || !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Fatalf("plan with nothing configured: err = %v, want the Anthropic key named as the fix", err)
	}

	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")
	plan, err := resolveBedrockTransportPlan()
	if err != nil || plan.Kind != bedrockTransportKindAnthropic {
		t.Fatalf("plan with ANTHROPIC_API_KEY = %+v, %v; want the Anthropic wire", plan.Kind, err)
	}
	if detail := plan.newTransport().detail(); !strings.Contains(detail, "Anthropic API") {
		t.Fatalf("transport detail = %q, want it named", detail)
	}

	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "oauth-token")
	if plan, err := resolveBedrockTransportPlan(); err != nil || plan.Kind != bedrockTransportKindAnthropic {
		t.Fatalf("plan with ANTHROPIC_AUTH_TOKEN = %+v, %v; want the Anthropic wire", plan.Kind, err)
	}
}

// The key is consulted only when nothing outranks it: an explicit request for
// AWS, or a configured cloud, still decides the wire.
func TestTransportPlanKeepsTheExplicitAndCloudWiresAheadOfTheKey(t *testing.T) {
	localReviewEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-test")

	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "1")
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAEXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	if plan, err := resolveBedrockTransportPlan(); err != nil || plan.Kind != bedrockTransportKindDirect {
		t.Fatalf("plan with GX_REVIEW_BEDROCK_DIRECT = %+v, %v; want direct AWS", plan.Kind, err)
	}

	t.Setenv("GX_REVIEW_BEDROCK_DIRECT", "")
	t.Setenv("GX_CLOUD_URL", "https://cloud.example.invalid")
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	if _, err := resolveBedrockTransportPlan(); err == nil || !strings.Contains(err.Error(), "gx auth login") {
		t.Fatalf("plan with a cloud configured and no session: err = %v, want the cloud's sign-in remedy rather than the key", err)
	}
}

// An openai: leg never rode the Bedrock plan, so a missing plan must not take
// it down with the others.
func TestOpenAILegsRunWithoutATransportPlan(t *testing.T) {
	localReviewEnv(t)
	t.Setenv("GX_REVIEW_AI", "")
	t.Setenv("GX_REVIEW_MODELS", "")
	t.Setenv("OPENAI_API_KEY", "sk-openai-test")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_A", "openai:gpt-5.2")
	t.Setenv("GX_REVIEW_BEDROCK_MODEL_B", "off")

	if reason := NoReviewerReason(); reason != "" {
		t.Fatalf("NoReviewerReason() = %q with an openai: leg and its key, want a reviewer", reason)
	}
}
