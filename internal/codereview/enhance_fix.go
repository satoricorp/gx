package codereview

import (
	"fmt"
	"strings"
)

// The enhance surface: one fix, written to be handed to a coding model.
//
// `gx review` answers "what is wrong with this change" — every finding, in two
// lanes, with a verdict. That is the right shape for a human deciding whether
// to ship, and the wrong shape for an agent: a model handed twelve findings
// spreads its attention across twelve, and the one that actually blocks the
// merge gets the same weight as the twelfth-ranked nit.
//
// `gx enhance` answers the other question — "what is the single most valuable
// thing to change right now" — and answers it in a form that can be pasted
// straight into a model with no editing: the problem, the code it lives in,
// the change to make, and how to know it worked.
//
// It deliberately does NOT re-rank. The review already ranks: the judge reads
// every candidate side by side and orders them, lanes decide what blocks, and
// BuildFixPlan drops findings with no action to take. The top of that plan is
// the answer, so enhance and review can never disagree about what matters
// most — a second ranking here would be a second source of truth.

// FixPrompt is the one fix enhance reports, plus the context needed to say
// honestly how much it leaves on the table.
type FixPrompt struct {
	// Finding is the fix's subject, carried whole so callers (JSON, tests)
	// can read anything the render chooses not to show.
	Finding Finding
	// Step is the finding's entry in the review's own fix plan — its order,
	// lane, and whether doing it flips the verdict.
	Step FixStep
	// Remaining is how many other actionable findings the review produced,
	// and how many of those block. A prompt that shows one fix while eleven
	// others wait must say so, or it reads as "this is all that is wrong".
	Remaining         int
	RemainingBlocking int
}

// TopFix returns the review's highest-value actionable finding.
//
// The ordering is BuildFixPlan's: blocking lane first, then the report's own
// significance order within a lane, skipping findings with no recommendation.
// Reported false when the review found nothing to do — which is a real answer,
// not a failure, and the caller says so rather than inventing work.
func TopFix(report Report) (FixPrompt, bool) {
	plan := BuildFixPlan(report.Findings)
	if len(plan) == 0 {
		return FixPrompt{}, false
	}
	top := plan[0]
	byID := make(map[string]Finding, len(report.Findings))
	for _, f := range report.Findings {
		byID[f.ID] = f
	}
	finding, ok := byID[top.FindingID]
	if !ok {
		// A plan step always comes from a finding; if that ever stops being
		// true, say nothing rather than render a fix with no subject.
		return FixPrompt{}, false
	}
	blocking := 0
	for _, step := range plan[1:] {
		if step.FlipsVerdict {
			blocking++
		}
	}
	return FixPrompt{
		Finding:           finding,
		Step:              top,
		Remaining:         len(plan) - 1,
		RemainingBlocking: blocking,
	}, true
}

