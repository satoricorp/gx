package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/satoricorp/lgtm/internal/authstore"
	"github.com/satoricorp/lgtm/internal/version"
)

const (
	githubDeviceCodeURL  = "https://github.com/login/device/code"
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
)

// AuthEndpoints override GitHub and Convex URLs for tests.
type AuthEndpoints struct {
	GitHubDeviceCodeURL  string
	GitHubAccessTokenURL string
	ConvexSiteURL        string
	CloudURL             string
}

func (e AuthEndpoints) deviceCodeURL() string {
	if e.GitHubDeviceCodeURL != "" {
		return e.GitHubDeviceCodeURL
	}
	return githubDeviceCodeURL
}

func (e AuthEndpoints) accessTokenURL() string {
	if e.GitHubAccessTokenURL != "" {
		return e.GitHubAccessTokenURL
	}
	return githubAccessTokenURL
}

func (e AuthEndpoints) convexSiteURL() string {
	return strings.TrimRight(strings.TrimSpace(e.ConvexSiteURL), "/")
}

func (e AuthEndpoints) cloudURL() string {
	return strings.TrimRight(strings.TrimSpace(e.CloudURL), "/")
}

type deviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Error           string `json:"error"`
}

type accessTokenResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Error                 string `json:"error"`
}

// GitHubTokenValidation reports whether GitHub accepts a bearer token.
type GitHubTokenValidation struct {
	Valid      bool
	Login      string
	UserID     int64
	StatusCode int
	Error      string
}

// CloudAPISessionValidation reports whether lgtm-cloud accepts a bearer token.
type CloudAPISessionValidation struct {
	Valid      bool
	UserID     string
	Login      string
	StatusCode int
	Error      string
}

// CompleteAuthResponse is returned by POST /cx/auth/complete.
type CompleteAuthResponse struct {
	UserID              string `json:"user_id"`
	Login               string `json:"login"`
	AvatarURL           string `json:"avatar_url,omitempty"`
	GitHubAppInstallURL string `json:"github_app_install_url,omitempty"`
	CLISessionToken     string `json:"cli_session_token,omitempty"`
	CLISessionExpiresAt int64  `json:"cli_session_expires_at,omitempty"`
}

type completeAuthRequest struct {
	GitHubAccessToken string `json:"github_access_token"`
	MachineID         string `json:"machine_id"`
	MachineName       string `json:"machine_name"`
	TLVersion         string `json:"lgtm_version"`
}

// LoginOptions configures lgtm auth login.
type LoginOptions struct {
	MachineName string
	Endpoints   AuthEndpoints
	HTTPClient  *http.Client
	Out         io.Writer
}

