package codereview

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	maxStaticToolOutputBytes = 12000
	// maxStaticToolFileArgs caps how many changed files a file-scoped tool is
	// handed so one sweeping change cannot blow up the command line.
	maxStaticToolFileArgs = 40
	// maxStaticToolConfigBytes caps config files read while detecting tools.
	maxStaticToolConfigBytes = 256 * 1024
)

type StaticToolResult struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output,omitempty"`
	Skipped  bool   `json:"skipped,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// staticToolEnv is what every runner detects against: the repo on disk, the
// facts already gathered for the review, and the files in this change.
type staticToolEnv struct {
	repoRoot     string
	facts        RepoFacts
	changedFiles []string
}

// staticToolCommand is a resolved invocation. bin is the executable gx runs;
// argv is what humans and the AI brief see, so argv[0] stays the plain tool
// name even when bin points at a repo-local binary.
type staticToolCommand struct {
	bin  string
	argv []string
}

// staticToolRunner is one entry in the static tool registry. detect inspects
// the repo layout, the files in review, and the installed binaries, returning
// the invocation to run from the repo root; ok is false when the ecosystem
// does not apply or the tool is not installed, and the runner is then skipped
// silently — a missing tool is never a review finding.
//
// Runners must be fast and scoped to the change. They are NOT side-effect
// free, and it is worth being exact about that because the rest of `gx review`
// is: outside Go the runners are type checkers and linters, but the Go runner
// is `go test`, which compiles and executes the reviewed checkout's own test
// binaries, and `cargo check` executes the crate's build.rs and proc macros.
// On someone else's checkout that is code from the change under review running
// on the reviewer's machine. staticToolChildEnv keeps the reviewer's
// credentials out of it; whether it should run unprompted at all is an open
// question, not something this comment should imply is settled.
//
// They also write: `cargo check` populates <repo>/target/, `tsc` writes
// *.tsbuildinfo when tsconfig sets incremental or composite, and every Go run
// shares a build cache at $TMPDIR/gx-review-gocache that nothing prunes.
type staticToolRunner struct {
	name     string
	progress string
	detect   func(env staticToolEnv) (staticToolCommand, bool)
	// wholeProject marks a runner whose cost is set by the repository rather
	// than by the change: a test suite or a whole-program type check takes the
	// same minutes for a one-line diff as for a hundred-file one. A fast review
	// skips these — measured on this repository, `go test` alone was 18 of a
	// 48-second review, more than half the budget for a signal the developer
	// can get faster by running the suite themselves.
	//
	// Linters and per-file checks are not marked: they are proportional to the
	// change and cheap enough to keep even when optimizing for wall clock.
	wholeProject bool
}

// staticToolRunners is the registry; adding an ecosystem is one entry plus its
// detect function. Order is the order results are reported in.
var staticToolRunners = []staticToolRunner{
	{name: "go test", progress: "Running go test", detect: goStaticTool("test"), wholeProject: true},
	{name: "go vet", progress: "Running go vet", detect: goStaticTool("vet")},
	{name: "tsc", progress: "Running tsc", detect: detectTypeScriptCompiler, wholeProject: true},
	{name: "eslint", progress: "Running eslint", detect: detectESLint},
	{name: "ruff", progress: "Running ruff", detect: detectRuff},
	{name: "mypy", progress: "Running mypy", detect: detectMypy, wholeProject: true},
	{name: "cargo check", progress: "Running cargo check", detect: detectCargoCheck, wholeProject: true},
	{name: "dart analyze", progress: "Running dart analyze", detect: detectDartAnalyze},
	{name: "flutter analyze", progress: "Running flutter analyze", detect: detectFlutterAnalyze, wholeProject: true},
	{name: "dotnet build", progress: "Running dotnet build", detect: detectDotnetBuild, wholeProject: true},
}

// staticToolScope is the file set every runner detects against. It is the
// change set, except for a whole-repo review that has no diff: there the
// repository stands in for it. Without that, `gx review --repo` on a clean
// tree runs no compiler, no test, and no linter — collectStaticToolResults
// returns before detection on an empty set — and then reports the repository
// clean, which is the silent pass --fail-on exists to prevent.
func staticToolScope(facts RepoFacts, opts Options, changed []string) []string {
	scope := normalizedChangedFiles(changed)
	if len(scope) > 0 || !opts.WholeRepo {
		return scope
	}
	return repoWideStaticToolScope(facts)
}

// collectStaticToolResults runs the detected checkers scoped to changed, the
// files this review resolved — the caller owns that resolution so the tools see
// the same change set as the rest of the review, working tree or ref range.
func collectStaticToolResults(ctx context.Context, repoRoot string, facts RepoFacts, opts Options, changed []string) []StaticToolResult {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_STATIC_TOOLS")), "0") {
		return nil
	}
	env := staticToolEnv{
		repoRoot:     repoRoot,
		facts:        facts,
		changedFiles: staticToolScope(facts, opts, changed),
	}
	if len(env.changedFiles) == 0 {
		return nil
	}
	type plannedStaticTool struct {
		runner  staticToolRunner
		command staticToolCommand
	}
	var planned []plannedStaticTool
	var skippedForSpeed []string
	for _, runner := range staticToolRunners {
		if opts.Fast && runner.wholeProject {
			skippedForSpeed = append(skippedForSpeed, runner.name)
			continue
		}
		command, ok := runner.detect(env)
		if !ok {
			continue
		}
		planned = append(planned, plannedStaticTool{runner: runner, command: command})
	}
	// Say what was not run. A review missing its test results must not read
	// like a review whose tests passed.
	if len(skippedForSpeed) > 0 {
		reviewProgress(opts, "Skipping "+strings.Join(skippedForSpeed, ", ")+" (fast review)")
	}
	if len(planned) == 0 {
		return nil
	}
	timeout := 90 * time.Second
	if opts.Deep {
		timeout = 180 * time.Second
	}
	out := make([]StaticToolResult, len(planned))
	if !opts.Deep {
		for i, tool := range planned {
			reviewProgress(opts, tool.runner.progress)
			out[i] = runStaticTool(ctx, repoRoot, timeout, tool.runner.name, tool.command)
		}
		return out
	}

	var wg sync.WaitGroup
	for i, tool := range planned {
		i, tool := i, tool
		wg.Add(1)
		go func() {
			defer wg.Done()
			reviewProgress(opts, tool.runner.progress)
			out[i] = runStaticTool(ctx, repoRoot, timeout, tool.runner.name, tool.command)
		}()
	}
	wg.Wait()
	return out
}

// goStaticTool runs `go <sub>` over the packages the change touches, falling
// back to the whole module when the change is wide enough to affect it.
func goStaticTool(sub string) func(staticToolEnv) (staticToolCommand, bool) {
	return func(env staticToolEnv) (staticToolCommand, bool) {
		if !hasDependencyFile(env.facts.DependencyFiles, "go.mod") {
			return staticToolCommand{}, false
		}
		pkgs, ok := goPackageArgsForChangedFiles(env.repoRoot, env.changedFiles)
		if !ok {
			return staticToolCommand{}, false
		}
		bin, ok := lookStaticTool(env.repoRoot, "go")
		if !ok {
			return staticToolCommand{}, false
		}
		return staticToolCommand{bin: bin, argv: append([]string{"go", sub}, pkgs...)}, true
	}
}

func goPackageArgsForChangedFiles(repoRoot string, files []string) ([]string, bool) {
	files = normalizedChangedFiles(files)
	if len(files) == 0 {
		return nil, false
	}
	dirs := map[string]struct{}{}
	for _, file := range files {
		if goStaticToolsRequireRepoWide(file) {
			return []string{"./..."}, true
		}
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		dir := filepath.ToSlash(filepath.Dir(file))
		if dir == "" {
			return []string{"./..."}, true
		}
		if dir == "." {
			dirs["."] = struct{}{}
			continue
		}
		if _, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(dir))); err != nil {
			return []string{"./..."}, true
		}
		dirs[dir] = struct{}{}
	}
	if len(dirs) == 0 {
		return nil, false
	}
	if len(dirs) > 20 {
		return []string{"./..."}, true
	}
	out := make([]string, 0, len(dirs))
	for dir := range dirs {
		if dir == "." {
			out = append(out, ".")
		} else {
			out = append(out, "./"+dir)
		}
	}
	sort.Strings(out)
	return out, true
}

func goStaticToolsRequireRepoWide(file string) bool {
	file = filepath.ToSlash(strings.TrimSpace(file))
	base := filepath.Base(file)
	switch base {
	case "go.mod", "go.sum", "go.work", "go.work.sum", "Makefile", "Taskfile.yml", "Taskfile.yaml", "Dockerfile", "magefile.go":
		return true
	}
	return strings.HasPrefix(file, ".github/workflows/") ||
		strings.HasPrefix(file, "build/") ||
		strings.HasPrefix(file, "scripts/")
}

// detectTypeScriptCompiler type-checks the whole project rather than the
// changed files: tsc ignores tsconfig.json as soon as files are named on the
// command line, which would check them with the wrong compiler options.
func detectTypeScriptCompiler(env staticToolEnv) (staticToolCommand, bool) {
	if !repoFileExists(env.repoRoot, "tsconfig.json") {
		return staticToolCommand{}, false
	}
	if !changedFilesInclude(env.changedFiles, typeScriptFileExtensions...) &&
		!changedFilesIncludeBase(env.changedFiles, "tsconfig.json") {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "tsc", nodeLocalBinDirs...)
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: []string{"tsc", "--noEmit"}}, true
}

func detectESLint(env staticToolEnv) (staticToolCommand, bool) {
	if !repoFileExists(env.repoRoot, eslintConfigFiles...) &&
		!repoFileContains(env.repoRoot, "package.json", `"eslintConfig"`) {
		return staticToolCommand{}, false
	}
	files := staticToolFileArgs(env.repoRoot, env.changedFiles, javaScriptFileExtensions...)
	if len(files) == 0 {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "eslint", nodeLocalBinDirs...)
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: append([]string{"eslint"}, files...)}, true
}

// detectRuff needs no manifest: ruff runs without configuration, so changed
// Python files plus an installed binary are signal enough.
func detectRuff(env staticToolEnv) (staticToolCommand, bool) {
	files := staticToolFileArgs(env.repoRoot, env.changedFiles, pythonFileExtensions...)
	if len(files) == 0 {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "ruff", pythonLocalBinDirs...)
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: append([]string{"ruff", "check"}, files...)}, true
}

// detectMypy only runs where mypy is configured: unconfigured runs over single
// files report import errors that say nothing about the change.
func detectMypy(env staticToolEnv) (staticToolCommand, bool) {
	if !mypyConfigured(env.repoRoot) {
		return staticToolCommand{}, false
	}
	files := staticToolFileArgs(env.repoRoot, env.changedFiles, pythonFileExtensions...)
	if len(files) == 0 {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "mypy", pythonLocalBinDirs...)
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: append([]string{"mypy"}, files...)}, true
}

func mypyConfigured(repoRoot string) bool {
	return repoFileExists(repoRoot, "mypy.ini", ".mypy.ini") ||
		repoFileContains(repoRoot, "pyproject.toml", "[tool.mypy]") ||
		repoFileContains(repoRoot, "setup.cfg", "[mypy]")
}

// detectCargoCheck stops before codegen and never runs the crate's tests or
// binaries, but "check" is not "runs nothing": build.rs scripts and proc
// macros execute during a check, and the run writes <repo>/target/. Cargo has
// no per-file scope, so the manifest plus any Rust or dependency change is the
// trigger.
func detectCargoCheck(env staticToolEnv) (staticToolCommand, bool) {
	if !repoFileExists(env.repoRoot, "Cargo.toml") {
		return staticToolCommand{}, false
	}
	if !changedFilesInclude(env.changedFiles, ".rs") &&
		!changedFilesIncludeBase(env.changedFiles, "Cargo.toml", "Cargo.lock") {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "cargo")
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: []string{"cargo", "check"}}, true
}

// detectDartAnalyze lints the changed Dart files of a pure Dart package. A
// pubspec that declares flutter is detectFlutterAnalyze's territory: `dart
// analyze` cannot resolve Flutter SDK imports, so running it there reports
// errors the repo's own tooling never would.
func detectDartAnalyze(env staticToolEnv) (staticToolCommand, bool) {
	if !repoFileExists(env.repoRoot, "pubspec.yaml") || pubspecDeclaresFlutter(env.repoRoot) {
		return staticToolCommand{}, false
	}
	files := staticToolFileArgs(env.repoRoot, env.changedFiles, dartFileExtensions...)
	if len(files) == 0 {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "dart")
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: append([]string{"dart", "analyze"}, files...)}, true
}

