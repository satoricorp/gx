package codereview

import "testing"

func TestResolveSourceLabelsUsesSourceRefsAndCatalog(t *testing.T) {
	brief := ReviewBrief{
		SourceRefs: []SourceRef{
			{ID: "R1", Kind: "resource", Publisher: "OWASP", Title: "hidden"},
			{ID: "L1", Kind: "local", Publisher: "local", File: "REVIEW.md", StartLine: 12, Title: "REVIEW.md"},
		},
		SourceCatalog: []SourceBrief{{ID: "google-eng-practices", Publisher: "Google"}},
	}
	resolved := resolveSourceLabels(brief, []string{"L1", "google-eng-practices", "unknown"})
	if len(resolved) != 2 {
		t.Fatalf("resolved = %#v", resolved)
	}
	if resolved[0].Kind != "local" || resolved[0].Opaque || resolved[0].File != "REVIEW.md" || resolved[0].StartLine != 12 {
		t.Fatalf("local resolved = %#v", resolved[0])
	}
	if resolved[1].Kind != "catalog" || !resolved[1].Opaque || resolved[1].Publisher != "Google" {
		t.Fatalf("catalog resolved = %#v", resolved[1])
	}
}

func TestResolvedSourceLabel(t *testing.T) {
	if got := ResolvedSourceLabel(ResolvedSource{Opaque: true, Publisher: "OWASP"}); got != "OWASP" {
		t.Fatalf("opaque label = %q", got)
	}
	if got := ResolvedSourceLabel(ResolvedSource{File: "REVIEW.md", StartLine: 4}); got != "REVIEW.md:4" {
		t.Fatalf("file label = %q", got)
	}
}
