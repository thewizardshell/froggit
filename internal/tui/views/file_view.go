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

// getFileStatusIndicator returns the status indicator (M, A, C, etc.) for a file
func getFileStatusIndicator(file git.FileItem) string {
	if strings.Contains(file.Status, "U") {
		return "C"
	} else if strings.Contains(file.Status, "A") {
		return "A"
	} else if strings.Contains(file.Status, "M") {
		return "M"
	} else if strings.Contains(file.Status, "D") {
		return "D"
	} else if strings.Contains(file.Status, "?") {
		return "?"
	}
	return " "
}

// getFileStatusStyle returns the appropriate style based on file status
func getFileStatusStyle(file git.FileItem, isSelected bool) lipgloss.Style {
	isConflict := strings.Contains(file.Status, "U")
	isAdded := strings.Contains(file.Status, "A")
	isModified := strings.Contains(file.Status, "M")
	isDeleted := strings.Contains(file.Status, "D")
	isUntracked := strings.Contains(file.Status, "?")

	if isSelected {
		if isConflict {
			return styles.SelectedConflictStyle
		} else if isAdded {
			return styles.SelectedAddedStyle
		} else if isModified {
			return styles.SelectedModifiedStyle
		} else if isDeleted {
			return styles.SelectedDeletedStyle
		} else if isUntracked {
			return styles.SelectedUntrackedStyle
		}
		return styles.SelectedStyle
	} else {
		if isConflict {
			return styles.ConflictFileStyle
		} else if isAdded {
			return styles.AddedFileStyle
		} else if isModified {
			return styles.ModifiedFileStyle
		} else if isDeleted {
			return styles.DeletedFileStyle
		} else if isUntracked {
			return styles.UntrackedFileStyle
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
		totalFiles := len(m.Files)
		viewHeight := m.FileViewHeight
		if viewHeight <= 0 {
			viewHeight = 15
		}

		if m.FileViewOffset > 0 {
			s.WriteString(styles.HelpStyle.Render(fmt.Sprintf("  ↑ %d more above\n", m.FileViewOffset)))
		}

		endIdx := m.FileViewOffset + viewHeight
		if endIdx > totalFiles {
			endIdx = totalFiles
		}

		for i := m.FileViewOffset; i < endIdx; i++ {
			file := m.Files[i]
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
			statusIndicator := getFileStatusIndicator(file)
			line := fmt.Sprintf("%s [%s] %s %s %s", cursor, staged, statusIndicator, icon, file.Name)
			s.WriteString(style.Render(line) + "\n")
		}

		remaining := totalFiles - endIdx
		if remaining > 0 {
			s.WriteString(styles.HelpStyle.Render(fmt.Sprintf("  ↓ %d more below\n", remaining)))
		}
	}

	controlsWidget := controls.NewFileViewControls(stagedCount > 0, len(m.Files) > 0, m.AdvancedMode)
	s.WriteString("\n" + controlsWidget.Render())

	return s.String()
}
