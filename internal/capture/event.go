package capture

import "encoding/json"

const (
	ToolClaude = "claude"
	ToolCodex  = "codex"
	ToolCursor = "cursor"

	KindEdit       = "edit"
	KindRead       = "read"
	KindToolCall   = "tool_call"
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
	Raw           map[string]json.RawMessage

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
