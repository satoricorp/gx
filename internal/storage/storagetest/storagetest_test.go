package storagetest_test

import (
	"context"
	"testing"

	"github.com/satoricorp/totality/internal/storage"
	"github.com/satoricorp/totality/internal/storage/storagetest"
	"github.com/satoricorp/totality/internal/totalitytest"
)

// TestDriftedRepoIdentityReproducesTheLiveTwoRowState proves the shape actually
// builds the state it claims to, so tests that depend on it are not quietly
// running against a single-row database again.
//
// The reference is the author's real ~/.totality/totality.db, where /Users/joe/git/tx has
// TWO repos rows: id 1 (root_path=/Users/joe/git/tx, git_common_dir the same,
// 17 stale changes) and id 16 (root_path="", git_common_dir=.../tx/.git, 21
// changes, all of them newer). Both rows are built here through
// storage.Store.UpsertRepo, so the state is one a real sequence of pushes
// produces rather than one a raw INSERT invented.
func TestDriftedRepoIdentityReproducesTheLiveTwoRowState(t *testing.T) {
	const root = "/Users/joe/git/tx"
	const commonDir = "/Users/joe/git/tx/.git"

	h := storagetest.New(t, storagetest.DriftedRepoIdentity(root, commonDir), storagetest.WeatheredNeighbourRepos())

	legacy := h.RepoRow(t, h.LegacyRepoID)
	if legacy.RootPath != root || legacy.GitCommonDir != root {
		t.Fatalf("legacy row = {root_path:%q git_common_dir:%q}, want both %q (the pre-git_common_dir backfill shape)",
			legacy.RootPath, legacy.GitCommonDir, root)
	}
	current := h.RepoRow(t, h.CurrentRepoID)
	if current.RootPath != "" || current.GitCommonDir != commonDir {
		t.Fatalf("current row = {root_path:%q git_common_dir:%q}, want {\"\", %q} — the live row 16 shape",
			current.RootPath, current.GitCommonDir, commonDir)
	}
	if n := h.Count(t, `SELECT count(*) FROM repos WHERE root_path = ? OR git_common_dir = ?`, root, commonDir); n != 2 {
		t.Fatalf("repos rows for this repository = %d, want 2", n)
	}

	// The write path resolves the current row, so the read path must too.
	storagetest.AssertResolvesToRepo(t, h, commonDir, root, h.CurrentRepoID)

	// And the fourth resolver, pinned as a known divergence rather than left
	// latent: storage.Store.FindRepoByRoot matches root_path ONLY, so on this
	// state it returns the stale row — the exact row reviewbundle used to read
	// and publish `sessions: []` from. Its one production caller is
	// internal/cli/doctor.go; if that ever grows into a write path, this
	// assertion is where the mismatch surfaces.
	byRoot, err := h.Store.FindRepoByRoot(context.Background(), root)
	if err != nil {
		t.Fatalf("FindRepoByRoot() error = %v", err)
	}
	if byRoot == nil || byRoot.ID != h.LegacyRepoID {
		t.Fatalf("FindRepoByRoot(%q) = %#v, want the stale row %d: this resolver is root-only by contract and disagrees with the other three on a drifted database",
			root, byRoot, h.LegacyRepoID)
	}
}

// TestDriftedRepoIdentityFromRealWorktreePaths runs the same shape against real
// git paths, where git_common_dir is genuinely <root>/.git rather than a string
// the test made up.
func TestDriftedRepoIdentityFromRealWorktreePaths(t *testing.T) {
	world := totalitytest.NewWorld(t)
	repo := world.NewRepo(t)
	h := storagetest.NewInWorld(t, world, storagetest.DriftedRepoIdentity(repo.Root, repo.GitCommonDir))

	storagetest.AssertResolvesToRepo(t, h, repo.GitCommonDir, repo.Root, h.CurrentRepoID)

	// A linked worktree shares the main checkout's common dir, so it must
	// resolve to the same row even though its own root has never been stored.
	worktree := repo.AddWorktree(t, t.TempDir()+"-linked", "feature")
	storagetest.AssertResolvesToRepo(t, h, worktree.GitCommonDir, worktree.Root, h.CurrentRepoID)
}

