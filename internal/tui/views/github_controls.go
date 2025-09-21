package view

import (
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/styles"
)

func RenderGitHubControlsView() string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render("🐙 GitHub Integration") + "\n\n")

	s.WriteString(styles.NormalStyle.Render("Connect and manage your GitHub repositories directly from Froggit.") + "\n")
	s.WriteString(styles.HelpStyle.Render("Browse, clone, and work with your GitHub projects seamlessly.") + "\n\n")

	features := []string{
		"📋 Browse your repositories",
		"📥 Clone projects locally",
		"🔄 Sync with remote changes",
	}

	for _, feature := range features {
		s.WriteString(styles.NormalStyle.Render("  " + feature) + "\n")
	}

	s.WriteString("\n")

	controlsWidget := controls.NewGitHubControlsViewControls()
	s.WriteString(controlsWidget.Render())

	return s.String()
}
