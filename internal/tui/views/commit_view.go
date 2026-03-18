package view

import (
	"fmt"
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
		s.WriteString(styles.InputStyle.Render(m.CommitMsg+"_") + "\n")

		charCount := len(m.CommitMsg)
		countText := fmt.Sprintf("  %d/72", charCount)
		if charCount > 72 {
			s.WriteString(styles.ErrorStyle.Render(countText))
		} else if charCount > 50 {
			s.WriteString(styles.WarningStyle.Render(countText))
		} else {
			s.WriteString(styles.HelpStyle.Render(countText))
		}
		s.WriteString("\n")
	}

	controlsWidget := controls.NewCommitViewControls(m.CopilotAvailable, m.IsGeneratingAI)
	s.WriteString(controlsWidget.Render())

	return s.String()
}
