package publication

import (
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
	}), prBodyCatalog{Files: []string{"internal/auth/session.go"}}, summaryContext, lexicalReach{})

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
