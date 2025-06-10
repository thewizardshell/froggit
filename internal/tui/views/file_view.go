package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/icons"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

// RenderFileView dibuja la vista de archivos modificados
func RenderFileView(m model.Model) string {
	var s strings.Builder

	stagedCount := 0
	unstagedCount := 0
	for _, file := range m.Files {
		if file.Staged {
			stagedCount++
		} else {
			unstagedCount++
		}
	}

	s.WriteString(styles.HeaderStyle.Render(" Git Status:") + "\n")
	s.WriteString(fmt.Sprintf(" Stage: %d files\n", stagedCount))
	s.WriteString(fmt.Sprintf(" Unstaged: %d files\n", unstagedCount))
	s.WriteString("\n")
	s.WriteString(styles.HeaderStyle.Render(" Modified files:") + "\n\n")

	if len(m.Files) == 0 {
		s.WriteString(styles.HelpStyle.Render("No modified files\n"))
	} else {
		for i, file := range m.Files {
			cursor := "  "
			if m.Cursor == i {
				cursor = ""
			}

			staged := " "
			if file.Staged {
				staged = "✓"
			}

			style := styles.NormalStyle
			if m.Cursor == i {
				style = styles.SelectedStyle
			}

			icon := icons.GetIconForFile(file.Name)
			line := fmt.Sprintf("%s [%s] %s %s", cursor, staged, icon, file.Name)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	s.WriteString("\n" + styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n")+
			styles.HelpStyle.Render("  [↑/↓] navigate  [space] stage/unstage  [a] stage all  [x] discard changes")+
			styles.HelpStyle.Render("  [c] commit  [b] branches  [m] remotes  [p] push  [f] fetch  [l] pull  [r] refresh  [q] quit"),
	))

	return s.String()
}
