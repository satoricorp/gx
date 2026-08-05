package hooks_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/hooks"
	"github.com/satoricorp/gx/internal/reviewbundle"
	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/vcs"
)

const coldStartConversation = "11111111-2222-3333-4444-555555555555"

const coldStartContent = "package alpha\n\nfunc AlphaOne() int {\n\treturn 41\n}\n"

// TestRunPushAttachesSessionsFromColdStart is the test the shipped defect would
// have failed. It seeds nothing: the `sessions` table is asserted empty before
// the push, so the only thing that can populate it is the push path itself.
//
// The retired commit-time matcher was gated on `sessions` already being
// non-empty while being the only writer of that table, so on a fresh install it
// could never write a first row and every published bundle carried
// `sessions: []` forever. The whole suite stayed green because every other
// session test seeded via store.WriteSession, a function with no production
// callers at all. WriteSession has since been deleted, seeding now goes through
// the writers production runs, and TestSeedersAreProductionWriters in
// internal/storage/storagetest fails on any storage method that drifts back
// into that position.
//
// This test is the cold-start half of the picture and stays exactly as it is: a
// brand-new install is a real case. The weathered half —  a database whose repo
// identity has drifted, carrying legacy capture rows and a fossil session — is
// in drifted_state_test.go, built with internal/storage/storagetest.
func TestRunPushAttachesSessionsFromColdStart(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")

	repo, base, head := coldStartRepo(t, map[string]string{"alpha.go": coldStartContent})
	refRange := base + ".." + head
	transcript := filepath.Join(home, ".claude", "projects", claudeProjectSlug(repo), coldStartConversation+".jsonl")
	writeClaudeTranscript(t, transcript, repo, coldStartConversation, "alpha.go", coldStartContent)

	// The premise. If this is ever non-zero the test has started seeding the
	// very table it exists to prove gets written, and it stops being a
	// cold-start test.
	if n := sessionRowCount(t, ctx); n != 0 {
		t.Fatalf("sessions before push = %d, want 0: this test must not seed the table it verifies", n)
	}

	outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		LocalRef: "refs/heads/main",
		RefRange: refRange,
		HeadSHA:  head,
		HomeDir:  home,
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

	// The artifact this very push wrote must already carry the session:
	// EnqueueAdoptedPublication runs reviewbundle.BuildPush synchronously in
	// this process, which is why the attach has to happen before it.
	sessions := outcome.Publication.Artifact.Sessions
	if len(sessions) != 1 {
		t.Fatalf("artifact sessions = %+v, want exactly one", sessions)
	}
	if sessions[0].ID != coldStartConversation {
		t.Fatalf("session id = %q, want the conversation UUID %q (identity comes from the transcript path)",
			sessions[0].ID, coldStartConversation)
	}
	// The console reads session.command first and only falls back to
	// session.source; without this field every session renders as `command=?`.
	if sessions[0].Command != "claude" {
		t.Fatalf("session command = %q, want claude", sessions[0].Command)
	}
	// Both are written, so the console's fallback also resolves. Asserted so a
	// change that drops either one fails here rather than silently degrading the
	// rendered session to `command=?`.
	if sessions[0].Source == nil || *sessions[0].Source != "claude" {
		t.Fatalf("session source = %v, want claude", sessions[0].Source)
	}
	// Compared against the symlink-resolved path: on macOS t.TempDir hands back
	// /var/folders/... while git and the repo resolver both report the real
	// /private/var/folders/... underneath it.
	if sessions[0].RepoRoot == nil || *sessions[0].RepoRoot != resolvedPath(t, repo) {
		t.Fatalf("session repo_root = %v, want %q", sessions[0].RepoRoot, resolvedPath(t, repo))
	}

	// And the same rows are visible to a freshly built bundle, not just to the
	// artifact the push happened to hold in memory.
	rebuilt := buildPushBundle(t, ctx, repo, refRange, head)
	if len(rebuilt.Sessions) != 1 || rebuilt.Sessions[0].ID != coldStartConversation {
		t.Fatalf("rebuilt bundle sessions = %+v, want the one attached session", rebuilt.Sessions)
	}
	if rebuilt.SchemaVersion != 2 {
		t.Fatalf("schema version = %d, want 2", rebuilt.SchemaVersion)
	}
	if rebuilt.Sessions[0].Command != "claude" {
		t.Fatalf("rebuilt session command = %q, want claude", rebuilt.Sessions[0].Command)
	}
}

