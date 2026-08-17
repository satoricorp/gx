package codereview

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

const (
	reviewPolicyPath            = "REVIEW.md"
	maxReviewPolicySummaryBytes = 8000
)

// REVIEW.md describes the repository to its reviewer. It does not configure the
// reviewer, and the line between those two is the reason this file is smaller
// than it used to be.
//
// Two features crossed it and were removed. The first scraped every URL out of
// REVIEW.md and fetched it, checking only that the scheme was http or https —
// no host check, and a client that followed redirects. On a cloud CI runner
// `http://169.254.169.254/latest/meta-data/iam/security-credentials/` satisfied
// every check in that path, and up to a megabyte of the response was
// summarized into the review report and the model prompt. The second let a line
// of REVIEW.md prose pick the reviewer and judge models, ahead of the operator's
// own environment variables, with no allowlist.
//
// Both gave the repository under review control over what the reviewer does,
// which is backwards: the repository is the subject, not the operator. Guidance
// that matters belongs in the file itself, where it is diffable and reviewable,
// rather than behind a URL whose contents can change after the line was written.
// Which model reviews is the operator's call, and stays in the environment.
type ReviewPolicy struct {
	Present     bool       `json:"present"`
	Path        string     `json:"path,omitempty"`
	ByteSize    int        `json:"byte_size,omitempty"`
	Text        string     `json:"text,omitempty"`
	Summarized  bool       `json:"summarized,omitempty"`
	RiskPaths   []RiskPath `json:"risk_paths,omitempty"`
	Diagnostics []string   `json:"diagnostics,omitempty"`
	// Rules are the named review rules the repository declared: every `##`
	// heading in REVIEW.md that is not one of the reserved sections. Exceptions
	// are the bullets under `## Exceptions`; Extends names a pack the repository
	// builds on (e.g. "gx:recommended@2"). See parseReviewRules for the grammar
	// and for the line REVIEW.md is not allowed to cross.
	Rules      []PolicyRule `json:"rules,omitempty"`
	Exceptions []string     `json:"exceptions,omitempty"`
	Extends    string       `json:"extends,omitempty"`
}

type RiskPath struct {
	Glob    string `json:"glob,omitempty"`
	Message string `json:"message,omitempty"`
	Raw     string `json:"raw,omitempty"`
}

// PolicyRule is one named rule as REVIEW.md (or a built-in pack written in the
// same grammar) declares it. ID is the namespaced, canonical spelling that
// findings carry as RuleID; Slug is the bare kebab name; Name is the heading
// verbatim; Text is the guidance the model is given.
type PolicyRule struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Text     string `json:"text"`
	Advisory bool   `json:"advisory,omitempty"`
}

func LoadReviewPolicy(repoRoot string) ReviewPolicy {
	data, err := os.ReadFile(filepath.Join(repoRoot, reviewPolicyPath))
	if err != nil {
		return ReviewPolicy{Present: false, Path: reviewPolicyPath}
	}
	text := string(data)
	summary, summarized := summarizeReviewText(text, maxReviewPolicySummaryBytes)
	rules, exceptions, extends, diagnostics := parseReviewRules(text)
	policy := ReviewPolicy{
		Present:    true,
		Path:       reviewPolicyPath,
		ByteSize:   len(data),
		Text:       summary,
		Summarized: summarized,
		RiskPaths:  parseReviewRiskPaths(text),
		Rules:      rules,
		Exceptions: exceptions,
		Extends:    extends,
	}
	if summarized {
		policy.Diagnostics = append(policy.Diagnostics, fmt.Sprintf("%s summarized from %d byte(s)", reviewPolicyPath, len(data)))
	}
	policy.Diagnostics = append(policy.Diagnostics, diagnostics...)
	return policy
}

// RuleDefs renders the repository's rules as the minimal form the brief and
// the gates request carry: the namespaced ID and a one-line summary.
func (p ReviewPolicy) RuleDefs() []RuleDef {
	if len(p.Rules) == 0 {
		return nil
	}
	defs := make([]RuleDef, 0, len(p.Rules))
	for _, rule := range p.Rules {
		defs = append(defs, RuleDef{ID: rule.ID, Summary: policyRuleSummary(rule.Text)})
	}
	return defs
}

