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
	// defaultJudgeMaxOutputTokens bounds one judge batch's reply.
	//
	// Sized against measurement, not taste, and re-measured whenever the reply
	// shape changes — because overrunning it is not a partial answer, it is an
	// unparseable one that costs every candidate in the batch its verdict.
	//
	// The measurement that set this: replaying the real 21-candidate judge
	// request captured from a live gx review against
	// us.anthropic.claude-sonnet-4-6, with the cap high enough that the reply's
	// natural length was what got measured, produced 7,019 to 8,326 output
	// tokens. Scaled to a full judgeBatchSize of 24 that is roughly 9.5K on the
	// worst run. The previous 12000 was set when the judge answered without an
	// analysis field and a full batch cost 2.1K to 3.1K tokens; the analysis
	// field roughly tripled the reply, which would have put the worst case at
	// about 80% of that cap — the same thin margin that truncated the 4000-token
	// cap before it, arrived at from the other direction.
	//
	// A cap is not a reservation: unused budget is neither billed nor waited on,
	// so the headroom is free and there is no reason to run close to the line.
	defaultJudgeMaxOutputTokens = 32000
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
	CandidateID string `json:"candidate_id"`
	// Analysis is where the judge does its thinking, and it is deliberately
	// read into a field nothing renders.
	//
	// The model has no hidden reasoning channel — bedrockRequestBody carries
	// anthropic_version, max_tokens, system and messages, and nothing else, so
	// extended thinking is not enabled and never was. Telling it to "reason
	// silently" therefore did not move the reasoning anywhere, it deleted it:
	// measured across 18 paired runs of one identical 10-candidate request,
	// output fell 24% and the confirm rate doubled from 10.0% to 20.6%, taking a
	// candidate whose stated mechanism the provided file content refutes from
	// wrong 18/18 to confirmed 2/18.
	//
	// A JSON object is emitted key by key, so a reasoning field written before
	// verdict is reasoning the verdict is conditioned on. That is what this
	// field buys, and why it stays even though no report ever prints it.
	Analysis         string  `json:"analysis"`
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
func judgeFromEnv() FindingJudge {
	if judgeDisabledFromEnv() {
		return nil
	}
	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		return unavailableReviewJudge{reason: err.Error()}
	}
	return bedrockReviewJudge{client: newBedrockReviewer(plan.newTransport(), resolveBedrockJudgeModel())}
}

