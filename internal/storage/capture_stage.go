package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StagedExtract is one pending extract upload row.
type StagedExtract struct {
	ID            string
	RepoRoot      string
	RefRange      string
	PayloadJSON   []byte
	CreatedAt     int64
	RevisionID    string
	ContentHash   string
	ShareableAt   int64
	AcceptorName  string
	AcceptorEmail string
	AttestedAt    int64
}

// StagedSession is one pending session upload row.
type StagedSession struct {
	ID            string
	SessionID     string
	Tool          string
	PayloadJSON   []byte
	RawBlob       []byte
	SourcePath    string
	BadLines      int
	CreatedAt     int64
	RevisionID    string
	ContentHash   string
	ShareableAt   int64
	AcceptorName  string
	AcceptorEmail string
	AttestedAt    int64
}

// CaptureAttestation records who marked capture rows shareable.
type CaptureAttestation struct {
	Name  string
	Email string
}

// CaptureStager writes capture payloads to local SQLite.
type CaptureStager interface {
	StageExtract(ctx context.Context, row StagedExtract) error
	StageSession(ctx context.Context, row StagedSession) error
	PendingExtracts(ctx context.Context) ([]StagedExtract, error)
	PendingSessions(ctx context.Context) ([]StagedSession, error)
	ShareableExtracts(ctx context.Context) ([]StagedExtract, error)
	ShareableSessions(ctx context.Context) ([]StagedSession, error)
	MarkExtractShareable(ctx context.Context, revisionIDs []string, att CaptureAttestation) (int, error)
	MarkSessionShareable(ctx context.Context, revisionIDs []string, att CaptureAttestation) (int, error)
	MarkExtractUploaded(ctx context.Context, id string) error
	MarkSessionUploaded(ctx context.Context, id string) error
	RawSessions(ctx context.Context) ([]StagedSession, error)
	SetExtractUploadError(ctx context.Context, id, message string) error
	SetSessionUploadError(ctx context.Context, id, message string) error
}

// CaptureStage persists capture rows in ~/.gx/gx.db.
type CaptureStage struct {
	db *sql.DB
}

// OpenCaptureStager opens the default gx database with capture tables.
func OpenCaptureStager(ctx context.Context) (*CaptureStage, error) {
	db, err := Open(ctx)
	if err != nil {
		return nil, err
	}
	return &CaptureStage{db: db}, nil
}

