package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Store struct {
	db          *sql.DB
	requestStmt *sql.Stmt
	respStmt    *sql.Stmt
}

type sqlExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func NewStore(ctx context.Context, db *sql.DB) (*Store, error) {
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
		requestStmt: requestStmt,
		respStmt:    respStmt,
	}, nil
}

func (s *Store) Close() error {
	var firstErr error
	for _, stmt := range []*sql.Stmt{s.requestStmt, s.respStmt} {
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

// WriteSession is deliberately absent.
//
// It was a bare INSERT into `sessions` with no production callers and fourteen
// test call sites, and it is how the session bug survived a green suite: every
// session test seeded rows through it while UpsertObservedSession — the only
// writer of that table production actually runs — could never bootstrap its
// first row. Tests now seed through UpsertObservedSession (transcript
// observations, the push path).
//
// TestSeedersAreProductionWriters in internal/storage/storagetest fails on any
// exported storage method that has test callers and no production ones, so
// reintroducing this shape is a test failure rather than a silent regression.

// The session_event_attributions ledger has no Go API any more. It existed only
// so the commit-time matcher would not credit the same transcript event to two
// changes: that matcher scored events by recency, so nothing else stopped a
// second commit from re-consuming an event the first had already claimed. The
// push-time writer (vcs.AttachSessionsFromHunkLinks) derives links from hunk
// content instead, so a session lands on exactly the changes whose hunks it
// authored — and change_sessions' UNIQUE(change_id, session_id) makes repeat
// pushes idempotent without a ledger. The table and its migration stay so
// existing databases keep opening and the session-deletion cleanups in db.go
// keep working.

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