// TestRunPushAttachIsIdempotent covers the re-push and the retried push. The
// writer is idempotent by construction (ON CONFLICT on sessions, INSERT OR
// IGNORE on change_sessions); duplicated change_sessions rows would show up in
// the bundle as the same session listed twice.
func TestRunPushAttachIsIdempotent(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")

	repo, base, head := coldStartRepo(t, map[string]string{"alpha.go": coldStartContent})
	refRange := base + ".." + head
	transcript := filepath.Join(home, ".claude", "projects", claudeProjectSlug(repo), coldStartConversation+".jsonl")
	writeClaudeTranscript(t, transcript, repo, coldStartConversation, "alpha.go", coldStartContent)

	opts := hooks.PushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		LocalRef: "refs/heads/main",
		RefRange: refRange,
		HeadSHA:  head,
		HomeDir:  home,
		Tools:    []string{"claude"},
	}
	for i := 0; i < 2; i++ {
		outcome, err := hooks.RunPush(ctx, opts)
		if err != nil {
			t.Fatalf("RunPush() #%d error = %v", i+1, err)
		}
		if outcome.AttachError != "" {
			t.Fatalf("AttachError #%d = %q", i+1, outcome.AttachError)
		}
	}

	if n := tableCount(t, ctx, `SELECT count(*) FROM change_sessions`); n != 1 {
		t.Fatalf("change_sessions rows = %d, want 1 after two pushes", n)
	}
	if n := tableCount(t, ctx, `SELECT count(*) FROM sessions`); n != 1 {
		t.Fatalf("sessions rows = %d, want 1 after two pushes", n)
	}
	bundle := buildPushBundle(t, ctx, repo, refRange, head)
	if len(bundle.Sessions) != 1 {
		t.Fatalf("bundle sessions = %+v, want one entry, not a duplicate", bundle.Sessions)
	}
}

