package codereview

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// gx review is the pre-ship exit gate: six fixed gate checks over
// the current change, each answering with PASS, FAIL, or SKIPPED, rolled into
// one ship / no-ship verdict. It is deliberately not a review. A review looks
// for whatever matters; the gate asks the same six questions every time, runs
// a real command wherever one exists, and reserves the model for the judgments
// no command can make — with the deterministic results handed to it as facts
// it cannot argue with.
//
// It layers on the review substrate the same way gate.go does: change
// resolution, diff collection, the static-tool registry, the read-only git
// guarantee, and the Bedrock plumbing are all reused rather than rebuilt.

// GateID names one gate. The set is fixed: the gate asks the same
// questions of every change, which is what makes its verdict comparable from
// run to run.
type GateID string

const (
	GateCorrectness   GateID = "correctness"
	GateSecurity      GateID = "security"
	GateCodeHealth    GateID = "code-health"
	GateBackPressure  GateID = "back-pressure"
	GateAccessibility GateID = "accessibility"
	GatePerformance   GateID = "performance"
)

// AllGateIDs is the fixed evaluation and rendering order.
func AllGateIDs() []GateID {
	return []GateID{
		GateCorrectness,
		GateSecurity,
		GateCodeHealth,
		GateBackPressure,
		GateAccessibility,
		GatePerformance,
	}
}

func gateTitle(id GateID) string {
	switch id {
	case GateCorrectness:
		return "Correctness"
	case GateSecurity:
		return "Security"
	case GateCodeHealth:
		return "Code health"
	case GateBackPressure:
		return "Back-pressure"
	case GateAccessibility:
		return "Accessibility"
	case GatePerformance:
		return "Performance"
	default:
		return string(id)
	}
}

// ParseGateIDs parses a comma-separated gate list, rejecting unknown names so
// a typo in --skip-gates cannot silently skip nothing.
func ParseGateIDs(value string) ([]GateID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	known := map[GateID]struct{}{}
	for _, id := range AllGateIDs() {
		known[id] = struct{}{}
	}
	var out []GateID
	for _, part := range strings.Split(value, ",") {
		id := GateID(strings.ToLower(strings.TrimSpace(part)))
		if id == "" {
			continue
		}
		if _, ok := known[id]; !ok {
			return nil, fmt.Errorf("unknown gate %q (want one of: %s)", part, joinGateIDs(AllGateIDs()))
		}
		out = append(out, id)
	}
	return out, nil
}

func joinGateIDs(ids []GateID) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, string(id))
	}
	return strings.Join(parts, ", ")
}

type GateStatus string

const (
	GatePass    GateStatus = "PASS"
	GateFail    GateStatus = "FAIL"
	GateSkipped GateStatus = "SKIPPED"
)

// GateResult is one gate's answer. Summary is the one-line evidence the
// renderer puts next to the status; Findings carry the concrete problems with
// recommendations when the gate failed (or found something worth saying while
// passing). Checks are the commands that actually ran, kept whole so the
// report can show exactly what the verdict rests on. Files are the source
// files the gate examined or that triggered it — always in the JSON report,
// rendered only under --verbose.
type GateResult struct {
	Gate       GateID             `json:"gate"`
	Title      string             `json:"title"`
	Status     GateStatus         `json:"status"`
	SkipReason string             `json:"skip_reason,omitempty"`
	Summary    string             `json:"summary,omitempty"`
	Evidence   []Evidence         `json:"evidence,omitempty"`
	Findings   []Finding          `json:"findings,omitempty"`
	Checks     []StaticToolResult `json:"checks,omitempty"`
	Files      []string           `json:"files,omitempty"`
}

// Reviews verdicts. "degraded" is its own outcome rather than a flavor of
// ship: every resolved gate passed, but gates that needed the model never got
// one, so the run is not entitled to say ship — the same honesty rule the
// review gate applies.
const (
	VerdictShip           = "ship"
	VerdictNoShip         = "no-ship"
	VerdictDegraded       = "degraded"
	VerdictNothingToCheck = "nothing-to-check"
)

type ReviewDiffStats struct {
	Files          int `json:"files"`
	AddedLines     int `json:"added_lines"`
	RemovedLines   int `json:"removed_lines"`
	GeneratedFiles int `json:"generated_files,omitempty"`
}

