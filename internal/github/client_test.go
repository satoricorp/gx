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
		_, _ = w.Write([]byte(`[{"html_url":"https://github.com/satoricorp/gx/pull/7","number":7,"body":"Published by GX."}]`))
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
	if pr.Number != 7 || pr.Body != "Published by GX." {
		t.Fatalf("FindPullRequest() metadata = %#v, want number and body", pr)
	}
	if gotAuth != "Bearer token-one" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
}

func TestGetPullRequestUsesGitHubAPI(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/repos/satoricorp/gx/pulls/11" {
			t.Fatalf("path = %q, want pull endpoint", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"html_url":"https://github.com/satoricorp/gx/pull/11","number":11,"body":"Published by GX."}`))
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	pr, err := client.GetPullRequest(context.Background(), GetPullRequestOptions{
		Owner:  "satoricorp",
		Repo:   "gx",
		Number: 11,
	})
	if err != nil {
		t.Fatalf("GetPullRequest() error = %v", err)
	}
	if pr == nil || pr.URL != "https://github.com/satoricorp/gx/pull/11" || pr.Number != 11 || pr.Body != "Published by GX." {
		t.Fatalf("GetPullRequest() = %#v", pr)
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
		Body:       "Published by GX.",
	})
	if err != nil {
		t.Fatalf("CreatePullRequest() error = %v", err)
	}
	if pr == nil || pr.URL != "https://github.com/satoricorp/gx/pull/8" {
		t.Fatalf("CreatePullRequest() = %#v", pr)
	}
}

func TestUpdatePullRequestUsesGitHubAPI(t *testing.T) {
	var gotPayload map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("method = %s, want PATCH", r.Method)
		}
		if r.URL.Path != "/repos/satoricorp/gx/pulls/8" {
			t.Fatalf("path = %q, want pull endpoint", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		_, _ = w.Write([]byte(`{"html_url":"https://github.com/satoricorp/gx/pull/8","number":8,"body":"Published by GX.\n\n## Summary\nDemo"}`))
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)

	client := NewClientWithToken("github.com", "token-one", server.Client())
	pr, err := client.UpdatePullRequest(context.Background(), UpdatePullRequestOptions{
		Owner:  "satoricorp",
		Repo:   "gx",
		Number: 8,
		Body:   "Published by GX.\n\n## Summary\nDemo",
	})
	if err != nil {
		t.Fatalf("UpdatePullRequest() error = %v", err)
	}
	if pr == nil || pr.URL != "https://github.com/satoricorp/gx/pull/8" || pr.Number != 8 {
		t.Fatalf("UpdatePullRequest() = %#v", pr)
	}
	if gotPayload["body"] != "Published by GX.\n\n## Summary\nDemo" {
		t.Fatalf("payload body = %q", gotPayload["body"])
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
