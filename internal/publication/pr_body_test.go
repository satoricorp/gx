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
			File:           "internal/github/client.go",
			Line:           2,
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

func TestBlobPermalink(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/42"
	got := blobPermalink(prURL, "deadbeef", "internal/foo.go", 10)
	want := "https://github.com/satoricorp/gx/blob/deadbeef/internal/foo.go#L10"
	if got != want {
		t.Fatalf("blobPermalink() = %q, want %q", got, want)
	}
	if got := blobPermalink(prURL, "deadbeef", "internal/foo.go", 0); got != "https://github.com/satoricorp/gx/blob/deadbeef/internal/foo.go" {
		t.Fatalf("file-only blobPermalink() = %q", got)
	}
	if got := blobPermalink("", "deadbeef", "internal/foo.go", 10); got != "" {
		t.Fatalf("empty prURL blobPermalink() = %q, want empty", got)
	}
}

func TestHunkLinkForFindingValidatedLine(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/1"
	hunks := []prHunkSummary{{
		File:     "internal/github/client.go",
		NewStart: 1,
		NewLines: 3,
		Link:     githubHunkLink(prURL, "internal/github/client.go", 1, 1, 3),
	}}
	finding := codereview.Finding{File: "internal/github/client.go", Line: 2}
	got := hunkLinkForFinding(prURL, hunks, finding)
	want := githubHunkLineLink(prURL, "internal/github/client.go", 2)
	if got != want {
		t.Fatalf("hunkLinkForFinding() = %q, want %q", got, want)
	}
}

func TestHunkLinkForFindingUnresolvableFile(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/1"
	hunks := []prHunkSummary{{
		File:     "internal/github/client.go",
		NewStart: 1,
		NewLines: 3,
		Link:     githubHunkLink(prURL, "internal/github/client.go", 1, 1, 3),
	}}
	finding := codereview.Finding{File: "internal/missing/file.go", Line: 2}
	if got := hunkLinkForFinding(prURL, hunks, finding); got != "" {
		t.Fatalf("hunkLinkForFinding() = %q, want empty link for unresolvable file", got)
	}
}

func TestFindingAttributionsMixedResolvedSources(t *testing.T) {
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Push: reviewbundle.PushPayload{HeadCommitID: "abc123"},
	})
	catalog := prBodyCatalog{Files: []string{"REVIEW.md"}}
	prURL := "https://github.com/satoricorp/gx/pull/9"
	artifact.Bundle.Push.GitHubPullRequestURL = &prURL
	attrs := findingAttributions(codereview.Finding{
		ResolvedSources: []codereview.ResolvedSource{
			{Kind: "resource", Publisher: "Google Engineering Practices", Opaque: true},
			{Kind: "local", Ref: "REVIEW.md", File: "REVIEW.md", Publisher: "this repo"},
		},
	}, prSummaryContext{}, artifact, catalog)
	if len(attrs) != 2 {
		t.Fatalf("attributions = %#v, want 2 entries", attrs)
	}
	if !attrs[0].Opaque || attrs[0].Label != "Google Engineering Practices" {
		t.Fatalf("opaque attribution = %#v", attrs[0])
	}
	if attrs[1].Opaque {
		t.Fatalf("local attribution should not be opaque: %#v", attrs[1])
	}
}

func TestFindingAttributionsNoEvidence(t *testing.T) {
	attrs := findingAttributions(codereview.Finding{
		Title:   "Review auth handling",
		Summary: "Check session expiry.",
	}, prSummaryContext{}, reviewbundle.NewArtifact(reviewbundle.Bundle{}), prBodyCatalog{})
	if len(attrs) != 0 {
		t.Fatalf("attributions = %#v, want none", attrs)
	}
	if got := renderAttributions(attrs); got != "" {
		t.Fatalf("renderAttributions() = %q, want empty", got)
	}
}

func TestFindingAttributionsOutOfDiffBlobPermalink(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/5"
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Push: reviewbundle.PushPayload{HeadCommitID: "sha999"},
	})
	artifact.Bundle.Push.GitHubPullRequestURL = &prURL
	catalog := prBodyCatalog{
		Revisions: []prRevisionSummary{{CommitID: "sha999"}},
		Files:     []string{"internal/publication/pr_body.go"},
		Hunks: []prHunkSummary{{
			File:     "internal/publication/pr_body.go",
			NewStart: 1,
			NewLines: 5,
		}},
	}
	attrs := findingAttributions(codereview.Finding{
		ResolvedSources: []codereview.ResolvedSource{{
			Kind:      "code",
			File:      "internal/storage/schema.go",
			StartLine: 42,
			Publisher: "indexed",
		}},
	}, prSummaryContext{}, artifact, catalog)
	if len(attrs) != 1 {
		t.Fatalf("attributions = %#v", attrs)
	}
	want := blobPermalink(prURL, "sha999", "internal/storage/schema.go", 42)
	if attrs[0].URL != want {
		t.Fatalf("URL = %q, want blob %q", attrs[0].URL, want)
	}
	rendered := renderAttributions(attrs)
	if !strings.Contains(rendered, want) {
		t.Fatalf("renderAttributions() = %q, want blob link", rendered)
	}
}

func TestFindingAttributionsEmptyPRURLPlainText(t *testing.T) {
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Push: reviewbundle.PushPayload{HeadCommitID: "sha999"},
	})
	catalog := prBodyCatalog{
		Files: []string{"internal/publication/pr_body.go"},
		Hunks: []prHunkSummary{{
			File:     "internal/publication/pr_body.go",
			NewStart: 1,
			NewLines: 5,
		}},
	}
	attrs := findingAttributions(codereview.Finding{
		ResolvedSources: []codereview.ResolvedSource{{
			Kind:      "code",
			File:      "internal/storage/schema.go",
			StartLine: 42,
			Publisher: "indexed",
		}},
	}, prSummaryContext{}, artifact, catalog)
	if len(attrs) != 1 || attrs[0].URL != "" {
		t.Fatalf("attributions = %#v, want plain text without URL", attrs)
	}
	rendered := renderAttributions(attrs)
	if strings.Contains(rendered, "](http") {
		t.Fatalf("renderAttributions() = %q, want no markdown links", rendered)
	}
	if !strings.Contains(rendered, "internal/storage/schema.go:42") {
		t.Fatalf("renderAttributions() = %q, want plain file:line label", rendered)
	}
}

func TestSpeculativeFindingsDroppedFromPRBody(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return fakePRSummaryReviewer{findings: []codereview.Finding{{
			Title:    "Maybe consider renaming",
			Summary:  "Speculative nit.",
			Strength: "Speculative",
			File:     "internal/github/client.go",
		}}}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	prURL := "https://github.com/satoricorp/gx/pull/14"
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
		Stack: []reviewbundle.StackPayload{{
			Patch: "diff --git a/internal/github/client.go b/internal/github/client.go\n--- a/internal/github/client.go\n+++ b/internal/github/client.go\n@@ -1 +1,2 @@\n package github\n+// change\n",
			Change: reviewbundle.ChangePayload{
				Files: []string{"internal/github/client.go"},
			},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if strings.Contains(body, "Maybe consider renaming") {
		t.Fatalf("speculative finding should be dropped:\n%s", body)
	}
}
