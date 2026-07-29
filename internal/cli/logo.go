package cli

import (
	"charm.land/lipgloss/v2"
)

const tlTagline = "automating version control and simplifying code review"
const tlLogoColor = "6"

const tlLogoRaw = `████████╗ ██████╗ ████████╗ █████╗ ██╗     ██╗████████╗██╗   ██╗
╚══██╔══╝██╔═══██╗╚══██╔══╝██╔══██╗██║     ██║╚══██╔══╝╚██╗ ██╔╝
   ██║   ██║   ██║   ██║   ███████║██║     ██║   ██║    ╚████╔╝
   ██║   ██║   ██║   ██║   ██╔══██║██║     ██║   ██║     ╚██╔╝
   ██║   ╚██████╔╝   ██║   ██║  ██║███████╗██║   ██║      ██║
   ╚═╝    ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝      ╚═╝   `

func renderStaticLogo() string {
	return logoText(tlLogoRaw)
}

func logoText(text string) string {
	if !enableColor() {
		return text
	}
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(tlLogoColor))
	return style.Render(text)
}
