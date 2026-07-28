package cli

import (
	"errors"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
	"github.com/satoricorp/gx/internal/vcs"
)

func gateExitCode(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return 0
	}
	var coded *vcs.CodedError
	if !errors.As(err, &coded) {
		t.Fatalf("expected a coded error, got %v", err)
	}
	return coded.Code
}

// The gate decided pass/fail purely from findings, so a review in which no
// model ran exited 0 with zero findings and the pull request merged reporting a
// review that never happened. The rendered report carries an "AI review
// unavailable" banner, but an exit code is the only thing a CI step reads.
func TestReviewGateFailsWhenNoModelRan(t *testing.T) {
	report := codereview.Report{
		Reviewed:        true,
		ReviewTarget:    "the working tree",
		ReviewMode:      codereview.ReviewModeWorkingTree,
		DegradedReasons: []string{"AI review unavailable (no reviewer configured)"},
	}
	err := reviewGateError(report, codereview.FailOnStrong)
	if got := gateExitCode(t, err); got != reviewDegradedExitCode {
		t.Fatalf("exit code = %d, want %d (err=%v)", got, reviewDegradedExitCode, err)
	}
}

// Reading part of the subject cannot answer "no findings at or above X" for the
// part that was never read.
func TestReviewGateFailsOnPartialCoverage(t *testing.T) {
	report := codereview.Report{
		Reviewed:     true,
		ReviewTarget: "the working tree",
		ReviewMode:   codereview.ReviewModeWorkingTree,
		Coverage: codereview.Coverage{
			Planned: true,
			Total:   10,
			Read:    3,
			Units:   "files",
		},
	}
	err := reviewGateError(report, codereview.FailOnStrong)
	if got := gateExitCode(t, err); got != reviewDegradedExitCode {
		t.Fatalf("exit code = %d, want %d (err=%v)", got, reviewDegradedExitCode, err)
	}
}

// A degraded review must not be reported as a findings failure: the operator
// response differs (retry the review vs fix the code), so the codes stay
// distinct even when both conditions hold.
func TestReviewGatePrefersDegradedOverFindings(t *testing.T) {
	report := codereview.Report{
		Reviewed:        true,
		ReviewTarget:    "the working tree",
		ReviewMode:      codereview.ReviewModeWorkingTree,
		DegradedReasons: []string{"2 of 4 parallel reviews failed"},
		Findings:        []codereview.Finding{{Strength: "Blocking"}},
	}
	err := reviewGateError(report, codereview.FailOnStrong)
	if got := gateExitCode(t, err); got != reviewDegradedExitCode {
		t.Fatalf("exit code = %d, want %d (err=%v)", got, reviewDegradedExitCode, err)
	}
}

// A complete review with nothing to report still passes: the gate must not
// become unusable.
func TestReviewGatePassesCleanCompleteReview(t *testing.T) {
	report := codereview.Report{
		Reviewed:     true,
		ReviewTarget: "the working tree",
		ReviewMode:   codereview.ReviewModeWorkingTree,
		Coverage: codereview.Coverage{
			Planned: true,
			Total:   4,
			Read:    4,
			Units:   "files",
		},
	}
	if err := reviewGateError(report, codereview.FailOnStrong); err != nil {
		t.Fatalf("clean review should pass the gate, got %v", err)
	}
}

// Without --fail-on the command is advisory, so a degraded run must not start
// failing builds that never opted into gating.
func TestReviewGateSilentWhenDisabled(t *testing.T) {
	report := codereview.Report{
		Reviewed:        true,
		DegradedReasons: []string{"AI review unavailable"},
	}
	if err := reviewGateError(report, codereview.FailOnNone); err != nil {
		t.Fatalf("disabled gate should return nil, got %v", err)
	}
}

// Nothing-to-review keeps its own exit code: "there was no diff" and "the
// review broke" are different operator problems.
func TestReviewGateKeepsNothingToReviewDistinct(t *testing.T) {
	report := codereview.Report{Reviewed: false, ReviewTarget: "the working tree"}
	err := reviewGateError(report, codereview.FailOnStrong)
	if got := gateExitCode(t, err); got != reviewNothingToReviewExitCode {
		t.Fatalf("exit code = %d, want %d (err=%v)", got, reviewNothingToReviewExitCode, err)
	}
}
