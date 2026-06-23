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
