package cloud

import (
	"os"
	"testing"
)

func TestCloudURLUsesDefaultWhenUnset(t *testing.T) {
	os.Unsetenv("GX_CLOUD_URL")
	if got := CloudURL(); got != "" {
		t.Fatalf("CloudURL() = %q, want empty (cloud disabled without env or baked endpoint)", got)
	}
}

func TestCloudURLExplicitDisable(t *testing.T) {
	for _, value := range []string{"", "0", "off", "OFF"} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("GX_CLOUD_URL", value)
			if got := CloudURL(); got != "" {
				t.Fatalf("CloudURL(%q) = %q, want disabled", value, got)
			}
		})
	}
}

func TestCloudURLCustom(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "https://api.example/gx/pr")
	if got := CloudURL(); got != "https://api.example" {
		t.Fatalf("CloudURL() = %q", got)
	}
}
