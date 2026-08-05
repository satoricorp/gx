package hooks_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/lgtm/internal/auth"
	"github.com/satoricorp/lgtm/internal/capture/extract"
	"github.com/satoricorp/lgtm/internal/hooks"
	"github.com/satoricorp/lgtm/internal/storage"
	"github.com/satoricorp/lgtm/internal/telemetry"
)

const shareableProbeSession = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
const shareableProbeContent = "package alpha\n\nfunc AlphaOne() int {\n\treturn 41\n}\n"

// TestRunPushMarksRowsItStagedWithoutRevisionIDs covers a repo where revision
// resolution returns nothing for the pushed range — no commit carries a lgtm
// trailer, so `revision_id IN (...)` matches nothing. The rows this very push
// staged must still become shareable and drain through SyncPending. Under the
// revision-only gate they stayed local forever.
//
// (Note for anyone reading this as a diagnosis of the field failure: a repo
// with the lgtm prepare-commit-msg hook installed does stamp trailers on plain
// `git commit`, and marking by revision demonstrably did work for some rows.
// This test pins the case where it cannot, which is the case the id gate
// exists for — not a claim about why any particular row was stuck.)
func TestRunPushMarksRowsItStagedWithoutRevisionIDs(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	t.Setenv("LGTM_DISABLE_BACKGROUND_WORKERS", "1")
	t.Setenv(hooks.SuppressAdoptedPublicationEnv, "1")

	home := t.TempDir()
	repo := initPlainGitRepo(t)
	transcript := filepath.Join(home, ".claude", "projects", claudeProjectSlug(repo), shareableProbeSession+".jsonl")
	writeClaudeProbeTranscript(t, transcript, repo, "alpha.go", shareableProbeContent)

	ctx := context.Background()
	base := gitRev(t, repo, "HEAD~1")
	head := gitRev(t, repo, "HEAD")
	outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		RefRange: base + ".." + head,
		HeadSHA:  head,
		HomeDir:  home,
		Tools:    []string{"claude"},
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if outcome.CaptureError != "" {
		t.Fatalf("CaptureError = %q", outcome.CaptureError)
	}
	// The premise: there is no revision to match on, and no ref-range gate either.
	if len(outcome.RevisionIDs) != 0 {
		t.Fatalf("RevisionIDs = %v, want none: this repo has no lgtm trailers", outcome.RevisionIDs)
	}
	if len(outcome.Result.StagedSessionIDs) != 1 {
		t.Fatalf("StagedSessionIDs = %v, want the one matched transcript", outcome.Result.StagedSessionIDs)
	}
	if len(outcome.Result.StagedExtractIDs) != 1 {
		t.Fatalf("StagedExtractIDs = %v, want one extract row", outcome.Result.StagedExtractIDs)
	}
	if outcome.ShareableExtract != 1 || outcome.ShareableSession != 1 {
		t.Fatalf("shareable = %d extracts / %d sessions, want 1 each",
			outcome.ShareableExtract, outcome.ShareableSession)
	}

	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sessions, err := stager.ShareableSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].ID != outcome.Result.StagedSessionIDs[0] {
		t.Fatalf("shareable sessions = %+v, want row %v", sessions, outcome.Result.StagedSessionIDs)
	}
	if sessions[0].RevisionID != "" {
		t.Fatalf("session revision id = %q, want empty: the marking must not have come from a revision match", sessions[0].RevisionID)
	}
	if sessions[0].ShareableAt == 0 || sessions[0].AttestedAt == 0 {
		t.Fatalf("session attestation = %+v, want shareable_at and attested_at stamped", sessions[0])
	}
	extracts, err := stager.ShareableExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(extracts) != 1 || extracts[0].ID != outcome.Result.StagedExtractIDs[0] {
		t.Fatalf("shareable extracts = %+v, want row %v", extracts, outcome.Result.StagedExtractIDs)
	}

	// And the drain the marking exists to unblock actually runs.
	var sessionBodies []string
	var extractCalls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch r.URL.Path {
		case "/v1/sessions":
			var req struct {
				SessionId string `json:"sessionId"`
				Content   string `json:"content"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// The shipped uploader's 400s were "sessionId, tool, and content
			// are required" — assert the body the server demands.
			if req.SessionId == "" || req.Content == "" {
				http.Error(w, `{"error":"sessionId, tool, and content are required"}`, http.StatusBadRequest)
				return
			}
			sessionBodies = append(sessionBodies, req.Content)
		case "/v1/extracts":
			extractCalls++
		default:
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	result, err := extract.SyncPending(ctx, stager, auth.Credentials{APIURL: srv.URL, Token: "t"}, telemetry.NopClient{})
	if err != nil {
		t.Fatal(err)
	}
	if result.SessionsUploaded != 1 || result.ExtractsUploaded != 1 {
		t.Fatalf("sync result = %+v, want one session and one extract uploaded", result)
	}
	if extractCalls != 1 {
		t.Fatalf("POST /v1/extracts calls = %d, want 1", extractCalls)
	}
	if len(sessionBodies) != 1 || !strings.Contains(sessionBodies[0], "alpha.go") {
		t.Fatalf("uploaded session content = %q, want the transcript bytes", sessionBodies)
	}
}

// initPlainGitRepo builds a repo whose commits carry no lgtm revision trailer,
// which is what the plain-Git workflow produces.
func initPlainGitRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	for _, args := range [][]string{
		{"init", "-b", "main"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test"},
	} {
		runGit(t, repo, args...)
	}
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "initial")
	if err := os.WriteFile(filepath.Join(repo, "alpha.go"), []byte(shareableProbeContent), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "alpha.go")
	runGit(t, repo, "commit", "-m", "add alpha")
	return repo
}

func claudeProjectSlug(repoRoot string) string {
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		abs = repoRoot
	}
	slug := strings.ReplaceAll(abs, string(filepath.Separator), "-")
	if !strings.HasPrefix(slug, "-") {
		slug = "-" + slug
	}
	return slug
}

// writeClaudeProbeTranscript writes a one-line Claude transcript whose Write
// tool call produced exactly the committed file, so the matcher links it.
func writeClaudeProbeTranscript(t *testing.T, path, repoRoot, file, content string) {
	t.Helper()
	line, err := json.Marshal(map[string]any{
		"type":      "assistant",
		"cwd":       repoRoot,
		"sessionId": shareableProbeSession,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"message": map[string]any{
			"role":  "assistant",
			"model": "claude-probe",
			"content": []map[string]any{{
				"type": "tool_use",
				"id":   "t1",
				"name": "Write",
				"input": map[string]any{
					"file_path": filepath.ToSlash(filepath.Join(repoRoot, file)),
					"content":   content,
				},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(line, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
