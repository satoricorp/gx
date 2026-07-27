package background

import (
	"strings"
	"testing"
)

// TestStartDetachedGXRefusesTestBinaries pins the guard that keeps the suite
// from forking itself. StartDetachedGX re-executes os.Executable(), which under
// `go test` is the test binary — so an unguarded spawn runs the whole package's
// tests again in a detached child, and those children spawn more. One
// `go test ./...` was observed leaving 401 live `hooks.test capture sync
// --quiet` processes behind and timing the package out.
func TestStartDetachedGXRefusesTestBinaries(t *testing.T) {
	if !hasTestBinaryName("/tmp/go-build123/b001/hooks.test") {
		t.Fatal(`hasTestBinaryName("…/hooks.test") = false, want true`)
	}
	if hasTestBinaryName("/usr/local/bin/gx") {
		t.Fatal(`hasTestBinaryName("/usr/local/bin/gx") = true, want false`)
	}
	// The running executable is itself a test binary, so the real call must
	// refuse rather than fork.
	err := StartDetachedGX("capture", "sync", "--quiet")
	if err == nil {
		t.Fatal("StartDetachedGX() error = nil inside a test binary, want a refusal instead of a forked test suite")
	}
	if !strings.Contains(err.Error(), "refusing to spawn") {
		t.Fatalf("StartDetachedGX() error = %v, want the refusal", err)
	}
}