// AllRuleDefs is the full set of named rules in force for a review: the
// built-in recommended pack first, then whatever the repository's REVIEW.md
// declares. A nil policy means only the pack applies.
func AllRuleDefs(policy *ReviewPolicy) []RuleDef {
	var defs []RuleDef
	if pack, err := RecommendedPack(); err == nil {
		for _, rule := range pack {
			defs = append(defs, RuleDef{ID: rule.ID, Summary: policyRuleSummary(rule.Text)})
		}
	}
	if policy != nil {
		defs = append(defs, policy.RuleDefs()...)
	}
	return defs
}

// policyRuleSummary is the first sentence of a rule's text: up to the first
// ". " (or the end of the first line), capped at 160 bytes.
func policyRuleSummary(text string) string {
	const maxSummaryBytes = 160
	summary := strings.TrimSpace(text)
	if i := strings.Index(summary, ". "); i >= 0 {
		summary = summary[:i+1]
	} else if i := strings.IndexByte(summary, '\n'); i >= 0 {
		summary = strings.TrimSpace(summary[:i])
	}
	if len(summary) > maxSummaryBytes {
		summary = trimReviewTextAtBoundary(summary, maxSummaryBytes)
	}
	return summary
}

func (p ReviewPolicy) ContextSnippets() []ContextSnippet {
	if !p.Present {
		return nil
	}
	var snippets []ContextSnippet
	if strings.TrimSpace(p.Text) != "" {
		snippets = append(snippets, ContextSnippet{
			Kind:      "review_policy",
			Ref:       p.Path,
			Text:      p.Text,
			Source:    "local",
			Publisher: "this repo",
			File:      p.Path,
		})
	}
	return snippets
}

func (p ReviewPolicy) QueryText() string {
	if !p.Present {
		return ""
	}
	if strings.TrimSpace(p.Text) == "" {
		return ""
	}
	return "repo REVIEW.md policy:\n" + strings.TrimSpace(p.Text)
}

// A risk-path line is `risk-path: <glob> — <why>`, and the separator is the
// whole difficulty. One pattern matching either dash found the wrong one:
// `risk-path: src/my-app/** — why` split at the hyphen inside `my-app`, giving
// the glob `src/my`, which matches nothing and silently drops the entry. Any
// repository with a hyphen in a directory name — web-ui, api-server, my-app —
// was affected, and nothing said so.
//
// So the two dashes get different rules. An em or en dash never appears in a
// path, so it separates with or without surrounding spaces. A plain hyphen does
// appear in paths, so it only separates when spaced, which is how anyone writes
// it anyway. Em dash is tried first: `a-b — c` has both, and the em dash is the
// one that was meant.
var (
	reviewRiskPathDashPattern   = regexp.MustCompile(`(?i)^\s*risk-path:\s*(.+?)\s*[—–]\s*(.+?)\s*$`)
	reviewRiskPathHyphenPattern = regexp.MustCompile(`(?i)^\s*risk-path:\s*(.+?)\s+-\s+(.+?)\s*$`)
)

// matchRiskPathLine returns the glob and message of a risk-path line.
func matchRiskPathLine(line string) []string {
	if match := reviewRiskPathDashPattern.FindStringSubmatch(line); len(match) == 3 {
		return match
	}
	return reviewRiskPathHyphenPattern.FindStringSubmatch(line)
}

func parseReviewRiskPaths(text string) []RiskPath {
	inSection := false
	seen := map[string]struct{}{}
	var out []RiskPath
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "#") {
			heading := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			inSection = strings.Contains(heading, "high-risk paths") || strings.Contains(heading, "high risk paths")
			continue
		}
		if lower == "high-risk paths" || lower == "high risk paths" {
			inSection = true
			continue
		}
		if !inSection {
			continue
		}
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(lower, "#") {
			inSection = false
			continue
		}
		match := matchRiskPathLine(line)
		if len(match) != 3 {
			continue
		}
		glob := strings.TrimSpace(match[1])
		message := strings.TrimSpace(match[2])
		if glob == "" || message == "" {
			continue
		}
		key := glob + "\x00" + message
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, RiskPath{
			Glob:    glob,
			Message: message,
			Raw:     trimmed,
		})
	}
	return out
}

