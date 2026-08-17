package codereview

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// constraintSignals are the deterministic facts every gate reads: diff volume,
// what looks generated, what dependencies arrived, which files are UI, and
// which paths make the performance gate apply. They are computed once from the
// change set and the diff snippets, before any command or model runs.
type constraintSignals struct {
	DiffStats          ConstraintsDiffStats
	ChangedLinesByFile map[string]int
	AddedLinesByFile   map[string][]constraintsAddedLine
	GeneratedFiles     []string
	WhitespaceDominant bool
	NewDependencies    []string
	ManifestsChanged   []string
	SourceFilesChanged []string
	TestFilesChanged   []string

	UIFiles      []string
	A11yFindings []Finding
	A11yHardFail bool
	PerfApplies  bool
	PerfTriggers []string
	PerfFiles    []string
	NonGenerated int // changed files that are not generated
	NonGenLines  int // changed lines outside generated files
}

// constraintsAddedLine is one "+" line of a diff with its post-change line
// number, so a finding about it can anchor to file:line.
type constraintsAddedLine struct {
	Number int
	Text   string
}

var constraintsUIExtensions = []string{".tsx", ".jsx", ".html", ".vue", ".svelte"}

// constraintsPerfPathPattern marks paths whose changes make the performance
// gate apply. Deliberately broad: a false "applies" costs one extra judgment in
// the shared model call, a false "skipped" costs the one gate Joe weighted up.
var constraintsPerfPathPattern = regexp.MustCompile(`(?i)(^|[/_.-])(perf|bench|hot|render|query|sql|cache|worker|handler|stream|middleware)([/_.-]|$|\.)`)

var constraintsPerfMessagePattern = regexp.MustCompile(`(?i)(perf|latenc|hot[- ]?path|throughput|load|slow)`)

// constraintsLineCounts is one file's exact added/removed counts from git
// numstat, untruncated. The diff snippets are capped at 32KB per file for the
// model's sake, so counting lines from them alone undercounts exactly the
// oversized files the size warning exists for.
type constraintsLineCounts struct {
	Added   int
	Removed int
}

func collectConstraintSignals(changed []string, diffs []DiffSnippet, policy ReviewPolicy, numstat map[string]constraintsLineCounts) constraintSignals {
	signals := constraintSignals{
		ChangedLinesByFile: map[string]int{},
		AddedLinesByFile:   map[string][]constraintsAddedLine{},
	}
	signals.DiffStats.Files = len(changed)
	for _, snippet := range diffs {
		added := constraintsAddedLinesForDiff(snippet.Diff)
		signals.AddedLinesByFile[snippet.File] = added
		addedCount, removedCount := countDiffLines(snippet.Diff)
		if strings.HasPrefix(snippet.Diff, diffUnavailableContentHeader) {
			// Content fallback: the whole file is the addition.
			addedCount, removedCount = len(added), 0
		}
		if counts, ok := numstat[snippet.File]; ok && counts.Added+counts.Removed > addedCount+removedCount {
			addedCount, removedCount = counts.Added, counts.Removed
		}
		signals.ChangedLinesByFile[snippet.File] = addedCount + removedCount
		signals.DiffStats.AddedLines += addedCount
		signals.DiffStats.RemovedLines += removedCount
	}
	for _, file := range changed {
		if isGeneratedReviewFile(file) {
			signals.GeneratedFiles = append(signals.GeneratedFiles, file)
		} else {
			signals.NonGenerated++
			signals.NonGenLines += signals.ChangedLinesByFile[file]
		}
		if isTestFile(file) {
			signals.TestFilesChanged = append(signals.TestFilesChanged, file)
		} else if isConstraintsSourceFile(file) {
			signals.SourceFilesChanged = append(signals.SourceFilesChanged, file)
		}
		if hasFileExtension(file, constraintsUIExtensions) {
			signals.UIFiles = append(signals.UIFiles, file)
		}
		if isDependencyFile(file) && !IsLockfilePath(file) {
			signals.ManifestsChanged = append(signals.ManifestsChanged, file)
		}
	}
	signals.DiffStats.GeneratedFiles = len(signals.GeneratedFiles)
	signals.WhitespaceDominant = whitespaceDominantDiff(diffs)
	signals.NewDependencies = constraintsNewDependencies(signals.ManifestsChanged, signals.AddedLinesByFile)
	signals.A11yFindings, signals.A11yHardFail = constraintsA11yLineFindings(signals.UIFiles, signals.AddedLinesByFile)
	signals.PerfApplies, signals.PerfTriggers, signals.PerfFiles = constraintsPerfApplicability(changed, signals.ManifestsChanged, policy)
	return signals
}

