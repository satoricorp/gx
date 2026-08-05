package hooks

import (
	"os"
	"path/filepath"
	"testing"
)

// TestShouldStartOutboxWorkerRetriesBacklogWithoutNewEnqueue covers the retry
// path that used to belong to `lgtm sync`: a push that enqueues nothing must
// still spawn the outbox worker when failed or pending items are sitting in
// the outbox from an earlier push.
func TestShouldStartOutboxWorkerRetriesBacklogWithoutNewEnqueue(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)

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
	failedItem := []byte(`{"id":"lgtm-context-deadbeef","status":"failed","attempts":1,"last_error":"upload lgtm cloud payload: status 500","created_at":1}`)
	if err := os.WriteFile(filepath.Join(outboxDir, "lgtm-context-deadbeef.json"), failedItem, 0o600); err != nil {
		t.Fatal(err)
	}
	if !shouldStartOutboxWorker(false) {
		t.Fatal("shouldStartOutboxWorker(false) = false with a failed item; failed uploads would never retry")
	}

	if err := os.Remove(filepath.Join(outboxDir, "lgtm-context-deadbeef.json")); err != nil {
		t.Fatal(err)
	}
	pendingItem := []byte(`{"id":"lgtm-context-cafebabe","status":"pending","created_at":2}`)
	if err := os.WriteFile(filepath.Join(outboxDir, "lgtm-context-cafebabe.json"), pendingItem, 0o600); err != nil {
		t.Fatal(err)
	}
	if !shouldStartOutboxWorker(false) {
		t.Fatal("shouldStartOutboxWorker(false) = false with a pending item; leftover pending uploads would never retry")
	}
}
