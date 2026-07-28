package codereview

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// maxContextSnippetBytes and maxDiffSnippetBytes are the same numbers the AI
// layer enforces, referenced rather than restated. They were 3200 and 20000
// here against 1200 and 6000 there, so this file spent the work of reading and
// rendering bytes that compactReviewBriefForAI then discarded — a mismatch
// nothing could observe except by wondering why a review missed the second half
// of every file it claimed to have read.
const maxContextSnippetBytes = maxAIContextSnippetBytes
const maxDiffSnippetBytes = maxAIDiffSnippetBytes

// maxDiffSnippetFiles bounds how many changed files are diffed. It was 10 (24
// deep), applied to an alphabetically sorted list, which is how a review of a
// 181-file branch read files starting with "a" and reported the rest clean.
// Ranking (rankFilesByImpact) decides the order now and the shard planner packs
// everything that survives, so this is a guard against a pathological change
// set rather than the working limit — and whatever it drops is counted and
// reported, never silently discarded.
const maxDiffSnippetFiles = 400
const maxDeepDiffSnippetFiles = 1200

type ReviewBrief struct {
	RepoRoot      string             `json:"repo_root"`
	Scope         string             `json:"scope"`
	Depth         string             `json:"depth"`
	ReviewProfile string             `json:"review_profile"`
	Focus         string             `json:"focus,omitempty"`
	ReviewPrompt  string             `json:"review_prompt,omitempty"`
	Triage        ChangeTriage       `json:"triage,omitempty"`
	Static        StaticSnapshot     `json:"static"`
	Hints         []ReviewHint       `json:"hints"`
	Context       []ContextSnippet   `json:"context"`
	SourceRefs    []SourceRef        `json:"source_refs,omitempty"`
	SourceCatalog []SourceBrief      `json:"source_catalog"`
	Rubric        ArchitectureRubric `json:"rubric"`
	// Evidence records which retrieval sources answered and which did not. It
	// travels in the brief so the model is told what it is missing, and so the
	// report can say the same thing to the reader.
	Evidence []EvidenceStatus `json:"evidence,omitempty"`
}

// DegradedEvidence lists the evidence sources that failed this review.
func (b ReviewBrief) DegradedEvidence() []string {
	return EvidenceWarnings(b.Evidence)
}

type StaticSnapshot struct {
	FileCount       int             `json:"file_count"`
	TestFileCount   int             `json:"test_file_count"`
	DependencyFiles []string        `json:"dependency_files"`
	Docs            []FilePresence  `json:"docs"`
	ADRFiles        []string        `json:"adr_files,omitempty"`
	Modules         []ModuleSummary `json:"modules"`
	ChangedFiles    []string        `json:"changed_files"`
	// ReviewRange is the git ref range the changed files and diffs come from,
	// empty when the subject is the working tree.
	ReviewRange  string             `json:"review_range,omitempty"`
	DiffSnippets []DiffSnippet      `json:"diff_snippets,omitempty"`
	ToolResults  []StaticToolResult `json:"static_tool_results,omitempty"`
	CodeQuality  []CodeQualityHint  `json:"code_quality_hints,omitempty"`
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
	Publisher   string `json:"publisher,omitempty"`
	Title       string `json:"title,omitempty"`
	URL         string `json:"url,omitempty"`
	File        string `json:"file,omitempty"`
	StartLine   int    `json:"start_line,omitempty"`
	EndLine     int    `json:"end_line,omitempty"`
	Commit      string `json:"commit,omitempty"`
	SessionID   string `json:"session_id,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
	ResponseID  string `json:"response_id,omitempty"`
	ChunkHash   string `json:"chunk_hash,omitempty"`
}

type SourceBrief struct {
	ID        string   `json:"id"`
	Publisher string   `json:"publisher,omitempty"`
	Scopes    []string `json:"scopes"`
}

type SourceRef struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Title      string `json:"title,omitempty"`
	URL        string `json:"url,omitempty"`
	Source     string `json:"source,omitempty"`
	Publisher  string `json:"publisher,omitempty"`
	File       string `json:"file,omitempty"`
	StartLine  int    `json:"start_line,omitempty"`
	EndLine    int    `json:"end_line,omitempty"`
	Commit     string `json:"commit,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	ResponseID string `json:"response_id,omitempty"`
	ChunkHash  string `json:"chunk_hash,omitempty"`
}

