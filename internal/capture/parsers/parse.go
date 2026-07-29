package parsers

import (
	"github.com/satoricorp/totality/internal/capture"
	"github.com/satoricorp/totality/internal/capture/parsers/claude"
	"github.com/satoricorp/totality/internal/capture/parsers/codex"
	cursorparser "github.com/satoricorp/totality/internal/capture/parsers/cursor"
)

// ParseBytes normalizes one tool's transcript bytes into session events and
// reports how many lines it could not read.
//
// It lives here rather than in the orchestrator so the upload path can
// re-derive events for a row whose staged payload has none without dragging in
// the orchestrator's git and VCS dependencies.
func ParseBytes(tool, sourcePath string, raw []byte, repoRoot string, inventory *capture.InventoryCollector) ([]capture.SessionEvent, int, error) {
	switch tool {
	case capture.ToolClaude:
		parser := &claude.Parser{Inventory: inventory}
		events, err := parser.ParseBytes(raw, sourcePath, repoRoot)
		return events, parser.LastBadLines, err
	case capture.ToolCodex:
		parser := &codex.Parser{Inventory: inventory}
		return parser.ParseBytes(raw, sourcePath, repoRoot)
	case capture.ToolCursor:
		parser := &cursorparser.Parser{Inventory: inventory}
		return parser.ParseBytes(raw, sourcePath, repoRoot)
	default:
		return nil, 0, nil
	}
}
