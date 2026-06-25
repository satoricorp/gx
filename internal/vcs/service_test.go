package vcs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/gx/internal/gxconfig"
	"github.com/satoricorp/gx/internal/storage"
)

type fakeRunner struct {
	outputs       map[string][]string
	stdoutOutputs map[string][]string
	errors        map[string][]error
	calls         []string
}

func TestStackMatchesDirectNameRequiresStackNameOrBookmarkRef(t *testing.T) {
	body := storage.Stack{
		Name:         "Update terminal theme",
		BookmarkName: "feature/update-terminal-theme",
	}

	for _, name := range []string{"Update terminal theme", "feature/update-terminal-theme"} {
		if !stackMatchesDirectName(body, name) {
			t.Fatalf("stackMatchesDirectName(%q) = false, want true", name)
		}
	}
	if stackMatchesDirectName(body, "update-terminal-theme") {
		t.Fatalf("stackMatchesDirectName() accepted bookmark slug without conventional prefix")
	}
}

func TestRepoAuthoringBaseRefNormalizesLegacyGXStackBookmark(t *testing.T) {
	base := "gx/source-stack"
	repo := RepoInfo{AuthoringBase: &base}

	if got := repo.authoringBaseRef(); got != "feature/source-stack" {
		t.Fatalf("authoringBaseRef() = %q, want conventional stack bookmark", got)
	}
}

func TestResolveExecutableFromFindsLocalBin(t *testing.T) {
	home := t.TempDir()
	exe := filepath.Join(home, ".local", "bin", "jj")
	if err := os.MkdirAll(filepath.Dir(exe), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(exe, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := resolveExecutableFrom("jj", home, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != exe {
		t.Fatalf("resolveExecutableFrom() = %q, want %q", got, exe)
	}
}

func TestResolveExecutableFromFindsConfiguredDirectory(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "jj")
	if err := os.WriteFile(exe, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := resolveExecutableFrom("jj", "", []string{dir})
	if err != nil {
		t.Fatal(err)
	}
	if got != exe {
		t.Fatalf("resolveExecutableFrom() = %q, want %q", got, exe)
	}
}

func runnerKey(dir, name string, args ...string) string {
	return dir + "|" + name + "|" + strings.Join(args, " ")
}

func assertRunnerCalled(t *testing.T, calls []string, want string) {
	t.Helper()
	for _, call := range calls {
		if call == want {
			return
		}
	}
	t.Fatalf("missing call %q in %v", want, calls)
}

func assertRunnerNotCalled(t *testing.T, calls []string, forbidden string) {
	t.Helper()
	for _, call := range calls {
		if call == forbidden {
			t.Fatalf("unexpected call %q in %v", forbidden, calls)
		}
	}
}

func defaultFakeOutput(name string, args []string) (string, bool) {
	switch name {
	case "jj":
		if len(args) >= 2 && args[0] == "bookmark" && args[1] == "list" {
			return "", true
		}
		if len(args) >= 2 && args[0] == "log" {
			for i := 0; i+1 < len(args); i++ {
				if args[i] != "-T" {
					continue
				}
				switch args[i+1] {
				case "empty":
					return "false\n", true
				case "change_id":
					return "chg-default\n", true
				case "commit_id":
					return "commit-default\n", true
				}
			}
		}
	case "git":
		if len(args) == 1 && args[0] == "branch" {
			return "", true
		}
	}
	return "", false
}

func (r *fakeRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	key := runnerKey(dir, name, args...)
	r.calls = append(r.calls, key)
	if errs, ok := r.errors[key]; ok && len(errs) > 0 {
		err := errs[0]
		r.errors[key] = errs[1:]
		return "", err
	}
	if outs, ok := r.outputs[key]; ok && len(outs) > 0 {
		out := outs[0]
		if len(outs) > 1 {
			r.outputs[key] = outs[1:]
		}
		return out, nil
	}
	if out, ok := defaultFakeOutput(name, args); ok {
		return out, nil
	}
	return "", fmt.Errorf("unexpected command: %s", key)
}

func (r *fakeRunner) RunStdout(ctx context.Context, dir, name string, args ...string) (string, error) {
	key := runnerKey(dir, name, args...)
	r.calls = append(r.calls, key)
	if errs, ok := r.errors[key]; ok && len(errs) > 0 {
		err := errs[0]
		r.errors[key] = errs[1:]
		return "", err
	}
	if outs, ok := r.stdoutOutputs[key]; ok && len(outs) > 0 {
		out := outs[0]
		if len(outs) > 1 {
			r.stdoutOutputs[key] = outs[1:]
		}
		return out, nil
	}
	if outs, ok := r.outputs[key]; ok && len(outs) > 0 {
		out := outs[0]
		if len(outs) > 1 {
			r.outputs[key] = outs[1:]
		}
		return out, nil
	}
	if out, ok := defaultFakeOutput(name, args); ok {
		return out, nil
	}
	return "", fmt.Errorf("unexpected command: %s", key)
}

func (r *fakeRunner) RunStream(ctx context.Context, dir, name string, args ...string) error {
	_, err := r.Run(ctx, dir, name, args...)
	return err
}

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(value string) string {
	return ansiPattern.ReplaceAllString(value, "")
}

func TestBranchFromPushArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want *string
	}{
		{name: "no args", args: nil, want: nil},
		{name: "plain branch", args: []string{"feature/demo"}, want: ptr("feature/demo")},
		{name: "head refspec", args: []string{"HEAD:refs/heads/feature/demo"}, want: ptr("feature/demo")},
		{name: "head target", args: []string{"HEAD"}, want: nil},
		{name: "flag-like", args: []string{"--force"}, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := branchFromPushArgs(tt.args)
			switch {
			case got == nil && tt.want == nil:
				return
			case got == nil || tt.want == nil:
				t.Fatalf("branchFromPushArgs(%v) = %v, want %v", tt.args, got, tt.want)
			case *got != *tt.want:
				t.Fatalf("branchFromPushArgs(%v) = %q, want %q", tt.args, *got, *tt.want)
			}
		})
	}
}

func TestPromptForModifySelection(t *testing.T) {
	candidates := []ChangeInfo{
		{ChangeID: "abc123change", CommitID: "11111111commit", Description: "first"},
		{ChangeID: "def456change", CommitID: "22222222commit", Description: "second"},
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr string
	}{
		{name: "default first", input: "\n", want: "abc123change"},
		{name: "numeric selection", input: "2\n", want: "def456change"},
		{name: "change prefix", input: "def4\n", want: "def456change"},
		{name: "commit prefix", input: "2222\n", want: "def456change"},
		{name: "cancel", input: "q\n", wantErr: context.Canceled.Error()},
		{name: "out of range", input: "3\n", wantErr: "out of range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			got, err := PromptForModifySelection(strings.NewReader(tt.input), &out, candidates)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("PromptForModifySelection() error = %v, want substring %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("PromptForModifySelection() unexpected error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("PromptForModifySelection() = %q, want %q", got, tt.want)
			}
			if !strings.Contains(out.String(), "Select change to edit:") {
				t.Fatalf("prompt output missing header: %q", out.String())
			}
		})
	}
}

func TestParseGitHubRemote(t *testing.T) {
	tests := []struct {
		name   string
		remote string
		host   string
		owner  string
		repo   string
		ok     bool
	}{
		{
			name:   "https remote",
			remote: "https://github.com/satoricorp/gx.git",
			host:   "github.com",
			owner:  "satoricorp",
			repo:   "gx",
			ok:     true,
		},
		{
			name:   "ssh remote",
			remote: "git@github.com:satoricorp/gx.git",
			host:   "github.com",
			owner:  "satoricorp",
			repo:   "gx",
			ok:     true,
		},
		{
			name:   "enterprise remote",
			remote: "ssh://git@github.example.com/platform/gx.git",
			host:   "github.example.com",
			owner:  "platform",
			repo:   "gx",
			ok:     true,
		},
		{
			name:   "non github remote",
			remote: "https://gitlab.com/satoricorp/gx.git",
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, owner, repo, ok := parseGitHubRemote(tt.remote)
			if ok != tt.ok {
				t.Fatalf("parseGitHubRemote(%q) ok = %v, want %v", tt.remote, ok, tt.ok)
			}
			if !tt.ok {
				return
			}
			if host != tt.host || owner != tt.owner || repo != tt.repo {
				t.Fatalf("parseGitHubRemote(%q) = (%q, %q, %q), want (%q, %q, %q)", tt.remote, host, owner, repo, tt.host, tt.owner, tt.repo)
			}
		})
	}
}

func TestGitHubPullRequestURL(t *testing.T) {
	url := githubPullRequestURL("git@github.com:satoricorp/gx.git", "main", "feature/demo")
	if url == nil {
		t.Fatalf("githubPullRequestURL() returned nil")
	}
	want := "https://github.com/satoricorp/gx/compare/main...feature%2Fdemo?quick_pull=1"
	if *url != want {
		t.Fatalf("githubPullRequestURL() = %q, want %q", *url, want)
	}
}

func TestEnsureGitHubPullRequestReturnsExistingPR(t *testing.T) {
	repoRoot := t.TempDir()
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/repos/satoricorp/gx/pulls" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"html_url":"https://github.com/satoricorp/gx/pull/7"}]`))
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")
	runner := &fakeRunner{}
	svc := NewServiceWithRunner(runner)
	remoteURL := "git@github.com:satoricorp/gx.git"

	got, status, warnings := svc.ensureGitHubPullRequest(context.Background(), RepoInfo{RootPath: repoRoot, RemoteURL: &remoteURL}, StackInfo{
		Name:    "Demo stack",
		BaseRef: "main",
	}, "feature/demo", nil)
	if got == nil || *got != "https://github.com/satoricorp/gx/pull/7" {
		t.Fatalf("ensureGitHubPullRequest() = %v, want existing PR URL", got)
	}
	if status != "existing" || len(warnings) != 0 {
		t.Fatalf("status=%q warnings=%#v, want existing without warnings", status, warnings)
	}
	if gotAuth != "Bearer token-one" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("runner calls = %#v, want no gh shellout", runner.calls)
	}
}

