package codereview

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// fileTokenPattern matches file-ish tokens (a path segment ending in an
// extension) so relevance matching compares whole filenames, not substrings.
var fileTokenPattern = regexp.MustCompile(`[A-Za-z0-9_.\-/]+\.[A-Za-z0-9]+`)

type Finding struct {
	ID       string          `json:"id"`
	Scopes   []string        `json:"scopes,omitempty"`
	Title    string          `json:"title"`
	Summary  string          `json:"summary,omitempty"`
	Benefit  string          `json:"benefit,omitempty"`
	Evidence []Evidence      `json:"evidence,omitempty"`
	Anchors  []FindingAnchor `json:"anchors,omitempty"`
	File     string          `json:"file,omitempty"`
	Line     int             `json:"line,omitempty"`

	// Kind is the reviewer's own classification of what the finding asks of
	// the reader: "defect" (the code does something wrong now), "hardening"
	// (correct today, fragile tomorrow), or "suggestion" (nothing is wrong).
	// Empty means the reviewer did not say. Measured on the review benchmark:
	// findings claiming a present defect matched a human-recorded issue five
	// times as often as hardening or suggestion findings, so this label is the
	// strongest single noise separator the review produces.
	Kind string `json:"kind,omitempty"`

	// JudgeVerdict and the fields after it carry the verification model's
	// assessment through to the report: "confirmed", "unverified" (the judge
	// abstained or never ran), or empty for findings that predate the judge.
	// The values were previously computed, used to order the report, and then
	// discarded — leaving a reader (or any downstream ranking) no way to tell
	// a finding the judge verified at 0.95 from one it waved through.
	JudgeVerdict    string  `json:"judge_verdict,omitempty"`
	JudgeImpact     string  `json:"judge_impact,omitempty"`
	JudgeSeverity   int     `json:"judge_severity,omitempty"`
	JudgeConfidence float64 `json:"judge_confidence,omitempty"`
	// JudgeRank is the judge's batch-relative reading order (1 = read first),
	// the ordering key for the report. Absolute judge scores cluster (measured:
	// confidence mass sits in 0.8-0.9), so the ordering comes from the one
	// model that read every candidate side by side. 0 means unranked.
	JudgeRank int `json:"judge_rank,omitempty"`

	// Corroboration names the independent reviewer legs that each raised this
	// finding on their own. Two flagship models converging on the same problem
	// is the strongest quality signal a multi-model panel produces, and it used
	// to be thrown away: near-duplicates were merged (when they were merged at
	// all) into an evidence footnote nobody sees without --verbose. It is a
	// first-class field so the renderer, the JSON report, the ranking, and the
	// recorded review history can all tell a corroborated finding from a
	// single-reviewer one.
	//
	// One entry means one reviewer raised it. len > 1 means genuine
	// cross-reviewer agreement.
	Corroboration []string `json:"corroboration,omitempty"`

	// MergedFindings are the other reviewers' versions of this same finding,
	// kept verbatim. Merging picks one copy for the reader to meet, and the
	// copies it does not pick often carry the more concrete fix — one leg wrote
	// "restrict egress" where the other wrote "call krun_set_port_map with an
	// explicit allowlist and add a test asserting an outbound connection to an
	// unlisted host fails". Dropping that was de-duplication deleting the most
	// actionable sentence in the review, so it is recorded here, rendered under
	// the survivor, and serialised in the JSON report.
	MergedFindings []MergedFinding `json:"merged_findings,omitempty"`

	ResolvedSources []ResolvedSource `json:"resolved_sources,omitempty"`
	Recommendation  string           `json:"recommendation,omitempty"`
	// Strength is the severity vocabulary gates read: "Blocking", "Strong",
	// "Worth exploring", or "Speculative".
	Strength         string   `json:"strength,omitempty"`
	SourceIDs        []string `json:"source_ids,omitempty"`
	SourcePublishers []string `json:"source_publishers,omitempty"`
}

type FindingAnchor struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

