package cli

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
)

type tuiDiffView struct {
	title string
	view  viewport.Model
}

func newTUIDiffView() tuiDiffView {
	return tuiDiffView{view: viewport.New(viewport.WithWidth(120), viewport.WithHeight(24))}
}

func (d *tuiDiffView) resize(width, height int) {
	if width <= 0 {
		width = 120
	}
	viewHeight := height - 5
	if viewHeight < 5 {
		viewHeight = 5
	}
	d.view.SetWidth(width)
	d.view.SetHeight(viewHeight)
}

func (d *tuiDiffView) open(title, diff string) {
	d.title = title
	d.view.SetContent(colorizeDiff(diff))
	d.view.GotoTop()
}

func (d tuiDiffView) render(commandLabel string) string {
	var out strings.Builder
	out.WriteString(commandLine(commandLabel, true))
	out.WriteString("\n\n")
	out.WriteString(valueText(d.title))
	out.WriteString("\n\n")
	out.WriteString(d.view.View())
	out.WriteString("\n")
	out.WriteString(muted("j/k scroll · pgup/pgdn page · h/l horizontal · esc revisions · q quit"))
	out.WriteString("\n")
	return out.String()
}

func colorizeDiff(diff string) string {
	if strings.Contains(diff, "\x1b[") {
		return diff
	}
	lines := strings.Split(diff, "\n")
	for index, line := range lines {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			lines[index] = accent(line)
		case strings.HasPrefix(line, "@@"):
			lines[index] = command(line)
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			lines[index] = muted(line)
		case strings.HasPrefix(line, "+"):
			lines[index] = success(line)
		case strings.HasPrefix(line, "-"):
			lines[index] = danger(line)
		default:
			lines[index] = valueText(line)
		}
	}
	return strings.Join(lines, "\n")
}
