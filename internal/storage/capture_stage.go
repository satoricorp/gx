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
//
// A row is a pointer to a transcript source file, not a copy of it: SourcePath
// locates the transcript and ContentHash/SourceBytes/SourceMTime fingerprint
// the bytes seen at staging time. RawBlob is only populated on legacy rows
// written before the pointer model; it is read for fallback but never written.
type StagedSession struct {
	ID            string
	SessionID     string
	Tool          string
	PayloadJSON   []byte
	RawBlob       []byte
	SourcePath    string
	SourceBytes   int64
	SourceMTime   int64
	BadLines      int
	CreatedAt     int64
	RevisionID    string
	ContentHash   string
	ShareableAt   int64
	AcceptorName  string
	AcceptorEmail string
	AttestedAt    int64

	// SessionCwd is the working directory the agent session ran in, and
	// SessionOrigin is the normalized git origin that directory resolved to.
	// Together they answer which repository a session belongs to without
	// reference to whether its work was ever committed or matched.
	//
	// SessionOrigin is stored normalized, never raw: a remote URL can carry a
	// token, and the normalized form has credentials stripped.
	SessionCwd    string
	SessionOrigin string
}

// MaxUploadAttempts bounds how many times one staged row is re-sent.
//
// Uploads run in a detached `lgtm capture sync` whose output goes to /dev/null,
// so a row that fails deterministically — a payload the server rejects with
// 400, say — used to be re-sent on every push forever, failing silently every
// time. Attempts are counted so a permanently broken row stops consuming
// bandwidth and starts being reported instead. Re-staging with new content
// resets the counter, because new bytes deserve a fresh verdict.
const MaxUploadAttempts = 5

// CaptureUploadFailures counts staged rows whose upload failed, and carries one
// example message. It exists so the failures a detached sync writes into the
// database are readable by something a human actually looks at.
type CaptureUploadFailures struct {
	Extracts       int
	Sessions       int
	ExhaustedRows  int
	LastError      string
	LastErrorTable string
}

// Total reports how many rows are currently in a failed state.
func (f CaptureUploadFailures) Total() int {
	return f.Extracts + f.Sessions
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
	MarkExtractsShareableByID(ctx context.Context, ids []string, att CaptureAttestation) (int, error)
	MarkSessionsShareableByID(ctx context.Context, ids []string, att CaptureAttestation) (int, error)
	MarkSessionsShareableForSourcesOf(ctx context.Context, ids []string, att CaptureAttestation) (int, error)
	SessionOriginsByID(ctx context.Context, ids []string) (map[string]string, error)
	MarkExtractUploaded(ctx context.Context, id string) error
	MarkSessionUploaded(ctx context.Context, id string) error
	RawSessions(ctx context.Context) ([]StagedSession, error)
	SetExtractUploadError(ctx context.Context, id, message string) error
	SetSessionUploadError(ctx context.Context, id, message string) error
	SetSessionSourceFingerprint(ctx context.Context, id, contentHash string, sourceBytes, sourceMTime int64) error
}

// CaptureStage persists capture rows in ~/.lgtm/lgtm.db.
type CaptureStage struct {
	db *sql.DB
}

// OpenCaptureStager opens the default lgtm database with capture tables.
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

