package codereview

import (
	_ "embed"
	"fmt"
	"strings"
	"sync"
)

// The recommended pack is the built-in set of named rules every review runs
// with. It is written in the REVIEW.md grammar and read by the same parser, so
// a repository that wants to see how a rule is phrased can open the file, and
// a rule added to the pack cannot behave differently from one a repository
// declares itself. Only the namespace differs: pack rules are
// gx:recommended/<slug>, repo rules are REVIEW.md/<slug>.
//
// Slugs are permanent once released — suppressions and history rows reference
// them — so a rule that needs renaming gets an alias, not an edit.

// RecommendedPackVersion is the pack revision a REVIEW.md `## Extends:
// gx:recommended@<n>` line refers to.
const RecommendedPackVersion = "2"

const recommendedPackNamespace = "gx:recommended/"

//go:embed packs/recommended.md
var recommendedPackSource string

var (
	recommendedPackOnce  sync.Once
	recommendedPackRules []PolicyRule
	recommendedPackErr   error
)

// RecommendedPack returns the built-in rules, parsed once. The error is
// reserved for a pack that fails its own grammar (a diagnostic from the
// parser, or no rules at all) — a build-time mistake, surfaced rather than
// silently reviewing with fewer rules.
func RecommendedPack() ([]PolicyRule, error) {
	recommendedPackOnce.Do(func() {
		rules, _, _, diagnostics := parseReviewRules(recommendedPackSource)
		if len(diagnostics) > 0 {
			recommendedPackErr = fmt.Errorf("gx:recommended pack did not parse cleanly: %s", strings.Join(diagnostics, "; "))
			return
		}
		if len(rules) == 0 {
			recommendedPackErr = fmt.Errorf("gx:recommended pack declares no rules")
			return
		}
		for i := range rules {
			rules[i].ID = recommendedPackNamespace + rules[i].Slug
		}
		recommendedPackRules = rules
	})
	if recommendedPackErr != nil {
		return nil, recommendedPackErr
	}
	return append([]PolicyRule(nil), recommendedPackRules...), nil
}
