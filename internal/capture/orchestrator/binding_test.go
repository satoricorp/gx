package orchestrator_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/totality/internal/capture"
	"github.com/satoricorp/totality/internal/capture/orchestrator"
	"github.com/satoricorp/totality/internal/storage"
)

// A staged session records which repository it belongs to, independent of
// whether any of its work matched a commit. That binding is what a gate can
// act on later: matching answers what a session produced, not what it belongs
// to, and it cannot answer for a session whose work was never committed or
// whose transcript had not finished being written.
func TestOrchestrator_RecordsTheRepositoryASessionBelongsTo(t *testing.T) {
	repo := initTestGitRepo(t)
	runGit(t, repo, "remote", "add", "origin", "git@github.com:satoricorp/totality.git")
	added := commitDistinctiveWork(t, repo)

	claudeDir := t.TempDir()
	writeBoundTranscript(t, claudeDir, repo, "internal/app/binding.go", added)

	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TOTALITY_CLOUD_URL", "")
	t.Setenv("TOTALITY_UPLOAD_TOKEN", "")

	ctx := context.Background()
	if _, err := orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot:  repo,
		Base:      "HEAD~1",
		Head:      "HEAD",
		Tools:     []string{capture.ToolClaude},
		ClaudeDir: claudeDir,
		StageOnly: true,
	}); err != nil {
		t.Fatalf("orchestrator.Run: %v", err)
	}

	origin, cwd := stagedBinding(t, ctx)
	if origin != "github.com/satoricorp/totality" {
		t.Fatalf("session_origin = %q, want github.com/satoricorp/totality", origin)
	}
	if cwd == "" {
		t.Fatal("session_cwd was not recorded")
	}
	// The raw remote never reaches the database; only the normalized identity,
	// which has any credentials stripped.
	if strings.Contains(origin, "git@") || strings.Contains(origin, ".git") {
		t.Fatalf("session_origin stored a raw remote URL: %q", origin)
	}
}

// A repository with no remote is unidentifiable, and unidentifiable must mean
// unbound rather than "matches anything". This is the personal-side-project
// case: its sessions have to stay local.
func TestOrchestrator_LeavesSessionsInRemotelessReposUnbound(t *testing.T) {
	repo := initTestGitRepo(t) // no origin added
	added := commitDistinctiveWork(t, repo)
	claudeDir := t.TempDir()
	writeBoundTranscript(t, claudeDir, repo, "internal/app/binding.go", added)

	t.Setenv("TOTALITY_HOME", t.TempDir())
	t.Setenv("TOTALITY_CLOUD_URL", "")

	ctx := context.Background()
	if _, err := orchestrator.Run(ctx, orchestrator.RunOptions{
		RepoRoot:  repo,
		Base:      "HEAD~1",
		Head:      "HEAD",
		Tools:     []string{capture.ToolClaude},
		ClaudeDir: claudeDir,
		StageOnly: true,
	}); err != nil {
		t.Fatalf("orchestrator.Run: %v", err)
	}

	if origin, _ := stagedBinding(t, ctx); origin != "" {
		t.Fatalf("a repository with no remote bound to %q", origin)
	}
}

// commitDistinctiveWork adds a commit whose content is long and specific
// enough for content matching to link it, so the test exercises the staging
// path rather than the "nothing matched" path.
func commitDistinctiveWork(t *testing.T, repo string) string {
	t.Helper()
	added := "package app\n\n" +
		"// ResolveBindingForSession maps a session working directory onto the\n" +
		"// repository identity that owns it, so attribution does not depend on\n" +
		"// whether the work was ever committed.\n" +
		"func ResolveBindingForSession(dir string) string { return dir }\n"
	path := filepath.Join(repo, "internal", "app", "binding.go")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(added), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "internal/app/binding.go")
	runGit(t, repo, "commit", "-m", "add binding resolver")
	return added
}

func stagedBinding(t *testing.T, ctx context.Context) (origin, cwd string) {
	t.Helper()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	rows, err := db.QueryContext(ctx, `SELECT COALESCE(session_origin,''), COALESCE(session_cwd,'') FROM capture_sessions ORDER BY created_at DESC`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	found := false
	for rows.Next() {
		var o, c string
		if err := rows.Scan(&o, &c); err != nil {
			t.Fatal(err)
		}
		found = true
		if o != "" || c != "" {
			return o, c
		}
	}
	if !found {
		t.Fatal("no staged sessions")
	}
	return "", ""
}

// writeBoundTranscript writes a Claude transcript whose entries carry the cwd,
// which is how a session names the directory it ran in.
func writeBoundTranscript(t *testing.T, claudeDir, cwd, relFile, content string) {
	t.Helper()
	slug := strings.ReplaceAll(filepath.ToSlash(filepath.Clean(cwd)), "/", "-")
	dir := filepath.Join(claudeDir, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	line := fmt.Sprintf(
		`{"type":"assistant","cwd":%q,"message":{"role":"assistant","model":"claude","content":[{"type":"tool_use","name":"Write","input":{"file_path":%q,"content":%q}}]},"timestamp":%q,"sessionId":"bound-session"}`,
		filepath.ToSlash(cwd), filepath.ToSlash(filepath.Join(cwd, relFile)), content, ts,
	)
	if err := os.WriteFile(filepath.Join(dir, "bound-session.jsonl"), []byte(line+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}
