package codereview

import "testing"

func TestParseFailOnLevelAcceptsStrengthSpellings(t *testing.T) {
	for input, want := range map[string]FailOnLevel{
		"":                FailOnNone,
		"none":            FailOnNone,
		"any":             FailOnAny,
		"Blocking":        FailOnBlocking,
		"strong":          FailOnStrong,
		"worth exploring": FailOnWorthExploring,
		"worth-exploring": FailOnWorthExploring,
		"speculative":     FailOnSpeculative,
	} {
		got, err := ParseFailOnLevel(input)
		if err != nil {
			t.Fatalf("ParseFailOnLevel(%q) error = %v", input, err)
		}
		if got != want {
			t.Fatalf("ParseFailOnLevel(%q) = %q, want %q", input, got, want)
		}
	}
	if _, err := ParseFailOnLevel("critical"); err == nil {
		t.Fatal("ParseFailOnLevel(critical) accepted an unknown level")
	}
}

func TestGateFailuresAppliesTheThreshold(t *testing.T) {
	report := Report{
		ReviewMode: ReviewModeRange,
		Reviewed:   true,
		Findings: []Finding{
			{ID: "a", Strength: "Blocking"},
			{ID: "b", Strength: "Strong"},
			{ID: "c", Strength: "Worth exploring"},
			{ID: "d", Strength: "Speculative"},
		},
	}
	for level, want := range map[FailOnLevel]int{
		FailOnNone:           0,
		FailOnBlocking:       1,
		FailOnStrong:         2,
		FailOnWorthExploring: 3,
		FailOnSpeculative:    4,
		FailOnAny:            4,
	} {
		if got := len(report.GateFailures(level)); got != want {
			t.Fatalf("GateFailures(%q) = %d findings, want %d", level, got, want)
		}
	}
}

// The blocking level is lane-aware: quorum decides, not strength alone. Every
// other level is unchanged and still reads strength.
func TestGateFailuresBlockingLevelReadsLanes(t *testing.T) {
	twoLegs := []string{"Bedrock A", "Bedrock B"}
	report := Report{
		ReviewMode: ReviewModeRange,
		Reviewed:   true,
		Findings: AssignLanes([]Finding{
			// Deterministic Blocking: no judge, strength alone → blocking lane.
			{ID: "tools.static-failure", Strength: "Blocking"},
			// Strong with quorum: two legs and confirmed → blocking lane.
			{ID: "quorum", Strength: "Strong", Corroboration: twoLegs, JudgeVerdict: "confirmed"},
			// Strong, one leg, confirmed → demoted to advisory.
			{ID: "one-leg", Strength: "Strong", Corroboration: []string{"Bedrock A"}, JudgeVerdict: "confirmed"},
			// Strong, two legs, unverified → demoted to advisory.
			{ID: "unverified", Strength: "Strong", Corroboration: twoLegs, JudgeVerdict: "unverified"},
			// Below Strong never blocks.
			{ID: "explore", Strength: "Worth exploring", Corroboration: twoLegs, JudgeVerdict: "confirmed"},
		}),
	}
	blocking := report.GateFailures(FailOnBlocking)
	if len(blocking) != 2 {
		t.Fatalf("GateFailures(blocking) = %d findings (%s), want 2 (the Blocking tool finding and the quorum finding)", len(blocking), findingIDs(blocking))
	}
	for _, f := range blocking {
		if f.ID != "tools.static-failure" && f.ID != "quorum" {
			t.Fatalf("GateFailures(blocking) included %q, want only lane-blocking findings", f.ID)
		}
	}
	// "strong" keeps strength-only semantics: every Strong-or-worse finding,
	// demoted or not.
	if got := len(report.GateFailures(FailOnStrong)); got != 4 {
		t.Fatalf("GateFailures(strong) = %d findings, want 4", got)
	}
	if got := len(report.GateFailures(FailOnAny)); got != 5 {
		t.Fatalf("GateFailures(any) = %d findings, want 5", got)
	}
}

// A report whose lanes were never assigned is gated by the same rule: LaneOf
// computes the lane on the fly. A demoted-in-effect finding must not fail the
// blocking gate just because Lane is empty.
func TestGateFailuresBlockingLevelComputesMissingLanes(t *testing.T) {
	report := Report{
		ReviewMode: ReviewModeRange,
		Reviewed:   true,
		Findings: []Finding{
			{ID: "one-leg", Strength: "Strong", Corroboration: []string{"Bedrock A"}, JudgeVerdict: "confirmed"},
			{ID: "tools.static-failure", Strength: "Blocking"},
		},
	}
	got := report.GateFailures(FailOnBlocking)
	if len(got) != 1 || got[0].ID != "tools.static-failure" {
		t.Fatalf("GateFailures(blocking) = %s, want only tools.static-failure", findingIDs(got))
	}
}

// A gate must never report "no issues at or above X" for a diff it never read.
func TestGateFailuresIgnoresAnUnreviewedReport(t *testing.T) {
	report := Report{
		ReviewMode: ReviewModeNone,
		Findings:   []Finding{{ID: "a", Strength: "Blocking"}},
	}
	if got := report.GateFailures(FailOnAny); got != nil {
		t.Fatalf("GateFailures() = %#v, want nil when nothing was reviewed", got)
	}
}

func TestFailOnLevelEnabled(t *testing.T) {
	if FailOnNone.Enabled() {
		t.Fatal("FailOnNone.Enabled() = true, want false")
	}
	for _, level := range []FailOnLevel{FailOnAny, FailOnBlocking, FailOnStrong, FailOnWorthExploring, FailOnSpeculative} {
		if !level.Enabled() {
			t.Fatalf("%q.Enabled() = false, want true", level)
		}
	}
}
