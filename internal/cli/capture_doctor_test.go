package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/cloud"
	"github.com/satoricorp/gx/internal/hooks"
	"github.com/satoricorp/gx/internal/storage"
)

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

func TestCaptureDoctorIssuesExplainAuthAndBacklog(t *testing.T) {
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
	wantCodes := []string{"upload_auth_missing", "capture_backlog"}
	if !equalDoctorStrings(codes, wantCodes) {
		t.Fatalf("captureDoctorIssues() codes = %#v, want %#v", codes, wantCodes)
	}

	summary := diagnoseSummary(doctorJSON{
		Capture: captureDoctorJSON{Issues: issues},
		Cursor:  cursorStatusJSON{Found: true},
	}).Summary
	if !strings.Contains(summary, "upload auth missing") || !strings.Contains(summary, "26 pending") {
		t.Fatalf("diagnoseSummary() = %q, want auth and backlog details", summary)
	}
}

func TestCaptureDoctorStatusValidatesResolvedTokenAgainstUploadAPI(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")
	t.Setenv("GX_API_URL", "")

	var gotPath string
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		if gotAuth != "Bearer gho_api_token" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad token"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"user_id":           "user_1",
			"github_user_login": "octocat",
		})
	}))
	defer server.Close()
	t.Setenv("GX_CLOUD_URL", server.URL)

	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		GitHubAccessToken: "gho_api_token",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	status := captureDoctorStatus(context.Background(), "")
	if !status.UploadAuthed {
		t.Fatalf("UploadAuthed = false, error = %q", status.UploadAuthError)
	}
	if status.UploadAPI != server.URL {
		t.Fatalf("UploadAPI = %q, want %q", status.UploadAPI, server.URL)
	}
	if gotPath != "/v1/auth/me" {
		t.Fatalf("auth validation path = %q, want /v1/auth/me", gotPath)
	}
	if gotAuth != "Bearer gho_api_token" {
		t.Fatalf("Authorization = %q, want bearer upload token", gotAuth)
	}
}

func TestCaptureDoctorStatusReportsAPIRejectedTokenAsInvalid(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")
	t.Setenv("GX_API_URL", "")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer gxcs_bad" {
			t.Fatalf("Authorization = %q, want CLI session token", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
	}))
	defer server.Close()
	t.Setenv("GX_CLOUD_URL", server.URL)

	if err := cloud.SaveCloudCredentials(cloud.CloudCredentials{
		CLISessionToken: "gxcs_bad",
	}); err != nil {
		t.Fatalf("SaveCloudCredentials() error = %v", err)
	}

	status := captureDoctorStatus(context.Background(), "")
	if status.UploadAuthed {
		t.Fatal("UploadAuthed = true, want false")
	}
	if !strings.Contains(status.UploadAuthError, "Unauthorized") {
		t.Fatalf("UploadAuthError = %q, want API rejection", status.UploadAuthError)
	}
	var foundInvalid bool
	for _, issue := range status.Issues {
		if issue.Code == "upload_auth_invalid" {
			foundInvalid = true
			break
		}
	}
	if !foundInvalid {
		t.Fatalf("issues = %#v, want upload_auth_invalid", status.Issues)
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
