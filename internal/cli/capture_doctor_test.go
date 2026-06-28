package cli

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/hooks"
	"github.com/satoricorp/gx/internal/storage"
)

func TestPrintCaptureDoctorHumanOutputUsesCompactRows(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("NO_COLOR", "1")
	status := captureDoctorJSON{
		HookApplicable:  true,
		HookInstalled:   true,
		RepoHooksOK:     true,
		RepoHooksTotal:  1,
		UploadAuthed:    true,
		UploadAPI:       "https://api.gx.run",
		PendingExtracts: 300,
		PendingSessions: 176,
		CursorReachable: true,
		DiskFreeGB:      153,
	}

	var out bytes.Buffer
	printCaptureDoctor(&out, status)
	text := out.String()
	for _, want := range []string{
		"Hooks installed",
		"ok",
		"Upload",
		"ok: https://api.gx.run",
		"Staging backlog",
		"476 pending",
		"Disk used",
		"0GB",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("doctor output missing %q in:\n%s", want, text)
		}
	}
	for _, old := range []string{
		"Capture",
		"Pre-push hook",
		"Registered repo hooks",
		"Upload credentials",
		"Cursor",
		"Cursor vscdb",
		"Disk free",
	} {
		if strings.Contains(text, old) {
			t.Fatalf("doctor output should not contain %q in:\n%s", old, text)
		}
	}
}

func TestCaptureHookStatusNotApplicableOutsideGitRepo(t *testing.T) {
	installed, applicable, resolvedRoot := captureHookStatus(context.Background(), t.TempDir())
	if installed {
		t.Fatal("expected hook not installed outside a git repo")
	}
	if applicable {
		t.Fatal("expected hook check not applicable outside a git repo")
	}
	if resolvedRoot != "" {
		t.Fatalf("resolvedRoot = %q, want empty", resolvedRoot)
	}
}

func TestCaptureHookStatusResolvesGitRoot(t *testing.T) {
	repo := initDoctorGitRepo(t)
	wantRepo, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}
	subdir := filepath.Join(repo, "nested", "dir")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}

	installed, applicable, resolvedRoot := captureHookStatus(context.Background(), subdir)
	if installed {
		t.Fatal("expected hook not installed before gx init")
	}
	if !applicable {
		t.Fatal("expected hook check applicable inside a git repo")
	}
	if resolvedRoot != wantRepo {
		t.Fatalf("resolvedRoot = %q, want %q", resolvedRoot, wantRepo)
	}

	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, GXPath: "/usr/local/bin/gx"}); err != nil {
		t.Fatal(err)
	}
	installed, applicable, resolvedRoot = captureHookStatus(context.Background(), subdir)
	if !installed || !applicable || resolvedRoot != wantRepo {
		t.Fatalf("captureHookStatus() = installed %v applicable %v root %q, want true true %q", installed, applicable, resolvedRoot, wantRepo)
	}
}

func TestCaptureRegisteredRepoHooksReportsInstalledAndMissingHooks(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	installedRepo := initDoctorGitRepo(t)
	missingRepo := initDoctorGitRepo(t)
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: installedRepo, GXPath: "/usr/local/bin/gx"}); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, ctx)
	defer store.Close()
	for _, repoRoot := range []string{installedRepo, missingRepo} {
		if err := store.RecordInitializedRepo(ctx, repoRoot, 1); err != nil {
			t.Fatalf("RecordInitializedRepo(%s) error = %v", repoRoot, err)
		}
	}

	got := captureRegisteredRepoHooks(ctx)
	if len(got) != 2 {
		t.Fatalf("captureRegisteredRepoHooks() returned %d repos, want 2: %#v", len(got), got)
	}
	var installed, missing int
	for _, repo := range got {
		if !repo.GitReachable {
			t.Fatalf("repo should be reachable: %#v", repo)
		}
		if repo.HookInstalled {
			installed++
		} else {
			missing++
		}
	}
	if installed != 1 || missing != 1 {
		t.Fatalf("registered hooks installed=%d missing=%d, want 1 and 1: %#v", installed, missing, got)
	}
}

func TestCaptureDoctorStatusFallsBackToInitializedRepoRegistry(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	repo := initDoctorGitRepo(t)
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, GXPath: "/usr/local/bin/gx"}); err != nil {
		t.Fatal(err)
	}
	wantRepo, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}

	store := openTestStore(t, ctx)
	if err := store.RecordInitializedRepo(ctx, repo, 100); err != nil {
		t.Fatalf("RecordInitializedRepo() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	outsideRepo := t.TempDir()
	if err := os.Chdir(outsideRepo); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	}()

	status := captureDoctorStatus(ctx, "")
	if !status.HookApplicable || !status.HookInstalled || status.RepoRoot != wantRepo {
		t.Fatalf("captureDoctorStatus() hook installed=%v applicable=%v root=%q, want true true %q", status.HookInstalled, status.HookApplicable, status.RepoRoot, wantRepo)
	}
}

