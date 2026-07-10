package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/satoricorp/gx/internal/providers"
)

// UsageBackfillResult reports what BackfillUsage touched.
type UsageBackfillResult struct {
	ResponsesScanned  int
	ResponsesUpdated  int
	SessionsRefreshed int
}

// BackfillUsage parses stored response bodies to populate the token columns
// on legacy responses rows that were captured before usage extraction, then
// recomputes the usage rollups (tokens and models_json) for stale sessions.
// It is idempotent and cheap once the backlog has been processed.
func (s *Store) BackfillUsage(ctx context.Context) (UsageBackfillResult, error) {
	var result UsageBackfillResult

	type usageUpdate struct {
		id           string
		sessionID    string
		usage        providers.Usage
		finishReason *string
	}
	var updates []usageUpdate

	rows, err := s.db.QueryContext(ctx, `
		SELECT resp.id, req.session_id, req.provider, req.endpoint, resp.response_body
		FROM responses resp
		JOIN requests req ON req.id = resp.request_id
		WHERE resp.input_tokens IS NULL
			AND resp.output_tokens IS NULL
			AND CAST(resp.response_body AS TEXT) LIKE '%usage%'
	`)
	if err != nil {
		return result, fmt.Errorf("list responses missing usage: %w", err)
	}
	for rows.Next() {
		var id, sessionID, provider, endpoint string
		var body []byte
		if err := rows.Scan(&id, &sessionID, &provider, &endpoint, &body); err != nil {
			rows.Close()
			return result, fmt.Errorf("scan response missing usage: %w", err)
		}
		result.ResponsesScanned++
		summary, ok := summarizeStoredResponse(provider, endpoint, body)
		if !ok {
			continue
		}
		updates = append(updates, usageUpdate{
			id:           id,
			sessionID:    sessionID,
			usage:        summary.Usage,
			finishReason: summary.FinishReason,
		})
	}
	if err := rows.Close(); err != nil {
		return result, fmt.Errorf("close responses missing usage: %w", err)
	}
	if err := rows.Err(); err != nil {
		return result, fmt.Errorf("iterate responses missing usage: %w", err)
	}

	staleSessions := map[string]struct{}{}
	if len(updates) > 0 {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return result, fmt.Errorf("begin usage backfill: %w", err)
		}
		defer tx.Rollback()
		for _, update := range updates {
			if _, err := tx.ExecContext(ctx, `
				UPDATE responses
				SET input_tokens = ?,
					output_tokens = ?,
					cache_read_tokens = ?,
					cache_write_tokens = ?,
					finish_reason = COALESCE(finish_reason, ?)
				WHERE id = ?
			`, update.usage.InputTokens, update.usage.OutputTokens, update.usage.CacheReadTokens, update.usage.CacheWriteTokens, update.finishReason, update.id); err != nil {
				return result, fmt.Errorf("backfill response usage: %w", err)
			}
			result.ResponsesUpdated++
			staleSessions[update.sessionID] = struct{}{}
		}
		if err := tx.Commit(); err != nil {
			return result, fmt.Errorf("commit usage backfill: %w", err)
		}
	}

	if err := s.collectSessionsMissingModels(ctx, staleSessions); err != nil {
		return result, err
	}
	for sessionID := range staleSessions {
		if err := s.refreshSessionUsage(ctx, sessionID); err != nil {
			return result, err
		}
		result.SessionsRefreshed++
	}
	return result, nil
}

// collectSessionsMissingModels adds sessions whose models_json rollup was
// never computed even though their requests recorded a model.
func (s *Store) collectSessionsMissingModels(ctx context.Context, sessionIDs map[string]struct{}) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT req.session_id
		FROM requests req
		JOIN sessions s ON s.id = req.session_id
		WHERE s.models_json = '[]'
			AND req.model IS NOT NULL AND TRIM(req.model) != ''
	`)
	if err != nil {
		return fmt.Errorf("list sessions missing models: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sessionID string
		if err := rows.Scan(&sessionID); err != nil {
			return fmt.Errorf("scan session missing models: %w", err)
		}
		sessionIDs[sessionID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate sessions missing models: %w", err)
	}
	return nil
}

func summarizeStoredResponse(provider, endpoint string, body []byte) (providers.Summary, bool) {
	if summary, err := providers.SummarizeJSON(provider, body); err == nil && usagePresent(summary.Usage) {
		return summary, true
	}
	if looksLikeEventStream(body) {
		if handler := providers.HandlerFor(provider, endpoint); handler != nil {
			if _, summary, err := handler.Assemble(bytes.NewReader(body)); err == nil && usagePresent(summary.Usage) {
				return summary, true
			}
		}
	}
	return providers.Summary{}, false
}

func usagePresent(usage providers.Usage) bool {
	return usage.InputTokens != nil || usage.OutputTokens != nil ||
		usage.CacheReadTokens != nil || usage.CacheWriteTokens != nil
}

func looksLikeEventStream(body []byte) bool {
	trimmed := bytes.TrimSpace(body)
	return bytes.HasPrefix(trimmed, []byte("data:")) ||
		bytes.HasPrefix(trimmed, []byte("event:")) ||
		bytes.Contains(body, []byte("\ndata:"))
}
