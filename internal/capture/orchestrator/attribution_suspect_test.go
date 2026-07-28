package orchestrator

import "testing"

// Capture parses the agent's transcript while the agent is still writing it, so
// a push made mid-session can read a file that is missing the very work being
// pushed. The run then reports zero coverage, which is indistinguishable from a
// change written by hand.
//
// Truncating a real 2.5MB transcript reproduces it: at ~70% of its eventual
// size a ref range scored 0.0% coverage with 0 sessions, at 80% it scored
// 33.3%, and at 90% it scored the full 50.0% — same code, same commits, same
// invocation. Committed work plus sessions on disk plus zero links is that
// shape, and it must not print as an answer.
func TestAttributionSuspect(t *testing.T) {
	for _, tc := range []struct {
		name                                 string
		eligibleHunks, discovered, hunkLinks int
		want                                 bool
	}{
		{"work and sessions but nothing linked", 6, 14, 0, true},
		{"linked something", 6, 14, 3, false},
		{"no sessions on disk at all", 6, 0, 0, false},
		{"nothing committed in range", 0, 14, 0, false},
		{"nothing anywhere", 0, 0, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := attributionSuspect(tc.eligibleHunks, tc.discovered, tc.hunkLinks); got != tc.want {
				t.Fatalf("attributionSuspect(%d, %d, %d) = %v, want %v",
					tc.eligibleHunks, tc.discovered, tc.hunkLinks, got, tc.want)
			}
		})
	}
}
