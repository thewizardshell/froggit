package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderRemoteView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render(" Remote repositories:") + "\n\n")

	if len(m.Remotes) == 0 {
		s.WriteString(styles.HelpStyle.Render("No remote repositories configured\n"))
	} else {
		for i, remote := range m.Remotes {
			cursor := "  "
			if m.Cursor == i {
				cursor = " "
			}

			style := styles.NormalStyle
			if m.Cursor == i {
				style = styles.SelectedStyle
			}

			line := fmt.Sprintf("%s %s", cursor, remote)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	controlsWidget := controls.NewRemoteViewControls()
	s.WriteString("\n" + controlsWidget.Render())

	return s.String()
}
