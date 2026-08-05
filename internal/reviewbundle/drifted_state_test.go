package reviewbundle

import (
	"context"
	"reflect"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/storage/storagetest"
	"github.com/satoricorp/gx/internal/vcs"
)

// TestBuildPushReadsTheDriftedRowThePushWroteTo is the bundle-side test that
// runs against the database shape a real machine has, rather than the pristine
// single-row one every other bundle test builds.
//
// The pristine tests cannot fail the way production did. seedRepo calls
// UpsertRepo without a GitCommonDir, so the backfill sets
// git_common_dir := root_path and the two identity keys become IDENTICAL — the
// one arrangement in which resolving by root_path and resolving by
// git_common_dir provably agree. On the author's machine they do not: one
// repository owns two rows, and the changes are split across them.
//
// What makes this stronger than TestBuildPushFindsSessionsWhenRepoRootPathDiverges
// (which is kept, and covers the same rule on a two-row-free database) is that
// the STALE row here is populated. It holds a change under the same revision id
// with a different session and different files, so a resolver that picks the
// wrong row publishes confidently wrong provenance instead of an empty list —
// and a test asserting only `len(sessions) == 1` would still pass.
func TestBuildPushReadsTheDriftedRowThePushWroteTo(t *testing.T) {
	ctx := context.Background()
	world := gxtest.NewWorld(t)
	// A real git repository, so git_common_dir is genuinely <root>/.git rather
	// than a string the test invented.
	repo := world.NewRepo(t)

	h := storagetest.NewInWorld(t, world,
		storagetest.DriftedRepoIdentity(repo.Root, repo.GitCommonDir),
		storagetest.WeatheredNeighbourRepos(),
	)

	// The stale row: what an older gx wrote before the identity key moved to
	// the git common dir. Live analogue: repos id 1, 17 changes, last touched
	// six days before the row that now receives every write.
	staleChange := h.SeedChange(t, h.LegacyRepoID, "gxr-drift", "commit-drift", "stale alpha", []string{"stale.go"})
	h.SeedObservedSession(t, staleChange, "stale-session", "codex", repo.Root)

	// The row a modern write lands on: root_path empty, identified only by the
	// git common dir. Live analogue: repos id 16, 21 changes, all recent.
	currentChange := h.SeedChange(t, h.CurrentRepoID, "gxr-drift", "commit-drift", "real alpha", []string{"real.go"})
	h.SeedObservedSession(t, currentChange, "real-session", "claude", repo.Root)

	// All three identity-aware resolvers must land on the same row: the write
	// side (UpsertRepo, which produced CurrentRepoID), the store's own
	// FindRepoByIdentity, and this package's findRepoID.
	storagetest.AssertResolvesToRepo(t, h, repo.GitCommonDir, repo.Root, h.CurrentRepoID)
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	defer db.Close()
	gotID, found, err := findRepoID(ctx, db, repo.GitCommonDir, repo.Root)
	if err != nil || !found {
		t.Fatalf("findRepoID() = (%d, %v, %v), want the current row", gotID, found, err)
	}
	if gotID != h.CurrentRepoID {
		t.Fatalf("findRepoID() = repos row %d, want %d (UpsertRepo's row); the read path resolved a different repository than the write path",
			gotID, h.CurrentRepoID)
	}

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-drift",
		Repo: vcs.RepoInfo{
			RootPath:     repo.Root,
			GitCommonDir: repo.GitCommonDir,
			Backend:      "git",
		},
		Commits: []vcs.PushedCommit{{
			CommitID:   "commit-drift",
			RevisionID: "gxr-drift",
			Message:    "real alpha",
			Files:      []string{"real.go"},
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Sessions) != 1 {
		t.Fatalf("bundle sessions = %#v, want exactly the session written under the current row", bundle.Sessions)
	}
	if bundle.Sessions[0].ID != "real-session" {
		t.Fatalf("bundle session = %q, want real-session; %q means the bundle read the stale repos row, which is the shape that published wrong provenance",
			bundle.Sessions[0].ID, bundle.Sessions[0].ID)
	}
	if len(bundle.Revisions) != 1 || !reflect.DeepEqual(bundle.Revisions[0].Files, []string{"real.go"}) {
		t.Fatalf("bundle revision files = %#v, want [real.go] from the current row", bundle.Revisions)
	}
}

// TestBuildPushFromLinkedWorktreeOnADriftedDatabase combines the two states
// that each independently produced an empty sessions list: a repository whose
// stored identity has drifted, and a push whose worktree root has never
// appeared in `repos` at all.
//
// gx's own demux-worktree flow and the yeet harness both push from linked
// worktrees, so this is not a hypothetical pairing.
func TestBuildPushFromLinkedWorktreeOnADriftedDatabase(t *testing.T) {
	ctx := context.Background()
	world := gxtest.NewWorld(t)
	main := world.NewRepo(t)
	worktree := main.AddWorktree(t, t.TempDir()+"-linked", "feature")

	h := storagetest.NewInWorld(t, world,
		storagetest.DriftedRepoIdentity(main.Root, main.GitCommonDir),
		storagetest.WeatheredNeighbourRepos(),
	)
	changeID := h.SeedChange(t, h.CurrentRepoID, "gxr-worktree", "commit-worktree", "gamma", []string{"gamma.go"})
	h.SeedObservedSession(t, changeID, "worktree-session", "claude", worktree.Root)

	// The worktree's own root is not in `repos` and never will be: linked
	// worktrees deliberately share the main repository's row via the common dir.
	if n := h.Count(t, `SELECT count(*) FROM repos WHERE root_path = ?`, worktree.Root); n != 0 {
		t.Fatalf("repos rows for the worktree root = %d, want 0", n)
	}

	bundle, err := BuildPush(ctx, vcs.PushResult{
		HeadCommitID: "commit-worktree",
		Repo: vcs.RepoInfo{
			RootPath:     worktree.Root,
			GitCommonDir: worktree.GitCommonDir,
			Backend:      "git",
		},
		Commits: []vcs.PushedCommit{{
			CommitID:   "commit-worktree",
			RevisionID: "gxr-worktree",
			Message:    "gamma",
			Files:      []string{"gamma.go"},
		}},
	})
	if err != nil {
		t.Fatalf("BuildPush() error = %v", err)
	}
	if len(bundle.Sessions) != 1 || bundle.Sessions[0].ID != "worktree-session" {
		t.Fatalf("bundle sessions = %#v, want [worktree-session]: the worktree root is not in `repos`, so only the git common dir can resolve this push",
			bundle.Sessions)
	}
}
