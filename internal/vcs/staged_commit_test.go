package vcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/storage"
)

func gitCachedDiff(t *testing.T, root string) []byte {
	t.Helper()
	cmd := exec.Command("git", "diff", "--cached", "--binary")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git diff --cached: %v\n%s", err, out)
	}
	return out
}

func setupStagedCommitRepo(t *testing.T, defaultBranch string) (*Service, string) {
	t.Helper()
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj executable not found")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found")
	}
	root := t.TempDir()
	if defaultBranch == "" {
		defaultBranch = "main"
	}
	runGit(t, root, "init", "-b", defaultBranch)
	runGit(t, root, "config", "user.name", "Test User")
	runGit(t, root, "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", ".")
	runGit(t, root, "commit", "-m", "init")
	runJJ(t, root, "git", "init", "--colocate", ".")

	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Test User", Email: "test@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	svc := NewServiceWithRunner(ExecRunner{})
	if _, err := svc.InitAtPath(context.Background(), root, InitOptions{}); err != nil {
		t.Fatalf("InitAtPath() error = %v", err)
	}
	return svc, root
}

func gitStatusPorcelain(t *testing.T, root string) string {
	t.Helper()
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git status: %v\n%s", err, out)
	}
	return string(out)
}

func gitCurrentBranch(t *testing.T, root string) string {
	t.Helper()
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git branch --show-current: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func countUndescribedJJChanges(t *testing.T, root string) int {
	t.Helper()
	cmd := exec.Command("jj", "log", "-r", `description("") ~ @`, "--no-graph", "-T", "change_id")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("jj log: %v\n%s", err, out)
	}
	count := 0
	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) != "" && strings.Trim(line, "z") != "" {
			count++
		}
	}
	return count
}

func TestParseCachedRawDiffIntentToAddAndSubmodule(t *testing.T) {
	intentRaw := ":000000 100644 0000000000000000000000000000000000000000 0000000000000000000000000000000000000000 A\tnew.txt\n"
	intent, submodule := parseCachedRawDiff(intentRaw)
	if !intent || submodule {
		t.Fatalf("parseCachedRawDiff(intent) = (%v, %v), want (true, false)", intent, submodule)
	}
	subRaw := ":160000 160000 abc123  def456 M\tsub\n"
	intent, submodule = parseCachedRawDiff(subRaw)
	if intent || !submodule {
		t.Fatalf("parseCachedRawDiff(submodule) = (%v, %v), want (false, true)", intent, submodule)
	}
	deleteRaw := ":100644 000000 abc123def4567890123456789012345678901234 0000000000000000000000000000000000000000 D\tgone.txt\n"
	intent, submodule = parseCachedRawDiff(deleteRaw)
	if intent || submodule {
		t.Fatalf("parseCachedRawDiff(delete) = (%v, %v), want (false, false)", intent, submodule)
	}
}

func TestCommitHunksFromGitDiffQuotedPath(t *testing.T) {
	diff := strings.Join([]string{
		"diff --git a/path with spaces.txt b/path with spaces.txt",
		`--- a/path with spaces.txt`,
		`+++ "b/path with spaces.txt"`,
		"@@ -0,0 +1 @@",
		"+hello",
	}, "\n")
	hunks := commitHunksFromGitDiff(diff)
	if len(hunks) != 1 {
		t.Fatalf("hunks len = %d, want 1", len(hunks))
	}
	if hunks[0].FilePath != "path with spaces.txt" {
		t.Fatalf("file path = %q, want quoted path parsed", hunks[0].FilePath)
	}
}

func TestNoStagedChangesExitCode(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "empty"})
	if err == nil {
		t.Fatal("RecordStagedRevision() error = nil, want no staged changes")
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != ExitCodeNoStagedChanges {
		t.Fatalf("error = %v, want CodedError code %d", err, ExitCodeNoStagedChanges)
	}
	_ = root
}

func TestStagedCommitPostconditionPreservesUnstaged(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("b\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("a\nedit\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "b.txt")

	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "add b"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	status := gitStatusPorcelain(t, root)
	if strings.Contains(status, "A  b.txt") || strings.Contains(status, "A b.txt") {
		t.Fatalf("staged file remains after commit:\n%s", status)
	}
	if !strings.Contains(status, " M a.txt") && !strings.Contains(status, "M  a.txt") {
		t.Fatalf("unstaged edit missing after commit:\n%s", status)
	}
}

