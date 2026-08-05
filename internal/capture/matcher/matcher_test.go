package matcher_test

import (
	"testing"
	"time"

	"github.com/satoricorp/lgtm/internal/capture"
	"github.com/satoricorp/lgtm/internal/capture/exclude"
	"github.com/satoricorp/lgtm/internal/capture/matcher"
)

func TestMatcher_Tier1Exact(t *testing.T) {
	events := []capture.SessionEvent{{
		Kind:     capture.KindEdit,
		FilePath: "internal/foo.go",
		NewText:  "func Bar() int {\n\treturn 42\n}\n",
		TS:       time.Unix(1_718_000_000, 0).UnixMilli(),
	}}
	hunks := []matcher.HunkRef{{
		Index:      0,
		FilePath:   "internal/foo.go",
		AddedLines: []string{"func Bar() int {", "\treturn 42", "}"},
		CommitTime: time.Unix(1_718_000_100, 0).UnixMilli(),
	}}

	result := matcher.Match(events, hunks, matcher.DefaultConfig())
	if len(result.Tier1HunkIndexes) != 1 {
		t.Fatalf("tier1 hunks = %d, want 1", len(result.Tier1HunkIndexes))
	}
	if len(result.Tier2HunkIndexes) != 0 {
		t.Fatalf("tier2 hunks = %d, want 0", len(result.Tier2HunkIndexes))
	}
}

func TestMatcher_Tier2Fuzzy(t *testing.T) {
	events := []capture.SessionEvent{{
		Kind:     capture.KindEdit,
		FilePath: "internal/foo.go",
		NewText:  "function configuration settings module handler helper utility",
		TS:       time.Unix(1_718_000_000, 0).UnixMilli(),
	}}
	hunks := []matcher.HunkRef{{
		Index:      0,
		FilePath:   "internal/foo.go",
		AddedLines: []string{"function configuration settings module handler helper extension"},
		CommitTime: time.Unix(1_718_000_100, 0).UnixMilli(),
	}}

	result := matcher.Match(events, hunks, matcher.DefaultConfig())
	if len(result.Tier1HunkIndexes) != 0 {
		t.Fatalf("tier1 hunks = %d, want 0", len(result.Tier1HunkIndexes))
	}
	if len(result.Tier2HunkIndexes) != 1 {
		t.Fatalf("tier2 hunks = %d, want 1", len(result.Tier2HunkIndexes))
	}
}

func TestMatcher_Tier3ExcludedFromAuthorship(t *testing.T) {
	commitTime := time.Unix(1_718_000_000, 0).UnixMilli()
	eventTime := commitTime + int64(5*time.Minute/time.Millisecond)
	events := []capture.SessionEvent{{
		Kind:     capture.KindEdit,
		FilePath: "internal/unrelated.go",
		NewText:  "totally different content here",
		TS:       eventTime,
	}}
	hunks := []matcher.HunkRef{{
		Index:      0,
		FilePath:   "internal/unrelated.go",
		AddedLines: []string{"package unrelated", "func Other() {}"},
		CommitTime: commitTime,
	}}

	result := matcher.Match(events, hunks, matcher.DefaultConfig())
	if result.Tier3Pairs == 0 {
		t.Fatal("expected tier3 diagnostic pairs")
	}
	if result.AuthorshipHunkCount() != 0 {
		t.Fatalf("authorship hunks = %d, want 0 (tier3 excluded)", result.AuthorshipHunkCount())
	}
}

func TestExcludeGenerated(t *testing.T) {
	m := exclude.NewMatcher(nil)
	cases := []struct {
		path string
		want bool
	}{
		{"go.sum", true},
		{"internal/foo_gen.go", true},
		{".lgtm/state.json", true},
		{"internal/foo.go", false},
	}
	for _, tc := range cases {
		if got := m.IsExcluded(tc.path); got != tc.want {
			t.Fatalf("IsExcluded(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}
