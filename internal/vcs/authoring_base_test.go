package vcs

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// A brand-new repository whose base branch has no commits (unborn HEAD) must
// fail fast with an actionable message instead of surfacing a raw jj revset
// error later in the pipeline (or hanging on an interactive retry prompt).
func TestEnsureAuthoringCheckoutFailsFastWhenBaseHasNoCommits(t *testing.T) {
	cwd := t.TempDir()
	branch := "main"
	repo := RepoInfo{RootPath: cwd, BranchName: &branch}

	runner := &fakeRunner{
		outputs: map[string][]string{},
		errors: map[string][]error{
			runnerKey(cwd, "jj", "log", "-r", "main", "--limit", "1", "--no-graph", "-T", "commit_id"): {
				fmt.Errorf("Error: Revision `main` doesn't exist"),
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	_, _, err := svc.ensureAuthoringCheckout(context.Background(), repo, "gx generate")
	if err == nil {
		t.Fatalf("ensureAuthoringCheckout() error = nil, want unborn-base failure")
	}
	message := err.Error()
	if !strings.Contains(message, "at least one commit") {
		t.Fatalf("ensureAuthoringCheckout() error = %q, want mention of missing commits", message)
	}
	if !strings.Contains(message, "initial commit") {
		t.Fatalf("ensureAuthoringCheckout() error = %q, want actionable initial-commit hint", message)
	}
	if !strings.Contains(message, "gx generate") {
		t.Fatalf("ensureAuthoringCheckout() error = %q, want command name", message)
	}
}

func TestRequireBaseRevisionAllowsResolvableBase(t *testing.T) {
	cwd := t.TempDir()
	repo := RepoInfo{RootPath: cwd}

	runner := &fakeRunner{
		outputs: map[string][]string{
			runnerKey(cwd, "jj", "log", "-r", "main", "--limit", "1", "--no-graph", "-T", "commit_id"): {"abc123\n"},
		},
	}
	svc := NewServiceWithRunner(runner)

	if err := svc.requireBaseRevision(context.Background(), repo, "main", "gx generate"); err != nil {
		t.Fatalf("requireBaseRevision() error = %v, want nil", err)
	}
}

func TestRequireBaseRevisionIgnoresTransientLookupErrors(t *testing.T) {
	cwd := t.TempDir()
	repo := RepoInfo{RootPath: cwd}

	runner := &fakeRunner{
		outputs: map[string][]string{},
		errors: map[string][]error{
			runnerKey(cwd, "jj", "log", "-r", "main", "--limit", "1", "--no-graph", "-T", "commit_id"): {
				fmt.Errorf("jj: some transient failure"),
			},
		},
	}
	svc := NewServiceWithRunner(runner)

	if err := svc.requireBaseRevision(context.Background(), repo, "main", "gx generate"); err != nil {
		t.Fatalf("requireBaseRevision() error = %v, want nil for transient lookup failure", err)
	}
}