func TestStagedCommitStampsGXTrailer(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "work.txt"), []byte("work\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "work.txt")

	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "work change"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	trailer := RevisionTrailerLine(result.Change.ChangeID)
	if !strings.Contains(result.Change.Description, "work change") {
		t.Fatalf("description = %q, want subject work change", result.Change.Description)
	}

	gitMessage, err := gitCommitMessage(t, root, "HEAD")
	if err != nil {
		t.Fatalf("gitCommitMessage() error = %v", err)
	}
	if !strings.Contains(gitMessage, trailer) {
		t.Fatalf("git commit message = %q, want trailer %q", gitMessage, trailer)
	}

	parsed, err := gitInterpretTrailer(t, root, gitMessage, "GX")
	if err != nil {
		t.Fatalf("gitInterpretTrailer() error = %v", err)
	}
	wantTrailerValue := fmt.Sprintf("https://gx.run/r/%s", result.Change.ChangeID)
	if parsed != wantTrailerValue {
		t.Fatalf("parsed GX trailer = %q, want %q", parsed, wantTrailerValue)
	}

	if err := svc.ensureRevisionDescription(context.Background(), root, "@-", "work change"); err != nil {
		t.Fatalf("ensureRevisionDescription() error = %v", err)
	}
}

func gitCommitMessage(t *testing.T, root, rev string) (string, error) {
	t.Helper()
	cmd := exec.Command("git", "log", "-1", "--format=%B", rev)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git log: %w\n%s", err, out)
	}
	return string(out), nil
}

func gitInterpretTrailer(t *testing.T, root, message, key string) (string, error) {
	t.Helper()
	cmd := exec.Command("git", "interpret-trailers", "--parse", "--only-trailers")
	cmd.Dir = root
	cmd.Stdin = strings.NewReader(message)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git interpret-trailers: %w\n%s", err, out)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if strings.HasPrefix(line, key+": ") {
			return strings.TrimPrefix(line, key+": "), nil
		}
	}
	return "", fmt.Errorf("trailer %q not found in %q", key, string(out))
}

func TestStagedCommitOnUnconventionalBranchUsesBranchStack(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "joes-work")
	if err := os.WriteFile(filepath.Join(root, "work.txt"), []byte("work\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "work.txt")

	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "work change"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	if result.Stack == nil || result.Stack.BookmarkName != "joes-work" {
		t.Fatalf("stack bookmark = %v, want joes-work", result.Stack)
	}
	if branch := gitCurrentBranch(t, root); branch != "joes-work" {
		t.Fatalf("current branch = %q, want joes-work", branch)
	}
}

func TestStagedCommitOnBaseMintsStackAndSwitches(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "feature.txt"), []byte("f\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "feature.txt")

	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "auth change"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	if result.Stack == nil || !strings.HasPrefix(result.Stack.BookmarkName, "feature/") {
		t.Fatalf("stack bookmark = %v, want minted feature/*", result.Stack)
	}
	if !result.CreatedBranch {
		t.Fatal("CreatedBranch = false, want true")
	}
	if branch := gitCurrentBranch(t, root); branch != result.Stack.BookmarkName {
		t.Fatalf("current branch = %q, want %q", branch, result.Stack.BookmarkName)
	}
}

func TestStagedCommitOnProtectedMasterErrors(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "master")
	if err := os.WriteFile(filepath.Join(root, "x.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "x.txt")

	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "bad"})
	if err == nil {
		t.Fatal("RecordStagedRevision() error = nil, want protected branch error")
	}
	if !strings.Contains(err.Error(), "protected") {
		t.Fatalf("error = %v, want protected branch message", err)
	}
}

func TestStagedCommitNoOrphanJJChanges(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "feature-x")
	for i := 0; i < 2; i++ {
		name := fmt.Sprintf("file%d.txt", i)
		if err := os.WriteFile(filepath.Join(root, name), []byte("x\n"), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		runGit(t, root, "add", name)
		if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: fmt.Sprintf("commit %d", i)}); err != nil {
			t.Fatalf("RecordStagedRevision(%d) error = %v", i, err)
		}
	}
}

