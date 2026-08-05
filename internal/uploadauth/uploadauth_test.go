package uploadauth

import (
	"testing"

	"github.com/satoricorp/lgtm/internal/auth"
	"github.com/satoricorp/lgtm/internal/cloud"
)

func TestLoadPrefersCloudCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LGTM_HOME", dir)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("LGTM_UPLOAD_TOKEN", "")
	t.Setenv("LGTM_API_URL", "")
	t.Setenv("LGTM_CLOUD_URL", "https://api.lgtm.test")

	if err := auth.SaveUpload(auth.UploadCredentials{
		APIURL: "https://legacy.lgtm.test",
		Token:  "legacy-token",
	}); err != nil {
		t.Fatal(err)
	}
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_cloud_token",
		CLISessionToken:   "tlcs_cloud_token",
	}); err != nil {
		t.Fatal(err)
	}

	got, ok := Load()
	if !ok {
		t.Fatal("Load returned false")
	}
	if got.APIURL != "https://api.lgtm.test" || got.Token != "tlcs_cloud_token" {
		t.Fatalf("Load = %+v, want cloud credentials", got)
	}
}

func TestLoadFallsBackToLegacyUploadCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("LGTM_HOME", dir)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("LGTM_UPLOAD_TOKEN", "")
	t.Setenv("LGTM_API_URL", "")

	if err := auth.SaveUpload(auth.UploadCredentials{
		APIURL: "https://legacy.lgtm.test",
		Token:  "legacy-token",
		OrgID:  "org-123",
	}); err != nil {
		t.Fatal(err)
	}

	got, ok := Load()
	if !ok {
		t.Fatal("Load returned false")
	}
	if got.APIURL != "https://legacy.lgtm.test" || got.Token != "legacy-token" || got.OrgID != "org-123" {
		t.Fatalf("Load = %+v, want legacy upload credentials", got)
	}
}