// detectFlutterAnalyze covers the pubspecs detectDartAnalyze declines. It has
// no per-file scope — `flutter analyze` walks the whole project — which is why
// the registry marks it wholeProject while plain dart analyze stays per-file.
func detectFlutterAnalyze(env staticToolEnv) (staticToolCommand, bool) {
	if !repoFileExists(env.repoRoot, "pubspec.yaml") || !pubspecDeclaresFlutter(env.repoRoot) {
		return staticToolCommand{}, false
	}
	if !changedFilesInclude(env.changedFiles, dartFileExtensions...) &&
		!changedFilesIncludeBase(env.changedFiles, "pubspec.yaml") {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "flutter")
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: []string{"flutter", "analyze"}}, true
}

// pubspecDeclaresFlutter matches the `flutter:` key anywhere in the manifest —
// the dependency, the top-level assets section, or an environment constraint.
// Any of them means the package expects the Flutter SDK's analyzer.
func pubspecDeclaresFlutter(repoRoot string) bool {
	return repoFileContains(repoRoot, "pubspec.yaml", "flutter:")
}

// detectDotnetBuild compiles the solution or project; like `cargo check` that
// is not "runs nothing" — MSBuild executes the repo's own targets and analyzers
// during a build, and the run writes bin/ and obj/ next to each project. There
// is no per-file scope, so any C# or project-shape change triggers it.
func detectDotnetBuild(env staticToolEnv) (staticToolCommand, bool) {
	if !dotnetManifestPresent(env) {
		return staticToolCommand{}, false
	}
	if !changedFilesInclude(env.changedFiles, ".cs", ".csproj", ".sln") &&
		!changedFilesIncludeBase(env.changedFiles, "Directory.Build.props") {
		return staticToolCommand{}, false
	}
	bin, ok := lookStaticTool(env.repoRoot, "dotnet")
	if !ok {
		return staticToolCommand{}, false
	}
	return staticToolCommand{bin: bin, argv: []string{"dotnet", "build", "--nologo"}}, true
}