type ReviewReport struct {
	RepoRoot     string          `json:"repo_root"`
	Intent       string          `json:"intent,omitempty"`
	IntentSource string          `json:"intent_source,omitempty"`
	ChangedFiles []string        `json:"changed_files,omitempty"`
	DiffStats    ReviewDiffStats `json:"diff_stats"`
	Gates        []GateResult    `json:"gates"`
	Verdict      string          `json:"verdict"`
	// DegradedReasons say why the verdict is worth less than it looks: an AI
	// gate that could not run, a truncated diff. They never appear on a
	// deliberate opt-out (GX_REVIEW_AI=0), only on a configured capability that
	// failed — mirroring how the review engine reports degradation.
	DegradedReasons []string `json:"degraded_reasons,omitempty"`
	AIModel         string   `json:"ai_model,omitempty"`
	AITransport     string   `json:"ai_transport,omitempty"`

	// Reviewed and the fields below answer "what did this run actually read?",
	// with the same vocabulary the review report uses.
	Reviewed     bool   `json:"reviewed"`
	ReviewMode   string `json:"review_mode,omitempty"`
	ReviewBase   string `json:"review_base,omitempty"`
	ReviewRange  string `json:"review_range,omitempty"`
	ReviewTarget string `json:"review_target,omitempty"`

	Verbose bool `json:"-"`
	Color   bool `json:"-"`
}

// FailedGates lists the gates that failed, in evaluation order.
func (r ReviewReport) FailedGates() []GateID {
	var out []GateID
	for _, gate := range r.Gates {
		if gate.Status == GateFail {
			out = append(out, gate.Gate)
		}
	}
	return out
}

type ReviewOptions struct {
	// Intent is the stated purpose of the change, usually the positional CLI
	// argument. Empty means it is derived from the branch and commit subjects.
	Intent string
	// Base reviews the range "<Base>...HEAD" instead of the working tree,
	// exactly like review's --base.
	Base      string
	SkipGates []GateID
	// Timeout bounds the whole run. Zero means defaultReviewTimeout.
	Timeout        time.Duration
	Verbose        bool
	ProgressWriter io.Writer
	Color          bool
}

const (
	defaultReviewTimeout = 5 * time.Minute
	// reviewRetrievalTimeout bounds the context-retrieval leg on its own:
	// grounding is a bonus, and a hung index must not stall a gate whose
	// deterministic half already finished.
	reviewRetrievalTimeout = 20 * time.Second
	// Code-health thresholds. A change past either of these is one a human
	// cannot meaningfully review in a sitting, which is a code-health violation
	// on its own — no model needed to argue it.
	reviewMaxReviewableLines = 1500
	reviewMaxReviewableFiles = 60
)

