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
	"sync"
)

const (
	defaultJudgeMaxOutputTokens = 4000
	maxJudgeCandidateBytes      = 24 * 1024
	// maxAdvisoryFindings caps the fallback path (no judge / judge error), where
	// we have no impact signal and rely on heuristic strength. The judged path is
	// gated by impact instead and is intentionally uncapped.
	maxAdvisoryFindings = 3
	// minSurfaceConfidence is the confidence floor for surfacing a confirmed
	// finding whose impact is only functional (not breaking).
	minSurfaceConfidence = 0.5
)

// Judge impact levels, in ascending order of real-world consequence.
const (
	impactNone       = "none"
	impactCosmetic   = "cosmetic"
	impactFunctional = "functional"
	impactBreaking   = "breaking"
)

// impactRank orders impact levels. Unknown/blank ranks as functional, so a
// malformed judge response surfaces the finding rather than silently hiding it.
func impactRank(impact string) int {
	switch strings.ToLower(strings.TrimSpace(impact)) {
	case impactBreaking:
		return 3
	case impactFunctional:
		return 2
	case impactCosmetic:
		return 1
	case impactNone:
		return 0
	default:
		return 2
	}
}

type FindingJudge interface {
	Judge(ctx context.Context, req judgeRequest) ([]judgeResult, error)
}

type judgeAvailabilityReporter interface {
	Available() bool
}

type bedrockReviewJudge struct {
	client *bedrockAnthropicReviewer
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
	Impact           string  `json:"impact"`
	Severity         int     `json:"severity"`
	Confidence       float64 `json:"confidence"`
	VerificationNote string  `json:"verification_note"`
}

// judgeFromEnvWithPolicy builds the verification model.
//
// The judge runs on Bedrock like the reviewers, but deliberately on a third
// model: it decides which candidate findings a human sees, and a model that
// grades its own output is not a filter. It is also the only model on the
// sequential path — both reviewers must finish before it starts — so it is the
// cheaper Sonnet rather than a third Opus.
//
// Precedence is env > REVIEW.md hint > default here, inverted from the reviewer
// legs. That asymmetry is pre-existing (GX_REVIEW_JUDGE_MODEL has always won
// over the policy hint) and is preserved rather than tidied: a repo pins which
// models review its code, but an operator debugging a verification failure
// needs to be able to override the judge from the environment without editing
// the repo's REVIEW.md.
//
// The judge takes the same transport precedence as the reviewer legs — local
// AWS credentials if present, gx-cloud otherwise — resolved by the same
// function, so a review cannot end up with reviewers on one wire and its judge
// on another and no way to tell from the output.
//
// There is no OpenAI judge any more, and no fallback judge of any kind. With
// neither transport reachable the judge is unavailable and says so; the engine
// then keeps every candidate finding under capAdvisoryFindings rather than
// silently treating an unreachable judge as a verdict. That fail-open path is
// what the batching work made visible and it is still the only degraded
// behavior here.
func judgeFromEnvWithPolicy(policy *ReviewPolicy) FindingJudge {
	if judgeDisabledFromEnv() {
		return nil
	}
	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		return unavailableReviewJudge{reason: err.Error()}
	}
	return bedrockReviewJudge{client: newBedrockReviewer(plan.newTransport(), resolveBedrockJudgeModel(policy))}
}

// resolveBedrockJudgeModel applies env > policy hint > default.
func resolveBedrockJudgeModel(policy *ReviewPolicy) string {
	hint := ""
	if policy != nil {
		hint = policy.JudgeModelHint()
	}
	return normalizeBedrockModelID(firstNonEmpty(
		os.Getenv("GX_REVIEW_JUDGE_MODEL"),
		hint,
		defaultBedrockJudgeModel,
	))
}

func judgeDisabledFromEnv() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_JUDGE")), "0")
}

func (j bedrockReviewJudge) Available() bool {
	return j.client != nil
}

func (j bedrockReviewJudge) Judge(ctx context.Context, req judgeRequest) ([]judgeResult, error) {
	if j.client == nil {
		return nil, fmt.Errorf("review judge is not configured")
	}
	content, err := j.client.completeJSON(ctx, judgeDeveloperPrompt(), mustJSON(req), defaultJudgeMaxOutputTokens)
	if err != nil {
		return nil, err
	}
	return parseJudgeResponse(content)
}

