package cloud

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/satoricorp/gx/internal/authstore"
	"github.com/satoricorp/gx/internal/storage"
)

type machineIDFile struct {
	ID string `json:"id"`
}

type credentialsFile struct {
	Cloud *CloudCredentials `json:"cloud,omitempty"`
}

// CloudCredentials holds gx cloud CLI session state.
type CloudCredentials struct {
	Token             string    `json:"token"`
	GitHubAccessToken string    `json:"github_access_token,omitempty"`
	UserID            string    `json:"user_id"`
	Login             string    `json:"login"`
	AvatarURL         string    `json:"avatar_url,omitempty"`
	SessionID         string    `json:"session_id"`
	MachineID         string    `json:"machine_id"`
	MachineName       string    `json:"machine_name"`
	ObtainedAt        time.Time `json:"obtained_at"`
}

func machineIDPath() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "machine_id.json"), nil
}

func credentialsPath() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

// DefaultMachineID returns the stable machine UUID, creating machine_id.json on first use.
func DefaultMachineID() (string, error) {
	path, err := machineIDPath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err == nil {
		var file machineIDFile
		if err := json.Unmarshal(data, &file); err != nil {
			return "", fmt.Errorf("parse machine id: %w", err)
		}
		if id := strings.TrimSpace(file.ID); id != "" {
			return id, nil
		}
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("read machine id: %w", err)
	}

	id := uuid.NewString()
	file := machineIDFile{ID: id}
	if err := writeJSONFile(path, file, 0o600); err != nil {
		return "", err
	}
	return id, nil
}

// LoadCloudCredentials reads cloud credentials from credentials.json.
func LoadCloudCredentials() (*CloudCredentials, error) {
	path, err := credentialsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read cloud credentials: %w", err)
	}
	var file credentialsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse cloud credentials: %w", err)
	}
	if file.Cloud == nil {
		return nil, nil
	}
	return file.Cloud, nil
}

// SaveCloudCredentials persists cloud credentials to credentials.json.
func SaveCloudCredentials(creds CloudCredentials) error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	file := credentialsFile{Cloud: &creds}
	return writeJSONFile(path, file, 0o600)
}

// ClearCloudCredentials removes the cloud section from credentials.json.
func ClearCloudCredentials() error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read cloud credentials: %w", err)
	}
	var file credentialsFile
	if err := json.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("parse cloud credentials: %w", err)
	}
	if file.Cloud == nil {
		return nil
	}
	file.Cloud = nil
	if err := writeJSONFile(path, file, 0o600); err != nil {
		return err
	}
	return nil
}

// BearerToken returns the bearer token for gx-cloud upload and revoke.
// The stored GitHub token is used only for GitHub API calls, not gx-cloud upload.
func BearerToken() (string, error) {
	creds, err := LoadCloudCredentials()
	if err != nil {
		return "", err
	}
	if creds == nil || strings.TrimSpace(creds.Token) == "" {
		return "", bearerTokenError()
	}
	token := strings.TrimSpace(creds.Token)
	if !validCliToken(token) {
		return "", fmt.Errorf("stored gx cloud token is invalid (run `gx auth logout` then `gx auth login` to refresh)")
	}
	return token, nil
}

func GitHubAccessToken() (string, error) {
	return authstore.GitHubAccessToken()
}

// CloudAPIToken returns the bearer token used by consolidated gx-cloud HTTP APIs.
// These Hono endpoints validate GitHub OAuth tokens directly. Fall back to the
// legacy Convex-issued token for older local scripts during the transition.
func CloudAPIToken() (string, error) {
	if token := cloudAPIKeyFromEnv(); token != "" {
		return token, nil
	}
	if token, err := GitHubAccessToken(); err == nil {
		return token, nil
	}
	return BearerToken()
}

func cloudAPIKeyFromEnv() string {
	if token := strings.TrimSpace(os.Getenv("GX_API_KEY")); token != "" {
		return token
	}
	return ""
}

func bearerTokenError() error {
	return fmt.Errorf("not logged in to gx cloud: run `gx auth login`")
}

func validCliToken(token string) bool {
	token = strings.TrimSpace(token)
	if !strings.HasPrefix(token, "gx_") {
		return false
	}
	rest := strings.TrimPrefix(token, "gx_")
	if len(rest) < 20 {
		return false
	}
	switch token {
	case "gx_baked", "gx_test_token", "gx_saved":
		return false
	}
	return true
}

func writeJSONFile(path string, v any, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create gx home dir: %w", err)
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", filepath.Base(path), err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, mode); err != nil {
		return fmt.Errorf("write %s: %w", filepath.Base(path), err)
	}
	return nil
}
