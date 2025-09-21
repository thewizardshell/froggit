package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderConfirmCloneRepoView(m model.Model) string {
	var s strings.Builder
	repo := m.RepoToClone

	s.WriteString(styles.HeaderStyle.Render("📥 Clone Repository") + "\n\n")

	repoName := styles.SuccessStyle.Render(fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name))
	s.WriteString(styles.NormalStyle.Render("Repository: ") + repoName + "\n")

	if repo.Description != "" {
		s.WriteString(styles.HelpStyle.Render("Description: " + repo.Description) + "\n")
	}


	s.WriteString("\n")
	s.WriteString(styles.NormalStyle.Render("This will clone the repository to your current directory.") + "\n\n")

	controlsWidget := controls.NewConfirmDialogControls()
	s.WriteString(controlsWidget.Render())

	return s.String()
}
