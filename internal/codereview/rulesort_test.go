package codereview

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

type stubRuleSorter struct {
	results []ruleSortResult
	err     error
	seen    ruleSortRequest
	calls   int
}

func (s *stubRuleSorter) Available() bool { return true }

func (s *stubRuleSorter) Sort(_ context.Context, req ruleSortRequest) ([]ruleSortResult, error) {
	s.calls++
	s.seen = req
	return s.results, s.err
}

func findingsForSort() []Finding {
	return []Finding{
		{ID: "ai.review.1", Title: "Error from the write is discarded", Summary: "the early return drops err", Kind: "defect", JudgeVerdict: "confirmed", JudgeSeverity: 4},
		{ID: "ai.review.2", Title: "Something unclassifiable", Summary: "no rule covers this", Kind: "suggestion", JudgeVerdict: "confirmed", JudgeSeverity: 1},
	}
}

// The catalog is what the sorter may choose from, and the review/* defect
// classes are deliberately not in it: they are the open-ended judgment
// categories detection used to carry, and two of them are catch-alls that would
// sit next to the specific rules and absorb them.
func TestRuleSortCatalogIsThePackAndNotTheDefectClasses(t *testing.T) {
	catalog, err := ruleSortCatalog(ReviewContext{})
	if err != nil {
		t.Fatalf("ruleSortCatalog() error = %v", err)
	}
	if len(catalog) == 0 {
		t.Fatal("catalog is empty")
	}
	for _, rule := range catalog {
		if strings.HasPrefix(rule.ID, reviewClassNamespace) {
			t.Errorf("catalog contains defect class %q; the sorter must choose pack rules only", rule.ID)
		}
		if !strings.HasPrefix(rule.ID, recommendedPackNamespace) {
			t.Errorf("catalog rule %q is not a pack rule", rule.ID)
		}
		if strings.TrimSpace(rule.Summary) == "" {
			t.Errorf("catalog rule %q has no summary for the model to match against", rule.ID)
		}
	}
	pack, err := RecommendedPack()
	if err != nil {
		t.Fatalf("RecommendedPack() error = %v", err)
	}
	if len(catalog) != len(pack) {
		t.Fatalf("catalog has %d rules, want the pack's %d", len(catalog), len(pack))
	}
}

func TestRuleSortLabelsFindingsAndLeavesEverythingElseAlone(t *testing.T) {
	findings := findingsForSort()
	before := findingsForSort()
	sorter := &stubRuleSorter{results: []ruleSortResult{
		{CandidateID: "ai.review.1", Reasoning: "the write's error is dropped on the early return", RuleID: "gx:recommended/no-swallowed-errors"},
		{CandidateID: "ai.review.2", Reasoning: "no rule names this", RuleID: ""},
	}}
	outcome := runRuleSort(context.Background(), sorter, ReviewContext{}, findings)
	if outcome.Labeled != 1 || outcome.Unlabeled != 1 {
		t.Fatalf("outcome labeled=%d unlabeled=%d, want 1 and 1", outcome.Labeled, outcome.Unlabeled)
	}
	if findings[0].RuleID != "gx:recommended/no-swallowed-errors" {
		t.Errorf("finding 0 RuleID = %q, want the pack rule", findings[0].RuleID)
	}
	if findings[1].RuleID != "" {
		t.Errorf("finding 1 RuleID = %q, want empty — no rule fit", findings[1].RuleID)
	}
	// Labeling must name a finding, never edit one. If this ever fails, the
	// pass has become detection running after the judge.
	for i := range findings {
		findings[i].RuleID = before[i].RuleID
		if !reflect.DeepEqual(findings[i], before[i]) {
			t.Errorf("finding %d changed beyond RuleID:\n got %+v\nwant %+v", i, findings[i], before[i])
		}
	}
}

// A model that answers with a plausible-looking name that is not a rule must
// not mint one. This is the mechanical half of "never force a label" — the
// prompt asks, and NormalizeRuleID enforces.
func TestRuleSortDropsInventedAndOutOfCatalogRules(t *testing.T) {
	for _, invented := range []string{
		"gx:recommended/no-bad-code",
		"review/wrong-logic",
		"wrong-logic",
		"made up entirely",
	} {
		findings := findingsForSort()
		sorter := &stubRuleSorter{results: []ruleSortResult{
			{CandidateID: "ai.review.1", RuleID: invented},
		}}
		if outcome := runRuleSort(context.Background(), sorter, ReviewContext{}, findings); outcome.Labeled != 0 {
			t.Errorf("rule_id %q was accepted; want dropped", invented)
		}
		if findings[0].RuleID != "" {
			t.Errorf("rule_id %q wrote %q onto the finding", invented, findings[0].RuleID)
		}
	}
}

// A retired slug is not an invented one: it was released, so it resolves to the
// rule that absorbed it rather than being dropped.
func TestRuleSortResolvesRetiredSlugs(t *testing.T) {
	findings := findingsForSort()
	sorter := &stubRuleSorter{results: []ruleSortResult{
		{CandidateID: "ai.review.1", RuleID: "gx:recommended/no-secrets-in-logs"},
	}}
	runRuleSort(context.Background(), sorter, ReviewContext{}, findings)
	if findings[0].RuleID != "gx:recommended/no-secrets" {
		t.Errorf("retired slug resolved to %q, want gx:recommended/no-secrets", findings[0].RuleID)
	}
}

// A failed sort must degrade to the output that shipped before the pass
// existed: same findings, no rule names, nothing dropped.
func TestRuleSortFailureLeavesFindingsIntact(t *testing.T) {
	findings := findingsForSort()
	sorter := &stubRuleSorter{err: context.DeadlineExceeded}
	outcome := runRuleSort(context.Background(), sorter, ReviewContext{}, findings)
	if outcome.Failed != 1 || outcome.Err == nil {
		t.Fatalf("outcome = %+v, want one failed batch carrying the error", outcome)
	}
	if len(findings) != 2 {
		t.Fatalf("findings len = %d, want 2 — a failed sort must not drop findings", len(findings))
	}
	for i := range findings {
		if findings[i].RuleID != "" {
			t.Errorf("finding %d got a rule from a failed sort", i)
		}
	}
}

// Detection must not be told the rules exist. This is the whole reason the pass
// is placed after verification, so it is worth a test that fails loudly if the
// pack ever finds its way back into the reviewer prompt.
func TestReviewerPromptDoesNotCarryRules(t *testing.T) {
	brief := ReviewBrief{Scope: "patch_focused", ReviewProfile: "patch_focused"}
	prompt := strings.Join(baseReviewDeveloperPromptLines(brief), "\n")
	if strings.Contains(prompt, "rule_id") {
		t.Error("reviewer prompt asks for rule_id; rule assignment belongs after verification")
	}
	for _, slug := range []string{"no-swallowed-errors", "variable-misuse", "boolean-polarity"} {
		if strings.Contains(prompt, slug) {
			t.Errorf("reviewer prompt names pack rule %q", slug)
		}
	}
}
