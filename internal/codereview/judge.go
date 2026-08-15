package codereview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
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
	// maxJudgeFileBytes bounds ONE FILE's excerpt, deduplicated across the
	// batch. It was previously a per-candidate budget, so the same file could
	// consume it once for every candidate that named it.
	maxJudgeFileBytes = 24 * 1024
	// maxJudgeBatchBytes bounds the whole batch's file content. Dedupe and
	// windowing should keep a normal batch far under this; it exists so a
	// review touching an unusual number of large files degrades by dropping
	// excerpts it can name as dropped, rather than by building a request whose
	// size only shows up as latency.
	maxJudgeBatchBytes = 256 * 1024
	// maxJudgeDiffBytes bounds ONE file's hunks, and maxJudgeDiffBatchBytes the
	// batch's. Budgeted separately from file content rather than sharing its
	// allowance, because the two answer different questions — the diff says what
	// changed, the content says what the code now is — and a large diff must not
	// be able to starve the excerpts that verify everything else.
	maxJudgeDiffBytes      = 16 * 1024
	maxJudgeDiffBatchBytes = 128 * 1024
	// judgeWindowContextLines is how much of a file either side of an anchored
	// line the judge is shown.
	//
	// The judge's question is always local — "does this line do what the
	// finding says" — so the answer is in the lines around it, not in the first
	// 600 lines of the file, which is what a head-truncated excerpt supplied
	// regardless of where the finding pointed.
	//
	// Raised from 80 because the window, not the byte budget, is what runs out
	// first: one anchor at 80 either side is ~160 lines, roughly 6KB against a
	// maxJudgeFileBytes allowance of 24KB, so the excerpt was stopping at about
	// a quarter of what it was allowed to send. That shortfall is not free — the
	// prompt instructs the judge to answer unverifiable rather than guess when a
	// decision needs code that fell in an omitted stretch, so every line the
	// window clips converts directly into abstentions. Measured on the
	// discourse benchmark repo before this change: 10 of 72 findings (14%) came
	// back explicitly unverifiable, and none were omissions.
	//
	// 250 either side stays inside the same per-file budget for a typical
	// source file while covering the callers, guards and helpers a body claim
	// usually depends on.
	judgeWindowContextLines = 250
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

type bedrockEnhanceJudge struct {
	client *bedrockAnthropicReviewer
}

type unavailableReviewJudge struct {
	reason string
}

// judgeRequest carries file content ONCE per batch rather than once per
// candidate.
//
// Findings cluster in files — several candidates about the same file is the
// common case, not the edge case — and the per-candidate layout re-read and
// re-sent that file's bytes for every one of them. A batch of 24 candidates
// spread over 4 files sent those 4 files 24 times. The judge gained nothing
// from the copies: it is one call, reading one prompt.
//
// RepoRoot is deliberately absent. It was an absolute path on the machine
// running the review ("/Users/…/git/…"), which the model cannot act on and
// which has no business crossing into a cloud call. Every path elsewhere in
// this request is repository-relative.
type judgeRequest struct {
	ChangedFiles []string              `json:"changed_files"`
	Files        []judgeContentSnippet `json:"files"`
	// Diffs are the hunks that produced the review, for the files this batch
	// names. Without them a whole class of finding is unverifiable by
	// construction: "this value changed from X to Y" cannot be checked against
	// post-change content, which shows only Y.
	//
	// Measured on the discourse benchmark repo with whole-file excerpts already
	// in place, five of eight abstentions said exactly this — "the original
	// pre-diff value cannot be confirmed", "we only see the new file content".
	// The reviewer read these hunks to raise the findings; the judge was being
	// asked to check its work without them.
	Diffs      []judgeDiffSnippet `json:"diffs,omitempty"`
	Candidates []judgeCandidate   `json:"candidates"`
}

// judgeDiffSnippet is one file's hunks as the judge sees them.
type judgeDiffSnippet struct {
	File string `json:"file"`
	Diff string `json:"diff"`
}

// judgeCandidate names its files; the content lives in judgeRequest.Files.
type judgeCandidate struct {
	ID               string   `json:"candidate_id"`
	Title            string   `json:"title"`
	Summary          string   `json:"summary"`
	Recommendation   string   `json:"recommendation"`
	Evidence         []string `json:"evidence"`
	SourcePublishers []string `json:"source_publishers,omitempty"`
	NamedFiles       []string `json:"named_files"`
}

