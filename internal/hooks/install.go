package hooks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	hookMarker          = "# gx lifecycle hooks"
	legacyPrePushMarker = "# gx capture pre-push hook"
)

// AgentHookConfig documents future Claude/Codex lifecycle hook settings (V1.1).
type AgentHookConfig struct {
	Tool    string `json:"tool"`
	Enabled bool   `json:"enabled"`
}

// InstallOptions configures git lifecycle hook installation.
type InstallOptions struct {
	RepoRoot string
	GXPath   string
}

// Install writes gx lifecycle hooks using git's hooks directory resolution.
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
	hooksDir, err := ResolveHooksDir(context.Background(), repoRoot)
	if err != nil {
		return err
	}
	// Under a machine-wide install (`gx init --global`) git resolves every
	// repo's hooks dir to the shared GX directory. Those scripts chain to repo
	// hooks; overwriting them with the per-repo variants would break chaining
	// for every repository on the machine.
	if globalDir, globalErr := GlobalHooksDir(); globalErr == nil && sameHooksPath(hooksDir, globalDir) {
		return nil
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}
	scripts := map[string]string{
		"prepare-commit-msg": prepareCommitMsgScript(gxPath),
		"post-commit":        postCommitScript(gxPath),
		"post-rewrite":       postRewriteScript(gxPath),
		"pre-push":           prePushScript(gxPath),
	}
	for name, script := range scripts {
		hookPath := filepath.Join(hooksDir, name)
		if data, readErr := os.ReadFile(hookPath); readErr == nil {
			if string(data) == script {
				continue
			}
			if !gxOwnedHook(name, string(data)) {
				return fmt.Errorf("refusing to overwrite existing %s hook at %s; preserve or chain it before retrying gx init", name, hookPath)
			}
		}
		if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
			return fmt.Errorf("install %s hook: %w", name, err)
		}
	}
	return nil
}

// IsInstalled reports whether gx lifecycle hooks are present.
func IsInstalled(repoRoot string) bool {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return false
	}
	hooksDir, err := ResolveHooksDir(context.Background(), repoRoot)
	if err != nil {
		return false
	}
	for _, name := range []string{"prepare-commit-msg", "post-commit", "post-rewrite", "pre-push"} {
		data, err := os.ReadFile(filepath.Join(hooksDir, name))
		if err != nil || !containsHookMarker(string(data)) {
			return false
		}
	}
	return true
}

func containsHookMarker(content string) bool {
	return strings.Contains(content, hookMarker)
}

func gxOwnedHook(name, content string) bool {
	if containsHookMarker(content) {
		return true
	}
	return name == "pre-push" && strings.Contains(content, legacyPrePushMarker)
}

func hookRepoArg() string {
	return `repo="$(git rev-parse --show-toplevel)"`
}

func prepareCommitMsgScript(gxPath string) string {
	if gxPath == "" {
		gxPath = "gx"
	}
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%q __hooks prepare-commit-msg --repo "$repo" --message-path "$1"
`, hookMarker, hookRepoArg(), gxPath)
}

func postCommitScript(gxPath string) string {
	if gxPath == "" {
		gxPath = "gx"
	}
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%q __hooks post-commit --repo "$repo" || {
  echo "gx post-commit metadata recording failed" >&2
  exit 0
}
`, hookMarker, hookRepoArg(), gxPath)
}

func postRewriteScript(gxPath string) string {
	if gxPath == "" {
		gxPath = "gx"
	}
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%q __hooks post-rewrite --repo "$repo" || {
  echo "gx post-rewrite metadata update failed" >&2
  exit 0
}
`, hookMarker, hookRepoArg(), gxPath)
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
