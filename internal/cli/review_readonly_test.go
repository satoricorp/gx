package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/cloud"
)

// recordingEndpoint stands in for GX Cloud and PostHog so the credential-gated
// paths actually execute instead of returning early. A `go test` binary is
// built without -ldflags, so the embedded buildconfig.PostHogKey and CloudURL
// are empty and every telemetry- and cloud-gated branch is dead code under
// test — which is exactly how a release binary came to create $GX_HOME during
// a read-only review without any test noticing. Pointing both at a live
// recorder makes those branches run and makes what they send observable.
type recordingEndpoint struct {
	URL string

	mu    sync.Mutex
	paths []string
}

func newRecordingEndpoint(t *testing.T) *recordingEndpoint {
	t.Helper()
	rec := &recordingEndpoint{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.mu.Lock()
		rec.paths = append(rec.paths, r.URL.Path)
		rec.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"id": "test", "status": "ok"})
	}))
	t.Cleanup(server.Close)
	rec.URL = server.URL
	return rec
}

func (r *recordingEndpoint) requestedPaths() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.paths...)
}

func (r *recordingEndpoint) requested(path string) bool {
	for _, got := range r.requestedPaths() {
		if got == path {
			return true
		}
	}
	return false
}

// setReviewCloudEnv wires cloud and telemetry to a recorder. Call it after
// setReviewGateEnv, which owns GX_HOME.
func setReviewCloudEnv(t *testing.T, rec *recordingEndpoint) {
	t.Helper()
	t.Setenv("GX_CLOUD_URL", rec.URL)
	t.Setenv("GX_POSTHOG_KEY", "phc-test-key")
	t.Setenv("GX_POSTHOG_HOST", rec.URL)
}

// staleGitStatCache backdates the working tree's mtimes so the index's stat
// cache no longer matches it. This is the ordinary state of a fresh clone or
// checkout — the CI case gx review is built for — and it is the condition
// under which git status and git diff rewrite .git/index. Without it the
// cache is already current, git has nothing to refresh, and an index-guard
// regression would sail past this test.
func staleGitStatCache(t *testing.T, root string) {
	t.Helper()
	backdated := time.Now().Add(-time.Hour)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		return os.Chtimes(path, backdated, backdated)
	})
	if err != nil {
		t.Fatalf("backdate working tree: %v", err)
	}
}

func gitIndexDigest(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".git", "index"))
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		t.Fatalf("read .git/index: %v", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// `gx review` has to be usable as a CI gate and on a checkout the reviewer
// does not own, so it must not auto-initialize anything: no git hooks, no GX
// home, no repo state. This test reviews a repo that has never run `gx init`
// and asserts the repo and the machine come out untouched.
func TestReviewNeverInitializesTheRepoOrTheMachine(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)
	recorder := newRecordingEndpoint(t)
	setReviewCloudEnv(t, recorder)

	// Point GX_HOME and HOME at paths that do not exist yet: anything that
	// opens the store or writes config has to create them, which makes the
	// mutation visible instead of silently landing in an existing directory.
	sandbox := t.TempDir()
	gxHome := filepath.Join(sandbox, "gx-home")
	fakeHome := filepath.Join(sandbox, "home")
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("HOME", fakeHome)

	hooksDir := filepath.Join(root, ".git", "hooks")
	before := hookDirEntries(t, hooksDir)
	staleGitStatCache(t, root)
	// Reviewing runs git status and git diff, which refresh the stat cache and
	// rewrite .git/index unless review keeps them off it.
	indexBefore := gitIndexDigest(t, root)

	out, err := runReviewCommand(t, "--base", "main")
	if err != nil {
		t.Fatalf("gx review error = %v\n%s", err, out)
	}
	if strings.TrimSpace(out) == "" {
		t.Fatalf("gx review produced no report in an uninitialized repo")
	}

	if after := hookDirEntries(t, hooksDir); !equalStrings(before, after) {
		t.Fatalf("gx review changed .git/hooks:\nbefore: %v\nafter:  %v", before, after)
	}
	for _, name := range before {
		data, readErr := os.ReadFile(filepath.Join(hooksDir, name))
		if readErr != nil {
			t.Fatalf("read hook %s: %v", name, readErr)
		}
		if strings.Contains(string(data), "gx ") || strings.Contains(string(data), "GX_") {
			t.Fatalf("gx review left a GX marker in .git/hooks/%s:\n%s", name, data)
		}
	}
	for _, path := range []string{gxHome, fakeHome, filepath.Join(root, ".gx")} {
		if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
			t.Fatalf("gx review created %s (stat error = %v), want it untouched", path, statErr)
		}
	}
	if after := gitIndexDigest(t, root); after != indexBefore {
		t.Fatalf("gx review rewrote .git/index:\nbefore: %s\nafter:  %s", indexBefore, after)
	}
	if recorder.requested("/v1/reported-logs") {
		t.Fatalf("a successful gx review uploaded a failure report: %v", recorder.requestedPaths())
	}
}

