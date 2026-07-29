// Package storagetest builds tl databases that look like the ones on real
// machines instead of the pristine one a fresh migration produces.
//
// Why this exists. Two production defects shipped past a fully green suite for
// the same reason: every test did `t.Setenv("TOTALITY_HOME", t.TempDir())` against a
// freshly migrated database, and on such a database the states the bugs lived
// in are not merely absent, they are unrepresentable.
//
//   - The `repos` table on the author's machine holds TWO rows for one
//     repository: an old one with root_path=/Users/joe/git/tl and
//     git_common_dir=/Users/joe/git/tl (the pre-git_common_dir backfill shape,
//     48% of live rows), and a newer one with root_path="" and
//     git_common_dir=/Users/joe/git/tl/.git carrying every recent change. The
//     write path resolved the second, the read path resolved the first, and
//     171 of 174 published bundles shipped `sessions: []`. Every bundle test
//     seeded its repo without a GitCommonDir, so UpsertRepo backfilled
//     git_common_dir := root_path and the two identity keys were IDENTICAL —
//     the one shape in which the divergence is mathematically invisible.
//   - `capture_sessions` holds 151 rows under the retired
//     "<revisionID>-<contentHash>" id scheme against 2 under the current
//     "src-<hash>" one, 121 rows with shareable_at NULL, and 0 ever uploaded,
//     while `sessions` is empty and joins to none of them.
//
// Everything here is written through the writers production uses. There is no
// exported raw-INSERT escape hatch on purpose: seeding through a function with
// no production callers is how a test proves a state that a real machine can
// never reach, which is exactly what store.WriteSession did for the session
// bug. TestSeedersAreProductionWriters enforces that rule mechanically.
//
// Layering. This package imports internal/storage and internal/totalitytest and
// nothing else from tl. internal/vcs imports internal/storage, so importing
// vcs here would both risk a cycle and drag vcs's package init into the fast
// test binaries of internal/storage, internal/reviewbundle and
// internal/provenance.
package storagetest

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/totalitytest"
)

// CommitSelfReportSessionID is the fossil session id left behind by the retired
// `tl commit` self-report. storage.Open deletes it and its links on every open;
// nothing else in the tree covers that migration.
const CommitSelfReportSessionID = "totality-commit-self-report"

// Harness is an open tl database plus the ids the shapes created in it.
type Harness struct {
	// Store is the production writer surface.
	Store *storage.Store
	// DB is a second handle used only for counting and inspection, mirroring
	// how the hook tests read the database back.
	DB *sql.DB
	// World is the filesystem/git world this database belongs to.
	World *totalitytest.World

	// LegacyRepoID is the pre-git_common_dir row: root_path set,
	// git_common_dir equal to it. Live analogue: repos id 1.
	LegacyRepoID int64
	// CurrentRepoID is the row a modern writer resolves and lands changes on:
	// git_common_dir set to <root>/.git, root_path empty. Live analogue:
	// repos id 16.
	CurrentRepoID int64
}

// Shape mutates a freshly opened database into one of the states observed on a
// real machine. Shapes are applied in the order they are passed.
type Shape func(t *testing.T, h *Harness)

// New opens an isolated tl database in its own totalitytest.World and applies shapes.
func New(t *testing.T, shapes ...Shape) *Harness {
	t.Helper()
	return NewInWorld(t, totalitytest.NewWorld(t), shapes...)
}

// NewInWorld opens the tl database belonging to an existing world, so a test
// can create its git repositories first and then shape the database around
// their real paths.
func NewInWorld(t *testing.T, world *totalitytest.World, shapes ...Shape) *Harness {
	t.Helper()
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

	inspect, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() (inspection handle) error = %v", err)
	}
	t.Cleanup(func() { _ = inspect.Close() })

	h := &Harness{Store: store, DB: inspect, World: world}
	for _, shape := range shapes {
		shape(t, h)
	}
	return h
}

// Reopen runs storage.Open again against the same file.
//
// This is the only way to exercise the migration and repair passes at
// db.go:82-152 against non-empty tables: in every other test they run on a
// database that was empty a microsecond earlier, so deleteCommitSelfReportSession,
// repairDanglingSessionLinks and the git_common_dir backfill are all no-ops.
func (h *Harness) Reopen(t *testing.T) {
	t.Helper()
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() (reopen) error = %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
}

