package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
)

// newReviewGateRepo builds a repo whose work is committed and whose tree is
// clean: the case that used to report a clean bill of health without reading
// anything.
func newReviewGateRepo(t *testing.T) string {
	t.Helper()
	root := initGitRepo(t)
	runGitTest(t, root, "config", "user.name", "Test User")
	runGitTest(t, root, "config", "user.email", "test@example.com")
	writeTestFile(t, root, "go.mod", "module example.com/repo\n")
	writeTestFile(t, root, "internal/app/app.go", "package app\n\nfunc Run() {}\n")
	gitAddTestFiles(t, root, "go.mod", "internal/app/app.go")
	runGitTest(t, root, "commit", "-m", "initial")
	return root
}

func commitOnBranch(t *testing.T, root string) {
	t.Helper()
	runGitTest(t, root, "checkout", "-b", "feature")
	writeTestFile(t, root, "internal/app/feature.go", "package app\n\nfunc Feature() error { return nil }\n")
	gitAddTestFiles(t, root, "internal/app/feature.go")
	runGitTest(t, root, "commit", "-m", "add feature")
}

func setReviewGateEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GX_HOME", t.TempDir())
	t.Setenv("GX_API_URL", "")
	t.Setenv("GX_UPLOAD_TOKEN", "")
	t.Setenv("GX_REVIEW_AI", "0")
	t.Setenv("GX_REVIEW_JUDGE", "0")
	t.Setenv("GX_REVIEW_STATIC_TOOLS", "0")
	t.Setenv("GX_REVIEW_RESOURCES", "0")
	t.Setenv("GX_REVIEW_INDEXED_CONTEXT", "0")
	t.Setenv("GX_REVIEW_HISTORY_CONTEXT", "0")
	denyReviewNetwork(t)
}

// denyReviewNetwork makes a review test hermetic.
//
// These tests inherit the developer's shell, and a developer working on gx has
// OPENAI_API_KEY and TURBOPUFFER_API_KEY exported. Every review network path is
// gated on those keys being present, so with them exported the gate tests were
// not testing the offline path at all: each one embedded the fixture repository
// through the real embeddings API and upserted it into the production
// TurboPuffer account. The namespace is derived from the fixture's temp
// directory (semantic.NamespaceForRepo), so every `go test ./internal/cli` run
// created a brand new `gx-local-*` namespace that nothing would ever delete.
//
// Clearing the credentials is the mechanism rather than a per-feature kill
// switch such as GX_REVIEW_CODE_INDEX=0, because it shuts every door at once:
// a retriever added tomorrow is off here for the same reason it is off in CI,
// without anyone having to remember to add its switch to this list.
//
// The endpoints then point at a tripwire instead of a dead port, so a path that
// dials out despite having no credentials fails the test by name rather than
// hanging on a connect timeout and being read as slowness.
func denyReviewNetwork(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"OPENAI_API_KEY",
		"GX_OPENAI_API_KEY",
		"TURBOPUFFER_API_KEY",
		"ANTHROPIC_API_KEY",
		"GX_TPUF_NAMESPACE",
	} {
		t.Setenv(key, "")
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("review test called %s %s; review tests must not reach a real API", r.Method, r.URL.Path)
		http.Error(w, "network denied in tests", http.StatusForbidden)
	}))
	t.Cleanup(server.Close)
	t.Setenv("GX_OPENAI_BASE_URL", server.URL)
	t.Setenv("GX_TPUF_BASE_URL", server.URL)
}

func runEnhanceCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(append([]string{"enhance"}, args...))
	err := cmd.Execute()
	return out.String(), err
}

