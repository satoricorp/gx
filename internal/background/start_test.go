package background

import (
	"strings"
	"testing"
)

// TestStartDetachedTotalityRefusesTestBinaries pins the guard that keeps the suite
// from forking itself. StartDetachedTotality re-executes os.Executable(), which under
// `go test` is the test binary — so an unguarded spawn runs the whole package's
// tests again in a detached child, and those children spawn more. One
// `go test ./...` was observed leaving 401 live `hooks.test capture sync
// --quiet` processes behind and timing the package out.
func TestStartDetachedTotalityRefusesTestBinaries(t *testing.T) {
	if !hasTestBinaryName("/tmp/go-build123/b001/hooks.test") {
		t.Fatal(`hasTestBinaryName("…/hooks.test") = false, want true`)
	}
	if hasTestBinaryName("/usr/local/bin/tx") {
		t.Fatal(`hasTestBinaryName("/usr/local/bin/tx") = true, want false`)
	}
	// The running executable is itself a test binary, so the real call must
	// refuse rather than fork.
	err := StartDetachedTotality("capture", "sync", "--quiet")
	if err == nil {
		t.Fatal("StartDetachedTotality() error = nil inside a test binary, want a refusal instead of a forked test suite")
	}
	if !strings.Contains(err.Error(), "refusing to spawn") {
		t.Fatalf("StartDetachedTotality() error = %v, want the refusal", err)
	}
}
