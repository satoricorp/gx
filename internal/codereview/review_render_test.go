package codereview

import (
	"strings"
	"testing"
)

// A fixture shaped like the checkout-retry mockup: one blocking finding both
// graders agreed on, one advisory finding, one demoted finding, and a story
// with a passed rule.
func renderFixtureReport() Report {
	twoLegs := []string{"Bedrock A", "Bedrock B"}
	findings := AssignLanes([]Finding{
		{
			ID: "bedrock-a.ai.review.1", RuleID: "gx:recommended/no-secrets-in-logs",
			Title:          "Secrets must not reach logs, traces, or error strings",
			Summary:        "req embeds PaymentToken, and %+v prints it — the raw token reaches your log sink once per retry.",
			Recommendation: "log req.ID and req.Amount instead of req.",
			Strength:       "Strong", Corroboration: twoLegs, JudgeVerdict: "confirmed", JudgeConfidence: 0.91,
			File: "internal/checkout/session.go", Line: 88,
			CodeExcerpt:      "    for attempt := 0; attempt < maxRetries; attempt++ {\n        log.Printf(\"charge attempt %d failed: %+v\", attempt, req)",
			CodeExcerptStart: 87,
			Example:          "-        log.Printf(\"charge attempt %d failed: %+v\", attempt, req)\n+        log.Printf(\"charge attempt %d failed: id=%s amount=%d\", attempt, req.ID, req.Amount)",
		},
		{
			ID: "bedrock-a.ai.review.2", RuleID: "gx:recommended/reuse-before-rewrite",
			Title:          "Don't add what the codebase already has",
			Summary:        "internal/backoff.Exponential already does this, with jitter.",
			Recommendation: "backoff.Exponential(attempt, backoff.WithCap(30*time.Second))",
			Strength:       "Worth exploring", Corroboration: twoLegs, JudgeVerdict: "confirmed", JudgeConfidence: 0.84,
			File: "internal/checkout/retry.go", Line: 24,
			DiffHunk: "+func backoffDelay(attempt int) time.Duration {\n+    return min(d, 30*time.Second)\n+}",
		},
		{
			ID: "bedrock-b.ai.review.3", RuleID: "gx:recommended/comments-match-code",
			Title:          "Comments must describe what the code now does",
			Summary:        "the comment predates this change — Charge retries now.",
			Recommendation: "update the comment, or delete it.",
			Strength:       "Strong", Corroboration: []string{"Bedrock B"}, JudgeVerdict: "confirmed", JudgeConfidence: 0.77,
			File: "internal/checkout/session.go", Line: 71,
		},
	})
	return Report{
		Reviewed:        true,
		ReviewMode:      ReviewModeRange,
		ReviewRange:     "origin/main...feature/checkout-retry",
		ReviewBase:      "origin/main",
		ChangedFiles:    []string{"a.go", "b.go", "c.go"},
		Prompt:          "make checkout survive gateway blips",
		Reviewer:        "heuristic+ai",
		ReviewModels:    []string{"claude-a", "claude-b"},
		ReviewTransport: "gx Cloud",
		Findings:        findings,
		Story: []StoryItem{
			{
				Headline: "Checkout retries on 5xx where it used to fail fast", Category: StoryBehaviorDelta, Materiality: "high",
				File: "internal/checkout/session.go", Line: 88,
				DiffHunk:    "-if err != nil {\n-    return nil, err\n+if err != nil && isRetryable(err) {\n+    continue",
				Consequence: "a failed charge that returned in ~200ms now takes up to 4 attempts over ~15s.",
				Asked:       "make checkout survive gateway blips",
				Chose:       "in-process retry; rejected a job queue as beyond stated scope",
				Watch:       []string{"p95 checkout latency"},
			},
			{
				Headline: "The idempotency key is what makes the retry safe", Category: StorySemanticShift, Materiality: "high",
				File: "internal/checkout/session.go", Consequence: "the key spans the loop, so the gateway dedupes retries.",
				PassedRule: "gx:recommended/retries-are-idempotent",
			},
		},
	}
}

