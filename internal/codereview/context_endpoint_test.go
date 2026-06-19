package codereview

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpointContextRetrieverAddsIndexedSnippets(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		var req contextQueryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Scope != "architecture" || len(req.Hints) != 1 {
			t.Fatalf("request = %#v", req)
		}
		_ = json.NewEncoder(w).Encode(contextQueryResponse{
			Snippets: []ContextSnippet{{Kind: "session", Ref: "s1", Source: "indexed", Text: "prior session confusion"}},
		})
	}))
	defer server.Close()

	snippets, err := (EndpointContextRetriever{URL: server.URL, Token: "token-one", Client: server.Client()}).Retrieve(
		context.Background(),
		"/repo",
		Options{Scope: "architecture"},
		RepoFacts{},
		[]ReviewHint{{Kind: "deepening_candidate", Title: "heavy"}},
	)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if gotAuth != "Bearer token-one" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if len(snippets) != 1 || snippets[0].Source != "indexed" {
		t.Fatalf("snippets = %#v", snippets)
	}
}

func TestCompositeContextRetrieverIgnoresEndpointFailures(t *testing.T) {
	snippets, err := (CompositeContextRetriever{Retrievers: []ContextRetriever{
		fakeRetriever{snippets: []ContextSnippet{{Kind: "repo_doc", Ref: "README.md", Source: "local", Text: "readme"}}},
		EndpointContextRetriever{URL: "http://127.0.0.1:1"},
	}}).Retrieve(context.Background(), "/repo", Options{}, RepoFacts{}, nil)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 || snippets[0].Source != "local" {
		t.Fatalf("snippets = %#v", snippets)
	}
}
