package vcs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestDetectMissingStackBaseRefsFindsMissingParentBase(t *testing.T) {
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
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if _, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "structural",
		BookmarkName: "feature/structural",
		BaseRef:      "feature/authoring",
		BaseCommitID: "base",
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	if _, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "healthy",
		BookmarkName: "feature/healthy",
		BaseRef:      "main",
		BaseCommitID: "base",
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack(healthy) error = %v", err)
	}

	repo := RepoInfo{RootPath: repoRoot, DefaultBranch: ptr("main")}
	revExists := func(rev string) string {
		return runnerKey(repoRoot, "jj", "log", "-r", rev, "-n", "1", "--no-graph", "-T", "change_id")
	}
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			revExists("main"):               {"mainchange\n", "mainchange\n"},
			revExists("feature/structural"): {"stackchange\n", "stackchange\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "feature/structural", "--limit", "1", "--no-graph", "-T", "commit_id"): {
				"stackcommit\n", "stackcommit\n",
			},
		},
		errors: map[string][]error{
			revExists("feature/authoring"):        {fmt.Errorf(`revision "feature/authoring" doesn't exist`), fmt.Errorf(`revision "feature/authoring" doesn't exist`)},
			revExists("feature/authoring@origin"): {fmt.Errorf(`revision "feature/authoring@origin" doesn't exist`), fmt.Errorf(`revision "feature/authoring@origin" doesn't exist`)},
		},
	}
	svc := NewServiceWithRunner(runner)

	issues, err := svc.missingStackBaseRefsForRepo(ctx, store, repo, repoID)
	if err != nil {
		t.Fatalf("missingStackBaseRefsForRepo() error = %v", err)
	}
	got := missingStackBaseRefStatusFromIssues(issues)
	if !got.NeedsRebaseOntoDefault || got.FixAction != MissingStackBaseFixAction {
		t.Fatalf("status = %#v, want rebase fix metadata", got)
	}
	if len(got.Issues) != 1 {
		t.Fatalf("issues = %#v, want one missing base ref", got.Issues)
	}
	issue := got.Issues[0]
	if issue.BookmarkName != "feature/structural" || issue.MissingBaseRef != "feature/authoring" || issue.DefaultBaseRef != "main" {
		t.Fatalf("issue = %#v, want structural on missing feature/authoring", issue)
	}
	if got.FixPrompt == "" {
		t.Fatal("expected fix prompt for agent/interactive approval")
	}
}

func TestStackMergeBaseSelectorsFallsBackToDefaultBranch(t *testing.T) {
	repoRoot := t.TempDir()
	runner := &fakeRunner{
		errors: map[string][]error{
			runnerKey(repoRoot, "jj", "log", "-r", "feature/authoring", "-n", "1", "--no-graph", "-T", "change_id"): {
				fmt.Errorf(`revision "feature/authoring" doesn't exist`),
			},
			runnerKey(repoRoot, "jj", "log", "-r", "feature/authoring@origin", "-n", "1", "--no-graph", "-T", "change_id"): {
				fmt.Errorf(`revision "feature/authoring@origin" doesn't exist`),
			},
		},
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "main", "-n", "1", "--no-graph", "-T", "change_id"): {
				"mainchange\n",
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	selectors := svc.stackMergeBaseSelectors(context.Background(), repoRoot, "feature/authoring", "origin", "main")
	if len(selectors) != 1 || selectors[0] != "main" {
		t.Fatalf("selectors = %#v, want default branch fallback", selectors)
	}
}

func TestRebaseStackOntoDefaultUnlockedUpdatesStoredBase(t *testing.T) {
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
	repoID, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	})
	if err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if _, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "structural",
		BookmarkName: "feature/structural",
		BaseRef:      "feature/authoring",
		BaseCommitID: "oldbase",
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	}); err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}

	rebaseRevset := stackRebaseSourceRevset("feature/structural", "feature/authoring", "main")
	changeTmpl := `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`
	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "rebase", "-s", rebaseRevset, "-d", "main"): {""},
		},
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "log", "-r", "main", "-n", "1", "--no-graph", "-T", "change_id"): {
				"mainchange\n", "mainchange\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "main", "--limit", "1", "--no-graph", "-T", "commit_id"): {
				"newbase\n", "newbase\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "feature/structural", "--limit", "1", "--no-graph", "-T", "commit_id"): {
				"stackcommit\n", "stackcommit\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", rebaseRevset, "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(main)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "feature/structural", "--no-graph", "-T", changeTmpl): {
				"stackchange|stackcommit|structural change|\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "feature/structural", "--name-only"): {
				"",
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	repo := RepoInfo{RootPath: repoRoot, DefaultBranch: ptr("main")}

	if err := svc.rebaseStackOntoDefaultUnlocked(ctx, store, repo, repoID, MissingStackBaseRef{
		Name:           "structural",
		BookmarkName:   "feature/structural",
		MissingBaseRef: "feature/authoring",
		DefaultBaseRef: "main",
	}, 2); err != nil {
		t.Fatalf("rebaseStackOntoDefaultUnlocked() error = %v", err)
	}
	assertRunnerCalled(t, runner.calls, runnerKey(repoRoot, "jj", "rebase", "-s", rebaseRevset, "-d", "main"))

	stack, err := store.FindStackByBookmark(ctx, repoID, "feature/structural")
	if err != nil {
		t.Fatalf("FindStackByBookmark() error = %v", err)
	}
	if stack == nil || stack.BaseRef != "main" || stack.BaseCommitID != "newbase" {
		t.Fatalf("stack after fix = %#v, want base main/newbase", stack)
	}
}

