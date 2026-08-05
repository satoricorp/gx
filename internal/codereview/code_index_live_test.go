package codereview

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/gxtest"
)

// requireLiveRetrievalCredentials skips unless real keys are present. These
// tests only read; they create nothing and delete nothing.
func requireLiveRetrievalCredentials(t *testing.T) {
	t.Helper()
	// Checking only for keys is not a guard: it arms this test on the machine
	// of everyone who has them exported, which is everyone working on gx.
	gxtest.RequireNoNetwork(t)
	if strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")) == "" {
		t.Skip("TURBOPUFFER_API_KEY is not set; skipping live retrieval test")
	}
	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) == "" {
		t.Skip("OPENAI_API_KEY is not set; skipping live retrieval test")
	}
}

func liveIndexStore(t *testing.T) turboPufferIndexStore {
	t.Helper()
	return newTurboPufferIndexStore(
		strings.TrimSpace(os.Getenv("TURBOPUFFER_API_KEY")),
		firstNonEmpty(os.Getenv("GX_TPUF_BASE_URL"), defaultReviewResourceBaseURL),
	)
}

// TestLiveProbeReadsRealNamespaceSchemas pins the two facts every query depends
// on: that a namespace that exists reports a width and full-text columns, and
// that one that does not is reported as absent rather than as a failure.
func TestLiveProbeReadsRealNamespaceSchemas(t *testing.T) {
	requireLiveRetrievalCredentials(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	store := liveIndexStore(t)

	probe := store.Probe(ctx, "repo-satoricorp-gx")
	if probe.Err != nil {
		t.Fatalf("probe error = %v", probe.Err)
	}
	if !probe.Exists || !probe.Usable() {
		t.Fatalf("probe = %+v, want a usable namespace", probe)
	}
	if probe.BodyField == "" || probe.SymbolField == "" || probe.Dimensions == 0 {
		t.Fatalf("probe = %+v, want a body column, a symbol column and a vector width", probe)
	}

	missing := store.Probe(ctx, "repo-satoricorp-definitely-not-a-repo")
	if missing.Err != nil {
		t.Fatalf("a missing namespace must not read as an error: %v", missing.Err)
	}
	if missing.Exists {
		t.Fatal("a namespace that does not exist reported as existing")
	}
}

// TestLiveCodeIndexRetrievalFindsThisRepositorysCode runs the real retriever
// against the real code index and asserts it comes back with this repository's
// own source, attributed to a file and a line.
func TestLiveCodeIndexRetrievalFindsThisRepositorysCode(t *testing.T) {
	requireLiveRetrievalCredentials(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	store := liveIndexStore(t)
	if !store.Probe(ctx, "repo-satoricorp-gx").Usable() {
		t.Skip("repo-satoricorp-gx is not available; skipping")
	}

	in := RetrieveInput{
		RepoRoot:     "/repo",
		Options:      Options{Scope: "security"},
		ChangedFiles: []string{"internal/codereview/context_endpoint.go"},
		DiffSnippets: []DiffSnippet{{
			File: "internal/codereview/context_endpoint.go",
			Diff: "@@\n+func contextRetrieverFromEnv() ContextRetriever {\n+\tdedupeContextSnippets(out)\n",
		}},
		Evidence: &EvidenceLog{},
	}
	retriever := CodeIndexRetriever{
		Store:      store,
		Namespaces: []codeIndexTarget{{Namespace: "repo-satoricorp-gx", Origin: "gx Cloud code index"}},
		Limit:      8,
	}
	snippets, err := retriever.Retrieve(ctx, in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) == 0 {
		t.Fatalf("no snippets; evidence = %#v", in.Evidence.Statuses())
	}
	var sawTarget bool
	for _, snippet := range snippets {
		if snippet.Kind != "indexed_code" {
			t.Fatalf("snippet kind = %q, want indexed_code", snippet.Kind)
		}
		if snippet.File == "" || snippet.StartLine <= 0 {
			t.Fatalf("snippet %+v has no anchor; a finding could not cite it", snippet)
		}
		if strings.Contains(snippet.File, "context_endpoint.go") {
			sawTarget = true
		}
	}
	if !sawTarget {
		var files []string
		for _, snippet := range snippets {
			files = append(files, snippet.File)
		}
		t.Fatalf("hybrid retrieval did not surface the changed file's own module; got %v", files)
	}
	statuses := in.Evidence.Statuses()
	if len(statuses) != 1 || statuses[0].State != EvidenceOK {
		t.Fatalf("evidence = %#v, want one healthy source", statuses)
	}
}

// TestLiveMissingNamespaceDegradesVisibly is the yeet case: a repository gx has
// never indexed must produce a stated absence, not a silently thinner review.
func TestLiveMissingNamespaceDegradesVisibly(t *testing.T) {
	requireLiveRetrievalCredentials(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	in := RetrieveInput{
		RepoRoot:     "/repo",
		Options:      Options{Scope: "security"},
		ChangedFiles: []string{"main.go"},
		DiffSnippets: []DiffSnippet{{File: "main.go", Diff: "@@\n+func runGuest() error {\n"}},
		Evidence:     &EvidenceLog{},
	}
	retriever := CodeIndexRetriever{
		Store:      liveIndexStore(t),
		Namespaces: []codeIndexTarget{{Namespace: "repo-satoricorp-yeet", Origin: "gx Cloud code index"}},
	}
	snippets, err := retriever.Retrieve(ctx, in)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(snippets) != 0 {
		t.Fatalf("snippets = %#v, want none", snippets)
	}
	warnings := EvidenceWarnings(in.Evidence.Statuses())
	if len(warnings) != 1 || !strings.Contains(warnings[0], "repo-satoricorp-yeet") {
		t.Fatalf("warnings = %v, want the un-indexed repository stated", warnings)
	}
}
