package review

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/satoricorp/gx/internal/authoring"
	"github.com/satoricorp/gx/internal/capture/matcher"
)

// Review runs graded local review with optional server context.
func Review(ctx context.Context, repoRoot string) (GradedReport, error) {
	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return GradedReport{}, err
	}
	base, head := detectRefRange(ctx, repoRoot)
	changed := gitChangedFiles(ctx, repoRoot)
	captureResult, _ := matchCaptureForReview(ctx, repoRoot, base, head)

	in := localInputs{
		RepoRoot:     repoRoot,
		Base:         base,
		Head:         head,
		ChangedFiles: changed,
		CommitMsg:    gitHeadMessage(ctx, repoRoot),
		HunkLinks:    captureResult.HunkLinks,
		HunkCoverage: captureResult.HunkCoverage,
	}

	sections := gradeLocalSections(ctx, in)
	report := GradedReport{
		Repo:     repoRoot,
		RefRange: base + ".." + head,
		Sections: sections,
	}

	client, online := newRemoteClient()
	if !online {
		return applyOfflineDegrade(report, fmt.Errorf("GX API URL and token not configured")), nil
	}

	remote, err := client.fetchContext(ctx, repoRoot, base, head)
	if err != nil {
		return applyOfflineDegrade(report, err), nil
	}

	report.Sections = append(report.Sections,
		gradeRulesSection(remote, false),
		gradeCollisionsSection(remote, false),
	)
	return orderSections(report), nil
}

func applyOfflineDegrade(report GradedReport, reason error) GradedReport {
	report.Offline = true
	report.Sections = append(report.Sections,
		gradeRulesSection(RemoteContext{}, true),
		gradeCollisionsSection(RemoteContext{}, true),
	)
	if reason != nil {
		report.Sections = append(report.Sections, Section{
			Name:    "server",
			Grade:   GradeSkip,
			Summary: fmt.Sprintf("Server context unavailable: %v", reason),
		})
	}
	return orderSections(report)
}

func orderSections(report GradedReport) GradedReport {
	order := []string{
		SectionIntentMatch,
		SectionStructure,
		SectionRules,
		SectionBlastRadius,
		SectionCollisions,
		SectionDoneNess,
		SectionBugPass,
	}
	byName := map[string]Section{}
	var extras []Section
	for _, section := range report.Sections {
		if section.Name == "server" {
			extras = append(extras, section)
			continue
		}
		byName[section.Name] = section
	}
	var ordered []Section
	for _, name := range order {
		if section, ok := byName[name]; ok {
			ordered = append(ordered, section)
		}
	}
	report.Sections = append(ordered, extras...)
	return report
}

type reviewCaptureResult struct {
	HunkLinks    []matcher.HunkLink
	HunkCoverage float64
}

func matchCaptureForReview(ctx context.Context, repoRoot, base, head string) (reviewCaptureResult, error) {
	hunks, err := collectDiffHunks(ctx, repoRoot, base, head)
	if err != nil || len(hunks) == 0 {
		return reviewCaptureResult{}, err
	}
	result, err := authoring.MatchWorkingCopy(ctx, repoRoot, hunks, nil)
	if err != nil {
		return reviewCaptureResult{}, err
	}
	return reviewCaptureResult{
		HunkLinks:    result.HunkLinks,
		HunkCoverage: result.HunkCoverage,
	}, nil
}

func collectDiffHunks(ctx context.Context, repoRoot, base, head string) ([]authoring.HunkRange, error) {
	diff, err := gitDiff(ctx, repoRoot, base, head)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(diff) == "" {
		diff, _ = gitDiffWorkingTree(ctx, repoRoot)
	}
	return parseDiffHunks(diff), nil
}

func gitDiffWorkingTree(ctx context.Context, repoRoot string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "diff", "HEAD")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	return string(out), err
}

func parseDiffHunks(diff string) []authoring.HunkRange {
	lines := strings.Split(diff, "\n")
	var hunks []authoring.HunkRange
	var currentFile string
	index := 1
	for _, line := range lines {
		if strings.HasPrefix(line, "+++ b/") {
			currentFile = strings.TrimPrefix(line, "+++ b/")
			continue
		}
		if strings.HasPrefix(line, "@@") && currentFile != "" {
			id := fmt.Sprintf("h%d", index)
			hunks = append(hunks, authoring.HunkRange{ID: id, File: currentFile, Patch: line + "\n"})
			index++
		}
	}
	return hunks
}