// dotnetManifestPresent looks for a solution or project file. Unlike go.mod or
// Cargo.toml the name is not fixed, so repoFileExists cannot answer this: the
// repo scan's dependency files cover projects anywhere in the tree, and a root
// glob still detects the conventional layout when the caller has no facts.
func dotnetManifestPresent(env staticToolEnv) bool {
	if changedFilesInclude(env.facts.DependencyFiles, dotnetManifestExtensions...) ||
		changedFilesInclude(env.changedFiles, dotnetManifestExtensions...) {
		return true
	}
	for _, pattern := range []string{"*.sln", "*.csproj"} {
		if matches, err := filepath.Glob(filepath.Join(env.repoRoot, pattern)); err == nil && len(matches) > 0 {
			return true
		}
	}
	return false
}

var (
	typeScriptFileExtensions = []string{".ts", ".tsx", ".mts", ".cts"}
	dartFileExtensions       = []string{".dart"}
	dotnetManifestExtensions = []string{".sln", ".csproj"}
	javaScriptFileExtensions = []string{".js", ".jsx", ".mjs", ".cjs", ".ts", ".tsx", ".mts", ".cts"}
	pythonFileExtensions     = []string{".py", ".pyi"}
	nodeLocalBinDirs         = []string{"node_modules/.bin"}
	pythonLocalBinDirs       = []string{".venv/bin", "venv/bin"}
	eslintConfigFiles        = []string{
		"eslint.config.js", "eslint.config.mjs", "eslint.config.cjs",
		"eslint.config.ts", "eslint.config.mts", "eslint.config.cts",
		".eslintrc", ".eslintrc.js", ".eslintrc.cjs", ".eslintrc.mjs",
		".eslintrc.json", ".eslintrc.yml", ".eslintrc.yaml",
	}
)

