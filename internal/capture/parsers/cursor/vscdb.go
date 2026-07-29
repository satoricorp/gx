package cursor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/totality/internal/capture"
)

// unknownEndGrace is how long past createdAt a composer with no lastUpdatedAt
// is assumed to have stayed active. Cursor maintains lastUpdatedAt on live
// sessions, so rows without it are sessions that died young: measured across
// a real 1,703-composer database, no such composer had bubbles more than 1.0
// days after createdAt. A week of grace on top of the orchestrator's ±24h
// slack leaves no realistic way to prune a composer whose bubbles were in
// the window.
const unknownEndGrace = 7 * 24 * time.Hour

// metaMayOverlapWindow decides from composerData alone whether a composer
// could have bubbles inside [since, until]. It exists so bubble rows of
// sessions that ended long before the window are never read; whatever it
// keeps is still window-checked bubble-by-bubble by the caller. Unknown
// timestamps always pass, mirroring inTimeWindow. The one deliberate delta
// from the scan-everything era: a composer whose own metadata bounds its
// life strictly outside the window is pruned without reading its bubbles,
// so a stray timestamp-less bubble inside a long-dead session no longer
// resurrects it.
func metaMayOverlapWindow(meta composerMeta, since, until time.Time) bool {
	if since.IsZero() && until.IsZero() {
		return true
	}
	if !until.IsZero() && meta.createdAt > 0 && meta.createdAt > until.UnixMilli() {
		return false
	}
	end := meta.lastUpdatedAt
	if end <= 0 && meta.createdAt > 0 {
		end = meta.createdAt + unknownEndGrace.Milliseconds()
	}
	if !since.IsZero() && end > 0 && end < since.UnixMilli() {
		return false
	}
	return true
}

// openReadOnly opens the live database without copying it. mode=ro plus a
// busy timeout is the WAL-safe way to read while Cursor is running: readers
// take shared locks and see committed WAL content without ever blocking the
// writer. The previous DSN's immutable=1 asserts the file cannot change —
// true only of a private copy — which is what forced the whole-file copy
// this parser no longer does.
func openReadOnly(path string) (*sql.DB, error) {
	dsn := "file:" + path + "?mode=ro&_pragma=busy_timeout(5000)"
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
	// Range predicates instead of LIKE: the key column's UNIQUE index serves
	// prefix ranges, while LIKE runs a full table scan (the default
	// case-insensitive LIKE cannot use the index). ';' is ':'+1, so the
	// half-open range covers exactly the 'composerData:' prefix.
	rows, err := db.QueryContext(ctx, `SELECT key, value FROM cursorDiskKV WHERE key >= 'composerData:' AND key < 'composerData;'`)
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

// bubbleComposerIDs enumerates every composer that has at least one bubble,
// in key order, reading only the key index — no bubble payloads. This is what
// lets window pruning happen before any large row is touched, and it also
// finds composers that have bubbles but no composerData row (16 of 1,236 in
// the database this was built against), which metadata-driven selection alone
// would miss.
func bubbleComposerIDs(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT key FROM cursorDiskKV WHERE key >= 'bubbleId:' AND key < 'bubbleId;'`)
	if err != nil {
		return nil, fmt.Errorf("query bubble keys: %w", err)
	}
	defer rows.Close()

	var out []string
	seen := map[string]struct{}{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		composerID, _, ok := parseBubbleKey(key)
		if !ok {
			continue
		}
		if _, dup := seen[composerID]; dup {
			continue
		}
		seen[composerID] = struct{}{}
		out = append(out, composerID)
	}
	return out, rows.Err()
}

// loadComposerBubbles fetches one composer's bubbles through the key index.
func loadComposerBubbles(ctx context.Context, db *sql.DB, composerID string) ([]bubble, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT key, value FROM cursorDiskKV WHERE key >= 'bubbleId:'||?||':' AND key < 'bubbleId:'||?||';'`,
		composerID, composerID)
	if err != nil {
		return nil, fmt.Errorf("query bubbles: %w", err)
	}
	defer rows.Close()

	var out []bubble
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		keyComposer, bubbleID, ok := parseBubbleKey(key)
		if !ok || keyComposer != composerID {
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
		out = append(out, bubble{
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
