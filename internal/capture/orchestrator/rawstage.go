package orchestrator

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/parsers"
	"github.com/satoricorp/gx/internal/storage"
)

// sourceFingerprint describes one transcript file's bytes at one instant, so
// a later reader can tell whether the file grew or was rewritten.
type sourceFingerprint struct {
	hash    string
	bytes   int64
	mtimeMS int64
}

func fingerprintBytes(raw []byte, path string) sourceFingerprint {
	fp := sourceFingerprint{
		hash:  storage.PayloadContentHash(raw),
		bytes: int64(len(raw)),
	}
	if info, err := os.Stat(path); err == nil {
		fp.mtimeMS = info.ModTime().UnixMilli()
	}
	return fp
}

// IngestRawSession parses one transcript and stages its pointer row: parsed
// events plus source path and content fingerprint, never the raw bytes. The
// upload path re-reads the transcript from disk at upload time.
func IngestRawSession(
	ctx context.Context,
	stager storage.CaptureStager,
	tool string,
	sourcePath string,
	raw []byte,
	repoRoot string,
	revisionIDs []string,
) (StagedSession, int, error) {
	sessionID := parsers.SourceSessionID(tool, sourcePath)
	events, badLines, err := parseRawBytes(tool, sourcePath, raw, repoRoot, nil)
	if err != nil {
		return StagedSession{}, badLines, err
	}
	// The per-source identity wins over any in-file sessionId: Claude stamps
	// subagent transcripts with the parent conversation's id.
	for i := range events {
		events[i].SessionID = sessionID
	}
	session := StagedSession{
		SessionID: sessionID,
		Tool:      tool,
		Events:    redactEvents(events),
	}
	if _, err := stageSourceSession(ctx, stager, session, sourcePath, fingerprintBytes(raw, sourcePath), badLines, revisionIDs); err != nil {
		return session, badLines, err
	}
	return session, badLines, nil
}

// stageMatchedSource stages the capture_sessions rows for one source whose
// events matched this push's hunks, and returns the ids of the rows it wrote.
// JSONL sources get exactly one pointer row. A Cursor vscdb holds many sessions
// and its bytes are a SQLite database, not a transcript, so it stages one
// parsed-events row per contained session with no source pointer for the upload
// path to re-read.
func stageMatchedSource(
	ctx context.Context,
	stager storage.CaptureStager,
	source parsers.DiscoveredSession,
	events []capture.SessionEvent,
	revisionIDs []string,
) ([]string, error) {
	if source.Kind == parsers.SessionKindCursorVSCDB {
		var staged []string
		for _, session := range groupSessions(events) {
			id, err := stageSourceSession(ctx, stager, session, "", sourceFingerprint{}, 0, revisionIDs)
			if err != nil {
				return staged, err
			}
			staged = append(staged, id)
		}
		return staged, nil
	}
	session := StagedSession{
		SessionID: source.SessionID,
		Tool:      source.Tool,
		Events:    events,
	}
	// The fingerprint read is the only time staging touches the transcript
	// bytes; they are hashed and dropped, not stored. A read failure still
	// stages the pointer — the upload path reports the missing source.
	var fp sourceFingerprint
	if raw, err := os.ReadFile(source.Path); err == nil {
		fp = fingerprintBytes(raw, source.Path)
	}
	id, err := stageSourceSession(ctx, stager, session, source.Path, fp, 0, revisionIDs)
	if err != nil {
		return nil, err
	}
	return []string{id}, nil
}

// stageSourceSession writes the single capture_sessions row for one staged
// session. The row is keyed by source identity, so re-staging on a later push
// updates in place instead of accumulating copies.
//
// A push usually spans several revisions and one session usually backs many
// of them. That many-to-many already lives in the extract's hunk links (each
// link names its commit), so the session row carries only one revision from
// the pushed range: RunPush marks rows shareable with the full pushed
// revision set, and membership of any one revision is enough to be swept in.
// The old row-per-revision model duplicated multi-megabyte blobs and said
// nothing the hunk links do not.
// It returns the row id it wrote so callers can address the row directly,
// without re-deriving it or matching on revision.
func stageSourceSession(
	ctx context.Context,
	stager storage.CaptureStager,
	session StagedSession,
	sourcePath string,
	fp sourceFingerprint,
	badLines int,
	revisionIDs []string,
) (string, error) {
	payload, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	row := storage.StagedSession{
		SessionID:   session.SessionID,
		Tool:        session.Tool,
		PayloadJSON: payload,
		SourcePath:  sourcePath,
		SourceBytes: fp.bytes,
		SourceMTime: fp.mtimeMS,
		BadLines:    badLines,
		RevisionID:  firstNonEmpty(revisionIDs),
		ContentHash: fp.hash,
	}
	row.ID = storage.SessionSourceRowID(session.Tool, sourcePath, session.SessionID)
	if row.ID == "" {
		row.ID = uuid.NewString()
	}
	if err := stager.StageSession(ctx, row); err != nil {
		return "", err
	}
	return row.ID, nil
}

func firstNonEmpty(values []string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func ParseRawBytes(tool, sourcePath string, raw []byte, repoRoot string, inventory *capture.InventoryCollector) ([]capture.SessionEvent, int, error) {
	return parseRawBytes(tool, sourcePath, raw, repoRoot, inventory)
}

func parseRawBytes(tool, sourcePath string, raw []byte, repoRoot string, inventory *capture.InventoryCollector) ([]capture.SessionEvent, int, error) {
	return parsers.ParseBytes(tool, sourcePath, raw, repoRoot, inventory)
}
