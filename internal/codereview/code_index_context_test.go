package codereview

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/satoricorp/lgtm/internal/semantic"
)

// fakeIndexStore serves probes and query rows from memory so retrieval, fusion
// and evidence reporting can be tested without TurboPuffer.
type fakeIndexStore struct {
	probes map[string]indexProbe
	rows   map[string][]indexRow
	errs   map[string]error

	// Namespaces are queried concurrently, so the record of what was asked is
	// guarded.
	mu      sync.Mutex
	queries map[string]indexQuery
}

func (s *fakeIndexStore) query(namespace string) indexQuery {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queries[namespace]
}

func (s *fakeIndexStore) Probe(_ context.Context, namespace string) indexProbe {
	if probe, ok := s.probes[namespace]; ok {
		probe.Namespace = namespace
		if probe.Filterable == nil {
			probe.Filterable = map[string]bool{}
		}
		return probe
	}
	return indexProbe{Namespace: namespace, Filterable: map[string]bool{}}
}

func (s *fakeIndexStore) Query(_ context.Context, namespace string, req indexQuery) ([]indexRow, error) {
	s.mu.Lock()
	if s.queries == nil {
		s.queries = map[string]indexQuery{}
	}
	s.queries[namespace] = req
	s.mu.Unlock()
	if err := s.errs[namespace]; err != nil {
		return nil, err
	}
	return s.rows[namespace], nil
}

func staticEmbedder(dimensions int) embedderFactory {
	return func(width int) (reviewResourceEmbedder, string, bool) {
		if width != dimensions {
			return nil, "", false
		}
		return embedderFunc(func(context.Context, []string) ([][]float32, error) {
			return [][]float32{make([]float32, width)}, nil
		}), "fake", true
	}
}

type embedderFunc func(ctx context.Context, inputs []string) ([][]float32, error)

func (f embedderFunc) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	return f(ctx, inputs)
}

func codeIndexTestInput() RetrieveInput {
	return RetrieveInput{
		RepoRoot:     "/repo",
		Options:      Options{Scope: "security"},
		ChangedFiles: []string{"internal/codereview/engine.go"},
		DiffSnippets: []DiffSnippet{{
			File: "internal/codereview/engine.go",
			Diff: "@@\n+func contextRetrieverFromEnv() ContextRetriever {\n-\tolderHelper()\n",
		}},
		Evidence: &EvidenceLog{},
	}
}

