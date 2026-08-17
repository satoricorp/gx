package codereview

import "strings"

// Lanes.
//
// Every finding lands in one of two lanes. Blocking stops the change; advisory
// is reported and does not. The lane is derived from two things the pipeline
// already records — the finding's strength and its quorum — and is never asked
// of a model directly, because a model asked "is this blocking?" says yes far
// more often than two independent models converge on the same problem.
//
// Strength decides the default lane. Quorum decides whether a finding of
// blocking strength is allowed to keep it. Quorum is two facts together: at
// least two independent reviewer legs raised the finding on their own
// (Finding.Corroboration), and the verification model confirmed it
// (Finding.JudgeVerdict). One flagship model's opinion, however strongly
// worded, is advisory until a second one agrees and the judge signs off.
//
// Materiality (how much observable behavior moved) is a separate axis on the
// Finding and plays no part here on purpose.
const (
	LaneBlocking = "blocking"
	LaneAdvisory = "advisory"
)

// laneQuorumLegs is how many independent reviewer legs must have raised a
// finding before its strength alone can make it blocking.
const laneQuorumLegs = 2

// AssignLanes sets Lane and DemotedFrom on every finding and returns the
// slice. It is idempotent: the lane is a pure function of strength, quorum,
// and verdict, so calling it twice computes the same answer.
//
// The rule:
//
//   - Strength ranks at or above Strong (strengthRank <= 1) AND the finding
//     went through the judge (JudgeVerdict != "") AND at least one reviewer
//     leg raised it (Corroboration non-empty): blocking iff quorum holds —
//     len(Corroboration) >= laneQuorumLegs and JudgeVerdict == "confirmed".
//     Strength that qualifies without quorum is advisory with
//     DemotedFrom = LaneBlocking, so the report can say "one grader flagged ·
//     demoted" rather than quietly filing a Strong finding under advisory.
//
//   - No judge verdict, or no reviewer leg: strength alone decides, no
//     demotion, because there is no quorum to fail. Only the literal
//     "Blocking" strength lands in the blocking lane here; "Strong" is
//     advisory. That is what keeps --fast (one leg, no judge) from blocking
//     on a single model's say-so where the full panel would demand two legs
//     plus confirmation, and it keeps deterministic rules meaning what their
//     strength says: tools.static-failure is Blocking and blocks, a missing
//     README is Strong and does not.
//
//   - Anything below Strong is advisory, whatever the quorum.
//
// Recognizing the "no quorum to fail" path. Deterministic findings — the
// rules in defaultRules() and every gates.* finding — never carry
// Corroboration: only the AI panel sets it (ai.go stamps the raising leg on
// each finding, and dedupe merges the legs into one list). They also mostly
// never carry a JudgeVerdict: tools.* Blocking findings are split out before
// the judge runs (splitBlockingToolFindings), and gates findings never
// see it at all. The remaining deterministic rules can reach the judge on a
// whole-repo review, which is why the test is "no verdict OR no leg" rather
// than the verdict alone: a deterministic finding the judge happened to
// confirm has no legs to count, so it takes the strength-only path instead of
// being reported as demoted from a lane it could never have qualified for.
// This is simpler and more reliable than an ID-prefix allowlist ("tools.",
// "gates.", ...) that every new deterministic rule would have to join.
func AssignLanes(findings []Finding) []Finding {
	for i := range findings {
		lane, demotedFrom := laneFor(findings[i])
		findings[i].Lane = lane
		findings[i].DemotedFrom = demotedFrom
	}
	return findings
}

// LaneOf is the effective lane of a finding: the recorded Lane when
// AssignLanes has run, otherwise the lane it would assign. Gates read this so
// a report that reached them without lanes (an older saved report, a test
// fixture) is still judged by the same rule.
func LaneOf(f Finding) string {
	if lane := strings.TrimSpace(f.Lane); lane != "" {
		return lane
	}
	lane, _ := laneFor(f)
	return lane
}

// laneFor computes (lane, demotedFrom) for one finding. See AssignLanes.
func laneFor(f Finding) (string, string) {
	if strengthRank(f.Strength) > strengthRank("Strong") {
		return LaneAdvisory, ""
	}
	verdict := strings.ToLower(strings.TrimSpace(f.JudgeVerdict))
	if verdict == "" || len(f.Corroboration) == 0 {
		// No quorum to fail: strength alone decides.
		if strengthRank(f.Strength) == strengthRank("Blocking") {
			return LaneBlocking, ""
		}
		return LaneAdvisory, ""
	}
	if len(f.Corroboration) >= laneQuorumLegs && verdict == "confirmed" {
		return LaneBlocking, ""
	}
	return LaneAdvisory, LaneBlocking
}
