package providers

import (
	"bytes"
	"testing"
)

func TestAnthropicAssemble(t *testing.T) {
	stream := []byte("event: message_start\n" +
		"data: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_1\",\"type\":\"message\",\"role\":\"assistant\",\"content\":[],\"model\":\"claude-sonnet-4\",\"usage\":{\"input_tokens\":10}}}\n\n" +
		"event: content_block_start\n" +
		"data: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\",\"text\":\"\"}}\n\n" +
		"event: content_block_delta\n" +
		"data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"hello\"}}\n\n" +
		"event: message_delta\n" +
		"data: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":5}}\n\n" +
		"event: message_stop\n" +
		"data: {\"type\":\"message_stop\"}\n\n")

	body, summary, err := (Anthropics{}).Assemble(bytes.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(body, []byte(`"text":"hello"`)) {
		t.Fatalf("expected assembled text, got %s", body)
	}
	if summary.Usage.InputTokens == nil || *summary.Usage.InputTokens != 10 {
		t.Fatalf("expected input tokens, got %+v", summary.Usage)
	}
	if summary.Usage.OutputTokens == nil || *summary.Usage.OutputTokens != 5 {
		t.Fatalf("expected output tokens, got %+v", summary.Usage)
	}
	if summary.FinishReason == nil || *summary.FinishReason != "end_turn" {
		t.Fatalf("expected finish reason, got %+v", summary)
	}
}
