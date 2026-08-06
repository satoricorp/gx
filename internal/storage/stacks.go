package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) UpsertStack(ctx context.Context, stack Stack) (int64, error) {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO stacks (
			repo_id, name, bookmark_name, base_ref, base_commit_id, head_change_id, head_commit_id,
			remote_name, remote_ref, github_pr_url, status, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(repo_id, bookmark_name) DO UPDATE SET
			name = excluded.name,
			base_ref = excluded.base_ref,
			base_commit_id = excluded.base_commit_id,
			head_change_id = excluded.head_change_id,
			head_commit_id = excluded.head_commit_id,
			remote_name = excluded.remote_name,
			remote_ref = excluded.remote_ref,
			github_pr_url = excluded.github_pr_url,
			status = excluded.status,
			updated_at = excluded.updated_at
	`,
		stack.RepoID,
		stack.Name,
		stack.BookmarkName,
		stack.BaseRef,
		stack.BaseCommitID,
		stack.HeadChangeID,
		stack.HeadCommitID,
		stack.RemoteName,
		stack.RemoteRef,
		stack.GitHubPRURL,
		stack.Status,
		stack.CreatedAt,
		stack.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("upsert stack: %w", err)
	}
	var id int64
	if err := s.db.QueryRowContext(ctx, `SELECT id FROM stacks WHERE repo_id = ? AND bookmark_name = ?`, stack.RepoID, stack.BookmarkName).Scan(&id); err != nil {
		return 0, fmt.Errorf("lookup stack id: %w", err)
	}
	return id, nil
}

func (s *Store) RenameStack(ctx context.Context, repoID int64, bookmarkName, normalizedName string, updatedAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE stacks
		SET name = ?, updated_at = ?
		WHERE repo_id = ? AND bookmark_name = ?
	`, normalizedName, updatedAt, repoID, bookmarkName)
	if err != nil {
		return fmt.Errorf("rename stack: %w", err)
	}
	return nil
}

func (s *Store) RenameStackBaseRef(ctx context.Context, repoID int64, oldRef, newRef string, updatedAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE stacks
		SET base_ref = ?, updated_at = ?
		WHERE repo_id = ? AND base_ref = ?
	`, newRef, updatedAt, repoID, oldRef)
	if err != nil {
		return fmt.Errorf("rename stack base ref: %w", err)
	}
	return nil
}

func (s *Store) PrunePublishedStack(ctx context.Context, repoID, stackID int64, publishRef string, updatedAt int64) error {
	publishRef = strings.TrimSpace(publishRef)
	if repoID == 0 || stackID == 0 || publishRef == "" {
		return nil
	}
	gx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin prune published stack: %w", err)
	}
	defer gx.Rollback()

	if _, err := gx.ExecContext(ctx, `
		DELETE FROM change_bookmarks
		WHERE bookmark_name = ?
			AND change_id IN (
				SELECT change_id FROM stack_changes WHERE stack_id = ?
			)
	`, publishRef, stackID); err != nil {
		return fmt.Errorf("delete stack change bookmarks: %w", err)
	}
	if _, err := gx.ExecContext(ctx, `
		DELETE FROM pushes
		WHERE repo_id = ? AND branch_name = ?
	`, repoID, publishRef); err != nil {
		return fmt.Errorf("delete stack pushes: %w", err)
	}
	if _, err := gx.ExecContext(ctx, `
		DELETE FROM stack_changes
		WHERE stack_id = ?
	`, stackID); err != nil {
		return fmt.Errorf("delete pruned stack changes: %w", err)
	}
	if _, err := gx.ExecContext(ctx, `
		DELETE FROM stacks
		WHERE repo_id = ? AND id = ?
	`, repoID, stackID); err != nil {
		return fmt.Errorf("delete pruned stack: %w", err)
	}
	return gx.Commit()
}

func (s *Store) MarkStackStatus(ctx context.Context, stackID int64, status string, updatedAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE stacks
		SET status = ?, updated_at = ?
		WHERE id = ?
	`, status, updatedAt, stackID)
	if err != nil {
		return fmt.Errorf("mark stack status: %w", err)
	}
	return nil
}

func (s *Store) AddChangeToStack(ctx context.Context, stackID, changeID int64, createdAt int64) error {
	var next int
	if err := s.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(position), 0) + 1 FROM stack_changes WHERE stack_id = ?`, stackID).Scan(&next); err != nil {
		return fmt.Errorf("next stack change position: %w", err)
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO stack_changes (stack_id, change_id, jj_change_id, commit_id, position, created_at)
		SELECT ?, c.id, c.jj_change_id, c.current_commit_id, ?, ?
		FROM changes c
		WHERE c.id = ?
		ON CONFLICT(stack_id, change_id) DO UPDATE SET
			jj_change_id = excluded.jj_change_id,
			commit_id = excluded.commit_id
	`, stackID, next, createdAt, changeID)
	if err != nil {
		return fmt.Errorf("add change to stack: %w", err)
	}
	return nil
}