// constraintsNumstat reads exact per-file line counts from git. Empty refRange
// means the working tree (unstaged plus staged, summed); untracked files do
// not appear and fall back to the snippet-derived counts. Any git failure
// returns nil — the snippet counts are the graceful floor.
func constraintsNumstat(ctx context.Context, repoRoot, refRange string) map[string]constraintsLineCounts {
	argSets := [][]string{{"diff", "--numstat", "--no-ext-diff"}, {"diff", "--cached", "--numstat", "--no-ext-diff"}}
	if refRange = strings.TrimSpace(refRange); refRange != "" {
		argSets = [][]string{{"diff", "--numstat", "--no-ext-diff", refRange}}
	}
	out := map[string]constraintsLineCounts{}
	for _, args := range argSets {
		cmd := gitCommand(ctx, repoRoot, args...)
		output, err := cmd.Output()
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(output), "\n") {
			fields := strings.SplitN(strings.TrimSpace(line), "\t", 3)
			if len(fields) != 3 {
				continue
			}
			added, addErr := strconv.Atoi(fields[0])
			removed, remErr := strconv.Atoi(fields[1])
			if addErr != nil || remErr != nil {
				continue // binary files report "-"
			}
			file := filepath.ToSlash(strings.TrimSpace(fields[2]))
			counts := out[file]
			counts.Added += added
			counts.Removed += removed
			out[file] = counts
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// isConstraintsSourceFile reports whether a changed file is code rather than
// docs, config, or assets — the denominator for the "tests accompany sources"
// signal.
func isConstraintsSourceFile(rel string) bool {
	switch strings.ToLower(filepath.Ext(rel)) {
	case ".go", ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".py", ".rs", ".cs",
		".dart", ".java", ".rb", ".php", ".c", ".h", ".cc", ".cpp", ".hpp", ".swift",
		".kt", ".vue", ".svelte":
		return true
	default:
		return false
	}
}

// constraintsAddedLinesForDiff returns the added lines of one snippet. A
// content-fallback snippet (an untracked file — the commonest agent output)
// has no hunks: every line of it is new, and skipping it would exempt exactly
// the files most worth scanning.
func constraintsAddedLinesForDiff(diff string) []constraintsAddedLine {
	if strings.HasPrefix(diff, diffUnavailableContentHeader) {
		content := strings.TrimPrefix(strings.TrimPrefix(diff, diffUnavailableContentHeader), "\n")
		if strings.TrimSpace(content) == "" {
			return nil
		}
		lines := strings.Split(content, "\n")
		out := make([]constraintsAddedLine, 0, len(lines))
		for i, line := range lines {
			out = append(out, constraintsAddedLine{Number: i + 1, Text: line})
		}
		return out
	}
	return parseConstraintsAddedLines(diff)
}

// constraintsHunkContext is how many hunk lines surround the discussed line
// when a finding shows its diff window.
const constraintsHunkContext = 2

// constraintsDiffHunkForLine returns the unified-diff window around a
// post-change line, hunk header included, or "" when the line is not part of
// this diff — a finding about untouched code has no hunk to show.
func constraintsDiffHunkForLine(diff string, target int) string {
	if target <= 0 || strings.HasPrefix(diff, diffUnavailableContentHeader) {
		return ""
	}
	type diffHunk struct {
		header string
		lines  []string
		// nums holds each line's post-change line number; 0 for removed lines,
		// which exist only on the pre-change side.
		nums []int
	}
	var hunks []diffHunk
	newLine := 0
	inHunk := false
	for _, raw := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(raw, "@@"):
			newLine = parseHunkNewStart(raw)
			inHunk = newLine > 0
			if inHunk {
				hunks = append(hunks, diffHunk{header: raw})
			}
		case !inHunk:
		case strings.HasPrefix(raw, "-"):
			hunk := &hunks[len(hunks)-1]
			hunk.lines = append(hunk.lines, raw)
			hunk.nums = append(hunk.nums, 0)
		default:
			hunk := &hunks[len(hunks)-1]
			hunk.lines = append(hunk.lines, raw)
			hunk.nums = append(hunk.nums, newLine)
			newLine++
		}
	}
	for _, hunk := range hunks {
		targetIndex := -1
		for i, num := range hunk.nums {
			if num == target {
				targetIndex = i
				break
			}
		}
		if targetIndex < 0 {
			continue
		}
		start := targetIndex - constraintsHunkContext
		if start < 0 {
			start = 0
		}
		end := targetIndex + constraintsHunkContext + 1
		if end > len(hunk.lines) {
			end = len(hunk.lines)
		}
		out := []string{hunk.header}
		if start > 0 {
			out = append(out, "…")
		}
		out = append(out, hunk.lines[start:end]...)
		if end < len(hunk.lines) {
			out = append(out, "…")
		}
		return strings.Join(out, "\n")
	}
	return ""
}