// CheckReview runs the exit gate against repoRoot's current change.
func CheckReview(ctx context.Context, repoRoot string, opts ReviewOptions) (ReviewReport, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return ReviewReport{}, fmt.Errorf("repo root is required")
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = defaultReviewTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	// Same guarantee as review: checking a checkout must not modify it.
	ctx, releaseIndex := withScratchGitIndex(ctx, repoRoot)
	defer releaseIndex()

	gateProgress(opts, "Resolving the current change")
	facts, err := scanRepo(ctx, repoRoot, "")
	if err != nil {
		return ReviewReport{}, err
	}
	policy := LoadReviewPolicy(repoRoot)
	changes := resolveChangeSet(ctx, repoRoot, opts.Base)
	report := ReviewReport{
		RepoRoot:     repoRoot,
		Reviewed:     changes.Reviewed(),
		ReviewMode:   changes.Mode,
		ReviewBase:   changes.Base,
		ReviewRange:  changes.Range,
		ReviewTarget: changes.Target,
		ChangedFiles: changes.Files,
		Verbose:      opts.Verbose,
		Color:        opts.Color,
	}
	if !changes.Reviewed() {
		// Nothing was inspected: no tools, no model, and a verdict that says
		// so. This is never a pass.
		report.Verdict = VerdictNothingToCheck
		return report, nil
	}

	reviewOpts := Options{ReviewPolicy: &policy, Color: opts.Color}
	diffs := collectDiffSnippets(ctx, repoRoot, changes.Files, reviewOpts, changes.Range)
	signals := collectReviewSignals(changes.Files, diffs, policy, reviewNumstat(ctx, repoRoot, changes.Range))
	report.DiffStats = signals.DiffStats
	report.Intent, report.IntentSource = resolveReviewIntent(ctx, repoRoot, opts.Intent, changes)

	skipped := map[GateID]struct{}{}
	for _, id := range opts.SkipGates {
		skipped[id] = struct{}{}
	}

	// The three expensive legs run together: the project's own checks, the
	// dependency audits, and context retrieval. Retrieval rides in the shadow
	// of the tool run — its snippets ground the one model call, so it adds no
	// wall clock of its own.
	var (
		toolResults  []StaticToolResult
		auditResults []StaticToolResult
		retrieved    []ContextSnippet
	)
	needTools := !bothGatesSkipped(skipped, GateCorrectness, GateCodeHealth)
	judge, judgeModel, judgeTransport, judgeUnavailable := reviewJudgeFactory()
	report.AIModel = judgeModel
	report.AITransport = judgeTransport
	var wg sync.WaitGroup
	if needTools {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gateProgress(opts, "Running project checks")
			toolResults = collectReviewToolResults(ctx, repoRoot, facts, reviewOpts, changes.Files)
		}()
	}
	if _, skip := skipped[GateSecurity]; !skip {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gateProgress(opts, "Auditing dependencies")
			auditResults = collectReviewAuditResults(ctx, repoRoot, changes.Files)
		}()
	}
	if judge != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			retrieved = collectReviewContext(ctx, repoRoot, facts, reviewOpts, changes, diffs)
		}()
	}
	wg.Wait()
	if err := ctx.Err(); err == context.Canceled {
		return report, err
	}

	gates := map[GateID]GateResult{
		GateCorrectness: reviewCorrectnessGate(toolResults, changes.Files),
		GateSecurity:    reviewSecurityGate(repoRoot, changes.Files, diffs, auditResults),
	}
	gates[GateCodeHealth] = reviewCodeHealthGate(toolResults, signals, changes.Files)
	gates[GateBackPressure] = GateResult{Gate: GateBackPressure, Title: gateTitle(GateBackPressure), Status: GatePass}
	gates[GateAccessibility] = reviewAccessibilityGate(signals)
	gates[GatePerformance] = reviewPerformanceGate(signals)

	// The judgment gates: one model call over the diff, grounded in the
	// deterministic results above and whatever retrieval found.
	aiGates := reviewAIGates(gates, skipped)
	if len(aiGates) > 0 {
		switch {
		case judge == nil && judgeUnavailable == "":
			// Deliberate opt-out (GX_REVIEW_AI=0): the AI gates are skipped and
			// the run is not degraded — the operator asked for exactly this.
			markReviewAIGatesSkipped(gates, aiGates, "AI judgment disabled (GX_REVIEW_AI=0)")
		case judge == nil:
			markReviewAIGatesSkipped(gates, aiGates, "AI judgment unavailable: "+judgeUnavailable)
			report.DegradedReasons = append(report.DegradedReasons, judgeUnavailable)
		default:
			gateProgress(opts, "Asking the AI judge")
			degraded := runReviewAIPass(ctx, judge, reviewAIPassInput{
				report:    &report,
				gates:     gates,
				aiGates:   aiGates,
				signals:   signals,
				diffs:     diffs,
				tools:     toolResults,
				policy:    policy,
				retrieved: retrieved,
				repoRoot:  repoRoot,
				changes:   changes,
			})
			report.DegradedReasons = append(report.DegradedReasons, degraded...)
		}
	}

	for _, id := range opts.SkipGates {
		result := gates[id]
		result.Gate = id
		result.Title = gateTitle(id)
		result.Status = GateSkipped
		result.SkipReason = "skipped by flag"
		result.Summary = ""
		result.Findings = nil
		result.Checks = nil
		gates[id] = result
	}

	diffsByFile := make(map[string]string, len(diffs))
	for _, snippet := range diffs {
		diffsByFile[snippet.File] = snippet.Diff
	}
	for _, id := range AllGateIDs() {
		result := gates[id]
		if result.Gate == "" {
			result.Gate = id
			result.Title = gateTitle(id)
			result.Status = GateSkipped
			result.SkipReason = "not evaluated"
		}
		// A finding about the change shows its diff hunk; only findings about
		// lines the change did not touch fall back to the current source.
		attachReviewDiffHunks(diffsByFile, result.Findings)
		attachReviewCodeExcerpts(repoRoot, result.Findings)
		report.Gates = append(report.Gates, result)
	}
	report.Verdict = reviewVerdict(report)
	return report, nil
}

const (
	// reviewExcerptContext is how many lines surround a finding's line in
	// its code excerpt; reviewExcerptMaxLineBytes keeps a minified or
	// generated line from turning the excerpt into a wall.
	reviewExcerptContext      = 2
	reviewExcerptMaxLineBytes = 200
)

// attachReviewDiffHunks gives each finding whose line is part of the
// change its unified-diff window, extracted from the same snippets the gates
// read.
func attachReviewDiffHunks(diffsByFile map[string]string, findings []Finding) {
	for i := range findings {
		finding := &findings[i]
		if finding.DiffHunk != "" || finding.File == "" || finding.Line <= 0 {
			continue
		}
		diff, ok := diffsByFile[finding.File]
		if !ok {
			continue
		}
		finding.DiffHunk = reviewDiffHunkForLine(diff, finding.Line)
	}
}

