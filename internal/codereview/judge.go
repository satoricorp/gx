package codereview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultJudgeMaxOutputTokens = 4000
	maxJudgeCandidateBytes      = 24 * 1024
	maxAdvisoryFindings         = 3
)

type FindingJudge interface {
	Judge(ctx context.Context, req judgeRequest) ([]judgeResult, error)
}

type judgeAvailabilityReporter interface {
	Available() bool
}

type openAIReviewJudge struct {
	client *responsesAIReviewer
}

type fallbackReviewJudge struct {
	primary  FindingJudge
	fallback FindingJudge
}

type unavailableReviewJudge struct {
	reason string
}

type judgeRequest struct {
	RepoRoot     string           `json:"repo_root"`
	ChangedFiles []string         `json:"changed_files"`
	Candidates   []judgeCandidate `json:"candidates"`
}

type judgeCandidate struct {
	ID                  string                `json:"candidate_id"`
	Title               string                `json:"title"`
	Summary             string                `json:"summary"`
	Recommendation      string                `json:"recommendation"`
	Evidence            []string              `json:"evidence"`
	SourcePublishers    []string              `json:"source_publishers,omitempty"`
	NamedFiles          []string              `json:"named_files"`
	FileContentSnippets []judgeContentSnippet `json:"file_content_snippets"`
}

type judgeContentSnippet struct {
	File string `json:"file"`
	Text string `json:"text"`
}

type judgeResponse struct {
	Results []judgeResult `json:"results"`
}

type judgeResult struct {
	CandidateID      string  `json:"candidate_id"`
	Verdict          string  `json:"verdict"`
	Severity         int     `json:"severity"`
	Confidence       float64 `json:"confidence"`
	VerificationNote string  `json:"verification_note"`
}

func judgeFromEnvWithPolicy(policy *ReviewPolicy) FindingJudge {
	if judgeDisabledFromEnv() {
		return nil
	}
	model := strings.TrimSpace(os.Getenv("GX_REVIEW_JUDGE_MODEL"))
	if model == "" && policy != nil {
		model = policy.OpenAIModelHint()
	}
	model = firstNonEmpty(model, os.Getenv("GX_REVIEW_OPENAI_MODEL"), os.Getenv("GX_REVIEW_MODEL"), os.Getenv("OPENAI_MODEL"), defaultReviewModel)

	direct, directErr := directOpenAIReviewerFromEnv(model)
	cloudReviewer := cloudOpenAIReviewerFromEnv(model)
	var directJudge FindingJudge
	if directClient, ok := direct.(*responsesAIReviewer); ok {
		directJudge = openAIReviewJudge{client: directClient}
	}
	var cloudJudge FindingJudge
	if cloudClient, ok := cloudReviewer.(*responsesAIReviewer); ok {
		cloudJudge = openAIReviewJudge{client: cloudClient}
	}
	if directJudge != nil && cloudJudge != nil {
		return fallbackReviewJudge{primary: directJudge, fallback: cloudJudge}
	}
	if directJudge != nil {
		return directJudge
	}
	if directErr != nil {
		return unavailableReviewJudge{reason: directErr.Error()}
	}
	if cloudJudge != nil {
		return cloudJudge
	}
	return unavailableReviewJudge{reason: "review judge is not configured"}
}

func judgeDisabledFromEnv() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_JUDGE")), "0")
}

func (j openAIReviewJudge) Available() bool {
	return j.client != nil
}

func (j openAIReviewJudge) Judge(ctx context.Context, req judgeRequest) ([]judgeResult, error) {
	if j.client == nil {
		return nil, fmt.Errorf("review judge is not configured")
	}
	content, err := j.client.completeJSON(ctx, judgeDeveloperPrompt(), mustJSON(req), defaultJudgeMaxOutputTokens)
	if err != nil {
		return nil, err
	}
	var parsed judgeResponse
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("decode judge JSON: %w", err)
	}
	return parsed.Results, nil
}

func (j fallbackReviewJudge) Available() bool {
	return judgeAvailable(j.primary) || judgeAvailable(j.fallback)
}

func (j fallbackReviewJudge) Judge(ctx context.Context, req judgeRequest) ([]judgeResult, error) {
	results, err := j.primary.Judge(ctx, req)
	if err == nil {
		return results, nil
	}
	fallbackResults, fallbackErr := j.fallback.Judge(ctx, req)
	if fallbackErr == nil {
		return fallbackResults, nil
	}
	return nil, fmt.Errorf("%w; fallback judge failed: %v", err, fallbackErr)
}

