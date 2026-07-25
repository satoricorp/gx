package cli

import (
	"bytes"
	"context"
	"encoding/json"
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
}

func runReviewCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(append([]string{"review"}, args...))
	err := cmd.Execute()
	return out.String(), err
}

func TestReviewJSONReportsTheResolvedCommitRange(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runReviewCommand(t, "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx review --json error = %v\n%s", err, out)
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

func TestReviewJSONMarksNothingToReview(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runReviewCommand(t, "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx review --json error = %v\n%s", err, out)
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

func TestReviewNothingToReviewDoesNotClaimACleanReview(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runReviewCommand(t, "--no-publish")
	if err != nil {
		t.Fatalf("gx review error = %v\n%s", err, out)
	}
	if strings.Contains(out, "No material issues found in this change.") {
		t.Fatalf("gx review claimed a clean review without inspecting anything:\n%s", out)
	}
	if !strings.Contains(out, "Nothing to review") {
		t.Fatalf("gx review output missing the nothing-to-review outcome:\n%s", out)
	}
}

func TestReviewFailOnExitsNonZeroForSurvivingFindings(t *testing.T) {
	root := newReviewGateRepo(t)
	writeTestFile(t, root, "package.json", "{\n  \"name\": \"example\"\n}\n")
	gitAddTestFiles(t, root, "package.json")
	runGitTest(t, root, "commit", "-m", "add package.json")
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runReviewCommand(t, "--scope", "dependencies", "--fail-on", "strong", "--no-publish")
	if err == nil {
		t.Fatalf("gx review --fail-on strong exited 0 with findings:\n%s", out)
	}
	if code := ExitCode(err); code != reviewFindingsExitCode {
		t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewFindingsExitCode, err)
	}

	// The same review under a stricter threshold has nothing blocking.
	if out, err := runReviewCommand(t, "--scope", "dependencies", "--fail-on", "blocking", "--no-publish"); err != nil {
		t.Fatalf("gx review --fail-on blocking error = %v\n%s", err, out)
	}
	// And the default gate never fails.
	if out, err := runReviewCommand(t, "--scope", "dependencies", "--no-publish"); err != nil {
		t.Fatalf("gx review without --fail-on error = %v\n%s", err, out)
	}
}

// An explicit gate treats "never looked" as a failure: passing a diff nobody
// opened is the silent false pass --fail-on exists to prevent.
func TestReviewFailOnExitsDistinctlyWhenNothingWasReviewed(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runReviewCommand(t, "--fail-on", "any", "--no-publish")
	if err == nil {
		t.Fatalf("gx review --fail-on any exited 0 without reviewing anything:\n%s", out)
	}
	if code := ExitCode(err); code != reviewNothingToReviewExitCode {
		t.Fatalf("ExitCode() = %d, want %d (error: %v)", code, reviewNothingToReviewExitCode, err)
	}
	if !strings.Contains(err.Error(), "nothing was reviewed") {
		t.Fatalf("error = %v, want the nothing-to-review message", err)
	}

	// Without a gate the same run stays exit 0 for interactive use.
	if _, err := runReviewCommand(t, "--no-publish"); err != nil {
		t.Fatalf("gx review without --fail-on error = %v", err)
	}
}

func TestReviewBaseFlagSelectsTheRange(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	// A dirty file the explicit base must ignore.
	writeTestFile(t, root, "internal/app/scratch.go", "package app\n")
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runReviewCommand(t, "--base", "main", "--json", "--no-publish")
	if err != nil {
		t.Fatalf("gx review --base error = %v\n%s", err, out)
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

func TestReviewRejectsUnknownFailOnLevel(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	_, err := runReviewCommand(t, "--fail-on", "critical", "--no-publish")
	if err == nil {
		t.Fatal("gx review accepted an unknown --fail-on level")
	}
	if !strings.Contains(err.Error(), "unsupported --fail-on level") {
		t.Fatalf("error = %v, want an unsupported-level error", err)
	}
}
