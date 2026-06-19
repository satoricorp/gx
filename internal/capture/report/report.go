package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/satoricorp/gx/internal/capture/matcher"
)

// RepoStats holds per-repository spike metrics.
type RepoStats struct {
	Root                string  `json:"root"`
	RefRange            string  `json:"ref_range"`
	CommitCount         int     `json:"commit_count"`
	EligibleHunks       int     `json:"eligible_hunks"`
	EligibleAgentEvents int     `json:"eligible_agent_events"`
	Tier1Hunks          int     `json:"tier1_hunks"`
	Tier2Hunks          int     `json:"tier2_hunks"`
	Tier3Pairs          int     `json:"tier3_pairs"`
	HunkCoverage        float64 `json:"hunk_coverage"`
	AgentPrecision      float64 `json:"agent_precision"`
	ExcludedFiles       int     `json:"excluded_files"`
	Verdict             string  `json:"verdict"`
}

// AggregateStats summarizes all repos.
type AggregateStats struct {
	EligibleHunks  int     `json:"eligible_hunks"`
	Tier1Hunks     int     `json:"tier1_hunks"`
	Tier2Hunks     int     `json:"tier2_hunks"`
	Tier3Pairs     int     `json:"tier3_pairs"`
	HunkCoverage   float64 `json:"hunk_coverage"`
	GateVerdict    string  `json:"gate_verdict"`
}

// MatchReport is the JSON spike output.
type MatchReport struct {
	SchemaVersion int            `json:"schema_version"`
	GeneratedAt   string         `json:"generated_at"`
	Repos         []RepoStats    `json:"repos"`
	Aggregate     AggregateStats `json:"aggregate"`
}

// BuildInput carries raw counts for one repo.
type BuildInput struct {
	Root                string
	RefRange            string
	CommitCount         int
	EligibleHunks       int
	EligibleAgentEvents int
	MatchedAgentEvents  int
	MatchResult         matcher.Result
	ExcludedFiles       int
}

// BuildRepoStats computes per-repo metrics and verdict.
func BuildRepoStats(in BuildInput) RepoStats {
	t1 := len(in.MatchResult.Tier1HunkIndexes)
	t2 := len(in.MatchResult.Tier2HunkIndexes)
	coverage := 0.0
	if in.EligibleHunks > 0 {
		coverage = float64(t1+t2) / float64(in.EligibleHunks)
	}
	precision := 0.0
	if in.EligibleAgentEvents > 0 {
		precision = float64(in.MatchedAgentEvents) / float64(in.EligibleAgentEvents)
	}
	return RepoStats{
		Root:                in.Root,
		RefRange:            in.RefRange,
		CommitCount:         in.CommitCount,
		EligibleHunks:       in.EligibleHunks,
		EligibleAgentEvents: in.EligibleAgentEvents,
		Tier1Hunks:          t1,
		Tier2Hunks:          t2,
		Tier3Pairs:          in.MatchResult.Tier3Pairs,
		HunkCoverage:        coverage,
		AgentPrecision:      precision,
		ExcludedFiles:       in.ExcludedFiles,
		Verdict:             VerdictForCoverage(coverage),
	}
}

// VerdictForCoverage maps hunk coverage to gate verdict.
func VerdictForCoverage(coverage float64) string {
	switch {
	case coverage >= 0.70:
		return "proceed"
	case coverage >= 0.50:
		return "investigate"
	default:
		return "stop"
	}
}

// BuildAggregate combines repo stats into aggregate gate metrics.
func BuildAggregate(repos []RepoStats) AggregateStats {
	var eligible, t1, t2, t3 int
	for _, r := range repos {
		eligible += r.EligibleHunks
		t1 += r.Tier1Hunks
		t2 += r.Tier2Hunks
		t3 += r.Tier3Pairs
	}
	coverage := 0.0
	if eligible > 0 {
		coverage = float64(t1+t2) / float64(eligible)
	}
	return AggregateStats{
		EligibleHunks: eligible,
		Tier1Hunks:    t1,
		Tier2Hunks:    t2,
		Tier3Pairs:    t3,
		HunkCoverage:  coverage,
		GateVerdict:   VerdictForCoverage(coverage),
	}
}

// NewMatchReport builds the full report document.
func NewMatchReport(repos []RepoStats) MatchReport {
	return MatchReport{
		SchemaVersion: 1,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Repos:         repos,
		Aggregate:     BuildAggregate(repos),
	}
}

// WriteJSON encodes the report to w.
func WriteJSON(w io.Writer, report MatchReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

// PrintTable renders a human-readable summary table.
func PrintTable(w io.Writer, report MatchReport) {
	fmt.Fprintln(w, "Repository          Hunks  Agent  T1   T2   T3   Coverage  Verdict")
	for _, r := range report.Repos {
		root := r.Root
		if len(root) > 18 {
			root = "..." + root[len(root)-15:]
		}
		fmt.Fprintf(w, "%-20s %5d %5d %3d %3d %3d %8.1f%%  %s\n",
			root, r.EligibleHunks, r.EligibleAgentEvents,
			r.Tier1Hunks, r.Tier2Hunks, r.Tier3Pairs,
			r.HunkCoverage*100, r.Verdict)
	}
	fmt.Fprintln(w, strings.Repeat("─", 65))
	fmt.Fprintf(w, "AGGREGATE             %5d       %3d %3d %3d %8.1f%%  %s\n",
		report.Aggregate.EligibleHunks,
		report.Aggregate.Tier1Hunks,
		report.Aggregate.Tier2Hunks,
		report.Aggregate.Tier3Pairs,
		report.Aggregate.HunkCoverage*100,
		report.Aggregate.GateVerdict)
}
