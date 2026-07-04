package cursor

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/satoricorp/gx/internal/capture"
	ingestcursor "github.com/satoricorp/gx/internal/ingest/cursor"
)

// Parser normalizes Cursor composer bubbles from state.vscdb into SessionEvents.
type Parser struct {
	Inventory    *capture.InventoryCollector
	Since        time.Time
	Until        time.Time
	contentCache map[string]string
}

func (p *Parser) Tool() string { return capture.ToolCursor }

// ParseFile reads Cursor's global state.vscdb (copy-on-read, WAL-safe).
func (p *Parser) ParseFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	if strings.HasSuffix(path, ".jsonl") {
		return p.parseTranscriptFile(path, repoRoot)
	}
	if p.Inventory != nil {
		p.Inventory.RecordSampleFile(capture.ToolCursor, path)
	}

	dbPath, cleanup, err := copyOnRead(path)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	db, err := openReadOnly(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	p.contentCache = map[string]string{}
	composers, err := loadComposers(context.Background(), db)
	if err != nil {
		return nil, err
	}
	bubbles, err := loadBubbles(context.Background(), db)
	if err != nil {
		return nil, err
	}
	workspaces := loadWorkspaceFolders()

	var events []capture.SessionEvent
	for composerID, bs := range bubbles {
		if len(bs) == 0 {
			continue
		}
		cwd, detectedRepo := pickWorkspace(workspaces, bs)
		if repoRoot != "" && detectedRepo != "" {
			absRepo, err := filepath.Abs(repoRoot)
			if err == nil {
				absDetected, err := filepath.Abs(detectedRepo)
				if err == nil && absDetected != absRepo {
					continue
				}
			}
		} else if repoRoot != "" && cwd != "" {
			absRepo, err := filepath.Abs(repoRoot)
			if err == nil {
				absCwd, err := filepath.Abs(cwd)
				if err == nil && !strings.HasPrefix(absCwd, absRepo+string(filepath.Separator)) && absCwd != absRepo {
					continue
				}
			}
		}

		meta := composers[composerID]
		sessionID := "cursor-" + composerID
		lastUpdated := composerLastSeen(meta, bs)
		if !inTimeWindow(lastUpdated, p.Since, p.Until) && !bubblesInWindow(bs, p.Since, p.Until) {
			continue
		}

		for _, b := range bs {
			if !inTimeWindow(b.createdAt, p.Since, p.Until) {
				continue
			}
			parsed := p.parseBubble(db, b, sessionID, repoRoot)
			events = append(events, parsed...)
		}
	}
	return events, nil
}

func (p *Parser) parseTranscriptFile(path string, repoRoot string) ([]capture.SessionEvent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if p.Inventory != nil {
		p.Inventory.RecordSampleFile(capture.ToolCursor, path)
	}

	sessionID := transcriptSessionID(path)
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
		return nil, err
	}
	return events, nil
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

func copyOnRead(vscdbPath string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "cursor-vscdb-*")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	base := filepath.Base(vscdbPath)
	for _, suffix := range []string{"", "-wal", "-shm"} {
		src := vscdbPath + suffix
		if _, err := os.Stat(src); err != nil {
			continue
		}
		dstName := base + suffix
		if err := copyFile(src, filepath.Join(dir, dstName)); err != nil {
			cleanup()
			return "", nil, err
		}
	}
	return filepath.Join(dir, base), cleanup, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func openReadOnly(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?mode=ro&immutable=1"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open cursor vscdb: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping cursor vscdb: %w", err)
	}
	return db, nil
}

type composerMeta struct {
	name          string
	createdAt     int64
	lastUpdatedAt int64
}

func loadComposers(ctx context.Context, db *sql.DB) (map[string]composerMeta, error) {
	rows, err := db.QueryContext(ctx, `SELECT key, value FROM cursorDiskKV WHERE key LIKE 'composerData:%'`)
	if err != nil {
		return nil, fmt.Errorf("query composers: %w", err)
	}
	defer rows.Close()

	out := map[string]composerMeta{}
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		composerID := strings.TrimPrefix(key, "composerData:")
		var doc struct {
			Name          string `json:"name"`
			CreatedAt     int64  `json:"createdAt"`
			LastUpdatedAt int64  `json:"lastUpdatedAt"`
		}
		_ = json.Unmarshal(value, &doc)
		out[composerID] = composerMeta{
			name:          doc.Name,
			createdAt:     doc.CreatedAt,
			lastUpdatedAt: doc.LastUpdatedAt,
		}
		recordBubblePaths(nil, value, "")
	}
	return out, rows.Err()
}

