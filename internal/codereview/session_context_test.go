package codereview

import (
	"context"
	"strings"
	"testing"
)

func sessionTestInput() RetrieveInput {
	return RetrieveInput{
		RepoRoot:     "/repo",
		Options:      Options{Scope: "security", Prompt: "did we ever decide to skip the retry?"},
		ChangedFiles: []string{"internal/capture/orchestrator/rawstage.go"},
		DiffSnippets: []DiffSnippet{{
			File: "internal/capture/orchestrator/rawstage.go",
			Diff: "@@\n+func stageRawEvent(ctx context.Context) error {\n",
		}},
		Evidence: &EvidenceLog{},
	}
}

func TestSessionContextRetrieverEmitsSessionSnippets(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"totality-org-repo-v2": {
				Exists: true, Dimensions: 512, BodyField: "text",
				Filterable: map[string]bool{"source_kind": true, "repo_full_name": true},
			},
		},
		rows: map[string][]indexRow{
			"totality-org-repo-v2": {{
				"source_kind":    "published_session_context",
				"source_id":      "sess-1#0",
				"session_id":     "sess-1",
				"request_id":     "req-9",
				"repo_full_name": "satoricorp/totality",
				"branch_name":    "tx-review",
				"text":           "The agent was told to keep retries idempotent and left the backoff for later.",
			}},
		},
	}
	in := sessionTestInput()
	retriever := SessionContextRetriever{
		Store:       store,
		Namespaces:  []codeIndexTarget{{Namespace: "totality-org-repo-v2", Origin: "org session namespace"}},
		Limit:       5,
		EmbedderFor: func(width int) (reviewResourceEmbedder, string, bool) { return staticEmbedder(width)(width) },
	}
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 1 {
		t.Fatalf("snippets = %#v, want one session snippet", snippets)
	}
	snippet := snippets[0]
	if snippet.Kind != "indexed_session" {
		t.Fatalf("kind = %q, want indexed_session", snippet.Kind)
	}
	if snippet.SessionID != "sess-1" || snippet.RequestID != "req-9" {
		t.Fatalf("snippet = %+v, want the session identity preserved", snippet)
	}
	if snippet.Ref != "sess-1/req-9" {
		t.Fatalf("ref = %q", snippet.Ref)
	}
	if !strings.Contains(snippet.Text, "idempotent") {
		t.Fatalf("text = %q, want the transcript body", snippet.Text)
	}
	// The snippet must survive the mapping into a citable source reference,
	// otherwise the model can read it but cannot attribute a finding to it.
	refs := sourceRefsFromContextSnippets(labelContextSnippets(snippets))
	if len(refs) != 1 || refs[0].Kind != "session" || refs[0].SessionID != "sess-1" {
		t.Fatalf("source refs = %#v, want one session ref", refs)
	}
}

func TestSessionContextRetrieverFiltersByKindAndRepo(t *testing.T) {
	store := &fakeIndexStore{
		probes: map[string]indexProbe{
			"scoped":   {Exists: true, Dimensions: 512, BodyField: "text", Filterable: map[string]bool{"source_kind": true, "repo_full_name": true}},
			"unscoped": {Exists: true, Dimensions: 512, BodyField: "text"},
		},
	}
	in := sessionTestInput()
	retriever := SessionContextRetriever{
		Store: store,
		Namespaces: []codeIndexTarget{
			{Namespace: "scoped", Origin: "org session namespace"},
			{Namespace: "unscoped", Origin: "legacy session namespace"},
		},
		EmbedderFor: func(width int) (reviewResourceEmbedder, string, bool) { return staticEmbedder(width)(width) },
	}
	if _, err := retriever.Retrieve(context.Background(), in); err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if store.query("scoped").Filters == nil {
		t.Fatal("a namespace declaring source_kind and repo_full_name must be filtered on both")
	}
	if store.query("unscoped").Filters != nil {
		t.Fatalf("filters = %v, want none for a namespace declaring neither attribute", store.query("unscoped").Filters)
	}
}

func TestSessionFiltersUseOnlyDeclaredAttributes(t *testing.T) {
	both := sessionFilters(indexProbe{Filterable: map[string]bool{"source_kind": true, "repo_full_name": true}}, "satoricorp/totality")
	conditions, ok := both.([]any)
	if !ok || len(conditions) != 2 || conditions[0] != "And" {
		t.Fatalf("filters = %#v, want a two-condition And", both)
	}
	onlyKind := sessionFilters(indexProbe{Filterable: map[string]bool{"source_kind": true}}, "satoricorp/totality")
	clause, ok := onlyKind.([]any)
	if !ok || clause[0] != "source_kind" {
		t.Fatalf("filters = %#v, want a bare source_kind clause", onlyKind)
	}
	if sessionFilters(indexProbe{}, "satoricorp/totality") != nil {
		t.Fatal("a namespace declaring neither attribute must not be filtered")
	}
}

func TestSessionContextRetrieverReportsMissingNamespace(t *testing.T) {
	in := sessionTestInput()
	retriever := SessionContextRetriever{
		Store:       &fakeIndexStore{},
		Namespaces:  []codeIndexTarget{{Namespace: "totality-sessions-dev", Origin: "legacy session namespace"}},
		EmbedderFor: func(int) (reviewResourceEmbedder, string, bool) { return nil, "", false },
	}
	snippets, err := retriever.Retrieve(context.Background(), in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 0 {
		t.Fatalf("snippets = %#v, want none", snippets)
	}
	warnings := EvidenceWarnings(in.Evidence.Statuses())
	if len(warnings) != 1 {
		t.Fatalf("warnings = %v, want one line for the source", warnings)
	}
	// The namespace is not in the warning. A reader is told what they lost and
	// what to do; an internal namespace id is neither. It stays in the
	// structured evidence, which is where anyone diagnosing a lookup looks.
	statuses := in.Evidence.Statuses()
	if len(statuses) != 1 || statuses[0].Namespace != "totality-sessions-dev" {
		t.Fatalf("statuses = %#v, want the searched namespace recorded", statuses)
	}
}
