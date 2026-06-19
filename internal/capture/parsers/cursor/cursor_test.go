package cursor_test

import (
	"database/sql"
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
