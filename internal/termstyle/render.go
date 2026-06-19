package termstyle

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"charm.land/lipgloss/v2/tree"
)

// StackRevision is one entry in a gx stack listing.
type StackRevision struct {
	Index       int
	Description string
	Active      bool
	Published   bool
}

// RenderStack formats stack metadata and revisions for terminal output.
func RenderStack(containerName string, publishedCount int, revisions []StackRevision) string {
	initTheme()
	if len(revisions) == 0 {
		return Muted("No GX stack recorded yet.")
	}

	lines := []string{
		Section("Stack"),
		LabelValue("Container", containerName),
		LabelValue("Revisions", fmt.Sprintf("%d total, %d published", len(revisions), publishedCount)),
		"",
	}

	if !Enabled() {
		for _, revision := range revisions {
			lines = append(lines, plainStackLine(revision))
		}
		return strings.Join(lines, "\n")
	}

	borderColor := lipgloss.Color(colorBorder)

	root := tree.Root(containerName).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(borderColor)).
		RootStyle(styleSection).
		ItemStyle(styleValue)

	for _, revision := range revisions {
		label := stackRevisionLabel(revision)
		if revision.Active {
			root = root.Child(styleActive.Render(label))
		} else {
			root = root.Child(label)
		}
	}
	lines = append(lines, Render(root.String()))
	return strings.Join(lines, "\n")
}

func plainStackLine(revision StackRevision) string {
	marker := " "
	if revision.Active {
		marker = "*"
	}
	line := fmt.Sprintf("  %s %d. %s", marker, revision.Index, revision.Description)
	tags := stackTags(revision)
	if len(tags) > 0 {
		line += "  [" + strings.Join(tags, ", ") + "]"
	}
	return line
}

func stackRevisionLabel(revision StackRevision) string {
	prefix := fmt.Sprintf("%d. %s", revision.Index, revision.Description)
	tags := stackTags(revision)
	if len(tags) == 0 {
		return prefix
	}
	tagText := strings.Join(tags, ", ")
	if !Enabled() {
		return prefix + "  [" + tagText + "]"
	}
	return Render(
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			styleValue.Render(prefix),
			"  ",
			styleTag.Render("["+tagText+"]"),
		),
	)
}

func stackTags(revision StackRevision) []string {
	var tags []string
	if revision.Published {
		tags = append(tags, "published")
	}
	if revision.Active {
		tags = append(tags, "active")
	}
	return tags
}

// DoctorRow is one health-check line for gx doctor.
type DoctorRow struct {
	Check  string
	Status string
}

// RenderDoctorTable formats doctor checks as a bordered table.
func RenderDoctorTable(rows []DoctorRow) string {
	initTheme()
	if len(rows) == 0 {
		return ""
	}
	if !Enabled() {
		var lines []string
		for _, row := range rows {
			lines = append(lines, row.Check+"  "+row.Status)
		}
		return strings.Join(lines, "\n")
	}

	borderColor := lipgloss.Color(colorBorder)

	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorSecondary)).Padding(0, 1)
	checkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorValue)).Padding(0, 1)
	cellStyle := lipgloss.NewStyle().Padding(0, 1)

	tableRows := make([][]string, len(rows))
	for i, row := range rows {
		tableRows[i] = []string{row.Check, row.Status}
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(borderColor)).
		Headers("CHECK", "STATUS").
		Rows(tableRows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			if col == 0 {
				return checkStyle
			}
			dataRow := row - 1
			if dataRow < 0 || dataRow >= len(tableRows) {
				return cellStyle
			}
			status := tableRows[dataRow][1]
			switch {
			case strings.HasPrefix(status, "ok"):
				return cellStyle.Inherit(styleOK)
			case strings.HasPrefix(status, "warn"):
				return cellStyle.Inherit(styleWarn)
			default:
				return cellStyle.Inherit(styleValue)
			}
		})

	return Render(t.String())
}
