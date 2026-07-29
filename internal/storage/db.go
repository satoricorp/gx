package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"os/exec"
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

// DefaultDir is where tl keeps its machine-wide state: $TOTALITY_HOME when set,
// otherwise ~/.totality.
//
// The override is trimmed because every other reader of TOTALITY_HOME trims it —
// internal/auth, internal/semantic, internal/capture, internal/hooks — and this
// one did not. A whitespace-only value put the database in a directory named
// two spaces while the rest of tl carried on using ~/.totality, which is the
// reader-and-writer-disagree shape this codebase has already paid for more than
// once.
func DefaultDir() (string, error) {
	if override := strings.TrimSpace(os.Getenv("TOTALITY_HOME")); override != "" {
		return override, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".totality"), nil
}

func DefaultDBPath() (string, error) {
	dir, err := DefaultDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "totality.db"), nil
}

func Open(ctx context.Context) (*sql.DB, error) {
	dbPath, err := DefaultDBPath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create tl dir: %w", err)
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
	// After the column migration, never in schema.sql: repo_root arrived as a
	// migration column, so an older database reaches schema.sql without it.
	if err := ensureSessionIndexes(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate sessions indexes: %w", err)
	}
	if err := ensureChangeSessionProvenanceTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate change_session_provenance table: %w", err)
	}
	if err := ensureSessionContextsTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate session_contexts table: %w", err)
	}
	if err := ensureSessionEventAttributionsTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate session_event_attributions table: %w", err)
	}
	// Immediately before the repair pass so its five stale change_sessions
	// rows cascade off in the same open.
	if err := deleteCommitSelfReportSession(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("remove tl commit self-report session: %w", err)
	}
	if err := repairDanglingSessionLinks(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("repair session links: %w", err)
	}
	if err := repairChangeRows(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("repair change rows: %w", err)
	}
	if err := ensureRepoColumns(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate repos table: %w", err)
	}
	if err := ensureRepoGitCommonDirColumn(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate repos git_common_dir: %w", err)
	}
	if err := ensureCaptureSessionBindingColumns(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate capture_sessions binding columns: %w", err)
	}
	if err := ensureChangeBookmarksTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate change_bookmarks table: %w", err)
	}
	if err := ensureStacksTables(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate stacks tables: %w", err)
	}
	if err := ensureDemuxProposalsTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate demux_proposals table: %w", err)
	}
	if err := ensureChangeDemuxEvidenceTable(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate change_demux_evidence table: %w", err)
	}
	if err := ensureSemanticLabelTables(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate semantic label tables: %w", err)
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

func ensureRepoGitCommonDirColumn(ctx context.Context, db *sql.DB) error {
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
	if !columns["git_common_dir"] {
		if _, err := db.ExecContext(ctx, `ALTER TABLE repos ADD COLUMN git_common_dir TEXT`); err != nil {
			return err
		}
	}
	repos, err := db.QueryContext(ctx, `
		SELECT id, root_path
		FROM repos
		WHERE git_common_dir IS NULL OR git_common_dir = ''
	`)
	if err != nil {
		return err
	}
	type repoPath struct {
		id   int64
		root string
	}
	var pending []repoPath
	for repos.Next() {
		var repo repoPath
		if err := repos.Scan(&repo.id, &repo.root); err != nil {
			repos.Close()
			return err
		}
		pending = append(pending, repo)
	}
	if err := repos.Close(); err != nil {
		return err
	}
	for _, repo := range pending {
		commonDir := filepath.Clean(repo.root)
		cmd := exec.CommandContext(ctx, "git", "-C", repo.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
		if out, commandErr := cmd.Output(); commandErr == nil && strings.TrimSpace(string(out)) != "" {
			commonDir = filepath.Clean(strings.TrimSpace(string(out)))
		}
		if _, err := db.ExecContext(ctx, `UPDATE repos SET git_common_dir = ? WHERE id = ?`, commonDir, repo.id); err != nil {
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

func ensureSemanticLabelTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS semantic_labels (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_id INTEGER NOT NULL REFERENCES repos(id),
			label TEXT NOT NULL,
			aliases_json TEXT NOT NULL,
			sources_json TEXT NOT NULL,
			seen_count INTEGER NOT NULL,
			accepted_count INTEGER NOT NULL,
			confidence REAL NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL,
			UNIQUE(repo_id, label)
		);
		CREATE INDEX IF NOT EXISTS idx_semantic_labels_repo ON semantic_labels(repo_id, updated_at DESC);
		CREATE VIRTUAL TABLE IF NOT EXISTS semantic_label_fts USING fts5(
			repo_id UNINDEXED,
			label,
			aliases_json,
			sources_json,
			tokenize='unicode61'
		);
		CREATE TABLE IF NOT EXISTS semantic_label_links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_id INTEGER NOT NULL REFERENCES repos(id),
			label_id INTEGER NOT NULL REFERENCES semantic_labels(id),
			demux_proposal_id TEXT,
			revision_proposal_id TEXT,
			change_id INTEGER REFERENCES changes(id),
			stack_bookmark TEXT,
			source TEXT NOT NULL,
			score REAL NOT NULL,
			accepted INTEGER NOT NULL,
			evidence_json TEXT NOT NULL,
			created_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_semantic_label_links_repo ON semantic_label_links(repo_id, created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_semantic_label_links_change ON semantic_label_links(change_id);
		CREATE INDEX IF NOT EXISTS idx_semantic_label_links_proposal ON semantic_label_links(demux_proposal_id);
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

func ensureSessionContextsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS session_contexts (
			session_id TEXT PRIMARY KEY REFERENCES sessions(id) ON DELETE CASCADE,
			tool TEXT NOT NULL,
			model TEXT,
			format TEXT NOT NULL,
			content_json BLOB NOT NULL,
			captured_at INTEGER NOT NULL
		);
	`)
	return err
}

func ensureSessionEventAttributionsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS session_event_attributions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_id INTEGER NOT NULL REFERENCES repos(id),
			tool TEXT NOT NULL,
			session_id TEXT NOT NULL REFERENCES sessions(id),
			event_fingerprint TEXT NOT NULL,
			change_id INTEGER NOT NULL REFERENCES changes(id),
			stack_bookmark TEXT,
			attributed_via TEXT NOT NULL,
			confidence REAL,
			created_at INTEGER NOT NULL,
			UNIQUE(repo_id, tool, session_id, event_fingerprint)
		);
		CREATE INDEX IF NOT EXISTS idx_session_event_attributions_repo ON session_event_attributions(repo_id);
		CREATE INDEX IF NOT EXISTS idx_session_event_attributions_change ON session_event_attributions(change_id);
		CREATE INDEX IF NOT EXISTS idx_session_event_attributions_session ON session_event_attributions(session_id);
	`)
	return err
}

func ensureSessionIndexes(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_sessions_repo_root ON sessions(repo_root);
	`)
	return err
}

// commitSelfReportSessionID is the fossil left by the retired `tl commit`
// self-report. It has cwd=” and repo_root=NULL, so it can never match a repo
// and only ever contributed noise to the change_sessions table. Nothing writes
// it any more, so removing it on open is a one-way cleanup.
const commitSelfReportSessionID = "totality-commit-self-report"

func deleteCommitSelfReportSession(ctx context.Context, db *sql.DB) error {
	for _, stmt := range []string{
		`DELETE FROM change_session_provenance WHERE session_id = ?`,
		`DELETE FROM session_event_attributions WHERE session_id = ?`,
		`DELETE FROM change_sessions WHERE session_id = ?`,
		`DELETE FROM session_contexts WHERE session_id = ?`,
		`DELETE FROM sessions WHERE id = ?`,
	} {
		if _, err := db.ExecContext(ctx, stmt, commitSelfReportSessionID); err != nil {
			return err
		}
	}
	return nil
}

func repairDanglingSessionLinks(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `
		DELETE FROM change_session_provenance
		WHERE change_id NOT IN (SELECT id FROM changes)
			OR session_id NOT IN (SELECT id FROM sessions)
			OR NOT EXISTS (
				SELECT 1
				FROM change_sessions cs
				WHERE cs.change_id = change_session_provenance.change_id
					AND cs.session_id = change_session_provenance.session_id
			)
	`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `
		DELETE FROM change_sessions
		WHERE change_id NOT IN (SELECT id FROM changes)
			OR session_id NOT IN (SELECT id FROM sessions)
	`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `
		DELETE FROM session_event_attributions
		WHERE repo_id NOT IN (SELECT id FROM repos)
			OR change_id NOT IN (SELECT id FROM changes)
			OR session_id NOT IN (SELECT id FROM sessions)
	`); err != nil {
		return err
	}
	return nil
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
		_, err = db.ExecContext(ctx, `CREATE INDEX IF NOT EXISTS idx_change_bookmarks_name ON change_bookmarks(bookmark_name)`)
		return err
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
			CREATE INDEX IF NOT EXISTS idx_change_bookmarks_name ON change_bookmarks(bookmark_name);
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
		"source":             "TEXT",
		"process_name":       "TEXT",
		"parent_pid":         "INTEGER",
		"last_seen_at":       "INTEGER",
		"end_reason":         "TEXT",
		"repo_root":          "TEXT",
		"models_json":        "TEXT NOT NULL DEFAULT '[]'",
		"input_tokens":       "INTEGER NOT NULL DEFAULT 0",
		"output_tokens":      "INTEGER NOT NULL DEFAULT 0",
		"cache_read_tokens":  "INTEGER NOT NULL DEFAULT 0",
		"cache_write_tokens": "INTEGER NOT NULL DEFAULT 0",
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

// ensureCaptureSessionBindingColumns adds the columns that record which
// repository a session belongs to. They arrived after capture_sessions
// existed, so an older database reaches schema.sql without them.
func ensureCaptureSessionBindingColumns(ctx context.Context, db *sql.DB) error {
	if err := ensureColumn(ctx, db, "capture_sessions", "session_cwd", "TEXT"); err != nil {
		return err
	}
	return ensureColumn(ctx, db, "capture_sessions", "session_origin", "TEXT")
}
