package claude

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture"
)

// Parser normalizes Claude Code session JSONL into SessionEvents.
type Parser struct {
	Inventory *capture.InventoryCollector
}

func (p *Parser) Tool() string { return capture.ToolClaude }

// ParseFile reads one Claude session JSONL and returns normalized events.
func (p *Parser) ParseFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if p.Inventory != nil {
		p.Inventory.RecordSampleFile(capture.ToolClaude, path)
	}

	sessionID := strings.TrimSuffix(filepath.Base(path), ".jsonl")
	var events []capture.SessionEvent
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}
		recordShallowPaths(p.Inventory, capture.ToolClaude, raw, "")
		parsed, err := parseClaudeLine(raw, sessionID, repoRoot)
		if err != nil {
			return nil, fmt.Errorf("parse line: %w", err)
		}
		events = append(events, parsed...)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func parseClaudeLine(raw map[string]json.RawMessage, sessionID, repoRoot string) ([]capture.SessionEvent, error) {
	var events []capture.SessionEvent
	ts := parseTimestamp(raw["timestamp"])
	if sid := rawString(raw["sessionId"]); sid != "" {
		sessionID = sid
	}

	var msg map[string]json.RawMessage
	if err := json.Unmarshal(raw["message"], &msg); err != nil {
		return events, nil
	}
	recordShallowPaths(nil, capture.ToolClaude, msg, "message")

	model := rawString(msg["model"])
	var content []json.RawMessage
	if err := json.Unmarshal(msg["content"], &content); err != nil {
		return events, nil
	}

	for _, blockRaw := range content {
		var block map[string]json.RawMessage
		if err := json.Unmarshal(blockRaw, &block); err != nil {
			continue
		}
		blockType := rawString(block["type"])
		switch blockType {
		case "tool_use":
			ev := capture.SessionEvent{
				SessionID: sessionID,
				Tool:      capture.ToolClaude,
				Model:     model,
				TS:        ts,
				Kind:      capture.KindToolCall,
				Raw:       cloneRaw(block),
			}
			name := rawString(block["name"])
			var input map[string]json.RawMessage
			_ = json.Unmarshal(block["input"], &input)
			switch name {
			case "Write", "write_file":
				ev.Kind = capture.KindEdit
				ev.FilePath = relPath(rawString(input["file_path"]), repoRoot)
				ev.NewText = rawString(input["content"])
			case "Edit", "StrReplace", "MultiEdit":
				ev.Kind = capture.KindEdit
				ev.FilePath = relPath(rawString(input["file_path"]), repoRoot)
				ev.OldText = rawString(input["old_string"])
				ev.NewText = rawString(input["new_string"])
				if ev.NewText == "" {
					ev.NewText = rawString(input["newText"])
				}
			default:
				ev.Kind = capture.KindToolCall
			}
			events = append(events, ev)
		case "text":
			events = append(events, capture.SessionEvent{
				SessionID: sessionID,
				Tool:      capture.ToolClaude,
				Model:     model,
				TS:        ts,
				Kind:      capture.KindMessage,
				PromptContext: truncate(rawString(block["text"]), 512),
				Raw:       cloneRaw(block),
			})
		case "tool_result":
			events = append(events, capture.SessionEvent{
				SessionID: sessionID,
				Tool:      capture.ToolClaude,
				Model:     model,
				TS:        ts,
				Kind:      capture.KindToolResult,
				Raw:       cloneRaw(block),
			})
		}
	}
	return events, nil
}

func recordShallowPaths(inv *capture.InventoryCollector, tool string, obj map[string]json.RawMessage, prefix string) {
	if inv == nil {
		return
	}
	for key, val := range obj {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		inv.RecordPath(tool, path, inferJSONType(val), mapField(path), "")
		if strings.HasPrefix(string(val), "{") {
			var nested map[string]json.RawMessage
			if err := json.Unmarshal(val, &nested); err == nil {
				for nk := range nested {
					inv.RecordPath(tool, path+"."+nk, "unknown", mapField(path+"."+nk), "shallow")
				}
			}
		}
		if strings.HasPrefix(string(val), "[") {
			inv.RecordPath(tool, path+"[]", "array", mapField(path), "shallow")
		}
	}
}

func mapField(path string) string {
	switch {
	case strings.Contains(path, "sessionId"):
		return "session_id"
	case strings.Contains(path, "model"):
		return "model"
	case strings.Contains(path, "timestamp"):
		return "ts"
	case strings.Contains(path, "file_path"):
		return "file_path"
	case strings.Contains(path, "content") && strings.Contains(path, "input"):
		return "new_text"
	case strings.Contains(path, "old_string"):
		return "old_text"
	case strings.Contains(path, "tool_use"):
		return "kind"
	default:
		return ""
	}
}

func inferJSONType(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	switch {
	case s == "null":
		return "null"
	case strings.HasPrefix(s, "{"):
		return "object"
	case strings.HasPrefix(s, "["):
		return "array"
	case strings.HasPrefix(s, `"`):
		return "string"
	default:
		return "number"
	}
}

func parseTimestamp(raw json.RawMessage) int64 {
	s := rawString(raw)
	if s == "" {
		return 0
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UnixMilli()
		}
	}
	return 0
}

func rawString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return strings.Trim(string(raw), `"`)
}

func relPath(path, repoRoot string) string {
	path = filepath.ToSlash(strings.TrimSpace(path))
	if path == "" {
		return ""
	}
	if repoRoot != "" {
		absRepo, err := filepath.Abs(repoRoot)
		if err == nil {
			absRepo = filepath.ToSlash(absRepo)
			absPath, err := filepath.Abs(path)
			if err == nil {
				absPath = filepath.ToSlash(absPath)
				if rel, err := filepath.Rel(absRepo, absPath); err == nil && !strings.HasPrefix(rel, "..") {
					return rel
				}
			}
			prefix := absRepo + "/"
			if strings.HasPrefix(path, prefix) {
				return strings.TrimPrefix(path, prefix)
			}
		}
	}
	return strings.TrimPrefix(path, "./")
}

func cloneRaw(obj map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		out[k] = v
	}
	return out
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
