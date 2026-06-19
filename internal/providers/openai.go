package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type OpenAIChatCompletions struct{}

func (OpenAIChatCompletions) Provider() string { return "openai" }

type OpenAIResponses struct{}

func (OpenAIResponses) Provider() string { return "openai" }

func (OpenAIResponses) Assemble(events io.Reader) ([]byte, Summary, error) {
	parsed, err := readSSEEvents(events)
	if err != nil {
		return nil, Summary{}, err
	}

	var summary Summary
	var last map[string]any
	for _, event := range parsed {
		var payload map[string]any
		if err := decodeJSON(event.Data, &payload); err != nil {
			return nil, Summary{}, err
		}
		if response, _ := payload["response"].(map[string]any); response != nil {
			last = response
			mergeOpenAIResponseSummary(&summary, response)
			continue
		}
		last = payload
		mergeOpenAIResponseSummary(&summary, payload)
	}
	if last == nil {
		last = map[string]any{}
	}
	body, err := json.Marshal(last)
	if err != nil {
		return nil, Summary{}, fmt.Errorf("marshal openai response: %w", err)
	}
	return body, summary, nil
}

func (OpenAIChatCompletions) Assemble(events io.Reader) ([]byte, Summary, error) {
	parsed, err := readSSEEvents(events)
	if err != nil {
		return nil, Summary{}, err
	}

	type toolCallState struct {
		ID       string
		Type     string
		Name     string
		Args     string
		Sequence int
	}
	type choiceState struct {
		Index        int
		Role         string
		Content      string
		FinishReason *string
		ToolCalls    map[int]*toolCallState
	}

	choices := map[int]*choiceState{}
	root := map[string]any{
		"object": "chat.completion",
	}
	var summary Summary

	for _, event := range parsed {
		var chunk map[string]any
		if err := decodeJSON(event.Data, &chunk); err != nil {
			return nil, Summary{}, err
		}
		for _, key := range []string{"id", "model", "system_fingerprint", "created"} {
			if value, ok := chunk[key]; ok {
				root[key] = value
			}
		}
		if usage, _ := chunk["usage"].(map[string]any); usage != nil {
			summary.Usage = openAIUsage(usage)
			root["usage"] = usage
		}
		rawChoices, _ := chunk["choices"].([]any)
		for _, rawChoice := range rawChoices {
			choiceMap, _ := rawChoice.(map[string]any)
			index := asInt(choiceMap["index"])
			state := choices[index]
			if state == nil {
				state = &choiceState{
					Index:     index,
					Role:      "assistant",
					ToolCalls: map[int]*toolCallState{},
				}
				choices[index] = state
			}
			delta, _ := choiceMap["delta"].(map[string]any)
			if role := asString(delta["role"]); role != "" {
				state.Role = role
			}
			if text := asString(delta["content"]); text != "" {
				state.Content += text
			}
			if finish := asString(choiceMap["finish_reason"]); finish != "" {
				state.FinishReason = stringPtr(finish)
				if summary.FinishReason == nil && index == 0 {
					summary.FinishReason = stringPtr(finish)
				}
			}
			rawToolCalls, _ := delta["tool_calls"].([]any)
			for _, rawToolCall := range rawToolCalls {
				toolCall, _ := rawToolCall.(map[string]any)
				toolIndex := asInt(toolCall["index"])
				entry := state.ToolCalls[toolIndex]
				if entry == nil {
					entry = &toolCallState{Type: "function", Sequence: toolIndex}
					state.ToolCalls[toolIndex] = entry
				}
				if id := asString(toolCall["id"]); id != "" {
					entry.ID = id
				}
				if typ := asString(toolCall["type"]); typ != "" {
					entry.Type = typ
				}
				if fn, _ := toolCall["function"].(map[string]any); fn != nil {
					if name := asString(fn["name"]); name != "" {
						entry.Name = name
					}
					if args := asString(fn["arguments"]); args != "" {
						entry.Args += args
					}
				}
			}
		}
	}

	order := make([]int, 0, len(choices))
	for idx := range choices {
		order = append(order, idx)
	}
	sort.Ints(order)

	outChoices := make([]any, 0, len(order))
	for _, idx := range order {
		state := choices[idx]
		message := map[string]any{
			"role":    state.Role,
			"content": state.Content,
		}
		if len(state.ToolCalls) > 0 {
			toolOrder := make([]int, 0, len(state.ToolCalls))
			for toolIdx := range state.ToolCalls {
				toolOrder = append(toolOrder, toolIdx)
			}
			sort.Ints(toolOrder)
			toolCalls := make([]any, 0, len(toolOrder))
			for _, toolIdx := range toolOrder {
				call := state.ToolCalls[toolIdx]
				toolCalls = append(toolCalls, map[string]any{
					"id":   call.ID,
					"type": call.Type,
					"function": map[string]any{
						"name":      call.Name,
						"arguments": call.Args,
					},
				})
			}
			message["tool_calls"] = toolCalls
		}

		outChoices = append(outChoices, map[string]any{
			"index":         state.Index,
			"message":       message,
			"finish_reason": derefString(state.FinishReason),
		})
	}
	root["choices"] = outChoices

	body, err := json.Marshal(root)
	if err != nil {
		return nil, Summary{}, fmt.Errorf("marshal openai completion: %w", err)
	}
	return body, summary, nil
}

func summarizeOpenAIJSON(body []byte) (Summary, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return Summary{}, err
	}

	var summary Summary
	if usage, _ := payload["usage"].(map[string]any); usage != nil {
		summary.Usage = openAIUsage(usage)
	}
	if choices, _ := payload["choices"].([]any); len(choices) > 0 {
		if first, _ := choices[0].(map[string]any); first != nil {
			if finish := asString(first["finish_reason"]); finish != "" {
				summary.FinishReason = stringPtr(finish)
			}
		}
	}
	return summary, nil
}

func mergeOpenAIResponseSummary(summary *Summary, payload map[string]any) {
	if usage, _ := payload["usage"].(map[string]any); usage != nil {
		summary.Usage = openAIUsage(usage)
	}
	if finish := asString(payload["finish_reason"]); finish != "" {
		summary.FinishReason = stringPtr(finish)
	}
	if status := asString(payload["status"]); status != "" && summary.FinishReason == nil {
		summary.FinishReason = stringPtr(status)
	}
}

func openAIUsage(usage map[string]any) Usage {
	return Usage{
		InputTokens:  firstInt(usage["prompt_tokens"], usage["input_tokens"]),
		OutputTokens: firstInt(usage["completion_tokens"], usage["output_tokens"]),
	}
}
