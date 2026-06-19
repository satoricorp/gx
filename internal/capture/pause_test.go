package capture

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsPaused(t *testing.T) {
	t.Run("env", func(t *testing.T) {
		t.Setenv("GX_CAPTURE_PAUSED", "1")
		if !IsPaused("") {
			t.Fatal("expected paused from env")
		}
	})

	t.Run("pause-capture file", func(t *testing.T) {
		home := t.TempDir()
		gxHome := filepath.Join(home, ".gx")
		if err := os.MkdirAll(gxHome, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(gxHome, pauseFileName), []byte("1\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if !IsPaused(home) {
			t.Fatal("expected paused from pause-capture file")
		}
	})

	t.Run("legacy capture-paused file", func(t *testing.T) {
		home := t.TempDir()
		gxHome := filepath.Join(home, ".gx")
		if err := os.MkdirAll(gxHome, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(gxHome, legacyPauseFileName), []byte("1\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if !IsPaused(home) {
			t.Fatal("expected paused from capture-paused file")
		}
	})

	t.Run("not paused", func(t *testing.T) {
		home := t.TempDir()
		if IsPaused(home) {
			t.Fatal("expected not paused")
		}
	})
}
