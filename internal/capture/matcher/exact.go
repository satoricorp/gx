package matcher

import (
	"hash/fnv"
	"strings"
)

func tier1Match(newText, addedText string, n int) bool {
	if n <= 0 {
		n = 5
	}
	newText = normalizeExact(newText)
	addedText = normalizeExact(addedText)
	if newText == "" || addedText == "" {
		return false
	}
	need := ngramSet(newText, n)
	if len(need) == 0 {
		return false
	}
	have := ngramSet(addedText, n)
	for hash := range need {
		if _, ok := have[hash]; !ok {
			return false
		}
	}
	return true
}

func normalizeExact(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r == '\n' || r == ' ' || r == '\t' {
			b.WriteByte(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func ngramSet(text string, n int) map[uint64]struct{} {
	tokens := strings.Fields(text)
	set := map[uint64]struct{}{}
	if len(tokens) >= n {
		for i := 0; i+n <= len(tokens); i++ {
			window := strings.Join(tokens[i:i+n], " ")
			set[hashString(window)] = struct{}{}
		}
		return set
	}
	if len(text) >= n {
		for i := 0; i+n <= len(text); i++ {
			set[hashString(text[i:i+n])] = struct{}{}
		}
	}
	return set
}

// containmentScore is the fraction of needle's n-grams present in haystack.
//
// Asymmetric on purpose, and that is what separates it from fuzzyScore: a
// command is usually far larger than the hunk it wrote — a heredoc holding a
// whole file, a script doing several edits at once — so Jaccard similarity
// scores it near zero however completely the hunk is contained. The question
// here is only "is the hunk inside the command", so the command's other
// content must not count against it.
func containmentScore(needle, haystack string, n int) float64 {
	needle = normalizeExact(needle)
	haystack = normalizeExact(haystack)
	if needle == "" || haystack == "" {
		return 0
	}
	need := ngramSet(needle, n)
	if len(need) == 0 {
		return 0
	}
	have := ngramSet(haystack, n)
	hit := 0
	for hash := range need {
		if _, ok := have[hash]; ok {
			hit++
		}
	}
	return float64(hit) / float64(len(need))
}

// countTokens counts whitespace-separated tokens in already-normalized text.
func countTokens(normalized string) int {
	return len(strings.Fields(normalized))
}

func hashString(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
