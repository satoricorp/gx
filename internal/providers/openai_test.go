package providers

import (
	"bytes"
	"testing"
)

func TestOpenAIAssemble(t *testing.T) {
	stream := []byte("data: {\"id\":\"chatcmpl_1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"gpt-4.1\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"hel\"},\"finish_reason\":null}]}\n\n" +
		"data: {\"id\":\"chatcmpl_1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"gpt-4.1\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"lo\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":7,\"completion_tokens\":2}}\n\n" +
		"data: [DONE]\n\n")

	body, summary, err := (OpenAIChatCompletions{}).Assemble(bytes.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(body, []byte(`"content":"hello"`)) {
		t.Fatalf("expected assembled content, got %s", body)
	}
	if summary.Usage.InputTokens == nil || *summary.Usage.InputTokens != 7 {
		t.Fatalf("expected prompt tokens, got %+v", summary.Usage)
	}
	if summary.Usage.OutputTokens == nil || *summary.Usage.OutputTokens != 2 {
		t.Fatalf("expected completion tokens, got %+v", summary.Usage)
	}
	if summary.FinishReason == nil || *summary.FinishReason != "stop" {
		t.Fatalf("expected finish reason, got %+v", summary)
	}
}

func TestOpenAIResponsesAssembleUsage(t *testing.T) {
	stream := []byte("event: response.in_progress\n" +
		"data: {\"type\":\"response.in_progress\",\"response\":{\"id\":\"resp_1\",\"status\":\"in_progress\",\"usage\":{\"input_tokens\":50,\"output_tokens\":5}}}\n\n" +
		"event: response.completed\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_1\",\"status\":\"completed\",\"finish_reason\":\"stop\",\"usage\":{\"input_tokens\":100,\"output_tokens\":25,\"total_tokens\":125}}}\n\n")

	body, summary, err := (OpenAIResponses{}).Assemble(bytes.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(body, []byte(`"id":"resp_1"`)) {
		t.Fatalf("expected assembled response body, got %s", body)
	}
	if summary.Usage.InputTokens == nil || *summary.Usage.InputTokens != 100 {
		t.Fatalf("expected input tokens, got %+v", summary.Usage)
	}
	if summary.Usage.OutputTokens == nil || *summary.Usage.OutputTokens != 25 {
		t.Fatalf("expected output tokens, got %+v", summary.Usage)
	}
	if summary.FinishReason == nil || *summary.FinishReason != "stop" {
		t.Fatalf("expected finish reason, got %+v", summary)
	}
}

func TestSummarizeOpenAICompactionUsage(t *testing.T) {
	body := []byte(`{
		"id": "resp_compact_1",
		"object": "response.compaction",
		"usage": {
			"input_tokens": 236389,
			"output_tokens": 1652,
			"total_tokens": 238041
		}
	}`)

	summary, err := summarizeOpenAIJSON(body)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Usage.InputTokens == nil || *summary.Usage.InputTokens != 236389 {
		t.Fatalf("expected input tokens, got %+v", summary.Usage)
	}
	if summary.Usage.OutputTokens == nil || *summary.Usage.OutputTokens != 1652 {
		t.Fatalf("expected output tokens, got %+v", summary.Usage)
	}
}
