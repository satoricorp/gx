package codereview

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Diff signals the review engine reads: which lines a change added, the hunk
// around a line, and the code a finding points at.
//
// These lived in the pre-ship gate's engine, which is why they carried a
// `gates` prefix — but the review engine was already the heavier user of them
// (story.go for a story item's hunk, render_code.go for the fence, the secret
// scan for added lines), so removing the gate command left them here rather
// than deleting them with it.

// addedLine is one "+" line of a diff with its post-change line
// number, so a finding about it can anchor to file:line.
type addedLine struct {
	Number int
	Text   string
}

// addedLinesForDiff returns the added lines of one snippet. A
// content-fallback snippet (an untracked file — the commonest agent output)
// has no hunks: every line of it is new, and skipping it would exempt exactly
// the files most worth scanning.
func addedLinesForDiff(diff string) []addedLine {
	if strings.HasPrefix(diff, diffUnavailableContentHeader) {
		content := strings.TrimPrefix(strings.TrimPrefix(diff, diffUnavailableContentHeader), "\n")
		if strings.TrimSpace(content) == "" {
			return nil
		}
		lines := strings.Split(content, "\n")
		out := make([]addedLine, 0, len(lines))
		for i, line := range lines {
			out = append(out, addedLine{Number: i + 1, Text: line})
		}
		return out
	}
	return parseAddedLines(diff)
}

// diffHunkContext is how many hunk lines surround the discussed line
// when a finding shows its diff window.
const (
	diffHunkContext = 2
	// excerptContext / excerptMaxLineBytes bound the source window shown for a
	// finding whose line the change did not touch.
	excerptContext      = 2
	excerptMaxLineBytes = 200
)

// diffHunkForLine returns the unified-diff window around a
// post-change line, hunk header included, or "" when the line is not part of
// this diff — a finding about untouched code has no hunk to show.

// diffHunkContext is how many hunk lines surround the discussed line
// when a finding shows its diff window.

// diffHunkForLine returns the unified-diff window around a
// post-change line, hunk header included, or "" when the line is not part of
// this diff — a finding about untouched code has no hunk to show.

// diffHunkForLine returns the unified-diff window around a
// post-change line, hunk header included, or "" when the line is not part of
// this diff — a finding about untouched code has no hunk to show.
func diffHunkForLine(diff string, target int) string {
	if target <= 0 || strings.HasPrefix(diff, diffUnavailableContentHeader) {
		return ""
	}
	type diffHunk struct {
		header string
		lines  []string
		// nums holds each line's post-change line number; 0 for removed lines,
		// which exist only on the pre-change side.
		nums []int
	}
	var hunks []diffHunk
	newLine := 0
	inHunk := false
	for _, raw := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(raw, "@@"):
			newLine = parseHunkNewStart(raw)
			inHunk = newLine > 0
			if inHunk {
				hunks = append(hunks, diffHunk{header: raw})
			}
		case !inHunk:
		case strings.HasPrefix(raw, "-"):
			hunk := &hunks[len(hunks)-1]
			hunk.lines = append(hunk.lines, raw)
			hunk.nums = append(hunk.nums, 0)
		default:
			hunk := &hunks[len(hunks)-1]
			hunk.lines = append(hunk.lines, raw)
			hunk.nums = append(hunk.nums, newLine)
			newLine++
		}
	}
	for _, hunk := range hunks {
		targetIndex := -1
		for i, num := range hunk.nums {
			if num == target {
				targetIndex = i
				break
			}
		}
		if targetIndex < 0 {
			continue
		}
		start := targetIndex - diffHunkContext
		if start < 0 {
			start = 0
		}
		end := targetIndex + diffHunkContext + 1
		if end > len(hunk.lines) {
			end = len(hunk.lines)
		}
		out := []string{hunk.header}
		if start > 0 {
			out = append(out, "…")
		}
		out = append(out, hunk.lines[start:end]...)
		if end < len(hunk.lines) {
			out = append(out, "…")
		}
		return strings.Join(out, "\n")
	}
	return ""
}

