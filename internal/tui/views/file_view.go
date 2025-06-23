package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/icons"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

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

	s.WriteString(styles.HeaderStyle.Render("  Git Status:") + "\n")
	s.WriteString(fmt.Sprintf("  Staged: %d files\n", stagedCount))
	s.WriteString(fmt.Sprintf("  Unstaged: %d files\n", unstagedCount))

	if m.HasRemoteChanges {
		s.WriteString(styles.WarningStyle.Render("  New commits are available on the remote please pull\n"))
	}
	s.WriteString("\n")
	s.WriteString(styles.HeaderStyle.Render(" Modified files:") + "\n\n")

	if len(m.Files) == 0 {
		s.WriteString(styles.HelpStyle.Render("No modified files\n"))
	} else {
		for i, file := range m.Files {
			cursor := "  "
			if m.Cursor == i {
				cursor = ""
			}

			staged := " "
			if file.Staged {
				staged = "✓"
			}

			style := styles.NormalStyle
			if m.Cursor == i {
				style = styles.SelectedStyle
			}

			icon := icons.GetIconForFile(file.Name)
			line := fmt.Sprintf("%s [%s] %s %s", cursor, staged, icon, file.Name)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	var controls []string
	if m.AdvancedMode {
		controls = []string{
			"[↑/↓] navigate",
			"[L] log graph",
			"[esc] exit advanced",
			"[?] help",
		}
		s.WriteString(styles.WarningStyle.Render("\n\n🚧 Advanced features (merge, stash, rebase) coming soon!\n"))
	} else {
		controls = []string{"[↑/↓] navigate", "[space] stage/unstage"}
		if len(m.Files) > 0 {
			controls = append(controls, "[x] discard changes")
		}
		if stagedCount > 0 {
			controls = append(controls, "[c] commit")
		}
		controls = append(controls, "[?] help")
	}

	controlsLine := "  " + strings.Join(controls, "  ")

	s.WriteString("\n" + styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n")+
			styles.HelpStyle.Render(controlsLine),
	))

	return s.String()
}
