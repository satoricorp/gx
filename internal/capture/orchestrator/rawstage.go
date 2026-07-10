package orchestrator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/parsers"
	"github.com/satoricorp/gx/internal/capture/parsers/claude"
	"github.com/satoricorp/gx/internal/capture/parsers/codex"
	cursorparser "github.com/satoricorp/gx/internal/capture/parsers/cursor"
	"github.com/satoricorp/gx/internal/storage"
)

// IngestRawSession stores a raw session blob before parsing and updates parsed payload.
func IngestRawSession(
	ctx context.Context,
	stager storage.CaptureStager,
	tool string,
	sourcePath string,
	raw []byte,
	repoRoot string,
	revisionIDs []string,
) (StagedSession, int, error) {
	sessionID := sessionIDFromPath(tool, sourcePath)
	if err := stageRawBlob(ctx, stager, tool, sessionID, sourcePath, raw, revisionIDs); err != nil {
		return StagedSession{}, 0, err
	}

	events, badLines, err := parseRawBytes(tool, sourcePath, raw, repoRoot, nil)
	if err != nil {
		return StagedSession{}, badLines, err
	}
	session := StagedSession{
		SessionID: sessionID,
		Tool:      tool,
		Events:    redactEvents(events),
	}
	if err := stageParsedSession(ctx, stager, session, sourcePath, raw, badLines, revisionIDs); err != nil {
		return session, badLines, err
	}
	return session, badLines, nil
}

// StageDiscoveredRaw stores raw blobs for discovered JSONL sessions before parsing.
func StageDiscoveredRaw(ctx context.Context, stager storage.CaptureStager, sessions []parsers.DiscoveredSession, revisionIDs []string) error {
	for _, session := range sessions {
		if session.Kind != parsers.SessionKindJSONL {
			continue
		}
		raw, err := os.ReadFile(session.Path)
		if err != nil {
			continue
		}
		if err := stageRawBlob(ctx, stager, session.Tool, session.SessionID, session.Path, raw, revisionIDs); err != nil {
			return err
		}
	}
	return nil
}

func stageRawBlob(
	ctx context.Context,
	stager storage.CaptureStager,
	tool, sessionID, sourcePath string,
	raw []byte,
	revisionIDs []string,
) error {
	ids := revisionIDs
	if len(ids) == 0 {
		ids = []string{""}
	}
	contentHash := storage.PayloadContentHash(raw)
	for _, revisionID := range ids {
		row := storage.StagedSession{
			SessionID:   sessionID,
			Tool:        tool,
			RawBlob:     raw,
			SourcePath:  sourcePath,
			RevisionID:  revisionID,
			ContentHash: contentHash,
		}
		row.ID = storage.CaptureRowID(revisionID, contentHash)
		if row.ID == "" {
			row.ID = uuid.NewString()
		}
		if err := stager.StageSession(ctx, row); err != nil {
			return err
		}
	}
	return nil
}

func stageParsedSession(
	ctx context.Context,
	stager storage.CaptureStager,
	session StagedSession,
	sourcePath string,
	raw []byte,
	badLines int,
	revisionIDs []string,
) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ids := revisionIDs
	if len(ids) == 0 {
		ids = []string{""}
	}
	contentHash := storage.PayloadContentHash(raw)
	if len(raw) == 0 {
		contentHash = storage.PayloadContentHash(payload)
	}
	for _, revisionID := range ids {
		row := storage.StagedSession{
			SessionID:   session.SessionID,
			Tool:        session.Tool,
			PayloadJSON: payload,
			RawBlob:     raw,
			SourcePath:  sourcePath,
			BadLines:    badLines,
			RevisionID:  revisionID,
			ContentHash: contentHash,
		}
		row.ID = storage.CaptureRowID(revisionID, contentHash)
		if row.ID == "" {
			row.ID = uuid.NewString()
		}
		if err := stager.StageSession(ctx, row); err != nil {
			return err
		}
	}
	return nil
}

func ParseRawBytes(tool, sourcePath string, raw []byte, repoRoot string, inventory *capture.InventoryCollector) ([]capture.SessionEvent, int, error) {
	return parseRawBytes(tool, sourcePath, raw, repoRoot, inventory)
}

func parseRawBytes(tool, sourcePath string, raw []byte, repoRoot string, inventory *capture.InventoryCollector) ([]capture.SessionEvent, int, error) {
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

func sessionIDFromPath(tool, sourcePath string) string {
	switch tool {
	case capture.ToolClaude, capture.ToolCodex:
		return strings.TrimSuffix(filepath.Base(sourcePath), ".jsonl")
	case capture.ToolCursor:
		stem := strings.TrimSuffix(filepath.Base(sourcePath), ".jsonl")
		if stem != "" {
			return "cursor:" + stem
		}
	}
	return strings.TrimSuffix(filepath.Base(sourcePath), ".jsonl")
}
