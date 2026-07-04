package parsers

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/capture"
)

func TestDiscoverSessionsFindsCursorAgentTranscripts(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "Users", "joe", "git", "gx")
	root := filepath.Join(home, ".cursor", "projects", "Users-joe-git-gx", "agent-transcripts", "parent-1")
	mainPath := filepath.Join(root, "parent-1.jsonl")
	subPath := filepath.Join(root, "subagents", "sub-1.jsonl")
	oldPath := filepath.Join(home, ".cursor", "projects", "Users-joe-git-gx", "agent-transcripts", "old", "old.jsonl")
	for _, path := range []string{mainPath, subPath, oldPath} {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte("{}\n"), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	now := time.Now()
	if err := os.Chtimes(mainPath, now, now); err != nil {
		t.Fatalf("chtimes main: %v", err)
	}
	if err := os.Chtimes(subPath, now.Add(-time.Minute), now.Add(-time.Minute)); err != nil {
		t.Fatalf("chtimes sub: %v", err)
	}
	old := now.Add(-30 * 24 * time.Hour)
	if err := os.Chtimes(oldPath, old, old); err != nil {
		t.Fatalf("chtimes old: %v", err)
	}

	sessions, err := DiscoverSessions(DiscoverOptions{
		HomeDir:     home,
		RepoRoot:    repoRoot,
		Since:       now.Add(-time.Hour),
		Until:       now.Add(time.Hour),
		CursorVSCDB: filepath.Join(home, "missing-state.vscdb"),
		Tools:       []string{capture.ToolCursor},
	})
	if err != nil {
		t.Fatalf("DiscoverSessions: %v", err)
	}
	if len(sessions) != 2 {
		t.Fatalf("sessions = %#v, want main + subagent", sessions)
	}
	if sessions[0].Path != mainPath || sessions[0].Kind != SessionKindJSONL || sessions[0].SessionID != "cursor:parent-1" {
		t.Fatalf("main session = %#v", sessions[0])
	}
	if sessions[1].Path != subPath || sessions[1].SessionID != "cursor:parent-1:subagent:sub-1" {
		t.Fatalf("subagent session = %#v", sessions[1])
	}
}
