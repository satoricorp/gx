package publication

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/satoricorp/lgtm/internal/reviewbundle"
	"github.com/satoricorp/lgtm/internal/semantic"
	"github.com/satoricorp/lgtm/internal/vcs"
)

func TestPublishBundleUploadsAndIndexes(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	indexer := semantic.NewIndexer(
		semantic.Config{Enabled: true, BatchSize: 2, MaxChunkBytes: 12000},
		fakeEmbedder{vectors: [][]float32{{0.1, 0.2}}},
		&fakeVectorStore{},
	)
	uploader := &fakeUploader{reviewURL: "http://lgtm.test/review/1"}
	publisher := NewPublisherWithIndexer(uploader, indexer)

	result, err := publisher.PublishBundle(context.Background(), bundleWithTranscriptSource())
	if err != nil {
		t.Fatalf("PublishBundle() error = %v", err)
	}
	if !result.Uploaded || result.ReviewURL != "http://lgtm.test/review/1" {
		t.Fatalf("result = %#v, want upload result", result)
	}
	if result.ArtifactPath == "" || result.Artifact.ReviewID != "review-test" {
		t.Fatalf("artifact result = path %q artifact %#v", result.ArtifactPath, result.Artifact)
	}
	data, err := os.ReadFile(result.ArtifactPath)
	if err != nil {
		t.Fatalf("ReadFile(local artifact) error = %v", err)
	}
	var artifact reviewbundle.Artifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		t.Fatalf("unmarshal local artifact: %v", err)
	}
	if !reflect.DeepEqual(artifact, result.Artifact) {
		t.Fatalf("local artifact = %#v, want canonical published artifact %#v", artifact, result.Artifact)
	}
	if len(uploader.uploaded) != 1 || uploader.uploaded[0].ReviewID != "" || uploader.uploaded[0].ReviewURL != "" {
		t.Fatalf("uploaded artifact = %#v, want pre-cloud canonical artifact", uploader.uploaded)
	}
	if uploader.uploaded[0].SchemaVersion != reviewbundle.SchemaVersion || len(uploader.uploaded[0].Revisions) == 0 || uploader.uploaded[0].Revisions[0].ReviewContext == nil {
		t.Fatalf("uploaded artifact missing review data: %#v", uploader.uploaded[0])
	}
	if !result.SemanticIndexed || result.SemanticChunks != 1 || result.SemanticIndexError != "" {
		t.Fatalf("semantic result = %#v, want indexed chunk", result)
	}
	if !uploader.called {
		t.Fatal("uploader was not called")
	}
}

func TestPublishBundleSemanticErrorDoesNotBlockUpload(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	indexer := semantic.NewIndexer(
		semantic.Config{Enabled: true, BatchSize: 2, MaxChunkBytes: 12000},
		fakeEmbedder{err: errors.New("embedding unavailable")},
		&fakeVectorStore{},
	)
	publisher := NewPublisherWithIndexer(&fakeUploader{reviewURL: "http://lgtm.test/review/1"}, indexer)

	result, err := publisher.PublishBundle(context.Background(), bundleWithTranscriptSource())
	if err != nil {
		t.Fatalf("PublishBundle() error = %v", err)
	}
	if !result.Uploaded {
		t.Fatalf("result = %#v, want upload despite semantic error", result)
	}
	if result.SemanticIndexError != "embedding unavailable" {
		t.Fatalf("SemanticIndexError = %q, want embedding unavailable", result.SemanticIndexError)
	}
}

