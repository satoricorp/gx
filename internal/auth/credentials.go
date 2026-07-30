package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Credentials hold CLI upload auth for the Totality server.
type Credentials struct {
	APIURL string `json:"api_url"`
	Token  string `json:"token"`
	OrgID  string `json:"org_id,omitempty"`
}

// DefaultAPIURL is the local dev server port from server/src/index.ts.
const DefaultAPIURL = "http://localhost:3201"

// Load reads upload credentials from env or ~/.totality/upload.json.
func Load() (Credentials, bool) {
	upload, ok := LoadUpload()
	if !ok {
		return Credentials{}, false
	}
	return Credentials{
		APIURL: upload.APIURL,
		Token:  upload.Token,
		OrgID:  upload.OrgID,
	}, true
}

// HasUploadCredentials reports whether upload to the server is configured.
func HasUploadCredentials() bool {
	_, ok := Load()
	return ok
}

// Save persists upload credentials (used by tx login).
func Save(creds Credentials) error {
	return SaveUpload(UploadCredentials{
		APIURL: creds.APIURL,
		Token:  creds.Token,
		OrgID:  creds.OrgID,
	})
}

func uploadCredentialsPath() (string, error) {
	if home := strings.TrimSpace(os.Getenv("TOTALITY_HOME")); home != "" {
		return filepath.Join(home, "upload.json"), nil
	}
	return UploadPath()
}

// Legacy helper for tests reading raw file shape.
func readUploadFile() (Credentials, error) {
	path, err := uploadCredentialsPath()
	if err != nil {
		return Credentials{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Credentials{}, err
	}
	var c Credentials
	if err := json.Unmarshal(data, &c); err != nil {
		return Credentials{}, err
	}
	return c, nil
}
