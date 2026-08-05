package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/satoricorp/lgtm/internal/agentprovenance"
)

func (s *Store) UpsertChange(ctx context.Context, change Change) (int64, error) {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO changes (
			repo_id, jj_change_id, current_commit_id, description, parent_change_id, status, first_seen_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(repo_id, jj_change_id) DO UPDATE SET
			current_commit_id = excluded.current_commit_id,
			description = excluded.description,
			parent_change_id = excluded.parent_change_id,
			status = excluded.status,
			updated_at = excluded.updated_at
	`,
		change.RepoID,
		change.JJChangeID,
		change.CurrentCommitID,
		change.Description,
		change.ParentChangeID,
		change.Status,
		change.FirstSeenAt,
		change.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("upsert change: %w", err)
	}
	var id int64
	if err := s.db.QueryRowContext(ctx, `SELECT id FROM changes WHERE repo_id = ? AND jj_change_id = ?`, change.RepoID, change.JJChangeID).Scan(&id); err != nil {
		return 0, fmt.Errorf("lookup change id: %w", err)
	}
	return id, nil
}

func (s *Store) FindChangeByJJChangeID(ctx context.Context, repoID int64, jjChangeID string) (*Change, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, jj_change_id, current_commit_id, description, parent_change_id, status, first_seen_at, updated_at
		FROM changes
		WHERE repo_id = ? AND jj_change_id = ?
	`, repoID, jjChangeID)
	var change Change
	var parentChangeID sql.NullString
	if err := row.Scan(
		&change.ID,
		&change.RepoID,
		&change.JJChangeID,
		&change.CurrentCommitID,
		&change.Description,
		&parentChangeID,
		&change.Status,
		&change.FirstSeenAt,
		&change.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find change by jj change id: %w", err)
	}
	if parentChangeID.Valid {
		change.ParentChangeID = &parentChangeID.String
	}
	return &change, nil
}

func (s *Store) WriteChangeRevision(ctx context.Context, rev ChangeRevision) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO change_revisions (change_id, jj_commit_id, jj_operation_id, changed_files_json, created_at)
		VALUES (?, ?, ?, ?, ?)
	`,
		rev.ChangeID,
		rev.JJCommitID,
		rev.JJOperationID,
		rev.ChangedFiles,
		rev.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert change revision: %w", err)
	}
	return nil
}

func (s *Store) WriteChangeSessions(ctx context.Context, changeID int64, sessionIDs []string, createdAt int64) error {
	for _, sessionID := range sessionIDs {
		if sessionID == "" {
			continue
		}
		if _, err := s.db.ExecContext(ctx, `
			INSERT OR IGNORE INTO change_sessions (change_id, session_id, created_at)
			VALUES (?, ?, ?)
		`, changeID, sessionID, createdAt); err != nil {
			return fmt.Errorf("insert change session: %w", err)
		}
		if err := s.WriteChangeSessionProvenance(ctx, changeID, sessionID, createdAt); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) WriteChangeSessionProvenance(ctx context.Context, changeID int64, sessionID string, createdAt int64) error {
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.id, s.command, s.source, s.process_name, COALESCE(r.provider, ''), COALESCE(r.model, '')
		FROM sessions s
		LEFT JOIN requests r ON r.session_id = s.id
		WHERE s.id = ?
		ORDER BY r.created_at ASC, r.id ASC
	`, sessionID)
	if err != nil {
		return fmt.Errorf("list change session provenance sources: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		source, err := scanAgentProvenanceSource(rows, createdAt)
		if err != nil {
			return err
		}
		record := agentprovenance.Resolve(source)
		if err := s.insertChangeSessionProvenance(ctx, changeID, record); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate change session provenance sources: %w", err)
	}
	return nil
}

func scanAgentProvenanceSource(rows *sql.Rows, createdAt int64) (agentprovenance.Source, error) {
	var source agentprovenance.Source
	var sourceValue, processName sql.NullString
	if err := rows.Scan(&source.SessionID, &source.Command, &sourceValue, &processName, &source.Provider, &source.ModelID); err != nil {
		return agentprovenance.Source{}, fmt.Errorf("scan change session provenance source: %w", err)
	}
	if sourceValue.Valid {
		source.Source = &sourceValue.String
	}
	if processName.Valid {
		source.ProcessName = &processName.String
	}
	source.CreatedAt = createdAt
	return source, nil
}

func (s *Store) insertChangeSessionProvenance(ctx context.Context, changeID int64, record agentprovenance.Record) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO change_session_provenance (
			change_id, session_id, agent_tool, provider, model_id, source, process_name, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, changeID, record.SessionID, record.AgentTool, record.Provider, record.ModelID, record.Source, record.ProcessName, record.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert change session provenance: %w", err)
	}
	return nil
}

