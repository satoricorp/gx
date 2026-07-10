package storage

import (
	"context"
	"database/sql"
	"testing"
)

func TestBackfillUsagePopulatesResponsesAndSessionRollups(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	// Simulate legacy capture rows written before usage extraction existed:
	// raw inserts bypass WriteRequest/WriteResponse so no rollup runs.
	if _, err := db.ExecContext(ctx, `
		INSERT INTO sessions (id, created_at, command, cwd, gx_version)
		VALUES ('legacy-session', 10, 'codex', '/repo', 'old')
	`); err != nil {
		t.Fatalf("insert session: %v", err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO requests (id, session_id, created_at, provider, endpoint, method, model, request_body, request_headers)
		VALUES
			('req-json', 'legacy-session', 11, 'openai', '/v1/responses', 'POST', 'gpt-5.1-code', X'7B7D', '{}'),
			('req-sse', 'legacy-session', 12, 'openai', '/v1/responses', 'POST', 'gpt-5.1-code', X'7B7D', '{}'),
			('req-models', 'legacy-session', 13, 'openai', '/v1/models', 'GET', NULL, X'7B7D', '{}')
	`); err != nil {
		t.Fatalf("insert requests: %v", err)
	}
	jsonBody := `{"id":"resp_1","object":"response","status":"completed","usage":{"input_tokens":56294,"input_tokens_details":{"cached_tokens":53120},"output_tokens":416,"total_tokens":56710}}`
	sseBody := "event: response.created\n" +
		"data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_2\",\"status\":\"in_progress\"}}\n\n" +
		"event: response.completed\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_2\",\"status\":\"completed\",\"usage\":{\"input_tokens\":100,\"input_tokens_details\":{\"cached_tokens\":80},\"output_tokens\":25}}}\n\n"
	modelsBody := `{"object":"list","data":[{"id":"gpt-5.1-code"}]}`
	insertResponse := `
		INSERT INTO responses (id, request_id, created_at, completed_at, status_code, response_body, response_headers, is_streaming, duration_ms)
		VALUES (?, ?, ?, ?, 200, ?, '{}', ?, 1)
	`
	if _, err := db.ExecContext(ctx, insertResponse, "resp-json", "req-json", 14, 15, []byte(jsonBody), 0); err != nil {
		t.Fatalf("insert json response: %v", err)
	}
	if _, err := db.ExecContext(ctx, insertResponse, "resp-sse", "req-sse", 16, 17, []byte(sseBody), 1); err != nil {
		t.Fatalf("insert sse response: %v", err)
	}
	if _, err := db.ExecContext(ctx, insertResponse, "resp-models", "req-models", 18, 19, []byte(modelsBody), 0); err != nil {
		t.Fatalf("insert models response: %v", err)
	}

	result, err := store.BackfillUsage(ctx)
	if err != nil {
		t.Fatalf("BackfillUsage() error = %v", err)
	}
	if result.ResponsesScanned != 2 || result.ResponsesUpdated != 2 {
		t.Fatalf("result = %+v, want 2 scanned / 2 updated", result)
	}
	if result.SessionsRefreshed != 1 {
		t.Fatalf("result = %+v, want 1 session refreshed", result)
	}

	assertResponseUsage := func(id string, input, output, cacheRead int, finish string) {
		t.Helper()
		var gotInput, gotOutput, gotCacheRead sql.NullInt64
		var gotFinish sql.NullString
		if err := db.QueryRowContext(ctx, `
			SELECT input_tokens, output_tokens, cache_read_tokens, finish_reason
			FROM responses WHERE id = ?
		`, id).Scan(&gotInput, &gotOutput, &gotCacheRead, &gotFinish); err != nil {
			t.Fatalf("select response %s: %v", id, err)
		}
		if !gotInput.Valid || int(gotInput.Int64) != input {
			t.Fatalf("%s input_tokens = %+v, want %d", id, gotInput, input)
		}
		if !gotOutput.Valid || int(gotOutput.Int64) != output {
			t.Fatalf("%s output_tokens = %+v, want %d", id, gotOutput, output)
		}
		if !gotCacheRead.Valid || int(gotCacheRead.Int64) != cacheRead {
			t.Fatalf("%s cache_read_tokens = %+v, want %d", id, gotCacheRead, cacheRead)
		}
		if !gotFinish.Valid || gotFinish.String != finish {
			t.Fatalf("%s finish_reason = %+v, want %q", id, gotFinish, finish)
		}
	}
	assertResponseUsage("resp-json", 56294, 416, 53120, "completed")
	assertResponseUsage("resp-sse", 100, 25, 80, "completed")

	var modelsInput sql.NullInt64
	if err := db.QueryRowContext(ctx, `SELECT input_tokens FROM responses WHERE id = 'resp-models'`).Scan(&modelsInput); err != nil {
		t.Fatalf("select models response: %v", err)
	}
	if modelsInput.Valid {
		t.Fatalf("resp-models input_tokens = %+v, want NULL", modelsInput)
	}

	var modelsJSON string
	var input, output, cacheRead, cacheWrite int
	if err := db.QueryRowContext(ctx, `
		SELECT models_json, input_tokens, output_tokens, cache_read_tokens, cache_write_tokens
		FROM sessions WHERE id = 'legacy-session'
	`).Scan(&modelsJSON, &input, &output, &cacheRead, &cacheWrite); err != nil {
		t.Fatalf("select session rollup: %v", err)
	}
	if modelsJSON != `["gpt-5.1-code"]` {
		t.Fatalf("models_json = %s", modelsJSON)
	}
	if input != 56394 || output != 441 || cacheRead != 53200 || cacheWrite != 0 {
		t.Fatalf("session tokens = %d/%d/%d/%d, want 56394/441/53200/0", input, output, cacheRead, cacheWrite)
	}

	// A second run has nothing left to scan or refresh.
	again, err := store.BackfillUsage(ctx)
	if err != nil {
		t.Fatalf("BackfillUsage() second run error = %v", err)
	}
	if again.ResponsesScanned != 0 || again.ResponsesUpdated != 0 || again.SessionsRefreshed != 0 {
		t.Fatalf("second run = %+v, want all zero", again)
	}
}
