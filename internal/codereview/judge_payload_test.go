package codereview

import "testing"

// The judge verifies findings against file content. When it receives none it
// cannot verify anything, answers in prose instead of its JSON verdict block,
// and every finding in that batch loses its verdict — the review then reports
// nothing while looking like a completed review. These tests pin the inputs
// that must produce file content.

func judgePayloadContext(changed []string) ReviewContext {
	return ReviewContext{
		Brief: ReviewBrief{
			RepoRoot: "/repo",
			Static:   StaticSnapshot{ChangedFiles: changed},
		},
		Facts: RepoFacts{Files: changed},
	}
}

// A finding whose file is carried ONLY in the structured field — no path in the
// prose. Models legitimately do this: the schema has a `file` field, so filling
// it and not repeating the path in the sentence is correct behavior.
func TestNamedFilesUsesStructuredFileField(t *testing.T) {
	ctx := judgePayloadContext([]string{"src/orders.ts"})
	finding := Finding{
		Title:          "Ownership check missing",
		Summary:        "Any caller can read another user's order.",
		Recommendation: "Compare the caller against the record owner before returning it.",
		File:           "src/orders.ts",
		Line:           14,
	}

	files := namedFilesForFinding(ctx, finding)
	if len(files) == 0 {
		t.Fatal("no named files: the judge will be sent null file content and cannot verify this finding")
	}
	if files[0] != "src/orders.ts" {
		t.Fatalf("namedFilesForFinding() = %v, want src/orders.ts", files)
	}
}

// Anchors are the other structured carrier of location.
func TestNamedFilesUsesAnchors(t *testing.T) {
	ctx := judgePayloadContext([]string{"server/src/auth.ts"})
	finding := Finding{
		Title:   "Timing-unsafe token comparison",
		Summary: "The comparison short-circuits on the first differing byte.",
		Anchors: []FindingAnchor{{File: "server/src/auth.ts", Line: 4}},
	}

	if files := namedFilesForFinding(ctx, finding); len(files) == 0 {
		t.Fatal("no named files from anchors: the judge cannot verify this finding")
	}
}

// The prose path must still work, and must not be dropped just because the
// repository scan produced no Facts.Files.
func TestNamedFilesFromProseSurvivesEmptyFacts(t *testing.T) {
	ctx := ReviewContext{
		Brief: ReviewBrief{
			RepoRoot: "/repo",
			Static:   StaticSnapshot{ChangedFiles: []string{"src/orders.ts"}},
		},
	}
	finding := Finding{
		Title:          "Race in incrementTotal",
		Recommendation: "In `src/orders.ts` line 31, serialize the read-modify-write.",
	}

	if files := namedFilesForFinding(ctx, finding); len(files) == 0 {
		t.Fatal("prose path was not recognised")
	}
}
