package hooks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
)

const (
	// globalHookMarker identifies scripts written by `gx init --global`. Every
	// global script also carries hookMarker so the existing GX-ownership checks
	// (gxOwnedHook / IsInstalled) keep recognizing them.
	globalHookMarker = "# gx global lifecycle hooks"

	// hooksPathConfigKey is the git config key a global install points at the
	// GX hooks directory.
	hooksPathConfigKey = "core.hooksPath"

	// RepoEnabledConfigKey is the per-repo opt-out. `git config gx.enabled false`
	// inside a repository stops GX lifecycle work there while leaving the
	// machine-wide install (and the repo's own hooks) intact.
	RepoEnabledConfigKey = "gx.enabled"
)

// globalHookSpec describes one script installed into the GX global hooks
// directory.
//
// A global core.hooksPath makes git ignore every repository's own
// .git/hooks/*, so GX has to shim more than the four hooks it cares about:
// each script here chains to the repository's own hook of the same name.
type globalHookSpec struct {
	name string
	// stdin marks hooks git feeds data on stdin. Those scripts must drain
	// stdin even when there is nothing to do, and must replay it to the
	// chained repo hook when GX consumed it first.
	stdin bool
	// gx names the GX lifecycle behavior this script runs before chaining.
	// Empty means the script is a pure forwarder.
	gx string
}