// SessionSourceRowID builds the stable capture_sessions key for one transcript
// source: one row per (tool, source file, per-source session identity).
// Re-staging the same source on a later push hits ON CONFLICT and updates the
// row in place instead of accumulating a copy per push or per revision.
// sessionID is part of the key because one source file can hold several
// sessions (Cursor's state.vscdb); for JSONL sources it is path-derived and
// the key degenerates to one row per file.
func SessionSourceRowID(tool, sourcePath, sessionID string) string {
	tool = strings.TrimSpace(tool)
	sourcePath = strings.TrimSpace(sourcePath)
	sessionID = strings.TrimSpace(sessionID)
	if sourcePath == "" && sessionID == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(tool + "\x00" + sourcePath + "\x00" + sessionID))
	return "src-" + hex.EncodeToString(sum[:12])
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
	// An unchanged content hash means the exact bytes already uploaded, so the
	// upload state survives the restage. A changed hash means the source grew
	// (or was rewritten) and the row must upload again.
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO capture_sessions (
			id, session_id, tool, payload_json, raw_blob, source_path, source_bytes, source_mtime, bad_lines, created_at,
			revision_id, content_hash, shareable_at, acceptor_name, acceptor_email, attested_at,
			session_cwd, session_origin
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			session_id = excluded.session_id,
			tool = excluded.tool,
			payload_json = excluded.payload_json,
			raw_blob = COALESCE(excluded.raw_blob, capture_sessions.raw_blob),
			source_path = COALESCE(NULLIF(excluded.source_path, ''), capture_sessions.source_path),
			source_bytes = COALESCE(excluded.source_bytes, capture_sessions.source_bytes),
			source_mtime = COALESCE(excluded.source_mtime, capture_sessions.source_mtime),
			bad_lines = CASE WHEN excluded.bad_lines > 0 THEN excluded.bad_lines ELSE capture_sessions.bad_lines END,
			created_at = excluded.created_at,
			revision_id = excluded.revision_id,
			content_hash = excluded.content_hash,
			uploaded_at = CASE WHEN excluded.content_hash IS capture_sessions.content_hash THEN capture_sessions.uploaded_at ELSE NULL END,
			upload_error = CASE WHEN excluded.content_hash IS capture_sessions.content_hash THEN capture_sessions.upload_error ELSE NULL END,
			upload_attempts = CASE WHEN excluded.content_hash IS capture_sessions.content_hash THEN COALESCE(capture_sessions.upload_attempts, 0) ELSE 0 END,
			session_cwd = COALESCE(NULLIF(excluded.session_cwd, ''), capture_sessions.session_cwd),
			session_origin = COALESCE(NULLIF(excluded.session_origin, ''), capture_sessions.session_origin)
	`, row.ID, row.SessionID, row.Tool, row.PayloadJSON, nullBytes(row.RawBlob), nullString(row.SourcePath), nullInt64(row.SourceBytes), nullInt64(row.SourceMTime), nullInt(row.BadLines), created,
		nullString(row.RevisionID), nullString(row.ContentHash),
		nullInt64(row.ShareableAt), nullString(row.AcceptorName), nullString(row.AcceptorEmail), nullInt64(row.AttestedAt),
		nullString(row.SessionCwd), nullString(row.SessionOrigin))
	if err != nil {
		return fmt.Errorf("insert capture_session: %w", err)
	}
	return nil
}

func (s *CaptureStage) RawSessions(ctx context.Context) ([]StagedSession, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+sessionColumns+`
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
		SELECT `+sessionColumns+`
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

// ShareableExtracts lists extract rows this machine may still upload. Rows that
// have already burned MaxUploadAttempts are left out: they are not going to
// start working, and re-sending them hides the working rows behind noise.
func (s *CaptureStage) ShareableExtracts(ctx context.Context) ([]StagedExtract, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, repo_root, ref_range, payload_json, created_at,
			COALESCE(revision_id, ''), COALESCE(content_hash, ''), COALESCE(shareable_at, 0),
			COALESCE(acceptor_name, ''), COALESCE(acceptor_email, ''), COALESCE(attested_at, 0)
		FROM capture_extracts
		WHERE shareable_at IS NOT NULL AND uploaded_at IS NULL
			AND COALESCE(upload_attempts, 0) < ?
		ORDER BY shareable_at ASC, created_at ASC
	`, MaxUploadAttempts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanExtractRows(rows)
}

// ShareableSessions lists session rows this machine may still upload, subject to
// the same attempt ceiling as ShareableExtracts.
func (s *CaptureStage) ShareableSessions(ctx context.Context) ([]StagedSession, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT `+sessionColumns+`
		FROM capture_sessions
		WHERE shareable_at IS NOT NULL AND uploaded_at IS NULL
			AND COALESCE(upload_attempts, 0) < ?
		ORDER BY shareable_at ASC, created_at ASC
	`, MaxUploadAttempts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSessionRows(rows)
}

// SessionOriginsByID returns the repository each staged session is bound to,
// keyed by row id. Rows with no binding are present with an empty origin, so a
// caller can tell "bound elsewhere" from "not bound at all" — those mean
// different things to a gate.
func (s *CaptureStage) SessionOriginsByID(ctx context.Context, ids []string) (map[string]string, error) {
	out := map[string]string{}
	if len(ids) == 0 {
		return out, nil
	}
	list, args := placeholdersForStrings(ids)
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, COALESCE(session_origin, '') FROM capture_sessions WHERE id IN (`+list+`)`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, origin string
		if err := rows.Scan(&id, &origin); err != nil {
			return nil, err
		}
		out[id] = origin
	}
	return out, rows.Err()
}

// UploadFailures reports the staged rows whose last upload attempt failed.
//
// Nothing else reads upload_error. The sync that writes it is detached with
// both its streams pointed at /dev/null, so without this the only record of a
// failing upload was a column no code path ever looked at — which is how 32
// rows could 400 on every push for weeks while every push printed a clean
// capture line.
func (s *CaptureStage) UploadFailures(ctx context.Context) (CaptureUploadFailures, error) {
	var out CaptureUploadFailures
	for _, table := range []string{"capture_extracts", "capture_sessions"} {
		var count int
		var exhausted int
		var lastError sql.NullString
		query := fmt.Sprintf(`
			SELECT COUNT(*),
				SUM(CASE WHEN COALESCE(upload_attempts, 0) >= ? THEN 1 ELSE 0 END),
				MAX(upload_error)
			FROM %s
			WHERE upload_error IS NOT NULL AND uploaded_at IS NULL
		`, table)
		var exhaustedNull sql.NullInt64
		if err := s.db.QueryRowContext(ctx, query, MaxUploadAttempts).Scan(&count, &exhaustedNull, &lastError); err != nil {
			return out, err
		}
		exhausted = int(exhaustedNull.Int64)
		if table == "capture_extracts" {
			out.Extracts = count
		} else {
			out.Sessions = count
		}
		out.ExhaustedRows += exhausted
		if count > 0 && lastError.Valid && out.LastError == "" {
			out.LastError = lastError.String
			out.LastErrorTable = table
		}
	}
	return out, nil
}

// CaptureUploadFailureCounts opens the default database and reports the upload
// backlog's failure state.
func CaptureUploadFailureCounts(ctx context.Context) (CaptureUploadFailures, error) {
	stager, err := OpenCaptureStager(ctx)
	if err != nil {
		return CaptureUploadFailures{}, err
	}
	return stager.UploadFailures(ctx)
}

func (s *CaptureStage) MarkExtractShareable(ctx context.Context, revisionIDs []string, att CaptureAttestation) (int, error) {
	return s.markShareable(ctx, "capture_extracts", "revision_id", revisionIDs, att)
}

func (s *CaptureStage) MarkSessionShareable(ctx context.Context, revisionIDs []string, att CaptureAttestation) (int, error) {
	return s.markShareable(ctx, "capture_sessions", "revision_id", revisionIDs, att)
}

// MarkExtractsShareableByID marks the exact rows a capture run staged.
//
// Revision matching reaches a row only when that row carried a revision id from
// the pushed range, which makes it useless for the rows this very push just
// staged whenever revision resolution returns an empty set for the range. The
// staging row ids are known exactly, so they are the reliable gate; revision
// matching stays as a second pass for the backlog.
func (s *CaptureStage) MarkExtractsShareableByID(ctx context.Context, ids []string, att CaptureAttestation) (int, error) {
	return s.markShareable(ctx, "capture_extracts", "id", ids, att)
}

// MarkSessionsShareableByID marks the exact session rows a capture run staged.
func (s *CaptureStage) MarkSessionsShareableByID(ctx context.Context, ids []string, att CaptureAttestation) (int, error) {
	return s.markShareable(ctx, "capture_sessions", "id", ids, att)
}

// MarkSessionsShareableForSourcesOf marks the older rows that describe the same
// transcripts as ids.
//
// The row key changed from "<revisionID>-<contentHash>" to a source-derived
// "src-<hash>", so re-staging a transcript now writes a brand-new row and
// leaves its predecessor behind: same tool, same session, same file on disk,
// unreachable by id and — because its revision is already in pushed history —
// unreachable by revision too. Those predecessors would sit unshareable
// forever. They are the same evidence the caller is attesting to right now, not
// new evidence, so they inherit the same attestation. A source path is required
// to match, so rows without one (Cursor's vscdb) can never be swept in by an
// empty-string coincidence.
func (s *CaptureStage) MarkSessionsShareableForSourcesOf(ctx context.Context, ids []string, att CaptureAttestation) (int, error) {
	ids = uniqueNonEmpty(ids)
	if len(ids) == 0 {
		return 0, nil
	}
	now := time.Now().UnixMilli()
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]any, 0, len(ids)+4)
	args = append(args, now, att.Name, att.Email, now)
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf(`
		UPDATE capture_sessions
		SET shareable_at = COALESCE(shareable_at, ?),
			acceptor_name = COALESCE(acceptor_name, ?),
			acceptor_email = COALESCE(acceptor_email, ?),
			attested_at = COALESCE(attested_at, ?)
		WHERE uploaded_at IS NULL
			AND shareable_at IS NULL
			AND COALESCE(source_path, '') <> ''
			AND EXISTS (
				SELECT 1 FROM capture_sessions staged
				WHERE staged.id IN (%s)
					AND staged.tool = capture_sessions.tool
					AND staged.session_id = capture_sessions.session_id
					AND COALESCE(staged.source_path, '') = COALESCE(capture_sessions.source_path, '')
			)
	`, placeholders)
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

// StaleStagedRows counts rows that were staged before cutoff and never became
// shareable. Nothing will ever make them shareable — their revisions are long
// since pushed and no future run addresses them by id — so they are pure
// residue, and residue nobody can see is residue nobody ever clears.
func (s *CaptureStage) StaleStagedRows(ctx context.Context, cutoffMS int64) (int, error) {
	total := 0
	for _, table := range []string{"capture_extracts", "capture_sessions"} {
		var count int
		query := fmt.Sprintf(`
			SELECT COUNT(*) FROM %s
			WHERE shareable_at IS NULL AND uploaded_at IS NULL AND created_at < ?
		`, table)
		if err := s.db.QueryRowContext(ctx, query, cutoffMS).Scan(&count); err != nil {
			return total, err
		}
		total += count
	}
	return total, nil
}

// markShareable stamps attestation on rows selected by one column. Rows that
// are already shareable are excluded so the count is of newly marked rows and
// two passes over the same row (by id, then by revision) cannot double-count.
func (s *CaptureStage) markShareable(ctx context.Context, table, column string, values []string, att CaptureAttestation) (int, error) {
	values = uniqueNonEmpty(values)
	if len(values) == 0 {
		return 0, nil
	}
	now := time.Now().UnixMilli()
	placeholders := strings.Repeat("?,", len(values))
	placeholders = strings.TrimSuffix(placeholders, ",")
	args := make([]any, 0, len(values)+4)
	args = append(args, now, att.Name, att.Email, now)
	for _, id := range values {
		args = append(args, id)
	}
	query := fmt.Sprintf(`
		UPDATE %s
		SET shareable_at = COALESCE(shareable_at, ?),
			acceptor_name = COALESCE(acceptor_name, ?),
			acceptor_email = COALESCE(acceptor_email, ?),
			attested_at = COALESCE(attested_at, ?)
		WHERE %s IN (%s) AND uploaded_at IS NULL AND shareable_at IS NULL
	`, table, column, placeholders)
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
		SET uploaded_at = ?, upload_error = NULL, upload_attempts = 0
		WHERE id = ?
	`, time.Now().UnixMilli(), id)
	return err
}

