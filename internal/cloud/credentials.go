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

// CloudCredentials holds gx cloud auth state.
type CloudCredentials struct {
	GitHubAccessToken           string    `json:"github_access_token,omitempty"`
	GitHubAccessTokenExpiresAt  time.Time `json:"github_access_token_expires_at,omitempty"`
	GitHubRefreshToken          string    `json:"github_refresh_token,omitempty"`
	GitHubRefreshTokenExpiresAt time.Time `json:"github_refresh_token_expires_at,omitempty"`
	GitHubKeychainAccount       string    `json:"github_keychain_account,omitempty"`
	CLISessionToken             string    `json:"cli_session_token,omitempty"`
	CLISessionExpiresAt         time.Time `json:"cli_session_expires_at,omitempty"`
	UserID                      string    `json:"user_id"`
	Login                       string    `json:"login"`
	AvatarURL                   string    `json:"avatar_url,omitempty"`
	MachineID                   string    `json:"machine_id"`
	MachineName                 string    `json:"machine_name"`
	ObtainedAt                  time.Time `json:"obtained_at"`
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

func GitHubAccessToken() (string, error) {
	return authstore.GitHubAccessToken()
}

func GitHubAccessTokenWithSource() (string, string, error) {
	return authstore.GitHubAccessTokenWithSource()
}

// CloudAPIToken returns the bearer token used by consolidated gx-cloud HTTP APIs.
func CloudAPIToken() (string, error) {
	token, _, err := CloudAPITokenWithKind()
	return token, err
}

// CloudAPITokenWithKind returns the token plus its source for callers that need
// to decide whether it is safe/useful to inject into another process.
func CloudAPITokenWithKind() (string, string, error) {
	creds, err := LoadCloudCredentials()
	if err != nil {
		return "", "", err
	}
	if creds != nil && strings.TrimSpace(creds.CLISessionToken) != "" {
		if !creds.CLISessionExpiresAt.IsZero() && time.Now().UTC().After(creds.CLISessionExpiresAt.UTC()) {
			if strings.TrimSpace(os.Getenv("GX_MCP")) != "" {
				return "", "gx-cli", fmt.Errorf("gx cloud session expired: run `gx auth logout` then `gx auth login` in a terminal, then retry the MCP tool")
			}
			return "", "gx-cli", fmt.Errorf("gx cloud session expired: run `gx auth logout` then `gx auth login`")
		}
		return strings.TrimSpace(creds.CLISessionToken), "gx-cli", nil
	}
	if token, err := GitHubAccessToken(); err == nil {
		return token, "github", nil
	} else {
		return "", "", err
	}
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
