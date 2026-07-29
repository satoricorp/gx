package orchestrator_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/satoricorp/totality/internal/capture"
	"github.com/satoricorp/totality/internal/capture/matcher"
	"github.com/satoricorp/totality/internal/capture/orchestrator"
	"github.com/satoricorp/totality/internal/capture/redact"
	"github.com/satoricorp/totality/internal/storage"
)

func TestOrchestrator_StagesExtract(t *testing.T) {
	repo := initTestGitRepo(t)
	totalityHome := t.TempDir()
	t.Setenv("TOTALITY_HOME", totalityHome)

	ctx := context.Background()
	result, err := orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot:  repo,
		Base:      "HEAD~1",
		Head:      "HEAD",
		Tools:     []string{capture.ToolClaude},
		StageOnly: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.StagedExtractID == "" {
		t.Fatal("expected staged extract id")
	}

	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	pending, err := stager.PendingExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) == 0 {
		t.Fatal("expected pending extract row")
	}
}

func TestOrchestrator_BuildsHunkLinks(t *testing.T) {
	repo := initTestGitRepo(t)
	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")
	t.Setenv("TOTALITY_CLOUD_URL", "")

	ctx := context.Background()
	result, err := orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot:  repo,
		Base:      "HEAD~1",
		Head:      "HEAD",
		Tools:     []string{capture.ToolClaude},
		StageOnly: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.HunkLinks) == 0 {
		t.Fatal("expected hunk links for committed file")
	}
	for _, link := range result.HunkLinks {
		if link.HunkID == "" {
			t.Fatal("empty hunk id")
		}
		if link.Authorship != matcher.AuthorshipHuman && link.Authorship != matcher.AuthorshipAgent && link.Authorship != matcher.AuthorshipUnknown {
			t.Fatalf("invalid authorship %q", link.Authorship)
		}
	}
}

func TestOrchestrator_StagesWhenUploadIsUnauthorized(t *testing.T) {
	repo := initTestGitRepo(t)
	totalityHome := t.TempDir()
	t.Setenv("TOTALITY_HOME", totalityHome)
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("TOTALITY_CLOUD_URL", "")
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "bad-token")
	t.Setenv("TOTALITY_API_URL", "http://127.0.0.1:1")

	ctx := context.Background()
	result, err := orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot: repo,
		Base:     "HEAD~1",
		Head:     "HEAD",
		Tools:    []string{capture.ToolClaude},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.StagedExtractID == "" {
		t.Fatal("expected staged extract id")
	}
	if result.UploadError != "" {
		t.Fatalf("UploadError = %q, want empty because upload is deferred to capture sync", result.UploadError)
	}
}

func TestRedactIntegration(t *testing.T) {
	secret := "AKIAIOSFODNN7EXAMPLE"
	out := redact.Redact("key=" + secret)
	if out == "key="+secret {
		t.Fatal("secret not redacted")
	}
}

func initTestGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init", "-b", "main")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-m", "initial")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello\nworld\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-m", "second")
	time.Sleep(10 * time.Millisecond)
	return dir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