// Login runs GitHub device flow, stores the GitHub token with Convex via
// POST /cx/auth/complete, and saves the token locally for lgtm cloud API calls.
func Login(ctx context.Context, opts LoginOptions) (CloudCredentials, error) {
	clientID := GitHubClientID()
	if clientID == "" {
		return CloudCredentials{}, fmt.Errorf("cloud auth not configured (dev build? set GITHUB_CLIENT_ID)")
	}
	convexURL := opts.Endpoints.convexSiteURL()
	if convexURL == "" {
		convexURL = ConvexSiteURL()
	}
	if convexURL == "" {
		return CloudCredentials{}, fmt.Errorf("cloud auth not configured (dev build? set CONVEX_SITE_URL)")
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	device, err := requestGitHubDeviceCode(ctx, httpClient, opts.Endpoints.deviceCodeURL(), clientID)
	if err != nil {
		return CloudCredentials{}, err
	}

	verificationURI := device.VerificationURI
	if verificationURI == "" {
		verificationURI = "https://github.com/login/device"
	}
	if opts.Out != nil {
		fmt.Fprintf(opts.Out, "Visit %s\n", verificationURI)
		fmt.Fprintf(opts.Out, "Enter code: %s\n", device.UserCode)
	}

	githubToken, err := pollGitHubAccessToken(ctx, httpClient, opts.Endpoints.accessTokenURL(), clientID, device)
	if err != nil {
		return CloudCredentials{}, err
	}

	machineID, err := DefaultMachineID()
	if err != nil {
		return CloudCredentials{}, err
	}
	machineName := strings.TrimSpace(opts.MachineName)
	if machineName == "" {
		machineName, _ = os.Hostname()
	}

	complete, err := completeConvexAuth(ctx, httpClient, convexURL, completeAuthRequest{
		GitHubAccessToken: githubToken.AccessToken,
		MachineID:         machineID,
		MachineName:       machineName,
		TLVersion:         version.Current(),
	})
	if err != nil {
		return CloudCredentials{}, err
	}

	now := time.Now().UTC()
	creds := CloudCredentials{
		GitHubAccessToken:  githubToken.AccessToken,
		GitHubRefreshToken: githubToken.RefreshToken,
		CLISessionToken:    complete.CLISessionToken,
		UserID:             complete.UserID,
		Login:              complete.Login,
		AvatarURL:          complete.AvatarURL,
		MachineID:          machineID,
		MachineName:        machineName,
		ObtainedAt:         now,
	}
	if githubToken.ExpiresIn > 0 {
		creds.GitHubAccessTokenExpiresAt = now.Add(time.Duration(githubToken.ExpiresIn) * time.Second)
	}
	if githubToken.RefreshTokenExpiresIn > 0 {
		creds.GitHubRefreshTokenExpiresAt = now.Add(time.Duration(githubToken.RefreshTokenExpiresIn) * time.Second)
	}
	if complete.CLISessionExpiresAt > 0 {
		creds.CLISessionExpiresAt = time.UnixMilli(complete.CLISessionExpiresAt).UTC()
	}
	creds.GitHubKeychainAccount = authstore.GitHubKeychainAccount(complete.UserID, complete.Login)
	if strings.TrimSpace(creds.CLISessionToken) != "" {
		cloudURL := opts.Endpoints.cloudURL()
		if cloudURL == "" {
			cloudURL = CloudURL()
		}
		if cloudURL != "" {
			validation, err := validateCloudAPISessionAtURL(ctx, httpClient, cloudURL, creds.CLISessionToken)
			if err != nil {
				return CloudCredentials{}, err
			}
			if !validation.Valid {
				detail := strings.TrimSpace(validation.Error)
				if detail == "" {
					detail = "console API rejected CLI session"
				}
				return CloudCredentials{}, fmt.Errorf("verify lgtm console session: %s", detail)
			}
		}
	}
	if err := authstore.StoreGitHubToken(creds.GitHubKeychainAccount, authstore.GitHubToken{
		AccessToken:           githubToken.AccessToken,
		AccessTokenExpiresAt:  creds.GitHubAccessTokenExpiresAt,
		RefreshToken:          githubToken.RefreshToken,
		RefreshTokenExpiresAt: creds.GitHubRefreshTokenExpiresAt,
	}); err != nil {
		return CloudCredentials{}, fmt.Errorf("store GitHub token in keychain: %w", err)
	}
	creds.GitHubAccessToken = ""
	creds.GitHubAccessTokenExpiresAt = time.Time{}
	creds.GitHubRefreshToken = ""
	creds.GitHubRefreshTokenExpiresAt = time.Time{}
	if err := SaveCloudCredentials(creds); err != nil {
		return CloudCredentials{}, err
	}
	if opts.Out != nil && strings.TrimSpace(complete.GitHubAppInstallURL) != "" {
		fmt.Fprintf(opts.Out, "Install the lgtm GitHub App: %s\n", strings.TrimSpace(complete.GitHubAppInstallURL))
	}
	return creds, nil
}

// Logout clears local credentials.
func Logout(ctx context.Context, endpoints AuthEndpoints, httpClient *http.Client) error {
	_ = ctx
	_ = endpoints
	_ = httpClient
	return ClearCloudCredentials()
}

func requestGitHubDeviceCode(ctx context.Context, client *http.Client, deviceURL, clientID string) (deviceCodeResponse, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("scope", "repo")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, deviceURL, strings.NewReader(form.Encode()))
	if err != nil {
		return deviceCodeResponse{}, fmt.Errorf("create github device code request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return deviceCodeResponse{}, fmt.Errorf("request github device code: %w", err)
	}
	defer resp.Body.Close()

	var device deviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return deviceCodeResponse{}, fmt.Errorf("decode github device code response: %w", err)
	}
	if device.Error != "" {
		return deviceCodeResponse{}, fmt.Errorf("github device code: %s", device.Error)
	}
	if device.DeviceCode == "" || device.UserCode == "" {
		return deviceCodeResponse{}, fmt.Errorf("github device code response missing device_code or user_code")
	}
	return device, nil
}

