package codereview

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReviewResourceRetrieverQueriesBroadAndFilteredResources(t *testing.T) {
	store := &recordingReviewResourceStore{
		rows: []reviewResourceRow{{
			"source_id":       "owasp-sql-injection",
			"title":           "OWASP SQL Injection Prevention Cheat Sheet",
			"url":             "https://example.test/sql",
			"category":        "database",
			"authority":       "standard",
			"evidence_level":  "E1",
			"language_tags":   []any{"sql"},
			"risk_tag_values": []any{"sql-injection"},
			"chunk_kind":      "source",
			"chunk_index":     float64(2),
			"text":            "Prefer parameterized queries and prepared statements.",
		}},
	}
	retriever := ReviewResourceRetriever{
		Embedder:  fakeReviewResourceEmbedder{vector: []float32{0.1, 0.2}},
		Store:     store,
		Namespace: "gx-review-knowledge",
		Limit:     6,
	}

	snippets, err := retriever.Retrieve(context.Background(), RetrieveInput{
		RepoRoot: "/missing-repo",
		Options:  Options{Scope: "security"},
		Facts: RepoFacts{
			Files:           []string{"internal/auth/session.go", "db/migrations/001_add_users.sql"},
			DependencyFiles: []string{"package.json"},
		},
		Plan: ReviewExecutionPlan{RunReviewResources: true},
	})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(store.requests) != 2 {
		t.Fatalf("queries = %d, want broad + filtered", len(store.requests))
	}
	if !filterContains(store.requests[0].Filters, `"source_kind","Eq","review_knowledge"`) {
		t.Fatalf("broad filter = %#v", store.requests[0].Filters)
	}
	filtered := mustReviewResourceJSON(store.requests[1].Filters)
	for _, want := range []string{
		`"category","In"`,
		`"language_tags","ContainsAny"`,
		`"risk_tag_values","ContainsAny"`,
		`"review_tag_values","ContainsAny"`,
		`"database"`,
		`"sql"`,
	} {
		if !strings.Contains(filtered, want) {
			t.Fatalf("filtered query missing %s in %s", want, filtered)
		}
	}
	if len(snippets) != 1 {
		t.Fatalf("snippets = %#v", snippets)
	}
	snippet := snippets[0]
	if snippet.Kind != "review_resource" || snippet.Source != "turbopuffer:gx-review-knowledge" {
		t.Fatalf("snippet = %#v", snippet)
	}
	if !strings.Contains(snippet.Text, "OWASP SQL Injection") || !strings.Contains(snippet.Text, "parameterized queries") {
		t.Fatalf("snippet text = %q", snippet.Text)
	}
}

func TestReviewResourceQueryTextUsesPatchAndDeepReviewIntents(t *testing.T) {
	patchOpts := normalizeOptions(Options{})
	patchSignals := reviewResourceSignalSet{
		Files:      []string{"internal/auth/session.go"},
		Languages:  []string{"go"},
		RiskTags:   riskTagsForReview([]string{"internal/auth/session.go"}, nil, patchOpts),
		Categories: categoriesForReview(patchOpts),
		Intents:    reviewResourceIntents(patchOpts),
	}
	patchQuery := reviewResourceQueryText(patchOpts, patchSignals)
	for _, want := range []string{
		"profile: patch_focused",
		"changed-line bug and regression checks",
		"security authn authz secret token and webhook verification checks",
		"race concurrency retry and idempotency checks",
		"production observability",
		"internal/auth/session.go",
	} {
		if !strings.Contains(patchQuery, want) {
			t.Fatalf("patch query missing %q in:\n%s", want, patchQuery)
		}
	}

	deepOpts := normalizeOptions(Options{Deep: true})
	deepSignals := reviewResourceSignalSet{
		Files:      []string{"internal/review/engine.go"},
		Languages:  []string{"go"},
		RiskTags:   riskTagsForReview([]string{"internal/review/engine.go"}, nil, deepOpts),
		Categories: categoriesForReview(deepOpts),
		Intents:    reviewResourceIntents(deepOpts),
	}
	deepQuery := reviewResourceQueryText(deepOpts, deepSignals)
	for _, want := range []string{
		"profile: deep_full_spectrum",
		"Module Interface Depth locality and adapter architecture checks",
		"performance scalability and resource usage checks",
		"supply-chain dependency CI and deployment checks",
		"documentation onboarding and REVIEW.md policy checks",
		"observability",
		"performance",
	} {
		if !strings.Contains(deepQuery, want) {
			t.Fatalf("deep query missing %q in:\n%s", want, deepQuery)
		}
	}
}

