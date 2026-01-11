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

	if m.IsGeneratingAI {
		s.WriteString(styles.InputStyle.Render(m.SpinnerFrames[m.SpinnerIndex]+" Generating...") + "\n\n")
	} else {
		s.WriteString(styles.InputStyle.Render(m.CommitMsg+"_") + "\n\n")
	}

	controlsWidget := controls.NewCommitViewControls(m.CopilotAvailable, m.IsGeneratingAI)
	s.WriteString(controlsWidget.Render())

	return s.String()
}
