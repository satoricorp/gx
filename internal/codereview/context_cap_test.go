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