// TestDriftedRowSurvivesAnUnrelatedRepositoryRegistering pins the invariant a
// weathered database must keep: registering some other repository must not
// disturb the row this one's changes hang off.
//
// It also documents the boundary of that invariant, which only the weathered
// fixture makes visible. `repos.root_path` is UNIQUE and can be the empty
// string, and UpsertRepo's existing-row UPDATE has no root_path in its SET
// list, so once a row is inserted with an empty root_path it stays that way
// forever and there is at most ONE such row in the whole database. That row is
// therefore a single global slot: an UpsertRepo whose RootPath is absent lands
// on it either through the SELECT (`OR root_path = ""`) or, if the SELECT is
// guarded, through the INSERT's `ON CONFLICT(root_path) DO UPDATE`, which
// rewrites git_common_dir just the same. Guarding the SELECT alone changes
// nothing; closing it means either never minting empty-root rows or dropping
// the UNIQUE constraint, both of which are schema decisions rather than a
// resolver fix.
//
// So this test asserts the case that does hold and that every real repository
// falls into: a registrant with a root path of its own.
func TestDriftedRowSurvivesAnUnrelatedRepositoryRegistering(t *testing.T) {
	ctx := context.Background()
	world := totalitytest.NewWorld(t)
	victim := world.NewRepo(t)
	h := storagetest.NewInWorld(t, world, storagetest.DriftedRepoIdentity(victim.Root, victim.GitCommonDir))

	if row := h.RepoRow(t, h.CurrentRepoID); row.RootPath != "" {
		t.Fatalf("precondition: current row root_path = %q, want empty", row.RootPath)
	}
	if n := h.Count(t, `SELECT count(*) FROM repos WHERE root_path = ''`); n != 1 {
		t.Fatalf("empty-root_path rows = %d, want exactly 1: root_path is UNIQUE, so it is a single global slot", n)
	}

	intruder, err := h.Store.UpsertRepo(ctx, storage.Repo{
		RootPath:     "/somewhere/else",
		GitCommonDir: "/somewhere/else/.git",
		Backend:      "git",
		CreatedAt:    1_790_000_000_000,
		UpdatedAt:    1_790_000_000_000,
	})
	if err != nil {
		t.Fatalf("UpsertRepo(intruder) error = %v", err)
	}
	if intruder == h.CurrentRepoID || intruder == h.LegacyRepoID {
		t.Fatalf("an unrelated repository resolved to repos row %d, one of this repository's own rows", intruder)
	}
	if got := h.RepoRow(t, h.CurrentRepoID); got.GitCommonDir != victim.GitCommonDir || got.RootPath != "" {
		t.Fatalf("drifted row = {root_path:%q git_common_dir:%q} after an unrelated registration, want {\"\", %q}",
			got.RootPath, got.GitCommonDir, victim.GitCommonDir)
	}
	storagetest.AssertResolvesToRepo(t, h, victim.GitCommonDir, victim.Root, h.CurrentRepoID)
}

// TestFossilCommitSelfReportSessionIsRemovedOnReopen covers
// storage.deleteCommitSelfReportSession, which runs on every single open and
// has no other coverage in the tree. It can only do anything on a database an
// older binary wrote, and every other test opens a database that was empty a
// microsecond earlier.
func TestFossilCommitSelfReportSessionIsRemovedOnReopen(t *testing.T) {
	world := totalitytest.NewWorld(t)
	repo := world.NewRepo(t)
	h := storagetest.NewInWorld(t, world,
		storagetest.DriftedRepoIdentity(repo.Root, repo.GitCommonDir),
		storagetest.FossilCommitSelfReportSession(),
	)
	changeID := h.SeedChange(t, h.CurrentRepoID, "rev-1", "commit-1", "alpha", []string{"alpha.go"})
	// The fossil's five change_sessions rows are what the deletion actually
	// targets on a real machine.
	if err := h.Store.WriteChangeSessions(context.Background(), changeID,
		[]string{storagetest.CommitSelfReportSessionID}, 1); err != nil {
		t.Fatalf("link fossil session: %v", err)
	}
	h.SeedObservedSession(t, changeID, "real-session", "claude", repo.Root)

	if n := h.Count(t, `SELECT count(*) FROM sessions WHERE id = ?`, storagetest.CommitSelfReportSessionID); n != 1 {
		t.Fatalf("fossil session rows before reopen = %d, want 1", n)
	}

	h.Reopen(t)

	if n := h.Count(t, `SELECT count(*) FROM sessions WHERE id = ?`, storagetest.CommitSelfReportSessionID); n != 0 {
		t.Fatalf("fossil session rows after reopen = %d, want 0", n)
	}
	if n := h.Count(t, `SELECT count(*) FROM change_sessions WHERE session_id = ?`, storagetest.CommitSelfReportSessionID); n != 0 {
		t.Fatalf("fossil change_sessions rows after reopen = %d, want 0", n)
	}
	// The migration must be surgical: a genuine session on the same change has
	// to survive it.
	if n := h.Count(t, `SELECT count(*) FROM change_sessions WHERE session_id = 'real-session'`); n != 1 {
		t.Fatalf("real change_sessions rows after reopen = %d, want 1", n)
	}
}

