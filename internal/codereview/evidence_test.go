package codereview

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type namedFailingRetriever struct{ name string }

func (r namedFailingRetriever) EvidenceSource() string { return r.name }

func (namedFailingRetriever) Retrieve(context.Context, RetrieveInput) ([]ContextSnippet, error) {
	return nil, errors.New("dial tcp: connection refused")
}

type namedEmptyRetriever struct{ name string }

func (r namedEmptyRetriever) EvidenceSource() string { return r.name }

func (namedEmptyRetriever) Retrieve(context.Context, RetrieveInput) ([]ContextSnippet, error) {
	return nil, nil
}

func TestCompositeContextRetrieverRecordsFailuresInsteadOfSwallowingThem(t *testing.T) {
	log := &EvidenceLog{}
	snippets, err := (CompositeContextRetriever{Retrievers: []ContextRetriever{
		LocalContextRetriever{},
		namedFailingRetriever{name: "code index"},
		namedEmptyRetriever{name: "review knowledge"},
	}}).Retrieve(context.Background(), RetrieveInput{RepoRoot: t.TempDir(), Options: Options{}, Evidence: log})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 0 {
		t.Fatalf("snippets = %#v", snippets)
	}

	statuses := log.Statuses()
	if len(statuses) != 2 {
		t.Fatalf("statuses = %#v, want one per named source", statuses)
	}
	byName := map[string]EvidenceStatus{}
	for _, status := range statuses {
		byName[status.Source] = status
	}
	if byName["code index"].State != EvidenceUnavailable {
		t.Fatalf("code index status = %+v, want unavailable", byName["code index"])
	}
	if !strings.Contains(byName["code index"].Detail, "connection refused") {
		t.Fatalf("detail = %q, want the underlying error", byName["code index"].Detail)
	}
	if byName["review knowledge"].State != EvidenceEmpty {
		t.Fatalf("review knowledge status = %+v, want empty (it answered)", byName["review knowledge"])
	}

	warnings := EvidenceWarnings(statuses)
	if len(warnings) != 1 || !strings.HasPrefix(warnings[0], "code index:") {
		t.Fatalf("warnings = %v, want only the failed source", warnings)
	}
}

func TestCompositeContextRetrieverKeepsRetrieverReportedStatuses(t *testing.T) {
	log := &EvidenceLog{}
	retriever := CodeIndexRetriever{
		Store:       &fakeIndexStore{},
		Namespaces:  []codeIndexTarget{{Namespace: "repo-owner-repo", Origin: "GX Cloud code index"}},
		EmbedderFor: func(int) (reviewResourceEmbedder, string, bool) { return nil, "", false },
	}
	if _, err := (CompositeContextRetriever{Retrievers: []ContextRetriever{retriever}}).Retrieve(
		context.Background(),
		RetrieveInput{RepoRoot: "/repo", Options: Options{}, Evidence: log},
	); err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	statuses := log.Statuses()
	if len(statuses) != 1 {
		t.Fatalf("statuses = %#v, want the retriever's own per-namespace status only", statuses)
	}
	if statuses[0].Namespace != "repo-owner-repo" || statuses[0].State != EvidenceMissing {
		t.Fatalf("status = %+v, want the namespace-level missing report kept", statuses[0])
	}
}

func TestEvidenceLogToleratesNilReceiver(t *testing.T) {
	var log *EvidenceLog
	log.Record(EvidenceStatus{Source: "code index", State: EvidenceOK})
	if statuses := log.Statuses(); statuses != nil {
		t.Fatalf("Statuses() = %#v, want nil", statuses)
	}
	if log.count() != 0 {
		t.Fatal("count() on a nil log must be zero")
	}
}

func TestEvidenceStatusesAreOrderedDeterministically(t *testing.T) {
	log := &EvidenceLog{}
	log.Record(EvidenceStatus{Source: "review knowledge", State: EvidenceOK})
	log.Record(EvidenceStatus{Source: "code index", Namespace: "repo-b", State: EvidenceOK})
	log.Record(EvidenceStatus{Source: "code index", Namespace: "repo-a", State: EvidenceMissing})
	statuses := log.Statuses()
	got := []string{
		statuses[0].Source + "/" + statuses[0].Namespace,
		statuses[1].Source + "/" + statuses[1].Namespace,
		statuses[2].Source + "/" + statuses[2].Namespace,
	}
	want := []string{"code index/repo-a", "code index/repo-b", "review knowledge/"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v, want %v", got, want)
		}
	}
}