func TestEnsureGitHubPullRequestReusesStoredPR(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{}
	svc := NewServiceWithRunner(runner)
	remoteURL := "git@github.com:satoricorp/gx.git"
	prURL := "https://github.com/satoricorp/gx/pull/42"

	got, status, warnings := svc.ensureGitHubPullRequest(context.Background(), RepoInfo{RootPath: repoRoot, RemoteURL: &remoteURL}, StackInfo{
		Name:        "Demo stack",
		BaseRef:     "main",
		GitHubPRURL: &prURL,
	}, "feature/demo", nil)
	if got == nil || *got != prURL {
		t.Fatalf("ensureGitHubPullRequest() = %v, want stored PR URL", got)
	}
	if status != "stored" || len(warnings) != 0 {
		t.Fatalf("status=%q warnings=%#v, want stored without warnings", status, warnings)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("runner calls = %#v, want none", runner.calls)
	}
}

func TestEnsureGitHubPullRequestCreatesPRWhenMissing(t *testing.T) {
	repoRoot := t.TempDir()
	var createPayload map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`[]`))
		case http.MethodPost:
			if err := json.NewDecoder(r.Body).Decode(&createPayload); err != nil {
				t.Fatalf("decode create payload: %v", err)
			}
			_, _ = w.Write([]byte(`{"html_url":"https://github.com/satoricorp/gx/pull/8"}`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "main"): {
				"abc123\trefs/heads/main\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	remoteURL := "git@github.com:satoricorp/gx.git"

	got, status, warnings := svc.ensureGitHubPullRequest(context.Background(), RepoInfo{RootPath: repoRoot, RemoteURL: &remoteURL}, StackInfo{
		Name:    "Demo stack",
		BaseRef: "origin/main",
	}, "feature/demo", []PushedChange{{
		Change: ChangeInfo{ChangeID: "abc123", Description: "Add publish flow"},
	}})
	if got == nil || *got != "https://github.com/satoricorp/gx/pull/8" {
		t.Fatalf("ensureGitHubPullRequest() = %v, want created PR URL", got)
	}
	if status != "created" || len(warnings) != 0 {
		t.Fatalf("status=%q warnings=%#v, want created without warnings", status, warnings)
	}
	if createPayload["base"] != "main" || createPayload["head"] != "feature/demo" || createPayload["title"] != "Demo stack" {
		t.Fatalf("create payload = %#v", createPayload)
	}
}

func TestEnsureGitHubPullRequestResolvesInternalBaseRefToPublicBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	var createPayload map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`[]`))
		case http.MethodPost:
			if err := json.NewDecoder(r.Body).Decode(&createPayload); err != nil {
				t.Fatalf("decode create payload: %v", err)
			}
			_, _ = w.Write([]byte(`{"html_url":"https://github.com/satoricorp/gx/pull/9"}`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "gx/edit", "--no-graph", "-T", "change_id"): {
				"basechange\n",
			},
			runnerKey(repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {
				"feature/recovered-compose-batch|basechange\n" +
					"gx/base|basechange\n" +
					"gx/edit|basechange\n",
			},
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "feature/recovered-compose-batch"): {
				"abc123\trefs/heads/feature/recovered-compose-batch\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	remoteURL := "git@github.com:satoricorp/gx.git"

	got, status, warnings := svc.ensureGitHubPullRequest(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		DefaultBranch: ptr("main"),
		RemoteURL:     &remoteURL,
	}, StackInfo{
		Name:    "Demo stack",
		BaseRef: "gx/edit",
	}, "feature/demo", nil)
	if got == nil || *got != "https://github.com/satoricorp/gx/pull/9" {
		t.Fatalf("ensureGitHubPullRequest() = %v, want created PR URL", got)
	}
	if status != "created" || len(warnings) != 0 {
		t.Fatalf("status=%q warnings=%#v, want created without warnings", status, warnings)
	}
	if createPayload["base"] != "feature/recovered-compose-batch" {
		t.Fatalf("create payload base = %q, want public bookmark; payload = %#v", createPayload["base"], createPayload)
	}
}

func TestEnsureGitHubPullRequestWarnsWhenBaseBranchMissing(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	var posted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`[]`))
		case http.MethodPost:
			posted = true
			_, _ = w.Write([]byte(`{"html_url":"https://github.com/satoricorp/gx/pull/10"}`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "gx-internal-checkout-refs"): {
				"",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	remoteURL := "git@github.com:satoricorp/gx.git"

	got, status, warnings := svc.ensureGitHubPullRequest(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		DefaultBranch: ptr("main"),
		RemoteURL:     &remoteURL,
	}, StackInfo{
		Name:    "Demo stack",
		BaseRef: "gx-internal-checkout-refs",
	}, "feature/gx-internal-checkout-refs", nil)
	if got != nil {
		t.Fatalf("ensureGitHubPullRequest() = %v, want nil", got)
	}
	if status != "warning" || len(warnings) != 1 || !strings.Contains(warnings[0], "base branch") {
		t.Fatalf("status=%q warnings=%#v, want missing base warning", status, warnings)
	}
	if posted {
		t.Fatal("ensureGitHubPullRequest() posted despite missing base branch")
	}
}

func TestEnsureGitHubPullRequestWarnsWhenGXStackBaseNotPublished(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	var posted bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`[]`))
		case http.MethodPost:
			posted = true
			_, _ = w.Write([]byte(`{"html_url":"https://github.com/satoricorp/gx/pull/11"}`))
		default:
			t.Fatalf("method = %s", r.Method)
		}
	}))
	defer server.Close()
	t.Setenv("GX_GITHUB_API_URL", server.URL)
	t.Setenv("GH_TOKEN", "token-one")
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "feature/base-stack"): {
				"",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	remoteURL := "git@github.com:satoricorp/gx.git"

	got, status, warnings := svc.ensureGitHubPullRequest(context.Background(), RepoInfo{
		RootPath:  repoRoot,
		RemoteURL: &remoteURL,
	}, StackInfo{
		Name:    "Demo stack",
		BaseRef: "feature/base-stack",
	}, "feature/demo-stack", nil)
	if got != nil {
		t.Fatalf("ensureGitHubPullRequest() = %v, want nil", got)
	}
	if status != "warning" {
		t.Fatalf("status=%q warnings=%#v, want warning", status, warnings)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "Run `gx publish feature/base-stack` first, then retry `gx publish feature/demo-stack`") {
		t.Fatalf("warnings=%#v, want stack publish guidance", warnings)
	}
	if posted {
		t.Fatal("ensureGitHubPullRequest() posted despite missing base branch")
	}
}

func TestEnsureGitHubPullRequestSkipsNonGitHubRemote(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{}
	svc := NewServiceWithRunner(runner)
	remoteURL := "https://example.com/satoricorp/gx.git"

	got, status, warnings := svc.ensureGitHubPullRequest(context.Background(), RepoInfo{RootPath: repoRoot, RemoteURL: &remoteURL}, StackInfo{
		Name:    "Demo stack",
		BaseRef: "main",
	}, "feature/demo", nil)
	if got != nil {
		t.Fatalf("ensureGitHubPullRequest() = %v, want nil", got)
	}
	if status != "skipped" || len(warnings) != 0 {
		t.Fatalf("status=%q warnings=%#v, want skipped without warnings", status, warnings)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("runner calls = %#v, want none", runner.calls)
	}
}

func TestPublishRefUsesExplicitArgOrStackName(t *testing.T) {
	if got := publishRefFromArgs([]string{"feature/login"}); got != "feature/login" {
		t.Fatalf("publishRefFromArgs() = %q, want feature/login", got)
	}
	body := StackInfo{Name: "Login Flow", BookmarkName: "feature/login-flow"}
	if got := publishRefForStack(body); got != "feature/login-flow" {
		t.Fatalf("publishRefForStack() = %q, want feature/login-flow", got)
	}
}

func TestSwitchUsesStoredStackBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "chg123"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "login",
		BookmarkName: "login",
		BaseRef:      "main",
		BaseCommitID: "base123",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"): {cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"):      {""},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {""},
			runnerKey(cwd, "jj", "edit", "chg123"):                                           {"Working copy now at chg123\n"},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"):              {"false\n", "false\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "login", "-r", "@", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "commit_id"):          {"abc123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/login", "abc123"):                {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/login"):                {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                {""},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"chg123|abc123|login|base\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@", "--name-only"): {""},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey(cwd, "git", "branch", "--show-current"):                            {fmt.Errorf("detached")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)
	result, err := svc.Switch(context.Background(), "login")
	if err != nil {
		t.Fatalf("Switch() unexpected error = %v", err)
	}
	if result.Stack.BookmarkName != "login" || result.CurrentChange.ChangeID != "chg123" {
		t.Fatalf("Switch() = %#v", result)
	}
}