// A named rule is a `##` heading and the prose under it. Exactly two hashes:
// `#` is the document title, and `###` or deeper is structure inside a rule
// body — a sub-heading under a rule stays part of that rule, it does not open
// a new one. Three headings are reserved and never become rules: high-risk
// paths (parseReviewRiskPaths owns it), Exceptions (bullets to collect), and
// Extends (the pack this file builds on, either as the first line under the
// heading or inline after a colon).
//
// The line REVIEW.md may not cross, restated from the type comment above:
// this parser hands the reviewer rule names, guidance text, and exceptions —
// nothing else. The only thing a heading can set is the Advisory boolean via
// the `Advisory:` prefix. Severity beyond that, which model reviews, whether
// a gate passes, and anything to fetch are the operator's, and a rule body
// that contains what looks like a directive (`model:`, `severity:`, `fetch:`,
// `gate:`) is prose to the model, not a setting to this code. Nothing here
// interprets it, and nothing here should start to.
var (
	reviewRuleHeadingPattern  = regexp.MustCompile(`^##\s+(.+?)\s*$`)
	reviewRuleSlugPattern     = regexp.MustCompile(`[^a-z0-9]+`)
	reviewRuleAdvisoryPattern = regexp.MustCompile(`(?i)^advisory:\s*`)
)

const (
	reviewRuleNamespace = "REVIEW.md/"
	maxReviewRuleBytes  = 2000
)

// slugifyReviewRule lowercases the heading and folds every run of characters
// outside [a-z0-9] into one hyphen, trimming hyphens at either end.
func slugifyReviewRule(name string) string {
	slug := reviewRuleSlugPattern.ReplaceAllString(strings.ToLower(name), "-")
	return strings.Trim(slug, "-")
}

// stripReviewBullet removes one leading list marker (-, *, •) and the
// whitespace around it.
func stripReviewBullet(line string) string {
	line = strings.TrimSpace(line)
	for _, marker := range []string{"-", "*", "•"} {
		if strings.HasPrefix(line, marker) {
			return strings.TrimSpace(strings.TrimPrefix(line, marker))
		}
	}
	return line
}

func parseReviewRules(text string) (rules []PolicyRule, exceptions []string, extends string, diagnostics []string) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	type section int
	const (
		sectionNone section = iota
		sectionRule
		sectionRiskPaths
		sectionExceptions
		sectionExtends
	)
	current := sectionNone
	var (
		name     string
		advisory bool
		body     []string
	)
	seen := map[string]struct{}{}

	flush := func() {
		if current != sectionRule {
			return
		}
		slug := slugifyReviewRule(name)
		ruleText := trimReviewTextAtBoundary(normalizeReviewWhitespace(strings.Join(body, "\n")), maxReviewRuleBytes)
		switch {
		case slug == "":
			diagnostics = append(diagnostics, fmt.Sprintf("%s: rule heading %q has no usable name; skipped", reviewPolicyPath, name))
		case ruleText == "":
			diagnostics = append(diagnostics, fmt.Sprintf("%s: rule %q has no guidance text; skipped", reviewPolicyPath, name))
		default:
			if _, dup := seen[slug]; dup {
				diagnostics = append(diagnostics, fmt.Sprintf("%s: duplicate rule %q; keeping the first", reviewPolicyPath, slug))
				return
			}
			seen[slug] = struct{}{}
			rules = append(rules, PolicyRule{
				ID:       reviewRuleNamespace + slug,
				Slug:     slug,
				Name:     name,
				Text:     ruleText,
				Advisory: advisory,
			})
		}
	}

	for _, line := range strings.Split(text, "\n") {
		if match := reviewRuleHeadingPattern.FindStringSubmatch(line); len(match) == 2 {
			flush()
			heading := strings.TrimSpace(match[1])
			lower := strings.ToLower(heading)
			body = nil
			name = ""
			advisory = false
			switch {
			case lower == "high-risk paths" || lower == "high risk paths":
				current = sectionRiskPaths
			case lower == "exceptions":
				current = sectionExceptions
			case lower == "extends":
				current = sectionExtends
			case strings.HasPrefix(lower, "extends:"):
				current = sectionNone
				if inline := strings.TrimSpace(heading[len("extends:"):]); inline != "" && extends == "" {
					extends = inline
				}
			default:
				current = sectionRule
				if loc := reviewRuleAdvisoryPattern.FindStringIndex(heading); loc != nil {
					advisory = true
					heading = strings.TrimSpace(heading[loc[1]:])
				}
				name = heading
			}
			continue
		}
		switch current {
		case sectionRule:
			body = append(body, line)
		case sectionExceptions:
			if item := stripReviewBullet(line); item != "" {
				exceptions = append(exceptions, item)
			}
		case sectionExtends:
			if item := stripReviewBullet(line); item != "" {
				if extends == "" {
					extends = item
				}
				current = sectionNone
			}
		}
	}
	flush()
	return rules, exceptions, extends, diagnostics
}

