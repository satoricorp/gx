package codereview

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/satoricorp/gx/internal/storage"
	"github.com/satoricorp/gx/internal/termstyle"
)

const (
	DefaultScope    = "architecture"
	DefaultFormat   = "markdown"
	reviewMintANSI  = "\x1b[38;2;61;220;151m"
	reviewResetANSI = "\x1b[0m"
)

var supportedScopes = map[string]struct{}{
	"architecture":    {},
	"security":        {},
	"performance":     {},
	"onboarding":      {},
	"docs":            {},
	"dependencies":    {},
	"testing":         {},
	"maintainability": {},
}

var supportedFormats = map[string]struct{}{
	"markdown": {},
	"html":     {},
}

var baselineScopes = []string{"dependencies", "testing", "maintainability"}

type Options struct {
	Scope        string
	Format       string
	Deep         bool
	Since        string
	Focus        string
	Prompt       string
	Verbose      bool
	PatchFocused bool
	ReviewPolicy *ReviewPolicy

	ProgressWriter io.Writer
	Color          bool
}

type Report struct {
	RepoRoot          string
	Scope             string
	Format            string
	Deep              bool
	Since             string
	Focus             string
	Prompt            string
	BaselineScopes    []string
	Docs              []FilePresence
	DependencyFiles   []string
	TestFileCount     int
	TrackedFileCount  int
	ChangedFiles      []string
	ObservationLabels []string
	Findings          []Finding
	Sources           []Source
	SourceRefs        []SourceRef
	Reviewer          string
	ContextSnippets   int
	Verbose           bool
	Color             bool
}

type FilePresence struct {
	Path    string
	Present bool
}

func ValidateOptions(opts Options) error {
	opts = normalizeOptions(opts)
	if _, ok := supportedScopes[opts.Scope]; !ok {
		return fmt.Errorf("unsupported review scope %q", opts.Scope)
	}
	if _, ok := supportedFormats[opts.Format]; !ok {
		return fmt.Errorf("unsupported review format %q", opts.Format)
	}
	return nil
}

