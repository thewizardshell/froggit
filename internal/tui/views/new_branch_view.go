package view

import (
	"strings"

	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderNewBranchView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render("🌿 New Branch:") + "\n\n")
	s.WriteString(styles.InputStyle.Render(m.NewBranchName+"_") + "\n\n")

	s.WriteString(styles.BorderStyle.Render(
		styles.HelpStyle.Render("Type the branch name and press [Enter] to create\n") +
			styles.HelpStyle.Render("[Esc] to cancel"),
	))

	return s.String()
}
