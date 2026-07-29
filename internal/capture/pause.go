package capture

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	pauseFileName       = "pause-capture"
	legacyPauseFileName = "capture-paused"
)

// IsPaused reports whether capture hooks should no-op.
// Respects TOTALITY_CAPTURE_PAUSED env and flag files under ~/.totality/.
func IsPaused(homeDir string) bool {
	if pausedFromEnv() {
		return true
	}
	if strings.TrimSpace(homeDir) == "" {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return false
		}
	}
	totalityHome := totalityHomeDir(homeDir)
	for _, name := range []string{pauseFileName, legacyPauseFileName} {
		if _, err := os.Stat(filepath.Join(totalityHome, name)); err == nil {
			return true
		}
	}
	return false
}

func pausedFromEnv() bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv("TOTALITY_CAPTURE_PAUSED")))
	return value == "1" || value == "true" || value == "yes"
}

func totalityHomeDir(homeDir string) string {
	if custom := strings.TrimSpace(os.Getenv("TOTALITY_HOME")); custom != "" {
		return custom
	}
	return filepath.Join(homeDir, ".totality")
}
