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

// TestReviewReportsFailedVerificationAsDegraded pins that a failed verification
// is visible in the report. It used to collapse the findings to at most
// maxAdvisoryFindings while the report still read as a completed review.
func TestReviewReportsFailedVerificationAsDegraded(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
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