// Count runs a scalar count query against the inspection handle.
func (h *Harness) Count(t *testing.T, query string, args ...any) int {
	t.Helper()
	var n int
	if err := h.DB.QueryRowContext(context.Background(), query, args...).Scan(&n); err != nil {
		t.Fatalf("count %q: %v", query, err)
	}
	return n
}

// RepoRow reads one repos row by id.
func (h *Harness) RepoRow(t *testing.T, id int64) storage.Repo {
	t.Helper()
	var repo storage.Repo
	var commonDir sql.NullString
	err := h.DB.QueryRowContext(context.Background(),
		`SELECT id, root_path, git_common_dir, backend FROM repos WHERE id = ?`, id).
		Scan(&repo.ID, &repo.RootPath, &commonDir, &repo.Backend)
	if err != nil {
		t.Fatalf("read repos row %d: %v", id, err)
	}
	repo.GitCommonDir = commonDir.String
	return repo
}

// SeedChange writes a change and its revision through the writers the git
// record path uses, and returns the change row id.
func (h *Harness) SeedChange(t *testing.T, repoID int64, revisionID, commitID, description string, files []string) int64 {
	t.Helper()
	ctx := context.Background()
	id, err := h.Store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      revisionID,
		CurrentCommitID: commitID,
		Description:     description,
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange(%s) error = %v", revisionID, err)
	}
	if files == nil {
		files = []string{}
	}
	filesJSON, err := json.Marshal(files)
	if err != nil {
		t.Fatalf("marshal changed files: %v", err)
	}
	if err := h.Store.WriteChangeRevision(ctx, storage.ChangeRevision{
		ChangeID:      id,
		JJCommitID:    commitID,
		JJOperationID: "op-" + revisionID,
		ChangedFiles:  string(filesJSON),
		CreatedAt:     1,
	}); err != nil {
		t.Fatalf("WriteChangeRevision(%s) error = %v", revisionID, err)
	}
	return id
}

// SeedObservedSession writes a sessions row and links it to a change through
// the writers the push path uses: storage.Store.UpsertObservedSession is the
// only production writer of `sessions` (vcs.AttachSessionsFromHunkLinks), and
// WriteChangeSessions is the only production writer of `change_sessions`.
//
// Deliberately NOT store.WriteSession. That function has zero production
// callers and was the seeding path for every session test while the sole
// production writer could never bootstrap its first row — the suite proved a
// state the product could not reach.
func (h *Harness) SeedObservedSession(t *testing.T, changeID int64, sessionID, tool, repoRoot string) {
	t.Helper()
	ctx := context.Background()
	if err := h.Store.UpsertObservedSession(ctx, storage.Session{
		ID:        sessionID,
		CreatedAt: 1,
		Command:   tool,
		Cwd:       repoRoot,
		TLVersion: "test",
		Source:    &tool,
		RepoRoot:  &repoRoot,
	}); err != nil {
		t.Fatalf("UpsertObservedSession(%s) error = %v", sessionID, err)
	}
	if changeID == 0 {
		return
	}
	if err := h.Store.WriteChangeSessions(ctx, changeID, []string{sessionID}, 1); err != nil {
		t.Fatalf("WriteChangeSessions(%s) error = %v", sessionID, err)
	}
}

