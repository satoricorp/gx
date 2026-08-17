package codereview

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseGatesAddedLinesNumbersFromHunks(t *testing.T) {
	diff := strings.Join([]string{
		"@@ -10,3 +12,4 @@ func Run() {",
		" kept",
		"-removed",
		"+added one",
		" kept",
		"+added two",
	}, "\n")
	added := parseGatesAddedLines(diff)
	if len(added) != 2 {
		t.Fatalf("added = %#v, want two lines", added)
	}
	if added[0].Number != 13 || added[0].Text != "added one" {
		t.Fatalf("first added = %#v, want line 13 'added one'", added[0])
	}
	if added[1].Number != 15 || added[1].Text != "added two" {
		t.Fatalf("second added = %#v, want line 15 'added two'", added[1])
	}
}

func TestGatesAddedLinesForContentFallback(t *testing.T) {
	diff := diffUnavailableContentHeader + "\nline one\nline two\n"
	added := gatesAddedLinesForDiff(diff)
	if len(added) != 3 { // trailing newline yields an empty final line
		t.Fatalf("added = %#v, want every content line treated as added", added)
	}
	if added[0].Number != 1 || added[0].Text != "line one" || added[1].Number != 2 {
		t.Fatalf("added = %#v, want 1-based file line numbers", added)
	}
}

func TestGatesDiffHunkForLine(t *testing.T) {
	diff := strings.Join([]string{
		"diff --git a/app.go b/app.go",
		"@@ -10,3 +12,4 @@ func Run() {",
		" kept",
		"-removed",
		"+added one",
		" kept",
		"+added two",
	}, "\n")
	hunk := gatesDiffHunkForLine(diff, 13)
	if hunk == "" {
		t.Fatalf("no hunk for line 13")
	}
	for _, want := range []string{"@@ -10,3 +12,4 @@", "-removed", "+added one", "+added two"} {
		if !strings.Contains(hunk, want) {
			t.Fatalf("hunk missing %q:\n%s", want, hunk)
		}
	}
	if got := gatesDiffHunkForLine(diff, 500); got != "" {
		t.Fatalf("line outside the diff produced a hunk:\n%s", got)
	}
	if got := gatesDiffHunkForLine(diffUnavailableContentHeader+"\ncontent\n", 1); got != "" {
		t.Fatalf("content fallback produced a hunk:\n%s", got)
	}
}

func TestGatesDiffHunkForLineTrimsLongHunks(t *testing.T) {
	lines := []string{"@@ -1,40 +1,40 @@"}
	for i := 1; i <= 40; i++ {
		lines = append(lines, fmt.Sprintf(" line %d", i))
	}
	hunk := gatesDiffHunkForLine(strings.Join(lines, "\n"), 20)
	if !strings.Contains(hunk, "line 20") {
		t.Fatalf("trimmed hunk lost the target line:\n%s", hunk)
	}
	if strings.Contains(hunk, "line 1\n") || strings.Contains(hunk, "line 40") {
		t.Fatalf("hunk was not trimmed around the target:\n%s", hunk)
	}
	if strings.Count(hunk, "…") != 2 {
		t.Fatalf("trimmed hunk should mark both cut edges:\n%s", hunk)
	}
}

func TestGatesNewDependencies(t *testing.T) {
	addedByFile := map[string][]gatesAddedLine{
		"go.mod": {
			{Number: 5, Text: "\tgithub.com/some/dep v1.2.3"},
			{Number: 6, Text: "\tgolang.org/x/text v0.14.0 // indirect"},
			{Number: 7, Text: ")"},
		},
		"web/package.json": {
			{Number: 12, Text: "    \"left-pad\": \"^1.3.0\","},
			{Number: 13, Text: "    \"version\": \"2.0.0\","},
		},
	}
	deps := gatesNewDependencies([]string{"go.mod", "web/package.json"}, addedByFile)
	joined := strings.Join(deps, "\n")
	if !strings.Contains(joined, "go.mod: github.com/some/dep") {
		t.Fatalf("deps = %#v, want the go.mod addition", deps)
	}
	if !strings.Contains(joined, "package.json: left-pad") {
		t.Fatalf("deps = %#v, want the package.json addition", deps)
	}
	if strings.Contains(joined, "version") {
		t.Fatalf("deps = %#v, package.json's own version field is not a dependency", deps)
	}
}