func TestModifyReattachesVisibleGitHeadToStackBranch(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "targetchg"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "feature",
		BookmarkName: "feature/feature",
		BaseRef:      "main",
		BaseCommitID: "base123",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	changeTemplate := `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"): {cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {
				"feature/feature\n",
				"",
				"",
				"feature/feature\n",
			},
			runnerKey(cwd, "git", "rev-parse", "--git-path", "info/exclude"): {filepath.Join(cwd, ".git", "info", "exclude") + "\n"},
			runnerKey(cwd, "jj", "config", "get", "user.name"):               {"Joe Example\n"},
			runnerKey(cwd, "jj", "config", "get", "user.email"):              {"joe@example.com\n"},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {
				"false\n",
				"false\n",
				"false\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", changeTemplate): {
				"beforechg|beforecommit|before|basechg\n",
				"beforechg|beforecommit|before|basechg\n",
				"targetchg|targetcommit|target|basechg\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@", "--name-only"):                                     {""},
			runnerKey(cwd, "jj", "edit", "targetchg"):                                                  {"Working copy now at targetchg\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feature", "-r", "@", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(feature/feature)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"targetchg|targetcommit|target|basechg\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "commit_id"):           {"targetcommit\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feature", "targetcommit"): {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feature"):       {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                 {""},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"): {
				"op123\n",
			},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head")},
		},
	}
	svc := NewServiceWithRunner(runner)
	result, err := svc.Modify(context.Background(), "targetchg")
	if err != nil {
		t.Fatalf("Modify() unexpected error = %v\ncalls: %v", err, runner.calls)
	}
	if result.CurrentChange.ChangeID != "targetchg" {
		t.Fatalf("Modify() current change = %q, want targetchg", result.CurrentChange.ChangeID)
	}
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feature", "targetcommit"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feature"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "reset", "--mixed", "HEAD"))
	assertRunnerNotCalled(t, runner.calls, runnerKey(cwd, "git", "switch", "-f", "main"))
}

func TestResolveStoredRevisionRefUsesCommitForStoredChange(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	repoRoot := t.TempDir()
	ctx := context.Background()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	defer store.Close()
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if _, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "targetchg",
		CurrentCommitID: "targetcommit",
		Description:     "target",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	}); err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}

	svc := NewServiceWithRunner(&fakeRunner{})
	got, err := svc.resolveStoredRevisionRef(ctx, RepoInfo{RootPath: repoRoot, Backend: "jj"}, "targetchg")
	if err != nil {
		t.Fatalf("resolveStoredRevisionRef() error = %v", err)
	}
	if got != "targetcommit" {
		t.Fatalf("resolveStoredRevisionRef() = %q, want targetcommit", got)
	}
}

func TestEditRevisionReattachesVisibleGitHeadToStackBranch(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "targetchg"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "feature",
		BookmarkName: "feature/feature",
		BaseRef:      "main",
		BaseCommitID: "base123",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "edit", "targetchg"): {"Working copy now at targetchg\n"},
			runnerKey(cwd, "jj", "root"):              {cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {
				"",
				"",
				"",
				"feature/feature\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {
				"false\n",
				"false\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"targetchg|targetcommit|target|basechg\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@", "--name-only"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "change_id"): {
				"targetchg\n",
			},
			runnerKey(cwd, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {
				"feature/feature|targetchg\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "feature/feature", "--no-graph", "-T", "commit_id"): {"targetcommit\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feature", "targetcommit"):     {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feature"):           {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                     {""},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head")},
		},
	}
	svc := NewServiceWithRunner(runner)
	if err := svc.EditRevision(context.Background(), cwd, "targetchg"); err != nil {
		t.Fatalf("EditRevision() unexpected error = %v\ncalls: %v", err, runner.calls)
	}
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feature", "targetcommit"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feature"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "reset", "--mixed", "HEAD"))
}

func TestBaseResultTreatsEquivalentRefsAsOnBase(t *testing.T) {
	repoRoot := t.TempDir()
	base := "codex/gx-desktop-local-data-auth"
	current := "main"
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", current, "--no-graph", "-T", "change_id"): {"same-change\n"},
			runnerKey(repoRoot, "jj", "log", "-r", base, "--no-graph", "-T", "change_id"):    {"same-change\n"},
		},
	}
	svc := NewServiceWithRunner(runner)

	result := svc.baseResult(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		BranchName:    &current,
		AuthoringBase: &base,
	})

	if !result.OnBase {
		t.Fatalf("baseResult().OnBase = false for equivalent refs")
	}
}

func TestDeleteRevisionAbandonsJJChangeAndMarksMetadata(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	changeID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      "chg-delete",
		CurrentCommitID: "commit-delete",
		Description:     "delete me",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                     {cwd + "\n"},
			runnerKey(cwd, "jj", "abandon", "commit-delete"): {"Abandoned\n"},
		},
		stdoutOutputs: map[string][]string{
			runnerKey(cwd, "jj", "log", "-r", "commit-delete", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {"chg-delete|commit-delete|delete me|\n"},
			runnerKey(cwd, "jj", "diff", "-r", "commit-delete", "--name-only"): {"a.txt\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "branch", "--show-current"):                            {fmt.Errorf("detached")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head")},
		},
	}
	svc := NewServiceWithRunner(runner)
	result, err := svc.DeleteRevision(context.Background(), "chg-delete")
	if err != nil {
		t.Fatalf("DeleteRevision() error = %v", err)
	}
	if result.Change.ChangeID != "chg-delete" || result.Output != "Abandoned\n" {
		t.Fatalf("DeleteRevision() = %#v", result)
	}

	db, err = storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() reopen error = %v", err)
	}
	store, err = storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() reopen error = %v", err)
	}
	defer store.Close()
	change, err := store.FindChangeByJJChangeID(context.Background(), repoID, "chg-delete")
	if err != nil {
		t.Fatalf("FindChangeByJJChangeID() error = %v", err)
	}
	if change == nil || change.ID != changeID || change.Status != "abandoned" {
		t.Fatalf("stored change = %#v, want abandoned id %d", change, changeID)
	}
}

func TestDeleteStackAbandonsRevisionsDeletesBookmarkAndMetadata(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "chg-two"
	stackID, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "docs",
		BookmarkName: "feature/docs",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	firstID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      "chg-one",
		CurrentCommitID: "commit-one",
		Description:     "one",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange(one) error = %v", err)
	}
	secondID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      "chg-two",
		CurrentCommitID: "commit-two",
		Description:     "two",
		Status:          "draft",
		FirstSeenAt:     2,
		UpdatedAt:       2,
	})
	if err != nil {
		t.Fatalf("UpsertChange(two) error = %v", err)
	}
	if err := store.AddChangeToStack(context.Background(), stackID, firstID, 1); err != nil {
		t.Fatalf("AddChangeToStack(one) error = %v", err)
	}
	if err := store.AddChangeToStack(context.Background(), stackID, secondID, 2); err != nil {
		t.Fatalf("AddChangeToStack(two) error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                               {cwd + "\n"},
			runnerKey(cwd, "jj", "abandon", "chg-one"):                 {"Abandoned one\n"},
			runnerKey(cwd, "jj", "abandon", "chg-two"):                 {"Abandoned two\n"},
			runnerKey(cwd, "jj", "bookmark", "delete", "feature/docs"): {"Deleted bookmark\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "branch", "--show-current"):                            {fmt.Errorf("detached")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head")},
		},
	}
	svc := NewServiceWithRunner(runner)
	result, err := svc.DeleteStack(context.Background(), "feature/docs")
	if err != nil {
		t.Fatalf("DeleteStack() error = %v", err)
	}
	if result.Stack.BookmarkName != "feature/docs" || len(result.Revisions) != 2 {
		t.Fatalf("DeleteStack() = %#v", result)
	}

	db, err = storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() reopen error = %v", err)
	}
	store, err = storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() reopen error = %v", err)
	}
	defer store.Close()
	if stack, err := store.FindStackByBookmark(context.Background(), repoID, "feature/docs"); err != nil {
		t.Fatalf("FindStackByBookmark() error = %v", err)
	} else if stack != nil {
		t.Fatalf("stack still exists: %#v", stack)
	}
	for _, changeID := range []string{"chg-one", "chg-two"} {
		change, err := store.FindChangeByJJChangeID(context.Background(), repoID, changeID)
		if err != nil {
			t.Fatalf("FindChangeByJJChangeID(%s) error = %v", changeID, err)
		}
		if change == nil || change.Status != "abandoned" {
			t.Fatalf("change %s = %#v, want abandoned", changeID, change)
		}
	}
}

func TestSwitchForksMutableChangeWhenImmutable(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "published123"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "fork",
		BookmarkName: "feature/fork",
		BaseRef:      "main",
		BaseCommitID: "base123",
		HeadChangeID: &head,
		Status:       "published",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	immutableErr := fmt.Errorf("jj edit: exit status 1\nCommit abc is immutable")
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"): {cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"):              {""},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"):         {""},
			runnerKey(cwd, "jj", "new", "feature/fork"):                                              {"Working copy now at new456\n"},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"):                      {"true\n", "true\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/fork", "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "commit_id"):                  {"abc123\n"},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):                 {"abc123\n"},
			runnerKey(cwd, "jj", "log", "-r", "feature/fork", "--no-graph", "-T", "commit_id"):       {"abc123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/fork", "abc123"):                 {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/fork"):                 {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                        {""},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"new456|abc123|fork|published123\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@", "--name-only"): {""},
		},
		errors: map[string][]error{
			runnerKey(cwd, "jj", "edit", "feature/fork"):                                 {immutableErr},
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey(cwd, "git", "branch", "--show-current"):                            {fmt.Errorf("detached")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)
	result, err := svc.Switch(context.Background(), "s1")
	if err != nil {
		t.Fatalf("Switch() unexpected error = %v", err)
	}
	if result.CurrentChange.ChangeID != "new456" {
		t.Fatalf("Switch() change = %q, want new456", result.CurrentChange.ChangeID)
	}
	for _, got := range runner.calls {
		if got == runnerKey(cwd, "jj", "new", "feature/fork") {
			return
		}
	}
	t.Fatalf("Switch() missing jj new fallback in %v", runner.calls)
}

func TestSwitchBaseCreatesMutableChildOfBase(t *testing.T) {
	repoRoot := t.TempDir()
	main := "main"
	current := "feature/current-stack"
	tmpl := `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "new", "main"):                                        {"Working copy now at child456\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "main", "--no-graph", "-T", "commit_id"): {"base123\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"):        {"true\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):   {"base123\n", "base123\n"},
			runnerKey(repoRoot, "git", "branch", "--show-current"):                          {"main\n"},
			runnerKey(repoRoot, "git", "update-ref", "refs/heads/main", "base123"):          {""},
			runnerKey(repoRoot, "git", "symbolic-ref", "HEAD", "refs/heads/main"):           {""},
			runnerKey(repoRoot, "git", "reset", "--mixed", "HEAD"):                          {""},
			runnerKey(repoRoot, "jj", "root"):                                               {repoRoot + "\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", tmpl): {
				"child456|work123||basechange\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "@", "--name-only"): {""},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.switchBaseUnlocked(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		DefaultBranch: &main,
		AuthoringBase: &main,
		BranchName:    &current,
	}, "main")
	if err != nil {
		t.Fatalf("switchBaseUnlocked() unexpected error = %v", err)
	}
	if result.CurrentChange.ChangeID != "child456" {
		t.Fatalf("switchBaseUnlocked() change = %q, want child456", result.CurrentChange.ChangeID)
	}
	if result.Stack.BookmarkName != "main" || result.Stack.BaseCommitID != "base123" {
		t.Fatalf("switchBaseUnlocked() stack = %#v", result.Stack)
	}
	if result.Repo.BranchName == nil || *result.Repo.BranchName != "main" {
		t.Fatalf("switchBaseUnlocked() branch = %#v, want main", result.Repo.BranchName)
	}
	assertRunnerCalled(t, runner.calls, runnerKey(repoRoot, "git", "update-ref", "refs/heads/main", "base123"))
	assertRunnerCalled(t, runner.calls, runnerKey(repoRoot, "git", "symbolic-ref", "HEAD", "refs/heads/main"))
	assertRunnerCalled(t, runner.calls, runnerKey(repoRoot, "git", "reset", "--mixed", "HEAD"))
	for _, got := range runner.calls {
		if strings.Contains(got, "|jj|edit ") {
			t.Fatalf("switchBaseUnlocked() called jj edit: %v", runner.calls)
		}
	}
}

