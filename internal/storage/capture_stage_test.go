package storage_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestCaptureStageRoundTrip(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	payload := []byte(`{"refRange":"main..HEAD"}`)
	if err := stager.StageExtract(ctx, storage.StagedExtract{
		ID:          "ext-1",
		RepoRoot:    filepath.Join(t.TempDir(), "repo"),
		RefRange:    "main..HEAD",
		PayloadJSON: payload,
	}); err != nil {
		t.Fatal(err)
	}
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID:          "sess-row-1",
		SessionID:   "cursor-abc",
		Tool:        "cursor",
		PayloadJSON: []byte(`{"sessionID":"cursor-abc"}`),
	}); err != nil {
		t.Fatal(err)
	}
	extracts, err := stager.PendingExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(extracts) != 1 || string(extracts[0].PayloadJSON) != string(payload) {
		t.Fatalf("pending extracts = %+v", extracts)
	}
	sessions, err := stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].SessionID != "cursor-abc" {
		t.Fatalf("pending sessions = %+v", sessions)
	}
	if err := stager.MarkExtractUploaded(ctx, "ext-1"); err != nil {
		t.Fatal(err)
	}
	extracts, err = stager.PendingExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(extracts) != 0 {
		t.Fatalf("expected no pending extracts after upload, got %d", len(extracts))
	}
}
