package capture

import (
	"encoding/json"
	"strings"
)

const (
	ToolClaude = "claude"
	ToolCodex  = "codex"
	ToolCursor = "cursor"

	KindEdit       = "edit"
	KindRead       = "read"
	KindToolCall   = "tool_call"
	KindCommand    = "command"
	KindMessage    = "message"
	KindToolResult = "tool_result"
)

// SessionEvent is a normalized agent session record derived from JSONL logs.
type SessionEvent struct {
	SessionID     string
	Tool          string // "claude" | "codex" | "cursor"
	Model         string
	TS            int64 // unix milliseconds
	Kind          string
	FilePath      string
	OldText       string
	NewText       string
	PromptContext string
	// Command is the shell text of a command the agent ran, for tool calls that
	// run one. An agent that edits with `sed -i`, a heredoc or an inline script
	// changes files without ever emitting a file-edit event, so without this the
	// work is invisible to attribution even though the transcript recorded it in
	// full. It is deliberately not truncated: the match is against the text the
	// command wrote, so a clipped command matches nothing.
	Command       string
	Raw           map[string]json.RawMessage

	// Cwd is the working directory the agent was in when this event happened,
	// when the tool records one. It is what binds a session to a repository —
	// and because Claude writes it per entry rather than once per session, a
	// session that moves between repositories can be attributed per event
	// rather than wholesale.
	Cwd string
	// OriginURL is a git remote the tool recorded itself. Codex writes one in
	// its session_meta; most tools do not. It is a fallback for identifying a
	// repository whose checkout no longer exists on this machine.
	OriginURL string

	// SourceIndex attributes the event to the source file that produced it
	// (index into one parse run's source list). SessionID cannot serve here:
	// Claude subagent transcripts share the parent's sessionId. Only
	// meaningful within a single run, so never serialized.
	SourceIndex int `json:"-"`
}

// CommitHunk is one file's added lines from a commit in a ref range.
type CommitHunk struct {
	CommitSHA  string
	FilePath   string
	AddedLines []string
	CommitTime int64 // unix milliseconds
}

// IsEditEvent reports whether the event represents an agent file edit.
func (e SessionEvent) IsEditEvent() bool {
	switch e.Kind {
	case KindEdit:
		return true
	default:
		return false
	}
}

// IsCommandEvent reports whether the event ran a shell command whose text was
// recorded. It is kept apart from IsEditEvent because the two carry different
// evidence: an edit event names the file it wrote, while a command only shows
// the text that went through it, so a command can support authorship but never
// declare it.
func (e SessionEvent) IsCommandEvent() bool {
	return e.Kind == KindCommand && strings.TrimSpace(e.Command) != ""
}

// Directory reports the working directory the agent was in for this event, and
// RecordedOrigin the git remote the tool recorded itself. Together they satisfy
// the interface repobind uses to bind a session to a repository without that
// package depending on this one.
func (e SessionEvent) Directory() string { return e.Cwd }

// RecordedOrigin is the remote URL the tool wrote into its own session record,
// where it writes one.
func (e SessionEvent) RecordedOrigin() string { return e.OriginURL }