// resolveBedrockJudgeModel applies env > default. The reviewed repository does
// not get a say in which model verifies the findings against it.
func resolveBedrockJudgeModel() string {
	return normalizeBedrockModelID(firstNonEmpty(
		os.Getenv("GX_REVIEW_JUDGE_MODEL"),
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
	completion, err := j.client.completeJSON(ctx, judgeDeveloperPrompt(), mustJSON(req), defaultJudgeMaxOutputTokens)
	if err != nil {
		return nil, err
	}
	results, parseErr := parseJudgeResponse(completion.Text, judgeCandidateIDs(req.Candidates))
	if parseErr != nil && completion.truncated() {
		// A batch cut off at the cap has no complete verdict set, only a
		// prefix of one. Report it as a failure so the batch fails open and
		// says why, rather than reading a partial answer as a whole one and
		// silently dropping every candidate the model never reached.
		return nil, describeTruncatedCompletion("judge", defaultJudgeMaxOutputTokens, parseErr)
	}
	return results, parseErr
}

// parseJudgeResponse tolerates a model that wraps its JSON in prose or a fenced
// block. The deleted OpenAI path could rely on response_format:json_object to
// guarantee a bare object; Bedrock has no equivalent, so a model free to answer
// however it likes would otherwise come back as a parse error on whole batches.
//
// asked is the candidate IDs this reply is supposed to account for, and it is
// what makes the tolerance safe. Without it the rule was "first decodable object
// carrying a results key", and the judge prompt prints the literal shape
// {"results":[...]} to the model — so a model that restates its own output shape
// in a preamble hands the parser an empty result set that decodes perfectly.
// Measured: a reply reading `I will answer with {"results":[]} if none of the
// candidates check out.` followed by the real verdicts returned 0 verdicts and
// no error, which (applyJudgeResults being delete-only) drops every candidate in
// the batch with BatchesFailed=0, so nothing tells the reader. Scoring by how
// many of the asked candidates an object actually answers picks the answer over
// the quotation.
//
// There is deliberately no "the whole reply parses, take it" fast path any more.
// It looked like a strict-first optimization and was a third silent-drop route:
// json.Unmarshal of `{"overview":"looks fine"}` into judgeResponse succeeds with
// no results and no error, which is the same batch-wide silent drop wearing a
// different hat. Requiring the results key on every path costs one decode of an
// already-in-memory string.
func parseJudgeResponse(content string, asked []string) ([]judgeResult, error) {
	object := pickAnsweringJSONObject(content, judgeAnswerScore(asked))
	if object == "" {
		return nil, fmt.Errorf("decode judge JSON: no complete JSON object with a results field in response")
	}
	var parsed judgeResponse
	if err := json.Unmarshal([]byte(object), &parsed); err != nil {
		return nil, fmt.Errorf("decode judge JSON: %w", err)
	}
	return parsed.Results, nil
}

// judgeAnswerScore ranks a decoded object by how much of this request it
// answers: 0 for "not a verdict set at all", otherwise one point for carrying
// the results key plus one for each asked candidate it returns a verdict for.
//
// With no asked IDs — a caller that does not know what was requested — every
// verdict set scores 1 and the first one wins, which is the old behavior and the
// best that can be done without knowing the question.
func judgeAnswerScore(asked []string) func(map[string]json.RawMessage) int {
	want := make(map[string]struct{}, len(asked))
	for _, id := range asked {
		if id = strings.TrimSpace(id); id != "" {
			want[id] = struct{}{}
		}
	}
	return func(object map[string]json.RawMessage) int {
		raw, ok := object["results"]
		if !ok {
			return 0
		}
		if len(want) == 0 {
			return 1
		}
		var results []judgeResult
		if err := json.Unmarshal(raw, &results); err != nil {
			return 0
		}
		score := 1
		for _, result := range results {
			if _, ok := want[strings.TrimSpace(result.CandidateID)]; ok {
				score++
			}
		}
		return score
	}
}

// judgeCandidateIDs is what the batch asked about, in the order it asked.
func judgeCandidateIDs(candidates []judgeCandidate) []string {
	out := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if id := strings.TrimSpace(candidate.ID); id != "" {
			out = append(out, id)
		}
	}
	return out
}

// extractJSONObject returns the substring of content that best answers a request
// whose reply carries one of requiredKeys, preferring the object that carries
// the most of them. With no requiredKeys, the first object wins.
//
// It scans for a *decodable* object rather than slicing from the first "{" to
// the last "}", which is what this used to do and is why a real review ran
// degraded with "invalid character 'r' looking for beginning of object key
// string". Measured, not inferred: replaying a captured 24-candidate judge
// request against us.anthropic.claude-sonnet-4-6 reproduced it 1 run in 20. The
// model prefixed 5,439 characters of reasoning before its ```json fence, and
// that prose quoted source code — `if (existsSync(p)) { try { return ... } }` —
// so the first "{" in the response was 4KB before the JSON and the naive slice
// began mid-sentence. The JSON inside the fence was complete and valid the whole
// time; every finding in that batch lost its verdict to a substring bug.
//
// Key count is a weaker discriminator than the judge's, which counts answered
// candidate IDs (see judgeAnswerScore), and it is weaker on purpose: this is the
// caller for replies with no identifiers to match — the reviewer legs' overview
// and recommendations — where "the object carrying most of the shape" is the
// most the content supports. A preamble quoting a one-key fragment of the shape
// loses to the real reply, which carries several; a preamble quoting the whole
// shape would still win, and that is a known limit rather than a claim of
// safety.
//
// A truncated response has no complete object anywhere, so this returns "" and
// the caller reports a decode failure. That is deliberate: partial verdicts are
// worse than none, because a missing verdict drops a real finding silently.
func extractJSONObject(content string, requiredKeys ...string) string {
	return pickAnsweringJSONObject(content, func(object map[string]json.RawMessage) int {
		if len(requiredKeys) == 0 {
			return 1
		}
		score := 0
		for _, key := range requiredKeys {
			if _, ok := object[key]; ok {
				score++
			}
		}
		return score
	})
}

// pickAnsweringJSONObject returns the highest-scoring decodable JSON object in
// content, or "" if score rejects every one of them. Ties go to the earliest,
// so a model that simply answers is never second-guessed.
//
// json.Decoder is what makes the scan correct: it tracks strings and escapes, so
// a "}" inside a verification note ends nothing.
//
// An object that scores nothing is stepped over one byte at a time rather than
// skipped whole, so that an answer nested inside a wrapper — {"result":{"results":
// [...]}} — is still found. Skipping non-answers whole was faster and was a
// regression: it could turn a reply the old first-match extractor recovered into
// a decode failure. A scoring object is skipped whole, because its own children
// cannot be a better answer than it is.
func pickAnsweringJSONObject(content string, score func(map[string]json.RawMessage) int) string {
	best := ""
	bestScore := 0
	for offset := 0; offset < len(content); {
		index := strings.IndexByte(content[offset:], '{')
		if index < 0 {
			break
		}
		start := offset + index
		decoder := json.NewDecoder(strings.NewReader(content[start:]))
		var probe map[string]json.RawMessage
		if err := decoder.Decode(&probe); err != nil {
			offset = start + 1
			continue
		}
		points := score(probe)
		if points == 0 {
			offset = start + 1
			continue
		}
		if points > bestScore {
			best, bestScore = content[start:start+int(decoder.InputOffset())], points
		}
		offset = start + int(decoder.InputOffset())
	}
	return best
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

// judgeDeveloperPrompt states the output contract twice, at the top and at the
// bottom, because the model reads both ends hardest — and it moves the judge's
// reasoning inside the object rather than banning it.
//
// The contract is not decoration. Bedrock's InvokeModel has no
// response_format:json_object, and the two mechanisms that would enforce a shape
// from the request side — Anthropic structured outputs and forced tool_choice —
// are both unavailable here: the request body for the cloud transport is a fixed
// struct that GX Cloud's /gx/bedrock/fight normalizes, so a field only the
// direct-AWS path could send would leave the judge behaving differently
// depending on which wire the user is on, which is precisely the split this
// package refuses to ship. So the prompt is the enforcement, and
// parseJudgeResponse is the seatbelt.
//
// The sentences below are aimed at an observed failure, not a hypothetical one:
// the model answered a 24-candidate batch with 5,439 characters of prose
// reasoning and then a ```json fence. The prose quoted code containing braces,
// which broke the response parser outright.
//
// The first attempt at fixing that said "Reason silently, then emit only the
// object", and it was a regression wearing a fix's clothes. There is no silent
// channel to reason in — see judgeResult.Analysis — so the instruction did not
// relocate the reasoning, it removed it. A/B on one identical 10-candidate
// request, 18 paired runs, prompt the only variable: output tokens fell 24%
// (median 1313 -> 1005) and the confirm rate doubled, 10.0% -> 20.6% (Fisher
// exact p = 7.9e-3). A candidate claiming a bare `catch { }` swallows only JSON
// errors — refuted by the very snippet the judge was handed — went from wrong
// 18/18 to confirmed 2/18 at confidence 0.85, which is high enough to reach the
// reader. The judge was not formatting better, it was checking less.
//
// So the analysis field carries what the preamble used to, in the one place it
// can do its job: inside the object, before the verdict it justifies.
func judgeDeveloperPrompt() string {
	return strings.Join([]string{
		"You are GX Review Judge. Verify candidate findings against the provided real file content, and decide which ones a human must review before merge.",
		"OUTPUT CONTRACT: reply with one JSON object and nothing else. The first character of your reply must be { and the last must be }.",
		"Do not write anything before or after that object — no preamble, no commentary, no summary — and do not wrap it in a markdown code fence.",
		"Return JSON only with shape {\"results\":[{\"candidate_id\":string,\"analysis\":string,\"verdict\":\"confirmed|unverified|wrong\",\"impact\":\"breaking|functional|cosmetic|none\",\"severity\":1-5,\"confidence\":0-1,\"verification_note\":string}]}",
		"THINK IN THE analysis FIELD. Write it first, before the verdict of the same object, and use it to do the actual work: quote the lines of the provided file content that decide the claim, state what that code really does, and name any part of the claimed mechanism the file contradicts. Then let the verdict follow from it.",
		"Take as many sentences in analysis as the candidate needs. An analysis that only restates the claim is a verdict guessed rather than checked, and a wrong confirmation costs a reviewer more than a long analysis costs you.",
		"analysis is the only place reasoning may appear. Never write it outside the JSON object. Every result object must carry a non-empty analysis; a result without one has skipped the check the field exists to force.",
		"analysis and verification_note are JSON string values, so quote code inside them with backticks and never with a raw double quote. One unescaped \" makes the whole reply undecodable and costs every candidate in this batch its verdict, not just the one you were writing about.",
		"Verdict must be confirmed, unverified, or wrong. Only confirmed findings survive; mark weak or unsupported claims unverified or wrong.",
		"Confirm only when the named files and file content snippets support the title, summary, recommendation, and evidence. A finding whose conclusion may be right but whose stated mechanism the file content refutes is wrong, not confirmed — a reviewer acting on it would look for a bug that is not there.",
		"Read absolute claims literally. When a candidate says something never happens, always happens, or appears nowhere in a file, one counter-example in the provided content refutes it; look for that counter-example before you confirm, and say in the analysis whether you found one.",
		"If a candidate names no file, treat that as a strike and do not confirm unless static tool evidence conclusively proves it.",
		"impact rates the real-world consequence if the finding is acted on or ignored:",
		"- breaking: risks incorrect behavior, a crash, data loss, a security hole, or a broken build or test.",
		"- functional: affects behavior, an API or contract, or maintainability in a way a reviewer should weigh.",
		"- cosmetic: style, naming, wording, or clarity only — nothing that can break.",
		"- none: not a real issue.",
		"Reserve breaking and functional for changes a senior engineer would want to see before merge. If you are confident a finding cannot break anything, mark it cosmetic or none — GX will not surface those.",
		"Severity is 1-5. Confidence is 0-1. verification_note is one line for the reader, summarizing the analysis.",
		"Emit exactly one result object per candidate_id you were given, in the order given. A candidate you leave out is a finding GX must ship unverified, so leave none out.",
		"Reply with the JSON object alone. No preamble, no fence, no trailing remarks.",
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
// maxAdvisoryFindings (3) of them and said nothing about why.
//
// Re-measured on Bedrock after the judge moved there, and again after the
// analysis field landed: a real 21-candidate request now costs 7.0K to 8.3K
// output tokens, so a full batch of 24 sits at roughly 9.5K, under a third of
// defaultJudgeMaxOutputTokens. The size is kept at 24 rather than widened with
// the budget, because the batch is also the blast radius: one unparseable reply
// costs every candidate in its batch a verdict, and two concurrent batches cost
// less wall clock than one call twice the size.
const judgeBatchSize = 24

// judgeBatchOutcome is what one batch of candidates came back with.
type judgeBatchOutcome struct {
	// Judged are the candidates a verdict was returned for.
	Judged []Finding
	// Unjudged are candidates that got no verdict — because their batch failed,
	// or because a batch that succeeded simply did not mention them. They are
	// kept rather than dropped: "the judge could not be reached" and "the judge
	// rejected this" are different facts and must not produce the same review.
	Unjudged []Finding
	// Err is the first batch failure, for the degraded-reasons line.
	Err error
	// Batches / BatchesFailed describe the fan-out for the same line.
	Batches       int
	BatchesFailed int
	// Unanswered counts candidates a *successful* batch returned no verdict
	// for. It is separate from BatchesFailed because it is a separate fact with
	// a separate cause: the call worked, the JSON parsed, and the model just
	// left findings out. Without this count that shortfall is indistinguishable
	// from a judge that considered every candidate and rejected most of them.
	Unanswered int
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
		// A verdict set is checked against the candidates it was asked about
		// before it is applied. applyJudgeResults is delete-only, so a candidate
		// the model never mentioned is indistinguishable there from one it
		// rejected — and this model does leave entries out of list replies: the
		// de-duplicator, which runs the same Sonnet, reports exactly that
		// shortfall on real reviews. Without this split a reply covering 4 of 10
		// candidates dropped the other 6 as unconfirmed with BatchesFailed=0,
		// which reads to the engine as a clean, complete verification.
		//
		// Splitting rather than failing the batch keeps both halves honest: the
		// answered candidates get the verdicts the judge actually reached, and
		// the unanswered ones fall back to unjudged instead of being deleted by
		// a verdict nobody gave.
		verdicts := answeredCandidateIDs(out[i].results)
		var judged, silent []Finding
		for _, finding := range batch {
			if _, ok := verdicts[strings.TrimSpace(finding.ID)]; ok {
				judged = append(judged, finding)
				continue
			}
			silent = append(silent, finding)
		}
		outcome.Judged = append(outcome.Judged, applyJudgeResults(judged, out[i].results)...)
		if len(silent) > 0 {
			outcome.Unanswered += len(silent)
			outcome.Unjudged = append(outcome.Unjudged, silent...)
		}
	}
	return outcome
}

// answeredCandidateIDs is the set of candidates a verdict set speaks to. A
// verdict with no candidate_id speaks to nothing and is ignored, as it is in
// applyJudgeResults.
func answeredCandidateIDs(results []judgeResult) map[string]struct{} {
	out := make(map[string]struct{}, len(results))
	for _, result := range results {
		if id := strings.TrimSpace(result.CandidateID); id != "" {
			out[id] = struct{}{}
		}
	}
	return out
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