// parseConstraintsAddedLines walks a unified diff and returns the added lines
// with their post-change line numbers.
func parseConstraintsAddedLines(diff string) []constraintsAddedLine {
	var out []constraintsAddedLine
	line := 0
	inHunk := false
	for _, raw := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(raw, "@@"):
			line = parseHunkNewStart(raw)
			inHunk = line > 0
			// The hunk header names the first line of the section; the counter
			// advances as lines are consumed below.
		case !inHunk:
			continue
		case strings.HasPrefix(raw, "+++"), strings.HasPrefix(raw, "---"):
			continue
		case strings.HasPrefix(raw, "+"):
			out = append(out, constraintsAddedLine{Number: line, Text: strings.TrimPrefix(raw, "+")})
			line++
		case strings.HasPrefix(raw, "-"):
			// Removed lines do not advance the post-change counter.
		default:
			line++
		}
	}
	return out
}

// parseHunkNewStart reads the "+c[,d]" of "@@ -a,b +c,d @@", or 0.
func parseHunkNewStart(header string) int {
	index := strings.Index(header, "+")
	if index < 0 {
		return 0
	}
	rest := header[index+1:]
	end := strings.IndexAny(rest, ", @")
	if end < 0 {
		end = len(rest)
	}
	value, err := strconv.Atoi(strings.TrimSpace(rest[:end]))
	if err != nil || value < 0 {
		return 0
	}
	if value == 0 {
		// "+0,0" is an empty new side; there is no line 0.
		return 0
	}
	return value
}

func countDiffLines(diff string) (added, removed int) {
	for _, raw := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(raw, "+++"), strings.HasPrefix(raw, "---"):
		case strings.HasPrefix(raw, "+"):
			added++
		case strings.HasPrefix(raw, "-"):
			removed++
		}
	}
	return added, removed
}

