package codereview

import (
	"sort"
	"strings"
)

// Story is the review's human-facing "what changed that you now own" lane.
// Nothing in it is wrong: a finding asks for action, a story item ends in
// ownership. It answers the question a reviewer asks after the findings are
// dealt with — what did I just take on by merging this? — and it is what a
// reader remembers a change by six months later.

// Story item categories. Structural, not decorative: the category says what
// kind of ownership the reader is taking on.
const (
	StoryBehaviorDelta  = "behavior_delta"  // an existing path now acts differently
	StoryNewSurface     = "new_surface"     // endpoint, flag, env var, schema, event, dependency that now exists
	StorySemanticShift  = "semantic_shift"  // idempotency, ordering, cache lifetime, transaction scope, error contract changed meaning
	StoryDecision       = "decision"        // approach chosen and alternative rejected, from the session
	StoryExceptionTaken = "exception_taken" // a REVIEW.md exception/suppression this change adds — mandatory, never model-discretionary
	StoryCoverageMove   = "coverage_move"   // what became tested or untested
)

// Materiality levels. Materiality is how much observable behavior moved, not
// how confident anyone is about it.
const (
	StoryMaterialityHigh   = "high"
	StoryMaterialityMedium = "medium"
	StoryMaterialityLow    = "low"
)

// StoryItem is one thing the reader now owns.
type StoryItem struct {
	Headline    string   `json:"headline"` // past → present, one sentence
	Category    string   `json:"category"`
	Materiality string   `json:"materiality"`    // high | medium | low
	File        string   `json:"file,omitempty"` // primary anchor
	Line        int      `json:"line,omitempty"`
	Files       []string `json:"files,omitempty"`       // all touched, when more than one
	DiffHunk    string   `json:"diff_hunk,omitempty"`   // window around the anchor line, attached from the diff
	Consequence string   `json:"consequence"`           // "what changes for you", with numbers where possible
	Asked       string   `json:"asked,omitempty"`       // the user's words, from the session
	Chose       string   `json:"chose,omitempty"`       // decision + rejected alternative
	Watch       []string `json:"watch,omitempty"`       // what to monitor after merge
	PassedRule  string   `json:"passed_rule,omitempty"` // a rule that passed and still speaks (namespaced rule ID)
}

// normalizeStoryCategory folds a model-supplied category onto one of the six
// the schema defines, or "" when it names none of them. Same principle as
// normalizeFindingKind: unknown becomes empty rather than a guess, because
// "the model did not say" is information a renderer must be able to see.
func normalizeStoryCategory(s string) string {
	key := strings.ToLower(strings.TrimSpace(s))
	key = strings.NewReplacer("-", "_", " ", "_").Replace(key)
	switch key {
	case StoryBehaviorDelta, "behaviour_delta", "behavior_change", "behaviour_change", "behavior", "behaviour":
		return StoryBehaviorDelta
	case StoryNewSurface, "surface", "new_api", "new_endpoint", "new_flag":
		return StoryNewSurface
	case StorySemanticShift, "semantics", "semantic", "semantic_change":
		return StorySemanticShift
	case StoryDecision, "design_decision", "tradeoff", "trade_off":
		return StoryDecision
	case StoryExceptionTaken, "exception", "suppression":
		return StoryExceptionTaken
	case StoryCoverageMove, "coverage", "test_coverage", "coverage_change", "tests":
		return StoryCoverageMove
	}
	return ""
}

// normalizeMateriality folds a model-supplied materiality onto high, medium,
// or low, or "" when it is none of them.
func normalizeMateriality(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case StoryMaterialityHigh, "major", "critical":
		return StoryMaterialityHigh
	case StoryMaterialityMedium, "med", "moderate":
		return StoryMaterialityMedium
	case StoryMaterialityLow, "minor":
		return StoryMaterialityLow
	}
	return ""
}

// storyMaterialityRank orders materiality for display: high first, then
// medium, low, and finally items that did not say.
func storyMaterialityRank(materiality string) int {
	switch materiality {
	case StoryMaterialityHigh:
		return 0
	case StoryMaterialityMedium:
		return 1
	case StoryMaterialityLow:
		return 2
	}
	return 3
}

// sortStoryByMateriality orders items high, medium, low, then unlabeled. It is
// stable, so within a level the model's own order — which the prompt asks to
// be by materiality already — is preserved.
func sortStoryByMateriality(items []StoryItem) {
	sort.SliceStable(items, func(i, j int) bool {
		return storyMaterialityRank(items[i].Materiality) < storyMaterialityRank(items[j].Materiality)
	})
}

// MandatoryStoryItems returns the story items no model gets to decide. Today
// there is one: a change that edits the Exceptions section of its own
// REVIEW.md is a change that exempts itself from a rule, and a change never
// gets to do that silently. If REVIEW.md is among the changed files and its
// diff adds or removes lines inside `## Exceptions`, one exception_taken item
// is emitted; otherwise none.
//
// policy is accepted so a repository that names its policy file differently
// can be honored later; today it only confirms the path, and nil is fine.
func MandatoryStoryItems(policy *ReviewPolicy, changedFiles []string, diffs []DiffSnippet) []StoryItem {
	policyPath := reviewPolicyPath
	if policy != nil && strings.TrimSpace(policy.Path) != "" {
		policyPath = strings.TrimSpace(policy.Path)
	}
	if !storyChangedFilesInclude(changedFiles, policyPath) {
		return nil
	}
	for _, snippet := range diffs {
		if !storyPathMatches(snippet.File, policyPath) {
			continue
		}
		if reviewPolicyDiffTouchesExceptions(snippet.Diff) {
			return []StoryItem{{
				Headline:    "This change edits its own REVIEW.md exceptions",
				Category:    StoryExceptionTaken,
				Materiality: StoryMaterialityHigh,
				File:        policyPath,
				Consequence: "A change that exempts itself from a rule needs a human to have seen it. Confirm the exception is intended and scoped.",
			}}
		}
	}
	return nil
}