type bubble struct {
	id        string
	composer  string
	role      string
	text      string
	raw       []byte
	createdAt int64
}

func loadBubbles(ctx context.Context, db *sql.DB) (map[string][]bubble, error) {
	rows, err := db.QueryContext(ctx, `SELECT key, value FROM cursorDiskKV WHERE key LIKE 'bubbleId:%'`)
	if err != nil {
		return nil, fmt.Errorf("query bubbles: %w", err)
	}
	defer rows.Close()

	out := map[string][]bubble{}
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		composerID, bubbleID, ok := parseBubbleKey(key)
		if !ok {
			continue
		}
		var doc struct {
			Type       int    `json:"type"`
			Text       string `json:"text"`
			CreatedAt  string `json:"createdAt"`
			TokenCount struct {
				InputTokens  int `json:"inputTokens"`
				OutputTokens int `json:"outputTokens"`
			} `json:"tokenCount"`
			TimingInfo struct {
				ClientRpcSendTime int64 `json:"clientRpcSendTime"`
				ClientEndTime     int64 `json:"clientEndTime"`
			} `json:"timingInfo"`
		}
		_ = json.Unmarshal(value, &doc)
		ts := parseBubbleTime(doc.CreatedAt, doc.TimingInfo.ClientRpcSendTime, doc.TimingInfo.ClientEndTime)
		out[composerID] = append(out[composerID], bubble{
			id:        bubbleID,
			composer:  composerID,
			role:      roleFromType(doc.Type),
			text:      doc.Text,
			raw:       append([]byte(nil), value...),
			createdAt: ts,
		})
	}
	return out, rows.Err()
}

func (p *Parser) parseBubble(db *sql.DB, b bubble, sessionID, repoRoot string) []capture.SessionEvent {
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(b.raw, &doc); err != nil {
		return nil
	}
	recordBubblePaths(p.Inventory, b.raw, "")

	ts := b.createdAt
	model := ""

	var events []capture.SessionEvent
	if toolRaw, ok := doc["toolFormerData"]; ok {
		var toolDoc map[string]json.RawMessage
		if err := json.Unmarshal(toolRaw, &toolDoc); err == nil {
			recordShallow(p.Inventory, toolDoc, "toolFormerData")
			events = append(events, p.parseToolFormer(db, toolDoc, sessionID, model, ts, repoRoot, doc)...)
		}
	}

	if len(events) == 0 && len(doc["codeBlocks"]) > 0 {
		events = append(events, parseCodeBlockEdits(doc, sessionID, model, ts, repoRoot)...)
	}

	if b.text != "" && len(events) == 0 {
		events = append(events, capture.SessionEvent{
			SessionID:     sessionID,
			Tool:          capture.ToolCursor,
			Model:         model,
			TS:            ts,
			Kind:          capture.KindMessage,
			PromptContext: truncate(b.text, 512),
		})
	}
	return events
}

