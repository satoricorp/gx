package cli

import (
	"os"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
)

// TestMain cuts this package off from the network.
//
// `gx enhance` is the command under test here, and its retrieval paths arm
// themselves from ambient credentials. setReviewGateEnv sandboxes the tests
// that call it, but it is opt-in, and the two tests in root_test.go that build
// their own environment were reaching the real embeddings API and the real
// TurboPuffer account on every run — which is the failure mode of any guard a
// test has to remember to ask for. Doing it here makes offline the default for
// the package and leaves the per-test helper as a way to be explicit.
func TestMain(m *testing.M) {
	egress := gxtest.DenyNetwork()
	os.Exit(gxtest.FailOnEgress(m.Run(), egress()))
}
