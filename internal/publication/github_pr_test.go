package publication

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestEnqueueArtifactUpdatesGitHubPullRequestBodyFromReviewBundle(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	branchName := "bug/fix-github-pr-summary"
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
				"body":     "Published by GX.\n\nRevisions:\n- old fallback body",
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

	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Push: reviewbundle.PushPayload{
			BranchName:           stringPtr("main"),
			HeadCommitID:         "commit-head",
			GitHubPullRequestURL: &prURL,
		},
		Stack: []reviewbundle.StackPayload{{
			BranchName:     branchName,
			BaseBranchName: "main",
			Patch:          "diff --git a/internal/publication/github_pr.go b/internal/publication/github_pr.go\n+added\n-old\n",
			Change: reviewbundle.ChangePayload{
				JJChangeID:      "change-one",
				CurrentCommitID: "abcdef123456",
				Description:     "render rich PR bodies from review bundles",
				Files: []string{
					"internal/github/client.go",
					"internal/github/client_test.go",
					"internal/publication/github_pr.go",
					"internal/publication/github_pr_test.go",
					"internal/publication/outbox.go",
					"internal/publication/pr_body.go",
					"internal/publication/publication.go",
				},
				ReviewContext: &reviewbundle.ReviewContextPayload{
					ProvenanceStatus:   "explicit",
					LinkedSessionCount: 1,
					ProvenanceSources: []reviewbundle.ReviewProvenanceSource{{
						SessionID: "session-one",
						Status:    "explicit",
						Source:    "change_sessions",
					}},
					AgentProvenance: []reviewbundle.ReviewAgentProvenance{{
						SessionID: "session-one",
						AgentTool: "codex",
						Provider:  "openai",
						ModelID:   "gpt-5",
					}},
					StructuralStatus: "available",
					ChangedSymbols: []reviewbundle.ReviewChangedSymbol{{
						File:   "internal/publication/github_pr.go",
						Symbol: "UpdateGitHubPullRequestBody",
						Kind:   "function",
					}},
					FeasibilityWarnings: []reviewbundle.ReviewFeasibilityWarning{{
						Severity: "warning",
						Source:   "structural_dependency",
						Message:  "publication must update GitHub after review bundle assembly",
					}},
					Risk: reviewbundle.RiskPayload{
						Level:   "high",
						Score:   80,
						Signals: []string{"structural_dependencies", "warning:structural_dependency"},
					},
				},
			},
			GitHubPullRequestURL: &prURL,
		}},
	})

	if _, err := EnqueueArtifact(context.Background(), artifact); err != nil {
		t.Fatalf("EnqueueArtifact() error = %v", err)
	}
	for _, want := range []string{
		"Publishes 1 GX revision(s) for `bug/fix-github-pr-summary` with high blast radius.",
		"## Review signals",
		"Blast radius: high",
		"Risk: high (80/100; structural dependencies, warning:structural dependency)",
		"Provenance: explicit (1 linked session(s)); attribution: codex/openai/gpt-5",
		"Structural context: available (1 changed symbol(s))",
		"Feasibility: 1 warning(s), 0 info item(s)",
		"[render rich PR bodies from review bundles](https://github.com/satoricorp/gx/commit/abcdef123456)",
		"`internal/publication/github_pr.go`",
		"and 1 more",
		"`internal/publication/github_pr.go:UpdateGitHubPullRequestBody`",
		"[Files changed](https://github.com/satoricorp/gx/pull/11/files)",
	} {
		if !strings.Contains(patchedBody, want) {
			t.Fatalf("patched body missing %q:\n%s", want, patchedBody)
		}
	}
	if strings.Contains(patchedBody, "and%201%20more") {
		t.Fatalf("patched body linked truncated file count as a path:\n%s", patchedBody)
	}
}

func TestUpdateGitHubPullRequestBodyPreservesNonGXBody(t *testing.T) {
	t.Setenv("GH_TOKEN", "token-one")
	prURL := "https://github.com/satoricorp/gx/pull/11"
	branchName := "feature/demo"
	patchCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"html_url": prURL,
				"number":   11,
				"body":     "human-authored body",
			})
		case http.MethodPatch:
			patchCalled = true
			t.Fatalf("unexpected PATCH for non-GX body")
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	updated, err := UpdateGitHubPullRequestBody(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Push: reviewbundle.PushPayload{
			BranchName:           &branchName,
			HeadCommitID:         "abc123",
			GitHubPullRequestURL: &prURL,
		},
		Change: &reviewbundle.ChangePayload{Description: "demo"},
	}))
	if err != nil {
		t.Fatalf("UpdateGitHubPullRequestBody() error = %v", err)
	}
	if updated || patchCalled {
		t.Fatalf("updated=%t patchCalled=%t, want no update", updated, patchCalled)
	}
}

func stringPtr(value string) *string {
	return &value
}
