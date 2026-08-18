package cli

import (
	"os"
	"os/exec"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
)

// TestMain cuts this package off from the network, and from the developer's
// own checkout.
//
// `gx review` is the command under test here, and its retrieval paths arm
// themselves from ambient credentials. setReviewGateEnv sandboxes the tests
// that call it, but it is opt-in, and the two tests in root_test.go that build
// their own environment were reaching the real embeddings API and the real
// TurboPuffer account on every run — which is the failure mode of any guard a
// test has to remember to ask for. Doing it here makes offline the default for
// the package and leaves the per-test helper as a way to be explicit.
//
// useScratchRepo is the same argument applied to the filesystem. See its
// comment for what a test was writing into the repository it ran from.
func TestMain(m *testing.M) {
	egress := gxtest.DenyNetwork()
	leaveScratchRepo := useScratchRepo()
	code := m.Run()
	// os.Exit skips deferred calls, so the scratch repo is released here.
	leaveScratchRepo()
	os.Exit(gxtest.FailOnEgress(code, egress()))
}

// useScratchRepo runs the package's tests from a throwaway git repository, and
// returns the func that puts the working directory back.
//
// Every command built by NewRoot auto-initializes the repository it finds,
// and initializing installs gx's lifecycle hooks. A test that executes a
// command without moving somewhere first finds the repository the test binary
// is running inside — the developer's own gx checkout. That is not
// hypothetical: `gx index` in index_test.go rewrote .git/hooks on every
// `go test ./internal/cli/`, pointing all four hooks at the ephemeral test
// binary under /var/folders. The path stops existing when the run ends, so
// what the tests left behind was a set of hooks surviving on their fallback.
//
// Tests that need a repository with real contents still build one and call
// t.Chdir, which restores the previous directory when they finish. This only
// decides where the tests that do not care about the working directory land,
// and the answer is: never the repository the suite is running inside.
func useScratchRepo() func() {
	previous, err := os.Getwd()
	if err != nil {
		panic("gx tests: read working directory: " + err.Error())
	}
	dir, err := os.MkdirTemp("", "gx-cli-scratch-repo-")
	if err != nil {
		panic("gx tests: create scratch repo: " + err.Error())
	}
	if out, err := exec.Command("git", "-C", dir, "init", "--quiet").CombinedOutput(); err != nil {
		panic("gx tests: git init scratch repo: " + err.Error() + ": " + string(out))
	}
	if err := os.Chdir(dir); err != nil {
		panic("gx tests: enter scratch repo: " + err.Error())
	}
	return func() {
		// Back out before removing, so the removal cannot fail on a working
		// directory that no longer exists.
		_ = os.Chdir(previous)
		_ = os.RemoveAll(dir)
	}
}
