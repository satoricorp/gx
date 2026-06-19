package codex

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture"
)

// Parser normalizes Codex rollout JSONL into SessionEvents.
type Parser struct {
	Inventory *capture.InventoryCollector
}

func (p *Parser) Tool() string { return capture.ToolCodex }

// ParseFile reads one Codex rollout JSONL and returns normalized events.
func (p *Parser) ParseFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if p.Inventory != nil {
		p.Inventory.RecordSampleFile(capture.ToolCodex, path)
	}

	sessionID := strings.TrimSuffix(filepath.Base(path), ".jsonl")
	model := ""
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
		recordShallowPaths(p.Inventory, capture.ToolCodex, raw, "")
		if sid, _ := metaFromLine(raw); sid != "" {
			sessionID = sid
		}
		if m := modelFromLine(raw); m != "" {
			model = m
		}
		parsed := parseCodexLine(raw, sessionID, model, repoRoot)
		events = append(events, parsed...)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func metaFromLine(raw map[string]json.RawMessage) (sessionID, model string) {
	if rawString(raw["type"]) != "session_meta" {
		return "", ""
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw["payload"], &payload); err != nil {
		return "", ""
	}
	return rawString(payload["session_id"]), rawString(payload["model"])
}

func modelFromLine(raw map[string]json.RawMessage) string {
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(raw["payload"], &payload); err != nil {
		return ""
	}
	return rawString(payload["model"])
}

func parseCodexLine(raw map[string]json.RawMessage, sessionID, model, repoRoot string) []capture.SessionEvent {
	ts := parseTimestamp(raw["timestamp"], raw["timestamp_ms"])
	lineType := rawString(raw["type"])

	var payload map[string]json.RawMessage
	_ = json.Unmarshal(raw["payload"], &payload)
	recordShallowPaths(nil, capture.ToolCodex, payload, "payload")

	switch lineType {
	case "response_item", "event", "function_call":
		return parseCodexPayload(payload, sessionID, model, ts, repoRoot)
	default:
		if len(payload) > 0 {
			return parseCodexPayload(payload, sessionID, model, ts, repoRoot)
		}
	}
	return nil
}

func parseCodexPayload(payload map[string]json.RawMessage, sessionID, model string, ts int64, repoRoot string) []capture.SessionEvent {
	if payload == nil {
		return nil
	}
	itemType := rawString(payload["type"])
	name := rawString(payload["name"])
	if name == "" {
		name = rawString(payload["tool"])
	}

	switch {
	case name == "apply_patch" || itemType == "apply_patch":
		text := rawString(payload["arguments"])
		if text == "" {
			text = rawString(payload["input"])
		}
		return parseApplyPatch(text, sessionID, model, ts, repoRoot, payload)
	case name == "write_file" || name == "Write":
		var args map[string]json.RawMessage
		if err := json.Unmarshal(payload["arguments"], &args); err != nil {
			_ = json.Unmarshal(payload["input"], &args)
		}
		filePath := relPath(rawString(args["path"]), repoRoot)
		if filePath == "" {
			filePath = relPath(rawString(args["file_path"]), repoRoot)
		}
		return []capture.SessionEvent{{
			SessionID: sessionID,
			Tool:      capture.ToolCodex,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			NewText:   rawString(args["content"]),
			Raw:       cloneRaw(payload),
		}}
	case itemType == "function_call":
		if name == "" {
			return nil
		}
		return parseCodexPayload(map[string]json.RawMessage{
			"name":      payload["name"],
			"arguments": payload["arguments"],
		}, sessionID, model, ts, repoRoot)
	}
	return nil
}

var patchFileRE = regexp.MustCompile(`(?m)^\*\*\* (?:Add|Update) File: (.+)$`)

func parseApplyPatch(text, sessionID, model string, ts int64, repoRoot string, payload map[string]json.RawMessage) []capture.SessionEvent {
	var events []capture.SessionEvent
	lines := strings.Split(text, "\n")
	var currentFile string
	var added []string
	flush := func() {
		if currentFile == "" || len(added) == 0 {
			return
		}
		events = append(events, capture.SessionEvent{
			SessionID: sessionID,
			Tool:      capture.ToolCodex,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  relPath(currentFile, repoRoot),
			NewText:   strings.Join(added, "\n"),
			Raw:       cloneRaw(payload),
		})
		added = nil
	}
	for _, line := range lines {
		if m := patchFileRE.FindStringSubmatch(line); len(m) == 2 {
			flush()
			currentFile = strings.TrimSpace(m[1])
			continue
		}
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			added = append(added, strings.TrimPrefix(line, "+"))
		}
	}
	flush()
	return events
}

func recordShallowPaths(inv *capture.InventoryCollector, tool string, obj map[string]json.RawMessage, prefix string) {
	if inv == nil || obj == nil {
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
	}
}

func mapField(path string) string {
	switch {
	case strings.Contains(path, "session_id"):
		return "session_id"
	case strings.Contains(path, "model"):
		return "model"
	case strings.Contains(path, "timestamp"):
		return "ts"
	case strings.Contains(path, "path"):
		return "file_path"
	case strings.Contains(path, "content"):
		return "new_text"
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

func parseTimestamp(primary, alt json.RawMessage) int64 {
	if ms := rawInt64(alt); ms > 0 {
		return ms
	}
	s := rawString(primary)
	if s == "" {
		return 0
	}
	layouts := []string{time.RFC3339Nano, time.RFC3339}
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

func rawInt64(raw json.RawMessage) int64 {
	if len(raw) == 0 {
		return 0
	}
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return n
	}
	return 0
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