// TestRunPushAttachesSubagentsAsDistinctSessions pins the identity rule that
// broke an earlier attempt: Claude stamps every subagent transcript with the
// PARENT conversation's sessionId, so identity must come from the file path.
// Reading the in-file id would collapse all three transcripts below into one
// session and lose the subagents' evidence entirely.
func TestRunPushAttachesSubagentsAsDistinctSessions(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")

	const otherConversation = "99999999-8888-7777-6666-555555555555"
	const alpha = "package alpha\n\nfunc AlphaOne() int {\n\treturn 41\n}\n"
	const beta = "package beta\n\nfunc BetaTwo() string {\n\treturn \"beta\"\n}\n"
	const gamma = "package gamma\n\nfunc GammaThree() bool {\n\treturn true\n}\n"
	const delta = "package delta\n\nfunc DeltaFour() float64 {\n\treturn 4.0\n}\n"

	repo, base, head := coldStartRepo(t, map[string]string{
		"alpha.go": alpha,
		"beta.go":  beta,
		"gamma.go": gamma,
		"delta.go": delta,
	})
	refRange := base + ".." + head

	projects := filepath.Join(home, ".claude", "projects", claudeProjectSlug(repo))
	// Parent conversation, plus two subagents that carry the PARENT's id in-file.
	writeClaudeTranscript(t, filepath.Join(projects, coldStartConversation+".jsonl"),
		repo, coldStartConversation, "alpha.go", alpha)
	writeClaudeTranscript(t, filepath.Join(projects, coldStartConversation, "subagents", "agent-a.jsonl"),
		repo, coldStartConversation, "beta.go", beta)
	writeClaudeTranscript(t, filepath.Join(projects, coldStartConversation, "subagents", "agent-b.jsonl"),
		repo, coldStartConversation, "gamma.go", gamma)
	// A second conversation with a subagent of the same name: namespacing by
	// conversation is what stops these two agent-a files colliding.
	writeClaudeTranscript(t, filepath.Join(projects, otherConversation, "subagents", "agent-a.jsonl"),
		repo, otherConversation, "delta.go", delta)

	outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		LocalRef: "refs/heads/main",
		RefRange: refRange,
		HeadSHA:  head,
		HomeDir:  home,
		Tools:    []string{"claude"},
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if outcome.AttachError != "" {
		t.Fatalf("AttachError = %q", outcome.AttachError)
	}

	want := map[string]bool{
		coldStartConversation:                       true,
		coldStartConversation + ":subagent:agent-a": true,
		coldStartConversation + ":subagent:agent-b": true,
		otherConversation + ":subagent:agent-a":     true,
	}
	got := map[string]bool{}
	for _, session := range buildPushBundle(t, ctx, repo, refRange, head).Sessions {
		got[session.ID] = true
	}
	for id := range want {
		if !got[id] {
			t.Fatalf("bundle sessions = %v, missing %q", got, id)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("bundle sessions = %v, want exactly %v", got, want)
	}
}

// TestRunPushAttachesLinksToTheirOwnCommits covers a push carrying more than
// one commit: each session must land on the change row for the commit its hunk
// belongs to, never on whichever change was resolved first. Human-authored and
// temporal-tier links must not create sessions at all.
func TestRunPushAttachesLinksToTheirOwnCommits(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")

	const secondConversation = "abababab-cdcd-efef-0101-232323232323"
	const alpha = "package alpha\n\nfunc AlphaOne() int {\n\treturn 41\n}\n"
	const beta = "package beta\n\nfunc BetaTwo() string {\n\treturn \"beta\"\n}\n"

	repo := t.TempDir()
	initBareGitRepo(t, repo)
	if _, err := vcs.NewService().InitAtPath(ctx, repo, vcs.InitOptions{}); err != nil {
		t.Fatal(err)
	}
	base := gitRev(t, repo, "HEAD")
	firstSHA := commitWithTrailer(t, repo, map[string]string{"alpha.go": alpha}, "add alpha")
	// A commit with no agent transcript behind it at all: its hunks are
	// human-authored and must contribute no session.
	humanSHA := commitWithTrailer(t, repo, map[string]string{"human.go": "package human\n\nvar Hand = 7\n"}, "hand written")
	secondSHA := commitWithTrailer(t, repo, map[string]string{"beta.go": beta}, "add beta")

	projects := filepath.Join(home, ".claude", "projects", claudeProjectSlug(repo))
	writeClaudeTranscript(t, filepath.Join(projects, coldStartConversation+".jsonl"),
		repo, coldStartConversation, "alpha.go", alpha)
	writeClaudeTranscript(t, filepath.Join(projects, secondConversation+".jsonl"),
		repo, secondConversation, "beta.go", beta)

	outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot: repo,
		Remote:   "origin",
		LocalRef: "refs/heads/main",
		RefRange: base + ".." + secondSHA,
		HeadSHA:  secondSHA,
		HomeDir:  home,
		Tools:    []string{"claude"},
	})
	if err != nil {
		t.Fatalf("RunPush() error = %v", err)
	}
	if outcome.AttachError != "" {
		t.Fatalf("AttachError = %q", outcome.AttachError)
	}

	if got := sessionIDsForCommit(t, ctx, repo, firstSHA); len(got) != 1 || got[0] != coldStartConversation {
		t.Fatalf("sessions on first commit = %v, want [%s]", got, coldStartConversation)
	}
	if got := sessionIDsForCommit(t, ctx, repo, secondSHA); len(got) != 1 || got[0] != secondConversation {
		t.Fatalf("sessions on second commit = %v, want [%s]", got, secondConversation)
	}
	if got := sessionIDsForCommit(t, ctx, repo, humanSHA); len(got) != 0 {
		t.Fatalf("sessions on human-authored commit = %v, want none", got)
	}
}