// globalHookSpecs lists every hook `gx init --global` installs.
//
// Deliberately excluded, because their mere existence changes what git does
// and a shim that exits 0 would silently replace the built-in behavior:
// push-to-checkout, proc-receive, and fsmonitor-watchman (the latter is
// addressed by core.fsmonitor path, not through core.hooksPath). git-p4 hooks
// (p4-*) are excluded as well; they only matter to git-p4 users.
func globalHookSpecs() []globalHookSpec {
	return []globalHookSpec{
		// GX lifecycle hooks: GX work first (never fatal), then the repo hook.
		{name: "prepare-commit-msg", gx: "prepare-commit-msg"},
		{name: "post-commit", gx: "post-commit"},
		{name: "post-rewrite", gx: "post-rewrite", stdin: true},
		{name: "pre-push", gx: "pre-push", stdin: true},

		// Pure forwarders: GX does nothing, but the repo's hook must still run.
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

// GlobalHooksDir returns the machine-wide GX hooks directory: $GX_HOME/hooks
// when GX_HOME is set, otherwise ~/.gx/hooks.
func GlobalHooksDir() (string, error) {
	dir, err := storage.DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "hooks"), nil
}

// GlobalInstallOptions configures a machine-wide hook install.
type GlobalInstallOptions struct {
	// HooksDir overrides the GX hooks directory. Defaults to GlobalHooksDir().
	HooksDir string
	// GXPath overrides the gx binary baked into the scripts.
	GXPath string
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
// core.hooksPath pointing at a GX directory with scripts missing.
func (s GlobalState) NeedsRepair() bool {
	if s.Enabled {
		return !s.Installed
	}
	return s.Installed
}

// InstallGlobal writes the chaining lifecycle hooks into the GX hooks
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

	gxPath := strings.TrimSpace(opts.GXPath)
	if gxPath == "" {
		gxPath, err = os.Executable()
		if err != nil {
			return result, fmt.Errorf("resolve gx binary: %w", err)
		}
	}

	current, err := globalHooksPath(ctx)
	if err != nil {
		return result, err
	}
	result.PreviousPath = current
	if current != "" && !sameHooksPath(current, hooksDir) && !gxOwnedHooksDir(current) {
		return result, fmt.Errorf(
			"git config --global %s is already set to %s, which is not managed by GX; "+
				"GX chains to each repository's own hooks but will not take over another global hooks directory. "+
				"Move those hooks into %s (they will be chained from there) or unset the config with "+
				"`git config --global --unset %s`, then re-run `gx init --global`",
			hooksPathConfigKey, current, hooksDir, hooksPathConfigKey)
	}

	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return result, err
	}
	for _, spec := range globalHookSpecs() {
		result.Scripts = append(result.Scripts, spec.name)
		script := globalHookScript(spec, gxPath, hooksDir)
		path := filepath.Join(hooksDir, spec.name)
		if data, readErr := os.ReadFile(path); readErr == nil {
			if string(data) == script {
				continue
			}
			if !gxOwnedHook(spec.name, string(data)) {
				return result, fmt.Errorf("refusing to overwrite non-GX script at %s; move it aside before retrying `gx init --global`", path)
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
	// Only call a foreign hooksPath a conflict once GX global hooks actually
	// exist; otherwise every husky/lefthook user would see a phantom problem.
	if !state.Enabled && current != "" && len(state.MissingScripts) < len(GlobalHookNames()) {
		state.Conflict = current
	}
	return state, nil
}

// EnabledForRepo reports whether GX lifecycle work should run in repoRoot.
//
// A repository opts out with `git config gx.enabled false`, which is how a
// user excludes one repo from a machine-wide `gx init --global` install. Unset
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

// gxOwnedHooksDir reports whether dir holds GX-written global hook scripts.
func gxOwnedHooksDir(dir string) bool {
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

// globalHookPreamble emits the shared header plus gx_resolve_local_hook, which
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
# Installed by `+"`gx init --global`"+`. A global core.hooksPath makes git ignore
# every repository's own .git/hooks, so each script here does GX's work (if any)
# and then runs the repository's own hook of the same name.
# Opt a repository out with: git config %[3]s false

gx_global_hooks_dir=%[4]s
gx_local_hook=""

gx_resolve_local_hook() {
  gx_local_hook=""
  gx_hook_name="$1"
  # Masking the global/system config keeps this from resolving to the GX dir.
  gx_dir="$(GIT_CONFIG_GLOBAL=/dev/null GIT_CONFIG_SYSTEM=/dev/null git config --get %[5]s 2>/dev/null)"
  if [ -z "$gx_dir" ]; then
    gx_dir="$(git rev-parse --path-format=absolute --git-common-dir 2>/dev/null)"
    [ -n "$gx_dir" ] || return 0
    gx_dir="$gx_dir/hooks"
  fi
  case "$gx_dir" in
    /*) ;;
    *) gx_dir="$(git rev-parse --show-toplevel 2>/dev/null)/$gx_dir" ;;
  esac
  # Never chain back into this directory: that would recurse forever.
  if [ "$gx_dir" = "$gx_global_hooks_dir" ]; then
    return 0
  fi
  case "$gx_dir" in
    %[4]s/*) return 0 ;;
  esac
  gx_candidate="$gx_dir/$gx_hook_name"
  [ -f "$gx_candidate" ] && [ -x "$gx_candidate" ] || return 0
  # A repo that ran plain `+"`gx init`"+` has GX's own hook here; running it would
  # repeat the work this script just did.
  if grep -q '%[1]s' "$gx_candidate" 2>/dev/null; then
    return 0
  fi
  gx_local_hook="$gx_candidate"
}
`, hookMarker, globalHookMarker, RepoEnabledConfigKey, quotedDir, hooksPathConfigKey)
}

// globalGXResolver emits shell that resolves gx at run time: the pinned
// install-time path when it is still an executable file, else gx from PATH.
// Unlike the per-repo hooks, a global script must never exit when gx is
// missing — it still has to chain to the repository's own hook — so gx_bin is
// left empty and the GX work is skipped instead.
func globalGXResolver(gxPath string) string {
	if strings.TrimSpace(gxPath) == "" {
		gxPath = "gx"
	}
	return fmt.Sprintf(`gx_bin=%s
case "$gx_bin" in /*) ;; *) gx_bin="" ;; esac
if [ ! -f "$gx_bin" ] || [ ! -x "$gx_bin" ]; then
  gx_bin="$(command -v gx 2>/dev/null)" || gx_bin=""
fi`, shellSingleQuote(gxPath))
}

// globalHookScript renders one global hook script.
func globalHookScript(spec globalHookSpec, gxPath, hooksDir string) string {
	preamble := globalHookPreamble(hooksDir)
	resolver := globalGXResolver(gxPath)
	switch spec.gx {
	case "prepare-commit-msg":
		return preamble + fmt.Sprintf(`
%s
gx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$gx_bin" ] && [ -n "$gx_repo" ]; then
  "$gx_bin" __hooks prepare-commit-msg --repo "$gx_repo" --message-path "$1" || {
    echo "gx prepare-commit-msg failed; commit continues without a GX trailer" >&2
  }
fi
gx_resolve_local_hook prepare-commit-msg
if [ -n "$gx_local_hook" ]; then
  exec "$gx_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-commit":
		return preamble + fmt.Sprintf(`
%s
gx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$gx_bin" ] && [ -n "$gx_repo" ]; then
  "$gx_bin" __hooks post-commit --repo "$gx_repo" || {
    echo "gx post-commit metadata recording failed" >&2
  }
fi
gx_resolve_local_hook post-commit
if [ -n "$gx_local_hook" ]; then
  exec "$gx_local_hook" "$@"
fi
exit 0
`, resolver)
	case "post-rewrite":
		return preamble + fmt.Sprintf(`
%s
gx_stdin="$(mktemp "${TMPDIR:-/tmp}/gx-post-rewrite.XXXXXX" 2>/dev/null)" || gx_stdin=""
if [ -n "$gx_stdin" ]; then
  trap 'rm -f "$gx_stdin"' EXIT
  cat > "$gx_stdin"
fi
gx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$gx_bin" ] && [ -n "$gx_repo" ] && [ -n "$gx_stdin" ]; then
  "$gx_bin" __hooks post-rewrite --repo "$gx_repo" < "$gx_stdin" || {
    echo "gx post-rewrite metadata update failed" >&2
  }
fi
gx_resolve_local_hook post-rewrite
if [ -n "$gx_local_hook" ]; then
  if [ -n "$gx_stdin" ]; then
    "$gx_local_hook" "$@" < "$gx_stdin"
  else
    "$gx_local_hook" "$@" < /dev/null
  fi
  exit $?
fi
exit 0
`, resolver)
	case "pre-push":
		return preamble + fmt.Sprintf(`
%s
gx_remote="$1"
gx_url="$2"
gx_stdin="$(mktemp "${TMPDIR:-/tmp}/gx-pre-push.XXXXXX" 2>/dev/null)" || gx_stdin=""
if [ -n "$gx_stdin" ]; then
  trap 'rm -f "$gx_stdin"' EXIT
  cat > "$gx_stdin"
else
  # Always drain stdin so git never blocks writing refs to this hook.
  cat > /dev/null
fi
gx_repo="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -n "$gx_bin" ] && [ -n "$gx_repo" ] && [ -n "$gx_stdin" ]; then
  while read gx_local_ref gx_local_sha gx_remote_ref gx_remote_sha
  do
    if [ "$gx_local_sha" = "0000000000000000000000000000000000000000" ]; then
      continue
    fi
    if [ "$gx_remote_sha" = "0000000000000000000000000000000000000000" ]; then
      gx_range="$gx_local_sha"
    else
      gx_range="${gx_remote_sha}..${gx_local_sha}"
    fi
    "$gx_bin" capture push --remote "$gx_remote" --ref-range "$gx_range" --local-ref "$gx_local_ref" --head-sha "$gx_local_sha" --repo "$gx_repo" || true
  done < "$gx_stdin"
fi
gx_resolve_local_hook pre-push
if [ -n "$gx_local_hook" ]; then
  if [ -n "$gx_stdin" ]; then
    "$gx_local_hook" "$@" < "$gx_stdin"
  else
    "$gx_local_hook" "$@" < /dev/null
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
gx_resolve_local_hook %s
if [ -n "$gx_local_hook" ]; then
  exec "$gx_local_hook" "$@"
fi
%sexit 0
`, spec.name, drain)
}
