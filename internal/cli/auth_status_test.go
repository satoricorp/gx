package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/satoricorp/gx/internal/cloud"
)

func TestBuildAuthStatusValidatesResolvedTokenAgainstGXAPI(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	var gotAPIAuth string
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/me" {
			t.Fatalf("path = %s, want /v1/auth/me", r.URL.Path)
		}
		gotAPIAuth = r.Header.Get("Authorization")
		if gotAPIAuth != "Bearer gho_api_token" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad token"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"user_id":           "user_1",
			"github_user_login": "octocat",
		})
	}))
	defer apiServer.Close()
	t.Setenv("GX_CLOUD_URL", apiServer.URL)

	githubServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer gho_api_token" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Bad credentials"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"login": "octocat",
			"id":    12345,
		})
	}))
	defer githubServer.Close()
	t.Setenv("GX_GITHUB_USER_URL", githubServer.URL)

	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_api_token",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	status, err := buildAuthStatusJSON(context.Background())
	if err != nil {
		t.Fatalf("buildAuthStatusJSON() error = %v", err)
	}
	if status.AuthKind != "github" || !status.LoggedIn {
		t.Fatalf("status auth kind/logged in = %q/%v, want github/true", status.AuthKind, status.LoggedIn)
	}
	if status.APIValid == nil || !*status.APIValid {
		t.Fatalf("APIValid = %v, APIError = %q, want valid", status.APIValid, status.APIError)
	}
	if gotAPIAuth != "Bearer gho_api_token" {
		t.Fatalf("API Authorization = %q, want resolved bearer token", gotAPIAuth)
	}
}
