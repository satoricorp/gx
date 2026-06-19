package clitui

import (
	"os"
	"sync"

	"charm.land/lipgloss/v2"
)

// Paper CLI palette from Console — PR Console / CLI Terminal artboards.
const (
	colorPrompt    = "#71717A"
	colorCommand   = "#6366F1"
	colorMuted     = "#A1A1AA"
	colorSuccess   = "#4ADE80"
	colorDanger    = "#F87171"
	colorAccent    = "#6366F1"
	colorSecondary = "#A1A1AA"
	colorText      = "#FAFAFA"
	colorValue     = "#D4D4D8"
	colorDivider   = "#1F1F22"
)

var (
	initOnce sync.Once

	stylePrompt    lipgloss.Style
	styleCommand   lipgloss.Style
	styleMuted     lipgloss.Style
	styleSuccess   lipgloss.Style
	styleDanger    lipgloss.Style
	styleAccent    lipgloss.Style
	styleSecondary lipgloss.Style
	styleText      lipgloss.Style
	styleValue     lipgloss.Style
	styleHint      lipgloss.Style
	styleDivider   lipgloss.Style
)

func initTheme() {
	initOnce.Do(func() {
		stylePrompt = lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrompt))
		styleCommand = lipgloss.NewStyle().Foreground(lipgloss.Color(colorCommand))
		styleMuted = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
		styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSuccess))
		styleDanger = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDanger))
		styleAccent = lipgloss.NewStyle().Foreground(lipgloss.Color(colorAccent))
		styleSecondary = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondary))
		styleText = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText))
		styleValue = lipgloss.NewStyle().Foreground(lipgloss.Color(colorValue))
		styleHint = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMuted))
		styleDivider = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDivider))
	})
}

func enabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return true
}

func paint(style lipgloss.Style, text string) string {
	initTheme()
	if text == "" || !enabled() {
		return text
	}
	return lipgloss.Sprint(style.Render(text))
}

func prompt(text string) string    { return paint(stylePrompt, text) }
func command(text string) string   { return paint(styleCommand, text) }
func muted(text string) string     { return paint(styleMuted, text) }
func success(text string) string   { return paint(styleSuccess, text) }
func accent(text string) string    { return paint(styleAccent, text) }
func secondary(text string) string { return paint(styleSecondary, text) }
func text(text string) string      { return paint(styleText, text) }
func value(text string) string     { return paint(styleValue, text) }
func hint(text string) string      { return paint(styleHint, text) }