// Dependency-manifest parsing. These are heuristics feeding the back-pressure
// judgment ("a new dependency for a one-line change is bad work"), never a
// verdict on their own, so rough matching is acceptable.
var (
	constraintsGoModDepPattern   = regexp.MustCompile(`^\s*([A-Za-z0-9._~\-/]+)\s+v\d`)
	constraintsNodeDepPattern    = regexp.MustCompile(`"((?:@[A-Za-z0-9._-]+/)?[A-Za-z0-9._-]+)"\s*:\s*"[~^]?\d`)
	constraintsCargoDepPattern   = regexp.MustCompile(`^([A-Za-z0-9_-]+)\s*=\s*(?:"|\{)`)
	constraintsPipDepPattern     = regexp.MustCompile(`^([A-Za-z0-9._-]+)\s*(?:[=<>!~;\[]|$)`)
	constraintsPubspecDepPattern = regexp.MustCompile(`^\s{2}([a-z0-9_]+):`)
)

func constraintsNewDependencies(manifests []string, addedByFile map[string][]constraintsAddedLine) []string {
	seen := map[string]struct{}{}
	var out []string
	record := func(manifest, name string) {
		name = strings.TrimSpace(name)
		if name == "" {
			return
		}
		entry := filepath.Base(manifest) + ": " + name
		if _, ok := seen[entry]; ok {
			return
		}
		seen[entry] = struct{}{}
		out = append(out, entry)
	}
	for _, manifest := range manifests {
		base := filepath.Base(manifest)
		for _, added := range addedByFile[manifest] {
			text := strings.TrimSpace(added.Text)
			if text == "" || strings.HasPrefix(text, "#") || strings.HasPrefix(text, "//") {
				continue
			}
			switch {
			case base == "go.mod":
				if match := constraintsGoModDepPattern.FindStringSubmatch(text); match != nil && strings.Contains(match[1], "/") {
					record(manifest, match[1])
				}
			case base == "package.json":
				if match := constraintsNodeDepPattern.FindStringSubmatch(text); match != nil && match[1] != "version" && match[1] != "name" {
					record(manifest, match[1])
				}
			case base == "Cargo.toml":
				if match := constraintsCargoDepPattern.FindStringSubmatch(text); match != nil && match[1] != "version" && match[1] != "name" && match[1] != "edition" {
					record(manifest, match[1])
				}
			case base == "requirements.txt":
				if match := constraintsPipDepPattern.FindStringSubmatch(text); match != nil {
					record(manifest, match[1])
				}
			case base == "pubspec.yaml":
				if match := constraintsPubspecDepPattern.FindStringSubmatch(added.Text); match != nil && match[1] != "sdk" && match[1] != "flutter" {
					record(manifest, match[1])
				}
			case base == "pyproject.toml":
				if strings.Contains(text, "==") || strings.Contains(text, ">=") {
					if match := constraintsPipDepPattern.FindStringSubmatch(strings.Trim(text, `"',`)); match != nil {
						record(manifest, match[1])
					}
				}
			}
		}
	}
	sort.Strings(out)
	return out
}

// Accessibility line checks over the added lines of changed UI files. Only
// violations a single line can prove are hard failures; a clickable div might
// have its keyboard handler three lines down, so that one stays a finding for
// the model to confirm rather than a verdict.
var (
	constraintsPositiveTabindexPattern = regexp.MustCompile(`(?i)tabindex\s*[=:]\s*["'{]?\s*([1-9]\d*)`)
	constraintsImgTagPattern           = regexp.MustCompile(`(?i)<img\b[^>]*>`)
	constraintsAltAttrPattern          = regexp.MustCompile(`(?i)\balt\s*=`)
	constraintsClickableDivPattern     = regexp.MustCompile(`(?i)<(div|span)\b[^>]*onclick`)
	constraintsKeyboardHintPattern     = regexp.MustCompile(`(?i)(onkey|role\s*=|tabindex)`)
)

