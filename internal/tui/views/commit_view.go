package view

import (
	"strings"

	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderCommitView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render(" Commit message:") + "\n\n")
	s.WriteString(styles.InputStyle.Render(m.CommitMsg+"_") + "\n\n")

	s.WriteString(styles.BorderStyle.Render(
		styles.HelpStyle.Render("Type your message and press [Enter] to confirm\n") +
			styles.HelpStyle.Render("[Esc] to cancel"),
	))

	return s.String()
}