// A gate that fails is the tool working, not a crash. `gx review --fail-on`
// returning findings must not ship log tails and repo identity to GX Cloud,
// and must not mint a machine ID — a CI job that fails the gate on every run
// would otherwise upload on every run, from a machine that never ran gx.
func TestReviewGateFailureNeverAutoReportsOrWritesGXHome(t *testing.T) {
	root := newReviewGateRepo(t)
	writeTestFile(t, root, "package.json", "{\n  \"name\": \"example\"\n}\n")
	gitAddTestFiles(t, root, "package.json")
	runGitTest(t, root, "commit", "-m", "add package.json")
	t.Chdir(root)
	setReviewGateEnv(t)
	recorder := newRecordingEndpoint(t)
	setReviewCloudEnv(t, recorder)

	sandbox := t.TempDir()
	gxHome := filepath.Join(sandbox, "gx-home")
	fakeHome := filepath.Join(sandbox, "home")
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("HOME", fakeHome)

	staleGitStatCache(t, root)
	indexBefore := gitIndexDigest(t, root)

	out, err := runReviewCommand(t, "--scope", "dependencies", "--fail-on", "strong", "--no-publish")
	if err == nil {
		t.Fatalf("gx review --fail-on strong exited 0 with findings:\n%s", out)
	}
	if code := ExitCode(err); code != reviewFindingsExitCode {
		t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewFindingsExitCode, err)
	}

	if recorder.requested("/v1/reported-logs") {
		t.Fatalf("a failing gate auto-reported to GX Cloud: %v", recorder.requestedPaths())
	}
	for _, path := range []string{gxHome, fakeHome, filepath.Join(root, ".gx")} {
		if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
			t.Fatalf("a failing gate created %s (stat error = %v), want it untouched", path, statErr)
		}
	}
	if after := gitIndexDigest(t, root); after != indexBefore {
		t.Fatalf("a failing gate rewrote .git/index:\nbefore: %s\nafter:  %s", indexBefore, after)
	}
}

// signIn writes credentials into an existing GX_HOME. Uploading a failure
// report needs a cloud token — without one ReportLogs fails before it sends
// anything, which would make an "it did not upload" assertion pass for the
// wrong reason. A signed-in machine is where the auto-report actually fires,
// so that is where the gate has to be proven silent.
func signIn(t *testing.T, gxHome string) {
	t.Helper()
	if err := os.MkdirAll(gxHome, 0o755); err != nil {
		t.Fatalf("create gx home: %v", err)
	}
	// credentials.json nests the cloud section under "cloud"; a flat object
	// parses without error and leaves the token empty, which would quietly
	// make this a signed-out test again.
	creds := map[string]any{
		"cloud": map[string]any{
			"cli_session_token":      "test-session-token",
			"cli_session_expires_at": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339Nano),
			"user_id":                "user-1",
			"login":                  "tester",
			"machine_id":             "machine-1",
			"machine_name":           "test-machine",
		},
	}
	data, err := json.Marshal(creds)
	if err != nil {
		t.Fatalf("marshal credentials: %v", err)
	}
	if err := os.WriteFile(filepath.Join(gxHome, "credentials.json"), data, 0o600); err != nil {
		t.Fatalf("write credentials: %v", err)
	}
	// Prove the machine really is signed in. Without this, a change to the
	// credentials format turns every "it did not upload" assertion below into
	// a test that passes because the upload could not have happened anyway.
	if _, err := cloud.CloudAPIToken(); err != nil {
		t.Fatalf("test credentials do not authenticate: %v", err)
	}
}

