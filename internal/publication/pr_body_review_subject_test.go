package publication

import (
	"context"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/reviewbundle"
)

// The PR-summary pipeline shares internal/codereview with `gx review`, so a
// change to how the local command picks its review subject could in principle
// move every PR summary GX writes. It cannot: this pipeline's subject is the
// pushed bundle, and it never routes through the local change-set resolution
// (resolveChangeSet / reviewChangeSet / normalizeOptions), which only
// codereview.Engine.Review reaches.
//
// The observable form of that claim: dirtying the working tree is exactly what
// makes the local review switch subject, and it leaves this body byte for byte
// unchanged. (The blast-radius reach analysis does read files under the repo
// root, so the edits below deliberately leave every HelloWorld call site
// standing: what is pinned here is the review subject, not file reads.)
func TestPullRequestBodyIsIndependentOfTheWorkingTreeState(t *testing.T) {
	oldReviewer := prSummaryReviewerFromEnvWithInfo
	oldContext := collectPRSummaryContext
	defer func() {
		prSummaryReviewerFromEnvWithInfo = oldReviewer
		collectPRSummaryContext = oldContext
	}()
	prSummaryReviewerFromEnvWithInfo = func() (codereview.AIReviewer, codereview.ReviewerInfo) {
		return nil, codereview.ReviewerInfo{}
	}
	collectPRSummaryContext = func(_ context.Context, _ reviewbundle.Artifact, _ prBodyCatalog) (prSummaryContext, error) {
		return prSummaryContext{}, nil
	}
	t.Setenv("GX_PR_BLAST_RADIUS", "1")

	root, headSHA := initLexicalReachRepo(t)
	artifact := reachTestArtifact(root, headSHA)

	clean, err := GitHubPullRequestBodyFromArtifact(context.Background(), artifact)
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}

	// Exactly the state that makes `gx review` prefer the working tree over a
	// commit range: a tracked file edited but not committed, plus an untracked
	// file. A local review would switch subject here; a PR summary must not.
	writeReachFile(t, root, "consumer.go", "package main\n\nfunc use() {\n\tHelloWorld()\n\t// uncommitted edit\n}\n")
	writeReachFile(t, root, "scratch-notes.txt", "uncommitted scratch\n")

	dirty, err := GitHubPullRequestBodyFromArtifact(context.Background(), artifact)
	if err != nil {
		t.Fatalf("GitHubPullRequestBodyFromArtifact() error = %v", err)
	}
	if dirty != clean {
		t.Fatalf("PR summary changed with the working tree state\n--- clean\n%s\n--- dirty\n%s", clean, dirty)
	}
	if strings.Contains(dirty, "scratch-notes.txt") {
		t.Fatalf("PR summary leaked an uncommitted file:\n%s", dirty)
	}
}