func TestCodeIndexRetrieverFusesNamespacesAndBuildsHybridLegs(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"lgtm-org-repo-v2": {
				Exists: true, Dimensions: 3072, BodyField: "text", SymbolField: "symbol",
				Filterable: map[string]bool{"source_kind": true},
			},
			"repo-owner-repo": {
				Exists: true, Dimensions: 1536, BodyField: "content", SymbolField: "symbol",
			},
		},
		rows: map[string][]indexRow{
			"lgtm-org-repo-v2": {
				{"file_path": "internal/codereview/engine.go", "start_line": 1.0, "end_line": 20.0, "text": "fresh chunk", "symbol_name": "Review"},
			},
			"repo-owner-repo": {
				{"file_path": "internal/codereview/other.go", "start_line": 5.0, "end_line": 9.0, "content": "console chunk", "symbol": "Other"},
				{"file_path": "internal/codereview/engine.go", "start_line": 1.0, "end_line": 20.0, "content": "a much longer console copy of the same chunk"},
			},
		},
	}
	in := codeIndexTestInput()
	retriever := CodeIndexRetriever{
		Store: store,
		Namespaces: []codeIndexTarget{
			{Namespace: "lgtm-org-repo-v2", Origin: "lgtm code index"},
			{Namespace: "repo-owner-repo", Origin: "lgtm Cloud code index"},
		},
		Limit:       10,
		EmbedderFor: func(width int) (reviewResourceEmbedder, string, bool) { return staticEmbedder(width)(width) },
	}
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}

	// Each namespace is queried with the field names its own schema declares.
	lgtmQuery := store.query("lgtm-org-repo-v2")
	if lgtmQuery.BodyField != "text" || lgtmQuery.SymbolField != "symbol" || len(lgtmQuery.Vector) != 3072 {
		t.Fatalf("lgtm query = %+v, want text/symbol legs and a 3072-wide vector", lgtmQuery)
	}
	if lgtmQuery.Legs() != 3 {
		t.Fatalf("lgtm query legs = %d, want 3", lgtmQuery.Legs())
	}
	consoleQuery := store.query("repo-owner-repo")
	if consoleQuery.BodyField != "content" || len(consoleQuery.Vector) != 1536 {
		t.Fatalf("console query = %+v, want a content leg and a 1536-wide vector", consoleQuery)
	}
	if strings.TrimSpace(consoleQuery.SymbolQuery) == "" {
		t.Fatal("symbol query is empty; identifiers from the diff never reached the lexical leg")
	}
	if !strings.Contains(consoleQuery.SymbolQuery, "contextRetrieverFromEnv") {
		t.Fatalf("symbol query = %q, want the changed identifier", consoleQuery.SymbolQuery)
	}
	// The lgtm namespace holds more than code, so it is filtered; the console one
	// does not declare source_kind and must not be filtered on it.
	if lgtmQuery.Filters == nil {
		t.Fatal("lgtm query has no source_kind filter")
	}
	if consoleQuery.Filters != nil {
		t.Fatalf("console query filters = %v, want none for a namespace without source_kind", consoleQuery.Filters)
	}

	// The same chunk from two namespaces fuses into one snippet, keeping the
	// longer body.
	if len(snippets) != 2 {
		t.Fatalf("snippets = %d, want 2 fused snippets: %#v", len(snippets), snippets)
	}
	if snippets[0].Ref != "internal/codereview/engine.go:1" {
		t.Fatalf("top snippet = %q, want the chunk both namespaces returned", snippets[0].Ref)
	}
	if !strings.Contains(snippets[0].Text, "a much longer console copy") {
		t.Fatalf("fused snippet kept the shorter body: %q", snippets[0].Text)
	}
	if snippets[0].Kind != "indexed_code" {
		t.Fatalf("snippet kind = %q, want indexed_code", snippets[0].Kind)
	}

	statuses := in.Evidence.Statuses()
	if len(statuses) != 2 {
		t.Fatalf("statuses = %#v, want one per namespace", statuses)
	}
	for _, status := range statuses {
		if status.State != EvidenceOK {
			t.Fatalf("status %+v, want ok", status)
		}
	}
	if warnings := EvidenceWarnings(statuses); len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none", warnings)
	}
}

func TestCodeIndexRetrieverReportsMissingNamespace(t *testing.T) {
	store := &fakeIndexStore{}
	in := codeIndexTestInput()
	retriever := CodeIndexRetriever{
		Store: store,
		Namespaces: []codeIndexTarget{
			{Namespace: "lgtm-local-yeet-v2", Origin: "lgtm code index"},
			{Namespace: "repo-owner-yeet", Origin: "lgtm Cloud code index"},
		},
		EmbedderFor: func(int) (reviewResourceEmbedder, string, bool) { return nil, "", false },
	}
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 0 {
		t.Fatalf("snippets = %#v, want none", snippets)
	}
	// One warning for the source, not one per candidate namespace: the reader
	// needs to know the code index could not be read, and how many names were
	// tried on the way to finding that out is detail for the verbose listing.
	warnings := EvidenceWarnings(in.Evidence.Statuses())
	if len(warnings) != 1 {
		t.Fatalf("warnings = %v, want one warning for the code index as a whole", warnings)
	}
	// The warning has to say how to fix it, and the fix is the website. It used
	// to name `lgtm index`, which is a hidden maintenance command that fills one
	// developer's namespace from one developer's checkout — a worse, manual
	// copy of the index lgtm Cloud maintains from the GitHub App on merge.
	joined := strings.Join(warnings, " ")
	if !strings.Contains(joined, "https://lgtm.cx/repositories") {
		t.Fatalf("warnings = %v, want the missing index to say where to get it built", warnings)
	}
	if strings.Contains(joined, "lgtm index") {
		t.Fatalf("warnings = %v, want users sent to the console rather than the hidden `lgtm index`", warnings)
	}
	if len(in.Evidence.Statuses()) != 2 {
		t.Fatalf("statuses = %#v, want one per namespace so the verbose listing still names each", in.Evidence.Statuses())
	}
}

