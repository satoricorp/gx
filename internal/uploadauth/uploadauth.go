package uploadauth

import (
	"strings"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/cloud"
)

// Load resolves capture upload credentials from current cloud auth first, then
// the legacy upload env/file shape used by early capture builds.
func Load() (auth.Credentials, bool) {
	if token, err := cloud.CloudAPIToken(); err == nil {
		if apiURL := cloud.CloudURL(); strings.TrimSpace(apiURL) != "" && strings.TrimSpace(token) != "" {
			return auth.Credentials{
				APIURL: apiURL,
				Token:  token,
			}, true
		}
	}
	return auth.Load()
}

func HasCredentials() bool {
	_, ok := Load()
	return ok
}
