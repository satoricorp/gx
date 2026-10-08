package codereview

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// anthropicCredential is how a review reaches the Anthropic API directly: an
// API key, or an OAuth access token (`ant auth print-credentials
// --access-token`) for callers who sign in rather than mint keys.
type anthropicCredential struct {
	apiKey    string
	authToken string
}

// anthropicCredentialFromEnv reads the credential, or reports false when the
// environment carries none.
func anthropicCredentialFromEnv() (anthropicCredential, bool) {
	if key := strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY")); key != "" {
		return anthropicCredential{apiKey: key}, true
	}
	if token := strings.TrimSpace(os.Getenv("ANTHROPIC_AUTH_TOKEN")); token != "" {
		return anthropicCredential{authToken: token}, true
	}
	return anthropicCredential{}, false
}

// anthropicTransport implements bedrockTransport against the Anthropic
// Messages API, so a review runs on the reviewer's own Anthropic account with
// no gx Cloud and no AWS account in the path.
//
// The panel's models are named as Bedrock inference profiles everywhere else
// (defaults, presets, GX_REVIEW_BEDROCK_MODEL_A/_B), and normalizeBedrockModelID
// turns even a bare `claude-opus-5-5` into one. Rather than fork every one of
// those names, this transport maps the profile back to the API's own ID at the
// wire, so the same model runs whichever wire carries it.
type anthropicTransport struct {
	client anthropic.Client
}

func newAnthropicTransport(cred anthropicCredential, opts ...option.RequestOption) *anthropicTransport {
	base := []option.RequestOption{
		// The SDK's own retries cover what the Bedrock wire's retry loop does
		// (429, 5xx, connection errors), and honour retry-after on top. A
		// model call is minutes of work, so a few seconds' wait to save one is
		// always the right trade.
		option.WithMaxRetries(len(bedrockRetryBackoffs)),
		option.WithRequestTimeout(bedrockDirectTimeout),
	}
	if cred.apiKey != "" {
		base = append(base, option.WithAPIKey(cred.apiKey))
	} else {
		base = append(base, option.WithAuthToken(cred.authToken))
	}
	return &anthropicTransport{client: anthropic.NewClient(append(base, opts...)...)}
}

func (t *anthropicTransport) detail() string { return "Anthropic API (api.anthropic.com)" }

func (t *anthropicTransport) complete(ctx context.Context, model, system, input string, maxOutputTokens int) (bedrockCompletion, error) {
	if t == nil {
		return bedrockCompletion{}, fmt.Errorf("Anthropic reviewer has no transport")
	}
	id, err := anthropicAPIModelID(model)
	if err != nil {
		return bedrockCompletion{}, err
	}
	if maxOutputTokens <= 0 {
		maxOutputTokens = defaultReviewMaxOutputTokens
	}
	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(id),
		MaxTokens: int64(maxOutputTokens),
		Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(input))},
	}
	if strings.TrimSpace(system) != "" {
		params.System = []anthropic.TextBlockParam{{Text: system}}
	}
	resp, err := t.client.Messages.New(ctx, params)
	if err != nil {
		var apiErr *anthropic.Error
		if errors.As(err, &apiErr) {
			return bedrockCompletion{}, fmt.Errorf("Anthropic model %s returned %d: %s", id, apiErr.StatusCode, firstLine([]byte(apiErr.RawJSON())))
		}
		return bedrockCompletion{}, fmt.Errorf("call Anthropic model %s: %w", id, err)
	}
	var text strings.Builder
	for _, block := range resp.Content {
		if block, ok := block.AsAny().(anthropic.TextBlock); ok {
			text.WriteString(block.Text)
		}
	}
	reply := strings.TrimSpace(text.String())
	if resp.StopReason == anthropic.StopReasonRefusal {
		// A declined review is not a review that found nothing; the leg reports
		// it as a failure so the run reads as degraded.
		return bedrockCompletion{}, fmt.Errorf("Anthropic model %s declined to review this change (stop_reason refusal)", id)
	}
	if reply == "" {
		return bedrockCompletion{}, fmt.Errorf("Anthropic model %s returned no content (stop_reason %s)", id, resp.StopReason)
	}
	return bedrockCompletion{Text: reply, StopReason: string(resp.StopReason)}, nil
}

// bedrockProfileVersionSuffix is the model version Bedrock appends to an
// Anthropic model ID: -v1:0, -v1, -v2:0. The Anthropic API names the same
// model without it.
var bedrockProfileVersionSuffix = regexp.MustCompile(`-v\d+(:\d+)?$`)

// anthropicAPIModelID maps a Bedrock model ID to the Anthropic API's ID for
// the same model:
//
//	us.anthropic.claude-sonnet-4-6                → claude-sonnet-4-6
//	us.anthropic.claude-haiku-4-5-20251001-v1:0   → claude-haiku-4-5-20251001
//	arn:aws:bedrock:…:inference-profile/us.anthropic.claude-opus-4-6-v1 → claude-opus-4-6
//
// An ID that names no Anthropic model (zai.glm-5, nvidia.nemotron-…, an
// application inference profile ARN) has no Anthropic API equivalent, and
// saying so here is better than a 404 that names a model the caller never
// wrote.
func anthropicAPIModelID(model string) (string, error) {
	model = strings.TrimSpace(model)
	name := model
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	if i := strings.Index(name, "anthropic."); i >= 0 {
		name = name[i+len("anthropic."):]
	}
	name = bedrockProfileVersionSuffix.ReplaceAllString(name, "")
	if !strings.HasPrefix(name, "claude-") {
		return "", fmt.Errorf("model %q is not an Anthropic model, so the Anthropic API cannot run it; point GX_REVIEW_BEDROCK_MODEL_A/_B (or GX_REVIEW_MODELS) at a claude-* model, or set %s=1 to run it on your own AWS account", model, bedrockDirectEnvVar)
	}
	return name, nil
}
