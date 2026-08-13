package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
)

func runConstraintsCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewRoot(context.Background())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(append([]string{"constraints"}, args...))
	err := cmd.Execute()
	return out.String(), err
}

// fixtureAWSKey is split mid-literal so secret scanners (GitHub push
// protection included) never see a contiguous key shape in this source file.
var fixtureAWSKey = "AKIA" + "ABCDEFGHIJKLMNOP"

func commitConstraintsSecret(t *testing.T, root string) {
	t.Helper()
	runGitTest(t, root, "checkout", "-b", "feature")
	writeTestFile(t, root, "internal/app/creds.go", "package app\n\nconst awsKey = \""+fixtureAWSKey+"\"\n")
	gitAddTestFiles(t, root, "internal/app/creds.go")
	runGitTest(t, root, "commit", "-m", "add creds")
}

func TestConstraintsExitsThreeOnNoShip(t *testing.T) {
	root := newReviewGateRepo(t)
	commitConstraintsSecret(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runConstraintsCommand(t)
	if err == nil {
		t.Fatalf("expected a no-ship error, output:\n%s", out)
	}
	if code := ExitCode(err); code != reviewFindingsExitCode {
		t.Fatalf("ExitCode() = %d, want %d (%v)", code, reviewFindingsExitCode, err)
	}
	if !strings.Contains(out, "Verdict: NO-SHIP") {
		t.Fatalf("output missing the verdict line:\n%s", out)
	}
}

func TestConstraintsExitsZeroOnShip(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runConstraintsCommand(t)
	if err != nil {
		t.Fatalf("gx constraints error = %v\n%s", err, out)
	}
	if !strings.Contains(out, "Verdict: SHIP") {
		t.Fatalf("output missing the ship verdict:\n%s", out)
	}
}

func TestConstraintsExitsFourOnNothingToCheck(t *testing.T) {
	root := newReviewGateRepo(t)
	t.Chdir(root)
	setReviewGateEnv(t)

	_, err := runConstraintsCommand(t)
	if code := ExitCode(err); code != reviewNothingToReviewExitCode {
		t.Fatalf("ExitCode() = %d, want %d (%v)", code, reviewNothingToReviewExitCode, err)
	}
}

func TestConstraintsReportOnlyExitsZero(t *testing.T) {
	root := newReviewGateRepo(t)
	commitConstraintsSecret(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runConstraintsCommand(t, "--report-only")
	if err != nil {
		t.Fatalf("--report-only must exit clean, got %v\n%s", err, out)
	}
	if !strings.Contains(out, "Verdict: NO-SHIP") {
		t.Fatalf("--report-only lost the verdict:\n%s", out)
	}
}

func TestConstraintsJSONRoundTrips(t *testing.T) {
	root := newReviewGateRepo(t)
	commitConstraintsSecret(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runConstraintsCommand(t, "--json", "--report-only")
	if err != nil {
		t.Fatalf("gx constraints --json error = %v\n%s", err, out)
	}
	var report codereview.ConstraintsReport
	if decodeErr := json.Unmarshal([]byte(out), &report); decodeErr != nil {
		t.Fatalf("json.Unmarshal() error = %v, output:\n%s", decodeErr, out)
	}
	if report.Verdict != codereview.VerdictNoShip {
		t.Fatalf("verdict = %q, want no-ship", report.Verdict)
	}
	if len(report.Gates) != len(codereview.AllGateIDs()) {
		t.Fatalf("gates = %d, want %d", len(report.Gates), len(codereview.AllGateIDs()))
	}
	if strings.Contains(out, "Verdict:") && !strings.Contains(out, "\"verdict\"") {
		t.Fatalf("--json emitted the text render:\n%s", out)
	}
	// Files ride in the JSON unconditionally; --verbose only affects text.
	for _, gate := range report.Gates {
		if gate.Gate == codereview.GateSecurity && len(gate.Files) == 0 {
			t.Fatalf("security gate files missing from JSON: %#v", gate)
		}
	}
}

func TestConstraintsMarkdownEmitsTheTable(t *testing.T) {
	root := newReviewGateRepo(t)
	commitConstraintsSecret(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	out, err := runConstraintsCommand(t, "--md", "--report-only")
	if err != nil {
		t.Fatalf("gx constraints --md error = %v\n%s", err, out)
	}
	for _, want := range []string{"| # | Gate | Status | Evidence |", "❌ FAIL", "**Verdict: NO-SHIP"} {
		if !strings.Contains(out, want) {
			t.Fatalf("markdown output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("markdown output carries ANSI escapes:\n%s", out)
	}
}

func TestConstraintsRejectsUnknownSkipGate(t *testing.T) {
	root := newReviewGateRepo(t)
	commitOnBranch(t, root)
	t.Chdir(root)
	setReviewGateEnv(t)

	_, err := runConstraintsCommand(t, "--skip-gates", "correctness,typo")
	if err == nil || !strings.Contains(err.Error(), "unknown constraints gate") {
		t.Fatalf("err = %v, want the unknown-gate rejection", err)
	}
}
