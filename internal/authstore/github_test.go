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

	"github.com/zalando/go-keyring"
)

func TestGitHubAccessTokenRefreshesExpiredKeychainToken(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("TOTALITY_HOME", home)
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
	t.Setenv("TOTALITY_GITHUB_ACCESS_TOKEN_URL", server.URL)

	account := GitHubKeychainAccount("user_1", "joe")
	writeTestCredentials(t, home, cloudCredentials{
		UserID:                "user_1",
		Login:                 "joe",
		GitHubKeychainAccount: account,
	})
	if err := StoreGitHubToken(account, GitHubToken{
		AccessToken:           "ghu_old",
		AccessTokenExpiresAt:  time.Now().Add(-time.Minute),
		RefreshToken:          "ghr_old",
		RefreshTokenExpiresAt: time.Now().Add(time.Hour),
	}); err != nil {
		t.Fatalf("StoreGitHubToken() error = %v", err)
	}

	token, err := GitHubAccessToken()
	if err != nil {
		t.Fatalf("GitHubAccessToken() error = %v", err)
	}
	if token != "ghu_new" || gotRefreshToken != "ghr_old" {
		t.Fatalf("token = %q refresh=%q, want refreshed", token, gotRefreshToken)
	}

	stored, err := loadKeychainToken(account)
	if err != nil {
		t.Fatalf("loadKeychainToken() error = %v", err)
	}
	if stored.AccessToken != "ghu_new" || stored.RefreshToken != "ghr_new" {
		t.Fatalf("keychain token = %+v, want refreshed tokens", stored)
	}
	if stored.AccessTokenExpiresAt.IsZero() || stored.RefreshTokenExpiresAt.IsZero() {
		t.Fatalf("keychain token missing expiry metadata: %+v", stored)
	}
}

func TestGitHubAccessTokenExpiredStoredTokenPromptsLoginWithoutRefresh(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("TOTALITY_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	writeTestCredentials(t, home, cloudCredentials{
		GitHubAccessToken:          "ghu_old",
		GitHubAccessTokenExpiresAt: time.Now().Add(-time.Minute),
	})

	_, err := GitHubAccessToken()
	if err == nil {
		t.Fatal("expected expired token error")
	}
	message := err.Error()
	if !strings.Contains(message, "stored GitHub token expired") || !strings.Contains(message, "tl auth login") {
		t.Fatalf("error = %q, want re-login instruction", message)
	}
}

func TestGitHubAccessTokenEnvOverridesExpiredStoredToken(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("TOTALITY_HOME", home)
	t.Setenv("GH_TOKEN", "env-token")
	t.Setenv("GITHUB_TOKEN", "")
	writeTestCredentials(t, home, cloudCredentials{
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

func TestGitHubAccessTokenMigratesLegacyTokenToKeychain(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("TOTALITY_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	writeTestCredentials(t, home, cloudCredentials{
		GitHubAccessToken: "ghu_legacy",
		UserID:            "user_1",
		Login:             "joe",
		CLISessionToken:   "tlcs_test",
	})

	token, source, err := GitHubAccessTokenWithSource()
	if err != nil {
		t.Fatalf("GitHubAccessTokenWithSource() error = %v", err)
	}
	if token != "ghu_legacy" || source != "legacy" {
		t.Fatalf("GitHubAccessTokenWithSource() = (%q, %q), want legacy token", token, source)
	}

	creds := readCredentials(t, home)
	if creds.GitHubAccessToken != "" || creds.GitHubRefreshToken != "" {
		t.Fatalf("legacy GitHub secrets were not scrubbed: %+v", creds)
	}
	if creds.GitHubKeychainAccount != GitHubKeychainAccount("user_1", "joe") || creds.CLISessionToken != "tlcs_test" {
		t.Fatalf("credentials metadata = %+v", creds)
	}
	stored, err := loadKeychainToken(creds.GitHubKeychainAccount)
	if err != nil {
		t.Fatalf("loadKeychainToken() error = %v", err)
	}
	if stored.AccessToken != "ghu_legacy" {
		t.Fatalf("keychain token = %+v, want legacy token", stored)
	}
}

func writeTestCredentials(t *testing.T, home string, creds cloudCredentials) {
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
