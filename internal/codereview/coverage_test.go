package codereview

import (
	"strings"
	"testing"
)

// "No material issues found in the repository" is a claim about the
// repository. A review that read twelve of its 273 source files has not earned
// that sentence, and the reader has no way to tell it apart from one that read
// everything.
func TestNoFindingsWordingDoesNotOverclaimCoverage(t *testing.T) {
	partial := Coverage{Planned: true, Subject: "the repository", Units: "files", Total: 273, Read: 12}
	message := noFindingsMessage(ReviewModeRepo, partial)
	if strings.Contains(message, "No material issues found in the repository.") {
		t.Fatalf("partial review still claims the repository was reviewed: %q", message)
	}
	for _, want := range []string{"not a clean bill of health", "Read 12 of 273 files"} {
		if !strings.Contains(message, want) {
			t.Fatalf("message %q missing %q", message, want)
		}
	}

	complete := Coverage{Planned: true, Subject: "the repository", Units: "files", Total: 273, Read: 273}
	message = noFindingsMessage(ReviewModeRepo, complete)
	if !strings.HasPrefix(message, "No material issues found in the repository.") {
		t.Fatalf("a complete review lost its plain wording: %q", message)
	}
	if !strings.Contains(message, "Read all 273 files") {
		t.Fatalf("a complete review does not state its coverage: %q", message)
	}
}

// With no coverage computed at all — the AI never ran — the wording is the
// wording it has always been. A zero value must not read as "read nothing".
func TestNoFindingsWordingIsUnchangedWithoutCoverage(t *testing.T) {
	if got := noFindingsMessage(ReviewModeWorkingTree, Coverage{}); got != "No material issues found in this change." {
		t.Fatalf("message = %q", got)
	}
	if got := noFindingsMessage(ReviewModeRepo, Coverage{}); got != "No material issues found in the repository." {
		t.Fatalf("message = %q", got)
	}
}

// Findings do not excuse partial coverage: two problems found in a third of a
// change is not two problems found in the change.
func TestPartialCoverageIsWarnedAboveTheFindings(t *testing.T) {
	rendered := RenderMarkdown(Report{
		ReviewMode:   ReviewModeRange,
		ReviewTarget: "`main...HEAD`",
		Reviewed:     true,
		Coverage:     Coverage{Planned: true, Subject: "the change", Units: "changed files", Total: 181, Read: 60, Shards: 3, ShardsFailed: 1},
		Findings: []Finding{{
			Title: "A real problem", Summary: "s", Recommendation: "r", Strength: "Worth exploring",
		}},
	})
	if !strings.Contains(rendered, "partial coverage") {
		t.Fatalf("report does not warn about partial coverage:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Read 60 of 181 changed files") {
		t.Fatalf("report does not state what was read:\n%s", rendered)
	}
	if !strings.Contains(rendered, "parallel reviews failed") {
		t.Fatalf("report hides the failed shard:\n%s", rendered)
	}
}

func TestCompleteCoverageDoesNotWarn(t *testing.T) {
	rendered := RenderMarkdown(Report{
		ReviewMode:   ReviewModeRange,
		ReviewTarget: "`main...HEAD`",
		Reviewed:     true,
		Coverage:     Coverage{Planned: true, Subject: "the change", Units: "changed files", Total: 4, Read: 4, Shards: 1},
	})
	if strings.Contains(rendered, "partial coverage") {
		t.Fatalf("a complete review warned about coverage:\n%s", rendered)
	}
	if !strings.Contains(rendered, "Read all 4 changed files") {
		t.Fatalf("report does not state its coverage:\n%s", rendered)
	}
}
