package gxtest_test

import (
	"os"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
)

// TestMain is also the subject of TestGitHelpersIgnoreAnEnclosingRepository,
// which re-runs this binary with an enclosing repository's GIT_DIR in its
// environment and relies on this call to keep the child out of it.
func TestMain(m *testing.M) {
	gxtest.DetachFromEnclosingGit()
	os.Exit(m.Run())
}
