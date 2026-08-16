package codereview

import "strings"

// Fix plan.
//
// A report lists findings in the order a reader should think about them. A
// fix plan lists the same findings in the order an agent should act on them:
// one step per finding, one action per step, blocking lane first because
// those are the steps that flip the verdict. It is a projection of the report
// — nothing here is decided that the findings do not already carry — so the
// JSON reader gets a to-do list without re-deriving lanes and ordering.

// FixStep is one entry in the fix plan.
type FixStep struct {
	Order     int    `json:"order"`
	FindingID string `json:"finding_id"`
	RuleID    string `json:"rule_id,omitempty"`
	// Action is Finding.Recommendation, trimmed: the one thing to do.
	Action string `json:"action"`
	File   string `json:"file,omitempty"`
	Line   int    `json:"line,omitempty"`
	Lane   string `json:"lane"`
	// FlipsVerdict is true iff the step is in the blocking lane: doing it
	// moves the review from failing the gate to passing it.
	FlipsVerdict bool `json:"flips_verdict"`
}

// BuildFixPlan orders blocking findings first, then advisory; within a lane it
// preserves the report's finding order, which is already significance-ranked.
// Findings with an empty Recommendation are skipped: a step with no action is
// not a step. One entry per finding, one action per entry. Order is 1-based
// and contiguous over the returned steps.
func BuildFixPlan(findings []Finding) []FixStep {
	var blocking, advisory []FixStep
	for _, finding := range findings {
		action := strings.TrimSpace(finding.Recommendation)
		if action == "" {
			continue
		}
		lane := LaneOf(finding)
		step := FixStep{
			FindingID:    finding.ID,
			RuleID:       strings.TrimSpace(finding.RuleID),
			Action:       action,
			File:         strings.TrimSpace(finding.File),
			Line:         finding.Line,
			Lane:         lane,
			FlipsVerdict: lane == LaneBlocking,
		}
		if step.FlipsVerdict {
			blocking = append(blocking, step)
		} else {
			advisory = append(advisory, step)
		}
	}
	if len(blocking)+len(advisory) == 0 {
		return nil
	}
	steps := append(blocking, advisory...)
	for i := range steps {
		steps[i].Order = i + 1
	}
	return steps
}