func (s *Store) FindStackByBookmark(ctx context.Context, repoID int64, bookmarkName string) (*Stack, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, name, bookmark_name, base_ref, base_commit_id, head_change_id, head_commit_id,
			remote_name, remote_ref, github_pr_url, status, created_at, updated_at
		FROM stacks
		WHERE repo_id = ? AND bookmark_name = ?
	`, repoID, bookmarkName)
	return scanStack(row)
}

func (s *Store) ListStacksByRepoID(ctx context.Context, repoID int64) ([]Stack, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_id, name, bookmark_name, base_ref, base_commit_id, head_change_id, head_commit_id,
			remote_name, remote_ref, github_pr_url, status, created_at, updated_at
		FROM stacks
		WHERE repo_id = ?
		ORDER BY updated_at DESC, id DESC
	`, repoID)
	if err != nil {
		return nil, fmt.Errorf("list stacks by repo: %w", err)
	}
	defer rows.Close()
	stacks := []Stack{}
	for rows.Next() {
		stack, err := scanStackRows(rows)
		if err != nil {
			return nil, err
		}
		stacks = append(stacks, stack)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stacks: %w", err)
	}
	return stacks, nil
}

func (s *Store) ListChangesByStackID(ctx context.Context, stackID int64) ([]Change, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.repo_id, c.jj_change_id, c.current_commit_id, c.description, c.parent_change_id, c.status, c.first_seen_at, c.updated_at
		FROM changes c
		INNER JOIN stack_changes bc ON bc.change_id = c.id
		WHERE bc.stack_id = ?
		ORDER BY bc.position ASC, bc.id ASC
	`, stackID)
	if err != nil {
		return nil, fmt.Errorf("list changes by stack: %w", err)
	}
	defer rows.Close()
	return scanChanges(rows)
}

func (s *Store) ListChangesByStackIDs(ctx context.Context, stackIDs []int64) (map[int64][]Change, error) {
	ids := uniquePositiveInt64s(stackIDs)
	changesByStack := make(map[int64][]Change, len(ids))
	for _, id := range ids {
		changesByStack[id] = nil
	}
	if len(ids) == 0 {
		return changesByStack, nil
	}
	placeholders, args := placeholdersForInt64s(ids)
	rows, err := s.db.QueryContext(ctx, `
		SELECT bc.stack_id, c.id, c.repo_id, c.jj_change_id, c.current_commit_id, c.description, c.parent_change_id, c.status, c.first_seen_at, c.updated_at
		FROM stack_changes bc
		INNER JOIN changes c ON bc.change_id = c.id
		WHERE bc.stack_id IN (`+placeholders+`)
		ORDER BY bc.stack_id ASC, bc.position ASC, bc.id ASC
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("list changes by stacks: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var stackID int64
		var change Change
		var parentChangeID sql.NullString
		if err := rows.Scan(
			&stackID,
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
			return nil, fmt.Errorf("scan stacked change: %w", err)
		}
		if parentChangeID.Valid {
			change.ParentChangeID = &parentChangeID.String
		}
		changesByStack[stackID] = append(changesByStack[stackID], change)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stacked changes: %w", err)
	}
	return changesByStack, nil
}

type stackScanner interface {
	Scan(dest ...any) error
}

func scanStack(row stackScanner) (*Stack, error) {
	stack, err := scanStackValue(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &stack, nil
}

func scanStackRows(row stackScanner) (Stack, error) {
	stack, err := scanStackValue(row)
	if err != nil {
		return Stack{}, err
	}
	return stack, nil
}

func scanStackValue(row stackScanner) (Stack, error) {
	var stack Stack
	var headChangeID sql.NullString
	var headCommitID sql.NullString
	var remoteName sql.NullString
	var remoteRef sql.NullString
	var githubPRURL sql.NullString
	if err := row.Scan(
		&stack.ID,
		&stack.RepoID,
		&stack.Name,
		&stack.BookmarkName,
		&stack.BaseRef,
		&stack.BaseCommitID,
		&headChangeID,
		&headCommitID,
		&remoteName,
		&remoteRef,
		&githubPRURL,
		&stack.Status,
		&stack.CreatedAt,
		&stack.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return Stack{}, err
		}
		return Stack{}, fmt.Errorf("scan stack: %w", err)
	}
	if headChangeID.Valid {
		stack.HeadChangeID = &headChangeID.String
	}
	if headCommitID.Valid {
		stack.HeadCommitID = &headCommitID.String
	}
	if remoteName.Valid {
		stack.RemoteName = &remoteName.String
	}
	if remoteRef.Valid {
		stack.RemoteRef = &remoteRef.String
	}
	if githubPRURL.Valid {
		stack.GitHubPRURL = &githubPRURL.String
	}
	return stack, nil
}
