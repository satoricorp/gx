package codereview

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The judge is shown file content so it can check a claim about a location.
// These tests pin the two properties that decide whether it gets the right
// bytes and how many of them: the excerpt is built around the line the finding
// points at, and each file is sent once per batch however many candidates name
// it.

// judgeExcerptRepo writes files into a temp repo and returns a ReviewContext
// rooted at it.
func judgeExcerptRepo(t *testing.T, files map[string]string) ReviewContext {
	t.Helper()
	root := t.TempDir()
	var changed []string
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
		changed = append(changed, name)
	}
	return ReviewContext{
		Brief: ReviewBrief{
			RepoRoot: root,
			Static:   StaticSnapshot{ChangedFiles: changed},
		},
		Facts: RepoFacts{Files: changed},
	}
}

// numberedLines builds a file whose every line names its own number, so a test
// can tell which part of it an excerpt came from.
func numberedLines(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "content of line %d\n", i)
	}
	return b.String()
}

// A finding about line 900 must be shown line 900. The head-truncated excerpt
// this replaced sent the first 24KB of the file regardless of where the finding
// pointed, so a claim about the end of a large file was verified against a
// region that could not contain it.
func TestJudgeExcerptWindowsAroundTheAnchorLine(t *testing.T) {
	ctx := judgeExcerptRepo(t, map[string]string{"src/big.go": numberedLines(1000)})
	finding := Finding{
		ID:      "f1",
		Title:   "Bug near the end of the file",
		Summary: "The guard is missing.",
		File:    "src/big.go",
		Line:    900,
	}

	req := buildJudgeRequest(ctx, []Finding{finding})

	if len(req.Files) != 1 {
		t.Fatalf("Files = %d, want 1", len(req.Files))
	}
	text := req.Files[0].Text
	if !strings.Contains(text, "content of line 900") {
		t.Fatal("excerpt does not contain the anchored line; the judge cannot check the claim")
	}
	if strings.Contains(text, "content of line 1\n") {
		t.Error("excerpt still starts at the top of the file rather than windowing on the anchor")
	}
	if !strings.Contains(text, "omitted") {
		t.Error("excerpt drops a stretch of the file without saying so; the judge cannot tell absent-from-excerpt from absent-from-file")
	}
}

// The line numbers are what make a claim about a location checkable.
func TestJudgeExcerptCarriesRealLineNumbers(t *testing.T) {
	ctx := judgeExcerptRepo(t, map[string]string{"src/big.go": numberedLines(1000)})
	finding := Finding{ID: "f1", File: "src/big.go", Line: 900}

	req := buildJudgeRequest(ctx, []Finding{finding})

	if !strings.Contains(req.Files[0].Text, "900| content of line 900") {
		t.Fatalf("excerpt is not line-numbered against the real file:\n%s", firstLines(req.Files[0].Text, 3))
	}
}

// Several findings in one file is the ordinary case, and it used to send that
// file's bytes once per finding.
func TestJudgeRequestSendsEachFileOnce(t *testing.T) {
	ctx := judgeExcerptRepo(t, map[string]string{"src/orders.ts": numberedLines(200)})
	findings := []Finding{
		{ID: "f1", Title: "One", File: "src/orders.ts", Line: 10},
		{ID: "f2", Title: "Two", File: "src/orders.ts", Line: 20},
		{ID: "f3", Title: "Three", File: "src/orders.ts", Line: 30},
	}

	req := buildJudgeRequest(ctx, findings)

	if len(req.Files) != 1 {
		t.Fatalf("Files = %d for 3 findings in 1 file, want 1", len(req.Files))
	}
	if len(req.Candidates) != 3 {
		t.Fatalf("Candidates = %d, want 3", len(req.Candidates))
	}
	for _, candidate := range req.Candidates {
		if len(candidate.NamedFiles) == 0 {
			t.Fatalf("candidate %s names no file, so it cannot reach the shared content", candidate.ID)
		}
	}
}