func storyChangedFilesInclude(files []string, path string) bool {
	for _, file := range files {
		if storyPathMatches(file, path) {
			return true
		}
	}
	return false
}

// storyPathMatches compares a changed-file path with the policy path, tolerant
// of a leading "./" and of backslashes on Windows checkouts.
func storyPathMatches(file, path string) bool {
	clean := func(s string) string {
		s = strings.TrimSpace(strings.ReplaceAll(s, "\\", "/"))
		s = strings.TrimPrefix(s, "./")
		return s
	}
	return clean(file) == clean(path)
}

// reviewPolicyDiffTouchesExceptions reports whether a REVIEW.md diff adds or
// removes lines inside its `## Exceptions` section. It walks the diff tracking
// which section each line falls under — a hunk header's trailing context, a
// context line, or an added/removed heading all move the cursor — and answers
// yes the first time an added or removed line lands while the cursor is in
// Exceptions. Adding or removing the heading itself counts: creating the
// section is taking an exception.
//
// A snippet that carries file content instead of a diff (an untracked
// REVIEW.md) is treated as all-added, matching how the constraints signals
// read the same header: a brand-new REVIEW.md with an Exceptions section is a
// change that adds exceptions.
func reviewPolicyDiffTouchesExceptions(diff string) bool {
	diff = strings.ReplaceAll(diff, "\r\n", "\n")
	if strings.HasPrefix(diff, diffUnavailableContentHeader) {
		for _, line := range strings.Split(diff, "\n") {
			if section, ok := reviewPolicySectionHeading(line); ok && section == "exceptions" {
				return true
			}
		}
		return false
	}
	inExceptions := false
	for _, raw := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(raw, "@@"):
			// The hunk header resets the cursor; git may append the enclosing
			// heading as function context after the second @@.
			inExceptions = false
			if idx := strings.LastIndex(raw, "@@"); idx >= 0 && idx+2 <= len(raw) {
				if section, ok := reviewPolicySectionHeading(strings.TrimSpace(raw[idx+2:])); ok {
					inExceptions = section == "exceptions"
				}
			}
			continue
		case strings.HasPrefix(raw, "+++ ") || strings.HasPrefix(raw, "--- "):
			continue
		case strings.HasPrefix(raw, "diff ") || strings.HasPrefix(raw, "index "):
			inExceptions = false
			continue
		}
		var body string
		changed := false
		switch {
		case strings.HasPrefix(raw, "+"), strings.HasPrefix(raw, "-"):
			body = raw[1:]
			changed = true
		case strings.HasPrefix(raw, " "):
			body = raw[1:]
		default:
			body = raw
		}
		if section, ok := reviewPolicySectionHeading(body); ok {
			inExceptions = section == "exceptions"
			if changed && inExceptions {
				return true
			}
			continue
		}
		if changed && inExceptions && strings.TrimSpace(body) != "" {
			return true
		}
	}
	return false
}

// reviewPolicySectionHeading returns the lowercased text of a `## Heading`
// line, and whether the line is one. Only level-two headings move the section
// cursor, which is the same grammar parseReviewRules reads.
func reviewPolicySectionHeading(line string) (string, bool) {
	match := reviewRuleHeadingPattern.FindStringSubmatch(line)
	if len(match) != 2 {
		return "", false
	}
	return strings.ToLower(strings.TrimSpace(match[1])), true
}

// attachStoryDiffHunks gives each anchored story item its unified-diff window,
// extracted from the same snippets the reviewer read. Same helper the
// constraints report uses for findings, so both lanes show the same window.
func attachStoryDiffHunks(diffsByFile map[string]string, items []StoryItem) {
	for i := range items {
		item := &items[i]
		if item.DiffHunk != "" || item.File == "" || item.Line <= 0 {
			continue
		}
		diff, ok := diffsByFile[item.File]
		if !ok {
			continue
		}
		item.DiffHunk = reviewDiffHunkForLine(diff, item.Line)
	}
}

// assembleReportStory builds the Story lane a Report carries: the mandatory
// items first — they are the ones no model decides — then whatever the model
// contributed, then diff hunks for everything with an anchor. Returns nil when
// there is nothing to say, so a report without a story serializes without one.
func assembleReportStory(policy *ReviewPolicy, changedFiles []string, diffs []DiffSnippet, model []StoryItem) []StoryItem {
	story := append(MandatoryStoryItems(policy, changedFiles, diffs), model...)
	if len(story) == 0 {
		return nil
	}
	diffsByFile := make(map[string]string, len(diffs))
	for _, snippet := range diffs {
		diffsByFile[snippet.File] = snippet.Diff
	}
	attachStoryDiffHunks(diffsByFile, story)
	return story
}