// MergedFinding is one reviewer's version of a finding that de-duplication
// folded into another. It exists so a merge never costs the reader a fix.
type MergedFinding struct {
	ID             string `json:"id,omitempty"`
	Reviewer       string `json:"reviewer,omitempty"`
	Title          string `json:"title,omitempty"`
	Summary        string `json:"summary,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

type Evidence struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Rule interface {
	ID() string
	Scopes() []string
	Evaluate(ReviewContext) []Finding
}

type ruleFunc struct {
	id       string
	scopes   []string
	evaluate func(ReviewContext) []Finding
}

func (r ruleFunc) ID() string { return r.id }

func (r ruleFunc) Scopes() []string { return r.scopes }

func (r ruleFunc) Evaluate(ctx ReviewContext) []Finding {
	if r.evaluate == nil {
		return nil
	}
	return r.evaluate(ctx)
}

func defaultRules() []Rule {
	return []Rule{
		ruleFunc{id: "tools.static-failure", scopes: []string{"testing", "maintainability", "dependencies", "security"}, evaluate: staticToolFailureFindings},
		ruleFunc{id: "quality.ignored-results", scopes: []string{"maintainability", "testing", "architecture"}, evaluate: ignoredResultFindings},
		ruleFunc{id: "security.quality-hints", scopes: []string{"security", "maintainability", "testing"}, evaluate: securityQualityHintFindings},
		ruleFunc{id: "architecture.domain-language-drift", scopes: []string{"architecture", "maintainability", "docs"}, evaluate: domainLanguageDriftFindings},
		ruleFunc{id: "architecture.implementation-heavy-module", scopes: []string{"architecture", "maintainability", "testing"}, evaluate: implementationHeavyModuleFindings},
		ruleFunc{id: "architecture.generic-package-name", scopes: []string{"architecture", "maintainability"}, evaluate: genericPackageNameFindings},
		ruleFunc{id: "architecture.package-without-tests", scopes: []string{"architecture", "testing", "maintainability"}, evaluate: packageWithoutTestsFindings},
		ruleFunc{id: "onboarding.missing-agents", scopes: []string{"onboarding", "maintainability"}, evaluate: oneFinding(missingAgentsFinding)},
		ruleFunc{id: "docs.missing-readme", scopes: []string{"docs", "onboarding", "maintainability"}, evaluate: oneFinding(missingReadmeFinding)},
		ruleFunc{id: "testing.no-tests", scopes: []string{"testing", "maintainability"}, evaluate: oneFinding(noTestsFinding)},
		ruleFunc{id: "dependencies.no-manifest", scopes: []string{"dependencies", "maintainability"}, evaluate: oneFinding(noDependencyManifestFinding)},
		ruleFunc{id: "dependencies.javascript-unlocked", scopes: []string{"dependencies", "security", "maintainability"}, evaluate: oneFinding(jsDependencyLockFinding)},
	}
}

func domainLanguageDriftFindings(ctx ReviewContext) []Finding {
	terms := internalTermsFromContext(ctx.Brief.Context)
	if len(terms) == 0 || strings.TrimSpace(ctx.Brief.RepoRoot) == "" {
		return nil
	}
	var evidence []Evidence
	for _, file := range ctx.Facts.Files {
		if !publicSurfaceFile(file) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(ctx.Brief.RepoRoot, filepath.FromSlash(file)))
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		for index, line := range lines {
			searchText := publicSearchText(file, line)
			if searchText == "" {
				continue
			}
			lower := strings.ToLower(searchText)
			for _, term := range terms {
				if !strings.Contains(lower, strings.ToLower(term)) {
					continue
				}
				evidence = append(evidence, Evidence{
					Label: fmt.Sprintf("%s:%d", file, index+1),
					Value: strings.TrimSpace(line),
				})
				if len(evidence) >= 8 {
					return []Finding{domainLanguageDriftFinding(terms, evidence)}
				}
				break
			}
		}
	}
	if len(evidence) == 0 {
		return nil
	}
	return []Finding{domainLanguageDriftFinding(terms, evidence)}
}

func domainLanguageDriftFinding(terms []string, evidence []Evidence) Finding {
	return Finding{
		ID:      "architecture.domain-language-drift",
		Scopes:  []string{"architecture", "maintainability", "docs"},
		Title:   "Keep internal terms out of public seams",
		Summary: fmt.Sprintf("`CONTEXT.md` marks %s as internal implementation language, but public-facing files still expose that vocabulary.", formatTermList(terms, 3)),
		Benefit: "Improves review and agent reliability by keeping public Interfaces aligned with the product nouns users and agents should learn.",
		Evidence: append(evidence, Evidence{
			Label: "Decision rule",
			Value: "Use product terms at public seams; keep internal terms in implementation code, ADRs, and migration notes.",
		}),
		Recommendation: "Rename public docs, CLI help, and MCP descriptions to use the canonical product terms from `CONTEXT.md`. Leave internal helper names alone until the public seam is clean.",
		Strength:       "Strong",
		SourceIDs:      []string{"go-package-names", "google-eng-practices"},
	}
}

func internalTermsFromContext(snippets []ContextSnippet) []string {
	seen := map[string]struct{}{}
	var terms []string
	for _, snippet := range snippets {
		if snippet.Ref != "CONTEXT.md" || snippet.Kind != "domain_doc" {
			continue
		}
		inInternalSection := false
		for _, raw := range strings.Split(snippet.Text, "\n") {
			line := strings.TrimSpace(raw)
			if strings.HasPrefix(line, "## ") {
				inInternalSection = strings.Contains(strings.ToLower(line), "internal") &&
					strings.Contains(strings.ToLower(line), "term")
				continue
			}
			if !inInternalSection || !strings.HasPrefix(line, "### ") {
				continue
			}
			term := strings.TrimSpace(strings.TrimPrefix(line, "### "))
			term = strings.Trim(term, "`*_ ")
			if term == "" {
				continue
			}
			key := strings.ToLower(term)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			terms = append(terms, term)
		}
	}
	return terms
}

