package codereview

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/satoricorp/gx/internal/termstyle"
)

// Shared code-window renderers for any report. The review renderer and the
// gates renderer both use them to show the code a finding is about: the
// diff hunk when the finding is about the change, otherwise a syntax-highlighted
// source excerpt with a line-number gutter. Colors ride on a plain bool and
// termstyle.Enabled(), so --json/--md/piped output stays plain.

func colorize(color bool, paint func(string) string, text string) string {
	if !color || text == "" || !termstyle.Enabled() {
		return text
	}
	return paint(text)
}

// highlightCodeForTerminal syntax-highlights an excerpt for the terminal.
// Plain text comes back on any failure — highlighting is presentation, never a
// reason to lose the code.
func highlightCodeForTerminal(file, code string) string {
	lexer := lexers.Match(filepath.Base(file))
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		return code
	}
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}
	var b strings.Builder
	if err := formatter.Format(&b, style, iterator); err != nil {
		return code
	}
	return strings.TrimRight(b.String(), "\n")
}

// renderCodeLines renders whatever code a finding carries for the terminal:
// the diff hunk when the finding is about the change, otherwise the source
// excerpt.
func renderCodeLines(color bool, finding Finding, indent string) []string {
	if strings.TrimSpace(finding.DiffHunk) != "" {
		return renderHunkWindow(color, finding.File, finding.DiffHunk, finding.Line, indent)
	}
	return renderExcerptLines(color, finding, indent)
}

// renderHunkLines renders a unified-diff window with the coloring a terminal
// reader expects: additions green, removals red, the hunk header and trims
// muted.
func renderHunkLines(color bool, hunk, indent string) []string {
	return renderHunkLinesFor(color, "", hunk, indent)
}

// renderHunkLinesFor is renderHunkLines with a file name for syntax
// highlighting. The sign column (+, -, or space) is painted green, red, or
// muted; the code after it is highlighted with the file's lexer when color is
// on, exactly as an excerpt would be, so a diff and an excerpt of the same
// code look the same. Hunk headers (@@ … @@) are dropped: they are addressing
// for tools, and the finding already carries file:line. Tabs are expanded so
// a terminal never shows a raw control character in place of indentation.
func renderHunkLinesFor(color bool, file, hunk, indent string) []string {
	var out []string
	for _, line := range strings.Split(hunk, "\n") {
		if strings.HasPrefix(line, "@@") {
			continue
		}
		if line == "…" {
			out = append(out, indent+colorize(color, termstyle.Muted, line))
			continue
		}
		if line == "" {
			out = append(out, indent)
			continue
		}
		// A leading NUL is renderHunkWindow's "this is the discussed line"
		// sentinel; it is not part of the diff. Strip it, remember it, and
		// swap the last two indent columns for the → gutter, so the code
		// column stays aligned with the unmarked lines around it.
		marked := strings.HasPrefix(line, "\x00")
		line = strings.TrimPrefix(line, "\x00")
		sign, body := line[:1], clipCodeLine(strings.ReplaceAll(line[1:], "\t", "    "), renderCodeWidth-len(indent)-1)
		if color && file != "" {
			body = highlightCodeForTerminal(file, body)
		}
		var painted string
		switch sign {
		case "+":
			painted = colorize(color, termstyle.Success, "+") + body
		case "-":
			painted = colorize(color, termstyle.Danger, "-") + body
		default:
			painted = " " + body
		}
		lead := indent
		if marked && len(indent) >= 2 {
			lead = indent[:len(indent)-2] + colorize(color, termstyle.Danger, "→ ")
		}
		out = append(out, lead+painted)
	}
	return out
}

// renderExcerptLines renders a finding's excerpt with a line-number gutter and
// a marker on the discussed line, syntax-highlighted when the report is in
// color.
func renderExcerptLines(color bool, finding Finding, indent string) []string {
	if strings.TrimSpace(finding.CodeExcerpt) == "" {
		return nil
	}
	// Clip and expand tabs per line BEFORE highlighting, so a cut never lands
	// inside an escape sequence and the lexer sees plain text.
	gutterWidth := len(indent) + len("  0000 | ")
	var clean []string
	for _, line := range strings.Split(finding.CodeExcerpt, "\n") {
		clean = append(clean, clipCodeLine(strings.ReplaceAll(line, "\t", "    "), renderCodeWidth-gutterWidth))
	}
	code := strings.Join(clean, "\n")
	if color && termstyle.Enabled() {
		code = highlightCodeForTerminal(finding.File, code)
	}
	var out []string
	for i, line := range strings.Split(code, "\n") {
		number := finding.CodeExcerptStart + i
		gutter := colorize(color, termstyle.Muted, fmt.Sprintf("  %4d | ", number))
		if number == finding.Line {
			gutter = colorize(color, termstyle.Danger, "→ ") +
				colorize(color, termstyle.Muted, fmt.Sprintf("%4d | ", number))
		}
		out = append(out, indent+gutter+line)
	}
	return out
}

// fenceLanguage tags a markdown code fence so renderers highlight the
// excerpt. Empty is fine: an untagged fence still renders as code.
func fenceLanguage(file string) string {
	if language := qualityLanguage(file); language != "" {
		return language
	}
	switch strings.ToLower(filepath.Ext(file)) {
	case ".html":
		return "html"
	case ".vue":
		return "vue"
	case ".svelte":
		return "svelte"
	case ".sh", ".bash":
		return "bash"
	case ".yml", ".yaml":
		return "yaml"
	case ".json":
		return "json"
	case ".toml":
		return "toml"
	case ".mod":
		return "go"
	default:
		return ""
	}
}

// renderCodeWidth is the column code lines are clipped at. A code line is
// never wrapped — a wrapped line of code reads as two lines of code — so a
// long one is cut and marked with an ellipsis. Prose has its own wrapping in
// the review renderer.
const renderCodeWidth = 100

// clipCodeLine truncates a code line to width runes with a trailing ellipsis.
// It runs before highlighting so the cut never lands inside an escape code.
func clipCodeLine(line string, width int) string {
	if width < 8 {
		width = 8
	}
	runes := []rune(line)
	if len(runes) <= width {
		return line
	}
	return string(runes[:width-1]) + "…"
}

// renderHunkWindow is renderHunkLinesFor with the discussed line marked. A
// hunk window from diffHunkForLine carries its @@ header, which
// names the post-change start line; from that and the +/space lines the
// post-change number of every line is known, so the one the finding is about
// gets the same → an excerpt gives it. Without the marker a ten-line window
// of closing braces gave the eye nowhere to land. When the header is absent
// (a model's example diff) no line is marked.
func renderHunkWindow(color bool, file, hunk string, target int, indent string) []string {
	lines := strings.Split(hunk, "\n")
	newLine := 0
	haveNumbers := false
	for _, raw := range lines {
		if strings.HasPrefix(raw, "@@") {
			if n := parseHunkNewStart(raw); n > 0 {
				newLine, haveNumbers = n, true
			}
			break
		}
	}
	// Recompute the post-change number as we walk, marking the target.
	var marked []string
	num := newLine
	for _, raw := range lines {
		switch {
		case strings.HasPrefix(raw, "@@"), raw == "…", raw == "":
			marked = append(marked, raw)
		case strings.HasPrefix(raw, "-"):
			marked = append(marked, raw)
		default:
			if haveNumbers && target > 0 && num == target {
				marked = append(marked, "\x00"+raw) // sentinel: mark this line
			} else {
				marked = append(marked, raw)
			}
			num++
		}
	}
	return renderHunkLinesFor(color, file, strings.Join(marked, "\n"), indent)
}