func TestStackForRevisionUsesRevisionOwner(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	baseHead := "basechange"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "Compact AI review brief payloads",
		BookmarkName: "feature/compact-ai-review-brief-payloads",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &baseHead,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack(base) error = %v", err)
	}
	editHead := "editchange"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "edits",
		BookmarkName: "feature/edits",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &editHead,
		Status:       "draft",
		CreatedAt:    2,
		UpdatedAt:    2,
	}); err != nil {
		t.Fatalf("UpsertStack(edit) error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "zprnwwy", "--no-graph", "-T", "change_id"): {"editchange\n"},
			runnerKey(repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {
				"feature/compact-ai-review-brief-payloads|basechange\n" +
					"feature/edits|editchange\n" +
					"main|basechange\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	got, found, err := svc.stackForRevision(context.Background(), RepoInfo{
		RootPath: repoRoot,
		Backend:  "jj",
	}, "zprnwwy")
	if err != nil {
		t.Fatalf("stackForRevision() error = %v", err)
	}
	if !found {
		t.Fatal("stackForRevision() found = false")
	}
	if got.BookmarkName != "feature/edits" {
		t.Fatalf("stackForRevision() bookmark = %q, want feature/edits", got.BookmarkName)
	}
}

func TestNewRevisionFromCreatesMutableChildOfRevision(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "new", "main"): {"Working copy now at child\n"},
		},
	}
	svc := NewServiceWithRunner(runner)

	if err := svc.NewRevisionFrom(context.Background(), repoRoot, "main"); err != nil {
		t.Fatalf("NewRevisionFrom() error = %v", err)
	}

	for _, got := range runner.calls {
		if got == runnerKey(repoRoot, "jj", "new", "main") {
			return
		}
	}
	t.Fatalf("NewRevisionFrom() calls = %#v, want jj new main", runner.calls)
}

func TestResolveCurrentStackCreatesNewStackFromBaseBranch(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	oldHead := "oldchange"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "old",
		BookmarkName: "feature/old",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &oldHead,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"false\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"newchange|newcommit||base\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "@", "--name-only"):                        {"new.txt\n"},
			runnerKey(repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl):            {""},
			runnerKey(repoRoot, "git", "rev-parse", "main"):                                    {"basecommit\n"},
			runnerKey(repoRoot, "jj", "bookmark", "set", "feature/change-newchang", "-r", "@"): {""},
		},
	}
	svc := NewServiceWithRunner(runner)
	body, err := svc.resolveCurrentStack(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		BranchName:    ptr("main"),
	}, true)
	if err != nil {
		t.Fatalf("resolveCurrentStack() unexpected error = %v", err)
	}
	if body.BookmarkName != "feature/change-newchang" {
		t.Fatalf("resolveCurrentStack() bookmark = %q, want new stack", body.BookmarkName)
	}
}

func TestLoadStoredStackInfosHidesEmptyRevisions(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	ctx := context.Background()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	defer store.Close()
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "empty-change"
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "docs",
		BookmarkName: "feature/docs",
		BaseRef:      "",
		BaseCommitID: "base",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	nonEmptyID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "non-empty-change",
		CurrentCommitID: "non-empty-commit",
		Description:     "real revision",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange(non-empty) error = %v", err)
	}
	emptyID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "empty-change",
		CurrentCommitID: "empty-commit",
		Description:     "empty revision",
		Status:          "draft",
		FirstSeenAt:     2,
		UpdatedAt:       2,
	})
	if err != nil {
		t.Fatalf("UpsertChange(empty) error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, nonEmptyID, 1); err != nil {
		t.Fatalf("AddChangeToStack(non-empty) error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, emptyID, 2); err != nil {
		t.Fatalf("AddChangeToStack(empty) error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {""},
			runnerKey(repoRoot, "jj", "log", "-r", "non-empty-commit", "--no-graph", "-T", "empty"): {
				"false\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "empty-commit", "--no-graph", "-T", "empty"): {
				"true\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	stacks, err := svc.loadStoredStackInfos(ctx, store, RepoInfo{RootPath: repoRoot, Backend: "jj"}, repoID)
	if err != nil {
		t.Fatalf("loadStoredStackInfos() error = %v", err)
	}
	if len(stacks) != 1 {
		t.Fatalf("stored stacks = %d, want 1", len(stacks))
	}
	if len(stacks[0].Revisions) != 1 {
		t.Fatalf("stored revisions = %#v, want only non-empty revision", stacks[0].Revisions)
	}
	if stacks[0].Revisions[0].ChangeID != "non-empty-change" {
		t.Fatalf("revision change = %q, want non-empty-change", stacks[0].Revisions[0].ChangeID)
	}
}

func TestResolveCurrentStackIgnoresLatestCursorGXRequest(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	oldHead := "oldchange"
	oldStackID, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "old",
		BookmarkName: "feature/old",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &oldHead,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    10,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	oldChangeID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      oldHead,
		CurrentCommitID: "oldcommit",
		Description:     "old",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(context.Background(), oldStackID, oldChangeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:        "cursor-old",
		CreatedAt: 1,
		Command:   "cursor",
		Cwd:       repoRoot,
		GXVersion: "test",
		Source:    ptr("cursor"),
		RepoRoot:  ptr(repoRoot),
	}); err != nil {
		t.Fatalf("WriteSession(cursor-old) error = %v", err)
	}
	if err := store.WriteChangeSessions(context.Background(), oldChangeID, []string{"cursor-old"}, 1); err != nil {
		t.Fatalf("WriteChangeSessions() error = %v", err)
	}
	lastSeen := int64(20)
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:         "cursor-current",
		CreatedAt:  20,
		Command:    "cursor",
		Cwd:        repoRoot,
		GXVersion:  "test",
		Source:     ptr("cursor"),
		LastSeenAt: &lastSeen,
		RepoRoot:   ptr(repoRoot),
	}); err != nil {
		t.Fatalf("WriteSession(cursor-current) error = %v", err)
	}
	if _, err := store.UpsertCursorMessage(context.Background(), storage.CursorMessage{
		ID:        "msg-current",
		SessionID: "cursor-current",
		CreatedAt: 21,
		Role:      "user",
		Text:      "save using gx",
		RawJSON:   []byte(`{}`),
	}); err != nil {
		t.Fatalf("UpsertCursorMessage() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	bookmark := "feature/old"
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"false\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"newchange|newcommit||base\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "@", "--name-only"): {"new.txt\n"},
		},
	}
	svc := NewServiceWithRunner(runner)
	body, err := svc.resolveCurrentStackWithDescription(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		BranchName:    ptr("feature/old"),
	}, true, "feat one")
	if err != nil {
		t.Fatalf("resolveCurrentStackWithDescription() unexpected error = %v", err)
	}
	if body.BookmarkName != bookmark {
		t.Fatalf("resolveCurrentStackWithDescription() bookmark = %q, want %q", body.BookmarkName, bookmark)
	}
	if body.BaseRef != "main" {
		t.Fatalf("resolveCurrentStackWithDescription() base = %q, want main", body.BaseRef)
	}
}

func TestResolveCurrentStackForPublishPrefersAttachedStackOverCursorSession(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	workHead := "workchange"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "document-gx-pr",
		BookmarkName: "feature/document-gx-pr",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &workHead,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    20,
	}); err != nil {
		t.Fatalf("UpsertStack(work) error = %v", err)
	}
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:        "cursor-current",
		CreatedAt: 20,
		Command:   "cursor",
		Cwd:       repoRoot,
		GXVersion: "test",
		Source:    ptr("cursor"),
		RepoRoot:  ptr(repoRoot),
	}); err != nil {
		t.Fatalf("WriteSession() error = %v", err)
	}
	if _, err := store.UpsertCursorMessage(context.Background(), storage.CursorMessage{
		ID:        "msg-current",
		SessionID: "cursor-current",
		CreatedAt: 21,
		Role:      "user",
		Text:      "save using gx",
		RawJSON:   []byte(`{}`),
	}); err != nil {
		t.Fatalf("UpsertCursorMessage() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"):  {"true\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "@-", "--no-graph", "-T", "empty"): {"false\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				workHead + "|workcommit||base\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "@-", "--name-only"): {""},
			runnerKey(repoRoot, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {
				"feature/document-gx-pr|" + workHead + "\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	body, err := svc.resolveCurrentStackForPublish(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		BranchName:    ptr("feature/document-gx-pr"),
	})
	if err != nil {
		t.Fatalf("resolveCurrentStackForPublish() unexpected error = %v", err)
	}
	if body.BookmarkName != "feature/document-gx-pr" {
		t.Fatalf("resolveCurrentStackForPublish() bookmark = %q, want feature/document-gx-pr", body.BookmarkName)
	}
}

func TestResolveCurrentStackForReadIgnoresLegacyBaseCheckout(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "parentchange"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "edits",
		BookmarkName: "feature/edits",
		BaseRef:      "main",
		BaseCommitID: "basecommit",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}
	current := "main"
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", current, "--no-graph", "-T", "change_id"): {"basechange\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "main", "--no-graph", "-T", "change_id"):  {"basechange\n"},
		},
	}
	svc := NewServiceWithRunner(runner)

	_, found, err := svc.resolveCurrentStackForRead(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		BranchName:    &current,
	})
	if err != nil {
		t.Fatalf("resolveCurrentStackForRead() error = %v", err)
	}
	if found {
		t.Fatal("resolveCurrentStackForRead() found stack on legacy gx/base checkout")
	}
}

