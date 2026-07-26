package parsers

import (
	"os"
	"path/filepath"
	"strings"
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

// claudeEditLine renders one Claude transcript line whose edit targets filePath.
func claudeEditLine(sessionID, filePath string) string {
	return `{"type":"assistant","cwd":"/somewhere/else","sessionId":"` + sessionID +
		`","message":{"role":"assistant","model":"claude","content":[{"type":"tool_use","id":"t1","name":"Write","input":{"file_path":"` +
		filepath.ToSlash(filePath) + `","content":"package foo\n"}}]},"timestamp":"2025-06-01T10:01:00.000Z"}` + "\n"
}

func writeClaudeTranscript(t *testing.T, path, body string, mod time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatalf("chtimes %s: %v", path, err)
	}
}

func discoveredPaths(sessions []DiscoveredSession) []string {
	paths := make([]string, 0, len(sessions))
	for _, s := range sessions {
		paths = append(paths, s.Path)
	}
	return paths
}

func TestDiscoverSessionsFindsClaudeSessionsFiledUnderAnotherProject(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "gx-xrepo", "repo")
	sibling := filepath.Join(string(filepath.Separator), "gx-xrepo", "repo-cloud")
	other := filepath.Join(string(filepath.Separator), "gx-xrepo", "elsewhere")
	projects := filepath.Join(home, ".claude", "projects")
	ownDir := filepath.Join(projects, repoSlug(repoRoot))
	otherDir := filepath.Join(projects, repoSlug(other))

	now := time.Now()
	stale := now.Add(-30 * 24 * time.Hour)

	ownPath := filepath.Join(ownDir, "own.jsonl")
	crossPath := filepath.Join(otherDir, "cross.jsonl")
	unrelatedPath := filepath.Join(otherDir, "unrelated.jsonl")
	siblingPath := filepath.Join(otherDir, "sibling.jsonl")
	stalePath := filepath.Join(otherDir, "stale.jsonl")

	writeClaudeTranscript(t, ownPath, claudeEditLine("own", filepath.Join(repoRoot, "a.go")), now)
	writeClaudeTranscript(t, crossPath, claudeEditLine("cross", filepath.Join(repoRoot, "internal", "b.go")), now.Add(-time.Minute))
	writeClaudeTranscript(t, unrelatedPath, claudeEditLine("unrelated", filepath.Join(other, "c.go")), now.Add(-time.Minute))
	writeClaudeTranscript(t, siblingPath, claudeEditLine("sibling", filepath.Join(sibling, "d.go")), now.Add(-time.Minute))
	writeClaudeTranscript(t, stalePath, claudeEditLine("stale", filepath.Join(repoRoot, "e.go")), stale)

	sessions, err := DiscoverSessions(DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-time.Hour),
		Until:    now.Add(time.Hour),
		Tools:    []string{capture.ToolClaude},
	})
	if err != nil {
		t.Fatalf("DiscoverSessions: %v", err)
	}
	got := discoveredPaths(sessions)
	want := []string{ownPath, crossPath}
	if len(got) != len(want) {
		t.Fatalf("discovered %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("discovered[%d] = %q, want %q (all: %v)", i, got[i], want[i], got)
		}
	}
	for _, s := range sessions {
		if s.Tool != capture.ToolClaude || s.Kind != SessionKindJSONL {
			t.Fatalf("session = %#v", s)
		}
	}
	if sessions[1].SessionID != "cross" {
		t.Fatalf("cross session id = %q", sessions[1].SessionID)
	}
}

func TestDiscoverSessionsDeduplicatesClaudeSessionsAcrossProjects(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "gx-xrepo", "repo")
	other := filepath.Join(string(filepath.Separator), "gx-xrepo", "elsewhere")
	projects := filepath.Join(home, ".claude", "projects")

	now := time.Now()
	body := claudeEditLine("shared", filepath.Join(repoRoot, "a.go"))
	ownPath := filepath.Join(projects, repoSlug(repoRoot), "shared.jsonl")
	dupPath := filepath.Join(projects, repoSlug(other), "shared.jsonl")
	writeClaudeTranscript(t, ownPath, body, now)
	writeClaudeTranscript(t, dupPath, body, now)

	sessions, err := DiscoverSessions(DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-time.Hour),
		Until:    now.Add(time.Hour),
		Tools:    []string{capture.ToolClaude},
	})
	if err != nil {
		t.Fatalf("DiscoverSessions: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("discovered %v, want only the repo-slug copy", discoveredPaths(sessions))
	}
	if sessions[0].Path != ownPath {
		t.Fatalf("discovered %q, want %q", sessions[0].Path, ownPath)
	}
}

