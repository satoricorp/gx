package cloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateGitHubAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer gho_ok" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"message": "Bad credentials"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"login": "octocat",
			"id":    12345,
		})
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_USER_URL", server.URL)

	valid, err := ValidateGitHubAccessToken(context.Background(), server.Client(), "gho_ok")
	if err != nil {
		t.Fatalf("ValidateGitHubAccessToken() valid error = %v", err)
	}
	if !valid.Valid || valid.Login != "octocat" || valid.UserID != 12345 {
		t.Fatalf("valid result = %+v", valid)
	}

	invalid, err := ValidateGitHubAccessToken(context.Background(), server.Client(), "gho_bad")
	if err != nil {
		t.Fatalf("ValidateGitHubAccessToken() invalid error = %v", err)
	}
	if invalid.Valid || invalid.StatusCode != http.StatusUnauthorized || invalid.Error != "Bad credentials" {
		t.Fatalf("invalid result = %+v", invalid)
	}
}

func TestValidateCloudAPISession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/me" {
			t.Fatalf("path = %s, want /v1/auth/me", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer gxcs_ok" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"user_id":           "user_1",
			"github_user_login": "octocat",
		})
	}))
	defer server.Close()
	t.Setenv("GX_CLOUD_URL", server.URL)

	valid, err := ValidateCloudAPISession(context.Background(), server.Client(), "gxcs_ok")
	if err != nil {
		t.Fatalf("ValidateCloudAPISession() valid error = %v", err)
	}
	if !valid.Valid || valid.UserID != "user_1" || valid.Login != "octocat" {
		t.Fatalf("valid result = %+v", valid)
	}

	invalid, err := ValidateCloudAPISession(context.Background(), server.Client(), "gxcs_bad")
	if err != nil {
		t.Fatalf("ValidateCloudAPISession() invalid error = %v", err)
	}
	if invalid.Valid || invalid.StatusCode != http.StatusUnauthorized || invalid.Error != "Unauthorized" {
		t.Fatalf("invalid result = %+v", invalid)
	}
}
