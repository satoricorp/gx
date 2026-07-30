package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

func (s *Store) UpsertSession(ctx context.Context, session Session) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sessions (
			id, created_at, ended_at, command, cwd, client_pid, exit_code, tx_version,
			source, process_name, parent_pid, last_seen_at, end_reason, repo_root
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			ended_at = COALESCE(excluded.ended_at, sessions.ended_at),
			command = CASE WHEN excluded.command != '' THEN excluded.command ELSE sessions.command END,
			cwd = CASE WHEN excluded.cwd != '' THEN excluded.cwd ELSE sessions.cwd END,
			client_pid = COALESCE(excluded.client_pid, sessions.client_pid),
			exit_code = COALESCE(excluded.exit_code, sessions.exit_code),
			tx_version = CASE WHEN excluded.tx_version != '' THEN excluded.tx_version ELSE sessions.tx_version END,
			source = COALESCE(excluded.source, sessions.source),
			process_name = COALESCE(excluded.process_name, sessions.process_name),
			parent_pid = COALESCE(excluded.parent_pid, sessions.parent_pid),
			last_seen_at = COALESCE(excluded.last_seen_at, sessions.last_seen_at),
			end_reason = COALESCE(excluded.end_reason, sessions.end_reason),
			repo_root = COALESCE(excluded.repo_root, sessions.repo_root)
	`,
		session.ID,
		session.CreatedAt,
		session.EndedAt,
		session.Command,
		session.Cwd,
		session.ClientPID,
		session.ExitCode,
		session.TLVersion,
		session.Source,
		session.ProcessName,
		session.ParentPID,
		session.LastSeenAt,
		session.EndReason,
		session.RepoRoot,
	)
	if err != nil {
		return fmt.Errorf("upsert session: %w", err)
	}
	return nil
}

// UpsertObservedSession records a session that was observed as a transcript on
// disk, filling in blanks on an existing row but never overwriting what is
// already there.
//
// UpsertSession is the wrong writer for this. Its ON CONFLICT clause lets a
// non-empty incoming value win, and a transcript observation only knows the
// tool name. Session IDs derived from a transcript path can collide with IDs
// an earlier, richer writer minted: Cursor's state.vscdb leg keeps the
// parser-minted "cursor-<composerID>", byte-identical to the ids the retired
// cursor ingest wrote, and real databases still hold its rows with the
// composer title as the command ("cursor: Refactor the auth middleware").
// An UpsertSession here would replace that title with a bare "cursor",
// permanently — nothing rewrites command afterwards.
func (s *Store) UpsertObservedSession(ctx context.Context, session Session) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sessions (id, created_at, command, cwd, tx_version, source, repo_root)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			command = CASE WHEN TRIM(COALESCE(sessions.command, '')) = '' THEN excluded.command ELSE sessions.command END,
			cwd = CASE WHEN TRIM(COALESCE(sessions.cwd, '')) = '' THEN excluded.cwd ELSE sessions.cwd END,
			tx_version = CASE WHEN TRIM(COALESCE(sessions.tx_version, '')) = '' THEN excluded.tx_version ELSE sessions.tx_version END,
			source = CASE WHEN TRIM(COALESCE(sessions.source, '')) = '' THEN excluded.source ELSE sessions.source END,
			repo_root = CASE WHEN TRIM(COALESCE(sessions.repo_root, '')) = '' THEN excluded.repo_root ELSE sessions.repo_root END
	`,
		session.ID,
		session.CreatedAt,
		session.Command,
		session.Cwd,
		session.TLVersion,
		session.Source,
		session.RepoRoot,
	)
	if err != nil {
		return fmt.Errorf("upsert observed session: %w", err)
	}
	return nil
}