// TestCodeIndexAnsweredByOneNamespaceIsNotDegraded is the regression test for
// the review's most misleading line. Probing several candidate namespaces means
// most of them miss by design; a review whose index answered must never open
// with "review evidence unavailable … findings are based on the change and the
// checkout only".
func TestCodeIndexAnsweredByOneNamespaceIsNotDegraded(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"lgtm-local-satoricorp-yeet-v2": {Exists: true, Dimensions: 1536, BodyField: "text", SymbolField: "symbol"},
		},
		rows: map[string][]indexRow{
			"lgtm-local-satoricorp-yeet-v2": {
				{"file_path": "src/bridge.ts", "start_line": 1.0, "end_line": 20.0, "text": "chunk", "symbol_name": "Bridge"},
			},
		},
	}
	in := codeIndexTestInput()
	retriever := CodeIndexRetriever{
		Store: store,
		Namespaces: []codeIndexTarget{
			{Namespace: "lgtm-local-satoricorp-yeet-v2", Origin: semantic.NamespaceOriginPrimary},
			{Namespace: "lgtm-local-yeet-8d862445e7e4-v2", Origin: semantic.NamespaceOriginPreRemote},
		},
		Limit:       10,
		EmbedderFor: func(width int) (reviewResourceEmbedder, string, bool) { return staticEmbedder(width)(width) },
	}
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) == 0 {
		t.Fatal("no snippets from the namespace that exists")
	}
	if warnings := EvidenceWarnings(in.Evidence.Statuses()); len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none when the code index answered", warnings)
	}
	// The namespaces that missed are still on the record, named, so the verbose
	// evidence line says truthfully what was searched.
	var namespaces []string
	for _, status := range in.Evidence.Statuses() {
		namespaces = append(namespaces, status.Namespace)
	}
	if len(namespaces) != 2 {
		t.Fatalf("statuses = %v, want every searched namespace reported", namespaces)
	}
}

func TestCodeIndexRetrieverReportsQueryFailure(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"repo-owner-repo": {Exists: true, Dimensions: 1536, BodyField: "content", SymbolField: "symbol"},
		},
		errs: map[string]error{"repo-owner-repo": errors.New("status 503 Service Unavailable")},
	}
	in := codeIndexTestInput()
	retriever := CodeIndexRetriever{
		Store:       store,
		Namespaces:  []codeIndexTarget{{Namespace: "repo-owner-repo", Origin: "lgtm Cloud code index"}},
		EmbedderFor: func(width int) (reviewResourceEmbedder, string, bool) { return staticEmbedder(width)(width) },
	}
	if _, err := retriever.Retrieve(context.Background(), in); err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	warnings := EvidenceWarnings(in.Evidence.Statuses())
	if len(warnings) != 1 || !strings.Contains(warnings[0], "503") {
		t.Fatalf("warnings = %v, want the query failure surfaced", warnings)
	}
}

func TestCodeIndexRetrieverDegradesToLexicalWithoutAnEmbedder(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"repo-owner-repo": {Exists: true, Dimensions: 999, BodyField: "content", SymbolField: "symbol"},
		},
		rows: map[string][]indexRow{
			"repo-owner-repo": {{"file_path": "a.go", "start_line": 1.0, "end_line": 4.0, "content": "chunk"}},
		},
	}
	in := codeIndexTestInput()
	retriever := CodeIndexRetriever{
		Store:       store,
		Namespaces:  []codeIndexTarget{{Namespace: "repo-owner-repo", Origin: "lgtm Cloud code index"}},
		EmbedderFor: func(int) (reviewResourceEmbedder, string, bool) { return nil, "", false },
	}
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 {
		t.Fatalf("snippets = %#v, want the lexical legs to still answer", snippets)
	}
	query := store.query("repo-owner-repo")
	if len(query.Vector) != 0 || query.Legs() != 2 {
		t.Fatalf("query = %+v, want two lexical legs and no vector", query)
	}
	var sawNote bool
	for _, status := range in.Evidence.Statuses() {
		if strings.Contains(status.Detail, "lexical") {
			sawNote = true
		}
	}
	if !sawNote {
		t.Fatalf("statuses = %#v, want the missing vector leg reported", in.Evidence.Statuses())
	}
}

// A namespace that exists but exposes nothing this reader can query is not the
// same as one that answered and had nothing; reporting it as empty would credit
// it with an answer it never gave.
func TestCodeIndexRetrieverDistinguishesUnqueryableFromEmpty(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"unqueryable": {Exists: true, Dimensions: 999},
			"answered":    {Exists: true, Dimensions: 1536, BodyField: "content", SymbolField: "symbol"},
		},
	}
	in := codeIndexTestInput()
	retriever := CodeIndexRetriever{
		Store: store,
		Namespaces: []codeIndexTarget{
			{Namespace: "unqueryable", Origin: "lgtm code index"},
			{Namespace: "answered", Origin: "lgtm Cloud code index"},
		},
		EmbedderFor: func(int) (reviewResourceEmbedder, string, bool) { return nil, "", false },
	}
	if _, err := retriever.Retrieve(context.Background(), in); err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	byNamespace := map[string]EvidenceStatus{}
	for _, status := range in.Evidence.Statuses() {
		byNamespace[status.Namespace] = status
	}
	if byNamespace["unqueryable"].State != EvidenceUnavailable {
		t.Fatalf("unqueryable status = %+v, want unavailable", byNamespace["unqueryable"])
	}
	if byNamespace["answered"].State != EvidenceEmpty {
		t.Fatalf("answered status = %+v, want empty", byNamespace["answered"])
	}
}

