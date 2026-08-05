package cloud

import (
	"os"
	"testing"

	"github.com/satoricorp/lgtm/internal/buildconfig"
)

func TestGitHubClientIDEnvOverridesDefault(t *testing.T) {
	t.Setenv("GITHUB_CLIENT_ID", " env-client ")
	buildconfig.GitHubClientID = "baked-client"
	t.Cleanup(func() { buildconfig.GitHubClientID = "" })

	if got := GitHubClientID(); got != "env-client" {
		t.Fatalf("GitHubClientID() = %q, want env-client", got)
	}
}

func TestGitHubClientIDUsesBakedDefault(t *testing.T) {
	t.Setenv("GITHUB_CLIENT_ID", "")
	buildconfig.GitHubClientID = "baked-client"
	t.Cleanup(func() { buildconfig.GitHubClientID = "" })

	if got := GitHubClientID(); got != "baked-client" {
		t.Fatalf("GitHubClientID() = %q, want baked-client", got)
	}
}

func TestGitHubClientIDEmptyWhenUnset(t *testing.T) {
	t.Setenv("GITHUB_CLIENT_ID", "")
	buildconfig.GitHubClientID = ""
	if got := GitHubClientID(); got != "" {
		t.Fatalf("GitHubClientID() = %q, want empty", got)
	}
}

func TestConvexSiteURLEnvOverridesDefault(t *testing.T) {
	t.Setenv("CONVEX_SITE_URL", " https://env.convex.site/ ")
	buildconfig.ConvexSiteURL = "https://baked.convex.site"
	t.Cleanup(func() { buildconfig.ConvexSiteURL = "" })

	if got := ConvexSiteURL(); got != "https://env.convex.site" {
		t.Fatalf("ConvexSiteURL() = %q, want trimmed env URL", got)
	}
}

func TestConvexSiteURLUsesBakedDefault(t *testing.T) {
	t.Setenv("CONVEX_SITE_URL", "")
	buildconfig.ConvexSiteURL = "https://baked.convex.site/"
	t.Cleanup(func() { buildconfig.ConvexSiteURL = "" })

	if got := ConvexSiteURL(); got != "https://baked.convex.site" {
		t.Fatalf("ConvexSiteURL() = %q, want trimmed baked URL", got)
	}
}

func TestConvexSiteURLEmptyWhenUnset(t *testing.T) {
	t.Setenv("CONVEX_SITE_URL", "")
	buildconfig.ConvexSiteURL = ""
	if got := ConvexSiteURL(); got != "" {
		t.Fatalf("ConvexSiteURL() = %q, want empty", got)
	}
}

func TestCloudURLEnvOverridesDefault(t *testing.T) {
	t.Setenv("LGTM_CLOUD_URL", "http://localhost:3201")
	buildconfig.CloudURL = "https://api.example.com/lgtm/pr"
	t.Cleanup(func() { buildconfig.CloudURL = "" })

	if got := CloudURL(); got != "http://localhost:3201" {
		t.Fatalf("CloudURL() = %q, want env override", got)
	}
}

func TestCloudURLUsesBakedDefault(t *testing.T) {
	os.Unsetenv("LGTM_CLOUD_URL")
	buildconfig.CloudURL = "https://api.example.com/lgtm/pr"
	t.Cleanup(func() { buildconfig.CloudURL = "" })

	if got := CloudURL(); got != "https://api.example.com" {
		t.Fatalf("CloudURL() = %q, want baked default", got)
	}
}

func TestCloudURLNormalizesLegacyPublishSuffix(t *testing.T) {
	t.Setenv("LGTM_CLOUD_URL", "http://localhost:3200/lgtm/pr/")
	buildconfig.CloudURL = ""
	t.Cleanup(func() { buildconfig.CloudURL = "" })

	if got := CloudURL(); got != "http://localhost:3200" {
		t.Fatalf("CloudURL() = %q, want legacy suffix stripped", got)
	}
}

func TestCloudURLDisabledWhenUnset(t *testing.T) {
	os.Unsetenv("LGTM_CLOUD_URL")
	buildconfig.CloudURL = ""
	t.Cleanup(func() { buildconfig.CloudURL = "" })

	if got := CloudURL(); got != "" {
		t.Fatalf("CloudURL() = %q, want empty (cloud disabled without env or baked endpoint)", got)
	}
	if CloudConfigured() {
		t.Fatal("CloudConfigured() = true, want false")
	}
}
