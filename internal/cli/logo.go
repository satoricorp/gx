package cli

import (
	"charm.land/lipgloss/v2"
)

const lgtmTagline = "automating version control and simplifying code review"
const lgtmLogoColor = "6"

const lgtmLogoRaw = `████████╗ ██████╗ ████████╗ █████╗ ██╗     ██╗████████╗██╗   ██╗
╚══██╔══╝██╔═══██╗╚══██╔══╝██╔══██╗██║     ██║╚══██╔══╝╚██╗ ██╔╝
   ██║   ██║   ██║   ██║   ███████║██║     ██║   ██║    ╚████╔╝
   ██║   ██║   ██║   ██║   ██╔══██║██║     ██║   ██║     ╚██╔╝
   ██║   ╚██████╔╝   ██║   ██║  ██║███████╗██║   ██║      ██║
   ╚═╝    ╚═════╝    ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝      ╚═╝   `

func renderStaticLogo() string {
	return logoText(lgtmLogoRaw)
}

func logoText(text string) string {
	if !enableColor() {
		return text
	}
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(lgtmLogoColor))
	return style.Render(text)
}