// judgeContentSnippet is one file's content as the judge sees it: a
// line-numbered excerpt built around the lines the batch's findings point at,
// not the head of the file. Text carries `N| ` prefixes so a claim about a
// specific line can be checked against that line.
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
	// Rank is the judge's batch-relative reading order: 1 is the candidate a
	// maintainer most needs to see. It exists because the judge's absolute
	// scores cluster — measured on 30 PRs, confirmed confidence bunches in
	// 0.8-0.9 and severity floors cost recall without buying precision — while
	// the judge is the only participant that has read every candidate in the
	// batch and can say which matters MORE. Relative orderings are the one
	// signal absolute scoring cannot Goodhart into uniformity. 0 = unranked
	// (older judge reply, or a parse that dropped it); sorting treats those as
	// last. Ranks are only comparable within one batch.
	Rank int `json:"rank"`
	// NeededFiles is how an abstaining judge asks for what it lacked: the
	// repo-relative paths (or bare file names) whose absence forced an
	// unverifiable verdict. The traced abstentions already named their missing
	// files this precisely in prose — "app/models/blocked_email.rb was not
	// provided" — so this field just makes the request machine-readable, and a
	// second round supplies the files and re-asks. Empty on any decided
	// verdict.
	NeededFiles []string `json:"needed_files"`
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
	return bedrockEnhanceJudge{client: newBedrockReviewer(plan.newTransport(), resolveBedrockJudgeModel())}
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

func (j bedrockEnhanceJudge) Available() bool {
	return j.client != nil
}

func (j bedrockEnhanceJudge) Judge(ctx context.Context, req judgeRequest) ([]judgeResult, error) {
	if j.client == nil {
		return nil, fmt.Errorf("review judge is not configured")
	}
	input := mustJSON(req)
	asked := judgeCandidateIDs(req.Candidates)
	// One retry, for malformed replies only. A batch whose reply carries no
	// results object loses every candidate in it at once, and the failure is
	// not deterministic — it is the model wandering out of its output contract,
	// which a fresh completion of the same request generally does not repeat.
	// Measured on the benchmark runs before this: roughly one batch per 10-30
	// PRs failed this way. Transport-level failures are retried a layer down
	// (see bedrockRetryBackoffs); this loop is only for a call that succeeded
	// and answered in the wrong shape.
	var lastParseErr error
	for attempt := 0; attempt < 2; attempt++ {
		completion, err := j.client.completeJSON(ctx, judgeDeveloperPrompt(), input, defaultJudgeMaxOutputTokens)
		if err != nil {
			return nil, err
		}
		results, parseErr := parseJudgeResponse(completion.Text, asked)
		if parseErr == nil {
			return results, nil
		}
		if completion.truncated() {
			// A batch cut off at the cap has no complete verdict set, only a
			// prefix of one — and the cap binds the same way on a retry, so
			// re-asking spends a full batch's tokens to reproduce the failure.
			// Report it so the batch fails open and says why.
			return nil, describeTruncatedCompletion("judge", defaultJudgeMaxOutputTokens, parseErr)
		}
		lastParseErr = parseErr
	}
	return nil, lastParseErr
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
// struct that gx Cloud's /gx/bedrock/fight normalizes, so a field only the
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
		"You are gx Review Judge. Verify candidate findings against the provided real file content, and decide which ones a human must review before merge.",
		"INPUT SHAPE: file content is in the top-level `files` array, once per file, shared by every candidate in this batch. A candidate's `named_files` lists which of those files it concerns — look them up there. A file named by a candidate but absent from `files` was not provided at all.",
		"Each file's `text` is a line-numbered excerpt, not the whole file. Every line is prefixed with its real 1-based line number and `| `, so a claim about a line number can be checked against that exact line. The excerpt is built around the lines the findings point at.",
		"A line reading `[... N line(s) omitted ...]` means that stretch of the file was not sent. Absent from the excerpt is NOT absent from the file: never treat an omitted stretch as proof that something does not exist. If deciding a candidate needs code that fell in an omitted stretch, answer unverifiable rather than guessing from what you were shown.",
		"The top-level `diffs` array holds the changes under review, as unified diff hunks, keyed by file. `files` shows what the code IS NOW; `diffs` shows what CHANGED to make it that way. A candidate whose claim is about a change — a value that moved, a line that was removed, a guard that used to be there — is decided by the diff, so read it there rather than answering unverifiable because the current content alone cannot show a before state.",
		"A file with no entry in `diffs` had no hunks provided. That is not evidence the file is unchanged; treat it the same as an omitted stretch.",
		"When your verdict is unverifiable because a SPECIFIC file, class, or template you can name was not provided — an implementation the candidate's claim depends on, a caller that would rescue the exception, the template that binds the variable — set `needed_files` to the repo-relative paths (bare file names are acceptable when you do not know the directory). They will be fetched and the candidate re-asked with them present. Use it only for that: an unverifiable verdict with an empty needed_files means no specific file would settle the claim.",
		"OUTPUT CONTRACT: reply with one JSON object and nothing else. The first character of your reply must be { and the last must be }.",
		"Do not write anything before or after that object — no preamble, no commentary, no summary — and do not wrap it in a markdown code fence.",
		"Return JSON only with shape {\"results\":[{\"candidate_id\":string,\"analysis\":string,\"verdict\":\"confirmed|unverified|wrong\",\"impact\":\"breaking|functional|cosmetic|none\",\"severity\":1-5,\"confidence\":0-1,\"verification_note\":string,\"rank\":int}]}",
		"THINK IN THE analysis FIELD. Write it first, before the verdict of the same object, and use it to do the actual work: quote the lines of the provided file content that decide the claim, state what that code really does, and name any part of the claimed mechanism the file contradicts. Then let the verdict follow from it.",
		"Take as many sentences in analysis as the candidate needs. An analysis that only restates the claim is a verdict guessed rather than checked, and a wrong confirmation costs a reviewer more than a long analysis costs you.",
		"analysis is the only place reasoning may appear. Never write it outside the JSON object. Every result object must carry a non-empty analysis; a result without one has skipped the check the field exists to force.",
		"analysis and verification_note are JSON string values, so quote code inside them with backticks and never with a raw double quote. One unescaped \" makes the whole reply undecodable and costs every candidate in this batch its verdict, not just the one you were writing about.",
		"Verdict must be confirmed, unverified, wrong, or unverifiable. confirmed survives; unverified means you checked and the claim is weak; wrong means the file content refutes it; unverifiable means you could not check it at all.",
		"Use unverifiable ONLY when the file content you needed was not provided. Do not mark such a candidate wrong or unverified: those say you checked, and gx deletes them. unverifiable surfaces the finding to the reviewer as unverified and reports the review as degraded, which is the honest outcome when the evidence never reached you. Rejecting what you could not check silently discards real findings.",
		"Confirm only when the named files and the excerpts in `files` support the title, summary, recommendation, and evidence. A finding whose conclusion may be right but whose stated mechanism the file content refutes is wrong, not confirmed — a reviewer acting on it would look for a bug that is not there.",
		"Read absolute claims literally. When a candidate says something never happens, always happens, or appears nowhere in a file, one counter-example in the provided content refutes it; look for that counter-example before you confirm, and say in the analysis whether you found one.",
		"If a candidate names no file, treat that as a strike and do not confirm unless static tool evidence conclusively proves it.",
		"impact rates the real-world consequence if the finding is acted on or ignored:",
		"- breaking: risks incorrect behavior, a crash, data loss, a security hole, or a broken build or test.",
		"- functional: affects behavior, an API or contract, or maintainability in a way a reviewer should weigh.",
		"- cosmetic: style, naming, wording, or clarity only — nothing that can break.",
		"- none: not a real issue.",
		"Reserve breaking and functional for changes a senior engineer would want to see before merge. If you are confident a finding cannot break anything, mark it cosmetic or none — gx will not surface those.",
		"Severity is 1-5. Confidence is 0-1. verification_note is one line for the reader, summarizing the analysis.",
		"rank is this batch's reading order: 1 for the candidate a maintainer most needs to see before merge, 2 for the next, and so on, one distinct rank per candidate with no ties — refuted and benign candidates go last. You have read every candidate in this batch, so rank them against EACH OTHER: which single finding matters most, which next. This is a different judgment from severity or confidence — those score each candidate alone and tend to cluster; the ordering is what decides what a reader sees first, so weigh real-world consequence, how sure you are, and how actionable the finding is.",
		"Emit exactly one result object per candidate_id you were given, in the order given. A candidate you leave out is a finding gx must ship unverified, so leave none out.",
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
	// files in first-named order, so the request is deterministic for a given
	// batch and diffable between runs.
	var ordered []string
	seen := map[string]struct{}{}
	hints := map[string][]int{}
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
		for _, file := range candidate.NamedFiles {
			if _, ok := seen[file]; !ok {
				seen[file] = struct{}{}
				ordered = append(ordered, file)
			}
		}
		for file, lines := range findingLineHints(finding) {
			hints[file] = append(hints[file], lines...)
		}
		candidates = append(candidates, candidate)
	}
	return judgeRequest{
		ChangedFiles: normalizedChangedFiles(ctx.Brief.Static.ChangedFiles),
		Files:        judgeFileContentSnippets(ctx.Brief.RepoRoot, ordered, hints),
		Diffs:        judgeDiffSnippets(ctx.Brief.Static.DiffSnippets, ordered),
		Candidates:   candidates,
	}
}

