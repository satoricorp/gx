package publication

import (
	"fmt"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestPRReviewBriefLabelsContextAndSourceRefs(t *testing.T) {
	summaryContext := prSummaryContext{
		Snippets: []codereview.ContextSnippet{
			{Kind: "repo_doc", Ref: "REVIEW.md", Source: "local", Text: "Review auth carefully."},
			{Kind: "review_resource", Ref: "owasp#1", Source: "turbopuffer:gx-review-knowledge", Text: "Validate auth inputs."},
		},
	}
	brief := prReviewBrief(reviewbundle.NewArtifact(reviewbundle.Bundle{
		Repo: reviewbundle.RepoPayload{RootPath: "/repo"},
	}), prBodyCatalog{Files: []string{"internal/auth/session.go"}}, summaryContext, lexicalReach{}, codereview.ReviewPolicy{})

	if brief.ReviewProfile != "pr_summary" {
		t.Fatalf("ReviewProfile = %q, want pr_summary", brief.ReviewProfile)
	}
	if len(brief.Context) != 2 || brief.Context[0].SourceLabel != "L1" || brief.Context[1].SourceLabel != "R1" {
		t.Fatalf("Context = %#v", brief.Context)
	}
	if len(brief.SourceRefs) != 2 {
		t.Fatalf("SourceRefs = %#v", brief.SourceRefs)
	}
	if len(brief.SourceCatalog) == 0 {
		t.Fatal("SourceCatalog is empty")
	}
}

func TestPRSummaryDiffSnippetsKeepHighImpactHunksOverTheCap(t *testing.T) {
	// Fill the brief past its cap with low-impact hunks, then bury a
	// risk-path hunk at the very end. Ranking must still surface it.
	var hunks []prHunkSummary
	for i := 0; i < maxPRSummaryDiffSnippets+10; i++ {
		hunks = append(hunks, prHunkSummary{
			File:     fmt.Sprintf("internal/widget/widget%02d.go", i),
			Patch:    fmt.Sprintf("@@ -1 +1,2 @@\n+// widget %d\n", i),
			NewLines: 2,
		})
	}
	buried := prHunkSummary{
		File:     "internal/auth/session.go",
		Patch:    "@@ -1 +1,2 @@\n+// rotate credentials\n",
		NewLines: 2,
	}
	hunks = append(hunks, buried)

	catalog := prBodyCatalog{Hunks: hunks}
	policy := codereview.ReviewPolicy{RiskPaths: []codereview.RiskPath{{
		Glob:    "internal/auth/**",
		Message: "Auth changes need a careful look",
	}}}

	snippets := prSummaryDiffSnippets(catalog, policy)
	if len(snippets) != maxPRSummaryDiffSnippets {
		t.Fatalf("len(snippets) = %d, want %d", len(snippets), maxPRSummaryDiffSnippets)
	}
	if snippets[0].File != buried.File {
		t.Fatalf("snippets[0].File = %q, want the risk-path hunk %q ranked first", snippets[0].File, buried.File)
	}
}

func TestPRSummaryDiffSnippetsRankCodeAheadOfTests(t *testing.T) {
	catalog := prBodyCatalog{Hunks: []prHunkSummary{
		{File: "internal/app/app_test.go", Patch: "@@ -1 +1,9 @@\n+// big test\n", NewLines: 9},
		{File: "internal/app/app.go", Patch: "@@ -1 +1,2 @@\n+// code\n", NewLines: 2},
	}}
	snippets := prSummaryDiffSnippets(catalog, codereview.ReviewPolicy{})
	if len(snippets) != 2 || snippets[0].File != "internal/app/app.go" {
		t.Fatalf("snippets = %#v, want the non-test hunk first", snippets)
	}
}
