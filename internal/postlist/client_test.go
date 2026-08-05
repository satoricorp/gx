package postlist

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpsertIdentity(t *testing.T) {
	var gotMethod string
	var gotConflict string
	var gotPrefer string
	var gotAuth string
	var payload []Identity

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotConflict = r.URL.Query().Get("on_conflict")
		gotPrefer = r.Header.Get("prefer")
		gotAuth = r.Header.Get("authorization")
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		apiKey:     "secret",
		httpClient: server.Client(),
	}

	err := client.UpsertIdentity(context.Background(), Identity{
		Name:      "Joe Example",
		Email:     "joe@example.com",
		TLVersion: "dev",
	})
	if err != nil {
		t.Fatalf("UpsertIdentity() unexpected error = %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method = %s, want POST", gotMethod)
	}
	if gotConflict != "email" {
		t.Fatalf("on_conflict = %q, want email", gotConflict)
	}
	if gotPrefer != "resolution=merge-duplicates,return=minimal" {
		t.Fatalf("prefer = %q", gotPrefer)
	}
	if gotAuth != "Bearer secret" {
		t.Fatalf("authorization = %q", gotAuth)
	}
	if len(payload) != 1 {
		t.Fatalf("payload len = %d, want 1", len(payload))
	}
	if payload[0].Name != "Joe Example" || payload[0].Email != "joe@example.com" || payload[0].Source != "lgtm" {
		t.Fatalf("unexpected payload: %+v", payload[0])
	}
}