// lookStaticTool resolves a tool binary, preferring the repo's conventional
// local bin directories over PATH, because a repo's pinned tsc or eslint is
// the one whose result is worth reporting — a newer tsc on PATH reports errors
// the repo's CI never sees.
//
// The directory list is fixed, so the repo cannot redirect the lookup to an
// arbitrary path. That is the only guarantee here: the *executable* found at
// node_modules/.bin/tsc or .venv/bin/ruff is content the checkout supplies, so
// on a checkout the reviewer does not own this is repo-chosen code being
// exec'd. It runs with the scrubbed environment from staticToolChildEnv.
func lookStaticTool(repoRoot, name string, localBinDirs ...string) (string, bool) {
	for _, dir := range localBinDirs {
		path := filepath.Join(repoRoot, filepath.FromSlash(dir), name)
		if !filepath.IsAbs(path) {
			abs, err := filepath.Abs(path)
			if err != nil {
				continue
			}
			path = abs
		}
		info, err := os.Stat(path)
		if err != nil || info.IsDir() || info.Mode().Perm()&0o111 == 0 {
			continue
		}
		return path, true
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}

// staticToolFileArgs turns changed files into safe arguments: paths that
// escape the repo root, live in vendored or generated trees, no longer exist,
// or do not match the requested extensions are dropped, and every argument
// keeps a "./" prefix so a file name can never be read as a flag.
func staticToolFileArgs(repoRoot string, files []string, exts ...string) []string {
	var out []string
	for _, file := range files {
		rel, ok := staticToolRelPath(repoRoot, file)
		if !ok || !hasFileExtension(rel, exts) {
			continue
		}
		out = append(out, "./"+rel)
		if len(out) >= maxStaticToolFileArgs {
			break
		}
	}
	sort.Strings(out)
	return out
}

func staticToolRelPath(repoRoot, file string) (string, bool) {
	file = strings.TrimSpace(file)
	if file == "" || filepath.IsAbs(file) {
		return "", false
	}
	file = filepath.ToSlash(filepath.Clean(filepath.FromSlash(file)))
	file = strings.TrimPrefix(file, "./")
	if file == "." || file == ".." || strings.HasPrefix(file, "../") || strings.HasPrefix(file, "/") {
		return "", false
	}
	if skipDir(file) {
		return "", false
	}
	info, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(file)))
	if err != nil || info.IsDir() {
		return "", false
	}
	return file, true
}

