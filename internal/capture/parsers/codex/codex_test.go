package codex_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/lgtm/internal/capture"
	"github.com/satoricorp/lgtm/internal/capture/parsers/codex"
)

func TestParseCodex_Golden(t *testing.T) {
	fixture := filepath.Join("..", "testdata", "codex", "v1-minimal.jsonl")
	parser := &codex.Parser{}
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
	if edits[0].SessionID != "sess-codex-001" {
		t.Fatalf("session id = %q", edits[0].SessionID)
	}
	if edits[0].Tool != capture.ToolCodex {
		t.Fatalf("tool = %q", edits[0].Tool)
	}
	if edits[0].FilePath != "internal/bar.go" {
		t.Fatalf("patch path = %q", edits[0].FilePath)
	}
	if !strings.Contains(edits[0].NewText, "func Baz") {
		t.Fatalf("patch content missing Baz: %q", edits[0].NewText)
	}
	if edits[1].FilePath != "internal/baz.go" {
		t.Fatalf("write path = %q", edits[1].FilePath)
	}
}
