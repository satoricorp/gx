package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

type changeRow struct {
	ID              int64
	RepoID          int64
	JJChangeID      string
	CurrentCommitID string
	Description     string
	ParentChangeID  sql.NullString
	Status          string
	FirstSeenAt     int64
	UpdatedAt       int64
}

func DefaultDir() (string, error) {
	if override := os.Getenv("GX_HOME"); override != "" {
		return override, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".gx"), nil
}

func DefaultDBPath() (string, error) {
	dir, err := DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "gx.db"), nil
}

func Open(ctx context.Context) (*sql.DB, error) {
	dbPath, err := DefaultDBPath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create gx dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA busy_timeout=5000;",
	}
	for _, stmt := range pragmas {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("apply pragma %q: %w", stmt, err)
		}
	}

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("read schema: %w", err)
	}
	if _, err := db.ExecContext(ctx, string(schema)); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	if err := ensureSessionColumns(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate sessions table: %w", err)
	}
	if err := ensureChangeSessionProvenanceTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate change_session_provenance table: %w", err)
	}
	if err := repairChangeRows(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("repair change rows: %w", err)
	}
	if err := ensureRepoColumns(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate repos table: %w", err)
	}
	if err := ensureChangeBookmarksTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate change_bookmarks table: %w", err)
	}
	if err := ensureStacksTables(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate stacks tables: %w", err)
	}
	if err := ensureCloudBookmarksTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate cloud_bookmarks table: %w", err)
	}
	if err := ensureDemuxProposalsTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate demux_proposals table: %w", err)
	}
	if err := ensureChangeDemuxEvidenceTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate change_demux_evidence table: %w", err)
	}
	if err := ensureCaptureStagingTables(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate capture staging tables: %w", err)
	}
	return db, nil
}

func ensureRepoColumns(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `PRAGMA table_info(repos)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns := map[string]bool{}
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		columns[name] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if !columns["authoring_base_ref"] {
		if _, err := db.ExecContext(ctx, `ALTER TABLE repos ADD COLUMN authoring_base_ref TEXT`); err != nil {
			return err
		}
	}
	return nil
}

func ensureDemuxProposalsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS demux_proposals (
			id TEXT PRIMARY KEY,
			repo_id INTEGER NOT NULL REFERENCES repos(id),
			base_change_id TEXT NOT NULL,
			status TEXT NOT NULL,
			payload_json TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			applied_at INTEGER
		);
		CREATE INDEX IF NOT EXISTS idx_demux_proposals_repo ON demux_proposals(repo_id, status, created_at);
	`)
	return err
}

func ensureChangeDemuxEvidenceTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS change_demux_evidence (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			change_id INTEGER NOT NULL REFERENCES changes(id),
			demux_proposal_id TEXT NOT NULL REFERENCES demux_proposals(id),
			revision_proposal_id TEXT NOT NULL,
			intent TEXT NOT NULL,
			files_json TEXT NOT NULL,
			hunk_ids_json TEXT NOT NULL,
			use_hunks INTEGER NOT NULL,
			confidence REAL NOT NULL,
			provenance_status TEXT NOT NULL,
			evidence_json TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			UNIQUE(change_id)
		);
		CREATE INDEX IF NOT EXISTS idx_change_demux_evidence_proposal ON change_demux_evidence(demux_proposal_id);
	`)
	return err
}

func ensureChangeSessionProvenanceTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS change_session_provenance (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			change_id INTEGER NOT NULL REFERENCES changes(id),
			session_id TEXT NOT NULL REFERENCES sessions(id),
			agent_tool TEXT NOT NULL,
			provider TEXT NOT NULL,
			model_id TEXT NOT NULL,
			source TEXT,
			process_name TEXT,
			created_at INTEGER NOT NULL,
			UNIQUE(change_id, session_id, agent_tool, provider, model_id)
		);
		CREATE INDEX IF NOT EXISTS idx_change_session_provenance_change ON change_session_provenance(change_id);
		CREATE INDEX IF NOT EXISTS idx_change_session_provenance_session ON change_session_provenance(session_id);
	`)
	return err
}

func ensureStacksTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS stacks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_id INTEGER NOT NULL REFERENCES repos(id),
			name TEXT NOT NULL,
			bookmark_name TEXT NOT NULL,
			base_ref TEXT NOT NULL,
			base_commit_id TEXT NOT NULL,
			head_change_id TEXT,
			head_commit_id TEXT,
			remote_name TEXT,
			remote_ref TEXT,
			github_pr_url TEXT,
			status TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			UNIQUE(repo_id, bookmark_name)
		);
		CREATE INDEX IF NOT EXISTS idx_stacks_repo ON stacks(repo_id);
		CREATE INDEX IF NOT EXISTS idx_stacks_head_change ON stacks(repo_id, head_change_id);
		CREATE TABLE IF NOT EXISTS stack_changes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stack_id INTEGER NOT NULL REFERENCES stacks(id),
			change_id INTEGER NOT NULL REFERENCES changes(id),
			jj_change_id TEXT,
			commit_id TEXT,
			position INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			UNIQUE(stack_id, change_id)
		);
		CREATE INDEX IF NOT EXISTS idx_stack_changes_stack ON stack_changes(stack_id, position);
		CREATE INDEX IF NOT EXISTS idx_stack_changes_change ON stack_changes(change_id);
	`)
	if err != nil {
		return err
	}
	if err := ensureColumn(ctx, db, "stacks", "head_commit_id", "TEXT"); err != nil {
		return err
	}
	if err := ensureColumn(ctx, db, "stack_changes", "jj_change_id", "TEXT"); err != nil {
		return err
	}
	if err := ensureColumn(ctx, db, "stack_changes", "commit_id", "TEXT"); err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `
		UPDATE stacks
		SET head_commit_id = (
			SELECT current_commit_id
			FROM changes
			WHERE changes.repo_id = stacks.repo_id
				AND changes.jj_change_id = stacks.head_change_id
		)
		WHERE (head_commit_id IS NULL OR head_commit_id = '')
			AND head_change_id IS NOT NULL
			AND head_change_id != '';
		UPDATE stack_changes
		SET jj_change_id = (
				SELECT jj_change_id
				FROM changes
				WHERE changes.id = stack_changes.change_id
			),
			commit_id = (
				SELECT current_commit_id
				FROM changes
				WHERE changes.id = stack_changes.change_id
			)
		WHERE jj_change_id IS NULL OR jj_change_id = ''
			OR commit_id IS NULL OR commit_id = '';
		CREATE INDEX IF NOT EXISTS idx_stacks_head_commit ON stacks(repo_id, head_commit_id);
	`)
	return err
}

func ensureColumn(ctx context.Context, db *sql.DB, table, name, typ string) error {
	rows, err := db.QueryContext(ctx, `PRAGMA table_info(`+table+`)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var columnName, columnType string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &columnName, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if columnName == name {
			return rows.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, name, typ))
	return err
}

func ensureChangeBookmarksTable(ctx context.Context, db *sql.DB) error {
	var name string
	err := db.QueryRowContext(ctx, `
		SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'change_bookmarks'
	`).Scan(&name)
	if err == nil {
		return nil
	}
	if err != sql.ErrNoRows {
		return err
	}
	_, err = db.ExecContext(ctx, `
		CREATE TABLE change_bookmarks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			change_id INTEGER NOT NULL REFERENCES changes(id),
			bookmark_name TEXT NOT NULL,
			remote_name TEXT,
			remote_ref TEXT,
			last_pushed_commit_id TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			UNIQUE(change_id, bookmark_name)
		);
		CREATE INDEX IF NOT EXISTS idx_change_bookmarks_change ON change_bookmarks(change_id);
	`)
	return err
}