func (s *Store) UpsertSessionContext(ctx context.Context, session Session, context SessionContext) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin session context: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO sessions (
			id, created_at, ended_at, command, cwd, client_pid, exit_code, tx_version,
			source, process_name, parent_pid, last_seen_at, end_reason, repo_root
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			ended_at = COALESCE(excluded.ended_at, sessions.ended_at),
			command = CASE WHEN excluded.command != '' THEN excluded.command ELSE sessions.command END,
			cwd = CASE WHEN excluded.cwd != '' THEN excluded.cwd ELSE sessions.cwd END,
			client_pid = COALESCE(excluded.client_pid, sessions.client_pid),
			exit_code = COALESCE(excluded.exit_code, sessions.exit_code),
			tx_version = CASE WHEN excluded.tx_version != '' THEN excluded.tx_version ELSE sessions.tx_version END,
			source = COALESCE(excluded.source, sessions.source),
			process_name = COALESCE(excluded.process_name, sessions.process_name),
			parent_pid = COALESCE(excluded.parent_pid, sessions.parent_pid),
			last_seen_at = COALESCE(excluded.last_seen_at, sessions.last_seen_at),
			end_reason = COALESCE(excluded.end_reason, sessions.end_reason),
			repo_root = COALESCE(excluded.repo_root, sessions.repo_root)
	`, session.ID, session.CreatedAt, session.EndedAt, session.Command, session.Cwd, session.ClientPID, session.ExitCode, session.TLVersion, session.Source, session.ProcessName, session.ParentPID, session.LastSeenAt, session.EndReason, session.RepoRoot); err != nil {
		return fmt.Errorf("upsert matched session: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO session_contexts (session_id, tool, model, format, content_json, captured_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(session_id) DO UPDATE SET
			tool = excluded.tool,
			model = excluded.model,
			format = excluded.format,
			content_json = excluded.content_json,
			captured_at = excluded.captured_at
	`, context.SessionID, context.Tool, context.Model, context.Format, context.ContentJSON, context.CapturedAt); err != nil {
		return fmt.Errorf("upsert session context: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit session context: %w", err)
	}
	if err := s.refreshSessionUsage(ctx, context.SessionID); err != nil {
		return err
	}
	return nil
}

func (s *Store) WriteRequest(ctx context.Context, req Request) error {
	_, err := s.requestStmt.ExecContext(
		ctx,
		req.ID,
		req.SessionID,
		req.CreatedAt,
		req.Provider,
		req.Endpoint,
		req.Method,
		req.Model,
		req.RequestBody,
		req.RequestHeaders,
	)
	if err != nil {
		return fmt.Errorf("insert request: %w", err)
	}
	if err := s.refreshSessionUsage(ctx, req.SessionID); err != nil {
		return err
	}
	return nil
}