func RenderMarkdown(report Report) string {
	var b strings.Builder
	fmt.Fprintln(&b, reviewTitle(report, "## Recommendations"))
	if len(report.Findings) == 0 {
		fmt.Fprintln(&b, "- No material issues found in this change.")
	} else {
		for index, finding := range report.Findings {
			fmt.Fprintf(&b, "%s\n", reviewTitle(report, fmt.Sprintf("### %d. %s", index+1, finding.Title)))
			fmt.Fprintln(&b)
			fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Why:**"), finding.Summary)
			if strings.TrimSpace(finding.Benefit) != "" {
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Benefit:**"), finding.Benefit)
			}
			fmt.Fprintln(&b)
			fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Do next:**"), finding.Recommendation)
			if attributions := renderFindingAttributions(report, finding); len(attributions) > 0 {
				fmt.Fprintln(&b)
				fmt.Fprintf(&b, "%s %s\n", reviewLabel(report, "**Attribution:**"), strings.Join(attributions, " · "))
			}
			if report.Verbose && len(finding.Evidence) > 0 {
				fmt.Fprintln(&b)
				fmt.Fprintln(&b, reviewLabel(report, "**Evidence:**"))
				for _, evidence := range finding.Evidence {
					fmt.Fprintf(&b, "- %s: %s\n", evidence.Label, evidence.Value)
				}
			}
			fmt.Fprintln(&b)
		}
	}
	fmt.Fprintln(&b)

	if len(report.Sources) > 0 {
		fmt.Fprintln(&b, reviewTitle(report, "## Sources"))
		for _, source := range report.Sources {
			fmt.Fprintf(&b, "- %s\n", renderSourceCatalogEntry(source))
		}
		fmt.Fprintln(&b)
	}

	if len(report.SourceRefs) > 0 {
		fmt.Fprintln(&b, reviewTitle(report, "## Context Sources"))
		for _, sourceRef := range report.SourceRefs {
			fmt.Fprintf(&b, "- %s\n", renderSourceRef(sourceRef))
		}
		fmt.Fprintln(&b)
	}

	if report.Verbose {
		fmt.Fprintln(&b, reviewTitle(report, "## Repo Facts"))
		fmt.Fprintf(&b, "- Tracked/source files scanned: `%d`\n", report.TrackedFileCount)
		fmt.Fprintf(&b, "- Test files: `%d`\n", report.TestFileCount)
		fmt.Fprintf(&b, "- Context snippets: `%d`\n", report.ContextSnippets)
		if len(report.DependencyFiles) == 0 {
			fmt.Fprintln(&b, "- Dependency manifests: none detected")
		} else {
			fmt.Fprintf(&b, "- Dependency manifests: `%s`\n", strings.Join(report.DependencyFiles, "`, `"))
		}
		if report.Since != "" {
			fmt.Fprintf(&b, "- Since: `%s`\n", report.Since)
		} else {
			fmt.Fprintln(&b, "- Since: `forever`")
		}
		fmt.Fprintln(&b)

		fmt.Fprintln(&b, reviewTitle(report, "## Docs"))
		for _, doc := range report.Docs {
			status := "missing"
			if doc.Present {
				status = "present"
			}
			fmt.Fprintf(&b, "- `%s`: %s\n", doc.Path, status)
		}
		fmt.Fprintln(&b)

		fmt.Fprintln(&b, reviewTitle(report, "## Changed Files"))
		if len(report.ChangedFiles) == 0 {
			fmt.Fprintln(&b, "- none detected")
		} else {
			for _, file := range report.ChangedFiles {
				fmt.Fprintf(&b, "- `%s`\n", file)
			}
		}
		fmt.Fprintln(&b)
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

func renderFindingAttributions(report Report, finding Finding) []string {
	if len(finding.SourceIDs) == 0 {
		return nil
	}
	sources := sourceMap(report.Sources)
	var out []string
	for _, id := range finding.SourceIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if source, ok := sources[id]; ok {
			out = append(out, renderSourceCatalogEntry(source))
			continue
		}
		out = append(out, fmt.Sprintf("`%s`", id))
	}
	return out
}

func sourceMap(sources []Source) map[string]Source {
	out := make(map[string]Source, len(sources))
	for _, source := range sources {
		if strings.TrimSpace(source.ID) == "" {
			continue
		}
		out[source.ID] = source
	}
	return out
}

func renderSourceCatalogEntry(source Source) string {
	title := strings.TrimSpace(source.Title)
	if title == "" {
		title = strings.TrimSpace(source.ID)
	}
	if title == "" {
		return "unknown source"
	}
	if url := strings.TrimSpace(source.URL); url != "" {
		return fmt.Sprintf("[%s](%s)", title, url)
	}
	if id := strings.TrimSpace(source.ID); id != "" {
		return fmt.Sprintf("%s (`%s`)", title, id)
	}
	return title
}

func renderSourceRef(ref SourceRef) string {
	label := strings.TrimSpace(ref.ID)
	if label == "" {
		label = strings.TrimSpace(ref.Kind)
	}
	if label == "" {
		label = "context"
	}
	title := strings.TrimSpace(ref.Title)
	if title == "" {
		title = strings.TrimSpace(ref.File)
	}
	if title == "" {
		title = label
	}
	var details []string
	if url := strings.TrimSpace(ref.URL); url != "" {
		details = append(details, fmt.Sprintf("[%s](%s)", title, url))
	} else if ref.File != "" {
		if ref.StartLine > 0 {
			details = append(details, fmt.Sprintf("`%s:%d`", ref.File, ref.StartLine))
		} else {
			details = append(details, fmt.Sprintf("`%s`", ref.File))
		}
	} else {
		details = append(details, fmt.Sprintf("`%s`", title))
	}
	if ref.Source != "" && ref.Source != "local" {
		details = append(details, fmt.Sprintf("source=%s", ref.Source))
	}
	return fmt.Sprintf("`%s` %s", label, strings.Join(details, " · "))
}

func reviewTitle(report Report, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return reviewMintANSI + text + reviewResetANSI
}

func reviewLabel(report Report, text string) string {
	if !report.Color || text == "" || !termstyle.Enabled() {
		return text
	}
	return reviewMintANSI + text + reviewResetANSI
}

func normalizeOptions(opts Options) Options {
	opts.Prompt = strings.TrimSpace(opts.Prompt)
	opts.Focus = strings.TrimSpace(opts.Focus)
	scope := strings.ToLower(strings.TrimSpace(opts.Scope))
	if scope == "" {
		opts.Scope = DefaultScope
		if !opts.Deep && opts.Prompt == "" {
			opts.PatchFocused = true
		}
	} else {
		opts.Scope = scope
	}
	opts.Format = strings.ToLower(strings.TrimSpace(opts.Format))
	if opts.Format == "" {
		opts.Format = DefaultFormat
	}
	return opts
}

func baselineFor(scope string) []string {
	var out []string
	for _, baseline := range baselineScopes {
		if baseline != scope {
			out = append(out, baseline)
		}
	}
	return out
}

func depthLabel(deep bool) string {
	if deep {
		return "deep"
	}
	return "shallow"
}

func reviewProfile(opts Options) string {
	if opts.Deep {
		return "deep_full_spectrum"
	}
	if opts.PatchFocused {
		return "patch_focused"
	}
	if strings.TrimSpace(opts.Prompt) != "" {
		return "prompt_directed"
	}
	return "scope_focused"
}

type RepoFacts struct {
	Docs             []FilePresence
	ADRFiles         []string
	DependencyFiles  []string
	Files            []string
	GoPackages       []PackageFact
	TestFileCount    int
	TrackedFileCount int
}

type PackageFact struct {
	Path       string
	GoFiles    int
	TestFiles  int
	PublicName bool
}

type LocalScanner struct{}

func (LocalScanner) Scan(_ context.Context, repoRoot string, focus string) (RepoFacts, error) {
	return scanRepo(repoRoot, focus)
}

func scanRepo(repoRoot, focus string) (RepoFacts, error) {
	docs := []FilePresence{
		{Path: "README.md", Present: exists(repoRoot, "README.md")},
		{Path: "AGENTS.md", Present: exists(repoRoot, "AGENTS.md")},
		{Path: "CONTEXT.md", Present: exists(repoRoot, "CONTEXT.md")},
		{Path: "docs/", Present: exists(repoRoot, "docs")},
	}
	facts := RepoFacts{Docs: docs}
	if files, ok := gitTrackedFiles(repoRoot); ok {
		for _, rel := range files {
			if focus != "" && !inFocus(rel, focus) {
				continue
			}
			if skipFile(rel) {
				continue
			}
			facts.addFile(rel)
		}
		sort.Strings(facts.DependencyFiles)
		facts.ADRFiles = adrFiles(facts.Files)
		facts.finalize()
		return facts, nil
	}
	err := filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if entry.IsDir() {
			if skipDir(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if focus != "" && !inFocus(rel, focus) {
			return nil
		}
		if skipFile(rel) {
			return nil
		}
		facts.addFile(rel)
		return nil
	})
	sort.Strings(facts.DependencyFiles)
	facts.ADRFiles = adrFiles(facts.Files)
	facts.finalize()
	return facts, err
}

func (f *RepoFacts) addFile(rel string) {
	f.TrackedFileCount++
	f.Files = append(f.Files, rel)
	if isDependencyFile(rel) {
		f.DependencyFiles = append(f.DependencyFiles, rel)
	}
	if isTestFile(rel) {
		f.TestFileCount++
	}
}

func (f *RepoFacts) finalize() {
	sort.Strings(f.Files)
	f.GoPackages = goPackages(f.Files)
}

func goPackages(files []string) []PackageFact {
	byDir := map[string]PackageFact{}
	for _, file := range files {
		if !strings.HasSuffix(file, ".go") {
			continue
		}
		dir := filepath.ToSlash(filepath.Dir(file))
		if dir == "." {
			dir = "."
		}
		fact := byDir[dir]
		fact.Path = dir
		if isTestFile(file) {
			fact.TestFiles++
		} else {
			fact.GoFiles++
		}
		base := filepath.Base(dir)
		fact.PublicName = !strings.HasPrefix(base, "internal") && dir != "."
		byDir[dir] = fact
	}
	var out []PackageFact
	for _, fact := range byDir {
		out = append(out, fact)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})
	return out
}

