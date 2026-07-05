package vcs

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/gxconfig"
)

func TestJJStackRevsetExcludesContainerAncestors(t *testing.T) {
	got := jjStackRevset("main")
	if !strings.Contains(got, "~ancestors(main)") {
		t.Fatalf("jjStackRevset() = %q, want container exclusion", got)
	}
}

func TestCommitRecoversDetachedWithJJBookmark(t *testing.T) {
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
			runnerKey(cwd, "jj", "root"): {cwd + "\n", cwd + "\n", cwd + "\n"},
			runnerKey(cwd, "git", "branch", "--show-current"): {
				"",
				"",
				"",
			},
			runnerKey(cwd, "jj", "bookmark", "list", "-T", jjBookmarkListTmpl): {
				"main|workchange\n",
				"main|workchange\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "empty"): {
				"true\n",
				"true\n",
				"true\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@", "--no-graph", "-T", "change_id"): {
				"workchange\n",
				"workchange\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", `change_id ++ "|" ++ commit_id ++ "|" ++ description.first_line() ++ "|" ++ parents.map(|c| c.change_id()).join(",") ++ "\n"`): {
				"workchange|abc123|feat one|parent1\n",
				"workchange|abc123|feat one|parent1\n",
				"workchange|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "diff", "-r", "@-", "--name-only"): {
				"a.txt\n",
				"a.txt\n",
			},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.name", "Joe Example"):                  {"", ""},
			runnerKey(cwd, "jj", "config", "set", "--user", "user.email", "joe@example.com"):             {"", ""},
			runnerKey(cwd, "jj", "commit", "-m", "feat one"):                                             {"Committed\n"},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-"):                      {""},
			runnerKey(cwd, "jj", "bookmark", "set", "feature/feat-one", "-r", "@-", "--allow-backwards"): {""},
			runnerKey(cwd, "jj", "log", "-r", "@-", "--no-graph", "-T", "commit_id"):                     {"abc123\n"},
			runnerKey(cwd, "jj", "log", "-r", "feature/feat-one", "--no-graph", "-T", "commit_id"):       {"abc123\n"},
			runnerKey(cwd, "git", "update-ref", "refs/heads/feature/feat-one", "abc123"):                 {""},
			runnerKey(cwd, "git", "symbolic-ref", "HEAD", "refs/heads/feature/feat-one"):                 {""},
			runnerKey(cwd, "git", "reset", "--mixed", "HEAD"):                                            {""},
			runnerKey(cwd, "jj", "log", "-r", "mutable() & ~empty() & ~hidden() & ancestors(@) & ~ancestors(feature/feat-one)", "--reversed", "--no-graph", "-T", jjStackLineTmpl): {
				"workchange|abc123|feat one|parent1\n",
			},
			runnerKey(cwd, "jj", "log", "-r", "workchange", "-n", "1", "--no-graph", "-T", "hidden"): {
				"false\n",
			},
			runnerKey(cwd, "jj", "op", "log", "-n", "1", "--no-graph", "-T", "id"): {"op123\n"},
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
}

func TestIsGitIndexLock(t *testing.T) {
	if !isGitIndexLock(fmt.Errorf("fatal: Unable to create '/repo/.git/index.lock': File exists.")) {
		t.Fatalf("isGitIndexLock() = false, want true")
	}
}

func TestWithGitRetryRetriesIndexLock(t *testing.T) {
	attempts := 0
	err := withGitRetry(context.Background(), "git op", func() error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("Unable to create index.lock: File exists")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("withGitRetry() unexpected error = %v", err)
	}
	if attempts != 2 {
		t.Fatalf("withGitRetry() attempts = %d, want 2", attempts)
	}
}

type flakyLockRunner struct {
	failures int
	calls    int
}

func (r *flakyLockRunner) Run(ctx context.Context, dir, name string, args ...string) (string, error) {
	r.calls++
	if r.calls <= r.failures {
		return "", fmt.Errorf("jj %s: Failed to reset Git HEAD state: The lockfile at '.git/index.lock' might need manual deletion", name)
	}
	return "", nil
}

func (r *flakyLockRunner) RunStdout(ctx context.Context, dir, name string, args ...string) (string, error) {
	return r.Run(ctx, dir, name, args...)
}

func (r *flakyLockRunner) RunStream(ctx context.Context, dir, name string, args ...string) error {
	_, err := r.Run(ctx, dir, name, args...)
	return err
}

func TestRunJJGitBackedRetriesGitIndexLock(t *testing.T) {
	runner := &flakyLockRunner{failures: 2}
	svc := NewServiceWithRunner(runner)
	if _, err := svc.runJJGitBacked(context.Background(), "/repo", "bookmark", "set", "feature/x", "-r", "@-"); err != nil {
		t.Fatalf("runJJGitBacked() error = %v", err)
	}
	if runner.calls != 3 {
		t.Fatalf("runJJGitBacked() calls = %d, want 3", runner.calls)
	}
}