// judgeDiffSnippets picks the hunks for the files this batch names, in the same
// order as Files so the two read together.
//
// Only the batch's own files: a diff for a file no candidate mentions is spend
// with nothing to verify against, and the batch is deliberately small so that
// its request stays small.
func judgeDiffSnippets(available []DiffSnippet, ordered []string) []judgeDiffSnippet {
	if len(available) == 0 || len(ordered) == 0 {
		return nil
	}
	byFile := make(map[string]string, len(available))
	for _, snippet := range available {
		if file := normalizeReviewPath(snippet.File); file != "" {
			byFile[file] = snippet.Diff
		}
	}
	remaining := maxJudgeDiffBatchBytes
	var out []judgeDiffSnippet
	for _, file := range ordered {
		diff := strings.TrimSpace(byFile[file])
		if diff == "" {
			continue
		}
		// Truncating at a hunk boundary keeps every hunk that survives readable
		// as a diff; cutting mid-hunk would hand the judge a fragment whose
		// line numbers no longer add up.
		diff = truncateAtHunkBoundary(diff, min(maxJudgeDiffBytes, remaining))
		if strings.TrimSpace(diff) == "" {
			continue
		}
		remaining -= len(diff)
		out = append(out, judgeDiffSnippet{File: file, Diff: diff})
		if remaining <= 0 {
			break
		}
	}
	return out
}

