package matcher

import (
	"time"

	"github.com/satoricorp/gx/internal/capture"
)

const (
	TierExact    = 1
	TierFuzzy    = 2
	TierTemporal = 3
	// TierCommand is a hunk whose added lines appear inside the text of a
	// command the agent ran. It is authorship evidence — the agent produced
	// those bytes — but weaker than an edit event, which names the file it
	// wrote. A command only shows text passing through it, so the content match
	// is doing all the work and the file path is inferred from the hunk rather
	// than declared by the tool. Numbered after the temporal tier because it
	// was added later and the values are reported; it is stronger than tier 3,
	// not weaker, and unlike tier 3 it counts as authorship.
	TierCommand = 4
)

// commandMatchMinTokens is the smallest hunk a command may claim: enough tokens
// for tier1Match to use real token n-grams.
//
// Below n tokens tier1Match falls back to character shingles, which are
// permissive enough that a command merely mentioning a short line would claim
// it — `grep -n "return nil"` must not make the agent the author of every
// `return nil` in the diff. At or above it, claiming requires every n-token
// window of the hunk to appear in the command, which is a much harder
// coincidence than it sounds: normalizeExact strips punctuation, so a six-line
// Go function is only about seven tokens, and a command has to carry all of
// them in order.
func commandMatchMinTokens(cfg Config) int {
	if cfg.NgramSize > 0 {
		return cfg.NgramSize
	}
	return 5
}

// Config controls matcher thresholds.
type Config struct {
	FuzzyThreshold float64
	TemporalWindow time.Duration
	NgramSize      int
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
	// CommandHunkIndexes are hunks claimed from the text of a command the agent
	// ran rather than from a file-edit event. See TierCommand.
	CommandHunkIndexes map[int]struct{}
	Tier3Pairs         int
	MatchedEvents    map[int]int // event index -> best tier (1 or 2)
	// RelevantEvents records every event whose content matched some hunk
	// (tier 1 or 2), before hunk-claim dedup. MatchedEvents only keeps the
	// single claiming event per hunk, which under-reports a session whose
	// edits landed on hunks another session's event claimed first.
	RelevantEvents map[int]int // event index -> best tier (1 or 2)
	Outcomes       []MatchOutcome
}

// AuthorshipHunkCount returns hunks attributed to an agent: matched by content
// from a file-edit event (tiers 1 and 2) or from a command the agent ran.
// Tier 3 is excluded — proximity in time is not authorship.
func (r Result) AuthorshipHunkCount() int {
	return len(r.Tier1HunkIndexes) + len(r.Tier2HunkIndexes) + len(r.CommandHunkIndexes)
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
		Tier1HunkIndexes:   map[int]struct{}{},
		Tier2HunkIndexes:   map[int]struct{}{},
		CommandHunkIndexes: map[int]struct{}{},
		MatchedEvents:      map[int]int{},
		RelevantEvents:     map[int]int{},
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
			result.RelevantEvents[ai] = bestTier
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

	// Last, so that it sees the tier 1 and 2 sets already folded from Outcomes
	// and can skip every hunk they claimed. Running it earlier reads two empty
	// sets and re-claims hunks that had a real edit event behind them.
	matchCommands(events, hunks, cfg, &result)
	return result
}

func eligibleEvent(ev capture.SessionEvent) bool {
	return ev.IsEditEvent() && ev.FilePath != "" && len(trimSpace(ev.NewText)) > 0
}

// matchCommands claims hunks whose added lines appear inside the text of a
// command the agent ran.
//
// The direction is the inverse of the edit-event match and that is the whole
// idea. An edit event carries the text it wrote, so the question is whether
// that text is in the hunk. A command carries whatever passed through it —
// often a heredoc holding an entire file — so the question is whether the hunk
// is in the command. tier1Match already answers "is every n-gram of the first
// argument present in the second", so it answers both with the arguments
// swapped.
//
// It runs only over hunks no edit event claimed. A hunk with a real edit event
// behind it has better evidence already, and this pass must never downgrade it.
// There is no file-path join available — a command does not say what it wrote —
// so the temporal window does the narrowing that a path would, and
// minCommandMatchTokens keeps the content match from being satisfied by a
// passing mention.
func matchCommands(events []capture.SessionEvent, hunks []HunkRef, cfg Config, result *Result) {
	for hi, hunk := range hunks {
		if _, ok := result.Tier1HunkIndexes[hi]; ok {
			continue
		}
		if _, ok := result.Tier2HunkIndexes[hi]; ok {
			continue
		}
		addedText := joinLines(hunk.AddedLines)
		if countTokens(normalizeExact(addedText)) < commandMatchMinTokens(cfg) {
			continue
		}
		for ai, ev := range events {
			if !ev.IsCommandEvent() {
				continue
			}
			if !temporalMatch(ev.TS, hunk.CommitTime, cfg.TemporalWindow) {
				continue
			}
			// A threshold rather than total containment. Measured against a
			// real push: of the hunks this pass exists to catch, one scored
			// 1.00 and the next 0.99 — a single n-gram short, because the
			// committed bytes are the command's text after gofmt and after
			// whatever else the file already held. Demanding every n-gram
			// rejected work the agent plainly did. The bar is the same
			// FuzzyThreshold tier 2 uses, so both inexact tiers relax by the
			// same amount rather than by two separately tuned numbers.
			score := containmentScore(addedText, ev.Command, cfg.NgramSize)
			if score < cfg.FuzzyThreshold {
				continue
			}
			result.CommandHunkIndexes[hi] = struct{}{}
			result.RelevantEvents[ai] = TierCommand
			if _, seen := result.MatchedEvents[ai]; !seen {
				result.MatchedEvents[ai] = TierCommand
			}
			result.Outcomes = append(result.Outcomes, MatchOutcome{
				EventIndex: ai,
				HunkIndex:  hi,
				Tier:       TierCommand,
				Score:      score,
			})
			break
		}
	}
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
