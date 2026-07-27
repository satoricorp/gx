package publication

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
	"github.com/satoricorp/gx/internal/reviewsource"
)

func TestGitHubPullRequestBodyOmitsNotableChangesWhenEmpty(t *testing.T) {
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
		Revisions: []reviewbundle.RevisionPayload{{
			BranchName:           "docs/demo",
			CommitID:             "abc123",
			Description:          "clarify docs",
			Files:                []string{"docs/demo.md"},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "> 👀 **Quick scan** — standard change; AI review unavailable.") {
		t.Fatalf("body missing quick scan verdict banner for unavailable AI:\n%s", body)
	}
	if strings.Contains(body, "## Notable Changes") {
		t.Fatalf("body should omit Notable Changes when nothing is notable:\n%s", body)
	}
	if strings.Contains(body, "No specific high-impact review targets surfaced") {
		t.Fatalf("body should not render legacy filler line:\n%s", body)
	}
	if strings.Contains(body, "Attribution:") {
		t.Fatalf("zero-target body unexpectedly rendered attribution:\n%s", body)
	}
}

func TestGitHubPullRequestBodyCapsNotableChanges(t *testing.T) {
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
			Line:           i + 2,
		})
	}
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{findings: findings}, codereview.ReviewerInfo{}
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
		Revisions: []reviewbundle.RevisionPayload{{
			BranchName: "feature/github-client",
			Patch: strings.Join([]string{
				"diff --git a/internal/github/client.go b/internal/github/client.go",
				"--- a/internal/github/client.go",
				"+++ b/internal/github/client.go",
				"@@ -1,2 +1,14 @@",
				" package github",
				"+func UpdatePullRequest01() {}",
				"+func UpdatePullRequest02() {}",
				"+func UpdatePullRequest03() {}",
				"+func UpdatePullRequest04() {}",
				"+func UpdatePullRequest05() {}",
				"+func UpdatePullRequest06() {}",
				"+func UpdatePullRequest07() {}",
				"+func UpdatePullRequest08() {}",
				"+func UpdatePullRequest09() {}",
				"+func UpdatePullRequest10() {}",
				"+func UpdatePullRequest11() {}",
				"+func UpdatePullRequest12() {}",
				"",
			}, "\n"),
			CommitID:             "abc123",
			Description:          "update github client",
			Files:                []string{"internal/github/client.go"},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if count := strings.Count(body, "\n- "); count != maxPRNotableChanges {
		t.Fatalf("Notable Changes item count = %d, want %d:\n%s", count, maxPRNotableChanges, body)
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

func TestGitHubPullRequestFilesURL(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://github.com/satoricorp/gx/pull/71", "https://github.com/satoricorp/gx/pull/71/files"},
		{"https://github.com/satoricorp/gx/pull/71/", "https://github.com/satoricorp/gx/pull/71/files"},
		{"https://github.com/satoricorp/gx/pull/71/changes", "https://github.com/satoricorp/gx/pull/71/files"},
		{"https://github.com/satoricorp/gx/pull/71/files", "https://github.com/satoricorp/gx/pull/71/files"},
	}
	for _, tc := range cases {
		if got := githubPullRequestFilesURL(tc.in); got != tc.want {
			t.Fatalf("githubPullRequestFilesURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestGitHubHunkLineLinkUsesFilesTabAnchor(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/71/changes"
	file := "internal/reviewbundle/bundle.go"
	line := 910
	got := githubHunkLineLink(prURL, file, line)
	want := "https://github.com/satoricorp/gx/pull/71/files#diff-222c707db70208988323245aee1da2520b10ae584ecdf8136e08c1fd6b1cae65R910"
	if got != want {
		t.Fatalf("githubHunkLineLink() = %q, want %q", got, want)
	}
}

func TestHunkLinkForFindingValidatedLine(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/1"
	hunks := []prHunkSummary{{
		File:        "internal/github/client.go",
		NewStart:    1,
		NewLines:    3,
		ChangedLine: 1,
		Link:        githubHunkLink(prURL, "internal/github/client.go", 1, 1, 3),
	}}
	finding := codereview.Finding{File: "internal/github/client.go", Line: 99}
	got := hunkLinkForFinding(prURL, hunks, finding)
	if got != hunks[0].Link {
		t.Fatalf("hunkLinkForFinding() = %q, want hunk link %q", got, hunks[0].Link)
	}
	if strings.Contains(got, "#diff-") && !strings.Contains(got, "R") && !strings.Contains(got, "L") {
		t.Fatalf("hunkLinkForFinding() = %q, want line-anchored link", got)
	}
}

func TestPRBodyHasNoBareFileDiffLinks(t *testing.T) {
	bareFileLinkRE := regexp.MustCompile(`#diff-[0-9a-f]{64}\)`)
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prURL := "https://github.com/satoricorp/gx/pull/20"
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{findings: []codereview.Finding{{
			Title:    "Verify auth handling",
			Summary:  "Session token validation changed.",
			Strength: "Strong",
			File:     "internal/auth/session.go",
			Line:     99,
		}}}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
		Revisions: []reviewbundle.RevisionPayload{{
			Patch:                "diff --git a/internal/auth/session.go b/internal/auth/session.go\n--- a/internal/auth/session.go\n+++ b/internal/auth/session.go\n@@ -1 +1,2 @@\n package auth\n+func Validate() {}\n",
			Description:          "add session validation",
			Files:                []string{"internal/auth/session.go"},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if bareFileLinkRE.MatchString(body) {
		t.Fatalf("body contains bare file-level diff link:\n%s", body)
	}
}

func TestNotableChangesRendersPlainTextWithoutPRURL(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Revisions: []reviewbundle.RevisionPayload{{
			Patch: strings.Join([]string{
				"diff --git a/internal/app/handler.go b/internal/app/handler.go",
				"--- a/internal/app/handler.go",
				"+++ b/internal/app/handler.go",
				"@@ -1 +1,2 @@",
				" package app",
				"+func Handle() {}",
			}, "\n"),
			Description: "update handler",
			Files:       []string{"internal/app/handler.go"},
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "## Notable Changes") {
		t.Fatalf("body missing Notable Changes section:\n%s", body)
	}
	if strings.Contains(body, `href="http`) {
		t.Fatalf("body should not contain links without PR URL:\n%s", body)
	}
}

func TestNotableChangesIncludesTopHunkFill(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prURL := "https://github.com/satoricorp/gx/pull/21"
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
		Revisions: []reviewbundle.RevisionPayload{{
			Patch: strings.Join([]string{
				"diff --git a/internal/app/handler.go b/internal/app/handler.go",
				"--- a/internal/app/handler.go",
				"+++ b/internal/app/handler.go",
				"@@ -1 +1,2 @@",
				" package app",
				"+func Handle() {}",
			}, "\n"),
			Description:          "update handler",
			Files:                []string{"internal/app/handler.go"},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, `target="_blank"`) || !strings.Contains(body, "update handler: handler.go") {
		t.Fatalf("body missing top-hunk fill item:\n%s", body)
	}
}

func TestNotableChangesDropsUnanchoredAIEntry(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prURL := "https://github.com/satoricorp/gx/pull/22"
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{
			notableChanges: []codereview.NotableChange{
				{File: "internal/missing.go", Line: 1, Note: "This should not render"},
				{File: "internal/app/handler.go", Line: 2, Note: "Adds request handling"},
			},
		}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
		Revisions: []reviewbundle.RevisionPayload{{
			Patch: strings.Join([]string{
				"diff --git a/internal/app/handler.go b/internal/app/handler.go",
				"--- a/internal/app/handler.go",
				"+++ b/internal/app/handler.go",
				"@@ -1 +1,2 @@",
				" package app",
				"+func Handle() {}",
			}, "\n"),
			Description:          "update handler",
			Files:                []string{"internal/app/handler.go"},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if strings.Contains(body, "This should not render") {
		t.Fatalf("unanchored notable change should be dropped:\n%s", body)
	}
	if !strings.Contains(body, "Adds request handling") {
		t.Fatalf("anchored notable change missing:\n%s", body)
	}
}

func TestDedupeNotableChangesByTitle(t *testing.T) {
	items := dedupeNotableChanges([]prNotableChange{
		{
			Title: "Verify Publication code coordinates GX Cloud",
			Link:  "https://github.com/example/gx/pull/1/files#diff-aaaR10",
			Score: 700,
		},
		{
			Title: "Verify Publication code coordinates GX Cloud",
			Link:  "https://github.com/example/gx/pull/1/files#diff-bbbR20",
			Score: 750,
		},
		{
			Title: "Stop linkifying whole free-form detail sentences",
			Link:  "https://github.com/example/gx/pull/1/files#diff-cccR30",
			Score: 1000,
		},
	})
	if len(items) != 2 {
		t.Fatalf("dedupeNotableChanges() len = %d, want 2", len(items))
	}
	var publicationLink string
	for _, item := range items {
		if item.Title == "Verify Publication code coordinates GX Cloud" {
			publicationLink = item.Link
		}
	}
	if publicationLink != "https://github.com/example/gx/pull/1/files#diff-bbbR20" {
		t.Fatalf("kept duplicate title with higher score link = %q", publicationLink)
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
	if got := renderAttributions(attrs, attributionLinkContext{}); got != "" {
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
	linkCtx := newAttributionLinkContext(artifact, catalog, nil)
	rendered := renderAttributions(attrs, linkCtx)
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
	rendered := renderAttributions(attrs, attributionLinkContext{})
	if strings.Contains(rendered, `href="http`) {
		t.Fatalf("renderAttributions() = %q, want no links", rendered)
	}
	if !strings.Contains(rendered, "internal/storage/schema.go:42") {
		t.Fatalf("renderAttributions() = %q, want plain file:line label", rendered)
	}
}

func TestRenderAttributionsLinksInPRFilePath(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/21"
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Push: reviewbundle.PushPayload{HeadCommitID: "abc", GitHubPullRequestURL: &prURL},
	})
	catalog := prBodyCatalog{
		Hunks: []prHunkSummary{{
			File:        "internal/github/client.go",
			NewStart:    2,
			NewLines:    1,
			ChangedLine: 2,
			Link:        githubHunkLineLink(prURL, "internal/github/client.go", 2),
		}},
	}
	linkCtx := newAttributionLinkContext(artifact, catalog, nil)
	rendered := renderAttributions([]prAttribution{{
		Kind: "codebase",
		Ref:  "internal/github/client.go",
	}}, linkCtx)
	want := githubHunkLineLink(prURL, "internal/github/client.go", 2)
	if !strings.Contains(rendered, want) {
		t.Fatalf("renderAttributions() = %q, want diff link %q", rendered, want)
	}
}

func TestRenderAttributionsLinksPRNumberReference(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/99"
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Push: reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
	})
	linkCtx := newAttributionLinkContext(artifact, prBodyCatalog{}, nil)
	rendered := renderAttributions([]prAttribution{{
		Kind:  "previous_review",
		Label: "follow-up from PR #68",
	}}, linkCtx)
	want := "https://github.com/satoricorp/gx/pull/68"
	if !strings.Contains(rendered, want) {
		t.Fatalf("renderAttributions() = %q, want PR link %q", rendered, want)
	}
}

func TestRenderAttributionsLinksIndexedDocURL(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/11"
	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Push: reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
	})
	snippets := []codereview.ContextSnippet{{
		Kind:      "review_resource",
		Ref:       "google-eng-practices",
		Publisher: "Google Engineering Practices",
		URL:       "https://google.github.io/eng-practices/review/",
	}}
	linkCtx := newAttributionLinkContext(artifact, prBodyCatalog{}, snippets)
	rendered := renderAttributions([]prAttribution{{
		Kind:   "review_resource",
		Label:  "Google Engineering Practices",
		Opaque: true,
	}}, linkCtx)
	want := "https://google.github.io/eng-practices/review/"
	if !strings.Contains(rendered, want) {
		t.Fatalf("renderAttributions() = %q, want doc URL %q", rendered, want)
	}
}

func TestLinkifyAttributionTextOwnerRepoPRReference(t *testing.T) {
	prURL := "https://github.com/satoricorp/gx/pull/1"
	linkCtx := newAttributionLinkContext(reviewbundle.NewArtifact(reviewbundle.Bundle{
		Push: reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
	}), prBodyCatalog{}, nil)
	got := linkCtx.linkifyAttributionText("see satoricorp/gx#68 for context")
	want := "https://github.com/satoricorp/gx/pull/68"
	if !strings.Contains(got, want) {
		t.Fatalf("linkifyAttributionText() = %q, want %q", got, want)
	}
}

func TestReviewVerdictDocsOnlyWithSuccessfulAI(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "> ✅ **No review needed** — documentation-only change (1 file(s)); nothing needs human eyes.") {
		t.Fatalf("body missing no-review verdict banner:\n%s", body)
	}
}

func TestReviewVerdictDocsOnlyWithReviewerError(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{failCount: 2, err: fmt.Errorf("model unavailable")}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "> 👀 **Quick scan** — standard change; AI review unavailable.") {
		t.Fatalf("body missing quick scan verdict banner:\n%s", body)
	}
	if strings.Contains(body, "> ✅ **No review needed**") {
		t.Fatalf("body should not grant no-review when AI failed:\n%s", body)
	}
}

func TestReviewVerdictStrongFinding(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{findings: []codereview.Finding{{
			Title:    "Verify auth handling",
			Summary:  "Session token validation changed.",
			Strength: "Strong",
			File:     "internal/auth/session.go",
			Line:     2,
		}}}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	prURL := "https://github.com/satoricorp/gx/pull/20"
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push:          reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
		Revisions: []reviewbundle.RevisionPayload{{
			Patch:                "diff --git a/internal/auth/session.go b/internal/auth/session.go\n--- a/internal/auth/session.go\n+++ b/internal/auth/session.go\n@@ -1 +1,2 @@\n package auth\n+func Validate() {}\n",
			Files:                []string{"internal/auth/session.go"},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "> 🔴 **Requires Deep Review** — security-sensitive change: internal/auth.") {
		t.Fatalf("body missing deep review verdict banner:\n%s", body)
	}
}

func TestTriageChangeFromCatalogTestsOnly(t *testing.T) {
	catalog := buildPRBodyCatalog(reviewbundle.NewArtifact(reviewbundle.Bundle{
		Revisions: []reviewbundle.RevisionPayload{{
			Patch: strings.Join([]string{
				"diff --git a/internal/app/app_test.go b/internal/app/app_test.go",
				"--- a/internal/app/app_test.go",
				"+++ b/internal/app/app_test.go",
				"@@ -1 +1,2 @@",
				" package app",
				"+func TestMore() {}",
			}, "\n"),
			Files: []string{"internal/app/app_test.go"},
		}},
	}))
	triage := triageChangeFromCatalog(catalog)
	if triage.Class != "tests-only" {
		t.Fatalf("Class = %q, want tests-only", triage.Class)
	}
}

func docsOnlyPRArtifact() reviewbundle.Artifact {
	prURL := "https://github.com/satoricorp/gx/pull/12"
	return reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         "abc123",
			GitHubPullRequestURL: &prURL,
		},
		Revisions: []reviewbundle.RevisionPayload{{
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
			CommitID:             "abc123",
			Description:          "clarify docs",
			Files:                []string{"docs/demo.md"},
			GitHubPullRequestURL: &prURL,
		}},
	})
}

func TestSpeculativeFindingsDroppedFromPRBody(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{findings: []codereview.Finding{{
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
		Revisions: []reviewbundle.RevisionPayload{{
			Patch:                "diff --git a/internal/github/client.go b/internal/github/client.go\n--- a/internal/github/client.go\n+++ b/internal/github/client.go\n@@ -1 +1,2 @@\n package github\n+// change\n",
			Files:                []string{"internal/github/client.go"},
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

func TestOpeningSummaryRendersOverviewFirst(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	overview := "This change clarifies publication ordering so reviewers see PR summaries before side effects complete."
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{overview: overview}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	verdictIdx := strings.Index(body, "> ✅ **No review needed**")
	if verdictIdx < 0 {
		verdictIdx = strings.Index(body, "> 👀 **Quick scan**")
	}
	if verdictIdx < 0 {
		verdictIdx = strings.Index(body, "> 🔴 **Requires Deep Review**")
	}
	overviewIdx := strings.Index(body, overview)
	if overviewIdx < 0 {
		t.Fatalf("body missing overview:\n%s", body)
	}
	if verdictIdx >= 0 && overviewIdx < verdictIdx {
		t.Fatalf("overview should follow verdict line:\n%s", body)
	}
	if strings.Contains(body, "This PR changes") {
		t.Fatalf("body should not use template change summary when overview present:\n%s", body)
	}
	opening := openingParagraphFromPRBody(body)
	if opening != overview {
		t.Fatalf("opening = %q, want exactly overview %q", opening, overview)
	}
	if strings.Contains(opening, "Blast radius") {
		t.Fatalf("opening should not contain blast radius prose:\n%s", opening)
	}
}

func TestOpeningSummaryFallsBackWithoutOverview(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "This PR changes") {
		t.Fatalf("body should fall back to template change summary:\n%s", body)
	}
	opening := openingParagraphFromPRBody(body)
	want := changeSummarySentence(docsOnlyPRArtifact(), buildPRBodyCatalog(docsOnlyPRArtifact()))
	if opening != want {
		t.Fatalf("opening = %q, want template %q", opening, want)
	}
	if strings.Contains(opening, "Blast radius") {
		t.Fatalf("opening should not contain blast radius prose:\n%s", opening)
	}
}

func TestOpeningSummaryBlastRadiusNotBeforeFirstHeading(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	t.Setenv("GX_PR_BLAST_RADIUS", "1")
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{Sources: []string{"codebase"}}, nil
	}
	root, headSHA := initLexicalReachRepo(t)
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reachTestArtifact(root, headSHA))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	opening := openingParagraphFromPRBody(body)
	if strings.Contains(opening, "Blast radius") {
		t.Fatalf("opening should not contain blast radius prose:\n%s", opening)
	}
	if !strings.Contains(body, "## Blast Radius") || !strings.Contains(body, "MEDIUM:") {
		t.Fatalf("body should render blast radius lead in Blast Radius section:\n%s", body)
	}
}

func TestOpeningSummaryContextOnlyInFooter(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{Models: []string{"test-model"}}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{Sources: []string{"codebase", "session"}}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if strings.Contains(openingParagraphFromPRBody(body), "Context used:") {
		t.Fatalf("opening should not contain context clause:\n%s", body)
	}
	if !strings.Contains(body, "*Generated by GX — model test-model; context: codebase, session.*") {
		t.Fatalf("body missing context in provenance footer:\n%s", body)
	}
}

func openingParagraphFromPRBody(body string) string {
	rest := strings.TrimPrefix(body, githubPRBodyMarker+"\n\n")
	if blank := strings.Index(rest, "\n\n"); blank >= 0 {
		rest = rest[blank+2:]
	}
	if idx := strings.Index(rest, "\n\n## "); idx >= 0 {
		return strings.TrimSpace(rest[:idx])
	}
	if idx := strings.Index(rest, "\n\n*Generated by GX"); idx >= 0 {
		return strings.TrimSpace(rest[:idx])
	}
	return strings.TrimSpace(rest)
}

func TestOpeningSummaryStripsURLsFromOverview(t *testing.T) {
	got := sanitizedOverviewParagraph("See https://example.com/docs for context on the auth flow.", true)
	if strings.Contains(got, "https://") || strings.Contains(got, "example.com") {
		t.Fatalf("sanitizedOverviewParagraph() = %q, want URL stripped", got)
	}
	if !strings.Contains(got, "auth flow") {
		t.Fatalf("sanitizedOverviewParagraph() = %q, want non-URL text preserved", got)
	}
}

func TestOpeningSummaryFallsBackWhenAIFailed(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{
			failCount: 2,
			overview:  "This overview should not render when AI failed.",
			err:       fmt.Errorf("model unavailable"),
		}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if strings.Contains(body, "This overview should not render") {
		t.Fatalf("body should not render overview when AI failed:\n%s", body)
	}
	if !strings.Contains(body, "This PR changes") {
		t.Fatalf("body should fall back to template change summary:\n%s", body)
	}
}

func TestHunkImpactConfiguredRiskPath(t *testing.T) {
	root := t.TempDir()
	reviewMD := strings.Join([]string{
		"## high-risk paths",
		"risk-path: internal/auth/** — auth changes can leak or misuse credentials",
	}, "\n")
	if err := os.WriteFile(filepath.Join(root, "REVIEW.md"), []byte(reviewMD), 0o644); err != nil {
		t.Fatalf("write REVIEW.md: %v", err)
	}
	policy := codereview.LoadReviewPolicy(context.Background(), root)
	catalog := prBodyCatalog{}
	hunk := prHunkSummary{File: "internal/auth/session.go"}
	score, title, detail := hunkImpact(catalog, hunk, policy)
	if score != 700 {
		t.Fatalf("score = %d, want 700", score)
	}
	if !strings.Contains(detail, "internal/auth/**") || !strings.Contains(detail, "session.go") {
		t.Fatalf("detail = %q, want matched glob and file", detail)
	}
	if !strings.Contains(detail, "auth changes can leak or misuse credentials") {
		t.Fatalf("detail = %q, want configured message", detail)
	}
	if title == "" {
		t.Fatalf("title = %q, want non-empty title derived from message", title)
	}
}

func TestHunkImpactGenericFallbackWithoutReviewPolicy(t *testing.T) {
	policy := codereview.ReviewPolicy{}
	catalog := prBodyCatalog{}
	hunk := prHunkSummary{File: "pkg/session/token.go"}
	score, title, detail := hunkImpact(catalog, hunk, policy)
	if score != 650 {
		t.Fatalf("score = %d, want 650 for generic auth pattern", score)
	}
	if !strings.Contains(title, "token.go") {
		t.Fatalf("title = %q, want basename", title)
	}
	if !strings.Contains(detail, "auth|token|secret|credential") || !strings.Contains(detail, "pkg/session/token.go") {
		t.Fatalf("detail = %q, want generic pattern and file path", detail)
	}
}

func TestGitHubPullRequestBodyUsesConfiguredRiskPath(t *testing.T) {
	root := t.TempDir()
	reviewMD := strings.Join([]string{
		"## high-risk paths",
		"risk-path: internal/billing/** — billing changes can mischarge customers",
	}, "\n")
	if err := os.WriteFile(filepath.Join(root, "REVIEW.md"), []byte(reviewMD), 0o644); err != nil {
		t.Fatalf("write REVIEW.md: %v", err)
	}
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return nil, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	prURL := "https://github.com/example/acme/pull/1"
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: root},
		Push:          reviewbundle.PushPayload{GitHubPullRequestURL: &prURL},
		Revisions: []reviewbundle.RevisionPayload{{
			Patch: strings.Join([]string{
				"diff --git a/internal/billing/charge.go b/internal/billing/charge.go",
				"--- a/internal/billing/charge.go",
				"+++ b/internal/billing/charge.go",
				"@@ -1 +1,2 @@",
				" package billing",
				"+func Charge() {}",
			}, "\n"),
			Files:                []string{"internal/billing/charge.go"},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "billing changes can mischarge customers") {
		t.Fatalf("body missing configured risk-path message:\n%s", body)
	}
	if !strings.Contains(body, "internal/billing/**") {
		t.Fatalf("body missing matched glob:\n%s", body)
	}
}

func TestReviewPRSummaryFindingsRetriesAfterTransientError(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldWait := prSummaryReviewRetryWait
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		prSummaryReviewRetryWait = oldWait
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewRetryWait = func(context.Context, time.Duration) error { return nil }
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{failCount: 1}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if !strings.Contains(body, "> ✅ **No review needed** — documentation-only change (1 file(s)); nothing needs human eyes.") {
		t.Fatalf("body missing no-review verdict banner after retry:\n%s", body)
	}
	if strings.Contains(body, "heuristics only") {
		t.Fatalf("body should not show heuristics-only footer after successful retry:\n%s", body)
	}
}

func TestProvenanceFooterHeuristicsOnlyWhenAllAttemptsFail(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldWait := prSummaryReviewRetryWait
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		prSummaryReviewRetryWait = oldWait
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewRetryWait = func(context.Context, time.Duration) error { return nil }
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{failCount: 2, err: fmt.Errorf("model unavailable")}, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	want := "*Generated by GX — heuristics only (AI unavailable).*"
	if !strings.Contains(body, want) {
		t.Fatalf("body missing heuristics-only footer:\n%s", body)
	}
}

func TestProvenanceFooterListsModels(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{Models: []string{"test-model"}}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{Sources: []string{"codebase", "session"}}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	want := "*Generated by GX — model test-model; context: codebase, session.*"
	if !strings.Contains(body, want) {
		t.Fatalf("body missing model provenance footer, want %q:\n%s", want, body)
	}
}

func TestReviewVerdictTierStrings(t *testing.T) {
	triage := codereview.ChangeTriage{Class: "docs-only"}
	stats := prBodyStats{FileCount: 1, MaxRiskLevel: "low"}
	if verdict, reason := reviewVerdict(triage, stats, nil, true, []string{"docs/readme.md"}, []string{"docs"}); verdict != verdictNoReview || reason != "documentation-only change (1 file(s)); nothing needs human eyes" {
		t.Fatalf("reviewVerdict() = %q, %q; want no-review tier", verdict, reason)
	}

	triage = codereview.ChangeTriage{Class: "code"}
	if verdict, reason := reviewVerdict(triage, stats, nil, false, nil, nil); verdict != verdictQuickScan || reason != "standard change; AI review unavailable" {
		t.Fatalf("reviewVerdict() = %q, %q; want quick scan when AI unavailable", verdict, reason)
	}
	if verdict, reason := reviewVerdict(triage, stats, nil, true, nil, nil); verdict != verdictQuickScan || reason != "standard change; skim the notable changes" {
		t.Fatalf("reviewVerdict() = %q, %q; want quick scan when AI succeeds", verdict, reason)
	}

	triage = codereview.ChangeTriage{Class: "security-sensitive"}
	if verdict, reason := reviewVerdict(triage, stats, nil, true, []string{"internal/auth/session.go"}, []string{"internal/auth"}); verdict != verdictDeepReview || reason != "security-sensitive change: internal/auth" {
		t.Fatalf("reviewVerdict() = %q, %q; want deep review for security-sensitive", verdict, reason)
	}

	highStats := prBodyStats{MaxRiskLevel: "high", RiskSignals: []string{"lexical_reach:30_refs_in_10_files"}}
	if verdict, reason := reviewVerdict(codereview.ChangeTriage{Class: "code"}, highStats, nil, true, nil, nil); verdict != verdictDeepReview || reason != "changes are referenced widely across the codebase" {
		t.Fatalf("reviewVerdict() = %q, %q; want deep review for lexical reach", verdict, reason)
	}

	highStats = prBodyStats{MaxRiskLevel: "high", RiskSignals: []string{"auth_change"}}
	if verdict, reason := reviewVerdict(codereview.ChangeTriage{Class: "code"}, highStats, nil, true, nil, nil); verdict != verdictDeepReview || reason != "high-risk change signals" {
		t.Fatalf("reviewVerdict() = %q, %q; want deep review for high risk", verdict, reason)
	}

	findings := []codereview.Finding{{Strength: "Strong"}}
	if verdict, reason := reviewVerdict(codereview.ChangeTriage{Class: "code"}, stats, findings, true, nil, nil); verdict != verdictDeepReview || reason != "review findings need human judgment" {
		t.Fatalf("reviewVerdict() = %q, %q; want deep review for strong findings", verdict, reason)
	}
}

func TestRenderVerdictBannerFormat(t *testing.T) {
	tests := []struct {
		verdict string
		reason  string
		want    string
	}{
		{verdictNoReview, "documentation-only change (1 file(s)); nothing needs human eyes", "> ✅ **No review needed** — documentation-only change (1 file(s)); nothing needs human eyes."},
		{verdictQuickScan, "standard change; skim the notable changes", "> 👀 **Quick scan** — standard change; skim the notable changes."},
		{verdictDeepReview, "security-sensitive change: internal/auth", "> 🔴 **Requires Deep Review** — security-sensitive change: internal/auth."},
	}
	for _, tc := range tests {
		if got := renderVerdictBanner(tc.verdict, tc.reason); got != tc.want {
			t.Fatalf("renderVerdictBanner(%q, %q) = %q, want %q", tc.verdict, tc.reason, got, tc.want)
		}
	}
}

// A reader has to be able to tell a context-rich summary from a diff-only one,
// and the note must never fire on a PR that really was captured.
func TestSessionContextInformedSummarySignals(t *testing.T) {
	linkedContext := func(context *reviewbundle.ReviewContextPayload) reviewbundle.Artifact {
		artifact := docsOnlyPRArtifact()
		artifact.Bundle.Revisions[0].ReviewContext = context
		return artifact
	}
	tests := []struct {
		name           string
		artifact       reviewbundle.Artifact
		summaryContext prSummaryContext
		wantNote       bool
	}{
		{
			name:           "nothing captured",
			artifact:       docsOnlyPRArtifact(),
			summaryContext: prSummaryContext{Sources: []string{"codebase"}},
			wantNote:       true,
		},
		{
			name:     "session transcript reached the model",
			artifact: docsOnlyPRArtifact(),
			summaryContext: prSummaryContext{
				Snippets: []codereview.ContextSnippet{{Kind: "session_transcript", Source: "artifact", Ref: "abc123:sess-1:req-1"}},
				Sources:  []string{"codebase", "session"},
			},
			wantNote: false,
		},
		{
			name:     "indexed session hit is not this branch's session",
			artifact: docsOnlyPRArtifact(),
			summaryContext: prSummaryContext{
				Snippets: []codereview.ContextSnippet{{Kind: "session_transcript", Source: "turbopuffer:gx-sessions", Ref: "other-branch"}},
				Sources:  []string{"indexed code/session"},
			},
			wantNote: true,
		},
		{
			name:           "sessions linked but no transcript survived into the bundle",
			artifact:       linkedContext(&reviewbundle.ReviewContextPayload{ProvenanceStatus: reviewsource.StatusLinked, LinkedSessionCount: 1}),
			summaryContext: prSummaryContext{},
			wantNote:       false,
		},
		{
			name:           "repo-local sessions still count as captured",
			artifact:       linkedContext(&reviewbundle.ReviewContextPayload{ProvenanceStatus: reviewsource.StatusRepoLocal, LinkedSessionCount: 2}),
			summaryContext: prSummaryContext{},
			wantNote:       false,
		},
		{
			name:           "review context that matched no session",
			artifact:       linkedContext(&reviewbundle.ReviewContextPayload{ProvenanceStatus: reviewsource.StatusAbsent}),
			summaryContext: prSummaryContext{},
			wantNote:       true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			note := missingSessionContextNote(tc.artifact, tc.summaryContext)
			if gotNote := note != ""; gotNote != tc.wantNote {
				t.Fatalf("missingSessionContextNote() = %q, want note = %v", note, tc.wantNote)
			}
			if tc.wantNote && !strings.Contains(note, "pre-push hook") {
				t.Fatalf("note does not say what causes it: %q", note)
			}
		})
	}
}

func TestGitHubPullRequestBodyPlacesMissingSessionNoteAboveTheFooter(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{}, codereview.ReviewerInfo{}
	}

	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{Sources: []string{"codebase"}}, nil
	}
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	noteIdx := strings.Index(body, "*No captured coding session was linked to this branch")
	footerIdx := strings.LastIndex(body, "*Generated by GX")
	if noteIdx < 0 {
		t.Fatalf("body missing the missing-session-context note:\n%s", body)
	}
	if noteIdx >= footerIdx {
		t.Fatalf("note should sit above the provenance footer:\n%s", body)
	}
	if between := body[noteIdx:footerIdx]; strings.Count(between, "\n\n") != 1 {
		t.Fatalf("note should be the paragraph immediately above the footer, got %q", between)
	}
	if strings.Contains(openingParagraphFromPRBody(body), "No captured coding session") {
		t.Fatalf("note should not crowd the opening summary:\n%s", body)
	}

	// With session context the body says nothing extra.
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{
			Snippets: []codereview.ContextSnippet{{Kind: "session_transcript", Source: "artifact", Ref: "abc123:sess-1:req-1"}},
			Sources:  []string{"codebase", "session"},
		}, nil
	}
	body, err = GitHubPullRequestBodyFromArtifact(context.Background(), docsOnlyPRArtifact())
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if strings.Contains(body, "No captured coding session") {
		t.Fatalf("body should stay quiet when session context is present:\n%s", body)
	}
}
