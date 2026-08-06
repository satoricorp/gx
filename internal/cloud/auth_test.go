package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/satoricorp/gx/internal/buildconfig"
)

func TestLoginDeviceFlowAndComplete(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GITHUB_CLIENT_ID", "test-client")
	t.Setenv("CONVEX_SITE_URL", "")

	var gotComplete completeAuthRequest
	convex := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cx/auth/complete" || r.Method != http.MethodPost {
			t.Fatalf("unexpected convex request: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotComplete); err != nil {
			t.Fatalf("decode complete body: %v", err)
		}
		_ = json.NewEncoder(w).Encode(CompleteAuthResponse{
			UserID:              "user_1",
			Login:               "joe",
			AvatarURL:           "https://avatars.githubusercontent.com/u/1?v=4",
			CLISessionToken:     "gxcs_login",
			CLISessionExpiresAt: time.Now().Add(90 * 24 * time.Hour).UnixMilli(),
		})
	}))
	defer convex.Close()

	console := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/me" || r.Method != http.MethodGet {
			t.Fatalf("unexpected console request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer gxcs_login" {
			t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"user_id":           "user_1",
			"github_user_login": "joe",
		})
	}))
	defer console.Close()

	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login/device/code":
			if r.Method != http.MethodPost {
				t.Fatalf("device code method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(deviceCodeResponse{
				DeviceCode:      "device-code-1",
				UserCode:        "ABCD-1234",
				VerificationURI: "https://github.example/login/device",
				ExpiresIn:       60,
				Interval:        0,
			})
		case "/login/oauth/access_token":
			if r.Method != http.MethodPost {
				t.Fatalf("access token method = %s", r.Method)
			}
			_ = json.NewEncoder(w).Encode(accessTokenResponse{
				AccessToken:           "ghu_test",
				ExpiresIn:             8 * 60 * 60,
				RefreshToken:          "ghr_test",
				RefreshTokenExpiresIn: 6 * 30 * 24 * 60 * 60,
			})
		default:
			t.Fatalf("unexpected github path: %s", r.URL.Path)
		}
	}))
	defer github.Close()

	var out bytes.Buffer
	creds, err := Login(context.Background(), LoginOptions{
		MachineName: "work-laptop",
		Endpoints: AuthEndpoints{
			GitHubDeviceCodeURL:  github.URL + "/login/device/code",
			GitHubAccessTokenURL: github.URL + "/login/oauth/access_token",
			ConvexSiteURL:        convex.URL,
			CloudURL:             console.URL,
		},
		HTTPClient: github.Client(),
		Out:        &out,
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if creds.Login != "joe" || creds.GitHubAccessToken != "" || creds.CLISessionToken != "gxcs_login" || creds.AvatarURL != "https://avatars.githubusercontent.com/u/1?v=4" {
		t.Fatalf("unexpected creds: %+v", creds)
	}
	if gotComplete.GitHubAccessToken != "ghu_test" {
		t.Fatalf("github token = %q", gotComplete.GitHubAccessToken)
	}
	if creds.GitHubKeychainAccount == "" {
		t.Fatalf("missing github keychain account: %+v", creds)
	}
	if gotComplete.MachineName != "work-laptop" {
		t.Fatalf("machine name = %q", gotComplete.MachineName)
	}
	if gotComplete.MachineID == "" {
		t.Fatal("expected machine id in complete request")
	}
	if !strings.Contains(out.String(), "ABCD-1234") {
		t.Fatalf("output missing user code: %q", out.String())
	}

	loaded, err := LoadCloudCredentials()
	if err != nil {
		t.Fatalf("LoadCloudCredentials() error = %v", err)
	}
	if loaded == nil || loaded.GitHubAccessToken != "" || loaded.GitHubRefreshToken != "" || loaded.CLISessionToken != "gxcs_login" || loaded.GitHubKeychainAccount == "" {
		t.Fatalf("saved credentials = %+v", loaded)
	}
	token, source, err := GitHubAccessTokenWithSource()
	if err != nil {
		t.Fatalf("GitHubAccessTokenWithSource() error = %v", err)
	}
	if token != "ghu_test" || source != "keychain" {
		t.Fatalf("GitHubAccessTokenWithSource() = (%q, %q), want keychain token", token, source)
	}
}

func TestLoginRejectsUnverifiedConsoleSession(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GITHUB_CLIENT_ID", "test-client")
	t.Setenv("CONVEX_SITE_URL", "")

	convex := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(CompleteAuthResponse{
			UserID:          "user_1",
			Login:           "joe",
			CLISessionToken: "gxcs_bad",
		})
	}))
	defer convex.Close()

	console := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "wrong environment"})
	}))
	defer console.Close()

	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login/device/code":
			_ = json.NewEncoder(w).Encode(deviceCodeResponse{
				DeviceCode: "device-code-1",
				UserCode:   "ABCD-1234",
				ExpiresIn:  60,
			})
		case "/login/oauth/access_token":
			_ = json.NewEncoder(w).Encode(accessTokenResponse{AccessToken: "ghu_test"})
		default:
			t.Fatalf("unexpected github path: %s", r.URL.Path)
		}
	}))
	defer github.Close()

	_, err := Login(context.Background(), LoginOptions{
		Endpoints: AuthEndpoints{
			GitHubDeviceCodeURL:  github.URL + "/login/device/code",
			GitHubAccessTokenURL: github.URL + "/login/oauth/access_token",
			ConvexSiteURL:        convex.URL,
			CloudURL:             console.URL,
		},
		HTTPClient: github.Client(),
	})
	if err == nil {
		t.Fatal("expected login to reject unverified console session")
	}
	if !strings.Contains(err.Error(), "verify gx console session") || !strings.Contains(err.Error(), "wrong environment") {
		t.Fatalf("error = %v, want console verification failure", err)
	}
	if creds, loadErr := LoadCloudCredentials(); loadErr != nil || creds != nil {
		t.Fatalf("credentials should not be saved after failed console verification: creds=%+v err=%v", creds, loadErr)
	}
}

