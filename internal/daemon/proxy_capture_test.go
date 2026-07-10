package daemon

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

type staticSessionResolver struct {
	sessionID string
}

func (r staticSessionResolver) ResolveRequestSession(*http.Request) (RequestSession, error) {
	return RequestSession{ID: r.sessionID, Capture: true}, nil
}

func (staticSessionResolver) Active() int { return 1 }

type stubTransport struct {
	status      int
	contentType string
	body        string
}

func (t stubTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: t.status,
		Header:     http.Header{"Content-Type": []string{t.contentType}},
		Body:       io.NopCloser(strings.NewReader(t.body)),
	}, nil
}

func newCaptureTestProxy(t *testing.T, transport http.RoundTripper) (*Proxy, *sql.DB) {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if err := store.WriteSession(ctx, storage.Session{
		ID:        "capture-session",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       "/repo",
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("WriteSession() error = %v", err)
	}
	proxy := NewProxyWithResolver(staticSessionResolver{sessionID: "capture-session"}, store)
	proxy.httpClient = &http.Client{Transport: transport}
	return proxy, db
}

func TestProxyCapturesOpenAIResponsesUsage(t *testing.T) {
	upstreamBody := `{"id":"resp_1","object":"response","model":"gpt-5.1-codex","status":"completed",` +
		`"usage":{"input_tokens":56294,"input_tokens_details":{"cached_tokens":53120},` +
		`"output_tokens":416,"output_tokens_details":{"reasoning_tokens":162},"total_tokens":56710}}`
	proxy, db := newCaptureTestProxy(t, stubTransport{
		status:      http.StatusOK,
		contentType: "application/json",
		body:        upstreamBody,
	})

	req := httptest.NewRequest(http.MethodPost, "http://proxy/v1/responses", bytes.NewReader([]byte(`{"model":"gpt-5.1-codex"}`)))
	req.Header.Set("Authorization", "Bearer sk-test")
	recorder := httptest.NewRecorder()
	proxy.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("proxy status = %d", recorder.Code)
	}

	ctx := context.Background()
	var input, output, cacheRead sql.NullInt64
	var finish sql.NullString
	if err := db.QueryRowContext(ctx, `
		SELECT resp.input_tokens, resp.output_tokens, resp.cache_read_tokens, resp.finish_reason
		FROM responses resp
		JOIN requests req ON req.id = resp.request_id
		WHERE req.session_id = 'capture-session'
	`).Scan(&input, &output, &cacheRead, &finish); err != nil {
		t.Fatalf("select captured response: %v", err)
	}
	if !input.Valid || input.Int64 != 56294 {
		t.Fatalf("input_tokens = %+v, want 56294", input)
	}
	if !output.Valid || output.Int64 != 416 {
		t.Fatalf("output_tokens = %+v, want 416", output)
	}
	if !cacheRead.Valid || cacheRead.Int64 != 53120 {
		t.Fatalf("cache_read_tokens = %+v, want 53120", cacheRead)
	}
	if !finish.Valid || finish.String != "completed" {
		t.Fatalf("finish_reason = %+v, want completed", finish)
	}

	var modelsJSON string
	var sessionInput, sessionOutput, sessionCacheRead int
	if err := db.QueryRowContext(ctx, `
		SELECT models_json, input_tokens, output_tokens, cache_read_tokens
		FROM sessions WHERE id = 'capture-session'
	`).Scan(&modelsJSON, &sessionInput, &sessionOutput, &sessionCacheRead); err != nil {
		t.Fatalf("select session rollup: %v", err)
	}
	if modelsJSON != `["gpt-5.1-codex"]` {
		t.Fatalf("models_json = %s", modelsJSON)
	}
	if sessionInput != 56294 || sessionOutput != 416 || sessionCacheRead != 53120 {
		t.Fatalf("session tokens = %d/%d/%d, want 56294/416/53120", sessionInput, sessionOutput, sessionCacheRead)
	}
}

func TestProxyCapturesOpenAIResponsesStreamingUsage(t *testing.T) {
	stream := "event: response.created\n" +
		"data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_1\",\"status\":\"in_progress\"}}\n\n" +
		"event: response.completed\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_1\",\"status\":\"completed\"," +
		"\"usage\":{\"input_tokens\":100,\"input_tokens_details\":{\"cached_tokens\":80},\"output_tokens\":25}}}\n\n"
	proxy, db := newCaptureTestProxy(t, stubTransport{
		status:      http.StatusOK,
		contentType: "text/event-stream",
		body:        stream,
	})

	req := httptest.NewRequest(http.MethodPost, "http://proxy/v1/responses", bytes.NewReader([]byte(`{"model":"gpt-5.1-codex","stream":true}`)))
	req.Header.Set("Authorization", "Bearer sk-test")
	recorder := httptest.NewRecorder()
	proxy.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("proxy status = %d", recorder.Code)
	}

	var input, output, cacheRead sql.NullInt64
	var finish sql.NullString
	var isStreaming int
	if err := db.QueryRowContext(context.Background(), `
		SELECT resp.input_tokens, resp.output_tokens, resp.cache_read_tokens, resp.finish_reason, resp.is_streaming
		FROM responses resp
		JOIN requests req ON req.id = resp.request_id
		WHERE req.session_id = 'capture-session'
	`).Scan(&input, &output, &cacheRead, &finish, &isStreaming); err != nil {
		t.Fatalf("select captured response: %v", err)
	}
	if isStreaming != 1 {
		t.Fatalf("is_streaming = %d, want 1", isStreaming)
	}
	if !input.Valid || input.Int64 != 100 {
		t.Fatalf("input_tokens = %+v, want 100", input)
	}
	if !output.Valid || output.Int64 != 25 {
		t.Fatalf("output_tokens = %+v, want 25", output)
	}
	if !cacheRead.Valid || cacheRead.Int64 != 80 {
		t.Fatalf("cache_read_tokens = %+v, want 80", cacheRead)
	}
	if !finish.Valid || finish.String != "completed" {
		t.Fatalf("finish_reason = %+v, want completed", finish)
	}
}