func TestGatesPerfApplicability(t *testing.T) {
	policy := ReviewPolicy{RiskPaths: []RiskPath{{Glob: "internal/api/**", Message: "latency-sensitive request path"}}}
	applies, triggers, files := gatesPerfApplicability(
		[]string{"internal/api/routes.go", "docs/README.md"}, nil, policy)
	if !applies {
		t.Fatalf("perf gate should apply to a risk-path match")
	}
	if len(triggers) == 0 || !strings.Contains(triggers[0], "REVIEW.md risk-path") {
		t.Fatalf("triggers = %#v, want the risk-path named", triggers)
	}
	if len(files) != 1 || files[0] != "internal/api/routes.go" {
		t.Fatalf("files = %#v, want the matched file", files)
	}

	applies, triggers, _ = gatesPerfApplicability([]string{"internal/cache/store.go"}, nil, ReviewPolicy{})
	if !applies || len(triggers) == 0 {
		t.Fatalf("perf gate should apply to a cache path, triggers = %#v", triggers)
	}

	applies, _, _ = gatesPerfApplicability([]string{"docs/README.md"}, nil, ReviewPolicy{})
	if applies {
		t.Fatalf("perf gate should not apply to a docs-only change")
	}

	applies, triggers, _ = gatesPerfApplicability([]string{"go.mod"}, []string{"go.mod"}, ReviewPolicy{})
	if !applies || !strings.Contains(strings.Join(triggers, " "), "dependency manifest") {
		t.Fatalf("perf gate should apply on a manifest change, triggers = %#v", triggers)
	}
}

func TestGatesA11yLineFindings(t *testing.T) {
	addedByFile := map[string][]gatesAddedLine{
		"web/App.tsx": {
			{Number: 4, Text: `    <img src="/logo.png">`},
			{Number: 9, Text: `    <div tabindex="3">ok</div>`},
			{Number: 14, Text: `    <div onClick={submit}>Save</div>`},
			{Number: 20, Text: `    <img src="/a.png" alt="A diagram">`},
			{Number: 25, Text: `    <div onClick={go} onKeyDown={go} role="button">Go</div>`},
		},
	}
	findings, hardFail := gatesA11yLineFindings([]string{"web/App.tsx"}, addedByFile)
	if !hardFail {
		t.Fatalf("missing alt and positive tabindex must be hard failures")
	}
	if len(findings) != 3 {
		t.Fatalf("findings = %d, want img-without-alt, positive tabindex, and clickable div: %#v", len(findings), findings)
	}
	strengths := map[string]int{}
	for _, finding := range findings {
		strengths[finding.Strength]++
	}
	if strengths["Strong"] != 2 || strengths["Worth exploring"] != 1 {
		t.Fatalf("strengths = %#v, want two Strong and one Worth exploring", strengths)
	}
}

func TestCollectGateSignalsCountsContentFallback(t *testing.T) {
	diffs := []DiffSnippet{{
		File: "internal/app/new.go",
		Diff: diffUnavailableContentHeader + "\npackage app\n\nfunc New() {}\n",
	}}
	signals := collectGateSignals([]string{"internal/app/new.go"}, diffs, ReviewPolicy{}, nil)
	if signals.DiffStats.AddedLines < 3 {
		t.Fatalf("AddedLines = %d, want the untracked file's lines counted", signals.DiffStats.AddedLines)
	}
	if signals.DiffStats.RemovedLines != 0 {
		t.Fatalf("RemovedLines = %d, want 0 for new content", signals.DiffStats.RemovedLines)
	}
}
