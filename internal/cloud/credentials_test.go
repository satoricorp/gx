package cloud

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultMachineIDCreatesOnce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	id1, err := DefaultMachineID()
	if err != nil {
		t.Fatalf("DefaultMachineID() error = %v", err)
	}
	if id1 == "" {
		t.Fatal("expected non-empty machine id")
	}

	info, err := os.Stat(filepath.Join(home, "machine_id.json"))
	if err != nil {
		t.Fatalf("Stat(machine_id.json) error = %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("machine_id.json mode = %o, want 0600", info.Mode().Perm())
	}

	id2, err := DefaultMachineID()
	if err != nil {
		t.Fatalf("DefaultMachineID() second call error = %v", err)
	}
	if id1 != id2 {
		t.Fatalf("machine id changed: %q -> %q", id1, id2)
	}
}

func TestSaveLoadClearCloudCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	if _, err := DefaultMachineID(); err != nil {
		t.Fatalf("DefaultMachineID() error = %v", err)
	}

	obtained := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	want := CloudCredentials{
		Token:             "gx_test_token",
		APIKey:            "gx_api_key",
		GitHubAccessToken: "gho_test_token",
		UserID:            "user_1",
		Login:             "joe",
		AvatarURL:         "https://avatars.githubusercontent.com/u/1?v=4",
		SessionID:         "sess_1",
		MachineID:         "machine-1",
		MachineName:       "test-host",
		ObtainedAt:        obtained,
	}
	if err := SaveCloudCredentials(want); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	info, err := os.Stat(filepath.Join(home, "credentials.json"))
	if err != nil {
		t.Fatalf("Stat(credentials.json) error = %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("credentials.json mode = %o, want 0600", info.Mode().Perm())
	}

	got, err := LoadCloudCredentials()
	if err != nil {
		t.Fatalf("LoadCloudCredentials() error = %v", err)
	}
	if got == nil {
		t.Fatal("expected credentials")
	}
	if got.Token != want.Token || got.APIKey != want.APIKey || got.GitHubAccessToken != want.GitHubAccessToken || got.Login != want.Login || got.AvatarURL != want.AvatarURL || got.MachineID != want.MachineID {
		t.Fatalf("LoadCloudCredentials() = %+v, want %+v", got, want)
	}

	if err := ClearCloudCredentials(); err != nil {
		t.Fatalf("ClearCloudCredentials() error = %v", err)
	}
	cleared, err := LoadCloudCredentials()
	if err != nil {
		t.Fatalf("LoadCloudCredentials() after clear error = %v", err)
	}
	if cleared != nil {
		t.Fatalf("expected nil credentials after clear, got %+v", cleared)
	}

	machineData, err := os.ReadFile(filepath.Join(home, "machine_id.json"))
	if err != nil {
		t.Fatalf("machine_id.json should remain: %v", err)
	}
	var machine machineIDFile
	if err := json.Unmarshal(machineData, &machine); err != nil {
		t.Fatalf("parse machine_id.json: %v", err)
	}
	if machine.ID == "" {
		t.Fatal("machine_id.json should not be cleared")
	}
}

func TestGitHubAccessTokenResolution(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if _, err := GitHubAccessToken(); err == nil {
		t.Fatal("expected error when github token is missing")
	}

	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_abcdefghijklmnopqrstuvwxyz", GitHubAccessToken: "stored-token"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, err := GitHubAccessToken()
	if err != nil {
		t.Fatalf("GitHubAccessToken() stored error = %v", err)
	}
	if token != "stored-token" {
		t.Fatalf("GitHubAccessToken() = %q, want stored-token", token)
	}

	t.Setenv("GH_TOKEN", "env-token")
	token, err = GitHubAccessToken()
	if err != nil {
		t.Fatalf("GitHubAccessToken() env error = %v", err)
	}
	if token != "env-token" {
		t.Fatalf("GitHubAccessToken() = %q, want env-token", token)
	}
}

func TestBearerTokenResolution(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	if _, err := BearerToken(); err == nil {
		t.Fatal("expected error when not logged in")
	}

	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, err := BearerToken()
	if err != nil {
		t.Fatalf("BearerToken() from file error = %v", err)
	}
	if token != "gx_abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("BearerToken() = %q, want stored session token", token)
	}
}

func TestCloudAPITokenPrefersGitHubToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if err := SaveCloudCredentials(CloudCredentials{
		Token:             "gx_abcdefghijklmnopqrstuvwxyz",
		GitHubAccessToken: "gho_test_token",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, err := CloudAPIToken()
	if err != nil {
		t.Fatalf("CloudAPIToken() error = %v", err)
	}
	if token != "gho_test_token" {
		t.Fatalf("CloudAPIToken() = %q, want GitHub token", token)
	}
}

func TestSaveAPIKeyCredentialsAndCloudAPIToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_API_KEY", "")

	creds, err := SaveAPIKeyCredentials(" gx_api_secret ", "test-machine")
	if err != nil {
		t.Fatalf("SaveAPIKeyCredentials() error = %v", err)
	}
	if creds.APIKey != "gx_api_secret" || creds.Login != "api-key" || creds.MachineName != "test-machine" {
		t.Fatalf("saved api key credentials = %+v", creds)
	}

	token, kind, err := CloudAPITokenWithKind()
	if err != nil {
		t.Fatalf("CloudAPITokenWithKind() error = %v", err)
	}
	if token != "gx_api_secret" || kind != "api_key" {
		t.Fatalf("CloudAPITokenWithKind() = (%q, %q), want api key", token, kind)
	}
}

func TestCloudAPITokenPrefersGXAPIKey(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "github-token")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_API_KEY", "api-key")

	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, err := CloudAPIToken()
	if err != nil {
		t.Fatalf("CloudAPIToken() error = %v", err)
	}
	if token != "api-key" {
		t.Fatalf("CloudAPIToken() = %q, want GX API key", token)
	}
}

func TestCloudAPITokenPrefersStoredAPIKeyOverGitHubToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "github-token")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_API_KEY", "")

	if err := SaveCloudCredentials(CloudCredentials{
		APIKey:            "stored-api-key",
		GitHubAccessToken: "stored-github-token",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, kind, err := CloudAPITokenWithKind()
	if err != nil {
		t.Fatalf("CloudAPITokenWithKind() error = %v", err)
	}
	if token != "stored-api-key" || kind != "api_key" {
		t.Fatalf("CloudAPITokenWithKind() = (%q, %q), want stored api key", token, kind)
	}
}

func TestCloudAPITokenFallsBackToLegacyBearerToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, err := CloudAPIToken()
	if err != nil {
		t.Fatalf("CloudAPIToken() error = %v", err)
	}
	if token != "gx_abcdefghijklmnopqrstuvwxyz" {
		t.Fatalf("CloudAPIToken() = %q, want legacy bearer token", token)
	}
}

func TestBearerTokenRejectsInvalidCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	if err := SaveCloudCredentials(CloudCredentials{Token: "gx_baked"}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	if _, err := BearerToken(); err == nil {
		t.Fatal("expected invalid stored token error")
	}
}
