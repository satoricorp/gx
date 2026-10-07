package gxtest_test

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/satoricorp/gx/internal/gxtest"
)

// enclosedChildEnv marks the re-run of this test binary that
// TestGitHelpersIgnoreAnEnclosingRepository makes from inside a repository.
const enclosedChildEnv = "GXTEST_ENCLOSED_CHILD"

// TestGitHelpersIgnoreAnEnclosingRepository runs this test binary the way a
// git hook in a linked worktree runs it, with that worktree's GIT_DIR and
// GIT_INDEX_FILE inherited, and fails if anything the child's tests do with
// git lands in the enclosing repository instead of their own temp repos.
//
// The variables are inherited at process start rather than set with
// t.Setenv, because that is how the CLI suite met them on 2026-08-15. The
// child repeats what wrote `user.name = Test User` and `core.bare = true`
// into gx's own .git/config, through gxtest and through the bare
// exec.Command that package-local helpers use.
func TestGitHelpersIgnoreAnEnclosingRepository(t *testing.T) {
	world := gxtest.NewWorld(t)
	main := world.NewRepo(t)
	enclosing := main.AddWorktree(t, filepath.Join(t.TempDir(), "enclosing"), "enclosing")

	// Without this control, a git that ignored GIT_DIR would let the
	// comparison below pass vacuously.
	probe := exec.Command("git", "-C", t.TempDir(), "config", "gxtest.probe", "redirected")
	probe.Env = append(os.Environ(), "GIT_DIR="+enclosing.GitDir)
	if out, err := probe.CombinedOutput(); err != nil {
		t.Fatalf("control git config under GIT_DIR: %v\n%s", err, out)
	}
	if got := gxtest.GitOutput(t, main.Root, "config", "--get", "gxtest.probe"); got != "redirected" {
		t.Fatalf("GIT_DIR did not redirect `git config` into the enclosing repository (got %q), so this test cannot detect a leak", got)
	}

	before := repositoryState(t, enclosing)
	child := exec.Command(os.Args[0], "-test.run=^TestEnclosedChild$", "-test.count=1", "-test.v")
	child.Env = append(os.Environ(),
		enclosedChildEnv+"=1",
		"GIT_DIR="+enclosing.GitDir,
		"GIT_INDEX_FILE="+filepath.Join(enclosing.GitDir, "index"),
		// What a commit hook also receives when the commit was made with -c
		// and an identity: it outranks the identity the tests configure.
		"GIT_CONFIG_PARAMETERS='user.name'='Enclosing'",
		"GIT_AUTHOR_NAME=Enclosing",
		"GIT_AUTHOR_EMAIL=enclosing@example.com",
		"GIT_AUTHOR_DATE=@0 +0000",
	)
	out, err := child.CombinedOutput()
	// The leak is checked first: an escaped fixture usually breaks the child
	// too, and what it wrote is the more useful half of the report.
	if after := repositoryState(t, enclosing); after != before {
		t.Fatalf("a test binary run inside a hook wrote into the enclosing repository\n--- before\n%s\n--- after\n%s\n--- child output\n%s", before, after, out)
	}
	if err != nil {
		t.Fatalf("test binary run inside the enclosing repository failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "--- PASS: TestEnclosedChild") {
		t.Fatalf("child run did not execute TestEnclosedChild:\n%s", out)
	}
}

// TestEnclosedChild is the child half of
// TestGitHelpersIgnoreAnEnclosingRepository; run on its own it skips.
func TestEnclosedChild(t *testing.T) {
	if os.Getenv(enclosedChildEnv) == "" {
		t.Skip("runs only as the child of TestGitHelpersIgnoreAnEnclosingRepository")
	}
	repo := gxtest.NewWorld(t).NewRepo(t)
	if want := filepath.Join(repo.Root, ".git"); repo.GitCommonDir != want {
		t.Fatalf("NewRepo resolved its git dir to %s, want its own %s", repo.GitCommonDir, want)
	}

	// The incident's line, in runGitTest's shape: cmd.Dir rather than -C.
	cmd := exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = repo.Root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git config user.name: %v\n%s", err, out)
	}
	gxtest.Git(t, t.TempDir(), "init", "--bare", ".")
	repo.Commit(t, map[string]string{"second.txt": "second\n"}, "second")

	if got := gxtest.GitOutput(t, repo.Root, "config", "--local", "user.name"); got != "Test User" {
		t.Fatalf("temp repo user.name = %q, want Test User", got)
	}
	if got := gxtest.GitOutput(t, repo.Root, "log", "-1", "--format=%an <%ae>"); got != "Test User <test@example.com>" {
		t.Fatalf("temp repo commit author = %q, want the identity the test configured", got)
	}
	if got := gxtest.GitOutput(t, repo.Root, "log", "-1", "--format=%at"); got == "0" {
		t.Fatal("temp repo commit took the enclosing commit's GIT_AUTHOR_DATE")
	}
}

