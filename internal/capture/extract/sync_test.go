package extract_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/capture/extract"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/telemetry"
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

	t.Setenv("GX_HOME", t.TempDir())
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
		ID: "e1", RepoRoot: "/repo", RefRange: "main..HEAD", PayloadJSON: extractPayload,
	}); err != nil {
		t.Fatal(err)
	}
	sessionPayload, _ := json.Marshal(map[string]any{
		"sessionID": "s1",
		"tool":      "cursor",
	})
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID: "s1-row", SessionID: "s1", Tool: "cursor", PayloadJSON: sessionPayload,
	}); err != nil {
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
