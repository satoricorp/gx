package review

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/capture/matcher"
)

type localInputs struct {
	RepoRoot     string
	Base         string
	Head         string
	ChangedFiles []string
	CommitMsg    string
	HunkLinks    []matcher.HunkLink
	HunkCoverage float64
}

func gradeLocalSections(ctx context.Context, in localInputs) []Section {
	return []Section{
		gradeIntentMatch(in),
		gradeStructure(in),
		gradeBlastRadius(in),
		gradeDoneNess(in),
		gradeBugPass(ctx, in),
	}
}

func gradeIntentMatch(in localInputs) Section {
	section := Section{Name: SectionIntentMatch, Grade: GradePass, Summary: "Intent aligns with changed files."}
	msg := strings.ToLower(strings.TrimSpace(in.CommitMsg))
	if msg == "" {
		section.Grade = GradeWarn
		section.Summary = "No commit message on HEAD; intent match uses file names only."
	}
	matched := 0
	for _, file := range in.ChangedFiles {
		base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		if base != "" && strings.Contains(msg, strings.ToLower(base)) {
			matched++
			section.Evidence = append(section.Evidence, fmt.Sprintf("message mentions `%s`", base))
		}
	}
	if len(in.ChangedFiles) > 0 && matched == 0 && msg != "" {
		section.Grade = GradeWarn
		section.Summary = "Commit message does not mention any changed file stems."
	}
	if len(in.ChangedFiles) == 0 {
		section.Grade = GradeSkip
		section.Summary = "No changed files to compare against intent."
	}
	return section
}

func gradeStructure(in localInputs) Section {
	section := Section{Name: SectionStructure, Grade: GradePass, Summary: "Change structure looks reasonable."}
	dirs := map[string]int{}
	for _, file := range in.ChangedFiles {
		dir := filepath.Dir(file)
		if dir == "." {
			dir = "(root)"
		}
		dirs[dir]++
	}
	if len(dirs) > 5 {
		section.Grade = GradeWarn
		section.Summary = fmt.Sprintf("Cross-cutting edit spans %d directories.", len(dirs))
	}
	for dir, count := range dirs {
		section.Evidence = append(section.Evidence, fmt.Sprintf("%s: %d file(s)", dir, count))
	}
	if len(in.ChangedFiles) == 0 {
		section.Grade = GradeSkip
		section.Summary = "No changed files."
	}
	return section
}

func gradeBlastRadius(in localInputs) Section {
	section := Section{Name: SectionBlastRadius, Grade: GradePass, Summary: "Limited blast radius."}
	files := uniqueFilesFromLinks(in.HunkLinks)
	if len(files) == 0 {
		files = in.ChangedFiles
	}
	if len(files) > 15 {
		section.Grade = GradeWarn
		section.Summary = fmt.Sprintf("Wide fan-out across %d files.", len(files))
	} else if len(files) > 0 {
		section.Summary = fmt.Sprintf("Touches %d file(s).", len(files))
	}
	if in.HunkCoverage > 0 {
		section.Evidence = append(section.Evidence, fmt.Sprintf("matcher coverage: %.0f%%", in.HunkCoverage*100))
	}
	for _, file := range files {
		if isGeneratedPath(file) {
			continue
		}
		if len(section.Evidence) >= 8 {
			section.Evidence = append(section.Evidence, fmt.Sprintf("…and %d more", len(files)-8))
			break
		}
		section.Evidence = append(section.Evidence, file)
	}
	if len(files) == 0 {
		section.Grade = GradeSkip
		section.Summary = "No files in diff."
	}
	return section
}

