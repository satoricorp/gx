package cli

import (
	"charm.land/lipgloss/v2"
)

const txTagline = "automating version control and simplifying code review"
const txLogoColor = "6"

const txLogoRaw = `████████╗ ██████╗ ████████╗ █████╗ ██╗     ██╗████████╗██╗   ██╗
╚══██╔══╝██╔═══██╗╚══██╔══╝██╔══██╗██║     ██║╚══██╔══╝╚██╗ ██╔╝
   ██║   ██║   ██║   ██║   ███████║██║     ██║   ██║    ╚████╔╝
   ██║   ██║   ██║   ██║   ██╔══██║██║     ██║   ██║     ╚██╔╝
   ██║   ╚██████╔╝   ██║   ██║  ██║███████╗██║   ██║      ██║
   ╚═╝    ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝      ╚═╝   `

func renderStaticLogo() string {
	return logoText(txLogoRaw)
}

func logoText(text string) string {
	if !enableColor() {
		return text
	}
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(txLogoColor))
	return style.Render(text)
}
