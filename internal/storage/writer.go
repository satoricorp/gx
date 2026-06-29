package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/agentprovenance"
)

type StorageWriter interface {
	WriteSession(ctx context.Context, s Session) error
	WriteRequest(ctx context.Context, r Request) error
	WriteResponse(ctx context.Context, resp Response) error
}

type Store struct {
	db          *sql.DB
	sessionStmt *sql.Stmt
	endStmt     *sql.Stmt
	requestStmt *sql.Stmt
	respStmt    *sql.Stmt
}

func NewStore(ctx context.Context, db *sql.DB) (*Store, error) {
	sessionStmt, err := db.PrepareContext(ctx, `
		INSERT INTO sessions (
			id, created_at, ended_at, command, cwd, client_pid, exit_code, gx_version,
			source, process_name, parent_pid, last_seen_at, end_reason, repo_root
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare session insert: %w", err)
	}
	endStmt, err := db.PrepareContext(ctx, `
		UPDATE sessions
		SET ended_at = ?, exit_code = ?
		WHERE id = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare session end: %w", err)
	}
	requestStmt, err := db.PrepareContext(ctx, `
		INSERT INTO requests (id, session_id, created_at, provider, endpoint, method, model, request_body, request_headers)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare request insert: %w", err)
	}
	respStmt, err := db.PrepareContext(ctx, `
		INSERT INTO responses (
			id, request_id, created_at, completed_at, status_code, response_body, response_headers,
			is_streaming, duration_ms, provider_request_id, finish_reason,
			input_tokens, output_tokens, cache_read_tokens, cache_write_tokens, error
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("prepare response insert: %w", err)
	}

	return &Store{
		db:          db,
		sessionStmt: sessionStmt,
		endStmt:     endStmt,
		requestStmt: requestStmt,
		respStmt:    respStmt,
	}, nil
}

func (s *Store) Close() error {
	var firstErr error
	for _, stmt := range []*sql.Stmt{s.sessionStmt, s.endStmt, s.requestStmt, s.respStmt} {
		if stmt == nil {
			continue
		}
		if err := stmt.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if err := s.db.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

func (s *Store) WriteSession(ctx context.Context, session Session) error {
	_, err := s.sessionStmt.ExecContext(
		ctx,
		session.ID,
		session.CreatedAt,
		session.EndedAt,
		session.Command,
		session.Cwd,
		session.ClientPID,
		session.ExitCode,
		session.GXVersion,
		session.Source,
		session.ProcessName,
		session.ParentPID,
		session.LastSeenAt,
		session.EndReason,
		session.RepoRoot,
	)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (s *Store) TouchSession(ctx context.Context, sessionID string, lastSeenAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE sessions
		SET last_seen_at = ?
		WHERE id = ?
	`, lastSeenAt, sessionID)
	if err != nil {
		return fmt.Errorf("touch session: %w", err)
	}
	return nil
}

func (s *Store) EndSession(ctx context.Context, sessionID string, endedAt int64, exitCode int) error {
	_, err := s.endStmt.ExecContext(ctx, endedAt, exitCode, sessionID)
	if err != nil {
		return fmt.Errorf("end session: %w", err)
	}
	return nil
}

func (s *Store) WriteRequest(ctx context.Context, req Request) error {
	_, err := s.requestStmt.ExecContext(
		ctx,
		req.ID,
		req.SessionID,
		req.CreatedAt,
		req.Provider,
		req.Endpoint,
		req.Method,
		req.Model,
		req.RequestBody,
		req.RequestHeaders,
	)
	if err != nil {
		return fmt.Errorf("insert request: %w", err)
	}
	return nil
}

func (s *Store) WriteResponse(ctx context.Context, resp Response) error {
	streaming := 0
	if resp.IsStreaming {
		streaming = 1
	}
	_, err := s.respStmt.ExecContext(
		ctx,
		resp.ID,
		resp.RequestID,
		resp.CreatedAt,
		resp.CompletedAt,
		resp.StatusCode,
		resp.ResponseBody,
		resp.ResponseHeaders,
		streaming,
		resp.DurationMS,
		resp.ProviderRequestID,
		resp.FinishReason,
		resp.InputTokens,
		resp.OutputTokens,
		resp.CacheReadTokens,
		resp.CacheWriteTokens,
		resp.Error,
	)
	if err != nil {
		return fmt.Errorf("insert response: %w", err)
	}
	return nil
}

func (s *Store) UpsertRepo(ctx context.Context, repo Repo) (int64, error) {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO repos (root_path, backend, default_remote, default_branch, authoring_base_ref, remote_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(root_path) DO UPDATE SET
			backend = excluded.backend,
			default_remote = excluded.default_remote,
			default_branch = excluded.default_branch,
			authoring_base_ref = COALESCE(excluded.authoring_base_ref, repos.authoring_base_ref),
			remote_url = excluded.remote_url,
			updated_at = excluded.updated_at
	`,
		repo.RootPath,
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

func (s *Store) RenameStackBookmark(ctx context.Context, repoID int64, oldName, newName, normalizedName string, updatedAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE stacks
		SET bookmark_name = ?,
			name = CASE
				WHEN name = ? THEN ?
				ELSE name
			END,
			updated_at = ?
		WHERE repo_id = ? AND bookmark_name = ?
	`, newName, oldName, normalizedName, updatedAt, repoID, oldName)
	if err != nil {
		return fmt.Errorf("rename stack bookmark: %w", err)
	}
	return nil
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

func (s *Store) DeleteStack(ctx context.Context, stackID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin delete stack: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM stack_changes WHERE stack_id = ?`, stackID); err != nil {
		return fmt.Errorf("delete stack changes: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM stacks WHERE id = ?`, stackID); err != nil {
		return fmt.Errorf("delete stack: %w", err)
	}
	return tx.Commit()
}

func (s *Store) PrunePublishedStack(ctx context.Context, repoID, stackID int64, publishRef string, updatedAt int64) error {
	publishRef = strings.TrimSpace(publishRef)
	if repoID == 0 || stackID == 0 || publishRef == "" {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin prune published stack: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM change_bookmarks
		WHERE bookmark_name = ?
			AND change_id IN (
				SELECT change_id FROM stack_changes WHERE stack_id = ?
			)
	`, publishRef, stackID); err != nil {
		return fmt.Errorf("delete stack change bookmarks: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM pushes
		WHERE repo_id = ? AND branch_name = ?
	`, repoID, publishRef); err != nil {
		return fmt.Errorf("delete stack pushes: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM stack_changes
		WHERE stack_id = ?
	`, stackID); err != nil {
		return fmt.Errorf("delete pruned stack changes: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM stacks
		WHERE repo_id = ? AND id = ?
	`, repoID, stackID); err != nil {
		return fmt.Errorf("delete pruned stack: %w", err)
	}
	return tx.Commit()
}

func (s *Store) MarkChangeStatus(ctx context.Context, changeID int64, status string, updatedAt int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE changes
		SET status = ?, updated_at = ?
		WHERE id = ?
	`, status, updatedAt, changeID)
	if err != nil {
		return fmt.Errorf("mark change status: %w", err)
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

func (s *Store) ListChangeRevisionsByChangeID(ctx context.Context, changeID int64) ([]ChangeRevision, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, change_id, jj_commit_id, jj_operation_id, changed_files_json, created_at
		FROM change_revisions
		WHERE change_id = ?
		ORDER BY created_at ASC, id ASC
	`, changeID)
	if err != nil {
		return nil, fmt.Errorf("list change revisions: %w", err)
	}
	defer rows.Close()
	revisions := []ChangeRevision{}
	for rows.Next() {
		var revision ChangeRevision
		if err := rows.Scan(
			&revision.ID,
			&revision.ChangeID,
			&revision.JJCommitID,
			&revision.JJOperationID,
			&revision.ChangedFiles,
			&revision.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan change revision: %w", err)
		}
		revisions = append(revisions, revision)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change revisions: %w", err)
	}
	return revisions, nil
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

func (s *Store) WriteChangeDemuxEvidence(ctx context.Context, evidence ChangeDemuxEvidence) error {
	useHunks := 0
	if evidence.UseHunks {
		useHunks = 1
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO change_demux_evidence (
			change_id, demux_proposal_id, revision_proposal_id, intent, files_json,
			hunk_ids_json, use_hunks, confidence, provenance_status, evidence_json, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(change_id) DO UPDATE SET
			demux_proposal_id = excluded.demux_proposal_id,
			revision_proposal_id = excluded.revision_proposal_id,
			intent = excluded.intent,
			files_json = excluded.files_json,
			hunk_ids_json = excluded.hunk_ids_json,
			use_hunks = excluded.use_hunks,
			confidence = excluded.confidence,
			provenance_status = excluded.provenance_status,
			evidence_json = excluded.evidence_json,
			created_at = excluded.created_at
	`,
		evidence.ChangeID,
		evidence.DemuxProposalID,
		evidence.RevisionProposalID,
		evidence.Intent,
		evidence.FilesJSON,
		evidence.HunkIDsJSON,
		useHunks,
		evidence.Confidence,
		evidence.ProvenanceStatus,
		evidence.EvidenceJSON,
		evidence.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("write change demux evidence: %w", err)
	}
	return nil
}

func (s *Store) ListChangeDemuxEvidenceByProposal(ctx context.Context, proposalID string) ([]ChangeDemuxEvidence, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, change_id, demux_proposal_id, revision_proposal_id, intent, files_json,
			hunk_ids_json, use_hunks, confidence, provenance_status, evidence_json, created_at
		FROM change_demux_evidence
		WHERE demux_proposal_id = ?
		ORDER BY id ASC
	`, proposalID)
	if err != nil {
		return nil, fmt.Errorf("list change demux evidence: %w", err)
	}
	defer rows.Close()

	var out []ChangeDemuxEvidence
	for rows.Next() {
		var evidence ChangeDemuxEvidence
		var useHunks int
		if err := rows.Scan(
			&evidence.ID,
			&evidence.ChangeID,
			&evidence.DemuxProposalID,
			&evidence.RevisionProposalID,
			&evidence.Intent,
			&evidence.FilesJSON,
			&evidence.HunkIDsJSON,
			&useHunks,
			&evidence.Confidence,
			&evidence.ProvenanceStatus,
			&evidence.EvidenceJSON,
			&evidence.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan change demux evidence: %w", err)
		}
		evidence.UseHunks = useHunks != 0
		out = append(out, evidence)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate change demux evidence: %w", err)
	}
	return out, nil
}

func (s *Store) FindAttachableSessionsForRepo(ctx context.Context, repoRoot string, limit int) ([]string, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}
	childCWD := strings.TrimRight(repoRoot, "/") + "/%"
	rows, err := s.db.QueryContext(ctx, `
		SELECT s.id
		FROM sessions s
		WHERE (
			s.repo_root = ?
			OR s.cwd = ?
			OR (s.repo_root IS NULL AND s.cwd LIKE ?)
		)
		AND NOT EXISTS (
			SELECT 1 FROM change_sessions cs WHERE cs.session_id = s.id
		)
		ORDER BY COALESCE(s.last_seen_at, s.ended_at, s.created_at) DESC, s.created_at DESC
		LIMIT ?
	`, repoRoot, repoRoot, childCWD, limit)
	if err != nil {
		return nil, fmt.Errorf("find attachable sessions for repo: %w", err)
	}
	defer rows.Close()

	var sessionIDs []string
	for rows.Next() {
		var sessionID string
		if err := rows.Scan(&sessionID); err != nil {
			return nil, fmt.Errorf("scan attachable session: %w", err)
		}
		sessionIDs = append(sessionIDs, sessionID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate attachable sessions: %w", err)
	}
	return sessionIDs, nil
}

func (s *Store) UpsertDemuxProposal(ctx context.Context, proposal DemuxProposal) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO demux_proposals (
			id, repo_id, base_change_id, status, payload_json, created_at, updated_at, applied_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			repo_id = excluded.repo_id,
			base_change_id = excluded.base_change_id,
			status = excluded.status,
			payload_json = excluded.payload_json,
			updated_at = excluded.updated_at,
			applied_at = excluded.applied_at
	`,
		proposal.ID,
		proposal.RepoID,
		proposal.BaseChangeID,
		proposal.Status,
		proposal.PayloadJSON,
		proposal.CreatedAt,
		proposal.UpdatedAt,
		proposal.AppliedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert demux proposal: %w", err)
	}
	return nil
}

func (s *Store) FindDemuxProposal(ctx context.Context, id string) (*DemuxProposal, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, base_change_id, status, payload_json, created_at, updated_at, applied_at
		FROM demux_proposals
		WHERE id = ?
	`, id)
	return scanDemuxProposal(row)
}

func (s *Store) FindLatestDemuxProposal(ctx context.Context, repoID int64, status string) (*DemuxProposal, error) {
	query := `
		SELECT id, repo_id, base_change_id, status, payload_json, created_at, updated_at, applied_at
		FROM demux_proposals
		WHERE repo_id = ?
	`
	args := []any{repoID}
	if strings.TrimSpace(status) != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC, updated_at DESC LIMIT 1`
	row := s.db.QueryRowContext(ctx, query, args...)
	return scanDemuxProposal(row)
}

func (s *Store) ListDemuxProposals(ctx context.Context, repoID int64, status string, limit int) ([]DemuxProposal, error) {
	query := `
		SELECT id, repo_id, base_change_id, status, payload_json, created_at, updated_at, applied_at
		FROM demux_proposals
		WHERE repo_id = ?
	`
	args := []any{repoID}
	if strings.TrimSpace(status) != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC, updated_at DESC`
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list demux proposals: %w", err)
	}
	defer rows.Close()
	var proposals []DemuxProposal
	for rows.Next() {
		proposal, err := scanDemuxProposal(rows)
		if err != nil {
			return nil, err
		}
		if proposal != nil {
			proposals = append(proposals, *proposal)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate demux proposals: %w", err)
	}
	return proposals, nil
}

func (s *Store) DeletePendingDemuxProposalsForRepo(ctx context.Context, repoID int64) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM change_demux_evidence
		WHERE demux_proposal_id IN (
			SELECT id FROM demux_proposals WHERE repo_id = ? AND status = 'pending'
		)
	`, repoID)
	if err != nil {
		return fmt.Errorf("delete pending demux evidence: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		DELETE FROM demux_proposals
		WHERE repo_id = ? AND status = 'pending'
	`, repoID)
	if err != nil {
		return fmt.Errorf("delete pending demux proposals: %w", err)
	}
	return nil
}

func (s *Store) FilterExistingSessionIDs(ctx context.Context, sessionIDs []string) ([]string, error) {
	if len(sessionIDs) == 0 {
		return nil, nil
	}
	seen := map[string]struct{}{}
	ordered := make([]string, 0, len(sessionIDs))
	for _, sessionID := range sessionIDs {
		sessionID = strings.TrimSpace(sessionID)
		if sessionID == "" {
			continue
		}
		if _, ok := seen[sessionID]; ok {
			continue
		}
		seen[sessionID] = struct{}{}
		ordered = append(ordered, sessionID)
	}
	if len(ordered) == 0 {
		return nil, nil
	}
	placeholders := strings.Repeat("?,", len(ordered))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]any, len(ordered))
	for i, sessionID := range ordered {
		args[i] = sessionID
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id FROM sessions WHERE id IN (`+placeholders+`)
	`, args...)
	if err != nil {
		return nil, fmt.Errorf("filter existing session ids: %w", err)
	}
	defer rows.Close()
	existing := map[string]struct{}{}
	for rows.Next() {
		var sessionID string
		if err := rows.Scan(&sessionID); err != nil {
			return nil, fmt.Errorf("scan existing session id: %w", err)
		}
		existing[sessionID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate existing session ids: %w", err)
	}
	filtered := make([]string, 0, len(ordered))
	for _, sessionID := range ordered {
		if _, ok := existing[sessionID]; ok {
			filtered = append(filtered, sessionID)
		}
	}
	return filtered, nil
}

func (s *Store) UpdateDemuxProposalStatus(ctx context.Context, id, status string, updatedAt int64, appliedAt *int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE demux_proposals
		SET status = ?, updated_at = ?, applied_at = ?
		WHERE id = ?
	`, status, updatedAt, appliedAt, id)
	if err != nil {
		return fmt.Errorf("update demux proposal status: %w", err)
	}
	return nil
}

func (s *Store) UpdateSessionWorkspace(ctx context.Context, sessionID, cwd, repoRoot string) error {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil
	}
	cwd = strings.TrimSpace(cwd)
	repoRoot = strings.TrimSpace(repoRoot)
	if cwd == "" && repoRoot == "" {
		return nil
	}
	_, err := s.db.ExecContext(ctx, `
		UPDATE sessions
		SET
			cwd = CASE WHEN ? != '' THEN ? ELSE cwd END,
			repo_root = CASE WHEN ? != '' THEN ? ELSE repo_root END
		WHERE id = ?
	`, cwd, cwd, repoRoot, repoRoot, sessionID)
	if err != nil {
		return fmt.Errorf("update session workspace: %w", err)
	}
	return nil
}

func (s *Store) WriteModifyEvent(ctx context.Context, ev ModifyEvent) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO modify_events (repo_id, target_change_id, previous_current_change_id, jj_operation_id, created_at)
		VALUES (?, ?, ?, ?, ?)
	`,
		ev.RepoID,
		ev.TargetChangeID,
		ev.PreviousCurrentChangeID,
		ev.JJOperationID,
		ev.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert modify event: %w", err)
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

func (s *Store) WritePush(ctx context.Context, push Push) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO pushes (repo_id, remote_name, branch_name, head_commit_id, current_change_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		push.RepoID,
		push.RemoteName,
		push.BranchName,
		push.HeadCommitID,
		push.CurrentChangeID,
		push.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert push: %w", err)
	}
	return nil
}

func (s *Store) UpsertCursorSession(ctx context.Context, session Session) (bool, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO sessions (
			id, created_at, ended_at, command, cwd, client_pid, exit_code, gx_version,
			source, process_name, parent_pid, last_seen_at, end_reason, repo_root
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		session.ID,
		session.CreatedAt,
		session.EndedAt,
		session.Command,
		session.Cwd,
		session.ClientPID,
		session.ExitCode,
		session.GXVersion,
		session.Source,
		session.ProcessName,
		session.ParentPID,
		session.LastSeenAt,
		session.EndReason,
		session.RepoRoot,
	)
	if err != nil {
		return false, fmt.Errorf("upsert cursor session: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("cursor session rows affected: %w", err)
	}
	return affected > 0, nil
}

func (s *Store) UpsertCursorMessage(ctx context.Context, msg CursorMessage) (bool, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO cursor_messages (
			id, session_id, created_at, role, text, raw_json, input_tokens, output_tokens
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		msg.ID,
		msg.SessionID,
		msg.CreatedAt,
		msg.Role,
		msg.Text,
		msg.RawJSON,
		msg.InputTokens,
		msg.OutputTokens,
	)
	if err != nil {
		return false, fmt.Errorf("upsert cursor message: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("cursor message rows affected: %w", err)
	}
	return affected > 0, nil
}

func (s *Store) CountCursorSessions(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM sessions WHERE source = 'cursor'`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count cursor sessions: %w", err)
	}
	return count, nil
}

