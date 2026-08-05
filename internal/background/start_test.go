package background

import (
	"strings"
	"testing"
)

// TestStartDetachedGxRefusesTestBinaries pins the guard that keeps the suite
// from forking itself. StartDetachedGx re-executes os.Executable(), which under
// `go test` is the test binary — so an unguarded spawn runs the whole package's
// tests again in a detached child, and those children spawn more. One
// `go test ./...` was observed leaving 401 live `hooks.test capture sync
// --quiet` processes behind and timing the package out.
func TestStartDetachedGxRefusesTestBinaries(t *testing.T) {
	if !hasTestBinaryName("/tmp/go-build123/b001/hooks.test") {
		t.Fatal(`hasTestBinaryName("…/hooks.test") = false, want true`)
	}
	if hasTestBinaryName("/usr/local/bin/gx") {
		t.Fatal(`hasTestBinaryName("/usr/local/bin/gx") = true, want false`)
	}
	// The running executable is itself a test binary, so the real call must
	// refuse rather than fork.
	err := StartDetachedGx("capture", "sync", "--quiet")
	if err == nil {
		t.Fatal("StartDetachedGx() error = nil inside a test binary, want a refusal instead of a forked test suite")
	}
	if !strings.Contains(err.Error(), "refusing to spawn") {
		t.Fatalf("StartDetachedGx() error = %v, want the refusal", err)
	}
}