func TestContextBudgetDoesNotLetOneSourceStarveTheOthers(t *testing.T) {
	brief := ReviewBrief{Depth: "shallow"}
	for i := 0; i < maxAIContextSnippets*2; i++ {
		brief.Context = append(brief.Context, ContextSnippet{
			Kind: "indexed_code", Ref: "code" + string(rune('a'+i%26)) + string(rune('a'+i/26)), Text: "code",
		})
	}
	for i := 0; i < 6; i++ {
		brief.Context = append(brief.Context, ContextSnippet{
			Kind: "indexed_session", Ref: "sess" + string(rune('a'+i)), Text: "session",
		})
	}
	for i := 0; i < 6; i++ {
		brief.Context = append(brief.Context, ContextSnippet{
			Kind: "code_review_history", Ref: "prior" + string(rune('a'+i)), Text: "prior finding",
		})
	}

	kinds := map[string]int{}
	for _, snippet := range compactReviewBriefForAI(brief).Context {
		kinds[snippet.Kind]++
	}
	for _, kind := range []string{"indexed_code", "indexed_session", "code_review_history"} {
		if kinds[kind] == 0 {
			t.Fatalf("kind %q was cut entirely by the context budget: %v", kind, kinds)
		}
	}
	if kinds["indexed_code"] < 6 {
		t.Fatalf("indexed code should still dominate the budget: %v", kinds)
	}
}

func TestRenderMarkdownStatesUnavailableEvidence(t *testing.T) {
	report := Report{
		Reviewed:   true,
		ReviewMode: ReviewModeWorkingTree,
		Findings:   nil,
		Evidence: []EvidenceStatus{
			{Source: "code index", Namespace: "repo-satoricorp-yeet", State: EvidenceMissing, Detail: "GX Cloud has never indexed this repository"},
			{Source: "review knowledge", State: EvidenceOK, Snippets: 8},
		},
	}
	markdown := RenderMarkdown(report)
	if !strings.Contains(markdown, "review evidence unavailable") {
		t.Fatalf("markdown does not warn about missing evidence:\n%s", markdown)
	}
	if !strings.Contains(markdown, "repo-satoricorp-yeet") {
		t.Fatalf("markdown does not name the missing namespace:\n%s", markdown)
	}
	if strings.Contains(markdown, "review knowledge") {
		t.Fatalf("a healthy source must not appear in the warning:\n%s", markdown)
	}

	// A review with every source healthy reads clean, so the warning stays
	// meaningful.
	healthy := report
	healthy.Evidence = []EvidenceStatus{{Source: "code index", Namespace: "repo-satoricorp-gx", State: EvidenceOK, Snippets: 12}}
	if strings.Contains(RenderMarkdown(healthy), "review evidence unavailable") {
		t.Fatal("a fully informed review must not carry the degraded-evidence warning")
	}
}

func TestRenderMarkdownListsEveryEvidenceSourceWhenVerbose(t *testing.T) {
	report := Report{
		Reviewed:   true,
		ReviewMode: ReviewModeWorkingTree,
		Verbose:    true,
		Evidence: []EvidenceStatus{
			{Source: "code index", Namespace: "repo-satoricorp-gx", State: EvidenceOK, Snippets: 12},
			{Source: "indexed sessions", State: EvidenceDisabled, Detail: "TURBOPUFFER_API_KEY is not set"},
		},
	}
	markdown := RenderMarkdown(report)
	for _, want := range []string{"Evidence: code index (repo-satoricorp-gx): ok, 12 snippet(s)", "Evidence: indexed sessions: disabled"} {
		if !strings.Contains(markdown, want) {
			t.Fatalf("markdown missing %q:\n%s", want, markdown)
		}
	}
}
