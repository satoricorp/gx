package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Anthropics struct{}

func (Anthropics) Provider() string { return "anthropic" }

func (Anthropics) Assemble(events io.Reader) ([]byte, Summary, error) {
	parsed, err := readSSEEvents(events)
	if err != nil {
		return nil, Summary{}, err
	}

	message := map[string]any{
		"content": []any{},
	}
	content := map[int]map[string]any{}
	jsonBuffers := map[int]*bytes.Buffer{}
	var summary Summary

	for _, event := range parsed {
		var payload map[string]any
		if err := decodeJSON(event.Data, &payload); err != nil {
			return nil, Summary{}, err
		}
		typ, _ := payload["type"].(string)
		switch typ {
		case "message_start":
			src, _ := payload["message"].(map[string]any)
			for k, v := range src {
				message[k] = v
			}
			if usage, _ := src["usage"].(map[string]any); usage != nil {
				summary.Usage = anthropicUsage(usage)
			}
		case "content_block_start":
			idx := asInt(payload["index"])
			block, _ := payload["content_block"].(map[string]any)
			clone := cloneMap(block)
			content[idx] = clone
			if blockType, _ := block["type"].(string); blockType == "tool_use" {
				jsonBuffers[idx] = &bytes.Buffer{}
			}
		case "content_block_delta":
			idx := asInt(payload["index"])
			block := content[idx]
			if block == nil {
				continue
			}
			delta, _ := payload["delta"].(map[string]any)
			switch delta["type"] {
			case "text_delta", "thinking_delta":
				appendString(block, fieldForDelta(delta["type"]), asString(delta["text"], delta["thinking"]))
			case "signature_delta":
				block["signature"] = asString(delta["signature"])
			case "input_json_delta", "tool_use_delta":
				if buf := jsonBuffers[idx]; buf != nil {
					buf.WriteString(asString(delta["partial_json"], delta["json"]))
				}
			}
		case "content_block_stop":
			idx := asInt(payload["index"])
			block := content[idx]
			if block == nil {
				continue
			}
			if buf := jsonBuffers[idx]; buf != nil && buf.Len() > 0 {
				var input any
				if err := json.Unmarshal(buf.Bytes(), &input); err == nil {
					block["input"] = input
				}
			}
		case "message_delta":
			if delta, _ := payload["delta"].(map[string]any); delta != nil {
				if finish := asString(delta["stop_reason"]); finish != "" {
					message["stop_reason"] = finish
					summary.FinishReason = stringPtr(finish)
				}
				if stopSequence := delta["stop_sequence"]; stopSequence != nil {
					message["stop_sequence"] = stopSequence
				}
			}
			if usage, _ := payload["usage"].(map[string]any); usage != nil {
				mergeUsage(&summary.Usage, anthropicUsage(usage))
				message["usage"] = usage
			}
		}
	}

	if len(content) > 0 {
		ordered := make([]any, 0, len(content))
		for idx := 0; idx < len(content); idx++ {
			if block, ok := content[idx]; ok {
				ordered = append(ordered, block)
			}
		}
		message["content"] = ordered
	}
	body, err := json.Marshal(message)
	if err != nil {
		return nil, Summary{}, fmt.Errorf("marshal anthropic message: %w", err)
	}
	return body, summary, nil
}

func summarizeAnthropicJSON(body []byte) (Summary, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return Summary{}, err
	}
	var summary Summary
	if usage, _ := payload["usage"].(map[string]any); usage != nil {
		summary.Usage = anthropicUsage(usage)
	}
	if finish := asString(payload["stop_reason"]); finish != "" {
		summary.FinishReason = stringPtr(finish)
	}
	return summary, nil
}

func anthropicUsage(usage map[string]any) Usage {
	return Usage{
		InputTokens:      optInt(usage["input_tokens"]),
		OutputTokens:     optInt(usage["output_tokens"]),
		CacheReadTokens:  optInt(usage["cache_read_input_tokens"]),
		CacheWriteTokens: optInt(usage["cache_creation_input_tokens"]),
	}
}