func TestRenderReviewTextPlainHasEverySectionInOrder(t *testing.T) {
	out := RenderReviewText(renderFixtureReport()) // Color=false: plain
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("plain render must carry no ANSI:\n%s", out)
	}
	order := []string{
		"gx review — origin/main...feature/checkout-retry",
		`intent: "make checkout survive gateway blips"`,
		"scope",
		"reviewers",
		"findings",
		"BLOCKING",
		"no-secrets-in-logs  gx:recommended  ·  internal/checkout/session.go:88",
		"Secrets must not reach logs",
		"Why  req embeds PaymentToken",
		"Fix  log req.ID and req.Amount",
		"-        log.Printf(\"charge attempt %d failed: %+v\", attempt, req)",
		"+        log.Printf(\"charge attempt %d failed: id=%s amount=%d\", attempt, req.ID, req.Amount)",
		"both graders agreed · judge confirmed",
		"ADVISORY",
		"Don't add what the codebase already has",
		"Comments must describe what the code now does",
		"one grader flagged · judge confirmed · demoted from blocking",
		"to silence a rule where it doesn't apply",
		"FIX PLAN",
		// The action leads; the rule and location are the muted line under it.
		"1. log req.ID and req.Amount instead of req.",
		"no-secrets-in-logs · internal/checkout/session.go:88",
		"WORTH KNOWING",
		"1  Checkout retries on 5xx where it used to fail fast   materiality: high",
		"What changes for you: a failed charge",
		`asked  "make checkout survive gateway blips"`,
		"chose  in-process retry",
		"watch  p95 checkout latency",
		"retries-are-idempotent passed on this",
		"Verdict: NO-SHIP — 1 blocking finding(s) (no-secrets-in-logs)",
		"Next: work the fix plan, then rerun `gx review`.",
	}
	pos := 0
	for _, want := range order {
		idx := strings.Index(out[pos:], want)
		if idx < 0 {
			t.Fatalf("missing or out of order: %q\n--- render ---\n%s", want, out)
		}
		pos += idx + len(want)
	}
	// The seam is the last two lines, exactly.
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 2 || !strings.HasPrefix(lines[len(lines)-2], "Verdict: ") || !strings.HasPrefix(lines[len(lines)-1], "Next: ") {
		t.Fatalf("Verdict/Next must be the last two lines, got:\n%s", strings.Join(lines[len(lines)-3:], "\n"))
	}
}

func TestRenderReviewTextFixPlanBlockingFirst(t *testing.T) {
	out := RenderReviewText(renderFixtureReport())
	plan := out[strings.Index(out, "FIX PLAN"):strings.Index(out, "WORTH KNOWING")]
	// Step one is an action, so the rule it came from is the reference line
	// beneath it — before step two either way.
	first := strings.Index(plan, "no-secrets-in-logs")
	second := strings.Index(plan, "2. ")
	if first < 0 || second < 0 || first > second {
		t.Fatalf("fix plan must list the blocking finding first:\n%s", plan)
	}
}

// Every step carries its action, not just the rule that raised it: the plan is
// the section meant to be worked top to bottom.
func TestRenderReviewTextFixPlanCarriesTheAction(t *testing.T) {
	out := RenderReviewText(renderFixtureReport())
	plan := out[strings.Index(out, "FIX PLAN"):strings.Index(out, "WORTH KNOWING")]
	for _, want := range []string{
		"1. log req.ID and req.Amount instead of req.",
		"2. backoff.Exponential(attempt, backoff.WithCap(30*time.Second))",
		"3. update the comment, or delete it.",
	} {
		if !strings.Contains(plan, want) {
			t.Fatalf("fix plan must carry the action %q:\n%s", want, plan)
		}
	}
}

