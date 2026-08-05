package hooks_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/lgtm/internal/hooks"
)

// isolateGlobalGitConfig points HOME, LGTM_HOME and git's global/system config
// files at a throwaway directory so these tests can exercise
// `git config --global` without ever touching the developer's real config.
func isolateGlobalGitConfig(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	configPath := filepath.Join(home, "gitconfig")
	if err := os.WriteFile(configPath, []byte("[lgtm]\n\tisolationprobe = temp\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("LGTM_HOME", filepath.Join(home, "lgtm"))
	t.Setenv("GIT_CONFIG_GLOBAL", configPath)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	// Refuse to run if this git ignores GIT_CONFIG_GLOBAL (added in git 2.32):
	// every write below would otherwise land in the real ~/.gitconfig.
	out, err := exec.Command("git", "config", "--global", "--get", "lgtm.isolationprobe").Output()
	if err != nil || strings.TrimSpace(string(out)) != "temp" {
		t.Skipf("git does not honor GIT_CONFIG_GLOBAL, refusing to touch the real global config: %v", err)
	}
	return configPath
}

func globalHooksPathValue(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "config", "--global", "--get", "core.hooksPath").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func TestInstallGlobalWritesChainingScripts(t *testing.T) {
	isolateGlobalGitConfig(t)
	hooksDir := filepath.Join(t.TempDir(), "hooks")

	result, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{
		HooksDir:     hooksDir,
		LgtmPath: "/usr/local/bin/lgtm",
	})
	if err != nil {
		t.Fatalf("InstallGlobal() error = %v", err)
	}
	if !result.ConfigUpdated {
		t.Fatal("InstallGlobal() did not set core.hooksPath")
	}
	if got := globalHooksPathValue(t); got != hooksDir {
		t.Fatalf("core.hooksPath = %q, want %q", got, hooksDir)
	}

	// Every hook git may run must exist, otherwise a global core.hooksPath
	// silently disables the repository's own hook of that name.
	for _, name := range []string{
		"prepare-commit-msg", "post-commit", "post-rewrite", "pre-push",
		"pre-commit", "commit-msg", "pre-merge-commit", "post-checkout",
		"post-merge", "pre-rebase", "applypatch-msg", "reference-transaction",
	} {
		if _, err := os.Stat(filepath.Join(hooksDir, name)); err != nil {
			t.Fatalf("missing global hook %s: %v", name, err)
		}
	}

	for _, tc := range []struct {
		name  string
		want  []string
		avoid []string
	}{
		{
			name: "prepare-commit-msg",
			want: []string{
				"# lgtm lifecycle hooks",
				"# lgtm global lifecycle hooks",
				"__hooks prepare-commit-msg",
				"command -v lgtm",
				"lgtm_resolve_local_hook prepare-commit-msg",
				`exec "$lgtm_local_hook" "$@"`,
			},
		},
		{
			name: "post-commit",
			want: []string{
				"__hooks post-commit",
				"lgtm_resolve_local_hook post-commit",
				`exec "$lgtm_local_hook" "$@"`,
			},
		},
		{
			name: "post-rewrite",
			want: []string{
				"__hooks post-rewrite",
				`cat > "$lgtm_stdin"`,
				"lgtm_resolve_local_hook post-rewrite",
				`"$lgtm_local_hook" "$@" < "$lgtm_stdin"`,
				"exit $?",
			},
		},
		{
			name: "pre-push",
			want: []string{
				"capture push",
				"command -v lgtm",
				`cat > "$lgtm_stdin"`,
				`done < "$lgtm_stdin"`,
				"lgtm_resolve_local_hook pre-push",
				`"$lgtm_local_hook" "$@" < "$lgtm_stdin"`,
				"exit $?",
			},
		},
		{
			name: "pre-commit",
			want: []string{
				"lgtm_resolve_local_hook pre-commit",
				`exec "$lgtm_local_hook" "$@"`,
			},
			// Pure forwarders must not run lgtm work or read stdin git never sends.
			avoid: []string{"__hooks", "capture push", "cat > /dev/null"},
		},
		{
			name: "reference-transaction",
			want: []string{
				"lgtm_resolve_local_hook reference-transaction",
				`exec "$lgtm_local_hook" "$@"`,
				"cat > /dev/null",
			},
			avoid: []string{"__hooks", "capture push"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(hooksDir, tc.name))
			if err != nil {
				t.Fatal(err)
			}
			content := string(data)
			for _, want := range tc.want {
				if !strings.Contains(content, want) {
					t.Fatalf("%s hook missing %q:\n%s", tc.name, want, content)
				}
			}
			for _, avoid := range tc.avoid {
				if strings.Contains(content, avoid) {
					t.Fatalf("%s hook should not contain %q:\n%s", tc.name, avoid, content)
				}
			}
			// The resolver must never use `git rev-parse --git-path hooks/...`:
			// that honors core.hooksPath and resolves back to this script.
			if strings.Contains(content, "--git-path hooks") {
				t.Fatalf("%s hook resolves repo hooks through core.hooksPath:\n%s", tc.name, content)
			}
		})
	}
}

