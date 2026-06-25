package publication

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

func TestGitHubPullRequestBodyIncludesHeadCommitMissingFromStackMetadata(t *testing.T) {
	root := t.TempDir()
	gitTest(t, root, "init")
	gitTest(t, root, "config", "user.email", "test@example.com")
	gitTest(t, root, "config", "user.name", "Test User")
	writeTestFile(t, root, "internal/publication/pr_body.go", "package publication\n\nfunc Old() {}\n")
	gitTest(t, root, "add", ".")
	gitTest(t, root, "commit", "-m", "base")

	writeTestFile(t, root, "internal/publication/pr_body.go", "package publication\n\nfunc Old() {}\nfunc StackChange() {}\n")
	gitTest(t, root, "add", ".")
	gitTest(t, root, "commit", "-m", "render rich PR bodies from review bundles")
	stackCommit := strings.TrimSpace(gitTest(t, root, "rev-parse", "HEAD"))
	stackPatch := gitTest(t, root, "show", "--format=", "--no-ext-diff", "HEAD")

	writeTestFile(t, root, "internal/publication/pr_body.go", "package publication\n\nfunc Old() {}\nfunc StackChange() {}\nfunc HeadRepair() {}\n")
	gitTest(t, root, "add", ".")
	gitTest(t, root, "commit", "-m", "replace PR body fallback with context summaries")
	headCommit := strings.TrimSpace(gitTest(t, root, "rev-parse", "HEAD"))

	oldReviewer := prSummaryReviewerFromEnv
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnv = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnv = func() codereview.AIReviewer { return nil }
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{Sources: []string{"codebase"}}, nil
	}

	prURL := "https://github.com/satoricorp/gx/pull/11"
	body, err := GitHubPullRequestBodyFromArtifact(context.Background(), reviewbundle.NewArtifact(reviewbundle.Bundle{
		Event:         "gx.pr",
		SchemaVersion: reviewbundle.SchemaVersion,
		Repo:          reviewbundle.RepoPayload{RootPath: root},
		Push: reviewbundle.PushPayload{
			HeadCommitID:         headCommit,
			GitHubPullRequestURL: &prURL,
		},
		Stack: []reviewbundle.StackPayload{{
			BranchName: "bug/fix-github-pr-summary",
			Patch:      stackPatch,
			Change: reviewbundle.ChangePayload{
				CurrentCommitID: stackCommit,
				Description:     "render rich PR bodies from review bundles",
				Files:           []string{"internal/publication/pr_body.go"},
			},
			GitHubPullRequestURL: &prURL,
		}},
	}))
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	for _, want := range []string{
		"replace PR body fallback with context summaries",
		"2 GX revisions",
		githubHunkLink(prURL, "internal/publication/pr_body.go", 4, 4, 1),
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q:\n%s", want, body)
		}
	}
}

func gitTest(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, string(out))
	}
	return string(out)
}

func writeTestFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}
