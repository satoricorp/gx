package capture

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsPaused(t *testing.T) {
	t.Run("env", func(t *testing.T) {
		t.Setenv("TOTALITY_CAPTURE_PAUSED", "1")
		if !IsPaused("") {
			t.Fatal("expected paused from env")
		}
	})

	t.Run("pause-capture file", func(t *testing.T) {
		t.Setenv("TOTALITY_CAPTURE_PAUSED", "")
		home := t.TempDir()
		totalityHome := filepath.Join(home, ".totality")
		t.Setenv("TOTALITY_HOME", totalityHome)
		if err := os.MkdirAll(totalityHome, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(totalityHome, pauseFileName), []byte("1\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if !IsPaused(home) {
			t.Fatal("expected paused from pause-capture file")
		}
	})

	t.Run("legacy capture-paused file", func(t *testing.T) {
		t.Setenv("TOTALITY_CAPTURE_PAUSED", "")
		home := t.TempDir()
		totalityHome := filepath.Join(home, ".totality")
		t.Setenv("TOTALITY_HOME", totalityHome)
		if err := os.MkdirAll(totalityHome, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(totalityHome, legacyPauseFileName), []byte("1\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if !IsPaused(home) {
			t.Fatal("expected paused from capture-paused file")
		}
	})

	t.Run("not paused", func(t *testing.T) {
		t.Setenv("TOTALITY_CAPTURE_PAUSED", "")
		home := t.TempDir()
		totalityHome := filepath.Join(home, ".totality")
		if err := os.MkdirAll(totalityHome, 0o700); err != nil {
			t.Fatal(err)
		}
		t.Setenv("TOTALITY_HOME", totalityHome)
		if IsPaused(home) {
			t.Fatal("expected not paused")
		}
	})
}
