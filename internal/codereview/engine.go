package codereview

import (
	"context"
	"fmt"
	"io"
	"strings"
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

	reviewProgress(opts, "Scanning repository")
	facts, err := scanner.Scan(ctx, repoRoot, strings.TrimSpace(opts.Focus))
	if err != nil {
		return Report{}, err
	}
	policy := LoadReviewPolicy(ctx, repoRoot)
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
		reviewer = reviewerFromEnvWithPolicy(&policy)
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
			} else if len(result.Findings) > 0 {
				findings = mergeFindings(findings, result.Findings)
				reviewerLabel = "heuristic+ai"
			}
			if coverage.ShardsFailed > 0 && coverage.ShardsFailed < coverage.Shards && aiConfigured {
				degradedReasons = append(degradedReasons, fmt.Sprintf("%d of %d parallel reviews failed", coverage.ShardsFailed, coverage.Shards))
			}
		}
	}
	findings = validateFindingAnchors(reviewContext, findings)
	findings = filterPatchFocusedFindings(reviewContext, findings)
	blocking, advisory := splitBlockingToolFindings(findings)
	judge := e.judge
	if judge == nil {
		judge = judgeFromEnvWithPolicy(&policy)
	}
	if !judgeDisabledFromEnv() && judgeAvailable(judge) && len(advisory) > 0 {
		reviewProgress(opts, "Verifying review findings")
		candidates := prepareFindingsForJudge(reviewContext, advisory)
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
			// not a verdict.
			advisory = append(advisory, capAdvisoryFindings(outcome.Unjudged)...)
		}
		if outcome.BatchesFailed > 0 {
			// A failed verification used to be invisible: findings silently
			// collapsed to at most maxAdvisoryFindings and the report still
			// read as a completed review. Say so instead.
			degradedReasons = append(degradedReasons, fmt.Sprintf(
				"%d of %d finding verification batches failed: %v",
				outcome.BatchesFailed, outcome.Batches, outcome.Err))
		}
	} else {
		advisory = capAdvisoryFindings(prepareFindingsForJudge(reviewContext, advisory))
	}
	findings = append(blocking, advisory...)

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
	if opts.ProgressWriter == nil || strings.TrimSpace(message) == "" {
		return
	}
	_, _ = io.WriteString(opts.ProgressWriter, strings.TrimSpace(message)+"\n")
}
