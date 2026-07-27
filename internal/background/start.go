package background

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// StartDetachedGX launches the gx binary with args in a detached background process.
//
// It refuses to launch anything that is not the gx binary. os.Executable()
// resolves to the *test* binary under `go test`, so an unguarded spawn
// re-executes the entire test suite as a detached child — which performs more
// pushes, which spawn more copies. Observed in practice: one `go test ./...`
// left 401 live `hooks.test capture sync --quiet` processes and blew the
// package timeout. Tests that mean to exercise the kickoff set
// GX_DISABLE_BACKGROUND_WORKERS; this guard is for the ones that forget.
func StartDetachedGX(args ...string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if runningUnderGoTest(exe) {
		return fmt.Errorf("refusing to spawn %q: not the gx binary", filepath.Base(exe))
	}
	cmd := exec.Command(exe, args...)
	if devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0); err == nil {
		cmd.Stdin = devNull
		cmd.Stdout = devNull
		cmd.Stderr = devNull
	}
	cmd.Env = os.Environ()
	configureDetachedCommand(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}
	return nil
}

// hasTestBinaryName reports whether a path looks like a binary `go test` built.
func hasTestBinaryName(exe string) bool {
	base := filepath.Base(exe)
	return strings.HasSuffix(base, ".test") || strings.HasSuffix(base, ".test.exe")
}

// runningUnderGoTest reports whether this process is a Go test binary. The
// registered testing flags are the backstop for a binary built with
// `go test -c -o <name>`, whose name gives nothing away.
func runningUnderGoTest(exe string) bool {
	return hasTestBinaryName(exe) || flag.Lookup("test.v") != nil
}