func publicSurfaceFile(file string) bool {
	if isTestFile(file) {
		return false
	}
	switch file {
	case "README.md", "AGENTS.md", "mcp/README.md":
		return true
	}
	if strings.HasPrefix(file, "skills/") && strings.HasSuffix(file, ".md") {
		return true
	}
	if strings.HasPrefix(file, "mcp/src/tools/") && strings.HasSuffix(file, ".ts") {
		return true
	}
	return false
}

func publicSearchText(file string, line string) string {
	if strings.HasPrefix(file, "mcp/src/tools/") {
		if !strings.Contains(line, ".describe(") &&
			!strings.Contains(line, "description:") &&
			!strings.Contains(line, "nextActions") &&
			!strings.Contains(line, "return ") {
			return ""
		}
	}
	if strings.HasSuffix(file, ".go") || strings.HasSuffix(file, ".ts") {
		return quotedLineText(line)
	}
	return line
}

func quotedLineText(line string) string {
	var out strings.Builder
	for i := 0; i < len(line); i++ {
		quote := line[i]
		if quote != '"' && quote != '`' {
			continue
		}
		for j := i + 1; j < len(line); j++ {
			if line[j] == quote {
				if out.Len() > 0 {
					out.WriteByte(' ')
				}
				out.WriteString(line[i+1 : j])
				i = j
				break
			}
		}
	}
	return out.String()
}

func formatTermList(terms []string, limit int) string {
	if len(terms) == 0 {
		return "internal terms"
	}
	if limit <= 0 || len(terms) <= limit {
		return "`" + strings.Join(terms, "`, `") + "`"
	}
	return fmt.Sprintf("`%s` and %d more", strings.Join(terms[:limit], "`, `"), len(terms)-limit)
}

func ignoredResultFindings(ctx ReviewContext) []Finding {
	var ignored []CodeQualityHint
	for _, hint := range ctx.Brief.Static.CodeQuality {
		if hint.Kind == "ignored_result" {
			ignored = append(ignored, hint)
		}
	}
	if len(ignored) == 0 {
		return nil
	}
	limit := 5
	if len(ignored) < limit {
		limit = len(ignored)
	}
	var evidence []Evidence
	for _, hint := range ignored[:limit] {
		evidence = append(evidence, Evidence{
			Label: fmt.Sprintf("%s:%d", hint.File, hint.Line),
			Value: hint.Text,
		})
	}
	first := ignored[0]
	return []Finding{{
		ID:      "quality.ignored-results",
		Scopes:  []string{"maintainability", "testing", "architecture"},
		Title:   "Review ignored errors in production code",
		Summary: fmt.Sprintf("`%s:%d` discards a result in production code; these are often intentional, but they can also hide failed recovery, parsing, or state transitions.", first.File, first.Line),
		Benefit: "Improves error handling and observability by turning silent authoring failures into explicit handled states or documented best-effort paths.",
		Evidence: append(evidence, Evidence{
			Label: "Why this matters",
			Value: "A discarded result should either be proven harmless, handled, or wrapped behind a Module Interface that makes the failure mode explicit.",
		}),
		Recommendation: fmt.Sprintf("Start with `%s:%d`, then the next four ignored results listed in evidence. For each one, either return the error, log the best-effort failure with context, or add a small test proving the ignored result is harmless.", first.File, first.Line),
		Strength:       "Worth exploring",
		SourceIDs:      []string{"google-eng-practices"},
	}}
}

