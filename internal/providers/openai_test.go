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

func TestSummarizeOpenAIResponsesJSONUsage(t *testing.T) {
	body := []byte(`{
		"id": "resp_1",
		"object": "response",
		"status": "completed",
		"usage": {
			"input_tokens": 56294,
			"input_tokens_details": {"cached_tokens": 53120},
			"output_tokens": 416,
			"output_tokens_details": {"reasoning_tokens": 162},
			"total_tokens": 56710
		}
	}`)

	summary, err := summarizeOpenAIJSON(body)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Usage.InputTokens == nil || *summary.Usage.InputTokens != 56294 {
		t.Fatalf("expected input tokens, got %+v", summary.Usage)
	}
	if summary.Usage.OutputTokens == nil || *summary.Usage.OutputTokens != 416 {
		t.Fatalf("expected output tokens, got %+v", summary.Usage)
	}
	if summary.Usage.CacheReadTokens == nil || *summary.Usage.CacheReadTokens != 53120 {
		t.Fatalf("expected cached tokens, got %+v", summary.Usage)
	}
	if summary.FinishReason == nil || *summary.FinishReason != "completed" {
		t.Fatalf("expected finish reason from status, got %+v", summary)
	}
}

func TestSummarizeOpenAIChatCompletionsCachedTokens(t *testing.T) {
	body := []byte(`{
		"id": "chatcmpl_1",
		"object": "chat.completion",
		"choices": [{"index": 0, "message": {"role": "assistant", "content": "hi"}, "finish_reason": "stop"}],
		"usage": {
			"prompt_tokens": 100,
			"prompt_tokens_details": {"cached_tokens": 75},
			"completion_tokens": 10
		}
	}`)

	summary, err := summarizeOpenAIJSON(body)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Usage.CacheReadTokens == nil || *summary.Usage.CacheReadTokens != 75 {
		t.Fatalf("expected cached tokens, got %+v", summary.Usage)
	}
	if summary.FinishReason == nil || *summary.FinishReason != "stop" {
		t.Fatalf("expected finish reason, got %+v", summary)
	}
}

func TestOpenAIResponsesAssembleStatusFinishReason(t *testing.T) {
	stream := []byte("event: response.created\n" +
		"data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_1\",\"status\":\"in_progress\"}}\n\n" +
		"event: response.completed\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_1\",\"status\":\"completed\",\"usage\":{\"input_tokens\":100,\"input_tokens_details\":{\"cached_tokens\":80},\"output_tokens\":25}}}\n\n")

	_, summary, err := (OpenAIResponses{}).Assemble(bytes.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	if summary.FinishReason == nil || *summary.FinishReason != "completed" {
		t.Fatalf("expected final status as finish reason, got %+v", summary.FinishReason)
	}
	if summary.Usage.CacheReadTokens == nil || *summary.Usage.CacheReadTokens != 80 {
		t.Fatalf("expected cached tokens, got %+v", summary.Usage)
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