func (j unavailableReviewJudge) Available() bool {
	return false
}

func (j unavailableReviewJudge) Judge(context.Context, judgeRequest) ([]judgeResult, error) {
	return nil, fmt.Errorf("%s", j.reason)
}

func judgeAvailable(judge FindingJudge) bool {
	if judge == nil {
		return false
	}
	if reporter, ok := judge.(judgeAvailabilityReporter); ok {
		return reporter.Available()
	}
	return true
}

func judgeDeveloperPrompt() string {
	return strings.Join([]string{
		"You are GX Review Judge. Verify candidate findings against the provided real file content.",
		"Return JSON only with shape {\"results\":[{\"candidate_id\":string,\"verdict\":\"confirmed|unverified|wrong\",\"severity\":1-5,\"confidence\":0-1,\"verification_note\":string}]}",
		"Verdict must be confirmed, unverified, or wrong.",
		"Drop weak claims by marking them unverified or wrong; only confirmed findings should survive.",
		"Confirm only when the named files and file content snippets support the title, summary, recommendation, and evidence.",
		"If a candidate names no file, treat that as a strike and do not confirm unless static tool evidence conclusively proves it.",
		"Severity is 1-5. Confidence is 0-1. Include a one-line verification note.",
	}, "\n")
}

func prepareFindingsForJudge(ctx ReviewContext, findings []Finding) []Finding {
	return mergeNearDuplicateFindings(findings)
}

func buildJudgeRequest(ctx ReviewContext, findings []Finding) judgeRequest {
	candidates := make([]judgeCandidate, 0, len(findings))
	sourcePublishers := sourcePublisherMap(ctx.Sources)
	for _, finding := range findings {
		candidate := judgeCandidate{
			ID:               finding.ID,
			Title:            finding.Title,
			Summary:          finding.Summary,
			Recommendation:   finding.Recommendation,
			Evidence:         findingEvidenceStrings(finding.Evidence),
			SourcePublishers: uniqueStrings(append(append([]string{}, finding.SourcePublishers...), findingSourcePublishers(finding.SourceIDs, sourcePublishers)...)),
		}
		candidate.NamedFiles = namedFilesForFinding(ctx, finding)
		candidate.FileContentSnippets = judgeFileContentSnippets(ctx.Brief.RepoRoot, candidate.NamedFiles)
		candidates = append(candidates, candidate)
	}
	return judgeRequest{
		RepoRoot:     ctx.Brief.RepoRoot,
		ChangedFiles: normalizedChangedFiles(ctx.Brief.Static.ChangedFiles),
		Candidates:   candidates,
	}
}

func applyJudgeResults(findings []Finding, results []judgeResult) []Finding {
	byID := map[string]judgeResult{}
	for _, result := range results {
		id := strings.TrimSpace(result.CandidateID)
		if id == "" {
			continue
		}
		byID[id] = result
	}
	var out []Finding
	for _, finding := range findings {
		result, ok := byID[finding.ID]
		if !ok || !strings.EqualFold(strings.TrimSpace(result.Verdict), "confirmed") {
			continue
		}
		finding.Strength = strengthFromJudge(result.Severity, result.Confidence)
		if note := strings.TrimSpace(result.VerificationNote); note != "" {
			finding.Evidence = append(finding.Evidence, Evidence{Label: "Judge verification", Value: note})
		}
		out = append(out, finding)
	}
	return out
}

func strengthFromJudge(severity int, confidence float64) string {
	if severity >= 4 && confidence >= 0.70 {
		return "Strong"
	}
	if severity <= 2 || confidence < 0.50 {
		return "Speculative"
	}
	return "Worth exploring"
}

func capAdvisoryFindings(findings []Finding) []Finding {
	out := append([]Finding(nil), findings...)
	sort.SliceStable(out, func(i, j int) bool {
		left := strengthRank(out[i].Strength)
		right := strengthRank(out[j].Strength)
		if left != right {
			return left < right
		}
		return out[i].ID < out[j].ID
	})
	if len(out) > maxAdvisoryFindings {
		out = out[:maxAdvisoryFindings]
	}
	return out
}

func splitBlockingToolFindings(findings []Finding) ([]Finding, []Finding) {
	var blocking []Finding
	var advisory []Finding
	for _, finding := range findings {
		if strings.HasPrefix(strings.TrimSpace(finding.ID), "tools.") && finding.Strength == "Blocking" {
			blocking = append(blocking, finding)
			continue
		}
		advisory = append(advisory, finding)
	}
	return blocking, advisory
}

