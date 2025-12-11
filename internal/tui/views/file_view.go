package view

import (
	"fmt"
	"strings"

	"froggit/internal/git"
	"froggit/internal/tui/controls"
	"froggit/internal/tui/icons"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// getFileStatusStyle returns the appropriate style based on file status
func getFileStatusStyle(file git.FileItem, isSelected bool) lipgloss.Style {
	// Check if file is added (A) or modified (M)
	isAdded := strings.Contains(file.Status, "A")
	isModified := strings.Contains(file.Status, "M")

	if isSelected {
		if isAdded {
			return styles.SelectedAddedStyle
		} else if isModified {
			return styles.SelectedModifiedStyle
		}
		return styles.SelectedStyle
	} else {
		if isAdded {
			return styles.AddedFileStyle
		} else if isModified {
			return styles.ModifiedFileStyle
		}
		return styles.NormalStyle
	}
}

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

	s.WriteString(styles.HeaderStyle.Render("  Git Status:") + "\n")
	s.WriteString(fmt.Sprintf("  Staged: %d files\n", stagedCount))
	s.WriteString(fmt.Sprintf("  Unstaged: %d files\n", unstagedCount))

	if m.HasRemoteChanges {
		s.WriteString(styles.WarningStyle.Render("  New commits are available on the remote please pull\n"))
	}
	s.WriteString("\n")
	s.WriteString(styles.HeaderStyle.Render(" Modified files:") + "\n\n")

	if len(m.Files) == 0 {
		s.WriteString(styles.HelpStyle.Render("No modified files\n"))
	} else {
		for i, file := range m.Files {
			cursor := "  "
			if m.Cursor == i {
				cursor = ""
			}

			staged := " "
			if file.Staged {
				staged = "✓"
			}

			isSelected := m.Cursor == i
			style := getFileStatusStyle(file, isSelected)

			icon := icons.GetIconForFile(file.Name)
			line := fmt.Sprintf("%s [%s] %s %s", cursor, staged, icon, file.Name)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	controlsWidget := controls.NewFileViewControls(stagedCount > 0, len(m.Files) > 0, m.AdvancedMode)
	s.WriteString("\n" + controlsWidget.Render())

	return s.String()
}
