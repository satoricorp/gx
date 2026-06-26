package authstore

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/buildconfig"
	"github.com/satoricorp/gx/internal/storage"
)

const (
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
	tokenRefreshSkew     = 5 * time.Minute
)

type credentialsFile struct {
	Cloud *cloudCredentials `json:"cloud"`
}

type cloudCredentials struct {
	GitHubAccessToken           string    `json:"github_access_token,omitempty"`
	GitHubAccessTokenExpiresAt  time.Time `json:"github_access_token_expires_at,omitempty"`
	GitHubRefreshToken          string    `json:"github_refresh_token,omitempty"`
	GitHubRefreshTokenExpiresAt time.Time `json:"github_refresh_token_expires_at,omitempty"`
}

type accessTokenResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
	Message               string `json:"message"`
}

// GitHubAccessToken resolves the token shared by desktop auth, CLI auth, and
// GitHub API clients.
func GitHubAccessToken() (string, error) {
	for _, name := range []string{"GH_TOKEN", "GITHUB_TOKEN"} {
		if token := strings.TrimSpace(os.Getenv(name)); token != "" {
			return token, nil
		}
	}
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "credentials.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", tokenError()
		}
		return "", fmt.Errorf("read github credentials: %w", err)
	}
	var file credentialsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return "", fmt.Errorf("parse github credentials: %w", err)
	}
	if file.Cloud == nil || strings.TrimSpace(file.Cloud.GitHubAccessToken) == "" {
		return "", tokenError()
	}
	token, updated, err := validStoredToken(file.Cloud, time.Now().UTC(), nil)
	if err != nil {
		return "", err
	}
	if updated {
		out, err := json.MarshalIndent(file, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal credentials.json: %w", err)
		}
		out = append(out, '\n')
		if err := os.WriteFile(path, out, 0o600); err != nil {
			return "", fmt.Errorf("write credentials.json: %w", err)
		}
	}
	return token, nil
}

func validStoredToken(creds *cloudCredentials, now time.Time, client *http.Client) (string, bool, error) {
	token := strings.TrimSpace(creds.GitHubAccessToken)
	if creds.GitHubAccessTokenExpiresAt.IsZero() || now.Before(creds.GitHubAccessTokenExpiresAt.Add(-tokenRefreshSkew)) {
		return token, false, nil
	}
	refreshToken := strings.TrimSpace(creds.GitHubRefreshToken)
	if refreshToken == "" {
		return "", false, reloginError("stored GitHub token expired and no refresh token is available")
	}
	if !creds.GitHubRefreshTokenExpiresAt.IsZero() && !now.Before(creds.GitHubRefreshTokenExpiresAt) {
		return "", false, reloginError("stored GitHub refresh token expired")
	}
	refreshed, err := refreshGitHubAccessToken(refreshToken, client)
	if err != nil {
		return "", false, reloginError("stored GitHub token refresh failed: " + err.Error())
	}
	if strings.TrimSpace(refreshed.AccessToken) == "" {
		return "", false, reloginError("GitHub refresh response did not include an access token")
	}
	creds.GitHubAccessToken = strings.TrimSpace(refreshed.AccessToken)
	if refreshed.ExpiresIn > 0 {
		creds.GitHubAccessTokenExpiresAt = now.Add(time.Duration(refreshed.ExpiresIn) * time.Second)
	} else {
		creds.GitHubAccessTokenExpiresAt = time.Time{}
	}
	if strings.TrimSpace(refreshed.RefreshToken) != "" {
		creds.GitHubRefreshToken = strings.TrimSpace(refreshed.RefreshToken)
	}
	if refreshed.RefreshTokenExpiresIn > 0 {
		creds.GitHubRefreshTokenExpiresAt = now.Add(time.Duration(refreshed.RefreshTokenExpiresIn) * time.Second)
	}
	return creds.GitHubAccessToken, true, nil
}

func refreshGitHubAccessToken(refreshToken string, client *http.Client) (accessTokenResponse, error) {
	clientID := strings.TrimSpace(buildconfig.GitHubClientIDOrEnv())
	if clientID == "" {
		return accessTokenResponse{}, fmt.Errorf("missing GITHUB_CLIENT_ID")
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	endpoint := strings.TrimSpace(os.Getenv("GX_GITHUB_ACCESS_TOKEN_URL"))
	if endpoint == "" {
		endpoint = githubAccessTokenURL
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("create github refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return accessTokenResponse{}, fmt.Errorf("request github refresh token: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	var tokenResp accessTokenResponse
	if err := json.Unmarshal(raw, &tokenResp); err != nil {
		return accessTokenResponse{}, fmt.Errorf("decode github refresh response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		detail := firstNonEmpty(tokenResp.ErrorDescription, tokenResp.Error, tokenResp.Message, strings.TrimSpace(string(raw)), resp.Status)
		return accessTokenResponse{}, fmt.Errorf("%s", detail)
	}
	if strings.TrimSpace(tokenResp.Error) != "" {
		detail := firstNonEmpty(tokenResp.ErrorDescription, tokenResp.Error)
		return accessTokenResponse{}, fmt.Errorf("%s", detail)
	}
	return tokenResp, nil
}

func tokenError() error {
	if strings.TrimSpace(os.Getenv("GX_MCP")) != "" {
		return fmt.Errorf("github token is not configured for MCP: run `gx auth login` in a terminal, then retry the MCP tool")
	}
	return fmt.Errorf("github token is not configured: run `gx auth login` or set GH_TOKEN/GITHUB_TOKEN")
}

func reloginError(reason string) error {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "stored GitHub token is invalid"
	}
	if strings.TrimSpace(os.Getenv("GX_MCP")) != "" {
		return fmt.Errorf("%s; run `gx auth login` in a terminal, then retry the MCP tool", reason)
	}
	return fmt.Errorf("%s; run `gx auth login` to refresh stored credentials or set GH_TOKEN/GITHUB_TOKEN", reason)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