func sourcePublisherMap(sources []Source) map[string]string {
	out := map[string]string{}
	for _, source := range sources {
		id := strings.TrimSpace(source.ID)
		if id == "" {
			continue
		}
		publisher := strings.TrimSpace(source.Publisher)
		if publisher == "" {
			publisher = strings.TrimSpace(source.Title)
		}
		if publisher == "" {
			publisher = id
		}
		out[id] = publisher
	}
	return out
}

func findingSourcePublishers(ids []string, publishers map[string]string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, id := range ids {
		publisher := strings.TrimSpace(publishers[strings.TrimSpace(id)])
		if publisher == "" {
			continue
		}
		if _, ok := seen[publisher]; ok {
			continue
		}
		seen[publisher] = struct{}{}
		out = append(out, publisher)
	}
	return out
}

func findingEvidenceStrings(evidence []Evidence) []string {
	var out []string
	for _, item := range evidence {
		text := strings.TrimSpace(strings.Join([]string{item.Label, item.Value}, ": "))
		text = strings.Trim(text, ": ")
		if text != "" {
			out = append(out, text)
		}
	}
	return out
}

func namedFilesForFinding(ctx ReviewContext, finding Finding) []string {
	known := knownReviewFiles(ctx)
	changed := normalizedChangedFiles(ctx.Brief.Static.ChangedFiles)
	seen := map[string]struct{}{}
	var parsed []string
	for _, text := range []string{finding.Title, finding.Summary, finding.Recommendation, evidenceText(finding.Evidence)} {
		for _, file := range parseReviewFilePaths(text, known, ctx.Brief.RepoRoot) {
			if _, ok := seen[file]; ok {
				continue
			}
			seen[file] = struct{}{}
			parsed = append(parsed, file)
		}
	}
	changedSet := map[string]struct{}{}
	for _, file := range changed {
		changedSet[file] = struct{}{}
	}
	sort.SliceStable(parsed, func(i, j int) bool {
		_, leftChanged := changedSet[parsed[i]]
		_, rightChanged := changedSet[parsed[j]]
		if leftChanged != rightChanged {
			return leftChanged
		}
		return parsed[i] < parsed[j]
	})
	return parsed
}

func knownReviewFiles(ctx ReviewContext) map[string]struct{} {
	out := map[string]struct{}{}
	for _, file := range ctx.Facts.Files {
		file = filepath.ToSlash(strings.TrimSpace(file))
		if file != "" {
			out[file] = struct{}{}
		}
	}
	for _, file := range ctx.Brief.Static.ChangedFiles {
		file = filepath.ToSlash(strings.TrimSpace(file))
		if file != "" {
			out[file] = struct{}{}
		}
	}
	return out
}

var reviewPathPattern = regexp.MustCompile("`([^`]+)`|([A-Za-z0-9_./-]+\\.[A-Za-z0-9][A-Za-z0-9_-]*(?::\\d+)?)")

func parseReviewFilePaths(text string, known map[string]struct{}, repoRoot string) []string {
	var out []string
	seen := map[string]struct{}{}
	for _, match := range reviewPathPattern.FindAllStringSubmatch(text, -1) {
		raw := match[1]
		if raw == "" {
			raw = match[2]
		}
		file := normalizeReviewPath(raw)
		if file == "" {
			continue
		}
		if len(known) > 0 {
			if _, ok := known[file]; !ok {
				continue
			}
		} else if repoRoot != "" {
			if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(file))); err != nil {
				continue
			}
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		out = append(out, file)
	}
	return out
}

func normalizeReviewPath(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.Trim(raw, "`'\".,;:()[]{}")
	if raw == "" || strings.Contains(raw, "://") {
		return ""
	}
	if index := strings.LastIndex(raw, ":"); index > 1 {
		if _, err := strconv.Atoi(raw[index+1:]); err == nil {
			raw = raw[:index]
		}
	}
	raw = filepath.ToSlash(filepath.Clean(filepath.FromSlash(raw)))
	raw = strings.TrimPrefix(raw, "./")
	if raw == "." || strings.HasPrefix(raw, "../") || filepath.IsAbs(raw) {
		return ""
	}
	if !strings.Contains(raw, ".") {
		return ""
	}
	return raw
}