func TestCodeIndexRetrieverReportsStaleIndex(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"repo-owner-repo": {
				Exists: true, Dimensions: 1536, BodyField: "content", SymbolField: "symbol",
				LastWrite: time.Now().Add(-30 * 24 * time.Hour),
			},
		},
		rows: map[string][]indexRow{
			"repo-owner-repo": {{"file_path": "a.go", "start_line": 1.0, "end_line": 4.0, "content": "chunk"}},
		},
	}
	in := codeIndexTestInput()
	retriever := CodeIndexRetriever{
		Store:       store,
		Namespaces:  []codeIndexTarget{{Namespace: "repo-owner-repo", Origin: "lgtm Cloud code index"}},
		EmbedderFor: func(width int) (reviewResourceEmbedder, string, bool) { return staticEmbedder(width)(width) },
	}
	if _, err := retriever.Retrieve(context.Background(), in); err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	statuses := in.Evidence.Statuses()
	if len(statuses) != 1 || !strings.Contains(statuses[0].Detail, "day(s) ago") {
		t.Fatalf("statuses = %#v, want the index age reported", statuses)
	}
}

func TestChangedIdentifiersRanksAndSplits(t *testing.T) {
	diffs := []DiffSnippet{{
		File: "internal/codereview/brief.go",
		Diff: strings.Join([]string{
			"--- a/internal/codereview/brief.go",
			"+++ b/internal/codereview/brief.go",
			"@@",
			"+func BuildReviewBrief(ctx context.Context) error {",
			"+\treturn buildReviewBrief(ctx)",
			" \tunchanged := notCounted",
			"-\told_helper_name()",
		}, "\n"),
	}}
	tokens := changedIdentifiers(diffs, []string{"internal/codereview/brief.go"}, 40)
	joined := " " + strings.Join(tokens, " ") + " "
	for _, want := range []string{" BuildReviewBrief ", " Build ", " Review ", " Brief ", " old_helper_name ", " helper ", " brief "} {
		if !strings.Contains(joined, want) {
			t.Fatalf("tokens = %v, missing %q", tokens, strings.TrimSpace(want))
		}
	}
	if strings.Contains(joined, " notCounted ") {
		t.Fatalf("tokens = %v, context lines must not contribute", tokens)
	}
	for _, unwanted := range []string{" func ", " return ", " error ", " the "} {
		if strings.Contains(joined, unwanted) {
			t.Fatalf("tokens = %v, keyword %q should be filtered", tokens, strings.TrimSpace(unwanted))
		}
	}
	if len(changedIdentifiers(diffs, nil, 3)) != 3 {
		t.Fatal("limit is not applied")
	}
}

func TestIndexQueryPayloadShapes(t *testing.T) {
	single, err := indexQueryPayload(indexQuery{BodyField: "content", BodyQuery: "hello", Limit: 5})
	if err != nil {
		t.Fatalf("indexQueryPayload() error = %v", err)
	}
	if _, multi := single["queries"]; multi {
		t.Fatalf("single leg produced a multi-query payload: %#v", single)
	}
	if limit, ok := single["limit"].(int); !ok || limit != 5 {
		t.Fatalf("single leg limit = %#v, want the scalar form", single["limit"])
	}

	fused, err := indexQueryPayload(indexQuery{
		BodyField: "content", BodyQuery: "hello",
		SymbolField: "symbol", SymbolQuery: "Thing",
		Vector: make([]float32, 4), Limit: 5,
	})
	if err != nil {
		t.Fatalf("indexQueryPayload() error = %v", err)
	}
	queries, ok := fused["queries"].([]any)
	if !ok || len(queries) != 3 {
		t.Fatalf("fused payload = %#v, want three ranked legs", fused)
	}
	rerank, ok := fused["rerank_by"].([]any)
	if !ok || len(rerank) != 1 || rerank[0] != "RRF" {
		t.Fatalf("rerank_by = %#v, want RRF", fused["rerank_by"])
	}
	for _, leg := range queries {
		limit, ok := leg.(map[string]any)["limit"].(map[string]any)
		if !ok || limit["total"] != 5 {
			t.Fatalf("fused leg limit = %#v, want the object form", leg)
		}
	}

	if _, err := indexQueryPayload(indexQuery{Limit: 5}); err == nil {
		t.Fatal("a query with no legs must be an error, not an empty request")
	}
}

