package codereview

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// The reading half of a whole-repo review.
//
// Setting ReviewModeRepo only changes the label on a report. What makes a
// review of the repository different from a review of a diff is that the model
// is shown the repository: an inventory of what is in it, and the contents of
// its most review-worthy source files. Without these snippets `--repo` would
// hand the reviewer the same brief as a diff review and ask it a wider
// question, which is how a whole-repo review ends up answering "no material
// issues" about code it never saw.
//
// Everything here is language-agnostic on purpose. The module summaries and
// module_file snippets elsewhere in the brief are Go-only, so in a TypeScript,
// Python, or Ruby repository they contribute nothing and the diff would be the
// only code in the brief.
const (
	// maxWholeRepoInventoryPaths caps the file listing. It is a map of the
	// repository, not its contents, so it can be much wider than the snippets
	// that carry code, and it has its own budget on the way to the model.
	maxWholeRepoInventoryPaths = 4000
	// maxWholeRepoSourceFiles is the ceiling on how many source files a
	// whole-repo review reads.
	//
	// It was 12 (30 with --deep). Reviewing 12 of a repository's 273 source
	// files and reporting "no material issues found in the repository" is a
	// sentence about a repository that was not read. Coverage is the point of
	// `--repo`, so the number is now large enough to mean "all of it" for any
	// repository lgtm reviews, and the files are spread across shards rather than
	// crammed into one prompt. It survives only as a guard against a repository
	// far larger than that, and whatever it drops is counted and reported.
	maxWholeRepoSourceFiles = 5000
	// maxWholeRepoDeepSourceFiles is the --deep budget. Depth changes how hard
	// each file is looked at, not how many are looked at, so it matches.
	maxWholeRepoDeepSourceFiles = 5000
	// maxWholeRepoSourceFileBytes is the largest single file read whole. Files
	// above it are split across several snippets, in order, with every line
	// present — the alternative is truncation, which is sampling with extra
	// steps.
	maxWholeRepoSourceFileBytes = 1500000
	// maxWholeRepoToolScopeFiles caps the file list handed to static tools when
	// the whole repository is the subject.
	maxWholeRepoToolScopeFiles = 5000
)

// wholeRepoContextSnippets builds the repository-as-subject context: one
// inventory snippet plus the highest-ranked source files that are not already
// in the brief. It returns nil unless the caller asked for a whole-repo
// review, so every other review path — including the PR summary pipeline,
// which never sets WholeRepo — is byte-for-byte unaffected.
func wholeRepoContextSnippets(repoRoot string, facts RepoFacts, opts Options, changed []string, existing []ContextSnippet) []ContextSnippet {
	if !opts.WholeRepo {
		return nil
	}
	var out []ContextSnippet
	if snippet, ok := repoInventorySnippet(facts); ok {
		out = append(out, snippet)
	}
	// Skip anything the reviewer is already being shown: the diff carries the
	// changed files, and the retriever may have supplied docs or module files.
	skip := map[string]struct{}{}
	for _, snippet := range existing {
		if ref := strings.TrimSpace(snippet.Ref); ref != "" {
			skip[ref] = struct{}{}
		}
	}
	for _, file := range normalizedChangedFiles(changed) {
		skip[file] = struct{}{}
	}
	limit := maxWholeRepoSourceFiles
	if opts.Deep {
		limit = maxWholeRepoDeepSourceFiles
	}
	for _, rel := range rankedRepoSourceFiles(facts.Files, opts, skip, limit) {
		out = append(out, wholeRepoSourceSnippets(repoRoot, rel)...)
	}
	return out
}

