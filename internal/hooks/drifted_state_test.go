package hooks_test

import (
	"context"
	"testing"

	"github.com/satoricorp/lgtm/internal/hooks"
	"github.com/satoricorp/lgtm/internal/storage"
	"github.com/satoricorp/lgtm/internal/storage/storagetest"
	"github.com/satoricorp/lgtm/internal/lgtmtest"
	"github.com/satoricorp/lgtm/internal/vcs"
)

const driftConversation = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"

const driftContent = "package gamma\n\nfunc GammaThree() int {\n\treturn 3\n}\n"

// TestRunPushPublishesSessionsOnAWeatheredDatabase drives a real push through a
// database in the state the author's machine is actually in, rather than
// through a database that was created microseconds earlier.
//
// TestRunPushAttachesSessionsFromColdStart, which is kept, covers the other
// real case: a brand-new install where nothing exists yet. Neither subsumes the
// other. What this test adds is every condition a fresh temp-dir database
// cannot express:
//
//   - `repos` holds TWO rows for this repository, the older one carrying
//     git_common_dir == root_path (the backfill-fallback shape, 48% of live
//     rows) and the newer one carrying root_path == "" — a row UpsertRepo's
//     UPDATE can never repair, because its SET list has no root_path.
//   - Neighbouring rows point at directories that no longer exist, carry
//     backend="" (a value no production code writes), and carry a sticky
//     backend='jj' from before the git-native pivot.
//   - `capture_sessions` is dominated by the retired
//     "<revisionID>-<contentHash>" id scheme, none of it ever uploaded, most of
//     it not even marked shareable.
//   - A `lgtm-commit-self-report` fossil session written by an older binary is
//     sitting in `sessions`, so storage.Open's repair passes actually have work
//     to do for once.
//
// The load-bearing assertion is the repo_id one: the push must land its change
// on the empty-root_path row, which is what forces the bundle read to cross the
// write/read identity boundary rather than trivially agreeing with it.
func TestRunPushPublishesSessionsOnAWeatheredDatabase(t *testing.T) {
	ctx := context.Background()
	world := lgtmtest.NewWorld(t)
	repo := world.NewRepo(t)
	if _, err := vcs.NewService().InitAtPath(ctx, repo.Root, vcs.InitOptions{}); err != nil {
		t.Fatal(err)
	}

	h := storagetest.NewInWorld(t, world,
		storagetest.DriftedRepoIdentity(repo.Root, repo.GitCommonDir),
		storagetest.WeatheredNeighbourRepos(),
		storagetest.LegacyCaptureSessions(
			storagetest.LegacyStagedSession{SessionID: "old-conv-a", Tool: "claude", RevisionID: "revlegacyaaaa"},
			storagetest.LegacyStagedSession{SessionID: "old-conv-b", Tool: "codex", RevisionID: "revlegacybbbb"},
		),
		storagetest.FossilCommitSelfReportSession(),
	)

	base := repo.Rev(t, "HEAD")
	commit := repo.Commit(t, map[string]string{"gamma.go": driftContent}, "add gamma")
	world.WriteClaudeTranscript(t, repo, lgtmtest.Transcript{
		ConversationID: driftConversation,
		File:           "gamma.go",
		Content:        driftContent,
	})

	refRange := repo.RefRange(t, base, commit.SHA)
	outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot: repo.Root,
		Remote:   "origin",
		LocalRef: "refs/heads/main",
		RefRange: refRange,
		HeadSHA:  commit.SHA,
		HomeDir:  world.Home,
		Tools:    []string{"claude"},
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v, want nil so the hook never blocks a push", err)
	}
	if outcome.CaptureError != "" {
		t.Fatalf("CaptureError = %q", outcome.CaptureError)
	}
	if outcome.AttachError != "" {
		t.Fatalf("AttachError = %q", outcome.AttachError)
	}
	if outcome.AttachedSessions != 1 {
		t.Fatalf("AttachedSessions = %d, want 1", outcome.AttachedSessions)
	}

	// The premise this test rests on. If the change did NOT land on the
	// empty-root_path row then the drift was not actually in play and the
	// session assertions below prove nothing.
	var changeRepoID int64
	if err := h.DB.QueryRowContext(ctx,
		`SELECT repo_id FROM changes WHERE current_commit_id = ?`, commit.SHA).Scan(&changeRepoID); err != nil {
		t.Fatalf("read repo_id for the pushed commit: %v", err)
	}
	if changeRepoID != h.CurrentRepoID {
		t.Fatalf("pushed change landed on repos row %d, want the drifted row %d (root_path='') — this test is not exercising the drift",
			changeRepoID, h.CurrentRepoID)
	}
	if row := h.RepoRow(t, changeRepoID); row.RootPath != "" {
		t.Fatalf("the row the push wrote to has root_path=%q, want empty", row.RootPath)
	}

	sessions := outcome.Publication.Artifact.Sessions
	if len(sessions) != 1 || sessions[0].ID != driftConversation {
		t.Fatalf("artifact sessions = %+v, want [%s]: the bundle resolved a different repos row than the attach did",
			sessions, driftConversation)
	}
	if sessions[0].Command != "claude" {
		t.Fatalf("session command = %q, want claude", sessions[0].Command)
	}

	rebuilt := buildPushBundle(t, ctx, repo.Root, refRange, commit.SHA)
	if len(rebuilt.Sessions) != 1 || rebuilt.Sessions[0].ID != driftConversation {
		t.Fatalf("rebuilt bundle sessions = %+v, want the attached session", rebuilt.Sessions)
	}

	// storage.Open ran several times inside the push, so the fossil an older
	// binary left behind must be gone — while the session this push wrote is
	// untouched. Nothing else in the tree covers that migration.
	if n := h.Count(t, `SELECT count(*) FROM sessions WHERE id = ?`, storagetest.CommitSelfReportSessionID); n != 0 {
		t.Fatalf("fossil %q rows after push = %d, want 0", storagetest.CommitSelfReportSessionID, n)
	}
	if n := h.Count(t, `SELECT count(*) FROM sessions WHERE id = ?`, driftConversation); n != 1 {
		t.Fatalf("pushed session rows = %d, want 1", n)
	}

	// The legacy-scheme capture rows belong to revisions this push did not
	// carry, so marking this push's revisions shareable must leave them alone.
	// On a real machine 151 of 153 rows look like this and none has uploaded.
	if n := h.Count(t, `SELECT count(*) FROM capture_sessions WHERE id NOT LIKE 'src-%' AND shareable_at IS NULL`); n != 2 {
		t.Fatalf("untouched legacy capture rows = %d, want 2", n)
	}
}

