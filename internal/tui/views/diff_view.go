package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

var (
	diffAddStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.Green))
	diffRemoveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.Red))
	diffHunkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.Cyan))
)

func RenderDiffView(m model.Model) string {
	const viewport = 20

	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("  Diff Preview:") + "\n\n")

	total := len(m.DiffLines)
	if total == 0 {
		sb.WriteString(styles.HelpStyle.Render("No diff to display\n"))
	} else {
		start := m.DiffViewOffset
		if start > total-viewport {
			start = max(0, total-viewport)
		}
		end := min(total, start+viewport)

		for i := start; i < end; i++ {
			line := m.DiffLines[i]
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				sb.WriteString(diffAddStyle.Render(line) + "\n")
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				sb.WriteString(diffRemoveStyle.Render(line) + "\n")
			} else if strings.HasPrefix(line, "@@") {
				sb.WriteString(diffHunkStyle.Render(line) + "\n")
			} else {
				sb.WriteString(styles.NormalStyle.Render(line) + "\n")
			}
		}
	}

	position := fmt.Sprintf("%d/%d", min(m.DiffViewOffset+1, total), total)
	sb.WriteString("\n" + styles.HelpStyle.Render(position))

	controlsWidget := controls.NewDiffViewControls()
	sb.WriteString("\n" + controlsWidget.Render())

	return sb.String()
}
