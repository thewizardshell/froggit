package view

import (
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderCommitView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render(" Commit message:") + "\n\n")
	s.WriteString(styles.InputStyle.Render(m.CommitMsg+"_") + "\n\n")

	controlsWidget := controls.NewCommitViewControls()
	s.WriteString(controlsWidget.Render())

	return s.String()
}