// TestRunPushFromLinkedWorktreeOnAWeatheredDatabase is the linked-worktree push
// with the identity drift underneath it.
//
// TestRunPushFromLinkedWorktreePublishesSessions covers the worktree case on a
// clean database, where exactly one repos row exists. Here the worktree pushes
// into a database that already holds a stale row for the same repository — the
// combination that made the bug survive: the worktree root is not in `repos` at
// all, and the row that is has the wrong root_path.
func TestRunPushFromLinkedWorktreeOnAWeatheredDatabase(t *testing.T) {
	ctx := context.Background()
	world := lgtmtest.NewWorld(t)
	main := world.NewRepo(t)
	if _, err := vcs.NewService().InitAtPath(ctx, main.Root, vcs.InitOptions{}); err != nil {
		t.Fatal(err)
	}

	h := storagetest.NewInWorld(t, world,
		storagetest.DriftedRepoIdentity(main.Root, main.GitCommonDir),
		storagetest.WeatheredNeighbourRepos(),
	)

	worktree := main.AddWorktree(t, t.TempDir()+"-linked", "feature")
	base := worktree.Rev(t, "HEAD")
	commit := worktree.Commit(t, map[string]string{"gamma.go": driftContent}, "add gamma")
	world.WriteClaudeTranscript(t, worktree, lgtmtest.Transcript{
		ConversationID: driftConversation,
		File:           "gamma.go",
		Content:        driftContent,
	})

	refRange := worktree.RefRange(t, base, commit.SHA)
	outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot: worktree.Root,
		Remote:   "origin",
		LocalRef: "refs/heads/feature",
		RefRange: refRange,
		HeadSHA:  commit.SHA,
		HomeDir:  world.Home,
		Tools:    []string{"claude"},
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if outcome.AttachError != "" {
		t.Fatalf("AttachError = %q", outcome.AttachError)
	}
	if outcome.AttachedSessions != 1 {
		t.Fatalf("AttachedSessions = %d, want 1", outcome.AttachedSessions)
	}

	if n := h.Count(t, `SELECT count(*) FROM repos WHERE root_path = ?`, worktree.Root); n != 0 {
		t.Fatalf("repos rows for the worktree root = %d, want 0: linked worktrees share the main repo's row on purpose", n)
	}

	sessions := outcome.Publication.Artifact.Sessions
	if len(sessions) != 1 || sessions[0].ID != driftConversation {
		t.Fatalf("artifact sessions = %+v, want [%s]", sessions, driftConversation)
	}
	rebuilt := buildPushBundle(t, ctx, worktree.Root, refRange, commit.SHA)
	if len(rebuilt.Sessions) != 1 || rebuilt.Sessions[0].ID != driftConversation {
		t.Fatalf("rebuilt bundle sessions = %+v, want the attached session", rebuilt.Sessions)
	}
}

