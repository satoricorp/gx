package hooks

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestShouldStartOutboxWorkerRetriesBacklogWithoutNewEnqueue covers the retry
// path that used to belong to `gx sync`: a push that enqueues nothing must
// still spawn the outbox worker when failed or pending items are sitting in
// the outbox from an earlier push.
func TestShouldStartOutboxWorkerRetriesBacklogWithoutNewEnqueue(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	if shouldStartOutboxWorker(false) {
		t.Fatal("shouldStartOutboxWorker(false) = true with an empty outbox")
	}
	if !shouldStartOutboxWorker(true) {
		t.Fatal("shouldStartOutboxWorker(true) = false; a fresh enqueue must spawn the worker")
	}

	outboxDir := filepath.Join(home, "publish-outbox")
	if err := os.MkdirAll(outboxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	failedItem := []byte(`{"id":"gx-context-deadbeef","status":"failed","attempts":1,"last_error":"upload gx cloud payload: status 500","created_at":1}`)
	if err := os.WriteFile(filepath.Join(outboxDir, "gx-context-deadbeef.json"), failedItem, 0o600); err != nil {
		t.Fatal(err)
	}
	if !shouldStartOutboxWorker(false) {
		t.Fatal("shouldStartOutboxWorker(false) = false with a failed item; failed uploads would never retry")
	}

	if err := os.Remove(filepath.Join(outboxDir, "gx-context-deadbeef.json")); err != nil {
		t.Fatal(err)
	}
	pendingItem := []byte(`{"id":"gx-context-cafebabe","status":"pending","created_at":2}`)
	if err := os.WriteFile(filepath.Join(outboxDir, "gx-context-cafebabe.json"), pendingItem, 0o600); err != nil {
		t.Fatal(err)
	}
	if !shouldStartOutboxWorker(false) {
		t.Fatal("shouldStartOutboxWorker(false) = false with a pending item; leftover pending uploads would never retry")
	}
}

// TestRunPushOnALocalBuildPublishesNothing: with no gx Cloud configured a
// push used to queue its review artifact anyway and then warn, on every push,
// that it would never upload. The outbox nothing drains grew without bound.
// Capture is the whole of a local push.
func TestRunPushOnALocalBuildPublishesNothing(t *testing.T) {
	repoRoot := t.TempDir()
	home := t.TempDir()
	t.Setenv("GX_HOME", home)
	t.Setenv("GX_CLOUD_URL", "")
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")
	initGitRepo(t, repoRoot)

	head := strings.TrimSpace(gitOutputInHookTest(t, repoRoot, "rev-parse", "HEAD"))
	outcome, err := RunPush(context.Background(), PushOptions{
		RepoRoot:    repoRoot,
		Remote:      "origin",
		LocalRef:    "refs/heads/main",
		HeadSHA:     head,
		SkipCapture: true,
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if outcome.Publication.Queued || outcome.PublicationError != "" {
		t.Fatalf("publication = %+v (error %q), want none attempted with no cloud", outcome.Publication, outcome.PublicationError)
	}
	entries, err := os.ReadDir(filepath.Join(home, "publish-outbox"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("publish-outbox holds %d item(s) after a local push, want none", len(entries))
	}
}
