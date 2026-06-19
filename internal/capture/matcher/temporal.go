package matcher

import "time"

func temporalMatch(eventTS, commitTS int64, window time.Duration) bool {
	if eventTS <= 0 || commitTS <= 0 {
		return false
	}
	ev := time.UnixMilli(eventTS)
	commit := time.UnixMilli(commitTS)
	diff := ev.Sub(commit)
	if diff < 0 {
		diff = -diff
	}
	return diff <= window
}