func hasFileExtension(file string, exts []string) bool {
	lower := strings.ToLower(file)
	for _, ext := range exts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func changedFilesInclude(files []string, exts ...string) bool {
	for _, file := range files {
		if hasFileExtension(file, exts) {
			return true
		}
	}
	return false
}

func changedFilesIncludeBase(files []string, bases ...string) bool {
	for _, file := range files {
		for _, base := range bases {
			if filepath.Base(filepath.ToSlash(file)) == base {
				return true
			}
		}
	}
	return false
}

func repoFileExists(repoRoot string, names ...string) bool {
	for _, name := range names {
		info, err := os.Stat(filepath.Join(repoRoot, filepath.FromSlash(name)))
		if err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

func repoFileContains(repoRoot, name, needle string) bool {
	path := filepath.Join(repoRoot, filepath.FromSlash(name))
	info, err := os.Stat(path)
	if err != nil || info.IsDir() || info.Size() > maxStaticToolConfigBytes {
		return false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return bytes.Contains(data, []byte(needle))
}

func runStaticTool(ctx context.Context, repoRoot string, timeout time.Duration, name string, command staticToolCommand) StaticToolResult {
	if strings.TrimSpace(command.bin) == "" || len(command.argv) == 0 {
		return StaticToolResult{Name: name, Skipped: true, Reason: "empty command"}
	}
	toolCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(toolCtx, command.bin, command.argv[1:]...)
	cmd.Dir = repoRoot
	// Not os.Environ(): `go test` here runs the reviewed checkout's own test
	// binaries, so this environment is handed to code from the change under
	// review. See staticToolChildEnv for what survives the allowlist and why.
	cmd.Env = staticToolChildEnv(os.Environ())
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	text := output.String()
	if len(text) > maxStaticToolOutputBytes {
		text = text[:maxStaticToolOutputBytes] + "\n[truncated]\n"
	}
	result := StaticToolResult{
		Name:    name,
		Command: strings.Join(command.argv, " "),
		Output:  strings.TrimSpace(text),
	}
	if err != nil {
		result.ExitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		if toolCtx.Err() == context.DeadlineExceeded {
			result.Skipped = true
			result.Reason = "timed out"
		}
	}
	return result
}

func hasDependencyFile(files []string, base string) bool {
	for _, file := range files {
		if filepath.Base(file) == base {
			return true
		}
	}
	return false
}
