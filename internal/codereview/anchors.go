package codereview

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func validateFindingAnchors(ctx ReviewContext, findings []Finding) []Finding {
	if len(findings) == 0 {
		return findings
	}
	out := make([]Finding, 0, len(findings))
	for _, finding := range findings {
		finding.Anchors = validatedAnchorsForFinding(ctx, finding)
		out = append(out, finding)
	}
	return out
}

func validatedAnchorsForFinding(ctx ReviewContext, finding Finding) []FindingAnchor {
	if len(finding.Anchors) == 0 || strings.TrimSpace(ctx.Brief.RepoRoot) == "" {
		return nil
	}
	changed := changedFileSet(ctx.Brief.Static.ChangedFiles)
	var inHunks []FindingAnchor
	var inChanged []FindingAnchor
	var other []FindingAnchor
	seen := map[string]struct{}{}
	for _, anchor := range finding.Anchors {
		anchor.File = normalizeAnchorFile(anchor.File)
		key := fmt.Sprintf("%s:%d", anchor.File, anchor.Line)
		if anchor.File == "" || anchor.Line <= 0 {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		lineText, ok := anchorLineText(ctx.Brief.RepoRoot, anchor)
		if !ok {
			continue
		}
		if !anchorRelevantToFinding(finding, anchor, lineText) {
			continue
		}
		if anchorMapsToReviewedHunk(context.Background(), ctx.Brief.RepoRoot, ctx.Brief.Static.ReviewRange, anchor) {
			inHunks = append(inHunks, anchor)
			continue
		}
		if _, ok := changed[anchor.File]; ok {
			inChanged = append(inChanged, anchor)
			continue
		}
		other = append(other, anchor)
	}
	switch {
	case len(inHunks) > 0:
		return sortedAnchors(inHunks)
	case len(inChanged) > 0:
		return sortedAnchors(inChanged)
	default:
		return sortedAnchors(other)
	}
}

func normalizeAnchorFile(file string) string {
	file = filepath.ToSlash(strings.TrimSpace(file))
	file = strings.TrimPrefix(file, "./")
	file = strings.TrimPrefix(file, "/")
	if file == "" {
		return ""
	}
	clean := filepath.ToSlash(filepath.Clean(file))
	if clean == "." {
		return ""
	}
	return clean
}

func anchorLineText(repoRoot string, anchor FindingAnchor) (string, bool) {
	if strings.Contains(anchor.File, "..") {
		return "", false
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(anchor.File)))
	if err != nil {
		return "", false
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	line := 0
	for scanner.Scan() {
		line++
		if line == anchor.Line {
			return strings.TrimSpace(scanner.Text()), true
		}
	}
	return "", false
}

func anchorRelevantToFinding(finding Finding, anchor FindingAnchor, lineText string) bool {
	evidence := strings.ToLower(filepath.ToSlash(evidenceText(finding.Evidence)))
	if strings.TrimSpace(evidence) == "" {
		evidence = strings.ToLower(filepath.ToSlash(strings.Join([]string{
			finding.Title,
			finding.Summary,
			finding.Recommendation,
		}, "\n")))
	}
	file := strings.ToLower(filepath.ToSlash(anchor.File))
	if strings.Contains(evidence, file) ||
		strings.Contains(evidence, strings.ToLower(fmt.Sprintf("%s:%d", anchor.File, anchor.Line))) ||
		strings.Contains(evidence, strings.ToLower(filepath.Base(anchor.File))) {
		return true
	}
	return sharesUsefulToken(evidence, strings.ToLower(lineText))
}

func sharesUsefulToken(haystack, text string) bool {
	for _, token := range strings.FieldsFunc(text, func(r rune) bool {
		return !(r == '_' || r == '-' || r == '.' || r == '/' || r == ':' || r >= '0' && r <= '9' || r >= 'a' && r <= 'z')
	}) {
		token = strings.Trim(token, "`\"'()[]{}.,;:")
		if len(token) < 5 {
			continue
		}
		if strings.Contains(haystack, token) {
			return true
		}
	}
	return false
}

func changedFileSet(files []string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, file := range normalizedChangedFiles(files) {
		out[file] = struct{}{}
	}
	return out
}

func sortedAnchors(anchors []FindingAnchor) []FindingAnchor {
	out := append([]FindingAnchor(nil), anchors...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].File != out[j].File {
			return out[i].File < out[j].File
		}
		return out[i].Line < out[j].Line
	})
	return out
}

func AnchorMapsToChangedHunk(ctx context.Context, repoRoot string, anchor FindingAnchor) bool {
	return anchorMapsToReviewedHunk(ctx, repoRoot, "", anchor)
}

// anchorMapsToReviewedHunk checks the anchor against the diff the review
// actually read. refRange is empty for a working-tree review, in which case the
// working-tree diff is tried first and the last commit second.
func anchorMapsToReviewedHunk(ctx context.Context, repoRoot, refRange string, anchor FindingAnchor) bool {
	anchor.File = normalizeAnchorFile(anchor.File)
	if strings.TrimSpace(repoRoot) == "" || anchor.File == "" || anchor.Line <= 0 {
		return false
	}
	var diff string
	if strings.TrimSpace(refRange) != "" {
		diff = rangeFileDiff(ctx, repoRoot, refRange, anchor.File)
	} else {
		diff = fileDiff(ctx, repoRoot, anchor.File)
		if strings.TrimSpace(diff) == "" {
			diff = committedFileDiff(ctx, repoRoot, anchor.File)
		}
	}
	if strings.TrimSpace(diff) == "" {
		return false
	}
	return diffContainsNewLine(diff, anchor.Line)
}

func committedFileDiff(ctx context.Context, repoRoot, file string) string {
	cmd := exec.CommandContext(ctx, "git", "diff", "--no-ext-diff", "HEAD^", "HEAD", "--", file)
	cmd.Dir = repoRoot
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	_ = cmd.Run()
	return out.String()
}

func diffContainsNewLine(diff string, line int) bool {
	newLine := 0
	for _, raw := range strings.Split(diff, "\n") {
		if strings.HasPrefix(raw, "@@") {
			newLine = parseNewHunkStart(raw)
			continue
		}
		if newLine <= 0 || strings.HasPrefix(raw, `\ No newline`) {
			continue
		}
		switch {
		case strings.HasPrefix(raw, "+"):
			if newLine == line {
				return true
			}
			newLine++
		case strings.HasPrefix(raw, "-"):
		default:
			if newLine == line {
				return true
			}
			newLine++
		}
	}
	return false
}

func parseNewHunkStart(header string) int {
	idx := strings.Index(header, " +")
	if idx < 0 {
		return 0
	}
	rest := header[idx+2:]
	end := strings.Index(rest, " ")
	if end >= 0 {
		rest = rest[:end]
	}
	rest = strings.TrimPrefix(rest, "+")
	if comma := strings.Index(rest, ","); comma >= 0 {
		rest = rest[:comma]
	}
	var line int
	_, _ = fmt.Sscanf(rest, "%d", &line)
	return line
}
