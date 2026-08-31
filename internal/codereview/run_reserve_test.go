package codereview

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/cloud"
)

// reserveCloudRun is the up-front metering check. These tests pin the one
// rule that matters: only a spent allowance (402) stops the review; every
// other answer lets it proceed, because the server gates each model call too.

func withCloudForReserve(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	t.Setenv(bedrockDirectEnvVar, "")
	t.Setenv("GX_HOME", t.TempDir())
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_reserve",
		CLISessionToken:   "gxcs_reserve",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	t.Setenv("GX_CLOUD_URL", server.URL)
}

func TestReserveCloudRunStopsOnSpentRuns(t *testing.T) {
	withCloudForReserve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/runs/reserve" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusPaymentRequired)
		_, _ = w.Write([]byte(`{"error":"payment_required","reason":"free_runs_exhausted","message":"You've used all 7 free gx runs. Subscribe to keep using gx Cloud AI: https://gx.run/checkout","checkout_url":"https://gx.run/checkout"}`))
	})

	ctx := cloud.WithRunKey(context.Background(), cloud.NewRunKey())
	err := reserveCloudRun(ctx, Options{})
	if !cloud.IsPaymentRequired(err) {
		t.Fatalf("reserveCloudRun() error = %v, want PaymentRequiredError", err)
	}
	if !strings.Contains(err.Error(), "https://gx.run/checkout") {
		t.Fatalf("error should carry the checkout URL, got %q", err.Error())
	}
}

func TestReserveCloudRunProceedsAndReportsRemaining(t *testing.T) {
	var gotKey string
	withCloudForReserve(t, func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get(cloud.RunHeader)
		_, _ = w.Write([]byte(`{"allowed":true,"reason":"free_run","used":2,"limit":7,"remaining":5,"checkout_url":"https://gx.run/checkout"}`))
	})

	var progress bytes.Buffer
	ctx := cloud.WithRunKey(context.Background(), "run_progress")
	if err := reserveCloudRun(ctx, Options{ProgressWriter: &progress}); err != nil {
		t.Fatalf("reserveCloudRun() error = %v", err)
	}
	if gotKey != "run_progress" {
		t.Fatalf("%s = %q", cloud.RunHeader, gotKey)
	}
	if !strings.Contains(progress.String(), "5 free runs left") {
		t.Fatalf("progress = %q, want the remaining-runs note", progress.String())
	}
}

func TestReserveCloudRunIgnoresTransportFailures(t *testing.T) {
	withCloudForReserve(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"entitlement_unavailable"}`))
	})

	ctx := cloud.WithRunKey(context.Background(), cloud.NewRunKey())
	if err := reserveCloudRun(ctx, Options{}); err != nil {
		t.Fatalf("a 503 must not stop the review (the server gates each call), got %v", err)
	}
}

func TestReserveCloudRunSkipsDirectAWS(t *testing.T) {
	called := false
	withCloudForReserve(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusPaymentRequired)
	})
	t.Setenv(bedrockDirectEnvVar, "1")

	ctx := cloud.WithRunKey(context.Background(), cloud.NewRunKey())
	if err := reserveCloudRun(ctx, Options{}); err != nil {
		t.Fatalf("direct AWS reviews are not metered, got %v", err)
	}
	if called {
		t.Fatal("direct AWS review must not reach gx Cloud")
	}
}
