package authstore

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGitHubAccessTokenRefreshesExpiredStoredToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GITHUB_CLIENT_ID", "client-id")

	var gotRefreshToken string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		if r.PostForm.Get("client_id") != "client-id" || r.PostForm.Get("grant_type") != "refresh_token" {
			t.Fatalf("form = %v", r.PostForm)
		}
		gotRefreshToken = r.PostForm.Get("refresh_token")
		_ = json.NewEncoder(w).Encode(accessTokenResponse{
			AccessToken:           "ghu_new",
			ExpiresIn:             8 * 60 * 60,
			RefreshToken:          "ghr_new",
			RefreshTokenExpiresIn: 6 * 30 * 24 * 60 * 60,
		})
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_ACCESS_TOKEN_URL", server.URL)

	writeCredentials(t, home, cloudCredentials{
		GitHubAccessToken:           "ghu_old",
		GitHubAccessTokenExpiresAt:  time.Now().Add(-time.Minute),
		GitHubRefreshToken:          "ghr_old",
		GitHubRefreshTokenExpiresAt: time.Now().Add(time.Hour),
	})

	token, err := GitHubAccessToken()
	if err != nil {
		t.Fatalf("GitHubAccessToken() error = %v", err)
	}
	if token != "ghu_new" || gotRefreshToken != "ghr_old" {
		t.Fatalf("token = %q refresh=%q, want refreshed", token, gotRefreshToken)
	}

	creds := readCredentials(t, home)
	if creds.GitHubAccessToken != "ghu_new" || creds.GitHubRefreshToken != "ghr_new" {
		t.Fatalf("stored creds = %+v, want refreshed tokens", creds)
	}
	if creds.GitHubAccessTokenExpiresAt.IsZero() || creds.GitHubRefreshTokenExpiresAt.IsZero() {
		t.Fatalf("stored creds missing expiry metadata: %+v", creds)
	}
}

func TestGitHubAccessTokenExpiredStoredTokenPromptsLoginWithoutRefresh(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	writeCredentials(t, home, cloudCredentials{
		GitHubAccessToken:          "ghu_old",
		GitHubAccessTokenExpiresAt: time.Now().Add(-time.Minute),
	})

	_, err := GitHubAccessToken()
	if err == nil {
		t.Fatal("expected expired token error")
	}
	message := err.Error()
	if !strings.Contains(message, "stored GitHub token expired") || !strings.Contains(message, "gx auth login") {
		t.Fatalf("error = %q, want re-login instruction", message)
	}
}

func TestGitHubAccessTokenEnvOverridesExpiredStoredToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "env-token")
	t.Setenv("GITHUB_TOKEN", "")
	writeCredentials(t, home, cloudCredentials{
		GitHubAccessToken:          "ghu_old",
		GitHubAccessTokenExpiresAt: time.Now().Add(-time.Minute),
	})

	token, err := GitHubAccessToken()
	if err != nil {
		t.Fatalf("GitHubAccessToken() error = %v", err)
	}
	if token != "env-token" {
		t.Fatalf("GitHubAccessToken() = %q, want env-token", token)
	}
}

func writeCredentials(t *testing.T, home string, creds cloudCredentials) {
	t.Helper()
	path := filepath.Join(home, "credentials.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	file := credentialsFile{Cloud: &creds}
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent: %v", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func readCredentials(t *testing.T, home string) cloudCredentials {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(home, "credentials.json"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var file credentialsFile
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if file.Cloud == nil {
		t.Fatal("missing cloud credentials")
	}
	return *file.Cloud
}
