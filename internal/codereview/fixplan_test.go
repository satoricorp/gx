package codereview

import (
	"strings"
	"testing"
)

// The fix plan is the report's findings re-cut for an agent: blocking steps
// first (they flip the verdict), advisory after, and inside each lane the
// report's own order — which is already significance-ranked — is kept.
func TestBuildFixPlanOrdersBlockingFirstThenAdvisoryInReportOrder(t *testing.T) {
	findings := []Finding{
		{ID: "ai.review.1", Recommendation: "advisory one", Lane: LaneAdvisory, File: "a.go", Line: 10},
		{ID: "ai.review.2", RuleID: "review/wrong-logic", Recommendation: "blocking one", Lane: LaneBlocking, File: "b.go", Line: 20},
		{ID: "ai.review.3", Recommendation: "advisory two", Lane: LaneAdvisory, File: "c.go", Line: 30},
		{ID: "ai.review.4", Recommendation: "blocking two", Lane: LaneBlocking, File: "d.go", Line: 40},
	}
	plan := BuildFixPlan(findings)
	if len(plan) != 4 {
		t.Fatalf("expected 4 steps, got %d: %+v", len(plan), plan)
	}
	wantIDs := []string{"ai.review.2", "ai.review.4", "ai.review.1", "ai.review.3"}
	for i, step := range plan {
		if step.FindingID != wantIDs[i] {
			t.Errorf("step %d: finding %q, want %q", i, step.FindingID, wantIDs[i])
		}
		if step.Order != i+1 {
			t.Errorf("step %d: order %d, want %d", i, step.Order, i+1)
		}
	}
	if !plan[0].FlipsVerdict || !plan[1].FlipsVerdict {
		t.Errorf("blocking steps must flip the verdict: %+v", plan[:2])
	}
	if plan[2].FlipsVerdict || plan[3].FlipsVerdict {
		t.Errorf("advisory steps must not flip the verdict: %+v", plan[2:])
	}
	if plan[0].RuleID != "review/wrong-logic" || plan[0].File != "b.go" || plan[0].Line != 20 || plan[0].Lane != LaneBlocking {
		t.Errorf("step 1 lost its identity: %+v", plan[0])
	}
	if plan[0].Action != "blocking one" {
		t.Errorf("step 1 action = %q", plan[0].Action)
	}
}

// A step with nothing to do is not a step: findings with an empty (or
// whitespace) Recommendation are dropped, and the surviving steps renumber so
// Order stays contiguous.
func TestBuildFixPlanSkipsFindingsWithoutRecommendation(t *testing.T) {
	findings := []Finding{
		{ID: "one", Recommendation: "   ", Lane: LaneBlocking},
		{ID: "two", Recommendation: "  do it  ", Lane: LaneAdvisory},
		{ID: "three", Lane: LaneBlocking},
	}
	plan := BuildFixPlan(findings)
	if len(plan) != 1 {
		t.Fatalf("expected 1 step, got %d: %+v", len(plan), plan)
	}
	if plan[0].FindingID != "two" || plan[0].Order != 1 || plan[0].Action != "do it" {
		t.Errorf("unexpected surviving step: %+v", plan[0])
	}
	if got := BuildFixPlan(nil); got != nil {
		t.Errorf("empty input should yield nil, got %+v", got)
	}
}

// A report that never had lanes assigned (an older saved report, a fixture)
// still gets a lane per step, derived by the same rule the gate uses.
func TestBuildFixPlanDerivesLaneWhenUnassigned(t *testing.T) {
	findings := []Finding{
		{ID: "weak", Strength: "Speculative", Recommendation: "maybe"},
		{ID: "hard", Strength: "Blocking", Recommendation: "must"},
	}
	plan := BuildFixPlan(findings)
	if len(plan) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(plan))
	}
	if plan[0].FindingID != "hard" || !plan[0].FlipsVerdict || plan[0].Lane != LaneBlocking {
		t.Errorf("blocking-strength finding should lead and flip: %+v", plan[0])
	}
	if plan[1].FindingID != "weak" || plan[1].FlipsVerdict || !strings.EqualFold(plan[1].Lane, LaneAdvisory) {
		t.Errorf("speculative finding should trail as advisory: %+v", plan[1])
	}
}
