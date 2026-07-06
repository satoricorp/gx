package vcs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/satoricorp/gx/internal/storage"
)

func TestRebaseMissingStackBaseRefsContinuesOnFailure(t *testing.T) {
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
	for _, bookmark := range []string{"feature/one", "feature/two"} {
		if _, err := store.UpsertStack(ctx, storage.Stack{
			RepoID:       repoID,
			Name:         bookmark,
			BookmarkName: bookmark,
			BaseRef:      "feature/main",
			BaseCommitID: "base",
			Status:       "draft",
			CreatedAt:    1,
			UpdatedAt:    1,
		}); err != nil {
			t.Fatalf("UpsertStack(%s) error = %v", bookmark, err)
		}
	}

	revsetOne := stackRebaseSourceRevset("feature/one", "feature/main", "main", false)
	revsetTwo := stackRebaseSourceRevset("feature/two", "feature/main", "main", false)
	revExists := func(rev string) string {
		return runnerKey(repoRoot, "jj", "log", "-r", rev, "-n", "1", "--no-graph", "-T", "change_id")
	}
	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "jj", "root"): {repoRoot + "\n", repoRoot + "\n"},
			runnerKey(repoRoot, "jj", "log", "-r", "main", "-n", "1", "--no-graph", "-T", "change_id"): {
				"mainchange\n", "mainchange\n", "mainchange\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "main", "--limit", "1", "--no-graph", "-T", "commit_id"): {
				"newbase\n", "newbase\n", "newbase\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "feature/one", "--limit", "1", "--no-graph", "-T", "commit_id"): {
				"onecommit\n", "onecommit\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", revsetOne, "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"", "",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(main)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"", "",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "@", "--no-graph", "-T", "change_id"): {
				"onechange\n",
			},
			runnerKey(repoRoot, "jj", "log", "-r", "feature/one", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"onechange|onecommit|one|\n",
			},
			runnerKey(repoRoot, "jj", "diff", "-r", "feature/one", "--name-only"): {
				"",
			},
		},
		outputs: map[string][]string{
			runnerKey(repoRoot, "jj", "rebase", "-s", revsetOne, "-d", "main"): {""},
			runnerKey(repoRoot, "git", "rev-parse", "main"):                    {"newbase\n", "newbase\n"},
		},
		errors: map[string][]error{
			revExists("feature/main"):        {fmt.Errorf(`revision "feature/main" doesn't exist`), fmt.Errorf(`revision "feature/main" doesn't exist`), fmt.Errorf(`revision "feature/main" doesn't exist`), fmt.Errorf(`revision "feature/main" doesn't exist`)},
			revExists("feature/main@origin"): {fmt.Errorf(`revision "feature/main@origin" doesn't exist`), fmt.Errorf(`revision "feature/main@origin" doesn't exist`), fmt.Errorf(`revision "feature/main@origin" doesn't exist`), fmt.Errorf(`revision "feature/main@origin" doesn't exist`)},
			runnerKey(repoRoot, "jj", "rebase", "-s", revsetTwo, "-d", "main"): {
				fmt.Errorf(`Revision "feature/main" doesn't exist`),
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	result, err := svc.RebaseMissingStackBaseRefs(ctx, []MissingStackBaseRef{
		{Name: "one", BookmarkName: "feature/one", MissingBaseRef: "feature/main", DefaultBaseRef: "main"},
		{Name: "two", BookmarkName: "feature/two", MissingBaseRef: "feature/main", DefaultBaseRef: "main"},
	})
	if err != nil {
		t.Fatalf("RebaseMissingStackBaseRefs() error = %v", err)
	}
	if len(result.Fixed) != 1 || result.Fixed[0].BookmarkName != "feature/one" {
		t.Fatalf("Fixed = %#v, want feature/one", result.Fixed)
	}
	if len(result.Failed) != 1 || result.Failed[0].Issue.BookmarkName != "feature/two" {
		t.Fatalf("Failed = %#v, want feature/two", result.Failed)
	}
}

func TestReconcileStoredStacksRemoteStateMarksFullyPushedStackPublished(t *testing.T) {
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
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "send-stack",
		BookmarkName: "feature/send-stack",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadCommitID: ptr("remote-head"),
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "change-one",
		CurrentCommitID: "remote-head",
		Description:     "send stack",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}

	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "feature/send-stack"): {
				"remote-head\trefs/heads/feature/send-stack\n",
			},
		},
		errors: map[string][]error{
			runnerKey(repoRoot, "git", "merge-base", "--is-ancestor", "remote-head", "main"):        {fmt.Errorf("not merged")},
			runnerKey(repoRoot, "git", "merge-base", "--is-ancestor", "remote-head", "origin/main"): {fmt.Errorf("not merged")},
		},
	}
	svc := NewServiceWithRunner(runner)
	repo := RepoInfo{RootPath: repoRoot, DefaultBranch: ptr("main"), DefaultRemote: ptr("origin")}

	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		t.Fatalf("ListStacksByRepoID() error = %v", err)
	}
	if err := svc.reconcileStoredStacksRemoteState(ctx, store, repo, repoID, stacks); err != nil {
		t.Fatalf("reconcileStoredStacksRemoteState() error = %v", err)
	}
	stack, err := store.FindStackByBookmark(ctx, repoID, "feature/send-stack")
	if err != nil {
		t.Fatalf("FindStackByBookmark() error = %v", err)
	}
	if stack == nil || stack.Status != "published" {
		t.Fatalf("stack after reconcile = %#v, want published", stack)
	}
	if stack.RemoteRef == nil || *stack.RemoteRef != "refs/heads/feature/send-stack" {
		t.Fatalf("stack remote_ref = %v, want refs/heads/feature/send-stack", stack.RemoteRef)
	}
}

