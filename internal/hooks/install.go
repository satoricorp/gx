package hooks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	hookMarker          = "# tx lifecycle hooks"
	legacyPrePushMarker = "# tx capture pre-push hook"
)

// AgentHookConfig documents future Claude/Codex lifecycle hook settings (V1.1).
type AgentHookConfig struct {
	Tool    string `json:"tool"`
	Enabled bool   `json:"enabled"`
}

// InstallOptions configures git lifecycle hook installation.
type InstallOptions struct {
	RepoRoot     string
	TotalityPath string
}

// Install writes tx lifecycle hooks using git's hooks directory resolution.
func Install(opts InstallOptions) error {
	repoRoot, err := filepath.Abs(opts.RepoRoot)
	if err != nil {
		return err
	}
	txPath := opts.TotalityPath
	if txPath == "" {
		txPath, err = installTotalityPath()
		if err != nil {
			return fmt.Errorf("resolve tx binary: %w", err)
		}
	}
	hooksDir, err := ResolveHooksDir(context.Background(), repoRoot)
	if err != nil {
		return err
	}
	// Under a machine-wide install (`tx init --global`) git resolves every
	// repo's hooks dir to the shared Totality directory. Those scripts chain to repo
	// hooks; overwriting them with the per-repo variants would break chaining
	// for every repository on the machine.
	if globalDir, globalErr := GlobalHooksDir(); globalErr == nil && sameHooksPath(hooksDir, globalDir) {
		return nil
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}
	scripts := map[string]string{
		"prepare-commit-msg": prepareCommitMsgScript(txPath),
		"post-commit":        postCommitScript(txPath),
		"post-rewrite":       postRewriteScript(txPath),
		"pre-push":           prePushScript(txPath),
	}
	for name, script := range scripts {
		hookPath := filepath.Join(hooksDir, name)
		if data, readErr := os.ReadFile(hookPath); readErr == nil {
			if string(data) == script {
				continue
			}
			if !totalityOwnedHook(name, string(data)) {
				return fmt.Errorf("refusing to overwrite existing %s hook at %s; preserve or chain it before retrying tx init", name, hookPath)
			}
		}
		if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
			return fmt.Errorf("install %s hook: %w", name, err)
		}
	}
	return nil
}

// IsInstalled reports whether tx lifecycle hooks are present.
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

func totalityOwnedHook(name, content string) bool {
	if containsHookMarker(content) {
		return true
	}
	return name == "pre-push" && strings.Contains(content, legacyPrePushMarker)
}

func hookRepoArg() string {
	return `repo="$(git rev-parse --show-toplevel)"`
}

// hookResolveTotality emits shell that resolves the tx binary at run time: the
// pinned install-time path when it is still an executable file, else tx from
// PATH. Pinned paths go stale — reinstalls move the binary, and a tx run from
// a temporary location pins that location — and a hook that cannot find tx
// must skip Totality work silently rather than break every commit and push.
func hookResolveTotality(txPath string) string {
	if txPath == "" {
		txPath = "tx"
	}
	return fmt.Sprintf(`tx_bin=%q
case "$tx_bin" in /*) ;; *) tx_bin="" ;; esac
if [ ! -f "$tx_bin" ] || [ ! -x "$tx_bin" ]; then
  tx_bin="$(command -v tx 2>/dev/null)" || tx_bin=""
fi
[ -n "$tx_bin" ] || exit 0`, txPath)
}

func prepareCommitMsgScript(txPath string) string {
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%s
"$tx_bin" __hooks prepare-commit-msg --repo "$repo" --message-path "$1"
`, hookMarker, hookResolveTotality(txPath), hookRepoArg())
}

func postCommitScript(txPath string) string {
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%s
"$tx_bin" __hooks post-commit --repo "$repo" || {
  echo "tx post-commit metadata recording failed" >&2
  exit 0
}
`, hookMarker, hookResolveTotality(txPath), hookRepoArg())
}

func postRewriteScript(txPath string) string {
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%s
"$tx_bin" __hooks post-rewrite --repo "$repo" || {
  echo "tx post-rewrite metadata update failed" >&2
  exit 0
}
`, hookMarker, hookResolveTotality(txPath), hookRepoArg())
}

func prePushScript(txPath string) string {
	return fmt.Sprintf(`#!/bin/sh
%s
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
  "$tx_bin" capture push --remote "$remote" --ref-range "$range" --local-ref "$local_ref" --head-sha "$local_sha" --repo "$(git rev-parse --show-toplevel)" || true
done
exit 0
`, hookMarker, hookResolveTotality(txPath))
}
