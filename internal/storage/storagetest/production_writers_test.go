package storagetest_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// allowlistPath holds the storage methods that are currently seeded only by
// tests. It is a ratchet: entries may be deleted, never added.
const allowlistPath = "testdata/test_only_writers.txt"

// TestSeedersAreProductionWriters fails when a test seeds the database through
// a storage method that production never calls.
//
// This is the guard chosen over porting the individual call sites, and the
// reason is that porting fixes one function while the rule fixes the class.
// store.WriteSession had zero production callers and fourteen test call sites;
// meanwhile the only real writer of `sessions`, UpsertObservedSession, could
// never bootstrap its first row. Every session test seeded around that, so the
// suite proved a state the product could not reach and 171 of 174 published
// bundles shipped `sessions: []`. Porting WriteSession's call sites would have
// left thirteen more functions in the same position — WriteRequest,
// WriteResponse, UpsertSession, UpsertSessionContext, the demux writers — each
// one able to hide the next defect exactly the same way. A mechanical check
// makes the next one impossible to introduce silently, and turns each existing
// one into a listed, dated debt instead of an invisible assumption.
//
// The rule: seed through the writer production uses, or do not seed. Where the
// production writer genuinely cannot express the row a test needs (daemon-era
// columns such as process_name and last_seen_at, which nothing writes any
// more), that is the finding, not the workaround — the entry stays on the
// allowlist until the test is rewritten or the column is dropped.
//
// Known limitation: callers are matched by method name alone, so an unrelated
// `x.Close()` counts as a call to Store.Close. That only ever makes this test
// more lenient, never more strict.
func TestSeedersAreProductionWriters(t *testing.T) {
	root := repoRoot(t)
	methods := storeMethods(t, filepath.Join(root, "internal", "storage"))
	if len(methods) == 0 {
		t.Fatal("found no exported *Store/*CaptureStage methods; the scan is broken")
	}

	prodCallers, testCallers := callerCounts(t, root, methods)

	var offenders []string
	for name := range methods {
		if prodCallers[name] == 0 && testCallers[name] > 0 {
			offenders = append(offenders, name)
		}
	}
	sort.Strings(offenders)

	allowed := readAllowlist(t, allowlistPath)
	var unlisted []string
	for _, name := range offenders {
		if !allowed[name] {
			unlisted = append(unlisted, name)
		}
	}
	if len(unlisted) > 0 {
		t.Fatalf("these storage methods are called only by tests, so seeding through them proves a state production cannot reach:\n  %s\n\n"+
			"Seed through the writer production uses instead (for `sessions` that is UpsertObservedSession, the writer vcs.AttachSessionsFromHunkLinks calls).\n"+
			"If the method is genuinely unreachable-but-needed, add it to internal/storage/storagetest/%s with a reason.",
			strings.Join(unlisted, "\n  "), allowlistPath)
	}

	// The ratchet. An allowlist entry that is no longer an offender must be
	// removed, otherwise the list decays into a permanent exemption for
	// functions that have since been fixed or deleted.
	offending := map[string]bool{}
	for _, name := range offenders {
		offending[name] = true
	}
	var stale []string
	for name := range allowed {
		if !offending[name] {
			stale = append(stale, name)
		}
	}
	sort.Strings(stale)
	if len(stale) > 0 {
		t.Fatalf("stale entries in %s — these methods now have production callers (or no test callers) and must be removed from the allowlist:\n  %s",
			allowlistPath, strings.Join(stale, "\n  "))
	}
}

// storeMethods returns the exported methods declared on *Store and
// *CaptureStage in the storage package.
func storeMethods(t *testing.T, dir string) map[string]bool {
	t.Helper()
	methods := map[string]bool{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	fset := token.NewFileSet()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", entry.Name(), err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) != 1 {
				continue
			}
			star, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			ident, ok := star.X.(*ast.Ident)
			if !ok || (ident.Name != "Store" && ident.Name != "CaptureStage") {
				continue
			}
			if fn.Name.IsExported() {
				methods[fn.Name.Name] = true
			}
		}
	}
	return methods
}

// callerCounts walks every Go file in the module and counts calls to each
// method name, split by whether the calling file is a test.
func callerCounts(t *testing.T, root string, methods map[string]bool) (prod, tests map[string]int) {
	t.Helper()
	prod = map[string]int{}
	tests = map[string]int{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "vendor", "node_modules", "testdata":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return nil
		}
		counts := prod
		if strings.HasSuffix(path, "_test.go") {
			counts = tests
		}
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if methods[sel.Sel.Name] {
				counts[sel.Sel.Name]++
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return prod, tests
}

func readAllowlist(t *testing.T, path string) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read allowlist %s: %v", path, err)
	}
	allowed := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, _, _ := strings.Cut(line, " ")
		allowed[strings.TrimSpace(name)] = true
	}
	return allowed
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find go.mod above the test package")
		}
		dir = parent
	}
}