func TestApplyIndexSchemaReadsWidthAndFullTextColumns(t *testing.T) {
	probe := indexProbe{}
	applyIndexSchema(&probe, map[string]indexSchemaField{
		"vector":     {Type: "[1536]f32"},
		"content":    {FullTextSearch: map[string]any{"stemming": false}},
		"symbol":     {FullTextSearch: map[string]any{"stemming": false}},
		"file_path":  {Type: "string", Filterable: true},
		"created_at": {Type: "int", Filterable: true},
	})
	if probe.Dimensions != 1536 {
		t.Fatalf("Dimensions = %d, want 1536", probe.Dimensions)
	}
	if probe.BodyField != "content" || probe.SymbolField != "symbol" {
		t.Fatalf("probe = %+v, want content/symbol full-text columns", probe)
	}
	if !probe.Has("file_path") || probe.Has("source_kind") {
		t.Fatalf("filterable set = %v", probe.Filterable)
	}
	// TurboPuffer rejects the entire query when include_attributes names an
	// attribute the namespace does not declare, so the wish list is narrowed to
	// the schema. This is not a refinement: naming `branch_name` against the
	// console code index returns HTTP 400 and no rows at all.
	include := probe.IncludeAttributes([]string{"content", "file_path", "branch_name", "symbol"})
	if strings.Join(include, ",") != "content,file_path,symbol" {
		t.Fatalf("IncludeAttributes = %v, want only the declared attributes", include)
	}
	// An unreadable schema means "ask for everything" rather than "ask for a
	// guess".
	if got := (indexProbe{}).IncludeAttributes([]string{"content"}); got != nil {
		t.Fatalf("IncludeAttributes with no known schema = %v, want nil", got)
	}
}

func TestIndexQueryPayloadAsksForAllAttributesWhenSchemaIsUnknown(t *testing.T) {
	payload, err := indexQueryPayload(indexQuery{BodyField: "content", BodyQuery: "hello", Limit: 3})
	if err != nil {
		t.Fatalf("indexQueryPayload() error = %v", err)
	}
	if payload["include_attributes"] != true {
		t.Fatalf("include_attributes = %#v, want true", payload["include_attributes"])
	}
}

func TestEmbedModelMatchesNamespaceWidth(t *testing.T) {
	t.Setenv("LGTM_REVIEW_INDEX_EMBED_MODEL", "")
	t.Setenv("LGTM_OPENAI_EMBEDDING_MODEL", "")
	t.Setenv("LGTM_EMBEDDING_DIMENSIONS", "")
	for width, want := range map[int]string{
		3072: "text-embedding-3-large",
		1536: "text-embedding-3-small",
		512:  "text-embedding-3-small",
	} {
		model, ok := embedModelForDimensions(width)
		if !ok || model != want {
			t.Fatalf("embedModelForDimensions(%d) = %q,%v want %q", width, model, ok, want)
		}
	}
	if _, ok := embedModelForDimensions(777); ok {
		t.Fatal("an unknown width must not be served by a mismatched model")
	}
}

func TestFuseRankedRowsPrefersRowsSeveralListsAgree_On(t *testing.T) {
	first := []indexRow{
		{"file_path": "a.go", "start_line": 1.0, "content": "a"},
		{"file_path": "b.go", "start_line": 1.0, "content": "b"},
	}
	second := []indexRow{
		{"file_path": "c.go", "start_line": 1.0, "content": "c"},
		{"file_path": "b.go", "start_line": 1.0, "content": "bb"},
	}
	fused := fuseRankedRows([][]indexRow{first, second}, codeIndexRowKey)
	if len(fused) != 3 {
		t.Fatalf("fused = %#v, want three distinct rows", fused)
	}
	if codeIndexRowFile(fused[0]) != "b.go" {
		t.Fatalf("fused[0] = %v, want the row both lists ranked", fused[0])
	}
	if codeIndexRowBody(fused[0]) != "bb" {
		t.Fatalf("fused[0] body = %q, want the longer copy", codeIndexRowBody(fused[0]))
	}
}
