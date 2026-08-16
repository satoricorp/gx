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
// constraints renderer both use them to show the code a finding is about: the
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
		return renderHunkLines(color, finding.DiffHunk, indent)
	}
	return renderExcerptLines(color, finding, indent)
}

// renderHunkLines renders a unified-diff window with the coloring a terminal
// reader expects: additions green, removals red, the hunk header and trims
// muted.
func renderHunkLines(color bool, hunk, indent string) []string {
	var out []string
	for _, line := range strings.Split(hunk, "\n") {
		painted := line
		switch {
		case strings.HasPrefix(line, "@@"), line == "…":
			painted = colorize(color, termstyle.Muted, line)
		case strings.HasPrefix(line, "+"):
			painted = colorize(color, termstyle.Success, line)
		case strings.HasPrefix(line, "-"):
			painted = colorize(color, termstyle.Danger, line)
		}
		out = append(out, indent+painted)
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
	code := finding.CodeExcerpt
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
