package codereview

import "sort"

// lineSpan is an inclusive 1-based line range inside a single file.
//
// Retrieval returns one span per indexed chunk, and the same file is routinely
// covered by several chunks whose ranges touch or overlap. Emitting them as-is
// spends one context slot per chunk and shows the model the same lines twice,
// so spans for a file are merged before they become snippets.
type lineSpan struct {
	Start int
	End   int
}

// mergeLineSpans collapses overlapping and adjacent spans into the smallest set
// of spans covering the same lines. The input order is irrelevant; the result is
// sorted by Start.
func mergeLineSpans(spans []lineSpan) []lineSpan {
	if len(spans) < 2 {
		return spans
	}
	sort.Slice(spans, func(i, j int) bool {
		if spans[i].Start != spans[j].Start {
			return spans[i].Start < spans[j].Start
		}
		return spans[i].End < spans[j].End
	})

	merged := spans[:1]
	for _, span := range spans[1:] {
		last := &merged[len(merged)-1]
		if span.Start <= last.End {
			if span.End > last.End {
				last.End = span.End
			}
			continue
		}
		merged = append(merged, span)
	}
	return merged
}

// spanCoverage reports how many distinct lines a set of merged spans covers. It
// is used to decide whether a file is well enough covered by retrieval to skip
// reading the whole file from disk.
func spanCoverage(spans []lineSpan) int {
	total := 0
	for _, span := range mergeLineSpans(spans) {
		total += span.End - span.Start
	}
	return total
}
