package codereview

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/satoricorp/lgtm/internal/lgtmtest"
)

// liveJudgeRequest is the fixed three-candidate probe the live tests share: one
// finding the snippet supports, one it refutes, and one that is real but
// cosmetic.
func liveJudgeRequest() judgeRequest {
	return judgeRequest{
		RepoRoot:     "/repo",
		ChangedFiles: []string{"internal/app/app.go", "internal/app/poll.go"},
		Candidates: []judgeCandidate{
			{
				ID:             "live.supported",
				Title:          "Error from Close is discarded",
				Summary:        "`internal/app/app.go` calls `defer f.Close()` on a file it wrote, so a failed flush is lost.",
				Recommendation: "Check the error from Close and return it.",
				NamedFiles:     []string{"internal/app/app.go"},
				FileContentSnippets: []judgeContentSnippet{{
					File: "internal/app/app.go",
					Text: "package app\n\nfunc Write(path string) error {\n\tf, err := os.Create(path)\n\tif err != nil {\n\t\treturn err\n\t}\n\tdefer f.Close()\n\t_, err = f.WriteString(\"data\")\n\treturn err\n}\n",
				}},
			},
			{
				// The mechanism this claims is contradicted by the snippet it is
				// handed: ReadFile's error is checked, so nothing escapes. See
				// TestLiveJudgeRejectsAClaimTheFileContentRefutes.
				ID:             "live.refuted",
				Title:          "Poll loop crashes when the answer file disappears",
				Summary:        "`internal/app/poll.go` stats the file and then reads it. Only the JSON decode error is handled, so the read error from a file deleted in between escapes the loop and crashes the process.",
				Recommendation: "Handle the read error explicitly instead of letting it escape.",
				Evidence:       []string{"Reviewer: Bedrock A", "internal/app/poll.go:12"},
				NamedFiles:     []string{"internal/app/poll.go"},
				FileContentSnippets: []judgeContentSnippet{{
					File: "internal/app/poll.go",
					Text: "package app\n\nfunc poll(path string) map[string]any {\n\tfor {\n\t\tif _, err := os.Stat(path); err == nil {\n\t\t\tdata, err := os.ReadFile(path)\n\t\t\tif err == nil {\n\t\t\t\tvar out map[string]any\n\t\t\t\tif json.Unmarshal(data, &out) == nil {\n\t\t\t\t\treturn out\n\t\t\t\t}\n\t\t\t}\n\t\t\t// a torn read between Stat and rename settling — the next poll gets it whole\n\t\t}\n\t\ttime.Sleep(300 * time.Millisecond)\n\t}\n}\n",
				}},
			},
			{
				ID:             "live.cosmetic",
				Title:          "Package name could be shorter",
				Summary:        "The package in `internal/app/app.go` is named `app`, which is generic.",
				Recommendation: "Rename the package.",
				NamedFiles:     []string{"internal/app/app.go"},
				FileContentSnippets: []judgeContentSnippet{{
					File: "internal/app/app.go",
					Text: "package app\n",
				}},
			},
		},
	}
}

// liveJudgeVerdicts runs one real judge call and returns the parsed reply, or
// skips when this machine cannot reach Bedrock.
func liveJudgeVerdicts(t *testing.T) (bedrockCompletion, []judgeResult, judgeRequest) {
	t.Helper()
	// Billable Bedrock inference, so reaching it has to be asked for rather
	// than inherited from whichever AWS profile the shell happens to carry.
	lgtmtest.RequireNoNetwork(t)
	if _, err := bedrockCredentialsFromEnv(); err != nil {
		t.Skipf("no AWS credentials; skipping live judge test (%v)", err)
	}
	plan, err := resolveBedrockTransportPlan()
	if err != nil {
		t.Skipf("no Bedrock transport; skipping live judge test (%v)", err)
	}
	// Built directly rather than through judgeFromEnvWithPolicy, which TestMain
	// disables for the rest of the suite.
	client := newBedrockReviewer(plan.newTransport(), resolveBedrockJudgeModel())

	request := liveJudgeRequest()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	completion, err := client.completeJSON(ctx, judgeDeveloperPrompt(), mustJSON(request), defaultJudgeMaxOutputTokens)
	if err != nil {
		t.Fatalf("judge call error = %v", err)
	}
	if completion.truncated() {
		t.Fatalf("judge reply stopped at the %d-token cap on a %d-candidate batch; the budget is wrong",
			defaultJudgeMaxOutputTokens, len(request.Candidates))
	}
	results, err := parseJudgeResponse(completion.Text, judgeCandidateIDs(request.Candidates))
	if err != nil {
		t.Fatalf("parseJudgeResponse() error = %v on a live reply", err)
	}
	return completion, results, request
}

