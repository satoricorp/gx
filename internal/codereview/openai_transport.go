package codereview

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// openAIModelPrefix marks a reviewer leg that runs on the OpenAI API instead
// of Bedrock: GX_REVIEW_BEDROCK_MODEL_B="openai:gpt-5.2". The prefix keeps the
// two namespaces from colliding — "gpt-5.2" alone would otherwise be sent to
// bedrock-runtime, which would reject it with an error pointing at the wrong
// layer.
const (
	openAIModelPrefix     = "openai:"
	openAIChatCompletions = "https://api.openai.com/v1/chat/completions"

	// openAIReasoningHeadroomTokens pads max_completion_tokens above the
	// caller's output budget. OpenAI's reasoning models bill their hidden
	// reasoning against the same completion cap as the visible answer, so a
	// cap sized for the answer alone (defaultReviewMaxOutputTokens) can be
	// consumed entirely by reasoning, returning an empty reply with
	// finish_reason "length". The visible-answer budget the caller asked for
	// is preserved; this only makes room for the thinking on top of it.
	openAIReasoningHeadroomTokens = 16000

	openAITimeout = 10 * time.Minute
)

func modelUsesOpenAI(model string) bool {
	return strings.HasPrefix(strings.TrimSpace(model), openAIModelPrefix)
}

// openAITransport implements bedrockTransport against the OpenAI chat
// completions API, so an OpenAI model can sit in a reviewer leg without the
// panel, prompts, parsing, or failure reporting knowing the difference.
type openAITransport struct {
	apiKey string
	url    string
	client *http.Client
}

func newOpenAITransport() (*openAITransport, error) {
	key := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if key == "" {
		return nil, fmt.Errorf("an openai: reviewer leg needs OPENAI_API_KEY in the environment")
	}
	return &openAITransport{
		apiKey: key,
		url:    openAIChatCompletions,
		client: &http.Client{Timeout: openAITimeout},
	}, nil
}

func (t *openAITransport) detail() string { return "OpenAI API (api.openai.com)" }

type openAIChatRequest struct {
	Model               string          `json:"model"`
	Messages            []openAIMessage `json:"messages"`
	MaxCompletionTokens int             `json:"max_completion_tokens"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func (t *openAITransport) complete(ctx context.Context, model, system, input string, maxOutputTokens int) (bedrockCompletion, error) {
	if t == nil {
		return bedrockCompletion{}, fmt.Errorf("OpenAI reviewer has no transport")
	}
	name := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(model), openAIModelPrefix))
	if name == "" {
		return bedrockCompletion{}, fmt.Errorf("OpenAI request has no model after the %q prefix", openAIModelPrefix)
	}
	if maxOutputTokens <= 0 {
		maxOutputTokens = defaultReviewMaxOutputTokens
	}
	var messages []openAIMessage
	if strings.TrimSpace(system) != "" {
		messages = append(messages, openAIMessage{Role: "system", Content: system})
	}
	messages = append(messages, openAIMessage{Role: "user", Content: input})
	body, err := json.Marshal(openAIChatRequest{
		Model:               name,
		Messages:            messages,
		MaxCompletionTokens: maxOutputTokens + openAIReasoningHeadroomTokens,
	})
	if err != nil {
		return bedrockCompletion{}, fmt.Errorf("marshal OpenAI request: %w", err)
	}
	// Same retry contract as the Bedrock wire: throttles and server errors are
	// transient and a model call is minutes of work, so waiting seconds to
	// save one is always the right trade.
	var lastErr error
	for attempt := 0; ; attempt++ {
		completion, retryable, attemptErr := t.completeOnce(ctx, name, body)
		if attemptErr == nil {
			return completion, nil
		}
		lastErr = attemptErr
		if !retryable || attempt >= len(bedrockRetryBackoffs) {
			return bedrockCompletion{}, lastErr
		}
		select {
		case <-ctx.Done():
			return bedrockCompletion{}, lastErr
		case <-time.After(bedrockRetryBackoffs[attempt]):
		}
	}
}

func (t *openAITransport) completeOnce(ctx context.Context, model string, body []byte) (bedrockCompletion, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.url, bytes.NewReader(body))
	if err != nil {
		return bedrockCompletion{}, false, fmt.Errorf("create OpenAI request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	client := t.client
	if client == nil {
		client = &http.Client{Timeout: openAITimeout}
	}
	resp, err := client.Do(req)
	if err != nil {
		return bedrockCompletion{}, ctx.Err() == nil, fmt.Errorf("call OpenAI model %s: %w", model, err)
	}
	defer resp.Body.Close()
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		retryable := resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
		return bedrockCompletion{}, retryable, fmt.Errorf("OpenAI model %s returned %s: %s", model, resp.Status, firstLine(raw))
	}
	if readErr != nil {
		return bedrockCompletion{}, true, fmt.Errorf("read OpenAI response for %s: %w", model, readErr)
	}
	var decoded openAIChatResponse
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return bedrockCompletion{}, false, fmt.Errorf("decode OpenAI response for %s: %w", model, err)
	}
	if decoded.Error != nil {
		return bedrockCompletion{}, false, fmt.Errorf("OpenAI model %s error: %s", model, decoded.Error.Message)
	}
	if len(decoded.Choices) == 0 {
		return bedrockCompletion{}, false, fmt.Errorf("OpenAI model %s returned no choices", model)
	}
	choice := decoded.Choices[0]
	text := strings.TrimSpace(choice.Message.Content)
	if text == "" {
		return bedrockCompletion{}, false, fmt.Errorf("OpenAI model %s returned no content (finish_reason %s)", model, choice.FinishReason)
	}
	stop := choice.FinishReason
	// The parse site speaks Bedrock's dialect for truncation; translate so a
	// reply cut off at the cap reads the same on every wire.
	if strings.EqualFold(stop, "length") {
		stop = "max_tokens"
	}
	return bedrockCompletion{Text: text, StopReason: stop}, false, nil
}

// firstLine trims an error body to something an operator can read in a log
// line without the full JSON payload.
func firstLine(raw []byte) string {
	s := strings.TrimSpace(string(raw))
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 300 {
		s = s[:300]
	}
	return s
}