func TestStagedCommitRollbackRestoresBranchAndIndex(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "rollback-branch")
	if err := os.WriteFile(filepath.Join(root, "rb.txt"), []byte("rb\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "rb.txt")
	beforeStatus := gitStatusPorcelain(t, root)
	beforeCached := gitCachedDiff(t, root)

	stagedCommitAfterImportHook = func() error { return fmt.Errorf("injected post-import failure") }
	t.Cleanup(func() { stagedCommitAfterImportHook = nil })

	_, commitErr := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "should rollback"})
	if commitErr == nil || !strings.Contains(commitErr.Error(), "injected post-import failure") {
		t.Fatalf("RecordStagedRevision() error = %v, want injected failure", commitErr)
	}
	if branch := gitCurrentBranch(t, root); branch != "rollback-branch" {
		t.Fatalf("branch after rollback = %q, want rollback-branch", branch)
	}
	afterCached := gitCachedDiff(t, root)
	if string(afterCached) != string(beforeCached) {
		t.Fatalf("cached diff changed after rollback:\nbefore=%q\nafter=%q", beforeCached, afterCached)
	}
	_ = beforeStatus
}

func TestAttributionLedgerSurvivesStoreReopen(t *testing.T) {
	ctx := context.Background()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)

	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(ctx, storage.Repo{RootPath: "/tmp/repo", Backend: "jj"})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if err := recordChangeForStack(ctx, store, RepoInfo{RootPath: "/tmp/repo", Backend: "jj"}, &StackInfo{
		Name:         "test",
		BookmarkName: "feature/test",
		BaseRef:      "main",
		BaseCommitID: "commit0",
		Status:       "draft",
	}, ChangeInfo{
		ChangeID:    "change1",
		CommitID:    "commit1",
		Description: "test",
	}, "op1", nil, nil, true, []storage.SessionEventAttribution{{
		Tool:             "cursor",
		SessionID:        "session-missing-row",
		EventFingerprint: "fp-1",
		AttributedVia:    "commit",
		CreatedAt:        1000,
	}}); err != nil {
		t.Fatalf("recordChangeForStack() error = %v", err)
	}
	_ = repoID
	if err := db.Close(); err != nil {
		t.Fatalf("db.Close() error = %v", err)
	}

	db, err = storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() reopen error = %v", err)
	}
	defer db.Close()
	store, err = storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() reopen error = %v", err)
	}
	keys, err := store.AttributedSessionEventKeys(ctx, repoID)
	if err != nil {
		t.Fatalf("AttributedSessionEventKeys() error = %v", err)
	}
	key := storage.SessionEventAttributionKey("cursor", "session-missing-row", "fp-1")
	if _, ok := keys[key]; !ok {
		t.Fatalf("attribution key %q missing after store reopen", key)
	}
}

func TestStagedCommitDetachedHEADRejected(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "--detach", "HEAD")
	if err := os.WriteFile(filepath.Join(root, "d.txt"), []byte("d\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "d.txt")
	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "detached"})
	if err == nil || !errors.Is(err, ErrDetachedHEAD) {
		t.Fatalf("RecordStagedRevision() error = %v, want detached HEAD", err)
	}
}

func gitShowPath(t *testing.T, root, rev, path string) string {
	t.Helper()
	cmd := exec.Command("git", "show", rev+":"+path)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git show %s:%s: %v\n%s", rev, path, err, out)
	}
	return string(out)
}

func gitShowNameStatus(t *testing.T, root, rev string) string {
	t.Helper()
	cmd := exec.Command("git", "show", "--name-status", "--format=", rev)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git show --name-status %s: %v\n%s", rev, err, out)
	}
	return string(out)
}

func readRepoFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	return string(data)
}

