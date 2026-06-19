package cli

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/storage"
)

type cloudSyncSummary struct {
	Listed   int
	Updated  int
	CaughtUp int
}

func syncCloudMetadata(
	ctx context.Context,
	engine *authoring.Engine,
	repo authoring.RepoInfo,
	out io.Writer,
) (*cloudSyncSummary, error) {
	client := cloud.NewClient()
	if client == nil || !cloud.CloudConfigured() {
		return nil, nil
	}

	repoFullName := cloud.RepoFullNameFromRemoteURL(derefString(repo.RemoteURL))
	if repoFullName == "" {
		fmt.Fprintln(out, muted("Skipping gx cloud bookmark sync: could not resolve repo from remote URL."))
		return nil, nil
	}

	bookmarks, err := client.ListBookmarks(ctx, cloud.ListBookmarksOptions{
		RepoFullName: repoFullName,
	})
	if err != nil {
		return nil, err
	}

	db, err := storage.Open(ctx)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	store, err := storage.NewStore(ctx, db)
	if err != nil {
		return nil, err
	}

	now := time.Now().UnixMilli()
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repo.RootPath,
		Backend:       repo.Backend,
		DefaultRemote: repo.DefaultRemote,
		DefaultBranch: repo.DefaultBranch,
		RemoteURL:     repo.RemoteURL,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return nil, err
	}

	summary := &cloudSyncSummary{Listed: len(bookmarks)}

	for _, bookmark := range bookmarks {
		prevRevision, hasPrev, err := storage.GetCloudBookmarkRevision(ctx, db, bookmark.ID)
		if err != nil {
			return nil, err
		}

		state := storage.CloudBookmarkState{
			PostgresBookmarkID: bookmark.ID,
			RepoID:             repoID,
			RepoFullName:       bookmark.RepoFullName,
			BranchName:         bookmark.BranchName,
			Revision:           bookmark.Revision,
			HeadCommitID:       nullString(bookmark.HeadCommitID),
			RemoteHeadSha:      nullString(bookmark.RemoteHeadSha),
			MergeStatus:        bookmark.MergeStatus,
			UpdatedAtMs:        bookmark.UpdatedAtMs,
			SyncedAt:           now,
		}
		if err := storage.UpsertCloudBookmark(ctx, db, state); err != nil {
			return nil, err
		}

		needsCatchUp := !hasPrev || bookmark.Revision > prevRevision
		if !needsCatchUp {
			continue
		}
		summary.Updated++

		if bookmark.RemoteHeadSha == nil || strings.TrimSpace(*bookmark.RemoteHeadSha) == "" {
			continue
		}
		if err := engine.SyncCloudBookmarkTip(ctx, repo, bookmark.BranchName); err != nil {
			fmt.Fprintf(out, "%s %s\n", danger("Cloud catch-up failed"), muted(fmt.Sprintf("%s: %v", bookmark.BranchName, err)))
			continue
		}
		summary.CaughtUp++
	}

	return summary, nil
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func nullString(value *string) sql.NullString {
	if value == nil || strings.TrimSpace(*value) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}
