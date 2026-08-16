package cli

import "testing"

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
