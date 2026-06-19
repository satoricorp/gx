package codereview

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const maxContextSnippetBytes = 3200

type ReviewBrief struct {
	RepoRoot      string             `json:"repo_root"`
	Scope         string             `json:"scope"`
	Depth         string             `json:"depth"`
	Focus         string             `json:"focus,omitempty"`
	Static        StaticSnapshot     `json:"static"`
	Hints         []ReviewHint       `json:"hints"`
	Context       []ContextSnippet   `json:"context"`
	SourceCatalog []SourceBrief      `json:"source_catalog"`
	Rubric        ArchitectureRubric `json:"rubric"`
}

type StaticSnapshot struct {
	FileCount       int                `json:"file_count"`
	TestFileCount   int                `json:"test_file_count"`
	DependencyFiles []string           `json:"dependency_files"`
	Docs            []FilePresence     `json:"docs"`
	Modules         []ModuleSummary    `json:"modules"`
	ChangedFiles    []string           `json:"changed_files"`
	ToolResults     []StaticToolResult `json:"static_tool_results,omitempty"`
	CodeQuality     []CodeQualityHint  `json:"code_quality_hints,omitempty"`
}

type ModuleSummary struct {
	Path      string `json:"path"`
	GoFiles   int    `json:"go_files"`
	TestFiles int    `json:"test_files"`
	Reason    string `json:"reason,omitempty"`
}

type ReviewHint struct {
	Kind    string   `json:"kind"`
	Title   string   `json:"title"`
	Modules []string `json:"modules,omitempty"`
	Why     string   `json:"why"`
}