// TestLiveJudgeHonorsTheBareJSONOutputContract asks the real judge model whether
// it still obeys the output contract.
//
// This is one of two properties no fixture can defend. judgeDeveloperPrompt is
// the primary defense against the failure that made a real review run degraded —
// a reply that opens with paragraphs of reasoning, quotes braces from source
// code, and buries its JSON in a ```json fence — and whether a prompt still
// works is a fact about the model, not about this repository. A model update can
// withdraw that behavior without a line of lgtm changing.
//
// Measured before and after the contract landed, replaying a real 24-candidate
// batch against us.anthropic.claude-sonnet-4-6 twenty times each: before, every
// reply was fenced, all but one carried a prose preamble, and one failed to
// parse outright; after, all twenty were bare objects and none failed.
//
// It only reads, and it skips without credentials, so the offline suite is
// unaffected.
func TestLiveJudgeHonorsTheBareJSONOutputContract(t *testing.T) {
	completion, results, request := liveJudgeVerdicts(t)

	text := strings.TrimSpace(completion.Text)
	if !strings.HasPrefix(text, "{") || !strings.HasSuffix(text, "}") {
		t.Errorf("judge no longer answers with a bare JSON object — the output contract has stopped working.\nfirst 200 chars: %q\nlast 80 chars: %q",
			text[:min(200, len(text))], text[max(0, len(text)-80):])
	}
	if len(results) != len(request.Candidates) {
		t.Fatalf("judge returned %d verdicts for %d candidates", len(results), len(request.Candidates))
	}
	seen := map[string]bool{}
	for _, result := range results {
		seen[strings.TrimSpace(result.CandidateID)] = true
		if strings.TrimSpace(result.Verdict) == "" {
			t.Errorf("verdict for %q is empty", result.CandidateID)
		}
	}
	for _, candidate := range request.Candidates {
		if !seen[candidate.ID] {
			t.Errorf("no verdict for candidate %q", candidate.ID)
		}
	}
}

// TestLiveJudgeStillReasonsAboutEveryCandidate is the other property a fixture
// cannot defend, and the more expensive one to lose.
//
// The output contract's first draft bought clean JSON by telling the model to
// "reason silently", which on this request shape has no channel to happen in and
// so deleted the reasoning outright — 24% fewer output tokens and double the
// confirm rate. The analysis field is what replaced it. A model that starts
// returning it empty, or a prompt edit that drops it, would look exactly like a
// healthy review while verifying nothing, so it is asserted against the live
// model rather than trusted.
// The assertion is deliberately about the batch rather than each candidate.
// Before the prompt said every result must carry one, 1 analysis in 36 came back
// empty while its verdict was still correct; after, 0 in 45, with the shortest
// at 234 characters and the median at 447. A per-candidate assertion would
// therefore be a live test that fails a few percent of the time on a judge that
// is working, and the failure actually worth catching — reasoning switched off,
// by a prompt edit or a model change — takes every analysis with it, not one.
// So the floor is two thirds of the batch: 100% under a healthy model, 0% under
// the regression, and unbothered by a single dropped field.
func TestLiveJudgeStillReasonsAboutEveryCandidate(t *testing.T) {
	_, results, _ := liveJudgeVerdicts(t)

	// A length floor rather than "non-empty": the failure being guarded is a
	// model that emits the field and stops thinking, and one clause is
	// indistinguishable from that.
	const substantive = 120
	reasoned := 0
	for _, result := range results {
		analysis := strings.TrimSpace(result.Analysis)
		if len(analysis) >= substantive {
			reasoned++
			continue
		}
		t.Logf("analysis for %q is only %d chars: %q", result.CandidateID, len(analysis), analysis)
	}
	if reasoned*3 < len(results)*2 {
		t.Fatalf("only %d of %d verdicts carried a substantive analysis; the judge has stopped reasoning about what it verifies",
			reasoned, len(results))
	}
}

// TestLiveJudgeRejectsAClaimTheFileContentRefutes is the regression for the
// judgement the output contract cost.
//
// live.refuted asserts that only the JSON decode error is handled and the read
// error escapes; the snippet it is handed checks that error. The equivalent
// candidate in the yeet A/B — a claim that a bare catch swallows only JSON
// errors, against a snippet whose catch is bare — was confirmed 1 run in 18
// under "reason silently", at impact functional and confidence 0.85, which
// clears minSurfaceConfidence and reaches the reader as a verified finding
// pointing at a bug that is not there. With the reasoning restored to the
// analysis field it was rejected in all 18. This fixture was rejected in all 27
// live runs measured while writing the test.
func TestLiveJudgeRejectsAClaimTheFileContentRefutes(t *testing.T) {
	_, results, _ := liveJudgeVerdicts(t)

	for _, result := range results {
		if strings.TrimSpace(result.CandidateID) != "live.refuted" {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(result.Verdict), "confirmed") {
			t.Fatalf("judge confirmed a finding whose stated mechanism the provided file content refutes.\nanalysis: %s\nnote: %s",
				result.Analysis, result.VerificationNote)
		}
		return
	}
	t.Fatal("no verdict for live.refuted")
}
