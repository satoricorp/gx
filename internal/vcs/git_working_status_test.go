package vcs

import "testing"

func TestParseGitWorkingStatus(t *testing.T) {
	out := " M unstaged.go\nA staged.go\n?? untracked.txt\nMM both.go\nR  old.go -> new.go\n"
	got := parseGitWorkingStatus(out)
	if len(got.Staged) != 3 || got.Staged[0] != "staged.go" || got.Staged[1] != "both.go" || got.Staged[2] != "new.go" {
		t.Fatalf("Staged = %v, want [staged.go both.go new.go]", got.Staged)
	}
	if len(got.Unstaged) != 2 || got.Unstaged[0] != "unstaged.go" || got.Unstaged[1] != "both.go" {
		t.Fatalf("Unstaged = %v, want [unstaged.go both.go]", got.Unstaged)
	}
	if len(got.Untracked) != 1 || got.Untracked[0] != "untracked.txt" {
		t.Fatalf("Untracked = %v, want [untracked.txt]", got.Untracked)
	}
}
