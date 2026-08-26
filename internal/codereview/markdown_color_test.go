package codereview

import (
	"strings"
	"testing"
)

// A GitHub comment, a stored summary, and a piped --md render have one thing
// in common: none of them is a terminal. Markdown must never carry an escape
// sequence, whatever the report says about colour.
func TestRenderMarkdownNeverEmitsEscapeSequences(t *testing.T) {
	// Force colour on. Without this the test passes for the wrong reason:
	// termstyle.Enabled() is false when stdout is not a terminal, which is
	// every `go test` run, and that is precisely why this shipped.
	t.Setenv("FORCE_COLOR", "1")
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")

	report := Report{
		Color:           true,
		RepoRoot:        "/repo",
		Reviewer:        "heuristic+ai",
		DegradedReasons: []string{"one reviewer did not run"},
		Findings: []Finding{{
			ID: "f1", Title: "A finding", Summary: "Something is wrong",
			Recommendation: "Fix it", File: "a.go", Line: 3, Strength: "Strong",
		}},
	}

	out := RenderMarkdown(report)
	if strings.ContainsRune(out, '\x1b') {
		t.Fatalf("RenderMarkdown emitted an ANSI escape:\n%q", out)
	}
	// The content must still be there — stripping colour is not stripping text.
	if !strings.Contains(out, "A finding") {
		t.Fatalf("RenderMarkdown dropped the finding:\n%s", out)
	}
}
