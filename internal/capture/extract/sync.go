package extract

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/satoricorp/totality/internal/auth"
	"github.com/satoricorp/totality/internal/capture"
	"github.com/satoricorp/totality/internal/capture/matcher"
	"github.com/satoricorp/totality/internal/capture/parsers"
	"github.com/satoricorp/totality/internal/capture/redact"
	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/telemetry"
	"github.com/satoricorp/totality/internal/version"
)

// SyncResult summarizes a capture staging drain.
type SyncResult struct {
	ExtractsUploaded int
	SessionsUploaded int
	ExtractErrors    int
	SessionErrors    int
}

type stagedExtractPayload struct {
	RefRange  string             `json:"refRange"`
	RepoRoot  string             `json:"repoRoot"`
	Head      string             `json:"head"`
	HunkLinks []matcher.HunkLink `json:"hunkLinks"`
	FileStats interface{}        `json:"fileStats"`
}

type stagedSessionPayload struct {
	SessionID string                 `json:"sessionID"`
	Tool      string                 `json:"tool"`
	Events    []capture.SessionEvent `json:"events"`
}

// SyncPending uploads all pending capture_extracts and capture_sessions rows.
func SyncPending(
	ctx context.Context,
	stager storage.CaptureStager,
	creds auth.Credentials,
	telemetryClient telemetry.Client,
) (SyncResult, error) {
	client := NewClient(creds.APIURL, creds.Token)
	var result SyncResult

	extracts, err := stager.ShareableExtracts(ctx)
	if err != nil {
		return result, fmt.Errorf("list pending extracts: %w", err)
	}
	for _, row := range extracts {
		if err := uploadStagedExtract(ctx, client, row); err != nil {
			_ = stager.SetExtractUploadError(ctx, row.ID, err.Error())
			result.ExtractErrors++
			continue
		}
		if err := stager.MarkExtractUploaded(ctx, row.ID); err != nil {
			return result, fmt.Errorf("mark extract uploaded: %w", err)
		}
		result.ExtractsUploaded++
	}

	sessions, err := stager.ShareableSessions(ctx)
	if err != nil {
		return result, fmt.Errorf("list pending sessions: %w", err)
	}
	for _, row := range sessions {
		// A single vanished source must never sink the whole batch: the row
		// is error-marked and the remaining sessions still upload. Downstream
		// already degrades a session-less push to "session context
		// unavailable".
		if err := uploadStagedSession(ctx, client, stager, row, telemetryClient); err != nil {
			_ = stager.SetSessionUploadError(ctx, row.ID, err.Error())
			result.SessionErrors++
			continue
		}
		if err := stager.MarkSessionUploaded(ctx, row.ID); err != nil {
			return result, fmt.Errorf("mark session uploaded: %w", err)
		}
		result.SessionsUploaded++
	}

	return result, nil
}

func uploadStagedExtract(ctx context.Context, client *Client, row storage.StagedExtract) error {
	var staged stagedExtractPayload
	if err := json.Unmarshal(row.PayloadJSON, &staged); err != nil {
		return fmt.Errorf("decode extract payload: %w", err)
	}
	head := staged.Head
	if head == "" {
		head = "HEAD"
	}
	req := ExtractRequest{
		RepoRoot:   row.RepoRoot,
		RefRange:   row.RefRange,
		HeadCommit: head,
		TlVersion:  version.Current(),
		HunkLinks:  staged.HunkLinks,
		FileStats:  staged.FileStats,
	}
	return client.postJSON(ctx, "/v1/extracts", req)
}

func uploadStagedSession(
	ctx context.Context,
	client *Client,
	stager storage.CaptureStager,
	row storage.StagedSession,
	telemetryClient telemetry.Client,
) error {
	var staged stagedSessionPayload
	if err := json.Unmarshal(row.PayloadJSON, &staged); err != nil {
		return fmt.Errorf("decode session payload: %w", err)
	}
	payload := SessionPayload{
		SessionID: staged.SessionID,
		Tool:      staged.Tool,
		Events:    staged.Events,
	}
	if payload.SessionID == "" {
		payload.SessionID = row.SessionID
	}
	if payload.Tool == "" {
		payload.Tool = row.Tool
	}
	content, err := sessionContent(ctx, stager, row, payload)
	if err != nil {
		return err
	}
	body := SessionRequest{
		SessionId:    payload.SessionID,
		Tool:         payload.Tool,
		Model:        SessionModel(payload),
		Content:      content,
		CapturedAtMs: row.CreatedAt,
	}
	if err := client.postJSON(ctx, "/v1/sessions", body); err != nil {
		return err
	}
	if telemetryClient != nil {
		telemetryClient.EmitSessionUploaded(ctx, telemetry.SessionUploadedProps{
			SessionID: payload.SessionID,
			Tool:      payload.Tool,
			Bytes:     len(content),
		})
	}
	return nil
}

