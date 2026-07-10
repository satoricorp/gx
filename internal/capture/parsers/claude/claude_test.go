package claude_test

import (
	"path/filepath"
	"strings"
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

func TestParseClaude_SkipsBadLines(t *testing.T) {
	parser := &claude.Parser{}
	raw := strings.Join([]string{
		`{"type":"assistant","message":{"role":"assistant","model":"claude","content":[{"type":"tool_use","name":"Write","input":{"file_path":"a.go","content":"ok"}}]},"timestamp":"2025-06-01T10:01:00.000Z","sessionId":"sess-1"}`,
		`not-json`,
		`{"type":"assistant","message":{"role":"assistant","model":"claude","content":[{"type":"tool_use","name":"Write","input":{"file_path":"b.go","content":"also-ok"}}]},"timestamp":"2025-06-01T10:02:00.000Z","sessionId":"sess-1"}`,
	}, "\n")
	events, err := parser.ParseBytes([]byte(raw), "sess-1.jsonl", "/tmp/repo")
	if err != nil {
		t.Fatalf("ParseBytes: %v", err)
	}
	if parser.LastBadLines != 1 {
		t.Fatalf("bad lines = %d, want 1", parser.LastBadLines)
	}
	var edits int
	for _, ev := range events {
		if ev.IsEditEvent() {
			edits++
		}
	}
	if edits != 2 {
		t.Fatalf("edit events = %d, want 2", edits)
	}
}
