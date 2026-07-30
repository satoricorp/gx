package parsers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/totality/internal/capture"
)

func TestDiscoverFindsCursorAgentTranscripts(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "Users", "joe", "git", "tx")
	root := filepath.Join(home, ".cursor", "projects", "Users-joe-git-tx", "agent-transcripts", "parent-1")
	mainPath := filepath.Join(root, "parent-1.jsonl")
	subPath := filepath.Join(root, "subagents", "sub-1.jsonl")
	oldPath := filepath.Join(home, ".cursor", "projects", "Users-joe-git-tx", "agent-transcripts", "old", "old.jsonl")
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

	discovery, err := Discover(DiscoverOptions{
		HomeDir:     home,
		RepoRoot:    repoRoot,
		Since:       now.Add(-time.Hour),
		Until:       now.Add(time.Hour),
		CursorVSCDB: filepath.Join(home, "missing-state.vscdb"),
		Tools:       []string{capture.ToolCursor},
	})
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	sessions := discovery.Sessions
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

func TestDiscoverFindsClaudeSessionsFiledUnderAnotherProject(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "totality-xrepo", "repo")
	sibling := filepath.Join(string(filepath.Separator), "totality-xrepo", "repo-cloud")
	other := filepath.Join(string(filepath.Separator), "totality-xrepo", "elsewhere")
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

	discovery, err := Discover(DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-time.Hour),
		Until:    now.Add(time.Hour),
		Tools:    []string{capture.ToolClaude},
	})
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	sessions := discovery.Sessions
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

func TestDiscoverDeduplicatesClaudeSessionsAcrossProjects(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "totality-xrepo", "repo")
	other := filepath.Join(string(filepath.Separator), "totality-xrepo", "elsewhere")
	projects := filepath.Join(home, ".claude", "projects")

	now := time.Now()
	body := claudeEditLine("shared", filepath.Join(repoRoot, "a.go"))
	ownPath := filepath.Join(projects, repoSlug(repoRoot), "shared.jsonl")
	dupPath := filepath.Join(projects, repoSlug(other), "shared.jsonl")
	writeClaudeTranscript(t, ownPath, body, now)
	writeClaudeTranscript(t, dupPath, body, now)

	discovery, err := Discover(DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-time.Hour),
		Until:    now.Add(time.Hour),
		Tools:    []string{capture.ToolClaude},
	})
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	sessions := discovery.Sessions
	if len(sessions) != 1 {
		t.Fatalf("discovered %v, want only the repo-slug copy", discoveredPaths(sessions))
	}
	if sessions[0].Path != ownPath {
		t.Fatalf("discovered %q, want %q", sessions[0].Path, ownPath)
	}
}

func TestFileMentionsRepoSpansChunkBoundary(t *testing.T) {
	dir := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "totality-xrepo", "repo")
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

func TestDiscoverFindsClaudeSubagentTranscripts(t *testing.T) {
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

	discovery, err := Discover(DiscoverOptions{
		HomeDir:   home,
		RepoRoot:  repoRoot,
		ClaudeDir: claudeDir,
		Since:     time.Now().Add(-time.Hour),
		Until:     time.Now().Add(time.Hour),
		Tools:     []string{capture.ToolClaude},
	})
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}
	sessions := discovery.Sessions

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

// makeUnreadable creates dir and removes every permission bit, so reading it
// fails the way a real permission-denied transcript directory does. Root
// bypasses mode bits, so the caller is skipped there rather than lied to.
func makeUnreadable(t *testing.T, dir string) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions, so unreadable directories cannot be simulated")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
	if err := os.Chmod(dir, 0); err != nil {
		t.Fatalf("chmod %s: %v", dir, err)
	}
	// TempDir cleanup has to be able to walk back in.
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	if _, err := os.ReadDir(dir); err == nil {
		t.Skipf("%s is still readable despite mode 0; filesystem does not enforce permissions", dir)
	}
}

// TestDiscoverSurvivesUnreadableCursorTranscripts is the "one bad leg must not
// kill the others" regression: a single tool's unavailability used to abort the
// whole sweep, so Claude and Codex transcripts were never discovered either.
func TestDiscoverSurvivesUnreadableCursorTranscripts(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "totality-nocursor", "repo")
	now := time.Now()

	claudePath := filepath.Join(home, ".claude", "projects", repoSlug(repoRoot), "claude-1.jsonl")
	writeClaudeTranscript(t, claudePath, claudeEditLine("claude-1", filepath.Join(repoRoot, "a.go")), now)
	codexPath := filepath.Join(home, ".codex", "sessions", "2026", "07", "rollout-2026-07-26.jsonl")
	writeClaudeTranscript(t, codexPath, `{"cwd":"`+filepath.ToSlash(repoRoot)+`"}`+"\n", now)
	cursorDir := filepath.Join(home, ".cursor", "projects", cursorRepoSlug(repoRoot), "agent-transcripts")
	makeUnreadable(t, cursorDir)

	opts := DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-time.Hour),
		Until:    now.Add(time.Hour),
	}
	discovery, err := Discover(opts)
	if err != nil {
		t.Fatalf("Discover() error = %v, want nil: one unavailable tool must not kill the others", err)
	}
	got := map[string]string{}
	for _, session := range discovery.Sessions {
		got[session.Tool] = session.Path
	}
	if got[capture.ToolClaude] != claudePath {
		t.Fatalf("claude source = %q, want %q (all: %v)", got[capture.ToolClaude], claudePath, discoveredPaths(discovery.Sessions))
	}
	if got[capture.ToolCodex] != codexPath {
		t.Fatalf("codex source = %q, want %q (all: %v)", got[capture.ToolCodex], codexPath, discoveredPaths(discovery.Sessions))
	}

	// The failure is reported, not swallowed.
	if len(discovery.Problems) != 1 || discovery.Problems[0].Tool != capture.ToolCursor {
		t.Fatalf("problems = %v, want one cursor problem", discovery.ProblemStrings())
	}
	if !strings.Contains(discovery.Problems[0].String(), cursorDir) {
		t.Fatalf("problem = %q, want it to name the unreadable directory %q", discovery.Problems[0].String(), cursorDir)
	}

}

