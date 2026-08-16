package codereview

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	scanner      Scanner
	catalog      Catalog
	rules        []Rule
	retriever    ContextRetriever
	reviewer     AIReviewer
	judge        FindingJudge
	autoReviewer bool
}

type Scanner interface {
	Scan(ctx context.Context, repoRoot string, focus string) (RepoFacts, error)
}

type Catalog interface {
	SourcesForScopes(scopes []string) []Source
}

type ReviewContext struct {
	Options      Options
	Facts        RepoFacts
	Sources      []Source
	ActiveScopes []string
	Brief        ReviewBrief
	Triage       ChangeTriage
	Plan         ReviewExecutionPlan
}

func NewEngine() *Engine {
	return &Engine{
		scanner:      LocalScanner{},
		catalog:      StaticCatalog{},
		rules:        defaultRules(),
		retriever:    contextRetrieverFromEnv(),
		autoReviewer: true,
	}
}

func NewEngineWith(scanner Scanner, catalog Catalog, rules []Rule) *Engine {
	return &Engine{scanner: scanner, catalog: catalog, rules: rules}
}

func NewEngineWithReviewer(scanner Scanner, catalog Catalog, rules []Rule, retriever ContextRetriever, reviewer AIReviewer) *Engine {
	return &Engine{scanner: scanner, catalog: catalog, rules: rules, retriever: retriever, reviewer: reviewer}
}

