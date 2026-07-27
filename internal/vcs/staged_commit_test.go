package vcs

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/gxconfig"
)

func setupStagedCommitRepo(t *testing.T, defaultBranch string) (*Service, string) {
	t.Helper()
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

func gitHeadCommit(t *testing.T, root string) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse HEAD: %v\n%s", err, out)
	}
	return strings.TrimSpace(string(out))
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

// commitViaHooks commits the staged selection with plain git and records it
// exactly as GX's prepare-commit-msg and post-commit hooks do.
func commitViaHooks(t *testing.T, svc *Service, root, message string) CommitResult {
	t.Helper()
	stamped, err := PrepareCommitMessageHook(message)
	if err != nil {
		t.Fatalf("PrepareCommitMessageHook() error = %v", err)
	}
	runGit(t, root, "commit", "-m", stamped)
	ctx := context.Background()
	repo, err := svc.ResolveGXRepoAtPath(ctx, root)
	if err != nil {
		t.Fatalf("ResolveGXRepoAtPath() error = %v", err)
	}
	result, err := svc.RecordGitCommit(ctx, repo, gitHeadCommit(t, root), PendingCommitContext{
		WorktreeRoot: repo.RootPath,
		GitCommonDir: repo.GitCommonDir,
	})
	if err != nil {
		t.Fatalf("RecordGitCommit() error = %v", err)
	}
	return result
}

func TestRecordedCommitPreservesUnstagedWork(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "b.txt"), []byte("b\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("a\nedit\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "b.txt")

	commitViaHooks(t, svc, root, "add b")
	status := gitStatusPorcelain(t, root)
	if strings.Contains(status, "A  b.txt") || strings.Contains(status, "A b.txt") {
		t.Fatalf("staged file remains after commit:\n%s", status)
	}
	if !strings.Contains(status, " M a.txt") && !strings.Contains(status, "M  a.txt") {
		t.Fatalf("unstaged edit missing after commit:\n%s", status)
	}
}

