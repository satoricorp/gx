package matcher

import (
	"time"

	"github.com/satoricorp/gx/internal/capture"
)

const (
	TierExact    = 1
	TierFuzzy    = 2
	TierTemporal = 3
)

// Config controls matcher thresholds.
type Config struct {
	FuzzyThreshold  float64
	TemporalWindow  time.Duration
	NgramSize       int
}

// DefaultConfig returns WP-0 default matcher settings.
func DefaultConfig() Config {
	return Config{
		FuzzyThreshold: 0.7,
		TemporalWindow: 30 * time.Minute,
		NgramSize:      5,
	}
}

// HunkRef identifies one eligible commit hunk.
type HunkRef struct {
	Index      int
	CommitSHA  string
	FilePath   string
	AddedLines []string
	CommitTime int64
}

// MatchOutcome records how one agent event matched hunks.
type MatchOutcome struct {
	EventIndex int
	HunkIndex  int
	Tier       int
	Score      float64
}

// Result aggregates tiered match statistics.
type Result struct {
	Tier1HunkIndexes map[int]struct{}
	Tier2HunkIndexes map[int]struct{}
	Tier3Pairs       int
	MatchedEvents    map[int]int // event index -> best tier (1 or 2)
	Outcomes         []MatchOutcome
}

// AuthorshipHunkCount returns hunks matched by Tier 1 or 2 only.
func (r Result) AuthorshipHunkCount() int {
	return len(r.Tier1HunkIndexes) + len(r.Tier2HunkIndexes)
}

// Match runs three-tier matching over eligible events and hunks.
func Match(events []capture.SessionEvent, hunks []HunkRef, cfg Config) Result {
	if cfg.FuzzyThreshold <= 0 {
		cfg = DefaultConfig()
	}
	if cfg.NgramSize <= 0 {
		cfg.NgramSize = 5
	}
	if cfg.TemporalWindow <= 0 {
		cfg.TemporalWindow = 30 * time.Minute
	}

	result := Result{
		Tier1HunkIndexes: map[int]struct{}{},
		Tier2HunkIndexes: map[int]struct{}{},
		MatchedEvents:    map[int]int{},
	}

	claimed := map[int]int{} // hunk index -> tier (lower wins)

	for ai, ev := range events {
		if !eligibleEvent(ev) {
			continue
		}
		bestTier := 0
		bestHunk := -1
		bestScore := 0.0

		for hi, hunk := range hunks {
			if hunk.FilePath != ev.FilePath {
				continue
			}
			addedText := joinLines(hunk.AddedLines)
			if tier1Match(ev.NewText, addedText, cfg.NgramSize) {
				if bestTier == 0 || bestTier > TierExact {
					bestTier = TierExact
					bestHunk = hi
					bestScore = 1.0
				}
				continue
			}
			score := fuzzyScore(ev.NewText, addedText)
			if score > cfg.FuzzyThreshold && (bestTier == 0 || bestTier > TierFuzzy) {
				bestTier = TierFuzzy
				bestHunk = hi
				bestScore = score
			}
		}

		if bestTier == TierExact || bestTier == TierFuzzy {
			if prev, ok := claimed[bestHunk]; !ok || bestTier < prev {
				claimed[bestHunk] = bestTier
				result.MatchedEvents[ai] = bestTier
				result.Outcomes = append(result.Outcomes, MatchOutcome{
					EventIndex: ai,
					HunkIndex:  bestHunk,
					Tier:       bestTier,
					Score:      bestScore,
				})
			}
			continue
		}

		for hi, hunk := range hunks {
			if hunk.FilePath != ev.FilePath {
				continue
			}
			if temporalMatch(ev.TS, hunk.CommitTime, cfg.TemporalWindow) {
				result.Tier3Pairs++
				result.Outcomes = append(result.Outcomes, MatchOutcome{
					EventIndex: ai,
					HunkIndex:  hi,
					Tier:       TierTemporal,
				})
			}
		}
	}

	for _, outcome := range result.Outcomes {
		if outcome.Tier == TierExact {
			result.Tier1HunkIndexes[outcome.HunkIndex] = struct{}{}
			delete(result.Tier2HunkIndexes, outcome.HunkIndex)
		} else if outcome.Tier == TierFuzzy {
			if _, inT1 := result.Tier1HunkIndexes[outcome.HunkIndex]; !inT1 {
				result.Tier2HunkIndexes[outcome.HunkIndex] = struct{}{}
			}
		}
	}
	return result
}

func eligibleEvent(ev capture.SessionEvent) bool {
	return ev.IsEditEvent() && ev.FilePath != "" && len(trimSpace(ev.NewText)) > 0
}

func joinLines(lines []string) string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = trimSpace(line)
	}
	joined := ""
	for i, line := range out {
		if i > 0 {
			joined += "\n"
		}
		joined += line
	}
	return joined
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
