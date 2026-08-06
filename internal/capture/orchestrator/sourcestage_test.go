package orchestrator_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/orchestrator"
	"github.com/satoricorp/gx/internal/storage"
)

const (
	probeParentSession = "11111111-2222-3333-4444-555555555555"
	probeFileAContent  = "package alpha\n\nfunc AlphaOne() int {\n\treturn 41\n}\n"
	probeFileBContent  = "package beta\n\nfunc BetaTwo() int {\n\treturn 42\n}\n"
)

// initProbeRepo commits two files in one commit so the push range holds two
// distinct hunks, one for the parent transcript and one for the subagent.
func initProbeRepo(t *testing.T) string {
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
	for name, content := range map[string]string{
		"alpha.go": probeFileAContent,
		"beta.go":  probeFileBContent,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runGit(t, dir, "add", "alpha.go", "beta.go")
	runGit(t, dir, "commit", "-m", "add alpha and beta")
	time.Sleep(10 * time.Millisecond)
	return dir
}

// probeEditLine renders one Claude transcript line writing content to path.
// Every line carries the PARENT conversation's sessionId, as Claude does even
// inside subagent transcripts.
func probeEditLine(repoRoot, file, content string) string {
	line := map[string]any{
		"type":      "assistant",
		"cwd":       repoRoot,
		"sessionId": probeParentSession,
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
	}
	data, err := json.Marshal(line)
	if err != nil {
		panic(err)
	}
	return string(data) + "\n"
}

func writeTranscript(t *testing.T, path string, lines ...string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "")), 0o644); err != nil {
		t.Fatal(err)
	}
}

func repoClaudeSlug(repoRoot string) string {
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

func runProbe(t *testing.T, ctx context.Context, repo, claudeDir string) orchestrator.Result {
	t.Helper()
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
	return result
}

func pendingByID(t *testing.T, ctx context.Context, stager storage.CaptureStager) map[string]storage.StagedSession {
	t.Helper()
	rows, err := stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]storage.StagedSession{}
	for _, row := range rows {
		out[row.ID] = row
	}
	return out
}

// TestSourceStaging_PointerRowsForMatchedSourcesOnly demonstrates the pointer
// staging model end to end:
//
//	(a) sessions are staged only after matching — an unmatched transcript
//	    stages nothing;
//	(b) one row per source file, with a subagent transcript distinct from its
//	    parent even though both carry the parent's in-file sessionId;
//	(c) re-staging a grown source updates the row in place: the row count
//	    stays stable and the changed row becomes pending again.
func TestSourceStaging_PointerRowsForMatchedSourcesOnly(t *testing.T) {
	repo := initProbeRepo(t)
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()

	claudeDir := filepath.Join(t.TempDir(), "projects")
	slugDir := filepath.Join(claudeDir, repoClaudeSlug(repo))
	parentPath := filepath.Join(slugDir, probeParentSession+".jsonl")
	subagentPath := filepath.Join(slugDir, probeParentSession, "subagents", "agent-worker.jsonl")
	unmatchedPath := filepath.Join(slugDir, "99999999-8888-7777-6666-555555555555.jsonl")

	// Parent edits alpha.go, its subagent edits beta.go, and a third session
	// edits a file this push never touched.
	writeTranscript(t, parentPath, probeEditLine(repo, "alpha.go", probeFileAContent))
	writeTranscript(t, subagentPath, probeEditLine(repo, "beta.go", probeFileBContent))
	writeTranscript(t, unmatchedPath, probeEditLine(repo, "gamma.go", "package gamma\n\nfunc Gamma() int {\n\treturn 43\n}\n"))

	result := runProbe(t, ctx, repo, claudeDir)
	if result.StagedSessions != 2 {
		t.Fatalf("StagedSessions = %d, want 2 (parent + subagent, not the unmatched session)", result.StagedSessions)
	}

	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	rows := pendingByID(t, ctx, stager)
	if len(rows) != 2 {
		t.Fatalf("pending session rows = %d, want 2: %+v", len(rows), rows)
	}

	var parentRow, subagentRow storage.StagedSession
	for _, row := range rows {
		switch row.SourcePath {
		case parentPath:
			parentRow = row
		case subagentPath:
			subagentRow = row
		default:
			t.Fatalf("unexpected staged source %q (unmatched sessions must not stage)", row.SourcePath)
		}
	}
	if parentRow.ID == "" || subagentRow.ID == "" {
		t.Fatalf("missing parent or subagent row: %+v", rows)
	}
	if parentRow.ID == subagentRow.ID {
		t.Fatal("parent and subagent staged into the same row despite sharing a sessionId")
	}
	if parentRow.SessionID != probeParentSession {
		t.Fatalf("parent SessionID = %q, want %q", parentRow.SessionID, probeParentSession)
	}
	wantSubagentID := probeParentSession + ":subagent:agent-worker"
	if subagentRow.SessionID != wantSubagentID {
		t.Fatalf("subagent SessionID = %q, want %q", subagentRow.SessionID, wantSubagentID)
	}
	for name, row := range map[string]storage.StagedSession{"parent": parentRow, "subagent": subagentRow} {
		if len(row.RawBlob) != 0 {
			t.Fatalf("%s row stored %d raw bytes, want a pointer only", name, len(row.RawBlob))
		}
		if row.ContentHash == "" || row.SourceBytes == 0 || row.SourceMTime == 0 {
			t.Fatalf("%s row missing fingerprint: hash=%q bytes=%d mtime=%d", name, row.ContentHash, row.SourceBytes, row.SourceMTime)
		}
		info, err := os.Stat(row.SourcePath)
		if err != nil {
			t.Fatalf("%s source: %v", name, err)
		}
		if row.SourceBytes != info.Size() {
			t.Fatalf("%s fingerprint bytes = %d, want %d", name, row.SourceBytes, info.Size())
		}
	}

	// Re-running the same push replaces rows instead of accumulating.
	runProbe(t, ctx, repo, claudeDir)
	again := pendingByID(t, ctx, stager)
	if len(again) != 2 {
		t.Fatalf("pending rows after re-run = %d, want 2", len(again))
	}
	if _, ok := again[parentRow.ID]; !ok {
		t.Fatalf("parent row id changed across identical re-runs: %v", keysOf(again))
	}

	// A grown source re-stages in place: same row id, fresh fingerprint, and
	// the row uploads again while the untouched subagent row stays uploaded.
	if err := stager.MarkSessionUploaded(ctx, parentRow.ID); err != nil {
		t.Fatal(err)
	}
	if err := stager.MarkSessionUploaded(ctx, subagentRow.ID); err != nil {
		t.Fatal(err)
	}
	grown := probeEditLine(repo, "alpha.go", probeFileAContent) + probeEditLine(repo, "alpha.go", probeFileAContent)
	writeTranscript(t, parentPath, grown)

	runProbe(t, ctx, repo, claudeDir)
	afterGrowth := pendingByID(t, ctx, stager)
	if len(afterGrowth) != 1 {
		t.Fatalf("pending rows after growth = %v, want only the grown parent row", keysOf(afterGrowth))
	}
	updated, ok := afterGrowth[parentRow.ID]
	if !ok {
		t.Fatalf("grown source staged a new row %v instead of updating %s", keysOf(afterGrowth), parentRow.ID)
	}
	if updated.ContentHash == parentRow.ContentHash {
		t.Fatal("grown source kept the old content hash")
	}
	if updated.SourceBytes <= parentRow.SourceBytes {
		t.Fatalf("grown source bytes = %d, want > %d", updated.SourceBytes, parentRow.SourceBytes)
	}
}