// wholeRepoSourceSnippets reads one source file for a whole-repo review,
// splitting it across snippets when it exceeds the per-snippet budget.
//
// readSnippet truncates, which is the right call for a document being consulted
// and the wrong one for a file being reviewed: the bug is as likely to be on
// line 900 as on line 9. Splitting keeps every line, and the part header tells
// the model the ranges are consecutive so it does not read part 2 as a separate
// file.
func wholeRepoSourceSnippets(repoRoot, rel string) []ContextSnippet {
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
	if err != nil {
		return nil
	}
	if len(data) > maxWholeRepoSourceFileBytes {
		data = data[:maxWholeRepoSourceFileBytes]
	}
	lines := strings.Split(string(data), "\n")
	var parts []struct {
		text  string
		start int
		end   int
	}
	var current strings.Builder
	start := 1
	for i, line := range lines {
		if current.Len() > 0 && current.Len()+len(line)+1 > maxAIContextSnippetBytes {
			parts = append(parts, struct {
				text  string
				start int
				end   int
			}{text: current.String(), start: start, end: i})
			current.Reset()
			start = i + 1
		}
		current.WriteString(line)
		current.WriteString("\n")
	}
	if current.Len() > 0 {
		parts = append(parts, struct {
			text  string
			start int
			end   int
		}{text: current.String(), start: start, end: len(lines)})
	}
	if len(parts) == 0 {
		return nil
	}
	out := make([]ContextSnippet, 0, len(parts))
	for i, part := range parts {
		ref := rel
		text := part.text
		if len(parts) > 1 {
			ref = fmt.Sprintf("%s (part %d/%d)", rel, i+1, len(parts))
			text = fmt.Sprintf("%s lines %d-%d of %d:\n%s", rel, part.start, part.end, len(lines), part.text)
		}
		out = append(out, ContextSnippet{
			Kind:      "repo_source_file",
			Ref:       ref,
			Text:      text,
			Source:    "local",
			Publisher: "this repo",
			File:      rel,
			StartLine: part.start,
			EndLine:   part.end,
		})
	}
	return out
}

// repoInventorySnippet states what the repository contains. A reviewer that is
// told the repository is the subject needs to know its shape before it can say
// anything about it, and the listing is what makes a finding about a file that
// was not read at least askable rather than invented.
func repoInventorySnippet(facts RepoFacts) (ContextSnippet, bool) {
	if len(facts.Files) == 0 {
		return ContextSnippet{}, false
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Whole-repo review inventory: %d tracked files, %d test files.\n", facts.TrackedFileCount, facts.TestFileCount)
	if languages := repoLanguageCounts(facts.Files); languages != "" {
		fmt.Fprintf(&b, "File types: %s.\n", languages)
	}
	if len(facts.DependencyFiles) > 0 {
		fmt.Fprintf(&b, "Dependency manifests: %s.\n", strings.Join(facts.DependencyFiles, ", "))
	}
	files := facts.Files
	truncated := 0
	if len(files) > maxWholeRepoInventoryPaths {
		truncated = len(files) - maxWholeRepoInventoryPaths
		files = files[:maxWholeRepoInventoryPaths]
	}
	b.WriteString("Files:\n")
	for _, file := range files {
		b.WriteString(file)
		b.WriteString("\n")
	}
	if truncated > 0 {
		fmt.Fprintf(&b, "[%d more files not listed]\n", truncated)
	}
	return ContextSnippet{
		Kind:      "repo_inventory",
		Ref:       "repository inventory",
		Text:      b.String(),
		Source:    "local",
		Publisher: "this repo",
	}, true
}

func repoLanguageCounts(files []string) string {
	counts := map[string]int{}
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file))
		if ext == "" {
			continue
		}
		counts[ext]++
	}
	type extCount struct {
		ext   string
		count int
	}
	ranked := make([]extCount, 0, len(counts))
	for ext, count := range counts {
		ranked = append(ranked, extCount{ext: ext, count: count})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].count != ranked[j].count {
			return ranked[i].count > ranked[j].count
		}
		return ranked[i].ext < ranked[j].ext
	})
	if len(ranked) > 8 {
		ranked = ranked[:8]
	}
	parts := make([]string, 0, len(ranked))
	for _, entry := range ranked {
		parts = append(parts, fmt.Sprintf("%s %d", entry.ext, entry.count))
	}
	return strings.Join(parts, ", ")
}

