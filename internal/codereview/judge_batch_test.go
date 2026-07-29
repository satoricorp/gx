package codereview

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

// recordingJudge confirms every candidate it is given and records the batches.
type recordingJudge struct {
	mu    sync.Mutex
	sizes []int
	// failFrom makes every batch whose first candidate index is >= failFrom
	// fail, so a test can exercise a partial verification failure.
	failAfterCalls int
	calls          int
}

func (j *recordingJudge) Available() bool { return true }

func (j *recordingJudge) Judge(_ context.Context, req judgeRequest) ([]judgeResult, error) {
	j.mu.Lock()
	j.sizes = append(j.sizes, len(req.Candidates))
	j.calls++
	call := j.calls
	j.mu.Unlock()
	if j.failAfterCalls > 0 && call > j.failAfterCalls {
		return nil, errors.New("judge batch failed")
	}
	results := make([]judgeResult, 0, len(req.Candidates))
	for _, candidate := range req.Candidates {
		results = append(results, judgeResult{
			CandidateID: candidate.ID,
			Verdict:     "confirmed",
			Impact:      impactBreaking,
			Severity:    4,
			Confidence:  0.9,
		})
	}
	return results, nil
}

func judgeBatchCandidates(n int) []Finding {
	out := make([]Finding, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, judgeTestFinding(fmt.Sprintf("openai.ai.review.%d", i), "internal/app/app.go"))
	}
	return out
}

// TestRunJudgeSplitsCandidatesIntoBoundedBatches pins the property that makes
// verification survivable at scale: no single judge call may carry more
// candidates than its output budget can answer for.
func TestRunJudgeSplitsCandidatesIntoBoundedBatches(t *testing.T) {
	judge := &recordingJudge{}
	candidates := judgeBatchCandidates(judgeBatchSize*2 + 5)

	outcome := runJudge(context.Background(), judge, ReviewContext{}, candidates)

	if outcome.Batches != 3 {
		t.Fatalf("Batches = %d, want 3 for %d candidates at batch size %d",
			outcome.Batches, len(candidates), judgeBatchSize)
	}
	for _, size := range judge.sizes {
		if size > judgeBatchSize {
			t.Fatalf("a judge call carried %d candidates, want at most %d", size, judgeBatchSize)
		}
	}
	if len(outcome.Judged) != len(candidates) {
		t.Fatalf("Judged = %d, want all %d candidates verified", len(outcome.Judged), len(candidates))
	}
	if outcome.BatchesFailed != 0 {
		t.Fatalf("BatchesFailed = %d, want 0", outcome.BatchesFailed)
	}
}

// TestRunJudgeKeepsCandidatesFromFailedBatches is the regression that matters
// most: a batch that could not be verified must not be reported the same way as
// a batch the judge rejected. Before batching, one truncated response dropped
// every finding at once.
func TestRunJudgeKeepsCandidatesFromFailedBatches(t *testing.T) {
	judge := &recordingJudge{failAfterCalls: 1}
	candidates := judgeBatchCandidates(judgeBatchSize * 2)

	outcome := runJudge(context.Background(), judge, ReviewContext{}, candidates)

	if outcome.BatchesFailed == 0 {
		t.Fatal("BatchesFailed = 0, want the failing batch to be reported")
	}
	if len(outcome.Unjudged) == 0 {
		t.Fatal("Unjudged is empty, want the failed batch's candidates kept rather than dropped")
	}
	if len(outcome.Judged)+len(outcome.Unjudged) != len(candidates) {
		t.Fatalf("Judged(%d) + Unjudged(%d) = %d, want every one of the %d candidates accounted for",
			len(outcome.Judged), len(outcome.Unjudged),
			len(outcome.Judged)+len(outcome.Unjudged), len(candidates))
	}
	if outcome.Err == nil {
		t.Fatal("Err is nil, want the batch failure surfaced for the degraded-reasons line")
	}
}

// partialJudge answers, successfully and in valid JSON, for only the first
// `answers` candidates of every batch — which is what this model actually does
// to list replies. The de-duplicator, running the same Sonnet, reports that
// shortfall on real reviews; before this was fixed the judge could not see it.
type partialJudge struct{ answers int }