func keysOf(rows map[string]storage.StagedSession) []string {
	out := make([]string, 0, len(rows))
	for id := range rows {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// TestIngestRawSessionStagesPointer covers the transcript-hook entry point:
// a single ingested transcript stages one pointer row, keyed by source, and
// re-ingesting the same (grown) file updates it in place.
func TestIngestRawSessionStagesPointer(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}

	repo := t.TempDir()
	dir := t.TempDir()
	path := filepath.Join(dir, probeParentSession, "subagents", "agent-scout.jsonl")
	writeTranscript(t, path, probeEditLine(repo, "alpha.go", probeFileAContent))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	session, badLines, err := orchestrator.IngestRawSession(ctx, stager, capture.ToolClaude, path, raw, repo, nil)
	if err != nil {
		t.Fatal(err)
	}
	if badLines != 0 {
		t.Fatalf("badLines = %d", badLines)
	}
	wantID := probeParentSession + ":subagent:agent-scout"
	if session.SessionID != wantID {
		t.Fatalf("SessionID = %q, want %q (path identity, not the in-file parent id)", session.SessionID, wantID)
	}

	rows := pendingByID(t, ctx, stager)
	if len(rows) != 1 {
		t.Fatalf("pending rows = %d, want 1", len(rows))
	}
	var row storage.StagedSession
	for _, r := range rows {
		row = r
	}
	if row.SourcePath != path || len(row.RawBlob) != 0 {
		t.Fatalf("row = %+v, want pointer to %s with no raw blob", row, path)
	}

	grown := append(raw, []byte(probeEditLine(repo, "beta.go", probeFileBContent))...)
	if err := os.WriteFile(path, grown, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := orchestrator.IngestRawSession(ctx, stager, capture.ToolClaude, path, grown, repo, nil); err != nil {
		t.Fatal(err)
	}
	rows = pendingByID(t, ctx, stager)
	if len(rows) != 1 {
		t.Fatalf("pending rows after re-ingest = %d, want the same single row", len(rows))
	}
	for id, r := range rows {
		if id != row.ID {
			t.Fatalf("re-ingest created row %s, want update of %s", id, row.ID)
		}
		if r.SourceBytes != int64(len(grown)) {
			t.Fatalf("SourceBytes = %d, want %d", r.SourceBytes, len(grown))
		}
	}
}