func gradeDoneNess(in localInputs) Section {
	section := Section{Name: SectionDoneNess, Grade: GradeWarn, Summary: "No test files in diff."}
	testFiles := 0
	todoHits := 0
	for _, file := range in.ChangedFiles {
		if isTestFile(file) {
			testFiles++
			section.Evidence = append(section.Evidence, "test file: "+file)
		}
	}
	if testFiles > 0 {
		section.Grade = GradePass
		section.Summary = fmt.Sprintf("%d test file(s) updated.", testFiles)
	}
	for _, file := range in.ChangedFiles {
		if strings.Contains(strings.ToUpper(file), "TODO") {
			todoHits++
		}
	}
	if todoHits > 0 {
		section.Evidence = append(section.Evidence, fmt.Sprintf("%d path(s) contain TODO", todoHits))
	}
	if len(in.ChangedFiles) == 0 {
		section.Grade = GradeSkip
		section.Summary = "No changed files."
	}
	return section
}

func gradeBugPass(ctx context.Context, in localInputs) Section {
	section := Section{Name: SectionBugPass, Grade: GradePass, Summary: "No obvious defect patterns in diff."}
	if len(in.ChangedFiles) == 0 {
		section.Grade = GradeSkip
		section.Summary = "No diff to scan."
		return section
	}
	diff, err := gitDiff(ctx, in.RepoRoot, in.Base, in.Head)
	if err != nil {
		section.Grade = GradeWarn
		section.Summary = "Could not read diff for bug pass."
		return section
	}
	patterns := []struct {
		label string
		needle string
	}{
		{"panic(nil)", "panic(nil)"},
		{"ignored error", "_ = err"},
		{"ignored error", "_, _ ="},
	}
	hits := 0
	for _, pattern := range patterns {
		if strings.Contains(diff, pattern.needle) {
			hits++
			section.Evidence = append(section.Evidence, "found "+pattern.label)
		}
	}
	if hits > 0 {
		section.Grade = GradeWarn
		section.Summary = fmt.Sprintf("%d obvious defect pattern(s) in diff.", hits)
	}
	return section
}

func uniqueFilesFromLinks(links []matcher.HunkLink) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, link := range links {
		parts := strings.SplitN(link.HunkID, ":", 3)
		if len(parts) < 2 {
			continue
		}
		file := parts[1]
		if _, ok := seen[file]; ok {
			continue
		}
		seen[file] = struct{}{}
		out = append(out, file)
	}
	return out
}

func isGeneratedPath(file string) bool {
	lower := strings.ToLower(file)
	return strings.Contains(lower, "generated") ||
		strings.HasSuffix(lower, ".pb.go") ||
		strings.Contains(lower, "node_modules/")
}

func isTestFile(file string) bool {
	lower := strings.ToLower(file)
	base := filepath.Base(lower)
	return strings.HasSuffix(base, "_test.go") ||
		strings.Contains(lower, "/test/") ||
		strings.Contains(lower, "/tests/") ||
		strings.Contains(base, ".test.") ||
		strings.Contains(base, ".spec.")
}

func gitDiff(ctx context.Context, repoRoot, base, head string) (string, error) {
	rangeSpec := base + ".." + head
	cmd := exec.CommandContext(ctx, "git", "diff", rangeSpec)
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		cmd = exec.CommandContext(ctx, "git", "diff", "--cached")
		cmd.Dir = repoRoot
		out, err = cmd.Output()
		if err != nil {
			cmd = exec.CommandContext(ctx, "git", "diff")
			cmd.Dir = repoRoot
			out, err = cmd.Output()
		}
	}
	return string(out), err
}

func gitChangedFiles(ctx context.Context, repoRoot string) []string {
	cmd := exec.CommandContext(ctx, "git", "status", "--short")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	var files []string
	for _, line := range strings.Split(string(out), "\n") {
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
	return files
}

func gitHeadMessage(ctx context.Context, repoRoot string) string {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%s")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func detectRefRange(ctx context.Context, repoRoot string) (base, head string) {
	head = "HEAD"
	for _, candidate := range []string{"main", "master", "gx/main", "origin/main"} {
		cmd := exec.CommandContext(ctx, "git", "rev-parse", "--verify", candidate+"^{commit}")
		cmd.Dir = repoRoot
		if err := cmd.Run(); err == nil {
			base = candidate
			return base, head
		}
	}
	base = "HEAD~1"
	return base, head
}