func (s *CaptureStage) MarkSessionUploaded(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_sessions
		SET uploaded_at = ?, upload_error = NULL, upload_attempts = 0
		WHERE id = ?
	`, time.Now().UnixMilli(), id)
	return err
}

// SetExtractUploadError records a failed attempt and counts it, so a row that
// can never succeed eventually stops being retried instead of failing silently
// on every push for the rest of time.
func (s *CaptureStage) SetExtractUploadError(ctx context.Context, id, message string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_extracts
		SET upload_error = ?, upload_attempts = COALESCE(upload_attempts, 0) + 1
		WHERE id = ?
	`, message, id)
	return err
}

// SetSessionUploadError records a failed attempt and counts it.
func (s *CaptureStage) SetSessionUploadError(ctx context.Context, id, message string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_sessions
		SET upload_error = ?, upload_attempts = COALESCE(upload_attempts, 0) + 1
		WHERE id = ?
	`, message, id)
	return err
}

// SetSessionSourceFingerprint records the fingerprint of the source bytes that
// were actually read at upload time, which may be newer than what staging saw:
// transcripts grow while their session is still running.
func (s *CaptureStage) SetSessionSourceFingerprint(ctx context.Context, id, contentHash string, sourceBytes, sourceMTime int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE capture_sessions
		SET content_hash = ?, source_bytes = ?, source_mtime = ?
		WHERE id = ?
	`, nullString(contentHash), nullInt64(sourceBytes), nullInt64(sourceMTime), id)
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

