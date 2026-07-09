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
			"source_id":        "owasp-sql-injection",
			"publisher":        "OWASP",
			"title":            "OWASP SQL Injection Prevention Cheat Sheet",
			"source_url":       "https://cheatsheetseries.owasp.org/cheatsheets/SQL_Injection_Prevention_Cheat_Sheet.html",
			"tier":             "security",
			"authority":        float64(1),
			"precedence_group": "security-spine",
			"language_tags":    []any{"sql"},
			"risk_tag_values":  []any{"sql-injection"},
			"chunk_id":         "owasp-sql-injection::queries::0002",
			"chunk_index":      float64(2),
			"body":             "Prefer parameterized queries and prepared statements.",
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
		Options:  normalizeOptions(Options{Scope: "security"}),
		Facts: RepoFacts{
			Files:           []string{"internal/auth/session.go", "db/migrations/001_add_users.sql"},
			DependencyFiles: []string{"package.json"},
		},
		ChangedFiles: []string{"internal/auth/session.go", "db/migrations/001_add_users.sql"},
	})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(store.requests) != 2 {
		t.Fatalf("queries = %d, want broad + filtered", len(store.requests))
	}
	if !containsString(store.requests[0].IncludeAttributes, "publisher") {
		t.Fatalf("include attributes = %#v, want publisher", store.requests[0].IncludeAttributes)
	}
	if !filterContains(store.requests[0].Filters, `"source_kind","Eq","review_corpus"`) {
		t.Fatalf("broad filter = %#v", store.requests[0].Filters)
	}
	filtered := mustReviewResourceJSON(store.requests[1].Filters)
	for _, want := range []string{
		`"tier","In"`,
		`"language_tags","ContainsAny"`,
		`"risk_tag_values","ContainsAny"`,
		`"review_tag_values","ContainsAny"`,
		`"security"`,
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
	if snippet.Publisher != "OWASP" {
		t.Fatalf("snippet publisher = %q, want OWASP", snippet.Publisher)
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
		Text:              "code health",
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
	if _, ok := got["queries"].([]any); !ok {
		t.Fatalf("queries = %#v", got["queries"])
	}
	if _, ok := got["rerank_by"].([]any); !ok {
		t.Fatalf("rerank_by = %#v", got["rerank_by"])
	}
	queries := got["queries"].([]any)
	if len(queries) != 2 {
		t.Fatalf("queries = %#v", queries)
	}
	first, ok := queries[0].(map[string]any)
	if !ok {
		t.Fatalf("first query = %#v", queries[0])
	}
	limit, ok := first["limit"].(map[string]any)
	if !ok || limit["total"] != float64(4) {
		t.Fatalf("limit = %#v", first["limit"])
	}
	if !filterContains(first["filters"], `"source_kind","Eq","review_corpus"`) {
		t.Fatalf("filters = %#v", first["filters"])
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func mustReviewResourceJSON(value any) string {
	body, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(body)
}
