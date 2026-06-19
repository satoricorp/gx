package providers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Usage struct {
	InputTokens      *int
	OutputTokens     *int
	CacheReadTokens  *int
	CacheWriteTokens *int
}

type Summary struct {
	Usage        Usage
	FinishReason *string
}

type StreamHandler interface {
	Assemble(events io.Reader) (body []byte, summary Summary, err error)
	Provider() string
}

type sseEvent struct {
	Event string
	Data  string
}

func HandlerFor(provider, endpoint string) StreamHandler {
	switch provider {
	case "anthropic":
		if endpoint == "/v1/messages" {
			return Anthropics{}
		}
	case "openai":
		if endpoint == "/v1/chat/completions" {
			return OpenAIChatCompletions{}
		}
		if endpoint == "/v1/responses" {
			return OpenAIResponses{}
		}
	}
	return nil
}

func SummarizeJSON(provider string, body []byte) (Summary, error) {
	switch provider {
	case "anthropic":
		return summarizeAnthropicJSON(body)
	case "openai":
		return summarizeOpenAIJSON(body)
	default:
		return Summary{}, nil
	}
}

func readSSEEvents(r io.Reader) ([]sseEvent, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	normalized := bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n"))
	chunks := bytes.Split(normalized, []byte("\n\n"))
	events := make([]sseEvent, 0, len(chunks))
	for _, chunk := range chunks {
		if len(bytes.TrimSpace(chunk)) == 0 {
			continue
		}
		var ev sseEvent
		var dataLines []string
		for _, line := range bytes.Split(chunk, []byte("\n")) {
			switch {
			case bytes.HasPrefix(line, []byte("event:")):
				ev.Event = strings.TrimSpace(string(bytes.TrimPrefix(line, []byte("event:"))))
			case bytes.HasPrefix(line, []byte("data:")):
				dataLines = append(dataLines, strings.TrimSpace(string(bytes.TrimPrefix(line, []byte("data:")))))
			}
		}
		ev.Data = strings.Join(dataLines, "\n")
		if ev.Data == "" {
			continue
		}
		events = append(events, ev)
	}
	return events, nil
}

func intPtr(v int) *int {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func decodeJSON(data string, dst any) error {
	if data == "" || data == "[DONE]" {
		return nil
	}
	if err := json.Unmarshal([]byte(data), dst); err != nil {
		return fmt.Errorf("decode event json: %w", err)
	}
	return nil
}

func expectMap(v any) (map[string]any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("expected object")
	}
	return m, nil
}