func TestGlobalHookScriptsAreValidShell(t *testing.T) {
	isolateGlobalGitConfig(t)
	hooksDir := filepath.Join(t.TempDir(), "hooks")
	if _, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{
		HooksDir:     hooksDir,
		LgtmPath: "/usr/local/bin/lgtm",
	}); err != nil {
		t.Fatal(err)
	}
	for _, name := range hooks.GlobalHookNames() {
		path := filepath.Join(hooksDir, name)
		if out, err := exec.Command("sh", "-n", path).CombinedOutput(); err != nil {
			t.Fatalf("sh -n %s: %v\n%s", name, err, out)
		}
	}
}

func TestInstallGlobalRefusesForeignHooksPath(t *testing.T) {
	isolateGlobalGitConfig(t)
	foreign := t.TempDir()
	if out, err := exec.Command("git", "config", "--global", "core.hooksPath", foreign).CombinedOutput(); err != nil {
		t.Fatalf("seed core.hooksPath: %v\n%s", err, out)
	}
	hooksDir := filepath.Join(t.TempDir(), "hooks")

	_, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{HooksDir: hooksDir, LgtmPath: "/bin/lgtm"})
	if err == nil {
		t.Fatal("InstallGlobal() error = nil, want refusal")
	}
	for _, want := range []string{"not managed by lgtm", foreign} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("InstallGlobal() error = %q, want mention of %q", err, want)
		}
	}
	if got := globalHooksPathValue(t); got != foreign {
		t.Fatalf("core.hooksPath = %q, want it left at %q", got, foreign)
	}
	if _, statErr := os.Stat(filepath.Join(hooksDir, "pre-push")); statErr == nil {
		t.Fatal("InstallGlobal() wrote scripts despite refusing")
	}
}

func TestInstallGlobalIsIdempotent(t *testing.T) {
	isolateGlobalGitConfig(t)
	hooksDir := filepath.Join(t.TempDir(), "hooks")
	opts := hooks.GlobalInstallOptions{HooksDir: hooksDir, LgtmPath: "/usr/local/bin/lgtm"}

	first, err := hooks.InstallGlobal(t.Context(), opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Written) != len(first.Scripts) {
		t.Fatalf("first install wrote %d of %d scripts", len(first.Written), len(first.Scripts))
	}

	second, err := hooks.InstallGlobal(t.Context(), opts)
	if err != nil {
		t.Fatalf("second InstallGlobal() error = %v", err)
	}
	if second.ConfigUpdated {
		t.Fatal("second InstallGlobal() rewrote core.hooksPath")
	}
	if len(second.Written) != 0 {
		t.Fatalf("second InstallGlobal() rewrote scripts: %v", second.Written)
	}
	if got := globalHooksPathValue(t); got != hooksDir {
		t.Fatalf("core.hooksPath = %q, want %q", got, hooksDir)
	}
}

