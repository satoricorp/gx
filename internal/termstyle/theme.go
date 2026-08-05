package termstyle

import (
	"os"
	"strings"
	"sync"

	"charm.land/lipgloss/v2"
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
		styleSection = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
		styleCommand = lipgloss.NewStyle().Foreground(lipgloss.Color(colorCommand))
		styleMuted = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
		styleValue = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValue))
		styleURL = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent)).Underline(true)
	})
}

// Enabled reports whether gx should emit color and styled output.
func Enabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	return true
}

func paint(style lipgloss.Style, text string) string {
	initTheme()
	if text == "" || !Enabled() {
		return text
	}
	return Render(style.Render(text))
}

// Render applies lipgloss downsampling when writing styled strings to variables.
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
