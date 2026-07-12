package publication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return &fakePRSummaryReviewer{findings: []codereview.Finding{{
			Title: "Verify GitHub PR summary update ordering",
		}}}, codereview.ReviewerInfo{}
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

	if _, err := EnqueueArtifact(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL)), QueueAttestation{}); err != nil {
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
	if drain.PRSummariesUpdated != 0 || drain.PRSummariesFailed != 0 {
		t.Fatalf("drain = %#v, want no CLI PR body updates (server owns summaries)", drain)
	}
	if patchedBody != "" {
		t.Fatalf("DrainQueuedUploads() should not rewrite GitHub PR bodies:\n%s", patchedBody)
	}
	if !uploader.called {
		t.Fatal("uploader was not called during drain")
	}
}

func TestDrainQueuedUploadsContinuesUploadWhenPRSummaryFails(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) { return nil, codereview.ReviewerInfo{} }
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

	if _, err := EnqueueArtifact(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL)), QueueAttestation{}); err != nil {
		t.Fatalf("EnqueueArtifact() error = %v", err)
	}
	uploader := &fakeUploader{}
	drain, err := DrainQueuedUploads(context.Background(), uploader, 10)
	if err != nil {
		t.Fatalf("DrainQueuedUploads() error = %v", err)
	}
	if drain.Uploaded != 1 || drain.Failed != 0 || drain.PRSummariesFailed != 0 {
		t.Fatalf("drain = %#v uploader.called=%t, want upload with no CLI PR body rewrite", drain, uploader.called)
	}
}

func TestDrainQueuedUploadsDoesNotRewriteGitHubPullRequestBody(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	var patchedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			var payload map[string]string
			_ = json.NewDecoder(r.Body).Decode(&payload)
			patchedBody = payload["body"]
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"html_url": prURL, "number": 11, "body": ""})
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	if _, err := EnqueueArtifact(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL)), QueueAttestation{}); err != nil {
		t.Fatalf("EnqueueArtifact() error = %v", err)
	}
	uploader := &fakeUploader{}
	if _, err := DrainQueuedUploads(context.Background(), uploader, 10); err != nil {
		t.Fatalf("DrainQueuedUploads() error = %v", err)
	}
	if patchedBody != "" {
		t.Fatalf("DrainQueuedUploads() rewrote PR body; server owns summaries:\n%s", patchedBody)
	}
	if !uploader.called {
		t.Fatal("uploader was not called during drain")
	}
}

func TestUpdateGitHubPullRequestBodyIsNoOp(t *testing.T) {
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	var patched bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			patched = true
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"html_url": prURL,
			"number":   11,
			"body":     "Please preserve these author notes.",
		})
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	updated, err := UpdateGitHubPullRequestBody(context.Background(), reviewbundle.NewArtifact(prSummaryTestBundle(prURL)))
	if err != nil {
		t.Fatalf("UpdateGitHubPullRequestBody() error = %v", err)
	}
	if updated {
		t.Fatal("UpdateGitHubPullRequestBody() updated = true, want no-op false")
	}
	if patched {
		t.Fatal("UpdateGitHubPullRequestBody() should not PATCH GitHub PR bodies")
	}
}

type fakePRSummaryReviewer struct {
	findings         []codereview.Finding
	overview         string
	downstreamImpact string
	notableChanges   []codereview.NotableChange
	err              error
	failCount        int
	attempts         int
}

func (f *fakePRSummaryReviewer) Review(context.Context, codereview.ReviewBrief) ([]codereview.Finding, error) {
	f.attempts++
	if f.attempts <= f.failCount {
		if f.err != nil {
			return nil, f.err
		}
		return nil, fmt.Errorf("transient reviewer failure")
	}
	return f.findings, nil
}

func (f *fakePRSummaryReviewer) ReviewForSummary(_ context.Context, _ codereview.ReviewBrief) (codereview.PRSummaryReview, error) {
	f.attempts++
	if f.attempts <= f.failCount {
		if f.err != nil {
			return codereview.PRSummaryReview{}, f.err
		}
		return codereview.PRSummaryReview{}, fmt.Errorf("transient reviewer failure")
	}
	return codereview.PRSummaryReview{
		Overview:         f.overview,
		DownstreamImpact: f.downstreamImpact,
		NotableChanges:   f.notableChanges,
		Findings:         f.findings,
	}, nil
}

func (f *fakePRSummaryReviewer) ReviewWithOverview(_ context.Context, _ codereview.ReviewBrief) (string, []codereview.Finding, error) {
	f.attempts++
	if f.attempts <= f.failCount {
		if f.err != nil {
			return "", nil, f.err
		}
		return "", nil, fmt.Errorf("transient reviewer failure")
	}
	return f.overview, f.findings, nil
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
