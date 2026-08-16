package codereview

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestNormalizeStoryCategory(t *testing.T) {
	cases := map[string]string{
		"behavior_delta":   StoryBehaviorDelta,
		" Behavior-Delta ": StoryBehaviorDelta,
		"behaviour_delta":  StoryBehaviorDelta,
		"new_surface":      StoryNewSurface,
		"New Surface":      StoryNewSurface,
		"semantic_shift":   StorySemanticShift,
		"decision":         StoryDecision,
		"exception_taken":  StoryExceptionTaken,
		"coverage_move":    StoryCoverageMove,
		"coverage":         StoryCoverageMove,
		"":                 "",
		"finding":          "",
		"bug":              "",
	}
	for in, want := range cases {
		if got := normalizeStoryCategory(in); got != want {
			t.Errorf("normalizeStoryCategory(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNormalizeMateriality(t *testing.T) {
	cases := map[string]string{
		"high":     StoryMaterialityHigh,
		" HIGH ":   StoryMaterialityHigh,
		"medium":   StoryMaterialityMedium,
		"Med":      StoryMaterialityMedium,
		"low":      StoryMaterialityLow,
		"":         "",
		"certain":  "",
		"unknown":  "",
		"strong":   "",
		"moderate": StoryMaterialityMedium,
	}
	for in, want := range cases {
		if got := normalizeMateriality(in); got != want {
			t.Errorf("normalizeMateriality(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSortStoryByMaterialityIsOrderedAndStable(t *testing.T) {
	items := []StoryItem{
		{Headline: "low-1", Materiality: "low"},
		{Headline: "none-1"},
		{Headline: "medium-1", Materiality: "medium"},
		{Headline: "high-1", Materiality: "high"},
		{Headline: "low-2", Materiality: "low"},
		{Headline: "high-2", Materiality: "high"},
		{Headline: "medium-2", Materiality: "medium"},
	}
	sortStoryByMateriality(items)
	var got []string
	for _, item := range items {
		got = append(got, item.Headline)
	}
	want := []string{"high-1", "high-2", "medium-1", "medium-2", "low-1", "low-2", "none-1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sortStoryByMateriality order = %v, want %v", got, want)
	}
}

const reviewPolicyExceptionsDiff = `diff --git a/REVIEW.md b/REVIEW.md
index 1111111..2222222 100644
--- a/REVIEW.md
+++ b/REVIEW.md
@@ -10,6 +10,7 @@
 Some rule guidance.

 ## Exceptions

 - internal/legacy/** is exempt from wrong-logic until the rewrite lands
+- internal/billing/retry.go is exempt from missing-idempotency-key

 ## High-risk paths
`

const reviewPolicyRuleOnlyDiff = `diff --git a/REVIEW.md b/REVIEW.md
index 1111111..2222222 100644
--- a/REVIEW.md
+++ b/REVIEW.md
@@ -1,6 +1,7 @@
 # Review policy

 ## Idempotency keys

 Every write that can be retried carries an idempotency key.
+Keys are derived from the request body, never from a timestamp.

 ## Exceptions

 - internal/legacy/** is exempt from wrong-logic until the rewrite lands
`

func TestMandatoryStoryItemsEmitsExceptionTakenWhenExceptionsChange(t *testing.T) {
	policy := &ReviewPolicy{Present: true, Path: "REVIEW.md"}
	items := MandatoryStoryItems(policy, []string{"internal/billing/retry.go", "REVIEW.md"}, []DiffSnippet{
		{File: "internal/billing/retry.go", Diff: "@@ -1 +1 @@\n-a\n+b\n"},
		{File: "REVIEW.md", Diff: reviewPolicyExceptionsDiff},
	})
	if len(items) != 1 {
		t.Fatalf("MandatoryStoryItems() = %#v, want exactly one item", items)
	}
	item := items[0]
	if item.Category != StoryExceptionTaken || item.Materiality != StoryMaterialityHigh || item.File != "REVIEW.md" {
		t.Fatalf("item = %#v, want exception_taken/high on REVIEW.md", item)
	}
	if item.Headline == "" || item.Consequence == "" {
		t.Fatalf("item = %#v, want headline and consequence", item)
	}
}

func TestMandatoryStoryItemsSilentWhenExceptionsUntouched(t *testing.T) {
	// A REVIEW.md edit outside the Exceptions section is a policy change the
	// reader can see in the diff; it is not the change exempting itself.
	if items := MandatoryStoryItems(nil, []string{"REVIEW.md"}, []DiffSnippet{{File: "REVIEW.md", Diff: reviewPolicyRuleOnlyDiff}}); len(items) != 0 {
		t.Fatalf("MandatoryStoryItems() = %#v, want none for a rule-only edit", items)
	}
	// REVIEW.md not among the changed files: nothing to say, even if a stale
	// snippet is around.
	if items := MandatoryStoryItems(nil, []string{"main.go"}, []DiffSnippet{{File: "REVIEW.md", Diff: reviewPolicyExceptionsDiff}}); len(items) != 0 {
		t.Fatalf("MandatoryStoryItems() = %#v, want none when REVIEW.md is not changed", items)
	}
	// No diff for the file at all.
	if items := MandatoryStoryItems(nil, []string{"REVIEW.md"}, nil); len(items) != 0 {
		t.Fatalf("MandatoryStoryItems() = %#v, want none without a diff", items)
	}
	// Nil policy and no changes.
	if items := MandatoryStoryItems(nil, nil, nil); items != nil {
		t.Fatalf("MandatoryStoryItems() = %#v, want nil", items)
	}
}

func TestMandatoryStoryItemsSeesRemovedExceptionsAndNewSection(t *testing.T) {
	removed := strings.Replace(reviewPolicyExceptionsDiff,
		"+- internal/billing/retry.go is exempt from missing-idempotency-key",
		"-- internal/billing/retry.go is exempt from missing-idempotency-key", 1)
	if items := MandatoryStoryItems(nil, []string{"REVIEW.md"}, []DiffSnippet{{File: "REVIEW.md", Diff: removed}}); len(items) != 1 {
		t.Fatalf("MandatoryStoryItems() = %#v, want one item for a removed exception", items)
	}
	// Adding the section heading itself is taking an exception.
	added := "@@ -1,3 +1,6 @@\n # Review policy\n \n+## Exceptions\n+\n+- cmd/** is exempt from missing-error-check\n"
	if items := MandatoryStoryItems(nil, []string{"REVIEW.md"}, []DiffSnippet{{File: "REVIEW.md", Diff: added}}); len(items) != 1 {
		t.Fatalf("MandatoryStoryItems() = %#v, want one item for a new Exceptions section", items)
	}
	// An untracked REVIEW.md arrives as file content; one with an Exceptions
	// section is a change that adds exceptions.
	content := diffUnavailableContentHeader + "\n# Review policy\n\n## Exceptions\n\n- cmd/** is exempt\n"
	if items := MandatoryStoryItems(nil, []string{"REVIEW.md"}, []DiffSnippet{{File: "REVIEW.md", Diff: content}}); len(items) != 1 {
		t.Fatalf("MandatoryStoryItems() = %#v, want one item for a new REVIEW.md with exceptions", items)
	}
	// Function context on the hunk header places the hunk in Exceptions even
	// when the heading itself is outside the window.
	funcCtx := "@@ -12,3 +12,4 @@ ## Exceptions\n \n - internal/legacy/** is exempt\n+- internal/billing/** is exempt\n"
	if items := MandatoryStoryItems(nil, []string{"REVIEW.md"}, []DiffSnippet{{File: "REVIEW.md", Diff: funcCtx}}); len(items) != 1 {
		t.Fatalf("MandatoryStoryItems() = %#v, want one item when the hunk header names Exceptions", items)
	}
}

func TestAIStoryToStoryItemsNormalizesAndDrops(t *testing.T) {
	brief := ReviewBrief{Rules: []RuleDef{{ID: "REVIEW.md/idempotency-keys", Summary: "keys"}}}
	raw := []aiStoryItem{
		{
			Headline:    "  Checkout retries on 5xx where it used to fail fast ",
			Category:    "Behavior-Delta",
			Materiality: "HIGH",
			File:        " internal/checkout/client.go ",
			Line:        json.RawMessage(`"42"`),
			Consequence: " Callers see up to 3 attempts; p99 latency can rise by 2x under upstream failure. ",
			Asked:       " make checkout resilient ",
			Chose:       " bounded retry with jitter over a circuit breaker ",
			Watch:       []string{" checkout_retry_total ", "", "upstream 5xx rate"},
			PassedRule:  "review/wrong-logic",
		},
		{
			Headline:    "Idempotency key now derived from the request body",
			Category:    "semantic_shift",
			Materiality: "medium",
			File:        "internal/billing/retry.go",
			Line:        json.RawMessage(`7`),
			Consequence: "Retries of the same body are deduplicated; retries with a new timestamp no longer double-charge.",
			PassedRule:  "REVIEW.md/Idempotency-Keys",
		},
		{
			Headline:    "Has an invented rule",
			Category:    "new_surface",
			Materiality: "low",
			Consequence: "A flag exists.",
			PassedRule:  "review/made-up-rule",
		},
		{
			Headline:    "Model claims an exception",
			Category:    "exception_taken",
			Materiality: "low",
			Consequence: "It should not get to.",
		},
		{Headline: "", Category: "decision", Materiality: "high", Consequence: "no headline"},
		{Headline: "no consequence", Category: "decision", Materiality: "high", Consequence: "  "},
		{Headline: "unknown labels", Category: "finding", Materiality: "certain", Consequence: "kept, unlabeled"},
	}
	items := aiStoryToStoryItems(raw, brief)
	if len(items) != 5 {
		t.Fatalf("aiStoryToStoryItems() = %d items, want 5:\n%#v", len(items), items)
	}
	first := items[0]
	if first.Headline != "Checkout retries on 5xx where it used to fail fast" ||
		first.Category != StoryBehaviorDelta || first.Materiality != StoryMaterialityHigh ||
		first.File != "internal/checkout/client.go" || first.Line != 42 ||
		first.Asked != "make checkout resilient" || first.Chose != "bounded retry with jitter over a circuit breaker" ||
		first.PassedRule != "review/wrong-logic" {
		t.Fatalf("first = %#v", first)
	}
	if !reflect.DeepEqual(first.Watch, []string{"checkout_retry_total", "upstream 5xx rate"}) {
		t.Fatalf("first.Watch = %#v", first.Watch)
	}
	if items[1].PassedRule != "REVIEW.md/idempotency-keys" || items[1].Line != 7 || items[1].Category != StorySemanticShift {
		t.Fatalf("second = %#v, want the REVIEW.md rule id accepted in its canonical spelling", items[1])
	}
	if items[2].PassedRule != "" {
		t.Fatalf("third.PassedRule = %q, want an invented rule dropped", items[2].PassedRule)
	}
	if items[3].Category != "" {
		t.Fatalf("fourth.Category = %q, want exception_taken reserved for mandatory items", items[3].Category)
	}
	if items[4].Category != "" || items[4].Materiality != "" || items[4].Headline != "unknown labels" {
		t.Fatalf("fifth = %#v, want kept with empty labels", items[4])
	}
	if got := aiStoryToStoryItems(nil, brief); got != nil {
		t.Fatalf("aiStoryToStoryItems(nil) = %#v, want nil", got)
	}
}

func TestParsePRSummaryReviewIncludesStory(t *testing.T) {
	summary, err := ParsePRSummaryReview(`{
		"overview":"Adds retry to checkout.",
		"notable_changes":[{"file":"internal/checkout/client.go","line":42,"note":"Retry loop."}],
		"story":[
			{"headline":"Checkout retries on 5xx where it used to fail fast","category":"behavior_delta","materiality":"high","file":"internal/checkout/client.go","line":42,"consequence":"Up to 3 attempts per call.","asked":"make checkout resilient","chose":"retry over breaker","watch":["retry rate"],"passed_rule":"wrong-logic"},
			{"headline":"","category":"decision","materiality":"low","consequence":"dropped"}
		],
		"recommendations":[]
	}`, ReviewBrief{})
	if err != nil {
		t.Fatalf("ParsePRSummaryReview() error = %v", err)
	}
	if len(summary.NotableChanges) != 1 {
		t.Fatalf("NotableChanges = %#v, want the existing lane untouched", summary.NotableChanges)
	}
	if len(summary.Story) != 1 {
		t.Fatalf("Story = %#v, want one item", summary.Story)
	}
	item := summary.Story[0]
	if item.Headline != "Checkout retries on 5xx where it used to fail fast" || item.Category != StoryBehaviorDelta ||
		item.Materiality != StoryMaterialityHigh || item.File != "internal/checkout/client.go" || item.Line != 42 ||
		item.Consequence != "Up to 3 attempts per call." || item.Asked != "make checkout resilient" ||
		item.Chose != "retry over breaker" || item.PassedRule != "review/wrong-logic" ||
		!reflect.DeepEqual(item.Watch, []string{"retry rate"}) {
		t.Fatalf("Story[0] = %#v", item)
	}
}

func TestParsePRSummaryReviewStorySurvivesFencedJSON(t *testing.T) {
	summary, err := ParsePRSummaryReview("Here you go:\n```json\n{\"story\":[{\"headline\":\"H\",\"category\":\"decision\",\"materiality\":\"low\",\"consequence\":\"C\"}],\"recommendations\":[]}\n```", ReviewBrief{})
	if err != nil {
		t.Fatalf("ParsePRSummaryReview() error = %v", err)
	}
	if len(summary.Story) != 1 || summary.Story[0].Headline != "H" {
		t.Fatalf("Story = %#v, want the fenced story parsed", summary.Story)
	}
}

func TestMultiAIReviewerCarriesStoryOrderedByMateriality(t *testing.T) {
	reviewer := multiAIReviewer{reviewers: []namedAIReviewer{
		{name: "bedrock-a", label: "Bedrock A", reviewer: cannedAIReviewer{payload: `{"recommendations":[]}`}},
		{name: "bedrock-b", label: "Bedrock B", reviewer: cannedAIReviewer{payload: `{
			"story":[
				{"headline":"low first","category":"decision","materiality":"low","consequence":"c"},
				{"headline":"then high","category":"behavior_delta","materiality":"high","consequence":"c"}
			],
			"recommendations":[]
		}`}},
	}}
	summary, err := reviewer.ReviewForSummary(context.Background(), ReviewBrief{})
	if err != nil {
		t.Fatalf("ReviewForSummary() error = %v", err)
	}
	if len(summary.Story) != 2 || summary.Story[0].Headline != "then high" || summary.Story[1].Headline != "low first" {
		t.Fatalf("Story = %#v, want the non-empty leg's story sorted high first", summary.Story)
	}
}

func TestAssembleReportStoryPutsMandatoryFirstAndAttachesHunks(t *testing.T) {
	diffs := []DiffSnippet{
		{File: "REVIEW.md", Diff: reviewPolicyExceptionsDiff},
		{File: "internal/billing/retry.go", Diff: "@@ -1,3 +1,4 @@\n package billing\n \n+// key from body\n func Retry() {}\n"},
	}
	model := []StoryItem{{Headline: "Key from body", Category: StorySemanticShift, Materiality: "high", File: "internal/billing/retry.go", Line: 3, Consequence: "c"}}
	story := assembleReportStory(nil, []string{"REVIEW.md", "internal/billing/retry.go"}, diffs, model)
	if len(story) != 2 {
		t.Fatalf("assembleReportStory() = %#v, want mandatory + model", story)
	}
	if story[0].Category != StoryExceptionTaken {
		t.Fatalf("story[0] = %#v, want the mandatory item first", story[0])
	}
	if !strings.Contains(story[1].DiffHunk, "+// key from body") {
		t.Fatalf("story[1].DiffHunk = %q, want the hunk around line 3", story[1].DiffHunk)
	}
	if got := assembleReportStory(nil, []string{"main.go"}, nil, nil); got != nil {
		t.Fatalf("assembleReportStory() with nothing to say = %#v, want nil", got)
	}
}

func TestReportStoryOmittedFromJSONWhenEmpty(t *testing.T) {
	data, err := json.Marshal(Report{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"story"`) {
		t.Fatalf("Report JSON = %s, want no story key when empty", data)
	}
	data, err = json.Marshal(Report{Story: []StoryItem{{Headline: "H", Category: StoryDecision, Materiality: "low", Consequence: "C"}}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"story":[{"headline":"H"`) {
		t.Fatalf("Report JSON = %s, want story serialized", data)
	}
}
