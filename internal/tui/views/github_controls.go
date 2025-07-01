package view

import (
	"froggit/internal/tui/styles"
)

func RenderGitHubControlsView() string {
	return styles.BorderStyle.Render(
		styles.HeaderStyle.Render(" GitHub Controls ") + "\n\n" +
			styles.HelpStyle.Render("[r] List repositories\n") +
			styles.HelpStyle.Render("[esc] Back"),
	)
}
