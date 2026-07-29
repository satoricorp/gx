package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) LatestPushByRepoID(ctx context.Context, repoID int64) (*Push, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, remote_name, branch_name, head_commit_id, current_change_id, created_at
		FROM pushes
		WHERE repo_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, repoID)

	var push Push
	var remoteName sql.NullString
	var branchName sql.NullString
	var currentChangeID sql.NullInt64
	if err := row.Scan(
		&push.ID,
		&push.RepoID,
		&remoteName,
		&branchName,
		&push.HeadCommitID,
		&currentChangeID,
		&push.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("latest push by repo: %w", err)
	}
	if remoteName.Valid {
		push.RemoteName = &remoteName.String
	}
	if branchName.Valid {
		push.BranchName = &branchName.String
	}
	if currentChangeID.Valid {
		push.CurrentChangeID = &currentChangeID.Int64
	}
	return &push, nil
}

func (s *Store) LatestPushesByBranchNames(ctx context.Context, repoID int64, branchNames []string) (map[string]Push, error) {
	names := uniqueNonEmptyStrings(branchNames)
	pushesByBranch := make(map[string]Push, len(names))
	if len(names) == 0 {
		return pushesByBranch, nil
	}
	placeholders, branchArgs := placeholdersForStrings(names)
	args := append([]any{repoID}, branchArgs...)
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_id, remote_name, branch_name, head_commit_id, current_change_id, created_at
		FROM pushes
		WHERE repo_id = ? AND branch_name IN (`+placeholders+`)
		ORDER BY branch_name ASC, created_at DESC, id DESC
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("latest pushes by branch names: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var push Push
		var remoteName sql.NullString
		var branchName sql.NullString
		var currentChangeID sql.NullInt64
		if err := rows.Scan(
			&push.ID,
			&push.RepoID,
			&remoteName,
			&branchName,
			&push.HeadCommitID,
			&currentChangeID,
			&push.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan push by branch: %w", err)
		}
		if remoteName.Valid {
			push.RemoteName = &remoteName.String
		}
		if branchName.Valid {
			push.BranchName = &branchName.String
		}
		if currentChangeID.Valid {
			push.CurrentChangeID = &currentChangeID.Int64
		}
		branch := ""
		if push.BranchName != nil {
			branch = strings.TrimSpace(*push.BranchName)
		}
		if branch == "" {
			continue
		}
		if _, exists := pushesByBranch[branch]; !exists {
			pushesByBranch[branch] = push
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pushes by branch: %w", err)
	}
	return pushesByBranch, nil
}