func (s *Store) CountCursorMessages(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT count(*) FROM cursor_messages`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count cursor messages: %w", err)
	}
	return count, nil
}

func (s *Store) AgentLedgerSummary(ctx context.Context, agent string) (AgentLedgerSummary, error) {
	summary := AgentLedgerSummary{Agent: agent}
	var where string
	switch agent {
	case "cursor":
		where = `COALESCE(s.source, '') = 'cursor'`
	case "codex":
		where = `(LOWER(COALESCE(s.command, '')) LIKE '%codex%' OR LOWER(COALESCE(s.process_name, '')) LIKE '%codex%')`
	case "claude":
		where = `(LOWER(COALESCE(s.command, '')) LIKE '%claude%' OR LOWER(COALESCE(s.process_name, '')) LIKE '%claude%')`
	default:
		return summary, fmt.Errorf("unsupported agent %q", agent)
	}

	row := s.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT
			COALESCE(
				(SELECT s2.cwd FROM sessions s2 WHERE %s AND s2.cwd != '' ORDER BY COALESCE(s2.last_seen_at, s2.ended_at, s2.created_at) DESC LIMIT 1),
				''
			),
			COALESCE(COUNT(DISTINCT s.id), 0),
			COALESCE(SUM(COALESCE(resp.input_tokens, 0) + COALESCE(resp.output_tokens, 0)), 0),
			COALESCE(COUNT(DISTINCT COALESCE(NULLIF(s.repo_root, ''), s.cwd)), 0),
			MAX(COALESCE(s.last_seen_at, s.ended_at, s.created_at))
		FROM sessions s
		LEFT JOIN requests req ON req.session_id = s.id
		LEFT JOIN responses resp ON resp.request_id = req.id
		WHERE %s
	`, strings.ReplaceAll(where, "s.", "s2."), where))

	var filepath string
	var lastSeen sql.NullInt64
	err := row.Scan(&filepath, &summary.Calls, &summary.Tokens, &summary.Files, &lastSeen)
	if err == sql.ErrNoRows {
		return summary, nil
	}
	if err != nil {
		return summary, fmt.Errorf("agent ledger summary %s: %w", agent, err)
	}
	summary.Filepath = filepath
	if lastSeen.Valid {
		value := lastSeen.Int64
		summary.LastSeenAt = &value
	}
	missingTokens, err := s.agentLedgerMissingTokens(ctx, where)
	if err != nil {
		return summary, err
	}
	summary.Tokens += missingTokens
	return summary, nil
}

