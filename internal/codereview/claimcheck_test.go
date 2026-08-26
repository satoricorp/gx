package codereview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The cases here are the real ones from a user's triage of a 20-finding
// review: which of them the repository could have settled before the finding
// shipped, and which must survive because it cannot.
func repoWith(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestDropsAFindingTheFileDisproves(t *testing.T) {
	// "playwright missing from pyproject.toml" — it was on line 36.
	root := repoWith(t, map[string]string{
		"pyproject.toml": strings.Repeat("# padding\n", 35) + `playwright = ">=1.49.0"` + "\n",
	})
	findings := []Finding{{
		Title:   "Playwright dependency is not declared",
		Summary: "`playwright` is missing from the project dependencies, so the browser probe cannot run in CI.",
		File:    "pyproject.toml",
	}}

	kept, notes := DropContradictedFindings(root, findings)
	if len(kept) != 0 {
		t.Fatalf("kept %d findings, want the contradicted one dropped", len(kept))
	}
	if len(notes) != 1 || !strings.Contains(notes[0], "pyproject.toml:36") {
		t.Fatalf("notes = %v, want the line that disproves it", notes)
	}
}

// The finding that must survive: the negation governs the None check, not
// `None` itself, and `None` appears in every Python file that mentions it.
// Reading this as "None is absent" would delete a real defect.
func TestKeepsAFindingWhoseNegationIsNotAboutTheSymbol(t *testing.T) {
	root := repoWith(t, map[string]string{
		"live_audit.py": "value = node.get_attribute('href') if node else None\n",
	})
	findings := []Finding{{
		Title:   "Missing None check",
		Summary: "`get_attribute` is called without a `None` check, so a missing node raises.",
		File:    "live_audit.py",
	}}

	kept, notes := DropContradictedFindings(root, findings)
	if len(kept) != 1 {
		t.Fatalf("kept %d findings, want the finding preserved", len(kept))
	}
	if len(notes) != 0 {
		t.Fatalf("notes = %v, want none", notes)
	}
}

func TestKeepsAClaimTheRepositoryConfirms(t *testing.T) {
	root := repoWith(t, map[string]string{"pyproject.toml": "[project]\nname = \"x\"\n"})
	findings := []Finding{{
		Title:   "Playwright dependency is not declared",
		Summary: "`playwright` is missing from the project dependencies.",
		File:    "pyproject.toml",
	}}

	kept, _ := DropContradictedFindings(root, findings)
	if len(kept) != 1 {
		t.Fatal("a true absence claim was dropped")
	}
}

func TestKeepsWhatItCannotCheck(t *testing.T) {
	root := repoWith(t, map[string]string{"a.go": "package a\n"})
	cases := []Finding{
		// No anchor: nowhere to look.
		{Title: "x", Summary: "`somethingGone` is missing", File: ""},
		// The file does not exist.
		{Title: "x", Summary: "`somethingGone` is missing", File: "nope.go"},
		// No claim shape at all.
		{Title: "x", Summary: "this function is hard to read", File: "a.go"},
	}
	kept, notes := DropContradictedFindings(root, cases)
	if len(kept) != len(cases) {
		t.Fatalf("kept %d of %d, want all preserved", len(kept), len(cases))
	}
	if len(notes) != 0 {
		t.Fatalf("notes = %v, want none", notes)
	}
}

// A finding's file field is model output. It must not be able to read a path
// outside the repository, whatever it contains.
func TestNeverReadsOutsideTheRepository(t *testing.T) {
	root := repoWith(t, map[string]string{"a.go": "package a\n"})
	findings := []Finding{{
		Title:   "x",
		Summary: "`root` is missing",
		File:    "../../../../../../etc/passwd",
	}}
	kept, notes := DropContradictedFindings(root, findings)
	if len(kept) != 1 || len(notes) != 0 {
		t.Fatalf("kept=%d notes=%v, want the finding preserved and nothing read", len(kept), notes)
	}
}

// Whole-word matching: a claim about `pid` is not settled by `rapid`.
func TestSubstringIsNotAMatch(t *testing.T) {
	root := repoWith(t, map[string]string{"x.ts": "const rapidly = 1;\n"})
	findings := []Finding{{
		Title:   "x",
		Summary: "`pid` is missing from the stop response.",
		File:    "x.ts",
	}}
	kept, _ := DropContradictedFindings(root, findings)
	if len(kept) != 1 {
		t.Fatal("a substring match dropped a finding that stands")
	}
}
