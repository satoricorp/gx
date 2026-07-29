package cursor

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/totality/internal/capture"
)

func (p *Parser) parseTranscriptFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	events, _, err := p.ParseBytes(data, path, repoRoot)
	return events, err
}

// ParseBytes normalizes Cursor transcript JSONL bytes.
func (p *Parser) ParseBytes(data []byte, sourcePath, repoRoot string) ([]capture.SessionEvent, int, error) {
	if p.Inventory != nil && strings.TrimSpace(sourcePath) != "" {
		p.Inventory.RecordSampleFile(capture.ToolCursor, sourcePath)
	}

	sessionID := transcriptSessionID(sourcePath)
	model := ""
	badLines := 0
	var events []capture.SessionEvent
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			badLines++
			continue
		}
		recordShallow(p.Inventory, raw, "")
		if sid := rawString(raw["sessionId"]); sid != "" {
			sessionID = "cursor:" + sid
		}
		if m := modelFromTranscriptLine(raw); m != "" {
			model = m
		}
		events = append(events, parseTranscriptLine(raw, sessionID, model, repoRoot)...)
	}
	if err := scanner.Err(); err != nil {
		return nil, badLines, err
	}
	return events, badLines, nil
}

func transcriptSessionID(path string) string {
	stem := strings.TrimSuffix(filepath.Base(path), ".jsonl")
	parent := filepath.Base(filepath.Dir(path))
	if parent == "subagents" {
		root := filepath.Base(filepath.Dir(filepath.Dir(path)))
		if root != "" {
			return "cursor:" + root + ":subagent:" + stem
		}
	}
	if stem != "" {
		return "cursor:" + stem
	}
	return "cursor:" + filepath.Base(path)
}

func parseTranscriptLine(raw map[string]json.RawMessage, sessionID, model, repoRoot string) []capture.SessionEvent {
	ts := parseTranscriptTimestamp(raw)
	var msg map[string]json.RawMessage
	if err := json.Unmarshal(raw["message"], &msg); err != nil {
		return nil
	}
	recordShallow(nil, msg, "message")
	if m := rawString(msg["model"]); m != "" {
		model = m
	}
	var content []json.RawMessage
	if err := json.Unmarshal(msg["content"], &content); err != nil {
		return nil
	}
	var events []capture.SessionEvent
	for _, blockRaw := range content {
		var block map[string]json.RawMessage
		if err := json.Unmarshal(blockRaw, &block); err != nil {
			continue
		}
		switch rawString(block["type"]) {
		case "tool_use":
			events = append(events, parseTranscriptToolUse(block, sessionID, model, ts, repoRoot)...)
		case "text":
			text := rawString(block["text"])
			if strings.TrimSpace(text) == "" {
				continue
			}
			events = append(events, capture.SessionEvent{
				SessionID:     sessionID,
				Tool:          capture.ToolCursor,
				Model:         model,
				TS:            ts,
				Kind:          capture.KindMessage,
				PromptContext: truncate(text, 512),
				Raw:           cloneRaw(block),
			})
		case "tool_result":
			events = append(events, capture.SessionEvent{
				SessionID: sessionID,
				Tool:      capture.ToolCursor,
				Model:     model,
				TS:        ts,
				Kind:      capture.KindToolResult,
				Raw:       cloneRaw(block),
			})
		}
	}
	return events
}

func parseTranscriptToolUse(block map[string]json.RawMessage, sessionID, model string, ts int64, repoRoot string) []capture.SessionEvent {
	name := normalizeToolName(rawString(block["name"]))
	var input map[string]json.RawMessage
	_ = json.Unmarshal(block["input"], &input)
	recordShallow(nil, input, "input")

	filePath := relPath(firstNonEmpty(
		rawString(input["path"]),
		rawString(input["file_path"]),
		rawString(input["relativeWorkspacePath"]),
	), repoRoot)
	raw := cloneRaw(block)
	switch name {
	case "strreplace", "edit", "edit_file":
		return []capture.SessionEvent{{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			OldText:   firstNonEmpty(rawString(input["old_string"]), rawString(input["oldString"])),
			NewText:   firstNonEmpty(rawString(input["new_string"]), rawString(input["newString"]), rawString(input["content"]), rawString(input["contents"])),
			Raw:       raw,
		}}
	case "multiedit":
		return parseTranscriptMultiEdit(input, sessionID, model, ts, repoRoot, raw)
	case "write", "write_file":
		return []capture.SessionEvent{{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			NewText:   firstNonEmpty(rawString(input["content"]), rawString(input["contents"])),
			Raw:       raw,
		}}
	case "applypatch", "apply_patch":
		return parsePatchText(rawString(input["patch"]), sessionID, model, ts, repoRoot, raw)
	case "shell":
		return parsePatchText(patchFromShellCommand(rawString(input["command"])), sessionID, model, ts, repoRoot, raw)
	}
	return []capture.SessionEvent{{
		SessionID: sessionID,
		Tool:      capture.ToolCursor,
		Model:     model,
		TS:        ts,
		Kind:      capture.KindToolCall,
		FilePath:  filePath,
		Raw:       raw,
	}}
}

func parseTranscriptMultiEdit(input map[string]json.RawMessage, sessionID, model string, ts int64, repoRoot string, raw map[string]json.RawMessage) []capture.SessionEvent {
	filePath := relPath(firstNonEmpty(rawString(input["path"]), rawString(input["file_path"])), repoRoot)
	var edits []struct {
		OldString string `json:"old_string"`
		NewString string `json:"new_string"`
		OldCamel  string `json:"oldString"`
		NewCamel  string `json:"newString"`
	}
	if err := json.Unmarshal(input["edits"], &edits); err != nil {
		return nil
	}
	var events []capture.SessionEvent
	for _, edit := range edits {
		newText := firstNonEmpty(edit.NewString, edit.NewCamel)
		if strings.TrimSpace(newText) == "" {
			continue
		}
		events = append(events, capture.SessionEvent{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			OldText:   firstNonEmpty(edit.OldString, edit.OldCamel),
			NewText:   newText,
			Raw:       raw,
		})
	}
	return events
}

func patchFromShellCommand(command string) string {
	if strings.TrimSpace(command) == "" || !strings.Contains(command, "*** Begin Patch") {
		return ""
	}
	start := strings.Index(command, "*** Begin Patch")
	end := strings.Index(command[start:], "*** End Patch")
	if end < 0 {
		return command[start:]
	}
	end += start + len("*** End Patch")
	return command[start:end]
}

func modelFromTranscriptLine(raw map[string]json.RawMessage) string {
	var msg map[string]json.RawMessage
	if err := json.Unmarshal(raw["message"], &msg); err != nil {
		return ""
	}
	return rawString(msg["model"])
}

func parseTranscriptTimestamp(raw map[string]json.RawMessage) int64 {
	if ms := rawInt64(raw["timestamp_ms"]); ms > 0 {
		return ms
	}
	for _, key := range []string{"timestamp", "createdAt", "created_at"} {
		if ts := parseTimeString(rawString(raw[key])); ts > 0 {
			return ts
		}
	}
	return 0
}

func parseTimeString(s string) int64 {
	if s == "" {
		return 0
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05.000Z"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UnixMilli()
		}
	}
	return 0
}