// rankedRepoSourceFiles picks which of the repository's files are worth
// reading in full. A review prompt or a --focus path is the strongest signal
// available locally about what the caller cares about, so files whose path
// matches it come first: an agent asking "is there a race in the queue code?"
// should get the queue code, not the alphabetically first twelve files.
func rankedRepoSourceFiles(files []string, opts Options, skip map[string]struct{}, limit int) []string {
	if limit <= 0 {
		return nil
	}
	terms := repoQueryTerms(opts)
	type candidate struct {
		path   string
		score  int
		weight int
	}
	var candidates []string
	dirWeight := map[string]int{}
	for _, file := range files {
		if !isRepoSourceFile(file) {
			continue
		}
		dirWeight[filepath.ToSlash(filepath.Dir(file))]++
		if _, ok := skip[file]; ok {
			continue
		}
		candidates = append(candidates, file)
	}
	policy := reviewPolicyOf(opts)
	ranked := make([]candidate, 0, len(candidates))
	for _, file := range candidates {
		score := 0
		lower := strings.ToLower(file)
		for _, term := range terms {
			if strings.Contains(lower, term) {
				score += 3
			}
		}
		if !isTestFile(file) {
			score++
		}
		// The same path-impact signal the diff ranking and the PR summary use.
		// It decides which shard a file lands in rather than whether it is read
		// at all now, but reading the auth code in the first shard rather than
		// the seventh is still the difference between a review that leads with
		// what matters and one that gets there eventually.
		if impact := FileImpact(file, policy); impact.Matched() {
			score += 2
		}
		ranked = append(ranked, candidate{path: file, score: score, weight: dirWeight[filepath.ToSlash(filepath.Dir(file))]})
	}
	// Without a query the ranking still has to mean something. Directory
	// weight is the repository's own answer to "where is the code": a package
	// with twenty files is more of what this repository is than the
	// alphabetically first file in it.
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score != ranked[j].score {
			return ranked[i].score > ranked[j].score
		}
		if ranked[i].weight != ranked[j].weight {
			return ranked[i].weight > ranked[j].weight
		}
		return ranked[i].path < ranked[j].path
	})
	// There used to be a three-files-per-directory cap here, so that a
	// twelve-file budget saw a spread of the repository rather than one large
	// package. With every file read the cap has nothing left to protect, and a
	// cap on which files are read is exactly the sampling this exists to
	// remove. The ranking now only decides which shard a file lands in.
	var out []string
	for _, entry := range ranked {
		out = append(out, entry.path)
		if len(out) >= limit {
			break
		}
	}
	return out
}

// repoQueryTerms are the words worth matching paths against: what the caller
// asked and where they pointed.
func repoQueryTerms(opts Options) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, text := range []string{opts.Prompt, opts.Focus} {
		for _, term := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
			return !('a' <= r && r <= 'z') && !('0' <= r && r <= '9')
		}) {
			if len(term) < 3 || isRepoQueryStopword(term) {
				continue
			}
			if _, ok := seen[term]; ok {
				continue
			}
			seen[term] = struct{}{}
			out = append(out, term)
		}
	}
	return out
}

func isRepoQueryStopword(term string) bool {
	switch term {
	case "the", "and", "for", "are", "any", "does", "did", "how", "why", "what", "where", "when", "which", "who", "with", "from", "this", "that", "there", "here", "code", "review", "repo", "file", "files", "have", "has", "was", "were", "can", "should", "would", "could", "will", "our", "its", "into", "over", "out", "not", "all", "you", "your":
		return true
	default:
		return false
	}
}

// isRepoSourceFile keeps the selection to code. Docs, manifests, and lockfiles
// already have their own snippet kinds, and generated or vendored files teach a
// reviewer nothing about the repository its authors maintain.
func isRepoSourceFile(file string) bool {
	dir := filepath.ToSlash(filepath.Dir(file))
	if skipDir(dir) {
		return false
	}
	for _, part := range strings.Split(dir, "/") {
		// Dot directories are tool and build output — .build, .venv, .yarn —
		// and testdata and vendored trees are code this repository did not
		// write. A review of them is a review of somebody else's repository.
		if (strings.HasPrefix(part, ".") && part != ".") || part == "testdata" || part == "third_party" || part == "Pods" {
			return false
		}
	}
	lower := strings.ToLower(file)
	base := filepath.Base(lower)
	for _, marker := range []string{".pb.go", "_generated.go", ".gen.go", "_gen.go", ".generated.ts", ".min.js", ".d.ts"} {
		if strings.HasSuffix(base, marker) {
			return false
		}
	}
	switch filepath.Ext(lower) {
	case ".go", ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".py", ".rb", ".rs", ".java", ".kt", ".kts",
		".swift", ".m", ".mm", ".c", ".h", ".cc", ".cpp", ".hpp", ".cs", ".php", ".scala", ".ex", ".exs",
		".erl", ".hs", ".lua", ".pl", ".r", ".dart", ".sh", ".bash", ".zsh", ".sql", ".tf", ".proto":
		return true
	default:
		return false
	}
}

