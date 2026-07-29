package uploadauth

import (
	"strings"

	"github.com/satoricorp/totality/internal/auth"
	"github.com/satoricorp/totality/internal/cloud"
)

// Load resolves capture upload credentials from current cloud auth first, then
// the legacy upload env/file shape used by early capture builds.
func Load() (auth.Credentials, bool) {
	creds, _, ok := LoadWithKind()
	return creds, ok
}

// LoadWithKind resolves capture upload credentials and reports whether they
// came from the current GitHub-backed cloud auth or legacy upload config.
func LoadWithKind() (auth.Credentials, string, bool) {
	if token, kind, err := cloud.CloudAPITokenWithKind(); err == nil {
		if apiURL := cloud.CloudURL(); strings.TrimSpace(apiURL) != "" && strings.TrimSpace(token) != "" {
			return auth.Credentials{
				APIURL: apiURL,
				Token:  token,
			}, kind, true
		}
	}
	if creds, ok := auth.Load(); ok {
		return creds, "legacy", true
	}
	return auth.Credentials{}, "", false
}

func HasCredentials() bool {
	_, ok := Load()
	return ok
}
