package cli

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/codereview"
)

// The gx Cloud history row stores a 1-10 confidence derived from Strength. The
// real Strength vocabulary is Blocking / Strong / Worth exploring / Speculative;
// the mapping used to know only high/medium/low and filed every real strength
// under the default, so Blocking and Speculative findings recorded identical
// confidence.
func TestReviewFindingConfidenceMapsTheStrengthVocabulary(t *testing.T) {
	cases := []struct {
		strength string
		want     int
	}{
		{"Blocking", 9},
		{"blocking", 9},
		{"Strong", 8},
		{"Worth exploring", 5},
		{"worth-exploring", 5},
		{"Speculative", 3},
		{"  Speculative  ", 3},
		// Aliases kept for older producers.
		{"high", 8},
		{"medium", 6},
		{"moderate", 6},
		{"low", 4},
		{"weak", 4},
		// Unknown and empty fall to the middle.
		{"", 6},
		{"critical", 6},
	}
	for _, tc := range cases {
		t.Run(tc.strength, func(t *testing.T) {
			if got := reviewFindingConfidence(tc.strength); got != tc.want {
				t.Fatalf("reviewFindingConfidence(%q) = %d, want %d", tc.strength, got, tc.want)
			}
		})
	}
}

// A finding that names its rule is identified by the rule and the location,
// not by the prose around it. AI findings carry a positional ID and
// model-written Title/Summary, so hashing those made the same rule firing on
// the same line in two runs look like two unrelated findings.
func TestReviewFindingFingerprintUsesRuleAndLocationWhenRuleIDSet(t *testing.T) {
	base := codereview.Finding{
		ID:      "ai.review.3",
		RuleID:  "review/wrong-logic",
		Title:   "Off-by-one in the loop bound",
		Summary: "The loop stops one short.",
	}
	rerun := codereview.Finding{
		ID:      "ai.review.7",
		RuleID:  "review/wrong-logic",
		Title:   "Loop bound is off by one",
		Summary: "Reads one element fewer than intended.",
	}
	got := reviewFindingFingerprint("acme/widgets", base, "pkg/loop.go", 42, "review")
	if again := reviewFindingFingerprint("acme/widgets", rerun, "pkg/loop.go", 42, "correctness"); again != got {
		t.Fatalf("same rule+file+line with different prose/ID/category should share a fingerprint: %s vs %s", got, again)
	}
	if other := reviewFindingFingerprint("acme/widgets", base, "pkg/loop.go", 43, "review"); other == got {
		t.Fatalf("a different line must yield a different fingerprint")
	}
	if other := reviewFindingFingerprint("acme/widgets", base, "pkg/other.go", 42, "review"); other == got {
		t.Fatalf("a different file must yield a different fingerprint")
	}
	if other := reviewFindingFingerprint("acme/gadgets", base, "pkg/loop.go", 42, "review"); other == got {
		t.Fatalf("a different repo must yield a different fingerprint")
	}
	// Whitespace around the rule is not identity.
	padded := base
	padded.RuleID = "  review/wrong-logic  "
	if again := reviewFindingFingerprint("acme/widgets", padded, "pkg/loop.go", 42, "review"); again != got {
		t.Fatalf("padded RuleID should share the fingerprint")
	}
}

