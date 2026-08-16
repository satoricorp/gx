package codereview

import "testing"

// The measured shape this guards against: a code-index answer where five of
// the 48 rows were consecutive chunks of one planning document, spending a
// third of the context budget on a file that was not near the change.
func TestCapRetrievedChunksPerFileKeepsTwoPerFileInRankOrder(t *testing.T) {
	in := []ContextSnippet{
		{Kind: "indexed_code", File: "PLAN.md", Ref: "PLAN.md:171"},
		{Kind: "indexed_code", File: "PLAN.md", Ref: "PLAN.md:86"},
		{Kind: "indexed_code", File: "src/a.ts", Ref: "src/a.ts:1"},
		{Kind: "indexed_code", File: "PLAN.md", Ref: "PLAN.md:341"},
		{Kind: "indexed_code", File: "PLAN.md", Ref: "PLAN.md:1"},
		{Kind: "indexed_code", File: "src/b.ts", Ref: "src/b.ts:40"},
		{Kind: "indexed_code", File: "PLAN.md", Ref: "PLAN.md:256"},
	}
	got := capRetrievedChunksPerFile(in, 2)
	want := []string{"PLAN.md:171", "PLAN.md:86", "src/a.ts:1", "src/b.ts:40"}
	if len(got) != len(want) {
		t.Fatalf("kept %d snippets, want %d: %+v", len(got), len(want), refsOf(got))
	}
	for i, ref := range want {
		if got[i].Ref != ref {
			t.Errorf("position %d = %q, want %q (order must follow the retriever's ranking)", i, got[i].Ref, ref)
		}
	}
}

// Local documents were chosen by name, not by similarity, and a doc split into
// parts is meant to be read whole; snippets with no File cannot be grouped.
// Neither is capped.
func TestCapRetrievedChunksPerFileLeavesLocalAndFilelessSnippetsAlone(t *testing.T) {
	in := []ContextSnippet{
		{Kind: "local_doc", File: "README.md", Ref: "README.md:1"},
		{Kind: "local_doc", File: "README.md", Ref: "README.md:80"},
		{Kind: "local_doc", File: "README.md", Ref: "README.md:160"},
		{Kind: "indexed_session", Ref: "session-1"},
		{Kind: "indexed_session", Ref: "session-2"},
		{Kind: "indexed_session", Ref: "session-3"},
		{Kind: "code_review_history", File: "src/a.ts", Ref: "h1"},
		{Kind: "code_review_history", File: "src/a.ts", Ref: "h2"},
		{Kind: "code_review_history", File: "src/a.ts", Ref: "h3"},
	}
	got := capRetrievedChunksPerFile(in, 2)
	// 3 local + 3 file-less sessions pass through; review history for one
	// file is a retrieved kind and is capped to 2.
	if len(got) != 8 {
		t.Fatalf("kept %d snippets, want 8: %v", len(got), refsOf(got))
	}
	if got[7].Ref != "h2" {
		t.Errorf("last kept = %q, want h2 (third history row for the same file dropped)", got[7].Ref)
	}
}

func TestCapRetrievedChunksPerFileZeroDisables(t *testing.T) {
	in := []ContextSnippet{
		{Kind: "indexed_code", File: "x", Ref: "1"},
		{Kind: "indexed_code", File: "x", Ref: "2"},
		{Kind: "indexed_code", File: "x", Ref: "3"},
	}
	if got := capRetrievedChunksPerFile(in, 0); len(got) != 3 {
		t.Fatalf("limit 0 must disable the cap, kept %d", len(got))
	}
}

func refsOf(snippets []ContextSnippet) []string {
	out := make([]string, 0, len(snippets))
	for _, s := range snippets {
		out = append(out, s.Ref)
	}
	return out
}

// A change review ships whole "code quality" files only when the change
// touches them; the hints themselves still travel in static.code_quality. A
// whole-repo review, or one with no known change set, keeps the old behavior.
func TestCodeQualityFilesAreScopedToTheChange(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	flagged := "package app\nfunc Run() { result, _ := maybe(); _ = result }\nfunc maybe() (int, error) { return 0, nil }\n"
	writeFile(t, root, "internal/app/app.go", flagged)
	writeFile(t, root, "internal/other/other.go", flagged)
	facts := RepoFacts{
		DependencyFiles: []string{"go.mod"},
		Files:           []string{"go.mod", "internal/app/app.go", "internal/other/other.go"},
	}
	qualityRefs := func(in RetrieveInput) []string {
		snippets, err := LocalContextRetriever{}.Retrieve(t.Context(), in)
		if err != nil {
			t.Fatalf("Retrieve() error = %v", err)
		}
		var out []string
		for _, s := range snippets {
			if s.Kind == "code_quality_file" {
				out = append(out, s.Ref)
			}
		}
		return out
	}

	// Change review touching only app.go: other.go's whole file is not shipped.
	got := qualityRefs(RetrieveInput{RepoRoot: root, Options: Options{Scope: DefaultScope}, Facts: facts, ChangedFiles: []string{"internal/app/app.go"}})
	if len(got) != 1 || got[0] != "internal/app/app.go" {
		t.Fatalf("change review quality files = %v, want only the changed app.go", got)
	}
	// Whole-repo review: both.
	got = qualityRefs(RetrieveInput{RepoRoot: root, Options: Options{Scope: DefaultScope, WholeRepo: true}, Facts: facts, ChangedFiles: []string{"internal/app/app.go"}})
	if len(got) != 2 {
		t.Fatalf("whole-repo review quality files = %v, want both flagged files", got)
	}
	// No known change set: both (nothing to scope to).
	got = qualityRefs(RetrieveInput{RepoRoot: root, Options: Options{Scope: DefaultScope}, Facts: facts})
	if len(got) != 2 {
		t.Fatalf("unscoped review quality files = %v, want both flagged files", got)
	}
}
