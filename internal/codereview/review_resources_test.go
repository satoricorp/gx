package codereview

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
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
	if !containsString(store.requests[0].IncludeAttributes, "body") {
		t.Fatalf("include attributes = %#v, want body", store.requests[0].IncludeAttributes)
	}
	if !filterContains(store.requests[0].Filters, `"source_kind","Eq","review_corpus"`) {
		t.Fatalf("broad filter = %#v", store.requests[0].Filters)
	}
	filtered := mustReviewResourceJSON(store.requests[1].Filters)
	for _, want := range []string{
		`"tier","In"`,
		`"language_tags","ContainsAny"`,
		`"security"`,
		`"sql"`,
	} {
		if !strings.Contains(filtered, want) {
			t.Fatalf("filtered query missing %s in %s", want, filtered)
		}
	}
	// The corpus schema declares only tier and language_tags filterable among
	// the signal attributes; undeclared ones would 400 the whole query.
	for _, reject := range []string{`"framework_tags"`, `"risk_tag_values"`, `"review_tag_values"`} {
		if strings.Contains(filtered, reject) {
			t.Fatalf("filtered query has undeclared filter attribute %s in %s", reject, filtered)
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

// The tag strings here are the contract with scripts/review-knowledge/sources.yaml:
// reviewResourceSignalFilter matches them against language_tags, so a tag that
// does not appear verbatim under a `languages:` key can never match a corpus
// source, and that source stays reachable only through the unfiltered leg.
func TestLanguageTagsForFiles(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{name: "typescript implies javascript", files: []string{"web/app.tsx"}, want: []string{"javascript", "typescript"}},
		{name: "javascript", files: []string{"web/app.mjs"}, want: []string{"javascript"}},
		{name: "python", files: []string{"svc/main.py"}, want: []string{"python"}},
		{name: "go", files: []string{"internal/review/engine.go"}, want: []string{"go"}},
		{name: "rust", files: []string{"src/lib.rs"}, want: []string{"rust"}},
		{name: "java", files: []string{"src/Main.java"}, want: []string{"java"}},
		{name: "c", files: []string{"src/parse.c"}, want: []string{"c"}},
		{name: "cpp", files: []string{"src/parse.cpp"}, want: []string{"cpp"}},
		{name: "ambiguous header", files: []string{"src/parse.h"}, want: []string{"c", "cpp"}},
		{name: "csharp", files: []string{"src/Program.cs"}, want: []string{"csharp"}},
		{name: "sql", files: []string{"db/migrations/001_add_users.sql"}, want: []string{"sql"}},
		{name: "shell", files: []string{"scripts/deploy.sh"}, want: []string{"shell"}},
		{name: "shell basenames", files: []string{"Makefile", "justfile"}, want: []string{"shell"}},
		{name: "php", files: []string{"public/index.php"}, want: []string{"php"}},
		{name: "kotlin", files: []string{"app/Main.kt"}, want: []string{"kotlin"}},
		{name: "swift", files: []string{"ios/AppDelegate.swift"}, want: []string{"swift"}},
		{name: "ruby", files: []string{"app/models/user.rb"}, want: []string{"ruby"}},
		{name: "ruby rake task", files: []string{"lib/tasks/import.rake"}, want: []string{"ruby"}},
		{name: "ruby gemspec", files: []string{"gx.gemspec"}, want: []string{"ruby"}},
		{name: "ruby basenames", files: []string{"Rakefile", "Gemfile", "Gemfile.lock"}, want: []string{"ruby"}},
		{name: "dart", files: []string{"lib/main.dart"}, want: []string{"dart"}},
		{name: "terraform", files: []string{"infra/main.tf"}, want: []string{"terraform"}},
		{name: "terraform vars", files: []string{"infra/prod.tfvars"}, want: []string{"terraform"}},
		{name: "html", files: []string{"public/index.html"}, want: []string{"html"}},
		{name: "htm", files: []string{"public/legacy.htm"}, want: []string{"html"}},
		{name: "dockerfile", files: []string{"Dockerfile"}, want: []string{"docker", "shell"}},
		{name: "dockerfile variant", files: []string{"deploy/Dockerfile.prod"}, want: []string{"docker", "shell"}},
		{name: "dockerfile extension", files: []string{"deploy/api.dockerfile"}, want: []string{"docker", "shell"}},
		{name: "unmapped extension", files: []string{"README.md", "config.yaml"}, want: nil},
		// Kubernetes manifests are plain YAML with no distinguishing path or
		// extension, so the kubernetes corpus tag stays deliberately unmapped
		// here rather than guessed at from a filename.
		{name: "kubernetes manifest stays unmapped", files: []string{"k8s/deployment.yaml"}, want: nil},
		{name: "mixed change set", files: []string{"app/models/user.rb", "lib/main.dart", "infra/main.tf", "public/index.html"}, want: []string{"dart", "html", "ruby", "terraform"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := languageTagsForFiles(tt.files)
			if len(got) != len(tt.want) {
				t.Fatalf("languageTagsForFiles(%v) = %#v, want %#v", tt.files, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("languageTagsForFiles(%v) = %#v, want %#v", tt.files, got, tt.want)
				}
			}
		})
	}
}

// Guards the mapping above against the manifest itself: every tag the extension
// map can emit must exist verbatim as a `languages:` value in sources.yaml.
func TestLanguageTagsForFilesMatchManifestTags(t *testing.T) {
	manifest, err := os.ReadFile(filepath.Join("..", "..", "scripts", "review-knowledge", "sources.yaml"))
	if err != nil {
		t.Fatalf("read sources.yaml: %v", err)
	}
	declared := map[string]bool{}
	for _, match := range regexp.MustCompile(`languages:\s*\[([^\]]*)\]`).FindAllStringSubmatch(string(manifest), -1) {
		for _, tag := range strings.Split(match[1], ",") {
			if tag = strings.TrimSpace(tag); tag != "" {
				declared[tag] = true
			}
		}
	}
	if len(declared) == 0 {
		t.Fatalf("no languages: tags parsed from sources.yaml")
	}

	emitted := languageTagsForFiles([]string{
		"web/app.tsx", "web/app.mjs", "svc/main.py", "internal/review/engine.go",
		"src/lib.rs", "src/Main.java", "src/parse.c", "src/parse.cpp", "src/parse.h",
		"src/Program.cs", "db/001.sql", "scripts/deploy.sh", "public/index.php",
		"app/Main.kt", "ios/AppDelegate.swift", "app/models/user.rb", "lib/main.dart",
		"infra/main.tf", "public/index.html", "Dockerfile",
	})
	for _, tag := range emitted {
		if !declared[tag] {
			t.Errorf("languageTagsForFiles emits %q, which no sources.yaml source declares under languages:", tag)
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

// Every outcome of a corpus query has to leave a line in the evidence log. A
// review that skipped the corpus, one whose query failed, and one that read
// 15,999 chunks and found nothing used to be indistinguishable in the report:
// all three simply had no review-knowledge line.
func TestReviewResourceRetrieverRecordsEvidenceForEveryOutcome(t *testing.T) {
	rows := []reviewResourceRow{{
		"source_id": "owasp-sql-injection",
		"chunk_id":  "owasp-sql-injection::queries::0002",
		"body":      "Prefer parameterized queries and prepared statements.",
	}}
	input := func(plan ReviewExecutionPlan) RetrieveInput {
		return RetrieveInput{
			RepoRoot:     "/missing-repo",
			Options:      normalizeOptions(Options{Scope: "security"}),
			Plan:         plan,
			Facts:        RepoFacts{Files: []string{"internal/auth/session.go"}},
			ChangedFiles: []string{"internal/auth/session.go"},
		}
	}
	tests := []struct {
		name       string
		retriever  ReviewResourceRetriever
		in         RetrieveInput
		wantState  string
		wantDetail string
		wantErr    bool
	}{
		{
			name:      "rows found",
			retriever: ReviewResourceRetriever{Embedder: fakeReviewResourceEmbedder{vector: []float32{0.1}}, Store: &recordingReviewResourceStore{rows: rows}},
			in:        input(ReviewExecutionPlan{}),
			wantState: EvidenceOK,
		},
		{
			name:       "corpus matched nothing",
			retriever:  ReviewResourceRetriever{Embedder: fakeReviewResourceEmbedder{vector: []float32{0.1}}, Store: &recordingReviewResourceStore{}},
			in:         input(ReviewExecutionPlan{}),
			wantState:  EvidenceEmpty,
			wantDetail: "matched no rows",
		},
		{
			name:       "rows without body text",
			retriever:  ReviewResourceRetriever{Embedder: fakeReviewResourceEmbedder{vector: []float32{0.1}}, Store: &recordingReviewResourceStore{rows: []reviewResourceRow{{"source_id": "empty-source"}}}},
			in:         input(ReviewExecutionPlan{}),
			wantState:  EvidenceEmpty,
			wantDetail: "none carried body text",
		},
		{
			name:       "query failed",
			retriever:  ReviewResourceRetriever{Embedder: fakeReviewResourceEmbedder{vector: []float32{0.1}}, Store: failingReviewResourceStore{err: errors.New("status 400 Bad Request")}},
			in:         input(ReviewExecutionPlan{}),
			wantState:  EvidenceUnavailable,
			wantDetail: "400 Bad Request",
			wantErr:    true,
		},
		{
			name:       "embedding failed",
			retriever:  ReviewResourceRetriever{Embedder: failingReviewResourceEmbedder{err: errors.New("429 rate limited")}, Store: &recordingReviewResourceStore{rows: rows}},
			in:         input(ReviewExecutionPlan{}),
			wantState:  EvidenceUnavailable,
			wantDetail: "429 rate limited",
			wantErr:    true,
		},
		{
			name:       "skipped by triage",
			retriever:  ReviewResourceRetriever{Embedder: fakeReviewResourceEmbedder{vector: []float32{0.1}}, Store: &recordingReviewResourceStore{rows: rows}},
			in:         input(ReviewExecutionPlan{Triage: ChangeTriage{Class: "docs-only"}, RunReviewResources: false}),
			wantState:  EvidenceSkipped,
			wantDetail: "docs-only",
		},
		{
			name:       "switched off",
			retriever:  ReviewResourceRetriever{Disabled: "GX_REVIEW_RESOURCES=0"},
			in:         input(ReviewExecutionPlan{}),
			wantState:  EvidenceDisabled,
			wantDetail: "GX_REVIEW_RESOURCES=0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &EvidenceLog{}
			in := tt.in
			in.Evidence = log
			_, err := tt.retriever.Retrieve(context.Background(), in)
			if tt.wantErr != (err != nil) {
				t.Fatalf("Retrieve() error = %v, wantErr = %v", err, tt.wantErr)
			}
			statuses := log.Statuses()
			if len(statuses) != 1 {
				t.Fatalf("statuses = %#v, want exactly one review-knowledge line", statuses)
			}
			if statuses[0].Source != reviewKnowledgeEvidenceSource {
				t.Fatalf("source = %q", statuses[0].Source)
			}
			if statuses[0].State != tt.wantState {
				t.Fatalf("state = %q, want %q", statuses[0].State, tt.wantState)
			}
			if tt.wantDetail != "" && !strings.Contains(statuses[0].Detail, tt.wantDetail) {
				t.Fatalf("detail = %q, want it to mention %q", statuses[0].Detail, tt.wantDetail)
			}
		})
	}
}

// The narrowed query is a second read of the same source. When it fails after
// the broad one answered, the review is entitled to both facts: what it got,
// and that it did not get all of it.
func TestReviewResourceRetrieverReportsAPartialReadWhenTheNarrowedQueryFails(t *testing.T) {
	log := &EvidenceLog{}
	retriever := ReviewResourceRetriever{
		Embedder: fakeReviewResourceEmbedder{vector: []float32{0.1}},
		Store: &secondQueryFailingStore{rows: []reviewResourceRow{{
			"source_id": "owasp-sql-injection",
			"chunk_id":  "owasp-sql-injection::queries::0002",
			"body":      "Prefer parameterized queries.",
		}}, err: errors.New("status 400 Bad Request")},
	}
	snippets, err := retriever.Retrieve(context.Background(), RetrieveInput{
		RepoRoot:     "/missing-repo",
		Options:      normalizeOptions(Options{Scope: "security"}),
		Facts:        RepoFacts{Files: []string{"internal/auth/session.go"}},
		ChangedFiles: []string{"internal/auth/session.go"},
		Evidence:     log,
	})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 {
		t.Fatalf("snippets = %#v, want the broad query's result kept", snippets)
	}
	statuses := log.Statuses()
	if len(statuses) != 2 {
		t.Fatalf("statuses = %#v, want the answer and the failed narrowing", statuses)
	}
	warnings := EvidenceWarnings(statuses)
	if len(warnings) != 1 || !strings.Contains(warnings[0], "read in part") {
		t.Fatalf("warnings = %v, want a partial-read warning", warnings)
	}
}

type failingReviewResourceStore struct{ err error }

func (s failingReviewResourceStore) Query(context.Context, reviewResourceQuery) ([]reviewResourceRow, error) {
	return nil, s.err
}

// secondQueryFailingStore answers the broad query and fails the narrowed one.
type secondQueryFailingStore struct {
	rows    []reviewResourceRow
	err     error
	queries int
}

func (s *secondQueryFailingStore) Query(context.Context, reviewResourceQuery) ([]reviewResourceRow, error) {
	s.queries++
	if s.queries == 1 {
		return s.rows, nil
	}
	return nil, s.err
}

type failingReviewResourceEmbedder struct{ err error }

func (f failingReviewResourceEmbedder) Embed(context.Context, []string) ([][]float32, error) {
	return nil, f.err
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