func TestEnqueueArtifactQueuesAndDrainUploads(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())

	result, err := EnqueueArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "lgtm.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		TLVersion:     "test",
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo", Backend: "jj"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "abc123"},
	}), QueueAttestation{})
	if err != nil {
		t.Fatalf("EnqueueArtifact() error = %v", err)
	}
	if !result.Queued || result.QueueID == "" || result.ArtifactPath == "" {
		t.Fatalf("enqueue result = %#v, want queued artifact", result)
	}

	status, err := QueuedUploadStatus()
	if err != nil {
		t.Fatalf("QueuedUploadStatus() error = %v", err)
	}
	if status.Pending != 1 || status.Failed != 0 {
		t.Fatalf("status = %#v, want one pending upload", status)
	}

	uploader := &fakeUploader{}
	drain, err := DrainQueuedUploads(context.Background(), uploader, 10)
	if err != nil {
		t.Fatalf("DrainQueuedUploads() error = %v", err)
	}
	if drain.Uploaded != 1 || drain.Failed != 0 || !uploader.called {
		t.Fatalf("drain = %#v uploader.called=%t, want one uploaded", drain, uploader.called)
	}

	status, err = QueuedUploadStatus()
	if err != nil {
		t.Fatalf("QueuedUploadStatus() after drain error = %v", err)
	}
	if status.Pending != 0 || status.Failed != 0 || status.Uploading != 0 {
		t.Fatalf("status after drain = %#v, want empty outbox", status)
	}
}

func TestQueuedUploadStatusTreatsDeadUploaderAsPending(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())

	artifact := reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "lgtm.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		TLVersion:     "test",
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo", Backend: "jj"},
		Push:          reviewbundle.PushPayload{HeadCommitID: "deadbeef"},
	})
	item, err := newQueueItem(artifact, QueueAttestation{})
	if err != nil {
		t.Fatalf("newQueueItem() error = %v", err)
	}
	item.Status = outboxStatusRunning
	item.Attempts = 1
	item.LastAttemptAt = 1782258693368
	if _, err := writeQueueItem(item); err != nil {
		t.Fatalf("writeQueueItem() error = %v", err)
	}
	lockPath, err := uploadLockPath()
	if err != nil {
		t.Fatalf("uploadLockPath() error = %v", err)
	}
	if err := os.WriteFile(lockPath, []byte(`{"started_at":1782258693368,"pid":999999999}`+"\n"), 0o600); err != nil {
		t.Fatalf("write stale lock: %v", err)
	}

	status, err := QueuedUploadStatus()
	if err != nil {
		t.Fatalf("QueuedUploadStatus() error = %v", err)
	}
	if status.UploadRunning || status.Uploading != 0 || status.Pending != 1 {
		t.Fatalf("status = %#v, want dead uploader item reported pending", status)
	}
}

func TestPublishStackPreparesUploadsAndRecords(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	uploader := &fakeUploader{reviewURL: "http://lgtm.test/review/stack"}
	publisher := NewPublisherWithIndexer(uploader, nil)
	push := vcs.PushResult{
		Repo:         vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"},
		HeadCommitID: "abc123",
	}
	engine := &fakeStackPublisher{push: push}

	result, err := PublishStack(context.Background(), engine, []string{"feature/demo"}, vcs.PushOptions{Mode: vcs.PublishModeReviewOnly}, publisher)
	if err != nil {
		t.Fatalf("PublishStack() error = %v", err)
	}
	if !engine.prepareCalled || !engine.recordCalled {
		t.Fatalf("engine prepareCalled=%t recordCalled=%t, want both true", engine.prepareCalled, engine.recordCalled)
	}
	if !uploader.called || !result.Review.Uploaded || result.Review.ReviewURL != "http://lgtm.test/review/stack" {
		t.Fatalf("result review = %#v uploader.called=%t", result.Review, uploader.called)
	}
	if result.Push.HeadCommitID != "abc123" {
		t.Fatalf("push = %#v, want returned push", result.Push)
	}
}

func TestPublishStackUploadErrorReturnsPreparedPush(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	uploadErr := errors.New("upload unavailable")
	publisher := NewPublisherWithIndexer(&fakeUploader{err: uploadErr}, nil)
	push := vcs.PushResult{
		Repo:         vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"},
		HeadCommitID: "abc123",
	}
	engine := &fakeStackPublisher{push: push}

	result, err := PublishStack(context.Background(), engine, nil, vcs.PushOptions{Mode: vcs.PublishModeReviewOnly}, publisher)
	if !errors.Is(err, uploadErr) {
		t.Fatalf("PublishStack() error = %v, want %v", err, uploadErr)
	}
	if result.Push.HeadCommitID != "abc123" {
		t.Fatalf("push = %#v, want prepared push returned on upload error", result.Push)
	}
	if engine.recordCalled {
		t.Fatal("RecordPublish called after upload error")
	}
}

