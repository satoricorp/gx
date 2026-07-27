package vcs

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateRevisionID(t *testing.T) {
	id, err := GenerateRevisionID()
	if err != nil {
		t.Fatal(err)
	}
	if !ValidRevisionID(id) {
		t.Fatalf("GenerateRevisionID() = %q, invalid format", id)
	}
	second, err := GenerateRevisionID()
	if err != nil {
		t.Fatal(err)
	}
	if id == second {
		t.Fatalf("GenerateRevisionID() returned duplicate ids")
	}
}

func TestPrepareCommitMessageHookPreservesExistingTrailer(t *testing.T) {
	existing, err := GenerateRevisionID()
	if err != nil {
		t.Fatal(err)
	}
	message := "add feature\n\n" + RevisionTrailerLine(existing)
	got, err := PrepareCommitMessageHook(message)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, RevisionTrailerLine(existing)) {
		t.Fatalf("PrepareCommitMessageHook() = %q, want preserved trailer", got)
	}
	if len(ParseRevisionIDsFromMessage(got)) != 1 {
		t.Fatalf("PrepareCommitMessageHook() = %q, want exactly one trailer", got)
	}
}

func TestPrepareCommitMessageHookAddsTrailer(t *testing.T) {
	got, err := PrepareCommitMessageHook("add feature")
	if err != nil {
		t.Fatal(err)
	}
	ids := ParseRevisionIDsFromMessage(got)
	if len(ids) != 1 {
		t.Fatalf("PrepareCommitMessageHook() = %q, want one trailer, got %d", got, len(ids))
	}
	if !ValidRevisionID(ids[0]) {
		t.Fatalf("PrepareCommitMessageHook() trailer id = %q, invalid", ids[0])
	}
}

func TestParsePostRewriteMappings(t *testing.T) {
	mappings := ParsePostRewriteMappings("aaa111 bbb222\nccc333 ddd444\n")
	if len(mappings) != 2 {
		t.Fatalf("len = %d, want 2", len(mappings))
	}
	if mappings[0].OldOID != "aaa111" || mappings[0].NewOID != "bbb222" {
		t.Fatalf("first mapping = %+v", mappings[0])
	}
}

func TestNormalizeGitCommonDir(t *testing.T) {
	worktree := "/tmp/repo"
	got := NormalizeGitCommonDir(worktree, ".git")
	want := "/tmp/repo/.git"
	if got != want {
		t.Fatalf("NormalizeGitCommonDir() = %q, want %q", got, want)
	}
}