func TestGlobalStatusTracksInstallAndConflict(t *testing.T) {
	isolateGlobalGitConfig(t)
	dir, err := hooks.GlobalHooksDir()
	if err != nil {
		t.Fatal(err)
	}

	before, err := hooks.GlobalStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if before.Enabled || before.Installed || before.Conflict != "" || before.NeedsRepair() {
		t.Fatalf("GlobalStatus() before install = %+v, want empty state", before)
	}

	if _, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{LgtmPath: "/usr/local/bin/lgtm"}); err != nil {
		t.Fatal(err)
	}
	after, err := hooks.GlobalStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if after.HooksDir != dir || !after.Enabled || !after.Installed || after.NeedsRepair() {
		t.Fatalf("GlobalStatus() after install = %+v, want enabled and installed at %s", after, dir)
	}

	foreign := t.TempDir()
	if out, err := exec.Command("git", "config", "--global", "core.hooksPath", foreign).CombinedOutput(); err != nil {
		t.Fatalf("override core.hooksPath: %v\n%s", err, out)
	}
	conflicted, err := hooks.GlobalStatus(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if conflicted.Enabled || conflicted.Conflict != foreign || !conflicted.NeedsRepair() {
		t.Fatalf("GlobalStatus() with foreign hooksPath = %+v, want conflict %s", conflicted, foreign)
	}
}

func TestGlobalHookRunsRepoHookAndPropagatesExit(t *testing.T) {
	isolateGlobalGitConfig(t)
	hooksDir := filepath.Join(t.TempDir(), "hooks")
	if _, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{HooksDir: hooksDir, LgtmPath: "/bin/false"}); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	writeExecutable(t, filepath.Join(repo, ".git", "hooks", "pre-commit"), `#!/bin/sh
echo "repo pre-commit ran: $*"
exit 3
`)

	out, err := runHookScript(t, filepath.Join(hooksDir, "pre-commit"), repo, "", "--arg")
	if exitCode(t, err) != 3 {
		t.Fatalf("global pre-commit exit = %v (want 3), output:\n%s", err, out)
	}
	if !strings.Contains(out, "repo pre-commit ran: --arg") {
		t.Fatalf("repo hook did not run with forwarded args:\n%s", out)
	}
}

func TestGlobalHookSkipsLgtmOwnedRepoHook(t *testing.T) {
	isolateGlobalGitConfig(t)
	hooksDir := filepath.Join(t.TempDir(), "hooks")
	if _, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{HooksDir: hooksDir, LgtmPath: "/bin/true"}); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	// A repo that ran plain `lgtm init` keeps its own lgtm hooks; the global script
	// must not chain into them or lgtm work would run twice.
	writeExecutable(t, filepath.Join(repo, ".git", "hooks", "pre-commit"), `#!/bin/sh
# lgtm lifecycle hooks
echo "lgtm owned repo hook ran"
exit 5
`)

	out, err := runHookScript(t, filepath.Join(hooksDir, "pre-commit"), repo, "")
	if err != nil {
		t.Fatalf("global pre-commit error = %v, output:\n%s", err, out)
	}
	if strings.Contains(out, "lgtm owned repo hook ran") {
		t.Fatalf("global hook chained into a lgtm-owned repo hook:\n%s", out)
	}
}