func TestReconcileStoredStacksRemoteStateKeepsNewLocalWorkDraft(t *testing.T) {
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
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "cli",
		BookmarkName: "feature/cli",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadCommitID: ptr("local-head"),
		Status:       "draft",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "change-new",
		CurrentCommitID: "local-head",
		Description:     "new local work",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}

	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "feature/cli"): {
				"merged-remote-head\trefs/heads/feature/cli\n",
			},
		},
		errors: map[string][]error{
			runnerKey(repoRoot, "git", "merge-base", "--is-ancestor", "local-head", "merged-remote-head"): {
				fmt.Errorf("not an ancestor"),
			},
		},
	}
	svc := NewServiceWithRunner(runner)
	repo := RepoInfo{RootPath: repoRoot, DefaultBranch: ptr("main"), DefaultRemote: ptr("origin")}

	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		t.Fatalf("ListStacksByRepoID() error = %v", err)
	}
	if err := svc.reconcileStoredStacksRemoteState(ctx, store, repo, repoID, stacks); err != nil {
		t.Fatalf("reconcileStoredStacksRemoteState() error = %v", err)
	}
	stack, err := store.FindStackByBookmark(ctx, repoID, "feature/cli")
	if err != nil {
		t.Fatalf("FindStackByBookmark() error = %v", err)
	}
	if stack == nil || stack.Status != "draft" {
		t.Fatalf("stack after reconcile = %#v, want draft local work", stack)
	}
	if stack.RemoteRef == nil || *stack.RemoteRef != "refs/heads/feature/cli" {
		t.Fatalf("stack remote_ref = %v, want refs/heads/feature/cli", stack.RemoteRef)
	}
}

func TestReconcileStoredStacksRemoteStateMarksMergedWhenRemoteMatchesMain(t *testing.T) {
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
	stackID, err := store.UpsertStack(ctx, storage.Stack{
		RepoID:       repoID,
		Name:         "merged-stack",
		BookmarkName: "feature/merged-stack",
		BaseRef:      "main",
		BaseCommitID: "base",
		HeadCommitID: ptr("remote-head"),
		Status:       "published",
		CreatedAt:    1,
		UpdatedAt:    1,
	})
	if err != nil {
		t.Fatalf("UpsertStack() error = %v", err)
	}
	changeID, err := store.UpsertChange(ctx, storage.Change{
		RepoID:          repoID,
		JJChangeID:      "change-one",
		CurrentCommitID: "remote-head",
		Description:     "merged stack",
		Status:          "draft",
		FirstSeenAt:     1,
		UpdatedAt:       1,
	})
	if err != nil {
		t.Fatalf("UpsertChange() error = %v", err)
	}
	if err := store.AddChangeToStack(ctx, stackID, changeID, 1); err != nil {
		t.Fatalf("AddChangeToStack() error = %v", err)
	}

	runner := &fakeRunner{
		stdoutOutputs: map[string][]string{
			runnerKey(repoRoot, "git", "ls-remote", "--heads", "origin", "feature/merged-stack"): {
				"remote-head\trefs/heads/feature/merged-stack\n",
			},
		},
		outputs: map[string][]string{
			runnerKey(repoRoot, "git", "merge-base", "--is-ancestor", "remote-head", "remote-head"): {""},
			runnerKey(repoRoot, "git", "merge-base", "--is-ancestor", "remote-head", "main"):        {""},
		},
	}
	svc := NewServiceWithRunner(runner)
	repo := RepoInfo{RootPath: repoRoot, DefaultBranch: ptr("main"), DefaultRemote: ptr("origin")}

	stacks, err := store.ListStacksByRepoID(ctx, repoID)
	if err != nil {
		t.Fatalf("ListStacksByRepoID() error = %v", err)
	}
	if err := svc.reconcileStoredStacksRemoteState(ctx, store, repo, repoID, stacks); err != nil {
		t.Fatalf("reconcileStoredStacksRemoteState() error = %v", err)
	}
	stack, err := store.FindStackByBookmark(ctx, repoID, "feature/merged-stack")
	if err != nil {
		t.Fatalf("FindStackByBookmark() error = %v", err)
	}
	if stack == nil || stack.Status != "merged" {
		t.Fatalf("stack after reconcile = %#v, want merged", stack)
	}
}