func TestLoginPollsUntilAuthorized(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GITHUB_CLIENT_ID", "test-client")

	var pollCount int
	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login/device/code":
			_ = json.NewEncoder(w).Encode(deviceCodeResponse{
				DeviceCode: "device-code-2",
				UserCode:   "WXYZ-5678",
				ExpiresIn:  30,
				Interval:   0,
			})
		case "/login/oauth/access_token":
			pollCount++
			if pollCount < 2 {
				_ = json.NewEncoder(w).Encode(accessTokenResponse{Error: "authorization_pending"})
				return
			}
			_ = json.NewEncoder(w).Encode(accessTokenResponse{AccessToken: "gho_late"})
		}
	}))
	defer github.Close()

	convex := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(CompleteAuthResponse{
			Login: "jane", UserID: "u2",
		})
	}))
	defer convex.Close()

	_, err := Login(context.Background(), LoginOptions{
		Endpoints: AuthEndpoints{
			GitHubDeviceCodeURL:  github.URL + "/login/device/code",
			GitHubAccessTokenURL: github.URL + "/login/oauth/access_token",
			ConvexSiteURL:        convex.URL,
		},
		HTTPClient: github.Client(),
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if pollCount < 2 {
		t.Fatalf("poll count = %d, want >= 2", pollCount)
	}
}

func TestLoginUsesBakedDefaults(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GITHUB_CLIENT_ID", "")
	t.Setenv("CONVEX_SITE_URL", "")
	buildconfig.GitHubClientID = "baked-client"
	buildconfig.ConvexSiteURL = "https://baked.convex.site"
	t.Cleanup(func() {
		buildconfig.GitHubClientID = ""
		buildconfig.ConvexSiteURL = ""
	})

	var gotClientID string
	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login/device/code" {
			_ = r.ParseForm()
			gotClientID = r.PostForm.Get("client_id")
			_ = json.NewEncoder(w).Encode(deviceCodeResponse{
				DeviceCode: "d1",
				UserCode:   "CODE-1",
				ExpiresIn:  30,
			})
			return
		}
		_ = json.NewEncoder(w).Encode(accessTokenResponse{AccessToken: "gho_baked"})
	}))
	defer github.Close()

	convex := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(CompleteAuthResponse{
			Login: "baker", UserID: "u",
		})
	}))
	defer convex.Close()

	_, err := Login(context.Background(), LoginOptions{
		Endpoints: AuthEndpoints{
			GitHubDeviceCodeURL:  github.URL + "/login/device/code",
			GitHubAccessTokenURL: github.URL + "/login/oauth/access_token",
			ConvexSiteURL:        convex.URL,
		},
		HTTPClient: github.Client(),
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if gotClientID != "baked-client" {
		t.Fatalf("client_id = %q, want baked-client", gotClientID)
	}
}

func TestLoginNotConfiguredInDevBuild(t *testing.T) {
	t.Setenv("GITHUB_CLIENT_ID", "")
	t.Setenv("CONVEX_SITE_URL", "")
	buildconfig.GitHubClientID = ""
	buildconfig.ConvexSiteURL = ""
	t.Cleanup(func() {
		buildconfig.GitHubClientID = ""
		buildconfig.ConvexSiteURL = ""
	})

	_, err := Login(context.Background(), LoginOptions{})
	if err == nil {
		t.Fatal("expected error when cloud auth is not configured")
	}
	if !strings.Contains(err.Error(), "cloud auth not configured") {
		t.Fatalf("error = %v, want configuration hint", err)
	}
}

func TestLogoutRevokesAndClears(t *testing.T) {
	keyring.MockInit()
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("CONVEX_SITE_URL", "")

	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken: "gho_test",
		Login:             "joe",
		ObtainedAt:        time.Now().UTC(),
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	convex := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("logout should not call server, got %s %s", r.Method, r.URL.Path)
	}))
	defer convex.Close()

	if err := Logout(context.Background(), AuthEndpoints{ConvexSiteURL: convex.URL}, convex.Client()); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if creds, _ := LoadCloudCredentials(); creds != nil {
		t.Fatalf("expected cleared credentials, got %+v", creds)
	}
}

func TestRequestGitHubDeviceCodeForm(t *testing.T) {
	var gotForm url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		gotForm = r.PostForm
		_ = json.NewEncoder(w).Encode(deviceCodeResponse{
			DeviceCode: "d",
			UserCode:   "u",
			ExpiresIn:  10,
		})
	}))
	defer server.Close()

	_, err := requestGitHubDeviceCode(context.Background(), server.Client(), server.URL, "cid")
	if err != nil {
		t.Fatalf("requestGitHubDeviceCode() error = %v", err)
	}
	if gotForm.Get("client_id") != "cid" || gotForm.Get("scope") != "repo" {
		t.Fatalf("form = %v", gotForm)
	}
}