// DriftedRepoIdentity reproduces the two-row `repos` state one repository
// actually occupies on the author's machine.
//
//	id 1  root_path=/Users/joe/git/tl  git_common_dir=/Users/joe/git/tl
//	id 16 root_path=""                 git_common_dir=/Users/joe/git/tl/.git
//
// Row one is what the git_common_dir backfill leaves when `git rev-parse` fails
// (it falls back to the root path); 54 of 113 live rows are in that shape, and
// no test can produce it by accident because test repositories are always live
// git repos. Row two is what a modern UpsertRepo inserts when it matches
// neither key, and its empty root_path is permanent: UpsertRepo's existing-row
// UPDATE never rewrites root_path and no repair migration touches it.
//
// Both rows are created through storage.Store.UpsertRepo, so the state is one a
// real sequence of pushes produces rather than one a raw INSERT invented. Every
// dependent row a caller seeds afterwards should hang off CurrentRepoID, which
// is where a modern write lands.
func DriftedRepoIdentity(rootPath, gitCommonDir string) Shape {
	return func(t *testing.T, h *Harness) {
		t.Helper()
		ctx := context.Background()
		legacy, err := h.Store.UpsertRepo(ctx, storage.Repo{
			// GitCommonDir left empty on purpose: UpsertRepo then backfills it
			// to RootPath, which is the legacy row's defining shape.
			RootPath:  rootPath,
			Backend:   "git",
			CreatedAt: 1_700_000_000_000,
			UpdatedAt: 1_700_000_000_000,
		})
		if err != nil {
			t.Fatalf("UpsertRepo(legacy row) error = %v", err)
		}
		current, err := h.Store.UpsertRepo(ctx, storage.Repo{
			RootPath:     "",
			GitCommonDir: gitCommonDir,
			Backend:      "git",
			CreatedAt:    1_780_000_000_000,
			UpdatedAt:    1_780_000_000_000,
		})
		if err != nil {
			t.Fatalf("UpsertRepo(current row) error = %v", err)
		}
		if legacy == current {
			t.Fatalf("DriftedRepoIdentity produced one row (id %d) for root=%q common=%q; the drifted state needs two",
				legacy, rootPath, gitCommonDir)
		}
		h.LegacyRepoID = legacy
		h.CurrentRepoID = current
	}
}

// WeatheredNeighbourRepos adds the noise every real repos table carries, so a
// resolver has to step over it rather than being handed a single-row universe.
//
// Live shape being reproduced: 58 of 113 rows point at directories that no
// longer exist (torn-down sandboxes and expired /var/folders temp dirs), 32
// carry backend="" — a value no production code writes — and 24 still say
// backend='jj' on a git-native machine because UpsertRepo deliberately refuses
// to downgrade jj.
func WeatheredNeighbourRepos() Shape {
	return func(t *testing.T, h *Harness) {
		t.Helper()
		ctx := context.Background()
		for _, repo := range []storage.Repo{
			// A path that no longer exists on disk: any resolver that shells
			// out to git against a stored root_path fails on these.
			{RootPath: "/var/folders/zz/TestSomethingLongGone1234/001", Backend: "", CreatedAt: 1, UpdatedAt: 2},
			// Sticky jj on a machine that has been git-native for months.
			{RootPath: "/private/tmp/sandbox-a", Backend: "jj", CreatedAt: 3, UpdatedAt: 4},
			// A live-looking neighbour whose common dir must never be confused
			// with the repository under test.
			{RootPath: "/Users/someone/git/neighbour", GitCommonDir: "/Users/someone/git/neighbour/.git", Backend: "git", CreatedAt: 5, UpdatedAt: 6},
		} {
			if _, err := h.Store.UpsertRepo(ctx, repo); err != nil {
				t.Fatalf("UpsertRepo(neighbour %s) error = %v", repo.RootPath, err)
			}
		}
	}
}

// LegacyStagedSession is one capture_sessions row in the retired id scheme.
type LegacyStagedSession struct {
	SessionID  string
	Tool       string
	RevisionID string
	Payload    []byte
	// Shareable marks the row eligible for upload. Leave it false to reproduce
	// the dominant live state: 121 of 153 rows have shareable_at NULL, so they
	// are invisible to the ShareableSessions query production uploads through
	// while remaining fully visible to the PendingSessions query tests assert
	// through.
	Shareable bool
}

