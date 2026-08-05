package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const uploadFileName = "upload.json"

// UploadCredentials hold CLI upload auth for the lgtm capture server.
type UploadCredentials struct {
	APIURL string `json:"api_url"`
	Token  string `json:"token"`
	OrgID  string `json:"org_id,omitempty"`
}

// UploadPath returns ~/.lgtm/upload.json (or $LGTM_HOME/upload.json).
func UploadPath() (string, error) {
	if home := strings.TrimSpace(os.Getenv("LGTM_HOME")); home != "" {
		return filepath.Join(home, uploadFileName), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".lgtm", uploadFileName), nil
}

// LoadUpload reads upload credentials from env (overrides file) and upload.json.
func LoadUpload() (UploadCredentials, bool) {
	c := UploadCredentials{
		APIURL: strings.TrimRight(strings.TrimSpace(os.Getenv("LGTM_API_URL")), "/"),
		Token:  strings.TrimSpace(os.Getenv("LGTM_UPLOAD_TOKEN")),
	}
	if org := strings.TrimSpace(os.Getenv("LGTM_ORG_ID")); org != "" {
		c.OrgID = org
	}
	if c.APIURL == "" {
		c.APIURL = DefaultAPIURL
	}

	path, err := UploadPath()
	if err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			var file UploadCredentials
			if json.Unmarshal(data, &file) == nil {
				if c.Token == "" {
					c.Token = strings.TrimSpace(file.Token)
				}
				if url := strings.TrimRight(strings.TrimSpace(file.APIURL), "/"); url != "" && c.APIURL == DefaultAPIURL {
					c.APIURL = url
				} else if url := strings.TrimRight(strings.TrimSpace(file.APIURL), "/"); url != "" && c.APIURL == "" {
					c.APIURL = url
				}
				if c.OrgID == "" {
					c.OrgID = strings.TrimSpace(file.OrgID)
				}
			}
		}
	}

	if strings.TrimSpace(c.Token) == "" {
		return UploadCredentials{}, false
	}
	return c, true
}

// SaveUpload persists upload credentials to ~/.lgtm/upload.json (0600).
func SaveUpload(creds UploadCredentials) error {
	path, err := UploadPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create lgtm dir: %w", err)
	}
	creds.APIURL = strings.TrimRight(strings.TrimSpace(creds.APIURL), "/")
	if creds.APIURL == "" {
		creds.APIURL = DefaultAPIURL
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal upload creds: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write upload creds: %w", err)
	}
	return nil
}

// ClearUpload removes stored upload credentials.
func ClearUpload() error {
	path, err := UploadPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
