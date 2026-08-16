package codereview

import (
	"regexp"
	"strings"
	"testing"
)

// The pack's slugs are permanent: findings, suppressions, and history rows
// reference them by name. This list is the contract; a rule that disappears
// or is renamed fails here before it can orphan a reference.
var recommendedPackSlugs = map[string]bool{
	// Tier 1. The value is whether the rule is advisory.
	"no-secrets-in-code":     false,
	"no-secrets-in-logs":     false,
	"no-injection-sinks":     false,
	"no-swallowed-errors":    false,
	"tests-can-fail":         false,
	"no-placeholder-code":    false,
	"comments-match-code":    true,
	"no-unbounded-work":      true,
	"cleanup-on-failure":     true,
	"handles-missing-values": false,
	// Tier 2.
	"reuse-before-rewrite":    true,
	"no-duplicated-blocks":    true,
	"no-invented-packages":    false,
	"no-invented-apis":        false,
	"apis-used-as-documented": true,
	"endpoints-enforce-authz": false,
	"no-dead-code":            true,
	"no-speculative-layers":   true,
	"respects-module-seams":   true,
	"no-race-hazards":         false,
	"retries-are-idempotent":  false,
	"no-vulnerable-deps":      false,
	// Tier 3.
	"scope-matches-intent":     false,
	"claims-match-diff":        false,
	"no-weakened-checks":       false,
	"no-silent-regressions":    false,
	"no-unguarded-destruction": false,
}

var recommendedPackSlugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

func TestRecommendedPackParsesCleanly(t *testing.T) {
	rules, err := RecommendedPack()
	if err != nil {
		t.Fatalf("RecommendedPack() error = %v", err)
	}
	if len(rules) != len(recommendedPackSlugs) {
		t.Fatalf("pack has %d rules, want %d", len(rules), len(recommendedPackSlugs))
	}
	if _, _, _, diagnostics := parseReviewRules(recommendedPackSource); len(diagnostics) != 0 {
		t.Fatalf("pack parse diagnostics = %#v, want none", diagnostics)
	}
	if RecommendedPackVersion != "2" {
		t.Fatalf("RecommendedPackVersion = %q, want 2", RecommendedPackVersion)
	}
}

func TestRecommendedPackSlugsAreThePermanentSet(t *testing.T) {
	rules, err := RecommendedPack()
	if err != nil {
		t.Fatalf("RecommendedPack() error = %v", err)
	}
	seen := map[string]int{}
	for _, rule := range rules {
		seen[rule.Slug]++
		if !strings.HasPrefix(rule.ID, "gx:recommended/") {
			t.Errorf("rule %q ID = %q, want gx:recommended/ namespace", rule.Slug, rule.ID)
		}
		if rule.ID != "gx:recommended/"+rule.Slug {
			t.Errorf("rule ID %q does not match slug %q", rule.ID, rule.Slug)
		}
		if !recommendedPackSlugPattern.MatchString(rule.Slug) {
			t.Errorf("slug %q is not kebab-case", rule.Slug)
		}
		if len(rule.Slug) > 24 {
			t.Errorf("slug %q is %d chars, want at most 24", rule.Slug, len(rule.Slug))
		}
		if !ruleIDPattern.MatchString(rule.ID) {
			t.Errorf("rule ID %q does not match ruleIDPattern", rule.ID)
		}
		if strings.TrimSpace(rule.Text) == "" {
			t.Errorf("rule %q has no text", rule.Slug)
		}
		if strings.TrimSpace(rule.Name) == "" || strings.HasPrefix(strings.ToLower(rule.Name), "advisory:") {
			t.Errorf("rule %q name = %q, want the heading without the Advisory: prefix", rule.Slug, rule.Name)
		}
		wantAdvisory, known := recommendedPackSlugs[rule.Slug]
		if !known {
			t.Errorf("pack declares %q, which is not in the permanent slug set", rule.Slug)
			continue
		}
		if rule.Advisory != wantAdvisory {
			t.Errorf("rule %q advisory = %v, want %v", rule.Slug, rule.Advisory, wantAdvisory)
		}
	}
	for slug := range recommendedPackSlugs {
		if seen[slug] != 1 {
			t.Errorf("slug %q appears %d time(s) in the pack, want exactly once", slug, seen[slug])
		}
	}
}

func TestRecommendedPackRulesAreKnownRuleIDs(t *testing.T) {
	defs := AllRuleDefs(nil)
	known := KnownRuleIDs(defs...)
	if got := NormalizeRuleID("gx:recommended/no-swallowed-errors", known); got != "gx:recommended/no-swallowed-errors" {
		t.Fatalf("pack rule not accepted by NormalizeRuleID: got %q", got)
	}
	if got := NormalizeRuleID("GX:Recommended/No-Swallowed-Errors", known); got != "gx:recommended/no-swallowed-errors" {
		t.Fatalf("case-folded pack rule = %q, want canonical", got)
	}
	for _, def := range defs {
		if strings.TrimSpace(def.Summary) == "" {
			t.Errorf("rule %q has an empty summary", def.ID)
		}
		if len(def.Summary) > 160 {
			t.Errorf("rule %q summary is %d bytes, want at most 160", def.ID, len(def.Summary))
		}
	}
	// AllRuleDefs layers REVIEW.md rules after the pack.
	policy := &ReviewPolicy{Rules: []PolicyRule{{ID: "REVIEW.md/no-reinvented-utils", Slug: "no-reinvented-utils", Text: "Reuse the shared helper. More detail here."}}}
	all := AllRuleDefs(policy)
	if len(all) != len(defs)+1 {
		t.Fatalf("AllRuleDefs(policy) has %d rules, want %d", len(all), len(defs)+1)
	}
	last := all[len(all)-1]
	if last.ID != "REVIEW.md/no-reinvented-utils" || last.Summary != "Reuse the shared helper." {
		t.Fatalf("repo rule def = %#v", last)
	}
}