func TestGlobalPrePushHookFeedsLgtmAndRepoHookStdin(t *testing.T) {
	isolateGlobalGitConfig(t)
	work := t.TempDir()
	lgtmLog := filepath.Join(work, "lgtm.log")
	fakeLgtm := filepath.Join(work, "lgtm")
	writeExecutable(t, fakeLgtm, "#!/bin/sh\nprintf '%s\\n' \"$*\" >> \""+lgtmLog+"\"\n")

	hooksDir := filepath.Join(work, "hooks")
	if _, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{HooksDir: hooksDir, LgtmPath: fakeLgtm}); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	repoLog := filepath.Join(work, "repo-pre-push.log")
	writeExecutable(t, filepath.Join(repo, ".git", "hooks", "pre-push"), `#!/bin/sh
{ echo "args: $*"; cat; } >> "`+repoLog+`"
exit 9
`)

	stdin := "refs/heads/main aaaa111 refs/heads/main bbbb222\n"
	out, err := runHookScript(t, filepath.Join(hooksDir, "pre-push"), repo, stdin, "origin", "git@example.com:acme/app.git")
	if exitCode(t, err) != 9 {
		t.Fatalf("global pre-push exit = %v (want 9), output:\n%s", err, out)
	}

	lgtmArgs := readFile(t, lgtmLog)
	for _, want := range []string{
		"capture push",
		"--remote origin",
		"--ref-range bbbb222..aaaa111",
		"--local-ref refs/heads/main",
		"--head-sha aaaa111",
	} {
		if !strings.Contains(lgtmArgs, want) {
			t.Fatalf("lgtm capture push args missing %q:\n%s", want, lgtmArgs)
		}
	}

	repoStdin := readFile(t, repoLog)
	if !strings.Contains(repoStdin, "args: origin git@example.com:acme/app.git") {
		t.Fatalf("repo pre-push hook missing forwarded args:\n%s", repoStdin)
	}
	// lgtm consumed stdin first; the repo hook must still receive every ref.
	if !strings.Contains(repoStdin, "refs/heads/main aaaa111 refs/heads/main bbbb222") {
		t.Fatalf("repo pre-push hook did not receive replayed stdin:\n%s", repoStdin)
	}
}

// TestGlobalPrePushChainsWhenLgtmIsMissing pins the failure mode that matters
// most under core.hooksPath: when lgtm cannot be resolved at all, the global
// script must skip lgtm work silently but still replay stdin to the
// repository's own hook — exiting early would disable every repo hook on the
// machine.
func TestGlobalPrePushChainsWhenLgtmIsMissing(t *testing.T) {
	isolateGlobalGitConfig(t)
	work := t.TempDir()
	hooksDir := filepath.Join(work, "hooks")
	if _, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{
		HooksDir:     hooksDir,
		LgtmPath: filepath.Join(work, "missing", "lgtm"),
	}); err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	repoLog := filepath.Join(work, "repo-pre-push.log")
	writeExecutable(t, filepath.Join(repo, ".git", "hooks", "pre-push"), `#!/bin/sh
{ echo "args: $*"; cat; } >> "`+repoLog+`"
exit 7
`)

	// PATH holds git and the shell utilities the script needs, but no lgtm.
	pathDir := filepath.Join(work, "pathbin")
	for _, tool := range []string{"git", "mktemp", "cat", "rm", "grep"} {
		resolved, err := exec.LookPath(tool)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(pathDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(resolved, filepath.Join(pathDir, tool)); err != nil {
			t.Fatal(err)
		}
	}
	stdin := "refs/heads/main aaaa111 refs/heads/main bbbb222\n"
	cmd := exec.Command(filepath.Join(hooksDir, "pre-push"), "origin", "git@example.com:acme/app.git")
	cmd.Dir = repo
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Env = append(os.Environ(), "PATH="+pathDir)
	out, err := cmd.CombinedOutput()
	if exitCode(t, err) != 7 {
		t.Fatalf("exit = %v (want repo hook's 7), output:\n%s", err, out)
	}
	repoStdin := readFile(t, repoLog)
	if !strings.Contains(repoStdin, "refs/heads/main aaaa111 refs/heads/main bbbb222") {
		t.Fatalf("repo pre-push hook did not receive replayed stdin:\n%s", repoStdin)
	}
	if strings.Contains(string(out), "lgtm") {
		t.Fatalf("missing lgtm must be silent, got:\n%s", out)
	}
}