func securityQualityHintFindings(ctx ReviewContext) []Finding {
	changed := changedFileSet(ctx.Brief.Static.ChangedFiles)
	if len(changed) == 0 {
		return nil
	}
	var matched []CodeQualityHint
	for _, hint := range ctx.Brief.Static.CodeQuality {
		if !isSecurityQualityHintKind(hint.Kind) {
			continue
		}
		file := filepath.ToSlash(strings.TrimSpace(hint.File))
		if _, ok := changed[file]; !ok {
			continue
		}
		matched = append(matched, hint)
	}
	if len(matched) == 0 {
		return nil
	}
	limit := 3
	if len(matched) < limit {
		limit = len(matched)
	}
	var out []Finding
	for index, hint := range matched[:limit] {
		out = append(out, securityQualityHintFinding(hint, index+1))
	}
	return out
}

func isSecurityQualityHintKind(kind string) bool {
	switch kind {
	case "sql_interpolation", "shell_injection", "unsafe_regex", "hardcoded_secret":
		return true
	default:
		return false
	}
}

func securityQualityHintFinding(hint CodeQualityHint, index int) Finding {
	title, recommendation := securityQualityHintCopy(hint.Kind)
	summary := strings.TrimSpace(hint.Reason)
	if summary == "" {
		summary = fmt.Sprintf("`%s:%d` matches a security-sensitive code pattern.", hint.File, hint.Line)
	}
	return Finding{
		ID:      fmt.Sprintf("security.quality-hints.%d", index),
		Scopes:  []string{"security", "maintainability", "testing"},
		Title:   title,
		Summary: summary,
		Benefit: "Reduces exploitability by catching common injection, secret exposure, and unsafe input handling in changed code without waiting for AI review.",
		Evidence: []Evidence{{
			Label: fmt.Sprintf("%s:%d", hint.File, hint.Line),
			Value: hint.Text,
		}},
		Anchors:        []FindingAnchor{{File: hint.File, Line: hint.Line}},
		Recommendation: recommendation,
		Strength:       "Strong",
		SourceIDs:      []string{"owasp-secure-coding"},
	}
}

func securityQualityHintCopy(kind string) (title, recommendation string) {
	switch kind {
	case "sql_interpolation":
		return "Avoid SQL string interpolation", "Replace interpolated SQL with parameterized queries or prepared statements, then add a test that proves user input cannot alter query structure."
	case "shell_injection":
		return "Avoid shell command injection", "Stop passing user input through a shell; use argv-based execution with an allowlisted command and arguments, then test malicious input cases."
	case "unsafe_regex":
		return "Validate regex patterns from user input", "Do not compile user-controlled patterns directly; validate length and syntax, use a safe matcher, or reject dangerous patterns before calling regexp.Compile."
	case "hardcoded_secret":
		return "Remove hardcoded secrets from source", "Move credentials to environment variables or a secret manager, rotate any exposed values, and add a scanner or test that fails on committed secrets."
	default:
		return "Review security-sensitive code in changed files", "Inspect the flagged line and replace the risky pattern with a safer alternative before merging."
	}
}

