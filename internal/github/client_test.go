package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindPullRequestUsesGitHubAPI(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if r.URL.Path != "/repos/satoricorp/gx/pulls" {
			t.Fatalf("path = %q, want /repos/satoricorp/gx/pulls", r.URL.Path)
		}
		if r.URL.Query().Get("head") != "satoricorp:feature/demo" {
			t.Fatalf("head query = %q", r.URL.Query().Get("head"))
		}
		_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/satoricorp/gx/pull/7","body":"existing body"}]`))
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	pr, err := client.FindPullRequest(context.Background(), CreatePullRequestOptions{
		Owner:      "satoricorp",
		Repo:       "gx",
		HeadBranch: "feature/demo",
	})
	if err != nil {
		t.Fatalf("FindPullRequest() error = %v", err)
	}
	if pr == nil || pr.URL != "https://github.com/satoricorp/gx/pull/7" || pr.Number != 7 || pr.Body != "existing body" {
		t.Fatalf("FindPullRequest() = %#v", pr)
	}
	if gotAuth != "Bearer token-one" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
}

func TestFindPullRequestByHeadCanIncludeMergedPRs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/satoricorp/gx/pulls" {
			t.Fatalf("path = %q, want /repos/satoricorp/gx/pulls", r.URL.Path)
		}
		if r.URL.Query().Get("state") != "all" {
			t.Fatalf("state query = %q, want all", r.URL.Query().Get("state"))
		}
		_, _ = w.Write([]byte(`[{"number":7,"html_url":"https://github.com/satoricorp/gx/pull/7","state":"closed","merged_at":"2026-06-28T12:00:00Z"}]`))
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	pr, err := client.FindPullRequestByHead(context.Background(), CreatePullRequestOptions{
		Owner:      "satoricorp",
		Repo:       "gx",
		HeadBranch: "feature/demo",
	}, "all")
	if err != nil {
		t.Fatalf("FindPullRequestByHead() error = %v", err)
	}
	if pr == nil || pr.Number != 7 || !pr.Merged || pr.State != "closed" {
		t.Fatalf("FindPullRequestByHead() = %#v, want merged closed PR", pr)
	}
}

func TestUpsertIssueCommentCreatesWhenMissing(t *testing.T) {
	var methods []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method+" "+r.URL.Path)
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/satoricorp/gx/issues/8/comments":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == http.MethodPost && r.URL.Path == "/repos/satoricorp/gx/issues/8/comments":
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if payload["body"] != "new body" {
				t.Fatalf("body = %q, want new body", payload["body"])
			}
			_, _ = w.Write([]byte(`{"id":12,"html_url":"https://github.com/satoricorp/gx/pull/8#issuecomment-12","body":"new body"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	comment, err := client.UpsertIssueComment(context.Background(), IssueCommentOptions{
		Owner:  "satoricorp",
		Repo:   "gx",
		Number: 8,
		Body:   "new body",
		Marker: "<!-- marker -->",
	})
	if err != nil {
		t.Fatalf("UpsertIssueComment() error = %v", err)
	}
	if comment == nil || comment.ID != 12 {
		t.Fatalf("comment = %#v, want created comment", comment)
	}
	if got := strings.Join(methods, ","); got != "GET /repos/satoricorp/gx/issues/8/comments,POST /repos/satoricorp/gx/issues/8/comments" {
		t.Fatalf("requests = %q", got)
	}
}

func TestUpsertIssueCommentUpdatesExistingMarker(t *testing.T) {
	var patched bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/satoricorp/gx/issues/8/comments":
			_, _ = w.Write([]byte(`[{"id":22,"html_url":"https://github.com/satoricorp/gx/pull/8#issuecomment-22","body":"old\n<!-- marker -->"}]`))
		case r.Method == http.MethodPatch && r.URL.Path == "/repos/satoricorp/gx/issues/comments/22":
			patched = true
			var payload map[string]string
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if payload["body"] != "replacement" {
				t.Fatalf("body = %q, want replacement", payload["body"])
			}
			_, _ = w.Write([]byte(`{"id":22,"html_url":"https://github.com/satoricorp/gx/pull/8#issuecomment-22","body":"replacement"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	comment, err := client.UpsertIssueComment(context.Background(), IssueCommentOptions{
		Owner:  "satoricorp",
		Repo:   "gx",
		Number: 8,
		Body:   "replacement",
		Marker: "<!-- marker -->",
	})
	if err != nil {
		t.Fatalf("UpsertIssueComment() error = %v", err)
	}
	if comment == nil || comment.ID != 22 || !patched {
		t.Fatalf("comment = %#v patched=%t, want updated comment", comment, patched)
	}
}

func TestCreatePullRequestReviewCommentUsesGitHubAPI(t *testing.T) {
	var payload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/repos/satoricorp/gx/pulls/8/comments" {
			t.Fatalf("path = %q, want /repos/satoricorp/gx/pulls/8/comments", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		_, _ = w.Write([]byte(`{"id":44}`))
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	err := client.CreatePullRequestReviewComment(context.Background(), PullRequestReviewCommentOptions{
		Owner:    "satoricorp",
		Repo:     "gx",
		Number:   8,
		Body:     "review body",
		CommitID: "abc123",
		Path:     "main.go",
		Line:     12,
	})
	if err != nil {
		t.Fatalf("CreatePullRequestReviewComment() error = %v", err)
	}
	if payload["body"] != "review body" ||
		payload["commit_id"] != "abc123" ||
		payload["path"] != "main.go" ||
		int(payload["line"].(float64)) != 12 ||
		payload["side"] != "RIGHT" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestGitHubAuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad token", http.StatusUnauthorized)
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "bad-token", server.Client())
	_, err := client.FindPullRequest(context.Background(), CreatePullRequestOptions{Owner: "satoricorp", Repo: "gx", HeadBranch: "feature/demo"})
	if err == nil || !strings.Contains(err.Error(), "github auth failed") {
		t.Fatalf("FindPullRequest() error = %v, want auth failure", err)
	}
}

func TestNewClientUsesStoredGitHubToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	if err := os.WriteFile(filepath.Join(home, "credentials.json"), []byte(`{
  "cloud": {
    "token": "gx_saved_token_with_enough_length",
    "github_access_token": "stored-github-token"
  }
}`), 0o600); err != nil {
		t.Fatalf("write credentials: %v", err)
	}
	client, err := NewClient("github.com")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.token != "stored-github-token" {
		t.Fatalf("client token = %q", client.token)
	}
}