func (e *Engine) Review(ctx context.Context, repoRoot string, opts Options) (Report, error) {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return Report{}, fmt.Errorf("repo root is required")
	}
	// Reviewing a checkout must not modify it. Every git call below inherits
	// this scratch index, so a stat-cache refresh cannot rewrite .git/index or
	// take .git/index.lock out from under a concurrent git. It sits here
	// rather than on the package-level Review so that callers holding an
	// Engine get the same guarantee.
	ctx, releaseIndex := withScratchGitIndex(ctx, repoRoot)
	defer releaseIndex()
	explicitScope := strings.TrimSpace(opts.Scope) != ""
	opts = normalizeOptions(opts)
	if err := ValidateOptions(opts); err != nil {
		return Report{}, err
	}
	scanner := e.scanner
	if scanner == nil {
		scanner = LocalScanner{}
	}
	catalog := e.catalog
	if catalog == nil {
		catalog = StaticCatalog{}
	}
	rules := e.rules
	if rules == nil {
		rules = defaultRules()
	}
	retriever := e.retriever
	if retriever == nil {
		retriever = LocalContextRetriever{}
	}

	reviewTimingReset()
	reviewProgress(opts, "Scanning repository")
	facts, err := scanner.Scan(ctx, repoRoot, strings.TrimSpace(opts.Focus))
	if err != nil {
		return Report{}, err
	}
	policy := LoadReviewPolicy(repoRoot)
	opts.ReviewPolicy = &policy
	changes := reviewChangeSet(ctx, repoRoot, opts, facts.TrackedFileCount)
	changed := changes.Files
	diffs := collectDiffSnippets(ctx, repoRoot, changed, opts, changes.Range)
	triage := ChangeTriage{}
	if triageEnabled() {
		triage = TriageChange(changed, diffs, opts)
	}
	plan := reviewPlanFor(opts, triage)
	if !triageEnabled() {
		plan.ActiveScopes = activeScopeList(opts)
		plan.RunStaticTools = true
		plan.RunAI = true
		plan.RunReviewResources = true
		plan.RiskTags = riskTagsForReview(changed, facts.DependencyFiles, opts)
	}
	if !changes.Reviewed() {
		// Nothing was inspected: never spend an AI call or a tool run on it,
		// and never let the result read like a clean review.
		plan.RunAI = false
		plan.RunStaticTools = false
	}
	active := plan.ActiveScopes
	if len(active) == 0 {
		active = activeScopeList(opts)
		plan.ActiveScopes = active
	}
	sources := catalog.SourcesForScopes(active)
	if !changes.Reviewed() {
		return Report{
			RepoRoot:         repoRoot,
			Scope:            opts.Scope,
			Deep:             opts.Deep,
			Focus:            strings.TrimSpace(opts.Focus),
			Prompt:           strings.TrimSpace(opts.Prompt),
			BaselineScopes:   baselineFor(opts.Scope),
			Docs:             facts.Docs,
			DependencyFiles:  facts.DependencyFiles,
			TestFileCount:    facts.TestFileCount,
			TrackedFileCount: facts.TrackedFileCount,
			Reviewer:         "none",
			Verbose:          opts.Verbose,
			Color:            opts.Color,
		}.applyChangeSet(changes), nil
	}
	if shouldShortCircuitDocsOnly(opts, explicitScope, plan) {
		return Report{
			RepoRoot:          repoRoot,
			Scope:             opts.Scope,
			Deep:              opts.Deep,
			Focus:             strings.TrimSpace(opts.Focus),
			Prompt:            strings.TrimSpace(opts.Prompt),
			BaselineScopes:    baselineFor(opts.Scope),
			Docs:              facts.Docs,
			DependencyFiles:   facts.DependencyFiles,
			TestFileCount:     facts.TestFileCount,
			TrackedFileCount:  facts.TrackedFileCount,
			ChangedFiles:      changed,
			ObservationLabels: facts.observations(),
			Findings:          nil,
			Sources:           nil,
			Reviewer:          "triage",
			Verbose:           opts.Verbose,
			Color:             opts.Color,
			Triage:            triage,
			NoFindingsMessage: "No material issues in this change (documentation-only).",
		}.applyChangeSet(changes), nil
	}
	reviewProgress(opts, "Building review context")
	input := RetrieveInput{
		RepoRoot:     repoRoot,
		Options:      opts,
		Facts:        facts,
		Hints:        reviewHints(facts),
		Plan:         plan,
		ChangedFiles: changed,
		DiffSnippets: diffs,
		DiffRange:    changes.Range,
		Evidence:     &EvidenceLog{},
	}
	brief, err := BuildReviewBrief(ctx, input, sources, retriever)
	if err != nil {
		return Report{}, err
	}
	reviewContext := ReviewContext{
		Options:      opts,
		Facts:        facts,
		Sources:      sources,
		ActiveScopes: active,
		Brief:        brief,
		Triage:       triage,
		Plan:         plan,
	}
	reviewProgress(opts, "Checking fallback review rules")
	findings := evaluateFindings(reviewContext, rules)
	reviewerLabel := "heuristic fallback"
	var degradedReasons []string
	reviewer := e.reviewer
	autoLoadedReviewer := false
	if reviewer == nil && e.autoReviewer {
		reviewer = reviewerFromEnvFast(opts.Fast)
		autoLoadedReviewer = true
	}
	wantAIReview := plan.RunAI && (!autoLoadedReviewer || aiReviewRequestedFromEnv())
	aiConfigured := aiReviewRequestedFromEnv()
	// Captured before the review runs so the report can say which models and
	// which wire were used even when the call fails — that is precisely the
	// case where the reader needs to know.
	reviewerInfo := reviewerInfoFromReviewer(reviewer)
	var coverage Coverage
	if wantAIReview {
		if !reviewerAvailable(reviewer) {
			if aiConfigured {
				// Bedrock is the only review provider, so an unavailable
				// reviewer means nothing reviewed this change. "AI reviewer not
				// configured" named neither the provider nor the fix, and a
				// report with only deterministic findings reads like a clean
				// review; say which credential is missing instead.
				degradedReasons = append(degradedReasons, ReviewerUnavailableReason(reviewer))
			}
		} else {
			// The subject is split into as many calls as it takes to read all
			// of it, rather than trimmed to fit one. Ordinary changes plan a
			// single shard and behave exactly as before.
			shards, planned := planReviewShards(brief, opts)
			reviewProgress(opts, "Asking AI reviewer")
			result := runShardedReview(ctx, reviewer, shards, planned, opts)
			coverage = result.Coverage
			if result.Err != nil {
				if aiConfigured {
					degradedReasons = append(degradedReasons, formatReviewerDegradation(result.Err))
				}
			} else {
				// The label says which reviewers ran, not what they happened to
				// find. Setting it from the finding count meant a healthy review
				// that correctly found nothing reported "results are from
				// deterministic checks only" — the same sentence a real AI outage
				// produces, which teaches readers to discount it when it is true.
				findings = mergeFindings(findings, result.Findings)
				reviewerLabel = "heuristic+ai"
			}
			if coverage.ShardsFailed > 0 && coverage.ShardsFailed < coverage.Shards && aiConfigured {
				degradedReasons = append(degradedReasons, fmt.Sprintf("%d of %d parallel reviews failed", coverage.ShardsFailed, coverage.Shards))
			}
			// A reviewer leg that failed while the other answered produces a
			// complete-looking review from half the panel. Say so: it changes
			// how much the result is worth, and it is the only way a signing or
			// model-access fault on one leg ever becomes visible.
			for _, reason := range PartialReviewerFailures(reviewer) {
				degradedReasons = append(degradedReasons, "one reviewer did not run — "+reason)
			}
		}
	}
	findings = validateFindingAnchors(reviewContext, findings)
	findings = filterPatchFocusedFindings(reviewContext, findings)
	blocking, advisory := splitBlockingToolFindings(findings)
	judge := e.judge
	if judge == nil {
		judge = judgeFromEnv()
	}
	// The panel's own adjudicator, not a second one built from the environment:
	// the same "one plan, one wire" rule the reviewer legs and the judge follow.
	// It is nil whenever no AI panel ran, which is also when there are no
	// cross-reviewer duplicates to collapse.
	deduped := prepareFindingsForJudge(ctx, reviewerDuplicateAdjudicator(reviewer), advisory)
	advisory = deduped.Findings
	// Duplicates that could not be adjudicated are shipped, so say so. A review
	// that merged nothing because it could not ask must not look like a review
	// that found nothing to merge.
	if unresolved := deduped.Unresolved(); unresolved > 0 && reviewerLabel == "heuristic+ai" {
		reason := fmt.Sprintf("%d candidate duplicate pair(s) were not adjudicated, so near-duplicate findings may remain", unresolved)
		if deduped.Err != nil {
			reason += fmt.Sprintf(": %v", deduped.Err)
		}
		degradedReasons = append(degradedReasons, reason)
	}
	// Fast reviews skip verification. It is the single most expensive stage
	// after the review call itself (measured ~15s of a 47s review), and its
	// value is filtering a panel's disagreements — which a one-leg review does
	// not produce. The findings it would have graded are still reported; they
	// are reported unverified, which the report says.
	if opts.Fast {
		reviewProgress(opts, "Skipping verification (fast review)")
	}
	if !opts.Fast && !judgeDisabledFromEnv() && judgeAvailable(judge) && len(advisory) > 0 {
		reviewProgress(opts, "Verifying review findings")
		candidates := advisory
		// Verification runs in concurrent batches so that a large finding set
		// neither overruns the judge's output budget nor pays for each batch in
		// series. See judgeBatchSize for the measurements behind the size.
		outcome := runJudge(ctx, judge, reviewContext, candidates)
		// Precision filter: surface only confirmed, human-worthy findings,
		// impact-gated and ranked (not capped). No findings is a valid result.
		advisory = outcome.Judged
		if len(outcome.Unjudged) > 0 {
			// Unverified candidates fall back to the deduped, strength-capped
			// set rather than being dropped, because an unreachable judge is
			// not a verdict. Marked so a reader can tell them from findings the
			// judge actually checked.
			unjudged := capAdvisoryFindings(outcome.Unjudged, opts.MaxFindings)
			for i := range unjudged {
				unjudged[i].JudgeVerdict = "unverified"
			}
			advisory = append(advisory, unjudged...)
		}
		if outcome.BatchesFailed > 0 {
			// A failed verification used to be invisible: findings silently
			// collapsed to at most the findings cap and the report still
			// read as a completed review. Say so instead.
			degradedReasons = append(degradedReasons, fmt.Sprintf(
				"%d of %d finding verification batches failed: %v",
				outcome.BatchesFailed, outcome.Batches, outcome.Err))
		}
		if outcome.Omitted > 0 {
			// The call succeeded, the JSON parsed, and the model simply left
			// candidates out of its reply. That is a reliability failure and is
			// worth saying so plainly: re-asking generally answers them, so a
			// review reporting this is one call short of being fully judged.
			degradedReasons = append(degradedReasons, fmt.Sprintf(
				"%d candidate finding(s) were omitted from the verification model's reply, so they are reported unverified",
				outcome.Omitted))
		}
		if outcome.Abstained > 0 {
			// Reported separately because it is not the same problem. Here the
			// judge did answer — it said it could not check the claim, which
			// usually means it was not given the evidence to check it against.
			// Re-asking buys the same answer; supplying the file content is what
			// changes it.
			degradedReasons = append(degradedReasons, fmt.Sprintf(
				"%d candidate finding(s) got an explicit \"cannot verify\" from the verification model, so they are reported unverified",
				outcome.Abstained))
			if judgeTraceEnabled() {
				// Under trace, say what each one actually needed. An abstention
				// count says a review was not fully verified; only the reason
				// says whether that is fixable by sending more code.
				for _, note := range outcome.AbstentionNotes {
					degradedReasons = append(degradedReasons, "judge-trace: "+note)
				}
			}
		}
	} else {
		advisory = capAdvisoryFindings(advisory, opts.MaxFindings)
	}
	// Lanes are assigned here and nowhere else: this is the one line every
	// review path — full panel, --fast, judge disabled or unreachable, no AI at
	// all — passes through with its findings final. Earlier is too early (the
	// judge rewrites strength and verdicts, and --fast never reaches it); a hook
	// inside applyJudgeResults would miss the unjudged fallback and the Blocking
	// tool findings split out above. See lanes.go for the rule.
	findings = AssignLanes(append(blocking, advisory...))
	// The story lane. This path asks the reviewer for findings (Review), not
	// for a summary (ReviewForSummary), so the model's story items are not
	// available here — they travel in PRSummaryReview.Story for the pr_summary
	// callers. What is assembled here is the part no model decides: the
	// mandatory items, hunks attached from the same diffs the reviewer read.
	// Nil when there is nothing to say.
	story := assembleReportStory(opts.ReviewPolicy, changed, diffs, nil)

	return Report{
		RepoRoot:          repoRoot,
		Scope:             opts.Scope,
		Deep:              opts.Deep,
		Focus:             strings.TrimSpace(opts.Focus),
		Prompt:            strings.TrimSpace(opts.Prompt),
		BaselineScopes:    baselineFor(opts.Scope),
		Docs:              facts.Docs,
		DependencyFiles:   facts.DependencyFiles,
		TestFileCount:     facts.TestFileCount,
		TrackedFileCount:  facts.TrackedFileCount,
		ChangedFiles:      changed,
		ObservationLabels: facts.observations(),
		Findings:          findings,
		Story:             story,
		Sources:           sources,
		SourceRefs:        brief.SourceRefs,
		Reviewer:          reviewerLabel,
		ReviewModels:      reviewerInfo.Models,
		ReviewTransport:   reviewerInfo.Transport,
		ContextSnippets:   len(brief.Context),
		Verbose:           opts.Verbose,
		Color:             opts.Color,
		Triage:            triage,
		DegradedReasons:   degradedReasons,
		Evidence:          brief.Evidence,
		Coverage:          coverage,
	}.applyChangeSet(changes), nil
}

