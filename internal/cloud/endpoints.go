package cloud

import (
	"os"
	"strings"

	"github.com/satoricorp/gx/internal/buildconfig"
)

const localDefaultCloudURL = "http://localhost:3201"

func GitHubClientID() string {
	return buildconfig.GitHubClientIDOrEnv()
}

func ConvexSiteURL() string {
	return buildconfig.ConvexSiteURLFromEnv()
}

func CloudURL() string {
	raw, ok := os.LookupEnv("GX_CLOUD_URL")
	if !ok {
		if embedded := buildconfig.CloudURLFromEnvOrEmbedded(); embedded != "" {
			return normalizeCloudURL(embedded)
		}
		return localDefaultCloudURL
	}
	return normalizeCloudURL(raw)
}

func normalizeCloudURL(raw string) string {
	url := strings.TrimRight(strings.TrimSpace(raw), "/")
	if url == "" || url == "0" || strings.EqualFold(url, "off") {
		return ""
	}
	const legacyPublishSuffix = "/gx/pr"
	if strings.HasSuffix(strings.ToLower(url), legacyPublishSuffix) {
		return strings.TrimRight(url[:len(url)-len(legacyPublishSuffix)], "/")
	}
	return url
}

func CloudConfigured() bool {
	return CloudURL() != ""
}