func TestPostCommitHookRecordsPlainGitCommit(t *testing.T) {
	repo := initGitNativeTestRepo(t)
	t.Setenv("GX_HOME", t.TempDir())
	service := NewService()
	if _, err := service.InitAtPath(context.Background(), repo, InitOptions{}); err != nil {
		t.Fatal(err)
	}
	writeGitNativeTestFile(t, repo, "feature.txt", "feature\n")
	runGitNativeTestGit(t, repo, "add", "feature.txt")
	stamped, err := PrepareCommitMessageHook("add feature")
	if err != nil {
		t.Fatal(err)
	}
	runGitNativeTestGit(t, repo, "commit", "-m", stamped)

	// The post-commit hook runs without a pending context file; it must fall
	// back to the repo it was handed and still record the revision.
	if err := service.RunPostCommitHook(context.Background(), repo); err != nil {
		t.Fatal(err)
	}

	message := runGitNativeTestGit(t, repo, "log", "-1", "--format=%B")
	ids := ParseRevisionIDsFromMessage(message)
	if len(ids) != 1 || !ValidRevisionID(ids[0]) {
		t.Fatalf("commit message ids = %v, want exactly one valid revision id", ids)
	}
	headOID := strings.TrimSpace(runGitNativeTestGit(t, repo, "rev-parse", "HEAD"))
	repoInfo, err := service.ResolveGXRepoAtPath(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	store, err := openStore(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	repoRow, err := store.FindRepoByIdentity(context.Background(), repoInfo.GitCommonDir, repoInfo.RootPath)
	if err != nil || repoRow == nil {
		t.Fatalf("FindRepoByIdentity() = %+v, %v", repoRow, err)
	}
	change, err := store.FindChangeByJJChangeID(context.Background(), repoRow.ID, ids[0])
	if err != nil || change == nil {
		t.Fatalf("FindChangeByJJChangeID() = %+v, %v", change, err)
	}
	if change.CurrentCommitID != headOID {
		t.Fatalf("recorded commit = %q, want %q", change.CurrentCommitID, headOID)
	}
}

func TestPendingContextIsolatedByLinkedWorktree(t *testing.T) {
	repo := initGitNativeTestRepo(t)
	worktree := filepath.Join(t.TempDir(), "linked")
	runGitNativeTestGit(t, repo, "worktree", "add", "-b", "linked", worktree)
	service := NewService()
	mainPaths, err := service.ResolveGitPaths(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	linkedPaths, err := service.ResolveGitPaths(context.Background(), worktree)
	if err != nil {
		t.Fatal(err)
	}
	if mainPaths.CommonDir != linkedPaths.CommonDir {
		t.Fatalf("common dirs differ: %q != %q", mainPaths.CommonDir, linkedPaths.CommonDir)
	}
	if mainPaths.GitDir == linkedPaths.GitDir {
		t.Fatalf("git dirs unexpectedly equal: %q", mainPaths.GitDir)
	}
	mainPending := PendingCommitContext{WorktreeRoot: repo}
	linkedPending := PendingCommitContext{WorktreeRoot: worktree}
	if err := writePendingCommitContext(mainPaths.GitDir, mainPending); err != nil {
		t.Fatal(err)
	}
	if err := writePendingCommitContext(linkedPaths.GitDir, linkedPending); err != nil {
		t.Fatal(err)
	}
	mainGot, err := readPendingCommitContext(mainPaths.GitDir)
	if err != nil {
		t.Fatal(err)
	}
	linkedGot, err := readPendingCommitContext(linkedPaths.GitDir)
	if err != nil {
		t.Fatal(err)
	}
	if mainGot.WorktreeRoot != repo || linkedGot.WorktreeRoot != worktree {
		t.Fatalf("pending contexts crossed: main=%q linked=%q", mainGot.WorktreeRoot, linkedGot.WorktreeRoot)
	}
	info, err := os.Stat(pendingCommitContextPath(mainPaths.GitDir))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("pending context permissions = %o, want 600", info.Mode().Perm())
	}
}

func TestRecoverGitCommitByOIDUpdatesExistingRevision(t *testing.T) {
	repoRoot := initGitNativeTestRepo(t)
	t.Setenv("GX_HOME", t.TempDir())
	service := NewService()
	if _, err := service.InitAtPath(context.Background(), repoRoot, InitOptions{}); err != nil {
		t.Fatal(err)
	}
	revisionID, err := GenerateRevisionID()
	if err != nil {
		t.Fatal(err)
	}
	writeGitNativeTestFile(t, repoRoot, "feature.txt", "one\n")
	runGitNativeTestGit(t, repoRoot, "add", "feature.txt")
	message := StampRevisionTrailer("feature", revisionID)
	runGitNativeTestGit(t, repoRoot, "commit", "-m", message)
	oldOID := strings.TrimSpace(runGitNativeTestGit(t, repoRoot, "rev-parse", "HEAD"))
	repo, err := service.ResolveGXRepoAtPath(context.Background(), repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.RecordGitCommit(context.Background(), repo, oldOID, PendingCommitContext{WorktreeRoot: repo.RootPath, GitCommonDir: repo.GitCommonDir}); err != nil {
		t.Fatal(err)
	}
	writeGitNativeTestFile(t, repoRoot, "feature.txt", "two\n")
	runGitNativeTestGit(t, repoRoot, "add", "feature.txt")
	runGitNativeTestGit(t, repoRoot, "commit", "--amend", "-m", message)
	newOID := strings.TrimSpace(runGitNativeTestGit(t, repoRoot, "rev-parse", "HEAD"))

	store, err := openStore(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	repoRow, err := store.FindRepoByIdentity(context.Background(), repo.GitCommonDir, repo.RootPath)
	if err != nil || repoRow == nil {
		t.Fatalf("FindRepoByIdentity() = %+v, %v", repoRow, err)
	}
	if err := service.recoverGitCommitByOID(context.Background(), store, repo, *repoRow, newOID); err != nil {
		t.Fatal(err)
	}
	change, err := store.FindChangeByJJChangeID(context.Background(), repoRow.ID, revisionID)
	if err != nil {
		t.Fatal(err)
	}
	if change == nil || change.CurrentCommitID != newOID {
		t.Fatalf("recovered change = %+v, want commit %s", change, newOID)
	}
}

func initGitNativeTestRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runGitNativeTestGit(t, repo, "init", "-b", "main")
	runGitNativeTestGit(t, repo, "config", "user.name", "GX Test")
	runGitNativeTestGit(t, repo, "config", "user.email", "gx@example.com")
	writeGitNativeTestFile(t, repo, "README.md", "base\n")
	runGitNativeTestGit(t, repo, "add", "README.md")
	runGitNativeTestGit(t, repo, "commit", "-m", "base")
	return repo
}

func writeGitNativeTestFile(t *testing.T, repo, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(repo, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runGitNativeTestGit(t *testing.T, repo string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = repo
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}
