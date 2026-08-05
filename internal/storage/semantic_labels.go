package storage

import (
	"context"
	"fmt"
)

func (s *Store) UpsertSemanticLabel(ctx context.Context, label SemanticLabel) (int64, error) {
	return upsertSemanticLabel(ctx, s.db, label)
}

func upsertSemanticLabel(ctx context.Context, exec sqlExecutor, label SemanticLabel) (int64, error) {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO semantic_labels (
			repo_id, label, aliases_json, sources_json, seen_count, accepted_count,
			confidence, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(repo_id, label) DO UPDATE SET
			aliases_json = excluded.aliases_json,
			sources_json = excluded.sources_json,
			seen_count = semantic_labels.seen_count + excluded.seen_count,
			accepted_count = semantic_labels.accepted_count + excluded.accepted_count,
			confidence = CASE
				WHEN excluded.confidence > semantic_labels.confidence THEN excluded.confidence
				ELSE semantic_labels.confidence
			END,
			updated_at = excluded.updated_at
	`, label.RepoID, label.Label, label.AliasesJSON, label.SourcesJSON, label.SeenCount, label.AcceptedCount, label.Confidence, label.CreatedAt, label.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("upsert semantic label: %w", err)
	}
	var id int64
	if err := exec.QueryRowContext(ctx, `
		SELECT id FROM semantic_labels WHERE repo_id = ? AND label = ?
	`, label.RepoID, label.Label).Scan(&id); err != nil {
		return 0, fmt.Errorf("lookup semantic label: %w", err)
	}
	if err := upsertSemanticLabelFTS(ctx, exec, id, label); err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Store) WriteSemanticLabelLink(ctx context.Context, link SemanticLabelLink) error {
	return writeSemanticLabelLink(ctx, s.db, link)
}

func writeSemanticLabelLink(ctx context.Context, exec sqlExecutor, link SemanticLabelLink) error {
	accepted := 0
	if link.Accepted {
		accepted = 1
	}
	_, err := exec.ExecContext(ctx, `
		INSERT INTO semantic_label_links (
			repo_id, label_id, demux_proposal_id, revision_proposal_id, change_id,
			stack_bookmark, source, score, accepted, evidence_json, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, link.RepoID, link.LabelID, link.DemuxProposalID, link.RevisionProposalID, link.ChangeID, link.StackBookmark, link.Source, link.Score, accepted, link.EvidenceJSON, link.CreatedAt)
	if err != nil {
		return fmt.Errorf("write semantic label link: %w", err)
	}
	return nil
}

func upsertSemanticLabelFTS(ctx context.Context, exec sqlExecutor, labelID int64, label SemanticLabel) error {
	_, err := exec.ExecContext(ctx, `
		INSERT OR REPLACE INTO semantic_label_fts (
			rowid, repo_id, label, aliases_json, sources_json
		)
		VALUES (?, ?, ?, ?, ?)
	`, labelID, label.RepoID, label.Label, label.AliasesJSON, label.SourcesJSON)
	if err != nil {
		return fmt.Errorf("upsert semantic label fts: %w", err)
	}
	return nil
}

func (s *Store) WriteSemanticLabels(ctx context.Context, labels []SemanticLabelWrite) error {
	if len(labels) == 0 {
		return nil
	}
	lgtm, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin semantic labels: %w", err)
	}
	defer lgtm.Rollback()
	for _, write := range labels {
		labelID, err := upsertSemanticLabel(ctx, lgtm, write.Label)
		if err != nil {
			return err
		}
		link := write.Link
		link.LabelID = labelID
		if err := writeSemanticLabelLink(ctx, lgtm, link); err != nil {
			return err
		}
	}
	if err := lgtm.Commit(); err != nil {
		return fmt.Errorf("commit semantic labels: %w", err)
	}
	return nil
}

func (s *Store) ListSemanticLabelsByRepoID(ctx context.Context, repoID int64, limit int) ([]SemanticLabel, error) {
	if limit <= 0 {
		limit = 200
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_id, label, aliases_json, sources_json, seen_count, accepted_count,
			confidence, created_at, updated_at
		FROM semantic_labels
		WHERE repo_id = ?
		ORDER BY accepted_count DESC, seen_count DESC, updated_at DESC
		LIMIT ?
	`, repoID, limit)
	if err != nil {
		return nil, fmt.Errorf("list semantic labels: %w", err)
	}
	defer rows.Close()

	var out []SemanticLabel
	for rows.Next() {
		var label SemanticLabel
		if err := rows.Scan(
			&label.ID,
			&label.RepoID,
			&label.Label,
			&label.AliasesJSON,
			&label.SourcesJSON,
			&label.SeenCount,
			&label.AcceptedCount,
			&label.Confidence,
			&label.CreatedAt,
			&label.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan semantic label: %w", err)
		}
		out = append(out, label)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate semantic labels: %w", err)
	}
	return out, nil
}