// repositoryState is everything an escaped git command could have written: the
// shared config, the hooks, every ref, and the worktree's HEAD and index.
func repositoryState(t *testing.T, worktree *gxtest.Repo) string {
	t.Helper()
	var b strings.Builder
	config, err := os.ReadFile(filepath.Join(worktree.GitCommonDir, "config"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	fmt.Fprintf(&b, "config:\n%s\n", config)
	hooks, err := os.ReadDir(filepath.Join(worktree.GitCommonDir, "hooks"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("read hooks: %v", err)
	}
	for _, hook := range hooks {
		fmt.Fprintf(&b, "hook: %s\n", hook.Name())
	}
	fmt.Fprintf(&b, "refs:\n%s\n", gxtest.GitOutput(t, worktree.Root, "for-each-ref", "--format=%(refname) %(objectname)"))
	for _, name := range []string{"HEAD", "index"} {
		data, err := os.ReadFile(filepath.Join(worktree.GitDir, name))
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("read %s: %v", name, err)
		}
		fmt.Fprintf(&b, "%s: %x\n", name, sha256.Sum256(data))
	}
	return b.String()
}

// TestEnclosingGitEnvCoversGitsOwnList fails when the git on this machine
// knows a repository-local variable that DetachFromEnclosingGit would leave
// in place.
func TestEnclosingGitEnvCoversGitsOwnList(t *testing.T) {
	out, err := exec.Command("git", "rev-parse", "--local-env-vars").Output()
	if err != nil {
		t.Fatalf("git rev-parse --local-env-vars: %v", err)
	}
	names := strings.Fields(string(out))
	if len(names) == 0 {
		t.Fatal("git rev-parse --local-env-vars printed nothing")
	}
	cleared := map[string]bool{}
	for _, name := range gxtest.EnclosingGitEnv {
		cleared[name] = true
	}
	var missing []string
	for _, name := range names {
		if !cleared[name] {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("git ties commands to a repository through %s, which gxtest.EnclosingGitEnv does not clear; add them", strings.Join(missing, ", "))
	}
}

// TestEveryPackageThatRunsGitDetachesFromTheEnclosingGit makes
// DetachFromEnclosingGit a rule rather than a habit: a package whose test
// binary can run git must open its TestMain with it. That is a package whose
// tests exec git, or one that links, at any depth, a package of this module
// that does.
//
// The second half is not a technicality. cmd/gx's tests never mention git,
// and run inside a hook they installed gx's lifecycle hooks into the enclosing
// repository through cli.NewRoot's auto-init. The code under test spawns git
// with the process environment like everything else, which is also why the
// call has to come before anything else in TestMain.
func TestEveryPackageThatRunsGitDetachesFromTheEnclosingGit(t *testing.T) {
	root := moduleRoot(t)
	packages := scanModule(t, root)

	var runners []string
	for path, pkg := range packages {
		if pkg.execsGit != "" {
			runners = append(runners, path)
		}
	}
	for _, known := range []string{"internal/gxtest", "internal/vcs"} {
		if pkg := packages[known]; pkg == nil || pkg.execsGit == "" {
			t.Fatalf("the scan found no git exec in %s, which certainly has one; it is broken (found: %v)", known, runners)
		}
	}

	// exposed maps each package whose test binary can run git to the reason.
	exposed := map[string]string{}
	for path, pkg := range packages {
		if !pkg.hasTests {
			continue
		}
		if pkg.testExecsGit != "" {
			exposed[path] = pkg.testExecsGit
			continue
		}
		// The package under test is linked into its own test binary, whether
		// the tests are internal or external.
		start := []string{path}
		for imported := range pkg.testImports {
			start = append(start, imported)
		}
		if chain := chainToGit(packages, start); chain != nil {
			exposed[path] = "links " + strings.Join(chain, " → ") + ", which runs git"
		}
	}
	if _, ok := exposed["cmd/gx"]; !ok {
		t.Fatal("the scan did not reach git from cmd/gx, whose tests run it through internal/cli; the transitive half is broken")
	}

	var offenders []string
	for path, reason := range exposed {
		if !packages[path].detaches {
			offenders = append(offenders, fmt.Sprintf("%s (%s)", path, reason))
		}
	}
	sort.Strings(offenders)
	if len(offenders) > 0 {
		t.Fatalf("these packages' test binaries can run git, but their TestMain does not start with gxtest.DetachFromEnclosingGit():\n  %s\n\n"+
			"Run from a git hook in a linked worktree, or from `git rebase --exec`, their git commands write to the repository running the hook "+
			"instead of the ones the tests create. Open the package's TestMain with it:\n\n"+
			"\tfunc TestMain(m *testing.M) {\n\t\tgxtest.DetachFromEnclosingGit()\n\t\tos.Exit(m.Run())\n\t}",
			strings.Join(offenders, "\n  "))
	}
}

// modulePackage is what the scan needs to know about one directory of this
// module. Paths are relative to the module root, slash-separated.
type modulePackage struct {
	imports      map[string]bool // module packages its non-test files import
	testImports  map[string]bool // module packages its test files import
	execsGit     string          // a non-test file that execs git
	testExecsGit string          // a test file that execs git, and how
	hasTests     bool
	detaches     bool
}

// scanModule parses every Go file of the module rooted at root.
func scanModule(t *testing.T, root string) map[string]*modulePackage {
	t.Helper()
	gomod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	module := ""
	for _, line := range strings.Split(string(gomod), "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "module "); ok {
			module = strings.TrimSpace(rest)
			break
		}
	}
	if module == "" {
		t.Fatal("go.mod declares no module path")
	}

	packages := map[string]*modulePackage{}
	fset := token.NewFileSet()
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path == root {
				return nil
			}
			// Dot directories include .claude/worktrees, which under a main
			// checkout hold whole other copies of this module.
			switch name := entry.Name(); {
			case strings.HasPrefix(name, "."), name == "vendor", name == "node_modules", name == "testdata":
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(path, "go.mod")); err == nil {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		dir, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		dir = filepath.ToSlash(dir)
		pkg := packages[dir]
		if pkg == nil {
			pkg = &modulePackage{imports: map[string]bool{}, testImports: map[string]bool{}}
			packages[dir] = pkg
		}
		isTest := strings.HasSuffix(path, "_test.go")
		imports := pkg.imports
		if isTest {
			imports = pkg.testImports
			pkg.hasTests = true
		}
		for _, imp := range file.Imports {
			if rel, ok := strings.CutPrefix(strings.Trim(imp.Path.Value, `"`), module+"/"); ok && rel != dir {
				imports[rel] = true
			}
		}
		use := execsGit(file)
		switch {
		case use == "":
		case isTest && pkg.testExecsGit == "":
			pkg.testExecsGit = filepath.Base(path) + " " + use
		case !isTest && pkg.execsGit == "":
			pkg.execsGit = filepath.Base(path)
		}
		if isTest && testMainDetaches(file) {
			pkg.detaches = true
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return packages
}

// chainToGit returns the shortest import chain from one of start to a package
// whose own code execs git, or nil when there is none.
func chainToGit(packages map[string]*modulePackage, start []string) []string {
	sort.Strings(start)
	parent := map[string]string{}
	queue := []string{}
	for _, path := range start {
		if _, seen := parent[path]; !seen {
			parent[path] = ""
			queue = append(queue, path)
		}
	}
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		pkg := packages[path]
		if pkg == nil {
			continue
		}
		if pkg.execsGit != "" {
			var chain []string
			for at := path; at != ""; at = parent[at] {
				chain = append([]string{at}, chain...)
			}
			return chain
		}
		next := make([]string, 0, len(pkg.imports))
		for imported := range pkg.imports {
			next = append(next, imported)
		}
		sort.Strings(next)
		for _, imported := range next {
			if _, seen := parent[imported]; !seen {
				parent[imported] = path
				queue = append(queue, imported)
			}
		}
	}
	return nil
}

// execsGit names how file execs git, or returns "" when it does not.
func execsGit(file *ast.File) string {
	use := ""
	ast.Inspect(file, func(node ast.Node) bool {
		if use != "" {
			return false
		}
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "exec" {
			return true
		}
		nameArg := -1
		switch sel.Sel.Name {
		case "Command":
			nameArg = 0
		case "CommandContext":
			nameArg = 1
		}
		if nameArg < 0 || len(call.Args) <= nameArg {
			return true
		}
		if lit, ok := call.Args[nameArg].(*ast.BasicLit); ok && lit.Kind == token.STRING && lit.Value == `"git"` {
			use = fmt.Sprintf("runs exec.%s(\"git\", ...)", sel.Sel.Name)
		}
		return true
	})
	return use
}

// testMainDetaches reports whether file declares a TestMain whose first
// statement calls DetachFromEnclosingGit.
func testMainDetaches(file *ast.File) bool {
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv != nil || fn.Name.Name != "TestMain" || fn.Body == nil || len(fn.Body.List) == 0 {
			continue
		}
		stmt, ok := fn.Body.List[0].(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := stmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		switch fun := call.Fun.(type) {
		case *ast.SelectorExpr:
			if fun.Sel.Name == "DetachFromEnclosingGit" {
				return true
			}
		case *ast.Ident:
			if fun.Name == "DetachFromEnclosingGit" {
				return true
			}
		}
	}
	return false
}

func moduleRoot(t *testing.T) string {
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
