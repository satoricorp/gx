package hooks

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRunPushRespectsPauseFlag(t *testing.T) {
	home := t.TempDir()
	gxHome := filepath.Join(home, ".gx")
	if err := os.MkdirAll(gxHome, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gxHome, "pause-capture"), []byte("1\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	result, err := RunPush(context.Background(), PushOptions{
		RepoRoot: t.TempDir(),
		HomeDir:  home,
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if result.RefRange != "" || result.EligibleHunks != 0 {
		t.Fatalf("RunPush() = %+v, want empty result when paused", result)
	}
}

func TestParseRefRangeSingleSHALeavesBaseForOrchestrator(t *testing.T) {
	base, head := parseRefRange("abc123", "", "")
	if base != "" || head != "abc123" {
		t.Fatalf("parseRefRange() = (%q, %q), want empty base and abc123 head", base, head)
	}
}
