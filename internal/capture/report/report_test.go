package report_test

import (
	"bytes"
	"testing"

	"github.com/satoricorp/gx/internal/capture/matcher"
	"github.com/satoricorp/gx/internal/capture/report"
)

func TestVerdictForCoverage(t *testing.T) {
	cases := []struct {
		cov  float64
		want string
	}{
		{0.80, "proceed"},
		{0.70, "proceed"},
		{0.55, "investigate"},
		{0.30, "stop"},
	}
	for _, tc := range cases {
		if got := report.VerdictForCoverage(tc.cov); got != tc.want {
			t.Fatalf("coverage %.2f => %q, want %q", tc.cov, got, tc.want)
		}
	}
}

func TestBuildAggregate(t *testing.T) {
	repos := []report.RepoStats{
		{EligibleHunks: 10, Tier1Hunks: 7, Tier2Hunks: 1, Tier3Pairs: 3, HunkCoverage: 0.8, Verdict: "proceed"},
		{EligibleHunks: 10, Tier1Hunks: 3, Tier2Hunks: 2, Tier3Pairs: 5, HunkCoverage: 0.5, Verdict: "investigate"},
	}
	agg := report.BuildAggregate(repos)
	if agg.EligibleHunks != 20 {
		t.Fatalf("eligible = %d", agg.EligibleHunks)
	}
	if agg.Tier1Hunks != 10 || agg.Tier2Hunks != 3 {
		t.Fatalf("tier counts t1=%d t2=%d", agg.Tier1Hunks, agg.Tier2Hunks)
	}
	if agg.HunkCoverage != 0.65 {
		t.Fatalf("coverage = %f", agg.HunkCoverage)
	}
	if agg.GateVerdict != "investigate" {
		t.Fatalf("verdict = %q", agg.GateVerdict)
	}
}

func TestWriteJSONAndTable(t *testing.T) {
	result := matcher.Result{
		Tier1HunkIndexes: map[int]struct{}{0: {}},
		Tier2HunkIndexes: map[int]struct{}{},
		Tier3Pairs:       2,
	}
	stats := report.BuildRepoStats(report.BuildInput{
		Root:                "/tmp/repo",
		RefRange:            "main..HEAD",
		EligibleHunks:       1,
		EligibleAgentEvents: 1,
		MatchedAgentEvents:  1,
		MatchResult:         result,
	})
	doc := report.NewMatchReport([]report.RepoStats{stats})
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf, doc); err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"gate_verdict"`)) {
		t.Fatal("missing gate_verdict in json")
	}
	report.PrintTable(&buf, doc)
	if !bytes.Contains(buf.Bytes(), []byte("AGGREGATE")) {
		t.Fatal("missing aggregate row in table")
	}
}
