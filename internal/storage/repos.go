package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) UpsertRepo(ctx context.Context, repo Repo) (int64, error) {
	gitCommonDir := strings.TrimSpace(repo.GitCommonDir)
	if gitCommonDir == "" {
		gitCommonDir = strings.TrimSpace(repo.RootPath)
	}
	var existingID int64
	err := s.db.QueryRowContext(ctx, `
		SELECT id
		FROM repos
		WHERE git_common_dir = ? OR root_path = ?
		ORDER BY CASE WHEN git_common_dir = ? THEN 0 ELSE 1 END, updated_at DESC
		LIMIT 1
	`, gitCommonDir, repo.RootPath, gitCommonDir).Scan(&existingID)
	if err == nil {
		_, err = s.db.ExecContext(ctx, `
			UPDATE repos
			SET git_common_dir = ?,
				backend = CASE WHEN backend = 'jj' THEN backend ELSE ? END,
				default_remote = ?,
				default_branch = ?,
				authoring_base_ref = COALESCE(?, authoring_base_ref),
				remote_url = ?,
				updated_at = ?
			WHERE id = ?
		`, gitCommonDir, repo.Backend, repo.DefaultRemote, repo.DefaultBranch, repo.AuthoringBase, repo.RemoteURL, repo.UpdatedAt, existingID)
		if err != nil {
			return 0, fmt.Errorf("update repo: %w", err)
		}
		return existingID, nil
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("lookup repo identity: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO repos (root_path, git_common_dir, backend, default_remote, default_branch, authoring_base_ref, remote_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(root_path) DO UPDATE SET
			git_common_dir = COALESCE(NULLIF(excluded.git_common_dir, ''), repos.git_common_dir),
			backend = CASE WHEN repos.backend = 'jj' THEN repos.backend ELSE excluded.backend END,
			default_remote = excluded.default_remote,
			default_branch = excluded.default_branch,
			authoring_base_ref = COALESCE(excluded.authoring_base_ref, repos.authoring_base_ref),
			remote_url = excluded.remote_url,
			updated_at = excluded.updated_at
	`,
		repo.RootPath,
		gitCommonDir,
		repo.Backend,
		repo.DefaultRemote,
		repo.DefaultBranch,
		repo.AuthoringBase,
		repo.RemoteURL,
		repo.CreatedAt,
		repo.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("upsert repo: %w", err)
	}
	var id int64
	if gitCommonDir != "" {
		if err := s.db.QueryRowContext(ctx, `SELECT id FROM repos WHERE git_common_dir = ? ORDER BY updated_at DESC LIMIT 1`, gitCommonDir).Scan(&id); err == nil {
			return id, nil
		}
	}
	if err := s.db.QueryRowContext(ctx, `SELECT id FROM repos WHERE root_path = ?`, repo.RootPath).Scan(&id); err != nil {
		return 0, fmt.Errorf("lookup repo id: %w", err)
	}
	return id, nil
}

func (s *Store) SetRepoAuthoringBase(ctx context.Context, repoID int64, baseRef string, updatedAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE repos
		SET authoring_base_ref = ?, updated_at = ?
		WHERE id = ?
	`, baseRef, updatedAt, repoID)
	if err != nil {
		return fmt.Errorf("set repo authoring base: %w", err)
	}
	return nil
}

func (s *Store) RecordInitializedRepo(ctx context.Context, rootPath string, now int64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO initialized_repos (root_path, created_at, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(root_path) DO UPDATE SET
			updated_at = excluded.updated_at
	`, rootPath, now, now)
	if err != nil {
		return fmt.Errorf("record initialized repo: %w", err)
	}
	return nil
}

func (s *Store) IsInitializedRepo(ctx context.Context, rootPath string) (bool, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(1)
		FROM initialized_repos
		WHERE root_path = ?
	`, rootPath).Scan(&count); err != nil {
		return false, fmt.Errorf("lookup initialized repo: %w", err)
	}
	return count > 0, nil
}

func (s *Store) ListInitializedRepos(ctx context.Context) ([]InitializedRepo, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT root_path, created_at, updated_at
		FROM initialized_repos
		ORDER BY updated_at DESC, root_path ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list initialized repos: %w", err)
	}
	defer rows.Close()

	var repos []InitializedRepo
	for rows.Next() {
		var repo InitializedRepo
		if err := rows.Scan(&repo.RootPath, &repo.CreatedAt, &repo.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan initialized repo: %w", err)
		}
		repos = append(repos, repo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list initialized repos rows: %w", err)
	}
	return repos, nil
}

func (s *Store) FindRepoByRoot(ctx context.Context, rootPath string) (*Repo, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, root_path, git_common_dir, backend, default_remote, default_branch, authoring_base_ref, remote_url, created_at, updated_at
		FROM repos
		WHERE root_path = ?
	`, rootPath)

	var repo Repo
	var gitCommonDir sql.NullString
	var defaultRemote sql.NullString
	var defaultBranch sql.NullString
	var authoringBase sql.NullString
	var remoteURL sql.NullString
	if err := row.Scan(
		&repo.ID,
		&repo.RootPath,
		&gitCommonDir,
		&repo.Backend,
		&defaultRemote,
		&defaultBranch,
		&authoringBase,
		&remoteURL,
		&repo.CreatedAt,
		&repo.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find repo by root: %w", err)
	}
	if gitCommonDir.Valid {
		repo.GitCommonDir = gitCommonDir.String
	}
	if defaultRemote.Valid {
		repo.DefaultRemote = &defaultRemote.String
	}
	if defaultBranch.Valid {
		repo.DefaultBranch = &defaultBranch.String
	}
	if authoringBase.Valid {
		repo.AuthoringBase = &authoringBase.String
	}
	if remoteURL.Valid {
		repo.RemoteURL = &remoteURL.String
	}
	return &repo, nil
}

func (s *Store) FindRepoByIdentity(ctx context.Context, gitCommonDir, worktreeRoot string) (*Repo, error) {
	gitCommonDir = strings.TrimSpace(gitCommonDir)
	worktreeRoot = strings.TrimSpace(worktreeRoot)
	if gitCommonDir != "" {
		repo, err := s.findRepoByGitCommonDir(ctx, gitCommonDir)
		if err != nil || repo != nil {
			return repo, err
		}
	}
	if worktreeRoot != "" {
		return s.FindRepoByRoot(ctx, worktreeRoot)
	}
	return nil, nil
}

func (s *Store) findRepoByGitCommonDir(ctx context.Context, gitCommonDir string) (*Repo, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, root_path, git_common_dir, backend, default_remote, default_branch, authoring_base_ref, remote_url, created_at, updated_at
		FROM repos
		WHERE git_common_dir = ?
		ORDER BY updated_at DESC
		LIMIT 1
	`, gitCommonDir)
	return scanRepoRow(row)
}

