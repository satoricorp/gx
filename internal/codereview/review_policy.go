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
}

type RiskPath struct {
	Glob    string `json:"glob,omitempty"`
	Message string `json:"message,omitempty"`
	Raw     string `json:"raw,omitempty"`
}

func LoadReviewPolicy(repoRoot string) ReviewPolicy {
	data, err := os.ReadFile(filepath.Join(repoRoot, reviewPolicyPath))
	if err != nil {
		return ReviewPolicy{Present: false, Path: reviewPolicyPath}
	}
	text := string(data)
	summary, summarized := summarizeReviewText(text, maxReviewPolicySummaryBytes)
	policy := ReviewPolicy{
		Present:    true,
		Path:       reviewPolicyPath,
		ByteSize:   len(data),
		Text:       summary,
		Summarized: summarized,
		RiskPaths:  parseReviewRiskPaths(text),
	}
	if summarized {
		policy.Diagnostics = append(policy.Diagnostics, fmt.Sprintf("%s summarized from %d byte(s)", reviewPolicyPath, len(data)))
	}
	return policy
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
