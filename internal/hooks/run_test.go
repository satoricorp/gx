package hooks

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	cmds := [][]string{
		{"git", "init", "-b", "main"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	readme := filepath.Join(dir, "README.md")
	if err := os.WriteFile(readme, []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "add", "README.md")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v: %s", err, out)
	}
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v: %s", err, out)
	}
}

func TestRunPushRespectsPauseFlag(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()
	initGitRepo(t, repo)

	gxHome := filepath.Join(home, ".gx")
	t.Setenv("GX_HOME", gxHome)
	if err := os.MkdirAll(gxHome, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gxHome, "pause-capture"), []byte("1\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	outcome, err := RunPush(context.Background(), PushOptions{
		RepoRoot: repo,
		HomeDir:  home,
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if outcome.Result.RefRange != "" || outcome.Result.EligibleHunks != 0 {
		t.Fatalf("RunPush() = %+v, want empty result when paused", outcome)
	}
}

func TestParseRefRangeSingleSHALeavesBaseForOrchestrator(t *testing.T) {
	base, head := parseRefRange("abc123", "", "")
	if base != "" || head != "abc123" {
		t.Fatalf("parseRefRange() = (%q, %q), want empty base and abc123 head", base, head)
	}
}

func TestRunPushRecoversMissingRevision(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")
	t.Setenv(SuppressAdoptedPublicationEnv, "1")
	initGitRepo(t, repoRoot)
	service := vcs.NewService()
	if _, err := service.InitAtPath(context.Background(), repoRoot, vcs.InitOptions{}); err != nil {
		t.Fatal(err)
	}
	revisionID, err := vcs.GenerateRevisionID()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, "feature.txt"), []byte("feature\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitInHookTest(t, repoRoot, "add", "feature.txt")
	runGitInHookTest(t, repoRoot, "commit", "--no-verify", "-m", vcs.StampRevisionTrailer("feature", revisionID))
	commitOID := strings.TrimSpace(gitOutputInHookTest(t, repoRoot, "rev-parse", "HEAD"))

	outcome, err := RunPush(context.Background(), PushOptions{
		RepoRoot:    repoRoot,
		RefRange:    commitOID,
		HeadSHA:     commitOID,
		LocalRef:    "refs/heads/main",
		SkipCapture: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if outcome.RecoveryError != "" {
		t.Fatalf("recovery error = %q", outcome.RecoveryError)
	}
	if len(outcome.RevisionIDs) != 1 || outcome.RevisionIDs[0] != revisionID {
		t.Fatalf("revision ids = %v, want %q", outcome.RevisionIDs, revisionID)
	}
	repo, err := service.ResolveGXRepoAtPath(context.Background(), repoRoot)
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
	change, err := store.FindChangeByCommitID(context.Background(), repoRow.ID, commitOID)
	if err != nil || change == nil {
		t.Fatalf("FindChangeByCommitID() = %+v, %v", change, err)
	}
}

func TestRunPushReportsRevisionScanFailure(t *testing.T) {
	repoRoot := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	initGitRepo(t, repoRoot)

	outcome, err := RunPush(context.Background(), PushOptions{
		RepoRoot:    repoRoot,
		RefRange:    "missing-ref..HEAD",
		SkipCapture: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(outcome.RecoveryError, "scan GX revision trailers") {
		t.Fatalf("recovery error = %q, want revision scan failure", outcome.RecoveryError)
	}
}

func runGitInHookTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func gitOutputInHookTest(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}
