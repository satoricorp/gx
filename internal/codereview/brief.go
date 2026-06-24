package codereview

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const maxContextSnippetBytes = 3200
const maxDiffSnippetBytes = 5000
const maxDiffSnippetFiles = 10
const maxDeepDiffSnippetFiles = 24

type ReviewBrief struct {
	RepoRoot      string             `json:"repo_root"`
	Scope         string             `json:"scope"`
	Depth         string             `json:"depth"`
	ReviewProfile string             `json:"review_profile"`
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
	ADRFiles        []string           `json:"adr_files,omitempty"`
	Modules         []ModuleSummary    `json:"modules"`
	ChangedFiles    []string           `json:"changed_files"`
	DiffSnippets    []DiffSnippet      `json:"diff_snippets,omitempty"`
	ToolResults     []StaticToolResult `json:"static_tool_results,omitempty"`
	CodeQuality     []CodeQualityHint  `json:"code_quality_hints,omitempty"`
}

type ModuleSummary struct {
	Path      string `json:"path"`
	GoFiles   int    `json:"go_files"`
	TestFiles int    `json:"test_files"`
	Reason    string `json:"reason,omitempty"`
}

type DiffSnippet struct {
	File string `json:"file"`
	Diff string `json:"diff"`
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
	Kind        string `json:"kind"`
	Ref         string `json:"ref"`
	Text        string `json:"text"`
	Source      string `json:"source,omitempty"`
	SourceLabel string `json:"source_label,omitempty"`
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
	contextSnippets = labelContextSnippets(contextSnippets)
	changed := changedFiles(ctx, repoRoot)
	return ReviewBrief{
		RepoRoot:      repoRoot,
		Scope:         opts.Scope,
		Depth:         depthLabel(opts.Deep),
		ReviewProfile: reviewProfile(opts),
		Focus:         strings.TrimSpace(opts.Focus),
		Static: StaticSnapshot{
			FileCount:       facts.TrackedFileCount,
			TestFileCount:   facts.TestFileCount,
			DependencyFiles: facts.DependencyFiles,
			Docs:            facts.Docs,
			ADRFiles:        facts.ADRFiles,
			Modules:         moduleSummaries(facts),
			ChangedFiles:    changed,
			DiffSnippets:    collectDiffSnippets(ctx, repoRoot, changed, opts.Deep),
			ToolResults:     collectStaticToolResults(ctx, repoRoot, facts, opts),
			CodeQuality:     collectCodeQualityHints(repoRoot, facts, opts),
		},
		Hints:         hints,
		Context:       contextSnippets,
		SourceCatalog: sourceBriefs(sources),
		Rubric:        reviewRubric(opts),
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
	return hints
}

func (LocalContextRetriever) Retrieve(_ context.Context, repoRoot string, opts Options, facts RepoFacts, hints []ReviewHint) ([]ContextSnippet, error) {
	var snippets []ContextSnippet
	for _, doc := range []struct {
		path string
		kind string
	}{
		{path: "CONTEXT.md", kind: "domain_doc"},
		{path: "REVIEW.md", kind: "repo_doc"},
		{path: "AGENTS.md", kind: "repo_doc"},
		{path: "README.md", kind: "repo_doc"},
	} {
		if snippet, ok := readSnippet(repoRoot, doc.path, doc.kind); ok {
			snippets = append(snippets, snippet)
		}
	}
	for _, rel := range facts.ADRFiles {
		if snippet, ok := readSnippet(repoRoot, rel, "adr"); ok {
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

func collectDiffSnippets(ctx context.Context, repoRoot string, files []string, deep bool) []DiffSnippet {
	files = normalizedChangedFiles(files)
	limit := maxDiffSnippetFiles
	if deep {
		limit = maxDeepDiffSnippetFiles
	}
	if len(files) > limit {
		files = files[:limit]
	}
	var snippets []DiffSnippet
	for _, file := range files {
		diff := fileDiff(ctx, repoRoot, file)
		if strings.TrimSpace(diff) == "" {
			diff = fileContentSnippet(repoRoot, file)
		}
		diff = truncateReviewText(diff, maxDiffSnippetBytes)
		if strings.TrimSpace(diff) == "" {
			continue
		}
		snippets = append(snippets, DiffSnippet{File: file, Diff: diff})
	}
	return snippets
}

func fileDiff(ctx context.Context, repoRoot, file string) string {
	var parts []string
	for _, args := range [][]string{
		{"diff", "--no-ext-diff", "--", file},
		{"diff", "--cached", "--no-ext-diff", "--", file},
	} {
		cmd := exec.CommandContext(ctx, "git", args...)
		cmd.Dir = repoRoot
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		_ = cmd.Run()
		if text := strings.TrimSpace(out.String()); text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, "\n")
}

func fileContentSnippet(repoRoot, file string) string {
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(file)))
	if err != nil {
		return ""
	}
	return "No git diff was available for this changed file. Current file content:\n" + string(data)
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

func labelContextSnippets(snippets []ContextSnippet) []ContextSnippet {
	counts := map[string]int{}
	out := make([]ContextSnippet, 0, len(snippets))
	for _, snippet := range snippets {
		prefix := contextLabelPrefix(snippet)
		counts[prefix]++
		label := fmt.Sprintf("%s%d", prefix, counts[prefix])
		snippet.SourceLabel = label
		header := fmt.Sprintf("[%s] kind=%s ref=%s source=%s", label, strings.TrimSpace(snippet.Kind), strings.TrimSpace(snippet.Ref), strings.TrimSpace(snippet.Source))
		if !strings.HasPrefix(strings.TrimSpace(snippet.Text), "["+label+"]") {
			snippet.Text = strings.TrimSpace(header + "\n" + snippet.Text)
		}
		out = append(out, snippet)
	}
	return out
}

func contextLabelPrefix(snippet ContextSnippet) string {
	switch {
	case strings.HasPrefix(strings.TrimSpace(snippet.Source), "turbopuffer:"):
		return "R"
	case strings.EqualFold(strings.TrimSpace(snippet.Source), "local"):
		return "L"
	default:
		return "C"
	}
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

func reviewRubric(opts Options) ArchitectureRubric {
	if opts.Deep {
		return deepReviewRubric()
	}
	if opts.PatchFocused {
		return patchReviewRubric()
	}
	return scopedReviewRubric(opts.Scope)
}

func patchReviewRubric() ArchitectureRubric {
	return ArchitectureRubric{
		Goal: "Review the current patch for concrete regressions, security issues, data correctness problems, race/idempotency risks, error-handling gaps, missing tests, and observability gaps.",
		Questions: []string{
			"What behavior changed in static.diff_snippets and static.changed_files?",
			"Can the changed code fail for realistic inputs, missing state, retries, concurrency, timeouts, or partial external failures?",
			"Does the patch change authn/authz, token handling, secrets, webhook verification, database writes, migrations, file IO, shell execution, network calls, or deploy/CI behavior?",
			"What test would fail before the correct fix and pass after it?",
			"Would an operator have enough logs, errors, metrics, or status to diagnose this changed behavior in production?",
			"Is any architecture recommendation directly required to fix a bug or review risk introduced by this patch?",
		},
		Reject: []string{
			"Do not report repo-wide architecture, naming, docs, or cleanup advice unless the changed diff directly creates the risk.",
			"Do not turn file counts, missing docs, missing tests, or large Modules directly into findings.",
			"Do not recommend broad refactors when a localized fix or test would address the changed behavior.",
			"Do not expose source URLs or source titles.",
			"If there is no concrete patch-grounded issue, return no recommendations.",
		},
		Output: "Return only concrete patch-grounded recommendations with title, summary, benefit, recommendation, optional evidence, and strength.",
	}
}

func scopedReviewRubric(scope string) ArchitectureRubric {
	return ArchitectureRubric{
		Goal: "Review the requested scope for concrete findings tied to code, tool output, local project policy, or retrieved review resources.",
		Questions: []string{
			"Do static tools reveal failing tests, vet warnings, compile errors, or dependency problems?",
			"Which code looks brittle, non-idiomatic, insecure, undertested, or likely to break under realistic inputs?",
			"Which findings are specific to the requested scope `" + strings.TrimSpace(scope) + "`?",
			"Which recommendation has a clear first file to edit and verification command to run?",
		},
		Reject: []string{
			"Do not turn static hints directly into findings.",
			"Do not call something best-practice advice unless it is tied to code, tool output, or retrieved context.",
			"Do not expose source URLs or source titles.",
			"Do not produce generic audit facts.",
		},
		Output: "Return only concrete scoped recommendations with title, summary, benefit, recommendation, optional evidence, and strength.",
	}
}

func deepReviewRubric() ArchitectureRubric {
	return ArchitectureRubric{
		Goal: "Run a full-spectrum review across bugs, security, data integrity, race/idempotency, architecture, testing, observability, performance, dependencies, documentation, and operability.",
		Questions: []string{
			"Do static tools reveal failing tests, vet warnings, compile errors, or dependency problems?",
			"Which changed paths are high risk for authn/authz, token lifecycle, secrets, webhook verification, data writes, migrations, retries, concurrency, shell execution, network calls, deploys, or CI?",
			"Where does the patch need tests for failure modes, idempotency, rollback, retries, permissions, or compatibility?",
			"Where does understanding one concept require bouncing across many Modules?",
			"Where is the Interface nearly as complex as the Implementation?",
			"Apply the deletion test: if deleting the Module removes complexity, it may be shallow; if complexity spreads to callers, it earns its keep.",
			"Where do tests bypass the Interface or verify implementation details?",
			"Where does a seam have only one adapter and exist mostly as indirection?",
			"Which dependency category applies: in-process, local-substitutable, remote but owned, or true external?",
			"Does the change introduce performance, scalability, or resource-use risk?",
			"Would production observability through logs, errors, metrics, status, or traces identify the failure if this change regresses?",
		},
		Reject: []string{
			"Do not turn static hints directly into findings.",
			"Do not call something best-practice advice unless it is tied to code, tool output, or retrieved context.",
			"Do not report missing docs unless it blocks a concrete Module, Interface, operation, or review policy recommendation.",
			"Do not expose source URLs or source titles.",
			"Do not use component, service, API, or boundary when Module, Interface, or seam fits.",
		},
		Output: "Return only concrete recommendations with title, summary, benefit, recommendation, optional evidence, and strength.",
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
