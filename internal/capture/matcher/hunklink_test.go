package matcher_test

import (
	"testing"

	"github.com/satoricorp/lgtm/internal/capture"
	"github.com/satoricorp/lgtm/internal/capture/matcher"
)

func TestBuildHunkLinks_AgentMatch(t *testing.T) {
	events := []capture.SessionEvent{{
		SessionID: "sess-1",
		Tool:      capture.ToolCursor,
		Model:     "gpt-4",
		Kind:      capture.KindEdit,
		FilePath:  "internal/foo.go",
		NewText:   "func Bar() {}",
	}}
	hunks := []capture.CommitHunk{{
		CommitSHA:  "abc123",
		FilePath:   "internal/foo.go",
		AddedLines: []string{"func Bar() {}"},
	}}
	refs := []matcher.HunkRef{{Index: 0, FilePath: "internal/foo.go", AddedLines: hunks[0].AddedLines}}
	result := matcher.Match(events, refs, matcher.DefaultConfig())
	links := matcher.BuildHunkLinks(hunks, events, result)
	if len(links) != 1 {
		t.Fatalf("links = %d, want 1", len(links))
	}
	if links[0].Authorship != matcher.AuthorshipAgent {
		t.Fatalf("authorship = %q, want agent", links[0].Authorship)
	}
	if links[0].SessionID != "sess-1" {
		t.Fatalf("sessionID = %q", links[0].SessionID)
	}
	if links[0].Tier != matcher.TierExact {
		t.Fatalf("tier = %d, want 1", links[0].Tier)
	}
}

func TestBuildHunkLinks_UnmatchedHuman(t *testing.T) {
	hunks := []capture.CommitHunk{{
		CommitSHA:  "abc123",
		FilePath:   "internal/foo.go",
		AddedLines: []string{"manual edit"},
	}}
	refs := []matcher.HunkRef{{Index: 0, FilePath: "internal/foo.go", AddedLines: hunks[0].AddedLines}}
	result := matcher.Match(nil, refs, matcher.DefaultConfig())
	links := matcher.BuildHunkLinks(hunks, nil, result)
	if links[0].Authorship != matcher.AuthorshipHuman {
		t.Fatalf("authorship = %q, want human", links[0].Authorship)
	}
}