func (p *Parser) parseToolFormer(db *sql.DB, toolDoc map[string]json.RawMessage, sessionID, model string, ts int64, repoRoot string, bubbleDoc map[string]json.RawMessage) []capture.SessionEvent {
	name := normalizeToolName(rawString(toolDoc["name"]))
	status := rawString(toolDoc["status"])
	if status != "" && status != "completed" {
		return nil
	}

	var params map[string]json.RawMessage
	if rawParams := toolDoc["params"]; len(rawParams) > 0 {
		if err := json.Unmarshal(rawParams, &params); err != nil {
			var paramsStr string
			if err2 := json.Unmarshal(rawParams, &paramsStr); err2 == nil {
				_ = json.Unmarshal([]byte(paramsStr), &params)
			}
		}
	}
	recordShallow(nil, params, "toolFormerData.params")

	filePath := relPath(rawString(params["relativeWorkspacePath"]), repoRoot)
	if filePath == "" {
		filePath = relPath(filePathFromRawArgs(rawString(toolDoc["rawArgs"])), repoRoot)
	}

	switch name {
	case "search_replace", "strreplace":
		return []capture.SessionEvent{{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			OldText:   rawString(params["oldString"]),
			NewText:   rawString(params["newString"]),
			Raw:       cloneRaw(toolDoc),
		}}
	case "edit_file":
		newText := firstNonEmpty(
			rawString(params["streamingContent"]),
			rawString(params["contents"]),
			p.resolveEditTextFromResult(db, toolDoc["result"]),
			codeBlockContent(bubbleDoc, filePath, repoRoot),
		)
		oldText := p.resolveBeforeTextFromResult(db, toolDoc["result"])
		return []capture.SessionEvent{{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			OldText:   oldText,
			NewText:   newText,
			Raw:       cloneRaw(toolDoc),
		}}
	case "write":
		newText := codeBlockContent(bubbleDoc, filePath, repoRoot)
		if newText == "" {
			newText = rawString(params["contents"])
		}
		return []capture.SessionEvent{{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			NewText:   newText,
			Raw:       cloneRaw(toolDoc),
		}}
	case "apply_patch":
		patch := patchFromRawArgs(rawString(toolDoc["rawArgs"]))
		if patch == "" {
			patch = codeBlockContent(bubbleDoc, filePath, repoRoot)
		}
		return parsePatchText(patch, sessionID, model, ts, repoRoot, toolDoc)
	case "multiedit":
		if edits := parseMultiEditRawArgs(rawString(toolDoc["rawArgs"]), sessionID, model, ts, repoRoot, toolDoc); len(edits) > 0 {
			return edits
		}
		content := codeBlockContent(bubbleDoc, filePath, repoRoot)
		if filePath != "" && content != "" {
			return []capture.SessionEvent{{
				SessionID: sessionID,
				Tool:      capture.ToolCursor,
				Model:     model,
				TS:        ts,
				Kind:      capture.KindEdit,
				FilePath:  filePath,
				NewText:   content,
				Raw:       cloneRaw(toolDoc),
			}}
		}
	}

	return []capture.SessionEvent{{
		SessionID: sessionID,
		Tool:      capture.ToolCursor,
		Model:     model,
		TS:        ts,
		Kind:      capture.KindToolCall,
		FilePath:  filePath,
		Raw:       cloneRaw(toolDoc),
	}}
}

