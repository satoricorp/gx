package review

import (
	"fmt"
	"strings"
)

// Grade is a section outcome.
type Grade string

const (
	GradePass Grade = "pass"
	GradeWarn Grade = "warn"
	GradeFail Grade = "fail"
	GradeSkip Grade = "skip"
)

const (
	SectionIntentMatch  = "intent_match"
	SectionStructure    = "structure"
	SectionRules        = "rules"
	SectionBlastRadius  = "blast_radius"
	SectionCollisions   = "collisions"
	SectionDoneNess     = "done_ness"
	SectionBugPass      = "bug_pass"
)

// Section is one graded review dimension.
type Section struct {
	Name     string   `json:"name"`
	Grade    Grade    `json:"grade"`
	Summary  string   `json:"summary"`
	Evidence []string `json:"evidence,omitempty"`
}

// GradedReport is the local gx review artifact.
type GradedReport struct {
	Repo     string    `json:"repo"`
	RefRange string    `json:"ref_range"`
	Sections []Section `json:"sections"`
	Offline  bool      `json:"offline"`
}

// RenderMarkdown formats the graded report for stdout.
func RenderMarkdown(report GradedReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# GX review\n\n")
	fmt.Fprintf(&b, "Repo: `%s`\n", report.Repo)
	fmt.Fprintf(&b, "Range: `%s`\n", report.RefRange)
	if report.Offline {
		fmt.Fprintf(&b, "\n_Server context unavailable — local grading only._\n")
	}
	fmt.Fprintf(&b, "\n## Dimensions\n\n")
	for _, section := range report.Sections {
		fmt.Fprintf(&b, "### %s — %s\n\n", humanSectionName(section.Name), section.Grade)
		fmt.Fprintf(&b, "%s\n", section.Summary)
		if len(section.Evidence) > 0 {
			fmt.Fprintf(&b, "\n")
			for _, item := range section.Evidence {
				fmt.Fprintf(&b, "- %s\n", item)
			}
		}
		fmt.Fprintf(&b, "\n")
	}
	return b.String()
}

func humanSectionName(name string) string {
	switch name {
	case SectionIntentMatch:
		return "Intent match"
	case SectionStructure:
		return "Structure"
	case SectionRules:
		return "Rules"
	case SectionBlastRadius:
		return "Blast radius"
	case SectionCollisions:
		return "Collisions"
	case SectionDoneNess:
		return "Done-ness"
	case SectionBugPass:
		return "Bug pass"
	default:
		return strings.ReplaceAll(name, "_", " ")
	}
}
