package vcs

import "testing"

func TestResolveStackRevision(t *testing.T) {
	units := []UnitSummary{
		{Index: 1, ChangeID: "aaa111", CommitID: "commit111"},
		{Index: 2, ChangeID: "bbb222", CommitID: "commit222"},
	}

	got, err := ResolveStackRevision("r2", units)
	if err != nil {
		t.Fatalf("ResolveStackRevision(r2) error = %v", err)
	}
	if got != "bbb222" {
		t.Fatalf("ResolveStackRevision(r2) = %q, want bbb222", got)
	}

	got, err = ResolveStackRevision("aaa", units)
	if err != nil {
		t.Fatalf("ResolveStackRevision(prefix) error = %v", err)
	}
	if got != "aaa111" {
		t.Fatalf("ResolveStackRevision(prefix) = %q, want aaa111", got)
	}
}
