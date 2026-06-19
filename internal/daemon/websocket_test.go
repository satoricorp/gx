package daemon

import "testing"

func TestExtractSummaryFromWebSocketFrames(t *testing.T) {
	frames := []websocketFrame{
		{
			Type:    "text",
			Payload: `{"type":"response.created","response":{"id":"resp_1","status":"in_progress"}}`,
		},
		{
			Type:    "text",
			Payload: `{"type":"response.completed","response":{"id":"resp_1","status":"completed","finish_reason":"stop","usage":{"input_tokens":17618,"input_tokens_details":{"cached_tokens":0},"output_tokens":74,"output_tokens_details":{"reasoning_tokens":66},"total_tokens":17692}}}`,
		},
	}

	summary := extractSummaryFromWebSocketFrames(frames)
	if summary.Usage.InputTokens == nil || *summary.Usage.InputTokens != 17618 {
		t.Fatalf("expected input tokens, got %+v", summary.Usage)
	}
	if summary.Usage.OutputTokens == nil || *summary.Usage.OutputTokens != 74 {
		t.Fatalf("expected output tokens, got %+v", summary.Usage)
	}
	if summary.FinishReason == nil || *summary.FinishReason != "stop" {
		t.Fatalf("expected finish reason, got %+v", summary)
	}
}