// repoWideStaticToolScope is what the static tools run over when a whole-repo
// review has no diff to scope them to. Without it a `--repo` review of a clean
// tree runs no tools at all — collectStaticToolResults returns early on an
// empty change set — and reports the repository clean without ever compiling
// it, which is exactly the silent pass --fail-on exists to prevent.
func repoWideStaticToolScope(facts RepoFacts) []string {
	var out []string
	for _, file := range facts.Files {
		if !isRepoSourceFile(file) && !isDependencyFile(file) {
			continue
		}
		out = append(out, file)
		if len(out) >= maxWholeRepoToolScopeFiles {
			break
		}
	}
	return out
}

// wholeRepoPromptLines are appended to the review instructions only when the
// brief's profile is whole_repo. They are additive for exactly that reason:
// the PR summary pipeline shares this prompt, so a whole-repo instruction that
// applied to every profile would change every summary lgtm writes on push.
func wholeRepoPromptLines() []string {
	return []string{
		"review_profile is whole_repo: the repository itself is the subject of this review, not the current patch. The caller asked about the codebase.",
		"In whole_repo reviews, repo_inventory lists what the repository contains and repo_source_file snippets carry its code. Treat both as primary evidence alongside static.diff_snippets.",
		"In whole_repo reviews, static.diff_snippets are the change in focus, not the boundary of the review: repo-wide findings grounded in the repository files you were shown are valid and expected, and the instruction to prefer changed lines over repo-wide advice does not apply here.",
		"In whole_repo reviews, when review_prompt is present, answer it from the repository files you were shown. If the provided context does not contain what the question is about, say so plainly in a recommendation instead of returning nothing.",
		"In whole_repo reviews at depth deep, apply the deep_full_spectrum checklist across the repository rather than only around the change in focus.",
	}
}

// wholeRepoReviewRubric is the rubric that matches the whole_repo profile. The
// scoped rubric it replaces asks about "the requested scope", which is the
// wrong question when the subject is the repository.
func wholeRepoReviewRubric(opts Options) ArchitectureRubric {
	scope := strings.TrimSpace(opts.Scope)
	if scope == "" {
		scope = DefaultScope
	}
	questions := []string{
		"What does repo_inventory say this repository is, and which Modules carry the most weight?",
		"Which repository code shown in repo_source_file snippets is insecure, brittle, undertested, or likely to break under realistic inputs, whether or not the current diff touches it?",
		"How does the change in focus in static.diff_snippets interact with the rest of the repository?",
		"Which findings are specific to the requested scope `" + scope + "` without ignoring security, dependencies, testing, and maintainability?",
		"Which recommendation has a clear first file to edit and verification command to run?",
	}
	if prompt := strings.TrimSpace(opts.Prompt); prompt != "" {
		questions = append([]string{"What does review_prompt ask about, and which of the repository files provided answer it?"}, questions...)
	}
	return ArchitectureRubric{
		Goal:      "Review the repository itself for concrete findings tied to the repository files provided, the current change in focus, tool output, local project policy, or retrieved review resources.",
		Questions: questions,
		Reject: []string{
			"Do not restrict findings to changed lines: repo-wide findings are the point of this review.",
			"Do not report a finding you cannot ground in something you were shown: a repository file, a diff snippet, tool output, or context.",
			"Do not turn file counts, missing docs, or missing tests directly into findings.",
			"Do not expose source URLs or source titles.",
		},
		Output: "Return only concrete repository-grounded recommendations with title, summary, benefit, recommendation, optional evidence, and strength.",
	}
}
