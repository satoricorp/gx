package vcs

import (
	"context"
	"database/sql"
	"testing"

	"github.com/satoricorp/lgtm/internal/capture/matcher"
	"github.com/satoricorp/lgtm/internal/storage"
)

func newSessionLinkStore(t *testing.T) (*storage.Store, *sql.DB, context.Context) {
	t.Helper()
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store, db, ctx
}

func seedAttachChange(t *testing.T, store *storage.Store, ctx context.Context, root, sha string) int64 {
	t.Helper()
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath: root, GitCommonDir: root + "/.git", Backend: "git", CreatedAt: 1, UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID: repoID, JJChangeID: "rev-" + sha, CurrentCommitID: sha,
		Description: "seeded", Status: "draft", FirstSeenAt: 1, UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	return changeID
}

func sessionCommand(t *testing.T, db *sql.DB, ctx context.Context, id string) string {
	t.Helper()
	var command string
	if err := db.QueryRowContext(ctx, `SELECT command FROM sessions WHERE id = ?`, id).Scan(&command); err != nil {
		t.Fatalf("read session %q: %v", id, err)
	}
	return command
}

func changeSessionCount(t *testing.T, db *sql.DB, ctx context.Context, query string, args ...any) int {
	t.Helper()
	var count int
	if err := db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		t.Fatalf("count change_sessions: %v", err)
	}
	return count
}

// TestAttachSessionsPreservesExistingSessionMetadata pins the push-time writer
// to filling blanks rather than overwriting.
//
// Cursor's state.vscdb transcripts keep the id the parser mints,
// "cursor-<composerID>", because the path-derived identity override in
// internal/capture/parsers only applies to JSONL sources. Databases written
// while the cursor ingest existed still hold rows under exactly that id with
// the composer title as the command, so an attach lands on those rows. An
// UpsertSession here replaced the title with a bare "cursor", irrecoverably:
// nothing rewrites command afterwards. The richer pre-existing row is seeded
// through UpsertObservedSession, the same writer the attach path uses.
func TestAttachSessionsPreservesExistingSessionMetadata(t *testing.T) {
	store, db, ctx := newSessionLinkStore(t)
	const sha = "1111111111111111111111111111111111111111"
	const root = "/repo"
	changeID := seedAttachChange(t, store, ctx, root, sha)

	ingestSource := "cursor"
	ingestRoot := "/repo"
	if err := store.UpsertObservedSession(ctx, storage.Session{
		ID:        "cursor-abc123",
		CreatedAt: 500,
		Command:   "cursor: Refactor the auth middleware",
		Cwd:       "/repo",
		TLVersion: "test",
		Source:    &ingestSource,
		RepoRoot:  &ingestRoot,
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}

	attached, err := NewService().AttachSessionsFromHunkLinks(ctx, RepoInfo{
		RootPath: root, GitCommonDir: root + "/.git", Backend: "git",
	}, []matcher.HunkLink{{
		HunkID:     sha + ":auth.go:1-4",
		SessionID:  "cursor-abc123",
		Tier:       matcher.TierExact,
		Authorship: matcher.AuthorshipAgent,
		Tool:       "cursor",
	}})
	if err != nil {
		t.Fatalf("AttachSessionsFromHunkLinks() error = %v", err)
	}
	if attached != 1 {
		t.Fatalf("attached = %d, want 1", attached)
	}

	if got := sessionCommand(t, db, ctx, "cursor-abc123"); got != "cursor: Refactor the auth middleware" {
		t.Fatalf("session command = %q, want the composer title preserved", got)
	}

	links := changeSessionCount(t, db, ctx,
		`SELECT COUNT(1) FROM change_sessions WHERE change_id = ? AND session_id = ?`,
		changeID, "cursor-abc123")
	if links != 1 {
		t.Fatalf("change_sessions rows = %d, want 1", links)
	}
}

// TestAttachSessionsFillsBlanksOnNewSessions is the other half: a session the
// push observes for the first time must still get its command, cwd and
// repo_root, or the bundle renders `command=?`.
func TestAttachSessionsFillsBlanksOnNewSessions(t *testing.T) {
	store, db, ctx := newSessionLinkStore(t)
	const sha = "2222222222222222222222222222222222222222"
	const root = "/repo"
	seedAttachChange(t, store, ctx, root, sha)

	repo := RepoInfo{RootPath: root, GitCommonDir: root + "/.git", Backend: "git"}
	links := []matcher.HunkLink{{
		HunkID:     sha + ":main.go:1-4",
		SessionID:  "conversation-uuid",
		Tier:       matcher.TierExact,
		Authorship: matcher.AuthorshipAgent,
		Tool:       "claude",
	}}
	if _, err := NewService().AttachSessionsFromHunkLinks(ctx, repo, links); err != nil {
		t.Fatalf("AttachSessionsFromHunkLinks() error = %v", err)
	}
	if got := sessionCommand(t, db, ctx, "conversation-uuid"); got != "claude" {
		t.Fatalf("session command = %q, want claude", got)
	}

	// A second push must not duplicate the link.
	if _, err := NewService().AttachSessionsFromHunkLinks(ctx, repo, links); err != nil {
		t.Fatalf("second AttachSessionsFromHunkLinks() error = %v", err)
	}
	count := changeSessionCount(t, db, ctx,
		`SELECT COUNT(1) FROM change_sessions WHERE session_id = ?`, "conversation-uuid")
	if count != 1 {
		t.Fatalf("change_sessions rows = %d, want 1 after two pushes", count)
	}
}