func TestEnhanceJSONReportsTheResolvedCommitRange(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance --json error = %v\n%s", err, out)
	}
	var report codereview.Report
	if decodeErr := json.Unmarshal([]byte(out), &report); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, out)
	}
	if !report.Reviewed {
		t.Fatalf("reviewed = false, want the committed range to be reviewed:\n%s", out)
	}
	if report.ReviewMode != codereview.ReviewModeRange || report.ReviewRange != "main...HEAD" {
		t.Fatalf("review_mode/review_range = %q/%q, want range/main...HEAD:\n%s", report.ReviewMode, report.ReviewRange, out)
	}
	if len(report.ChangedFiles) != 1 || report.ChangedFiles[0] != "internal/app/feature.go" {
		t.Fatalf("changed_files = %#v, want the committed file", report.ChangedFiles)
	}
	if strings.Contains(out, "## Recommendations") {
		t.Fatalf("--json emitted markdown as well:\n%s", out)
	}
	// The gate contract: every key a CI step reads must be present.
	var raw map[string]any
	if decodeErr := json.Unmarshal([]byte(out), &raw); decodeErr != nil {
		t.Fatalf("json.Unmarshal(raw) error = %v", decodeErr)
	}
	for _, key := range []string{"reviewed", "review_mode", "review_range", "findings"} {
		if _, ok := raw[key]; !ok {
			t.Fatalf("JSON report missing %q key:\n%s", key, out)
		}
	}
}

func TestEnhanceJSONMarksNothingToReview(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance --json error = %v\n%s", err, out)
	}
	var report codereview.Report
	if decodeErr := json.Unmarshal([]byte(out), &report); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, out)
	}
	if report.Reviewed || report.ReviewMode != codereview.ReviewModeNone {
		t.Fatalf("reviewed/review_mode = %v/%q, want false/none:\n%s", report.Reviewed, report.ReviewMode, out)
	}
	if strings.TrimSpace(report.ReviewTarget) == "" {
		t.Fatalf("review_target is empty, want the refs gx looked at:\n%s", out)
	}
}

func TestEnhanceNothingToReviewDoesNotClaimACleanReview(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance error = %v\n%s", err, out)
	}
	if strings.Contains(out, "No material issues found in this change.") {
		t.Fatalf("gx enhance claimed a clean review without inspecting anything:\n%s", out)
	}
	if !strings.Contains(out, "Nothing to review") {
		t.Fatalf("gx enhance output missing the nothing-to-review outcome:\n%s", out)
	}
}

func TestEnhanceFailOnExitsNonZeroForSurvivingFindings(t *testing.T) {
	root := newReviewGateRepo(t)
	writeTestFile(t, root, "package.json", "{\n  \"name\": \"example\"\n}\n")
	gitAddTestFiles(t, root, "package.json")
	runGitTest(t, root, "commit", "-m", "add package.json")
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--scope", "dependencies", "--fail-on", "strong", "--no-publish")
	if err == nil {
		t.Fatalf("gx enhance --fail-on strong exited 0 with findings:\n%s", out)
	}
	if code := ExitCode(err); code != reviewFindingsExitCode {
		t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewFindingsExitCode, err)
	}

	// The same review under a stricter threshold has nothing blocking.
	if out, err := runEnhanceCommand(t, "--scope", "dependencies", "--fail-on", "blocking", "--no-publish"); err != nil {
		t.Fatalf("gx enhance --fail-on blocking error = %v\n%s", err, out)
	}
	// And the default gate never fails.
	if out, err := runEnhanceCommand(t, "--scope", "dependencies", "--no-publish"); err != nil {
		t.Fatalf("gx enhance without --fail-on error = %v\n%s", err, out)
	}
}

// An explicit gate treats "never looked" as a failure: passing a diff nobody
// opened is the silent false pass --fail-on exists to prevent.
func TestEnhanceFailOnExitsDistinctlyWhenNothingWasReviewed(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--fail-on", "any", "--no-publish")
	if err == nil {
		t.Fatalf("gx enhance --fail-on any exited 0 without reviewing anything:\n%s", out)
	}
	if code := ExitCode(err); code != reviewNothingToReviewExitCode {
		t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewNothingToReviewExitCode, err)
	}
	if !strings.Contains(err.Error(), "nothing was reviewed") {
		t.Fatalf("error = %v, want the nothing-to-review message", err)
	}

	// Without a gate the same run stays exit 0 for interactive use.
	if _, err := runEnhanceCommand(t, "--no-publish"); err != nil {
		t.Fatalf("gx enhance without --fail-on error = %v", err)
	}
}

