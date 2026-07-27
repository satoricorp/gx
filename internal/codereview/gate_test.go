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
