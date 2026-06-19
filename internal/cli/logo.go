package cli

import (
	"charm.land/lipgloss/v2"
)

const gxTagline = "automating version control and simplifying code review"
const gxLogoColor = "6"

const gxLogoRaw = `   .aMMMMP dMP dMP
  dMP"    dMK.dMP
 dMP MMP".dMMMK"
dMP.dMP dMP"AMF
VMMMP" dMP dMP`

func renderStaticLogo() string {
	return logoText(gxLogoRaw)
}

func logoText(text string) string {
	if !enableColor() {
		return text
	}
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(gxLogoColor))
	return style.Render(text)
}
