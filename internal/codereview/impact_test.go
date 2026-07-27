package codereview

import (
	"strings"
	"testing"
)

// The ranking the review command did not have. Sorting changed files by path
// and reading the first N means an auth change is reviewed if and only if its
// path sorts early, which is not a property anyone would choose.
func TestChangedFilesAreRankedByImpactNotAlphabet(t *testing.T) {
	files := []string{
		"zzz/last.go",
		"internal/auth/token.go",
		"aaa/first.go",
		"internal/db/migration.go",
		"internal/auth/token_test.go",
	}
	ranked := rankFilesByImpact(files, ReviewPolicy{})
	if ranked[0] != "internal/db/migration.go" {
		t.Fatalf("first ranked file = %q; the schema change should lead", ranked[0])
	}
	if ranked[1] != "internal/auth/token.go" {
		t.Fatalf("second ranked file = %q, want the auth change", ranked[1])
	}
	if ranked[len(ranked)-1] == "internal/auth/token.go" {
		t.Fatalf("ranking put the auth change last: %#v", ranked)
	}
	// A test file scores on its path like any other, and still sorts behind its
	// non-test peer at the same score.
	if ranked[2] != "internal/auth/token_test.go" {
		t.Fatalf("ranked = %#v; the auth test should follow the auth source", ranked)
	}
	// Below the scored files, the alphabetical order survives as the tiebreak
	// that makes runs reproducible.
	if ranked[3] != "aaa/first.go" || ranked[4] != "zzz/last.go" {
		t.Fatalf("unscored files lost the stable alphabetical tiebreak: %#v", ranked)
	}
}

func TestConfiguredRiskPathsOutrankGenericPatterns(t *testing.T) {
	policy := ReviewPolicy{RiskPaths: []RiskPath{{Glob: "internal/billing/**", Message: "billing changes move money"}}}
	ranked := rankFilesByImpact([]string{"internal/auth/token.go", "internal/billing/charge.go"}, policy)
	if ranked[0] != "internal/billing/charge.go" {
		t.Fatalf("ranked = %#v; the project's own configured risk path must win", ranked)
	}
	impact := FileImpact("internal/billing/charge.go", policy)
	if impact.Score != 700 || !strings.Contains(impact.Detail, "billing changes move money") {
		t.Fatalf("impact = %+v", impact)
	}
	// And the risk-path title still reads the way the PR summary renders it.
	if impact.Title != "Verify change in charge.go" {
		t.Fatalf("risk-path title = %q", impact.Title)
	}
}

// The extraction has to be a move, not an edit: internal/publication ranks the
// hunks a PR summary describes with this, and every score and sentence it
// returns is load-bearing for that output.
func TestFileImpactPreservesThePRSummaryScores(t *testing.T) {
	for _, tc := range []struct {
		file  string
		score int
		title string
	}{
		{file: "internal/auth/token.go", score: 650, title: "Verify auth handling in token.go"},
		{file: "server/migrations/001.sql", score: 670, title: "Verify schema change in 001.sql"},
		{file: ".github/workflows/ci.yml", score: 660, title: "Verify deployment change in ci.yml"},
		{file: "internal/queue/queue.go", score: 0, title: ""},
	} {
		impact := FileImpact(tc.file, ReviewPolicy{})
		if impact.Score != tc.score || impact.Title != tc.title {
			t.Fatalf("FileImpact(%q) = %+v, want score %d title %q", tc.file, impact, tc.score, tc.title)
		}
	}
}

// A generated lockfile must never outrank the code under review. It used to:
// dependency_manifest sat at 6 and repo_source_file/module_file at 8, so when
// the context budget bound, go.sum survived and the source did not.
func TestLockfilesNeverEvictSourceCode(t *testing.T) {
	source := ContextSnippet{Kind: "repo_source_file", Ref: "internal/queue/queue.go", Text: "code"}
	module := ContextSnippet{Kind: "module_file", Ref: "internal/app/app.go", Text: "code"}
	changed := ContextSnippet{Kind: "changed_file", Ref: "internal/app/edit.go", Text: "code"}
	priorFinding := ContextSnippet{Kind: "code_review_history", Ref: "review#1", Text: "previously found"}
	lockfile := ContextSnippet{Kind: "dependency_lockfile", Ref: "go.sum", Text: "hashes"}
	manifest := ContextSnippet{Kind: "dependency_manifest", Ref: "go.mod", Text: "module"}

	for _, snippet := range []ContextSnippet{source, module, changed, priorFinding} {
		if contextSnippetPriority(snippet) >= contextSnippetPriority(lockfile) {
			t.Fatalf("%s (%d) does not outrank a lockfile (%d)", snippet.Kind, contextSnippetPriority(snippet), contextSnippetPriority(lockfile))
		}
		if contextSnippetPriority(snippet) >= contextSnippetPriority(manifest) {
			t.Fatalf("%s (%d) does not outrank a dependency manifest (%d)", snippet.Kind, contextSnippetPriority(snippet), contextSnippetPriority(manifest))
		}
	}
	if contextSnippetPriority(manifest) >= contextSnippetPriority(lockfile) {
		t.Fatalf("a declared manifest does not outrank its generated lockfile")
	}

	// And end to end through the budget: one slot, source and a lockfile.
	got := compactContextSnippets([]ContextSnippet{lockfile, source}, 1)
	if len(got) != 1 || got[0].Kind != "repo_source_file" {
		t.Fatalf("compacted = %#v; the lockfile evicted the source", got)
	}
}

// go.sum and go.mod are both "dependency files" and are not the same evidence.
// Splitting the kind is what lets the budget keep the one a reviewer can act on.
func TestLockfilesGetTheirOwnSnippetKind(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "go.sum", "example.com/dep v1.0.0 h1:abc=\n")

	snippets, err := LocalContextRetriever{}.Retrieve(t.Context(), RetrieveInput{
		RepoRoot: root,
		Options:  Options{Scope: DefaultScope},
		Facts:    RepoFacts{DependencyFiles: []string{"go.mod", "go.sum"}},
	})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	kinds := map[string]string{}
	for _, snippet := range snippets {
		kinds[snippet.Ref] = snippet.Kind
	}
	if kinds["go.mod"] != "dependency_manifest" {
		t.Fatalf("go.mod kind = %q", kinds["go.mod"])
	}
	if kinds["go.sum"] != "dependency_lockfile" {
		t.Fatalf("go.sum kind = %q; it shares a budget slot with the manifest", kinds["go.sum"])
	}
}

// The two budgets were 3200 in brief.go and 1200 in ai.go: two thirds of every
// file this package read was assembled and then discarded one layer down.
func TestTheSnippetBudgetIsOneNumber(t *testing.T) {
	if maxContextSnippetBytes != maxAIContextSnippetBytes {
		t.Fatalf("context snippet budgets disagree: brief %d vs ai %d", maxContextSnippetBytes, maxAIContextSnippetBytes)
	}
	if maxDiffSnippetBytes != maxAIDiffSnippetBytes {
		t.Fatalf("diff snippet budgets disagree: brief %d vs ai %d", maxDiffSnippetBytes, maxAIDiffSnippetBytes)
	}
}