func (s *Store) agentLedgerMissingTokens(ctx context.Context, where string) (int, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT resp.response_body
		FROM sessions s
		JOIN requests req ON req.session_id = s.id
		JOIN responses resp ON resp.request_id = req.id
		WHERE %s
			AND resp.input_tokens IS NULL
			AND resp.output_tokens IS NULL
	`, where))
	if err != nil {
		return 0, fmt.Errorf("agent ledger missing tokens: %w", err)
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		var body []byte
		if err := rows.Scan(&body); err != nil {
			return 0, fmt.Errorf("scan missing token response: %w", err)
		}
		total += responseBodyTokens(body)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterate missing token responses: %w", err)
	}
	return total, nil
}

func responseBodyTokens(body []byte) int {
	latest := 0
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
		if line == "" || line == "[DONE]" || !strings.HasPrefix(line, "{") {
			continue
		}
		if tokens := responseJSONTokens([]byte(line)); tokens > 0 {
			latest = tokens
		}
	}
	if latest > 0 {
		return latest
	}
	return responseJSONTokens(body)
}

func responseJSONTokens(body []byte) int {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0
	}
	if response, _ := payload["response"].(map[string]any); response != nil {
		if tokens := usageTokens(response["usage"]); tokens > 0 {
			return tokens
		}
	}
	return usageTokens(payload["usage"])
}

func usageTokens(value any) int {
	usage, _ := value.(map[string]any)
	if usage == nil {
		return 0
	}
	return usageTokenValue(usage, "input_tokens", "prompt_tokens") +
		usageTokenValue(usage, "output_tokens", "completion_tokens")
}

func usageTokenValue(usage map[string]any, keys ...string) int {
	for _, key := range keys {
		switch value := usage[key].(type) {
		case float64:
			return int(value)
		case int:
			return value
		}
	}
	return 0
}

func (s *Store) FindRepoByRoot(ctx context.Context, rootPath string) (*Repo, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, root_path, backend, default_remote, default_branch, authoring_base_ref, remote_url, created_at, updated_at
		FROM repos
		WHERE root_path = ?
	`, rootPath)

	var repo Repo
	var defaultRemote sql.NullString
	var defaultBranch sql.NullString
	var authoringBase sql.NullString
	var remoteURL sql.NullString
	if err := row.Scan(
		&repo.ID,
		&repo.RootPath,
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
		SELECT id, root_path, backend, default_remote, default_branch, authoring_base_ref, remote_url, created_at, updated_at
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
		var defaultRemote sql.NullString
		var defaultBranch sql.NullString
		var authoringBase sql.NullString
		var remoteURL sql.NullString
		if err := rows.Scan(
			&repo.ID,
			&repo.RootPath,
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

type demuxProposalScanner interface {
	Scan(dest ...any) error
}

func scanDemuxProposal(row demuxProposalScanner) (*DemuxProposal, error) {
	var proposal DemuxProposal
	var appliedAt sql.NullInt64
	if err := row.Scan(
		&proposal.ID,
		&proposal.RepoID,
		&proposal.BaseChangeID,
		&proposal.Status,
		&proposal.PayloadJSON,
		&proposal.CreatedAt,
		&proposal.UpdatedAt,
		&appliedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan demux proposal: %w", err)
	}
	if appliedAt.Valid {
		proposal.AppliedAt = &appliedAt.Int64
	}
	return &proposal, nil
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

func (s *Store) FindStackByHeadChange(ctx context.Context, repoID int64, headChangeID string) (*Stack, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, name, bookmark_name, base_ref, base_commit_id, head_change_id, head_commit_id,
			remote_name, remote_ref, github_pr_url, status, created_at, updated_at
		FROM stacks
		WHERE repo_id = ? AND head_change_id = ?
		ORDER BY updated_at DESC, id DESC
		LIMIT 1
	`, repoID, headChangeID)
	return scanStack(row)
}

func (s *Store) LatestStackByRepoID(ctx context.Context, repoID int64) (*Stack, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, name, bookmark_name, base_ref, base_commit_id, head_change_id, head_commit_id,
			remote_name, remote_ref, github_pr_url, status, created_at, updated_at
		FROM stacks
		WHERE repo_id = ?
		ORDER BY updated_at DESC, id DESC
		LIMIT 1
	`, repoID)
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

func (s *Store) ResetAllPublished(ctx context.Context) (int, error) {
	now := time.Now().UnixMilli()
	result, err := s.db.ExecContext(ctx, `
		UPDATE stacks
		SET status = 'draft', updated_at = ?
		WHERE status != 'draft'
	`, now)
	if err != nil {
		return 0, fmt.Errorf("reset stack publish state: %w", err)
	}
	stacksUpdated, _ := result.RowsAffected()

	if _, err := s.db.ExecContext(ctx, `DELETE FROM pushes`); err != nil {
		return 0, fmt.Errorf("clear push history: %w", err)
	}
	return int(stacksUpdated), nil
}

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

func (s *Store) LatestPushByBranchName(ctx context.Context, repoID int64, branchName string) (*Push, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, repo_id, remote_name, branch_name, head_commit_id, current_change_id, created_at
		FROM pushes
		WHERE repo_id = ? AND branch_name = ?
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, repoID, branchName)

	var push Push
	var remoteName sql.NullString
	var storedBranchName sql.NullString
	var currentChangeID sql.NullInt64
	if err := row.Scan(
		&push.ID,
		&push.RepoID,
		&remoteName,
		&storedBranchName,
		&push.HeadCommitID,
		&currentChangeID,
		&push.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("latest push by branch: %w", err)
	}
	if remoteName.Valid {
		push.RemoteName = &remoteName.String
	}
	if storedBranchName.Valid {
		push.BranchName = &storedBranchName.String
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

func uniquePositiveInt64s(values []int64) []int64 {
	seen := map[int64]struct{}{}
	out := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func uniqueNonEmptyStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

func placeholdersForInt64s(values []int64) (string, []any) {
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(values)), ",")
	args := make([]any, len(values))
	for i, value := range values {
		args[i] = value
	}
	return placeholders, args
}

func placeholdersForStrings(values []string) (string, []any) {
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(values)), ",")
	args := make([]any, len(values))
	for i, value := range values {
		args[i] = value
	}
	return placeholders, args
}