// parseAddedLines walks a unified diff and returns the added lines
// with their post-change line numbers.

// parseAddedLines walks a unified diff and returns the added lines
// with their post-change line numbers.
func parseAddedLines(diff string) []addedLine {
	var out []addedLine
	line := 0
	inHunk := false
	for _, raw := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(raw, "@@"):
			line = parseHunkNewStart(raw)
			inHunk = line > 0
			// The hunk header names the first line of the section; the counter
			// advances as lines are consumed below.
		case !inHunk:
			continue
		case strings.HasPrefix(raw, "+++"), strings.HasPrefix(raw, "---"):
			continue
		case strings.HasPrefix(raw, "+"):
			out = append(out, addedLine{Number: line, Text: strings.TrimPrefix(raw, "+")})
			line++
		case strings.HasPrefix(raw, "-"):
			// Removed lines do not advance the post-change counter.
		default:
			line++
		}
	}
	return out
}

// parseHunkNewStart reads the "+c[,d]" of "@@ -a,b +c,d @@", or 0.

// parseHunkNewStart reads the "+c[,d]" of "@@ -a,b +c,d @@", or 0.
func parseHunkNewStart(header string) int {
	index := strings.Index(header, "+")
	if index < 0 {
		return 0
	}
	rest := header[index+1:]
	end := strings.IndexAny(rest, ", @")
	if end < 0 {
		end = len(rest)
	}
	value, err := strconv.Atoi(strings.TrimSpace(rest[:end]))
	if err != nil || value < 0 {
		return 0
	}
	if value == 0 {
		// "+0,0" is an empty new side; there is no line 0.
		return 0
	}
	return value
}

// attachDiffHunks gives each finding whose line is part of the
// change its unified-diff window, extracted from the same snippets the gates
// read.
func attachDiffHunks(diffsByFile map[string]string, findings []Finding) {
	for i := range findings {
		finding := &findings[i]
		if finding.DiffHunk != "" || finding.File == "" || finding.Line <= 0 {
			continue
		}
		diff, ok := diffsByFile[finding.File]
		if !ok {
			continue
		}
		finding.DiffHunk = diffHunkForLine(diff, finding.Line)
	}
}

// attachCodeExcerpts reads the source lines each finding points at,
// so the report shows the code being discussed rather than only naming it.
// Findings that already carry a diff hunk are left alone — the hunk is the
// better evidence. Best-effort by design: an unreadable file or a stale line
// number just leaves the excerpt empty.

// attachCodeExcerpts reads the source lines each finding points at,
// so the report shows the code being discussed rather than only naming it.
// Findings that already carry a diff hunk are left alone — the hunk is the
// better evidence. Best-effort by design: an unreadable file or a stale line
// number just leaves the excerpt empty.
func attachCodeExcerpts(repoRoot string, findings []Finding) {
	for i := range findings {
		finding := &findings[i]
		if finding.DiffHunk != "" || finding.CodeExcerpt != "" || finding.File == "" || finding.Line <= 0 {
			continue
		}
		rel, ok := staticToolRelPath(repoRoot, finding.File)
		if !ok {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		if finding.Line > len(lines) {
			continue
		}
		start := finding.Line - excerptContext
		if start < 1 {
			start = 1
		}
		end := finding.Line + excerptContext
		if end > len(lines) {
			end = len(lines)
		}
		excerpt := make([]string, 0, end-start+1)
		for _, line := range lines[start-1 : end] {
			if len(line) > excerptMaxLineBytes {
				line = line[:excerptMaxLineBytes] + "…"
			}
			excerpt = append(excerpt, line)
		}
		finding.CodeExcerpt = strings.Join(excerpt, "\n")
		finding.CodeExcerptStart = start
	}
}

// gatesVerdict rolls the gates up. Real failures outrank everything;
// a degraded run with no failures must not claim ship.

func staticToolsDisabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("GX_REVIEW_STATIC_TOOLS")), "0")
}