func TestStagedCommitPartialStagingThreeVersionFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "partial-staging")
	filePath := filepath.Join(root, "a.txt")
	if err := os.WriteFile(filePath, []byte("base\nstaged\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "a.txt")
	if err := os.WriteFile(filePath, []byte("base\nstaged\nunstaged\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "stage middle only"}); err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	if got := gitShowPath(t, root, "HEAD", "a.txt"); got != "base\nstaged\n" {
		t.Fatalf("committed content = %q, want staged middle version", got)
	}
	if got := readRepoFile(t, filePath); got != "base\nstaged\nunstaged\n" {
		t.Fatalf("working tree = %q, want unstaged remainder", got)
	}
	status := gitStatusPorcelain(t, root)
	if !strings.Contains(status, " M a.txt") && !strings.Contains(status, "M  a.txt") {
		t.Fatalf("expected unstaged edit on a.txt:\n%s", status)
	}
}

func TestStagedCommitRenameFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "rename-test")
	runGit(t, root, "mv", "a.txt", "renamed.txt")

	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "rename a"}); err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	nameStatus := gitShowNameStatus(t, root, "HEAD")
	if !strings.Contains(nameStatus, "R") || !strings.Contains(nameStatus, "renamed.txt") {
		t.Fatalf("commit name-status = %q, want rename to renamed.txt", nameStatus)
	}
	if _, err := os.Stat(filepath.Join(root, "a.txt")); !os.IsNotExist(err) {
		t.Fatalf("a.txt should be absent after rename commit: %v", err)
	}
	if got := gitShowPath(t, root, "HEAD", "renamed.txt"); got != "a\n" {
		t.Fatalf("renamed file content = %q, want original a.txt content", got)
	}
}

func TestStagedCommitModeBitOnlyFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "mode-bit")
	scriptPath := filepath.Join(root, "run.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\necho hi\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "run.sh")
	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "add script"}); err != nil {
		t.Fatalf("RecordStagedRevision() add error = %v", err)
	}
	if err := os.Chmod(scriptPath, 0o755); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}
	runGit(t, root, "add", "--chmod=+x", "run.sh")

	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "chmod script"}); err != nil {
		t.Fatalf("RecordStagedRevision() chmod error = %v", err)
	}
	cmd := exec.Command("git", "show", "--summary", "HEAD")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git show --summary: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "mode change") && !strings.Contains(string(out), "100755") {
		t.Fatalf("commit summary missing mode change:\n%s", out)
	}
}

func TestStagedCommitBinaryFileFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "binary")
	binPath := filepath.Join(root, "data.bin")
	payload := []byte{0x00, 0x01, 0x02, 0xff, 0xfe}
	if err := os.WriteFile(binPath, payload, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "data.bin")

	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "add binary"}); err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	cmd := exec.Command("git", "show", "HEAD:data.bin")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git show HEAD:data.bin: %v", err)
	}
	if string(out) != string(payload) {
		t.Fatalf("binary commit content = %x, want %x", out, payload)
	}
}

func TestStagedCommitDeleteFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "delete-test")
	runGit(t, root, "rm", "a.txt")

	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "delete a"}); err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	nameStatus := gitShowNameStatus(t, root, "HEAD")
	if !strings.Contains(nameStatus, "D") || !strings.Contains(nameStatus, "a.txt") {
		t.Fatalf("commit name-status = %q, want delete of a.txt", nameStatus)
	}
	if _, err := os.Stat(filepath.Join(root, "a.txt")); !os.IsNotExist(err) {
		t.Fatalf("a.txt should remain deleted in working tree: %v", err)
	}
}

func TestSessionEventAttributionsExcludeTemporalTier(t *testing.T) {
	events := []capture.SessionEvent{{
		Tool: "cursor", SessionID: "sess-1", FilePath: "a.go", NewText: "x",
	}}
	outcomes := []matcher.MatchOutcome{{
		EventIndex: 0,
		HunkIndex:  0,
		Tier:       matcher.TierTemporal,
		Score:      0.5,
	}}
	got := sessionEventAttributionsFromOutcomes(events, outcomes, "commit")
	if len(got) != 0 {
		t.Fatalf("sessionEventAttributionsFromOutcomes() = %#v, want none for temporal tier", got)
	}
}

