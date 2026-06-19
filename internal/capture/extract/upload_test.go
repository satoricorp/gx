package extract_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/satoricorp/gx/internal/capture/extract"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/telemetry"
)

func TestClientUploadRun(t *testing.T) {
	var extractBody extract.ExtractRequest
	var sessionBodies []extract.SessionRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		body, _ := io.ReadAll(r.Body)
		switch r.URL.Path {
		case "/v1/extracts":
			if err := json.Unmarshal(body, &extractBody); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"eventId":"e1","hunkLinksInserted":1}`))
		case "/v1/sessions":
			var s extract.SessionRequest
			if err := json.Unmarshal(body, &s); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			sessionBodies = append(sessionBodies, s)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"sessionRawId":"r1","promotionStatus":"done"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client := extract.NewClient(srv.URL, "test-token")
	err := client.UploadRun(
		context.Background(),
		"/repo",
		"main..HEAD",
		"abc123",
		"0.1.0",
		[]matcher.HunkLink{{
			HunkID:     "abc:file.go:1-2",
			SessionID:  "s1",
			Tier:       1,
			Confidence: 1,
			Authorship: "agent",
			Tool:       "cursor",
		}},
		nil,
		[]extract.SessionPayload{{
			SessionID: "s1",
			Tool:      "cursor",
		}},
		telemetry.NopClient{},
	)
	if err != nil {
		t.Fatalf("UploadRun: %v", err)
	}
	if extractBody.HeadCommit != "abc123" {
		t.Fatalf("head commit = %q", extractBody.HeadCommit)
	}
	if len(sessionBodies) != 1 {
		t.Fatalf("sessions uploaded = %d", len(sessionBodies))
	}
}
