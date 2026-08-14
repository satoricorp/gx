package codereview

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Engine.Review is the only entry point that resolves a review subject:
// resolveChangeSet, reviewChangeSet, and normalizeOptions are reachable from
// nowhere else. That matters because this package is shared. The PR-summary
// pipeline in internal/publication uses its types, its brief, and its AI
// reviewer, but takes its subject from the pushed bundle rather than from the
// working tree, so changing how `gx enhance` picks a subject cannot move the
// summaries gx writes on every push.
//
// Pinning the caller set keeps that true by construction: wiring the engine
// into another pipeline becomes a deliberate act with a failing test attached,
// rather than a silent change to every PR summary.
func TestReviewEngineIsReachedOnlyFromTheReviewCommand(t *testing.T) {
	root := filepath.Join("..", "..")
	// "codereview.NewEngine" with no open paren on purpose: it also covers
	// NewEngineWith and NewEngineWithReviewer, which return the same *Engine
	// and run the same subject resolution. Matching the exact call shape would
	// have left the constructor a publish-path wiring would actually reach —
	// publication already builds its own AI reviewer — outside the guard.
	entryPoints := []string{"codereview.Review(", "codereview.NewEngine"}
	assertEntryPointsCoverEngineConstructors(t, entryPoints)

	callers := map[string]struct{}{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".gocache", "node_modules", "dist", "build", ".next", "vendor", "testdata":
				return fs.SkipDir
			case ".claude":
				// Agent worktrees nest a whole checkout under the repo root.
				// Walking into one counts its copy of internal/cli as a second
				// caller, so this guard fails for a reason that has nothing to
				// do with the wiring it exists to protect.
				return fs.SkipDir
			}
			return nil
		}
		// Production wiring only: a test may drive the engine from anywhere.
		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(data)
		for _, entryPoint := range entryPoints {
			if !strings.Contains(text, entryPoint) {
				continue
			}
			rel, relErr := filepath.Rel(root, filepath.Dir(path))
			if relErr != nil {
				return relErr
			}
			callers[filepath.ToSlash(rel)] = struct{}{}
			break
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir() error = %v", err)
	}

	got := make([]string, 0, len(callers))
	for pkg := range callers {
		got = append(got, pkg)
	}
	sort.Strings(got)

	want := []string{"internal/cli"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("packages calling the review engine = %v, want %v.\n"+
			"A new caller inherits the review subject resolution, including the "+
			"working-tree-wins rule --repo exists to override. If that is "+
			"intended, add the package here and say what its subject is.", got, want)
	}
}

// A guard that misses a way in is worse than no guard: it reports safety it is
// not providing. Every exported way to build an engine has to be covered by
// the entry-point list, so adding a constructor cannot quietly open a door.
func assertEntryPointsCoverEngineConstructors(t *testing.T, entryPoints []string) {
	t.Helper()
	data, err := os.ReadFile("engine.go")
	if err != nil {
		t.Fatalf("ReadFile(engine.go) error = %v", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "func New") {
			continue
		}
		name, _, ok := strings.Cut(strings.TrimPrefix(line, "func "), "(")
		if !ok || !strings.HasSuffix(strings.TrimSpace(line), "*Engine {") {
			continue
		}
		call := "codereview." + name
		covered := false
		for _, entryPoint := range entryPoints {
			if strings.HasPrefix(call+"(", entryPoint) {
				covered = true
				break
			}
		}
		if !covered {
			t.Fatalf("constructor %s is not matched by any entry point in %v; "+
				"a caller using it would slip past this guard", call, entryPoints)
		}
	}
}