// TestRunPushFromLinkedWorktreePublishesSessions pins the identity a push is
// resolved by on both sides of the write/read boundary.
//
// A linked worktree has its own toplevel but shares the main repository's git
// common dir, and repos rows are keyed on that shared common dir on purpose
// (vcs.RepoIdentityKey): storage.Store.UpsertRepo matches git_common_dir first
// and its existing-row UPDATE never rewrites root_path. So the attach writes
// change_sessions under the MAIN repo's row while the worktree's own path never
// appears in `repos` at all. reviewbundle read those rows back by root_path,
// found nothing, and published `sessions: []` with no error on any surface —
// on a brand new database, for every push from a worktree. gx's own
// demux-worktree flow and the yeet harness both push from linked worktrees.
func TestRunPushFromLinkedWorktreePublishesSessions(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_DISABLE_BACKGROUND_WORKERS", "1")

	const worktreeConversation = "99999999-8888-7777-6666-555555555555"
	const gamma = "package gamma\n\nfunc GammaThree() int {\n\treturn 3\n}\n"

	// The main repository registers the repos row first, exactly as it would on
	// a machine that has pushed from the main checkout before.
	main := t.TempDir()
	initBareGitRepo(t, main)
	if _, err := vcs.NewService().InitAtPath(ctx, main, vcs.InitOptions{}); err != nil {
		t.Fatal(err)
	}
	commitWithTrailer(t, main, map[string]string{"seed.go": "package seed\n"}, "seed")
	if _, err := vcs.NewService().ResolveGxRepoAtPath(ctx, main); err != nil {
		t.Fatal(err)
	}

	worktree := filepath.Join(t.TempDir(), "linked")
	runGit(t, main, "worktree", "add", "-b", "feature", worktree)
	base := gitRev(t, worktree, "HEAD")
	head := commitWithTrailer(t, worktree, map[string]string{"gamma.go": gamma}, "add gamma")

	writeClaudeTranscript(t,
		filepath.Join(home, ".claude", "projects", claudeProjectSlug(worktree), worktreeConversation+".jsonl"),
		worktree, worktreeConversation, "gamma.go", gamma)

	refRange := base + ".." + head
	outcome, err := hooks.RunPush(ctx, hooks.PushOptions{
		RepoRoot: worktree,
		Remote:   "origin",
		LocalRef: "refs/heads/feature",
		RefRange: refRange,
		HeadSHA:  head,
		HomeDir:  home,
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

	// The write landed. The point of this test is that the READ finds it.
	sessions := outcome.Publication.Artifact.Sessions
	if len(sessions) != 1 || sessions[0].ID != worktreeConversation {
		t.Fatalf("artifact sessions = %+v, want [%s] — the bundle resolved a different repos row than the attach did",
			sessions, worktreeConversation)
	}
	if sessions[0].Command != "claude" {
		t.Fatalf("session command = %q, want claude", sessions[0].Command)
	}

	rebuilt := buildPushBundle(t, ctx, worktree, refRange, head)
	if len(rebuilt.Sessions) != 1 || rebuilt.Sessions[0].ID != worktreeConversation {
		t.Fatalf("rebuilt bundle sessions = %+v, want the attached session", rebuilt.Sessions)
	}
}

// coldStartRepo builds a repo with an initial commit plus one trailer-stamped
// commit adding files.
//
// Two commits matter: with commitCount==0 the orchestrator falls back to
// <head>~20, which does not resolve in a one-commit repo and makes the run bail
// with "empty ref range" before any matching happens.
func coldStartRepo(t *testing.T, files map[string]string) (repoRoot, baseSHA, headSHA string) {
	t.Helper()
	repo := t.TempDir()
	initBareGitRepo(t, repo)
	if _, err := vcs.NewService().InitAtPath(context.Background(), repo, vcs.InitOptions{}); err != nil {
		t.Fatal(err)
	}
	base := gitRev(t, repo, "HEAD")
	head := commitWithTrailer(t, repo, files, "add sources")
	return repo, base, head
}

// initBareGitRepo creates the repo and its initial commit.
func initBareGitRepo(t *testing.T, repo string) {
	t.Helper()
	for _, args := range [][]string{
		{"init", "-b", "main"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test"},
	} {
		runGit(t, repo, args...)
	}
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("# test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repo, "add", "README.md")
	runGit(t, repo, "commit", "--no-verify", "-m", "initial")
}

// commitWithTrailer writes files and commits them with a gx revision trailer,
// so RecoverMissingRevisions can create the `changes` row the sessions attach to.
func commitWithTrailer(t *testing.T, repo string, files map[string]string, subject string) string {
	t.Helper()
	revisionID, err := vcs.GenerateRevisionID()
	if err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		path := filepath.Join(repo, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		runGit(t, repo, "add", name)
	}
	runGit(t, repo, "commit", "--no-verify", "-m", vcs.StampRevisionTrailer(subject, revisionID))
	return gitRev(t, repo, "HEAD")
}

// writeClaudeTranscript writes a one-line Claude transcript whose Write tool
// call produced exactly the committed file, so the matcher links it at tier 1.
//
// inFileSessionID is written into the body deliberately: subagent transcripts
// carry the parent conversation's id there, and the pipeline must ignore it in
// favour of the path.
func writeClaudeTranscript(t *testing.T, path, repoRoot, inFileSessionID, file, content string) {
	t.Helper()
	line, err := json.Marshal(map[string]any{
		"type":      "assistant",
		"cwd":       repoRoot,
		"sessionId": inFileSessionID,
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"message": map[string]any{
			"role":  "assistant",
			"model": "claude-probe",
			"content": []map[string]any{{
				"type": "tool_use",
				"id":   "t1",
				"name": "Write",
				"input": map[string]any{
					"file_path": filepath.ToSlash(filepath.Join(repoRoot, file)),
					"content":   content,
				},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(line, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// buildPushBundle rebuilds the publish artifact through the same path the hook
// uses, so what it asserts is what a real push would upload. Going through
// EnqueueAdoptedPublication rather than calling reviewbundle.BuildPush directly
// is deliberate: BuildPush only lists sessions for commits it was handed a
// RevisionID for, and assembling that PushResult by hand in the test would let
// the test pass on a shape production never produces.
func buildPushBundle(t *testing.T, ctx context.Context, repoRoot, refRange, headSHA string) reviewbundle.Bundle {
	t.Helper()
	result, err := hooks.EnqueueAdoptedPublication(ctx, hooks.AdoptPushOptions{
		RepoRoot: repoRoot,
		Remote:   "origin",
		LocalRef: "refs/heads/main",
		HeadSHA:  headSHA,
		RefRange: refRange,
	})
	if err != nil {
		t.Fatalf("EnqueueAdoptedPublication() error = %v", err)
	}
	return result.Artifact.Bundle
}

func sessionIDsForCommit(t *testing.T, ctx context.Context, repoRoot, sha string) []string {
	t.Helper()
	db := openTestDB(t, ctx)
	defer db.Close()
	rows, err := db.QueryContext(ctx, `
		SELECT cs.session_id
		FROM change_sessions cs
		JOIN changes c ON c.id = cs.change_id
		WHERE c.current_commit_id = ?
		ORDER BY cs.session_id
	`, sha)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return ids
}

func sessionRowCount(t *testing.T, ctx context.Context) int {
	t.Helper()
	return tableCount(t, ctx, `SELECT count(*) FROM sessions`)
}

func tableCount(t *testing.T, ctx context.Context, query string) int {
	t.Helper()
	db := openTestDB(t, ctx)
	defer db.Close()
	var n int
	if err := db.QueryRowContext(ctx, query).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func resolvedPath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return resolved
}

func openTestDB(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
