package cloud

import (
	"os"
	"strings"

	"github.com/satoricorp/totality/internal/buildconfig"
)

func GitHubClientID() string {
	return buildconfig.GitHubClientIDOrEnv()
}

func ConvexSiteURL() string {
	return buildconfig.ConvexSiteURLFromEnv()
}

func CloudURL() string {
	raw, ok := os.LookupEnv("TOTALITY_CLOUD_URL")
	if !ok {
		// No env override and no baked endpoint means cloud is disabled.
		// A binary built without ldflags must not fall back to poking a
		// developer localhost port; set TOTALITY_CLOUD_URL for local server work.
		return normalizeCloudURL(buildconfig.CloudURLFromEnvOrEmbedded())
	}
	return normalizeCloudURL(raw)
}

func normalizeCloudURL(raw string) string {
	url := strings.TrimRight(strings.TrimSpace(raw), "/")
	if url == "" || url == "0" || strings.EqualFold(url, "off") {
		return ""
	}
	const legacyPublishSuffix = "/tl/pr"
	if strings.HasSuffix(strings.ToLower(url), legacyPublishSuffix) {
		return strings.TrimRight(url[:len(url)-len(legacyPublishSuffix)], "/")
	}
	return url
}

func CloudConfigured() bool {
	return CloudURL() != ""
}
