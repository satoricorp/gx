package cloud

import "testing"

func TestCloudBaseURLStripsGxPrSuffix(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "http://localhost:3200/gx/pr")
	if got := CloudBaseURL(); got != "http://localhost:3200" {
		t.Fatalf("CloudBaseURL() = %q, want http://localhost:3200", got)
	}
}

func TestRepoFullNameFromRemoteURL(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"git@github.com:satoricorp/gx.git", "satoricorp/gx"},
		{"https://github.com/example/widget", "example/widget"},
		{"https://github.com/example/widget.git", "example/widget"},
	}
	for _, tc := range tests {
		if got := RepoFullNameFromRemoteURL(tc.in); got != tc.want {
			t.Fatalf("RepoFullNameFromRemoteURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
