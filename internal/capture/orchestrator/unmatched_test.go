package orchestrator_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/lgtm/internal/capture"
	"github.com/satoricorp/lgtm/internal/capture/orchestrator"
	"github.com/satoricorp/lgtm/internal/storage"
)

// A session that worked in this repository is context worth keeping even when
// none of its text reached a commit.
//
// Matching answers what a session produced. It cannot answer for work that was
// explored and abandoned, and it cannot answer at all for a transcript the
// agent had not finished writing when capture read it — the case that reports
// zero coverage and is indistinguishable from a change nobody used an agent
// for. Binding answers a different question, one that is knowable for every
// session: where did this run?
func TestOrchestrator_StagesAnUnmatchedSessionBoundToThisRepo(t *testing.T) {
	repo := initTestGitRepo(t)
	runGit(t, repo, "remote", "add", "origin", "git@github.com:satoricorp/lgtm.git")
	commitDistinctiveWork(t, repo)

	// The transcript edits a file in this repository, but its content has
	// nothing to do with what was committed, so nothing can match.
	claudeDir := t.TempDir()
	writeUnmatchedTranscript(t, claudeDir, repo, "totally unrelated exploration that was never committed\n")

	t.Setenv("LGTM_HOME", t.TempDir())
	t.Setenv("LGTM_CLOUD_URL", "")
	t.Setenv("LGTM_UPLOAD_TOKEN", "")

	ctx := context.Background()
	result, err := orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot:  repo,
		Base:      "HEAD~1",
		Head:      "HEAD",
		Tools:     []string{capture.ToolClaude},
		ClaudeDir: claudeDir,
		StageOnly: true,
	})
	if err != nil {
		t.Fatalf("orchestrator.Run: %v", err)
	}
	if len(result.HunkLinks) != 0 && result.HunkCoverage != 0 {
		t.Fatalf("fixture matched after all (coverage %.1f%%); it no longer tests the unmatched path", result.HunkCoverage*100)
	}
	if result.StagedSessions == 0 {
		t.Fatal("an unmatched session bound to this repository was not staged")
	}

	origin, _ := stagedBinding(t, ctx)
	if origin != "github.com/satoricorp/lgtm" {
		t.Fatalf("staged session origin = %q, want github.com/satoricorp/lgtm", origin)
	}
}

// The same session, run somewhere else, earns nothing. Without matching as a
// filter, binding is the only thing keeping another repository's sessions out,
// so it has to hold on its own.
func TestOrchestrator_IgnoresAnUnmatchedSessionFromAnotherRepo(t *testing.T) {
	repo := initTestGitRepo(t)
	runGit(t, repo, "remote", "add", "origin", "git@github.com:satoricorp/lgtm.git")
	commitDistinctiveWork(t, repo)

	// A separate checkout, mentioning this repository's path so discovery still
	// finds it — which is exactly how another repository's sessions arrive.
	elsewhere := initTestGitRepo(t)
	runGit(t, elsewhere, "remote", "add", "origin", "git@github.com:joe/side-project.git")

	claudeDir := t.TempDir()
	// Discovery matches on the repository path with a trailing separator, so
	// name a file inside it — which is what a real cross-repo session does when
	// it reads or references another checkout.
	writeUnmatchedTranscript(t, claudeDir, elsewhere,
		"looked at "+filepath.Join(repo, "internal", "app", "binding.go")+" while working elsewhere\n")

	t.Setenv("LGTM_HOME", t.TempDir())
	t.Setenv("LGTM_CLOUD_URL", "")

	ctx := context.Background()
	result, err := orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot:  repo,
		Base:      "HEAD~1",
		Head:      "HEAD",
		Tools:     []string{capture.ToolClaude},
		ClaudeDir: claudeDir,
		StageOnly: true,
	})
	if err != nil {
		t.Fatalf("orchestrator.Run: %v", err)
	}
	// Without this the test proves nothing: if discovery never surfaced the
	// foreign session, it would be absent for reasons that have nothing to do
	// with the binding check.
	if result.EligibleEvents == 0 {
		t.Fatal("the foreign session was never discovered, so this test cannot show the binding check excluded it")
	}
	if result.StagedSessions != 0 {
		t.Fatalf("staged %d session(s) from another repository", result.StagedSessions)
	}
}

func writeUnmatchedTranscript(t *testing.T, claudeDir, cwd, content string) {
	t.Helper()
	slug := strings.ReplaceAll(filepath.ToSlash(filepath.Clean(cwd)), "/", "-")
	dir := filepath.Join(claudeDir, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	line := fmt.Sprintf(
		`{"type":"assistant","cwd":%q,"message":{"role":"assistant","model":"claude","content":[{"type":"tool_use","name":"Write","input":{"file_path":%q,"content":%q}}]},"timestamp":%q,"sessionId":"unmatched-session"}`,
		filepath.ToSlash(cwd), filepath.ToSlash(filepath.Join(cwd, "notes.md")), content, ts,
	)
	if err := os.WriteFile(filepath.Join(dir, "unmatched-session.jsonl"), []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

var _ = storage.StagedSession{}
