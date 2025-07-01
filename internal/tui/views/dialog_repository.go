package view

import (
	"fmt"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderConfirmCloneRepoView(m model.Model) string {
	repo := m.RepoToClone
	msg := fmt.Sprintf("Clone repository '%s/%s'?", repo.Owner.Login, repo.Name)
	return styles.BorderStyle.Render(
		styles.HeaderStyle.Render(" Confirm clone") + "\n\n" +
			styles.NormalStyle.Render(msg) + "\n\n" +
			styles.HelpStyle.Render("[y] Yes  [n] No"),
	)
}