func ensureSessionColumns(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `PRAGMA table_info(sessions)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := map[string]struct{}{}
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		existing[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	columns := map[string]string{
		"source":       "TEXT",
		"process_name": "TEXT",
		"parent_pid":   "INTEGER",
		"last_seen_at": "INTEGER",
		"end_reason":   "TEXT",
		"repo_root":    "TEXT",
	}
	for name, typ := range columns {
		if _, ok := existing[name]; ok {
			continue
		}
		if _, err := db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE sessions ADD COLUMN %s %s", name, typ)); err != nil {
			return err
		}
	}
	return nil
}

func repairChangeRows(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `
		SELECT id, repo_id, jj_change_id, current_commit_id, description, parent_change_id, status, first_seen_at, updated_at
		FROM changes
		WHERE jj_change_id LIKE 'Done importing changes from the underlying Git repo.%'
		ORDER BY id ASC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var dirtyRows []changeRow
	for rows.Next() {
		var row changeRow
		if err := rows.Scan(
			&row.ID,
			&row.RepoID,
			&row.JJChangeID,
			&row.CurrentCommitID,
			&row.Description,
			&row.ParentChangeID,
			&row.Status,
			&row.FirstSeenAt,
			&row.UpdatedAt,
		); err != nil {
			return err
		}
		dirtyRows = append(dirtyRows, row)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, dirty := range dirtyRows {
		cleanID := normalizeJJChangeID(dirty.JJChangeID)
		if cleanID == "" || cleanID == dirty.JJChangeID {
			continue
		}
		if err := repairChangeRow(ctx, db, dirty, cleanID); err != nil {
			return err
		}
	}

	return nil
}

func repairChangeRow(ctx context.Context, db *sql.DB, dirty changeRow, cleanJJChangeID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var canonicalID int64
	err = tx.QueryRowContext(ctx, `
		SELECT id
		FROM changes
		WHERE repo_id = ? AND jj_change_id = ?
	`, dirty.RepoID, cleanJJChangeID).Scan(&canonicalID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		if _, err = tx.ExecContext(ctx, `
			UPDATE changes
			SET jj_change_id = ?
			WHERE id = ?
		`, cleanJJChangeID, dirty.ID); err != nil {
			return err
		}
		return tx.Commit()
	}

	if canonicalID == dirty.ID {
		return tx.Commit()
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE changes
		SET current_commit_id = ?, description = ?, parent_change_id = ?, status = ?, updated_at = ?
		WHERE id = ?
	`, dirty.CurrentCommitID, dirty.Description, nullableStringValue(dirty.ParentChangeID), dirty.Status, dirty.UpdatedAt, canonicalID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE pushes
		SET current_change_id = ?
		WHERE current_change_id = ?
	`, canonicalID, dirty.ID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE modify_events
		SET target_change_id = ?
		WHERE target_change_id = ?
	`, canonicalID, dirty.ID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE modify_events
		SET previous_current_change_id = ?
		WHERE previous_current_change_id = ?
	`, canonicalID, dirty.ID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO change_revisions (change_id, jj_commit_id, jj_operation_id, changed_files_json, created_at)
		SELECT ?, jj_commit_id, jj_operation_id, changed_files_json, created_at
		FROM change_revisions
		WHERE change_id = ?
	`, canonicalID, dirty.ID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM change_revisions WHERE change_id = ?`, dirty.ID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO change_sessions (change_id, session_id, created_at)
		SELECT ?, session_id, created_at
		FROM change_sessions
		WHERE change_id = ?
	`, canonicalID, dirty.ID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM change_sessions WHERE change_id = ?`, dirty.ID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO change_session_provenance (
			change_id, session_id, agent_tool, provider, model_id, source, process_name, created_at
		)
		SELECT ?, session_id, agent_tool, provider, model_id, source, process_name, created_at
		FROM change_session_provenance
		WHERE change_id = ?
	`, canonicalID, dirty.ID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM change_session_provenance WHERE change_id = ?`, dirty.ID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO change_demux_evidence (
			change_id, demux_proposal_id, revision_proposal_id, intent, files_json,
			hunk_ids_json, use_hunks, confidence, provenance_status, evidence_json, created_at
		)
		SELECT ?, demux_proposal_id, revision_proposal_id, intent, files_json,
			hunk_ids_json, use_hunks, confidence, provenance_status, evidence_json, created_at
		FROM change_demux_evidence
		WHERE change_id = ?
	`, canonicalID, dirty.ID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM change_demux_evidence WHERE change_id = ?`, dirty.ID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM changes WHERE id = ?`, dirty.ID); err != nil {
		return err
	}

	return tx.Commit()
}

func normalizeJJChangeID(value string) string {
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, "Done importing changes from the underlying Git repo.") {
		return value
	}
	lines := strings.Split(value, "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

func nullableStringValue(value sql.NullString) any {
	if value.Valid {
		return value.String
	}
	return nil
}
