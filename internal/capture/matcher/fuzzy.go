package matcher

import (
	"strings"
	"unicode"
)

var stopTokens = map[string]struct{}{
	"this": {}, "that": {}, "with": {}, "from": {}, "have": {},
	"will": {}, "your": {}, "they": {}, "them": {}, "then": {},
	"when": {}, "what": {}, "were": {}, "been": {}, "into": {},
}

func fuzzyScore(a, b string) float64 {
	ta := tokenSet(a)
	tb := tokenSet(b)
	if len(ta) == 0 || len(tb) == 0 {
		return 0
	}
	inter := 0
	union := map[string]struct{}{}
	for t := range ta {
		union[t] = struct{}{}
		if _, ok := tb[t]; ok {
			inter++
		}
	}
	for t := range tb {
		union[t] = struct{}{}
	}
	if len(union) == 0 {
		return 0
	}
	return float64(inter) / float64(len(union))
}

func tokenSet(text string) map[string]struct{} {
	text = strings.ToLower(text)
	var tokens []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() >= 4 {
			tok := cur.String()
			if _, stop := stopTokens[tok]; !stop {
				tokens = append(tokens, tok)
			}
		}
		cur.Reset()
	}
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cur.WriteRune(r)
		} else {
			flush()
		}
	}
	flush()
	set := map[string]struct{}{}
	for _, t := range tokens {
		set[t] = struct{}{}
	}
	return set
}