func MatchRiskPathGlob(pattern, file string) bool {
	pattern = filepath.ToSlash(strings.TrimSpace(pattern))
	file = filepath.ToSlash(strings.TrimSpace(file))
	if pattern == "" || file == "" {
		return false
	}
	if !strings.Contains(pattern, "**") {
		matched, err := filepath.Match(pattern, file)
		return err == nil && matched
	}
	if pattern == "**" {
		return true
	}
	if strings.HasSuffix(pattern, "/**") {
		prefix := strings.TrimSuffix(pattern, "/**")
		if prefix == "" {
			return true
		}
		return file == prefix || strings.HasPrefix(file, prefix+"/")
	}
	if strings.HasPrefix(pattern, "**/") {
		suffix := strings.TrimPrefix(pattern, "**/")
		if suffix == "" {
			return true
		}
		return file == suffix || strings.HasSuffix(file, "/"+suffix) || strings.HasPrefix(file, suffix+"/")
	}
	left, right, ok := strings.Cut(pattern, "**")
	if !ok {
		return false
	}
	if left != "" && !strings.HasPrefix(file, left) {
		return false
	}
	rest := strings.TrimPrefix(file, left)
	if right != "" && !strings.HasSuffix(rest, right) {
		return false
	}
	return true
}

func summarizeReviewText(text string, maxBytes int) (string, bool) {
	text = normalizeReviewWhitespace(text)
	if maxBytes <= 0 || len(text) <= maxBytes {
		return text, false
	}
	paragraphs := strings.Split(text, "\n\n")
	var b strings.Builder
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		remaining := maxBytes - b.Len() - 80
		if remaining <= 0 {
			break
		}
		if len(paragraph) > remaining {
			paragraph = trimReviewTextAtBoundary(paragraph, remaining)
		}
		b.WriteString(paragraph)
		if b.Len() >= maxBytes-80 {
			break
		}
	}
	summary := strings.TrimSpace(b.String())
	if summary == "" {
		summary = trimReviewTextAtBoundary(text, maxBytes-80)
	}
	return strings.TrimSpace(summary) + "\n\n[summary generated from longer review context]", true
}

func trimReviewTextAtBoundary(text string, maxBytes int) string {
	if len(text) <= maxBytes {
		return text
	}
	if maxBytes <= 0 {
		return ""
	}
	cut := maxBytes
	for cut > 0 && !utf8Boundary(text[cut]) {
		cut--
	}
	if cut <= 0 {
		cut = maxBytes
	}
	trimmed := text[:cut]
	for i := len(trimmed) - 1; i >= 0 && i > len(trimmed)-160; i-- {
		if trimmed[i] == '.' || trimmed[i] == '\n' {
			return strings.TrimSpace(trimmed[:i+1])
		}
	}
	for i := len(trimmed) - 1; i >= 0 && i > len(trimmed)-80; i-- {
		if unicode.IsSpace(rune(trimmed[i])) {
			return strings.TrimSpace(trimmed[:i])
		}
	}
	return strings.TrimSpace(trimmed)
}

func utf8Boundary(b byte) bool {
	return b&0xC0 != 0x80
}

func normalizeReviewWhitespace(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")
	var out []string
	blank := false
	for _, line := range lines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if strings.TrimSpace(line) == "" {
			if !blank {
				out = append(out, "")
			}
			blank = true
			continue
		}
		out = append(out, line)
		blank = false
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}
