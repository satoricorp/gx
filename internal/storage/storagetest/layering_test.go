package storagetest_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// TestHarnessLayering keeps the test harness out of internal/vcs.
//
// The constraint is real in both directions. internal/vcs imports
// internal/storage, so a harness under internal/storage that imported vcs would
// cycle the moment an in-package storage test used it. And vcs pulls in the
// capture pipeline and the git subprocess layer, so importing it would add
// roughly fifteen seconds of package init to the test binaries of
// internal/storage, internal/reviewbundle and internal/provenance — the fast
// ones, which is exactly where a shared fixture has to stay cheap or nobody
// uses it.
//
// The layering is: internal/lgtmtest imports nothing from lgtm at all (which is why
// it carries its own copy of the revision trailer, pinned by
// TestLgtmTestRevisionTrailerMatchesProduction in internal/hooks), and
// internal/storage/storagetest imports only internal/storage and
// internal/lgtmtest. Checking internal/storage too closes the transitive hole:
// if storage ever imported vcs, storagetest would inherit it.
func TestHarnessLayering(t *testing.T) {
	root := repoRoot(t)
	const prefix = "github.com/satoricorp/lgtm/"

	for _, tc := range []struct {
		pkg     string
		allowed []string
	}{
		{pkg: "internal/lgtmtest", allowed: nil},
		{pkg: "internal/storage", allowed: []string{"internal/agentprovenance"}},
		{pkg: "internal/storage/storagetest", allowed: []string{"internal/storage", "internal/lgtmtest"}},
	} {
		allowed := map[string]bool{}
		for _, name := range tc.allowed {
			allowed[name] = true
		}
		for _, imported := range lgtmImports(t, filepath.Join(root, tc.pkg), prefix) {
			if !allowed[imported] {
				t.Errorf("%s imports %s, which the harness layering forbids; allowed: %v",
					tc.pkg, imported, tc.allowed)
			}
		}
	}
}

// lgtmImports returns the lgtm-internal packages a package's non-test files import.
func lgtmImports(t *testing.T, dir, prefix string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	found := map[string]bool{}
	fset := token.NewFileSet()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		for _, spec := range file.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				continue
			}
			if strings.HasPrefix(path, prefix) {
				found[strings.TrimPrefix(path, prefix)] = true
			}
		}
	}
	var out []string
	for name := range found {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