func (s *Store) WriteResponse(ctx context.Context, resp Response) error {
	streaming := 0
	if resp.IsStreaming {
		streaming = 1
	}
	if resp.ResponseBody == nil {
		// Empty upstream bodies arrive as nil slices, which the driver binds
		// as NULL and the NOT NULL column rejects, dropping the capture row.
		resp.ResponseBody = []byte{}
	}
	_, err := s.respStmt.ExecContext(
		ctx,
		resp.ID,
		resp.RequestID,
		resp.CreatedAt,
		resp.CompletedAt,
		resp.StatusCode,
		resp.ResponseBody,
		resp.ResponseHeaders,
		streaming,
		resp.DurationMS,
		resp.ProviderRequestID,
		resp.FinishReason,
		resp.InputTokens,
		resp.OutputTokens,
		resp.CacheReadTokens,
		resp.CacheWriteTokens,
		resp.Error,
	)
	if err != nil {
		return fmt.Errorf("insert response: %w", err)
	}
	sessionID, err := s.sessionIDForRequest(ctx, resp.RequestID)
	if err != nil {
		return err
	}
	if sessionID != "" {
		if err := s.refreshSessionUsage(ctx, sessionID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SessionContext(ctx context.Context, sessionID string) (*SessionContext, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT session_id, tool, model, format, content_json, captured_at
		FROM session_contexts
		WHERE session_id = ?
	`, sessionID)
	var context SessionContext
	var model sql.NullString
	if err := row.Scan(&context.SessionID, &context.Tool, &model, &context.Format, &context.ContentJSON, &context.CapturedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find session context: %w", err)
	}
	if model.Valid {
		context.Model = &model.String
	}
	return &context, nil
}

func (s *Store) CountCursorSessions(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM sessions WHERE source = 'cursor'`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count cursor sessions: %w", err)
	}
	return count, nil
}

func (s *Store) CountCursorMessages(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM cursor_messages`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count cursor messages: %w", err)
	}
	return count, nil
}

func (s *Store) sessionIDForRequest(ctx context.Context, requestID string) (string, error) {
	var sessionID string
	err := s.db.QueryRowContext(ctx, `SELECT session_id FROM requests WHERE id = ?`, requestID).Scan(&sessionID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("find request session: %w", err)
	}
	return sessionID, nil
}

func (s *Store) refreshSessionUsage(ctx context.Context, sessionID string) error {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil
	}
	models, err := s.sessionModels(ctx, sessionID)
	if err != nil {
		return err
	}
	modelsJSON, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("marshal session models: %w", err)
	}
	usage, err := s.sessionTokenUsage(ctx, sessionID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE sessions
		SET models_json = ?,
			input_tokens = ?,
			output_tokens = ?,
			cache_read_tokens = ?,
			cache_write_tokens = ?
		WHERE id = ?
	`, string(modelsJSON), usage.InputTokens, usage.OutputTokens, usage.CacheReadTokens, usage.CacheWriteTokens, sessionID)
	if err != nil {
		return fmt.Errorf("update session usage: %w", err)
	}
	return nil
}

func (s *Store) sessionModels(ctx context.Context, sessionID string) ([]string, error) {
	var models []string
	seen := map[string]struct{}{}
	add := func(model string) {
		model = strings.TrimSpace(model)
		if model == "" {
			return
		}
		if _, ok := seen[model]; ok {
			return
		}
		seen[model] = struct{}{}
		models = append(models, model)
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT model
		FROM requests
		WHERE session_id = ? AND model IS NOT NULL AND TRIM(model) != ''
		ORDER BY created_at ASC, id ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list request models: %w", err)
	}
	for rows.Next() {
		var model string
		if err := rows.Scan(&model); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan request model: %w", err)
		}
		add(model)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close request models: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate request models: %w", err)
	}

	rows, err = s.db.QueryContext(ctx, `
		SELECT model
		FROM session_contexts
		WHERE session_id = ? AND model IS NOT NULL AND TRIM(model) != ''
		ORDER BY captured_at ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list session context models: %w", err)
	}
	for rows.Next() {
		var model string
		if err := rows.Scan(&model); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan session context model: %w", err)
		}
		add(model)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close session context models: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate session context models: %w", err)
	}

	rows, err = s.db.QueryContext(ctx, `
		SELECT raw_json
		FROM cursor_messages
		WHERE session_id = ?
		ORDER BY created_at ASC, id ASC
	`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list cursor message models: %w", err)
	}
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			rows.Close()
			return nil, fmt.Errorf("scan cursor message model: %w", err)
		}
		for _, model := range modelsFromRawJSON(raw) {
			add(model)
		}
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close cursor message models: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cursor message models: %w", err)
	}
	return models, nil
}

func (s *Store) sessionTokenUsage(ctx context.Context, sessionID string) (tokenUsage, error) {
	var total tokenUsage
	rows, err := s.db.QueryContext(ctx, `
		SELECT resp.input_tokens, resp.output_tokens, resp.cache_read_tokens, resp.cache_write_tokens, resp.response_body
		FROM requests req
		JOIN responses resp ON resp.request_id = req.id
		WHERE req.session_id = ?
		ORDER BY resp.created_at ASC, resp.id ASC
	`, sessionID)
	if err != nil {
		return total, fmt.Errorf("list response usage: %w", err)
	}
	for rows.Next() {
		var input, output, cacheRead, cacheWrite sql.NullInt64
		var body []byte
		if err := rows.Scan(&input, &output, &cacheRead, &cacheWrite, &body); err != nil {
			rows.Close()
			return total, fmt.Errorf("scan response usage: %w", err)
		}
		var parsed tokenUsage
		if !input.Valid && !output.Valid {
			parsed = responseBodyUsage(body)
		}
		total.InputTokens += usageColumn(input, parsed.InputTokens)
		total.OutputTokens += usageColumn(output, parsed.OutputTokens)
		total.CacheReadTokens += usageColumn(cacheRead, parsed.CacheReadTokens)
		total.CacheWriteTokens += usageColumn(cacheWrite, parsed.CacheWriteTokens)
	}
	if err := rows.Close(); err != nil {
		return total, fmt.Errorf("close response usage: %w", err)
	}
	if err := rows.Err(); err != nil {
		return total, fmt.Errorf("iterate response usage: %w", err)
	}

	rows, err = s.db.QueryContext(ctx, `
		SELECT input_tokens, output_tokens
		FROM cursor_messages
		WHERE session_id = ?
	`, sessionID)
	if err != nil {
		return total, fmt.Errorf("list cursor message usage: %w", err)
	}
	for rows.Next() {
		var input, output sql.NullInt64
		if err := rows.Scan(&input, &output); err != nil {
			rows.Close()
			return total, fmt.Errorf("scan cursor message usage: %w", err)
		}
		total.InputTokens += usageColumn(input, 0)
		total.OutputTokens += usageColumn(output, 0)
	}
	if err := rows.Close(); err != nil {
		return total, fmt.Errorf("close cursor message usage: %w", err)
	}
	if err := rows.Err(); err != nil {
		return total, fmt.Errorf("iterate cursor message usage: %w", err)
	}
	return total, nil
}

func usageColumn(value sql.NullInt64, fallback int) int {
	if value.Valid {
		return int(value.Int64)
	}
	return fallback
}

func modelsFromRawJSON(raw []byte) []string {
	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	var models []string
	collectModels(payload, "", &models)
	return models
}

func collectModels(value any, key string, models *[]string) {
	switch v := value.(type) {
	case map[string]any:
		for childKey, childValue := range v {
			collectModels(childValue, childKey, models)
		}
	case []any:
		for _, childValue := range v {
			collectModels(childValue, key, models)
		}
	case string:
		if isModelKey(key) && strings.TrimSpace(v) != "" {
			*models = append(*models, v)
		}
	}
}

func isModelKey(key string) bool {
	key = strings.ToLower(strings.TrimSpace(key))
	switch key {
	case "model", "modelid", "model_id", "modelname", "model_name", "selectedmodel", "selected_model":
		return true
	default:
		return false
	}
}

type tokenUsage struct {
	InputTokens      int
	OutputTokens     int
	CacheReadTokens  int
	CacheWriteTokens int
}

func (u tokenUsage) total() int {
	return u.InputTokens + u.OutputTokens + u.CacheReadTokens + u.CacheWriteTokens
}

func (s *Store) AgentLedgerSummary(ctx context.Context, agent string) (AgentLedgerSummary, error) {
	summary := AgentLedgerSummary{Agent: agent}
	var where string
	switch agent {
	case "cursor":
		where = `COALESCE(s.source, '') = 'cursor'`
	case "codex":
		where = `(LOWER(COALESCE(s.command, '')) LIKE '%codex%' OR LOWER(COALESCE(s.process_name, '')) LIKE '%codex%')`
	case "claude":
		where = `(LOWER(COALESCE(s.command, '')) LIKE '%claude%' OR LOWER(COALESCE(s.process_name, '')) LIKE '%claude%')`
	default:
		return summary, fmt.Errorf("unsupported agent %q", agent)
	}

	row := s.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT
			COALESCE(
				(SELECT s2.cwd FROM sessions s2 WHERE %s AND s2.cwd != '' ORDER BY COALESCE(s2.last_seen_at, s2.ended_at, s2.created_at) DESC LIMIT 1),
				''
			),
			COALESCE(COUNT(DISTINCT s.id), 0),
			COALESCE(SUM(COALESCE(resp.input_tokens, 0) + COALESCE(resp.output_tokens, 0)), 0),
			COALESCE(COUNT(DISTINCT COALESCE(NULLIF(s.repo_root, ''), s.cwd)), 0),
			MAX(COALESCE(s.last_seen_at, s.ended_at, s.created_at))
		FROM sessions s
		LEFT JOIN requests req ON req.session_id = s.id
		LEFT JOIN responses resp ON resp.request_id = req.id
		WHERE %s
	`, strings.ReplaceAll(where, "s.", "s2."), where))

	var filepath string
	var lastSeen sql.NullInt64
	err := row.Scan(&filepath, &summary.Calls, &summary.Tokens, &summary.Files, &lastSeen)
	if err == sql.ErrNoRows {
		return summary, nil
	}
	if err != nil {
		return summary, fmt.Errorf("agent ledger summary %s: %w", agent, err)
	}
	summary.Filepath = filepath
	if lastSeen.Valid {
		value := lastSeen.Int64
		summary.LastSeenAt = &value
	}
	missingTokens, err := s.agentLedgerMissingTokens(ctx, where)
	if err != nil {
		return summary, err
	}
	summary.Tokens += missingTokens
	return summary, nil
}