func TestStackRebaseSourceRevsetExcludesMissingBase(t *testing.T) {
	got := stackRebaseSourceRevset("feature/structural", "feature/authoring", "main")
	want := `ancestors(feature/structural) & mutable() & ~empty() & ~hidden() & ~ancestors(feature/authoring) & ~ancestors(main)`
	if got != want {
		t.Fatalf("revset = %q, want %q", got, want)
	}
}

func TestMissingStackBaseRefStatusJSONFields(t *testing.T) {
	status := missingStackBaseRefStatusFromIssues([]MissingStackBaseRef{{
		StackID:        1,
		BookmarkName:   "feature/structural",
		MissingBaseRef: "feature/authoring",
		DefaultBaseRef: "main",
	}})
	if !status.NeedsRebaseOntoDefault {
		t.Fatal("expected needs_rebase_onto_default")
	}
	if status.FixAction != MissingStackBaseFixAction {
		t.Fatalf("fix_action = %q, want %q", status.FixAction, MissingStackBaseFixAction)
	}
	if status.FixPrompt == "" {
		t.Fatal("expected fix_prompt")
	}
}

func TestDetectMissingStackBaseRefsUsesConfiguredRepo(t *testing.T) {
	repoRoot := t.TempDir()
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}
	t.Setenv("GX_HOME", t.TempDir())
	prev, _ := os.Getwd()
	if err := os.Chdir(repoRoot); err != nil {
		t.Fatalf("Chdir() error = %v", err)
	}
	defer os.Chdir(prev)

	ctx := context.Background()
	db, err := storage.Open(ctx)
	if err != nil {
		t.Fatalf("storage.Open() error = %v", err)
	}
	store, err := storage.NewStore(ctx, db)
	if err != nil {
		t.Fatalf("storage.NewStore() error = %v", err)
	}
	if _, err := store.UpsertRepo(ctx, storage.Repo{
		RootPath:      repoRoot,
		Backend:       "jj",
		DefaultBranch: ptr("main"),
		CreatedAt:     1,
		UpdatedAt:     1,
	}); err != nil {
		t.Fatalf("UpsertRepo() error = %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("store.Close() error = %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(cwd, "jj", "root"): {cwd + "\n", cwd + "\n"},
		},
	}
	svc := NewServiceWithRunner(runner)
	if _, err := svc.DetectMissingStackBaseRefs(ctx); err != nil {
		t.Fatalf("DetectMissingStackBaseRefs() error = %v", err)
	}
}
