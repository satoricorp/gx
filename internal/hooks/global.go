package hooks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/totality/internal/storage"
)

const (
	// globalHookMarker identifies scripts written by `tx init --global`. Every
	// global script also carries hookMarker so the existing Totality-ownership checks
	// (totalityOwnedHook / IsInstalled) keep recognizing them.
	globalHookMarker = "# tx global lifecycle hooks"

	// hooksPathConfigKey is the git config key a global install points at the
	// Totality hooks directory.
	hooksPathConfigKey = "core.hooksPath"

	// RepoEnabledConfigKey is the per-repo opt-out. `git config tx.enabled false`
	// inside a repository stops Totality lifecycle work there while leaving the
	// machine-wide install (and the repo's own hooks) intact.
	RepoEnabledConfigKey = "tx.enabled"
)

// globalHookSpec describes one script installed into the Totality global hooks
// directory.
//
// A global core.hooksPath makes git ignore every repository's own
// .git/hooks/*, so Totality has to shim more than the four hooks it cares about:
// each script here chains to the repository's own hook of the same name.
type globalHookSpec struct {
	name string
	// stdin marks hooks git feeds data on stdin. Those scripts must drain
	// stdin even when there is nothing to do, and must replay it to the
	// chained repo hook when Totality consumed it first.
	stdin bool
	// tx names the Totality lifecycle behavior this script runs before chaining.
	// Empty means the script is a pure forwarder.
	tx string
}

// globalHookSpecs lists every hook `tx init --global` installs.
//
// Deliberately excluded, because their mere existence changes what git does
// and a shim that exits 0 would silently replace the built-in behavior:
// push-to-checkout, proc-receive, and fsmonitor-watchman (the latter is
// addressed by core.fsmonitor path, not through core.hooksPath). git-p4 hooks
// (p4-*) are excluded as well; they only matter to git-p4 users.
func globalHookSpecs() []globalHookSpec {
	return []globalHookSpec{
		// Totality lifecycle hooks: Totality work first (never fatal), then the repo hook.
		{name: "prepare-commit-msg", tx: "prepare-commit-msg"},
		{name: "post-commit", tx: "post-commit"},
		{name: "post-rewrite", tx: "post-rewrite", stdin: true},
		{name: "pre-push", tx: "pre-push", stdin: true},

		// Pure forwarders: Totality does nothing, but the repo's hook must still run.
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

// GlobalHooksDir returns the machine-wide Totality hooks directory: $TOTALITY_HOME/hooks
// when TOTALITY_HOME is set, otherwise ~/.totality/hooks.
func GlobalHooksDir() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "hooks"), nil
}

