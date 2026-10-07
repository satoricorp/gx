package codex_test

import (
	"os"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
)

// TestMain keeps this package's git commands in the repositories its tests
// build, even when the suite runs inside a git hook. See
// gxtest.DetachFromEnclosingGit.
func TestMain(m *testing.M) {
	gxtest.DetachFromEnclosingGit()
	os.Exit(m.Run())
}
