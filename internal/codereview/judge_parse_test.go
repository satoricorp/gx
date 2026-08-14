package codereview

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stubBedrockTransport returns a canned model reply, so a test can drive the
// whole judge path — request, parse, verdict — over bytes a real model actually
// produced.
type stubBedrockTransport struct {
	completion bedrockCompletion
	err        error
}

func (t stubBedrockTransport) complete(context.Context, string, string, string, int) (bedrockCompletion, error) {
	return t.completion, t.err
}

func (t stubBedrockTransport) detail() string { return "stub transport" }

func judgeFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}

func stubJudge(completion bedrockCompletion) bedrockReviewJudge {
	return bedrockReviewJudge{client: newBedrockReviewer(stubBedrockTransport{completion: completion}, "stub-model")}
}

// prosePreambleCandidateIDs is the candidate set the captured reply was actually
// answering, in the order it was asked. The parser is given the question as well
// as the answer, so the fixture tests have to carry it too.
var prosePreambleCandidateIDs = []string{
	"bedrock-a.ai.review.1.s2",
	"bedrock-a.ai.review.1.s3",
	"bedrock-a.ai.review.2.s1",
	"bedrock-a.ai.review.2.s3",
	"bedrock-a.ai.review.3.s2",
	"bedrock-a.ai.review.3.s3",
	"bedrock-a.ai.review.4.s1",
	"bedrock-a.ai.review.4.s2",
	"bedrock-a.ai.review.4.s3",
	"bedrock-a.ai.review.5.s1",
	"bedrock-a.ai.review.5.s2",
	"bedrock-a.ai.review.5.s3",
	"bedrock-b.ai.review.1.s2",
	"bedrock-b.ai.review.2.s1",
	"bedrock-b.ai.review.2.s2",
	"bedrock-b.ai.review.2.s3",
	"bedrock-b.ai.review.3.s1",
	"bedrock-b.ai.review.4.s1",
	"bedrock-a.ai.review.2.s1.pad23",
	"bedrock-a.ai.review.1.s2.pad21",
	"bedrock-a.ai.review.1.s3.pad22",
	"testing.no-tests",
	"dependencies.no-manifest",
	"onboarding.missing-agents",
}

func prosePreambleRequest() judgeRequest {
	candidates := make([]judgeCandidate, 0, len(prosePreambleCandidateIDs))
	for _, id := range prosePreambleCandidateIDs {
		candidates = append(candidates, judgeCandidate{ID: id})
	}
	return judgeRequest{Candidates: candidates}
}

// TestParseJudgeResponseReadsRealProsePreambleReply is the regression for the
// live degraded review. The fixture is not synthetic: it is the verbatim reply
// captured from us.anthropic.claude-sonnet-4-6 when a real 24-candidate batch
// from /Users/joe/git/yeet was replayed against it.
//
// The model answered with 5,439 characters of reasoning and then a ```json
// fence. Its JSON was complete and correct. What broke was the parser: the
// prose quoted source code containing braces, so slicing from the first "{" to
// the last "}" started mid-sentence and json.Unmarshal reported "invalid
// character 't' looking for beginning of object key string" — the same failure,
// down to the shape of the message, that a real review reported with 'r'.
func TestParseJudgeResponseReadsRealProsePreambleReply(t *testing.T) {
	content := judgeFixture(t, "judge-prose-preamble-response.txt")

	// Guard the fixture itself: if these stop being true the test no longer
	// covers the bug it was written for.
	if strings.HasPrefix(strings.TrimSpace(content), "{") {
		t.Fatal("fixture no longer starts with a prose preamble")
	}
	firstBrace := strings.IndexByte(content, '{')
	if firstBrace < 0 || !strings.HasPrefix(content[firstBrace:], "{ try {") {
		t.Fatalf("fixture's first brace is no longer the brace-bearing code quote in the prose: %q",
			content[firstBrace:min(firstBrace+24, len(content))])
	}

	results, err := parseJudgeResponse(content, prosePreambleCandidateIDs)
	if err != nil {
		t.Fatalf("parseJudgeResponse() error = %v, want the fenced object read past the prose", err)
	}
	if len(results) != 24 {
		t.Fatalf("parseJudgeResponse() returned %d results, want all 24 verdicts the model actually sent", len(results))
	}
	if results[0].CandidateID != "bedrock-a.ai.review.1.s2" || results[0].Verdict != "confirmed" {
		t.Fatalf("first result = %+v, want the model's real first verdict", results[0])
	}
	for i, result := range results {
		if strings.TrimSpace(result.CandidateID) == "" {
			t.Fatalf("result %d has no candidate_id, so the wrong object was parsed: %+v", i, result)
		}
	}
}