func (s *Store) UpsertChangeBookmark(ctx context.Context, bookmark ChangeBookmark) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO change_bookmarks (
			change_id, bookmark_name, remote_name, remote_ref, last_pushed_commit_id, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(change_id, bookmark_name) DO UPDATE SET
			remote_name = excluded.remote_name,
			remote_ref = excluded.remote_ref,
			last_pushed_commit_id = excluded.last_pushed_commit_id,
			updated_at = excluded.updated_at
	`,
		bookmark.ChangeID,
		bookmark.BookmarkName,
		bookmark.RemoteName,
		bookmark.RemoteRef,
		bookmark.LastPushedCommitID,
		bookmark.CreatedAt,
		bookmark.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert change bookmark: %w", err)
	}
	return nil
}

func (s *Store) ListChangeBookmarksByName(ctx context.Context, bookmarkName string) ([]ChangeBookmark, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, change_id, bookmark_name, remote_name, remote_ref, last_pushed_commit_id, created_at, updated_at
		FROM change_bookmarks
		WHERE bookmark_name = ?
	`, bookmarkName)
	if err != nil {
		return nil, fmt.Errorf("list change bookmarks: %w", err)
	}
	defer rows.Close()
	bookmarks := []ChangeBookmark{}
	for rows.Next() {
		var bookmark ChangeBookmark
		var remoteName sql.NullString
		var remoteRef sql.NullString
		if err := rows.Scan(
			&bookmark.ID,
			&bookmark.ChangeID,
			&bookmark.BookmarkName,
			&remoteName,
			&remoteRef,
			&bookmark.LastPushedCommitID,
			&bookmark.CreatedAt,
			&bookmark.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan change bookmark: %w", err)
		}
		if remoteName.Valid {
			bookmark.RemoteName = &remoteName.String
		}
		if remoteRef.Valid {
			bookmark.RemoteRef = &remoteRef.String
		}
		bookmarks = append(bookmarks, bookmark)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change bookmarks: %w", err)
	}
	return bookmarks, nil
}

func (s *Store) ListChangeBookmarksByNames(ctx context.Context, bookmarkNames []string) (map[string][]ChangeBookmark, error) {
	names := uniqueNonEmptyStrings(bookmarkNames)
	bookmarksByName := make(map[string][]ChangeBookmark, len(names))
	for _, name := range names {
		bookmarksByName[name] = nil
	}
	if len(names) == 0 {
		return bookmarksByName, nil
	}
	placeholders, args := placeholdersForStrings(names)
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, change_id, bookmark_name, remote_name, remote_ref, last_pushed_commit_id, created_at, updated_at
		FROM change_bookmarks
		WHERE bookmark_name IN (`+placeholders+`)
		ORDER BY bookmark_name ASC, id ASC
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("list change bookmarks by names: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var bookmark ChangeBookmark
		var remoteName sql.NullString
		var remoteRef sql.NullString
		if err := rows.Scan(
			&bookmark.ID,
			&bookmark.ChangeID,
			&bookmark.BookmarkName,
			&remoteName,
			&remoteRef,
			&bookmark.LastPushedCommitID,
			&bookmark.CreatedAt,
			&bookmark.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan change bookmark: %w", err)
		}
		if remoteName.Valid {
			bookmark.RemoteName = &remoteName.String
		}
		if remoteRef.Valid {
			bookmark.RemoteRef = &remoteRef.String
		}
		bookmarksByName[bookmark.BookmarkName] = append(bookmarksByName[bookmark.BookmarkName], bookmark)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change bookmarks: %w", err)
	}
	return bookmarksByName, nil
}

func (s *Store) FindChangeByCommitID(ctx context.Context, repoID int64, commitID string) (*Change, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, jj_change_id, current_commit_id, description, parent_change_id, status, first_seen_at, updated_at
		FROM changes
		WHERE repo_id = ? AND current_commit_id = ?
		ORDER BY updated_at DESC
		LIMIT 1
	`, repoID, commitID)
	change, err := scanChangeRow(row)
	if err != nil || change == nil {
		return change, err
	}
	return change, nil
}

func (s *Store) UpdateChangeCommitID(ctx context.Context, changeID int64, commitID string, updatedAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE changes
		SET current_commit_id = ?, updated_at = ?
		WHERE id = ?
	`, commitID, updatedAt, changeID)
	if err != nil {
		return fmt.Errorf("update change commit id: %w", err)
	}
	return nil
}

func scanChangeRow(row *sql.Row) (*Change, error) {
	var change Change
	var parent sql.NullString
	if err := row.Scan(
		&change.ID,
		&change.RepoID,
		&change.JJChangeID,
		&change.CurrentCommitID,
		&change.Description,
		&parent,
		&change.Status,
		&change.FirstSeenAt,
		&change.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan change row: %w", err)
	}
	if parent.Valid {
		change.ParentChangeID = &parent.String
	}
	return &change, nil
}

func (s *Store) ListChangesByRepoID(ctx context.Context, repoID int64) ([]Change, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_id, jj_change_id, current_commit_id, description, parent_change_id, status, first_seen_at, updated_at
		FROM changes
		WHERE repo_id = ?
		ORDER BY first_seen_at ASC, id ASC
	`, repoID)
	if err != nil {
		return nil, fmt.Errorf("list changes by repo: %w", err)
	}
	defer rows.Close()
	return scanChanges(rows)
}

func scanChanges(rows *sql.Rows) ([]Change, error) {
	changes := []Change{}
	for rows.Next() {
		var change Change
		var parentChangeID sql.NullString
		if err := rows.Scan(
			&change.ID,
			&change.RepoID,
			&change.JJChangeID,
			&change.CurrentCommitID,
			&change.Description,
			&parentChangeID,
			&change.Status,
			&change.FirstSeenAt,
			&change.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan change: %w", err)
		}
		if parentChangeID.Valid {
			change.ParentChangeID = &parentChangeID.String
		}
		changes = append(changes, change)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate changes: %w", err)
	}
	return changes, nil
}
