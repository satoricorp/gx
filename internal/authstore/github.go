package authstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

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
	data, err := os.ReadFile(filepath.Join(dir, "credentials.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", tokenError()
		}
		return "", fmt.Errorf("read github credentials: %w", err)
	}
	var file struct {
		Cloud *struct {
			GitHubAccessToken string `json:"github_access_token"`
		} `json:"cloud"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		return "", fmt.Errorf("parse github credentials: %w", err)
	}
	if file.Cloud == nil || strings.TrimSpace(file.Cloud.GitHubAccessToken) == "" {
		return "", tokenError()
	}
	return strings.TrimSpace(file.Cloud.GitHubAccessToken), nil
}

func tokenError() error {
	return fmt.Errorf("github token is not configured: run `gx auth login` or set GH_TOKEN/GITHUB_TOKEN")
}
