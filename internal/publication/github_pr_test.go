package publication

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestEnqueueArtifactQueuesWithoutUpdatingGitHubPullRequestBody(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	oldReviewer := prSummaryReviewerFromEnv
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnv = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnv = func() codereview.AIReviewer {
		return fakePRSummaryReviewer{findings: []codereview.Finding{{
			Title: "Verify GitHub PR summary update ordering",
		}}}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{Sources: []string{"codebase"}}, nil
	}
	var patchedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"html_url": prURL,
				"number":   11,
				"body":     "",
			})
		case http.MethodPatch:
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode patch payload: %v", err)
			}
			patchedBody = payload["body"]
			_ = json.NewEncoder(w).Encode(map[string]any{"html_url": prURL, "number": 11, "body": patchedBody})
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	if _, err := EnqueueArtifact(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL))); err != nil {
		t.Fatalf("EnqueueArtifact() error = %v", err)
	}
	if patchedBody != "" {
		t.Fatalf("EnqueueArtifact() updated GitHub PR body synchronously:\n%s", patchedBody)
	}

	uploader := &fakeUploader{}
	drain, err := DrainQueuedUploads(context.Background(), uploader, 10)
	if err != nil {
		t.Fatalf("DrainQueuedUploads() error = %v", err)
	}
	if drain.PRSummariesUpdated != 1 || drain.PRSummariesFailed != 0 {
		t.Fatalf("drain = %#v, want one PR summary updated", drain)
	}
	if !strings.Contains(patchedBody, githubPRBodyMarker) {
		t.Fatalf("patched body missing GX summary marker:\n%s", patchedBody)
	}
	if !uploader.called {
		t.Fatal("uploader was not called during drain")
	}
}

func TestDrainQueuedUploadsContinuesUploadWhenPRSummaryFails(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	oldReviewer := prSummaryReviewerFromEnv
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnv = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnv = func() codereview.AIReviewer { return nil }
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, errors.New("summary context unavailable")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET only when summary fails before patch", r.Method)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"html_url": prURL, "number": 11, "body": ""})
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	if _, err := EnqueueArtifact(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL))); err != nil {
		t.Fatalf("EnqueueArtifact() error = %v", err)
	}
	uploader := &fakeUploader{}
	drain, err := DrainQueuedUploads(context.Background(), uploader, 10)
	if err != nil {
		t.Fatalf("DrainQueuedUploads() error = %v", err)
	}
	if drain.Uploaded != 1 || drain.Failed != 0 || drain.PRSummariesFailed != 1 {
		t.Fatalf("drain = %#v uploader.called=%t, want upload despite summary failure", drain, uploader.called)
	}
}

func TestDrainQueuedUploadsUpdatesGitHubPullRequestBodyFromReviewBundle(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	oldReviewer := prSummaryReviewerFromEnv
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnv = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnv = func() codereview.AIReviewer {
		return fakePRSummaryReviewer{findings: []codereview.Finding{{
			Title:          "Verify GitHub PR summary update ordering",
			Summary:        "internal/publication/publication.go updates the GitHub PR body around publish side effects, so failures could leave reviewers without the important review targets.",
			Recommendation: "Review the publication hunk before merging and confirm GitHub update failures are visible.",
			Evidence:       []codereview.Evidence{{Label: "Changed hunk", Value: "internal/publication/publication.go"}},
			Strength:       "Strong",
			SourceIDs:      []string{"google-eng-practices"},
		}}}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{
			Sources: []string{"codebase", "session", "review resources"},
			Snippets: []codereview.ContextSnippet{{
				Kind:        "review_resource",
				Ref:         "google-eng-practices",
				Source:      "turbopuffer:gx-review-knowledge",
				SourceLabel: "Google Engineering Practices",
				Text:        "Review remote side effects before returning success.",
			}},
		}, nil
	}
	var patchedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path != "/repos/satoricorp/gx/pulls/11" {
				t.Fatalf("get path = %q, want pull endpoint", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"html_url": prURL,
				"number":   11,
				"body":     "",
			})
		case http.MethodPatch:
			if r.URL.Path != "/repos/satoricorp/gx/pulls/11" {
				t.Fatalf("patch path = %q, want pull endpoint", r.URL.Path)
			}
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode patch payload: %v", err)
			}
			patchedBody = payload["body"]
			_ = json.NewEncoder(w).Encode(map[string]any{
				"html_url": prURL,
				"number":   11,
				"body":     patchedBody,
			})
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	if _, err := EnqueueArtifact(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL))); err != nil {
		t.Fatalf("EnqueueArtifact() error = %v", err)
	}
	if patchedBody != "" {
		t.Fatalf("EnqueueArtifact() updated GitHub PR body before drain:\n%s", patchedBody)
	}
	uploader := &fakeUploader{}
	if _, err := DrainQueuedUploads(context.Background(), uploader, 10); err != nil {
		t.Fatalf("DrainQueuedUploads() error = %v", err)
	}
	for _, want := range []string{
		githubPRBodyMarker,
		"This PR changes",
		"Blast radius is high",
		"Context used: codebase, review resources, session.",
		"## Needs Review",
		"Verify GitHub PR summary update ordering",
		"Attribution:",
		"google-eng-practices",
		githubHunkLink(prURL, "internal/publication/publication.go", 180, 180, 9),
	} {
		if !strings.Contains(patchedBody, want) {
			t.Fatalf("patched body missing %q:\n%s", want, patchedBody)
		}
	}
	if strings.Contains(patchedBody, "Published by GX") {
		t.Fatalf("patched body contains legacy placeholder:\n%s", patchedBody)
	}
}

