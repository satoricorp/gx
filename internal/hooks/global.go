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
	// globalHookMarker identifies scripts written by `tl init --global`. Every
	// global script also carries hookMarker so the existing Totality-ownership checks
	// (totalityOwnedHook / IsInstalled) keep recognizing them.
	globalHookMarker = "# tl global lifecycle hooks"

	// hooksPathConfigKey is the git config key a global install points at the
	// Totality hooks directory.
	hooksPathConfigKey = "core.hooksPath"

	// RepoEnabledConfigKey is the per-repo opt-out. `git config tl.enabled false`
	// inside a repository stops Totality lifecycle work there while leaving the
	// machine-wide install (and the repo's own hooks) intact.
	RepoEnabledConfigKey = "tl.enabled"
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
	// tl names the Totality lifecycle behavior this script runs before chaining.
	// Empty means the script is a pure forwarder.
	tl string
}

// globalHookSpecs lists every hook `tl init --global` installs.
//
// Deliberately excluded, because their mere existence changes what git does
// and a shim that exits 0 would silently replace the built-in behavior:
// push-to-checkout, proc-receive, and fsmonitor-watchman (the latter is
// addressed by core.fsmonitor path, not through core.hooksPath). git-p4 hooks
// (p4-*) are excluded as well; they only matter to git-p4 users.
func globalHookSpecs() []globalHookSpec {
	return []globalHookSpec{
		// Totality lifecycle hooks: Totality work first (never fatal), then the repo hook.
		{name: "prepare-commit-msg", tl: "prepare-commit-msg"},
		{name: "post-commit", tl: "post-commit"},
		{name: "post-rewrite", tl: "post-rewrite", stdin: true},
		{name: "pre-push", tl: "pre-push", stdin: true},

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
	// TotalityPath overrides the tl binary baked into the scripts.
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

	tlPath := strings.TrimSpace(opts.TotalityPath)
	if tlPath == "" {
		tlPath, err = installTotalityPath()
		if err != nil {
			return result, fmt.Errorf("resolve tl binary: %w", err)
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
				"`git config --global --unset %s`, then re-run `tl init --global`",
			hooksPathConfigKey, current, hooksDir, hooksPathConfigKey)
	}

	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return result, err
	}
	for _, spec := range globalHookSpecs() {
		result.Scripts = append(result.Scripts, spec.name)
		script := globalHookScript(spec, tlPath, hooksDir)
		path := filepath.Join(hooksDir, spec.name)
		if data, readErr := os.ReadFile(path); readErr == nil {
			if string(data) == script {
				continue
			}
			if !totalityOwnedHook(spec.name, string(data)) {
				return result, fmt.Errorf("refusing to overwrite non-Totality script at %s; move it aside before retrying `tl init --global`", path)
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
// A repository opts out with `git config tl.enabled false`, which is how a
// user excludes one repo from a machine-wide `tl init --global` install. Unset
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

// globalHookPreamble emits the shared header plus tl_resolve_local_hook, which
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
# Installed by `+"`tl init --global`"+`. A global core.hooksPath makes git ignore
# every repository's own .git/hooks, so each script here does Totality's work (if any)
# and then runs the repository's own hook of the same name.
# Opt a repository out with: git config %[3]s false

tl_global_hooks_dir=%[4]s
tl_local_hook=""

tl_resolve_local_hook() {
  tl_local_hook=""
  tl_hook_name="$1"
  # Masking the global/system config keeps this from resolving to the Totality dir.
  tl_dir="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null git config --get %[5]s 2>/dev/null)"
  if [ -z "$tl_dir" ]; then
    tl_dir="$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null)"
    [ -n "$tl_dir" ] || return 0
    tl_dir="$tl_dir/hooks"
  fi
  case "$tl_dir" in
    /*) ;;
    *) tl_dir="$(git rev-parse --show-toplevel 2>/dev/null)/$tl_dir" ;;
  esac
  # Never chain back into this directory: that would recurse forever.
  if [ "$tl_dir" = "$tl_global_hooks_dir" ]; then
    return 0
  fi
  case "$tl_dir" in
    %[4]s/*) return 0 ;;
  esac
  tl_candidate="$tl_dir/$tl_hook_name"
  [ -f "$tl_candidate" ] && [ -x "$tl_candidate" ] || return 0
  # A repo that ran plain `+"`tl init`"+` has Totality's own hook here; running it would
  # repeat the work this script just did.
  if grep -q '%[1]s' "$tl_candidate" 2>/dev/null; then
    return 0
  fi
  tl_local_hook="$tl_candidate"
}
`, hookMarker, globalHookMarker, RepoEnabledConfigKey, quotedDir, hooksPathConfigKey)
}

// globalTotalityResolver emits shell that resolves tl at run time: the pinned
// install-time path when it is still an executable file, else tl from PATH.
// Unlike the per-repo hooks, a global script must never exit when tl is
// missing — it still has to chain to the repository's own hook — so tl_bin is
// left empty and the Totality work is skipped instead.
func globalTotalityResolver(tlPath string) string {
	if strings.TrimSpace(tlPath) == "" {
		tlPath = "tl"
	}
	return fmt.Sprintf(`tl_bin=%s
case "$tl_bin" in /*) ;; *) tl_bin="" ;; esac
if [ ! -f "$tl_bin" ] || [ ! -x "$tl_bin" ]; then
  tl_bin="$(command -v tl 2>/dev/null)" || tl_bin=""
fi`, shellSingleQuote(tlPath))
}

// globalHookScript renders one global hook script.
func globalHookScript(spec globalHookSpec, tlPath, hooksDir string) string {
	preamble := globalHookPreamble(hooksDir)
	resolver := globalTotalityResolver(tlPath)
	switch spec.tl {
	case "prepare-commit-msg":
		return preamble + fmt.Sprintf(`
%s
tl_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tl_bin" ] && [ -n "$tl_repo" ]; then
  "$tl_bin" __hooks prepare-commit-msg --repo "$tl_repo" --message-path "$1" || {
    echo "tl prepare-commit-msg failed; commit continues without a Totality trailer" >&2
  }
fi
tl_resolve_local_hook prepare-commit-msg
if [ -n "$tl_local_hook" ]; then
  exec "$tl_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-commit":
		return preamble + fmt.Sprintf(`
%s
tl_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tl_bin" ] && [ -n "$tl_repo" ]; then
  "$tl_bin" __hooks post-commit --repo "$tl_repo" || {
    echo "tl post-commit metadata recording failed" >&2
  }
fi
tl_resolve_local_hook post-commit
if [ -n "$tl_local_hook" ]; then
  exec "$tl_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-rewrite":
		return preamble + fmt.Sprintf(`
%s
tl_stdin="$(mktemp "${TMPDIR:-/tmp}/totality-post-rewrite.XXXXXX" 2>/dev/null)" || tl_stdin=""
if [ -n "$tl_stdin" ]; then
  trap 'rm -f "$tl_stdin"' EXIT
  cat > "$tl_stdin"
fi
tl_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tl_bin" ] && [ -n "$tl_repo" ] && [ -n "$tl_stdin" ]; then
  "$tl_bin" __hooks post-rewrite --repo "$tl_repo" < "$tl_stdin" || {
    echo "tl post-rewrite metadata update failed" >&2
  }
fi
tl_resolve_local_hook post-rewrite
if [ -n "$tl_local_hook" ]; then
  if [ -n "$tl_stdin" ]; then
    "$tl_local_hook" "$@" < "$tl_stdin"
  else
    "$tl_local_hook" "$@" < /dev/null
  fi
  exit $?
fi
exit 0
`, resolver)
	case "pre-push":
		return preamble + fmt.Sprintf(`
%s
tl_remote="$1"
tl_url="$2"
tl_stdin="$(mktemp "${TMPDIR:-/tmp}/totality-pre-push.XXXXXX" 2>/dev/null)" || tl_stdin=""
if [ -n "$tl_stdin" ]; then
  trap 'rm -f "$tl_stdin"' EXIT
  cat > "$tl_stdin"
else
  # Always drain stdin so git never blocks writing refs to this hook.
  cat > /dev/null
fi
tl_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$tl_bin" ] && [ -n "$tl_repo" ] && [ -n "$tl_stdin" ]; then
  while read tl_local_ref tl_local_sha tl_remote_ref tl_remote_sha
  do
    if [ "$tl_local_sha" = "0000000000000000000000000000000000000000" ]; then
      continue
    fi
    if [ "$tl_remote_sha" = "0000000000000000000000000000000000000000" ]; then
      tl_range="$tl_local_sha"
    else
      tl_range="${tl_remote_sha}..${tl_local_sha}"
    fi
    "$tl_bin" capture push --remote "$tl_remote" --ref-range "$tl_range" --local-ref "$tl_local_ref" --head-sha "$tl_local_sha" --repo "$tl_repo" || true
  done < "$tl_stdin"
fi
tl_resolve_local_hook pre-push
if [ -n "$tl_local_hook" ]; then
  if [ -n "$tl_stdin" ]; then
    "$tl_local_hook" "$@" < "$tl_stdin"
  else
    "$tl_local_hook" "$@" < /dev/null
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
tl_resolve_local_hook %s
if [ -n "$tl_local_hook" ]; then
  exec "$tl_local_hook" "$@"
fi
%sexit 0
`, spec.name, drain)
}