func judgeFileContentSnippets(repoRoot string, files []string) []judgeContentSnippet {
	if strings.TrimSpace(repoRoot) == "" {
		return nil
	}
	remaining := maxJudgeCandidateBytes
	var out []judgeContentSnippet
	for _, file := range files {
		if remaining <= 0 {
			break
		}
		path := filepath.Join(repoRoot, filepath.FromSlash(file))
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		text := string(data)
		if len(text) > remaining {
			text = text[:remaining] + "\n[truncated]\n"
		}
		remaining -= len(text)
		out = append(out, judgeContentSnippet{File: file, Text: text})
	}
	return out
}

func mergeNearDuplicateFindings(findings []Finding) []Finding {
	type groupedFinding struct {
		finding   Finding
		tokens    map[string]struct{}
		providers map[string]struct{}
		count     int
	}
	var groups []groupedFinding
	for _, finding := range findings {
		tokens := duplicateTokens(finding)
		provider := findingProvider(finding)
		merged := false
		for i := range groups {
			if nearDuplicateTokens(tokens, groups[i].tokens) {
				groups[i].finding = mergeFindingMetadata(groups[i].finding, finding)
				groups[i].providers[provider] = struct{}{}
				groups[i].count++
				merged = true
				break
			}
		}
		if !merged {
			groups = append(groups, groupedFinding{
				finding:   finding,
				tokens:    tokens,
				providers: map[string]struct{}{provider: {}},
				count:     1,
			})
		}
	}
	out := make([]Finding, 0, len(groups))
	for _, group := range groups {
		finding := group.finding
		if group.count > 1 {
			finding.Evidence = append(finding.Evidence, Evidence{
				Label: "Agreement",
				Value: fmt.Sprintf("Similar findings merged from %d providers/sources: %s", group.count, strings.Join(sortedSet(group.providers), ", ")),
			})
		}
		out = append(out, finding)
	}
	return out
}

func mergeResolvedSources(left, right []ResolvedSource) []ResolvedSource {
	seen := map[string]struct{}{}
	var out []ResolvedSource
	for _, src := range append(append([]ResolvedSource{}, left...), right...) {
		key := strings.TrimSpace(src.ID) + "|" + strings.TrimSpace(ResolvedSourceLabel(src))
		if key == "|" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, src)
	}
	return out
}

func mergeFindingMetadata(left Finding, right Finding) Finding {
	left.Evidence = append(left.Evidence, right.Evidence...)
	left.SourceIDs = uniqueStrings(append(left.SourceIDs, right.SourceIDs...))
	left.SourcePublishers = uniqueStrings(append(left.SourcePublishers, right.SourcePublishers...))
	left.ResolvedSources = mergeResolvedSources(left.ResolvedSources, right.ResolvedSources)
	if strengthRank(right.Strength) < strengthRank(left.Strength) {
		left.Strength = right.Strength
	}
	return left
}

func duplicateTokens(finding Finding) map[string]struct{} {
	text := normalizeDuplicateText(finding.Title + " " + firstSentence(finding.Summary))
	out := map[string]struct{}{}
	for _, token := range strings.Fields(text) {
		if len(token) < 3 || duplicateStopWords[token] {
			continue
		}
		out[token] = struct{}{}
	}
	return out
}

var duplicateStopWords = map[string]bool{
	"the": true, "and": true, "for": true, "with": true, "that": true, "this": true, "from": true, "into": true,
	"can": true, "are": true, "was": true, "were": true, "has": true, "have": true, "but": true, "not": true,
}

func normalizeDuplicateText(text string) string {
	text = strings.ToLower(text)
	var b strings.Builder
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '/' || r == '_' || r == '-' {
			b.WriteRune(r)
			continue
		}
		b.WriteByte(' ')
	}
	return b.String()
}

func firstSentence(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	for _, sep := range []string{".", "\n"} {
		if index := strings.Index(text, sep); index >= 0 {
			return strings.TrimSpace(text[:index])
		}
	}
	return text
}

func nearDuplicateTokens(left, right map[string]struct{}) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}
	overlap := 0
	for token := range left {
		if _, ok := right[token]; ok {
			overlap++
		}
	}
	smaller := len(left)
	if len(right) < smaller {
		smaller = len(right)
	}
	return float64(overlap)/float64(smaller) >= 0.85
}

func findingProvider(finding Finding) string {
	id := strings.TrimSpace(finding.ID)
	if id == "" {
		return "unknown"
	}
	if index := strings.Index(id, "."); index > 0 {
		return id[:index]
	}
	return id
}

func sortedSet(values map[string]struct{}) []string {
	out := make([]string, 0, len(values))
	for value := range values {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
