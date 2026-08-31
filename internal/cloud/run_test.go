package cloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func saveRunTestCredentials(t *testing.T) {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken: "gho_run",
		CLISessionToken:   "gxcs_run",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
}

func TestReserveRunSendsRunKeyAndDecodesReservation(t *testing.T) {
	saveRunTestCredentials(t)

	var gotHeader string
	var gotBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/runs/reserve" {
			t.Fatalf("path = %s, want /v1/runs/reserve", r.URL.Path)
		}
		gotHeader = r.Header.Get(RunHeader)
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(RunReservation{
			Allowed: true, Reason: "free_run", Used: 3, Limit: 7, Remaining: 4,
			CheckoutURL: "https://gx.run/checkout",
		})
	}))
	defer server.Close()

	ctx := WithRunKey(context.Background(), "run_test-1")
	client := &Client{url: server.URL, http: server.Client()}
	reservation, err := client.ReserveRun(ctx)
	if err != nil {
		t.Fatalf("ReserveRun() error = %v", err)
	}
	if gotHeader != "run_test-1" {
		t.Fatalf("%s header = %q, want run_test-1", RunHeader, gotHeader)
	}
	if gotBody["run_key"] != "run_test-1" || gotBody["kind"] != "review" {
		t.Fatalf("request body = %v", gotBody)
	}
	if !reservation.Allowed || reservation.Remaining != 4 {
		t.Fatalf("reservation = %+v", reservation)
	}
	if note := reservation.FreeRunsNote(); !strings.Contains(note, "4 free runs left") {
		t.Fatalf("FreeRunsNote() = %q", note)
	}
}

func TestReserveRunMaps402ToPaymentRequired(t *testing.T) {
	saveRunTestCredentials(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
		_, _ = w.Write([]byte(`{"error":"payment_required","reason":"free_runs_exhausted","message":"You've used all 7 free gx runs. Subscribe to keep using gx Cloud AI: https://gx.run/checkout","checkout_url":"https://gx.run/checkout"}`))
	}))
	defer server.Close()

	ctx := WithRunKey(context.Background(), NewRunKey())
	client := &Client{url: server.URL, http: server.Client()}
	_, err := client.ReserveRun(ctx)
	if !IsPaymentRequired(err) {
		t.Fatalf("ReserveRun() error = %v, want PaymentRequiredError", err)
	}
	if !strings.Contains(err.Error(), "https://gx.run/checkout") {
		t.Fatalf("error = %q, want the checkout URL", err.Error())
	}
}

func TestBedrockFightCarriesRunKeyAndMaps402(t *testing.T) {
	saveRunTestCredentials(t)

	var gotHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get(RunHeader)
		w.WriteHeader(http.StatusPaymentRequired)
		_, _ = w.Write([]byte(`{"error":"payment_required","checkout_url":"https://gx.run/checkout"}`))
	}))
	defer server.Close()

	ctx := WithRunKey(context.Background(), "run_fight")
	client := &Client{url: server.URL, http: server.Client()}
	_, err := client.BedrockFight(ctx, BedrockFightRequest{Model: "m", Messages: []BedrockMessage{{Role: "user", Content: "hi"}}})
	if gotHeader != "run_fight" {
		t.Fatalf("%s header = %q, want run_fight", RunHeader, gotHeader)
	}
	if !IsPaymentRequired(err) {
		t.Fatalf("BedrockFight() error = %v, want PaymentRequiredError", err)
	}
	if !strings.Contains(err.Error(), "Subscribe: https://gx.run/checkout") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestRunKeyContext(t *testing.T) {
	if got := RunKeyFrom(context.Background()); got != "" {
		t.Fatalf("RunKeyFrom(empty) = %q", got)
	}
	ctx := WithRunKey(context.Background(), "  ")
	if got := RunKeyFrom(ctx); got != "" {
		t.Fatalf("blank run key should not be stored, got %q", got)
	}
	ctx = WithRunKey(context.Background(), NewRunKey())
	if got := RunKeyFrom(ctx); !strings.HasPrefix(got, "run_") {
		t.Fatalf("RunKeyFrom() = %q", got)
	}
}
