package storage_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/satoricorp/totality/internal/storage"
)

func TestCaptureStageRoundTrip(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
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
		RevisionID:  "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID:          "sess-row-1",
		SessionID:   "cursor-abc",
		Tool:        "cursor",
		PayloadJSON: []byte(`{"sessionID":"cursor-abc"}`),
		RawBlob:     []byte(`{"type":"assistant"}` + "\n"),
		SourcePath:  "/tmp/session.jsonl",
		BadLines:    1,
		RevisionID:  "rev-1",
	}); err != nil {
		t.Fatal(err)
	}
	rawSessions, err := stager.RawSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rawSessions) != 1 || string(rawSessions[0].RawBlob) == "" {
		t.Fatalf("raw sessions = %+v", rawSessions)
	}
	extracts, err := stager.PendingExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(extracts) != 1 || string(extracts[0].PayloadJSON) != string(payload) {
		t.Fatalf("pending extracts = %+v", extracts)
	}
	if _, err := stager.MarkExtractShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{Name: "Test", Email: "test@example.com"}); err != nil {
		t.Fatal(err)
	}
	shareable, err := stager.ShareableExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(shareable) != 1 {
		t.Fatalf("shareable extracts = %+v", shareable)
	}
	sessions, err := stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].SessionID != "cursor-abc" {
		t.Fatalf("pending sessions = %+v", sessions)
	}
	if _, err := stager.MarkSessionShareable(ctx, []string{"rev-1"}, storage.CaptureAttestation{Name: "Test", Email: "test@example.com"}); err != nil {
		t.Fatal(err)
	}
	if err := stager.MarkExtractUploaded(ctx, extracts[0].ID); err != nil {
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

// TestStageSessionRestagePreservesUploadStateOnSameContent pins the upsert's
// upload-state rules for pointer rows: restaging identical content keeps the
// row uploaded, restaging changed content (a grown transcript) makes the same
// row pending again.
func TestStageSessionRestagePreservesUploadStateOnSameContent(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	row := storage.StagedSession{
		ID:          storage.SessionSourceRowID("claude", "/tmp/conv.jsonl", "conv"),
		SessionID:   "conv",
		Tool:        "claude",
		PayloadJSON: []byte(`{"sessionID":"conv"}`),
		SourcePath:  "/tmp/conv.jsonl",
		SourceBytes: 10,
		SourceMTime: 111,
		ContentHash: "hash-v1",
		RevisionID:  "rev-1",
	}
	if err := stager.StageSession(ctx, row); err != nil {
		t.Fatal(err)
	}
	if err := stager.MarkSessionUploaded(ctx, row.ID); err != nil {
		t.Fatal(err)
	}

	// Same content: still uploaded, no new pending work.
	if err := stager.StageSession(ctx, row); err != nil {
		t.Fatal(err)
	}
	pending, err := stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 0 {
		t.Fatalf("pending after same-content restage = %+v, want none", pending)
	}

	// Grown content: the one row flips back to pending with the new fingerprint.
	row.ContentHash = "hash-v2"
	row.SourceBytes = 20
	row.SourceMTime = 222
	if err := stager.StageSession(ctx, row); err != nil {
		t.Fatal(err)
	}
	pending, err = stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 1 || pending[0].ID != row.ID {
		t.Fatalf("pending after grown restage = %+v, want row %s", pending, row.ID)
	}
	if pending[0].ContentHash != "hash-v2" || pending[0].SourceBytes != 20 || pending[0].SourceMTime != 222 {
		t.Fatalf("restaged fingerprint = %+v", pending[0])
	}

	if err := stager.SetSessionSourceFingerprint(ctx, row.ID, "hash-v3", 30, 333); err != nil {
		t.Fatal(err)
	}
	pending, err = stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if pending[0].ContentHash != "hash-v3" || pending[0].SourceBytes != 30 || pending[0].SourceMTime != 333 {
		t.Fatalf("fingerprint after SetSessionSourceFingerprint = %+v", pending[0])
	}
}

// TestMarkSessionsShareableForSourcesOfAdoptsLegacyRows covers the backlog the
// row-key change stranded. The key moved from "<revisionID>-<contentHash>" to
// a source-derived "src-<hash>", so re-staging a transcript writes a new row
// and leaves the old one behind: unreachable by id (the id scheme changed) and
// unreachable by revision (its revision is already in pushed history). Every
// such row would stay shareable_at NULL for the rest of time.
func TestMarkSessionsShareableForSourcesOfAdoptsLegacyRows(t *testing.T) {
	t.Setenv("TOTALITY_HOME", t.TempDir())
	ctx := context.Background()
	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sourcePath := filepath.Join(t.TempDir(), "session.jsonl")
	rows := []storage.StagedSession{
		// The legacy row, keyed by a revision that has long since been pushed.
		{
			ID: "revoldaaaaaa-c0ffee", SessionID: "s1", Tool: "claude",
			PayloadJSON: []byte(`{"sessionID":"s1"}`), SourcePath: sourcePath, RevisionID: "revoldaaaaaa",
		},
		// The row this push staged, under the new source-derived key.
		{
			ID: storage.SessionSourceRowID("claude", sourcePath, "s1"), SessionID: "s1", Tool: "claude",
			PayloadJSON: []byte(`{"sessionID":"s1"}`), SourcePath: sourcePath,
		},
		// A different transcript entirely: it must not be swept in.
		{
			ID: "revoldaaaaaa-beefee", SessionID: "other", Tool: "claude",
			PayloadJSON: []byte(`{"sessionID":"other"}`), SourcePath: filepath.Join(t.TempDir(), "other.jsonl"),
		},
		// A Cursor vscdb row has no source path; an empty-string match must
		// never be enough to attest it.
		{
			ID: "vscdb-row", SessionID: "s1", Tool: "claude",
			PayloadJSON: []byte(`{"sessionID":"s1"}`),
		},
	}
	for _, row := range rows {
		if err := stager.StageSession(ctx, row); err != nil {
			t.Fatal(err)
		}
	}

	staged := []string{storage.SessionSourceRowID("claude", sourcePath, "s1")}
	att := storage.CaptureAttestation{Name: "Test", Email: "test@example.com"}
	if _, err := stager.MarkSessionsShareableByID(ctx, staged, att); err != nil {
		t.Fatal(err)
	}
	n, err := stager.MarkSessionsShareableForSourcesOf(ctx, staged, att)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("adopted %d rows, want exactly the legacy row for the same transcript", n)
	}

	shareable, err := stager.ShareableSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, row := range shareable {
		got[row.ID] = true
	}
	if !got["revoldaaaaaa-c0ffee"] {
		t.Fatalf("shareable = %v, want the stranded legacy row adopted", got)
	}
	if got["revoldaaaaaa-beefee"] || got["vscdb-row"] {
		t.Fatalf("shareable = %v, want only rows describing the same transcript source", got)
	}
}