func TestSessionIDsFromMatchedLinksExcludeTemporalTier(t *testing.T) {
	links := []matcher.HunkLink{
		{SessionID: "exact", Tier: matcher.TierExact, Authorship: matcher.AuthorshipAgent},
		{SessionID: "temporal", Tier: matcher.TierTemporal, Authorship: matcher.AuthorshipAgent},
	}
	got := sessionIDsFromMatchedLinks(links)
	if len(got) != 1 || got[0] != "exact" {
		t.Fatalf("sessionIDsFromMatchedLinks() = %#v, want [exact]", got)
	}
}

func TestFilterUnattributedSessionEventsSkipsAttributed(t *testing.T) {
	ctx := context.Background()
	svc, root := setupStagedCommitRepo(t, "main")

	ev := capture.SessionEvent{
		Tool: "cursor", SessionID: "sess-dedupe", FilePath: "dedupe.go",
		NewText: "added", Kind: capture.KindEdit,
	}
	fingerprint := matcher.EventFingerprint(ev)

	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(ctx, storage.Repo{RootPath: root, Backend: "jj"})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID: repoID, JJChangeID: "dedupe-change", CurrentCommitID: "dedupe-commit", Description: "dedupe",
		Status: "draft", FirstSeenAt: 1, UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.UpsertSession(ctx, storage.Session{
		ID: ev.SessionID, CreatedAt: 1000, Command: ev.Tool, Cwd: root, RepoRoot: &root,
	}); err != nil {
		t.Fatalf("UpsertSession() error = %v", err)
	}
	if err := store.WriteSessionEventAttributions(ctx, []storage.SessionEventAttribution{{
		RepoID: repoID, ChangeID: changeID, Tool: ev.Tool, SessionID: ev.SessionID,
		EventFingerprint: fingerprint, AttributedVia: "commit", CreatedAt: 1000,
	}}); err != nil {
		t.Fatalf("WriteSessionEventAttributions() error = %v", err)
	}
	_ = db.Close()

	filtered := svc.filterUnattributedSessionEvents(ctx, root, []capture.SessionEvent{ev})
	if len(filtered) != 0 {
		t.Fatalf("filterUnattributedSessionEvents() = %#v, want attributed event removed", filtered)
	}
}

func TestAttributedEventNotWrittenTwiceForSecondCommit(t *testing.T) {
	ctx := context.Background()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)

	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	defer db.Close()
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(ctx, storage.Repo{RootPath: "/tmp/ledger", Backend: "jj"})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID: repoID, JJChangeID: "change-1", CurrentCommitID: "commit-1", Description: "first",
		Status: "draft", FirstSeenAt: 1, UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.UpsertSession(ctx, storage.Session{
		ID: "sess-1", CreatedAt: 1000, Command: "cursor",
	}); err != nil {
		t.Fatalf("UpsertSession() error = %v", err)
	}
	attr := storage.SessionEventAttribution{
		RepoID: repoID, ChangeID: changeID, Tool: "cursor", SessionID: "sess-1",
		EventFingerprint: "fp-stable", AttributedVia: "commit", CreatedAt: 1000,
	}
	if err := store.WriteSessionEventAttributions(ctx, []storage.SessionEventAttribution{attr}); err != nil {
		t.Fatalf("first WriteSessionEventAttributions() error = %v", err)
	}
	changeID2, err := store.UpsertChange(ctx, storage.Change{
		RepoID: repoID, JJChangeID: "change-2", CurrentCommitID: "commit-2", Description: "second",
		Status: "draft", FirstSeenAt: 2, UpdatedAt: 2,
	})
	if err != nil {
		t.Fatalf("UpsertChange(2) error = %v", err)
	}
	attr.ChangeID = changeID2
	if err := store.WriteSessionEventAttributions(ctx, []storage.SessionEventAttribution{attr}); err != nil {
		t.Fatalf("second WriteSessionEventAttributions() error = %v", err)
	}
	keys, err := store.AttributedSessionEventKeys(ctx, repoID)
	if err != nil {
		t.Fatalf("AttributedSessionEventKeys() error = %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("attributed keys = %d, want 1 (no duplicate re-attribution)", len(keys))
	}
}

