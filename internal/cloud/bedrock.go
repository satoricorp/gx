package cloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// defaultBedrockTimeout matches the server's own ceiling on /gx/bedrock/fight.
// A review call is a model call, not an upload: the upload timeout is far too
// short for a flagship model reading a large brief, and a client that gives up
// first turns a completed review into a transport error.
const defaultBedrockTimeout = 5 * time.Minute

// BedrockMessage is one turn in the Anthropic messages shape.
type BedrockMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// BedrockFightRequest is the request /gx/bedrock/fight accepts.
//
// It is deliberately the Anthropic InvokeModel shape rather than Bedrock's
// Converse shape: the CLI's direct-to-AWS leg has to speak InvokeModel, and the
// server normalizes both. One shape here means one body builder and one parser
// for both transports.
type BedrockFightRequest struct {
	Model            string           `json:"model"`
	AnthropicVersion string           `json:"anthropic_version,omitempty"`
	System           string           `json:"system,omitempty"`
	Messages         []BedrockMessage `json:"messages"`
	MaxTokens        int              `json:"max_tokens,omitempty"`
}

// BedrockContentBlock is one block of the model's reply.
type BedrockContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// BedrockFightResponse is the Anthropic-compatible half of the endpoint's
// dual-shape reply. The Converse-native fields beside it are ignored here.
type BedrockFightResponse struct {
	Model      string                `json:"model"`
	Content    []BedrockContentBlock `json:"content"`
	StopReason string                `json:"stop_reason"`
}

// Text joins the reply's text blocks, which is all a JSON-returning call needs.
func (r BedrockFightResponse) Text() string {
	var b strings.Builder
	for _, block := range r.Content {
		if block.Type != "" && block.Type != "text" {
			continue
		}
		b.WriteString(block.Text)
	}
	return b.String()
}

// NewBedrockClient is a cloud client sized for model calls rather than uploads.
// It is nil when cloud is not configured, like NewClient.
func NewBedrockClient() *Client {
	url := CloudURL()
	if url == "" {
		return nil
	}
	return &Client{url: url, http: &http.Client{Timeout: bedrockTimeout()}}
}

func bedrockTimeout() time.Duration {
	raw := strings.TrimSpace(os.Getenv("GX_CLOUD_BEDROCK_TIMEOUT"))
	if raw == "" {
		return defaultBedrockTimeout
	}
	timeout, err := time.ParseDuration(raw)
	if err != nil || timeout <= 0 {
		return defaultBedrockTimeout
	}
	return timeout
}

// BedrockFight runs one model call on gx's Bedrock account.
//
// gx holds the AWS credentials, so a user with no AWS account still gets the
// same reviewers. The endpoint is a passthrough: retrieval, the brief, the
// two-reviewer fan-out and the judge all stay in the CLI.
func (c *Client) BedrockFight(ctx context.Context, reqBody BedrockFightRequest) (BedrockFightResponse, error) {
	var result BedrockFightResponse
	if c == nil || c.url == "" {
		return result, fmt.Errorf("gx cloud base URL is not configured")
	}
	if err := c.postJSON(ctx, "/gx/bedrock/fight", reqBody, &result); err != nil {
		return BedrockFightResponse{}, err
	}
	return result, nil
}
