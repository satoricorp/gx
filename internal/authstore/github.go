package authstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/satoricorp/totality/internal/buildconfig"
	"github.com/satoricorp/totality/internal/storage"
)

const (
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
	githubKeyringService = "tx"
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
	GitHubKeychainAccount       string    `json:"github_keychain_account,omitempty"`
	CLISessionToken             string    `json:"cli_session_token,omitempty"`
	CLISessionExpiresAt         time.Time `json:"cli_session_expires_at,omitempty"`
	UserID                      string    `json:"user_id,omitempty"`
	Login                       string    `json:"login,omitempty"`
	AvatarURL                   string    `json:"avatar_url,omitempty"`
	MachineID                   string    `json:"machine_id,omitempty"`
	MachineName                 string    `json:"machine_name,omitempty"`
	ObtainedAt                  time.Time `json:"obtained_at,omitempty"`
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

// GitHubToken is the secret GitHub OAuth state stored in the OS keychain.
type GitHubToken struct {
	AccessToken           string    `json:"access_token,omitempty"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at,omitempty"`
	RefreshToken          string    `json:"refresh_token,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty"`
}

// GitHubAccessToken resolves the token shared by desktop auth, CLI auth, and
// GitHub API clients.
func GitHubAccessToken() (string, error) {
	token, _, err := GitHubAccessTokenWithSource()
	return token, err
}

// GitHubAccessTokenWithSource resolves the token and reports where it came from:
// env, keychain, or legacy.
func GitHubAccessTokenWithSource() (string, string, error) {
	for _, name := range []string{"GH_TOKEN", "GITHUB_TOKEN"} {
		if token := strings.TrimSpace(os.Getenv(name)); token != "" {
			return token, "env", nil
		}
	}

	path, file, err := loadCredentials()
	if err != nil {
		return "", "", err
	}
	if file.Cloud == nil {
		return "", "", tokenError()
	}

	if account := githubKeychainAccountForCloud(file.Cloud); account != "" {
		token, err := loadKeychainToken(account)
		if err == nil {
			accessToken, updated, err := validStoredToken(token, time.Now().UTC(), nil)
			if err != nil {
				return "", "", err
			}
			if updated {
				if err := StoreGitHubToken(account, *token); err != nil {
					return "", "", fmt.Errorf("update GitHub token in keychain: %w", err)
				}
			}
			return accessToken, "keychain", nil
		}
		if !errors.Is(err, keyring.ErrNotFound) {
			return "", "", fmt.Errorf("read GitHub token from keychain: %w", err)
		}
	}

	legacy := legacyTokenFromCloud(file.Cloud)
	if strings.TrimSpace(legacy.AccessToken) == "" {
		return "", "", tokenError()
	}
	token, updated, err := validStoredToken(&legacy, time.Now().UTC(), nil)
	if err != nil {
		return "", "", err
	}
	_ = updated
	account := githubKeychainAccountForCloud(file.Cloud)
	if account == "" {
		account = GitHubKeychainAccount(file.Cloud.UserID, file.Cloud.Login)
	}
	if err := StoreGitHubToken(account, legacy); err != nil {
		return "", "", fmt.Errorf("migrate GitHub token to keychain: %w", err)
	}
	file.Cloud.GitHubKeychainAccount = account
	file.Cloud.GitHubAccessToken = ""
	file.Cloud.GitHubAccessTokenExpiresAt = time.Time{}
	file.Cloud.GitHubRefreshToken = ""
	file.Cloud.GitHubRefreshTokenExpiresAt = time.Time{}
	if err := writeCredentials(path, file); err != nil {
		return "", "", err
	}
	return token, "legacy", nil
}

func validStoredToken(creds *GitHubToken, now time.Time, client *http.Client) (string, bool, error) {
	token := strings.TrimSpace(creds.AccessToken)
	if creds.AccessTokenExpiresAt.IsZero() || now.Before(creds.AccessTokenExpiresAt.Add(-tokenRefreshSkew)) {
		return token, false, nil
	}
	refreshToken := strings.TrimSpace(creds.RefreshToken)
	if refreshToken == "" {
		return "", false, reloginError("stored GitHub token expired and no refresh token is available")
	}
	if !creds.RefreshTokenExpiresAt.IsZero() && !now.Before(creds.RefreshTokenExpiresAt) {
		return "", false, reloginError("stored GitHub refresh token expired")
	}
	refreshed, err := refreshGitHubAccessToken(refreshToken, client)
	if err != nil {
		return "", false, reloginError("stored GitHub token refresh failed: " + err.Error())
	}
	if strings.TrimSpace(refreshed.AccessToken) == "" {
		return "", false, reloginError("GitHub refresh response did not include an access token")
	}
	creds.AccessToken = strings.TrimSpace(refreshed.AccessToken)
	if refreshed.ExpiresIn > 0 {
		creds.AccessTokenExpiresAt = now.Add(time.Duration(refreshed.ExpiresIn) * time.Second)
	} else {
		creds.AccessTokenExpiresAt = time.Time{}
	}
	if strings.TrimSpace(refreshed.RefreshToken) != "" {
		creds.RefreshToken = strings.TrimSpace(refreshed.RefreshToken)
	}
	if refreshed.RefreshTokenExpiresIn > 0 {
		creds.RefreshTokenExpiresAt = now.Add(time.Duration(refreshed.RefreshTokenExpiresIn) * time.Second)
	}
	return creds.AccessToken, true, nil
}

// GitHubKeychainAccount returns the account name used for GitHub OAuth state.
func GitHubKeychainAccount(userID, login string) string {
	if id := strings.TrimSpace(userID); id != "" {
		return "github.com:" + id
	}
	if name := strings.TrimSpace(login); name != "" {
		return "github.com:" + name
	}
	return "github.com"
}

// StoreGitHubToken writes GitHub OAuth state to the OS keychain.
func StoreGitHubToken(account string, token GitHubToken) error {
	account = strings.TrimSpace(account)
	if account == "" {
		account = GitHubKeychainAccount("", "")
	}
	if strings.TrimSpace(token.AccessToken) == "" {
		return fmt.Errorf("missing GitHub access token")
	}
	raw, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshal GitHub token secret: %w", err)
	}
	if err := keyring.Set(githubKeyringService, account, string(raw)); err != nil {
		return err
	}
	return nil
}

