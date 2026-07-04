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
	RiskTags     []string
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
	if opts.Format == "html" {
		return Report{}, fmt.Errorf("html review output is not implemented yet")
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
	var triage ChangeTriage
	if triageEnabled() {
		triage = TriageChange(changed, diffs, opts)
	}
	plan := reviewPlanFor(opts, triage)
	if len(plan.RiskTags) == 0 {
		plan.RiskTags = riskTagsForReview(changed, facts.DependencyFiles, opts)
	}
	active := plan.ActiveScopes
	if len(active) == 0 {
		active = activeScopeList(opts)
		plan.ActiveScopes = active
	}
	if shouldShortCircuitDocsOnly(opts, explicitScope, plan) {
		return Report{
			RepoRoot:          repoRoot,
			Scope:             opts.Scope,
			Format:            opts.Format,
			Deep:              opts.Deep,
			Since:             strings.TrimSpace(opts.Since),
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
			Reviewer:          "triage",
			Verbose:           opts.Verbose,
			Color:             opts.Color,
			Triage:            triage,
			RiskTags:          plan.RiskTags,
			NoFindingsMessage: "No material issues in this change (documentation-only).",
		}, nil
	}
	sources := catalog.SourcesForScopes(active)
	reviewProgress(opts, "Building review context")
	retrieveInput := RetrieveInput{
		RepoRoot:     repoRoot,
		Options:      opts,
		Facts:        facts,
		Plan:         plan,
		ChangedFiles: changed,
		DiffSnippets: diffs,
	}
	brief, err := BuildReviewBrief(ctx, repoRoot, opts, facts, sources, retriever, retrieveInput)
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
		RiskTags:     plan.RiskTags,
	}
	reviewProgress(opts, "Checking fallback review rules")
	findings := evaluateFindings(reviewContext, rules)
	reviewerLabel := "heuristic fallback"
	reviewer := e.reviewer
	if reviewer == nil && e.autoReviewer {
		reviewer = reviewerFromEnvWithPolicy(&policy)
	}
	if reviewer != nil && plan.RunAI {
		reviewProgress(opts, "Asking AI reviewer")
		if aiFindings, err := reviewer.Review(ctx, brief); err == nil && len(aiFindings) > 0 {
			findings = mergeFindings(findings, aiFindings)
			reviewerLabel = "heuristic+ai"
		}
	}
	findings = filterPatchFocusedFindings(reviewContext, findings)

	return Report{
		RepoRoot:          repoRoot,
		Scope:             opts.Scope,
		Format:            opts.Format,
		Deep:              opts.Deep,
		Since:             strings.TrimSpace(opts.Since),
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
		RiskTags:          plan.RiskTags,
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