// staticToolFailureFindings turns failed tool runs into the one deterministic
// Blocking finding. Skipped results never qualify — that is where timeouts and
// environment failures land (the tool never ran because dependencies are not
// installed, or failed for a reason only the host can cause; see
// staticToolEnvironmentFailure) — so a bare clone does not get a published
// "fix your tools" comment about its own missing node_modules.
func staticToolFailureFindings(ctx ReviewContext) []Finding {
	var failed []StaticToolResult
	for _, result := range ctx.Brief.Static.ToolResults {
		if result.Skipped || result.ExitCode == 0 {
			continue
		}
		failed = append(failed, result)
	}
	if len(failed) == 0 {
		return nil
	}
	var evidence []Evidence
	for _, result := range failed {
		value := strings.TrimSpace(result.Output)
		if value == "" {
			value = fmt.Sprintf("%s exited with code %d", result.Command, result.ExitCode)
		}
		evidence = append(evidence, Evidence{Label: result.Name, Value: value})
	}
	first := failed[0]
	return []Finding{{
		ID:             "tools.static-failure",
		Scopes:         []string{"testing", "maintainability", "dependencies", "security"},
		Title:          "Fix static tool failures before architecture advice",
		Summary:        fmt.Sprintf("`%s` failed, so the review found concrete correctness or best-practice diagnostics before speculative structure work.", first.Command),
		Benefit:        "Restores a clean correctness baseline so later architecture recommendations are judged against working code instead of compile, vet, or test failures.",
		Evidence:       evidence,
		Recommendation: "Start by fixing the failing tool output, then rerun `gx review` so the reviewer can judge structure on a clean baseline.",
		Strength:       "Blocking",
		SourceIDs:      []string{"google-eng-practices"},
	}}
}

func oneFinding(fn func(RepoFacts) Finding) func(ReviewContext) []Finding {
	return func(ctx ReviewContext) []Finding {
		finding := fn(ctx.Facts)
		if finding.ID == "" {
			return nil
		}
		return []Finding{finding}
	}
}

