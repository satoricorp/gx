package cursor_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/parsers/cursor"
)

func TestParseCursor_Golden(t *testing.T) {
	fixture := filepath.Join("..", "testdata", "cursor", "minimal.vscdb")
	parser := &cursor.Parser{
		Since: time.Unix(0, 0),
		Until: time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	events, err := parser.ParseFile(fixture, "/tmp/repo")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	var edits []capture.SessionEvent
	for _, ev := range events {
		if ev.IsEditEvent() {
			edits = append(edits, ev)
		}
	}
	if len(edits) != 2 {
		t.Fatalf("edit events = %d, want 2", len(edits))
	}
	if edits[0].SessionID != "cursor-c1" {
		t.Fatalf("session id = %q", edits[0].SessionID)
	}
	if edits[0].Tool != capture.ToolCursor {
		t.Fatalf("tool = %q", edits[0].Tool)
	}
	if edits[0].FilePath != "internal/foo.go" {
		t.Fatalf("search_replace path = %q", edits[0].FilePath)
	}
	if edits[0].NewText != "return 43" {
		t.Fatalf("search_replace newText = %q", edits[0].NewText)
	}
	if edits[1].FilePath != "README.md" || edits[1].NewText == "" {
		t.Fatalf("write edit mismatch: path=%q newText empty=%v", edits[1].FilePath, edits[1].NewText == "")
	}
}

func TestParseCursorTranscriptJSONL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "transcript-1.jsonl")
	line := `{"timestamp":"2026-07-03T20:00:00Z","role":"assistant","message":{"model":"gpt-test","content":[{"type":"tool_use","name":"StrReplace","input":{"path":"/tmp/repo/internal/foo.go","old_string":"return 1","new_string":"return 2"}},{"type":"tool_use","name":"Write","input":{"path":"/tmp/repo/README.md","content":"hello"}}]}}`
	if err := os.WriteFile(path, []byte(line+"\n"), 0o600); err != nil {
		t.Fatalf("write transcript: %v", err)
	}

	parser := &cursor.Parser{}
	events, err := parser.ParseFile(path, "/tmp/repo")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	var edits []capture.SessionEvent
	for _, ev := range events {
		if ev.IsEditEvent() {
			edits = append(edits, ev)
		}
	}
	if len(edits) != 2 {
		t.Fatalf("edit events = %d, want 2: %#v", len(edits), edits)
	}
	if edits[0].Tool != capture.ToolCursor || edits[0].SessionID != "cursor:transcript-1" {
		t.Fatalf("first edit provenance = %s/%s", edits[0].Tool, edits[0].SessionID)
	}
	if edits[0].FilePath != "internal/foo.go" || edits[0].NewText != "return 2" {
		t.Fatalf("strreplace edit = path %q new %q", edits[0].FilePath, edits[0].NewText)
	}
	if edits[1].FilePath != "README.md" || edits[1].NewText != "hello" {
		t.Fatalf("patch edit = path %q new %q", edits[1].FilePath, edits[1].NewText)
	}
}

func buildGoldenVSCDB(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "minimal.vscdb")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE cursorDiskKV (key TEXT UNIQUE ON CONFLICT REPLACE, value BLOB)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	rows := []struct {
		key   string
		value string
	}{
		{
			key:   "composerData:c1",
			value: `{"composerId":"c1","name":"gx spike","createdAt":1700000000000,"lastUpdatedAt":1700000005000}`,
		},
		{
			key: "bubbleId:c1:b1",
			value: `{
				"bubbleId":"b1","type":2,
				"createdAt":"2024-11-14T12:00:03.000Z",
				"toolFormerData":{
					"name":"search_replace","status":"completed",
					"params":"{\"relativeWorkspacePath\":\"internal/foo.go\",\"oldString\":\"return 42\",\"newString\":\"return 43\"}"
				}
			}`,
		},
		{
			key: "bubbleId:c1:b2",
			value: `{
				"bubbleId":"b2","type":2,
				"createdAt":"2024-11-14T12:00:04.000Z",
				"toolFormerData":{"name":"write","status":"completed","params":"{\"relativeWorkspacePath\":\"README.md\"}"},
				"codeBlocks":[{"uri":{"path":"/tmp/repo/README.md"},"content":"# hello\n"}]
			}`,
		},
	}
	for _, r := range rows {
		if _, err := db.Exec(`INSERT INTO cursorDiskKV(key, value) VALUES(?, ?)`, r.key, r.value); err != nil {
			t.Fatalf("insert %s: %v", r.key, err)
		}
	}
	return path
}

func TestBuildGoldenFixture(t *testing.T) {
	if testing.Short() {
		t.Skip("fixture generation only")
	}
	_ = buildGoldenVSCDB(t)
}

func editBubble(newString string) string {
	return `{
		"bubbleId":"b","type":2,
		"toolFormerData":{
			"name":"search_replace","status":"completed",
			"params":"{\"relativeWorkspacePath\":\"internal/foo.go\",\"oldString\":\"a\",\"newString\":\"` + newString + `\"}"
		}
	}`
}

