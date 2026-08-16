package codereview

import "testing"

func TestAssignLanesDerivesLaneFromStrengthAndQuorum(t *testing.T) {
	twoLegs := []string{"Bedrock A", "Bedrock B"}
	oneLeg := []string{"Bedrock A"}
	cases := []struct {
		name        string
		finding     Finding
		wantLane    string
		wantDemoted string
	}{
		{
			name:     "strong, two legs, confirmed is blocking",
			finding:  Finding{ID: "a", Strength: "Strong", Corroboration: twoLegs, JudgeVerdict: "confirmed"},
			wantLane: LaneBlocking,
		},
		{
			name:        "strong, one leg, confirmed is demoted to advisory",
			finding:     Finding{ID: "b", Strength: "Strong", Corroboration: oneLeg, JudgeVerdict: "confirmed"},
			wantLane:    LaneAdvisory,
			wantDemoted: LaneBlocking,
		},
		{
			name:        "strong, two legs, unverified is demoted to advisory",
			finding:     Finding{ID: "c", Strength: "Strong", Corroboration: twoLegs, JudgeVerdict: "unverified"},
			wantLane:    LaneAdvisory,
			wantDemoted: LaneBlocking,
		},
		{
			name:     "worth exploring, two legs, confirmed is advisory without demotion",
			finding:  Finding{ID: "d", Strength: "Worth exploring", Corroboration: twoLegs, JudgeVerdict: "confirmed"},
			wantLane: LaneAdvisory,
		},
		{
			name:     "deterministic blocking-strength finding with no judge is blocking",
			finding:  Finding{ID: "tools.static-failure", Strength: "Blocking"},
			wantLane: LaneBlocking,
		},
		{
			name:     "deterministic strong finding with no judge is advisory, not demoted",
			finding:  Finding{ID: "docs.missing-readme", Strength: "Strong"},
			wantLane: LaneAdvisory,
		},
		{
			name:     "deterministic strong finding the judge confirmed has no legs to count and is not demoted",
			finding:  Finding{ID: "docs.missing-readme", Strength: "Strong", JudgeVerdict: "confirmed"},
			wantLane: LaneAdvisory,
		},
		{
			name:     "fast review: strong, one leg, no judge is advisory, not demoted",
			finding:  Finding{ID: "bedrock-a.ai.review.1", Strength: "Strong", Corroboration: oneLeg},
			wantLane: LaneAdvisory,
		},
		{
			name:     "speculative is advisory",
			finding:  Finding{ID: "e", Strength: "Speculative", Corroboration: twoLegs, JudgeVerdict: "confirmed"},
			wantLane: LaneAdvisory,
		},
		{
			name:     "empty strength is advisory",
			finding:  Finding{ID: "f"},
			wantLane: LaneAdvisory,
		},
		{
			name:     "verdict comparison is case-insensitive",
			finding:  Finding{ID: "g", Strength: "Strong", Corroboration: twoLegs, JudgeVerdict: "Confirmed"},
			wantLane: LaneBlocking,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := AssignLanes([]Finding{tc.finding})
			if len(got) != 1 {
				t.Fatalf("AssignLanes() returned %d findings, want 1", len(got))
			}
			if got[0].Lane != tc.wantLane {
				t.Fatalf("Lane = %q, want %q", got[0].Lane, tc.wantLane)
			}
			if got[0].DemotedFrom != tc.wantDemoted {
				t.Fatalf("DemotedFrom = %q, want %q", got[0].DemotedFrom, tc.wantDemoted)
			}
			// The lane must survive a second pass unchanged.
			again := AssignLanes(got)
			if again[0].Lane != tc.wantLane || again[0].DemotedFrom != tc.wantDemoted {
				t.Fatalf("second AssignLanes() = %q/%q, want %q/%q", again[0].Lane, again[0].DemotedFrom, tc.wantLane, tc.wantDemoted)
			}
			// LaneOf agrees with the recorded lane and with the computed one.
			if LaneOf(got[0]) != tc.wantLane {
				t.Fatalf("LaneOf(assigned) = %q, want %q", LaneOf(got[0]), tc.wantLane)
			}
			if LaneOf(tc.finding) != tc.wantLane {
				t.Fatalf("LaneOf(unassigned) = %q, want %q", LaneOf(tc.finding), tc.wantLane)
			}
		})
	}
}

// LaneOf trusts a recorded lane over recomputation, so a report whose lanes
// were assigned by the pipeline that produced it is gated on what it says.
func TestLaneOfPrefersTheRecordedLane(t *testing.T) {
	f := Finding{ID: "a", Strength: "Speculative", Lane: LaneBlocking}
	if got := LaneOf(f); got != LaneBlocking {
		t.Fatalf("LaneOf() = %q, want the recorded %q", got, LaneBlocking)
	}
}

func TestAssignLanesHandlesEmptyInput(t *testing.T) {
	if got := AssignLanes(nil); got != nil {
		t.Fatalf("AssignLanes(nil) = %#v, want nil", got)
	}
}