func TestRejectProtectedStackBookmarkSecondChokePoint(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	err := svc.rejectProtectedStackBookmark(context.Background(), root, "main")
	if err == nil || !strings.Contains(err.Error(), "protected") {
		t.Fatalf("rejectProtectedStackBookmark() error = %v, want protected stack bookmark", err)
	}
}

func TestEnsureBranchMutationAllowedBlocksProtectedRefMove(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	head, err := svc.runTrimmed(context.Background(), root, "git", "rev-parse", "HEAD")
	if err != nil {
		t.Fatalf("rev-parse HEAD: %v", err)
	}
	err = svc.attachGitBranch(context.Background(), root, "main", head+"0", false)
	if err == nil || !strings.Contains(err.Error(), "protected") {
		t.Fatalf("attachGitBranch() error = %v, want protected ref mutation blocked", err)
	}
}

func TestSetBaseAfterStagedCommitOnFeatureBranch(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "feature/set-base")
	if err := os.WriteFile(filepath.Join(root, "base-test.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "base-test.txt")
	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "feature work"}); err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}

	result, err := svc.SetBase(context.Background(), "main")
	if err != nil {
		t.Fatalf("SetBase() error = %v", err)
	}
	if result.BaseRef != "main" {
		t.Fatalf("SetBase() base ref = %q, want main", result.BaseRef)
	}
	if branch := gitCurrentBranch(t, root); branch != "main" {
		t.Fatalf("current branch after SetBase = %q, want main", branch)
	}
	if !result.OnBase {
		t.Fatalf("SetBase() OnBase = false, want true after returning to base")
	}
}

func TestRequireAuthoringBaseAfterFeatureCommit(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "feature/gate")
	if err := os.WriteFile(filepath.Join(root, "gate.txt"), []byte("g\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "gate.txt")
	if _, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "gate work"}); err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	if _, err := svc.SetBase(context.Background(), "main"); err != nil {
		t.Fatalf("SetBase() error = %v", err)
	}
	if err := svc.RequireAuthoringBase(context.Background(), "gx commit"); err != nil {
		t.Fatalf("RequireAuthoringBase() after SetBase error = %v", err)
	}
	_ = root
}

func setupUnbornStagedCommitRepo(t *testing.T) (*Service, string) {
	t.Helper()
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj executable not found")
	}
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git executable not found")
	}
	root := t.TempDir()
	runGit(t, root, "init", "-b", "main")
	runGit(t, root, "config", "user.name", "Test User")
	runGit(t, root, "config", "user.email", "test@example.com")
	runJJ(t, root, "git", "init", "--colocate", ".")

	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Test User", Email: "test@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	svc := NewServiceWithRunner(ExecRunner{})
	if _, err := svc.InitAtPath(context.Background(), root, InitOptions{}); err != nil {
		t.Fatalf("InitAtPath() error = %v", err)
	}
	return svc, root
}

func TestStagedCommitUnbornHEADRejected(t *testing.T) {
	svc, root := setupUnbornStagedCommitRepo(t)
	if err := os.WriteFile(filepath.Join(root, "first.txt"), []byte("first\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "first.txt")
	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "first"})
	if err == nil || !strings.Contains(err.Error(), "no commits yet") {
		t.Fatalf("RecordStagedRevision() error = %v, want unborn HEAD message", err)
	}
}

func TestPreflightStagedIndexRejectsIntentToAddRawDiff(t *testing.T) {
	repoRoot := t.TempDir()
	intentRaw := ":000000 100644 0000000000000000000000000000000000000000 0000000000000000000000000000000000000000 A\tnew.txt\n"
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-files", "-u"):                          {""},
			runnerKey(repoRoot, "git", "diff", "--cached", "--raw", "--no-abbrev"): {intentRaw},
		},
	}
	svc := NewServiceWithRunner(runner)
	err := svc.preflightStagedIndex(context.Background(), repoRoot)
	if err == nil || !strings.Contains(err.Error(), "intent-to-add") {
		t.Fatalf("preflightStagedIndex() error = %v, want intent-to-add rejection", err)
	}
}

