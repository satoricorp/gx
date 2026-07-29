package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) WriteChangeDemuxEvidence(ctx context.Context, evidence ChangeDemuxEvidence) error {
	if err := writeChangeDemuxEvidence(ctx, s.db, evidence); err != nil {
		return err
	}
	return nil
}

func writeChangeDemuxEvidence(ctx context.Context, exec sqlExecutor, evidence ChangeDemuxEvidence) error {
	useHunks := 0
	if evidence.UseHunks {
		useHunks = 1
	}
	_, err := exec.ExecContext(ctx, `
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

func (s *Store) WriteChangeDemuxEvidenceWithSemanticLabels(ctx context.Context, evidence ChangeDemuxEvidence, labels []SemanticLabelWrite) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin demux evidence semantic labels: %w", err)
	}
	defer tx.Rollback()

	if err := writeChangeDemuxEvidence(ctx, tx, evidence); err != nil {
		return err
	}
	for _, write := range labels {
		labelID, err := upsertSemanticLabel(ctx, tx, write.Label)
		if err != nil {
			return err
		}
		link := write.Link
		link.LabelID = labelID
		if err := writeSemanticLabelLink(ctx, tx, link); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit demux evidence semantic labels: %w", err)
	}
	return nil
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