type ArchitectureRubric struct {
	Goal      string   `json:"goal"`
	Questions []string `json:"questions"`
	Reject    []string `json:"reject"`
	Output    string   `json:"output"`
}

type ContextRetriever interface {
	Retrieve(ctx context.Context, in RetrieveInput) ([]ContextSnippet, error)
}

type LocalContextRetriever struct{}

func BuildReviewBrief(ctx context.Context, in RetrieveInput, sources []Source, retriever ContextRetriever) (ReviewBrief, error) {
	hints := in.Hints
	if hints == nil {
		hints = reviewHints(in.Facts)
		in.Hints = hints
	}
	opts := in.Options
	policy := opts.ReviewPolicy
	if policy == nil {
		loaded := LoadReviewPolicy(ctx, in.RepoRoot)
		policy = &loaded
		opts.ReviewPolicy = policy
		in.Options = opts
	}
	if in.Evidence == nil {
		in.Evidence = &EvidenceLog{}
	}
	contextSnippets, err := retriever.Retrieve(ctx, in)
	if err != nil {
		return ReviewBrief{}, err
	}
	contextSnippets = append(policy.ContextSnippets(), contextSnippets...)
	changed := normalizedChangedFiles(in.ChangedFiles)
	// A whole-repo review has to be given the repository, not just told that
	// the repository is the subject. Added here rather than inside a retriever
	// so it holds whichever retriever is configured.
	contextSnippets = append(contextSnippets, wholeRepoContextSnippets(in.RepoRoot, in.Facts, opts, changed, contextSnippets)...)
	contextSnippets = labelContextSnippets(contextSnippets)
	diffSnippets := in.DiffSnippets
	if diffSnippets == nil {
		diffSnippets = collectDiffSnippets(ctx, in.RepoRoot, changed, opts, in.DiffRange)
	}
	toolResults := []StaticToolResult(nil)
	if in.Plan.RunStaticTools || !reviewExecutionPlanConfigured(in.Plan) {
		toolResults = collectStaticToolResults(ctx, in.RepoRoot, in.Facts, opts, changed)
	}
	return ReviewBrief{
		RepoRoot:      in.RepoRoot,
		Scope:         opts.Scope,
		Depth:         depthLabel(opts.Deep),
		ReviewProfile: reviewProfile(opts),
		Focus:         strings.TrimSpace(opts.Focus),
		ReviewPrompt:  strings.TrimSpace(opts.Prompt),
		Triage:        in.Plan.Triage,
		Static: StaticSnapshot{
			FileCount:       in.Facts.TrackedFileCount,
			TestFileCount:   in.Facts.TestFileCount,
			DependencyFiles: in.Facts.DependencyFiles,
			Docs:            in.Facts.Docs,
			ADRFiles:        in.Facts.ADRFiles,
			Modules:         moduleSummaries(in.Facts),
			ChangedFiles:    changed,
			ReviewRange:     strings.TrimSpace(in.DiffRange),
			DiffSnippets:    diffSnippets,
			ToolResults:     toolResults,
			CodeQuality:     collectCodeQualityHints(in.RepoRoot, in.Facts, opts),
		},
		Hints:         hints,
		Context:       contextSnippets,
		SourceRefs:    sourceRefsFromContextSnippets(contextSnippets),
		SourceCatalog: sourceBriefs(sources),
		Rubric:        reviewRubric(opts),
		Evidence:      in.Evidence.Statuses(),
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

func (LocalContextRetriever) Retrieve(_ context.Context, in RetrieveInput) ([]ContextSnippet, error) {
	repoRoot := in.RepoRoot
	opts := in.Options
	facts := in.Facts
	hints := in.Hints
	var snippets []ContextSnippet
	for _, doc := range []struct {
		path string
		kind string
	}{
		{path: "CONTEXT.md", kind: "domain_doc"},
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
		// A declared manifest and its generated lockfile are not the same kind
		// of evidence and must not compete for the same context slot: go.mod
		// tells a reviewer what this project depends on, go.sum tells it
		// nothing it can act on. They were one kind, ranked above source code.
		kind := "dependency_manifest"
		if IsLockfilePath(rel) {
			kind = "dependency_lockfile"
		}
		if snippet, ok := readSnippet(repoRoot, rel, kind); ok {
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

// collectDiffSnippets renders the diff for each changed file. refRange selects
// the source: empty means the working tree (the interactive case), otherwise
// the files are diffed across that git ref range.
//
// Files are ordered by impact, not alphabetically. normalizedChangedFiles and
// rangeChangedFiles both sort by path, which is the right way to make a list
// reproducible and the wrong way to decide what a reviewer reads first: with a
// cap on the end of it, "sorted by path" means an auth change is reviewed if
// and only if it sorts early. rankFilesByImpact is the same ranking the PR
// summary uses to choose which hunks to describe.
func collectDiffSnippets(ctx context.Context, repoRoot string, files []string, opts Options, refRange string) []DiffSnippet {
	files = normalizedChangedFiles(files)
	files = rankFilesByImpact(files, reviewPolicyOf(opts))
	limit := maxDiffSnippetFiles
	if opts.Deep {
		limit = maxDeepDiffSnippetFiles
	}
	if len(files) > limit {
		files = files[:limit]
	}
	var snippets []DiffSnippet
	for _, file := range files {
		diff := diffForFile(ctx, repoRoot, refRange, file)
		if strings.TrimSpace(diff) == "" {
			diff = fileContentSnippet(repoRoot, file)
		}
		diff = truncateDiffText(diff, maxDiffSnippetBytes)
		if strings.TrimSpace(diff) == "" {
			continue
		}
		snippets = append(snippets, DiffSnippet{File: file, Diff: diff})
	}
	return snippets
}

func truncateDiffText(text string, limit int) string {
	return truncateAtHunkBoundary(text, limit)
}

// diffForFile picks the working-tree or ref-range diff for one file.
func diffForFile(ctx context.Context, repoRoot, refRange, file string) string {
	if strings.TrimSpace(refRange) == "" {
		return fileDiff(ctx, repoRoot, file)
	}
	return rangeFileDiff(ctx, repoRoot, refRange, file)
}

// rangeFileDiff diffs one file across a ref range, e.g. "main...HEAD".
func rangeFileDiff(ctx context.Context, repoRoot, refRange, file string) string {
	cmd := gitCommand(ctx, repoRoot, "diff", "--no-ext-diff", refRange, "--", file)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}

func fileDiff(ctx context.Context, repoRoot, file string) string {
	var parts []string
	for _, args := range [][]string{
		{"diff", "--no-ext-diff", "--", file},
		{"diff", "--cached", "--no-ext-diff", "--", file},
	} {
		cmd := gitCommand(ctx, repoRoot, args...)
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
	return ContextSnippet{Kind: kind, Ref: rel, Text: text, Source: "local", Publisher: "this repo"}, true
}

func LabelContextSnippets(snippets []ContextSnippet) []ContextSnippet {
	return labelContextSnippets(snippets)
}

func SourceRefsFromContextSnippets(snippets []ContextSnippet) []SourceRef {
	return sourceRefsFromContextSnippets(snippets)
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
		out = append(out, SourceBrief{ID: source.ID, Publisher: strings.TrimSpace(source.Publisher), Scopes: source.Scopes})
	}
	return out
}

func sourceRefsFromContextSnippets(snippets []ContextSnippet) []SourceRef {
	seen := map[string]struct{}{}
	var out []SourceRef
	for _, snippet := range snippets {
		ref := sourceRefFromContextSnippet(snippet)
		if ref.ID == "" {
			continue
		}
		if _, ok := seen[ref.ID]; ok {
			continue
		}
		seen[ref.ID] = struct{}{}
		out = append(out, ref)
	}
	return out
}

func sourceRefFromContextSnippet(snippet ContextSnippet) SourceRef {
	id := strings.TrimSpace(snippet.SourceLabel)
	if id == "" {
		id = strings.TrimSpace(snippet.Ref)
	}
	if id == "" {
		return SourceRef{}
	}
	return SourceRef{
		ID:         id,
		Kind:       sourceRefKind(snippet),
		Title:      sourceRefTitle(snippet),
		URL:        strings.TrimSpace(snippet.URL),
		Source:     strings.TrimSpace(snippet.Source),
		Publisher:  publisherForContextSnippet(snippet),
		File:       firstNonEmpty(snippet.File, snippet.Ref),
		StartLine:  snippet.StartLine,
		EndLine:    snippet.EndLine,
		Commit:     strings.TrimSpace(snippet.Commit),
		SessionID:  strings.TrimSpace(snippet.SessionID),
		RequestID:  strings.TrimSpace(snippet.RequestID),
		ResponseID: strings.TrimSpace(snippet.ResponseID),
		ChunkHash:  strings.TrimSpace(snippet.ChunkHash),
	}
}

func publisherForContextSnippet(snippet ContextSnippet) string {
	if publisher := strings.TrimSpace(snippet.Publisher); publisher != "" {
		return publisher
	}
	switch strings.TrimSpace(snippet.Kind) {
	case "domain_doc", "repo_doc", "adr", "dependency_manifest", "changed_file", "module_file", "code_quality_file", "repo_inventory", "repo_source_file":
		return "local"
	case "indexed_session", "session":
		return "session"
	case "indexed_code":
		return "indexed"
	case "review_policy", "review_reference":
		if publisher := publisherFromURLHost(snippet.URL); publisher != "" {
			return publisher
		}
		return "local"
	case "review_resource":
		if publisher := publisherFromURLHost(snippet.URL); publisher != "" {
			return publisher
		}
		return "unknown"
	default:
		if publisher := publisherFromURLHost(snippet.URL); publisher != "" {
			return publisher
		}
		if strings.EqualFold(strings.TrimSpace(snippet.Source), "local") {
			return "local"
		}
		if strings.HasPrefix(strings.TrimSpace(snippet.Source), "turbopuffer:") {
			return "indexed"
		}
		return "unknown"
	}
}

func sourceRefKind(snippet ContextSnippet) string {
	switch strings.TrimSpace(snippet.Kind) {
	case "indexed_code":
		return "code"
	case "indexed_session":
		return "session"
	case "review_resource":
		return "resource"
	case "review_policy":
		return "policy"
	case "review_reference":
		return "reference"
	case "domain_doc", "repo_doc", "adr", "dependency_manifest", "dependency_lockfile", "repo_inventory", "repo_source_file":
		return "local"
	default:
		return firstNonEmpty(snippet.Kind, "context")
	}
}

func sourceRefTitle(snippet ContextSnippet) string {
	if title := strings.TrimSpace(snippet.Title); title != "" {
		return title
	}
	if snippet.File != "" && snippet.StartLine > 0 {
		return fmt.Sprintf("%s:%d", snippet.File, snippet.StartLine)
	}
	if snippet.Ref != "" {
		return snippet.Ref
	}
	return snippet.Kind
}

func reviewRubric(opts Options) ArchitectureRubric {
	// WholeRepo comes first because it names the subject: depth, scope, and
	// prompt all describe how to review, and the rubric has to ask about the
	// repository before it asks how deeply to look at it.
	if opts.WholeRepo {
		return wholeRepoReviewRubric(opts)
	}
	if opts.Deep {
		return deepReviewRubric()
	}
	if opts.PatchFocused {
		return patchReviewRubric()
	}
	if strings.TrimSpace(opts.Prompt) != "" {
		return promptDirectedReviewRubric(opts)
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

func promptDirectedReviewRubric(opts Options) ArchitectureRubric {
	scope := strings.TrimSpace(opts.Scope)
	if scope == "" {
		scope = DefaultScope
	}
	return ArchitectureRubric{
		Goal: "Review the user's review_prompt against the current changes and the surrounding codebase, producing concrete findings where the prompt, diff, and repo context intersect.",
		Questions: []string{
			"What concern or behavior is the user asking about in review_prompt?",
			"How do static.diff_snippets and static.changed_files affect that concern?",
			"Which surrounding Modules, Interfaces, tests, docs, or local policies make the prompted concern safer or riskier?",
			"Which finding is specific to the requested scope `" + scope + "` while still answering review_prompt?",
			"Which recommendation has a clear first file to edit and verification command to run?",
		},
		Reject: []string{
			"Do not ignore the current diff; use it as evidence for why the prompted concern matters now.",
			"Do not limit the review to changed lines when surrounding code explains the risk or correct fix.",
			"Do not report generic repo-wide advice that does not answer review_prompt.",
			"Do not expose source URLs or source titles.",
			"If there is no concrete finding tied to review_prompt and repo evidence, return no recommendations.",
		},
		Output: "Return only concrete prompt-directed recommendations with title, summary, benefit, recommendation, optional evidence, and strength.",
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