// GlobalInstallOptions configures a machine-wide hook install.
type GlobalInstallOptions struct {
	// HooksDir overrides the Totality hooks directory. Defaults to GlobalHooksDir().
	HooksDir string
	// TotalityPath overrides the tx binary baked into the scripts.
	TotalityPath string
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
// core.hooksPath pointing at a Totality directory with scripts missing.
func (s GlobalState) NeedsRepair() bool {
	if s.Enabled {
		return !s.Installed
	}
	return s.Installed
}

// InstallGlobal writes the chaining lifecycle hooks into the Totality hooks
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

	txPath := strings.TrimSpace(opts.TotalityPath)
	if txPath == "" {
		txPath, err = installTotalityPath()
		if err != nil {
			return result, fmt.Errorf("resolve tx binary: %w", err)
		}
	}

	current, err := globalHooksPath(ctx)
	if err != nil {
		return result, err
	}
	result.PreviousPath = current
	if current != "" && !sameHooksPath(current, hooksDir) && !totalityOwnedHooksDir(current) {
		return result, fmt.Errorf(
			"git config --global %s is already set to %s, which is not managed by Totality; "+
				"Totality chains to each repository's own hooks but will not take over another global hooks directory. "+
				"Move those hooks into %s (they will be chained from there) or unset the config with "+
				"`git config --global --unset %s`, then re-run `tx init --global`",
			hooksPathConfigKey, current, hooksDir, hooksPathConfigKey)
	}

	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return result, err
	}
	for _, spec := range globalHookSpecs() {
		result.Scripts = append(result.Scripts, spec.name)
		script := globalHookScript(spec, txPath, hooksDir)
		path := filepath.Join(hooksDir, spec.name)
		if data, readErr := os.ReadFile(path); readErr == nil {
			if string(data) == script {
				continue
			}
			if !totalityOwnedHook(spec.name, string(data)) {
				return result, fmt.Errorf("refusing to overwrite non-Totality script at %s; move it aside before retrying `tx init --global`", path)
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
	// Only call a foreign hooksPath a conflict once Totality global hooks actually
	// exist; otherwise every husky/lefthook user would see a phantom problem.
	if !state.Enabled && current != "" && len(state.MissingScripts) < len(GlobalHookNames()) {
		state.Conflict = current
	}
	return state, nil
}

// EnabledForRepo reports whether Totality lifecycle work should run in repoRoot.
//
// A repository opts out with `git config tx.enabled false`, which is how a
// user excludes one repo from a machine-wide `tx init --global` install. Unset
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

// totalityOwnedHooksDir reports whether dir holds Totality-written global hook scripts.
func totalityOwnedHooksDir(dir string) bool {
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

// globalHookPreamble emits the shared header plus tx_resolve_local_hook, which
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
# Installed by `+"`tx init --global`"+`. A global core.hooksPath makes git ignore
# every repository's own .git/hooks, so each script here does Totality's work (if any)
# and then runs the repository's own hook of the same name.
# Opt a repository out with: git config %[3]s false

tx_global_hooks_dir=%[4]s
tx_local_hook=""

tx_resolve_local_hook() {
  tx_local_hook=""
  tx_hook_name="$1"
  # Masking the global/system config keeps this from resolving to the Totality dir.
  tx_dir="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null git config --get %[5]s 2>/dev/null)"
  if [ -z "$tx_dir" ]; then
    tx_dir="$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null)"
    [ -n "$tx_dir" ] || return 0
    tx_dir="$tx_dir/hooks"
  fi
  case "$tx_dir" in
    /*) ;;
    *) tx_dir="$(git rev-parse --show-toplevel 2>/dev/null)/$tx_dir" ;;
  esac
  # Never chain back into this directory: that would recurse forever.
  if [ "$tx_dir" = "$tx_global_hooks_dir" ]; then
    return 0
  fi
  case "$tx_dir" in
    %[4]s/*) return 0 ;;
  esac
  tx_candidate="$tx_dir/$tx_hook_name"
  [ -f "$tx_candidate" ] && [ -x "$tx_candidate" ] || return 0
  # A repo that ran plain `+"`tx init`"+` has Totality's own hook here; running it would
  # repeat the work this script just did.
  if grep -q '%[1]s' "$tx_candidate" 2>/dev/null; then
    return 0
  fi
  tx_local_hook="$tx_candidate"
}
`, hookMarker, globalHookMarker, RepoEnabledConfigKey, quotedDir, hooksPathConfigKey)
}

// globalTotalityResolver emits shell that resolves tx at run time: the pinned
// install-time path when it is still an executable file, else tx from PATH.
// Unlike the per-repo hooks, a global script must never exit when tx is
// missing — it still has to chain to the repository's own hook — so tx_bin is
// left empty and the Totality work is skipped instead.
func globalTotalityResolver(txPath string) string {
	if strings.TrimSpace(txPath) == "" {
		txPath = "tx"
	}
	return fmt.Sprintf(`tx_bin=%s
case "$tx_bin" in /*) ;; *) tx_bin="" ;; esac
if [ ! -f "$tx_bin" ] || [ ! -x "$tx_bin" ]; then
  tx_bin="$(command -v tx 2>/dev/null)" || tx_bin=""
fi`, shellSingleQuote(txPath))
}

// globalHookScript renders one global hook script.
func globalHookScript(spec globalHookSpec, txPath, hooksDir string) string {
	preamble := globalHookPreamble(hooksDir)
	resolver := globalTotalityResolver(txPath)
	switch spec.tx {
	case "prepare-commit-msg":
		return preamble + fmt.Sprintf(`
%s
tx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tx_bin" ] && [ -n "$tx_repo" ]; then
  "$tx_bin" __hooks prepare-commit-msg --repo "$tx_repo" --message-path "$1" || {
    echo "tx prepare-commit-msg failed; commit continues without a Totality trailer" >&2
  }
fi
tx_resolve_local_hook prepare-commit-msg
if [ -n "$tx_local_hook" ]; then
  exec "$tx_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-commit":
		return preamble + fmt.Sprintf(`
%s
tx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tx_bin" ] && [ -n "$tx_repo" ]; then
  "$tx_bin" __hooks post-commit --repo "$tx_repo" || {
    echo "tx post-commit metadata recording failed" >&2
  }
fi
tx_resolve_local_hook post-commit
if [ -n "$tx_local_hook" ]; then
  exec "$tx_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-rewrite":
		return preamble + fmt.Sprintf(`
%s
tx_stdin="$(mktemp "${TMPDIR:-/tmp}/totality-post-rewrite.XXXXXX" 2>/dev/null)" || tx_stdin=""
if [ -n "$tx_stdin" ]; then
  trap 'rm -f "$tx_stdin"' EXIT
  cat > "$tx_stdin"
fi
tx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tx_bin" ] && [ -n "$tx_repo" ] && [ -n "$tx_stdin" ]; then
  "$tx_bin" __hooks post-rewrite --repo "$tx_repo" < "$tx_stdin" || {
    echo "tx post-rewrite metadata update failed" >&2
  }
fi
tx_resolve_local_hook post-rewrite
if [ -n "$tx_local_hook" ]; then
  if [ -n "$tx_stdin" ]; then
    "$tx_local_hook" "$@" < "$tx_stdin"
  else
    "$tx_local_hook" "$@" < /dev/null
  fi
  exit $?
fi
exit 0
`, resolver)
	case "pre-push":
		return preamble + fmt.Sprintf(`
%s
tx_remote="$1"
tx_url="$2"
tx_stdin="$(mktemp "${TMPDIR:-/tmp}/totality-pre-push.XXXXXX" 2>/dev/null)" || tx_stdin=""
if [ -n "$tx_stdin" ]; then
  trap 'rm -f "$tx_stdin"' EXIT
  cat > "$tx_stdin"
else
  # Always drain stdin so git never blocks writing refs to this hook.
  cat > /dev/null
fi
tx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tx_bin" ] && [ -n "$tx_repo" ] && [ -n "$tx_stdin" ]; then
  while read tx_local_ref tx_local_sha tx_remote_ref tx_remote_sha
  do
    if [ "$tx_local_sha" = "0000000000000000000000000000000000000000" ]; then
      continue
    fi
    if [ "$tx_remote_sha" = "0000000000000000000000000000000000000000" ]; then
      tx_range="$tx_local_sha"
    else
      tx_range="${tx_remote_sha}..${tx_local_sha}"
    fi
    "$tx_bin" capture push --remote "$tx_remote" --ref-range "$tx_range" --local-ref "$tx_local_ref" --head-sha "$tx_local_sha" --repo "$tx_repo" || true
  done < "$tx_stdin"
fi
tx_resolve_local_hook pre-push
if [ -n "$tx_local_hook" ]; then
  if [ -n "$tx_stdin" ]; then
    "$tx_local_hook" "$@" < "$tx_stdin"
  else
    "$tx_local_hook" "$@" < /dev/null
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
tx_resolve_local_hook %s
if [ -n "$tx_local_hook" ]; then
  exec "$tx_local_hook" "$@"
fi
%sexit 0
`, spec.name, drain)
}
