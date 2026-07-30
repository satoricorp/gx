package telemetry_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/satoricorp/totality/internal/telemetry"
)

func TestNopClient(t *testing.T) {
	t.Setenv("TOTALITY_POSTHOG_KEY", "")
	client := telemetry.NewFromEnv()
	ctx := context.Background()
	client.EmitCaptureCoverage(ctx, telemetry.CaptureCoverageProps{HunkCoverage: 0.9})
	client.EmitMatchRate(ctx, telemetry.MatchRateProps{AgentPrecision: 0.8})
	client.EmitSessionUploaded(ctx, telemetry.SessionUploadedProps{SessionID: "s1"})
	client.EmitComposeRun(ctx, telemetry.ComposeRunProps{HunkCoverage: 0.5})
}

func TestClientImplImplementsInterface(t *testing.T) {
	var _ telemetry.Client = telemetry.NopClient{}
}

func TestClientSendsCaptureEventWithIdentity(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TOTALITY_POSTHOG_KEY", "test-posthog-key")

	requests := make(chan map[string]any, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/capture/" {
			t.Errorf("request path = %q, want /capture/", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("request method = %q, want POST", r.Method)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		requests <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	t.Setenv("TOTALITY_POSTHOG_HOST", server.URL+"/")

	client := telemetry.NewFromEnv()
	client.EmitCaptureCoverage(context.Background(), telemetry.CaptureCoverageProps{
		Repo:         "/tmp/repo",
		RefRange:     "main..HEAD",
		HunkCoverage: 0.75,
		Tier1:        2,
		Tier2:        1,
		Tools:        []string{"codex"},
	})

	var payload map[string]any
	select {
	case payload = <-requests:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for PostHog request")
	}

	if got := payload["api_key"]; got != "test-posthog-key" {
		t.Fatalf("api_key = %v, want test-posthog-key", got)
	}
	if got := payload["event"]; got != telemetry.EventCaptureCoverage {
		t.Fatalf("event = %v, want %s", got, telemetry.EventCaptureCoverage)
	}
	properties, ok := payload["properties"].(map[string]any)
	if !ok {
		t.Fatalf("properties = %T, want object", payload["properties"])
	}
	if properties["distinct_id"] == "" {
		t.Fatal("properties.distinct_id is empty")
	}
	if properties["machine_id"] == "" {
		t.Fatal("properties.machine_id is empty")
	}
	if properties["tx_version"] == "" {
		t.Fatal("properties.tx_version is empty")
	}
	if properties["entrypoint"] != "cli" {
		t.Fatalf("properties.entrypoint = %v, want cli", properties["entrypoint"])
	}
}