// parseJudgeResponse tolerates a model that wraps its JSON in prose or a fenced
// block. The deleted OpenAI path could rely on response_format:json_object to
// guarantee a bare object; Bedrock has no equivalent, so the same
// truncation-shaped failure the batching work fixed would otherwise come back as
// a parse error on every batch.
func parseJudgeResponse(content string) ([]judgeResult, error) {
	var parsed judgeResponse
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		return parsed.Results, nil
	}
	trimmed := extractJSONObject(content)
	if trimmed == "" {
		return nil, fmt.Errorf("decode judge JSON: no JSON object in response")
	}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil, fmt.Errorf("decode judge JSON: %w", err)
	}
	return parsed.Results, nil
}

func extractJSONObject(content string) string {
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start < 0 || end <= start {
		return ""
	}
	return content[start : end+1]
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
		"You are GX Review Judge. Verify candidate findings against the provided real file content, and decide which ones a human must review before merge.",
		"Return JSON only with shape {\"results\":[{\"candidate_id\":string,\"verdict\":\"confirmed|unverified|wrong\",\"impact\":\"breaking|functional|cosmetic|none\",\"severity\":1-5,\"confidence\":0-1,\"verification_note\":string}]}",
		"Verdict must be confirmed, unverified, or wrong. Only confirmed findings survive; mark weak or unsupported claims unverified or wrong.",
		"Confirm only when the named files and file content snippets support the title, summary, recommendation, and evidence.",
		"If a candidate names no file, treat that as a strike and do not confirm unless static tool evidence conclusively proves it.",
		"impact rates the real-world consequence if the finding is acted on or ignored:",
		"- breaking: risks incorrect behavior, a crash, data loss, a security hole, or a broken build or test.",
		"- functional: affects behavior, an API or contract, or maintainability in a way a reviewer should weigh.",
		"- cosmetic: style, naming, wording, or clarity only — nothing that can break.",
		"- none: not a real issue.",
		"Reserve breaking and functional for changes a senior engineer would want to see before merge. If you are confident a finding cannot break anything, mark it cosmetic or none — GX will not surface those.",
		"Severity is 1-5. Confidence is 0-1. Include a one-line verification note.",
	}, "\n")
}

