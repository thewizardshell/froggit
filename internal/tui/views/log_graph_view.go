package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

// RenderLogGraphView renders the interactive git log graph.
// It highlights the currently selected commit line.
func RenderLogGraphView(m model.Model) string {
	const viewport = 15 // number of lines to show at once

	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("  Git Log Graph:") + "\n\n")

	total := len(m.LogLines)
	if total == 0 {
		sb.WriteString(styles.HelpStyle.Render("No commits found\n"))
	} else {
		// Determine window start/end so cursor is always visible
		start := 0
		if m.Cursor >= viewport/2 {
			start = m.Cursor - viewport/2
		}
		if start+viewport > total {
			start = max(0, total-viewport)
		}
		end := min(total, start+viewport)

		for i := start; i < end; i++ {
			line := m.LogLines[i]
			cursor := "  "
			style := styles.NormalStyle
			if m.Cursor == i {
				cursor = "𓆏"
				style = styles.SelectedStyle
			}

			// Colourize parts
			parts := strings.SplitN(line, " ", 3)
			if len(parts) >= 2 {
				graph := styles.GraphSymbolStyle.Render(parts[0])
				hash := styles.CommitHashStyle.Render(parts[1])
				rest := ""
				if len(parts) == 3 {
					rest = parts[2]
				}
				line = fmt.Sprintf("%s %s %s", graph, hash, rest)
			}

			sb.WriteString(style.Render(cursor+" "+line) + "\n")
		}
	}

	// Controls with position info
	controls := fmt.Sprintf("[↑/↓] navigate  [esc] back   %d/%d", m.Cursor+1, total)
	sb.WriteString("\n" + styles.BorderStyle.Render(styles.HelpStyle.Render(controls)))

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