func TestFileMentionsRepoSpansChunkBoundary(t *testing.T) {
	dir := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "gx-xrepo", "repo")
	needles := repoPathNeedles(repoRoot)
	if len(needles) == 0 {
		t.Fatal("no needles for repo root")
	}
	needle := string(needles[0])

	// The exact read-window edge depends on the needle-sized carry, so sweep
	// every byte offset around the first two nominal chunk boundaries: a
	// straddling needle must be found wherever the seam actually lands.
	for _, base := range []int{repoScanChunkSize, 2 * repoScanChunkSize} {
		for shift := -len(needle); shift <= 2*len(needle); shift++ {
			offset := base + shift
			hit := filepath.Join(dir, "hit.jsonl")
			writeClaudeTranscript(t, hit, strings.Repeat("x", offset)+needle+"main.go\n", time.Now())
			if !fileMentionsRepo(hit, needles) {
				t.Fatalf("needle at offset %d (chunk %d, shift %+d) was missed", offset, base, shift)
			}
		}
	}

	miss := filepath.Join(dir, "miss.jsonl")
	writeClaudeTranscript(t, miss, strings.Repeat("y", 3*repoScanChunkSize)+"\n", time.Now())
	if fileMentionsRepo(miss, needles) {
		t.Fatal("file without the repo path matched")
	}

	empty := filepath.Join(dir, "empty.jsonl")
	writeClaudeTranscript(t, empty, "", time.Now())
	if fileMentionsRepo(empty, needles) {
		t.Fatal("empty file matched")
	}
}

func TestDiscoverSessionsFindsClaudeSubagentTranscripts(t *testing.T) {
	// An agent that delegates its edits leaves most of the evidence in its
	// subagents' transcripts, not in the parent conversation's.
	home := t.TempDir()
	repoRoot := filepath.Join(home, "work", "repo")
	claudeDir := filepath.Join(home, ".claude", "projects")
	mod := time.Now().Add(-time.Minute)
	edit := claudeEditLine("s1", filepath.Join(repoRoot, "main.go"))

	// Parent conversation filed under an unrelated project directory.
	other := filepath.Join(claudeDir, "-somewhere-else")
	writeClaudeTranscript(t, filepath.Join(other, "conv-1.jsonl"), edit, mod)
	// Its subagents, where the edits actually happened.
	writeClaudeTranscript(t, filepath.Join(other, "conv-1", "subagents", "agent-a.jsonl"), edit, mod)
	writeClaudeTranscript(t, filepath.Join(other, "conv-1", "subagents", "agent-b.jsonl"), edit, mod)
	// A second conversation reusing a subagent name must not collide.
	writeClaudeTranscript(t, filepath.Join(other, "conv-2", "subagents", "agent-a.jsonl"), edit, mod)

	sessions, err := DiscoverSessions(DiscoverOptions{
		HomeDir:   home,
		RepoRoot:  repoRoot,
		ClaudeDir: claudeDir,
		Since:     time.Now().Add(-time.Hour),
		Until:     time.Now().Add(time.Hour),
		Tools:     []string{capture.ToolClaude},
	})
	if err != nil {
		t.Fatalf("DiscoverSessions() error = %v", err)
	}

	got := map[string]bool{}
	for _, session := range sessions {
		got[session.SessionID] = true
	}
	for _, want := range []string{
		"conv-1",
		"conv-1:subagent:agent-a",
		"conv-1:subagent:agent-b",
		"conv-2:subagent:agent-a",
	} {
		if !got[want] {
			t.Fatalf("session IDs = %#v, missing %q", got, want)
		}
	}
	if len(sessions) != 4 {
		t.Fatalf("len(sessions) = %d, want 4 with no duplicates", len(sessions))
	}
}
