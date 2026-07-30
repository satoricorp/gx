package hooks_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/totality/internal/hooks"
	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/vcs"
)

func TestInstallLifecycleHooksUsesGitHooksDir(t *testing.T) {
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	runGitInRepo(t, repo, "config", "user.email", "dev@example.com")
	runGitInRepo(t, repo, "config", "user.name", "Dev")
	txPath := buildTotalityBinary(t)
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, TotalityPath: txPath}); err != nil {
		t.Fatal(err)
	}
	hooksDir, err := hooks.ResolveHooksDir(t.Context(), repo)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"prepare-commit-msg", "post-commit", "post-rewrite", "pre-push"} {
		data, err := os.ReadFile(filepath.Join(hooksDir, name))
		if err != nil {
			t.Fatalf("read %s hook: %v", name, err)
		}
		if !strings.Contains(string(data), "# tx lifecycle hooks") {
			t.Fatalf("%s hook missing marker:\n%s", name, data)
		}
	}
}

func TestPrepareCommitMsgHook(t *testing.T) {
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	messagePath := filepath.Join(repo, "COMMIT_EDITMSG")
	if err := os.WriteFile(messagePath, []byte("add feature\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := hooks.PrepareCommitMsg(hooks.PrepareCommitMsgOptions{
		RepoRoot:    repo,
		MessagePath: messagePath,
	}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(messagePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Totality: https://totality.sh/r/") {
		t.Fatalf("message = %q, want Totality trailer", string(data))
	}
}

func TestPlainGitCommitWithInstalledHooks(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("TOTALITY_HOME", t.TempDir())
	runGitInRepo(t, repo, "init", "-b", "main")
	runGitInRepo(t, repo, "config", "user.email", "dev@example.com")
	runGitInRepo(t, repo, "config", "user.name", "Dev")
	txPath := buildTotalityBinary(t)
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, TotalityPath: txPath}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "feature.txt"), []byte("feature\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInRepo(t, repo, "add", "feature.txt")
	runGitInRepo(t, repo, "commit", "-m", "add feature")
	message := gitOutput(t, repo, "log", "-1", "--format=%B")
	if strings.Count(message, "Totality: https://totality.sh/r/") != 1 {
		t.Fatalf("commit message = %q, want exactly one Totality trailer", message)
	}
}

func TestInstallRefusesForeignHook(t *testing.T) {
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	hooksDir, err := hooks.ResolveHooksDir(t.Context(), repo)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(hooksDir, "post-commit")
	if err := os.WriteFile(path, []byte("#!/bin/sh\necho foreign\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	err = hooks.Install(hooks.InstallOptions{RepoRoot: repo, TotalityPath: "/bin/tx"})
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("Install() error = %v, want refusal", err)
	}
}

func TestInstallUpgradesLegacyTotalityPrePushHook(t *testing.T) {
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	hooksDir, err := hooks.ResolveHooksDir(t.Context(), repo)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(hooksDir, "pre-push")
	legacy := "#!/bin/sh\n# tx capture pre-push hook\ntt capture push || true\n"
	if err := os.WriteFile(path, []byte(legacy), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, TotalityPath: "/bin/tx"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "# tx lifecycle hooks") || strings.Contains(string(data), "# tx capture pre-push hook") {
		t.Fatalf("legacy pre-push hook was not upgraded:\n%s", data)
	}
}

func TestPostLifecycleHooksWarnWithoutBlocking(t *testing.T) {
	repo := t.TempDir()
	runGitInRepo(t, repo, "init")
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, TotalityPath: "/missing/tx"}); err != nil {
		t.Fatal(err)
	}
	hooksDir, err := hooks.ResolveHooksDir(t.Context(), repo)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"post-commit", "post-rewrite"} {
		data, err := os.ReadFile(filepath.Join(hooksDir, name))
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		if !strings.Contains(content, "failed") || !strings.Contains(content, "exit 0") {
			t.Fatalf("%s does not warn and continue:\n%s", name, content)
		}
	}
}