// Findings without a rule keep the legacy layout byte for byte, so every row
// already recorded in gx Cloud keeps the fingerprint it was stored under.
func TestReviewFindingFingerprintLegacyLayoutUnchangedWithoutRuleID(t *testing.T) {
	finding := codereview.Finding{
		ID:      "ai.review.3",
		Title:   "Off-by-one in the loop bound",
		Summary: "The loop stops one short.",
	}
	legacy := sha256.Sum256([]byte(strings.Join([]string{
		"acme/widgets",
		"ai.review.3",
		"review",
		"pkg/loop.go",
		"42",
		"Off-by-one in the loop bound",
		"The loop stops one short.",
	}, "\x00")))
	want := fmt.Sprintf("%x", legacy[:16])
	if got := reviewFindingFingerprint("acme/widgets", finding, "pkg/loop.go", 42, "review"); got != want {
		t.Fatalf("legacy fingerprint changed: got %s want %s", got, want)
	}
	// A whitespace-only RuleID is "no rule" and takes the legacy path too.
	blank := finding
	blank.RuleID = "   "
	if got := reviewFindingFingerprint("acme/widgets", blank, "pkg/loop.go", 42, "review"); got != want {
		t.Fatalf("blank RuleID should take the legacy path: got %s want %s", got, want)
	}
	// The two layouts never collide for the same file/line/title: the "rule"
	// domain separator sits where the legacy layout has the finding ID.
	ruled := finding
	ruled.RuleID = "review/wrong-logic"
	if got := reviewFindingFingerprint("acme/widgets", ruled, "pkg/loop.go", 42, "review"); got == want {
		t.Fatalf("rule fingerprint must not equal the legacy fingerprint")
	}
	// Even a pathological rule ID that reproduces the legacy field sequence
	// cannot collide, because the "rule" separator precedes it.
	pathological := finding
	pathological.RuleID = "ai.review.3\x00review"
	if got := reviewFindingFingerprint("acme/widgets", pathological, "pkg/loop.go", 42, "review"); got == want {
		t.Fatalf("rule fingerprint collided with legacy layout via a crafted rule ID")
	}
}

// The history row carries the rule and lane so history can be read per rule
// without re-deriving them from prose. Keys are only present when the finding
// has them: older rows have no such keys and readers treat absence as unknown.
func TestReviewHistoryFindingsCarryRuleAndLaneInPayload(t *testing.T) {
	report := codereview.Report{
		ChangedFiles: []string{"pkg/loop.go"},
		Findings: []codereview.Finding{
			{
				ID:          "ai.review.1",
				RuleID:      "review/wrong-logic",
				Lane:        codereview.LaneAdvisory,
				DemotedFrom: codereview.LaneBlocking,
				Materiality: "high",
				Title:       "t",
				Summary:     "s",
				Strength:    "Strong",
			},
			{
				ID:       "ai.review.2",
				Title:    "no rule",
				Summary:  "s",
				Strength: "Speculative",
			},
		},
	}
	records := reviewHistoryFindings("acme/widgets", report)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	ruled := records[0]
	if ruled.RuleID != "review/wrong-logic" {
		t.Errorf("RuleID = %q, want review/wrong-logic", ruled.RuleID)
	}
	for key, want := range map[string]string{
		"rule_id":      "review/wrong-logic",
		"lane":         codereview.LaneAdvisory,
		"demoted_from": codereview.LaneBlocking,
		"materiality":  "high",
	} {
		if got, _ := ruled.Payload[key].(string); got != want {
			t.Errorf("payload[%q] = %q, want %q", key, got, want)
		}
	}
	// The pre-existing keys are still there.
	if got, _ := ruled.Payload["id"].(string); got != "ai.review.1" {
		t.Errorf("payload[id] = %q", got)
	}
	plain := records[1]
	if plain.RuleID != "" {
		t.Errorf("record without a rule should have empty RuleID, got %q", plain.RuleID)
	}
	for _, key := range []string{"rule_id", "lane", "demoted_from", "materiality"} {
		if _, present := plain.Payload[key]; present {
			t.Errorf("payload[%q] should be absent when the finding has none", key)
		}
	}
	// And the wire shape: ruleId appears only when set.
	encoded, err := json.Marshal(records)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), `"ruleId":"review/wrong-logic"`) {
		t.Errorf("encoded records missing ruleId: %s", encoded)
	}
	if strings.Count(string(encoded), `"ruleId"`) != 1 {
		t.Errorf("ruleId should be omitted when empty: %s", encoded)
	}
}