func constraintsA11yLineFindings(uiFiles []string, addedByFile map[string][]constraintsAddedLine) ([]Finding, bool) {
	var findings []Finding
	hardFail := false
	for _, file := range uiFiles {
		for _, added := range addedByFile[file] {
			text := added.Text
			for _, tag := range constraintsImgTagPattern.FindAllString(text, -1) {
				if !constraintsAltAttrPattern.MatchString(tag) {
					findings = append(findings, constraintsA11yFinding(file, added.Number,
						"Image without alt text",
						"An added <img> tag carries no alt attribute, so screen readers announce nothing for it.",
						"Add alt text describing the image, or alt=\"\" when it is decorative.",
						"Strong"))
					hardFail = true
				}
			}
			if match := constraintsPositiveTabindexPattern.FindStringSubmatch(text); match != nil {
				findings = append(findings, constraintsA11yFinding(file, added.Number,
					"Positive tabindex overrides focus order",
					fmt.Sprintf("An added element sets tabindex=%s, which hijacks the document's natural tab order for every keyboard user.", match[1]),
					"Use tabindex=\"0\" to join the natural order (or restructure the DOM so the order is right without it).",
					"Strong"))
				hardFail = true
			}
			if constraintsClickableDivPattern.MatchString(text) && !constraintsKeyboardHintPattern.MatchString(text) {
				findings = append(findings, constraintsA11yFinding(file, added.Number,
					"Click handler on a non-interactive element",
					"An added div/span handles onClick with no role, tabindex, or keyboard handler on the same line, which usually means keyboard users cannot reach it.",
					"Use a <button>, or add role, tabIndex={0}, and a keyboard handler.",
					"Worth exploring"))
			}
		}
	}
	return findings, hardFail
}

func constraintsA11yFinding(file string, line int, title, summary, recommendation, strength string) Finding {
	return Finding{
		ID:             "constraints.accessibility",
		Scopes:         []string{"maintainability"},
		Title:          title,
		Summary:        summary,
		Recommendation: recommendation,
		Strength:       strength,
		File:           file,
		Line:           line,
		Kind:           "defect",
	}
}

// constraintsPerfApplicability decides whether the performance gate applies:
// a REVIEW.md risk-path with a performance-flavored message, a path that names
// a hot-path concern, or any dependency-manifest change.
func constraintsPerfApplicability(changed []string, manifests []string, policy ReviewPolicy) (bool, []string, []string) {
	var triggers []string
	fileSet := map[string]struct{}{}
	for _, risk := range policy.RiskPaths {
		if !constraintsPerfMessagePattern.MatchString(risk.Message) {
			continue
		}
		for _, file := range changed {
			if MatchRiskPathGlob(risk.Glob, file) {
				triggers = append(triggers, fmt.Sprintf("REVIEW.md risk-path %s (%s)", risk.Glob, risk.Message))
				fileSet[file] = struct{}{}
				break
			}
		}
	}
	var keywordFiles []string
	for _, file := range changed {
		if isGeneratedReviewFile(file) {
			continue
		}
		if constraintsPerfPathPattern.MatchString(file) {
			keywordFiles = append(keywordFiles, file)
			fileSet[file] = struct{}{}
		}
	}
	if len(keywordFiles) > 0 {
		triggers = append(triggers, fmt.Sprintf("%d changed file(s) on perf-relevant paths", len(keywordFiles)))
	}
	if len(manifests) > 0 {
		triggers = append(triggers, "dependency manifest changed ("+strings.Join(baseNames(manifests), ", ")+")")
		for _, manifest := range manifests {
			fileSet[manifest] = struct{}{}
		}
	}
	files := make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	sort.Strings(files)
	return len(triggers) > 0, triggers, files
}

func baseNames(files []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, file := range files {
		base := filepath.Base(file)
		if _, ok := seen[base]; ok {
			continue
		}
		seen[base] = struct{}{}
		out = append(out, base)
	}
	return out
}

// constraintsCodeFiles filters the change to source-code files (tests
// included), for gate file attribution.
func constraintsCodeFiles(changed []string) []string {
	var out []string
	for _, file := range changed {
		if isConstraintsSourceFile(file) {
			out = append(out, file)
		}
	}
	return out
}
