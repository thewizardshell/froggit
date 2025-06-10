package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderBranchView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render("Branches:") + "\n\n")

	for i, branch := range m.Branches {
		cursor := "  "
		if m.Cursor == i {
			cursor = ""
		}

		current := " "
		if branch == m.CurrentBranch {
			current = "●"
		}

		style := styles.NormalStyle
		if m.Cursor == i {
			style = styles.SelectedStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, current, branch)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n" + styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n")+
			styles.HelpStyle.Render("  [↑/↓] navigate  [Enter] switch branch  [n] new branch  [d] delete branch  [Esc] back"),
	))

	return s.String()
}