// TestLgtmTestRevisionTrailerMatchesProduction pins internal/lgtmtest's copy of the
// lgtm trailer against the real one.
//
// lgtmtest cannot import internal/vcs — internal/vcs imports internal/storage,
// which internal/storage/storagetest builds on, and dragging vcs into the fast
// storage/reviewbundle test binaries is the thing the layering exists to
// prevent. So the trailer format is duplicated, and this assertion lives here,
// in a package that already imports both. Without it a drift in the trailer
// would silently stop the harness's commits from producing `changes` rows, and
// every session assertion built on them would pass vacuously.
func TestLgtmTestRevisionTrailerMatchesProduction(t *testing.T) {
	const revisionID = "AbCdEfGhIjKlMnOpQrS"
	if got, want := lgtmtest.RevisionTrailerLine(revisionID), vcs.RevisionTrailerLine(revisionID); got != want {
		t.Fatalf("lgtmtest.RevisionTrailerLine() = %q, want %q", got, want)
	}
	if got, want := lgtmtest.StampRevisionTrailer("add gamma", revisionID), vcs.StampRevisionTrailer("add gamma", revisionID); got != want {
		t.Fatalf("lgtmtest.StampRevisionTrailer() = %q, want %q", got, want)
	}
	// The ids the harness mints must be ones production would accept, or
	// RecoverMissingRevisions would drop the commits it stamps.
	if id := lgtmtest.NewRevisionID(t); !vcs.ValidRevisionID(id) {
		t.Fatalf("lgtmtest.NewRevisionID() = %q, which vcs.ValidRevisionID rejects", id)
	}
}

// TestLgtmTestWorldIsolatesLgtmHome pins the guard that the 32 stray rows in the
// author's real ~/.lgtm/lgtm.db exist for: four tests that never set LGTM_HOME wrote
// straight into the production database, and storage.DefaultDir still falls
// back to ~/.lgtm whenever the variable is unset.
func TestLgtmTestWorldIsolatesLgtmHome(t *testing.T) {
	world := lgtmtest.NewWorld(t)
	dir, err := storage.DefaultDir()
	if err != nil {
		t.Fatalf("storage.DefaultDir() error = %v", err)
	}
	if dir != world.LgtmHome {
		t.Fatalf("storage.DefaultDir() = %q, want the world's LGTM_HOME %q", dir, world.LgtmHome)
	}
}
