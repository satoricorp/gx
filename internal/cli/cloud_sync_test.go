package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

type cloudSyncNoopRunner struct{}

func (cloudSyncNoopRunner) Run(context.Context, string, string, ...string) (string, error) {
	return "", fmt.Errorf("unexpected runner call")
}

func (cloudSyncNoopRunner) RunStdout(context.Context, string, string, ...string) (string, error) {
	return "", fmt.Errorf("unexpected runner stdout call")
}

func (cloudSyncNoopRunner) RunStream(context.Context, string, string, ...string) error {
	return fmt.Errorf("unexpected runner stream call")
}

func TestSyncCloudMetadataDeletesMergedBookmarksAndPrunesPublishedStack(t *testing.T) {
	ctx := context.Background()
	repoRoot := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_sync",
		CLISessionToken:   "gxcs_sync",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	deleted := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer gxcs_sync" {
			t.Fatalf("authorization = %q", auth)
		}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/bookmarks":
			if got := r.URL.Query().Get("repo_full_name"); got != "satoricorp/gx" {
				t.Fatalf("repo_full_name = %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[{"id":"bookmark-1","repo_full_name":"satoricorp/gx","branch_name":"feature/merged","revision":2,"merge_status":"merged","updated_at_ms":20}]`)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/bookmarks/bookmark-1":
			deleted = true
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected cloud request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()
	t.Setenv("GX_CLOUD_URL", server.URL)

	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		RemoteURL: stringPtr("git@github.com:satoricorp/gx.git"),
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "merged",
		BookmarkName: "feature/merged",
		BaseRef:      "main",
		BaseCommitID: "base",
		RemoteRef:    stringPtr("refs/heads/feature/merged"),
		Status:       "published",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "merged-change",
		CurrentCommitID: "merged-commit",
		Description:     "merged change",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}
	if err := store.UpsertChangeBookmark(ctx, storage.ChangeBookmark{
		ChangeID:           changeID,
		BookmarkName:       "feature/merged",
		RemoteName:         stringPtr("origin"),
		RemoteRef:          stringPtr("refs/heads/feature/merged"),
		LastPushedCommitID: "merged-commit",
		CreatedAt:          1,
		UpdatedAt:          1,
	}); err != nil {
		t.Fatalf("UpsertChangeBookmark() error = %v", err)
	}
	if err := store.WritePush(ctx, storage.Push{
		RepoID:          repoID,
		RemoteName:      stringPtr("origin"),
		BranchName:      stringPtr("feature/merged"),
		HeadCommitID:    "merged-commit",
		CurrentChangeID: int64Ptr(changeID),
		CreatedAt:       1,
	}); err != nil {
		t.Fatalf("WritePush() error = %v", err)
	}
	if err := store.UpsertCloudBookmark(ctx, storage.CloudBookmarkState{
		PostgresBookmarkID: "bookmark-1",
		RepoID:             repoID,
		RepoFullName:       "satoricorp/gx",
		BranchName:         "feature/merged",
		Revision:           1,
		MergeStatus:        "open",
		UpdatedAtMs:        1,
		SyncedAt:           1,
	}); err != nil {
		t.Fatalf("UpsertCloudBookmark(open) error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	engine := authoring.NewEngineWithVCS(vcs.NewServiceWithRunner(cloudSyncNoopRunner{}))
	summary, err := syncCloudMetadata(ctx, engine, authoring.RepoInfo{
		RootPath:  repoRoot,
		Backend:   "jj",
		RemoteURL: stringPtr("git@github.com:satoricorp/gx.git"),
	}, io.Discard)
	if err != nil {
		t.Fatalf("syncCloudMetadata() error = %v", err)
	}
	if summary == nil || summary.Listed != 1 || summary.Removed != 1 {
		t.Fatalf("summary = %#v, want one listed and removed", summary)
	}
	if !deleted {
		t.Fatal("merged cloud bookmark was not deleted")
	}

	db, err = storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open(after sync) error = %v", err)
	}
	store, err = storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore(after sync) error = %v", err)
	}
	defer store.Close()
	stack, err := store.FindStackByBookmark(ctx, repoID, "feature/merged")
	if err != nil {
		t.Fatalf("FindStackByBookmark() error = %v", err)
	}
	if stack != nil {
		t.Fatalf("stack = %#v, want removed from local db", stack)
	}
	bookmarks, err := store.ListChangeBookmarksByName(ctx, "feature/merged")
	if err != nil {
		t.Fatalf("ListChangeBookmarksByName() error = %v", err)
	}
	if len(bookmarks) != 0 {
		t.Fatalf("change bookmarks = %#v, want none", bookmarks)
	}
	push, err := store.LatestPushByBranchName(ctx, repoID, "feature/merged")
	if err != nil {
		t.Fatalf("LatestPushByBranchName() error = %v", err)
	}
	if push != nil {
		t.Fatalf("push = %#v, want nil", push)
	}
	openCloudBookmarks, err := store.CountOpenCloudBookmarksByRepoID(ctx, repoID)
	if err != nil {
		t.Fatalf("CountOpenCloudBookmarksByRepoID() error = %v", err)
	}
	if openCloudBookmarks != 0 {
		t.Fatalf("open cloud bookmarks = %d, want 0", openCloudBookmarks)
	}
}

func stringPtr(value string) *string {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}
