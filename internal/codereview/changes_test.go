package codereview

import (
	"context"
	"strings"
	"testing"
)

// gitCommitMessage commits everything staged with an explicit message, so a
// test repo can carry more than one commit.
func gitCommitMessage(t *testing.T, root, message string) {
	t.Helper()
	runGit(t, root, "commit", "-m", message)
}

// initRepoOnMain gives every change-resolution test the same starting point: a
// repo on `main` with one committed file and a clean tree.
func initRepoOnMain(t *testing.T) string {
	t.Helper()
	root := initRepo(t)
	writeFile(t, root, "go.mod", "module example.com/repo\n")
	writeFile(t, root, "internal/app/app.go", "package app\n\nfunc Run() {}\n")
	gitAdd(t, root, "go.mod", "internal/app/app.go")
	gitCommitMessage(t, root, "initial")
	runGit(t, root, "branch", "-M", "main")
	return root
}

// commitOnFeatureBranch adds the committed work a branch review should find.
func commitOnFeatureBranch(t *testing.T, root string) {
	t.Helper()
	runGit(t, root, "checkout", "-b", "feature")
	writeFile(t, root, "internal/app/feature.go", "package app\n\nfunc Feature() error { return nil }\n")
	gitAdd(t, root, "internal/app/feature.go")
	gitCommitMessage(t, root, "add feature")
}

// (a) an explicit base wins over everything, including a dirty working tree.
func TestResolveChangeSetUsesExplicitBaseRange(t *testing.T) {
	root := initRepoOnMain(t)
	commitOnFeatureBranch(t, root)
	writeFile(t, root, "internal/app/scratch.go", "package app\n")

	changes := resolveChangeSet(context.Background(), root, "main")
	if changes.Mode != ReviewModeRange {
		t.Fatalf("Mode = %q, want %q", changes.Mode, ReviewModeRange)
	}
	if changes.Range != "main...HEAD" || changes.Base != "main" {
		t.Fatalf("Range/Base = %q/%q, want main...HEAD/main", changes.Range, changes.Base)
	}
	if len(changes.Files) != 1 || changes.Files[0] != "internal/app/feature.go" {
		t.Fatalf("Files = %#v, want the committed feature file only", changes.Files)
	}
}

// (b) a dirty working tree stays the interactive default.
func TestResolveChangeSetPrefersWorkingTreeWhenDirty(t *testing.T) {
	root := initRepoOnMain(t)
	commitOnFeatureBranch(t, root)
	writeFile(t, root, "internal/app/app.go", "package app\n\nfunc Run() { println(1) }\n")

	changes := resolveChangeSet(context.Background(), root, "")
	if changes.Mode != ReviewModeWorkingTree {
		t.Fatalf("Mode = %q, want %q", changes.Mode, ReviewModeWorkingTree)
	}
	if changes.Range != "" {
		t.Fatalf("Range = %q, want empty for a working-tree review", changes.Range)
	}
	if len(changes.Files) != 1 || changes.Files[0] != "internal/app/app.go" {
		t.Fatalf("Files = %#v, want the dirty file only", changes.Files)
	}
}

// (c) a clean tree falls back to the branch's own commits.
func TestResolveChangeSetFallsBackToDetectedBaseRange(t *testing.T) {
	root := initRepoOnMain(t)
	commitOnFeatureBranch(t, root)

	changes := resolveChangeSet(context.Background(), root, "")
	if changes.Mode != ReviewModeRange {
		t.Fatalf("Mode = %q, want %q", changes.Mode, ReviewModeRange)
	}
	if changes.Range != "main...HEAD" {
		t.Fatalf("Range = %q, want main...HEAD", changes.Range)
	}
	if len(changes.Files) != 1 || changes.Files[0] != "internal/app/feature.go" {
		t.Fatalf("Files = %#v, want the committed feature file", changes.Files)
	}
}

// (c) on the base branch itself there is no branch to compare, so the last
// commit is the change worth reviewing.
func TestResolveChangeSetReviewsLastCommitOnBaseBranch(t *testing.T) {
	root := initRepoOnMain(t)
	writeFile(t, root, "internal/app/second.go", "package app\n\nfunc Second() {}\n")
	gitAdd(t, root, "internal/app/second.go")
	gitCommitMessage(t, root, "second")

	changes := resolveChangeSet(context.Background(), root, "")
	if changes.Mode != ReviewModeRange {
		t.Fatalf("Mode = %q, want %q", changes.Mode, ReviewModeRange)
	}
	if changes.Range != "HEAD~1...HEAD" {
		t.Fatalf("Range = %q, want HEAD~1...HEAD", changes.Range)
	}
	if len(changes.Files) != 1 || changes.Files[0] != "internal/app/second.go" {
		t.Fatalf("Files = %#v, want the last commit's file", changes.Files)
	}
}

// (d) a single commit on the base branch really is nothing to review.
func TestResolveChangeSetReportsNothingToReview(t *testing.T) {
	root := initRepoOnMain(t)

	changes := resolveChangeSet(context.Background(), root, "")
	if changes.Mode != ReviewModeNone || changes.Reviewed() {
		t.Fatalf("Mode = %q, Reviewed = %v; want %q and false", changes.Mode, changes.Reviewed(), ReviewModeNone)
	}
	if !strings.Contains(changes.Target, "main") || !strings.Contains(changes.Target, "HEAD~1") {
		t.Fatalf("Target = %q, want the refs gx looked at", changes.Target)
	}
}