func TestBootstrapFromLinkedWorktree(t *testing.T) {
	primary := t.TempDir()
	totalityHome := t.TempDir()
	t.Setenv("TOTALITY_HOME", totalityHome)
	t.Setenv("TOTALITY_DISABLE_BACKGROUND_WORKERS", "1")
	runGitInRepo(t, primary, "init", "-b", "main")
	runGitInRepo(t, primary, "config", "user.email", "dev@example.com")
	runGitInRepo(t, primary, "config", "user.name", "Dev")
	if err := os.WriteFile(filepath.Join(primary, "README.md"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInRepo(t, primary, "add", "README.md")
	runGitInRepo(t, primary, "commit", "-m", "base")

	linked := filepath.Join(t.TempDir(), "linked")
	runGitInRepo(t, primary, "worktree", "add", "-b", "linked", linked)
	if err := os.WriteFile(filepath.Join(primary, "primary.txt"), []byte("primary staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInRepo(t, primary, "add", "primary.txt")
	primaryHead := strings.TrimSpace(gitOutput(t, primary, "rev-parse", "HEAD"))
	primaryIndex := gitOutput(t, primary, "diff", "--cached", "--binary")

	if err := os.WriteFile(filepath.Join(linked, "linked.txt"), []byte("linked\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInRepo(t, linked, "add", "linked.txt")
	// Any tx command run inside the linked worktree bootstraps the repo and
	// installs the lifecycle hooks; the commit itself is plain git.
	txPath := buildTotalityBinary(t)
	cmd := exec.Command(txPath, "doctor")
	cmd.Dir = linked
	cmd.Env = append(os.Environ(), "TOTALITY_HOME="+totalityHome, "TOTALITY_DISABLE_BACKGROUND_WORKERS=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("tx doctor: %v\n%s", err, out)
	}

	if !hooks.IsInstalled(linked) {
		t.Fatal("bootstrap did not install lifecycle hooks")
	}
	runGitInRepo(t, linked, "commit", "-m", "linked worktree commit")
	linkedMessage := gitOutput(t, linked, "log", "-1", "--format=%B")
	revisionIDs := vcs.ParseRevisionIDsFromMessage(linkedMessage)
	if len(revisionIDs) != 1 {
		t.Fatalf("linked commit message has revision ids %v", revisionIDs)
	}
	linkedOID := strings.TrimSpace(gitOutput(t, linked, "rev-parse", "HEAD"))

	service := vcs.NewService()
	repo, err := service.ResolveTotalityRepoAtPath(context.Background(), linked)
	if err != nil {
		t.Fatal(err)
	}
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	repoRow, err := store.FindRepoByIdentity(context.Background(), repo.GitCommonDir, repo.RootPath)
	if err != nil || repoRow == nil {
		t.Fatalf("FindRepoByIdentity() = %+v, %v", repoRow, err)
	}
	initialized, err := store.IsInitializedRepo(context.Background(), repoRow.RootPath)
	if err != nil || !initialized {
		t.Fatalf("IsInitializedRepo() = %v, %v", initialized, err)
	}
	change, err := store.FindChangeByCommitID(context.Background(), repoRow.ID, linkedOID)
	if err != nil || change == nil {
		t.Fatalf("FindChangeByCommitID() = %+v, %v", change, err)
	}
	if change.JJChangeID != revisionIDs[0] {
		t.Fatalf("stored revision = %q, want %q", change.JJChangeID, revisionIDs[0])
	}

	if err := os.WriteFile(filepath.Join(linked, "plain.txt"), []byte("plain\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInRepo(t, linked, "add", "plain.txt")
	runGitInRepo(t, linked, "commit", "-m", "plain linked commit")
	plainMessage := gitOutput(t, linked, "log", "-1", "--format=%B")
	plainRevisionIDs := vcs.ParseRevisionIDsFromMessage(plainMessage)
	if len(plainRevisionIDs) != 1 {
		t.Fatalf("plain commit message has revision ids %v", plainRevisionIDs)
	}
	plainOID := strings.TrimSpace(gitOutput(t, linked, "rev-parse", "HEAD"))

	if err := os.WriteFile(filepath.Join(linked, "plain.txt"), []byte("plain amended\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInRepo(t, linked, "add", "plain.txt")
	runGitInRepo(t, linked, "commit", "--amend", "--no-edit")
	amendedOID := strings.TrimSpace(gitOutput(t, linked, "rev-parse", "HEAD"))
	if amendedOID == plainOID {
		t.Fatal("amend did not change the Git commit OID")
	}
	amendedRevisionIDs := vcs.ParseRevisionIDsFromMessage(gitOutput(t, linked, "log", "-1", "--format=%B"))
	if len(amendedRevisionIDs) != 1 || amendedRevisionIDs[0] != plainRevisionIDs[0] {
		t.Fatalf("amend changed Totality revision identity: before=%v after=%v", plainRevisionIDs, amendedRevisionIDs)
	}
	plainChange, err := store.FindChangeByJJChangeID(context.Background(), repoRow.ID, plainRevisionIDs[0])
	if err != nil || plainChange == nil || plainChange.CurrentCommitID != amendedOID {
		t.Fatalf("amended change = %+v, %v; want commit %s", plainChange, err, amendedOID)
	}

	hooksDir, err := hooks.ResolveHooksDir(context.Background(), linked)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"post-commit", "post-rewrite"} {
		if err := os.Rename(filepath.Join(hooksDir, name), filepath.Join(hooksDir, name+".disabled")); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(linked, "plain.txt"), []byte("plain missed rewrite\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInRepo(t, linked, "add", "plain.txt")
	runGitInRepo(t, linked, "commit", "--amend", "--no-edit")
	missedOID := strings.TrimSpace(gitOutput(t, linked, "rev-parse", "HEAD"))
	staleChange, err := store.FindChangeByJJChangeID(context.Background(), repoRow.ID, plainRevisionIDs[0])
	if err != nil || staleChange == nil || staleChange.CurrentCommitID != amendedOID {
		t.Fatalf("change was not stale before recovery: %+v, %v", staleChange, err)
	}

	t.Setenv(hooks.SuppressAdoptedPublicationEnv, "1")
	outcome, err := hooks.RunPush(context.Background(), hooks.PushOptions{
		RepoRoot:    linked,
		RefRange:    missedOID,
		HeadSHA:     missedOID,
		LocalRef:    "refs/heads/linked",
		SkipCapture: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.RecoveryError != "" {
		t.Fatalf("pre-push recovery error = %q", outcome.RecoveryError)
	}
	recoveredChange, err := store.FindChangeByJJChangeID(context.Background(), repoRow.ID, plainRevisionIDs[0])
	if err != nil || recoveredChange == nil || recoveredChange.CurrentCommitID != missedOID {
		t.Fatalf("recovered change = %+v, %v; want commit %s", recoveredChange, err, missedOID)
	}
	if got := strings.TrimSpace(gitOutput(t, primary, "rev-parse", "HEAD")); got != primaryHead {
		t.Fatalf("primary HEAD moved: got %s want %s", got, primaryHead)
	}
	if got := gitOutput(t, primary, "diff", "--cached", "--binary"); got != primaryIndex {
		t.Fatalf("primary index changed:\n%s", got)
	}
}

func runGitInRepo(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func gitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

func buildTotalityBinary(t *testing.T) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "tx")
	cmd := exec.Command("go", "build", "-o", out, "./cmd/tx")
	cmd.Dir = mustRepoRoot(t)
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build tx: %v\n%s", err, combined)
	}
	return out
}

func mustRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatal("go.mod not found")
		}
		wd = parent
	}
}
