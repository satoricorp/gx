package termstyle

import (
	"fmt"
	"os"
	"strings"
	"sync"

	bubblesprogress "charm.land/bubbles/v2/progress"
	"charm.land/lipgloss/v2"
)

// Paper CLI palette from Console — PR Console / CLI Terminal artboards.
const (
	colorPrompt        = "#71717A"
	colorCommand       = "#6366F1"
	colorMuted         = "#A1A1AA"
	colorMint          = "#3DDC97"
	colorSuccess       = "#4ADE80"
	colorDanger        = "#F87171"
	colorWarning       = "#F59E0B"
	colorAccent        = "#6366F1"
	colorSecondary     = "#A1A1AA"
	colorText          = "#FAFAFA"
	colorValue         = "#D4D4D8"
	colorDivider       = "#1F1F22"
	colorBorder        = "#27272A"
	colorProgressEmpty = "#606060"
)

var (
	initOnce sync.Once

	// Semantic styles (initialized in initTheme).
	styleDanger  lipgloss.Style
	styleSuccess lipgloss.Style
	styleWarning lipgloss.Style
	styleSection lipgloss.Style
	styleCommand lipgloss.Style
	styleMint    lipgloss.Style
	styleMuted   lipgloss.Style
	styleAccent  lipgloss.Style
	styleValue   lipgloss.Style
	styleHint    lipgloss.Style
	styleOK      lipgloss.Style
	styleWarn    lipgloss.Style
	styleActive  lipgloss.Style
	styleTag     lipgloss.Style
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
		styleMint = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMint))
		styleMuted = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
		styleAccent = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
		styleValue = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValue))
		styleHint = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrompt))
		styleOK = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSuccess))
		styleWarn = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondary))
		styleActive = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
		styleTag = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrompt))
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

func Warning(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleWarning, text)
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

func LabelMint(label, value string) string {
	initTheme()
	if !Enabled() {
		return label + " " + value
	}
	lbl := styleMuted.Width(labelWidth).Render(label)
	val := styleMint.Render(value)
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

func Mint(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleMint, text)
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

func Accent(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleAccent, text)
}

func Value(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleValue, text)
}

func Secondary(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	if !Enabled() {
		return text
	}
	return paint(lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondary)), text)
}

func Hint(text string) string {
	initTheme()
	if text == "" {
		return text
	}
	return paint(styleHint, text)
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

// FormatStatus styles values that use gx ok:/warn: conventions.
func FormatStatus(value string) string {
	initTheme()
	if !Enabled() {
		return value
	}
	trimmed := strings.TrimSpace(value)
	switch {
	case strings.HasPrefix(trimmed, "ok:"):
		return paint(styleOK, trimmed)
	case strings.HasPrefix(trimmed, "warn:"):
		return paint(styleWarn, trimmed)
	case trimmed == "ok" || trimmed == "updated":
		return paint(styleOK, trimmed)
	default:
		return paint(styleValue, value)
	}
}

// LabelStatus is LabelValue with status coloring on the value.
func LabelStatus(label, value string) string {
	initTheme()
	if !Enabled() {
		return label + " " + value
	}
	lbl := styleMuted.Width(labelWidth).Render(label)
	val := FormatStatus(value)
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

// PromptLabel styles interactive prompt labels (gx init, modify).
func PromptLabel(text string) string {
	initTheme()
	if text == "" || !Enabled() {
		return text
	}
	return paint(styleCommand, text)
}

func Divider(width int) string {
	if width <= 0 {
		width = 40
	}
	return paint(lipgloss.NewStyle().Foreground(lipgloss.Color(colorDivider)), strings.Repeat("─", width))
}

func ProgressBar(percent float64) string {
	initTheme()
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}
	bar := bubblesprogress.New(
		bubblesprogress.WithWidth(40),
		bubblesprogress.WithColors(lipgloss.Color(colorAccent)),
	)
	bar.EmptyColor = lipgloss.Color(colorProgressEmpty)
	bar.PercentageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondary))
	if !Enabled() {
		return fmt.Sprintf("%3.0f%%", percent*100)
	}
	return Render(bar.ViewAs(percent))
}
