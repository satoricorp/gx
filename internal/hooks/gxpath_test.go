package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A hook must not pin a binary that lives somewhere the system can reclaim.
// This is not hypothetical: `go test` builds its binary under TMPDIR, so any
// test that installs hooks into a real repository pins a path like
// /var/folders/.../T/go-build123/b381/cli.test. Once that path is gone the
// hook falls back to PATH and still works — but the directory can be handed
// out again, and then the hook executes whatever now sits there.
func TestInstallGxPathRefusesTemporaryLocations(t *testing.T) {
	if isTemporaryPath(filepath.Join(os.TempDir(), "go-build123", "b381", "cli.test")) != true {
		t.Fatal("a binary under TMPDIR was not recognized as temporary")
	}
	for _, path := range []string{"/tmp/gx", "/var/tmp/gx"} {
		if !isTemporaryPath(path) {
			t.Fatalf("%s was not recognized as temporary", path)
		}
	}
	for _, path := range []string{"/usr/local/bin/gx", filepath.Join(os.Getenv("HOME"), ".local", "bin", "gx")} {
		if isTemporaryPath(path) {
			t.Fatalf("%s was wrongly treated as temporary", path)
		}
	}
	// The root directory must never count, or every path would be temporary.
	if isTemporaryPath("/") {
		t.Fatal("/ was treated as temporary")
	}
}

// Under `go test` the running binary is itself temporary, so the resolver
// declines to pin anything and the hook resolves gx from PATH.
func TestInstallGxPathDeclinesTheTestBinary(t *testing.T) {
	got, err := installGxPath()
	if err != nil {
		t.Fatalf("installGxPath: %v", err)
	}
	if got != "" {
		t.Fatalf("installGxPath() = %q under go test, want empty so the hook uses PATH", got)
	}
}

// And an empty path renders a hook that resolves gx from PATH rather than one
// that pins the empty string.
func TestHookResolveGxWithoutAPinnedPath(t *testing.T) {
	script := hookResolveGx("")
	if !strings.Contains(script, `command -v gx`) {
		t.Fatalf("resolver does not fall back to PATH:\n%s", script)
	}
	// Whatever it pins must not be an absolute path, so the script's own
	// `case "$gx_bin" in /*)` guard clears it and PATH decides.
	if strings.Contains(script, `gx_bin="/`) {
		t.Fatalf("resolver pinned an absolute path when given none:\n%s", script)
	}
}