// maxSessionContentBytes bounds one session upload.
//
// Transcripts grow while their session runs, and re-staging a grown transcript
// clears uploaded_at, so the row is re-sent on the next push. Without a ceiling
// that meant re-uploading an entire multi-megabyte transcript — one staged
// pointer in the field targeted a 30 MB file — on every push, forever, with a
// fresh full copy inserted server-side each time. Formatted event content is a
// small fraction of a raw transcript, and this caps what is left.
const maxSessionContentBytes = 1 << 20

// sessionContent resolves the wire content for one staged session row.
//
// Two invariants, both learned the hard way:
//
// Format. The server promoter parses sessions_raw.content as Totality's own
// line-oriented format (server/src/ingest/promote.ts, matching format.go) and
// discards any line that is not "[tool] kind …". Raw transcript JSONL starts
// with "{", so uploading source bytes promotes zero events — and promotion
// DELETEs the session's existing events first, so it actively destroys context
// rather than merely failing to add it. Content is therefore always formatted
// events, never file bytes.
//
// Redaction. redact.Redact is applied to events, so anything derived from event
// text has been through it. Staging redacts before writing payload_json, which
// makes the staged payload the preferred source. When a row has no events the
// bytes are re-parsed here and redacted here, so no path ships unredacted
// transcript text off the machine.
func sessionContent(ctx context.Context, stager storage.CaptureStager, row storage.StagedSession, payload SessionPayload) (string, error) {
	if content := FormatSessionContent(payload); content != "" {
		return capSessionContent(content), nil
	}
	// Rows staged with no events at all: every legacy row written before
	// staging populated them, which is exactly the shape behind the observed
	// `400 sessionId, tool, and content are required` responses.
	events, err := eventsFromSource(ctx, stager, row)
	if err != nil {
		return "", err
	}
	content := FormatSessionContent(SessionPayload{
		SessionID: payload.SessionID,
		Tool:      payload.Tool,
		Events:    redactEvents(events),
	})
	if content == "" {
		// Uploading this would earn a 400 on every attempt. Failing here makes
		// it a recorded, attempt-counted error instead of an invisible loop.
		return "", fmt.Errorf("session %s has no uploadable content: no events parsed from %s",
			payload.SessionID, sourceLabel(row))
	}
	return capSessionContent(content), nil
}

// eventsFromSource re-derives events for a row whose payload carries none, from
// the transcript on disk or from a legacy inline blob.
func eventsFromSource(ctx context.Context, stager storage.CaptureStager, row storage.StagedSession) ([]capture.SessionEvent, error) {
	if row.SourcePath != "" {
		raw, err := os.ReadFile(row.SourcePath)
		if err == nil {
			// A transcript that grew since staging is expected, not a
			// mismatch: the session kept running. Record what was actually
			// read so the fingerprint describes the uploaded content.
			hash := storage.PayloadContentHash(raw)
			if hash != row.ContentHash || int64(len(raw)) != row.SourceBytes {
				mtime := row.SourceMTime
				if info, statErr := os.Stat(row.SourcePath); statErr == nil {
					mtime = info.ModTime().UnixMilli()
				}
				_ = stager.SetSessionSourceFingerprint(ctx, row.ID, hash, int64(len(raw)), mtime)
			}
			return parseSessionBytes(row.Tool, row.SourcePath, raw)
		}
		if len(row.RawBlob) > 0 {
			return parseSessionBytes(row.Tool, row.SourcePath, row.RawBlob)
		}
		return nil, fmt.Errorf("session source missing: %s: %v", row.SourcePath, err)
	}
	if len(row.RawBlob) > 0 {
		return parseSessionBytes(row.Tool, row.SessionID, row.RawBlob)
	}
	return nil, fmt.Errorf("session %s has neither events, a source path, nor an inline blob", row.SessionID)
}

// parseSessionBytes normalizes transcript bytes with the tool's own parser.
// repoRoot is empty because a staged session row does not carry one; that only
// leaves edit paths absolute, which the promoter stores verbatim either way.
func parseSessionBytes(tool, sourcePath string, raw []byte) ([]capture.SessionEvent, error) {
	events, _, err := parsers.ParseBytes(tool, sourcePath, raw, "", nil)
	if err != nil {
		return nil, fmt.Errorf("parse session source %s: %w", sourcePath, err)
	}
	return events, nil
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

// capSessionContent truncates at a line boundary so the promoter still sees
// whole events, and says so in-band rather than silently shipping less.
func capSessionContent(content string) string {
	if len(content) <= maxSessionContentBytes {
		return content
	}
	cut := strings.LastIndexByte(content[:maxSessionContentBytes], '\n')
	if cut <= 0 {
		cut = maxSessionContentBytes
	}
	omitted := len(content) - cut
	return content[:cut] + fmt.Sprintf("\n[tl] truncated %d bytes omitted", omitted)
}

func sourceLabel(row storage.StagedSession) string {
	if row.SourcePath != "" {
		return row.SourcePath
	}
	return "session " + row.SessionID
}