func TestUpdateGitHubPullRequestBodyAppendsHumanBodyAsAuthorNotes(t *testing.T) {
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	oldReviewer := prSummaryReviewerFromEnv
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnv = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnv = func() codereview.AIReviewer { return nil }
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	var patchedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"html_url": prURL,
				"number":   11,
				"body":     "Please preserve these author notes.",
			})
		case http.MethodPatch:
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode patch payload: %v", err)
			}
			patchedBody = payload["body"]
			_ = json.NewEncoder(w).Encode(map[string]any{"html_url": prURL, "number": 11, "body": patchedBody})
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	updated, err := UpdateGitHubPullRequestBody(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL)))
	if err != nil {
		t.Fatalf("UpdateGitHubPullRequestBody() error = %v", err)
	}
	if !updated {
		t.Fatal("UpdateGitHubPullRequestBody() updated = false, want true")
	}
	for _, want := range []string{
		githubPRBodyMarker,
		githubPRAuthorNotesMarker,
		"## Author Notes",
		"Please preserve these author notes.",
	} {
		if !strings.Contains(patchedBody, want) {
			t.Fatalf("patched body missing %q:\n%s", want, patchedBody)
		}
	}
}

type fakePRSummaryReviewer struct {
	findings []codereview.Finding
}

func (f fakePRSummaryReviewer) Review(context.Context, codereview.ReviewBrief) ([]codereview.Finding, error) {
	return f.findings, nil
}

func prSummaryTestBundle(prURL string) reviewbundle.Bundle {
	return reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         "commit-head",
			GitHubPullRequestURL: &prURL,
		},
		Stack: []reviewbundle.StackPayload{{
			BranchName:     "feature/github-pr-summary",
			BaseBranchName: "main",
			Patch: strings.Join([]string{
				"diff --git a/internal/publication/publication.go b/internal/publication/publication.go",
				"--- a/internal/publication/publication.go",
				"+++ b/internal/publication/publication.go",
				"@@ -180,6 +180,9 @@ func (p *Publisher) PublishPush(ctx context.Context, push vcs.PushResult) (Result, error) {",
				"+\tif _, err := UpdateGitHubPullRequestBody(ctx, artifact); err != nil {",
				"+\t\treturn Result{}, err",
				"+\t}",
				" \treturn p.PublishArtifact(ctx, artifact)",
				"",
			}, "\n"),
			Change: reviewbundle.ChangePayload{
				JJChangeID:      "change-one",
				CurrentCommitID: "abcdef123456",
				Description:     "adopt GitHub PR summaries",
				Files: []string{
					"internal/publication/publication.go",
					"internal/publication/github_pr.go",
					"internal/publication/pr_body.go",
				},
				ReviewContext: &reviewbundle.ReviewContextPayload{
					Risk: reviewbundle.RiskPayload{
						Level:   "high",
						Score:   80,
						Signals: []string{"structural_dependencies", "warning:publication_order"},
					},
				},
			},
			GitHubPullRequestURL: &prURL,
		}},
	}
}
