package parsers

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture"
	cursorparser "github.com/satoricorp/gx/internal/capture/parsers/cursor"
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

// DiscoverSessions returns session sources overlapping the time window.
func DiscoverSessions(opts DiscoverOptions) ([]DiscoveredSession, error) {
	if opts.HomeDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		opts.HomeDir = home
	}
	tools := normalizeTools(opts.Tools)

	var sessions []DiscoveredSession
	if containsTool(tools, capture.ToolClaude) {
		claudeDir := opts.ClaudeDir
		if claudeDir == "" {
			claudeDir = filepath.Join(opts.HomeDir, ".claude", "projects")
		}
		claudeSlug := repoSlug(opts.RepoRoot)
		claudeProject := filepath.Join(claudeDir, claudeSlug)
		if entries, err := os.ReadDir(claudeProject); err == nil {
			for _, entry := range entries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
					continue
				}
				path := filepath.Join(claudeProject, entry.Name())
				if inWindow(path, opts.Since, opts.Until) {
					sessions = append(sessions, DiscoveredSession{
						Tool:      capture.ToolClaude,
						Path:      path,
						Kind:      SessionKindJSONL,
						SessionID: strings.TrimSuffix(entry.Name(), ".jsonl"),
					})
				}
			}
		}
	}

	if containsTool(tools, capture.ToolCodex) {
		codexDir := opts.CodexDir
		if codexDir == "" {
			codexDir = filepath.Join(opts.HomeDir, ".codex", "sessions")
		}
		codexPaths, err := discoverCodexSessions(codexDir, opts.Since, opts.Until, opts.RepoRoot)
		if err != nil {
			return nil, err
		}
		for _, path := range codexPaths {
			sessions = append(sessions, DiscoveredSession{
				Tool:      capture.ToolCodex,
				Path:      path,
				Kind:      SessionKindJSONL,
				SessionID: strings.TrimSuffix(filepath.Base(path), ".jsonl"),
			})
		}
	}

	if containsTool(tools, capture.ToolCursor) {
		cursorPaths, err := discoverCursorTranscriptSessions(opts.HomeDir, opts.RepoRoot, opts.Since, opts.Until)
		if err != nil {
			return nil, err
		}
		for _, path := range cursorPaths {
			sessions = append(sessions, DiscoveredSession{
				Tool:      capture.ToolCursor,
				Path:      path,
				Kind:      SessionKindJSONL,
				SessionID: cursorTranscriptSessionID(path),
			})
		}

		vscdb := opts.CursorVSCDB
		if vscdb == "" {
			path, err := cursorparser.DiscoverVSCDBPath(opts.HomeDir)
			if err != nil {
				return nil, err
			}
			vscdb = path
		}
		if _, err := os.Stat(vscdb); err == nil {
			sessions = append(sessions, DiscoveredSession{
				Tool:      capture.ToolCursor,
				Path:      vscdb,
				Kind:      SessionKindCursorVSCDB,
				SessionID: "state.vscdb",
			})
		}
	}

	return sessions, nil
}

func discoverCursorTranscriptSessions(homeDir, repoRoot string, since, until time.Time) ([]string, error) {
	if homeDir == "" || repoRoot == "" {
		return nil, nil
	}
	root := filepath.Join(homeDir, ".cursor", "projects", cursorRepoSlug(repoRoot), "agent-transcripts")
	if _, err := os.Stat(root); err != nil {
		return nil, nil
	}
	start := since.AddDate(0, 0, -7)
	end := until.AddDate(0, 0, 7)
	var paths []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
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
	sort.SliceStable(paths, func(i, j int) bool {
		return modTime(paths[i]).After(modTime(paths[j]))
	})
	return paths, err
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

func discoverCodexSessions(root string, since, until time.Time, repoRoot string) ([]string, error) {
	start := since.AddDate(0, 0, -7)
	end := until.AddDate(0, 0, 7)
	var paths []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
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
	return paths, err
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

// ParseAll parses discovered session sources with the given parsers.
func ParseAll(sessions []DiscoveredSession, repoRoot string, parsers []Parser) ([]capture.SessionEvent, error) {
	byTool := map[string]Parser{}
	for _, p := range parsers {
		byTool[p.Tool()] = p
	}
	var events []capture.SessionEvent
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
		events = append(events, parsed...)
	}
	return events, nil
}

// CountSessionsScanned returns session units scanned per tool.
func CountSessionsScanned(sessions []DiscoveredSession, events []capture.SessionEvent) int {
	if len(sessions) == 0 {
		return 0
	}
	hasCursor := false
	jsonlCount := 0
	for _, s := range sessions {
		if s.Tool == capture.ToolCursor {
			hasCursor = true
			continue
		}
		jsonlCount++
	}
	if !hasCursor {
		return jsonlCount
	}
	seen := map[string]struct{}{}
	for _, ev := range events {
		if ev.Tool != capture.ToolCursor || ev.SessionID == "" {
			continue
		}
		seen[ev.SessionID] = struct{}{}
	}
	return jsonlCount + len(seen)
}
