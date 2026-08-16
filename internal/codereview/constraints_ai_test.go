package codereview

import (
	"strings"
	"testing"
)

func TestParseConstraintsJudgeResponse(t *testing.T) {
	valid := `{"gates":[{"gate":"back-pressure","status":"pass","justification":"scoped","findings":[]}]}`
	response, err := parseConstraintsJudgeResponse(valid)
	if err != nil {
		t.Fatalf("parseConstraintsJudgeResponse() error = %v", err)
	}
	if len(response.Gates) != 1 || response.Gates[0].Gate != "back-pressure" {
		t.Fatalf("response = %#v", response)
	}

	wrapped := "Here is my assessment:\n```json\n" + valid + "\n```\nDone."
	if _, err := parseConstraintsJudgeResponse(wrapped); err != nil {
		t.Fatalf("prose-wrapped reply did not parse: %v", err)
	}

	if _, err := parseConstraintsJudgeResponse("no json here"); err == nil {
		t.Fatalf("reply without JSON must not parse")
	}
	if _, err := parseConstraintsJudgeResponse(`{"gates":[]}`); err == nil {
		t.Fatalf("a reply that judged no gates must not parse as success")
	}
}

func TestConstraintsFindingFromAIDefaults(t *testing.T) {
	finding := constraintsFindingFromAI(GateBackPressure, constraintsAIFinding{Title: "Out of scope"}, KnownRuleIDs())
	if finding.Strength != "Worth exploring" {
		t.Fatalf("default strength = %q, want Worth exploring", finding.Strength)
	}
	perf := constraintsFindingFromAI(GatePerformance, constraintsAIFinding{Title: "N+1 query"}, KnownRuleIDs())
	if perf.Strength != "Strong" {
		t.Fatalf("performance default strength = %q, want Strong (the weighted gate)", perf.Strength)
	}
	if len(perf.Scopes) != 1 || perf.Scopes[0] != "performance" {
		t.Fatalf("performance scopes = %#v", perf.Scopes)
	}
}

func TestConstraintsGateScopes(t *testing.T) {
	scopes := constraintsGateScopes([]GateID{GateCodeHealth, GateBackPressure, GatePerformance, GateAccessibility})
	joined := strings.Join(scopes, ",")
	for _, want := range []string{"maintainability", "testing", "architecture", "performance"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("scopes = %#v, missing %s", scopes, want)
		}
	}
}

func TestConstraintsDeveloperPromptNamesEveryAIGate(t *testing.T) {
	prompt := constraintsDeveloperPrompt()
	for _, gate := range []string{"code-health", "back-pressure", "accessibility", "performance"} {
		if !strings.Contains(prompt, gate) {
			t.Fatalf("developer prompt does not define the %s gate", gate)
		}
	}
	if !strings.Contains(prompt, "justification") {
		t.Fatalf("developer prompt must demand the performance justification")
	}
}