func (s *Store) agentLedgerMissingTokens(ctx context.Context, where string) (int, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT resp.response_body
		FROM sessions s
		JOIN requests req ON req.session_id = s.id
		JOIN responses resp ON resp.request_id = req.id
		WHERE %s
			AND resp.input_tokens IS NULL
			AND resp.output_tokens IS NULL
	`, where))
	if err != nil {
		return 0, fmt.Errorf("agent ledger missing tokens: %w", err)
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		var body []byte
		if err := rows.Scan(&body); err != nil {
			return 0, fmt.Errorf("scan missing token response: %w", err)
		}
		total += responseBodyTokens(body)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterate missing token responses: %w", err)
	}
	return total, nil
}

func responseBodyTokens(body []byte) int {
	return responseBodyUsage(body).total()
}

func responseBodyUsage(body []byte) tokenUsage {
	var latest tokenUsage
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
		if line == "" || line == "[DONE]" || !strings.HasPrefix(line, "{") {
			continue
		}
		if usage := responseJSONUsage([]byte(line)); usage.total() > 0 {
			latest = usage
		}
	}
	if latest.total() > 0 {
		return latest
	}
	return responseJSONUsage(body)
}

func responseJSONTokens(body []byte) int {
	return responseJSONUsage(body).total()
}

func responseJSONUsage(body []byte) tokenUsage {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return tokenUsage{}
	}
	if response, _ := payload["response"].(map[string]any); response != nil {
		if usage := usageFromValue(response["usage"]); usage.total() > 0 {
			return usage
		}
	}
	return usageFromValue(payload["usage"])
}

func usageTokens(value any) int {
	return usageFromValue(value).total()
}

func usageFromValue(value any) tokenUsage {
	usage, _ := value.(map[string]any)
	if usage == nil {
		return tokenUsage{}
	}
	out := tokenUsage{
		InputTokens:      usageTokenValue(usage, "input_tokens", "prompt_tokens"),
		OutputTokens:     usageTokenValue(usage, "output_tokens", "completion_tokens"),
		CacheReadTokens:  usageTokenValue(usage, "cache_read_tokens", "cache_read_input_tokens"),
		CacheWriteTokens: usageTokenValue(usage, "cache_write_tokens", "cache_creation_input_tokens"),
	}
	if out.CacheReadTokens == 0 {
		for _, key := range []string{"prompt_tokens_details", "input_tokens_details"} {
			if details, _ := usage[key].(map[string]any); details != nil {
				if cached := usageTokenValue(details, "cached_tokens"); cached > 0 {
					out.CacheReadTokens = cached
					break
				}
			}
		}
	}
	return out
}

func usageTokenValue(usage map[string]any, keys ...string) int {
	for _, key := range keys {
		switch value := usage[key].(type) {
		case float64:
			return int(value)
		case int:
			return value
		}
	}
	return 0
}
