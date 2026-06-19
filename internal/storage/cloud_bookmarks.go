package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type CloudBookmarkState struct {
	PostgresBookmarkID string
	RepoID             int64
	RepoFullName       string
	BranchName         string
	Revision           int
	HeadCommitID       sql.NullString
	RemoteHeadSha      sql.NullString
	MergeStatus        string
	UpdatedAtMs        int64
	SyncedAt           int64
}

func ensureCloudBookmarksTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS cloud_bookmarks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postgres_bookmark_id TEXT NOT NULL UNIQUE,
			repo_id INTEGER NOT NULL REFERENCES repos(id),
			repo_full_name TEXT NOT NULL,
			branch_name TEXT NOT NULL,
			revision INTEGER NOT NULL,
			head_commit_id TEXT,
			remote_head_sha TEXT,
			merge_status TEXT NOT NULL,
			updated_at_ms INTEGER NOT NULL,
			synced_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_cloud_bookmarks_repo ON cloud_bookmarks(repo_id, branch_name);
	`)
	return err
}

func UpsertCloudBookmark(ctx context.Context, db *sql.DB, state CloudBookmarkState) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO cloud_bookmarks (
			postgres_bookmark_id,
			repo_id,
			repo_full_name,
			branch_name,
			revision,
			head_commit_id,
			remote_head_sha,
			merge_status,
			updated_at_ms,
			synced_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(postgres_bookmark_id) DO UPDATE SET
			revision = excluded.revision,
			head_commit_id = excluded.head_commit_id,
			remote_head_sha = excluded.remote_head_sha,
			merge_status = excluded.merge_status,
			updated_at_ms = excluded.updated_at_ms,
			synced_at = excluded.synced_at
	`, state.PostgresBookmarkID, state.RepoID, state.RepoFullName, state.BranchName,
		state.Revision, nullableStringArg(state.HeadCommitID), nullableStringArg(state.RemoteHeadSha),
		state.MergeStatus, state.UpdatedAtMs, state.SyncedAt)
	if err != nil {
		return fmt.Errorf("upsert cloud bookmark: %w", err)
	}
	return nil
}

func GetCloudBookmarkRevision(ctx context.Context, db *sql.DB, postgresBookmarkID string) (int, bool, error) {
	var revision int
	err := db.QueryRowContext(ctx, `
		SELECT revision FROM cloud_bookmarks WHERE postgres_bookmark_id = ?
	`, postgresBookmarkID).Scan(&revision)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return revision, true, nil
}

func nullableStringArg(value sql.NullString) any {
	if value.Valid {
		return value.String
	}
	return nil
}
