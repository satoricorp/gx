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
	changed := reviewChangedFiles(ctx, repoRoot)
	diffs := collectDiffSnippets(ctx, repoRoot, changed, opts.Deep)
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
	active := plan.ActiveScopes
	if len(active) == 0 {
		active = activeScopeList(opts)
		plan.ActiveScopes = active
	}
	sources := catalog.SourcesForScopes(active)
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
		}, nil
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
	if wantAIReview {
		if !reviewerAvailable(reviewer) {
			if aiConfigured {
				degradedReasons = append(degradedReasons, "AI reviewer not configured")
			}
		} else {
			reviewProgress(opts, "Asking AI reviewer")
			aiFindings, err := reviewer.Review(ctx, brief)
			if err != nil {
				if aiConfigured {
					degradedReasons = append(degradedReasons, formatReviewerDegradation(err))
				}
			} else if len(aiFindings) > 0 {
				findings = mergeFindings(findings, aiFindings)
				reviewerLabel = "heuristic+ai"
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
		if results, err := judge.Judge(ctx, buildJudgeRequest(reviewContext, candidates)); err == nil {
			// Precision filter: surface only confirmed, human-worthy findings,
			// impact-gated and ranked (not capped). No findings is a valid result.
			advisory = applyJudgeResults(candidates, results)
		} else {
			// Judge call failed: fall back to the deduped, strength-capped set
			// rather than the pre-dedup findings.
			advisory = capAdvisoryFindings(candidates)
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
		ContextSnippets:   len(brief.Context),
		Verbose:           opts.Verbose,
		Color:             opts.Color,
		Triage:            triage,
		DegradedReasons:   degradedReasons,
	}, nil
}

func Review(ctx context.Context, repoRoot string, opts Options) (Report, error) {
	return NewEngine().Review(ctx, repoRoot, opts)
}

func activeScopeList(opts Options) []string {
	if opts.PatchFocused {
		return []string{"security", "dependencies", "testing", "maintainability"}
	}
	return append([]string{opts.Scope}, baselineFor(opts.Scope)...)
}

func reviewProgress(opts Options, message string) {
	if opts.ProgressWriter == nil || strings.TrimSpace(message) == "" {
		return
	}
	_, _ = io.WriteString(opts.ProgressWriter, strings.TrimSpace(message)+"\n")
}