// findingLineHints is where in each file this finding says the problem is,
// keyed by the same normalized path namedFilesForFinding produces so the two
// agree on what counts as the same file.
//
// The location was always on the finding — `file`/`line` and every anchor —
// and was thrown away at exactly the point it was useful. Reading it here is
// what lets the excerpt be built around the claim instead of around line 1.
func findingLineHints(finding Finding) map[string][]int {
	out := map[string][]int{}
	add := func(rawFile string, line int) {
		file := normalizeReviewPath(rawFile)
		if file == "" || line <= 0 {
			return
		}
		out[file] = append(out[file], line)
	}
	add(finding.File, finding.Line)
	for _, anchor := range finding.Anchors {
		add(anchor.File, anchor.Line)
	}
	return out
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
		finding.JudgeVerdict = "confirmed"
		finding.JudgeImpact = strings.TrimSpace(strings.ToLower(result.Impact))
		finding.JudgeSeverity = result.Severity
		finding.JudgeConfidence = result.Confidence
		finding.JudgeRank = result.Rank
		if note := strings.TrimSpace(result.VerificationNote); note != "" {
			finding.Evidence = append(finding.Evidence, Evidence{Label: "Judge verification", Value: note})
		}
		kept = append(kept, judged{finding: finding, result: result})
	}
	sort.SliceStable(kept, func(i, j int) bool {
		// The judge's batch-relative rank is recorded on the finding but does
		// NOT order the report. That was the plan — the one model that reads
		// every candidate side by side should out-order clustered absolute
		// scores — and it measured false: on the 30-PR benchmark, top-K by the
		// judge's stated rank scored below top-K by its own severity+confidence
		// at every K (F1 39 vs 42 at top-2, converging by top-4). The rank adds
		// noise, not signal, so it stays a recorded field for future
		// measurement and the ordering keeps the keys that won.
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
// the findings cap (see resolveMaxFindings) and said nothing about why.
//
// Re-measured on Bedrock after the judge moved there, and again after the
// analysis field landed: a real 21-candidate request now costs 7.0K to 8.3K
// output tokens, so a full batch of 24 sits at roughly 9.5K, under a third of
// defaultJudgeMaxOutputTokens. The size is kept well under what the budget
// would allow, because the batch is also the blast radius: one unparseable
// reply costs every candidate in its batch a verdict, and concurrent batches
// cost less wall clock than one call the size of all of them.
//
// 8 rather than 24, because output length is what the judge's wall clock is
// made of and batches already run concurrently. At 24 a review with 24 or
// fewer findings — the ordinary case — produced exactly ONE batch, so the
// concurrency below never engaged and the whole verification was one serial
// call emitting ~9.5K tokens. Splitting the same findings three ways runs
// three calls at once, each answering for a third as many candidates. On the
// numbers above the per-call output falls to roughly a third, and a batch
// carries almost no shared context (changed_files and the files its own
// candidates name), so the split duplicates very little.
const defaultJudgeBatchSize = 8

// maxConcurrentJudgeBatches caps how many judge calls are in flight at once.
//
// Smaller batches mean more of them, and the transport this runs over has no
// retry (see cloudBedrockTransport.complete), so a throttle is not a slow
// batch, it is a lost one. This bounds the fan-out a large review can aim at
// Bedrock while still leaving the ordinary review fully parallel.
const maxConcurrentJudgeBatches = 6

// resolveJudgeBatchSize applies env > default, so the size can be tuned
// against a real repository without a rebuild.
func resolveJudgeBatchSize() int {
	raw := strings.TrimSpace(os.Getenv("GX_REVIEW_JUDGE_BATCH_SIZE"))
	if raw == "" {
		return defaultJudgeBatchSize
	}
	size, err := strconv.Atoi(raw)
	if err != nil || size <= 0 {
		return defaultJudgeBatchSize
	}
	return size
}

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
	// Omitted and Abstained split what used to be one "unanswered" count.
	//
	// Both end up unverified, but they are opposite problems and the fix for one
	// is wasted effort on the other. Omitted means the reply never mentioned the
	// candidate — a reliability failure, and asking again usually answers it.
	// Abstained means the judge answered "I cannot check this" (see
	// judgeCouldNotVerify) — an honest verdict, and asking again just buys the
	// same answer at twice the price; what that case needs is better evidence,
	// not another call.
	//
	// Counting them together made a review that could not reach its evidence
	// look identical to a model dropping list entries, which is why this is
	// reported as two numbers.
	Omitted   int
	Abstained int
	// AbstentionNotes is what the judge said when it declined to decide, one
	// entry per abstaining candidate.
	//
	// "The judge could not verify this" is a useless fact on its own: it could
	// mean the excerpt clipped the code that decides the claim, that the claim
	// is about a change and only the post-change file was sent, or that the
	// claim is architectural and no file content could settle it. Those want
	// three different fixes, and the reasoning that distinguishes them is
	// already in judgeResult.Analysis — which nothing renders, so the answer
	// was being computed and dropped on every run.
	//
	// Collected unconditionally (they are strings already in memory) and
	// surfaced only under GX_REVIEW_JUDGE_TRACE, because an ordinary review
	// should not carry the judge's working out in its degraded reasons.
	AbstentionNotes []string
}

// Unanswered is every candidate a successful batch left without a usable
// verdict, by either route. Kept so callers that only care "was this review
// fully judged" do not have to know the difference.
func (o judgeBatchOutcome) Unanswered() int { return o.Omitted + o.Abstained }

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
	batchSize := resolveJudgeBatchSize()
	var batches [][]Finding
	for start := 0; start < len(candidates); start += batchSize {
		end := start + batchSize
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
	slots := make(chan struct{}, maxConcurrentJudgeBatches)
	var wg sync.WaitGroup
	for i, batch := range batches {
		i, batch := i, batch
		wg.Add(1)
		go func() {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()
			results, err := judge.Judge(ctx, buildJudgeRequest(reviewContext, batch))
			out[i] = batchResult{results: results, err: err}
		}()
	}
	wg.Wait()

	outcome := judgeBatchOutcome{Batches: len(batches)}
	var retriable []retriableAbstention
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
		mentioned := mentionedCandidateIDs(out[i].results)
		byID := map[string]judgeResult{}
		for _, result := range out[i].results {
			if id := strings.TrimSpace(result.CandidateID); id != "" {
				byID[id] = result
			}
		}
		var judged, silent []Finding
		for _, finding := range batch {
			id := strings.TrimSpace(finding.ID)
			if _, ok := verdicts[id]; ok {
				judged = append(judged, finding)
				continue
			}
			// Both are unverified, but only one is a failure. A candidate the
			// reply named and declined to decide was answered; one it never
			// named was dropped.
			if _, seen := mentioned[id]; seen {
				outcome.Abstained++
				outcome.AbstentionNotes = append(outcome.AbstentionNotes,
					abstentionNote(finding, out[i].results))
				if needs := byID[id].NeededFiles; len(needs) > 0 {
					retriable = append(retriable, retriableAbstention{finding: finding, needs: needs})
					continue
				}
			} else {
				outcome.Omitted++
			}
			silent = append(silent, finding)
		}
		outcome.Judged = append(outcome.Judged, applyJudgeResults(judged, out[i].results)...)
		outcome.Unjudged = append(outcome.Unjudged, silent...)
	}
	rejudgeWithRequestedFiles(ctx, judge, reviewContext, retriable, &outcome)
	return outcome
}

