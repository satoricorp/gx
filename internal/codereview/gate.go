package codereview

import (
	"fmt"
	"sort"
	"strings"
)

// FailOnLevel is the severity threshold an automated gate fails at. The levels
// are the finding strengths gx already speaks ("Blocking", "Strong", "Worth
// exploring", "Speculative"), plus "none" to disable the gate and "any" to
// catch every surviving finding.
type FailOnLevel string

const (
	FailOnNone           FailOnLevel = "none"
	FailOnAny            FailOnLevel = "any"
	FailOnSpeculative    FailOnLevel = "speculative"
	FailOnWorthExploring FailOnLevel = "worth-exploring"
	FailOnStrong         FailOnLevel = "strong"
	FailOnBlocking       FailOnLevel = "blocking"
)

// failOnRanks maps a level onto strengthRank, where 0 is the most severe.
// A level fails on every finding whose rank is at or above the threshold.
var failOnRanks = map[FailOnLevel]int{
	FailOnAny:            3,
	FailOnSpeculative:    3,
	FailOnWorthExploring: 2,
	FailOnStrong:         1,
	FailOnBlocking:       0,
}

// FailOnLevels lists the accepted values, most permissive first, for help text
// and error messages.
func FailOnLevels() []string {
	levels := []string{string(FailOnNone), string(FailOnAny)}
	named := make([]string, 0, len(failOnRanks))
	for level := range failOnRanks {
		if level == FailOnAny {
			continue
		}
		named = append(named, string(level))
	}
	sort.Slice(named, func(i, j int) bool {
		return failOnRanks[FailOnLevel(named[i])] > failOnRanks[FailOnLevel(named[j])]
	})
	return append(levels, named...)
}

// ParseFailOnLevel accepts the level names above, case-insensitively, and also
// tolerates the spaced strength spelling ("worth exploring").
func ParseFailOnLevel(value string) (FailOnLevel, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, " ", "-")
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if normalized == "" {
		return FailOnNone, nil
	}
	level := FailOnLevel(normalized)
	if level == FailOnNone {
		return FailOnNone, nil
	}
	if _, ok := failOnRanks[level]; ok {
		return level, nil
	}
	return "", fmt.Errorf("unsupported --fail-on level %q (want one of: %s)", value, strings.Join(FailOnLevels(), ", "))
}

// Enabled reports whether the level gates anything at all.
func (l FailOnLevel) Enabled() bool {
	_, ok := failOnRanks[l]
	return ok
}

// GateFailures returns the findings at or above the threshold. It returns
// nothing when the level is disabled, and nothing when the review never
// inspected any code: a gate must not report "no issues at or above X" for a
// diff it never read. Callers handle that case through Report.ReviewMode.
func (r Report) GateFailures(level FailOnLevel) []Finding {
	threshold, ok := failOnRanks[level]
	if !ok || r.ReviewMode == ReviewModeNone {
		return nil
	}
	var out []Finding
	for _, finding := range r.Findings {
		if strengthRank(finding.Strength) <= threshold {
			out = append(out, finding)
		}
	}
	return out
}