func pollGitHubAccessToken(ctx context.Context, client *http.Client, tokenURL, clientID string, device deviceCodeResponse) (accessTokenResponse, error) {
	interval := device.Interval
	if interval <= 0 {
		interval = 5
	}
	deadline := time.Now().Add(time.Duration(device.ExpiresIn) * time.Second)
	if device.ExpiresIn <= 0 {
		deadline = time.Now().Add(15 * time.Minute)
	}

	for {
		if time.Now().After(deadline) {
			return accessTokenResponse{}, fmt.Errorf("github device authorization timed out")
		}

		form := url.Values{}
		form.Set("client_id", clientID)
		form.Set("device_code", device.DeviceCode)
		form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
		if err != nil {
			return accessTokenResponse{}, fmt.Errorf("create github access token request: %w", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return accessTokenResponse{}, fmt.Errorf("poll github access token: %w", err)
		}

		var tokenResp accessTokenResponse
		decodeErr := json.NewDecoder(resp.Body).Decode(&tokenResp)
		resp.Body.Close()
		if decodeErr != nil {
			return accessTokenResponse{}, fmt.Errorf("decode github access token response: %w", decodeErr)
		}

		if tokenResp.AccessToken != "" {
			return tokenResp, nil
		}

		switch tokenResp.Error {
		case "authorization_pending":
			// continue polling
		case "slow_down":
			interval++
		case "incorrect_device_code", "expired_token":
			return accessTokenResponse{}, fmt.Errorf("github device authorization was rejected or expired; run `lgtm auth login` again and use the newest code")
		case "":
			return accessTokenResponse{}, fmt.Errorf("github access token response missing access_token")
		default:
			return accessTokenResponse{}, fmt.Errorf("github access token: %s", tokenResp.Error)
		}

		select {
		case <-ctx.Done():
			return accessTokenResponse{}, ctx.Err()
		case <-time.After(time.Duration(interval) * time.Second):
		}
	}
}

// ValidateGitHubAccessToken checks the token against GitHub's user endpoint.
func ValidateGitHubAccessToken(ctx context.Context, client *http.Client, token string) (GitHubTokenValidation, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return GitHubTokenValidation{Valid: false, Error: "missing token"}, nil
	}
	userURL := strings.TrimSpace(os.Getenv("LGTM_GITHUB_USER_URL"))
	if userURL == "" {
		userURL = "https://api.github.com/user"
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userURL, nil)
	if err != nil {
		return GitHubTokenValidation{}, fmt.Errorf("create github token validation request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "lgtm/"+version.Current())
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return GitHubTokenValidation{}, fmt.Errorf("validate github token: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		Login   string `json:"login"`
		ID      int64  `json:"id"`
		Message string `json:"message"`
	}
	_ = json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		message := strings.TrimSpace(body.Message)
		if message == "" {
			message = resp.Status
		}
		return GitHubTokenValidation{
			Valid:      false,
			StatusCode: resp.StatusCode,
			Error:      message,
		}, nil
	}
	if body.ID == 0 {
		return GitHubTokenValidation{
			Valid:      false,
			StatusCode: resp.StatusCode,
			Error:      "GitHub response missing user id",
		}, nil
	}
	return GitHubTokenValidation{
		Valid:      true,
		Login:      body.Login,
		UserID:     body.ID,
		StatusCode: resp.StatusCode,
	}, nil
}

// ValidateCloudAPISession checks the token against lgtm-cloud's auth endpoint.
func ValidateCloudAPISession(ctx context.Context, client *http.Client, token string) (CloudAPISessionValidation, error) {
	return validateCloudAPISessionAtURL(ctx, client, CloudURL(), token)
}

func validateCloudAPISessionAtURL(ctx context.Context, client *http.Client, cloudURL string, token string) (CloudAPISessionValidation, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return CloudAPISessionValidation{Valid: false, Error: "missing token"}, nil
	}
	meURL := cloudURLWithPath(cloudURL, "/v1/auth/me")
	if meURL == "" {
		return CloudAPISessionValidation{Valid: false, Error: "lgtm cloud URL is not configured"}, nil
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, meURL, nil)
	if err != nil {
		return CloudAPISessionValidation{}, fmt.Errorf("create lgtm api auth validation request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "lgtm/"+version.Current())

	resp, err := client.Do(req)
	if err != nil {
		return CloudAPISessionValidation{}, fmt.Errorf("validate lgtm api token: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		UserID          string `json:"user_id"`
		GitHubUserLogin string `json:"github_user_login"`
		Error           string `json:"error"`
	}
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	_ = json.Unmarshal(raw, &body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		message := strings.TrimSpace(body.Error)
		if message == "" {
			message = strings.TrimSpace(string(raw))
		}
		if message == "" {
			message = resp.Status
		}
		return CloudAPISessionValidation{
			Valid:      false,
			StatusCode: resp.StatusCode,
			Error:      message,
		}, nil
	}
	if strings.TrimSpace(body.UserID) == "" {
		return CloudAPISessionValidation{
			Valid:      false,
			StatusCode: resp.StatusCode,
			Error:      "lgtm API response missing user id",
		}, nil
	}
	return CloudAPISessionValidation{
		Valid:      true,
		UserID:     body.UserID,
		Login:      body.GitHubUserLogin,
		StatusCode: resp.StatusCode,
	}, nil
}

func completeConvexAuth(ctx context.Context, client *http.Client, convexURL string, body completeAuthRequest) (CompleteAuthResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return CompleteAuthResponse{}, fmt.Errorf("marshal auth complete payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, convexURL+"/cx/auth/complete", bytes.NewReader(payload))
	if err != nil {
		return CompleteAuthResponse{}, fmt.Errorf("create auth complete request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "lgtm/"+version.Current())

	resp, err := client.Do(req)
	if err != nil {
		return CompleteAuthResponse{}, fmt.Errorf("auth complete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return CompleteAuthResponse{}, fmt.Errorf("auth complete: status %s: %s", resp.Status, strings.TrimSpace(string(msg)))
	}

	var complete CompleteAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&complete); err != nil {
		return CompleteAuthResponse{}, fmt.Errorf("decode auth complete response: %w", err)
	}
	return complete, nil
}