// retriableAbstention is an abstention the judge itself said how to fix: the
// candidate plus the files whose absence forced the unverifiable verdict.
type retriableAbstention struct {
	finding Finding
	needs   []string
}

// maxJudgeRequestedFiles bounds how many judge-requested files one second
// round fetches. The traced abstentions each named one or two files; a reply
// requesting many more is fishing, not verifying.
const maxJudgeRequestedFiles = 8

// rejudgeWithRequestedFiles runs one supplemental round for abstentions that
// named their missing files.
//
// The first round's abstentions are precise about what they lack — "the view
// file was not provided", naming it — and every named file sits in the local
// checkout the review is already reading. Fetching it and re-asking converts
// an unverifiable verdict into a real one in a single extra call: either a
// confirmation with evidence, or — the case that pays for precision — a
// refutation of a claim the judge previously had to wave through unverified.
// One round only; a candidate still unverifiable with the files it asked for
// stays unverified, and the failure mode is the status quo.
func rejudgeWithRequestedFiles(ctx context.Context, judge FindingJudge,
	reviewContext ReviewContext, retriable []retriableAbstention, outcome *judgeBatchOutcome) {
	if len(retriable) == 0 {
		return
	}
	repoRoot := reviewContext.Brief.RepoRoot
	var findings []Finding
	requested := map[string]struct{}{}
	for _, r := range retriable {
		findings = append(findings, r.finding)
		for _, need := range r.needs {
			if len(requested) >= maxJudgeRequestedFiles {
				break
			}
			if resolved := resolveRequestedFile(ctx, repoRoot, need); resolved != "" {
				requested[resolved] = struct{}{}
			}
		}
	}
	fallBack := func() {
		outcome.Unjudged = append(outcome.Unjudged, findings...)
	}
	if len(requested) == 0 {
		fallBack()
		return
	}
	request := buildJudgeRequest(reviewContext, findings)
	present := map[string]struct{}{}
	for _, snippet := range request.Files {
		present[snippet.File] = struct{}{}
	}
	hints := map[string][]int{}
	var extra []string
	for file := range requested {
		if _, ok := present[file]; !ok {
			extra = append(extra, file)
		}
	}
	sort.Strings(extra)
	request.Files = append(request.Files, judgeFileContentSnippets(repoRoot, extra, hints)...)

	results, err := judge.Judge(ctx, request)
	if err != nil {
		fallBack()
		return
	}
	verdicts := answeredCandidateIDs(results)
	var judged, still []Finding
	for _, finding := range findings {
		if _, ok := verdicts[strings.TrimSpace(finding.ID)]; ok {
			judged = append(judged, finding)
			continue
		}
		still = append(still, finding)
	}
	answered := applyJudgeResults(judged, results)
	outcome.Judged = append(outcome.Judged, answered...)
	outcome.Unjudged = append(outcome.Unjudged, still...)
	// The candidates that got real verdicts this round are no longer
	// abstentions; the rest stay counted from round one.
	outcome.Abstained -= len(judged)
	if outcome.Abstained < 0 {
		outcome.Abstained = 0
	}
}

