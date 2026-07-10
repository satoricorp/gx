package vcs

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/commitcontext"
	"github.com/satoricorp/gx/internal/storage"
)

func TestStagedCommitPersistsAgentDeclaredSelfReport(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "work.txt"), []byte("work\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "work.txt")

	t.Setenv("GX_SESSION_ID", "mcp-session-1")
	result, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{
		Message: "work change",
		SelfReport: commitcontext.SelfReport{
			TaskSummary: "add commit surface",
			CommandsRun: []string{"go test ./..."},
			TestsRun:    []string{"internal/vcs"},
		},
	})
	if err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}
	trailer := RevisionTrailerLine(result.Change.ChangeID)
	gitMessage, err := gitCommitMessage(t, root, "HEAD")
	if err != nil {
		t.Fatalf("gitCommitMessage() error = %v", err)
	}
	if !strings.Contains(gitMessage, trailer) {
		t.Fatalf("git commit message = %q, want trailer %q", gitMessage, trailer)
	}

	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	defer db.Close()

	var agentTool, provider, modelID string
	var source sql.NullString
	var sessionID string
	if err := db.QueryRowContext(context.Background(), `
		SELECT session_id, agent_tool, provider, model_id, source
		FROM change_session_provenance
		WHERE agent_tool = ?
	`, commitcontext.AgentDeclaredTool).Scan(&sessionID, &agentTool, &provider, &modelID, &source); err != nil {
		t.Fatalf("select change_session_provenance: %v", err)
	}
	if sessionID != "mcp-session-1" {
		t.Fatalf("session_id = %q, want mcp-session-1", sessionID)
	}
	if agentTool != commitcontext.AgentDeclaredTool || provider != commitcontext.AgentDeclaredProvider || modelID != commitcontext.AgentDeclaredModelID {
		t.Fatalf("provenance = %s/%s/%s", agentTool, provider, modelID)
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(source.String), &payload); err != nil {
		t.Fatalf("unmarshal provenance source: %v", err)
	}
	if payload["task_summary"] != "add commit surface" {
		t.Fatalf("task_summary = %#v", payload["task_summary"])
	}
}

func TestCommitContextMissingFileStillAllowsCommit(t *testing.T) {
	svc, root := setupStagedCommitRepo(t, "main")
	if err := os.WriteFile(filepath.Join(root, "work.txt"), []byte("work\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, root, "add", "work.txt")

	_, err := svc.RecordStagedRevision(context.Background(), StagedRevisionOptions{Message: "work change"})
	if err != nil {
		t.Fatalf("RecordStagedRevision() error = %v", err)
	}

	missingPath := filepath.Join(root, "missing-context.json")
	_, loadErr := commitcontext.LoadFile(missingPath)
	if loadErr == nil {
		t.Fatal("LoadFile() error = nil, want missing file error")
	}
	if !strings.Contains(loadErr.Error(), "commit context file not found") {
		t.Fatalf("error = %q, want missing file message", loadErr)
	}
}
