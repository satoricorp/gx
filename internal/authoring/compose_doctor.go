package authoring

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	cursoringest "github.com/satoricorp/gx/internal/ingest/cursor"
	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/parsers"
)

// ComposeDoctorReport summarizes compose health checks.
type ComposeDoctorReport struct {
	JJOK              bool     `json:"jj_ok"`
	JJVersion         string   `json:"jj_version,omitempty"`
	JJRepoOK          bool     `json:"jj_repo_ok"`
	ParsersReachable  []string `json:"parsers_reachable,omitempty"`
	ParsersMissing    []string `json:"parsers_missing,omitempty"`
	LastHunkCoverage  float64  `json:"last_hunk_coverage,omitempty"`
	LastCoverageFound bool     `json:"last_coverage_found"`
	FeasibilityHint   string   `json:"feasibility_hint,omitempty"`
	Warnings          []string `json:"warnings,omitempty"`
	OK                bool     `json:"ok"`
}

// RunComposeDoctor checks jj provisioning, capture parsers, and last matcher coverage.
func (e *Engine) RunComposeDoctor(ctx context.Context) (ComposeDoctorReport, error) {
	report := ComposeDoctorReport{}
	version, err := e.vcs.JJVersion(ctx)
	if err != nil {
		report.Warnings = append(report.Warnings, "jj not available: "+err.Error())
	} else {
		report.JJOK = true
		report.JJVersion = version
	}
	repo, repoErr := e.vcs.ResolveJJRepo(ctx)
	if repoErr != nil {
		report.Warnings = append(report.Warnings, "not a jj repo: "+repoErr.Error())
	} else {
		report.JJRepoOK = true
		if coverage, ok := readLastComposeHunkCoverage(repo.RootPath); ok {
			report.LastHunkCoverage = coverage
			report.LastCoverageFound = true
		}
	}
	report.ParsersReachable, report.ParsersMissing = captureParserChecks(repo.RootPath)
	if len(report.ParsersMissing) > 0 {
		report.Warnings = append(report.Warnings, "capture parsers missing: "+strings.Join(report.ParsersMissing, ", "))
	}
	if pending, err := e.LoadLatestPendingDemuxProposal(ctx); err == nil && len(pending.FeasibilityWarnings) > 0 {
		report.FeasibilityHint = fmt.Sprintf("latest compose proposal has %d feasibility warning(s)", len(pending.FeasibilityWarnings))
	}
	report.OK = report.JJOK && report.JJRepoOK && len(report.ParsersMissing) == 0
	return report, nil
}

func captureParserChecks(repoRoot string) (reachable, missing []string) {
	tools := []string{capture.ToolClaude, capture.ToolCodex, capture.ToolCursor}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, tools
	}
	now := time.Now()
	discovered, err := parsers.DiscoverSessions(parsers.DiscoverOptions{
		HomeDir:  home,
		RepoRoot: repoRoot,
		Since:    now.Add(-30 * 24 * time.Hour),
		Until:    now.Add(time.Hour),
		Tools:    tools,
	})
	if err != nil {
		return nil, tools
	}
	found := map[string]struct{}{}
	for _, session := range discovered {
		found[session.Tool] = struct{}{}
	}
	if path, err := cursoringest.DefaultVSCDBPath(); err == nil {
		if _, statErr := os.Stat(path); statErr == nil {
			found[capture.ToolCursor] = struct{}{}
		}
	}
	for _, tool := range tools {
		if _, ok := found[tool]; ok {
			reachable = append(reachable, tool)
		} else {
			missing = append(missing, tool)
		}
	}
	return reachable, missing
}