// A repository with no commits and no files: nothing to review is still the
// honest answer, and an explicit gate still refuses to pass on it. --repo is
// included on purpose — a flag that could turn the one non-passing outcome
// into a pass by asking for a bigger subject would be a gate with a bypass.
func TestEnhanceOnAnEmptyRepoNeverPassesAGate(t *testing.T) {
	for _, flags := range [][]string{nil, {"--repo"}} {
		name := "plain"
		if len(flags) > 0 {
			name = strings.Join(flags, " ")
		}
		t.Run(name, func(t *testing.T) {
			root := initGitRepo(t)
			runGitTest(t, root, "config", "user.name", "Test User")
			runGitTest(t, root, "config", "user.email", "test@example.com")
			t.Chdir(root)
			setReviewGateEnv(t)

			out, err := runEnhanceCommand(t, append(append([]string{}, flags...), "--json", "--no-publish")...)
			if err != nil {
				t.Fatalf("gx enhance error = %v\n%s", err, out)
			}
			var report codereview.Report
			if decodeErr := json.Unmarshal([]byte(out), &report); decodeErr != nil {
				t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, out)
			}
			if report.Reviewed || report.ReviewMode != codereview.ReviewModeNone {
				t.Fatalf("reviewed/review_mode = %v/%q, want false/%q:\n%s", report.Reviewed, report.ReviewMode, codereview.ReviewModeNone, out)
			}

			gateOut, gateErr := runEnhanceCommand(t, append(append([]string{}, flags...), "--fail-on", "any", "--no-publish")...)
			if gateErr == nil {
				t.Fatalf("gx enhance --fail-on any passed on an empty repo:\n%s", gateOut)
			}
			if code := ExitCode(gateErr); code != reviewNothingToReviewExitCode {
				t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewNothingToReviewExitCode, gateErr)
			}
		})
	}
}

// The gx Cloud history row and the PostHog event both label the run, and they
// derive that label from the same place so they cannot disagree with each
// other — or, more importantly, with the summary text stored beside them. A
// run whose summary says "Reviewed the repository" must not be filed as a
// patch review.
func TestEnhanceRunModeNamesTheSubjectItReviewed(t *testing.T) {
	cases := []struct {
		name          string
		prompt        string
		scopeExplicit bool
		deep          bool
		wholeRepo     bool
		want          string
	}{
		{name: "default", want: "patch"},
		{name: "scope", scopeExplicit: true, want: "scope"},
		{name: "prompt", prompt: "why is auth slow?", want: "prompt"},
		{name: "deep", deep: true, want: "deep"},
		{name: "whole repo", wholeRepo: true, want: "repo"},
		// The subject wins over the depth and the lens: those are recorded as
		// their own fields, the mode is the one label.
		{name: "whole repo and deep", wholeRepo: true, deep: true, want: "repo"},
		{name: "whole repo and prompt", wholeRepo: true, prompt: "is there a race?", want: "repo"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := reviewRunMode(tc.prompt, tc.scopeExplicit, tc.deep, tc.wholeRepo); got != tc.want {
				t.Fatalf("reviewRunMode() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestEnhanceBaseFlagSelectsTheRange(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	// A dirty file the explicit base must ignore.
	writeTestFile(t, root, "internal/app/scratch.go", "package app\n")
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--base", "main", "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance --base error = %v\n%s", err, out)
	}
	var report codereview.Report
	if decodeErr := json.Unmarshal([]byte(out), &report); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, out)
	}
	if report.ReviewRange != "main...HEAD" || report.ReviewBase != "main" {
		t.Fatalf("review_range/review_base = %q/%q, want main...HEAD/main", report.ReviewRange, report.ReviewBase)
	}
	if len(report.ChangedFiles) != 1 || report.ChangedFiles[0] != "internal/app/feature.go" {
		t.Fatalf("changed_files = %#v, want the range files only", report.ChangedFiles)
	}
}

// The reported case, end to end through the CLI: uncommitted work used to make
// the diff the only reachable subject, so there was no way to ask for a review
// of the codebase. --repo is that way.
func TestEnhanceRepoFlagReviewsTheRepositoryWithADirtyTree(t *testing.T) {
	root := newReviewGateRepo(t)
	writeTestFile(t, root, "internal/app/app.go", "package app\n\nfunc Run() { println(1) }\n")
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--repo", "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance --repo error = %v\n%s", err, out)
	}
	var report codereview.Report
	if decodeErr := json.Unmarshal([]byte(out), &report); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, out)
	}
	if !report.Reviewed || report.ReviewMode != codereview.ReviewModeRepo {
		t.Fatalf("reviewed/review_mode = %v/%q, want true/%q:\n%s", report.Reviewed, report.ReviewMode, codereview.ReviewModeRepo, out)
	}
	if !strings.Contains(report.ReviewTarget, "the repository") {
		t.Fatalf("review_target = %q, want it to name the repository", report.ReviewTarget)
	}
	// The uncommitted work stays in focus rather than being discarded.
	if len(report.ChangedFiles) != 1 || report.ChangedFiles[0] != "internal/app/app.go" {
		t.Fatalf("changed_files = %#v, want the uncommitted file", report.ChangedFiles)
	}

	// Without the flag the same tree is a working-tree review, unchanged.
	plain, err := runEnhanceCommand(t, "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance error = %v\n%s", err, plain)
	}
	var plainReport codereview.Report
	if decodeErr := json.Unmarshal([]byte(plain), &plainReport); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, plain)
	}
	if plainReport.ReviewMode != codereview.ReviewModeWorkingTree {
		t.Fatalf("review_mode = %q, want %q for a plain review of a dirty tree", plainReport.ReviewMode, codereview.ReviewModeWorkingTree)
	}
}

