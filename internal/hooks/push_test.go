package hooks_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/hooks"
	"github.com/satoricorp/gx/internal/publication"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

func TestRunPushAlwaysSucceedsWithNetworkUnavailable(t *testing.T) {
	repo := initPushHookRepo(t)
	base := gitRev(t, repo, "HEAD~1")
	head := gitRev(t, repo, "HEAD")
	refRange := base + ".." + head

	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("GX_API_URL", "http://127.0.0.1:1")
	t.Setenv("GX_TOKEN", "test-token")

	outcome, err := hooks.RunPush(context.Background(), hooks.PushOptions{
		RepoRoot:    repo,
		Remote:      "origin",
		RefRange:    refRange,
		HeadSHA:     head,
		SkipCapture: true,
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v, want nil so hook never blocks", err)
	}
	if len(outcome.RevisionIDs) != 1 || outcome.RevisionIDs[0] != "revpushaaaaa" {
		t.Fatalf("RevisionIDs = %v, want [revpushaaaaa]", outcome.RevisionIDs)
	}
}

func TestRunPushMarksOnlyPushedRevisionsShareable(t *testing.T) {
	repo := initPushHookRepo(t)
	ctx := context.Background()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)

	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	payload := []byte(`{"refRange":"main..HEAD"}`)
	for _, revisionID := range []string{"revpushaaaaa", "revnotpushed"} {
		if err := stager.StageExtract(ctx, storage.StagedExtract{
			RepoRoot:    repo,
			RefRange:    "main..HEAD",
			PayloadJSON: payload,
			RevisionID:  revisionID,
		}); err != nil {
			t.Fatal(err)
		}
	}

	base := gitRev(t, repo, "HEAD~1")
	head := gitRev(t, repo, "HEAD")
	if _, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot:    repo,
		RefRange:    base + ".." + head,
		HeadSHA:     head,
		SkipCapture: true,
	}); err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}

	stager, err = storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	shareable, err := stager.ShareableExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(shareable) != 1 || shareable[0].RevisionID != "revpushaaaaa" {
		t.Fatalf("shareable extracts = %+v, want only revpushaaaaa", shareable)
	}
	pending, err := stager.PendingExtracts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 2 {
		t.Fatalf("pending extracts = %d, want 2", len(pending))
	}
}

func TestRunPushIdempotentCaptureRows(t *testing.T) {
	ctx := context.Background()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")

	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatal(err)
	}
	payload := []byte(`{"sessionID":"s1"}`)
	row := storage.StagedSession{
		SessionID:   "s1",
		Tool:        "cursor",
		PayloadJSON: payload,
		RevisionID:  "revpushaaaaa",
	}
	if err := stager.StageSession(ctx, row); err != nil {
		t.Fatal(err)
	}
	firstID := storage.CaptureRowID(row.RevisionID, storage.PayloadContentHash(payload))
	if err := stager.StageSession(ctx, row); err != nil {
		t.Fatal(err)
	}
	sessions, err := stager.PendingSessions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 || sessions[0].ID != firstID {
		t.Fatalf("sessions = %+v, want one idempotent row %s", sessions, firstID)
	}
}

func TestEnqueueAdoptedPublicationIsIdempotent(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")
	repo := t.TempDir()
	head := "abc123"
	opts := hooks.AdoptPushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		LocalRef: "refs/heads/feature/demo",
		HeadSHA:  head,
	}
	first, err := hooks.EnqueueAdoptedPublication(context.Background(), opts)
	if err != nil {
		t.Fatalf("EnqueueAdoptedPublication() error = %v", err)
	}
	second, err := hooks.EnqueueAdoptedPublication(context.Background(), opts)
	if err != nil {
		t.Fatalf("EnqueueAdoptedPublication() second error = %v", err)
	}
	if !first.Queued || !second.Queued {
		t.Fatalf("queued = %t/%t, want both queued", first.Queued, second.Queued)
	}
	if first.QueueID != second.QueueID {
		t.Fatalf("queue ids = %q vs %q, want identical", first.QueueID, second.QueueID)
	}
	status, err := publication.QueuedUploadStatus()
	if err != nil {
		t.Fatal(err)
	}
	if status.Pending != 1 {
		t.Fatalf("pending uploads = %d, want 1", status.Pending)
	}
}