// resolveRequestedFile turns a judge-named file into a repo-relative path it is
// safe to read: inside the repository, existing, and matched by exact path
// first, then by unique basename via the git index. The judge's request is
// model output — a path that escapes the root or matches nothing is dropped,
// never guessed at.
func resolveRequestedFile(ctx context.Context, repoRoot, need string) string {
	need = normalizeReviewPath(strings.TrimSpace(need))
	if need == "" || strings.Contains(need, "..") || path.IsAbs(need) {
		return ""
	}
	if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(need))); err == nil {
		return need
	}
	// Bare or wrong-directory name: let the git index find it, and accept the
	// match only when it is unambiguous.
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "ls-files", "--", "*/"+path.Base(need), path.Base(need))
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var matches []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			matches = append(matches, line)
		}
	}
	if len(matches) == 1 {
		return normalizeReviewPath(matches[0])
	}
	return ""
}

// answeredCandidateIDs is the set of candidates a verdict set speaks to. A
// verdict with no candidate_id speaks to nothing and is ignored, as it is in
// applyJudgeResults. A verdict that says the judge could not check the claim is
// likewise not an answer — see judgeCouldNotVerify.
func answeredCandidateIDs(results []judgeResult) map[string]struct{} {
	out := make(map[string]struct{}, len(results))
	for _, result := range results {
		if judgeCouldNotVerify(result) {
			continue
		}
		if id := strings.TrimSpace(result.CandidateID); id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

// maxAbstentionNoteChars keeps one note readable in a degraded-reasons list.
// The judge's analysis runs to paragraphs; what identifies the cause is the
// first sentence or two, and the rest restates the finding.
const maxAbstentionNoteChars = 400

// abstentionNote pairs a declined candidate with what the judge said about it,
// so the reason can be read rather than inferred from a count.
func abstentionNote(finding Finding, results []judgeResult) string {
	id := strings.TrimSpace(finding.ID)
	for _, result := range results {
		if strings.TrimSpace(result.CandidateID) != id {
			continue
		}
		// VerificationNote is the judge's own summary of why it stopped; the
		// analysis is the longer reasoning behind it. Prefer the summary and
		// fall back, because either one answers the question and neither is
		// guaranteed to be filled in.
		reason := strings.TrimSpace(result.VerificationNote)
		if reason == "" {
			reason = strings.TrimSpace(result.Analysis)
		}
		if reason == "" {
			reason = "(no reason given)"
		}
		if len(reason) > maxAbstentionNoteChars {
			reason = reason[:maxAbstentionNoteChars] + "…"
		}
		return fmt.Sprintf("%s [%s] %q: %s",
			id, strings.TrimSpace(result.Verdict), strings.TrimSpace(finding.Title), reason)
	}
	return fmt.Sprintf("%s %q: (verdict not found in reply)", id, strings.TrimSpace(finding.Title))
}

// mentionedCandidateIDs is every candidate the reply named at all, including the
// ones it declined to decide. Set-differenced against answeredCandidateIDs it
// separates "the judge abstained" from "the judge never mentioned it", which is
// the difference between an evidence problem and a reliability one.
func mentionedCandidateIDs(results []judgeResult) map[string]struct{} {
	out := make(map[string]struct{}, len(results))
	for _, result := range results {
		if id := strings.TrimSpace(result.CandidateID); id != "" {
			out[id] = struct{}{}
		}
	}
	return out
}

// judgeCouldNotVerify reports a verdict that declines to decide.
//
// "I could not check this" and "I checked this and it is wrong" are different
// answers, and applyJudgeResults cannot tell them apart: it keeps `confirmed`
// and drops everything else, so an honest abstention deletes a real finding as
// efficiently as a refutation does. Routing abstentions here sends them to the
// unjudged fallback instead, where they surface unverified and are counted in
// the degraded reasons — the same treatment as a candidate the judge never
// mentioned, which is exactly what an abstention is.
//
// This matters most in the case that motivated it: a judge given no file
// content cannot verify anything, and the failure is silent because rejecting
// every candidate looks identical in the output to a review that found nothing.
func judgeCouldNotVerify(result judgeResult) bool {
	switch strings.ToLower(strings.TrimSpace(result.Verdict)) {
	case "unverifiable", "insufficient_evidence", "insufficient evidence",
		"cannot_verify", "cannot verify", "abstain", "indeterminate":
		return true
	}
	// Deliberately NOT "unverified": in this contract that means "I checked and
	// the claim is weak", which is a real judgement and should drop the finding.
	// Conflating it with abstention would surface every weak claim the judge
	// correctly filtered out.
	return false
}

func capAdvisoryFindings(findings []Finding, limit int) []Finding {
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
	limit = resolveMaxFindings(limit)
	if len(out) > limit {
		out = out[:limit]
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

// namedFilesForFinding is what decides whether the judge gets any source to
// verify against, so it reads the finding's STRUCTURED location first and only
// then falls back to scraping paths out of its prose.
//
// Prose-only was the original implementation and it silently starved the judge.
// A finding carries `file`/`line` and `anchors` precisely so the location does
// not have to be restated in a sentence, and a model that fills those fields
// and writes "the ownership check is missing" — correct, idiomatic output —
// produced no named files, hence no file content, hence a judge asked to verify
// a claim about source it was never shown. It then answers in prose instead of
// its verdict block and the whole batch loses its verdicts. The symptom is a
// review that reports nothing while reading as complete, which is the exact
// failure --fail-on exists to prevent.
func namedFilesForFinding(ctx ReviewContext, finding Finding) []string {
	known := knownReviewFiles(ctx)
	changed := normalizedChangedFiles(ctx.Brief.Static.ChangedFiles)
	seen := map[string]struct{}{}
	var parsed []string
	// Structured fields first: they are unambiguous, and they are what the
	// finding schema asks the model to fill in.
	for _, candidate := range structuredFindingFiles(finding) {
		file := normalizeReviewPath(candidate)
		if file == "" {
			continue
		}
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		parsed = append(parsed, file)
	}
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

// judgeFileContentSnippets reads each file ONCE for the whole batch and returns
// the excerpt the judge is shown, in the order the files were first named.
//
// files must already be deduplicated; hints maps a file to the lines the
// batch's findings point at, which is what decides which part of it is worth
// sending.
func judgeFileContentSnippets(repoRoot string, files []string, hints map[string][]int) []judgeContentSnippet {
	if strings.TrimSpace(repoRoot) == "" {
		return nil
	}
	remaining := maxJudgeBatchBytes
	var out []judgeContentSnippet
	for _, file := range files {
		if remaining <= 0 {
			break
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(file)))
		if err != nil {
			continue
		}
		text := judgeFileExcerpt(string(data), hints[file], min(maxJudgeFileBytes, remaining))
		if text == "" {
			continue
		}
		remaining -= len(text)
		out = append(out, judgeContentSnippet{File: file, Text: text})
	}
	return out
}

// judgeLineRange is an inclusive, 1-based span of a file.
type judgeLineRange struct{ start, end int }

// fittedLineRanges picks the widest window whose spans still fit the budget.
//
// Without this the width is fixed and the budget silently decides how much of
// each span survives — which, because spans render from their start, means it
// decides whether the anchored line is reached at all. Halving until it fits
// keeps the anchor and trades away only surrounding context, which is the right
// direction: context around a line the judge cannot see is worth nothing.
//
// The floor is not zero. A window of 0 is still one line — the anchor itself —
// and an excerpt showing only the lines under dispute beats one showing their
// neighbours instead.
func fittedLineRanges(fileLines []string, lines []int, limit int) []judgeLineRange {
	window := judgeWindowContextLines
	for {
		spans := judgeLineRangesWithWindow(lines, len(fileLines), window)
		if window == 0 || renderedSpansSize(fileLines, spans) <= limit {
			return spans
		}
		window /= 2
	}
}

// renderedSpansSize is what these spans cost once numbered, omission markers
// included — they are content and have carried an excerpt past its cap before.
func renderedSpansSize(fileLines []string, spans []judgeLineRange) int {
	size, shown := 0, 0
	for _, span := range spans {
		if span.start > shown+1 {
			size += len(fmt.Sprintf("[... %d line(s) omitted ...]\n", span.start-shown-1))
		}
		for n := span.start; n <= span.end; n++ {
			size += len(strconv.Itoa(n)) + len("| ") + len(fileLines[n-1]) + len("\n")
		}
		shown = max(shown, span.end)
	}
	if shown < len(fileLines) {
		size += len(fmt.Sprintf("[... %d line(s) omitted ...]\n", len(fileLines)-shown))
	}
	return size
}

// renderedLinesSize is how many bytes these lines occupy once numbered, which
// is what the budget is actually spent on — the `N| ` prefix is not free, and
// on a file of short lines it is a noticeable share of the total.
func renderedLinesSize(fileLines []string) int {
	size := 0
	for n, line := range fileLines {
		size += len(strconv.Itoa(n+1)) + len("| ") + len(line) + len("\n")
	}
	return size
}

// judgeLineRanges turns the lines a batch's findings pointed at into merged,
// in-bounds spans to show.
//
// With no usable line the whole file is one span, which the caller's budget
// then truncates — the old head-of-file behavior, kept only for the findings
// that genuinely name no line. A line past the end of the file counts as no
// line: a claim about line 5000 of a 100-line file is one the judge should see
// the whole file to reject.
func judgeLineRanges(lines []int, total int) []judgeLineRange {
	return judgeLineRangesWithWindow(lines, total, judgeWindowContextLines)
}

// judgeLineRangesWithWindow is judgeLineRanges at a caller-chosen width, so an
// excerpt can narrow its windows to fit a budget instead of being cut off.
//
// A span is centered on its anchor but rendered from its start, so a window
// wider than the budget can afford spends the whole allowance on the lines
// BEFORE the anchor and never reaches it — the excerpt drops precisely the line
// it exists to show. Widening the window made that reachable in practice, which
// is why the width is now chosen against the budget rather than fixed.
func judgeLineRangesWithWindow(lines []int, total, window int) []judgeLineRange {
	var anchors []int
	for _, line := range lines {
		if line > 0 && line <= total {
			anchors = append(anchors, line)
		}
	}
	if len(anchors) == 0 {
		return []judgeLineRange{{start: 1, end: total}}
	}
	sort.Ints(anchors)
	var out []judgeLineRange
	for _, line := range anchors {
		next := judgeLineRange{
			start: max(1, line-window),
			end:   min(total, line+window),
		}
		// Merge windows that touch or overlap, so two findings a few lines
		// apart produce one span rather than two copies of the same code.
		if n := len(out); n > 0 && next.start <= out[n-1].end+1 {
			out[n-1].end = max(out[n-1].end, next.end)
			continue
		}
		out = append(out, next)
	}
	return out
}

// judgeFileExcerpt renders the spans of one file worth showing, with 1-based
// line numbers.
//
// The numbers are the point: the judge is checking a claim about a location, so
// it has to be able to tell which line it is looking at. Omitted stretches are
// marked rather than silently spliced, so the model can tell "not in the file"
// apart from "not in the excerpt" — the distinction its unverifiable verdict
// depends on.
func judgeFileExcerpt(content string, lines []int, budget int) string {
	if budget <= 0 {
		return ""
	}
	fileLines := strings.Split(content, "\n")
	// A trailing newline splits into a final empty element that is not a line.
	if n := len(fileLines); n > 0 && fileLines[n-1] == "" {
		fileLines = fileLines[:n-1]
	}
	if len(fileLines) == 0 {
		return ""
	}
	const truncationNote = "[... truncated at the size limit ...]\n"
	// Every write goes through the budget, including the omission markers.
	// They are content: a file with many anchors emits many of them, and
	// checking only the source lines let an excerpt drift past its cap by their
	// combined size. The limit also holds back room for the closing note, so
	// the result honors the budget whether it ends by running out of file or by
	// running out of room.
	limit := budget - len(truncationNote)
	if limit <= 0 {
		return ""
	}
	var b strings.Builder
	shown := 0
	truncated := false
	write := func(text string) bool {
		if b.Len()+len(text) > limit {
			truncated = true
			return false
		}
		b.WriteString(text)
		return true
	}
	// Whole file when it fits, windows only when it does not.
	//
	// An omitted stretch is not a neutral saving: the prompt tells the judge to
	// answer unverifiable rather than guess when a decision needs code that
	// fell in one, so every marker is a potential abstention. Measured on the
	// discourse benchmark repo, half the abstentions cited exactly that — "in
	// omitted lines 365+", "may be in the omitted portion", "truncated before
	// that key appears" — while the excerpts themselves were running at roughly
	// a quarter of maxJudgeFileBytes. Spending the unused budget removes the
	// question rather than widening the window and hoping it now reaches.
	//
	// Windowing still earns its keep on files too big to send, which is where
	// it was doing real work all along.
	spans := fittedLineRanges(fileLines, lines, limit)
	if renderedLinesSize(fileLines) <= limit {
		spans = []judgeLineRange{{start: 1, end: len(fileLines)}}
	}
	for _, span := range spans {
		if span.start > shown+1 {
			if !write(fmt.Sprintf("[... %d line(s) omitted ...]\n", span.start-shown-1)) {
				break
			}
		}
		for n := span.start; n <= span.end; n++ {
			if !write(fmt.Sprintf("%d| %s\n", n, fileLines[n-1])) {
				break
			}
			shown = n
		}
		if truncated {
			break
		}
	}
	switch {
	case truncated:
		b.WriteString(truncationNote)
	case shown < len(fileLines):
		if !write(fmt.Sprintf("[... %d line(s) omitted ...]\n", len(fileLines)-shown)) {
			b.WriteString(truncationNote)
		}
	}
	return b.String()
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

// structuredFindingFiles lists the file paths a finding states outright, in the
// order a reader would trust them: the finding's own File field, then every
// anchor. These are deliberately NOT filtered against the known-file set the
// way prose paths are — prose needs that filter because scraping sentences
// yields false positives, whereas a path in a structured field is a claim the
// reviewer made on purpose. Filtering it would reintroduce the starvation this
// exists to fix whenever the repository scan is incomplete.
func structuredFindingFiles(finding Finding) []string {
	var out []string
	if file := strings.TrimSpace(finding.File); file != "" {
		out = append(out, file)
	}
	for _, anchor := range finding.Anchors {
		if file := strings.TrimSpace(anchor.File); file != "" {
			out = append(out, file)
		}
	}
	return out
}
