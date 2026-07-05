package publication

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestGitHubPullRequestBodyAllowsZeroNeedsReviewTargets(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) { return nil, codereview.ReviewerInfo{} }
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	prURL := "https://github.com/satoricorp/gx/pull/12"
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         "abc123",
			GitHubPullRequestURL: &prURL,
		},
		Stack: []reviewbundle.StackPayload{{
			BranchName: "docs/demo",
			Patch: strings.Join([]string{
				"diff --git a/docs/demo.md b/docs/demo.md",
				"--- a/docs/demo.md",
				"+++ b/docs/demo.md",
				"@@ -1,2 +1,3 @@",
				" # Demo",
				"+Small docs clarification.",
				"",
			}, "\n"),
			Change: reviewbundle.ChangePayload{
				CurrentCommitID: "abc123",
				Description:     "clarify docs",
				Files:           []string{"docs/demo.md"},
			},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "No specific high-impact review targets surfaced") {
		t.Fatalf("body did not render zero-target state:\n%s", body)
	}
	if strings.Contains(body, "Attribution:") {
		t.Fatalf("zero-target body unexpectedly rendered attribution:\n%s", body)
	}
}

func TestGitHubPullRequestBodyCapsNeedsReviewTargets(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	findings := make([]codereview.Finding, 0, 12)
	for i := 0; i < 12; i++ {
		findings = append(findings, codereview.Finding{
			Title:          fmt.Sprintf("Verify GitHub publish behavior %02d", i+1),
			Summary:        "internal/github/client.go changes GitHub publish behavior and could create a remote regression if wrong.",
			Recommendation: "Review the changed hunk before merging.",
			Evidence:       []codereview.Evidence{{Label: "Changed hunk", Value: "internal/github/client.go"}},
			Strength:       "Strong",
		})
	}
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return fakePRSummaryReviewer{findings: findings}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	prURL := "https://github.com/satoricorp/gx/pull/13"
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         "abc123",
			GitHubPullRequestURL: &prURL,
		},
		Stack: []reviewbundle.StackPayload{{
			BranchName: "feature/github-client",
			Patch: strings.Join([]string{
				"diff --git a/internal/github/client.go b/internal/github/client.go",
				"--- a/internal/github/client.go",
				"+++ b/internal/github/client.go",
				"@@ -1,2 +1,3 @@",
				" package github",
				"+func UpdatePullRequest() {}",
				"",
			}, "\n"),
			Change: reviewbundle.ChangePayload{
				CurrentCommitID: "abc123",
				Description:     "update github client",
				Files:           []string{"internal/github/client.go"},
			},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if count := strings.Count(body, "\n- "); count != maxPRReviewItems {
		t.Fatalf("Needs Review item count = %d, want %d:\n%s", count, maxPRReviewItems, body)
	}
	if strings.Contains(body, "Verify GitHub publish behavior 11") {
		t.Fatalf("body included item beyond cap:\n%s", body)
	}
}