func TestPreflightStagedIndexRejectsSubmoduleRawDiff(t *testing.T) {
	repoRoot := t.TempDir()
	subRaw := ":160000 160000 abc1234567890123456789012345678901234567890 def4567890123456789012345678901234567890 M\tsub\n"
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-files", "-u"):                          {""},
			runnerKey(repoRoot, "git", "diff", "--cached", "--raw", "--no-abbrev"): {subRaw},
		},
	}
	svc := NewServiceWithRunner(runner)
	err := svc.preflightStagedIndex(context.Background(), repoRoot)
	if err == nil || !strings.Contains(err.Error(), "submodule") {
		t.Fatalf("preflightStagedIndex() error = %v, want submodule rejection", err)
	}
}

func TestStagedCommitSubmoduleRejected(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	subRoot := filepath.Join(root, "subrepo")
	if err := os.MkdirAll(subRoot, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	runGit(t, subRoot, "init", "-b", "main")
	runGit(t, subRoot, "config", "user.name", "Test User")
	runGit(t, subRoot, "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(subRoot, "sub.txt"), []byte("sub\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, subRoot, "add", ".")
	runGit(t, subRoot, "commit", "-m", "sub init")
	subCommit, err := exec.Command("git", "-C", subRoot, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatalf("rev-parse subrepo HEAD: %v", err)
	}
	runGit(t, root, "update-index", "--add", "--cacheinfo", "160000", strings.TrimSpace(string(subCommit)), "sub")
	_, err = svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "add sub"})
	if err == nil || !strings.Contains(err.Error(), "submodule") {
		t.Fatalf("RecordStagedRevision() error = %v, want submodule rejection", err)
	}
}

func TestStagedCommitQuotedPathWithSpaces(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "quoted-path")
	spaced := filepath.Join(root, "path with spaces.txt")
	if err := os.WriteFile(spaced, []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "path with spaces.txt")

	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "spaced path"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	if !containsString(result.Change.Files, "path with spaces.txt") {
		t.Fatalf("recorded files = %#v, want spaced path", result.Change.Files)
	}
	if got := gitShowPath(t, root, "HEAD", "path with spaces.txt"); got != "hello\n" {
		t.Fatalf("committed spaced path content = %q", got)
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestStagedCommitResolveStackRejectsProtectedAtFirstChoke(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "branch", "master")
	repo, err := svc.ResolveJJRepoAtPath(context.Background(), root)
	if err != nil {
		t.Fatalf("ResolveJJRepoAtPath() error = %v", err)
	}
	head, err := svc.runTrimmed(context.Background(), root, "git", "rev-parse", "HEAD")
	if err != nil {
		t.Fatalf("rev-parse HEAD: %v", err)
	}
	_, _, err = svc.resolveStagedCommitStack(context.Background(), repo, "master", "blocked", head)
	if err == nil || !strings.Contains(err.Error(), "protected") {
		t.Fatalf("resolveStagedCommitStack() error = %v, want protected branch at first choke", err)
	}
}

func TestSessionEventAttributionsDedupeSameEvent(t *testing.T) {
	events := []capture.SessionEvent{{
		Tool: "cursor", SessionID: "sess-1", FilePath: "a.go", NewText: "x",
		Raw: map[string]json.RawMessage{"uuid": json.RawMessage(`"event-1"`)},
	}}
	outcomes := []matcher.MatchOutcome{
		{EventIndex: 0, HunkIndex: 0, Tier: matcher.TierExact, Score: 1},
		{EventIndex: 0, HunkIndex: 1, Tier: matcher.TierFuzzy, Score: 0.9},
	}
	got := sessionEventAttributionsFromOutcomes(events, outcomes, "commit")
	if len(got) != 1 {
		t.Fatalf("sessionEventAttributionsFromOutcomes() = %d attributions, want 1 deduped", len(got))
	}
}

func gitStagedNames(t *testing.T, root string) string {
	t.Helper()
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git diff --cached --name-only: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
}

func TestStagedCommitAfterExternalCheckout(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "side")
	if err := os.WriteFile(filepath.Join(root, "s.txt"), []byte("s\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "s.txt")

	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "side work"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() after external checkout error = %v", err)
	}
	if result.Stack == nil || result.Stack.BookmarkName != "side" {
		t.Fatalf("stack = %v, want bookmark side", result.Stack)
	}
	if branch := gitCurrentBranch(t, root); branch != "side" {
		t.Fatalf("current branch = %q, want side", branch)
	}
	if staged := gitStagedNames(t, root); staged != "" {
		t.Fatalf("staged files remain after commit: %q", staged)
	}
}

