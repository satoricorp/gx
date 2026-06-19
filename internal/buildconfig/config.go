package buildconfig

import (
	"os"
	"strings"
)

// Link-time defaults from .env via `just build` / `just install` (-ldflags -X).
var (
	GitHubClientID string
	ConvexSiteURL  string
	CloudURL       string
)

func envOrEmbedded(envKey, embedded string) string {
	if v := strings.TrimSpace(os.Getenv(envKey)); v != "" {
		return v
	}
	return strings.TrimSpace(embedded)
}

func GitHubClientIDOrEnv() string {
	return envOrEmbedded("GITHUB_CLIENT_ID", GitHubClientID)
}

func ConvexSiteURLFromEnv() string {
	return strings.TrimRight(envOrEmbedded("CONVEX_SITE_URL", ConvexSiteURL), "/")
}

func CloudURLFromEnvOrEmbedded() string {
	return strings.TrimSpace(envOrEmbedded("GX_CLOUD_URL", CloudURL))
}