func parseCodeBlockEdits(bubbleDoc map[string]json.RawMessage, sessionID, model string, ts int64, repoRoot string) []capture.SessionEvent {
	var blocks []struct {
		URI struct {
			Path string `json:"path"`
		} `json:"uri"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(bubbleDoc["codeBlocks"], &blocks); err != nil || len(blocks) == 0 {
		return nil
	}
	var events []capture.SessionEvent
	for _, block := range blocks {
		path := relPath(block.URI.Path, repoRoot)
		if path == "" || strings.TrimSpace(block.Content) == "" {
			continue
		}
		events = append(events, capture.SessionEvent{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  path,
			NewText:   block.Content,
		})
	}
	return events
}

var patchFileRE = regexp.MustCompile(`(?m)^\*\*\* (?:Add|Update) File: (.+)$`)

func parsePatchText(text, sessionID, model string, ts int64, repoRoot string, raw map[string]json.RawMessage) []capture.SessionEvent {
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
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  relPath(currentFile, repoRoot),
			NewText:   strings.Join(added, "\n"),
			Raw:       cloneRaw(raw),
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

func parseMultiEditRawArgs(rawArgs, sessionID, model string, ts int64, repoRoot string, toolDoc map[string]json.RawMessage) []capture.SessionEvent {
	if strings.TrimSpace(rawArgs) == "" {
		return nil
	}
	var doc struct {
		FilePath string `json:"file_path"`
		Edits    []struct {
			NewString string `json:"new_string"`
			OldString string `json:"old_string"`
		} `json:"edits"`
	}
	if err := json.Unmarshal([]byte(rawArgs), &doc); err != nil {
		return nil
	}
	filePath := relPath(doc.FilePath, repoRoot)
	var events []capture.SessionEvent
	for _, edit := range doc.Edits {
		if strings.TrimSpace(edit.NewString) == "" {
			continue
		}
		events = append(events, capture.SessionEvent{
			SessionID: sessionID,
			Tool:      capture.ToolCursor,
			Model:     model,
			TS:        ts,
			Kind:      capture.KindEdit,
			FilePath:  filePath,
			OldText:   edit.OldString,
			NewText:   edit.NewString,
			Raw:       cloneRaw(toolDoc),
		})
	}
	return events
}

func patchFromRawArgs(rawArgs string) string {
	if strings.TrimSpace(rawArgs) == "" {
		return ""
	}
	var doc struct {
		Patch string `json:"patch"`
	}
	if err := json.Unmarshal([]byte(rawArgs), &doc); err != nil {
		return ""
	}
	return doc.Patch
}

func filePathFromRawArgs(rawArgs string) string {
	if strings.TrimSpace(rawArgs) == "" {
		return ""
	}
	var doc struct {
		FilePath string `json:"file_path"`
	}
	_ = json.Unmarshal([]byte(rawArgs), &doc)
	return doc.FilePath
}

func codeBlockContent(bubbleDoc map[string]json.RawMessage, hintPath, repoRoot string) string {
	var blocks []struct {
		URI struct {
			Path string `json:"path"`
		} `json:"uri"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(bubbleDoc["codeBlocks"], &blocks); err != nil {
		return ""
	}
	absHint, _ := filepath.Abs(hintPath)
	for _, block := range blocks {
		if hintPath != "" {
			blockPath := relPath(block.URI.Path, repoRoot)
			if blockPath != hintPath && block.URI.Path != hintPath {
				if absHint != "" {
					absBlock, err := filepath.Abs(block.URI.Path)
					if err != nil || absBlock != absHint {
						continue
					}
				} else {
					continue
				}
			}
		}
		if strings.TrimSpace(block.Content) != "" {
			return block.Content
		}
	}
	if len(blocks) == 1 {
		return blocks[0].Content
	}
	return ""
}

func recordBubblePaths(inv *capture.InventoryCollector, raw []byte, prefix string) {
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(raw, &doc); err != nil {
		return
	}
	recordShallow(inv, doc, prefix)
}

func recordShallow(inv *capture.InventoryCollector, obj map[string]json.RawMessage, prefix string) {
	if inv == nil || obj == nil {
		return
	}
	for key, val := range obj {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		inv.RecordPath(capture.ToolCursor, path, inferJSONType(val), mapField(path), "")
	}
}

func mapField(path string) string {
	switch {
	case strings.Contains(path, "composerId"):
		return "session_id"
	case strings.Contains(path, "createdAt"):
		return "ts"
	case strings.Contains(path, "relativeWorkspacePath"):
		return "file_path"
	case strings.Contains(path, "newString"):
		return "new_text"
	case strings.Contains(path, "oldString"):
		return "old_text"
	case strings.Contains(path, "contents"):
		return "new_text"
	case strings.Contains(path, "streamingContent"):
		return "new_text"
	case strings.Contains(path, "content") && strings.Contains(path, "codeBlocks"):
		return "new_text"
	case strings.Contains(path, "name") && strings.Contains(path, "toolFormerData"):
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

func parseBubbleKey(key string) (composerID, bubbleID string, ok bool) {
	rest := strings.TrimPrefix(key, "bubbleId:")
	if rest == key {
		return "", "", false
	}
	sep := strings.IndexByte(rest, ':')
	if sep <= 0 || sep == len(rest)-1 {
		return "", "", false
	}
	return rest[:sep], rest[sep+1:], true
}

func roleFromType(t int) string {
	switch t {
	case 1:
		return "user"
	case 2:
		return "assistant"
	default:
		return "unknown"
	}
}

func parseBubbleTime(createdAt string, sendMS, endMS int64) int64 {
	if sendMS > 0 {
		return sendMS
	}
	if endMS > 0 {
		return endMS
	}
	if createdAt != "" {
		layouts := []string{time.RFC3339Nano, time.RFC3339}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, createdAt); err == nil {
				return t.UnixMilli()
			}
		}
	}
	return 0
}

func composerLastSeen(meta composerMeta, bs []bubble) int64 {
	latest := meta.lastUpdatedAt
	for _, b := range bs {
		if b.createdAt > latest {
			latest = b.createdAt
		}
	}
	return latest
}

func inTimeWindow(tsMS int64, since, until time.Time) bool {
	if tsMS <= 0 || since.IsZero() && until.IsZero() {
		return true
	}
	if !since.IsZero() && tsMS < since.UnixMilli() {
		return false
	}
	if !until.IsZero() && tsMS > until.UnixMilli() {
		return false
	}
	return true
}

func bubblesInWindow(bs []bubble, since, until time.Time) bool {
	for _, b := range bs {
		if inTimeWindow(b.createdAt, since, until) {
			return true
		}
	}
	return false
}

func loadWorkspaceFolders() map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	root := filepath.Join(home, "Library", "Application Support", "Cursor", "User", "workspaceStorage")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	out := map[string]string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, entry.Name(), "workspace.json"))
		if err != nil {
			continue
		}
		var doc struct {
			Folder string `json:"folder"`
		}
		if err := json.Unmarshal(data, &doc); err != nil || doc.Folder == "" {
			continue
		}
		path := uriToPath(doc.Folder)
		if !isMeaningfulFolder(path) {
			continue
		}
		out[entry.Name()] = path
	}
	return out
}