func TestCaptureDoctorStatusFallsBackToKnownRepoRegistryWhenInitializedReposEmpty(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	repo := initDoctorGitRepo(t)
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: repo, GXPath: "/usr/local/bin/gx"}); err != nil {
		t.Fatal(err)
	}
	wantRepo, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatal(err)
	}

	store := openTestStore(t, ctx)
	if _, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:  repo,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 100,
	}); err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	outsideRepo := t.TempDir()
	if err := os.Chdir(outsideRepo); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	}()

	status := captureDoctorStatus(ctx, "")
	if !status.HookApplicable || !status.HookInstalled || status.RepoRoot != wantRepo {
		t.Fatalf("captureDoctorStatus() hook installed=%v applicable=%v root=%q, want true true %q", status.HookInstalled, status.HookApplicable, status.RepoRoot, wantRepo)
	}
}

func TestCaptureRegisteredRepoHooksFallsBackToReachableKnownRepos(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	installedRepo := initDoctorGitRepo(t)
	missingRepo := initDoctorGitRepo(t)
	staleRepo := filepath.Join(t.TempDir(), "stale")
	if err := hooks.Install(hooks.InstallOptions{RepoRoot: installedRepo, GXPath: "/usr/local/bin/gx"}); err != nil {
		t.Fatal(err)
	}
	store := openTestStore(t, ctx)
	defer store.Close()
	for i, repoRoot := range []string{installedRepo, staleRepo, missingRepo} {
		if _, err := store.UpsertRepo(ctx, storage.Repo{
			RootPath:  repoRoot,
			Backend:   "jj",
			CreatedAt: int64(i + 1),
			UpdatedAt: int64(100 - i),
		}); err != nil {
			t.Fatalf("UpsertRepo(%s) error = %v", repoRoot, err)
		}
	}

	got := captureRegisteredRepoHooks(ctx)
	if len(got) != 2 {
		t.Fatalf("captureRegisteredRepoHooks() returned %d repos, want 2 reachable repos: %#v", len(got), got)
	}
	var installed, missing, unreachable int
	for _, repo := range got {
		if !repo.GitReachable {
			unreachable++
		}
		if repo.HookInstalled {
			installed++
		} else {
			missing++
		}
	}
	if installed != 1 || missing != 1 || unreachable != 0 {
		t.Fatalf("registered hooks installed=%d missing=%d unreachable=%d, want 1, 1, 0: %#v", installed, missing, unreachable, got)
	}
}

func TestCaptureDoctorIssuesExplainAuthOnly(t *testing.T) {
	status := captureDoctorJSON{
		HookApplicable:  true,
		HookInstalled:   true,
		RepoHooksOK:     true,
		UploadAuthed:    false,
		PendingExtracts: 14,
		PendingSessions: 12,
		CursorReachable: true,
		DiskFreeGB:      100,
	}

	issues := captureDoctorIssues(status)
	codes := make([]string, 0, len(issues))
	for _, issue := range issues {
		codes = append(codes, issue.Code)
		if strings.TrimSpace(issue.Action) == "" {
			t.Fatalf("issue %s missing action: %#v", issue.Code, issue)
		}
	}
	wantCodes := []string{"upload_auth_missing"}
	if !equalDoctorStrings(codes, wantCodes) {
		t.Fatalf("captureDoctorIssues() codes = %#v, want %#v", codes, wantCodes)
	}

	summary := diagnoseSummary(doctorJSON{
		Capture: captureDoctorJSON{Issues: issues},
		Cursor:  cursorStatusJSON{Found: true},
	}).Summary
	if !strings.Contains(summary, "upload auth missing") || strings.Contains(summary, "pending") {
		t.Fatalf("diagnoseSummary() = %q, want auth detail only", summary)
	}
}

func TestCaptureDoctorBacklogDoesNotAffectHealth(t *testing.T) {
	status := captureDoctorJSON{
		HookApplicable:  true,
		HookInstalled:   true,
		RepoHooksOK:     true,
		UploadAuthed:    true,
		PendingExtracts: 29,
		PendingSessions: 75,
		CursorReachable: true,
		DiskFreeGB:      100,
	}

	if !captureDoctorOK(status) {
		t.Fatal("captureDoctorOK() = false, want true for staged backlog")
	}
	if issues := captureDoctorIssues(status); len(issues) != 0 {
		t.Fatalf("captureDoctorIssues() = %#v, want no issues for staged backlog", issues)
	}
}

