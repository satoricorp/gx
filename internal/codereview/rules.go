package codereview

import (
	"regexp"
	"strings"
)

// Rule identity.
//
// A finding's ID (ai.review.3, bedrock-a.ai.review.3) is positional — "the
// third thing this leg said this run" — and exists so mergeFindings can
// de-duplicate within one review. It says nothing about which *rule* the
// finding is an instance of, so two runs that flag the same problem on the
// same line produce two unrelated findings, and nothing can ask "how often
// does this rule fire, and how often is it overridden."
//
// RuleID is that missing identity: a stable, namespaced name for the check a
// finding enforces. It is deliberately a separate field from ID because ID is
// the de-duplication key and two findings from one rule must both survive.
//
// Namespaces:
//
//	review/<class>          one of the defect classes the review prompt asks for
//	gx:recommended/<rule>   a rule from the built-in pack (added when the pack ships)
//	REVIEW.md/<rule>        a rule the repository wrote in its own REVIEW.md
//
// The name after the slash follows one grammar: kebab-case, 2-4 words, at most
// 24 characters, either `no-<bad-pattern>` or `<subject>-<verb>-<object>`.
// Names are permanent once released — suppressions and history rows reference
// them — so renaming after release needs an alias table, not an edit.

// reviewClassNamespace prefixes the open-ended defect classes the review prompt
// asks the model to prioritize. These are not pattern rules; they are the
// judgment classes that were already in the prompt, given names so that "no
// ID, no finding" holds before any pack exists.
const reviewClassNamespace = "review/"

// reviewDefectClasses are the classes the review prompt names, in the order
// the prompt lists them. The slug after review/ is what the model must return.
var reviewDefectClasses = []RuleDef{
	{ID: "review/wrong-logic", Summary: "a value, condition, or control flow that computes the wrong thing"},
	{ID: "review/api-contract-misuse", Summary: "a call a library, framework, or interface does not support, or that violates its documented semantics"},
	{ID: "review/missing-value-handling", Summary: "null, empty, or absent input on a path the code actually receives"},
	{ID: "review/security-reachable-path", Summary: "a security or auth issue demonstrable on a reachable path"},
	{ID: "review/race-or-idempotency", Summary: "a race or idempotency hazard with a nameable interleaving"},
	{ID: "review/data-correctness", Summary: "data written, read, or transformed incorrectly"},
	{ID: "review/locale-i18n", Summary: "locale or internationalization correctness"},
	{ID: "review/hardening", Summary: "correct today but fragile — the failure needs a hypothetical future change or an input nothing currently sends"},
}

// RuleDef is the minimum a rule needs to be referenced: its namespaced ID and
// a one-line summary the prompt can show the model. Pack and REVIEW.md rules
// carry more (full text, examples, exceptions) and are layered on later.
type RuleDef struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

// ruleIDPattern is the whole namespaced form: a namespace of letters, digits,
// dots and colons, a slash, then a kebab-case slug.
var ruleIDPattern = regexp.MustCompile(`^[a-z0-9.:_-]+/[a-z0-9]+(?:-[a-z0-9]+)*$`)

// KnownRuleIDs returns the set of rule IDs a finding is allowed to carry:
// the defect classes plus whatever extra rules the caller supplies (pack rules,
// REVIEW.md rules). Anything outside this set is dropped by NormalizeRuleID
// rather than passed through, so a model that invents a plausible-looking name
// cannot mint a rule.
//
// Keys are stored lowercased and matching is case-insensitive, because the
// canonical rendered form of a repo rule is "REVIEW.md/<slug>" while a model
// echoing it may fold the case. The value is the canonical spelling to restore.
func KnownRuleIDs(extra ...RuleDef) map[string]string {
	known := make(map[string]string, len(reviewDefectClasses)+len(extra))
	for _, def := range reviewDefectClasses {
		known[strings.ToLower(def.ID)] = def.ID
	}
	for _, def := range extra {
		if id := strings.TrimSpace(def.ID); id != "" {
			known[strings.ToLower(id)] = id
		}
	}
	// Retired pack slugs resolve to the rule that absorbed them, so a
	// suppression or a history row written against the old name keeps working
	// after a merge. Registered only when the successor is actually in force —
	// an alias to a rule this review is not running is not a rule either.
	for retired, current := range RecommendedPackAliases() {
		if canonical, ok := known[strings.ToLower(current)]; ok {
			known[strings.ToLower(retired)] = canonical
		}
	}
	return known
}

// NormalizeRuleID folds a model-supplied rule_id onto a known rule ID, or
// returns "" when it does not name one. Same principle as normalizeFindingKind:
// empty stays empty and unknown becomes empty. "The model did not name a rule"
// is information — a finding without a rule is an open-ended judgment, and it
// must not be dressed up as an enforced check.
//
// It accepts the bare slug for a review class (`wrong-logic`) as well as the
// namespaced form (`review/wrong-logic`), because the model is far more
// reliable at echoing a short slug than a slash-qualified one. The returned
// string is the canonical spelling registered in known, not the model's.
func NormalizeRuleID(raw string, known map[string]string) string {
	id := strings.ToLower(strings.TrimSpace(raw))
	if id == "" {
		return ""
	}
	if canonical, ok := known[id]; ok && ruleIDPattern.MatchString(id) {
		return canonical
	}
	if !strings.Contains(id, "/") {
		if canonical, ok := known[reviewClassNamespace+id]; ok {
			return canonical
		}
	}
	return ""
}

// ReviewDefectClasses exposes the classes for prompts and tests.
func ReviewDefectClasses() []RuleDef {
	return append([]RuleDef(nil), reviewDefectClasses...)
}

// reviewClassPromptLine renders the class list for the developer prompt: the
// slug the model must return, then the definition it must match.
func reviewClassPromptLine() string {
	var b strings.Builder
	b.WriteString("Set rule_id on every recommendation to the defect class it belongs to, using the exact slug: ")
	for i, def := range reviewDefectClasses {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(strings.TrimPrefix(def.ID, reviewClassNamespace))
		b.WriteString(" = ")
		b.WriteString(def.Summary)
	}
	b.WriteString(". When the finding is an instance of a named rule listed under rules, use that rule's exact id instead. Omit rule_id only when the finding fits none of these; never invent a name.")
	return b.String()
}