// attachReviewCodeExcerpts reads the source lines each finding points at,
// so the report shows the code being discussed rather than only naming it.
// Findings that already carry a diff hunk are left alone — the hunk is the
// better evidence. Best-effort by design: an unreadable file or a stale line
// number just leaves the excerpt empty.
func attachReviewCodeExcerpts(repoRoot string, findings []Finding) {
	for i := range findings {
		finding := &findings[i]
		if finding.DiffHunk != "" || finding.CodeExcerpt != "" || finding.File == "" || finding.Line <= 0 {
			continue
		}
		rel, ok := staticToolRelPath(repoRoot, finding.File)
		if !ok {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		if finding.Line > len(lines) {
			continue
		}
		start := finding.Line - reviewExcerptContext
		if start < 1 {
			start = 1
		}
		end := finding.Line + reviewExcerptContext
		if end > len(lines) {
			end = len(lines)
		}
		excerpt := make([]string, 0, end-start+1)
		for _, line := range lines[start-1 : end] {
			if len(line) > reviewExcerptMaxLineBytes {
				line = line[:reviewExcerptMaxLineBytes] + "…"
			}
			excerpt = append(excerpt, line)
		}
		finding.CodeExcerpt = strings.Join(excerpt, "\n")
		finding.CodeExcerptStart = start
	}
}

// reviewVerdict rolls the gates up. Real failures outrank everything;
// a degraded run with no failures must not claim ship.
func reviewVerdict(report ReviewReport) string {
	if len(report.FailedGates()) > 0 {
		return VerdictNoShip
	}
	if len(report.DegradedReasons) > 0 {
		return VerdictDegraded
	}
	return VerdictShip
}

// reviewAIGates lists the gates the model call must judge this run: the
// judgment gates that apply to this change and were not skipped by flag or
// already decided deterministically as FAIL. A gate the deterministic half
// already failed keeps its failure; the model can add findings to a passing
// gate but never overturn a command's exit code.
func reviewAIGates(gates map[GateID]GateResult, skipped map[GateID]struct{}) []GateID {
	var out []GateID
	for _, id := range []GateID{GateCodeHealth, GateBackPressure, GateAccessibility, GatePerformance} {
		if _, skip := skipped[id]; skip {
			continue
		}
		if gates[id].Status == GateSkipped {
			continue
		}
		out = append(out, id)
	}
	return out
}

func markReviewAIGatesSkipped(gates map[GateID]GateResult, ids []GateID, reason string) {
	for _, id := range ids {
		result := gates[id]
		if result.Status == GateFail {
			// A deterministic failure stands on its own; the missing model only
			// costs the judgment overlay.
			continue
		}
		if id == GateCodeHealth || id == GateAccessibility {
			// These two have a deterministic core that already answered; the
			// missing model is noted, not a reason to unsay it.
			result.Evidence = append(result.Evidence, Evidence{Label: "AI overlay", Value: reason})
			gates[id] = result
			continue
		}
		result.Status = GateSkipped
		result.SkipReason = reason
		gates[id] = result
	}
}

func bothGatesSkipped(skipped map[GateID]struct{}, ids ...GateID) bool {
	for _, id := range ids {
		if _, ok := skipped[id]; !ok {
			return false
		}
	}
	return true
}

// resolveReviewIntent picks what the change is supposed to be, most
// explicit source first: the caller's argument, then the branch name plus the
// range's commit subjects. Session context reaches the model through the
// retrieval pass instead of through this field, so a missing session can never
// block the gate.
func resolveReviewIntent(ctx context.Context, repoRoot, stated string, changes ChangeSet) (string, string) {
	stated = strings.TrimSpace(stated)
	if stated != "" {
		return stated, "stated"
	}
	var parts []string
	if branch := currentReviewBranch(ctx, repoRoot); branch != "" && branch != "HEAD" {
		parts = append(parts, "branch "+branch)
	}
	if subjects := reviewCommitSubjects(ctx, repoRoot, changes.Range); len(subjects) > 0 {
		parts = append(parts, "commits: "+strings.Join(subjects, "; "))
	}
	if changes.Mode == ReviewModeWorkingTree {
		parts = append(parts, "uncommitted working-tree change")
	}
	if len(parts) == 0 {
		return "", ""
	}
	return strings.Join(parts, " — "), "derived"
}

func currentReviewBranch(ctx context.Context, repoRoot string) string {
	cmd := gitCommand(ctx, repoRoot, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func reviewCommitSubjects(ctx context.Context, repoRoot, refRange string) []string {
	refRange = strings.TrimSpace(refRange)
	if refRange == "" {
		return nil
	}
	cmd := gitCommand(ctx, repoRoot, "log", "--format=%s", "--max-count=10", refRange)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	var subjects []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			subjects = append(subjects, line)
		}
	}
	return subjects
}

func gateProgress(opts ReviewOptions, message string) {
	if opts.ProgressWriter == nil || strings.TrimSpace(message) == "" {
		return
	}
	_, _ = io.WriteString(opts.ProgressWriter, strings.TrimSpace(message)+"\n")
}

func reviewStaticToolsDisabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_STATIC_TOOLS")), "0")
}
