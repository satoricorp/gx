package cli

import (
	"testing"
)

func TestRenderStaticLogoPlainWithoutColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	got := renderStaticLogo()
	if got != tlLogoRaw {
		t.Fatalf("renderStaticLogo() = %q, want raw logo", got)
	}
}
