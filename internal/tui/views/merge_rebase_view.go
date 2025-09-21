package view

import (
	"fmt"
	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
	"strings"
)

func RenderMergeView(m model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("[MERGE] Select a branch to merge into:") + "\n")
	for i, branch := range m.Branches {
		cursor := "  "
		if i == m.Cursor {
			cursor = "❯ "
		}

		checked := " "
		if m.DialogTarget == branch {
			checked = "✓"
		}

		line := fmt.Sprintf("%s[%s] %s", cursor, checked, branch)
		if i == m.Cursor {
			sb.WriteString(styles.SelectedStyle.Render(line) + "\n")
		} else {
			sb.WriteString(styles.NormalStyle.Render(line) + "\n")
		}
	}

	if len(m.LogLines) > 0 {
		sb.WriteString(styles.WarningStyle.Render("\nConflicts detected in:") + "\n")
		for _, file := range m.LogLines {
			sb.WriteString(styles.HelpStyle.Render("- " + file + "\n"))
		}
		sb.WriteString(styles.HelpStyle.Render("[P] Proceed (merge --continue)  [X] Cancel (merge --abort)\n"))
	}

	hasSelection := m.DialogTarget != ""
	controlsWidget := controls.NewMergeViewControls(hasSelection)
	sb.WriteString("\n" + controlsWidget.Render())

	return sb.String()
}

func RenderRebaseView(m model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("[REBASE] Select a branch to rebase onto:") + "\n")
	for i, branch := range m.Branches {
		cursor := "  "
		if i == m.Cursor {
			cursor = "❯ "
		}

		checked := " "
		if m.DialogTarget == branch {
			checked = "✓"
		}

		line := fmt.Sprintf("%s[%s] %s", cursor, checked, branch)
		if i == m.Cursor {
			sb.WriteString(styles.SelectedStyle.Render(line) + "\n")
		} else {
			sb.WriteString(styles.NormalStyle.Render(line) + "\n")
		}
	}

	if len(m.LogLines) > 0 {
		sb.WriteString(styles.WarningStyle.Render("\nConflicts detected in:") + "\n")
		for _, file := range m.LogLines {
			sb.WriteString(styles.HelpStyle.Render("- " + file + "\n"))
		}
		sb.WriteString(styles.HelpStyle.Render("[P] Proceed (rebase --continue)  [X] Cancel (rebase --abort)\n"))
	}

	hasSelection := m.DialogTarget != ""
	controlsWidget := controls.NewRebaseViewControls(hasSelection)
	sb.WriteString("\n" + controlsWidget.Render())

	return sb.String()
}