func TestPublishStackRecordErrorReturnsReviewResult(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	recordErr := errors.New("db unavailable")
	uploader := &fakeUploader{reviewURL: "http://lgtm.test/review/stack"}
	publisher := NewPublisherWithIndexer(uploader, nil)
	push := vcs.PushResult{
		Repo:         vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"},
		HeadCommitID: "abc123",
	}
	engine := &fakeStackPublisher{push: push, recordErr: recordErr}

	result, err := PublishStack(context.Background(), engine, nil, vcs.PushOptions{Mode: vcs.PublishModeReviewOnly}, publisher)
	if !errors.Is(err, recordErr) {
		t.Fatalf("PublishStack() error = %v, want %v", err, recordErr)
	}
	if !uploader.called || !result.Review.Uploaded {
		t.Fatalf("review = %#v uploader.called=%t, want uploaded review returned", result.Review, uploader.called)
	}
	if result.Push.HeadCommitID != "abc123" {
		t.Fatalf("push = %#v, want prepared push returned on record error", result.Push)
	}
}

func TestPublishAllStacksPreparesUploadsAndRecordsEachStack(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	uploader := &fakeUploader{reviewURL: "http://lgtm.test/review/stack"}
	publisher := NewPublisherWithIndexer(uploader, nil)
	engine := &fakeStackPublisher{allPushes: []vcs.PushResult{
		{Repo: vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"}, HeadCommitID: "abc123"},
		{Repo: vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"}, HeadCommitID: "def456"},
	}}

	result, err := PublishAllStacks(context.Background(), engine, nil, vcs.PushOptions{Mode: vcs.PublishModeReviewOnly}, publisher)
	if err != nil {
		t.Fatalf("PublishAllStacks() error = %v", err)
	}
	if !engine.prepareAllCalled {
		t.Fatal("PrepareAllPublishes was not called")
	}
	if engine.recordCount != 2 {
		t.Fatalf("recordCount = %d, want 2", engine.recordCount)
	}
	if uploader.calls != 2 {
		t.Fatalf("uploader calls = %d, want 2", uploader.calls)
	}
	if len(result.Stacks) != 2 {
		t.Fatalf("result stacks = %d, want 2", len(result.Stacks))
	}
}

func TestPublishAllStacksUploadErrorReturnsCompletedResults(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	uploadErr := errors.New("upload unavailable")
	uploader := &fakeUploader{reviewURL: "http://lgtm.test/review/stack", errAfterCalls: 2, err: uploadErr}
	publisher := NewPublisherWithIndexer(uploader, nil)
	engine := &fakeStackPublisher{allPushes: []vcs.PushResult{
		{Repo: vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"}, HeadCommitID: "abc123"},
		{Repo: vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"}, HeadCommitID: "def456"},
	}}

	result, err := PublishAllStacks(context.Background(), engine, nil, vcs.PushOptions{Mode: vcs.PublishModeReviewOnly}, publisher)
	if !errors.Is(err, uploadErr) {
		t.Fatalf("PublishAllStacks() error = %v, want %v", err, uploadErr)
	}
	if len(result.Stacks) != 2 {
		t.Fatalf("result stacks = %d, want first success plus failed push", len(result.Stacks))
	}
	if engine.recordCount != 1 {
		t.Fatalf("recordCount = %d, want 1", engine.recordCount)
	}
}

func TestPublishAllStacksRecordsPreparedStacksBeforeReturningPrepareError(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	prepareErr := errors.New("one stack failed")
	uploader := &fakeUploader{reviewURL: "http://lgtm.test/review/stack"}
	publisher := NewPublisherWithIndexer(uploader, nil)
	engine := &fakeStackPublisher{
		allPushes: []vcs.PushResult{
			{Repo: vcs.RepoInfo{RootPath: t.TempDir(), Backend: "jj"}, HeadCommitID: "abc123"},
		},
		prepareAllErr: prepareErr,
	}

	result, err := PublishAllStacks(context.Background(), engine, nil, vcs.PushOptions{Mode: vcs.PublishModeReviewOnly}, publisher)
	if !errors.Is(err, prepareErr) {
		t.Fatalf("PublishAllStacks() error = %v, want %v", err, prepareErr)
	}
	if len(result.Stacks) != 1 {
		t.Fatalf("result stacks = %d, want prepared stack result", len(result.Stacks))
	}
	if engine.recordCount != 1 {
		t.Fatalf("recordCount = %d, want 1", engine.recordCount)
	}
}

type fakeUploader struct {
	called        bool
	calls         int
	reviewURL     string
	err           error
	errAfterCalls int
	uploaded      []reviewbundle.Artifact
}

func (f *fakeUploader) UploadReviewArtifact(_ context.Context, artifact reviewbundle.Artifact) (reviewbundle.Artifact, error) {
	f.called = true
	f.calls++
	f.uploaded = append(f.uploaded, artifact)
	if f.err != nil && (f.errAfterCalls == 0 || f.calls >= f.errAfterCalls) {
		return reviewbundle.Artifact{}, f.err
	}
	artifact.ReviewID = "review-test"
	artifact.ReviewURL = f.reviewURL
	artifact.IndexStatus = "pending"
	return artifact, nil
}

type fakeStackPublisher struct {
	prepareCalled    bool
	prepareAllCalled bool
	recordCalled     bool
	recordCount      int
	push             vcs.PushResult
	allPushes        []vcs.PushResult
	err              error
	prepareAllErr    error
	recordErr        error
}

func (f *fakeStackPublisher) PreparePublish(ctx context.Context, args []string, opts vcs.PushOptions) (vcs.PushResult, error) {
	f.prepareCalled = true
	if f.err != nil {
		return vcs.PushResult{}, f.err
	}
	return f.push, nil
}

func (f *fakeStackPublisher) PrepareAllPublishes(ctx context.Context, args []string, opts vcs.PushOptions) ([]vcs.PushResult, error) {
	f.prepareAllCalled = true
	if f.err != nil {
		return nil, f.err
	}
	return f.allPushes, f.prepareAllErr
}

func (f *fakeStackPublisher) RecordPublish(ctx context.Context, result vcs.PushResult) error {
	f.recordCalled = true
	f.recordCount++
	if f.recordErr != nil {
		return f.recordErr
	}
	return nil
}

type fakeEmbedder struct {
	vectors [][]float32
	err     error
}

func (f fakeEmbedder) Embed(_ context.Context, _ []string) ([][]float32, error) {
	return f.vectors, f.err
}

type fakeVectorStore struct{}

func (f *fakeVectorStore) Upsert(_ context.Context, _ []semantic.VectorRow) error {
	return nil
}

func bundleWithTranscriptSource() reviewbundle.Bundle {
	responseID := "response-one"
	return reviewbundle.Bundle{
		Event:         "lgtm.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: "/repo"},
		Revisions: []reviewbundle.RevisionPayload{{
			RevisionID:  "change-one",
			Description: "feat alpha",
			ReviewContext: &reviewbundle.ReviewContextPayload{
				TranscriptSources: []reviewbundle.ReviewTranscriptSource{{
					SessionID:  "session-one",
					RequestID:  "request-one",
					ResponseID: &responseID,
					Provider:   "openai",
					Status:     "explicit",
				}},
			},
		}},
		Sessions: []reviewbundle.SessionPayload{{
			ID: "session-one",
			Requests: []reviewbundle.RequestPayload{{
				ID:             "request-one",
				Provider:       "openai",
				Method:         "POST",
				Endpoint:       "/v1/responses",
				RequestBody:    []byte(`{"input":"alpha"}`),
				RequestHeaders: "{}",
				Responses: []reviewbundle.ResponsePayload{{
					ID:              "response-one",
					ResponseBody:    []byte(`{"output":"done"}`),
					ResponseHeaders: "{}",
				}},
			}},
		}},
	}
}
