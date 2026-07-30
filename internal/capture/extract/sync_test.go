package extract_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/satoricorp/totality/internal/auth"
	"github.com/satoricorp/totality/internal/capture"
	"github.com/satoricorp/totality/internal/capture/extract"
	"github.com/satoricorp/totality/internal/capture/matcher"
	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/telemetry"
)

func TestSyncPending(t *testing.T) {
	var extractCount int
	var sessionCount int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sync-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		body, _ := io.ReadAll(r.Body)
		switch r.URL.Path {
		case "/v1/extracts":
			extractCount++
			var req extract.ExtractRequest
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if req.HeadCommit != "abc123" {
				http.Error(w, "bad head", http.StatusBadRequest)
				return
			}
		case "/v1/sessions":
			sessionCount++
		default:
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	t.Setenv("TOTALITY_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}

	extractPayload, _ := json.Marshal(map[string]any{
		"refRange":  "main..HEAD",
		"repoRoot":  "/repo",
		"head":      "abc123",
		"hunkLinks": []matcher.HunkLink{{HunkID: "x", Tier: 1, Confidence: 1, Authorship: "agent"}},
	})
	if err := stager.StageExtract(ctx, storage.StagedExtract{
		ID: "e1", RepoRoot: "/repo", RefRange: "main..HEAD", PayloadJSON: extractPayload, RevisionID: "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	sessionPayload, _ := json.Marshal(map[string]any{
		"sessionID": "s1",
		"tool":      "cursor",
		"events":    []capture.SessionEvent{{Tool: "cursor", Kind: capture.KindEdit, FilePath: "a.go", NewText: "package a\n"}},
	})
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID: "s1-row", SessionID: "s1", Tool: "cursor", PayloadJSON: sessionPayload, RevisionID: "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := stager.MarkExtractShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{}); err != nil {
		t.Fatal(err)
	}
	if _, err := stager.MarkSessionShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{}); err != nil {
		t.Fatal(err)
	}

	result, err := extract.SyncPending(ctx, stager, auth.Credentials{
		APIURL: srv.URL,
		Token:  "sync-token",
	}, telemetry.NopClient{})
	if err != nil {
		t.Fatal(err)
	}
	if result.ExtractsUploaded != 1 || result.SessionsUploaded != 1 {
		t.Fatalf("result = %+v", result)
	}
	if extractCount != 1 || sessionCount != 1 {
		t.Fatalf("server counts extract=%d session=%d", extractCount, sessionCount)
	}

	pending, err := stager.PendingExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 0 {
		t.Fatalf("expected drained extracts, got %d", len(pending))
	}
}

