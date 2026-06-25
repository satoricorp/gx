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
		_, _ = w.Write([]byte(`[{"html_url":"https://github.com/satoricorp/gx/pull/7"}]`))
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
	if pr == nil || pr.URL != "https://github.com/satoricorp/gx/pull/7" {
		t.Fatalf("FindPullRequest() = %#v", pr)
	}
	if gotAuth != "Bearer token-one" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
}

func TestCreatePullRequestUsesGitHubAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["base"] != "main" || payload["head"] != "feature/demo" || payload["title"] != "Demo" {
			t.Fatalf("payload = %#v", payload)
		}
		_, _ = w.Write([]byte(`{"html_url":"https://github.com/satoricorp/gx/pull/8"}`))
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	pr, err := client.CreatePullRequest(context.Background(), CreatePullRequestOptions{
		Owner:      "satoricorp",
		Repo:       "gx",
		BaseBranch: "main",
		HeadBranch: "feature/demo",
		Title:      "Demo",
		Body:       "",
	})
	if err != nil {
		t.Fatalf("CreatePullRequest() error = %v", err)
	}
	if pr == nil || pr.URL != "https://github.com/satoricorp/gx/pull/8" {
		t.Fatalf("CreatePullRequest() = %#v", pr)
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
