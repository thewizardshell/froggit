package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderConfirmDialog(m model.Model) string {
	var s strings.Builder
	var title, message, icon string

	switch m.DialogType {
	case "delete_branch":
		icon = "🗑️"
		title = "Delete Branch"
		message = fmt.Sprintf("Are you sure you want to delete branch '%s'?", styles.WarningStyle.Render(m.DialogTarget))
	case "delete_remote":
		icon = "🔗"
		title = "Remove Remote"
		message = fmt.Sprintf("Are you sure you want to remove remote '%s'?", styles.WarningStyle.Render(m.DialogTarget))
	case "discard_changes":
		icon = "⚠️"
		title = "Discard Changes"
		message = fmt.Sprintf("Are you sure you want to discard changes in '%s'?", styles.WarningStyle.Render(m.DialogTarget))
	case "drop_stash":
		icon = "💥"
		title = "Drop Stash"
		message = fmt.Sprintf("Are you sure you want to drop stash '%s'?", styles.WarningStyle.Render(m.DialogTarget))
	default:
		icon = "❓"
		title = "Confirm Action"
		message = "Are you sure you want to proceed?"
	}

	s.WriteString("\n\n")

	s.WriteString(styles.HeaderStyle.Render(icon + " " + title) + "\n\n")
	s.WriteString(styles.NormalStyle.Render(message) + "\n")
	s.WriteString(styles.HelpStyle.Render("This action cannot be undone.") + "\n\n")

	controlsWidget := controls.NewConfirmDialogControls()
	s.WriteString(controlsWidget.Render())

	return s.String()
}
