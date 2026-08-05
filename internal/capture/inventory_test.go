package capture_test

import (
	"testing"

	"github.com/satoricorp/gx/internal/capture"
)

func TestInventoryUnknownPaths(t *testing.T) {
	inv := capture.NewInventoryCollector()
	inv.RecordPath("claude", "uuid", "string", "", "")
	inv.RecordPath("claude", "sessionId", "string", "session_id", "")
	unknown := inv.UnknownPaths()
	if len(unknown["claude"]) != 1 || unknown["claude"][0] != "uuid" {
		t.Fatalf("unknown paths = %#v", unknown)
	}
}
