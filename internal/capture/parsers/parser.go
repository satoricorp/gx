package parsers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture"
)

// Parser normalizes one agent tool's session JSONL into SessionEvents.
type Parser interface {
	Tool() string
	ParseFile(path string, repoRoot string) ([]capture.SessionEvent, error)
}

// DiscoveredSession is one parseable session source (JSONL file or Cursor vscdb).
type DiscoveredSession struct {
	Tool      string
	Path      string
	Kind      string
	SessionID string
}

// SourceSessionID derives the per-source session identity from a transcript
// path, matching the IDs DiscoverSessions assigns. Identity must come from the
// path, not the transcript contents: Claude stamps subagent transcripts with
// the parent conversation's sessionId, so the in-file ID cannot distinguish
// the parent file from the subagent files beside it.
func SourceSessionID(tool, path string) string {
	stem := strings.TrimSuffix(filepath.Base(path), ".jsonl")
	switch tool {
	case capture.ToolCursor:
		return cursorTranscriptSessionID(path)
	case capture.ToolClaude:
		if filepath.Base(filepath.Dir(path)) == "subagents" {
			// Namespaced by conversation so two conversations cannot collide
			// on a subagent name and silently drop one another's evidence.
			conversation := filepath.Base(filepath.Dir(filepath.Dir(path)))
			if conversation != "" && conversation != "." {
				return conversation + ":subagent:" + stem
			}
		}
	}
	return stem
}

const (
	SessionKindJSONL       = "jsonl"
	SessionKindCursorVSCDB = "cursor_vscdb"
)

// DiscoverOptions controls session file discovery.
type DiscoverOptions struct {
	HomeDir     string
	RepoRoot    string
	Since       time.Time
	Until       time.Time
	ClaudeDir   string
	CodexDir    string
	CursorVSCDB string
	Tools       []string
}

// ToolProblem records one discovery leg that could not contribute its sources.
//
// A problem is reported, never returned as a fatal error: one unavailable tool
// must not blind capture to the others. Hard-failing on a single leg made
// discovery — and therefore all of capture — return nothing at all for the
// tools that were perfectly readable.
//
// A problem means evidence was skipped, so it is raised only when a source
// exists and could not be read. A tool that is simply not installed produces no
// problem: warning about it on every push would be permanent noise carrying no
// information, which is the same failure as silence.
type ToolProblem struct {
	Tool string
	Err  error
}

func (p ToolProblem) String() string {
	if p.Err == nil {
		return p.Tool
	}
	return p.Tool + ": " + p.Err.Error()
}

// Discovery is one discovery sweep: the sources found, plus the per-tool
// problems that kept a leg from contributing everything it could. Callers are
// expected to surface Problems rather than drop them.
type Discovery struct {
	Sessions []DiscoveredSession
	Problems []ToolProblem
}

// ProblemStrings renders the problems for user-facing warnings.
func (d Discovery) ProblemStrings() []string {
	if len(d.Problems) == 0 {
		return nil
	}
	out := make([]string, 0, len(d.Problems))
	for _, problem := range d.Problems {
		out = append(out, problem.String())
	}
	return out
}

// cursorVSCDBPath locates Cursor's global state.vscdb *inside* the home
// directory discovery was asked to sweep.
//
// Resolution must stay under opts.HomeDir. The previous resolver ignored its
// home argument and called os.UserHomeDir(), so a hermetic HomeDir — a test's
// temp dir, or another user's home — still resolved to the developer's own
// global Cursor database, which the orchestrator then opened and scanned inside
// the pre-push hook. That database is routinely multi-gigabyte.
//
// The path is returned unconditionally on every platform; whether Cursor is
// actually installed is decided by stat'ing it, not by GOOS. A tool that is not
// installed contributes nothing and is not a problem worth warning about.
func cursorVSCDBPath(home string) string {
	if strings.TrimSpace(home) == "" {
		return ""
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Cursor", "User", "globalStorage", "state.vscdb")
	case "windows":
		return filepath.Join(home, "AppData", "Roaming", "Cursor", "User", "globalStorage", "state.vscdb")
	default:
		return filepath.Join(home, ".config", "Cursor", "User", "globalStorage", "state.vscdb")
	}
}