func TestReviewResourceQueryTextIncludesPromptedReviewIntent(t *testing.T) {
	opts := normalizeOptions(Options{Prompt: "review auth rollback risk"})
	signals := reviewResourceSignalSet{
		Files:      []string{"internal/auth/session.go"},
		Languages:  []string{"go"},
		RiskTags:   riskTagsForReview([]string{"internal/auth/session.go"}, nil, opts),
		Categories: categoriesForReview(opts),
		Intents:    reviewResourceIntents(opts),
	}
	query := reviewResourceQueryText(opts, signals)
	for _, want := range []string{
		"profile: prompt_directed",
		"review prompt: review auth rollback risk",
		"user review_prompt checks",
		"broader codebase impact checks for the prompted concern",
		"security",
		"observability",
	} {
		if !strings.Contains(query, want) {
			t.Fatalf("prompted query missing %q in:\n%s", want, query)
		}
	}
}

func TestTurboPufferReviewResourceStoreBuildsQueryPayload(t *testing.T) {
	var gotPath string
	var gotAuth string
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"rows": []map[string]any{{
				"source_id": "google-review-standard",
				"title":     "Google Review Standard",
				"text":      "Protect long-term code health.",
			}},
		})
	}))
	defer server.Close()

	store := turboPufferReviewResourceStore{
		apiKey:     "test-tpuf",
		baseURL:    server.URL,
		namespace:  "gx-review-knowledge",
		httpClient: server.Client(),
	}
	rows, err := store.Query(context.Background(), reviewResourceQuery{
		Vector:            []float32{0.1, 0.2},
		Limit:             4,
		Filters:           reviewResourceBaseFilter(),
		IncludeAttributes: []string{"title", "text"},
	})
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if gotPath != "/v2/namespaces/gx-review-knowledge/query" {
		t.Fatalf("path = %s", gotPath)
	}
	if gotAuth != "Bearer test-tpuf" {
		t.Fatalf("auth = %q", gotAuth)
	}
	if _, ok := got["rank_by"].([]any); !ok {
		t.Fatalf("rank_by = %#v", got["rank_by"])
	}
	limit, ok := got["limit"].(map[string]any)
	if !ok || limit["total"] != float64(4) {
		t.Fatalf("limit = %#v", got["limit"])
	}
	if !filterContains(got["filters"], `"source_kind","Eq","review_knowledge"`) {
		t.Fatalf("filters = %#v", got["filters"])
	}
	if len(rows) != 1 || stringValue(rows[0]["title"]) != "Google Review Standard" {
		t.Fatalf("rows = %#v", rows)
	}
}

func TestCompactContextSnippetsKeepsReviewResourcesNearTop(t *testing.T) {
	snippets := []ContextSnippet{
		{Kind: "module_file", Ref: "a.go", Text: strings.Repeat("a", 10)},
		{Kind: "module_file", Ref: "b.go", Text: strings.Repeat("b", 10)},
		{Kind: "review_resource", Ref: "owasp#1", Text: "owasp guidance"},
		{Kind: "repo_doc", Ref: "REVIEW.md", Text: "repo policy"},
	}
	got := compactContextSnippets(snippets, 2)
	if len(got) != 2 {
		t.Fatalf("got %d snippets", len(got))
	}
	if got[0].Kind != "repo_doc" || got[1].Kind != "review_resource" {
		t.Fatalf("compacted snippets = %#v", got)
	}
}

type fakeReviewResourceEmbedder struct {
	vector []float32
}

func (f fakeReviewResourceEmbedder) Embed(_ context.Context, inputs []string) ([][]float32, error) {
	if len(inputs) != 1 {
		return nil, nil
	}
	return [][]float32{f.vector}, nil
}

type recordingReviewResourceStore struct {
	requests []reviewResourceQuery
	rows     []reviewResourceRow
}

func (s *recordingReviewResourceStore) Query(_ context.Context, req reviewResourceQuery) ([]reviewResourceRow, error) {
	s.requests = append(s.requests, req)
	return s.rows, nil
}

func filterContains(filter any, want string) bool {
	return strings.Contains(mustReviewResourceJSON(filter), want)
}

func mustReviewResourceJSON(value any) string {
	body, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(body)
}