// The tools row is the only place a passing checker is reported — a failure
// becomes the tools.static-failure finding, but a pass has nowhere else to go.
// A skipped run is neither: it must not read as a pass.
func TestRenderReviewTextLedgerReportsToolOutcomes(t *testing.T) {
	r := renderFixtureReport()
	r.Tools = []StaticToolResult{
		{Name: "go test", ExitCode: 0},
		{Name: "eslint", ExitCode: 1},
		{Name: "govulncheck", Skipped: true, Reason: "not installed"},
	}
	out := RenderReviewText(r)
	if !strings.Contains(out, "tools      go test ✓  eslint ✕  govulncheck —") {
		t.Fatalf("ledger must report each tool's outcome:\n%s", out)
	}
	// No tools detected is no row, not an empty one.
	r.Tools = nil
	if strings.Contains(RenderReviewText(r), "tools ") {
		t.Fatal("a repo with no detected checkers must get no tools row")
	}
}

func TestRenderReviewTextVerdictShapes(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Report)
		verdict string
	}{
		{"nothing to review", func(r *Report) { r.Reviewed = false; r.ReviewTarget = "the working tree" },
			"Verdict: NOTHING-TO-REVIEW — no change found (looked at the working tree)"},
		{"clean", func(r *Report) { r.Findings = nil; r.Story = nil }, "Verdict: SHIP — no findings"},
		{"advisory only", func(r *Report) { r.Findings = r.Findings[1:]; r.Story = nil }, "Verdict: SHIP — no blocking finding (2 advisory)"},
		{"degraded", func(r *Report) {
			r.Findings = nil
			r.Story = nil
			r.DegradedReasons = []string{"reviewer B did not answer"}
		}, "Verdict: DEGRADED — no blocking finding, but reviewer B did not answer (0 advisory)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := renderFixtureReport()
			tc.mutate(&r)
			got, _ := ReviewVerdictLines(r)
			if got != tc.verdict {
				t.Fatalf("verdict = %q\nwant      %q", got, tc.verdict)
			}
			out := RenderReviewText(r)
			if !strings.Contains(out, got) {
				t.Fatalf("render missing its own verdict line:\n%s", out)
			}
		})
	}
}

func TestRenderReviewTextColorIsGatedOnReportColor(t *testing.T) {
	t.Setenv("FORCE_COLOR", "1")
	t.Setenv("NO_COLOR", "")
	r := renderFixtureReport()
	r.Color = true
	colored := RenderReviewText(r)
	if !strings.Contains(colored, "\x1b[") {
		t.Fatalf("Color=true under FORCE_COLOR should paint:\n%s", colored)
	}
	r.Color = false
	if plain := RenderReviewText(r); strings.Contains(plain, "\x1b[") {
		t.Fatalf("Color=false must never paint even under FORCE_COLOR")
	}
}

func TestSuppressLineIsPasteReady(t *testing.T) {
	f := Finding{RuleID: "REVIEW.md/no-reinvented-utils", File: "internal/x/y.go"}
	got := SuppressLine(f)
	want := `- "no-reinvented-utils" doesn't apply in internal/x/y.go — <reason>.`
	if got != want {
		t.Fatalf("SuppressLine = %q, want %q", got, want)
	}
}

func TestWrapTextNeverSplitsWords(t *testing.T) {
	lines := wrapText("aaa bbb ccc ddd eee", 7)
	want := []string{"aaa bbb", "ccc ddd", "eee"}
	if strings.Join(lines, "|") != strings.Join(want, "|") {
		t.Fatalf("wrapText = %v, want %v", lines, want)
	}
	if got := wrapText("", 10); got != nil {
		t.Fatalf("wrapText(empty) = %v, want nil", got)
	}
}

