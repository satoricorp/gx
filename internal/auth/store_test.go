package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadUploadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOTALITY_HOME", dir)
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")
	t.Setenv("TOTALITY_API_URL", "")

	want := UploadCredentials{
		APIURL: "http://localhost:3201",
		Token:  "dev-token",
		OrgID:  "org-123",
	}
	if err := SaveUpload(want); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "upload.json")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("upload.json mode = %o, want 0600", info.Mode().Perm())
	}

	got, ok := LoadUpload()
	if !ok {
		t.Fatal("LoadUpload returned false")
	}
	if got.Token != want.Token || got.APIURL != want.APIURL || got.OrgID != want.OrgID {
		t.Fatalf("LoadUpload = %+v, want %+v", got, want)
	}

	t.Setenv("TOTALITY_UPLOAD_TOKEN", "env-override")
	got, ok = LoadUpload()
	if !ok || got.Token != "env-override" {
		t.Fatalf("env override token = %+v ok=%v", got, ok)
	}
}

func TestLoadUploadMissing(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TOTALITY_HOME", dir)
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")
	_, ok := LoadUpload()
	if ok {
		t.Fatal("expected missing credentials")
	}
}
