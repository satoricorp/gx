package termstyle

import (
	"os"
	"strings"
	"sync"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

// Paper CLI palette from Console — PR Console / CLI Terminal artboards.
const (
	colorPrompt  = "#71717A"
	colorCommand = "#6366F1"
	colorMuted   = "#A1A1AA"
	colorSuccess = "#4ADE80"
	colorDanger  = "#F87171"
	colorWarning = "#F59E0B"
	colorAccent  = "#6366F1"
	colorText    = "#FAFAFA"
	colorValue   = "#D4D4D8"
)

var (
	initOnce sync.Once

	// Semantic styles (initialized in initTheme).
	styleDanger  lipgloss.Style
	styleSuccess lipgloss.Style
	styleWarning lipgloss.Style
	styleSection lipgloss.Style
	styleCommand lipgloss.Style
	styleMuted   lipgloss.Style
	styleValue   lipgloss.Style
	styleURL     lipgloss.Style

	labelWidth = 12
)

func initTheme() {
	initOnce.Do(func() {
		styleDanger = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDanger))
		styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSuccess))
		styleWarning = lipgloss.NewStyle().Foreground(lipgloss.Color(colorWarning))
		styleSection = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Bold(true)
		styleCommand = lipgloss.NewStyle().Foreground(lipgloss.Color(colorCommand))
		styleMuted = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
		styleValue = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValue))
		styleURL = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)).Underline(true)
	})
}

// Enabled reports whether gx should emit color and styled output.
//
// The decision has three layers, in order:
//
//  1. NO_COLOR set, or TERM=dumb → never. This is the long-standing contract.
//  2. FORCE_COLOR set (any non-empty value except "0") → always. This is how a
//     CI log viewer, a pager, or a test asserts color through a pipe.
//  3. Otherwise → only when stdout is a terminal.
//
// Layer 3 is the one that used to be missing: color was decided from the
// environment alone, so `gx review | tee out.txt` wrote ANSI into the file.
// The interactive report — the accordion, the lanes, the code windows — must
// fall back to plain text the moment it is piped, and it needs the same fd
// answer the bubbletea loader gets from term.IsTerminal, not a different one.
//
// stdoutIsTerminal is a variable so tests can pin either answer without
// owning a real tty; the default reads os.Stdout.
func Enabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	if force := os.Getenv("FORCE_COLOR"); force != "" && force != "0" {
		return true
	}
	return stdoutIsTerminal()
}

// stdoutIsTerminal reports whether stdout is attached to a terminal. It is a
// package variable, not a function, so tests can substitute an answer.
var stdoutIsTerminal = func() bool {
	return term.IsTerminal(os.Stdout.Fd())
}

func paint(style lipgloss.Style, text string) string {
	initTheme()
	if text == "" || !Enabled() {
		return text
	}
	// Enabled() has already decided — NO_COLOR, TERM=dumb, FORCE_COLOR, or
	// an isatty check on stdout. Rendering the style directly honors that
	// decision. Routing through lipgloss.Sprint would re-detect a color
	// profile from the process's real stdout and downsample to no color
	// whenever that is a pipe, which silently stripped every painter under
	// FORCE_COLOR and in tests: only the raw-ANSI accents survived, and the
	// report came out mostly monochrome.
	return style.Render(text)
}

// Render applies lipgloss downsampling when writing styled strings to variables.
// It is kept for callers that assemble lipgloss layouts (JoinHorizontal and
// friends) and want the profile-aware output; the plain painters do not use it.
func Render(s string) string {
	initTheme()
	if !Enabled() || s == "" {
		return s
	}
	return lipgloss.Sprint(s)
}

func Danger(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleDanger, text)
}

func LabelWarning(label, value string) string {
	initTheme()
	if !Enabled() {
		return label + " " + value
	}
	lbl := styleMuted.Width(labelWidth).Render(label)
	val := styleWarning.Render(value)
	return Render(lipgloss.JoinHorizontal(lipgloss.Top, lbl, "  ", val))
}

func Success(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleSuccess, text)
}

// Warning paints text in the warning color: the advisory lane, a coverage
// caveat, anything reported but not blocking. It sits between Success and
// Danger in the same shape so a renderer can pick a painter by severity.
func Warning(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleWarning, text)
}

func Section(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleSection, text)
}

func Command(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleCommand, text)
}

func CommandLine(invocation string, accentCommand bool) string {
	initTheme()
	invocation = strings.TrimSpace(invocation)
	if invocation == "" {
		return ""
	}
	cmd := invocation
	if accentCommand {
		cmd = Command(invocation)
	} else {
		cmd = Value(invocation)
	}
	return paint(lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrompt)), "$") + " " + cmd
}

func Muted(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleMuted, text)
}

func Value(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleValue, text)
}

// LabelValue renders a fixed-width label with a value.
func LabelValue(label, value string) string {
	initTheme()
	if !Enabled() {
		return label + " " + value
	}
	lbl := styleMuted.Width(labelWidth).Render(label)
	val := styleValue.Render(value)
	return Render(lipgloss.JoinHorizontal(lipgloss.Top, lbl, "  ", val))
}

// Hyperlink renders a terminal hyperlink when supported.
func Hyperlink(url, text string) string {
	initTheme()
	if text == "" {
		text = url
	}
	if !Enabled() || url == "" {
		return text
	}
	return paint(styleURL.Hyperlink(url), text)
}
