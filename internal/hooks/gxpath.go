package hooks

import (
	"os"
	"path/filepath"
	"strings"
)

// installGxPath resolves the gx binary to pin into a hook script, and declines
// to pin one that lives in a temporary directory.
//
// Hooks pin the install-time path and fall back to gx on PATH when it is gone.
// That fallback handles a path going stale, but pinning a temp path is worse
// than pinning nothing: the operating system reuses those directories, so a
// path that is missing today can be a different executable tomorrow, and the
// hook would run it. A `go test` binary is the ordinary way this happens —
// os.Executable() inside a test is something like
// /var/folders/.../T/go-build123/b381/cli.test — and any test that installs
// hooks into a real repository leaves that behind in it.
//
// Returning empty means the hook resolves gx from PATH every time, which is
// the correct behavior when the running binary has no durable location.
func installGxPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if isTemporaryPath(exe) {
		return "", nil
	}
	return exe, nil
}

// isTemporaryPath reports whether path lives under a directory the system is
// free to reclaim. Both the resolved and unresolved spellings are checked
// because macOS reports TMPDIR as /var/folders while os.Executable() reports
// /private/var/folders.
func isTemporaryPath(path string) bool {
	candidates := spellings(path)
	if len(candidates) == 0 {
		return false
	}
	for _, dir := range []string{os.TempDir(), "/tmp", "/var/tmp"} {
		for _, root := range spellings(dir) {
			if root == string(filepath.Separator) {
				continue
			}
			prefix := strings.TrimSuffix(root, string(filepath.Separator)) + string(filepath.Separator)
			for _, candidate := range candidates {
				if candidate == root || strings.HasPrefix(candidate, prefix) {
					return true
				}
			}
		}
	}
	return false
}

// spellings returns the ways one path can be written: as given, and with
// symlinks resolved. Both are needed because a temp directory is /var/folders
// in the environment and /private/var/folders once resolved, and the binary
// being tested may be named either way. Resolution walks up to the deepest
// component that exists, since a path being tested need not exist at all.
func spellings(path string) []string {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	path = filepath.Clean(path)
	out := []string{path}

	rest := ""
	current := path
	for {
		if real, err := filepath.EvalSymlinks(current); err == nil {
			if joined := filepath.Clean(filepath.Join(real, rest)); joined != path {
				out = append(out, joined)
			}
			break
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		rest = filepath.Join(filepath.Base(current), rest)
		current = parent
	}
	return out
}