func TestResolveCurrentStackForReadDoesNotSyncOrRouteByCursorRequest(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	oldHead := "oldchange"
	oldStackID, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "old",
		BookmarkName: "feature/old",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &oldHead,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    10,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	oldChangeID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      oldHead,
		CurrentCommitID: "oldcommit",
		Description:     "old",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(context.Background(), oldStackID, oldChangeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:        "cursor-current",
		CreatedAt: 20,
		Command:   "cursor",
		Cwd:       repoRoot,
		GXVersion: "test",
		Source:    ptr("cursor"),
		RepoRoot:  ptr(repoRoot),
	}); err != nil {
		t.Fatalf("WriteSession() error = %v", err)
	}
	if _, err := store.UpsertCursorMessage(context.Background(), storage.CursorMessage{
		ID:        "msg-current",
		SessionID: "cursor-current",
		CreatedAt: 21,
		Role:      "user",
		Text:      "save using gx",
		RawJSON:   []byte(`{}`),
	}); err != nil {
		t.Fatalf("UpsertCursorMessage() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"newchange|newcommit||base\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	body, found, err := svc.resolveCurrentStackForRead(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
	})
	if err != nil {
		t.Fatalf("resolveCurrentStackForRead() unexpected error = %v", err)
	}
	if !found {
		t.Fatal("resolveCurrentStackForRead() found = false, want true")
	}
	if body.BookmarkName != "feature/old" {
		t.Fatalf("resolveCurrentStackForRead() bookmark = %q, want feature/old", body.BookmarkName)
	}
}