// PayloadContentHash returns a stable sha256 hex digest for capture payloads.
func PayloadContentHash(payload []byte) string {
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

// CaptureRowID builds the idempotent row key for revision-scoped capture uploads.
func CaptureRowID(revisionID, contentHash string) string {
	revisionID = strings.TrimSpace(revisionID)
	contentHash = strings.TrimSpace(contentHash)
	if revisionID == "" || contentHash == "" {
		return ""
	}
	short := contentHash
	if len(short) > 24 {
		short = short[:24]
	}
	return revisionID + "-" + short
}

func (s *CaptureStage) StageExtract(ctx context.Context, row StagedExtract) error {
	created := row.CreatedAt
	if created == 0 {
		created = time.Now().UnixMilli()
	}
	if row.ContentHash == "" {
		row.ContentHash = PayloadContentHash(row.PayloadJSON)
	}
	if row.ID == "" {
		row.ID = CaptureRowID(row.RevisionID, row.ContentHash)
	}
	if row.ID == "" {
		row.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO capture_extracts (
			id, repo_root, ref_range, payload_json, created_at,
			revision_id, content_hash, shareable_at, acceptor_name, acceptor_email, attested_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			repo_root = excluded.repo_root,
			ref_range = excluded.ref_range,
			payload_json = excluded.payload_json,
			created_at = excluded.created_at,
			revision_id = excluded.revision_id,
			content_hash = excluded.content_hash
	`, row.ID, row.RepoRoot, row.RefRange, row.PayloadJSON, created,
		nullString(row.RevisionID), nullString(row.ContentHash),
		nullInt64(row.ShareableAt), nullString(row.AcceptorName), nullString(row.AcceptorEmail), nullInt64(row.AttestedAt))
	if err != nil {
		return fmt.Errorf("insert capture_extract: %w", err)
	}
	return nil
}

func (s *CaptureStage) StageSession(ctx context.Context, row StagedSession) error {
	created := row.CreatedAt
	if created == 0 {
		created = time.Now().UnixMilli()
	}
	if row.ContentHash == "" {
		row.ContentHash = PayloadContentHash(row.PayloadJSON)
	}
	if row.ID == "" {
		row.ID = CaptureRowID(row.RevisionID, row.ContentHash)
	}
	if row.ID == "" {
		row.ID = uuid.NewString()
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO capture_sessions (
			id, session_id, tool, payload_json, raw_blob, source_path, bad_lines, created_at,
			revision_id, content_hash, shareable_at, acceptor_name, acceptor_email, attested_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			session_id = excluded.session_id,
			tool = excluded.tool,
			payload_json = excluded.payload_json,
			raw_blob = COALESCE(excluded.raw_blob, capture_sessions.raw_blob),
			source_path = COALESCE(NULLIF(excluded.source_path, ''), capture_sessions.source_path),
			bad_lines = CASE WHEN excluded.bad_lines > 0 THEN excluded.bad_lines ELSE capture_sessions.bad_lines END,
			created_at = excluded.created_at,
			revision_id = excluded.revision_id,
			content_hash = excluded.content_hash
	`, row.ID, row.SessionID, row.Tool, row.PayloadJSON, nullBytes(row.RawBlob), nullString(row.SourcePath), nullInt(row.BadLines), created,
		nullString(row.RevisionID), nullString(row.ContentHash),
		nullInt64(row.ShareableAt), nullString(row.AcceptorName), nullString(row.AcceptorEmail), nullInt64(row.AttestedAt))
	if err != nil {
		return fmt.Errorf("insert capture_session: %w", err)
	}
	return nil
}

func (s *CaptureStage) RawSessions(ctx context.Context) ([]StagedSession, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, session_id, tool, payload_json, COALESCE(raw_blob, ''), COALESCE(source_path, ''), COALESCE(bad_lines, 0), created_at,
			COALESCE(revision_id, ''), COALESCE(content_hash, ''), COALESCE(shareable_at, 0),
			COALESCE(acceptor_name, ''), COALESCE(acceptor_email, ''), COALESCE(attested_at, 0)
		FROM capture_sessions
		WHERE raw_blob IS NOT NULL AND length(raw_blob) > 0
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSessionRows(rows)
}

func (s *CaptureStage) PendingExtracts(ctx context.Context) ([]StagedExtract, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_root, ref_range, payload_json, created_at,
			COALESCE(revision_id, ''), COALESCE(content_hash, ''), COALESCE(shareable_at, 0),
			COALESCE(acceptor_name, ''), COALESCE(acceptor_email, ''), COALESCE(attested_at, 0)
		FROM capture_extracts
		WHERE uploaded_at IS NULL
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanExtractRows(rows)
}

func (s *CaptureStage) PendingSessions(ctx context.Context) ([]StagedSession, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, session_id, tool, payload_json, COALESCE(raw_blob, ''), COALESCE(source_path, ''), COALESCE(bad_lines, 0), created_at,
			COALESCE(revision_id, ''), COALESCE(content_hash, ''), COALESCE(shareable_at, 0),
			COALESCE(acceptor_name, ''), COALESCE(acceptor_email, ''), COALESCE(attested_at, 0)
		FROM capture_sessions
		WHERE uploaded_at IS NULL
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSessionRows(rows)
}

func (s *CaptureStage) ShareableExtracts(ctx context.Context) ([]StagedExtract, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_root, ref_range, payload_json, created_at,
			COALESCE(revision_id, ''), COALESCE(content_hash, ''), COALESCE(shareable_at, 0),
			COALESCE(acceptor_name, ''), COALESCE(acceptor_email, ''), COALESCE(attested_at, 0)
		FROM capture_extracts
		WHERE shareable_at IS NOT NULL AND uploaded_at IS NULL
		ORDER BY shareable_at ASC, created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanExtractRows(rows)
}

func (s *CaptureStage) ShareableSessions(ctx context.Context) ([]StagedSession, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, session_id, tool, payload_json, COALESCE(raw_blob, ''), COALESCE(source_path, ''), COALESCE(bad_lines, 0), created_at,
			COALESCE(revision_id, ''), COALESCE(content_hash, ''), COALESCE(shareable_at, 0),
			COALESCE(acceptor_name, ''), COALESCE(acceptor_email, ''), COALESCE(attested_at, 0)
		FROM capture_sessions
		WHERE shareable_at IS NOT NULL AND uploaded_at IS NULL
		ORDER BY shareable_at ASC, created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSessionRows(rows)
}

func (s *CaptureStage) MarkExtractShareable(ctx context.Context, revisionIDs []string, att CaptureAttestation) (int, error) {
	return s.markShareable(ctx, "capture_extracts", revisionIDs, att)
}

func (s *CaptureStage) MarkSessionShareable(ctx context.Context, revisionIDs []string, att CaptureAttestation) (int, error) {
	return s.markShareable(ctx, "capture_sessions", revisionIDs, att)
}

func (s *CaptureStage) markShareable(ctx context.Context, table string, revisionIDs []string, att CaptureAttestation) (int, error) {
	revisionIDs = uniqueNonEmpty(revisionIDs)
	if len(revisionIDs) == 0 {
		return 0, nil
	}
	now := time.Now().UnixMilli()
	placeholders := strings.Repeat("?,", len(revisionIDs))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]any, 0, len(revisionIDs)+4)
	args = append(args, now, att.Name, att.Email, now)
	for _, id := range revisionIDs {
		args = append(args, id)
	}
	query := fmt.Sprintf(`
		UPDATE %s
		SET shareable_at = COALESCE(shareable_at, ?),
			acceptor_name = COALESCE(acceptor_name, ?),
			acceptor_email = COALESCE(acceptor_email, ?),
			attested_at = COALESCE(attested_at, ?)
		WHERE revision_id IN (%s) AND uploaded_at IS NULL
	`, table, placeholders)
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func (s *CaptureStage) MarkExtractUploaded(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_extracts
		SET uploaded_at = ?, upload_error = NULL
		WHERE id = ?
	`, time.Now().UnixMilli(), id)
	return err
}

func (s *CaptureStage) MarkSessionUploaded(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_sessions
		SET uploaded_at = ?, upload_error = NULL
		WHERE id = ?
	`, time.Now().UnixMilli(), id)
	return err
}