// sessionColumns is the shared capture_sessions projection; keep in lockstep
// with scanSessionRows.
const sessionColumns = `id, session_id, tool, payload_json, COALESCE(raw_blob, ''), COALESCE(source_path, ''),
		COALESCE(source_bytes, 0), COALESCE(source_mtime, 0), COALESCE(bad_lines, 0), created_at,
		COALESCE(revision_id, ''), COALESCE(content_hash, ''), COALESCE(shareable_at, 0),
		COALESCE(acceptor_name, ''), COALESCE(acceptor_email, ''), COALESCE(attested_at, 0),
		COALESCE(session_cwd, ''), COALESCE(session_origin, '')`

func scanSessionRows(rows *sql.Rows) ([]StagedSession, error) {
	var out []StagedSession
	for rows.Next() {
		var row StagedSession
		if err := rows.Scan(
			&row.ID, &row.SessionID, &row.Tool, &row.PayloadJSON, &row.RawBlob, &row.SourcePath,
			&row.SourceBytes, &row.SourceMTime, &row.BadLines, &row.CreatedAt,
			&row.RevisionID, &row.ContentHash, &row.ShareableAt,
			&row.AcceptorName, &row.AcceptorEmail, &row.AttestedAt,
			&row.SessionCwd, &row.SessionOrigin,
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
			upload_error TEXT,
			upload_attempts INTEGER
		);
		CREATE INDEX IF NOT EXISTS idx_capture_extracts_pending ON capture_extracts(uploaded_at, created_at);

		CREATE TABLE IF NOT EXISTS capture_sessions (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			tool TEXT NOT NULL,
			payload_json BLOB NOT NULL,
			raw_blob BLOB,
			source_path TEXT,
			source_bytes INTEGER,
			source_mtime INTEGER,
			bad_lines INTEGER,
			created_at INTEGER NOT NULL,
			revision_id TEXT,
			content_hash TEXT,
			shareable_at INTEGER,
			acceptor_name TEXT,
			acceptor_email TEXT,
			attested_at INTEGER,
			uploaded_at INTEGER,
			upload_error TEXT,
			upload_attempts INTEGER
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
		{"capture_extracts", "upload_attempts", "INTEGER"},
		{"capture_sessions", "revision_id", "TEXT"},
		{"capture_sessions", "content_hash", "TEXT"},
		{"capture_sessions", "raw_blob", "BLOB"},
		{"capture_sessions", "source_path", "TEXT"},
		{"capture_sessions", "source_bytes", "INTEGER"},
		{"capture_sessions", "source_mtime", "INTEGER"},
		{"capture_sessions", "bad_lines", "INTEGER"},
		{"capture_sessions", "shareable_at", "INTEGER"},
		{"capture_sessions", "acceptor_name", "TEXT"},
		{"capture_sessions", "acceptor_email", "TEXT"},
		{"capture_sessions", "attested_at", "INTEGER"},
		{"capture_sessions", "upload_attempts", "INTEGER"},
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
