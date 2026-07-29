package uploadauth

import (
	"testing"

	"github.com/satoricorp/totality/internal/auth"
	"github.com/satoricorp/totality/internal/cloud"
)

func TestLoadPrefersCloudCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOTALITY_HOME", dir)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")
	t.Setenv("TOTALITY_API_URL", "")
	t.Setenv("TOTALITY_CLOUD_URL", "https://api.totality.test")

	if err := auth.SaveUpload(auth.UploadCredentials{
		APIURL: "https://legacy.totality.test",
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
	if got.APIURL != "https://api.totality.test" || got.Token != "tlcs_cloud_token" {
		t.Fatalf("Load = %+v, want cloud credentials", got)
	}
}

func TestLoadFallsBackToLegacyUploadCredentials(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOTALITY_HOME", dir)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")
	t.Setenv("TOTALITY_API_URL", "")

	if err := auth.SaveUpload(auth.UploadCredentials{
		APIURL: "https://legacy.totality.test",
		Token:  "legacy-token",
		OrgID:  "org-123",
	}); err != nil {
		t.Fatal(err)
	}

	got, ok := Load()
	if !ok {
		t.Fatal("Load returned false")
	}
	if got.APIURL != "https://legacy.totality.test" || got.Token != "legacy-token" || got.OrgID != "org-123" {
		t.Fatalf("Load = %+v, want legacy upload credentials", got)
	}
}
