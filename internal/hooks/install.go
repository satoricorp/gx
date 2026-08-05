package hooks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	hookMarker          = "# lgtm lifecycle hooks"
	legacyPrePushMarker = "# lgtm capture pre-push hook"
)

// AgentHookConfig documents future Claude/Codex lifecycle hook settings (V1.1).
type AgentHookConfig struct {
	Tool    string `json:"tool"`
	Enabled bool   `json:"enabled"`
}

// InstallOptions configures git lifecycle hook installation.
type InstallOptions struct {
	RepoRoot     string
	LgtmPath string
}

// Install writes lgtm lifecycle hooks using git's hooks directory resolution.
func Install(opts InstallOptions) error {
	repoRoot, err := filepath.Abs(opts.RepoRoot)
	if err != nil {
		return err
	}
	lgtmPath := opts.LgtmPath
	if lgtmPath == "" {
		lgtmPath, err = installLgtmPath()
		if err != nil {
			return fmt.Errorf("resolve lgtm binary: %w", err)
		}
	}
	hooksDir, err := ResolveHooksDir(context.Background(), repoRoot)
	if err != nil {
		return err
	}
	// Under a machine-wide install (`lgtm init --global`) git resolves every
	// repo's hooks dir to the shared lgtm directory. Those scripts chain to repo
	// hooks; overwriting them with the per-repo variants would break chaining
	// for every repository on the machine.
	if globalDir, globalErr := GlobalHooksDir(); globalErr == nil && sameHooksPath(hooksDir, globalDir) {
		return nil
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}
	scripts := map[string]string{
		"prepare-commit-msg": prepareCommitMsgScript(lgtmPath),
		"post-commit":        postCommitScript(lgtmPath),
		"post-rewrite":       postRewriteScript(lgtmPath),
		"pre-push":           prePushScript(lgtmPath),
	}
	for name, script := range scripts {
		hookPath := filepath.Join(hooksDir, name)
		if data, readErr := os.ReadFile(hookPath); readErr == nil {
			if string(data) == script {
				continue
			}
			if !lgtmOwnedHook(name, string(data)) {
				return fmt.Errorf("refusing to overwrite existing %s hook at %s; preserve or chain it before retrying lgtm init", name, hookPath)
			}
		}
		if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
			return fmt.Errorf("install %s hook: %w", name, err)
		}
	}
	return nil
}

// IsInstalled reports whether lgtm lifecycle hooks are present.
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

func lgtmOwnedHook(name, content string) bool {
	if containsHookMarker(content) {
		return true
	}
	return name == "pre-push" && strings.Contains(content, legacyPrePushMarker)
}

func hookRepoArg() string {
	return `repo="$(git rev-parse --show-toplevel)"`
}

// hookResolveLgtm emits shell that resolves the lgtm binary at run time: the
// pinned install-time path when it is still an executable file, else lgtm from
// PATH. Pinned paths go stale — reinstalls move the binary, and a lgtm run from
// a temporary location pins that location — and a hook that cannot find lgtm
// must skip lgtm work silently rather than break every commit and push.
func hookResolveLgtm(lgtmPath string) string {
	if lgtmPath == "" {
		lgtmPath = "lgtm"
	}
	return fmt.Sprintf(`lgtm_bin=%q
case "$lgtm_bin" in /*) ;; *) lgtm_bin="" ;; esac
if [ ! -f "$lgtm_bin" ] || [ ! -x "$lgtm_bin" ]; then
  lgtm_bin="$(command -v lgtm 2>/dev/null)" || lgtm_bin=""
fi
[ -n "$lgtm_bin" ] || exit 0`, lgtmPath)
}

func prepareCommitMsgScript(lgtmPath string) string {
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%s
"$lgtm_bin" __hooks prepare-commit-msg --repo "$repo" --message-path "$1"
`, hookMarker, hookResolveLgtm(lgtmPath), hookRepoArg())
}

func postCommitScript(lgtmPath string) string {
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%s
"$lgtm_bin" __hooks post-commit --repo "$repo" || {
  echo "lgtm post-commit metadata recording failed" >&2
  exit 0
}
`, hookMarker, hookResolveLgtm(lgtmPath), hookRepoArg())
}

func postRewriteScript(lgtmPath string) string {
	return fmt.Sprintf(`#!/bin/sh
%s
%s
%s
"$lgtm_bin" __hooks post-rewrite --repo "$repo" || {
  echo "lgtm post-rewrite metadata update failed" >&2
  exit 0
}
`, hookMarker, hookResolveLgtm(lgtmPath), hookRepoArg())
}

func prePushScript(lgtmPath string) string {
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
  "$lgtm_bin" capture push --remote "$remote" --ref-range "$range" --local-ref "$local_ref" --head-sha "$local_sha" --repo "$(git rev-parse --show-toplevel)" || true
done
exit 0
`, hookMarker, hookResolveLgtm(lgtmPath))
}