// TestSyncPendingReadsSourceFiles covers the pointer model's upload behavior
// for rows whose payload carries no events: the transcript on disk is re-read
// and re-parsed at upload time (with the row's fingerprint refreshed when the
// file grew), a legacy inline blob is parsed the same way, and one vanished
// source marks its row without sinking the batch.
func TestSyncPendingReadsSourceFiles(t *testing.T) {
	received := map[string]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sessions" {
			w.WriteHeader(http.StatusOK)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var req struct {
			SessionId string `json:"sessionId"`
			Content   string `json:"content"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		received[req.SessionId] = req.Content
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Setenv("TOTALITY_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}

	staleHash := storage.PayloadContentHash([]byte("older, shorter bytes"))
	sourcePath := filepath.Join(t.TempDir(), "session.jsonl")
	sourceContent := claudeEditLine("on-disk", "/repo/on-disk.go", "package ondisk\n")
	if err := os.WriteFile(sourcePath, []byte(sourceContent), 0o644); err != nil {
		t.Fatal(err)
	}
	sessionPayload := func(id string) []byte {
		payload, _ := json.Marshal(map[string]any{"sessionID": id, "tool": "claude"})
		return payload
	}
	rows := []storage.StagedSession{
		{
			ID: "row-on-disk", SessionID: "on-disk", Tool: "claude", PayloadJSON: sessionPayload("on-disk"),
			SourcePath: sourcePath, ContentHash: staleHash, SourceBytes: 20, RevisionID: "rev-1",
		},
		{
			ID: "row-missing", SessionID: "missing", Tool: "claude", PayloadJSON: sessionPayload("missing"),
			SourcePath: filepath.Join(t.TempDir(), "deleted.jsonl"), RevisionID: "rev-1",
		},
		{
			ID: "row-legacy", SessionID: "legacy", Tool: "claude", PayloadJSON: sessionPayload("legacy"),
			RawBlob: []byte(claudeEditLine("legacy", "/repo/legacy.go", "package legacy\n")), RevisionID: "rev-1",
		},
	}
	for _, row := range rows {
		if err := stager.StageSession(ctx, row); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := stager.MarkSessionShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{}); err != nil {
		t.Fatal(err)
	}

	result, err := extract.SyncPending(ctx, stager, auth.Credentials{APIURL: srv.URL, Token: "sync-token"}, telemetry.NopClient{})
	if err != nil {
		t.Fatal(err)
	}
	if result.SessionsUploaded != 2 || result.SessionErrors != 1 {
		t.Fatalf("result = %+v, want 2 uploads and 1 error for the missing source", result)
	}
	// The wire format is Totality's own line-oriented transcript, never the raw
	// transcript JSONL: the server promoter only parses "[tool] kind …" lines
	// and drops everything else, so raw bytes promote zero events.
	if got := received["on-disk"]; !strings.HasPrefix(got, "[claude] edit /repo/on-disk.go\n") {
		t.Fatalf("uploaded content = %q, want formatted events derived from the current file bytes", got)
	}
	if strings.Contains(received["on-disk"], `"type":"assistant"`) {
		t.Fatalf("uploaded content = %q, want no raw transcript JSONL on the wire", received["on-disk"])
	}
	if got := received["legacy"]; !strings.HasPrefix(got, "[claude] edit /repo/legacy.go\n") {
		t.Fatalf("legacy content = %q, want the inline blob parsed into formatted events", got)
	}
	if _, ok := received["missing"]; ok {
		t.Fatal("missing-source row must not upload")
	}

	pending, err := stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 1 || pending[0].ID != "row-missing" {
		t.Fatalf("pending after sync = %+v, want only the missing-source row", pending)
	}
	// The grown on-disk source refreshed the stored fingerprint at upload.
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var gotHash string
	var gotBytes int64
	if err := db.QueryRowContext(ctx,
		`SELECT content_hash, source_bytes FROM capture_sessions WHERE id = ?`, "row-on-disk",
	).Scan(&gotHash, &gotBytes); err != nil {
		t.Fatal(err)
	}
	fresh := storage.PayloadContentHash([]byte(sourceContent))
	if gotHash != fresh || gotBytes != int64(len(sourceContent)) {
		t.Fatalf("fingerprint = hash %q bytes %d, want %q / %d", gotHash, gotBytes, fresh, len(sourceContent))
	}
}

// claudeEditLine renders one Claude transcript line whose Write tool call
// produces exactly one edit event, so a parser sees real content in it.
func claudeEditLine(sessionID, filePath, content string) string {
	body, err := json.Marshal(map[string]any{
		"type":      "assistant",
		"sessionId": sessionID,
		"timestamp": "2026-07-26T10:01:00.000Z",
		"message": map[string]any{
			"role":  "assistant",
			"model": "claude-probe",
			"content": []map[string]any{{
				"type":  "tool_use",
				"id":    "t1",
				"name":  "Write",
				"input": map[string]any{"file_path": filePath, "content": content},
			}},
		},
	})
	if err != nil {
		panic(err)
	}
	return string(body) + "\n"
}

// TestSyncPendingUploadsPromotableRedactedContentForEventlessRow pins the whole
// contract of the eventless pointer row — the shape behind every observed
// `status 400: {"error":"sessionId, tool, and content are required"}`, staged
// with a payload whose events are null.
//
// Three things must hold at once, and shipping raw transcript bytes satisfied
// only the first:
//   - content is non-empty, so the server stops rejecting the body;
//   - content is Totality's line format, because the server promoter parses only
//     "[tool] kind …" lines and DELETEs the session's existing events before
//     inserting what it parsed — raw JSONL promotes zero events and destroys
//     whatever context was already there;
//   - content is redacted, because raw transcript bytes carry the secrets that
//     redact.Redact exists to strip.
//
// The row ordered ahead has a source that no longer exists, proving one
// vanished transcript marks its own row and the batch carries on.
func TestSyncPendingUploadsPromotableRedactedContentForEventlessRow(t *testing.T) {
	var bodies []extract.SessionRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var req extract.SessionRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// The real server's contract, reproduced.
		if req.SessionId == "" || req.Tool == "" || req.Content == "" {
			http.Error(w, `{"error":"sessionId, tool, and content are required"}`, http.StatusBadRequest)
			return
		}
		bodies = append(bodies, req)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Setenv("TOTALITY_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}

	const secret = "AKIAIOSFODNN7EXAMPLE"
	transcript := claudeEditLine("live", "/repo/live.go", "const key = \""+secret+"\"\n")
	sourcePath := filepath.Join(t.TempDir(), "live.jsonl")
	if err := os.WriteFile(sourcePath, []byte(transcript), 0o644); err != nil {
		t.Fatal(err)
	}
	// Exactly what pointer staging writes: no events in the payload at all.
	eventless := []byte(`{"sessionID":"live","tool":"claude","events":null}`)

	gonePath := filepath.Join(t.TempDir(), "gone.jsonl")
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID: "row-gone", SessionID: "gone", Tool: "claude", CreatedAt: 1,
		PayloadJSON: []byte(`{"sessionID":"gone","tool":"claude","events":null}`),
		SourcePath:  gonePath, RevisionID: "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID: "row-live", SessionID: "live", Tool: "claude", CreatedAt: 2,
		PayloadJSON: eventless, SourcePath: sourcePath, RevisionID: "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := stager.MarkSessionShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{}); err != nil {
		t.Fatal(err)
	}

	result, err := extract.SyncPending(ctx, stager, auth.Credentials{APIURL: srv.URL, Token: "t"}, telemetry.NopClient{})
	if err != nil {
		t.Fatal(err)
	}
	if result.SessionsUploaded != 1 || result.SessionErrors != 1 {
		t.Fatalf("result = %+v, want the live row uploaded after the missing one failed", result)
	}
	if len(bodies) != 1 {
		t.Fatalf("accepted bodies = %+v, want one", bodies)
	}
	content := bodies[0].Content
	if promotableEvents(content) == 0 {
		t.Fatalf("content = %q, want lines the server promoter can parse into events", content)
	}
	if !strings.HasPrefix(content, "[claude] edit /repo/live.go\n") {
		t.Fatalf("content = %q, want the formatted edit event", content)
	}
	if strings.Contains(content, secret) {
		t.Fatalf("content = %q, want the AWS key redacted before it leaves the machine", content)
	}
	if strings.Contains(content, `"type":"assistant"`) {
		t.Fatalf("content = %q, want no raw transcript JSONL on the wire", content)
	}

	// The vanished source degrades honestly: its row records why.
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var uploadError string
	var attempts int
	if err := db.QueryRowContext(ctx,
		`SELECT COALESCE(upload_error, ''), COALESCE(upload_attempts, 0) FROM capture_sessions WHERE id = ?`, "row-gone",
	).Scan(&uploadError, &attempts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(uploadError, gonePath) {
		t.Fatalf("upload_error = %q, want it to name the missing source %q", uploadError, gonePath)
	}
	if attempts != 1 {
		t.Fatalf("upload_attempts = %d, want the failed attempt counted so it cannot retry forever", attempts)
	}
}

// promotableEvents mirrors parseSessionContent in server/src/ingest/promote.ts:
// only lines matching "[tool] kind …" become session_events, and a body that
// yields zero of them wipes the session's context instead of adding to it.
func promotableEvents(content string) int {
	header := regexp.MustCompile(`^\[([^\]]+)\]\s+(\S+)(?:\s+(.*))?$`)
	count := 0
	for _, line := range strings.Split(content, "\n") {
		if header.MatchString(line) {
			count++
		}
	}
	return count
}

// TestSyncPendingCapsSessionContent bounds one upload. A transcript grows while
// its session runs and re-staging clears uploaded_at, so an uncapped row is
// re-sent whole on every push forever — one staged pointer in the field
// targeted a 30 MB file.
func TestSyncPendingCapsSessionContent(t *testing.T) {
	var uploaded string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var req extract.SessionRequest
		if err := json.Unmarshal(raw, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		uploaded = req.Content
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Setenv("TOTALITY_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}

	events := make([]capture.SessionEvent, 0, 400)
	for i := range 400 {
		events = append(events, capture.SessionEvent{
			Tool:     "claude",
			Kind:     capture.KindEdit,
			FilePath: fmt.Sprintf("big-%d.go", i),
			NewText:  strings.Repeat("x", 16*1024),
		})
	}
	payload, _ := json.Marshal(map[string]any{"sessionID": "huge", "tool": "claude", "events": events})
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID: "row-huge", SessionID: "huge", Tool: "claude", PayloadJSON: payload, RevisionID: "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := stager.MarkSessionShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{}); err != nil {
		t.Fatal(err)
	}
	if _, err := extract.SyncPending(ctx, stager, auth.Credentials{APIURL: srv.URL, Token: "t"}, telemetry.NopClient{}); err != nil {
		t.Fatal(err)
	}
	if len(uploaded) == 0 {
		t.Fatal("nothing uploaded")
	}
	if len(uploaded) > 2<<20 {
		t.Fatalf("uploaded %d bytes, want a bounded body", len(uploaded))
	}
	if !strings.Contains(uploaded, "[tx] truncated") {
		t.Fatalf("uploaded body was cut without saying so: %q", uploaded[max(0, len(uploaded)-120):])
	}
}

// TestShareableSessionsStopRetryingAfterMaxAttempts is the other half of making
// a failing upload visible: the 32 rows that 400'd in production were re-sent on
// every single push, forever, with nothing anywhere reporting it.
func TestShareableSessionsStopRetryingAfterMaxAttempts(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]any{"sessionID": "doomed", "tool": "claude"})
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID: "row-doomed", SessionID: "doomed", Tool: "claude", PayloadJSON: payload, RevisionID: "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := stager.MarkSessionShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{}); err != nil {
		t.Fatal(err)
	}
	for i := range storage.MaxUploadAttempts {
		rows, err := stager.ShareableSessions(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(rows) != 1 {
			t.Fatalf("attempt %d: shareable rows = %d, want the row still eligible", i, len(rows))
		}
		if err := stager.SetSessionUploadError(ctx, "row-doomed", `status 400: {"error":"sessionId, tool, and content are required"}`); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := stager.ShareableSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("shareable rows = %d after %d failures, want the row retired instead of retried forever",
			len(rows), storage.MaxUploadAttempts)
	}
	failures, err := stager.UploadFailures(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if failures.Sessions != 1 || failures.ExhaustedRows != 1 {
		t.Fatalf("failures = %+v, want the retired row reported, not forgotten", failures)
	}
	if !strings.Contains(failures.LastError, "sessionId, tool, and content are required") {
		t.Fatalf("failures.LastError = %q, want the server's reason readable", failures.LastError)
	}
}
