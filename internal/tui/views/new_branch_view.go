package view

import (
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderNewBranchView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render("🌿 New Branch:") + "\n\n")
	s.WriteString(styles.InputStyle.Render(m.NewBranchName+"_") + "\n\n")

	controlsWidget := controls.NewNewBranchViewControls()
	s.WriteString(controlsWidget.Render())

	return s.String()
}