// TestParseJudgeResponseIgnoresTheModelQuotingItsOwnOutputShape pins the fix for
// the extractor hijack.
//
// "First decodable object carrying a results key" is not a safe rule for this
// prompt, because the prompt prints the literal shape {"results":[...]} to the
// model — restating it in a preamble is the most natural thing a model could do
// with it. The quoted object decodes perfectly and answers nothing, and because
// applyJudgeResults is delete-only, adopting it drops every candidate in the
// batch with BatchesFailed=0, so no degraded line is ever written. Measured
// before the fix: 0 verdicts, nil error.
func TestParseJudgeResponseIgnoresTheModelQuotingItsOwnOutputShape(t *testing.T) {
	asked := []string{"bedrock-a.ai.review.1", "bedrock-b.ai.review.2"}
	reply := strings.Join([]string{
		`I will answer with {"results":[]} if none of the candidates check out.`,
		`Having reviewed them, two are real.`,
		`{"results":[` +
			`{"candidate_id":"bedrock-a.ai.review.1","analysis":"the file confirms it","verdict":"confirmed","impact":"breaking","severity":4,"confidence":0.9,"verification_note":"real"},` +
			`{"candidate_id":"bedrock-b.ai.review.2","analysis":"the file confirms it","verdict":"confirmed","impact":"breaking","severity":4,"confidence":0.9,"verification_note":"real"}` +
			`]}`,
	}, "\n")

	results, err := parseJudgeResponse(reply, asked)
	if err != nil {
		t.Fatalf("parseJudgeResponse() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("parseJudgeResponse() returned %d results, want the 2 real verdicts rather than the empty set quoted in the preamble", len(results))
	}
	for i, want := range asked {
		if results[i].CandidateID != want {
			t.Fatalf("result %d = %q, want %q", i, results[i].CandidateID, want)
		}
	}
}

// TestParseJudgeResponseRejectsAReplyWithNoVerdictsInIt covers the silent-drop
// route that the "if the whole reply parses, take it" fast path used to open:
// json.Unmarshal of a well-formed object with no results field succeeds, yields
// no verdicts, and reports no error — which is the same batch-wide silent drop
// as the hijack, reached without any preamble at all.
func TestParseJudgeResponseRejectsAReplyWithNoVerdictsInIt(t *testing.T) {
	results, err := parseJudgeResponse(`{"overview":"all of these look fine to me"}`, []string{"a", "b"})
	if err == nil {
		t.Fatalf("parseJudgeResponse() error = nil with %d results, want a reply carrying no verdict set to fail rather than silently confirm nothing", len(results))
	}
	if len(results) != 0 {
		t.Fatalf("parseJudgeResponse() returned %d results, want none", len(results))
	}
}

// TestParseJudgeResponseRejectsRealTruncatedReply pins the other half: a reply
// cut off at the output cap must fail, not half-succeed. The fixture is again a
// real capture — the same batch replayed with the cap lowered to 1200 tokens —
// and it stops mid-object with six complete verdicts already written.
//
// Returning those six would be the worst available outcome: applyJudgeResults
// is delete-only, so every candidate the model never reached would be dropped
// as unconfirmed, silently, by a batch that looked like it succeeded.
func TestParseJudgeResponseRejectsRealTruncatedReply(t *testing.T) {
	content := judgeFixture(t, "judge-truncated-response.txt")
	if strings.Contains(content, "\"onboarding.missing-agents\"") {
		t.Fatal("fixture is no longer truncated: it reached the last candidate")
	}

	results, err := parseJudgeResponse(content, prosePreambleCandidateIDs)
	if err == nil {
		t.Fatalf("parseJudgeResponse() error = nil with %d results, want a truncated reply to fail rather than yield a partial verdict set", len(results))
	}
	if len(results) != 0 {
		t.Fatalf("parseJudgeResponse() returned %d partial results, want none", len(results))
	}
}

// TestJudgeReportsTruncationAsTheCause pins that the failure a truncated batch
// reports names the output cap. "invalid character" reads like a model that
// cannot format JSON; it is really a budget that needs raising, and the two
// have different fixes.
func TestJudgeReportsTruncationAsTheCause(t *testing.T) {
	judge := stubJudge(bedrockCompletion{
		Text:       judgeFixture(t, "judge-truncated-response.txt"),
		StopReason: "max_tokens",
	})

	results, err := judge.Judge(context.Background(), prosePreambleRequest())
	if err == nil {
		t.Fatal("Judge() error = nil, want a truncated reply reported as a failure")
	}
	if len(results) != 0 {
		t.Fatalf("Judge() returned %d results alongside the error, want none", len(results))
	}
	if !strings.Contains(err.Error(), "output cap") {
		t.Fatalf("Judge() error = %v, want it to name the output cap as the cause", err)
	}
}

// TestJudgeReadsRealProsePreambleReplyEndToEnd runs the captured reply through
// the judge itself rather than the parser alone, so the fix is pinned at the
// boundary the engine actually calls.
func TestJudgeReadsRealProsePreambleReplyEndToEnd(t *testing.T) {
	judge := stubJudge(bedrockCompletion{
		Text:       judgeFixture(t, "judge-prose-preamble-response.txt"),
		StopReason: "end_turn",
	})

	results, err := judge.Judge(context.Background(), prosePreambleRequest())
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if len(results) != 24 {
		t.Fatalf("Judge() returned %d results, want 24", len(results))
	}
}

// TestUnparseableJudgeBatchKeepsItsFindings is the fail-open property, driven by
// a reply that genuinely cannot be parsed. A judge that cannot answer must not
// be able to delete a finding — "the judge could not be reached" and "the judge
// rejected this" are different facts and must not produce the same review.
func TestUnparseableJudgeBatchKeepsItsFindings(t *testing.T) {
	judge := stubJudge(bedrockCompletion{Text: "I could not evaluate these candidates.", StopReason: "end_turn"})
	candidates := judgeBatchCandidates(3)

	outcome := runJudge(context.Background(), judge, ReviewContext{}, candidates)

	if outcome.BatchesFailed != 1 || outcome.Batches != 1 {
		t.Fatalf("BatchesFailed/Batches = %d/%d, want 1/1", outcome.BatchesFailed, outcome.Batches)
	}
	if len(outcome.Judged) != 0 {
		t.Fatalf("Judged = %d, want no verdicts from an unparseable reply", len(outcome.Judged))
	}
	if len(outcome.Unjudged) != len(candidates) {
		t.Fatalf("Unjudged = %d, want all %d candidates kept rather than dropped", len(outcome.Unjudged), len(candidates))
	}
	if outcome.Err == nil {
		t.Fatal("Err is nil, want the parse failure carried to the degraded-reasons line")
	}
}

// TestEnhanceReportsTruncatedVerificationAsDegraded closes the loop: the reason a
// batch failed has to reach the human reading the review, or the review reads as
// a clean one that simply found nothing.
func TestEnhanceReportsTruncatedVerificationAsDegraded(t *testing.T) {
	t.Setenv("GX_REVIEW_JUDGE", "1")
	root := t.TempDir()
	writeFile(t, root, "internal/app/app.go", "package app\nfunc Run() {}\n")
	engine := judgeTestEngine(root, []Finding{
		judgeTestFinding("bedrock-a.ai.review.1", "internal/app/app.go"),
	})
	engine.judge = stubJudge(bedrockCompletion{
		Text:       judgeFixture(t, "judge-truncated-response.txt"),
		StopReason: "max_tokens",
	})

	report, err := engine.Review(context.Background(), root, Options{Scope: "architecture"})
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}
	if len(report.Findings) == 0 {
		t.Fatal("Findings is empty, want the unverified candidate kept when verification fails")
	}
	var reason string
	for _, candidate := range report.DegradedReasons {
		if strings.Contains(candidate, "verification") {
			reason = candidate
			break
		}
	}
	if reason == "" {
		t.Fatalf("DegradedReasons = %#v, want the failed verification reported", report.DegradedReasons)
	}
	if !strings.Contains(reason, "output cap") {
		t.Fatalf("degraded reason = %q, want it to name the output cap", reason)
	}
}

