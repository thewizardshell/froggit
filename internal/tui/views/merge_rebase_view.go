package view

import (
	"fmt"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
	"strings"
)

// RenderMergeView renders the interactive merge view.
func RenderMergeView(m model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("[MERGE] Select a branch to merge into:") + "\n")
	for i, branch := range m.Branches {
		cursor := "  "
		if i == m.Cursor {
			cursor = "❯ " // Clear visual cursor indicator
		}

		checked := " "
		if m.DialogTarget == branch {
			checked = "✓"
		}

		// Apply style based on selection
		line := fmt.Sprintf("%s[%s] %s", cursor, checked, branch)
		if i == m.Cursor {
			sb.WriteString(styles.SelectedStyle.Render(line) + "\n")
		} else {
			sb.WriteString(styles.NormalStyle.Render(line) + "\n")
		}
	}

	if m.Message != "" {
		sb.WriteString("\n" + styles.HelpStyle.Render(m.Message) + "\n")
	}
	if len(m.LogLines) > 0 {
		sb.WriteString(styles.WarningStyle.Render("\nConflicts detected in:") + "\n")
		for _, file := range m.LogLines {
			sb.WriteString(styles.HelpStyle.Render("- " + file + "\n"))
		}
		sb.WriteString(styles.HelpStyle.Render("[P] Proceed (merge --continue)  [X] Cancel (merge --abort)\n"))
	}

	controls := []string{
		"[↑/↓] Navigate",
		"[Space] Select",
		"[M] Merge",
		"[Esc] Back",
		"[?] Help",
	}
	controlsLine := "  " + strings.Join(controls, "  ")
	sb.WriteString("\n" + styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n")+
			styles.HelpStyle.Render(controlsLine),
	))

	return sb.String()
}

// RenderRebaseView renders the interactive rebase view.
func RenderRebaseView(m model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("[REBASE] Select a branch to rebase onto:") + "\n")
	for i, branch := range m.Branches {
		cursor := "  "
		if i == m.Cursor {
			cursor = "❯ " // Clear visual cursor indicator
		}

		checked := " "
		if m.DialogTarget == branch {
			checked = "✓"
		}

		// Apply style based on selection
		line := fmt.Sprintf("%s[%s] %s", cursor, checked, branch)
		if i == m.Cursor {
			sb.WriteString(styles.SelectedStyle.Render(line) + "\n")
		} else {
			sb.WriteString(styles.NormalStyle.Render(line) + "\n")
		}
	}

	if m.Message != "" {
		sb.WriteString("\n" + styles.HelpStyle.Render(m.Message) + "\n")
	}
	if len(m.LogLines) > 0 {
		sb.WriteString(styles.WarningStyle.Render("\nConflicts detected in:") + "\n")
		for _, file := range m.LogLines {
			sb.WriteString(styles.HelpStyle.Render("- " + file + "\n"))
		}
		sb.WriteString(styles.HelpStyle.Render("[P] Proceed (rebase --continue)  [X] Cancel (rebase --abort)\n"))
	}

	controls := []string{
		"[↑/↓] Navigate",
		"[Space] Select",
		"[R] Rebase",
		"[Esc] Back",
		"[?] Help",
	}
	controlsLine := "  " + strings.Join(controls, "  ")
	sb.WriteString("\n" + styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n")+
			styles.HelpStyle.Render(controlsLine),
	))

	return sb.String()
}
