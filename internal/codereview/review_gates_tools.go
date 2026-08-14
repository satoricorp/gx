package codereview

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

// The gate partitions one shared static-tool run: compile-and-test runners
// answer the correctness gate, linters answer code health. The names here must
// match staticToolRunners entries (plus the review-only extras below).
var (
	reviewCorrectnessTools = map[string]bool{
		"go test":         true,
		"tsc":             true,
		"cargo check":     true,
		"dotnet build":    true,
		"flutter analyze": true,
		"pytest":          true,
	}
	reviewLintTools = map[string]bool{
		"go vet":       true,
		"eslint":       true,
		"ruff":         true,
		"mypy":         true,
		"dart analyze": true,
	}
)

// reviewExtraToolRunners are runners the review gate adds on top of
// the review registry. They live here rather than in staticToolRunners so that
// adding the exit gate does not silently change what `gx enhance` executes.
var reviewExtraToolRunners = []staticToolRunner{
	{name: "pytest", progress: "Running pytest", detect: detectPytest, wholeProject: true},
}

// detectPytest runs the project's Python test suite when pytest is configured
// and installed. Configuration is required for the same reason mypy requires
// it: pytest pointed at an unconfigured checkout collects whatever it finds
// and reports import noise about the checkout rather than the change.
func detectPytest(env staticToolEnv) (staticToolCommand, bool) {
	if !pytestConfigured(env.repoRoot) {
		return staticToolCommand{}, false
	}
	if !changedFilesInclude(env.changedFiles, pythonFileExtensions...) {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "pytest", pythonLocalBinDirs...)
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: []string{"pytest", "-q"}}, true
}

func pytestConfigured(repoRoot string) bool {
	return repoFileExists(repoRoot, "pytest.ini") ||
		repoFileContains(repoRoot, "pyproject.toml", "[tool.pytest") ||
		repoFileContains(repoRoot, "setup.cfg", "[tool:pytest]")
}

// collectReviewToolResults is the shared run: the review registry plus
// the review-only extras, never in fast mode — the whole point of the
// gate is running the checks a fast review skips.
func collectReviewToolResults(ctx context.Context, repoRoot string, facts RepoFacts, opts Options, changed []string) []StaticToolResult {
	results := collectStaticToolResults(ctx, repoRoot, facts, opts, changed)
	if reviewStaticToolsDisabled() {
		return results
	}
	env := staticToolEnv{
		repoRoot:     repoRoot,
		facts:        facts,
		changedFiles: staticToolScope(facts, opts, changed),
	}
	if len(env.changedFiles) == 0 {
		return results
	}
	for _, runner := range reviewExtraToolRunners {
		command, ok := runner.detect(env)
		if !ok {
			continue
		}
		reviewProgress(opts, runner.progress)
		results = append(results, runStaticTool(ctx, repoRoot, 90*time.Second, runner.name, command))
	}
	return results
}

func partitionReviewTools(results []StaticToolResult, names map[string]bool) []StaticToolResult {
	var out []StaticToolResult
	for _, result := range results {
		if names[result.Name] {
			out = append(out, result)
		}
	}
	return out
}

// reviewToolLine summarizes one command for the gate's evidence line.
func reviewToolLine(result StaticToolResult) string {
	switch {
	case result.Skipped:
		reason := strings.TrimSpace(result.Reason)
		if reason == "" {
			reason = "skipped"
		}
		return result.Name + ": skipped (" + reason + ")"
	case result.ExitCode != 0:
		return result.Name + ": failed"
	default:
		return result.Name + ": ok"
	}
}

func reviewToolSummary(results []StaticToolResult) string {
	lines := make([]string, 0, len(results))
	for _, result := range results {
		lines = append(lines, reviewToolLine(result))
	}
	return strings.Join(lines, "; ")
}

// reviewToolFinding turns a failed command into a finding. The output
// excerpt is the evidence a reader needs to act without rerunning anything.
func reviewToolFinding(result StaticToolResult, scope string) Finding {
	excerpt := reviewOutputExcerpt(result.Output, 5)
	finding := Finding{
		ID:             "review." + strings.ReplaceAll(result.Name, " ", "-"),
		Scopes:         []string{scope},
		Title:          result.Name + " failed",
		Summary:        fmt.Sprintf("`%s` exited %d on this change.", result.Command, result.ExitCode),
		Recommendation: fmt.Sprintf("Run `%s` locally and fix the failures before shipping.", result.Command),
		Strength:       "Blocking",
		Kind:           "defect",
	}
	if excerpt != "" {
		finding.Evidence = append(finding.Evidence, Evidence{Label: "output", Value: excerpt})
	}
	return finding
}

func reviewOutputExcerpt(output string, maxLines int) string {
	var lines []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
		if len(lines) >= maxLines {
			break
		}
	}
	return strings.Join(lines, "\n")
}