func Review(ctx context.Context, repoRoot string, opts Options) (Report, error) {
	return NewEngine().Review(ctx, repoRoot, opts)
}

// reviewChangeSet resolves what this review looks at.
//
// Options.WholeRepo is the one instruction that names the subject outright, so
// it wins over everything, including a dirty working tree and an explicit
// base. It is the only way to reach a whole-repo review while there is a diff
// to be found.
//
// Otherwise the resolved change set stands as-is. A scope-, prompt-, or
// deep-directed review falls back to the repository only when nothing was
// found at all: those flags pick a lens and a depth, not a subject, so while a
// diff exists the diff remains what is being reviewed. An explicit base is the
// exception in the other direction: the caller named the range, so an empty
// range means exactly what it says.
// repoFileCount is what the scan found in the repository; a whole-repo review
// of a repository with nothing in it is still nothing to review.
func reviewChangeSet(ctx context.Context, repoRoot string, opts Options, repoFileCount int) ChangeSet {
	explicitBase := strings.TrimSpace(opts.Base) != ""
	changes := resolveChangeSet(ctx, repoRoot, opts.Base)
	if opts.WholeRepo {
		return wholeRepoChangeSet(changes, explicitBase, repoFileCount > 0)
	}
	if changes.Reviewed() || explicitBase || opts.PatchFocused {
		return changes
	}
	changes.Mode = ReviewModeRepo
	changes.Target = "the repository"
	return changes
}

