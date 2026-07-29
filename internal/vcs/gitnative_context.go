package vcs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const pendingCommitContextName = "totality-pending-commit.json"

// PendingCommitContext carries commit metadata consumed by lifecycle hooks.
// The post-commit hook tolerates a missing file and falls back to a zero value.
type PendingCommitContext struct {
	WorktreeRoot string `json:"worktree_root"`
	GitCommonDir string `json:"git_common_dir"`
	Branch       string `json:"branch,omitempty"`
}

func pendingCommitContextPath(gitDir string) string {
	return filepath.Join(gitDir, pendingCommitContextName)
}

func writePendingCommitContext(gitDir string, ctx PendingCommitContext) error {
	path := pendingCommitContextPath(gitDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(ctx)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func readPendingCommitContext(gitDir string) (PendingCommitContext, error) {
	path := pendingCommitContextPath(gitDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return PendingCommitContext{}, err
	}
	var ctx PendingCommitContext
	if err := json.Unmarshal(data, &ctx); err != nil {
		return PendingCommitContext{}, fmt.Errorf("decode pending commit context: %w", err)
	}
	return ctx, nil
}

func clearPendingCommitContext(gitDir string) {
	_ = os.Remove(pendingCommitContextPath(gitDir))
}
