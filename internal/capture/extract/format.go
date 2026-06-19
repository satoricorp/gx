package extract

import (
	"fmt"
	"strings"

	"github.com/satoricorp/gx/internal/capture"
)

// SessionPayload is one session blob ready for upload.
type SessionPayload struct {
	SessionID string
	Tool      string
	Events    []capture.SessionEvent
}

// FormatSessionContent builds a line-oriented redacted transcript for sessions_raw.
func FormatSessionContent(session SessionPayload) string {
	var lines []string
	for _, ev := range session.Events {
		line := formatEventLine(ev)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func formatEventLine(ev capture.SessionEvent) string {
	switch ev.Kind {
	case capture.KindEdit:
		return fmt.Sprintf("[%s] edit %s\n%s", ev.Tool, ev.FilePath, ev.NewText)
	case capture.KindMessage:
		if ev.PromptContext != "" {
			return fmt.Sprintf("[%s] message %s", ev.Tool, ev.PromptContext)
		}
		return fmt.Sprintf("[%s] message", ev.Tool)
	default:
		if ev.PromptContext != "" {
			return fmt.Sprintf("[%s] %s %s", ev.Tool, ev.Kind, ev.PromptContext)
		}
		return fmt.Sprintf("[%s] %s", ev.Tool, ev.Kind)
	}
}

// SessionModel picks the first non-empty model from session events.
func SessionModel(session SessionPayload) string {
	for _, ev := range session.Events {
		if strings.TrimSpace(ev.Model) != "" {
			return ev.Model
		}
	}
	return ""
}