func gitTrackedFiles(repoRoot string) ([]string, bool) {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = repoRoot
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, false
	}
	var files []string
	for _, line := range strings.Split(out.String(), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, filepath.ToSlash(line))
		}
	}
	return files, true
}

func inFocus(rel, focus string) bool {
	focus = strings.Trim(strings.TrimSpace(filepath.ToSlash(focus)), "/")
	return focus == "" || rel == focus || strings.HasPrefix(rel, focus+"/")
}

func (f RepoFacts) observations() []string {
	var out []string
	if !present(f.Docs, "README.md") {
		out = append(out, "`README.md` is missing.")
	}
	if !present(f.Docs, "AGENTS.md") {
		out = append(out, "`AGENTS.md` is missing.")
	}
	if !present(f.Docs, "CONTEXT.md") {
		out = append(out, "`CONTEXT.md` is missing.")
	}
	if len(f.DependencyFiles) == 0 {
		out = append(out, "No dependency manifest was detected.")
	}
	if f.TestFileCount == 0 {
		out = append(out, "No test files were detected.")
	}
	return out
}

func present(docs []FilePresence, path string) bool {
	for _, doc := range docs {
		if doc.Path == path {
			return doc.Present
		}
	}
	return false
}

func adrFiles(files []string) []string {
	var out []string
	for _, file := range files {
		lower := strings.ToLower(file)
		base := filepath.Base(lower)
		if strings.Contains(lower, "/adr/") || strings.HasPrefix(base, "adr-") || strings.HasPrefix(base, "adr_") {
			out = append(out, file)
		}
	}
	sort.Strings(out)
	return out
}

