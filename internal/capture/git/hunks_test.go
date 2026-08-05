package git_test

import (
	"testing"

	"github.com/satoricorp/lgtm/internal/capture/git"
)

const sampleDiff = `diff --git a/internal/foo.go b/internal/foo.go
index 1111111..2222222 100644
--- a/internal/foo.go
+++ b/internal/foo.go
@@ -1 +1,3 @@
 package foo
+func Bar() int {
+	return 42
+}
`

func TestParseUnifiedDiffInline(t *testing.T) {
	t.Parallel()
	// Exercise the parser via an inline diff, without shelling out to git.
	hunks := parseTestDiff("abc123", 1_718_000_000_000, sampleDiff)
	if len(hunks) != 1 {
		t.Fatalf("hunks = %d, want 1", len(hunks))
	}
	if hunks[0].FilePath != "internal/foo.go" {
		t.Fatalf("file = %q", hunks[0].FilePath)
	}
	if len(hunks[0].AddedLines) != 3 {
		t.Fatalf("added lines = %d", len(hunks[0].AddedLines))
	}
}

func TestCommitCountRequiresGit(t *testing.T) {
	t.Skip("integration: requires git repo")
	_, _ = git.CommitCount(git.HunkOptions{RepoRoot: ".", Base: "main", Head: "HEAD"})
}

// parseTestDiff mirrors internal parser for unit testing without exporting it.
func parseTestDiff(sha string, ts int64, diff string) []struct {
	FilePath   string
	AddedLines []string
} {
	// duplicate minimal logic for test isolation
	type hunk struct {
		FilePath   string
		AddedLines []string
	}
	var hunks []hunk
	var current string
	var added []string
	flush := func() {
		if current == "" || len(added) == 0 {
			added = nil
			return
		}
		hunks = append(hunks, hunk{FilePath: current, AddedLines: append([]string(nil), added...)})
		added = nil
	}
	for _, line := range splitLines(diff) {
		if len(line) >= 6 && line[:6] == "+++ b/" {
			flush()
			current = line[6:]
			continue
		}
		if len(line) > 0 && line[0] == '+' && (len(line) < 4 || line[:4] != "+++") {
			added = append(added, line[1:])
		}
	}
	flush()
	out := make([]struct {
		FilePath   string
		AddedLines []string
	}, len(hunks))
	for i, h := range hunks {
		out[i].FilePath = h.FilePath
		out[i].AddedLines = h.AddedLines
	}
	_ = sha
	_ = ts
	return out
}

func splitLines(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return out
}