// applyChangeSet records what the review read, so callers can tell a clean
// review apart from one that never opened a diff.
func (r Report) applyChangeSet(changes ChangeSet) Report {
	r.Reviewed = changes.Reviewed()
	r.ReviewMode = changes.Mode
	r.ReviewBase = changes.Base
	r.ReviewRange = changes.Range
	r.ReviewTarget = changes.Target
	return r
}

// patchFocusedScopes is what the default review looks for. Security leads it,
// and any wider review must keep it: asking for more of the repository must
// never buy a narrower set of concerns.
func patchFocusedScopes() []string {
	return []string{"security", "dependencies", "testing", "maintainability"}
}

func activeScopeList(opts Options) []string {
	if opts.WholeRepo {
		// The requested scope plus everything the default already covers.
		// Without the union, `--repo` would silently drop `security` — the
		// scope list gates rules, findings, and the model's source catalog —
		// which is a strange thing to lose by asking for a broader review.
		return dedupeScopes(append([]string{opts.Scope}, patchFocusedScopes()...))
	}
	if opts.PatchFocused {
		return patchFocusedScopes()
	}
	return append([]string{opts.Scope}, baselineFor(opts.Scope)...)
}

func dedupeScopes(scopes []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	return out
}

func reviewProgress(opts Options, message string) {
	if strings.TrimSpace(message) == "" {
		return
	}
	reviewTimingMark(message)
	if opts.ProgressWriter == nil {
		return
	}
	_, _ = io.WriteString(opts.ProgressWriter, strings.TrimSpace(message)+"\n")
}