// TestLegacyCaptureSessionsAreInvisibleToTheUploader pins the divergence
// between the query tests assert through and the query production uploads
// through.
//
// On the author's machine 121 of 153 capture_sessions rows have shareable_at
// NULL, and none has ever uploaded. PendingSessions (`WHERE uploaded_at IS
// NULL`) sees all of them and has zero production callers; ShareableSessions,
// which the uploader actually runs, sees none. Every staging test in the tree
// asserts through the first one.
func TestLegacyCaptureSessionsAreInvisibleToTheUploader(t *testing.T) {
	ctx := context.Background()
	h := storagetest.New(t, storagetest.LegacyCaptureSessions(
		storagetest.LegacyStagedSession{SessionID: "conv-a", Tool: "claude", RevisionID: "revaaaaaaaaa"},
		storagetest.LegacyStagedSession{SessionID: "conv-b", Tool: "claude", RevisionID: "revbbbbbbbbb"},
		storagetest.LegacyStagedSession{SessionID: "conv-c", Tool: "codex", RevisionID: "revccccccccc", Shareable: true},
	))

	stager, err := storage.OpenCaptureStager(ctx)
	if err != nil {
		t.Fatalf("OpenCaptureStager() error = %v", err)
	}

	pending, err := stager.PendingSessions(ctx)
	if err != nil {
		t.Fatalf("PendingSessions() error = %v", err)
	}
	if len(pending) != 3 {
		t.Fatalf("PendingSessions() = %d rows, want 3", len(pending))
	}
	shareable, err := stager.ShareableSessions(ctx)
	if err != nil {
		t.Fatalf("ShareableSessions() error = %v", err)
	}
	if len(shareable) != 1 || shareable[0].SessionID != "conv-c" {
		t.Fatalf("ShareableSessions() = %+v, want only the attested row; the other two are the 121-row live state the uploader never sees", shareable)
	}

	// Both id generations coexisting in one table is the live state: 151 legacy
	// rows against 2 current ones. Because the schemes hash different inputs,
	// restaging the same logical source under the current scheme cannot ON
	// CONFLICT against its legacy row — it inserts a second copy of the bytes.
	legacyID := storage.CaptureRowID("revaaaaaaaaa", storage.PayloadContentHash([]byte(`{"sessionID":"conv-a"}`)))
	if n := h.Count(t, `SELECT count(*) FROM capture_sessions WHERE id = ?`, legacyID); n != 1 {
		t.Fatalf("legacy row %q count = %d, want 1", legacyID, n)
	}
	if err := stager.StageSession(ctx, storage.StagedSession{
		ID:          storage.SessionSourceRowID("claude", h.World.Home+"/.claude/projects/conv-a.jsonl", "conv-a"),
		SessionID:   "conv-a",
		Tool:        "claude",
		PayloadJSON: []byte(`{"sessionID":"conv-a"}`),
		RevisionID:  "revaaaaaaaaa",
		SourcePath:  h.World.Home + "/.claude/projects/conv-a.jsonl",
	}); err != nil {
		t.Fatalf("StageSession(current scheme) error = %v", err)
	}
	if n := h.Count(t, `SELECT count(*) FROM capture_sessions WHERE session_id = 'conv-a'`); n != 2 {
		t.Fatalf("rows for conv-a = %d, want 2: the two id schemes cannot deduplicate against each other, which is why 109 live rows are byte-identical duplicates", n)
	}
}

// TestNoSessionsAllowsAColdStart is the positive control for the shape a
// cold-start test relies on.
func TestNoSessionsAllowsAColdStart(t *testing.T) {
	h := storagetest.New(t, storagetest.NoSessions())
	if n := h.Count(t, `SELECT count(*) FROM sessions`); n != 0 {
		t.Fatalf("sessions = %d, want 0", n)
	}
}
