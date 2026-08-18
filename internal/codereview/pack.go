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
const RecommendedPackVersion = "3"

const recommendedPackNamespace = "gx:recommended/"

// recommendedPackAliases maps a slug the pack no longer declares to the rule
// that absorbed it. v3 merged four v2 rules into two: the two secrets rules
// became one, and the two duplication rules became one.
//
// The merge is the reason this table exists rather than the rename being made
// in place. A slug is not just a heading — it is what a suppression in a
// repository's REVIEW.md was written against and what every history row already
// recorded. Deleting `no-secrets-in-code` would silently un-suppress a rule
// somebody deliberately turned off, and orphan every past finding that carried
// it. Aliases are one-way and permanent: an alias may never be reused as a live
// slug, and a rule that splits later needs its own new names, not a reversal of
// these.
var recommendedPackAliases = map[string]string{
	"no-secrets-in-code":   "no-secrets",
	"no-secrets-in-logs":   "no-secrets",
	"reuse-before-rewrite": "dont-repeat-yourself",
	"no-duplicated-blocks": "dont-repeat-yourself",
}

// RecommendedPackAliases returns the retired-to-current mapping in namespaced
// form, so a caller resolving a stored rule ID does not have to know that the
// pack's slugs are bare.
func RecommendedPackAliases() map[string]string {
	out := make(map[string]string, len(recommendedPackAliases))
	for retired, current := range recommendedPackAliases {
		out[recommendedPackNamespace+retired] = recommendedPackNamespace + current
	}
	return out
}

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