func TestEnsureRepoAtPathUsesExistingJJRepo(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey("/repo", "jj", "root"):                                        {"/repo\n", "/repo\n"},
			runnerKey("/repo", "git", "branch", "--show-current"):                   {"main\n", "main\n"},
			runnerKey("/repo", "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"false\n"},
			runnerKey("/repo", "jj", "bookmark", "set", "main", "-r", "@"):          {""},
		},
		errors: map[string][]error{
			runnerKey("/repo", "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey("/repo", "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey("/repo", "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey("/repo", "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "user.name"):                                 {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "--global", "user.name"):                     {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "user.email"):                                {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "--global", "user.email"):                    {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.EnsureRepoAtPath(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("EnsureRepoAtPath() unexpected error = %v", err)
	}
	if result.Initialized {
		t.Fatalf("EnsureRepoAtPath() initialized existing jj repo")
	}
	if result.Repo.RootPath != "/repo" {
		t.Fatalf("EnsureRepoAtPath() root = %q, want /repo", result.Repo.RootPath)
	}
	store, err := openStore(context.Background())
	if err != nil {
		t.Fatalf("openStore() error = %v", err)
	}
	defer store.Close()
	repo, err := store.FindRepoByRoot(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("FindRepoByRoot() error = %v", err)
	}
	if repo == nil {
		t.Fatal("expected gx init to register repo in GX storage")
	}
	initializedRepos, err := store.ListInitializedRepos(context.Background())
	if err != nil {
		t.Fatalf("ListInitializedRepos() error = %v", err)
	}
	if len(initializedRepos) != 1 || initializedRepos[0].RootPath != "/repo" {
		t.Fatalf("initialized repos = %#v, want /repo", initializedRepos)
	}
	for _, call := range runner.calls {
		if strings.Contains(call, "jj|git init .") {
			t.Fatalf("EnsureRepoAtPath() unexpectedly initialized repo: %v", runner.calls)
		}
	}
}

func TestEnsureRepoAtPathInitializesAtGitRoot(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey("/repo/subdir", "git", "rev-parse", "--show-toplevel"):        {"/repo\n"},
			runnerKey("/repo", "jj", "git", "init", "."):                            {"Initialized repo in \".\"\n"},
			runnerKey("/repo", "jj", "root"):                                        {"/repo\n", "/repo\n"},
			runnerKey("/repo", "git", "branch", "--show-current"):                   {"main\n", "main\n"},
			runnerKey("/repo", "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"false\n"},
			runnerKey("/repo", "jj", "bookmark", "set", "main", "-r", "@"):          {""},
		},
		errors: map[string][]error{
			runnerKey("/repo/subdir", "jj", "root"):                                          {fmt.Errorf("not a jj repo")},
			runnerKey("/repo", "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey("/repo", "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey("/repo", "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey("/repo", "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "user.name"):                                 {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "--global", "user.name"):                     {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "user.email"):                                {fmt.Errorf("missing")},
			runnerKey("/repo", "git", "config", "--global", "user.email"):                    {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.EnsureRepoAtPath(context.Background(), "/repo/subdir")
	if err != nil {
		t.Fatalf("EnsureRepoAtPath() unexpected error = %v", err)
	}
	if !result.Initialized {
		t.Fatalf("EnsureRepoAtPath() did not initialize git-backed repo")
	}
	if result.Repo.RootPath != "/repo" {
		t.Fatalf("EnsureRepoAtPath() root = %q, want /repo", result.Repo.RootPath)
	}
	if len(runner.calls) < 2 || runner.calls[1] != runnerKey("/repo/subdir", "git", "rev-parse", "--show-toplevel") {
		t.Fatalf("EnsureRepoAtPath() did not inspect git root first: %v", runner.calls)
	}
}

func TestEnsureRepoAtPathInitializesCurrentDirectoryWhenNoRepoExists(t *testing.T) {
	t.Setenv("GX_HOME", t.TempDir())
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey("/work", "jj", "git", "init", "."):                            {"Initialized repo in \".\"\n"},
			runnerKey("/work", "jj", "root"):                                        {"/work\n", "/work\n"},
			runnerKey("/work", "git", "branch", "--show-current"):                   {"main\n", "main\n"},
			runnerKey("/work", "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"false\n"},
			runnerKey("/work", "jj", "bookmark", "set", "main", "-r", "@"):          {""},
		},
		errors: map[string][]error{
			runnerKey("/work", "jj", "root"):                                                 {fmt.Errorf("not a jj repo")},
			runnerKey("/work", "git", "rev-parse", "--show-toplevel"):                        {fmt.Errorf("not a git repo")},
			runnerKey("/work", "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey("/work", "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey("/work", "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey("/work", "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
			runnerKey("/work", "git", "config", "user.name"):                                 {fmt.Errorf("missing")},
			runnerKey("/work", "git", "config", "--global", "user.name"):                     {fmt.Errorf("missing")},
			runnerKey("/work", "git", "config", "user.email"):                                {fmt.Errorf("missing")},
			runnerKey("/work", "git", "config", "--global", "user.email"):                    {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.EnsureRepoAtPath(context.Background(), "/work")
	if err != nil {
		t.Fatalf("EnsureRepoAtPath() unexpected error = %v", err)
	}
	if !result.Initialized {
		t.Fatalf("EnsureRepoAtPath() did not initialize current directory")
	}
	if result.Repo.RootPath != "/work" {
		t.Fatalf("EnsureRepoAtPath() root = %q, want /work", result.Repo.RootPath)
	}
}

func TestEnsureIdentityUsesGitConfigAndPersistsGXConfig(t *testing.T) {
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	repoRoot := t.TempDir()

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "git", "config", "user.name"):                                     {"Joe Example\n"},
			runnerKey(repoRoot, "git", "config", "user.email"):                                    {"joe@example.com\n"},
			runnerKey(repoRoot, "jj", "config", "set", "--user", "user.name", "Joe Example"):      {""},
			runnerKey(repoRoot, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {""},
		},
		errors: map[string][]error{
			runnerKey(repoRoot, "jj", "config", "get", "user.name"):  {fmt.Errorf("missing")},
			runnerKey(repoRoot, "jj", "config", "get", "user.email"): {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	name, email, changed, err := svc.ensureIdentity(context.Background(), repoRoot, InitOptions{})
	if err != nil {
		t.Fatalf("ensureIdentity() unexpected error = %v", err)
	}
	if !changed {
		t.Fatalf("ensureIdentity() did not report changes")
	}
	if name != "Joe Example" || email != "joe@example.com" {
		t.Fatalf("ensureIdentity() = %q <%s>", name, email)
	}
	data, err := os.ReadFile(gxHome + "/config.json")
	if err != nil {
		t.Fatalf("read gx config: %v", err)
	}
	text := string(data)
	if !strings.Contains(text, "Joe Example") || !strings.Contains(text, "joe@example.com") {
		t.Fatalf("gx config missing identity: %s", text)
	}
}

func TestCommitUpdatesStackBookmarkAndAttachesGitBranch(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                      {cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {"main\n", "", "main\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"): {
				"",
			},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {
				"",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"true\n", "true\n", "true\n"},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"chg123|abc123|feat one|parent1\n",
				"chg123|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@-", "--name-only"):                                      {"a.txt\n", "a.txt\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-"):                      {""},
			runnerKey(cwd, "jj", "commit", "-m", "feat one"):                                             {"Committed\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):                     {"abc123\n"},
			runnerKey(cwd, "jj", "log", "-r", "feature/feat-one", "--no-graph", "-T", "commit_id"):       {"abc123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feat-one", "abc123"):                 {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feat-one"):                 {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                            {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(feature/feat-one)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"chg123|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "chg123", "-n", "1", "--no-graph", "-T", "hidden"): {"false\n"},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"):               {"op123\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.Commit(context.Background(), "feat one")
	if err != nil {
		t.Fatalf("Commit() unexpected error = %v", err)
	}
	if result.Change.CommitID != "abc123" {
		t.Fatalf("Commit() commit id = %q", result.Change.CommitID)
	}
	wantBookmarkUpdate := runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-", "--allow-backwards")
	assertRunnerCalled(t, runner.calls, wantBookmarkUpdate)
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feat-one", "abc123"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feat-one"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "reset", "--mixed", "HEAD"))
}

func TestSplitCommitUsesJJSplitAndRecordsSelectedChange(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                      {cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {"main\n", "", "main\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"): {
				"",
			},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {
				"",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"true\n", "true\n", "true\n"},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"chg123|abc123|feat selected|parent1\n",
				"chg123|abc123|feat selected|parent1\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@-", "--name-only"):                                           {"src/app.ts\n", "src/app.ts\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-selected", "-r", "@-"):                      {""},
			runnerKey(cwd, "jj", "split", "-m", "feat selected", "src/app.ts"):                                {"Split\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-selected", "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):                          {"abc123\n"},
			runnerKey(cwd, "jj", "log", "-r", "feature/feat-selected", "--no-graph", "-T", "commit_id"):       {"abc123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feat-selected", "abc123"):                 {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feat-selected"):                 {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                                 {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(feature/feat-selected)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"chg123|abc123|feat selected|parent1\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "chg123", "-n", "1", "--no-graph", "-T", "hidden"): {"false\n"},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"):               {"op123\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.SplitCommit(context.Background(), SplitCommitOptions{
		Message:  "feat selected",
		Filesets: []string{"src/app.ts"},
	})
	if err != nil {
		t.Fatalf("SplitCommit() unexpected error = %v", err)
	}
	if result.Change.Description != "feat selected" {
		t.Fatalf("SplitCommit() description = %q", result.Change.Description)
	}
	want := runnerKey(cwd, "jj", "split", "-m", "feat selected", "src/app.ts")
	splitIndex := -1
	selectedReadIndex := -1
	reattachIndex := -1
	for index, got := range runner.calls {
		if got == want {
			splitIndex = index
		}
	}
	if splitIndex == -1 {
		t.Fatalf("SplitCommit() missing call %q in %v", want, runner.calls)
	}
	for index, got := range runner.calls {
		switch got {
		case runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`):
			if index > splitIndex && selectedReadIndex == -1 {
				selectedReadIndex = index
			}
		case runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-selected", "-r", "@-", "--allow-backwards"):
			if index > splitIndex {
				reattachIndex = index
			}
		}
	}
	if selectedReadIndex == -1 {
		t.Fatalf("SplitCommit() did not read selected @- change before recording")
	}
	if reattachIndex == -1 {
		t.Fatalf("SplitCommit() did not reattach stack")
	}
	if !(splitIndex < selectedReadIndex && selectedReadIndex < reattachIndex) {
		t.Fatalf("SplitCommit() selected-change read order = split:%d selected:%d reattach:%d; want split < selected < reattach\ncalls: %v", splitIndex, selectedReadIndex, reattachIndex, runner.calls)
	}
}

func TestRecordRevisionCommitsAndPersistsMetadata(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                      {cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {"main\n", "", "main\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"): {
				"",
			},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {
				"",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"true\n", "true\n", "true\n"},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"chg123|abc123|feat one|parent1\n",
				"chg123|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@-", "--name-only"):                                      {"a.txt\n", "a.txt\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-"):                      {""},
			runnerKey(cwd, "jj", "commit", "-m", "feat one"):                                             {"Committed\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):                     {"abc123\n"},
			runnerKey(cwd, "jj", "log", "-r", "feature/feat-one", "--no-graph", "-T", "commit_id"):       {"abc123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feat-one", "abc123"):                 {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feat-one"):                 {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                            {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(feature/feat-one)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"chg123|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "chg123", "-n", "1", "--no-graph", "-T", "hidden"): {"false\n"},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"):               {"op123\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.RecordRevision(context.Background(), RevisionOptions{Message: "feat one"})
	if err != nil {
		t.Fatalf("RecordRevision() unexpected error = %v", err)
	}
	if result.Change.CommitID != "abc123" {
		t.Fatalf("RecordRevision() commit id = %q", result.Change.CommitID)
	}

	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	defer db.Close()

	var description string
	var bookmark string
	err = db.QueryRowContext(context.Background(), `
		SELECT c.description, b.bookmark_name
		FROM changes c
		JOIN stack_changes sc ON sc.change_id = c.id
		JOIN stacks b ON b.id = sc.stack_id
	`).Scan(&description, &bookmark)
	if err != nil {
		t.Fatalf("query recorded revision: %v", err)
	}
	if description != "feat one" || bookmark != "feature/feat-one" {
		t.Fatalf("recorded revision = (%q, %q), want (feat one, feature/feat-one)", description, bookmark)
	}
}

func TestRecordRevisionFromStackCheckoutReattachesGitHeadToStackBranch(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "edit-chg"
	if _, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "edits",
		BookmarkName: "feature/edits",
		BaseRef:      "main",
		BaseCommitID: "base123",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                      {cwd + "\n", cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {"feature/edits\n", "", "", "feature/edits\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"): {
				"",
			},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {
				"",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"false\n", "true\n", "true\n"},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"edit-chg|edit-commit|editing|base123\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@", "--name-only"):             {"internal/cli/root.go\n"},
			runnerKey(cwd, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {"feature/edits|edit-chg\n"},
			runnerKey(cwd, "jj", "commit", "-m", "all loose"):                  {"Committed\n"},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"recorded-chg|a6fcefbe|all loose|edit-chg\n",
				"recorded-chg|a6fcefbe|all loose|edit-chg\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@-", "--name-only"):                                   {"internal/cli/root.go\n", "internal/cli/root.go\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/edits", "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):                  {"a6fcefbe\n"},
			runnerKey(cwd, "jj", "log", "-r", "feature/edits", "--no-graph", "-T", "commit_id"):       {"a6fcefbe\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/edits", "a6fcefbe"):               {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/edits"):                 {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                         {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(feature/edits)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"recorded-chg|a6fcefbe|all loose|edit-chg\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "recorded-chg", "-n", "1", "--no-graph", "-T", "hidden"): {"false\n"},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"):                     {"op123\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	_, err = svc.RecordRevision(context.Background(), RevisionOptions{Message: "all loose"})
	if err != nil {
		t.Fatalf("RecordRevision() unexpected error = %v", err)
	}
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "update-ref", "refs/heads/feature/edits", "a6fcefbe"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/edits"))
	assertRunnerCalled(t, runner.calls, runnerKey(cwd, "git", "reset", "--mixed", "HEAD"))
}

func TestPushRecordedStackRejectsEmptyStackBeforeGitSideEffects(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                                              {cwd + "\n"},
			runnerKey(cwd, "git", "remote"):                                           {"origin\n"},
			runnerKey(cwd, "git", "remote", "get-url", "origin"):                      {"git@github.com:example/repo.git\n"},
			runnerKey(cwd, "git", "branch", "--show-current"):                         {"feature/body\n"},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"):       {"false\n"},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "commit_id"):   {"head123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/body", "head123"): {""},
			runnerKey(cwd, "git", "push", "origin", "feature/body"):                   {""},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head")},
		},
	}
	svc := NewServiceWithRunner(runner)

	_, _, _, err = svc.pushRecordedStack(context.Background(), RepoInfo{
		RootPath:      cwd,
		Backend:       "jj",
		DefaultRemote: ptr("origin"),
		DefaultBranch: ptr("main"),
		BranchName:    ptr("feature/body"),
		RemoteURL:     ptr("git@github.com:example/repo.git"),
	}, "origin", "feature/body", StackInfo{BaseRef: "main"}, PushOptions{Mode: PublishModeReviewOnly})
	if err == nil || !strings.Contains(err.Error(), "no revisions on stack to publish") {
		t.Fatalf("pushRecordedStack() error = %v, want no revisions", err)
	}
	for _, call := range runner.calls {
		if call == runnerKey(cwd, "git", "update-ref", "refs/heads/feature/body", "head123") {
			t.Fatalf("pushRecordedStack() updated git ref before validating stack: %v", runner.calls)
		}
		if call == runnerKey(cwd, "git", "push", "origin", "feature/body") {
			t.Fatalf("pushRecordedStack() pushed before validating stack: %v", runner.calls)
		}
	}
}

func TestPushRecordedStackSkipsGitPushWhenRemoteAlreadyAtHead(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoID, err := store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultRemote: ptr("origin"),
		DefaultBranch: ptr("main"),
		RemoteURL:     ptr("git@github.com:example/repo.git"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	head := "chg123"
	stackID, err := store.UpsertStack(context.Background(), storage.Stack{
		RepoID:       repoID,
		Name:         "body",
		BookmarkName: "feature/body",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadChangeID: &head,
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	changeID, err := store.UpsertChange(context.Background(), storage.Change{
		RepoID:          repoID,
		JJChangeID:      head,
		CurrentCommitID: "head123",
		Description:     "body change",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(context.Background(), stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}
	remoteName := "origin"
	remoteRef := "refs/heads/feature/body"
	if err := store.UpsertChangeBookmark(context.Background(), storage.ChangeBookmark{
		ChangeID:           changeID,
		BookmarkName:       "feature/body",
		RemoteName:         &remoteName,
		RemoteRef:          &remoteRef,
		LastPushedCommitID: "old123",
		CreatedAt:          1,
		UpdatedAt:          1,
	}); err != nil {
		t.Fatalf("UpsertChangeBookmark() error = %v", err)
	}
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "feature/body", "--no-graph", "-T", "commit_id"): {
				"head123\n",
			},
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "feature/body"): {
				"head123\trefs/heads/feature/body\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "head123", "--git"): {
				"diff --git a/file b/file\n",
			},
		},
		outputs: map[string][]string{
			runnerKey(repoRoot, "git", "update-ref", "refs/heads/feature/body", "head123"): {""},
		},
	}
	svc := NewServiceWithRunner(runner)

	pushed, status, warnings, err := svc.pushRecordedStack(context.Background(), RepoInfo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultRemote: ptr("origin"),
		DefaultBranch: ptr("main"),
		BranchName:    ptr("feature/body"),
		RemoteURL:     ptr("git@github.com:example/repo.git"),
	}, "origin", "feature/body", StackInfo{ID: stackID, BookmarkName: "feature/body", BaseRef: "main"}, PushOptions{Mode: PublishModeReviewAndGit})
	if err != nil {
		t.Fatalf("pushRecordedStack() error = %v", err)
	}
	if status != "already up to date" || len(warnings) != 0 {
		t.Fatalf("status=%q warnings=%#v, want already up to date without warnings", status, warnings)
	}
	if len(pushed) != 1 {
		t.Fatalf("pushed len = %d, want 1", len(pushed))
	}
	for _, call := range runner.calls {
		if call == runnerKey(repoRoot, "git", "push", "origin", "feature/body") {
			t.Fatalf("pushRecordedStack() pushed despite matching remote head: %v", runner.calls)
		}
	}
}

func TestPushRecordedStackUsesHydratedRevisions(t *testing.T) {
	repoRoot := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "commit-two", "--no-graph", "-T", "commit_id"): {"commit-two\n"},
			runnerKey(repoRoot, "git", "update-ref", "refs/heads/feature/body", "commit-two"):     {""},
			runnerKey(repoRoot, "jj", "diff", "-r", "commit-one", "--git"):                        {"diff --git a/one b/one\n"},
			runnerKey(repoRoot, "jj", "diff", "-r", "commit-two", "--git"):                        {"diff --git a/two b/two\n"},
		},
	}
	svc := NewServiceWithRunner(runner)

	pushed, status, warnings, err := svc.pushRecordedStack(context.Background(), RepoInfo{
		RootPath: repoRoot,
		Backend:  "jj",
	}, "origin", "feature/body", StackInfo{
		ID:      999,
		BaseRef: "main",
		Revisions: []RevisionSummary{
			{ChangeID: "change-one", CommitID: "commit-one", Description: "one"},
			{ChangeID: "change-two", CommitID: "commit-two", Description: "two"},
		},
	}, PushOptions{Mode: PublishModeReviewOnly})
	if err != nil {
		t.Fatalf("pushRecordedStack() error = %v", err)
	}
	if len(pushed) != 2 {
		t.Fatalf("pushed revisions = %d, want 2", len(pushed))
	}
	if status != "not pushed" {
		t.Fatalf("git push status = %q, want not pushed", status)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %#v", warnings)
	}
}

func TestPublishResultRecordsLocalMetadataBeforeHook(t *testing.T) {
	repoRoot := t.TempDir()
	t.Setenv("GX_HOME", t.TempDir())
	result := PushResult{
		Repo: RepoInfo{
			RootPath:      repoRoot,
			Backend:       "jj",
			DefaultRemote: ptr("origin"),
			DefaultBranch: ptr("main"),
			RemoteURL:     ptr("git@github.com:example/repo.git"),
		},
		RemoteName:   ptr("origin"),
		GXStackRef:   "feature/body",
		HeadCommitID: "head123",
		CurrentChange: &ChangeInfo{
			ChangeID:    "chg123",
			CommitID:    "head123",
			Description: "body change",
		},
	}
	hookSawPush := false
	err := publishResult(context.Background(), result, func(PushResult) error {
		db, err := storage.Open(context.Background())
		if err != nil {
			return err
		}
		store, err := storage.NewStore(context.Background(), db)
		if err != nil {
			return err
		}
		repo, err := store.FindRepoByRoot(context.Background(), repoRoot)
		if err != nil {
			return err
		}
		if repo == nil {
			return fmt.Errorf("repo was not recorded before hook")
		}
		push, err := store.LatestPushByBranchName(context.Background(), repo.ID, "feature/body")
		if err != nil {
			return err
		}
		if push == nil || push.HeadCommitID != "head123" {
			return fmt.Errorf("push was not recorded before hook: %#v", push)
		}
		hookSawPush = true
		return nil
	})
	if err != nil {
		t.Fatalf("publishResult() error = %v", err)
	}
	if !hookSawPush {
		t.Fatal("publish hook did not run")
	}
}

func TestPushOptionsGitExportEnabled(t *testing.T) {
	tests := []struct {
		name string
		opts PushOptions
		want bool
	}{
		{
			name: "plain publish mode pushes to git",
			opts: PushOptions{Mode: PublishModeReviewAndGit},
			want: true,
		},
		{
			name: "no github publish mode skips git",
			opts: PushOptions{Mode: PublishModeReviewOnly},
			want: false,
		},
		{
			name: "legacy default still exports git",
			opts: PushOptions{},
			want: true,
		},
		{
			name: "legacy upload only skips git",
			opts: PushOptions{UploadOnly: true},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.gitExportEnabled(); got != tt.want {
				t.Fatalf("gitExportEnabled() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestCommitCreatesStackBookmarkWithoutActiveBranch(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"): {cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"):      {""},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {""},
			runnerKey(cwd, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {
				"\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"true\n", "true\n", "true\n"},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"chg|abc|desc|\n",
				"chg|abc|desc|\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@-", "--name-only"):                                      {"", ""},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-"):                      {""},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "commit", "-m", "feat one"):                                             {""},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):                     {"abc\n"},
			runnerKey(cwd, "jj", "log", "-r", "feature/feat-one", "--no-graph", "-T", "commit_id"):       {"abc\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feat-one", "abc"):                    {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feat-one"):                 {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                            {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(feature/feat-one)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"chg|abc|desc|\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "chg", "-n", "1", "--no-graph", "-T", "hidden"): {"false\n"},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"):            {"op\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head")},
			runnerKey(cwd, "git", "branch", "--show-current"):                            {fmt.Errorf("detached")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.Commit(context.Background(), "feat one")
	if err != nil {
		t.Fatalf("Commit() unexpected error = %v", err)
	}
	if result.Stack == nil || result.Stack.BookmarkName != "feature/feat-one" {
		t.Fatalf("Commit() stack = %#v, want feature/feat-one", result.Stack)
	}
}

func TestCommitUsesGXMentionSessionForStackBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	repoRootPtr := repoRoot
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:        "session-one",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       repoRoot,
		RepoRoot:  &repoRootPtr,
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("WriteSession(session-one) error = %v", err)
	}
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:        "session-two",
		CreatedAt: 2,
		Command:   "codex",
		Cwd:       repoRoot,
		RepoRoot:  &repoRootPtr,
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("WriteSession(session-two) error = %v", err)
	}
	if err := store.WriteRequest(context.Background(), storage.Request{
		ID:             "request-one",
		SessionID:      "session-one",
		CreatedAt:      3,
		Provider:       "openai",
		Endpoint:       "/v1/responses",
		Method:         "POST",
		RequestBody:    []byte(`{"input":"please run gx add -m \"feat one\""}`),
		RequestHeaders: "{}",
	}); err != nil {
		t.Fatalf("WriteRequest(request-one) error = %v", err)
	}
	if err := store.WriteRequest(context.Background(), storage.Request{
		ID:             "request-two",
		SessionID:      "session-two",
		CreatedAt:      4,
		Provider:       "openai",
		Endpoint:       "/v1/responses",
		Method:         "POST",
		RequestBody:    []byte(`{"input":"unrelated work"}`),
		RequestHeaders: "{}",
	}); err != nil {
		t.Fatalf("WriteRequest(request-two) error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}
	if err := gxconfig.Save(gxconfig.Config{
		User: gxconfig.User{Name: "Joe Example", Email: "joe@example.com"},
	}); err != nil {
		t.Fatalf("gxconfig.Save() error = %v", err)
	}
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	bookmark := "feature/feat-one"
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "root"):                      {cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {"main\n", "", "main\n"},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"): {
				"",
			},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {
				"",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {"true\n", "true\n", "true\n"},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"chg123|abc123|old desc|parent1\n",
				"chg123|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@-", "--name-only"):                            {"a.txt\n", "a.txt\n"},
			runnerKey(cwd, "jj", "bookmark", "set", bookmark, "-r", "@-"):                      {""},
			runnerKey(cwd, "jj", "commit", "-m", "feat one"):                                   {"Committed\n"},
			runnerKey(cwd, "jj", "bookmark", "set", bookmark, "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):           {"abc123\n"},
			runnerKey(cwd, "jj", "log", "-r", bookmark, "--no-graph", "-T", "commit_id"):       {"abc123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/"+bookmark, "abc123"):              {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/"+bookmark):              {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                  {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors("+bookmark+")", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"chg123|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "chg123", "-n", "1", "--no-graph", "-T", "hidden"): {"false\n"},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"):               {"op123\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote"), fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head"), fmt.Errorf("no origin head")},
			runnerKey(cwd, "jj", "config", "get", "user.name"):                           {fmt.Errorf("missing")},
			runnerKey(cwd, "jj", "config", "get", "user.email"):                          {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.Commit(context.Background(), "feat one")
	if err != nil {
		t.Fatalf("Commit() unexpected error = %v", err)
	}
	if result.Stack == nil || result.Stack.BookmarkName != bookmark {
		t.Fatalf("Commit() stack = %#v, want %s", result.Stack, bookmark)
	}
}

func TestCurrentChangeIgnoresJJStderrNoise(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"Done importing changes from the underlying Git repo.\nchg123|abc123|feat one|parent1\n",
			},
		},
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"chg123|abc123|feat one|parent1\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "@-", "--name-only"): {
				"a.txt\nb.txt\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	change, err := svc.CurrentChange(context.Background(), repoRoot, "@-")
	if err != nil {
		t.Fatalf("CurrentChange() unexpected error = %v", err)
	}
	if change.ChangeID != "chg123" {
		t.Fatalf("CurrentChange() change id = %q, want %q", change.ChangeID, "chg123")
	}
	if change.CommitID != "abc123" {
		t.Fatalf("CurrentChange() commit id = %q, want %q", change.CommitID, "abc123")
	}
	if change.Description != "feat one" {
		t.Fatalf("CurrentChange() description = %q, want %q", change.Description, "feat one")
	}
	if len(change.Files) != 2 || change.Files[0] != "a.txt" || change.Files[1] != "b.txt" {
		t.Fatalf("CurrentChange() files = %#v", change.Files)
	}
}

func TestRecordRevisionMetadataAttachesSessionForBookmark(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)

	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	_, err = store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:        "session-1",
		CreatedAt: 1,
		Command:   "codex",
		Cwd:       repoRoot,
		GXVersion: "test",
	}); err != nil {
		t.Fatalf("WriteSession(session-1) error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}
	t.Setenv("GX_SESSION_ID", "session-1")

	err = recordCommit(context.Background(), CommitResult{
		Repo: RepoInfo{
			RootPath: repoRoot,
			Backend:  "jj",
		},
		Stack: &StackInfo{
			BookmarkName: "feature/feat-one",
		},
		Change: ChangeInfo{
			ChangeID:    "chg123",
			CommitID:    "abc123",
			Description: "feat one",
			Files:       []string{"a.txt"},
		},
		OperationID: "op123",
	})
	if err != nil {
		t.Fatalf("recordCommit() error = %v", err)
	}

	db, err = storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	defer db.Close()

	var sessionID string
	err = db.QueryRowContext(context.Background(), `
		SELECT session_id
		FROM change_sessions
	`).Scan(&sessionID)
	if err != nil {
		t.Fatalf("query change_sessions: %v", err)
	}
	if sessionID != "session-1" {
		t.Fatalf("attached session = %q, want session-1", sessionID)
	}
}

func TestRecordRevisionMetadataAttachesRepoLocalSessionWithoutEnv(t *testing.T) {
	repoRoot := t.TempDir()
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)

	db, err := storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(context.Background(), db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	_, err = store.UpsertRepo(context.Background(), storage.Repo{
		RootPath:  repoRoot,
		Backend:   "jj",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if err := store.WriteSession(context.Background(), storage.Session{
		ID:        "cursor-session",
		CreatedAt: 10,
		Command:   "cursor",
		Cwd:       repoRoot,
		GXVersion: "test",
		Source:    ptr("cursor"),
		RepoRoot:  ptr(repoRoot),
	}); err != nil {
		t.Fatalf("WriteSession(cursor-session) error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	err = recordCommit(context.Background(), CommitResult{
		Repo: RepoInfo{
			RootPath: repoRoot,
			Backend:  "jj",
		},
		Stack: &StackInfo{
			BookmarkName: "feature/feat-one",
		},
		Change: ChangeInfo{
			ChangeID:    "chg123",
			CommitID:    "abc123",
			Description: "feat one",
			Files:       []string{"a.txt"},
		},
		OperationID: "op123",
	})
	if err != nil {
		t.Fatalf("recordCommit() error = %v", err)
	}

	db, err = storage.Open(context.Background())
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	defer db.Close()

	var sessionID string
	err = db.QueryRowContext(context.Background(), `
		SELECT session_id
		FROM change_sessions
	`).Scan(&sessionID)
	if err != nil {
		t.Fatalf("query change_sessions: %v", err)
	}
	if sessionID != "cursor-session" {
		t.Fatalf("attached session = %q, want cursor-session", sessionID)
	}
}

func TestPromptRequiredValueUsesDefaultOnBlank(t *testing.T) {
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	var out bytes.Buffer
	value, err := promptRequiredValue(strings.NewReader("\n"), &out, "gx name", "Joe Example")
	if err != nil {
		t.Fatalf("promptRequiredValue() unexpected error = %v", err)
	}
	if value != "Joe Example" {
		t.Fatalf("promptRequiredValue() = %q, want %q", value, "Joe Example")
	}
	if !strings.Contains(stripANSI(out.String()), "gx name [Joe Example]: ") {
		t.Fatalf("prompt output missing default: %q", out.String())
	}
}

func TestEnsureIdentityInteractivePromptsAndPrintsConfigNote(t *testing.T) {
	gxHome := t.TempDir()
	t.Setenv("GX_HOME", gxHome)
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	repoRoot := t.TempDir()

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "git", "config", "user.name"):                                     {"Joe Example\n"},
			runnerKey(repoRoot, "git", "config", "user.email"):                                    {"joe@example.com\n"},
			runnerKey(repoRoot, "jj", "config", "set", "--user", "user.name", "Joe Example"):      {""},
			runnerKey(repoRoot, "jj", "config", "set", "--user", "user.email", "joe@example.com"): {""},
		},
		errors: map[string][]error{
			runnerKey(repoRoot, "jj", "config", "get", "user.name"):  {fmt.Errorf("missing")},
			runnerKey(repoRoot, "jj", "config", "get", "user.email"): {fmt.Errorf("missing")},
		},
	}
	svc := NewServiceWithRunner(runner)

	var out bytes.Buffer
	name, email, changed, err := svc.ensureIdentity(context.Background(), repoRoot, InitOptions{
		Interactive: true,
		In:          strings.NewReader("\n\n"),
		Out:         &out,
	})
	if err != nil {
		t.Fatalf("ensureIdentity() unexpected error = %v", err)
	}
	if !changed {
		t.Fatalf("ensureIdentity() did not report changes")
	}
	if name != "Joe Example" || email != "joe@example.com" {
		t.Fatalf("ensureIdentity() = %q <%s>", name, email)
	}
	text := stripANSI(out.String())
	if !strings.Contains(text, "Your config is stored in "+gxHome+"/config.json") {
		t.Fatalf("missing local storage note: %q", text)
	}
	if !strings.Contains(text, "and by initializing with gx you share your email with gx.") {
		t.Fatalf("missing gx signup note: %q", text)
	}
	namePrompt := strings.Index(text, "gx name [Joe Example]: ")
	emailPrompt := strings.Index(text, "gx email [joe@example.com]: ")
	configNote := strings.Index(text, "Your config is stored in "+gxHome+"/config.json")
	signupNote := strings.Index(text, "and by initializing with gx you share your email with gx.")
	if namePrompt == -1 || emailPrompt == -1 {
		t.Fatalf("missing prompts: %q", text)
	}
	if !(configNote < signupNote && signupNote < namePrompt && namePrompt < emailPrompt) {
		t.Fatalf("unexpected prompt/note order: %q", text)
	}
}

func TestSplitCommitHunkRequiresPatchFile(t *testing.T) {
	svc := NewServiceWithRunner(&fakeRunner{})
	_, err := svc.SplitCommitByHunkPatch(context.Background(), SplitCommitOptions{
		Message: "feat",
		Hunk:    true,
	})
	if err == nil || !strings.Contains(err.Error(), "patch file is required") {
		t.Fatalf("SplitCommit() error = %v, want missing patch file error", err)
	}
}

func TestSyncDefaultsToOrigin(t *testing.T) {
	repoRoot := t.TempDir()
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "git", "rev-parse", "--show-toplevel"): {cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"):     {"main\n"},
			runnerKey(cwd, "git", "fetch", "--prune", "origin"):   {"Fetched\n"},
		},
		errors: map[string][]error{
			runnerKey(cwd, "git", "remote"):                                              {fmt.Errorf("no remote")},
			runnerKey(cwd, "git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short"): {fmt.Errorf("no origin head")},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.Sync(context.Background(), "")
	if err != nil {
		t.Fatalf("Sync() unexpected error = %v", err)
	}
	if result.RemoteName != "origin" {
		t.Fatalf("Sync() remote = %q, want origin", result.RemoteName)
	}
	want := runnerKey(cwd, "git", "fetch", "--prune", "origin")
	found := false
	for _, call := range runner.calls {
		if call == want {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Sync() missing call %q in %v", want, runner.calls)
	}
}

func TestIsSQLiteBusy(t *testing.T) {
	if !isSQLiteBusy(fmt.Errorf("database is locked (5) (SQLITE_BUSY)")) {
		t.Fatalf("isSQLiteBusy() = false, want true")
	}
	if isSQLiteBusy(fmt.Errorf("some other error")) {
		t.Fatalf("isSQLiteBusy() = true, want false")
	}
}

func TestWithBusyRetryRetriesAndSucceeds(t *testing.T) {
	attempts := 0
	err := withBusyRetry(context.Background(), "test op", func() error {
		attempts++
		if attempts < 3 {
			return fmt.Errorf("database is locked (5) (SQLITE_BUSY)")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("withBusyRetry() unexpected error = %v", err)
	}
	if attempts != 3 {
		t.Fatalf("withBusyRetry() attempts = %d, want 3", attempts)
	}
}

func TestWithBusyRetryReturnsBusyFailureAfterExhaustion(t *testing.T) {
	start := time.Now()
	err := withBusyRetry(context.Background(), "test op", func() error {
		return fmt.Errorf("database is locked (5) (SQLITE_BUSY)")
	})
	if err == nil {
		t.Fatalf("withBusyRetry() expected error")
	}
	if !strings.Contains(err.Error(), "failed after 5 retries") {
		t.Fatalf("withBusyRetry() error = %v", err)
	}
	if time.Since(start) < 100*time.Millisecond {
		t.Fatalf("withBusyRetry() did not appear to retry")
	}
}