func TestDoctorStatsSummarizesStacksAndAgents(t *testing.T) {
	ctx := context.Background()
	t.Setenv("GX_HOME", t.TempDir())
	store := openTestStore(t, ctx)

	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:  "/repo",
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	draftChangeID := insertDoctorStatsChange(t, ctx, store, repoID, "draft-change")
	draftStackID := insertDoctorStatsStack(t, ctx, store, repoID, "gx/draft-stack", "draft")
	if err := store.AddChangeToStack(ctx, draftStackID, draftChangeID, 2); err != nil {
		t.Fatalf("AddChangeToStack(draft) error = %v", err)
	}
	manualChangeID := insertDoctorStatsChange(t, ctx, store, repoID, "manual-change")
	manualStackID := insertDoctorStatsStack(t, ctx, store, repoID, "feature/manual-stack", "draft")
	if err := store.AddChangeToStack(ctx, manualStackID, manualChangeID, 3); err != nil {
		t.Fatalf("AddChangeToStack(manual) error = %v", err)
	}
	publishedChangeID := insertDoctorStatsChange(t, ctx, store, repoID, "published-change")
	publishedStackID := insertDoctorStatsStack(t, ctx, store, repoID, "published-stack", "published")
	if err := store.AddChangeToStack(ctx, publishedStackID, publishedChangeID, 4); err != nil {
		t.Fatalf("AddChangeToStack(published) error = %v", err)
	}
	if err := store.UpsertCloudBookmark(ctx, storage.CloudBookmarkState{
		PostgresBookmarkID: "cloud-open",
		RepoID:             repoID,
		RepoFullName:       "satoricorp/gx",
		BranchName:         "gx/published-stack",
		Revision:           1,
		MergeStatus:        "open",
		UpdatedAtMs:        1,
		SyncedAt:           1,
	}); err != nil {
		t.Fatalf("UpsertCloudBookmark(open) error = %v", err)
	}
	if err := store.UpsertCloudBookmark(ctx, storage.CloudBookmarkState{
		PostgresBookmarkID: "cloud-merged",
		RepoID:             repoID,
		RepoFullName:       "satoricorp/gx",
		BranchName:         "gx/merged-stack",
		Revision:           1,
		MergeStatus:        "merged",
		UpdatedAtMs:        1,
		SyncedAt:           1,
	}); err != nil {
		t.Fatalf("UpsertCloudBookmark(merged) error = %v", err)
	}
	mergedChangeID := insertDoctorStatsChange(t, ctx, store, repoID, "merged-change")
	mergedStackID := insertDoctorStatsStack(t, ctx, store, repoID, "merged-stack", "merged")
	if err := store.AddChangeToStack(ctx, mergedStackID, mergedChangeID, 5); err != nil {
		t.Fatalf("AddChangeToStack(merged) error = %v", err)
	}
	_ = insertDoctorStatsStack(t, ctx, store, repoID, "gx/empty-stack", "draft")
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	stats := doctorStats(ctx, doctorJSON{
		Capture: captureDoctorJSON{RepoRoot: "/repo"},
		Cursor:  cursorStatusJSON{Found: true},
		Ledger: []ledgerRowJSON{
			{Agent: "cursor", Status: "indexed", Calls: "15"},
			{Agent: "codex", Status: "indexed", Calls: "7"},
			{Agent: "claude", Status: "waiting", Calls: "0"},
		},
	})
	if stats.ApprovedStacksWaitingForPublish != 1 {
		t.Fatalf("approved waiting = %d, want 1", stats.ApprovedStacksWaitingForPublish)
	}
	if stats.PublishedStacksWaitingForReview != 1 {
		t.Fatalf("published waiting = %d, want 1", stats.PublishedStacksWaitingForReview)
	}
	if len(stats.Agents) != 3 {
		t.Fatalf("agents = %d, want 3", len(stats.Agents))
	}
	if stats.Agents[0].Agent != "cursor" || stats.Agents[0].Health != "green" || stats.Agents[0].Sessions != 15 {
		t.Fatalf("cursor stats = %#v, want green 15 sessions", stats.Agents[0])
	}
	if stats.Agents[2].Agent != "claude" || stats.Agents[2].Health != "yellow" {
		t.Fatalf("claude stats = %#v, want yellow", stats.Agents[2])
	}
}

func initDoctorGitRepo(t *testing.T) string {
	t.Helper()
	requireGit(t)
	repo := t.TempDir()
	runGit(t, repo, "init")
	return repo
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func openTestStore(t *testing.T, ctx context.Context) *storage.Store {
	t.Helper()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		_ = db.Close()
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	return store
}

func insertDoctorStatsChange(t *testing.T, ctx context.Context, store *storage.Store, repoID int64, name string) int64 {
	t.Helper()
	id, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      name,
		CurrentCommitID: name + "-commit",
		Description:     name,
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange(%s) error = %v", name, err)
	}
	return id
}

func insertDoctorStatsStack(t *testing.T, ctx context.Context, store *storage.Store, repoID int64, name, status string) int64 {
	t.Helper()
	id, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         name,
		BookmarkName: name,
		BaseRef:      "main",
		BaseCommitID: "base",
		Status:       status,
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack(%s) error = %v", name, err)
	}
	return id
}

func equalDoctorStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
