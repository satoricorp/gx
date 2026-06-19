package semantic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIEmbedderRequestsEmbeddings(t *testing.T) {
	var got openAIEmbeddingRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-openai" {
			t.Fatalf("authorization = %q", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/v1/embeddings" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"index": 0, "embedding": []float64{0.1, 0.2}},
				{"index": 1, "embedding": []float64{0.3, 0.4}},
			},
		})
	}))
	defer server.Close()

	embedder := NewOpenAIEmbedder(Config{
		OpenAIAPIKey:         "test-openai",
		OpenAIBaseURL:        server.URL,
		OpenAIEmbeddingModel: "text-embedding-3-small",
		EmbeddingDimensions:  2,
	})
	vectors, err := embedder.Embed(context.Background(), []string{"alpha", "beta"})
	if err != nil {
		t.Fatalf("Embed() error = %v", err)
	}
	if got.Model != "text-embedding-3-small" || got.Dimensions != 2 || got.EncodingFormat != "float" {
		t.Fatalf("request = %#v", got)
	}
	if len(got.Input) != 2 || got.Input[0] != "alpha" || got.Input[1] != "beta" {
		t.Fatalf("request inputs = %#v", got.Input)
	}
	if len(vectors) != 2 || len(vectors[0]) != 2 || vectors[0][0] != float32(0.1) {
		t.Fatalf("vectors = %#v", vectors)
	}
}