func evaluateFindings(ctx ReviewContext, rules []Rule) []Finding {
	active := activeScopes(ctx.ActiveScopes)
	knownSources := sourceIDSet(ctx.Sources)
	var out []Finding
	for _, rule := range rules {
		if rule == nil || !matchesScope(rule.Scopes(), active) {
			continue
		}
		for _, finding := range rule.Evaluate(ctx) {
			if finding.ID == "" || !matchesScope(finding.Scopes, active) {
				continue
			}
			finding.SourceIDs = filterSourceIDs(finding.SourceIDs, knownSources)
			out = append(out, finding)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if strengthRank(out[i].Strength) != strengthRank(out[j].Strength) {
			return strengthRank(out[i].Strength) < strengthRank(out[j].Strength)
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func filterPatchFocusedFindings(ctx ReviewContext, findings []Finding) []Finding {
	if !ctx.Options.PatchFocused || len(findings) == 0 {
		return findings
	}
	changed := normalizedChangedFiles(ctx.Brief.Static.ChangedFiles)
	var out []Finding
	for _, finding := range findings {
		if strings.HasPrefix(finding.ID, "tools.") {
			out = append(out, finding)
			continue
		}
		if len(changed) == 0 {
			continue
		}
		if findingMentionsChangedFile(finding, changed) {
			out = append(out, finding)
		}
	}
	return out
}

func normalizedChangedFiles(files []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, file := range files {
		file = strings.TrimSpace(filepath.ToSlash(file))
		if file == "" {
			continue
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		out = append(out, file)
	}
	return out
}

func findingMentionsChangedFile(finding Finding, changed []string) bool {
	text := strings.Join([]string{
		finding.Title,
		finding.Summary,
		finding.Benefit,
		finding.Recommendation,
		evidenceText(finding.Evidence),
		anchorsText(finding.Anchors),
		fileLineText(finding.File, finding.Line),
	}, "\n")

	changedFull := make(map[string]struct{}, len(changed))
	changedBase := make(map[string]struct{}, len(changed))
	for _, file := range changed {
		file = strings.TrimSpace(filepath.ToSlash(file))
		if file == "" {
			continue
		}
		changedFull[file] = struct{}{}
		changedBase[path.Base(file)] = struct{}{}
	}

	// Compare whole file tokens rather than substrings: an exact path match, or
	// a bare filename (no directory) that matches a changed file's basename.
	// This avoids "utils.go" matching "myutils.go" and matches findings that
	// reference a file by name only.
	for _, token := range fileTokenPattern.FindAllString(filepath.ToSlash(text), -1) {
		token = strings.Trim(token, "./")
		if token == "" {
			continue
		}
		if _, ok := changedFull[token]; ok {
			return true
		}
		if token == path.Base(token) {
			if _, ok := changedBase[token]; ok {
				return true
			}
		}
	}
	return false
}

func fileLineText(file string, line int) string {
	file = strings.TrimSpace(filepath.ToSlash(file))
	if file == "" || line <= 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func anchorsText(anchors []FindingAnchor) string {
	var b strings.Builder
	for _, anchor := range anchors {
		if strings.TrimSpace(anchor.File) == "" || anchor.Line <= 0 {
			continue
		}
		fmt.Fprintf(&b, "%s:%d\n", filepath.ToSlash(strings.TrimSpace(anchor.File)), anchor.Line)
	}
	return b.String()
}

func evidenceText(evidence []Evidence) string {
	var b strings.Builder
	for _, item := range evidence {
		if item.Label != "" {
			b.WriteString(item.Label)
			b.WriteByte('\n')
		}
		if item.Value != "" {
			b.WriteString(item.Value)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func activeScopes(scopes []string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, baseline := range scopes {
		out[baseline] = struct{}{}
	}
	return out
}

func matchesScope(scopes []string, active map[string]struct{}) bool {
	for _, scope := range scopes {
		if _, ok := active[scope]; ok {
			return true
		}
	}
	return false
}

func sourceIDSet(sources []Source) map[string]struct{} {
	out := map[string]struct{}{}
	for _, source := range sources {
		out[source.ID] = struct{}{}
	}
	return out
}

func filterSourceIDs(ids []string, known map[string]struct{}) []string {
	var out []string
	for _, id := range ids {
		if _, ok := known[id]; ok {
			out = append(out, id)
		}
	}
	return out
}

func strengthRank(strength string) int {
	switch strength {
	case "Blocking":
		return 0
	case "Strong":
		return 1
	case "Worth exploring":
		return 2
	default:
		return 3
	}
}

func mergeFindings(local []Finding, ai []Finding) []Finding {
	seen := map[string]struct{}{}
	out := make([]Finding, 0, len(local)+len(ai))
	for _, finding := range append(append([]Finding{}, local...), ai...) {
		id := strings.TrimSpace(finding.ID)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, finding)
	}
	sort.SliceStable(out, func(i, j int) bool {
		left := strengthRank(out[i].Strength)
		right := strengthRank(out[j].Strength)
		if left != right {
			return left < right
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func missingAgentsFinding(facts RepoFacts) Finding {
	if present(facts.Docs, "AGENTS.md") {
		return Finding{}
	}
	return Finding{
		ID:      "onboarding.missing-agents",
		Scopes:  []string{"onboarding", "maintainability"},
		Title:   "Missing agent instructions",
		Summary: "Agents do not have a repo-local instruction file for command, review, and version-control conventions.",
		Benefit: "Reduces failed agent runs and accidental workflow drift by making build, test, and gx/JJ conventions explicit.",
		Evidence: []Evidence{
			{Label: "File", Value: "`AGENTS.md` is missing"},
		},
		Recommendation: "Add `AGENTS.md` with the expected build/test commands, version-control workflow, and repo-specific constraints.",
		Strength:       "Worth exploring",
		SourceIDs:      []string{"google-eng-practices"},
	}
}

func genericPackageNameFindings(ctx ReviewContext) []Finding {
	genericNames := map[string]struct{}{
		"common": {}, "shared": {}, "utils": {}, "util": {}, "helpers": {}, "helper": {}, "types": {}, "models": {}, "lib": {},
	}
	var evidence []Evidence
	for _, pkg := range ctx.Facts.GoPackages {
		name := filepath.Base(pkg.Path)
		if _, ok := genericNames[name]; ok {
			evidence = append(evidence, Evidence{Label: "Package", Value: fmt.Sprintf("`%s`", pkg.Path)})
		}
	}
	if len(evidence) == 0 {
		return nil
	}
	return []Finding{{
		ID:      "architecture.generic-package-name",
		Scopes:  []string{"architecture", "maintainability"},
		Title:   "Generic package names reduce Interface clarity",
		Summary: "Generic package names make it harder to infer what Module Interface callers should rely on.",
		Benefit: "Improves readability and navigation by making package names describe domain concepts instead of storage buckets for shared code.",
		Evidence: append(evidence, Evidence{
			Label: "Why this matters",
			Value: "A package name should describe the concept behind its Interface, not just that code is shared.",
		}),
		Recommendation: "Review these packages and rename or split them around the domain concept they actually expose.",
		Strength:       "Worth exploring",
		SourceIDs:      []string{"go-package-names", "go-code-review-comments"},
	}}
}

func implementationHeavyModuleFindings(ctx ReviewContext) []Finding {
	var candidates []PackageFact
	for _, pkg := range ctx.Facts.GoPackages {
		if pkg.GoFiles < 6 || ignorablePackage(pkg.Path) {
			continue
		}
		candidates = append(candidates, pkg)
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].GoFiles != candidates[j].GoFiles {
			return candidates[i].GoFiles > candidates[j].GoFiles
		}
		return candidates[i].Path < candidates[j].Path
	})
	limit := 5
	if len(candidates) < limit {
		limit = len(candidates)
	}
	var evidence []Evidence
	for _, pkg := range candidates[:limit] {
		evidence = append(evidence, Evidence{
			Label: "Module",
			Value: fmt.Sprintf("`%s` has %d implementation file(s) and %d test file(s)", pkg.Path, pkg.GoFiles, pkg.TestFiles),
		})
	}
	moduleNames := packagePaths(candidates[:limit])
	firstModule := moduleNames[0]
	return []Finding{{
		ID:      "architecture.implementation-heavy-module",
		Scopes:  []string{"architecture", "maintainability", "testing"},
		Title:   "Look for deeper Interfaces in the largest Modules",
		Summary: fmt.Sprintf("%s are the largest Go Modules in this repo; inspect whether their Interfaces hide workflow complexity or force callers to know too much ordering, configuration, or error behavior.", formatModuleList(moduleNames, 3)),
		Benefit: "Improves readability and testability by moving workflow ordering into one Module Interface, which should reduce caller knowledge and make future fixes more local.",
		Evidence: append(evidence, Evidence{
			Label: "Deletion test",
			Value: "If deleting one of these Modules would scatter the same complexity across callers, deepen its Interface; if deletion removes mostly pass-through code, collapse it.",
		}),
		Recommendation: fmt.Sprintf("Start with `%s`: list the exported functions and the top two callers, then identify one workflow detail callers currently need to know. Move that detail behind a named package function or type, and add a package-level test that exercises the new Interface through the caller-facing path.", firstModule),
		Strength:       "Worth exploring",
		SourceIDs:      []string{"fowler-architecture", "google-eng-practices"},
	}}
}

func packageWithoutTestsFindings(ctx ReviewContext) []Finding {
	var missing []PackageFact
	for _, pkg := range ctx.Facts.GoPackages {
		if pkg.GoFiles == 0 || pkg.TestFiles > 0 {
			continue
		}
		if pkg.GoFiles < 2 || ignorablePackage(pkg.Path) {
			continue
		}
		missing = append(missing, pkg)
	}
	if len(missing) == 0 {
		return nil
	}
	sort.Slice(missing, func(i, j int) bool {
		if missing[i].GoFiles != missing[j].GoFiles {
			return missing[i].GoFiles > missing[j].GoFiles
		}
		return missing[i].Path < missing[j].Path
	})
	limit := 8
	if len(missing) < limit {
		limit = len(missing)
	}
	var evidence []Evidence
	for _, pkg := range missing[:limit] {
		evidence = append(evidence, Evidence{
			Label: "Package",
			Value: fmt.Sprintf("`%s` has %d implementation file(s) and no test files", pkg.Path, pkg.GoFiles),
		})
	}
	moduleNames := packagePaths(missing[:limit])
	return []Finding{{
		ID:             "architecture.package-without-tests",
		Scopes:         []string{"architecture", "testing", "maintainability"},
		Title:          "Test meaningful Modules through their package seam",
		Summary:        fmt.Sprintf("%s have non-trivial Implementation but no colocated tests, so review cannot tell whether their Interfaces are easy to exercise through the package seam.", formatModuleList(moduleNames, 3)),
		Benefit:        "Improves refactor safety and exposes shallow Interfaces early: awkward tests point directly at seams that should be simplified.",
		Evidence:       evidence,
		Recommendation: fmt.Sprintf("Add tests for `%s` through its public package Interface first. If that feels awkward, use the friction to reshape the Module before extracting more helpers.", moduleNames[0]),
		Strength:       "Worth exploring",
		SourceIDs:      []string{"fowler-test-pyramid", "google-eng-practices"},
	}}
}

func missingReadmeFinding(facts RepoFacts) Finding {
	if present(facts.Docs, "README.md") {
		return Finding{}
	}
	return Finding{
		ID:      "docs.missing-readme",
		Scopes:  []string{"docs", "onboarding", "maintainability"},
		Title:   "Missing README",
		Summary: "The repo lacks the standard entrypoint document new maintainers expect first.",
		Benefit: "Improves onboarding speed by giving humans and agents one place to find setup, purpose, and common commands.",
		Evidence: []Evidence{
			{Label: "File", Value: "`README.md` is missing"},
		},
		Recommendation: "Add a concise `README.md` covering purpose, setup, common commands, and where deeper documentation lives.",
		Strength:       "Strong",
		SourceIDs:      []string{"diataxis", "write-the-docs-guide"},
	}
}

func noTestsFinding(facts RepoFacts) Finding {
	if facts.TestFileCount > 0 {
		return Finding{}
	}
	return Finding{
		ID:      "testing.no-tests",
		Scopes:  []string{"testing", "maintainability"},
		Title:   "No test files detected",
		Summary: "The local scan did not find test files, so there is no obvious test surface for review or future changes.",
		Benefit: "Improves regression safety and makes later code review more specific because failures can be tied to executable behavior.",
		Evidence: []Evidence{
			{Label: "Test files", Value: "0"},
		},
		Recommendation: "Add tests around the highest-leverage public Interfaces before expanding review automation.",
		Strength:       "Strong",
		SourceIDs:      []string{"fowler-test-pyramid", "google-eng-practices"},
	}
}

func noDependencyManifestFinding(facts RepoFacts) Finding {
	if len(facts.DependencyFiles) > 0 {
		return Finding{}
	}
	return Finding{
		ID:      "dependencies.no-manifest",
		Scopes:  []string{"dependencies", "maintainability"},
		Title:   "No dependency manifest detected",
		Summary: "The local scan did not find a standard dependency manifest, which limits dependency and supply-chain review.",
		Benefit: "Improves dependency review and reproducibility by making the build inputs visible to gx, CI, and teammates.",
		Evidence: []Evidence{
			{Label: "Dependency manifests", Value: "none detected"},
		},
		Recommendation: "Make sure the repo commits the dependency manifest and lockfile used by its build tooling.",
		Strength:       "Worth exploring",
		SourceIDs:      []string{"openssf-scorecard", "openssf-best-practices"},
	}
}

func jsDependencyLockFinding(facts RepoFacts) Finding {
	packageDirs := map[string]struct{}{}
	lockDirs := map[string]struct{}{}
	for _, file := range facts.DependencyFiles {
		dir := filepath.ToSlash(filepath.Dir(file))
		if dir == "." {
			dir = ""
		}
		switch filepath.Base(file) {
		case "package.json":
			packageDirs[dir] = struct{}{}
		case "package-lock.json", "bun.lock", "bun.lockb", "pnpm-lock.yaml", "yarn.lock":
			lockDirs[dir] = struct{}{}
		}
	}
	var missing []string
	for dir := range packageDirs {
		if _, ok := lockDirs[dir]; !ok {
			if dir == "" {
				missing = append(missing, "package.json")
			} else {
				missing = append(missing, fmt.Sprintf("%s/package.json", dir))
			}
		}
	}
	if len(missing) == 0 {
		return Finding{}
	}
	sort.Strings(missing)
	return Finding{
		ID:      "dependencies.javascript-unlocked",
		Scopes:  []string{"dependencies", "security", "maintainability"},
		Title:   "JavaScript dependencies are not locked",
		Summary: "A `package.json` was detected without a lockfile in the same directory, so installs may resolve differently across machines and CI.",
		Benefit: "Improves reproducibility and supply-chain review by keeping dependency resolution stable across local machines and CI.",
		Evidence: []Evidence{
			{Label: "Package manifests without colocated lockfile", Value: strings.Join(missing, ", ")},
		},
		Recommendation: "Commit the lockfile produced by the package manager for each JavaScript package directory.",
		Strength:       "Strong",
		SourceIDs:      []string{"openssf-scorecard", "openssf-best-practices", "slsa"},
	}
}

func ignorablePackage(path string) bool {
	if path == "." {
		return true
	}
	base := filepath.Base(path)
	if strings.HasPrefix(path, "cmd/") || strings.HasPrefix(path, "test/") {
		return true
	}
	switch base {
	case "buildconfig", "version":
		return true
	default:
		return false
	}
}