// TestParseCursorPrunesByComposerMetadata pins the windowing rules that let
// the parser skip bubble rows of long-dead composers without reading them:
//   - metadata that bounds a composer's life strictly before the window
//     prunes it, even when a bubble carries no parseable timestamp (the one
//     deliberate delta from the scan-everything era, when such a bubble
//     resurrected the composer);
//   - a composer with no lastUpdatedAt stays eligible for a grace period
//     after createdAt rather than being treated as unbounded;
//   - a composer with bubbles but no composerData row is still parsed.
func TestParseCursorPrunesByComposerMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.vscdb")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE cursorDiskKV (key TEXT UNIQUE ON CONFLICT REPLACE, value BLOB)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	june2025 := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	jan8 := time.Date(2026, 1, 8, 0, 0, 0, 0, time.UTC).UnixMilli()
	rows := []struct{ key, value string }{
		// Ended months before the window; its bubble has no timestamp.
		{"composerData:dead", `{"composerId":"dead","createdAt":` + fmt.Sprint(june2025) + `,"lastUpdatedAt":` + fmt.Sprint(june2025+86400000) + `}`},
		{"bubbleId:dead:b1", editBubble("from the dead composer")},
		// No lastUpdatedAt; created 2 days before the window opens, so the
		// unknown-end grace must keep it, and its bubble is in the window.
		{"composerData:graced", `{"composerId":"graced","createdAt":` + fmt.Sprint(jan8) + `}`},
		{"bubbleId:graced:b1", `{"bubbleId":"b1","type":2,"createdAt":"2026-01-12T10:00:00.000Z",` +
			`"toolFormerData":{"name":"search_replace","status":"completed",` +
			`"params":"{\"relativeWorkspacePath\":\"internal/foo.go\",\"oldString\":\"a\",\"newString\":\"from the graced composer\"}"}}`},
		// Bubbles with no composerData row at all.
		{"bubbleId:orphan:b1", `{"bubbleId":"b1","type":2,"createdAt":"2026-01-12T11:00:00.000Z",` +
			`"toolFormerData":{"name":"search_replace","status":"completed",` +
			`"params":"{\"relativeWorkspacePath\":\"internal/foo.go\",\"oldString\":\"a\",\"newString\":\"from the orphan composer\"}"}}`},
	}
	for _, r := range rows {
		if _, err := db.Exec(`INSERT INTO cursorDiskKV(key, value) VALUES(?, ?)`, r.key, r.value); err != nil {
			t.Fatalf("insert %s: %v", r.key, err)
		}
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	parser := &cursor.Parser{
		Since: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		Until: time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
	}
	events, err := parser.ParseFile(path, "/tmp/repo")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	sessions := map[string]bool{}
	for _, ev := range events {
		sessions[ev.SessionID] = true
	}
	if sessions["cursor-dead"] {
		t.Fatalf("sessions = %v: composer whose metadata ends before the window must be pruned unread", sessions)
	}
	if !sessions["cursor-graced"] {
		t.Fatalf("sessions = %v: composer without lastUpdatedAt within the grace period must be kept", sessions)
	}
	if !sessions["cursor-orphan"] {
		t.Fatalf("sessions = %v: composer with bubbles but no composerData row must be kept", sessions)
	}
}

// TestParseCursorReadsLiveWALDatabase proves the no-copy path: rows that
// exist only in the -wal file of a database another connection still holds
// open — exactly the state Cursor leaves state.vscdb in while running — must
// be visible to ParseFile.
func TestParseCursorReadsLiveWALDatabase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.vscdb")
	writer, err := sql.Open("sqlite", "file:"+path+"?_pragma=journal_mode(WAL)")
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	defer writer.Close()

	if _, err := writer.Exec(`CREATE TABLE cursorDiskKV (key TEXT UNIQUE ON CONFLICT REPLACE, value BLOB)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	now := time.Now().UTC()
	inserts := []struct{ key, value string }{
		{"composerData:c1", `{"composerId":"c1","createdAt":` + fmt.Sprint(now.UnixMilli()) + `,"lastUpdatedAt":` + fmt.Sprint(now.UnixMilli()) + `}`},
		{"bubbleId:c1:b1", `{"bubbleId":"b1","type":2,"createdAt":"` + now.Format(time.RFC3339Nano) + `",` +
			`"toolFormerData":{"name":"search_replace","status":"completed",` +
			`"params":"{\"relativeWorkspacePath\":\"internal/foo.go\",\"oldString\":\"a\",\"newString\":\"b\"}"}}`},
	}
	for _, r := range inserts {
		if _, err := writer.Exec(`INSERT INTO cursorDiskKV(key, value) VALUES(?, ?)`, r.key, r.value); err != nil {
			t.Fatalf("insert %s: %v", r.key, err)
		}
	}
	if _, err := os.Stat(path + "-wal"); err != nil {
		t.Fatalf("expected a live -wal file, got %v: the fixture no longer reproduces Cursor's running state", err)
	}

	parser := &cursor.Parser{
		Since: now.Add(-time.Hour),
		Until: now.Add(time.Hour),
	}
	events, err := parser.ParseFile(path, "/tmp/repo")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	var edits []capture.SessionEvent
	for _, ev := range events {
		if ev.IsEditEvent() {
			edits = append(edits, ev)
		}
	}
	if len(edits) != 1 || edits[0].SessionID != "cursor-c1" || edits[0].NewText != "b" {
		t.Fatalf("edits = %#v, want the WAL-resident edit visible without copying", edits)
	}
}
