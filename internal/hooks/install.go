package hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const hookMarker = "# gx capture pre-push hook"

// AgentHookConfig documents future Claude/Codex lifecycle hook settings (V1.1).
type AgentHookConfig struct {
	Tool    string `json:"tool"`
	Enabled bool   `json:"enabled"`
}

// InstallOptions configures pre-push hook installation.
type InstallOptions struct {
	RepoRoot string
	GXPath   string
}

// Install writes .git/hooks/pre-push invoking gx capture push.
func Install(opts InstallOptions) error {
	repoRoot, err := filepath.Abs(opts.RepoRoot)
	if err != nil {
		return err
	}
	gxPath := opts.GXPath
	if gxPath == "" {
		gxPath, err = os.Executable()
		if err != nil {
			return fmt.Errorf("resolve gx binary: %w", err)
		}
	}
	gitDir := filepath.Join(repoRoot, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}
	hookPath := filepath.Join(gitDir, "hooks", "pre-push")
	if err := os.MkdirAll(filepath.Dir(hookPath), 0o755); err != nil {
		return err
	}
	script := prePushScript(gxPath)
	if data, err := os.ReadFile(hookPath); err == nil {
		if string(data) == script {
			return nil
		}
		// A gx-owned hook with different content is an older template
		// (missing --local-ref/--head-sha, or blocking on failure) or
		// points at a stale binary; regenerate it.
	}
	return os.WriteFile(hookPath, []byte(script), 0o755)
}

// IsInstalled reports whether the gx pre-push hook is present.
func IsInstalled(repoRoot string) bool {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return false
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, ".git", "hooks", "pre-push"))
	if err != nil {
		return false
	}
	return containsHookMarker(string(data))
}

func containsHookMarker(content string) bool {
	return strings.Contains(content, hookMarker)
}

func prePushScript(gxPath string) string {
	if gxPath == "" {
		gxPath = "gx"
	}
	return fmt.Sprintf(`#!/bin/sh
%s
remote="$1"
url="$2"
while read local_ref local_sha remote_ref remote_sha
do
  if [ "$local_sha" = "0000000000000000000000000000000000000000" ]; then
    continue
  fi
  if [ "$remote_sha" = "0000000000000000000000000000000000000000" ]; then
    range="$local_sha"
  else
    range="${remote_sha}..${local_sha}"
  fi
  %q capture push --remote "$remote" --ref-range "$range" --local-ref "$local_ref" --head-sha "$local_sha" --repo "$(git rev-parse --show-toplevel)" || true
done
exit 0
`, hookMarker, gxPath)
}