// reviewCorrectnessGate: the project's own tests and builds, taken as
// they answered. FAIL on any real failure; SKIPPED when nothing could run,
// with the reason — a gate must never read "passed" for a suite that never
// started.
func reviewCorrectnessGate(results []StaticToolResult, changed []string) GateResult {
	gate := GateResult{Gate: GateCorrectness, Title: gateTitle(GateCorrectness), Files: reviewCodeFiles(changed)}
	if reviewStaticToolsDisabled() {
		gate.Status = GateSkipped
		gate.SkipReason = "project checks disabled (GX_REVIEW_STATIC_TOOLS=0)"
		return gate
	}
	checks := partitionReviewTools(results, reviewCorrectnessTools)
	gate.Checks = checks
	if len(checks) == 0 {
		gate.Status = GateSkipped
		gate.SkipReason = "no supported test or build runner detected for this change"
		return gate
	}
	gate.Status = GatePass
	ran := 0
	for _, check := range checks {
		if check.Skipped {
			continue
		}
		ran++
		if check.ExitCode != 0 {
			gate.Status = GateFail
			gate.Findings = append(gate.Findings, reviewToolFinding(check, "testing"))
		}
	}
	if ran == 0 {
		// Every detected runner was skipped (unready checkout, timeout): that
		// is an unanswered question, not a pass.
		gate.Status = GateSkipped
		gate.SkipReason = "all detected runners were skipped: " + reviewToolSummary(checks)
		return gate
	}
	gate.Summary = reviewToolSummary(checks)
	return gate
}

// reviewCodeHealthGate is the deterministic core of code health: linter
// verdicts and the reviewability thresholds. The AI overlay may add judgment
// findings on top; it never overturns what is decided here.
func reviewCodeHealthGate(results []StaticToolResult, signals reviewSignals, changed []string) GateResult {
	gate := GateResult{Gate: GateCodeHealth, Title: gateTitle(GateCodeHealth), Status: GatePass}
	gate.Files = signals.SourceFilesChanged
	var summary []string

	if !reviewStaticToolsDisabled() {
		checks := partitionReviewTools(results, reviewLintTools)
		gate.Checks = checks
		for _, check := range checks {
			if check.Skipped {
				continue
			}
			if check.ExitCode != 0 {
				gate.Status = GateFail
				gate.Findings = append(gate.Findings, reviewToolFinding(check, "maintainability"))
			}
		}
		if len(checks) > 0 {
			summary = append(summary, reviewToolSummary(checks))
		}
	} else {
		gate.Evidence = append(gate.Evidence, Evidence{Label: "linters", Value: "disabled (GX_REVIEW_STATIC_TOOLS=0)"})
	}

	// Reviewability thresholds warn, they do not block: a large change still
	// gets checked and shipped on its merits — Joe's call — but the reader is
	// told the review read less carefully than a small change would get.
	if signals.NonGenLines > reviewMaxReviewableLines || signals.NonGenerated > reviewMaxReviewableFiles {
		gate.Findings = append(gate.Findings, Finding{
			ID:             "review.change-size",
			Scopes:         []string{"maintainability"},
			Title:          "Large change — review carefully",
			Summary:        fmt.Sprintf("%d changed lines across %d non-generated files is more than a human reviews carefully in one sitting; the other gates still ran, but weigh their PASSes accordingly.", signals.NonGenLines, signals.NonGenerated),
			Recommendation: "Consider splitting future changes of this size — one concern per change — so each gets a full-attention review.",
			Strength:       "Worth exploring",
			Kind:           "suggestion",
		})
		summary = append(summary, fmt.Sprintf("large change (%d lines)", signals.NonGenLines))
	}
	summary = append(summary, fmt.Sprintf("%d source file(s) changed, %d test file(s)", len(signals.SourceFilesChanged), len(signals.TestFilesChanged)))
	gate.Evidence = append(gate.Evidence, Evidence{
		Label: "tests accompany sources",
		Value: fmt.Sprintf("%d source file(s) changed alongside %d test file(s)", len(signals.SourceFilesChanged), len(signals.TestFilesChanged)),
	})
	gate.Summary = strings.Join(summary, "; ")
	return gate
}

// reviewAccessibilityGate applies only when UI files changed; otherwise
// it skips and says why. Its deterministic half is the added-line checks; the
// AI overlay confirms the softer candidates.
func reviewAccessibilityGate(signals reviewSignals) GateResult {
	gate := GateResult{Gate: GateAccessibility, Title: gateTitle(GateAccessibility)}
	if len(signals.UIFiles) == 0 {
		gate.Status = GateSkipped
		gate.SkipReason = "no UI files changed"
		return gate
	}
	gate.Files = signals.UIFiles
	gate.Findings = signals.A11yFindings
	gate.Status = GatePass
	if signals.A11yHardFail {
		gate.Status = GateFail
	}
	switch len(signals.A11yFindings) {
	case 0:
		gate.Summary = fmt.Sprintf("%d UI file(s) changed; static checks clean", len(signals.UIFiles))
	default:
		gate.Summary = fmt.Sprintf("%d UI file(s) changed; %d static accessibility issue(s)", len(signals.UIFiles), len(signals.A11yFindings))
	}
	return gate
}

// reviewPerformanceGate decides applicability deterministically; the
// verdict itself is the model's, which is why an applicable gate stays PASS
// only until the AI pass either justifies it or is found unavailable.
func reviewPerformanceGate(signals reviewSignals) GateResult {
	gate := GateResult{Gate: GatePerformance, Title: gateTitle(GatePerformance)}
	if !signals.PerfApplies {
		gate.Status = GateSkipped
		gate.SkipReason = "no perf-relevant paths changed"
		return gate
	}
	gate.Status = GatePass
	gate.Files = signals.PerfFiles
	sort.Strings(gate.Files)
	for _, trigger := range signals.PerfTriggers {
		gate.Evidence = append(gate.Evidence, Evidence{Label: "applies", Value: trigger})
	}
	gate.Summary = "applies (" + strings.Join(signals.PerfTriggers, "; ") + ")"
	return gate
}