// TestExtractJSONObjectSkipsProseAndPicksTheAnsweringObject covers the extractor
// used by the reviewer legs, which have no identifiers to match on and so rank
// by how much of the expected reply shape an object carries.
func TestExtractJSONObjectSkipsProseAndPicksTheAnsweringObject(t *testing.T) {
	tests := []struct {
		name    string
		content string
		keys    []string
		want    string
	}{
		{
			name:    "brace-bearing code quote before the answer",
			content: "Looking at it: `if (ok) { return null }` is fine.\n```json\n{\"recommendations\":[]}\n```",
			keys:    []string{"recommendations"},
			want:    `{"recommendations":[]}`,
		},
		{
			name:    "well-formed but wrong object before the answer",
			content: "The shape is {\"example\":true} — here it is: {\"recommendations\":[{\"title\":\"a\"}]}",
			keys:    []string{"recommendations"},
			want:    `{"recommendations":[{"title":"a"}]}`,
		},
		{
			name:    "a one-key quote of the shape loses to the real reply",
			content: "I will answer with {\"recommendations\":[]} if I find nothing.\n{\"overview\":\"ok\",\"recommendations\":[{\"title\":\"a\"}]}",
			keys:    []string{"recommendations", "overview", "notable_changes"},
			want:    `{"overview":"ok","recommendations":[{"title":"a"}]}`,
		},
		{
			// An answer wrapped in an object of its own must still be found.
			// Skipping a non-scoring object whole — rather than stepping into
			// it — turned this into a decode failure, and a decode failure on a
			// reviewer leg is a review that silently produces no AI findings.
			name:    "answer nested inside a wrapper object",
			content: `{"result":{"recommendations":[{"title":"a"}],"overview":"ok"}}`,
			keys:    []string{"recommendations", "overview"},
			want:    `{"recommendations":[{"title":"a"}],"overview":"ok"}`,
		},
		{
			name:    "brace inside a string value ends nothing",
			content: "{\"recommendations\":[{\"summary\":\"code reads } here\"}]}",
			keys:    []string{"recommendations"},
			want:    "{\"recommendations\":[{\"summary\":\"code reads } here\"}]}",
		},
		{
			name:    "trailing prose after the answer",
			content: "{\"recommendations\":[]}\n\nLet me know if you want more detail.",
			keys:    []string{"recommendations"},
			want:    `{"recommendations":[]}`,
		},
		{
			name:    "truncated object yields nothing",
			content: "```json\n{\"recommendations\":[{\"title\":\"a\"},{\"title\"",
			keys:    []string{"recommendations"},
			want:    "",
		},
		{
			name:    "no required key yields nothing",
			content: `{"unrelated":"all good"}`,
			keys:    []string{"recommendations"},
			want:    "",
		},
		{
			name:    "any of several keys is enough",
			content: "Here you go:\n{\"overview\":\"all good\"}",
			keys:    []string{"recommendations", "overview"},
			want:    `{"overview":"all good"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := extractJSONObject(test.content, test.keys...)
			if got != test.want {
				t.Fatalf("extractJSONObject() = %q, want %q", got, test.want)
			}
			if got == "" {
				return
			}
			var probe map[string]json.RawMessage
			if err := json.Unmarshal([]byte(got), &probe); err != nil {
				t.Fatalf("extractJSONObject() returned %q, which does not parse: %v", got, err)
			}
		})
	}
}

// TestJudgePromptStatesItsContract pins both halves of the prompt.
//
// The output contract is the belt and parseJudgeResponse the seatbelt: deleting
// those sentences puts the prose preamble back. The analysis field is the other
// half and the more expensive one to lose — it is where the judge's reasoning
// lives now, and an earlier version of this prompt deleted the reasoning
// outright by telling the model to "reason silently" when it has no channel to
// do that in. See judgeResult.Analysis for what that measured.
func TestJudgePromptStatesItsContract(t *testing.T) {
	prompt := judgeDeveloperPrompt()
	for _, want := range []string{
		"nothing else",
		"code fence",
		"No preamble",
		"THINK IN THE analysis FIELD",
		"before the verdict of the same object",
		"leave none out",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("judgeDeveloperPrompt() no longer says %q:\n%s", want, prompt)
		}
	}
	if strings.Contains(prompt, "Reason silently") {
		t.Fatalf("judgeDeveloperPrompt() tells the model to reason silently again, which deletes the reasoning:\n%s", prompt)
	}
	if strings.Index(prompt, "\"analysis\":string") > strings.Index(prompt, "\"verdict\":") {
		t.Fatal("the shape line no longer puts analysis before verdict, so the reasoning stops conditioning the verdict")
	}
}