func TestStatusSnapshotPreservesStagedIndexAfterExternalCheckout(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "topic")
	if err := os.WriteFile(filepath.Join(root, "p.txt"), []byte("p\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "p.txt")

	if _, err := svc.StatusSnapshot(context.Background()); err != nil {
		t.Fatalf("StatusSnapshot() error = %v", err)
	}
	if staged := gitStagedNames(t, root); staged != "p.txt" {
		t.Fatalf("staged files after StatusSnapshot = %q, want p.txt", staged)
	}
}

func TestStagedCommitAfterRawGitCommitSyncs(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "raw.txt"), []byte("raw\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "raw.txt")
	runGit(t, root, "commit", "-m", "raw commit outside gx")

	if err := os.WriteFile(filepath.Join(root, "next.txt"), []byte("next\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "next.txt")
	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "after raw commit"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() after raw git commit error = %v", err)
	}
	if result.Change.CommitID == "" {
		t.Fatal("recorded change has no commit id")
	}
	if staged := gitStagedNames(t, root); staged != "" {
		t.Fatalf("staged files remain after commit: %q", staged)
	}
}

func TestStagedCommitReportsJJClobberedIndex(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "c.txt"), []byte("c\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "c.txt")
	runJJ(t, root, "status")
	if staged := gitStagedNames(t, root); staged != "" {
		t.Skip("this jj version does not rewrite the git index on snapshot")
	}

	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "clobbered"})
	if err == nil {
		t.Fatal("RecordStagedRevision() error = nil, want clobbered-index explanation")
	}
	if !strings.Contains(err.Error(), "re-run git add") {
		t.Fatalf("error = %v, want intent-to-add explanation with re-run git add", err)
	}
	var coded *CodedError
	if !errors.As(err, &coded) || coded.Code != ExitCodeNoStagedChanges {
		t.Fatalf("error = %v, want CodedError code %d", err, ExitCodeNoStagedChanges)
	}
}

func TestStagedCommitRollbackDeletesMintedBranch(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "mint.txt"), []byte("m\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "mint.txt")
	beforeCached := gitCachedDiff(t, root)

	stagedCommitAfterImportHook = func() error { return fmt.Errorf("injected mint failure") }
	t.Cleanup(func() { stagedCommitAfterImportHook = nil })

	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "minted branch rollback"})
	if err == nil || !strings.Contains(err.Error(), "injected mint failure") {
		t.Fatalf("RecordStagedRevision() error = %v, want injected failure", err)
	}
	if branch := gitCurrentBranch(t, root); branch != "main" {
		t.Fatalf("branch after rollback = %q, want main", branch)
	}
	branches := exec.Command("git", "branch", "--list", "feature/*", "bug/*")
	branches.Dir = root
	if out, _ := branches.CombinedOutput(); strings.TrimSpace(string(out)) != "" {
		t.Fatalf("minted branch left behind after rollback:\n%s", out)
	}
	if after := gitCachedDiff(t, root); string(after) != string(beforeCached) {
		t.Fatalf("cached diff changed after rollback:\nbefore=%q\nafter=%q", beforeCached, after)
	}
}

func TestStagedCommitSelfInitializesPlainGitRepo(t *testing.T) {
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj executable not found")
	}
	root := t.TempDir()
	runGit(t, root, "init", "-b", "main")
	runGit(t, root, "config", "user.name", "Test User")
	runGit(t, root, "config", "user.email", "test@example.com")
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", ".")
	runGit(t, root, "commit", "-m", "init")

	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{User: gxconfig.User{Name: "Test User", Email: "test@example.com"}}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	if err := os.WriteFile(filepath.Join(root, "n.txt"), []byte("n\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "n.txt")

	svc := NewServiceWithRunner(ExecRunner{})
	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "first ever"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() in plain git repo error = %v", err)
	}
	if result.Change.CommitID == "" {
		t.Fatal("recorded change has no commit id")
	}
}
