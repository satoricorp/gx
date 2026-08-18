package hooks

import (
	"os"
	"path/filepath"
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

// TestCloudUploadBlockerCatchesTheUnbakedBinary covers the failure this
// pre-flight exists for. A binary built without the endpoint bake has no cloud
// URL; the upload worker reads that as "cloud disabled" and exits without
// touching the queue, so the outbox grows without bound while every push still
// reports success. The worker is detached with its output and exit code sent to
// /dev/null, so this has to be caught before the handoff or not at all.
func TestCloudUploadBlockerCatchesTheUnbakedBinary(t *testing.T) {
	t.Setenv("GX_CLOUD_URL", "")
	if blocker := cloudUploadBlocker(); blocker == "" {
		t.Fatal(`cloudUploadBlocker() = "" with no cloud endpoint; stranded uploads would go unreported`)
	}

	t.Setenv("GX_CLOUD_URL", "https://api.gx.run")
	if blocker := cloudUploadBlocker(); blocker != "" {
		t.Fatalf("cloudUploadBlocker() = %q with a configured endpoint; a healthy push would warn for nothing", blocker)
	}
}

// TestQueuedUploadBacklogCountsEverythingStranded: the count is what conveys
// scale in the warning — one waiting artifact is a hiccup, two dozen is a
// broken install — so items that already failed have to be in it alongside the
// ones still pending.
func TestQueuedUploadBacklogCountsEverythingStranded(t *testing.T) {
	home := t.TempDir()
	t.Setenv("GX_HOME", home)

	if got := queuedUploadBacklog(); got != 0 {
		t.Fatalf("queuedUploadBacklog() = %d with an empty outbox, want 0", got)
	}

	outboxDir := filepath.Join(home, "publish-outbox")
	if err := os.MkdirAll(outboxDir, 0o700); err != nil {
		t.Fatal(err)
	}
	stranded := map[string][]byte{
		"gx-context-deadbeef.json": []byte(`{"id":"gx-context-deadbeef","status":"pending","created_at":1}`),
		"gx-context-cafebabe.json": []byte(`{"id":"gx-context-cafebabe","status":"failed","attempts":1,"last_error":"upload gx cloud payload: status 500","created_at":2}`),
	}
	for name, body := range stranded {
		if err := os.WriteFile(filepath.Join(outboxDir, name), body, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if got := queuedUploadBacklog(); got != 2 {
		t.Fatalf("queuedUploadBacklog() = %d, want 2 (one pending, one failed)", got)
	}
}
