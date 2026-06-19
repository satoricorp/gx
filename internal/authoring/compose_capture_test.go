package authoring

import (
	"testing"

	"github.com/satoricorp/gx/internal/capture"
	"github.com/satoricorp/gx/internal/capture/matcher"
)

func TestSessionIDsFromHunkLinks(t *testing.T) {
	links := []matcher.HunkLink{
		{HunkID: "h1", SessionID: "s1", Tier: matcher.TierExact, Authorship: matcher.AuthorshipAgent},
		{HunkID: "h2", SessionID: "s1", Tier: matcher.TierFuzzy, Authorship: matcher.AuthorshipAgent},
		{HunkID: "h3", SessionID: "s2", Tier: matcher.TierTemporal, Authorship: matcher.AuthorshipUnknown},
		{HunkID: "h4", Authorship: matcher.AuthorshipHuman},
	}
	got := sessionIDsFromHunkLinks(links)
	if len(got) != 1 || got[0] != "s1" {
		t.Fatalf("sessionIDsFromHunkLinks() = %#v, want [s1]", got)
	}
}

func TestHunkLinksForRevision(t *testing.T) {
	links := []matcher.HunkLink{
		{HunkID: "working:internal/a.go:1-3"},
		{HunkID: "working:internal/b.go:1-2"},
	}
	got := hunkLinksForRevision(links, []string{"internal/a.go"})
	if len(got) != 1 || got[0].HunkID != links[0].HunkID {
		t.Fatalf("hunkLinksForRevision() = %#v", got)
	}
}

func TestAddedLinesFromPatch(t *testing.T) {
	patch := "@@ -1 +1,2 @@\n-old\n+new\n+line\n"
	got := addedLinesFromPatch(patch)
	want := []string{"new", "line"}
	if len(got) != len(want) {
		t.Fatalf("addedLinesFromPatch() = %#v, want %#v", got, want)
	}
}

func TestRemapHunkLinkIDs(t *testing.T) {
	demuxHunks := []HunkRange{{ID: "h1", File: "a.go", Patch: "+x\n"}}
	eligible := []capture.CommitHunk{{FilePath: "a.go"}}
	links := []matcher.HunkLink{{HunkID: "working:a.go:1-1"}}
	got := remapHunkLinkIDs(links, demuxHunks, eligible)
	if got[0].HunkID != "h1" {
		t.Fatalf("remapHunkLinkIDs() hunk id = %q, want h1", got[0].HunkID)
	}
}
