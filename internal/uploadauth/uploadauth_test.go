package uploadauth

import (
	"testing"

	"github.com/satoricorp/gx/internal/auth"
	"github.com/satoricorp/gx/internal/cloud"
)

func TestLoadPrefersCloudCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GX_HOME", dir)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")
	t.Setenv("GX_API_URL", "")
	t.Setenv("GX_CLOUD_URL", "https://api.gx.test")

	if err := auth.SaveUpload(auth.UploadCredentials{
		APIURL: "https://legacy.gx.test",
		Token:  "legacy-token",
	}); err != nil {
		t.Fatal(err)
	}
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_cloud_token",
		CLISessionToken:   "gxcs_cloud_token",
	}); err != nil {
		t.Fatal(err)
	}

	got, ok := Load()
	if !ok {
		t.Fatal("Load returned false")
	}
	if got.APIURL != "https://api.gx.test" || got.Token != "gxcs_cloud_token" {
		t.Fatalf("Load = %+v, want cloud credentials", got)
	}
}

func TestLoadFallsBackToLegacyUploadCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GX_HOME", dir)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")
	t.Setenv("GX_API_URL", "")

	if err := auth.SaveUpload(auth.UploadCredentials{
		APIURL: "https://legacy.gx.test",
		Token:  "legacy-token",
		OrgID:  "org-123",
	}); err != nil {
		t.Fatal(err)
	}

	got, ok := Load()
	if !ok {
		t.Fatal("Load returned false")
	}
	if got.APIURL != "https://legacy.gx.test" || got.Token != "legacy-token" || got.OrgID != "org-123" {
		t.Fatalf("Load = %+v, want legacy upload credentials", got)
	}
}