func pickWorkspace(workspaces map[string]string, bs []bubble) (cwd, repoRoot string) {
	if len(workspaces) == 0 {
		return "", ""
	}
	counts := map[string]int{}
	for _, b := range bs {
		for _, folder := range workspaces {
			if folder != "" && strings.Contains(string(b.raw), folder) {
				counts[folder]++
			}
		}
	}
	var best string
	bestN := 0
	for folder, n := range counts {
		if n > bestN {
			best = folder
			bestN = n
		}
	}
	if best == "" {
		return "", ""
	}
	return best, detectRepoRoot(best)
}

func detectRepoRoot(cwd string) string {
	for dir := cwd; dir != "" && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		if exists(filepath.Join(dir, ".git")) || exists(filepath.Join(dir, ".jj")) {
			return dir
		}
	}
	return ""
}

func uriToPath(folderURI string) string {
	if strings.HasPrefix(folderURI, "file://") {
		return strings.TrimPrefix(folderURI, "file://")
	}
	return folderURI
}

func isMeaningfulFolder(path string) bool {
	trimmed := strings.TrimRight(path, "/")
	if trimmed == "" {
		return false
	}
	if strings.Count(trimmed, "/") < 2 {
		return false
	}
	return true
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil {
		return int64(f)
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

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func normalizeToolName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.TrimSuffix(name, "_v2")
	return name
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func (p *Parser) resolveEditTextFromResult(db *sql.DB, raw json.RawMessage) string {
	beforeID, afterID := contentIDsFromResult(raw)
	if afterID == "" {
		return ""
	}
	before := p.loadComposerContent(db, beforeID)
	after := p.loadComposerContent(db, afterID)
	return diffAddedText(before, after)
}

func (p *Parser) resolveBeforeTextFromResult(db *sql.DB, raw json.RawMessage) string {
	beforeID, _ := contentIDsFromResult(raw)
	return p.loadComposerContent(db, beforeID)
}

func contentIDsFromResult(raw json.RawMessage) (beforeID, afterID string) {
	if len(raw) == 0 {
		return "", ""
	}
	var result struct {
		BeforeContentID string `json:"beforeContentId"`
		AfterContentID  string `json:"afterContentId"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		var resultStr string
		if err2 := json.Unmarshal(raw, &resultStr); err2 == nil {
			_ = json.Unmarshal([]byte(resultStr), &result)
		}
	}
	return result.BeforeContentID, result.AfterContentID
}

func (p *Parser) loadComposerContent(db *sql.DB, contentID string) string {
	contentID = strings.TrimSpace(contentID)
	if contentID == "" {
		return ""
	}
	if p.contentCache != nil {
		if cached, ok := p.contentCache[contentID]; ok {
			return cached
		}
	}
	key := contentID
	if !strings.HasPrefix(key, "composer.content.") {
		key = "composer.content." + strings.TrimPrefix(contentID, "composer.content.")
	}
	var value []byte
	err := db.QueryRowContext(context.Background(), `SELECT value FROM cursorDiskKV WHERE key = ?`, key).Scan(&value)
	if err != nil {
		return ""
	}
	text := string(value)
	if p.contentCache != nil {
		p.contentCache[contentID] = text
	}
	return text
}

func diffAddedText(before, after string) string {
	if strings.TrimSpace(after) == "" {
		return ""
	}
	if strings.TrimSpace(before) == "" {
		return after
	}
	bLines := strings.Split(before, "\n")
	aLines := strings.Split(after, "\n")
	start := 0
	for start < len(bLines) && start < len(aLines) && bLines[start] == aLines[start] {
		start++
	}
	endB, endA := len(bLines)-1, len(aLines)-1
	for endB >= start && endA >= start && bLines[endB] == aLines[endA] {
		endB--
		endA--
	}
	if start > endA {
		return ""
	}
	return strings.Join(aLines[start:endA+1], "\n")
}

// DefaultVSCDBPath returns Cursor's global state.vscdb on macOS.
func DefaultVSCDBPath() (string, error) {
	return ingestcursor.DefaultVSCDBPath()
}

// DiscoverVSCDBPath resolves the vscdb path for session discovery.
func DiscoverVSCDBPath(home string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("cursor capture spike is only supported on macOS in WP-0")
	}
	_ = home
	return DefaultVSCDBPath()
}
