package hooks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/lgtm/internal/storage"
)

const (
	// globalHookMarker identifies scripts written by `lgtm init --global`. Every
	// global script also carries hookMarker so the existing lgtm-ownership checks
	// (lgtmOwnedHook / IsInstalled) keep recognizing them.
	globalHookMarker = "# lgtm global lifecycle hooks"

	// hooksPathConfigKey is the git config key a global install points at the
	// lgtm hooks directory.
	hooksPathConfigKey = "core.hooksPath"

	// RepoEnabledConfigKey is the per-repo opt-out. `git config lgtm.enabled false`
	// inside a repository stops lgtm lifecycle work there while leaving the
	// machine-wide install (and the repo's own hooks) intact.
	RepoEnabledConfigKey = "lgtm.enabled"
)

// globalHookSpec describes one script installed into the lgtm global hooks
// directory.
//
// A global core.hooksPath makes git ignore every repository's own
// .git/hooks/*, so lgtm has to shim more than the four hooks it cares about:
// each script here chains to the repository's own hook of the same name.
type globalHookSpec struct {
	name string
	// stdin marks hooks git feeds data on stdin. Those scripts must drain
	// stdin even when there is nothing to do, and must replay it to the
	// chained repo hook when lgtm consumed it first.
	stdin bool
	// lgtm names the lgtm lifecycle behavior this script runs before chaining.
	// Empty means the script is a pure forwarder.
	lgtm string
}

// globalHookSpecs lists every hook `lgtm init --global` installs.
//
// Deliberately excluded, because their mere existence changes what git does
// and a shim that exits 0 would silently replace the built-in behavior:
// push-to-checkout, proc-receive, and fsmonitor-watchman (the latter is
// addressed by core.fsmonitor path, not through core.hooksPath). git-p4 hooks
// (p4-*) are excluded as well; they only matter to git-p4 users.
func globalHookSpecs() []globalHookSpec {
	return []globalHookSpec{
		// lgtm lifecycle hooks: lgtm work first (never fatal), then the repo hook.
		{name: "prepare-commit-msg", lgtm: "prepare-commit-msg"},
		{name: "post-commit", lgtm: "post-commit"},
		{name: "post-rewrite", lgtm: "post-rewrite", stdin: true},
		{name: "pre-push", lgtm: "pre-push", stdin: true},

		// Pure forwarders: lgtm does nothing, but the repo's hook must still run.
		{name: "applypatch-msg"},
		{name: "pre-applypatch"},
		{name: "post-applypatch"},
		{name: "pre-commit"},
		{name: "pre-merge-commit"},
		{name: "commit-msg"},
		{name: "pre-rebase"},
		{name: "post-checkout"},
		{name: "post-merge"},
		{name: "pre-auto-gc"},
		{name: "post-index-change"},
		{name: "sendemail-validate"},
		{name: "update"},
		{name: "post-update"},
		{name: "reference-transaction", stdin: true},
		{name: "pre-receive", stdin: true},
		{name: "post-receive", stdin: true},
	}
}

// GlobalHookNames returns the hook scripts a global install owns.
func GlobalHookNames() []string {
	specs := globalHookSpecs()
	names := make([]string, 0, len(specs))
	for _, spec := range specs {
		names = append(names, spec.name)
	}
	return names
}

// GlobalHooksDir returns the machine-wide lgtm hooks directory: $LGTM_HOME/hooks
// when LGTM_HOME is set, otherwise ~/.lgtm/hooks.
func GlobalHooksDir() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "hooks"), nil
}

// GlobalInstallOptions configures a machine-wide hook install.
type GlobalInstallOptions struct {
	// HooksDir overrides the lgtm hooks directory. Defaults to GlobalHooksDir().
	HooksDir string
	// LgtmPath overrides the lgtm binary baked into the scripts.
	LgtmPath string
}

// GlobalInstallResult reports what a machine-wide install changed.
type GlobalInstallResult struct {
	HooksDir      string   `json:"hooksDir"`
	Scripts       []string `json:"scripts"`
	Written       []string `json:"written,omitempty"`
	ConfigUpdated bool     `json:"configUpdated"`
	PreviousPath  string   `json:"previousPath,omitempty"`
}

// GlobalState reports the machine-wide install state.
type GlobalState struct {
	HooksDir       string   `json:"hooksDir"`
	ConfiguredPath string   `json:"configuredPath,omitempty"`
	Enabled        bool     `json:"enabled"`
	Installed      bool     `json:"installed"`
	MissingScripts []string `json:"missingScripts,omitempty"`
	Conflict       string   `json:"conflict,omitempty"`
}

