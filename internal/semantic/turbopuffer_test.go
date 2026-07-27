package semantic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTurboPufferClientUpsertsRows(t *testing.T) {
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-tpuf" {
			t.Fatalf("authorization = %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/v2/namespaces/gx-test" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewTurboPufferClient(Config{
		TurboPufferAPIKey:    "test-tpuf",
		TurboPufferBaseURL:   server.URL,
		TurboPufferNamespace: "gx-test",
		EmbeddingDimensions:  2,
	})
	err := client.Upsert(context.Background(), []VectorRow{{
		ID:     "gx-row",
		Vector: []float32{0.1, 0.2},
		Attributes: map[string]any{
			"text":         "alpha transcript",
			"session_id":   "session-one",
			"revision_id":  "change-one",
			"jj_change_id": "change-one",
		},
	}})
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if got["distance_metric"] != "cosine_distance" {
		t.Fatalf("distance_metric = %#v", got["distance_metric"])
	}
	rows, ok := got["upsert_rows"].([]any)
	if !ok || len(rows) != 1 {
		t.Fatalf("upsert_rows = %#v", got["upsert_rows"])
	}
	row := rows[0].(map[string]any)
	if row["id"] != "gx-row" || row["text"] != "alpha transcript" || row["session_id"] != "session-one" {
		t.Fatalf("row = %#v", row)
	}
	schema := got["schema"].(map[string]any)
	vectorSchema := schema["vector"].(map[string]any)
	if vectorSchema["type"] != "[2]f32" || vectorSchema["ann"] != true {
		t.Fatalf("vector schema = %#v", vectorSchema)
	}
	agentToolSchema := schema["agent_tool"].(map[string]any)
	if agentToolSchema["type"] != "string" || agentToolSchema["filterable"] != true {
		t.Fatalf("agent_tool schema = %#v, want filterable string", agentToolSchema)
	}
	// Both the new attribute and its legacy twin stay declared until a full
	// reindex retires jj_change_id.
	for _, field := range []string{"revision_id", "jj_change_id"} {
		fieldSchema, ok := schema[field].(map[string]any)
		if !ok || fieldSchema["type"] != "string" || fieldSchema["filterable"] != true {
			t.Fatalf("%s schema = %#v, want filterable string", field, schema[field])
		}
	}
}

func TestTurboPufferClientRejectsWrongDimension(t *testing.T) {
	client := NewTurboPufferClient(Config{EmbeddingDimensions: 2})
	err := client.Upsert(context.Background(), []VectorRow{{ID: "bad", Vector: []float32{0.1}}})
	if err == nil {
		t.Fatal("expected dimension error")
	}
}

func TestTurboPufferClientDeletesStaleCodeRows(t *testing.T) {
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/namespaces/gx-test" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewTurboPufferClient(Config{
		TurboPufferAPIKey:    "test-tpuf",
		TurboPufferBaseURL:   server.URL,
		TurboPufferNamespace: "gx-test",
		EmbeddingDimensions:  2,
	})
	if err := client.DeleteStaleCodeDocuments(context.Background(), "acme/widgets", "abc123"); err != nil {
		t.Fatalf("DeleteStaleCodeDocuments() error = %v", err)
	}
	filter, ok := got["delete_by_filter"].([]any)
	if !ok || len(filter) != 2 || filter[0] != "And" {
		t.Fatalf("delete_by_filter = %#v", got["delete_by_filter"])
	}
	conditions, ok := filter[1].([]any)
	if !ok || len(conditions) != 3 {
		t.Fatalf("conditions = %#v", filter[1])
	}
}
