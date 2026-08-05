package reparse

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/satoricorp/lgtm/internal/capture"
	"github.com/satoricorp/lgtm/internal/capture/orchestrator"
	"github.com/satoricorp/lgtm/internal/capture/redact"
	"github.com/satoricorp/lgtm/internal/storage"
	"github.com/satoricorp/lgtm/internal/telemetry"
)

// Result summarizes one reparse run.
type Result struct {
	Sessions int
	Updated  int
	Errors   int
}

// Run re-normalizes stored raw session blobs in capture_sessions.
func Run(ctx context.Context, stager storage.CaptureStager, repoRoot string, telemetryClient telemetry.Client) (Result, error) {
	if stager == nil {
		return Result{}, fmt.Errorf("capture stager required")
	}
	rows, err := stager.RawSessions(ctx)
	if err != nil {
		return Result{}, err
	}
	result := Result{Sessions: len(rows)}
	inventory := capture.NewInventoryCollector()
	for _, row := range rows {
		if len(row.RawBlob) == 0 {
			continue
		}
		sourcePath := row.SourcePath
		if sourcePath == "" {
			sourcePath = filepath.Join("raw", row.SessionID+".jsonl")
		}
		events, badLines, err := orchestrator.ParseRawBytes(row.Tool, sourcePath, row.RawBlob, repoRoot, inventory)
		if err != nil {
			result.Errors++
			continue
		}
		session := orchestrator.StagedSession{
			SessionID: row.SessionID,
			Tool:      row.Tool,
			Events:    redactEvents(events),
		}
		payload, err := json.Marshal(session)
		if err != nil {
			result.Errors++
			continue
		}
		row.PayloadJSON = payload
		row.BadLines = badLines
		// The stored blob was only needed for parsing; the upsert keeps the
		// existing raw_blob column, and blobs are never written anymore.
		row.RawBlob = nil
		if err := stager.StageSession(ctx, row); err != nil {
			result.Errors++
			continue
		}
		result.Updated++
	}
	emitSchemaDrift(ctx, telemetryClient, inventory)
	return result, nil
}

func redactEvents(events []capture.SessionEvent) []capture.SessionEvent {
	out := make([]capture.SessionEvent, len(events))
	for i, ev := range events {
		ev.OldText = redact.Redact(ev.OldText)
		ev.NewText = redact.Redact(ev.NewText)
		ev.PromptContext = redact.Redact(ev.PromptContext)
		out[i] = ev
	}
	return out
}

func emitSchemaDrift(ctx context.Context, client telemetry.Client, inventory *capture.InventoryCollector) {
	if inventory == nil {
		return
	}
	if client == nil {
		client = telemetry.NewFromEnv()
	}
	for tool, paths := range inventory.UnknownPaths() {
		if len(paths) == 0 {
			continue
		}
		client.EmitSchemaDrift(ctx, telemetry.SchemaDriftProps{
			Tool:  tool,
			Paths: paths,
			Count: len(paths),
		})
	}
}
