package cloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCloudBaseURLStripsGxPrSuffix(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "http://localhost:3200/gx/pr")
	if got := CloudBaseURL(); got != "http://localhost:3200" {
		t.Fatalf("CloudBaseURL() = %q, want http://localhost:3200", got)
	}
}

func TestRepoFullNameFromRemoteURL(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"git@github.com:satoricorp/gx.git", "satoricorp/gx"},
		{"https://github.com/example/widget", "example/widget"},
		{"https://github.com/example/widget.git", "example/widget"},
	}
	for _, tc := range tests {
		if got := RepoFullNameFromRemoteURL(tc.in); got != tc.want {
			t.Fatalf("RepoFullNameFromRemoteURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestDeleteBookmarkUsesV1Endpoint(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken: "gho_bookmark",
		CLISessionToken:   "gxcs_bookmark",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v1/bookmarks/bookmark-1" {
			t.Fatalf("path = %s, want /v1/bookmarks/bookmark-1", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer gxcs_bookmark" {
			t.Fatalf("authorization = %q", auth)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &Client{url: server.URL, http: server.Client()}
	if err := client.DeleteBookmark(context.Background(), "bookmark-1"); err != nil {
		t.Fatalf("DeleteBookmark() error = %v", err)
	}
}
