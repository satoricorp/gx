package codereview

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	reviewPolicyPath               = "REVIEW.md"
	maxReviewPolicySummaryBytes    = 8000
	maxReviewReferenceRawBytes     = 1 << 20
	maxReviewReferenceSummaryBytes = 5000
	reviewReferenceTimeout         = 8 * time.Second
)

type ReviewPolicy struct {
	Present     bool              `json:"present"`
	Path        string            `json:"path,omitempty"`
	ByteSize    int               `json:"byte_size,omitempty"`
	Text        string            `json:"text,omitempty"`
	Summarized  bool              `json:"summarized,omitempty"`
	URLs        []string          `json:"urls,omitempty"`
	References  []ReviewReference `json:"references,omitempty"`
	ModelHints  []ReviewModelHint `json:"model_hints,omitempty"`
	RiskPaths   []RiskPath        `json:"risk_paths,omitempty"`
	Diagnostics []string          `json:"diagnostics,omitempty"`
}

type RiskPath struct {
	Glob    string `json:"glob,omitempty"`
	Message string `json:"message,omitempty"`
	Raw     string `json:"raw,omitempty"`
}

type ReviewReference struct {
	URL        string `json:"url"`
	Status     string `json:"status"`
	Summary    string `json:"summary,omitempty"`
	Error      string `json:"error,omitempty"`
	ByteSize   int    `json:"byte_size,omitempty"`
	Summarized bool   `json:"summarized,omitempty"`
	Capped     bool   `json:"capped,omitempty"`
}

type ReviewModelHint struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Role     string `json:"role,omitempty"`
	Raw      string `json:"raw,omitempty"`
}

func LoadReviewPolicy(ctx context.Context, repoRoot string) ReviewPolicy {
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
		URLs:       extractReviewPolicyURLs(text),
		ModelHints: parseReviewModelHints(text),
		RiskPaths:  parseReviewRiskPaths(text),
	}
	if summarized {
		policy.Diagnostics = append(policy.Diagnostics, fmt.Sprintf("%s summarized from %d byte(s)", reviewPolicyPath, len(data)))
	}
	for _, ref := range fetchReviewReferences(ctx, policy.URLs) {
		policy.References = append(policy.References, ref)
		switch {
		case ref.Status != "ok":
			policy.Diagnostics = append(policy.Diagnostics, fmt.Sprintf("%s: %s", ref.URL, ref.Error))
		case ref.Capped:
			policy.Diagnostics = append(policy.Diagnostics, fmt.Sprintf("%s capped at %d byte(s) before summarization", ref.URL, maxReviewReferenceRawBytes))
		case ref.Summarized:
			policy.Diagnostics = append(policy.Diagnostics, fmt.Sprintf("%s summarized from %d byte(s)", ref.URL, ref.ByteSize))
		}
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
	for _, ref := range p.References {
		if ref.Status != "ok" || strings.TrimSpace(ref.Summary) == "" {
			continue
		}
		snippets = append(snippets, ContextSnippet{
			Kind:      "review_reference",
			Ref:       ref.URL,
			Text:      ref.Summary,
			Source:    "review.md-url",
			Publisher: publisherFromURLHost(ref.URL),
			URL:       ref.URL,
			Title:     ref.URL,
		})
	}
	return snippets
}

func (p ReviewPolicy) QueryText() string {
	if !p.Present {
		return ""
	}
	var parts []string
	if strings.TrimSpace(p.Text) != "" {
		parts = append(parts, "repo REVIEW.md policy:\n"+strings.TrimSpace(p.Text))
	}
	for _, ref := range p.References {
		if ref.Status != "ok" || strings.TrimSpace(ref.Summary) == "" {
			continue
		}
		parts = append(parts, "REVIEW.md referenced guidance "+ref.URL+":\n"+strings.TrimSpace(ref.Summary))
	}
	if len(p.ModelHints) > 0 {
		var hints []string
		for _, hint := range p.ModelHints {
			hints = append(hints, strings.TrimSpace(hint.Provider+":"+hint.Model))
		}
		parts = append(parts, "requested reviewer models: "+strings.Join(hints, ", "))
	}
	return strings.Join(parts, "\n\n")
}

func (p ReviewPolicy) OpenAIModelHint() string {
	for _, hint := range p.ModelHints {
		if hint.Provider == "openai" && strings.TrimSpace(hint.Model) != "" {
			return strings.TrimSpace(hint.Model)
		}
	}
	return ""
}

func (p ReviewPolicy) AnthropicModelHint() string {
	for _, hint := range p.ModelHints {
		if hint.Provider == "anthropic" && strings.TrimSpace(hint.Model) != "" {
			return strings.TrimSpace(hint.Model)
		}
	}
	return ""
}

func fetchReviewReferences(ctx context.Context, urls []string) []ReviewReference {
	if len(urls) == 0 {
		return nil
	}
	client := &http.Client{Timeout: reviewReferenceTimeout}
	out := make([]ReviewReference, 0, len(urls))
	for _, rawURL := range urls {
		out = append(out, fetchReviewReference(ctx, client, rawURL))
	}
	return out
}