func (s *CaptureStage) SetExtractUploadError(ctx context.Context, id, message string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_extracts SET upload_error = ? WHERE id = ?
	`, message, id)
	return err
}

func (s *CaptureStage) SetSessionUploadError(ctx context.Context, id, message string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_sessions SET upload_error = ? WHERE id = ?
	`, message, id)
	return err
}

// CaptureBacklogCounts reports unstaged upload backlog sizes.
type CaptureBacklogCounts struct {
	Extracts int
	Sessions int
}

// PendingCaptureCounts returns pending extract/session row counts.
func PendingCaptureCounts(ctx context.Context) (CaptureBacklogCounts, error) {
	stager, err := OpenCaptureStager(ctx)
	if err != nil {
		return CaptureBacklogCounts{}, err
	}
	extracts, err := stager.ShareableExtracts(ctx)
	if err != nil {
		return CaptureBacklogCounts{}, err
	}
	sessions, err := stager.ShareableSessions(ctx)
	if err != nil {
		return CaptureBacklogCounts{}, err
	}
	return CaptureBacklogCounts{
		Extracts: len(extracts),
		Sessions: len(sessions),
	}, nil
}

func scanExtractRows(rows *sql.Rows) ([]StagedExtract, error) {
	var out []StagedExtract
	for rows.Next() {
		var row StagedExtract
		if err := rows.Scan(
			&row.ID, &row.RepoRoot, &row.RefRange, &row.PayloadJSON, &row.CreatedAt,
			&row.RevisionID, &row.ContentHash, &row.ShareableAt,
			&row.AcceptorName, &row.AcceptorEmail, &row.AttestedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func scanSessionRows(rows *sql.Rows) ([]StagedSession, error) {
	var out []StagedSession
	for rows.Next() {
		var row StagedSession
		if err := rows.Scan(
			&row.ID, &row.SessionID, &row.Tool, &row.PayloadJSON, &row.RawBlob, &row.SourcePath, &row.BadLines, &row.CreatedAt,
			&row.RevisionID, &row.ContentHash, &row.ShareableAt,
			&row.AcceptorName, &row.AcceptorEmail, &row.AttestedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func ensureCaptureStagingTables(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS capture_extracts (
			id TEXT PRIMARY KEY,
			repo_root TEXT NOT NULL,
			ref_range TEXT NOT NULL,
			payload_json BLOB NOT NULL,
			created_at INTEGER NOT NULL,
			revision_id TEXT,
			content_hash TEXT,
			shareable_at INTEGER,
			acceptor_name TEXT,
			acceptor_email TEXT,
			attested_at INTEGER,
			uploaded_at INTEGER,
			upload_error TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_capture_extracts_pending ON capture_extracts(uploaded_at, created_at);

		CREATE TABLE IF NOT EXISTS capture_sessions (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			tool TEXT NOT NULL,
			payload_json BLOB NOT NULL,
			raw_blob BLOB,
			source_path TEXT,
			bad_lines INTEGER,
			created_at INTEGER NOT NULL,
			revision_id TEXT,
			content_hash TEXT,
			shareable_at INTEGER,
			acceptor_name TEXT,
			acceptor_email TEXT,
			attested_at INTEGER,
			uploaded_at INTEGER,
			upload_error TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_capture_sessions_pending ON capture_sessions(uploaded_at, created_at);
	`)
	if err != nil {
		return err
	}
	columns := []struct {
		table string
		name  string
		typ   string
	}{
		{"capture_extracts", "revision_id", "TEXT"},
		{"capture_extracts", "content_hash", "TEXT"},
		{"capture_extracts", "shareable_at", "INTEGER"},
		{"capture_extracts", "acceptor_name", "TEXT"},
		{"capture_extracts", "acceptor_email", "TEXT"},
		{"capture_extracts", "attested_at", "INTEGER"},
		{"capture_sessions", "revision_id", "TEXT"},
		{"capture_sessions", "content_hash", "TEXT"},
		{"capture_sessions", "raw_blob", "BLOB"},
		{"capture_sessions", "source_path", "TEXT"},
		{"capture_sessions", "bad_lines", "INTEGER"},
		{"capture_sessions", "shareable_at", "INTEGER"},
		{"capture_sessions", "acceptor_name", "TEXT"},
		{"capture_sessions", "acceptor_email", "TEXT"},
		{"capture_sessions", "attested_at", "INTEGER"},
	}
	for _, column := range columns {
		if err := ensureColumn(ctx, db, column.table, column.name, column.typ); err != nil {
			return err
		}
	}
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_capture_extracts_shareable ON capture_extracts(shareable_at, uploaded_at, created_at);
		CREATE INDEX IF NOT EXISTS idx_capture_sessions_shareable ON capture_sessions(shareable_at, uploaded_at, created_at);
	`)
	return err
}

func uniqueNonEmpty(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
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

func nullInt(value int) any {
	if value == 0 {
		return nil
	}
	return value
}

func nullBytes(value []byte) any {
	if len(value) == 0 {
		return nil
	}
	return value
}

func nullString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func nullInt64(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}
