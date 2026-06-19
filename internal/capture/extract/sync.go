package extract

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
	"github.com/satoricorp/gx/internal/version"
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

	extracts, err := stager.PendingExtracts(ctx)
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

	sessions, err := stager.PendingSessions(ctx)
	if err != nil {
		return result, fmt.Errorf("list pending sessions: %w", err)
	}
	for _, row := range sessions {
		if err := uploadStagedSession(ctx, client, row, telemetryClient); err != nil {
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
		GxVersion:  version.Current(),
		HunkLinks:  staged.HunkLinks,
		FileStats:  staged.FileStats,
	}
	return client.postJSON(ctx, "/v1/extracts", req)
}

func uploadStagedSession(
	ctx context.Context,
	client *Client,
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
	content := FormatSessionContent(payload)
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
