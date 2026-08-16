package codereview

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNormalizeRuleIDAcceptsKnownClassesInBothForms(t *testing.T) {
	known := KnownRuleIDs()
	cases := map[string]string{
		"review/wrong-logic":         "review/wrong-logic",
		"wrong-logic":                "review/wrong-logic",
		"  Wrong-Logic ":             "review/wrong-logic",
		"REVIEW/RACE-OR-IDEMPOTENCY": "review/race-or-idempotency",
	}
	for raw, want := range cases {
		if got := NormalizeRuleID(raw, known); got != want {
			t.Errorf("NormalizeRuleID(%q) = %q, want %q", raw, got, want)
		}
	}
}

func TestNormalizeRuleIDDropsUnknownAndInvented(t *testing.T) {
	known := KnownRuleIDs()
	for _, raw := range []string{
		"",
		"   ",
		"review/made-up-rule",
		"made-up-rule",
		"gx:recommended/no-swallowed-errors", // not known until a pack is loaded
		"review/Wrong Logic",                 // not kebab-case
		"review/",
	} {
		if got := NormalizeRuleID(raw, known); got != "" {
			t.Errorf("NormalizeRuleID(%q) = %q, want empty", raw, got)
		}
	}
}

func TestNormalizeRuleIDAcceptsExtraRulesFromBrief(t *testing.T) {
	known := KnownRuleIDs(
		RuleDef{ID: "gx:recommended/no-swallowed-errors"},
		RuleDef{ID: "REVIEW.md/no-reinvented-utils"},
	)
	if got := NormalizeRuleID("gx:recommended/no-swallowed-errors", known); got != "gx:recommended/no-swallowed-errors" {
		t.Errorf("pack rule not accepted: got %q", got)
	}
	// Repo rules are registered under the canonical "REVIEW.md/" spelling.
	// Matching is case-insensitive and the canonical spelling is what comes
	// back — a model that folds the case still lands on the registered ID.
	for _, raw := range []string{"REVIEW.md/no-reinvented-utils", "review.md/no-reinvented-utils", "  REVIEW.MD/No-Reinvented-Utils "} {
		if got := NormalizeRuleID(raw, known); got != "REVIEW.md/no-reinvented-utils" {
			t.Errorf("NormalizeRuleID(%q) = %q, want canonical REVIEW.md/no-reinvented-utils", raw, got)
		}
	}
}

func TestReviewClassPromptLineNamesEverySlug(t *testing.T) {
	line := reviewClassPromptLine()
	for _, def := range ReviewDefectClasses() {
		slug := strings.TrimPrefix(def.ID, reviewClassNamespace)
		if !strings.Contains(line, slug+" = ") {
			t.Errorf("prompt line missing class %q", slug)
		}
	}
	if !strings.Contains(line, "never invent a name") {
		t.Errorf("prompt line must forbid inventing rule names")
	}
}

func TestAIRecommendationsCarryRuleID(t *testing.T) {
	payload := `{"recommendations":[
	  {"title":"t1","summary":"s1","benefit":"b1","recommendation":"r1","kind":"defect","rule_id":"wrong-logic","strength":"Strong","file":"a.go","line":3},
	  {"title":"t2","summary":"s2","benefit":"b2","recommendation":"r2","kind":"suggestion","rule_id":"something-invented","strength":"Worth exploring"},
	  {"title":"t3","summary":"s3","benefit":"b3","recommendation":"r3","kind":"defect","strength":"Strong"}
	]}`
	var parsed aiReviewResponse
	if err := json.Unmarshal([]byte(payload), &parsed); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	findings := aiRecommendationsToFindings(parsed.Recommendations, ReviewBrief{})
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings, got %d", len(findings))
	}
	if findings[0].RuleID != "review/wrong-logic" {
		t.Errorf("finding 1 rule id = %q, want review/wrong-logic", findings[0].RuleID)
	}
	if findings[1].RuleID != "" {
		t.Errorf("invented rule id should be dropped, got %q", findings[1].RuleID)
	}
	if findings[2].RuleID != "" {
		t.Errorf("absent rule id should stay empty, got %q", findings[2].RuleID)
	}
	// The positional ID is untouched — RuleID is additive, never a replacement.
	if findings[0].ID != "ai.review.1" || findings[2].ID != "ai.review.3" {
		t.Errorf("positional IDs disturbed: %q %q", findings[0].ID, findings[2].ID)
	}
}

func TestConstraintsFindingFromAICarriesRuleID(t *testing.T) {
	known := KnownRuleIDs(AllRuleDefs(nil)...)
	f := constraintsFindingFromAI(GateBackPressure, constraintsAIFinding{
		Title: "t", Summary: "s", Recommendation: "r", RuleID: "api-contract-misuse",
	}, known)
	if f.RuleID != "review/api-contract-misuse" {
		t.Errorf("constraints finding rule id = %q", f.RuleID)
	}
	if f.ID != "constraints.back-pressure" {
		t.Errorf("constraints positional ID disturbed: %q", f.ID)
	}
	// A pack rule is accepted once the pack is in the known set, and a rule the
	// caller did not register is dropped rather than minted.
	packed := constraintsFindingFromAI(GateCodeHealth, constraintsAIFinding{Title: "t", RuleID: "gx:recommended/no-swallowed-errors"}, known)
	if packed.RuleID != "gx:recommended/no-swallowed-errors" {
		t.Errorf("pack rule id = %q, want gx:recommended/no-swallowed-errors", packed.RuleID)
	}
	unregistered := constraintsFindingFromAI(GateCodeHealth, constraintsAIFinding{Title: "t", RuleID: "REVIEW.md/not-declared"}, known)
	if unregistered.RuleID != "" {
		t.Errorf("undeclared repo rule id = %q, want empty", unregistered.RuleID)
	}
}
