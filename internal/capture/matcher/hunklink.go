package matcher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/satoricorp/totality/internal/capture"
)

const (
	AuthorshipAgent   = "agent"
	AuthorshipHuman   = "human"
	AuthorshipUnknown = "unknown"
)

// HunkLink records how one commit hunk links to an agent session.
type HunkLink struct {
	HunkID           string  `json:"hunkID"`
	SessionID        string  `json:"sessionID,omitempty"`
	EventFingerprint string  `json:"eventFingerprint,omitempty"`
	Tier             int     `json:"tier"`
	Confidence       float64 `json:"confidence"`
	Authorship       string  `json:"authorship"`
	Tool             string  `json:"tool,omitempty"`
	Model            string  `json:"model,omitempty"`
}

// HunkID builds a stable identifier for a commit hunk.
func HunkID(hunk capture.CommitHunk, index int) string {
	file := filepath.ToSlash(hunk.FilePath)
	end := len(hunk.AddedLines)
	if end == 0 {
		end = 1
	}
	return fmt.Sprintf("%s:%s:1-%d", hunk.CommitSHA, file, end)
}

// BuildHunkLinks maps matcher outcomes to per-hunk link records.
func BuildHunkLinks(
	eligible []capture.CommitHunk,
	events []capture.SessionEvent,
	result Result,
) []HunkLink {
	best := map[int]MatchOutcome{}
	for _, outcome := range result.Outcomes {
		prev, ok := best[outcome.HunkIndex]
		if !ok || outcome.Tier < prev.Tier {
			best[outcome.HunkIndex] = outcome
		}
	}

	links := make([]HunkLink, len(eligible))
	for i, hunk := range eligible {
		link := HunkLink{
			HunkID:     HunkID(hunk, i),
			Authorship: AuthorshipHuman,
			Tier:       0,
		}
		outcome, ok := best[i]
		if !ok {
			links[i] = link
			continue
		}
		link.Tier = outcome.Tier
		link.Confidence = outcome.Score
		if outcome.Tier == TierTemporal {
			link.Authorship = AuthorshipUnknown
		} else if outcome.EventIndex >= 0 && outcome.EventIndex < len(events) {
			ev := events[outcome.EventIndex]
			link.SessionID = ev.SessionID
			link.EventFingerprint = eventFingerprint(ev)
			link.Tool = ev.Tool
			link.Model = ev.Model
			link.Authorship = AuthorshipAgent
			if link.Confidence == 0 {
				link.Confidence = 1.0
			}
		}
		links[i] = link
	}
	return links
}

func eventFingerprint(ev capture.SessionEvent) string {
	for _, key := range []string{"uuid", "id", "event_id", "eventId", "message_id", "messageId"} {
		if raw, ok := ev.Raw[key]; ok && len(raw) > 0 {
			value := strings.Trim(strings.TrimSpace(string(raw)), `"`)
			if value != "" && value != "null" {
				return value
			}
		}
	}
	sum := sha256.Sum256([]byte(strings.Join([]string{
		ev.Tool,
		ev.SessionID,
		strconv.FormatInt(ev.TS, 10),
		filepath.ToSlash(ev.FilePath),
		hashText(ev.OldText),
		hashText(ev.NewText),
		// Without this every command event in one millisecond fingerprints
		// alike: they share a session and carry no file path or edit text.
		hashText(ev.Command),
	}, "\x00")))
	return hex.EncodeToString(sum[:])
}

func hashText(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
