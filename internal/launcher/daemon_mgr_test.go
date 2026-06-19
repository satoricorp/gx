package launcher

import (
	"testing"

	gxservice "github.com/satoricorp/gx/internal/service"
)

func TestDaemonManagerControlURLUsesStableControlURL(t *testing.T) {
	t.Setenv("GX_DAEMON_ADDR", "")
	manager := NewDaemonManager()

	got, err := manager.controlURL()
	if err != nil {
		t.Fatal(err)
	}
	if got != "http://"+gxservice.DefaultControl {
		t.Fatalf("controlURL() = %q, want %q", got, "http://"+gxservice.DefaultControl)
	}
}

func TestDaemonManagerControlURLHonorsExplicitOverride(t *testing.T) {
	t.Setenv("GX_DAEMON_ADDR", "127.0.0.1:49998")
	manager := NewDaemonManager()

	got, err := manager.controlURL()
	if err != nil {
		t.Fatal(err)
	}
	if got != "http://127.0.0.1:49998" {
		t.Fatalf("controlURL() = %q", got)
	}
}
