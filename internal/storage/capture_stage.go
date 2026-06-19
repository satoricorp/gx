package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// StagedExtract is one pending extract upload row.
type StagedExtract struct {
	ID          string
	RepoRoot    string
	RefRange    string
	PayloadJSON []byte
	CreatedAt   int64
}

// StagedSession is one pending session upload row.
type StagedSession struct {
	ID          string
	SessionID   string
	Tool        string
	PayloadJSON []byte
	CreatedAt   int64
}

// CaptureStager writes capture payloads to local SQLite.
type CaptureStager interface {
	StageExtract(ctx context.Context, row StagedExtract) error
	StageSession(ctx context.Context, row StagedSession) error
	PendingExtracts(ctx context.Context) ([]StagedExtract, error)
	PendingSessions(ctx context.Context) ([]StagedSession, error)
	MarkExtractUploaded(ctx context.Context, id string) error
	MarkSessionUploaded(ctx context.Context, id string) error
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

func (s *CaptureStage) StageExtract(ctx context.Context, row StagedExtract) error {
	created := row.CreatedAt
	if created == 0 {
		created = time.Now().UnixMilli()
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO capture_extracts (id, repo_root, ref_range, payload_json, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			repo_root = excluded.repo_root,
			ref_range = excluded.ref_range,
			payload_json = excluded.payload_json,
			created_at = excluded.created_at
	`, row.ID, row.RepoRoot, row.RefRange, row.PayloadJSON, created)
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
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO capture_sessions (id, session_id, tool, payload_json, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			session_id = excluded.session_id,
			tool = excluded.tool,
			payload_json = excluded.payload_json,
			created_at = excluded.created_at
	`, row.ID, row.SessionID, row.Tool, row.PayloadJSON, created)
	if err != nil {
		return fmt.Errorf("insert capture_session: %w", err)
	}
	return nil
}

func (s *CaptureStage) PendingExtracts(ctx context.Context) ([]StagedExtract, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_root, ref_range, payload_json, created_at
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
		SELECT id, session_id, tool, payload_json, created_at
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
	extracts, err := stager.PendingExtracts(ctx)
	if err != nil {
		return CaptureBacklogCounts{}, err
	}
	sessions, err := stager.PendingSessions(ctx)
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
		if err := rows.Scan(&row.ID, &row.RepoRoot, &row.RefRange, &row.PayloadJSON, &row.CreatedAt); err != nil {
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
		if err := rows.Scan(&row.ID, &row.SessionID, &row.Tool, &row.PayloadJSON, &row.CreatedAt); err != nil {
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
			uploaded_at INTEGER,
			upload_error TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_capture_extracts_pending ON capture_extracts(uploaded_at, created_at);

		CREATE TABLE IF NOT EXISTS capture_sessions (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			tool TEXT NOT NULL,
			payload_json BLOB NOT NULL,
			created_at INTEGER NOT NULL,
			uploaded_at INTEGER,
			upload_error TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_capture_sessions_pending ON capture_sessions(uploaded_at, created_at);
	`)
	return err
}
