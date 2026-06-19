package claude_test

import (
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/parsers/claude"
)

func TestParseClaude_Golden(t *testing.T) {
	fixture := filepath.Join("..", "testdata", "claude", "v1-minimal.jsonl")
	parser := &claude.Parser{}
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
	if edits[0].SessionID != "sess-claude-001" {
		t.Fatalf("session id = %q", edits[0].SessionID)
	}
	if edits[0].Tool != capture.ToolClaude {
		t.Fatalf("tool = %q", edits[0].Tool)
	}
	if edits[0].FilePath != "internal/foo.go" {
		t.Fatalf("write path = %q", edits[0].FilePath)
	}
	if edits[0].NewText == "" {
		t.Fatal("write newText empty")
	}
	if edits[1].OldText != "return 42" || edits[1].NewText != "return 43" {
		t.Fatalf("edit text mismatch: old=%q new=%q", edits[1].OldText, edits[1].NewText)
	}
}
