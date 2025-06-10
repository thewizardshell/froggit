package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderConfirmDialog(m model.Model) string {
	var s strings.Builder
	var message string

	switch m.DialogType {
	case "delete_branch":
		message = fmt.Sprintf("Are you sure you want to delete branch '%s'?", m.DialogTarget)
	case "discard_changes":
		message = fmt.Sprintf("Are you sure you want to discard changes in '%s'?", m.DialogTarget)
	}

	s.WriteString("\n\n")
	s.WriteString(styles.BorderStyle.Render(
		styles.HeaderStyle.Render(" Confirm action") + "\n\n" +
			styles.NormalStyle.Render(message) + "\n\n" +
			styles.HelpStyle.Render("[y] Yes  [n] No"),
	))

	return s.String()
}
