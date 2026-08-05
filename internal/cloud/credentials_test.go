package cloud

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultMachineIDCreatesOnce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)

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
	t.Setenv("LGTM_HOME", home)

	if _, err := DefaultMachineID(); err != nil {
		t.Fatalf("DefaultMachineID() error = %v", err)
	}

	obtained := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	want := CloudCredentials{
		GitHubAccessToken:   "gho_test_token",
		CLISessionToken:     "tlcs_test_token",
		CLISessionExpiresAt: obtained.Add(90 * 24 * time.Hour),
		UserID:              "user_1",
		Login:               "joe",
		AvatarURL:           "https://avatars.githubusercontent.com/u/1?v=4",
		MachineID:           "machine-1",
		MachineName:         "test-host",
		ObtainedAt:          obtained,
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
	if got.GitHubAccessToken != want.GitHubAccessToken || got.CLISessionToken != want.CLISessionToken || !got.CLISessionExpiresAt.Equal(want.CLISessionExpiresAt) || got.Login != want.Login || got.AvatarURL != want.AvatarURL || got.MachineID != want.MachineID {
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
	t.Setenv("LGTM_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if _, err := GitHubAccessToken(); err == nil {
		t.Fatal("expected error when github token is missing")
	}

	if err := SaveCloudCredentials(CloudCredentials{GitHubAccessToken: "stored-token"}); err != nil {
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

func TestCloudAPITokenPrefersCLISessionToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken: "gho_test_token",
		CLISessionToken:   "tlcs_test_token",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, err := CloudAPIToken()
	if err != nil {
		t.Fatalf("CloudAPIToken() error = %v", err)
	}
	if token != "tlcs_test_token" {
		t.Fatalf("CloudAPIToken() = %q, want CLI session token", token)
	}
}

func TestCloudAPITokenWithKindReportsCLISession(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken: "stored-github-token",
		CLISessionToken:   "stored-cli-token",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, kind, err := CloudAPITokenWithKind()
	if err != nil {
		t.Fatalf("CloudAPITokenWithKind() error = %v", err)
	}
	if token != "stored-cli-token" || kind != "lgtm-cli" {
		t.Fatalf("CloudAPITokenWithKind() = (%q, %q), want CLI session token", token, kind)
	}
}

func TestCloudAPITokenFallsBackToGitHubTokenForOldCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken: "stored-github-token",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	token, kind, err := CloudAPITokenWithKind()
	if err != nil {
		t.Fatalf("CloudAPITokenWithKind() error = %v", err)
	}
	if token != "stored-github-token" || kind != "github" {
		t.Fatalf("CloudAPITokenWithKind() = (%q, %q), want legacy GitHub token", token, kind)
	}
}

func TestCloudAPITokenRejectsExpiredCLISession(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if err := SaveCloudCredentials(CloudCredentials{
		GitHubAccessToken:   "stored-github-token",
		CLISessionToken:     "stored-cli-token",
		CLISessionExpiresAt: time.Now().Add(-time.Minute),
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}
	if _, _, err := CloudAPITokenWithKind(); err == nil {
		t.Fatal("expected expired CLI session error")
	}
}

func TestCloudAPITokenRequiresStoredToken(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")

	if _, err := CloudAPIToken(); err == nil {
		t.Fatal("expected error when GitHub token is missing")
	}
}

func TestCloudAPITokenRequiresStoredTokenForMCP(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("LGTM_MCP", "1")

	_, err := CloudAPIToken()
	if err == nil {
		t.Fatal("expected error when GitHub token is missing")
	}
	message := err.Error()
	if !strings.Contains(message, "run `lgtm auth login` in a terminal") {
		t.Fatalf("missing MCP login hint: %q", message)
	}
	if strings.Contains(message, "GH_TOKEN") || strings.Contains(message, "GITHUB_TOKEN") {
		t.Fatalf("MCP login hint should not prefer env token fallbacks: %q", message)
	}
}