// prepareFindingsForJudge collapses duplicates before verification.
//
// This runs in addition to the per-shard merge inside multiAIReviewer, and it
// is the one that catches cross-shard copies: a file that appears in two shards
// gets reviewed twice, and neither shard can see the other's findings.
func prepareFindingsForJudge(ctx context.Context, adjudicator duplicateAdjudicator, findings []Finding) dedupeOutcome {
	return dedupeFindings(ctx, adjudicator, findings)
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

// applyJudgeResults keeps only confirmed findings the judge considers worth a
// human's time before merge. It drops anything marked cosmetic or none (a change
// the judge is confident cannot break anything) and drops low-confidence
// non-breaking findings as noise; breaking findings always survive. Results are
// ranked by impact and left uncapped, so an empty result is a trustworthy
// "nothing here needs review".
func applyJudgeResults(findings []Finding, results []judgeResult) []Finding {
	byID := map[string]judgeResult{}
	for _, result := range results {
		id := strings.TrimSpace(result.CandidateID)
		if id == "" {
			continue
		}
		byID[id] = result
	}
	type judged struct {
		finding Finding
		result  judgeResult
	}
	var kept []judged
	for _, finding := range findings {
		result, ok := byID[finding.ID]
		if !ok || !strings.EqualFold(strings.TrimSpace(result.Verdict), "confirmed") {
			continue
		}
		rank := impactRank(result.Impact)
		if rank <= impactRank(impactCosmetic) {
			continue // confidently benign — the whole point is to not surface these
		}
		if rank == impactRank(impactFunctional) && result.Confidence < minSurfaceConfidence {
			continue // non-breaking and low confidence — noise
		}
		finding.Strength = strengthFromImpact(result.Impact)
		if note := strings.TrimSpace(result.VerificationNote); note != "" {
			finding.Evidence = append(finding.Evidence, Evidence{Label: "Judge verification", Value: note})
		}
		kept = append(kept, judged{finding: finding, result: result})
	}
	sort.SliceStable(kept, func(i, j int) bool {
		if ri, rj := impactRank(kept[i].result.Impact), impactRank(kept[j].result.Impact); ri != rj {
			return ri > rj
		}
		// Corroboration outranks the judge's own confidence, within an impact
		// class. Two independently prompted flagship models arriving at the
		// same problem from different context is stronger evidence that the
		// problem is real than one model's self-reported certainty about it,
		// and the judge — which sees candidates one at a time — has no way to
		// know the agreement happened.
		if li, lj := len(kept[i].finding.Corroboration), len(kept[j].finding.Corroboration); li != lj {
			return li > lj
		}
		if kept[i].result.Confidence != kept[j].result.Confidence {
			return kept[i].result.Confidence > kept[j].result.Confidence
		}
		if kept[i].result.Severity != kept[j].result.Severity {
			return kept[i].result.Severity > kept[j].result.Severity
		}
		return kept[i].finding.ID < kept[j].finding.ID
	})
	out := make([]Finding, 0, len(kept))
	for _, k := range kept {
		out = append(out, k.finding)
	}
	return out
}

func strengthFromImpact(impact string) string {
	if impactRank(impact) >= impactRank(impactBreaking) {
		return "Strong"
	}
	return "Worth exploring"
}

// judgeBatchSize is how many candidates one judge call may carry.
//
// The judge answers with one object per candidate, so its output length is
// linear in the candidate count while defaultJudgeMaxOutputTokens is fixed. It
// is the output budget, not the input, that binds. Measured against gpt-5.5
// with the production prompt: 3 candidates -> 640 output tokens, 6 -> 851,
// 12 -> 1671, 24 -> 2581, 36 -> 3635, and 48 -> the response stopped at
// max_output_tokens mid-object.
//
// A truncated response is not partial data, it is unparseable JSON, so Judge
// returns an error and every candidate loses its verdict at once. Before
// batching, a review that produced 48 candidate findings surfaced at most
// maxAdvisoryFindings (3) of them and said nothing about why. 24 keeps the
// worst measured batch at roughly 65% of the output budget, which leaves room
// for the longer verification notes a genuinely complex finding attracts.
const judgeBatchSize = 24

// judgeBatchOutcome is what one batch of candidates came back with.
type judgeBatchOutcome struct {
	// Judged are the candidates a verdict was returned for.
	Judged []Finding
	// Unjudged are candidates whose batch failed. They are kept rather than
	// dropped: "the judge could not be reached" and "the judge rejected this"
	// are different facts and must not produce the same review.
	Unjudged []Finding
	// Err is the first batch failure, for the degraded-reasons line.
	Err error
	// Batches / BatchesFailed describe the fan-out for the same line.
	Batches       int
	BatchesFailed int
}

// runJudge verifies candidates in concurrent batches.
//
// Concurrency is the point as much as batching. The judge is a second model
// call that used to run strictly after the reviewer finished, and it is a
// double-digit share of review latency; splitting it into batches that run at
// the same time turns a serial 2N-candidate call into one N-candidate call's
// worth of wall clock. Measured on the production prompt, a single 48-candidate
// call takes 47s and fails; two concurrent 24-candidate batches take about 40s
// and succeed.
func runJudge(ctx context.Context, judge FindingJudge, reviewContext ReviewContext, candidates []Finding) judgeBatchOutcome {
	if len(candidates) == 0 {
		return judgeBatchOutcome{}
	}
	var batches [][]Finding
	for start := 0; start < len(candidates); start += judgeBatchSize {
		end := start + judgeBatchSize
		if end > len(candidates) {
			end = len(candidates)
		}
		batches = append(batches, candidates[start:end])
	}
	type batchResult struct {
		results []judgeResult
		err     error
	}
	out := make([]batchResult, len(batches))
	var wg sync.WaitGroup
	for i, batch := range batches {
		i, batch := i, batch
		wg.Add(1)
		go func() {
			defer wg.Done()
			results, err := judge.Judge(ctx, buildJudgeRequest(reviewContext, batch))
			out[i] = batchResult{results: results, err: err}
		}()
	}
	wg.Wait()

	outcome := judgeBatchOutcome{Batches: len(batches)}
	for i, batch := range batches {
		if out[i].err != nil {
			outcome.BatchesFailed++
			if outcome.Err == nil {
				outcome.Err = out[i].err
			}
			outcome.Unjudged = append(outcome.Unjudged, batch...)
			continue
		}
		outcome.Judged = append(outcome.Judged, applyJudgeResults(batch, out[i].results)...)
	}
	return outcome
}

func capAdvisoryFindings(findings []Finding) []Finding {
	out := append([]Finding(nil), findings...)
	sort.SliceStable(out, func(i, j int) bool {
		left := strengthRank(out[i].Strength)
		right := strengthRank(out[j].Strength)
		if left != right {
			return left < right
		}
		// This path throws findings away, so corroboration matters more here
		// than anywhere else: with no judge verdict to rank by, "both reviewers
		// found it" is the best evidence available for which three survive.
		if li, lj := len(out[i].Corroboration), len(out[j].Corroboration); li != lj {
			return li > lj
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
