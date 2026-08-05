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
)

// An agent works in a linked worktree and the push runs from the main
// checkout. That is not an exotic setup — it is how every agent session in
// this repository runs — and capture used to score it at zero coverage with no
// error anywhere, because session paths were made relative to the pushing
// checkout:
//
//	session   .claude/worktrees/feature/internal/cli/x.go
//	commit    internal/cli/x.go
//
// The commits are identical either way; a worktree shares the object store. So
// a run that finds the transcript, parses every event out of it, and then
// links nothing is reporting a path bug as an absence of work.
func TestOrchestrator_AttributesEditsMadeInALinkedWorktree(t *testing.T) {
	main := initTestGitRepo(t)
	worktree := filepath.Join(main, ".claude", "worktrees", "feature")
	runGit(t, main, "worktree", "add", "-q", "-b", "feature", worktree)

	// Commit through the worktree, exactly as an agent session would.
	added := "package cli\n\nfunc WorktreeChange() string { return \"from a worktree\" }\n"
	writeFile(t, filepath.Join(worktree, "internal", "cli", "x.go"), added)
	runGit(t, worktree, "add", "internal/cli/x.go")
	runGit(t, worktree, "commit", "-m", "add x.go from the worktree")

	// A transcript that edited the file through its worktree path — the only
	// path the agent ever saw.
	claudeDir := t.TempDir()
	writeWorktreeTranscript(t, claudeDir, worktree, filepath.Join(worktree, "internal", "cli", "x.go"), added)

	t.Setenv("LGTM_HOME", t.TempDir())
	t.Setenv("LGTM_CLOUD_URL", "")
	t.Setenv("LGTM_UPLOAD_TOKEN", "")

	// Push from the main checkout, which is where `git push` runs.
	result, err := orchestrator.Run(context.Background(), orchestrator.RunOptions{
		RepoRoot:  main,
		Base:      "HEAD~1",
		Head:      "feature",
		Tools:     []string{capture.ToolClaude},
		ClaudeDir: claudeDir,
		StageOnly: false,
	})
	if err != nil {
		t.Fatalf("orchestrator.Run: %v", err)
	}

	if result.EligibleHunks == 0 {
		t.Fatalf("no eligible hunks; the test committed nothing to attribute")
	}
	if result.EligibleEvents == 0 {
		t.Fatalf("no events parsed from the worktree transcript; discovery or parsing regressed, not attribution")
	}
	if len(result.HunkLinks) == 0 {
		t.Fatalf("a worktree edit produced no hunk links (coverage %.1f%%): session paths did not match committed paths",
			result.HunkCoverage*100)
	}
	if result.HunkCoverage == 0 {
		t.Fatalf("hunk coverage is 0%% for a change made entirely in a worktree")
	}
	if result.StagedSessions == 0 {
		t.Fatalf("no sessions staged for a change made entirely in a worktree")
	}
}

// writeWorktreeTranscript writes a Claude transcript whose Write tool call
// names the file by its absolute path inside the worktree.
func writeWorktreeTranscript(t *testing.T, claudeDir, sessionCwd, absFile, content string) {
	t.Helper()
	// Claude files transcripts under a slug of the session's cwd, which for an
	// agent in a worktree is the worktree, not the repository.
	slug := strings.ReplaceAll(filepath.ToSlash(filepath.Clean(sessionCwd)), "/", "-")
	dir := filepath.Join(claudeDir, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	line := fmt.Sprintf(
		`{"type":"assistant","message":{"role":"assistant","model":"claude","content":[{"type":"tool_use","name":"Write","input":{"file_path":%q,"content":%q}}]},"timestamp":%q,"sessionId":"worktree-session"}`,
		filepath.ToSlash(absFile), content, ts,
	)
	if err := os.WriteFile(filepath.Join(dir, "worktree-session.jsonl"), []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