// TestDiscoverReportsUnreadableClaudeAndCodexRoots is the finding this whole
// reporting surface exists for. Claude and Codex hold essentially all the
// transcripts, and both legs used to be structurally incapable of returning an
// error: unreadable roots produced zero sessions, zero problems and a nil
// error, i.e. a sweep that looked exactly like a clean run over an empty
// machine. Every leg failing is additionally a hard error, which is only
// reachable because the legs can fail at all.
func TestDiscoverReportsUnreadableClaudeAndCodexRoots(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "totality-unreadable", "repo")
	now := time.Now()

	claudeDir := filepath.Join(home, ".claude", "projects")
	codexDir := filepath.Join(home, ".codex", "sessions")
	cursorDir := filepath.Join(home, ".cursor", "projects", cursorRepoSlug(repoRoot), "agent-transcripts")
	makeUnreadable(t, claudeDir)
	makeUnreadable(t, codexDir)
	makeUnreadable(t, cursorDir)

	discovery, err := Discover(DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-time.Hour),
		Until:    now.Add(time.Hour),
	})
	if len(discovery.Sessions) != 0 {
		t.Fatalf("sessions = %v, want none: nothing was readable", discoveredPaths(discovery.Sessions))
	}
	tools := map[string]bool{}
	for _, problem := range discovery.Problems {
		tools[problem.Tool] = true
	}
	for _, tool := range []string{capture.ToolClaude, capture.ToolCodex, capture.ToolCursor} {
		if !tools[tool] {
			t.Fatalf("problems = %v, want one for %s: an unreadable root must never read as 'this tool found nothing'",
				discovery.ProblemStrings(), tool)
		}
	}
	if err == nil {
		t.Fatalf("Discover() error = nil, want a hard error when every leg failed and nothing was discovered (problems: %v)",
			discovery.ProblemStrings())
	}
}

// TestDiscoverStaysInsideRequestedHome pins the hermeticity of a sweep. Cursor's
// database used to be resolved through os.UserHomeDir() regardless of the home
// discovery was asked about, so a temp HomeDir still returned — and the
// orchestrator then opened, inside the pre-push hook — the developer's real
// multi-gigabyte global Cursor database. A home with no agent tools installed
// must discover nothing, and must not warn about tools that are merely absent.
func TestDiscoverStaysInsideRequestedHome(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "totality-empty-home", "repo")
	now := time.Now()

	discovery, err := Discover(DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-time.Hour),
		Until:    now.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Discover() error = %v, want nil: nothing installed is not a failure", err)
	}
	for _, session := range discovery.Sessions {
		if !strings.HasPrefix(session.Path, home) {
			t.Fatalf("discovered %q, which is outside the requested home %q", session.Path, home)
		}
	}
	if len(discovery.Sessions) != 0 {
		t.Fatalf("sessions = %v, want none from an empty home", discoveredPaths(discovery.Sessions))
	}
	// A tool that is not installed is not a skipped tool: warning about it on
	// every push would be permanent, unactionable noise.
	if len(discovery.Problems) != 0 {
		t.Fatalf("problems = %v, want none: no agent tool is installed under this home", discovery.ProblemStrings())
	}
}

// TestDiscoverAcceptsMultiGigabyteCursorDatabase locks in that database size
// never gates discovery. A 256MB ceiling used to skip large databases because
// parsing began with a whole-file copy; the parser now reads the file in
// place with window-scoped index queries, so an established install's
// multi-gigabyte database — precisely the one holding the most history — must
// be discovered like any other, with no problem reported.
func TestDiscoverAcceptsMultiGigabyteCursorDatabase(t *testing.T) {
	home := t.TempDir()
	repoRoot := filepath.Join(string(filepath.Separator), "totality-bigcursor", "repo")
	now := time.Now()

	vscdb := filepath.Join(home, "state.vscdb")
	if err := os.WriteFile(vscdb, make([]byte, 4096), 0o600); err != nil {
		t.Fatal(err)
	}
	// Sparse-extend to 2GB: costs no real disk, but any reintroduced size
	// gate would trip on it.
	if err := os.Truncate(vscdb, 2<<30); err != nil {
		t.Fatal(err)
	}

	discovery, err := Discover(DiscoverOptions{
		HomeDir:     home,
		RepoRoot:    repoRoot,
		Since:       now.Add(-time.Hour),
		Until:       now.Add(time.Hour),
		CursorVSCDB: vscdb,
		Tools:       []string{capture.ToolCursor},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(discovery.Sessions) != 1 || discovery.Sessions[0].Kind != SessionKindCursorVSCDB {
		t.Fatalf("sessions = %v, want the multi-gigabyte database discovered", discoveredPaths(discovery.Sessions))
	}
	if len(discovery.Problems) != 0 {
		t.Fatalf("problems = %v, want none: size is not a reason to skip", discovery.ProblemStrings())
	}
}
