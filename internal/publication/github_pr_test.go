package publication

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/satoricorp/totality/internal/reviewbundle"
)

func TestEnqueueArtifactQueuesWithoutUpdatingGitHubPullRequestBody(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/totality/pull/11"
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
	t.Setenv("TOTALITY_GITHUB_API_URL", server.URL)

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

func TestDrainQueuedUploadsDoesNotRewriteGitHubPullRequestBody(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/totality/pull/11"
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
	t.Setenv("TOTALITY_GITHUB_API_URL", server.URL)

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
	prURL := "https://github.com/satoricorp/totality/pull/11"
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
	t.Setenv("TOTALITY_GITHUB_API_URL", server.URL)

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

func prSummaryTestBundle(prURL string) reviewbundle.Bundle {
	return reviewbundle.Bundle{
		Event:         "tx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         "commit-head",
			GitHubPullRequestURL: &prURL,
		},
		Revisions: []reviewbundle.RevisionPayload{{
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
			RevisionID:  "change-one",
			CommitID:    "abcdef123456",
			Description: "adopt GitHub PR summaries",
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
			GitHubPullRequestURL: &prURL,
		}},
	}
}