// LegacyCaptureSessions stages capture rows under the retired
// "<revisionID>-<contentHash[:24]>" id scheme.
//
// 151 of 153 live rows use it against 2 under the current "src-<hex>" scheme,
// and because the two schemes hash different inputs a restage under the new
// scheme can never ON CONFLICT against a legacy row — it inserts a second copy
// of the same bytes. No test in the tree mixes the schemes, so this is the
// state MarkSessionsShareableForSourcesOf exists to sweep and never sees.
//
// The rows go in through storage.CaptureStage.StageSession with an explicit ID
// built by storage.CaptureRowID, both production functions.
func LegacyCaptureSessions(rows ...LegacyStagedSession) Shape {
	return func(t *testing.T, h *Harness) {
		t.Helper()
		ctx := context.Background()
		stager, err := storage.OpenCaptureStager(ctx)
		if err != nil {
			t.Fatalf("OpenCaptureStager() error = %v", err)
		}
		for _, row := range rows {
			payload := row.Payload
			if len(payload) == 0 {
				payload = []byte(fmt.Sprintf(`{"sessionID":%q}`, row.SessionID))
			}
			hash := storage.PayloadContentHash(payload)
			id := storage.CaptureRowID(row.RevisionID, hash)
			if id == "" {
				t.Fatalf("LegacyCaptureSessions: revision %q and hash %q produce no legacy row id", row.RevisionID, hash)
			}
			if err := stager.StageSession(ctx, storage.StagedSession{
				ID:          id,
				SessionID:   row.SessionID,
				Tool:        row.Tool,
				PayloadJSON: payload,
				RevisionID:  row.RevisionID,
				ContentHash: hash,
				SourcePath:  filepath.Join(h.World.Home, ".claude", "projects", row.SessionID+".jsonl"),
			}); err != nil {
				t.Fatalf("StageSession(%s) error = %v", id, err)
			}
			if row.Shareable {
				if _, err := stager.MarkSessionsShareableByID(ctx, []string{id}, storage.CaptureAttestation{
					Name:  "Test",
					Email: "test@example.com",
				}); err != nil {
					t.Fatalf("MarkSessionsShareableByID(%s) error = %v", id, err)
				}
			}
		}
	}
}

// FossilCommitSelfReportSession writes the `totality-commit-self-report` sessions row
// an older tl binary left behind.
//
// It is the only fossil storage.Open actively deletes, the deletion runs on
// every single open, and it has zero test coverage anywhere in the tree — the
// live database shows it already fired (change_sessions is empty with
// sqlite_sequence=5, exactly the five rows the deletion targets). A database
// carrying rows an older binary wrote is the only case in which any of the
// migration passes does anything at all.
func FossilCommitSelfReportSession() Shape {
	return func(t *testing.T, h *Harness) {
		t.Helper()
		empty := ""
		if err := h.Store.UpsertObservedSession(context.Background(), storage.Session{
			ID:        CommitSelfReportSessionID,
			CreatedAt: 1,
			Command:   "tl commit",
			// cwd='' and repo_root NULL are what makes it a fossil: it can
			// never match a repository, so it only ever contributed noise.
			Cwd:       empty,
			TLVersion: "0.0.0",
		}); err != nil {
			t.Fatalf("seed %s error = %v", CommitSelfReportSessionID, err)
		}
	}
}

// NoSessions asserts the session tables are empty, and is the precondition a
// cold-start test needs to state out loud.
//
// internal/hooks/session_attach_test.go makes this assertion inline for exactly
// one test; as a shape it is reusable, so the next cold-start test cannot
// quietly start seeding the table it exists to verify gets written.
func NoSessions() Shape {
	return func(t *testing.T, h *Harness) {
		t.Helper()
		for _, table := range []string{"sessions", "change_sessions", "change_session_provenance"} {
			if n := h.Count(t, `SELECT count(*) FROM `+table); n != 0 {
				t.Fatalf("%s rows = %d at construction, want 0: a cold-start test must not seed the tables it verifies", table, n)
			}
		}
	}
}

// AssertResolvesToRepo fails unless the read-side identity resolver picks
// wantID for this (gitCommonDir, rootPath) pair.
//
// wantID is the id the WRITE path returned when the state was built, so this
// pins the two sides of the boundary the session bug fell through. Callers in
// other packages compare their own resolver against the same id;
// internal/reviewbundle does that for findRepoID.
func AssertResolvesToRepo(t *testing.T, h *Harness, gitCommonDir, rootPath string, wantID int64) {
	t.Helper()
	repo, err := h.Store.FindRepoByIdentity(context.Background(), gitCommonDir, rootPath)
	if err != nil {
		t.Fatalf("FindRepoByIdentity(%q, %q) error = %v", gitCommonDir, rootPath, err)
	}
	if repo == nil {
		t.Fatalf("FindRepoByIdentity(%q, %q) = nil, want repos row %d", gitCommonDir, rootPath, wantID)
	}
	if repo.ID != wantID {
		t.Fatalf("FindRepoByIdentity(%q, %q) = repos row %d (root_path=%q), want row %d — the read path resolved a different repository than the write path",
			gitCommonDir, rootPath, repo.ID, repo.RootPath, wantID)
	}
}