func loadKeychainToken(account string) (*GitHubToken, error) {
	raw, err := keyring.Get(githubKeyringService, account)
	if err != nil {
		return nil, err
	}
	var token GitHubToken
	if err := json.Unmarshal([]byte(raw), &token); err != nil {
		return nil, fmt.Errorf("parse GitHub token from keychain: %w", err)
	}
	if strings.TrimSpace(token.AccessToken) == "" {
		return nil, fmt.Errorf("GitHub token in keychain is missing access token")
	}
	return &token, nil
}

func githubKeychainAccountForCloud(creds *cloudCredentials) string {
	if creds == nil {
		return ""
	}
	if account := strings.TrimSpace(creds.GitHubKeychainAccount); account != "" {
		return account
	}
	if strings.TrimSpace(creds.UserID) != "" || strings.TrimSpace(creds.Login) != "" {
		return GitHubKeychainAccount(creds.UserID, creds.Login)
	}
	return ""
}

func legacyTokenFromCloud(creds *cloudCredentials) GitHubToken {
	if creds == nil {
		return GitHubToken{}
	}
	return GitHubToken{
		AccessToken:           strings.TrimSpace(creds.GitHubAccessToken),
		AccessTokenExpiresAt:  creds.GitHubAccessTokenExpiresAt,
		RefreshToken:          strings.TrimSpace(creds.GitHubRefreshToken),
		RefreshTokenExpiresAt: creds.GitHubRefreshTokenExpiresAt,
	}
}

func loadCredentials() (string, credentialsFile, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", credentialsFile{}, err
	}
	path := filepath.Join(dir, "credentials.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return path, credentialsFile{}, nil
		}
		return "", credentialsFile{}, fmt.Errorf("read github credentials: %w", err)
	}
	var file credentialsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return "", credentialsFile{}, fmt.Errorf("parse github credentials: %w", err)
	}
	return path, file, nil
}

func writeCredentials(path string, file credentialsFile) error {
	out, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal credentials.json: %w", err)
	}
	out = append(out, '\n')
	if err := os.WriteFile(path, out, 0o600); err != nil {
		return fmt.Errorf("write credentials.json: %w", err)
	}
	return nil
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
	endpoint := strings.TrimSpace(os.Getenv("TOTALITY_GITHUB_ACCESS_TOKEN_URL"))
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
	if strings.TrimSpace(os.Getenv("TOTALITY_MCP")) != "" {
		return fmt.Errorf("github token is not configured for MCP: run `tx auth login` in a terminal, then retry the MCP tool")
	}
	return fmt.Errorf("github token is not configured: run `tx auth login` or set GH_TOKEN/GITHUB_TOKEN")
}

func reloginError(reason string) error {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "stored GitHub token is invalid"
	}
	if strings.TrimSpace(os.Getenv("TOTALITY_MCP")) != "" {
		return fmt.Errorf("%s; run `tx auth login` in a terminal, then retry the MCP tool", reason)
	}
	return fmt.Errorf("%s; run `tx auth login` to refresh stored credentials or set GH_TOKEN/GITHUB_TOKEN", reason)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