func (partialJudge) Available() bool { return true }

func (j partialJudge) Judge(_ context.Context, req judgeRequest) ([]judgeResult, error) {
	var results []judgeResult
	for i, candidate := range req.Candidates {
		if i >= j.answers {
			break
		}
		results = append(results, judgeResult{
			CandidateID: candidate.ID,
			Verdict:     "confirmed",
			Impact:      impactBreaking,
			Severity:    4,
			Confidence:  0.9,
		})
	}
	return results, nil
}

// TestRunJudgeKeepsCandidatesTheJudgeNeverAnsweredFor is the delete-only
// invariant applied to the case that used to escape it.
//
// applyJudgeResults cannot tell "the judge rejected this" from "the judge never
// mentioned this" — both are simply absent from the verdict map — so a reply
// covering 4 of 10 candidates deleted the other 6 and reported BatchesFailed=0,
// which the engine reads as a clean, complete verification. Measured before the
// fix: candidates=10 judged=4 unjudged=0 failed=0 err=nil.
func TestRunJudgeKeepsCandidatesTheJudgeNeverAnsweredFor(t *testing.T) {
	candidates := judgeBatchCandidates(10)

	outcome := runJudge(context.Background(), partialJudge{answers: 4}, ReviewContext{}, candidates)

	if len(outcome.Judged) != 4 {
		t.Fatalf("Judged = %d, want the 4 candidates the judge actually ruled on", len(outcome.Judged))
	}
	if len(outcome.Unjudged) != 6 {
		t.Fatalf("Unjudged = %d, want the 6 candidates it never mentioned kept rather than deleted", len(outcome.Unjudged))
	}
	if outcome.Unanswered != 6 {
		t.Fatalf("Unanswered = %d, want 6 so the shortfall can be reported", outcome.Unanswered)
	}
	if len(outcome.Judged)+len(outcome.Unjudged) != len(candidates) {
		t.Fatalf("Judged(%d) + Unjudged(%d) does not account for all %d candidates",
			len(outcome.Judged), len(outcome.Unjudged), len(candidates))
	}
	// The batch itself worked. Reporting it as a failure would be the opposite
	// error: it would throw away four verdicts the judge did reach.
	if outcome.BatchesFailed != 0 || outcome.Err != nil {
		t.Fatalf("BatchesFailed/Err = %d/%v, want a successful batch with a shortfall, not a failed one",
			outcome.BatchesFailed, outcome.Err)
	}
}

// TestReviewReportsUnansweredCandidatesAsDegraded closes that loop at the report
// boundary: a finding shipped without a verdict must say it was shipped without
// a verdict.
func TestReviewReportsUnansweredCandidatesAsDegraded(t *testing.T) {
	t.Setenv("TOTALITY_REVIEW_JUDGE", "1")
	root := t.TempDir()
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	engine := judgeTestEngine(root, []Finding{
		judgeTestFinding("bedrock-a.ai.review.1", "internal/app/app.go"),
	})
	engine.judge = partialJudge{answers: 0}

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) != 1 {
		t.Fatalf("Findings = %#v, want the unanswered candidate kept rather than dropped as unconfirmed", report.Findings)
	}
	var reason string
	for _, candidate := range report.DegradedReasons {
		if strings.Contains(candidate, "no verdict") {
			reason = candidate
			break
		}
	}
	if reason == "" {
		t.Fatalf("DegradedReasons = %#v, want the missing verdict reported", report.DegradedReasons)
	}
}

// TestReviewReportsFailedVerificationAsDegraded pins that a failed verification
// is visible in the report. It used to collapse the findings to at most
// maxAdvisoryFindings while the report still read as a completed review.
func TestReviewReportsFailedVerificationAsDegraded(t *testing.T) {
	t.Setenv("TOTALITY_REVIEW_JUDGE", "1")
	root := t.TempDir()
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	engine := judgeTestEngine(root, []Finding{
		judgeTestFinding("openai.ai.review.1", "internal/app/app.go"),
	})
	engine.judge = failingJudge{}

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	var found bool
	for _, reason := range report.DegradedReasons {
		if strings.Contains(reason, "verification") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("DegradedReasons = %#v, want the failed finding verification reported", report.DegradedReasons)
	}
}