func TestPerRepoInstallUnaffectedByGlobalSupport(t *testing.T) {
	isolateGlobalGitConfig(t)
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, LgtmPath: "/usr/local/bin/lgtm"}); err != nil {
		t.Fatal(err)
	}
	if !hooks.IsInstalled(repo) {
		t.Fatal("per-repo install did not register lifecycle hooks")
	}
	for _, name := range []string{"prepare-commit-msg", "post-commit", "post-rewrite", "pre-push"} {
		content := readFile(t, filepath.Join(repo, ".git", "hooks", name))
		if !strings.Contains(content, "# lgtm lifecycle hooks") {
			t.Fatalf("per-repo %s hook missing marker:\n%s", name, content)
		}
		if strings.Contains(content, "# lgtm global lifecycle hooks") {
			t.Fatalf("per-repo %s hook picked up global chaining script:\n%s", name, content)
		}
	}
}

func TestPerRepoInstallDoesNotClobberGlobalHooks(t *testing.T) {
	isolateGlobalGitConfig(t)
	if _, err := hooks.InstallGlobal(t.Context(), hooks.GlobalInstallOptions{LgtmPath: "/usr/local/bin/lgtm"}); err != nil {
		t.Fatal(err)
	}
	globalDir, err := hooks.GlobalHooksDir()
	if err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")

	// git now resolves this repo's hooks dir to the shared lgtm directory; a
	// per-repo install must not replace the chaining scripts there.
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, LgtmPath: "/usr/local/bin/lgtm"}); err != nil {
		t.Fatalf("Install() error = %v", err)
	}
	for _, name := range []string{"pre-push", "post-commit"} {
		content := readFile(t, filepath.Join(globalDir, name))
		if !strings.Contains(content, "# lgtm global lifecycle hooks") {
			t.Fatalf("global %s hook was overwritten by a per-repo install:\n%s", name, content)
		}
		if !strings.Contains(content, "lgtm_resolve_local_hook "+name) {
			t.Fatalf("global %s hook lost its chaining logic:\n%s", name, content)
		}
	}
	if _, err := os.Stat(filepath.Join(repo, ".git", "hooks", "pre-push")); err == nil {
		t.Fatal("per-repo install wrote hooks that git would ignore under a global install")
	}
}

func TestEnabledForRepoHonorsOptOut(t *testing.T) {
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	for _, tc := range []struct {
		name  string
		value string
		want  bool
	}{
		{name: "unset", value: "", want: true},
		{name: "false", value: "false", want: false},
		{name: "true", value: "true", want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.value == "" {
				_ = exec.Command("git", "-C", repo, "config", "--unset", "lgtm.enabled").Run()
			} else {
				runGitInRepo(t, repo, "config", "lgtm.enabled", tc.value)
			}
			if got := hooks.EnabledForRepo(t.Context(), repo); got != tc.want {
				t.Fatalf("EnabledForRepo(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}

func TestPrepareCommitMsgRespectsRepoOptOut(t *testing.T) {
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	runGitInRepo(t, repo, "config", "lgtm.enabled", "false")
	messagePath := filepath.Join(repo, "COMMIT_EDITMSG")
	if err := os.WriteFile(messagePath, []byte("add feature\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := hooks.PrepareCommitMsg(hooks.PrepareCommitMsgOptions{RepoRoot: repo, MessagePath: messagePath}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, messagePath); strings.Contains(got, "lgtm: https://lgtm.cx/r/") {
		t.Fatalf("opted-out repo still got a lgtm trailer: %q", got)
	}
}

func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func runHookScript(t *testing.T, script, repo, stdin string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(script, args...)
	cmd.Dir = repo
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func exitCode(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("unexpected hook error: %v", err)
	}
	return exitErr.ExitCode()
}