// NeedsRepair reports whether a machine-wide install exists but is not fully
// wired up: scripts written without core.hooksPath pointing at them, or
// core.hooksPath pointing at a lgtm directory with scripts missing.
func (s GlobalState) NeedsRepair() bool {
	if s.Enabled {
		return !s.Installed
	}
	return s.Installed
}

// InstallGlobal writes the chaining lifecycle hooks into the lgtm hooks
// directory and points git's global core.hooksPath at it. It never needs to
// run inside a repository.
//
// It refuses to take over a core.hooksPath that some other tool owns; the
// caller is expected to surface that error verbatim.
func InstallGlobal(ctx context.Context, opts GlobalInstallOptions) (GlobalInstallResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	result := GlobalInstallResult{}
	hooksDir := strings.TrimSpace(opts.HooksDir)
	if hooksDir == "" {
		resolved, err := GlobalHooksDir()
		if err != nil {
			return result, err
		}
		hooksDir = resolved
	}
	hooksDir, err := filepath.Abs(hooksDir)
	if err != nil {
		return result, err
	}
	result.HooksDir = hooksDir

	lgtmPath := strings.TrimSpace(opts.LgtmPath)
	if lgtmPath == "" {
		lgtmPath, err = installLgtmPath()
		if err != nil {
			return result, fmt.Errorf("resolve lgtm binary: %w", err)
		}
	}

	current, err := globalHooksPath(ctx)
	if err != nil {
		return result, err
	}
	result.PreviousPath = current
	if current != "" && !sameHooksPath(current, hooksDir) && !lgtmOwnedHooksDir(current) {
		return result, fmt.Errorf(
			"git config --global %s is already set to %s, which is not managed by lgtm; "+
				"lgtm chains to each repository's own hooks but will not take over another global hooks directory. "+
				"Move those hooks into %s (they will be chained from there) or unset the config with "+
				"`git config --global --unset %s`, then re-run `lgtm init --global`",
			hooksPathConfigKey, current, hooksDir, hooksPathConfigKey)
	}

	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return result, err
	}
	for _, spec := range globalHookSpecs() {
		result.Scripts = append(result.Scripts, spec.name)
		script := globalHookScript(spec, lgtmPath, hooksDir)
		path := filepath.Join(hooksDir, spec.name)
		if data, readErr := os.ReadFile(path); readErr == nil {
			if string(data) == script {
				continue
			}
			if !lgtmOwnedHook(spec.name, string(data)) {
				return result, fmt.Errorf("refusing to overwrite non-lgtm script at %s; move it aside before retrying `lgtm init --global`", path)
			}
		}
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			return result, fmt.Errorf("install global %s hook: %w", spec.name, err)
		}
		result.Written = append(result.Written, spec.name)
	}

	if sameHooksPath(current, hooksDir) {
		return result, nil
	}
	if err := setGlobalHooksPath(ctx, hooksDir); err != nil {
		return result, err
	}
	result.ConfigUpdated = true
	return result, nil
}

// GlobalStatus reports the machine-wide install state without changing it.
func GlobalStatus(ctx context.Context) (GlobalState, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	state := GlobalState{}
	hooksDir, err := GlobalHooksDir()
	if err != nil {
		return state, err
	}
	state.HooksDir = hooksDir
	current, err := globalHooksPath(ctx)
	if err != nil {
		return state, err
	}
	state.ConfiguredPath = current
	state.Enabled = current != "" && sameHooksPath(current, hooksDir)

	state.Installed = true
	for _, name := range GlobalHookNames() {
		data, readErr := os.ReadFile(filepath.Join(hooksDir, name))
		if readErr != nil || !strings.Contains(string(data), globalHookMarker) {
			state.Installed = false
			state.MissingScripts = append(state.MissingScripts, name)
		}
	}
	// Only call a foreign hooksPath a conflict once lgtm global hooks actually
	// exist; otherwise every husky/lefthook user would see a phantom problem.
	if !state.Enabled && current != "" && len(state.MissingScripts) < len(GlobalHookNames()) {
		state.Conflict = current
	}
	return state, nil
}

// EnabledForRepo reports whether lgtm lifecycle work should run in repoRoot.
//
// A repository opts out with `git config lgtm.enabled false`, which is how a
// user excludes one repo from a machine-wide `lgtm init --global` install. Unset
// or unreadable means enabled.
func EnabledForRepo(ctx context.Context, repoRoot string) bool {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return true
	}
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "config", "--bool", "--get", RepoEnabledConfigKey)
	out, err := cmd.Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(out)) != "false"
}