func TestNormalizeExampleDiffKeepsSmallDiffsOnly(t *testing.T) {
	cases := map[string]string{
		"-old\n+new":                    "-old\n+new",
		"```diff\n-old\n+new\n```":      "-old\n+new",
		"  -old\n  +new  ":              "-old\n  +new",
		"just prose explaining the fix": "",
		"":                              "",
		strings.Repeat("+line\n", 20):   "",
	}
	for in, want := range cases {
		if got := normalizeExampleDiff(in); got != want {
			t.Errorf("normalizeExampleDiff(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFirstSentencesClampsProse(t *testing.T) {
	long := "The value is nil here. It is dereferenced on line 9. Callers never guard it. A fourth sentence."
	if got := firstSentences(long, 2); got != "The value is nil here. It is dereferenced on line 9. …" {
		t.Fatalf("firstSentences(2) = %q", got)
	}
	if got := firstSentences(long, 10); got != long {
		t.Fatalf("under the cap must be untouched, got %q", got)
	}
	// Abbreviations and paths do not end a sentence.
	tricky := "See e.g. internal/cli.go for the pattern. Then fix it."
	if got := firstSentences(tricky, 1); got != "See e.g. internal/cli.go for the pattern. …" {
		t.Fatalf("abbreviation split wrongly: %q", got)
	}
	// A backtick or quote can start the next sentence.
	code := "Guard the input. `req.ID` is the key. Done."
	if got := firstSentences(code, 1); got != "Guard the input. …" {
		t.Fatalf("backtick sentence start not recognized: %q", got)
	}
}

func TestRenderCodeLinesClipAndDropHunkHeaders(t *testing.T) {
	long := strings.Repeat("x", 300)
	f := Finding{File: "a.go", DiffHunk: "@@ -1,2 +1,2 @@\n-\tshort\n+" + long}
	lines := renderCodeLines(false, f, "      ")
	if len(lines) != 2 {
		t.Fatalf("expected 2 body lines (header dropped), got %d: %q", len(lines), lines)
	}
	if !strings.HasPrefix(lines[0], "      -    short") {
		t.Fatalf("tab should expand to spaces after the sign: %q", lines[0])
	}
	if !strings.HasSuffix(lines[1], "…") || len([]rune(lines[1])) > renderCodeWidth+2 {
		t.Fatalf("long code line should be clipped with an ellipsis: len=%d", len([]rune(lines[1])))
	}
	ex := Finding{File: "a.go", CodeExcerpt: "\tone\n" + long, CodeExcerptStart: 4, Line: 5}
	el := renderExcerptLines(false, ex, "      ")
	if !strings.Contains(el[0], "|     one") || !strings.HasSuffix(el[1], "…") {
		t.Fatalf("excerpt should expand tabs and clip: %q", el)
	}
}

func TestRenderHunkWindowMarksTheDiscussedLine(t *testing.T) {
	// @@ says the post-change side starts at 86. Lines: 86 ctx, 87 ctx, 88 add.
	hunk := "@@ -80,3 +86,3 @@\n     for attempt := 0; attempt < 3; attempt++ {\n         resp, err := charge()\n+        log.Printf(\"%+v\", req)"
	lines := renderHunkWindow(false, "a.go", hunk, 88, "      ")
	if len(lines) != 3 {
		t.Fatalf("expected 3 body lines, got %d: %q", len(lines), lines)
	}
	if strings.Contains(lines[0], "→") || strings.Contains(lines[1], "→") {
		t.Fatalf("only the target line is marked: %q", lines)
	}
	if !strings.HasPrefix(lines[2], "    → +") {
		t.Fatalf("target line should carry the arrow in the indent and keep its sign: %q", lines[2])
	}
	// Column alignment: unmarked lines are indent(6) + " " + body, so the
	// sign column is 6; the marked line is indent[:4] + "→ " + "+" + body,
	// which puts its "+" at rune 6 too — the code columns line up.
	// Count runes, not bytes: → is one column and three bytes.
	if idx := strings.IndexRune(string([]rune(lines[2])), '+'); len([]rune(lines[2][:strings.IndexRune(lines[2], '+')])) != 6 {
		t.Fatalf("marked line's sign should sit at column 6, got %d: %q", idx, lines[2])
	}
	// No header, no line numbers → nothing marked, nothing broken.
	if got := renderHunkWindow(false, "a.go", "-a\n+b", 1, "  "); strings.Contains(strings.Join(got, ""), "→") {
		t.Fatalf("headerless hunk must not mark: %q", got)
	}
}