// The same gate on a machine that *is* signed in, where the upload would
// actually go through. This is the configuration the auto-report was written
// for, so it is the one that proves the gate stays quiet.
func TestReviewGateFailureNeverAutoReportsWhenSignedIn(t *testing.T) {
	root := newReviewGateRepo(t)
	writeTestFile(t, root, "package.json", "{\n  \"name\": \"example\"\n}\n")
	gitAddTestFiles(t, root, "package.json")
	runGitTest(t, root, "commit", "-m", "add package.json")
	t.Chdir(root)
	setReviewGateEnv(t)
	recorder := newRecordingEndpoint(t)
	setReviewCloudEnv(t, recorder)

	gxHome := filepath.Join(t.TempDir(), "gx-home")
	t.Setenv("GX_HOME", gxHome)
	signIn(t, gxHome)

	out, err := runReviewCommand(t, "--scope", "dependencies", "--fail-on", "strong", "--no-publish")
	if err == nil {
		t.Fatalf("gx review --fail-on strong exited 0 with findings:\n%s", out)
	}
	if code := ExitCode(err); code != reviewFindingsExitCode {
		t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewFindingsExitCode, err)
	}
	if recorder.requested("/v1/reported-logs") {
		t.Fatalf("a failing gate uploaded logs and repo identity to GX Cloud: %v", recorder.requestedPaths())
	}
}

// `gx version` answers one question about the binary. A Dockerfile or CI step
// that runs it to check what it installed should not thereby acquire a machine
// ID and a $GX_HOME — the install event is not worth creating state on a
// machine that has not yet decided to use gx.
func TestVersionNeverWritesGXState(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)
	recorder := newRecordingEndpoint(t)
	setReviewCloudEnv(t, recorder)

	sandbox := t.TempDir()
	gxHome := filepath.Join(sandbox, "gx-home")
	fakeHome := filepath.Join(sandbox, "home")
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("HOME", fakeHome)

	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"version"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gx version error = %v", err)
	}
	if strings.TrimSpace(out.String()) == "" {
		t.Fatalf("gx version printed nothing")
	}
	for _, path := range []string{gxHome, fakeHome} {
		if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
			t.Fatalf("gx version created %s (stat error = %v), want it untouched", path, statErr)
		}
	}
}

// Not just gate verdicts: a genuine review failure uploads nothing either.
// Reviewing outside a git repo is a real error on a real code path — the kind
// the old auto-report existed to catch — so it is the case that proves the
// upload is gone rather than merely narrowed to exclude gates.
func TestReviewGenuineFailureNeverUploadsLogs(t *testing.T) {
	setReviewGateEnv(t)
	recorder := newRecordingEndpoint(t)
	setReviewCloudEnv(t, recorder)

	gxHome := filepath.Join(t.TempDir(), "gx-home")
	t.Setenv("GX_HOME", gxHome)
	signIn(t, gxHome)

	// Somewhere that is not a git repo at all.
	notARepo := t.TempDir()
	t.Chdir(notARepo)

	if _, err := runReviewCommand(t, "--no-publish"); err == nil {
		t.Fatalf("gx review succeeded outside a git repository")
	}
	if recorder.requested("/v1/reported-logs") {
		t.Fatalf("a failing gx review uploaded logs to GX Cloud: %v", recorder.requestedPaths())
	}
}

// The other gate outcome: nothing was reviewed. Same rule — a verdict, not a
// crash, so nothing gets uploaded.
func TestReviewNothingToReviewGateNeverAutoReports(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)
	recorder := newRecordingEndpoint(t)
	setReviewCloudEnv(t, recorder)

	sandbox := t.TempDir()
	t.Setenv("GX_HOME", filepath.Join(sandbox, "gx-home"))
	t.Setenv("HOME", filepath.Join(sandbox, "home"))

	out, err := runReviewCommand(t, "--fail-on", "any", "--no-publish")
	if err == nil {
		t.Fatalf("gx review --fail-on any exited 0 without reviewing anything:\n%s", out)
	}
	if code := ExitCode(err); code != reviewNothingToReviewExitCode {
		t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewNothingToReviewExitCode, err)
	}
	if recorder.requested("/v1/reported-logs") {
		t.Fatalf("the nothing-to-review gate auto-reported to GX Cloud: %v", recorder.requestedPaths())
	}
}

func hookDirEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("read %s: %v", dir, err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

func equalStrings(a, b []string) bool {
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
