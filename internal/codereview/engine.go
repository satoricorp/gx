package codereview

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type Engine struct {
	scanner   Scanner
	catalog   Catalog
	rules     []Rule
	retriever ContextRetriever
	reviewer  AIReviewer
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
}

func NewEngine() *Engine {
	return &Engine{
		scanner:   LocalScanner{},
		catalog:   StaticCatalog{},
		rules:     defaultRules(),
		retriever: contextRetrieverFromEnv(),
		reviewer:  reviewerFromEnv(),
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
	active := activeScopeList(opts.Scope)
	sources := catalog.SourcesForScopes(active)
	reviewProgress(opts, "Building review context")
	brief, err := BuildReviewBrief(ctx, repoRoot, opts, facts, sources, retriever)
	if err != nil {
		return Report{}, err
	}
	reviewContext := ReviewContext{
		Options:      opts,
		Facts:        facts,
		Sources:      sources,
		ActiveScopes: active,
		Brief:        brief,
	}
	reviewProgress(opts, "Checking fallback review rules")
	findings := evaluateFindings(reviewContext, rules)
	reviewerLabel := "heuristic fallback"
	if e.reviewer != nil {
		reviewProgress(opts, "Asking AI reviewer")
		if aiFindings, err := e.reviewer.Review(ctx, brief); err == nil && len(aiFindings) > 0 {
			findings = aiFindings
			reviewerLabel = "ai"
		}
	}

	return Report{
		RepoRoot:          repoRoot,
		Scope:             opts.Scope,
		Format:            opts.Format,
		Deep:              opts.Deep,
		Since:             strings.TrimSpace(opts.Since),
		Focus:             strings.TrimSpace(opts.Focus),
		BaselineScopes:    baselineFor(opts.Scope),
		Docs:              facts.Docs,
		DependencyFiles:   facts.DependencyFiles,
		TestFileCount:     facts.TestFileCount,
		TrackedFileCount:  facts.TrackedFileCount,
		ChangedFiles:      changedFiles(ctx, repoRoot),
		ObservationLabels: facts.observations(),
		Findings:          findings,
		Sources:           sources,
		Reviewer:          reviewerLabel,
		ContextSnippets:   len(brief.Context),
		Verbose:           opts.Verbose,
		Color:             opts.Color,
	}, nil
}

func Review(ctx context.Context, repoRoot string, opts Options) (Report, error) {
	return NewEngine().Review(ctx, repoRoot, opts)
}

func activeScopeList(scope string) []string {
	return append([]string{scope}, baselineFor(scope)...)
}

func reviewProgress(opts Options, message string) {
	if opts.ProgressWriter == nil || strings.TrimSpace(message) == "" {
		return
	}
	_, _ = io.WriteString(opts.ProgressWriter, strings.TrimSpace(message)+"\n")
}
