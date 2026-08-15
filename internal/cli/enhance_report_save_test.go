package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/codereview"
)

func TestSaveReviewReportWritesTimestampedJSON(t *testing.T) {
	repo := t.TempDir()
	startedAt := time.Date(2026, 8, 3, 21, 4, 5, 0, time.UTC)
	report := codereview.Report{Reviewed: true, RepoRoot: repo}

	path, err := saveReviewReport(repo, report, startedAt)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join("review", "review-findings-20260803-210405.json") {
		t.Fatalf("path = %q", path)
	}
	data, err := os.ReadFile(filepath.Join(repo, path))
	if err != nil {
		t.Fatal(err)
	}
	var decoded codereview.Report
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("saved report is not valid JSON: %v", err)
	}
	if !decoded.Reviewed {
		t.Fatal("round-tripped report lost Reviewed")
	}
}

func TestSaveReviewReportSkipsWhenNothingReviewed(t *testing.T) {
	repo := t.TempDir()
	path, err := saveReviewReport(repo, codereview.Report{Reviewed: false}, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if path != "" {
		t.Fatalf("path = %q, want empty", path)
	}
	if _, statErr := os.Stat(filepath.Join(repo, "review")); !os.IsNotExist(statErr) {
		t.Fatal("review directory should not be created for an empty report")
	}
}

func TestSaveReviewReportExcludesReviewDirFromGit(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, ".git", "info"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := saveReviewReport(repo, codereview.Report{Reviewed: true}, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repo, ".git", "info", "exclude"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "/review/\n") {
		t.Fatalf("exclude missing /review/ entry:\n%s", data)
	}

	// Idempotent: a second save must not duplicate the entry.
	if _, err := saveReviewReport(repo, codereview.Report{Reviewed: true}, time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(filepath.Join(repo, ".git", "info", "exclude"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(data), "/review/") != 1 {
		t.Fatalf("exclude entry duplicated:\n%s", data)
	}
}

func TestSaveReviewReportExcludesViaLinkedWorktreeGitdir(t *testing.T) {
	root := t.TempDir()
	worktree := filepath.Join(root, "wt")
	commonDir := filepath.Join(root, "main", ".git")
	worktreeGitDir := filepath.Join(commonDir, "worktrees", "wt")
	for _, dir := range []string{worktree, filepath.Join(commonDir, "info"), worktreeGitDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// Relative gitdir pointer, as git writes for worktrees under the same root.
	if err := os.WriteFile(filepath.Join(worktree, ".git"), []byte("gitdir: ../main/.git/worktrees/wt\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(worktreeGitDir, "commondir"), []byte("../..\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := saveReviewReport(worktree, codereview.Report{Reviewed: true}, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(commonDir, "info", "exclude"))
	if err != nil {
		t.Fatalf("exclude not written to the common dir git actually reads: %v", err)
	}
	if !strings.Contains(string(data), "/review/\n") {
		t.Fatalf("exclude missing /review/ entry:\n%s", data)
	}
}

func TestSaveReviewReportPathIsRepoRelative(t *testing.T) {
	repo := t.TempDir()
	path, err := saveReviewReport(repo, codereview.Report{Reviewed: true}, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.IsAbs(path) || !strings.HasPrefix(path, "review"+string(filepath.Separator)) {
		t.Fatalf("path = %q, want repo-relative under review/", path)
	}
}
