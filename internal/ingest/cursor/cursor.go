// Package cursor knows where Cursor keeps its local chat/composer database.
//
// Cursor stores every chat/composer/agent thread under
// ~/Library/Application Support/Cursor/User/globalStorage/state.vscdb. The
// sync that copied those threads into the gx store is retired; what remains
// is the path lookup gx doctor uses to report whether the database exists.
package cursor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// DefaultVSCDBPath returns the macOS location of Cursor's global state.vscdb.
func DefaultVSCDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("cursor ingest is only supported on macOS in v0")
	}
	return filepath.Join(home, "Library", "Application Support", "Cursor", "User", "globalStorage", "state.vscdb"), nil
}