type CodeQualityHint struct {
	Kind   string `json:"kind"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	Text   string `json:"text"`
	Reason string `json:"reason"`
}

type ContextSnippet struct {
	Kind   string `json:"kind"`
	Ref    string `json:"ref"`
	Text   string `json:"text"`
	Source string `json:"source,omitempty"`
}

type SourceBrief struct {
	ID     string   `json:"id"`
	Scopes []string `json:"scopes"`
}

type ArchitectureRubric struct {
	Goal      string   `json:"goal"`
	Questions []string `json:"questions"`
	Reject    []string `json:"reject"`
	Output    string   `json:"output"`
}

type ContextRetriever interface {
	Retrieve(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error)
}

type LocalContextRetriever struct{}

func BuildReviewBrief(ctx context.Context, repoRoot string, opts Options, facts RepoFacts, sources []Source, retriever ContextRetriever) (ReviewBrief, error) {
	hints := reviewHints(facts)
	contextSnippets, err := retriever.Retrieve(ctx, repoRoot, opts, facts, hints)
	if err != nil {
		return ReviewBrief{}, err
	}
	return ReviewBrief{
		RepoRoot: repoRoot,
		Scope:    opts.Scope,
		Depth:    depthLabel(opts.Deep),
		Focus:    strings.TrimSpace(opts.Focus),
		Static: StaticSnapshot{
			FileCount:       facts.TrackedFileCount,
			TestFileCount:   facts.TestFileCount,
			DependencyFiles: facts.DependencyFiles,
			Docs:            facts.Docs,
			Modules:         moduleSummaries(facts),
			ChangedFiles:    changedFiles(ctx, repoRoot),
			ToolResults:     collectStaticToolResults(ctx, repoRoot, facts, opts),
			CodeQuality:     collectCodeQualityHints(repoRoot, facts, opts),
		},
		Hints:         hints,
		Context:       contextSnippets,
		SourceCatalog: sourceBriefs(sources),
		Rubric:        architectureRubric(),
	}, nil
}

func reviewHints(facts RepoFacts) []ReviewHint {
	var hints []ReviewHint
	heavy := topModules(facts.GoPackages, func(pkg PackageFact) bool {
		return pkg.GoFiles >= 6 && !ignorablePackage(pkg.Path)
	}, 5)
	if len(heavy) > 0 {
		hints = append(hints, ReviewHint{
			Kind:    "deepening_candidate",
			Title:   "implementation-heavy Modules",
			Modules: packagePaths(heavy),
			Why:     "large implementations can be deep Modules, but they can also indicate leaked workflow knowledge or shallow Interfaces; inspect before judging.",
		})
	}
	untested := topModules(facts.GoPackages, func(pkg PackageFact) bool {
		return pkg.GoFiles >= 2 && pkg.TestFiles == 0 && !ignorablePackage(pkg.Path)
	}, 5)
	if len(untested) > 0 {
		hints = append(hints, ReviewHint{
			Kind:    "test_surface_friction",
			Title:   "non-trivial Modules without colocated tests",
			Modules: packagePaths(untested),
			Why:     "missing tests may be ordinary backlog, or it may mean the Interface is awkward to exercise through the package seam.",
		})
	}
	generic := genericPackages(facts.GoPackages)
	if len(generic) > 0 {
		hints = append(hints, ReviewHint{
			Kind:    "interface_naming",
			Title:   "generic Module names",
			Modules: packagePaths(generic),
			Why:     "generic names can hide the concept callers should rely on at the Interface.",
		})
	}
	if !present(facts.Docs, "CONTEXT.md") {
		hints = append(hints, ReviewHint{
			Kind:  "domain_vocabulary",
			Title: "missing domain vocabulary",
			Why:   "only mention this if missing vocabulary blocks a concrete Module or Interface recommendation.",
		})
	}
	return hints
}

func (LocalContextRetriever) Retrieve(_ context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error) {
	var snippets []ContextSnippet
	for _, doc := range []string{"CONTEXT.md", "AGENTS.md", "README.md"} {
		if snippet, ok := readSnippet(repoRoot, doc, "repo_doc"); ok {
			snippets = append(snippets, snippet)
		}
	}
	for _, rel := range facts.DependencyFiles {
		if snippet, ok := readSnippet(repoRoot, rel, "dependency_manifest"); ok {
			snippets = append(snippets, snippet)
		}
	}
	for _, module := range modulesFromHints(hints) {
		for _, rel := range moduleFiles(facts.Files, module, opts.Deep) {
			if snippet, ok := readSnippet(repoRoot, rel, "module_file"); ok {
				snippets = append(snippets, snippet)
			}
		}
	}
	for _, rel := range qualityFiles(collectCodeQualityHints(repoRoot, facts, opts), opts.Deep) {
		if snippet, ok := readSnippet(repoRoot, rel, "code_quality_file"); ok {
			snippets = append(snippets, snippet)
		}
	}
	return snippets, nil
}

func readSnippet(repoRoot, rel, kind string) (ContextSnippet, bool) {
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
	if err != nil {
		return ContextSnippet{}, false
	}
	text := string(data)
	if len(text) > maxContextSnippetBytes {
		text = text[:maxContextSnippetBytes] + "\n[truncated]\n"
	}
	return ContextSnippet{Kind: kind, Ref: rel, Text: text, Source: "local"}, true
}

func moduleFiles(files []string, module string, deep bool) []string {
	var out []string
	limit := 2
	if deep {
		limit = 8
	}
	for _, file := range files {
		if !strings.HasSuffix(file, ".go") || !strings.HasPrefix(file, module+"/") {
			continue
		}
		out = append(out, file)
	}
	sort.Slice(out, func(i, j int) bool {
		iTest := isTestFile(out[i])
		jTest := isTestFile(out[j])
		if iTest != jTest {
			return !iTest
		}
		return out[i] < out[j]
	})
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}

func moduleSummaries(facts RepoFacts) []ModuleSummary {
	var out []ModuleSummary
	for _, pkg := range topModules(facts.GoPackages, func(pkg PackageFact) bool {
		return !ignorablePackage(pkg.Path)
	}, 12) {
		reason := ""
		if pkg.GoFiles >= 6 {
			reason = "implementation-heavy"
		} else if pkg.TestFiles == 0 && pkg.GoFiles >= 2 {
			reason = "no colocated tests"
		}
		out = append(out, ModuleSummary{
			Path:      pkg.Path,
			GoFiles:   pkg.GoFiles,
			TestFiles: pkg.TestFiles,
			Reason:    reason,
		})
	}
	return out
}

func topModules(packages []PackageFact, keep func(PackageFact) bool, limit int) []PackageFact {
	var out []PackageFact
	for _, pkg := range packages {
		if keep(pkg) {
			out = append(out, pkg)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].GoFiles != out[j].GoFiles {
			return out[i].GoFiles > out[j].GoFiles
		}
		if out[i].TestFiles != out[j].TestFiles {
			return out[i].TestFiles < out[j].TestFiles
		}
		return out[i].Path < out[j].Path
	})
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}

func genericPackages(packages []PackageFact) []PackageFact {
	genericNames := map[string]struct{}{
		"common": {}, "shared": {}, "utils": {}, "util": {}, "helpers": {}, "helper": {}, "types": {}, "models": {}, "lib": {},
	}
	var out []PackageFact
	for _, pkg := range packages {
		if _, ok := genericNames[filepath.Base(pkg.Path)]; ok {
			out = append(out, pkg)
		}
	}
	return out
}

func packagePaths(packages []PackageFact) []string {
	out := make([]string, 0, len(packages))
	for _, pkg := range packages {
		out = append(out, pkg.Path)
	}
	return out
}

func modulesFromHints(hints []ReviewHint) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, hint := range hints {
		for _, module := range hint.Modules {
			if _, ok := seen[module]; ok {
				continue
			}
			seen[module] = struct{}{}
			out = append(out, module)
		}
	}
	return out
}

func qualityFiles(hints []CodeQualityHint, deep bool) []string {
	seen := map[string]struct{}{}
	limit := 5
	if deep {
		limit = 12
	}
	var out []string
	for _, hint := range hints {
		if _, ok := seen[hint.File]; ok {
			continue
		}
		seen[hint.File] = struct{}{}
		out = append(out, hint.File)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func sourceBriefs(sources []Source) []SourceBrief {
	out := make([]SourceBrief, 0, len(sources))
	for _, source := range sources {
		out = append(out, SourceBrief{ID: source.ID, Scopes: source.Scopes})
	}
	return out
}

func architectureRubric() ArchitectureRubric {
	return ArchitectureRubric{
		Goal: "Find bugs, brittle code, poorly structured Modules, best-practice gaps, and deepening opportunities that improve locality, leverage, testability, and AI navigability.",
		Questions: []string{
			"Do static tools reveal failing tests, vet warnings, compile errors, or dependency problems?",
			"Which code looks brittle, non-idiomatic, or likely to break under realistic inputs?",
			"Where does understanding one concept require bouncing across many Modules?",
			"Where is the Interface nearly as complex as the Implementation?",
			"Apply the deletion test: if deleting the Module removes complexity, it may be shallow; if complexity spreads to callers, it earns its keep.",
			"Where do tests bypass the Interface or verify implementation details?",
			"Where does a seam have only one adapter and exist mostly as indirection?",
			"Which dependency category applies: in-process, local-substitutable, remote but owned, or true external?",
		},
		Reject: []string{
			"Do not turn static hints directly into findings.",
			"Do not call something best-practice advice unless it is tied to code, tool output, or retrieved context.",
			"Do not report missing docs unless it blocks a concrete Module or Interface recommendation.",
			"Do not expose source URLs or source titles.",
			"Do not use component, service, API, or boundary when Module, Interface, or seam fits.",
		},
		Output: "Return only concrete recommendations with title, summary, recommendation, optional evidence, and strength.",
	}
}

func formatModuleList(modules []string, limit int) string {
	if len(modules) == 0 {
		return ""
	}
	if limit <= 0 || len(modules) <= limit {
		return "`" + strings.Join(modules, "`, `") + "`"
	}
	return fmt.Sprintf("`%s` and %d more", strings.Join(modules[:limit], "`, `"), len(modules)-limit)
}
