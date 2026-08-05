package storage

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRepoIdentityMigrationAndLinkedWorktreeUpsert(t *testing.T) {
	home := t.TempDir()
	t.Setenv("LGTM_HOME", home)
	repoRoot := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repoRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "init", "-b", "main")
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	dbPath := filepath.Join(home, "lgtm.db")
	legacy, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = legacy.Exec(`
		CREATE TABLE repos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			root_path TEXT NOT NULL UNIQUE,
			backend TEXT NOT NULL,
			default_remote TEXT,
			default_branch TEXT,
			authoring_base_ref TEXT,
			remote_url TEXT,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);
		INSERT INTO repos (root_path, backend, created_at, updated_at)
		VALUES (?, 'jj', 1, 1);
	`, repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}

	db, err := Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewStore(context.Background(), db)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	commonCommand := exec.Command("git", "rev-parse", "--path-format=absolute", "--git-common-dir")
	commonCommand.Dir = repoRoot
	commonOutput, err := commonCommand.Output()
	if err != nil {
		t.Fatal(err)
	}
	commonDir := filepath.Clean(strings.TrimSpace(string(commonOutput)))
	row, err := store.FindRepoByIdentity(context.Background(), commonDir, repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	if row == nil || row.GitCommonDir != commonDir || row.Backend != "jj" {
		t.Fatalf("migrated repo = %+v", row)
	}
	id, err := store.UpsertRepo(context.Background(), Repo{
		RootPath:     filepath.Join(t.TempDir(), "linked"),
		GitCommonDir: commonDir,
		Backend:      "git",
		UpdatedAt:    2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != row.ID {
		t.Fatalf("linked worktree repo id = %d, want %d", id, row.ID)
	}
	repos, err := store.ListRepos(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 || repos[0].Backend != "jj" {
		t.Fatalf("repos = %+v, want one preserved jj repo", repos)
	}
}

func TestOpenRepairsCorruptedJJChangeIDs(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()

	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() initial error = %v", err)
	}

	if _, err := db.ExecContext(ctx, `
		INSERT INTO repos (id, root_path, backend, created_at, updated_at)
		VALUES (1, '/repo', 'jj', 1, 1)
	`); err != nil {
		t.Fatalf("insert repo: %v", err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO changes (id, repo_id, jj_change_id, current_commit_id, description, parent_change_id, status, first_seen_at, updated_at)
		VALUES
		(1, 1, 'realchangeid', 'commit-a', 'feat one', NULL, 'draft', 1, 10),
		(2, 1, 'Done importing changes from the underlying Git repo.
realchangeid', 'commit-a', 'feat one', NULL, 'draft', 2, 20),
		(3, 1, 'Done importing changes from the underlying Git repo.
secondchangeid', 'commit-b', 'feat two', NULL, 'draft', 3, 30)
	`); err != nil {
		t.Fatalf("insert changes: %v", err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO pushes (id, repo_id, remote_name, branch_name, head_commit_id, current_change_id, created_at)
		VALUES (1, 1, 'origin', 'main', 'commit-a', 2, 100)
	`); err != nil {
		t.Fatalf("insert push: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("db.Close() error = %v", err)
	}

	db, err = Open(ctx)
	if err != nil {
		t.Fatalf("Open() repair error = %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM changes WHERE jj_change_id LIKE 'Done importing changes from the underlying Git repo.%'`).Scan(&count); err != nil {
		t.Fatalf("count corrupted rows: %v", err)
	}
	if count != 0 {
		t.Fatalf("corrupted change rows remain: %d", count)
	}

	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM changes WHERE repo_id = 1 AND jj_change_id = 'realchangeid'`).Scan(&count); err != nil {
		t.Fatalf("count canonical row: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one canonical realchangeid row, got %d", count)
	}

	var pushChangeID int64
	if err := db.QueryRowContext(ctx, `SELECT current_change_id FROM pushes WHERE id = 1`).Scan(&pushChangeID); err != nil {
		t.Fatalf("select repaired push: %v", err)
	}
	if pushChangeID != 1 {
		t.Fatalf("push current_change_id = %d, want 1", pushChangeID)
	}

	var normalized string
	if err := db.QueryRowContext(ctx, `SELECT jj_change_id FROM changes WHERE id = 3`).Scan(&normalized); err != nil {
		t.Fatalf("select normalized row: %v", err)
	}
	if normalized != "secondchangeid" {
		t.Fatalf("normalized change id = %q, want %q", normalized, "secondchangeid")
	}
}

func TestDemuxProposalRoundTrip(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	repoID, err := store.UpsertRepo(ctx, Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if err := store.UpsertDemuxProposal(ctx, DemuxProposal{
		ID:           "demux-test",
		RepoID:       repoID,
		BaseChangeID: "base-change",
		Status:       "pending",
		PayloadJSON:  `{"id":"demux-test"}`,
		CreatedAt:    10,
		UpdatedAt:    10,
	}); err != nil {
		t.Fatalf("UpsertDemuxProposal() error = %v", err)
	}
	got, err := store.FindDemuxProposal(ctx, "demux-test")
	if err != nil {
		t.Fatalf("FindDemuxProposal() error = %v", err)
	}
	if got == nil {
		t.Fatal("FindDemuxProposal() = nil")
	}
	if got.Status != "pending" || got.BaseChangeID != "base-change" || got.PayloadJSON != `{"id":"demux-test"}` {
		t.Fatalf("FindDemuxProposal() = %#v", got)
	}

	appliedAt := int64(20)
	if err := store.UpdateDemuxProposalStatus(ctx, "demux-test", "applied", 21, &appliedAt); err != nil {
		t.Fatalf("UpdateDemuxProposalStatus() error = %v", err)
	}
	got, err = store.FindDemuxProposal(ctx, "demux-test")
	if err != nil {
		t.Fatalf("FindDemuxProposal(applied) error = %v", err)
	}
	if got.Status != "applied" || got.AppliedAt == nil || *got.AppliedAt != appliedAt {
		t.Fatalf("applied proposal = %#v", got)
	}
}

func TestSessionContextRoundTrip(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	model := "gpt-5"
	if err := store.UpsertSessionContext(ctx, Session{
		ID:        "session-one",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       "/repo",
		TLVersion: "test",
	}, SessionContext{
		SessionID:   "session-one",
		Tool:        "codex",
		Model:       &model,
		Format:      "lgtm_session_events_v1",
		ContentJSON: []byte(`[{"session_id":"session-one","new_text":"redacted"}]`),
		CapturedAt:  2,
	}); err != nil {
		t.Fatalf("UpsertSessionContext() error = %v", err)
	}

	got, err := store.SessionContext(ctx, "session-one")
	if err != nil {
		t.Fatalf("SessionContext() error = %v", err)
	}
	if got == nil {
		t.Fatal("SessionContext() = nil")
	}
	if got.SessionID != "session-one" || got.Tool != "codex" || got.Model == nil || *got.Model != model || got.Format != "lgtm_session_events_v1" || string(got.ContentJSON) == "" || got.CapturedAt != 2 {
		t.Fatalf("SessionContext() = %#v", got)
	}
}

func TestOpenRepairsDanglingSessionLinks(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() initial error = %v", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = OFF`); err != nil {
		t.Fatalf("disable foreign keys: %v", err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO repos (id, root_path, backend, created_at, updated_at)
		VALUES (1, '/repo', 'jj', 1, 1);
		INSERT INTO changes (id, repo_id, jj_change_id, current_commit_id, description, parent_change_id, status, first_seen_at, updated_at)
		VALUES (1, 1, 'change-one', 'commit-one', 'one', NULL, 'draft', 1, 1);
		INSERT INTO sessions (id, created_at, command, cwd, lgtm_version)
		VALUES ('live-session', 1, 'codex', '/repo', 'test');
		INSERT INTO change_sessions (change_id, session_id, created_at)
		VALUES (1, 'live-session', 1), (999, 'live-session', 2), (1, 'missing-session', 3);
		INSERT INTO change_session_provenance (change_id, session_id, agent_tool, provider, model_id, source, process_name, created_at)
		VALUES
			(1, 'live-session', 'codex', '', '', NULL, NULL, 1),
			(999, 'live-session', 'codex', '', '', NULL, NULL, 2),
			(1, 'missing-session', 'codex', '', '', NULL, NULL, 3),
			(1, 'live-session-no-link', 'codex', '', '', NULL, NULL, 4)
	`); err != nil {
		t.Fatalf("insert dangling links: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("db.Close() error = %v", err)
	}

	db, err = Open(ctx)
	if err != nil {
		t.Fatalf("Open() repair error = %v", err)
	}
	defer db.Close()

	var count int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM change_sessions`).Scan(&count); err != nil {
		t.Fatalf("count change_sessions: %v", err)
	}
	if count != 1 {
		t.Fatalf("change_sessions count = %d, want 1", count)
	}
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM change_session_provenance`).Scan(&count); err != nil {
		t.Fatalf("count change_session_provenance: %v", err)
	}
	if count != 1 {
		t.Fatalf("change_session_provenance count = %d, want 1", count)
	}
}

func TestListReposOrdersByMostRecentlyUpdated(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	for _, repo := range []Repo{
		{RootPath: "/old", Backend: "jj", CreatedAt: 1, UpdatedAt: 10},
		{RootPath: "/new", Backend: "jj", CreatedAt: 2, UpdatedAt: 20},
	} {
		if _, err := store.UpsertRepo(ctx, repo); err != nil {
			t.Fatalf("UpsertRepo(%s) error = %v", repo.RootPath, err)
		}
	}

	got, err := store.ListRepos(ctx)
	if err != nil {
		t.Fatalf("ListRepos() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("ListRepos() returned %d repos, want 2", len(got))
	}
	if got[0].RootPath != "/new" || got[1].RootPath != "/old" {
		t.Fatalf("ListRepos() order = %#v, want /new then /old", got)
	}
}

func TestInitializedReposRoundTrip(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	if err := store.RecordInitializedRepo(ctx, "/old", 10); err != nil {
		t.Fatalf("RecordInitializedRepo(/old) error = %v", err)
	}
	if err := store.RecordInitializedRepo(ctx, "/new", 20); err != nil {
		t.Fatalf("RecordInitializedRepo(/new) error = %v", err)
	}
	if err := store.RecordInitializedRepo(ctx, "/old", 30); err != nil {
		t.Fatalf("RecordInitializedRepo(/old update) error = %v", err)
	}

	got, err := store.ListInitializedRepos(ctx)
	if err != nil {
		t.Fatalf("ListInitializedRepos() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("ListInitializedRepos() returned %d repos, want 2", len(got))
	}
	if got[0].RootPath != "/old" || got[0].CreatedAt != 10 || got[0].UpdatedAt != 30 {
		t.Fatalf("first initialized repo = %#v, want /old created 10 updated 30", got[0])
	}
	if got[1].RootPath != "/new" {
		t.Fatalf("second initialized repo = %#v, want /new", got[1])
	}
}

func TestStackTracksChangeAndCommitRefs(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	repoID, err := store.UpsertRepo(ctx, Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, Change{
		RepoID:          repoID,
		JJChangeID:      "change-a",
		CurrentCommitID: "commit-a",
		Description:     "feat a",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	headChange := "change-a"
	headCommit := "commit-a"
	stackID, err := store.UpsertStack(ctx, Stack{
		RepoID:       repoID,
		Name:         "feat a",
		BookmarkName: "feature/feat-a",
		BaseRef:      "main",
		BaseCommitID: "base-a",
		HeadChangeID: &headChange,
		HeadCommitID: &headCommit,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}

	if _, err := store.UpsertChange(ctx, Change{
		RepoID:          repoID,
		JJChangeID:      "change-a",
		CurrentCommitID: "commit-b",
		Description:     "feat a",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       2,
	}); err != nil {
		t.Fatalf("UpsertChange(rewrite) error = %v", err)
	}
	headCommit = "commit-b"
	if _, err := store.UpsertStack(ctx, Stack{
		RepoID:       repoID,
		Name:         "feat a",
		BookmarkName: "feature/feat-a",
		BaseRef:      "main",
		BaseCommitID: "base-a",
		HeadChangeID: &headChange,
		HeadCommitID: &headCommit,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    2,
	}); err != nil {
		t.Fatalf("UpsertStack(rewrite) error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 2); err != nil {
		t.Fatalf("AddChangeToStack(rewrite) error = %v", err)
	}

	stack, err := store.FindStackByBookmark(ctx, repoID, "feature/feat-a")
	if err != nil {
		t.Fatalf("FindStackByBookmark() error = %v", err)
	}
	if stack == nil || stack.HeadChangeID == nil || *stack.HeadChangeID != "change-a" || stack.HeadCommitID == nil || *stack.HeadCommitID != "commit-b" {
		t.Fatalf("stack refs = %#v", stack)
	}
	var stackChangeID, stackCommitID string
	if err := db.QueryRowContext(ctx, `
		SELECT jj_change_id, commit_id
		FROM stack_changes
		WHERE stack_id = ? AND change_id = ?
	`, stackID, changeID).Scan(&stackChangeID, &stackCommitID); err != nil {
		t.Fatalf("select stack change refs: %v", err)
	}
	if stackChangeID != "change-a" || stackCommitID != "commit-b" {
		t.Fatalf("stack change refs = %q/%q, want change-a/commit-b", stackChangeID, stackCommitID)
	}
}

func TestDeletePendingDemuxProposalsForRepo(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	repoID, err := store.UpsertRepo(ctx, Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	for _, proposal := range []DemuxProposal{
		{ID: "pending-a", RepoID: repoID, BaseChangeID: "change-a", Status: "pending", PayloadJSON: `{}`, CreatedAt: 10, UpdatedAt: 10},
		{ID: "pending-b", RepoID: repoID, BaseChangeID: "change-b", Status: "pending", PayloadJSON: `{}`, CreatedAt: 20, UpdatedAt: 20},
		{ID: "applied", RepoID: repoID, BaseChangeID: "change-c", Status: "applied", PayloadJSON: `{}`, CreatedAt: 30, UpdatedAt: 30},
	} {
		if err := store.UpsertDemuxProposal(ctx, proposal); err != nil {
			t.Fatalf("UpsertDemuxProposal(%s) error = %v", proposal.ID, err)
		}
	}
	if err := store.DeletePendingDemuxProposalsForRepo(ctx, repoID); err != nil {
		t.Fatalf("DeletePendingDemuxProposalsForRepo() error = %v", err)
	}
	list, err := store.ListDemuxProposals(ctx, repoID, "pending", 10)
	if err != nil {
		t.Fatalf("ListDemuxProposals(pending) error = %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("pending proposals after delete = %#v, want none", list)
	}
	got, err := store.FindDemuxProposal(ctx, "applied")
	if err != nil || got == nil || got.Status != "applied" {
		t.Fatalf("applied proposal = %#v, err = %v, want preserved", got, err)
	}
}

func TestFindLatestDemuxProposal(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	repoID, err := store.UpsertRepo(ctx, Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	for _, proposal := range []DemuxProposal{
		{ID: "old", RepoID: repoID, BaseChangeID: "base", Status: "pending", PayloadJSON: `{"id":"old"}`, CreatedAt: 10, UpdatedAt: 10},
		{ID: "applied", RepoID: repoID, BaseChangeID: "base", Status: "applied", PayloadJSON: `{"id":"applied"}`, CreatedAt: 30, UpdatedAt: 30},
		{ID: "new", RepoID: repoID, BaseChangeID: "base", Status: "pending", PayloadJSON: `{"id":"new"}`, CreatedAt: 20, UpdatedAt: 20},
	} {
		if err := store.UpsertDemuxProposal(ctx, proposal); err != nil {
			t.Fatalf("UpsertDemuxProposal(%s) error = %v", proposal.ID, err)
		}
	}

	got, err := store.FindLatestDemuxProposal(ctx, repoID, "pending")
	if err != nil {
		t.Fatalf("FindLatestDemuxProposal(pending) error = %v", err)
	}
	if got == nil || got.ID != "new" {
		t.Fatalf("FindLatestDemuxProposal(pending) = %#v, want new", got)
	}
	got, err = store.FindLatestDemuxProposal(ctx, repoID, "")
	if err != nil {
		t.Fatalf("FindLatestDemuxProposal(any) error = %v", err)
	}
	if got == nil || got.ID != "applied" {
		t.Fatalf("FindLatestDemuxProposal(any) = %#v, want applied", got)
	}
	list, err := store.ListDemuxProposals(ctx, repoID, "pending", 10)
	if err != nil {
		t.Fatalf("ListDemuxProposals(pending) error = %v", err)
	}
	if gotIDs := proposalIDs(list); !reflect.DeepEqual(gotIDs, []string{"new", "old"}) {
		t.Fatalf("ListDemuxProposals(pending) ids = %#v, want new, old", gotIDs)
	}
	list, err = store.ListDemuxProposals(ctx, repoID, "", 2)
	if err != nil {
		t.Fatalf("ListDemuxProposals(any) error = %v", err)
	}
	if gotIDs := proposalIDs(list); !reflect.DeepEqual(gotIDs, []string{"applied", "new"}) {
		t.Fatalf("ListDemuxProposals(any) ids = %#v, want applied, new", gotIDs)
	}
}

func proposalIDs(proposals []DemuxProposal) []string {
	out := make([]string, 0, len(proposals))
	for _, proposal := range proposals {
		out = append(out, proposal.ID)
	}
	return out
}

func TestChangeDemuxEvidenceRoundTrip(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	repoID, err := store.UpsertRepo(ctx, Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, Change{
		RepoID:          repoID,
		JJChangeID:      "change-1",
		CurrentCommitID: "commit-1",
		Description:     "split alpha",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.UpsertDemuxProposal(ctx, DemuxProposal{
		ID:           "demux-test",
		RepoID:       repoID,
		BaseChangeID: "base-change",
		Status:       "pending",
		PayloadJSON:  `{"id":"demux-test"}`,
		CreatedAt:    10,
		UpdatedAt:    10,
	}); err != nil {
		t.Fatalf("UpsertDemuxProposal() error = %v", err)
	}
	if err := store.WriteChangeDemuxEvidence(ctx, ChangeDemuxEvidence{
		ChangeID:           changeID,
		DemuxProposalID:    "demux-test",
		RevisionProposalID: "r1",
		Intent:             "split alpha",
		FilesJSON:          `["alpha.txt"]`,
		HunkIDsJSON:        `["h1"]`,
		UseHunks:           true,
		Confidence:         0.8,
		ProvenanceStatus:   "explicit",
		EvidenceJSON:       `{"revision":{"id":"r1"}}`,
		CreatedAt:          20,
	}); err != nil {
		t.Fatalf("WriteChangeDemuxEvidence() error = %v", err)
	}

	got, err := store.ListChangeDemuxEvidenceByProposal(ctx, "demux-test")
	if err != nil {
		t.Fatalf("ListChangeDemuxEvidenceByProposal() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("evidence rows = %#v, want 1", got)
	}
	if got[0].ChangeID != changeID || got[0].RevisionProposalID != "r1" || !got[0].UseHunks || got[0].HunkIDsJSON != `["h1"]` {
		t.Fatalf("evidence row = %#v", got[0])
	}
}

func TestWriteChangeSessionsPersistsAgentProvenance(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	repoID, err := store.UpsertRepo(ctx, Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, Change{
		RepoID:          repoID,
		JJChangeID:      "change-1",
		CurrentCommitID: "commit-1",
		Description:     "linked",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	source := "ambient"
	// UpsertObservedSession is the writer the push path uses, and the only
	// production writer of `sessions`. It cannot set process_name — no current
	// writer can, since the cursor ingest was retired — so the provenance row
	// this test expects carries a NULL process_name and the agent tool must
	// resolve from the command alone.
	if err := store.UpsertObservedSession(ctx, Session{
		ID:        "session-one",
		CreatedAt: 10,
		Command:   "codex exec",
		Cwd:       "/repo",
		TLVersion: "test",
		Source:    &source,
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}
	model := "gpt-5.1-code"
	if err := store.WriteRequest(ctx, Request{
		ID:             "request-one",
		SessionID:      "session-one",
		CreatedAt:      11,
		Provider:       "openai",
		Endpoint:       "/v1/responses",
		Method:         "POST",
		Model:          &model,
		RequestBody:    []byte(`{"input":"change main"}`),
		RequestHeaders: "{}",
	}); err != nil {
		t.Fatalf("WriteRequest() error = %v", err)
	}
	if err := store.WriteChangeSessions(ctx, changeID, []string{"session-one"}, 12); err != nil {
		t.Fatalf("WriteChangeSessions() error = %v", err)
	}

	var gotSession, gotAgent, gotProvider, gotModel string
	var gotSource, gotProcess sql.NullString
	var gotCreatedAt int64
	if err := store.db.QueryRowContext(ctx, `
		SELECT session_id, agent_tool, provider, model_id, source, process_name, created_at
		FROM change_session_provenance
		WHERE change_id = ?
	`, changeID).Scan(&gotSession, &gotAgent, &gotProvider, &gotModel, &gotSource, &gotProcess, &gotCreatedAt); err != nil {
		t.Fatalf("select change_session_provenance: %v", err)
	}
	if gotSession != "session-one" || gotAgent != "codex" || gotProvider != "openai" || gotModel != model || gotCreatedAt != 12 {
		t.Fatalf("provenance = %s/%s/%s/%s/%d, want session-one/codex/openai/%s/12", gotSession, gotAgent, gotProvider, gotModel, gotCreatedAt, model)
	}
	if gotSource.String != source || gotProcess.Valid {
		t.Fatalf("source/process = %#v/%#v, want %q/NULL", gotSource, gotProcess, source)
	}
}

func TestAgentLedgerSummaryBackfillsTokensFromRawResponseBodies(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	if err := store.UpsertObservedSession(ctx, Session{
		ID:        "codex-session",
		CreatedAt: 10,
		Command:   "codex",
		Cwd:       "/repo",
		TLVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}
	for _, requestID := range []string{"request-raw", "request-indexed"} {
		if err := store.WriteRequest(ctx, Request{
			ID:             requestID,
			SessionID:      "codex-session",
			CreatedAt:      11,
			Provider:       "openai",
			Endpoint:       "/v1/responses",
			Method:         "POST",
			RequestBody:    []byte(`{}`),
			RequestHeaders: `{}`,
		}); err != nil {
			t.Fatalf("WriteRequest(%s) error = %v", requestID, err)
		}
	}
	if err := store.WriteResponse(ctx, Response{
		ID:              "response-raw",
		RequestID:       "request-raw",
		CreatedAt:       12,
		CompletedAt:     13,
		StatusCode:      200,
		ResponseBody:    []byte("event: response.in_progress\ndata: {\"type\":\"response.in_progress\",\"response\":{\"usage\":{\"input_tokens\":50,\"output_tokens\":5}}}\n\nevent: response.completed\ndata: {\"type\":\"response.completed\",\"response\":{\"usage\":{\"input_tokens\":100,\"output_tokens\":25}}}\n\n"),
		ResponseHeaders: `{}`,
	}); err != nil {
		t.Fatalf("WriteResponse(raw) error = %v", err)
	}
	inputTokens := 7
	outputTokens := 3
	if err := store.WriteResponse(ctx, Response{
		ID:              "response-indexed",
		RequestID:       "request-indexed",
		CreatedAt:       14,
		CompletedAt:     15,
		StatusCode:      200,
		ResponseBody:    []byte(`{"usage":{"input_tokens":999,"output_tokens":999}}`),
		ResponseHeaders: `{}`,
		InputTokens:     &inputTokens,
		OutputTokens:    &outputTokens,
	}); err != nil {
		t.Fatalf("WriteResponse(indexed) error = %v", err)
	}

	summary, err := store.AgentLedgerSummary(ctx, "codex")
	if err != nil {
		t.Fatalf("AgentLedgerSummary() error = %v", err)
	}
	if summary.Tokens != 135 {
		t.Fatalf("AgentLedgerSummary().Tokens = %d, want 135", summary.Tokens)
	}
}

func TestSessionUsageAggregatesRequestModelsAndResponseTokens(t *testing.T) {
	t.Setenv("LGTM_HOME", t.TempDir())
	ctx := context.Background()
	db, err := Open(ctx)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	store, err := NewStore(ctx, db)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	if err := store.UpsertObservedSession(ctx, Session{
		ID:        "codex-session",
		CreatedAt: 10,
		Command:   "codex",
		Cwd:       "/repo",
		TLVersion: "test",
	}); err != nil {
		t.Fatalf("UpsertObservedSession() error = %v", err)
	}
	modelOne := "gpt-5.1-code"
	if err := store.WriteRequest(ctx, Request{
		ID:             "request-one",
		SessionID:      "codex-session",
		CreatedAt:      11,
		Provider:       "openai",
		Endpoint:       "/v1/responses",
		Method:         "POST",
		Model:          &modelOne,
		RequestBody:    []byte(`{"model":"gpt-5.1-code"}`),
		RequestHeaders: `{}`,
	}); err != nil {
		t.Fatalf("WriteRequest(request-one) error = %v", err)
	}
	modelTwo := "claude-sonnet-4-20250514"
	if err := store.WriteRequest(ctx, Request{
		ID:             "request-two",
		SessionID:      "codex-session",
		CreatedAt:      12,
		Provider:       "anthropic",
		Endpoint:       "/v1/messages",
		Method:         "POST",
		Model:          &modelTwo,
		RequestBody:    []byte(`{"model":"claude-sonnet-4-20250514"}`),
		RequestHeaders: `{}`,
	}); err != nil {
		t.Fatalf("WriteRequest(request-two) error = %v", err)
	}
	inputTokens := 10
	outputTokens := 5
	cacheReadTokens := 2
	cacheWriteTokens := 3
	if err := store.WriteResponse(ctx, Response{
		ID:               "response-one",
		RequestID:        "request-one",
		CreatedAt:        13,
		CompletedAt:      14,
		StatusCode:       200,
		ResponseBody:     []byte(`{}`),
		ResponseHeaders:  `{}`,
		InputTokens:      &inputTokens,
		OutputTokens:     &outputTokens,
		CacheReadTokens:  &cacheReadTokens,
		CacheWriteTokens: &cacheWriteTokens,
	}); err != nil {
		t.Fatalf("WriteResponse(response-one) error = %v", err)
	}
	if err := store.WriteResponse(ctx, Response{
		ID:              "response-two",
		RequestID:       "request-two",
		CreatedAt:       15,
		CompletedAt:     16,
		StatusCode:      200,
		ResponseBody:    []byte("data: {\"response\":{\"usage\":{\"input_tokens\":7,\"output_tokens\":4,\"cache_read_input_tokens\":1,\"cache_creation_input_tokens\":6}}}\n\n"),
		ResponseHeaders: `{}`,
	}); err != nil {
		t.Fatalf("WriteResponse(response-two) error = %v", err)
	}

	var modelsJSON string
	var gotInput, gotOutput, gotCacheRead, gotCacheWrite int
	if err := store.db.QueryRowContext(ctx, `
		SELECT models_json, input_tokens, output_tokens, cache_read_tokens, cache_write_tokens
		FROM sessions
		WHERE id = 'codex-session'
	`).Scan(&modelsJSON, &gotInput, &gotOutput, &gotCacheRead, &gotCacheWrite); err != nil {
		t.Fatalf("select session usage: %v", err)
	}
	if modelsJSON != `["gpt-5.1-code","claude-sonnet-4-20250514"]` {
		t.Fatalf("models_json = %s", modelsJSON)
	}
	if gotInput != 17 || gotOutput != 9 || gotCacheRead != 3 || gotCacheWrite != 9 {
		t.Fatalf("tokens = %d/%d/%d/%d, want 17/9/3/9", gotInput, gotOutput, gotCacheRead, gotCacheWrite)
	}
}
