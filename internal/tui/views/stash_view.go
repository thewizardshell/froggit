package view

import (
	"fmt"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
	"strings"
)

// RenderStashView renders the stash management view
func RenderStashView(m model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("Git Stash Manager") + "\n\n")

	if hasUnstagedChanges(m) {
		sb.WriteString(styles.WarningStyle.Render("⚠ You have unstaged changes that can be stashed") + "\n\n")
	}

	// Stash list
	if len(m.Stashes) == 0 {
		sb.WriteString(styles.HelpStyle.Render("No stashes found. Create one with [S] Save stash") + "\n\n")
	} else {
		sb.WriteString(styles.SubHeaderStyle.Render("Existing Stashes:") + "\n")

		for i, stash := range m.Stashes {
			cursor := "  "
			if i == m.Cursor {
				cursor = "❯ "
			}

			// Parse stash info
			stashInfo := parseStashInfo(stash)
			line := fmt.Sprintf("%s%s", cursor, stashInfo)

			if i == m.Cursor {
				sb.WriteString(styles.SelectedStyle.Render(line) + "\n")
			} else {
				sb.WriteString(styles.NormalStyle.Render(line) + "\n")
			}
		}
		sb.WriteString("\n")
	}

	if m.Message != "" {
		switch m.MessageType {
		case "success":
			sb.WriteString(styles.SuccessStyle.Render(m.Message) + "\n")
		case "error":
			sb.WriteString(styles.ErrorStyle.Render(m.Message) + "\n")
		case "warning":
			sb.WriteString(styles.WarningStyle.Render(m.Message) + "\n")
		default:
			sb.WriteString(styles.HelpStyle.Render(m.Message) + "\n")
		}
		sb.WriteString("\n")
	}

	if m.IsStashing {
		sb.WriteString(styles.HelpStyle.Render(fmt.Sprintf("%s Processing stash operation...", m.SpinnerFrames[m.SpinnerIndex])) + "\n\n")
	}

	// Controls
	controls := []string{}

	if hasUnstagedChanges(m) {
		controls = append(controls, "[S] Save stash")
	}

	if len(m.Stashes) > 0 {
		controls = append(controls,
			"[Enter] Apply stash",
			"[P] Pop stash",
			"[D] Drop stash",
			"[V] View stash",
		)
	}

	controls = append(controls, "[↑/↓] Navigate", "[Esc] Back", "[?] Help")

	controlsLine := "  " + strings.Join(controls, "  ")
	sb.WriteString(styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n") +
			styles.HelpStyle.Render(controlsLine),
	))

	return sb.String()
}

// RenderStashMessageView renders the stash message input view
func RenderStashMessageView(m model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderStyle.Render("💾 Save Stash") + "\n\n")

	sb.WriteString(styles.SubHeaderStyle.Render("Enter stash message (optional):") + "\n")

	// Input field with cursor
	inputDisplay := m.StashMessage
	if len(inputDisplay) == 0 {
		inputDisplay = "Work in progress..."
	}
	sb.WriteString(styles.InputStyle.Render(inputDisplay) + styles.CursorStyle.Render("│") + "\n\n")

	sb.WriteString(styles.HelpStyle.Render("📝 This will stash all your current changes") + "\n")
	sb.WriteString(styles.HelpStyle.Render("   including both staged and unstaged files") + "\n\n")

	// Status message
	if m.Message != "" {
		switch m.MessageType {
		case "error":
			sb.WriteString(styles.ErrorStyle.Render(m.Message) + "\n\n")
		default:
			sb.WriteString(styles.HelpStyle.Render(m.Message) + "\n\n")
		}
	}

	controls := []string{
		"[Enter] Save stash",
		"[Esc] Cancel",
		"[Backspace] Delete char",
	}
	controlsLine := "  " + strings.Join(controls, "  ")
	sb.WriteString(styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n") +
			styles.HelpStyle.Render(controlsLine),
	))

	return sb.String()
}

// hasUnstagedChanges checks if there are any changes that can be stashed
func hasUnstagedChanges(m model.Model) bool {
	return len(m.Files) > 0
}

// parseStashInfo formats stash information for display
func parseStashInfo(stashLine string) string {
	// Input: "stash@{0}: WIP on main: 1234567 commit message"
	// Output: "stash@{0}: WIP on main: commit message"

	parts := strings.SplitN(stashLine, ": ", 3)
	if len(parts) >= 3 {
		stashRef := parts[0]
		wipInfo := parts[1]
		commitInfo := parts[2]

		// Remove hash from commit info if present
		commitParts := strings.SplitN(commitInfo, " ", 2)
		if len(commitParts) >= 2 {
			commitMsg := commitParts[1]
			return fmt.Sprintf("%s: %s: %s", stashRef, wipInfo, commitMsg)
		}
		return fmt.Sprintf("%s: %s: %s", stashRef, wipInfo, commitInfo)
	}

	return stashLine
}
