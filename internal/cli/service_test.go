package cli

import (
	"testing"

	gxservice "github.com/satoricorp/gx/internal/service"
)

func TestDaemonControlURLUsesStableControlURL(t *testing.T) {
	t.Setenv("GX_DAEMON_ADDR", "")

	got, err := daemonControlURL()
	if err != nil {
		t.Fatal(err)
	}
	if got != "http://"+gxservice.DefaultControl {
		t.Fatalf("daemonControlURL() = %q, want %q", got, "http://"+gxservice.DefaultControl)
	}
}

func TestDaemonControlURLHonorsExplicitOverride(t *testing.T) {
	t.Setenv("GX_DAEMON_ADDR", "127.0.0.1:49998")

	got, err := daemonControlURL()
	if err != nil {
		t.Fatal(err)
	}
	if got != "http://127.0.0.1:49998" {
		t.Fatalf("daemonControlURL() = %q", got)
	}
}