// --repo works on a clean repo too, where a plain review has nothing to read.
func TestEnhanceRepoFlagReviewsTheRepositoryWithACleanTree(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--repo", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance --repo error = %v\n%s", err, out)
	}
	if strings.Contains(out, "Nothing to review") {
		t.Fatalf("gx enhance --repo reported nothing to review:\n%s", out)
	}
	if !strings.Contains(out, "Reviewed the repository") {
		t.Fatalf("gx enhance --repo does not say what it reviewed:\n%s", out)
	}

	// Saying it reviewed the repository is only honest if it read the
	// repository, so the report has to show context it could only have got
	// from there: this repo has no diff of any kind.
	jsonOut, err := runEnhanceCommand(t, "--repo", "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance --repo --json error = %v\n%s", err, jsonOut)
	}
	var report codereview.Report
	if decodeErr := json.Unmarshal([]byte(jsonOut), &report); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, jsonOut)
	}
	if len(report.ChangedFiles) != 0 {
		t.Fatalf("changed_files = %#v, want none for a clean repo", report.ChangedFiles)
	}
	if report.ContextSnippets == 0 {
		t.Fatalf("context_snippets = 0: --repo claimed to review the repository without reading it:\n%s", jsonOut)
	}

	// The same repo without the flag is still honest about reading nothing.
	plain, err := runEnhanceCommand(t, "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance error = %v\n%s", err, plain)
	}
	if !strings.Contains(plain, "Nothing to review") {
		t.Fatalf("gx enhance lost the nothing-to-review outcome:\n%s", plain)
	}
}

// --repo names the subject; --base still picks the diff. The report says which
// instruction did what instead of silently dropping one of them.
func TestEnhanceRepoFlagWinsOverBaseAndSaysSo(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runEnhanceCommand(t, "--repo", "--base", "main", "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx enhance --repo --base error = %v\n%s", err, out)
	}
	var report codereview.Report
	if decodeErr := json.Unmarshal([]byte(out), &report); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, out)
	}
	if report.ReviewMode != codereview.ReviewModeRepo {
		t.Fatalf("review_mode = %q, want %q: --repo sets the subject", report.ReviewMode, codereview.ReviewModeRepo)
	}
	if report.ReviewBase != "main" || report.ReviewRange != "main...HEAD" {
		t.Fatalf("review_base/review_range = %q/%q, want main/main...HEAD kept in focus", report.ReviewBase, report.ReviewRange)
	}
	for _, want := range []string{"the repository", "--base", "main...HEAD"} {
		if !strings.Contains(report.ReviewTarget, want) {
			t.Fatalf("review_target = %q, want it to mention %q", report.ReviewTarget, want)
		}
	}
}

func TestEnhanceRejectsUnknownFailOnLevel(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	_, err := runEnhanceCommand(t, "--fail-on", "critical", "--no-publish")
	if err == nil {
		t.Fatal("gx enhance accepted an unknown --fail-on level")
	}
	if !strings.Contains(err.Error(), "unsupported --fail-on level") {
		t.Fatalf("error = %v, want an unsupported-level error", err)
	}
}