// RenderFixPrompt writes the fix as a model-ready prompt.
//
// Plain text, never ANSI: the output's job is to survive a copy-paste or a
// pipe into another tool, and escape codes would ride along into the model's
// context as noise. Markdown headings because every coding model reads them,
// and a fenced block so the code arrives as code.
func RenderFixPrompt(report Report, fix FixPrompt) string {
	f := fix.Finding
	var b strings.Builder

	title := strings.TrimSpace(f.Title)
	if title == "" {
		title = "Fix this finding"
	}
	fmt.Fprintf(&b, "# %s\n\n", title)

	// The locator line: everything needed to find the code, on one line.
	var facts []string
	if where := findingLocation(f); where != "" {
		facts = append(facts, where)
	}
	if rule := strings.TrimSpace(f.RuleID); rule != "" {
		facts = append(facts, rule)
	}
	if fix.Step.FlipsVerdict {
		facts = append(facts, "blocking — this is what fails the review")
	} else {
		facts = append(facts, "advisory")
	}
	if len(f.Corroboration) > 1 {
		facts = append(facts, fmt.Sprintf("%d reviewers agreed", len(f.Corroboration)))
	}
	fmt.Fprintf(&b, "%s\n\n", strings.Join(facts, " · "))

	if intent := strings.TrimSpace(report.Prompt); intent != "" {
		fmt.Fprintf(&b, "The change was meant to: %s\n\n", intent)
	}

	if summary := strings.TrimSpace(f.Summary); summary != "" {
		fmt.Fprintf(&b, "## What's wrong\n\n%s\n\n", summary)
	}
	if benefit := strings.TrimSpace(f.Benefit); benefit != "" {
		fmt.Fprintf(&b, "## Why it matters\n\n%s\n\n", benefit)
	}

	// The code, as code. DiffHunk when the finding is about the change,
	// CodeExcerpt when it is about a line the change did not touch.
	if code, label := fixPromptCode(f); code != "" {
		lang := fenceLanguage(f.File)
		if label == "diff" {
			lang = "diff"
		}
		fmt.Fprintf(&b, "## The code\n\n```%s\n%s\n```\n\n", lang, strings.TrimRight(code, "\n"))
	} else if len(f.Evidence) > 0 {
		// Deterministic findings carry no code — they point at a failing
		// tool. Their evidence IS the actionable part: without it the prompt
		// says "eslint failed" and the model has to go run eslint to learn
		// what this is about. With it, the diagnostics are already in context.
		b.WriteString("## The failing output\n\n")
		for _, e := range f.Evidence {
			value := strings.TrimSpace(e.Value)
			if value == "" {
				continue
			}
			if label := strings.TrimSpace(e.Label); label != "" {
				fmt.Fprintf(&b, "%s:\n\n", label)
			}
			fmt.Fprintf(&b, "```\n%s\n```\n\n", value)
		}
	}

	if rec := strings.TrimSpace(f.Recommendation); rec != "" {
		fmt.Fprintf(&b, "## Change to make\n\n%s\n\n", rec)
	}
	if example := strings.TrimSpace(f.Example); example != "" {
		fmt.Fprintf(&b, "Sketch of the change:\n\n```diff\n%s\n```\n\n", example)
	}

	// Acceptance. Named concretely so the model has a finish line rather than
	// a vibe: the rule that fired, and the command that re-checks it.
	b.WriteString("## Done when\n\n")
	where := findingLocation(f)
	switch rule := strings.TrimSpace(f.RuleID); {
	case rule != "" && where != "":
		fmt.Fprintf(&b, "- The code at %s no longer trips `%s`.\n", where, rule)
	case rule != "":
		fmt.Fprintf(&b, "- The change no longer trips `%s`.\n", rule)
	case where != "":
		fmt.Fprintf(&b, "- The problem described above is gone from %s.\n", where)
	default:
		// No rule and no line — a deterministic finding about the change as a
		// whole. "gone from the code above" would point at nothing.
		b.WriteString("- The problem described above is fixed.\n")
	}
	b.WriteString("- Existing behavior and tests still pass — this is a fix, not a rewrite.\n")
	if fix.Step.FlipsVerdict && fix.RemainingBlocking == 0 {
		b.WriteString("- `gx review` reports no blocking findings.\n")
	} else {
		b.WriteString("- `gx review` no longer reports this finding.\n")
	}
	b.WriteString("\n")

	// What this prompt is NOT showing. A single fix presented alone reads as
	// the whole story; the review's own counts are the correction.
	b.WriteString("---\n\n")
	switch {
	case fix.Remaining == 0:
		b.WriteString("This is the only fix this review found.\n")
	case fix.RemainingBlocking > 0:
		fmt.Fprintf(&b, "%s left after this one, %d still blocking. Run `gx review` for all of them.\n",
			moreFixes(fix.Remaining), fix.RemainingBlocking)
	default:
		fmt.Fprintf(&b, "%s left after this one, none blocking. Run `gx review` for all of them.\n",
			moreFixes(fix.Remaining))
	}
	return b.String()
}

// fixPromptCode picks the best available code for the finding, and says which
// kind it is so the fence can be labelled honestly — a diff hunk rendered as
// Go is not Go, and a model asked to edit it will reproduce the +/- markers.
func fixPromptCode(f Finding) (code, kind string) {
	if hunk := strings.TrimSpace(f.DiffHunk); hunk != "" {
		return hunk, "diff"
	}
	if excerpt := strings.TrimSpace(f.CodeExcerpt); excerpt != "" {
		return excerpt, "source"
	}
	return "", ""
}

// moreFixes renders "1 more fix" / "3 more fixes".
func moreFixes(n int) string {
	if n == 1 {
		return "1 more fix"
	}
	return fmt.Sprintf("%d more fixes", n)
}
