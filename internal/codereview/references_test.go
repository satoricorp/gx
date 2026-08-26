package codereview

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// gitRepoWith builds a real repository: FindReferences shells out to git grep,
// which reads the index, so a bare directory would report nothing and the test
// would pass for the wrong reason.
func gitRepoWith(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	run := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "t@t.t")
	run("config", "user.name", "t")
	for name, body := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run("add", "-A")
	run("commit", "-qm", "init")
	return root
}

// The finding said onPaid was never used. Both call sites were right there.
func TestFindReferencesLocatesEveryCallSite(t *testing.T) {
	root := gitRepoWith(t, map[string]string{
		"page.tsx":               "return <PreviewGateBanner onPaid={onPaid} />;\n",
		"use-lifecycle-state.ts": "export function useLifecycle({ onPaid }) {}\n",
		"unrelated.ts":           "const x = 1;\n",
	})

	refs := FindReferences(context.Background(), root, "onPaid")
	if len(refs) != 2 {
		t.Fatalf("FindReferences() = %d refs, want 2: %+v", len(refs), refs)
	}
	if refs[0].File != "page.tsx" || refs[0].Line != 1 {
		t.Fatalf("first ref = %+v, want page.tsx:1", refs[0])
	}
	if !strings.Contains(refs[0].Text, "PreviewGateBanner") {
		t.Fatalf("ref text = %q, want the line itself", refs[0].Text)
	}
}

// Whole-word only: a claim about `pid` is not answered by `rapid`.
func TestFindReferencesIgnoresSubstrings(t *testing.T) {
	root := gitRepoWith(t, map[string]string{"x.ts": "const rapidly = 1;\n"})
	if refs := FindReferences(context.Background(), root, "pid"); len(refs) != 0 {
		t.Fatalf("FindReferences() = %+v, want no substring matches", refs)
	}
}

func TestAttachSymbolReferencesGivesTheFindingItsReceipts(t *testing.T) {
	root := gitRepoWith(t, map[string]string{
		"page.tsx":               "return <PreviewGateBanner onPaid={onPaid} />;\n",
		"use-lifecycle-state.ts": "export function useLifecycle({ onPaid }) {}\n",
	})
	findings := []Finding{{
		Title:   "`onPaid` is not threaded through every caller",
		Summary: "Callers may drop the handler.",
		File:    "page.tsx",
	}}

	out := AttachSymbolReferences(context.Background(), root, findings)
	var value string
	for _, evidence := range out[0].Evidence {
		if evidence.Label == referencesEvidenceLabel {
			value = evidence.Value
		}
	}
	if value == "" {
		t.Fatalf("no reference evidence attached: %+v", out[0].Evidence)
	}
	for _, want := range []string{"onPaid", "2 references", "page.tsx:1", "use-lifecycle-state.ts:1"} {
		if !strings.Contains(value, want) {
			t.Fatalf("evidence = %q, want it to contain %q", value, want)
		}
	}
}

// A name that appears everywhere identifies nothing; reporting an arbitrary
// handful of hundreds is worse than reporting none.
func TestFindReferencesStaysQuietForCommonNames(t *testing.T) {
	files := map[string]string{}
	for i := 0; i < referenceNoiseThreshold+5; i++ {
		files[filepath.Join("pkg", "f"+strings.Repeat("x", i%7)+string(rune('a'+i%26))+".go")] =
			"package pkg\n// handler\n"
	}
	root := gitRepoWith(t, files)
	if refs := FindReferences(context.Background(), root, "handler"); len(refs) != 0 {
		t.Fatalf("FindReferences() returned %d refs for a ubiquitous name", len(refs))
	}
}