// (a) an explicit base that resolves to an empty range is nothing to review,
// not a silent fallback to another source of changes.
func TestResolveChangeSetExplicitBaseWithEmptyRangeIsNothingToReview(t *testing.T) {
	root := initRepoOnMain(t)

	changes := resolveChangeSet(context.Background(), root, "HEAD")
	if changes.Mode != ReviewModeNone {
		t.Fatalf("Mode = %q, want %q", changes.Mode, ReviewModeNone)
	}
	if changes.Range != "HEAD...HEAD" {
		t.Fatalf("Range = %q, want the range that was asked for", changes.Range)
	}
}

func TestRangeDiffSnippetsComeFromTheRange(t *testing.T) {
	root := initRepoOnMain(t)
	commitOnFeatureBranch(t, root)

	changes := resolveChangeSet(context.Background(), root, "")
	snippets := collectDiffSnippets(context.Background(), root, changes.Files, Options{}, changes.Range)
	if len(snippets) != 1 {
		t.Fatalf("snippets = %#v, want one range diff", snippets)
	}
	if !strings.Contains(snippets[0].Diff, "+func Feature() error") {
		t.Fatalf("snippet = %#v, want the committed diff, not whole-file content", snippets[0])
	}
}

// The bug this fixes: a repo whose work is committed and whose tree is clean
// used to review nothing and report a clean bill of health.
func TestReviewOnCleanCommittedRepoReviewsTheCommitRange(t *testing.T) {
	root := initRepoOnMain(t)
	commitOnFeatureBranch(t, root)

	reviewer := &capturingReviewer{}
	engine := NewEngineWithReviewer(LocalScanner{}, StaticCatalog{}, defaultRules(), fakeRetriever{}, reviewer)
	report, err := engine.Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if !report.Reviewed || report.ReviewMode != ReviewModeRange {
		t.Fatalf("Reviewed/ReviewMode = %v/%q, want true/%q", report.Reviewed, report.ReviewMode, ReviewModeRange)
	}
	if report.ReviewRange != "main...HEAD" || report.ReviewBase != "main" {
		t.Fatalf("ReviewRange/ReviewBase = %q/%q, want main...HEAD/main", report.ReviewRange, report.ReviewBase)
	}
	if len(report.ChangedFiles) != 1 || report.ChangedFiles[0] != "internal/app/feature.go" {
		t.Fatalf("ChangedFiles = %#v, want the committed file", report.ChangedFiles)
	}
	if text := RenderMarkdown(report); strings.Contains(text, "Nothing to review") {
		t.Fatalf("RenderMarkdown() reported nothing to review for a committed change:\n%s", text)
	}
	if reviewer.brief.Static.ReviewRange != "main...HEAD" {
		t.Fatalf("brief review range = %q, want main...HEAD", reviewer.brief.Static.ReviewRange)
	}
	if len(reviewer.brief.Static.DiffSnippets) != 1 || !strings.Contains(reviewer.brief.Static.DiffSnippets[0].Diff, "+func Feature() error") {
		t.Fatalf("brief diff snippets = %#v, want the committed diff", reviewer.brief.Static.DiffSnippets)
	}
}

func TestReviewWithExplicitBaseReviewsThatRange(t *testing.T) {
	root := initRepoOnMain(t)
	commitOnFeatureBranch(t, root)
	// A dirty file the explicit base must ignore.
	writeFile(t, root, "internal/app/scratch.go", "package app\n")

	report, err := Review(context.Background(), root, Options{Base: "main"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.ReviewRange != "main...HEAD" {
		t.Fatalf("ReviewRange = %q, want main...HEAD", report.ReviewRange)
	}
	if len(report.ChangedFiles) != 1 || report.ChangedFiles[0] != "internal/app/feature.go" {
		t.Fatalf("ChangedFiles = %#v, want the range files only", report.ChangedFiles)
	}
}

func TestReviewReportsNothingToReviewInsteadOfCleanBill(t *testing.T) {
	root := initRepoOnMain(t)

	report, err := Review(context.Background(), root, Options{})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if report.Reviewed || report.ReviewMode != ReviewModeNone {
		t.Fatalf("Reviewed/ReviewMode = %v/%q, want false/%q", report.Reviewed, report.ReviewMode, ReviewModeNone)
	}
	if len(report.Findings) != 0 {
		t.Fatalf("Findings = %#v, want none when nothing was inspected", report.Findings)
	}
	text := RenderMarkdown(report)
	if strings.Contains(text, "No material issues found in this change.") {
		t.Fatalf("RenderMarkdown() reported a clean review for an uninspected repo:\n%s", text)
	}
	if !strings.Contains(text, "Nothing to review") || !strings.Contains(text, "not a clean review") {
		t.Fatalf("RenderMarkdown() missing the nothing-to-review outcome:\n%s", text)
	}
}

// A scope-directed review is a whole-repo review by request, so a repo with no
// diff at all is not "nothing to review" there.
func TestScopeDirectedReviewOnCleanRepoStaysWholeRepo(t *testing.T) {
	root := initRepo(t)
	writeFile(t, root, "package.json", "{\n  \"name\": \"example\"\n}\n")
	gitAdd(t, root, "package.json")
	gitCommitMessage(t, root, "initial")
	runGit(t, root, "branch", "-M", "main")

	report, err := Review(context.Background(), root, Options{Scope: "dependencies"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if !report.Reviewed || report.ReviewMode != ReviewModeRepo {
		t.Fatalf("Reviewed/ReviewMode = %v/%q, want true/%q", report.Reviewed, report.ReviewMode, ReviewModeRepo)
	}
	if !hasFinding(report.Findings, "dependencies.javascript-unlocked") {
		t.Fatalf("Findings = %#v, want repo-wide dependency finding", report.Findings)
	}
}