// unreadable turns a filesystem error into a discovery problem, or into nil
// when there is nothing to report.
//
// "This tool is not installed" (the directory does not exist) is not a problem:
// there is no evidence being skipped. "This tool is installed but I could not
// read it" is exactly the case that must never render as a clean empty sweep,
// because real transcripts are being dropped.
func unreadable(path string, err error) error {
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return fmt.Errorf("read %s: %w", path, err)
}

// Discover runs every requested tool's discovery leg independently. A leg that
// cannot run contributes nothing and records a problem; the others still
// return their sources. Only a sweep where every attempted leg failed and
// nothing at all was discovered is a hard error.
func Discover(opts DiscoverOptions) (Discovery, error) {
	if opts.HomeDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return Discovery{}, err
		}
		opts.HomeDir = home
	}
	tools := normalizeTools(opts.Tools)

	var discovery Discovery
	attempted, failed := 0, 0
	leg := func(tool string, discover func() ([]DiscoveredSession, error)) {
		if !containsTool(tools, tool) {
			return
		}
		attempted++
		sessions, err := discover()
		discovery.Sessions = append(discovery.Sessions, sessions...)
		if err == nil {
			return
		}
		discovery.Problems = append(discovery.Problems, ToolProblem{Tool: tool, Err: err})
		if len(sessions) == 0 {
			failed++
		}
	}

	leg(capture.ToolClaude, func() ([]DiscoveredSession, error) {
		claudeDir := opts.ClaudeDir
		if claudeDir == "" {
			claudeDir = filepath.Join(opts.HomeDir, ".claude", "projects")
		}
		return discoverClaudeSessions(claudeDir, opts.RepoRoot, opts.Since, opts.Until)
	})
	leg(capture.ToolCodex, func() ([]DiscoveredSession, error) { return discoverCodexLeg(opts) })
	leg(capture.ToolCursor, func() ([]DiscoveredSession, error) { return discoverCursorLeg(opts) })

	if attempted > 0 && failed == attempted && len(discovery.Sessions) == 0 {
		return discovery, fmt.Errorf("session discovery failed for every tool: %s",
			strings.Join(discovery.ProblemStrings(), "; "))
	}
	return discovery, nil
}

func discoverCodexLeg(opts DiscoverOptions) ([]DiscoveredSession, error) {
	codexDir := opts.CodexDir
	if codexDir == "" {
		codexDir = filepath.Join(opts.HomeDir, ".codex", "sessions")
	}
	// A directory-walk failure is partial by nature: whatever the walk reached
	// before the error is still real evidence, so it is returned alongside it.
	paths, err := discoverCodexSessions(codexDir, opts.Since, opts.Until, opts.RepoRoot)
	sessions := make([]DiscoveredSession, 0, len(paths))
	for _, path := range paths {
		sessions = append(sessions, DiscoveredSession{
			Tool:      capture.ToolCodex,
			Path:      path,
			Kind:      SessionKindJSONL,
			SessionID: strings.TrimSuffix(filepath.Base(path), ".jsonl"),
		})
	}
	if err != nil {
		return sessions, fmt.Errorf("walk %s: %w", codexDir, err)
	}
	return sessions, nil
}

// discoverCursorLeg has two independent sources — per-repo agent transcripts
// and Cursor's global state.vscdb — and either can fail without the other.
func discoverCursorLeg(opts DiscoverOptions) ([]DiscoveredSession, error) {
	var sessions []DiscoveredSession
	var problems []error

	paths, err := discoverCursorTranscriptSessions(opts.HomeDir, opts.RepoRoot, opts.Since, opts.Until)
	for _, path := range paths {
		sessions = append(sessions, DiscoveredSession{
			Tool:      capture.ToolCursor,
			Path:      path,
			Kind:      SessionKindJSONL,
			SessionID: cursorTranscriptSessionID(path),
		})
	}
	if err != nil {
		problems = append(problems, fmt.Errorf("agent transcripts: %w", err))
	}

	vscdb := opts.CursorVSCDB
	if vscdb == "" {
		vscdb = cursorVSCDBPath(opts.HomeDir)
	}
	if vscdb != "" {
		_, err := os.Stat(vscdb)
		switch {
		case err == nil:
			// Size is no cost concern: the parser reads the database in place
			// with window-scoped index queries, so a multi-gigabyte install is
			// discovered like any other. (A 256MB ceiling used to skip large
			// databases here, back when parsing began with a whole-file copy.)
			sessions = append(sessions, DiscoveredSession{
				Tool:      capture.ToolCursor,
				Path:      vscdb,
				Kind:      SessionKindCursorVSCDB,
				SessionID: "state.vscdb",
			})
		default:
			// Absent means Cursor is not installed under this home, which is
			// not worth a warning on every push. Anything else means the
			// database is there and unreadable, which is.
			if problem := unreadable(vscdb, err); problem != nil {
				problems = append(problems, fmt.Errorf("state.vscdb: %w", problem))
			}
		}
	}
	return sessions, errors.Join(problems...)
}