func scanRepoRow(row *sql.Row) (*Repo, error) {
	var repo Repo
	var gitCommonDir sql.NullString
	var defaultRemote sql.NullString
	var defaultBranch sql.NullString
	var authoringBase sql.NullString
	var remoteURL sql.NullString
	if err := row.Scan(
		&repo.ID,
		&repo.RootPath,
		&gitCommonDir,
		&repo.Backend,
		&defaultRemote,
		&defaultBranch,
		&authoringBase,
		&remoteURL,
		&repo.CreatedAt,
		&repo.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan repo row: %w", err)
	}
	if gitCommonDir.Valid {
		repo.GitCommonDir = gitCommonDir.String
	}
	if defaultRemote.Valid {
		repo.DefaultRemote = &defaultRemote.String
	}
	if defaultBranch.Valid {
		repo.DefaultBranch = &defaultBranch.String
	}
	if authoringBase.Valid {
		repo.AuthoringBase = &authoringBase.String
	}
	if remoteURL.Valid {
		repo.RemoteURL = &remoteURL.String
	}
	return &repo, nil
}

func (s *Store) ListRepos(ctx context.Context) ([]Repo, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, root_path, git_common_dir, backend, default_remote, default_branch, authoring_base_ref, remote_url, created_at, updated_at
		FROM repos
		ORDER BY updated_at DESC, id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list repos: %w", err)
	}
	defer rows.Close()

	var repos []Repo
	for rows.Next() {
		var repo Repo
		var gitCommonDir sql.NullString
		var defaultRemote sql.NullString
		var defaultBranch sql.NullString
		var authoringBase sql.NullString
		var remoteURL sql.NullString
		if err := rows.Scan(
			&repo.ID,
			&repo.RootPath,
			&gitCommonDir,
			&repo.Backend,
			&defaultRemote,
			&defaultBranch,
			&authoringBase,
			&remoteURL,
			&repo.CreatedAt,
			&repo.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan repo: %w", err)
		}
		if gitCommonDir.Valid {
			repo.GitCommonDir = gitCommonDir.String
		}
		if defaultRemote.Valid {
			repo.DefaultRemote = &defaultRemote.String
		}
		if defaultBranch.Valid {
			repo.DefaultBranch = &defaultBranch.String
		}
		if authoringBase.Valid {
			repo.AuthoringBase = &authoringBase.String
		}
		if remoteURL.Valid {
			repo.RemoteURL = &remoteURL.String
		}
		repos = append(repos, repo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list repos rows: %w", err)
	}
	return repos, nil
}