func globalHooksPath(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "config", "--global", "--get", hooksPathConfigKey)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			// Key not set.
			return "", nil
		}
		return "", fmt.Errorf("read git config --global %s: %w", hooksPathConfigKey, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func setGlobalHooksPath(ctx context.Context, hooksDir string) error {
	cmd := exec.CommandContext(ctx, "git", "config", "--global", hooksPathConfigKey, hooksDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git config --global %s: %w\n%s", hooksPathConfigKey, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// lgtmOwnedHooksDir reports whether dir holds lgtm-written global hook scripts.
func lgtmOwnedHooksDir(dir string) bool {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return false
	}
	for _, name := range GlobalHookNames() {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		if strings.Contains(string(data), globalHookMarker) {
			return true
		}
	}
	return false
}

func sameHooksPath(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a == "" || b == "" {
		return false
	}
	if filepath.Clean(a) == filepath.Clean(b) {
		return true
	}
	resolvedA, errA := filepath.EvalSymlinks(a)
	resolvedB, errB := filepath.EvalSymlinks(b)
	if errA != nil || errB != nil {
		return false
	}
	return filepath.Clean(resolvedA) == filepath.Clean(resolvedB)
}

// shellSingleQuote renders s as a literal POSIX shell single-quoted word.
func shellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// globalHookPreamble emits the shared header plus lgtm_resolve_local_hook, which
// finds the repository's own hook of the given name.
//
// It must not use `git rev-parse --git-path hooks/<name>`: that call honors
// core.hooksPath, so under a global install it resolves right back to this
// script. It reads the repo-scoped core.hooksPath with the global and system
// config files masked, and otherwise falls back to the common git dir.
func globalHookPreamble(hooksDir string) string {
	quotedDir := shellSingleQuote(filepath.Clean(hooksDir))
	return fmt.Sprintf(`#!/bin/sh
%[1]s
%[2]s
# Installed by `+"`lgtm init --global`"+`. A global core.hooksPath makes git ignore
# every repository's own .git/hooks, so each script here does lgtm's work (if any)
# and then runs the repository's own hook of the same name.
# Opt a repository out with: git config %[3]s false

lgtm_global_hooks_dir=%[4]s
lgtm_local_hook=""

lgtm_resolve_local_hook() {
  lgtm_local_hook=""
  lgtm_hook_name="$1"
  # Masking the global/system config keeps this from resolving to the lgtm dir.
  lgtm_dir="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null git config --get %[5]s 2>/dev/null)"
  if [ -z "$lgtm_dir" ]; then
    lgtm_dir="$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null)"
    [ -n "$lgtm_dir" ] || return 0
    lgtm_dir="$lgtm_dir/hooks"
  fi
  case "$lgtm_dir" in
    /*) ;;
    *) lgtm_dir="$(git rev-parse --show-toplevel 2>/dev/null)/$lgtm_dir" ;;
  esac
  # Never chain back into this directory: that would recurse forever.
  if [ "$lgtm_dir" = "$lgtm_global_hooks_dir" ]; then
    return 0
  fi
  case "$lgtm_dir" in
    %[4]s/*) return 0 ;;
  esac
  lgtm_candidate="$lgtm_dir/$lgtm_hook_name"
  [ -f "$lgtm_candidate" ] && [ -x "$lgtm_candidate" ] || return 0
  # A repo that ran plain `+"`lgtm init`"+` has lgtm's own hook here; running it would
  # repeat the work this script just did.
  if grep -q '%[1]s' "$lgtm_candidate" 2>/dev/null; then
    return 0
  fi
  lgtm_local_hook="$lgtm_candidate"
}
`, hookMarker, globalHookMarker, RepoEnabledConfigKey, quotedDir, hooksPathConfigKey)
}

// globalLgtmResolver emits shell that resolves lgtm at run time: the pinned
// install-time path when it is still an executable file, else lgtm from PATH.
// Unlike the per-repo hooks, a global script must never exit when lgtm is
// missing — it still has to chain to the repository's own hook — so lgtm_bin is
// left empty and the lgtm work is skipped instead.
func globalLgtmResolver(lgtmPath string) string {
	if strings.TrimSpace(lgtmPath) == "" {
		lgtmPath = "lgtm"
	}
	return fmt.Sprintf(`lgtm_bin=%s
case "$lgtm_bin" in /*) ;; *) lgtm_bin="" ;; esac
if [ ! -f "$lgtm_bin" ] || [ ! -x "$lgtm_bin" ]; then
  lgtm_bin="$(command -v lgtm 2>/dev/null)" || lgtm_bin=""
fi`, shellSingleQuote(lgtmPath))
}

// globalHookScript renders one global hook script.
func globalHookScript(spec globalHookSpec, lgtmPath, hooksDir string) string {
	preamble := globalHookPreamble(hooksDir)
	resolver := globalLgtmResolver(lgtmPath)
	switch spec.lgtm {
	case "prepare-commit-msg":
		return preamble + fmt.Sprintf(`
%s
lgtm_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$lgtm_bin" ] && [ -n "$lgtm_repo" ]; then
  "$lgtm_bin" __hooks prepare-commit-msg --repo "$lgtm_repo" --message-path "$1" || {
    echo "lgtm prepare-commit-msg failed; commit continues without a lgtm trailer" >&2
  }
fi
lgtm_resolve_local_hook prepare-commit-msg
if [ -n "$lgtm_local_hook" ]; then
  exec "$lgtm_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-commit":
		return preamble + fmt.Sprintf(`
%s
lgtm_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$lgtm_bin" ] && [ -n "$lgtm_repo" ]; then
  "$lgtm_bin" __hooks post-commit --repo "$lgtm_repo" || {
    echo "lgtm post-commit metadata recording failed" >&2
  }
fi
lgtm_resolve_local_hook post-commit
if [ -n "$lgtm_local_hook" ]; then
  exec "$lgtm_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-rewrite":
		return preamble + fmt.Sprintf(`
%s
lgtm_stdin="$(mktemp "${TMPDIR:-/tmp}/lgtm-post-rewrite.XXXXXX" 2>/dev/null)" || lgtm_stdin=""
if [ -n "$lgtm_stdin" ]; then
  trap 'rm -f "$lgtm_stdin"' EXIT
  cat > "$lgtm_stdin"
fi
lgtm_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$lgtm_bin" ] && [ -n "$lgtm_repo" ] && [ -n "$lgtm_stdin" ]; then
  "$lgtm_bin" __hooks post-rewrite --repo "$lgtm_repo" < "$lgtm_stdin" || {
    echo "lgtm post-rewrite metadata update failed" >&2
  }
fi
lgtm_resolve_local_hook post-rewrite
if [ -n "$lgtm_local_hook" ]; then
  if [ -n "$lgtm_stdin" ]; then
    "$lgtm_local_hook" "$@" < "$lgtm_stdin"
  else
    "$lgtm_local_hook" "$@" < /dev/null
  fi
  exit $?
fi
exit 0
`, resolver)
	case "pre-push":
		return preamble + fmt.Sprintf(`
%s
lgtm_remote="$1"
lgtm_url="$2"
lgtm_stdin="$(mktemp "${TMPDIR:-/tmp}/lgtm-pre-push.XXXXXX" 2>/dev/null)" || lgtm_stdin=""
if [ -n "$lgtm_stdin" ]; then
  trap 'rm -f "$lgtm_stdin"' EXIT
  cat > "$lgtm_stdin"
else
  # Always drain stdin so git never blocks writing refs to this hook.
  cat > /dev/null
fi
lgtm_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$lgtm_bin" ] && [ -n "$lgtm_repo" ] && [ -n "$lgtm_stdin" ]; then
  while read lgtm_local_ref lgtm_local_sha lgtm_remote_ref lgtm_remote_sha
  do
    if [ "$lgtm_local_sha" = "0000000000000000000000000000000000000000" ]; then
      continue
    fi
    if [ "$lgtm_remote_sha" = "0000000000000000000000000000000000000000" ]; then
      lgtm_range="$lgtm_local_sha"
    else
      lgtm_range="${lgtm_remote_sha}..${lgtm_local_sha}"
    fi
    "$lgtm_bin" capture push --remote "$lgtm_remote" --ref-range "$lgtm_range" --local-ref "$lgtm_local_ref" --head-sha "$lgtm_local_sha" --repo "$lgtm_repo" || true
  done < "$lgtm_stdin"
fi
lgtm_resolve_local_hook pre-push
if [ -n "$lgtm_local_hook" ]; then
  if [ -n "$lgtm_stdin" ]; then
    "$lgtm_local_hook" "$@" < "$lgtm_stdin"
  else
    "$lgtm_local_hook" "$@" < /dev/null
  fi
  exit $?
fi
exit 0
`, resolver)
	}

	drain := ""
	if spec.stdin {
		drain = "# Drain the data git feeds this hook so it never sees a short write.\ncat > /dev/null\n"
	}
	return preamble + fmt.Sprintf(`
lgtm_resolve_local_hook %s
if [ -n "$lgtm_local_hook" ]; then
  exec "$lgtm_local_hook" "$@"
fi
%sexit 0
`, spec.name, drain)
}