func TestEnqueueAdoptedPublicationBuildsV2Revisions(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	repo := initPushHookRepo(t)
	base := gitRev(t, repo, "HEAD~1")
	head := gitRev(t, repo, "HEAD")

	result, err := hooks.EnqueueAdoptedPublication(context.Background(), hooks.AdoptPushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		LocalRef: "refs/heads/main",
		HeadSHA:  head,
		RefRange: base + ".." + head,
	})
	if err != nil {
		t.Fatalf("EnqueueAdoptedPublication() error = %v", err)
	}
	if !result.Queued {
		t.Fatal("queued = false, want true")
	}
	artifact := result.Artifact
	if artifact.SchemaVersion != 2 {
		t.Fatalf("schema version = %d, want 2", artifact.SchemaVersion)
	}
	if len(artifact.Revisions) != 1 {
		t.Fatalf("revisions = %#v, want the one pushed commit", artifact.Revisions)
	}
	revision := artifact.Revisions[0]
	if revision.RevisionID != "revpushaaaaa" {
		t.Fatalf("revision id = %q, want trailer id revpushaaaaa", revision.RevisionID)
	}
	if revision.CommitID != head {
		t.Fatalf("commit id = %q, want %q", revision.CommitID, head)
	}
	if !strings.Contains(revision.Description, "feature") {
		t.Fatalf("description = %q, want commit message", revision.Description)
	}
	if !strings.Contains(revision.Patch, "feature.txt") {
		t.Fatalf("patch = %q, want per-commit diff", revision.Patch)
	}
	if len(revision.Files) != 1 || revision.Files[0] != "feature.txt" {
		t.Fatalf("files = %#v, want [feature.txt]", revision.Files)
	}
	if revision.BranchName != "main" {
		t.Fatalf("branch name = %q, want main", revision.BranchName)
	}
}

func TestEnqueueAdoptedPublicationResolvesHeadLocalRef(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	repo := initPushHookRepo(t)
	head := gitRev(t, repo, "HEAD")

	result, err := hooks.EnqueueAdoptedPublication(context.Background(), hooks.AdoptPushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		LocalRef: "HEAD",
		HeadSHA:  head,
	})
	if err != nil {
		t.Fatalf("EnqueueAdoptedPublication() error = %v", err)
	}
	if !result.Queued {
		t.Fatal("queued = false, want true")
	}
	branch := result.Artifact.Push.BranchName
	if branch == nil || *branch != "main" {
		t.Fatalf("push branch = %v, want main (HEAD local ref must resolve to the real branch)", branch)
	}
}

func initPushHookRepo(t *testing.T) string {
	t.Helper()
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")
	repo := t.TempDir()
	for _, args := range [][]string{
		{"git", "init", "-b", "main"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test"},
	} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = repo
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v: %s", err, out)
		}
	}
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "-m", "initial")
	if err := os.WriteFile(filepath.Join(repo, "feature.txt"), []byte("feature\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "feature.txt")
	message := "feature\n\n" + vcs.RevisionTrailerLine("revpushaaaaa")
	runGit(t, repo, "commit", "-m", message)
	return repo
}

func gitRev(t *testing.T, repo, ref string) string {
	t.Helper()
	cmd := exec.Command("git", "-C", repo, "rev-parse", ref)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse %s: %v: %s", ref, err, out)
	}
	return strings.TrimSpace(string(out))
}

func runGit(t *testing.T, repo string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(args, " "), err, out)
	}
}