func exists(root, rel string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
	return err == nil
}

func skipDir(rel string) bool {
	for _, part := range strings.Split(rel, "/") {
		switch part {
		case ".git", ".jj", ".gx", ".gocache", "node_modules", "dist", "build", ".next", "coverage", ".cache", ".turbo", "vendor":
			return true
		}
	}
	return false
}

func skipFile(rel string) bool {
	lower := strings.ToLower(rel)
	for _, suffix := range []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".svg", ".woff", ".woff2", ".ttf", ".eot", ".mp4", ".mp3", ".zip", ".tar", ".gz", ".pdf", ".exe", ".dll", ".so", ".dylib"} {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

func isDependencyFile(rel string) bool {
	switch filepath.Base(rel) {
	case "go.mod", "go.sum", "package.json", "package-lock.json", "bun.lock", "bun.lockb", "pnpm-lock.yaml", "yarn.lock", "Cargo.toml", "Cargo.lock", "requirements.txt", "pyproject.toml", "poetry.lock":
		return true
	default:
		return false
	}
}

func isTestFile(rel string) bool {
	lower := strings.ToLower(rel)
	base := filepath.Base(lower)
	return strings.HasSuffix(base, "_test.go") ||
		strings.Contains(lower, "/test/") ||
		strings.Contains(lower, "/tests/") ||
		strings.Contains(base, ".test.") ||
		strings.Contains(base, ".spec.")
}

func changedFiles(ctx context.Context, repoRoot string) []string {
	cmd := exec.CommandContext(ctx, "git", "status", "--short")
	cmd.Dir = repoRoot
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil
	}
	var files []string
	for _, line := range strings.Split(out.String(), "\n") {
		if len(line) < 4 {
			continue
		}
		file := strings.TrimSpace(line[3:])
		if strings.Contains(file, " -> ") {
			parts := strings.Split(file, " -> ")
			file = parts[len(parts)-1]
		}
		if file != "" {
			files = append(files, file)
		}
	}
	sort.Strings(files)
	return files
}

func reviewChangedFiles(ctx context.Context, repoRoot string) []string {
	files := changedFiles(ctx, repoRoot)
	if len(files) > 0 {
		return files
	}
	return gxRevisionChangedFiles(ctx, repoRoot)
}

func gxRevisionChangedFiles(ctx context.Context, repoRoot string) []string {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return nil
	}
	db, err := storage.Open(ctx)
	if err != nil {
		return nil
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, `
		WITH latest_revisions AS (
			SELECT cr.change_id, cr.changed_files_json
			FROM change_revisions cr
			JOIN (
				SELECT change_id, MAX(created_at) AS created_at
				FROM change_revisions
				GROUP BY change_id
			) latest
				ON latest.change_id = cr.change_id AND latest.created_at = cr.created_at
		)
		SELECT lr.changed_files_json
		FROM repos r
		JOIN stacks s ON s.repo_id = r.id
		JOIN stack_changes sc ON sc.stack_id = s.id
		JOIN changes c ON c.id = sc.change_id
		JOIN latest_revisions lr ON lr.change_id = c.id
		WHERE r.root_path = ? AND COALESCE(s.status, '') != 'merged'
		ORDER BY sc.position ASC, c.updated_at DESC
	`, repoRoot)
	if err != nil {
		return nil
	}
	defer rows.Close()

	seen := map[string]struct{}{}
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			return nil
		}
		var files []string
		if err := json.Unmarshal([]byte(raw), &files); err != nil {
			continue
		}
		for _, file := range files {
			file = strings.TrimSpace(filepath.ToSlash(file))
			if file == "" {
				continue
			}
			seen[file] = struct{}{}
		}
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	files := make([]string, 0, len(seen))
	for file := range seen {
		files = append(files, file)
	}
	sort.Strings(files)
	return files
}