// Findings a few lines apart must not produce two copies of the same code.
func TestJudgeLineRangesMergeOverlappingWindows(t *testing.T) {
	ranges := judgeLineRanges([]int{100, 110}, 1000)

	if len(ranges) != 1 {
		t.Fatalf("ranges = %v, want the two overlapping windows merged into one", ranges)
	}
	if ranges[0].start != 20 || ranges[0].end != 190 {
		t.Fatalf("merged range = %v, want {20 190}", ranges[0])
	}
}

// Windows far apart stay separate rather than dragging in everything between.
func TestJudgeLineRangesKeepDistantWindowsApart(t *testing.T) {
	ranges := judgeLineRanges([]int{100, 900}, 1000)

	if len(ranges) != 2 {
		t.Fatalf("ranges = %v, want 2 separate windows", ranges)
	}
}

// A finding that names no line still has to be verifiable, so the whole file
// (budget permitting) is the fallback rather than nothing.
func TestJudgeExcerptFallsBackToWholeFileWithoutAnchors(t *testing.T) {
	text := judgeFileExcerpt(numberedLines(20), nil, maxJudgeFileBytes)

	if !strings.Contains(text, "1| content of line 1") {
		t.Error("unanchored excerpt should start at line 1")
	}
	if !strings.Contains(text, "20| content of line 20") {
		t.Error("unanchored excerpt should reach the end of a small file")
	}
	if strings.Contains(text, "omitted") {
		t.Error("a whole small file was sent, so nothing should be reported omitted")
	}
}

// A line number past the end of the file is a claim the judge should be able to
// reject, which means seeing the file rather than an empty excerpt.
func TestJudgeExcerptIgnoresOutOfRangeAnchors(t *testing.T) {
	text := judgeFileExcerpt(numberedLines(20), []int{5000}, maxJudgeFileBytes)

	if !strings.Contains(text, "1| content of line 1") {
		t.Fatalf("out-of-range anchor produced an excerpt that omits the file:\n%s", text)
	}
}

// The request crosses into a cloud call, and repo_root was an absolute path on
// the reviewer's machine that the model could not act on.
func TestJudgeRequestCarriesNoLocalAbsolutePath(t *testing.T) {
	ctx := judgeExcerptRepo(t, map[string]string{"src/orders.ts": numberedLines(20)})
	req := buildJudgeRequest(ctx, []Finding{{ID: "f1", File: "src/orders.ts", Line: 5}})

	encoded, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(encoded), "repo_root") {
		t.Error("request still carries repo_root")
	}
	if strings.Contains(string(encoded), ctx.Brief.RepoRoot) {
		t.Errorf("request leaks the local repository path %q", ctx.Brief.RepoRoot)
	}
}

// The per-file budget must bound one file, and the batch budget the whole
// request, so a review over many large files cannot quietly build a giant call.
func TestJudgeExcerptRespectsItsBudget(t *testing.T) {
	text := judgeFileExcerpt(numberedLines(10000), nil, 2048)

	if len(text) > 2048 {
		t.Fatalf("excerpt is %d bytes, want it bounded by the 2048 budget", len(text))
	}
	if !strings.Contains(text, "truncated") {
		t.Error("excerpt was cut at the budget without saying so")
	}
}

// Omission markers are content, and a file with many scattered anchors emits
// one per gap. Counting only the source lines let their combined size carry the
// excerpt past its cap.
func TestJudgeExcerptCountsOmissionMarkersAgainstTheBudget(t *testing.T) {
	// Anchors far enough apart that every window is its own span, so each gap
	// between them produces a marker.
	var anchors []int
	for line := 200; line <= 9800; line += 400 {
		anchors = append(anchors, line)
	}

	for _, budget := range []int{512, 2048, 8192} {
		text := judgeFileExcerpt(numberedLines(10000), anchors, budget)
		if len(text) > budget {
			t.Errorf("budget %d: excerpt is %d bytes, over cap by %d",
				budget, len(text), len(text)-budget)
		}
	}
}

func firstLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	return strings.Join(lines, "\n")
}