// discoverClaudeSessions returns Claude transcripts that plausibly touched repoRoot.
//
// Claude Code files a transcript under the slug of the session's cwd, not the
// repo it edited, so an agent run from another checkout that edits this repo
// lands in a different project directory. The repo's own slug directory stays
// the fast path; every other project directory is scanned with a cheap
// substring pre-filter so cross-repo sessions still reach the matcher, which
// remains the layer that decides real attribution.
//
// A directory that exists but cannot be read is returned as an error rather
// than as an empty result: Claude holds most of the transcripts, so an
// unreadable ~/.claude/projects silently reduces capture to nothing while every
// summary still reads like a clean run.
func discoverClaudeSessions(claudeDir, repoRoot string, since, until time.Time) ([]DiscoveredSession, error) {
	primarySlug := repoSlug(repoRoot)
	seen := map[string]struct{}{}
	var problems []error

	sessions, err := claudeDirSessions(filepath.Join(claudeDir, primarySlug), since, until, nil, seen)
	if err != nil {
		problems = append(problems, err)
	}

	needles := repoPathNeedles(repoRoot)
	if len(needles) == 0 {
		return sessions, errors.Join(problems...)
	}
	entries, err := os.ReadDir(claudeDir)
	if err != nil {
		if problem := unreadable(claudeDir, err); problem != nil {
			problems = append(problems, problem)
		}
		return sessions, errors.Join(problems...)
	}
	var widened []DiscoveredSession
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == primarySlug {
			continue
		}
		found, err := claudeDirSessions(filepath.Join(claudeDir, entry.Name()), since, until, needles, seen)
		if err != nil {
			problems = append(problems, err)
		}
		widened = append(widened, found...)
	}
	sort.SliceStable(widened, func(i, j int) bool {
		return modTime(widened[i].Path).After(modTime(widened[j].Path))
	})
	return append(sessions, widened...), errors.Join(problems...)
}

// claudeDirSessions lists in-window transcripts in one project directory.
// A non-nil needles slice gates each file on mentioning the repo. seen
// de-duplicates by session ID so a transcript reachable from two directories
// is only reported once.
func claudeDirSessions(dir string, since, until time.Time, needles [][]byte, seen map[string]struct{}) ([]DiscoveredSession, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, unreadable(dir, err)
	}
	var sessions []DiscoveredSession
	var problems []error
	for _, entry := range entries {
		// A conversation's subagents write their own transcripts under
		// <conversation>/subagents, and an agent that delegates its edits
		// leaves most of the evidence there rather than in the parent.
		if entry.IsDir() {
			found, err := claudeSubagentSessions(
				filepath.Join(dir, entry.Name(), "subagents"),
				since, until, needles, seen,
			)
			if err != nil {
				problems = append(problems, err)
			}
			sessions = append(sessions, found...)
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if session, ok := claudeSessionAt(path, SourceSessionID(capture.ToolClaude, path), since, until, needles, seen); ok {
			sessions = append(sessions, session)
		}
	}
	return sessions, errors.Join(problems...)
}

