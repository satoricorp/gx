package codereview

import (
	"strings"
	"testing"
)

func TestParseReviewJudgeResponse(t *testing.T) {
	valid := `{"gates":[{"gate":"back-pressure","status":"pass","justification":"scoped","findings":[]}]}`
	response, err := parseReviewJudgeResponse(valid)
	if err != nil {
		t.Fatalf("parseReviewJudgeResponse() error = %v", err)
	}
	if len(response.Gates) != 1 || response.Gates[0].Gate != "back-pressure" {
		t.Fatalf("response = %#v", response)
	}

	wrapped := "Here is my assessment:\n```json\n" + valid + "\n```\nDone."
	if _, err := parseReviewJudgeResponse(wrapped); err != nil {
		t.Fatalf("prose-wrapped reply did not parse: %v", err)
	}

	if _, err := parseReviewJudgeResponse("no json here"); err == nil {
		t.Fatalf("reply without JSON must not parse")
	}
	if _, err := parseReviewJudgeResponse(`{"gates":[]}`); err == nil {
		t.Fatalf("a reply that judged no gates must not parse as success")
	}
}

func TestReviewFindingFromAIDefaults(t *testing.T) {
	finding := reviewFindingFromAI(GateBackPressure, reviewAIFinding{Title: "Out of scope"})
	if finding.Strength != "Worth exploring" {
		t.Fatalf("default strength = %q, want Worth exploring", finding.Strength)
	}
	perf := reviewFindingFromAI(GatePerformance, reviewAIFinding{Title: "N+1 query"})
	if perf.Strength != "Strong" {
		t.Fatalf("performance default strength = %q, want Strong (the weighted gate)", perf.Strength)
	}
	if len(perf.Scopes) != 1 || perf.Scopes[0] != "performance" {
		t.Fatalf("performance scopes = %#v", perf.Scopes)
	}
}

func TestReviewGateScopes(t *testing.T) {
	scopes := reviewGateScopes([]GateID{GateCodeHealth, GateBackPressure, GatePerformance, GateAccessibility})
	joined := strings.Join(scopes, ",")
	for _, want := range []string{"maintainability", "testing", "architecture", "performance"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("scopes = %#v, missing %s", scopes, want)
		}
	}
}

func TestReviewDeveloperPromptNamesEveryAIGate(t *testing.T) {
	prompt := gateDeveloperPrompt()
	for _, gate := range []string{"code-health", "back-pressure", "accessibility", "performance"} {
		if !strings.Contains(prompt, gate) {
			t.Fatalf("developer prompt does not define the %s gate", gate)
		}
	}
	if !strings.Contains(prompt, "justification") {
		t.Fatalf("developer prompt must demand the performance justification")
	}
}