func fetchReviewReference(ctx context.Context, client *http.Client, rawURL string) ReviewReference {
	item := ReviewReference{URL: rawURL, Status: "failed"}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		item.Error = "invalid URL"
		return item
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		item.Error = err.Error()
		return item
	}
	req.Header.Set("Accept", "text/plain,text/markdown,text/html;q=0.9,*/*;q=0.1")
	resp, err := client.Do(req)
	if err != nil {
		item.Error = err.Error()
		return item
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		item.Error = "status " + resp.Status
		return item
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxReviewReferenceRawBytes+1))
	item.ByteSize = len(body)
	if len(body) > maxReviewReferenceRawBytes {
		body = body[:maxReviewReferenceRawBytes]
		item.Capped = true
	}
	text := plainReviewReferenceText(resp.Header.Get("Content-Type"), string(body))
	if strings.TrimSpace(text) == "" {
		item.Error = "empty text content"
		return item
	}
	summary, summarized := summarizeReviewText(text, maxReviewReferenceSummaryBytes)
	item.Status = "ok"
	item.Error = ""
	item.Summary = summary
	item.Summarized = summarized || item.Capped
	return item
}

var (
	reviewURLPattern   = regexp.MustCompile(`https?://[^\s<>()"']+`)
	htmlTagPattern     = regexp.MustCompile(`(?s)<[^>]+>`)
	htmlScriptPattern  = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	htmlStylePattern   = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	reviewModelPattern = regexp.MustCompile(`(?i)\b(?:(openai|anthropic):)?((?:gpt|o)[A-Za-z0-9._-]*|(?:anthropic\.)?claude[A-Za-z0-9._-]*)\b`)
	reviewRiskPathPattern = regexp.MustCompile(`(?i)^\s*risk-path:\s*(.+?)\s*(?:—|-)\s*(.+?)\s*$`)
)

func extractReviewPolicyURLs(text string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, match := range reviewURLPattern.FindAllString(text, -1) {
		cleaned := strings.TrimRight(match, ".,;:)]}>\"'")
		if parsed, err := url.ParseRequestURI(cleaned); err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			continue
		}
		if _, ok := seen[cleaned]; ok {
			continue
		}
		seen[cleaned] = struct{}{}
		out = append(out, cleaned)
	}
	sort.Strings(out)
	return out
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
		match := reviewRiskPathPattern.FindStringSubmatch(line)
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

func parseReviewModelHints(text string) []ReviewModelHint {
	seen := map[string]struct{}{}
	var out []ReviewModelHint
	for _, line := range strings.Split(text, "\n") {
		lower := strings.ToLower(line)
		if !strings.Contains(lower, "model") && !strings.Contains(lower, "reviewer") && !strings.Contains(lower, "use ") {
			continue
		}
		for _, match := range reviewModelPattern.FindAllStringSubmatch(line, -1) {
			provider := strings.ToLower(strings.TrimSpace(match[1]))
			model := strings.TrimSpace(match[2])
			if provider == "" {
				provider = inferReviewModelProvider(model)
			}
			model = normalizeReviewModel(provider, model)
			if provider == "" || model == "" {
				continue
			}
			role := inferReviewModelRole(line)
			key := provider + "\x00" + model + "\x00" + role
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, ReviewModelHint{
				Provider: provider,
				Model:    model,
				Role:     role,
				Raw:      strings.TrimSpace(line),
			})
		}
	}
	return out
}

func inferReviewModelProvider(model string) string {
	lower := strings.ToLower(model)
	switch {
	case strings.HasPrefix(lower, "gpt") || strings.HasPrefix(lower, "o"):
		return "openai"
	case strings.HasPrefix(lower, "claude") || strings.HasPrefix(lower, "anthropic.claude"):
		return "anthropic"
	default:
		return ""
	}
}

func normalizeReviewModel(provider, model string) string {
	model = strings.TrimSpace(model)
	if provider == "anthropic" && strings.HasPrefix(strings.ToLower(model), "claude") {
		return "anthropic." + model
	}
	return model
}

func inferReviewModelRole(line string) string {
	lower := strings.ToLower(line)
	for _, role := range []string{"correctness", "security", "dissent", "performance", "testing", "maintainability", "docs", "dependencies"} {
		if strings.Contains(lower, role) {
			return role
		}
	}
	return "general"
}

func plainReviewReferenceText(contentType, text string) string {
	lower := strings.ToLower(contentType)
	if strings.Contains(lower, "html") || strings.Contains(text, "<html") || strings.Contains(text, "<body") {
		text = htmlScriptPattern.ReplaceAllString(text, " ")
		text = htmlStylePattern.ReplaceAllString(text, " ")
		text = htmlTagPattern.ReplaceAllString(text, " ")
		text = html.UnescapeString(text)
	}
	return normalizeReviewWhitespace(text)
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