func claudeSubagentSessions(dir string, since, until time.Time, needles [][]byte, seen map[string]struct{}) ([]DiscoveredSession, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, unreadable(dir, err)
	}
	var sessions []DiscoveredSession
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if session, ok := claudeSessionAt(path, SourceSessionID(capture.ToolClaude, path), since, until, needles, seen); ok {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func claudeSessionAt(path, sessionID string, since, until time.Time, needles [][]byte, seen map[string]struct{}) (DiscoveredSession, bool) {
	if _, ok := seen[sessionID]; ok {
		return DiscoveredSession{}, false
	}
	if !inWindow(path, since, until) {
		return DiscoveredSession{}, false
	}
	if needles != nil && !fileMentionsRepo(path, needles) {
		return DiscoveredSession{}, false
	}
	seen[sessionID] = struct{}{}
	return DiscoveredSession{
		Tool:      capture.ToolClaude,
		Path:      path,
		Kind:      SessionKindJSONL,
		SessionID: sessionID,
	}, true
}

// repoPathNeedles returns byte patterns whose presence means a transcript
// referenced a path inside repoRoot.
//
// The trailing separator matters: it keeps a sibling checkout such as
// /Users/joe/git/gx-cloud from matching /Users/joe/git/gx, and every session
// that actually edited the repo records at least one absolute path beneath the
// root. Windows transcripts store backslashes JSON-escaped, so both the raw and
// escaped spellings are included.
func repoPathNeedles(repoRoot string) [][]byte {
	if strings.TrimSpace(repoRoot) == "" {
		return nil
	}
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		abs = repoRoot
	}
	seen := map[string]struct{}{}
	var needles [][]byte
	add := func(s string) {
		if s == "" {
			return
		}
		if _, ok := seen[s]; ok {
			return
		}
		seen[s] = struct{}{}
		needles = append(needles, []byte(s))
	}
	addVariants := func(root string) {
		root = strings.TrimRight(root, `/\`)
		if root == "" {
			return
		}
		add(filepath.ToSlash(root) + "/")
		if native := filepath.FromSlash(root); native != filepath.ToSlash(root) {
			add(native + string(filepath.Separator))
			add(strings.ReplaceAll(native, `\`, `\\`) + `\\`)
		}
	}
	addVariants(abs)
	// A repo reached through a symlinked path (macOS /tmp, /var) is recorded
	// under its resolved name in transcripts.
	if resolved, err := filepath.EvalSymlinks(abs); err == nil {
		addVariants(resolved)
	}
	return needles
}

// repoScanChunkSize bounds the memory used while scanning one transcript.
const repoScanChunkSize = 256 * 1024

// fileMentionsRepo streams path looking for any needle. It never holds more
// than one chunk in memory, which keeps multi-megabyte transcripts cheap, and
// carries a needle-sized tail between chunks so a match spanning a chunk
// boundary is still found.
func fileMentionsRepo(path string, needles [][]byte) bool {
	if len(needles) == 0 {
		return false
	}
	overlap := 0
	for _, needle := range needles {
		if len(needle)-1 > overlap {
			overlap = len(needle) - 1
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, overlap+repoScanChunkSize)
	carry := 0
	for {
		n, err := io.ReadFull(f, buf[carry:])
		if n > 0 {
			window := buf[:carry+n]
			for _, needle := range needles {
				if bytes.Contains(window, needle) {
					return true
				}
			}
			if len(window) > overlap {
				carry = copy(buf, window[len(window)-overlap:])
			} else {
				carry = len(window)
			}
		}
		if err != nil {
			return false
		}
	}
}

func discoverCursorTranscriptSessions(homeDir, repoRoot string, since, until time.Time) ([]string, error) {
	if homeDir == "" || repoRoot == "" {
		return nil, nil
	}
	root := filepath.Join(homeDir, ".cursor", "projects", cursorRepoSlug(repoRoot), "agent-transcripts")
	if _, err := os.Stat(root); err != nil {
		return nil, unreadable(root, err)
	}
	start := since.AddDate(0, 0, -7)
	end := until.AddDate(0, 0, 7)
	var paths []string
	var problems []error
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if problem := unreadable(path, walkErr); problem != nil {
				problems = append(problems, problem)
			}
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".jsonl") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		mod := info.ModTime()
		if mod.Before(start) || mod.After(end) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if problem := unreadable(root, err); problem != nil {
		problems = append(problems, problem)
	}
	sort.SliceStable(paths, func(i, j int) bool {
		return modTime(paths[i]).After(modTime(paths[j]))
	})
	return paths, errors.Join(problems...)
}

func cursorRepoSlug(repoRoot string) string {
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		abs = repoRoot
	}
	vol := filepath.VolumeName(abs)
	abs = strings.TrimPrefix(abs, vol)
	abs = strings.Trim(abs, string(filepath.Separator))
	if abs == "" {
		return ""
	}
	return strings.ReplaceAll(abs, string(filepath.Separator), "-")
}

func cursorTranscriptSessionID(path string) string {
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

func normalizeTools(tools []string) []string {
	if len(tools) == 0 {
		return []string{capture.ToolClaude, capture.ToolCodex, capture.ToolCursor}
	}
	out := make([]string, 0, len(tools))
	seen := map[string]struct{}{}
	for _, tool := range tools {
		tool = strings.TrimSpace(strings.ToLower(tool))
		if tool == "" {
			continue
		}
		if _, ok := seen[tool]; ok {
			continue
		}
		seen[tool] = struct{}{}
		out = append(out, tool)
	}
	return out
}

func containsTool(tools []string, tool string) bool {
	for _, t := range tools {
		if t == tool {
			return true
		}
	}
	return false
}

// discoverCodexSessions walks ~/.codex/sessions. A subtree it cannot read is
// recorded and the walk continues, so one unreadable directory costs only its
// own transcripts — but it is never dropped, because "found nothing" and
// "could not look" are different answers and only one of them is fine.
func discoverCodexSessions(root string, since, until time.Time, repoRoot string) ([]string, error) {
	start := since.AddDate(0, 0, -7)
	end := until.AddDate(0, 0, 7)
	var paths []string
	var problems []error
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if problem := unreadable(path, walkErr); problem != nil {
				problems = append(problems, problem)
			}
			return nil
		}
		if d.IsDir() || !strings.HasPrefix(d.Name(), "rollout-") || !strings.HasSuffix(d.Name(), ".jsonl") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		mod := info.ModTime()
		if mod.Before(start) || mod.After(end) {
			return nil
		}
		if repoRoot != "" && !codexSessionMatchesRepo(path, repoRoot) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if problem := unreadable(root, err); problem != nil {
		problems = append(problems, problem)
	}
	return paths, errors.Join(problems...)
}

func codexSessionMatchesRepo(path, repoRoot string) bool {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return true
	}
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 16*1024)
	n, _ := f.Read(buf)
	content := string(buf[:n])
	return strings.Contains(content, repoRoot)
}

func inWindow(path string, since, until time.Time) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	mod := info.ModTime()
	return !mod.Before(since) && !mod.After(until)
}

func modTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func repoSlug(repoRoot string) string {
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		abs = repoRoot
	}
	slug := strings.ReplaceAll(abs, string(filepath.Separator), "-")
	if strings.HasPrefix(slug, "-") {
		return slug
	}
	return "-" + slug
}

// SourceEvents pairs one discovered source with the events parsed from it.
type SourceEvents struct {
	Source DiscoveredSession
	Events []capture.SessionEvent
}

// ParseAllPerSource parses each discovered source separately so callers can
// attribute every event to the file that produced it. Each event's SourceIndex
// is set to its source's position in the returned slice (which preserves the
// input order). For JSONL sources the events' SessionID is overridden with the
// source's path-derived identity: subagent transcripts carry the parent's
// in-file sessionId, so the source is the only trustworthy identity. A Cursor
// vscdb source holds many sessions, so its event SessionIDs are kept.
func ParseAllPerSource(sessions []DiscoveredSession, repoRoot string, parsers []Parser) ([]SourceEvents, error) {
	byTool := map[string]Parser{}
	for _, p := range parsers {
		byTool[p.Tool()] = p
	}
	var out []SourceEvents
	for _, session := range sessions {
		parser, ok := byTool[session.Tool]
		if !ok {
			continue
		}
		parsed, err := parser.ParseFile(session.Path, repoRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skip parse %s (%s): %v\n", session.Path, session.Tool, err)
			continue
		}
		index := len(out)
		for i := range parsed {
			parsed[i].SourceIndex = index
			if session.Kind == SessionKindJSONL && session.SessionID != "" {
				parsed[i].SessionID = session.SessionID
			}
		}
		out = append(out, SourceEvents{Source: session, Events: parsed})
	}
	return out, nil
}