// Review phase timing, enabled with GX_REVIEW_TIMING=1.
//
// Progress messages already mark every phase boundary, so timing them costs one
// clock read and answers the only question that matters when a review feels
// slow: which phase spent the time. Without it the answer is a stopwatch and a
// guess, and the guess is usually wrong — the first profile of a 2m30s review
// found the model calls were not the problem at all.
var (
	reviewTimingMu    sync.Mutex
	reviewTimingStart time.Time
	reviewTimingLast  time.Time
)

// judgeTraceEnabled adds the judge's own reason for each abstention to the
// degraded reasons, enabled with GX_REVIEW_JUDGE_TRACE=1.
//
// Off by default because it is diagnostic output, not review output: a reader
// wants to know a finding went unverified, not to read the paragraph the
// verification model wrote about it.
func judgeTraceEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GX_REVIEW_JUDGE_TRACE"))) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}

func reviewTimingEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GX_REVIEW_TIMING"))) {
	case "1", "true", "yes", "on":
		return true
	}
	return false
}

// reviewTimingReset starts a new timing run. Without it a process that reviews
// more than once — the MCP server, the test suite — would report every later
// review's phases as offsets from the first review's start.
func reviewTimingReset() {
	if !reviewTimingEnabled() {
		return
	}
	reviewTimingMu.Lock()
	defer reviewTimingMu.Unlock()
	reviewTimingStart, reviewTimingLast = time.Time{}, time.Time{}
}

func reviewTimingMark(message string) {
	if !reviewTimingEnabled() {
		return
	}
	now := time.Now()
	reviewTimingMu.Lock()
	defer reviewTimingMu.Unlock()
	if reviewTimingStart.IsZero() {
		reviewTimingStart, reviewTimingLast = now, now
	}
	fmt.Fprintf(os.Stderr, "[timing] %7.1fs (+%5.1fs) %s\n",
		now.Sub(reviewTimingStart).Seconds(), now.Sub(reviewTimingLast).Seconds(), strings.TrimSpace(message))
	reviewTimingLast = now
}