func TestRecordedCommitStampsGXTrailer(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "work.txt"), []byte("work\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "work.txt")

	result := commitViaHooks(t, svc, root, "work change")
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

func TestRecordedCommitOnUnconventionalBranchUsesBranchStack(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "joes-work")
	if err := os.WriteFile(filepath.Join(root, "work.txt"), []byte("work\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "work.txt")

	result := commitViaHooks(t, svc, root, "work change")
	if result.Stack == nil || result.Stack.BookmarkName != "joes-work" {
		t.Fatalf("stack bookmark = %v, want joes-work", result.Stack)
	}
	if branch := gitCurrentBranch(t, root); branch != "joes-work" {
		t.Fatalf("current branch = %q, want joes-work", branch)
	}
}

func TestRecordedCommitOnBaseStaysOnBase(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	beforeHead := gitHeadCommit(t, root)
	if err := os.WriteFile(filepath.Join(root, "feature.txt"), []byte("f\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "feature.txt")

	result := commitViaHooks(t, svc, root, "auth change")
	if result.Stack == nil || result.Stack.BookmarkName != "main" {
		t.Fatalf("stack bookmark = %v, want main (the checked-out branch)", result.Stack)
	}
	if branch := gitCurrentBranch(t, root); branch != "main" {
		t.Fatalf("current branch = %q, want main: recording must not move HEAD", branch)
	}
	if head := gitHeadCommit(t, root); head == beforeHead {
		t.Fatal("HEAD did not advance; expected the staged commit on main")
	}
}

func TestRecordedCommitOnProtectedMasterAdvancesBranch(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "master")
	beforeHead := gitHeadCommit(t, root)
	if err := os.WriteFile(filepath.Join(root, "x.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "x.txt")

	result := commitViaHooks(t, svc, root, "work on master")
	if branch := gitCurrentBranch(t, root); branch != "master" {
		t.Fatalf("current branch = %q, want master", branch)
	}
	if head := gitHeadCommit(t, root); head == beforeHead {
		t.Fatal("master did not advance")
	}
	if result.Stack == nil || result.Stack.BookmarkName != "master" {
		t.Fatalf("stack bookmark = %v, want master", result.Stack)
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

func TestRecordedCommitPartialStagingThreeVersionFidelity(t *testing.T) {
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

	commitViaHooks(t, svc, root, "stage middle only")
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

func TestRecordedCommitRenameFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "rename-test")
	runGit(t, root, "mv", "a.txt", "renamed.txt")

	result := commitViaHooks(t, svc, root, "rename a")
	if !containsString(result.Change.Files, "renamed.txt") {
		t.Fatalf("recorded files = %#v, want renamed.txt", result.Change.Files)
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

func TestRecordedCommitModeBitOnlyFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "mode-bit")
	scriptPath := filepath.Join(root, "run.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\necho hi\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "run.sh")
	commitViaHooks(t, svc, root, "add script")
	if err := os.Chmod(scriptPath, 0o755); err != nil {
		t.Fatalf("Chmod() error = %v", err)
	}
	runGit(t, root, "add", "--chmod=+x", "run.sh")

	commitViaHooks(t, svc, root, "chmod script")
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

func TestRecordedCommitBinaryFileFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "binary")
	binPath := filepath.Join(root, "data.bin")
	payload := []byte{0x00, 0x01, 0x02, 0xff, 0xfe}
	if err := os.WriteFile(binPath, payload, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "data.bin")

	commitViaHooks(t, svc, root, "add binary")
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

func TestRecordedCommitDeleteFidelity(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "delete-test")
	runGit(t, root, "rm", "a.txt")

	result := commitViaHooks(t, svc, root, "delete a")
	if !containsString(result.Change.Files, "a.txt") {
		t.Fatalf("recorded files = %#v, want deleted a.txt", result.Change.Files)
	}
	nameStatus := gitShowNameStatus(t, root, "HEAD")
	if !strings.Contains(nameStatus, "D") || !strings.Contains(nameStatus, "a.txt") {
		t.Fatalf("commit name-status = %q, want delete of a.txt", nameStatus)
	}
	if _, err := os.Stat(filepath.Join(root, "a.txt")); !os.IsNotExist(err) {
		t.Fatalf("a.txt should remain deleted in working tree: %v", err)
	}
}

func TestRejectProtectedStackBookmarkSecondChokePoint(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	err := svc.rejectProtectedStackBookmark(context.Background(), root, "main")
	if err == nil || !strings.Contains(err.Error(), "protected") {
		t.Fatalf("rejectProtectedStackBookmark() error = %v, want protected stack bookmark", err)
	}
}

func TestRecordedCommitQuotedPathWithSpaces(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "quoted-path")
	spaced := filepath.Join(root, "path with spaces.txt")
	if err := os.WriteFile(spaced, []byte("hello\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "path with spaces.txt")

	result := commitViaHooks(t, svc, root, "spaced path")
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

func TestRecordedCommitAfterExternalCheckout(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "side")
	if err := os.WriteFile(filepath.Join(root, "s.txt"), []byte("s\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "s.txt")

	result := commitViaHooks(t, svc, root, "side work")
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

func TestStackPreservesStagedIndexAfterExternalCheckout(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	runGit(t, root, "checkout", "-b", "topic")
	if err := os.WriteFile(filepath.Join(root, "p.txt"), []byte("p\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "p.txt")

	if _, err := svc.Stack(context.Background()); err != nil {
		t.Fatalf("Stack() error = %v", err)
	}
	if staged := gitStagedNames(t, root); staged != "p.txt" {
		t.Fatalf("staged files after Stack = %q, want p.txt", staged)
	}
}

func TestRecordedCommitAfterUnstampedGitCommitSyncs(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "raw.txt"), []byte("raw\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "raw.txt")
	runGit(t, root, "commit", "-m", "raw commit without a GX trailer")

	if err := os.WriteFile(filepath.Join(root, "next.txt"), []byte("next\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "next.txt")
	result := commitViaHooks(t, svc, root, "after raw commit")
	if result.Change.CommitID == "" {
		t.Fatal("recorded change has no commit id")
	}
	if staged := gitStagedNames(t, root); staged != "" {
		t.Fatalf("staged files remain after commit: %q", staged)
	}
}

func TestRecordedCommitSelfInitializesPlainGitRepo(t *testing.T) {
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
	result := commitViaHooks(t, svc, root, "first ever")
	if result.Change.CommitID == "" {
		t.Fatal("recorded change has no commit id")
	}
	if result.Stack == nil || result.Stack.BookmarkName != "main" {
		t.Fatalf("stack = %v, want the checked-out branch recorded in a fresh repo", result.Stack)
	}
}
