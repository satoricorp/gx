package provenance

import (
	"context"
	"reflect"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestParseSessionIDsDedupesAndSplitsCommonSeparators(t *testing.T) {
	got := ParseSessionIDs(" one,two  two\nthree\t")
	want := []string{"one", "two", "three"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseSessionIDs() = %#v, want %#v", got, want)
	}
}

func TestResolveIgnoresMissingEnvSessionIDs(t *testing.T) {
	t.Setenv("GX_SESSION_ID", "missing-session")
	t.Setenv("GX_SESSION_IDS", "")
	store := newTestStore(t)
	if err := store.UpsertObservedSession(context.Background(), storage.Session{
		ID:        "cursor-session",
		CreatedAt: 1,
		Command:   "cursor",
		Cwd:       "/repo",
		GXVersion: "test",
		RepoRoot:  strPtr("/repo"),
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}

	got, err := Resolve(context.Background(), store, "/repo")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got.Status != StatusRepoLocal || !reflect.DeepEqual(got.SessionIDs, []string{"cursor-session"}) {
		t.Fatalf("Resolve() = %#v, want repo-local fallback", got)
	}
}

func TestAttachPreferredFallsBackToRepoLocal(t *testing.T) {
	store := newTestStore(t)
	if err := store.UpsertObservedSession(context.Background(), storage.Session{
		ID:        "cursor-session",
		CreatedAt: 1,
		Command:   "cursor",
		Cwd:       "/repo",
		GXVersion: "test",
		RepoRoot:  strPtr("/repo"),
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      "change",
		CurrentCommitID: "commit",
		Description:     "change",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}

	got, err := AttachPreferred(context.Background(), store, "/repo", changeID, 3, []string{"stale-session"})
	if err != nil {
		t.Fatalf("AttachPreferred() error = %v", err)
	}
	if got.Status != StatusRepoLocal || !reflect.DeepEqual(got.SessionIDs, []string{"cursor-session"}) {
		t.Fatalf("AttachPreferred() = %#v, want repo-local fallback", got)
	}
}

func TestAttachSkipsMissingSessionsNeverFails(t *testing.T) {
	store := newTestStore(t)
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      "change",
		CurrentCommitID: "commit",
		Description:     "change",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}

	got, err := AttachPreferred(context.Background(), store, "/repo", changeID, 3, []string{"missing-one", "missing-two"})
	if err != nil {
		t.Fatalf("AttachPreferred() error = %v", err)
	}
	if got.Status != StatusAbsent || len(got.SessionIDs) != 0 {
		t.Fatalf("AttachPreferred() = %#v, want absent with no sessions", got)
	}
}

func TestResolvePrefersExplicitSessionIDs(t *testing.T) {
	t.Setenv("GX_SESSION_ID", "explicit-one")
	t.Setenv("GX_SESSION_IDS", "")
	store := newTestStore(t)
	if err := store.UpsertObservedSession(context.Background(), storage.Session{
		ID:        "explicit-one",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       "/repo",
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}

	got, err := Resolve(context.Background(), store, "/repo")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got.Status != StatusExplicit || !reflect.DeepEqual(got.SessionIDs, []string{"explicit-one"}) {
		t.Fatalf("Resolve() = %#v, want explicit session", got)
	}
}

func TestResolveReturnsRepoLocalSessionWhenNoExplicitEnv(t *testing.T) {
	t.Setenv("GX_SESSION_ID", "")
	t.Setenv("GX_SESSION_IDS", "")
	store := newTestStore(t)
	if err := store.UpsertObservedSession(context.Background(), storage.Session{
		ID:        "cursor-session",
		CreatedAt: 1,
		Command:   "cursor",
		Cwd:       "/repo",
		GXVersion: "test",
		RepoRoot:  strPtr("/repo"),
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}

	got, err := Resolve(context.Background(), store, "/repo")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if got.Status != StatusRepoLocal || !reflect.DeepEqual(got.SessionIDs, []string{"cursor-session"}) {
		t.Fatalf("Resolve() = %#v, want repo-local session", got)
	}
}

func TestAttachWritesResolvedSessions(t *testing.T) {
	t.Setenv("GX_SESSION_IDS", "session-one,session-two")
	store := newTestStore(t)
	if err := store.UpsertObservedSession(context.Background(), storage.Session{
		ID:        "session-one",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       "/repo",
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertObservedSession(session-one) error = %v", err)
	}
	if err := store.UpsertObservedSession(context.Background(), storage.Session{
		ID:        "session-two",
		CreatedAt: 2,
		Command:   "codex",
		Cwd:       "/repo",
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertObservedSession(session-two) error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      "change",
		CurrentCommitID: "commit",
		Description:     "change",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}

	got, err := Attach(context.Background(), store, "/repo", changeID, 3)
	if err != nil {
		t.Fatalf("Attach() error = %v", err)
	}
	if got.Status != StatusExplicit || !reflect.DeepEqual(got.SessionIDs, []string{"session-one", "session-two"}) {
		t.Fatalf("Attach() = %#v, want explicit sessions", got)
	}
}

func TestWithExplicitSessionEnvRestoresPreviousEnvironment(t *testing.T) {
	t.Setenv("GX_SESSION_ID", "old-one")
	t.Setenv("GX_SESSION_IDS", "old-many")

	got, err := WithExplicitSessionEnv([]string{"new-one", "new-two"}, func() ([]string, error) {
		return ExplicitSessionIDsFromEnv(), nil
	})
	if err != nil {
		t.Fatalf("WithExplicitSessionEnv() error = %v", err)
	}
	if !reflect.DeepEqual(got, []string{"new-one", "new-two"}) {
		t.Fatalf("WithExplicitSessionEnv() inner ids = %#v", got)
	}
	if got := ExplicitSessionIDsFromEnv(); !reflect.DeepEqual(got, []string{"old-many"}) {
		t.Fatalf("WithExplicitSessionEnv() restored ids = %#v, want old-many", got)
	}
}

// Sessions in this file are seeded with store.UpsertObservedSession, the writer
// vcs.AttachSessionsFromHunkLinks calls on the push path. store.WriteSession —
// what these tests used to call — has no production callers at all, so seeding
// through it proved provenance behaviour against `sessions` rows that no gx
// install could actually contain. See TestSeedersAreProductionWriters in
// internal/storage/storagetest for the rule and its ratchet.
func newTestStore(t *testing.T) *storage.Store {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })
	return store
}

func strPtr(value string) *string {
	return &value
}
