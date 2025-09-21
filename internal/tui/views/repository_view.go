package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderRepositoryListView(m model.Model) string {
	const viewport = 12 // number of repos to show at once
	var sb strings.Builder

	total := len(m.Repositories)
	headerText := "📚 GitHub Repositories"
	if total > 0 {
		headerText += fmt.Sprintf(" (%d/%d)", m.SelectedRepoIndex+1, total)
	}
	sb.WriteString(styles.HeaderStyle.Render(headerText) + "\n\n")

	if total == 0 {
		sb.WriteString(styles.HelpStyle.Render("🔍 No repositories found.") + "\n")
		sb.WriteString(styles.HelpStyle.Render("Make sure you have access to GitHub repositories.") + "\n\n")
	} else {
		start := 0
		if m.SelectedRepoIndex >= viewport/2 {
			start = m.SelectedRepoIndex - viewport/2
		}
		if start+viewport > total {
			start = max(0, total-viewport)
		}
		end := min(total, start+viewport)

		for i := start; i < end; i++ {
			repo := m.Repositories[i]
			selected := m.SelectedRepoIndex == i

			cursor := "  "
			if selected {
				cursor = "❯ "
			}

			icon := "📂"

			repoName := styles.SuccessStyle.Render(repo.Owner.Login + "/" + repo.Name)
			line := cursor + icon + " " + repoName


			if selected {
				sb.WriteString(styles.SelectedStyle.Render(line) + "\n")
			} else {
				sb.WriteString(styles.NormalStyle.Render(line) + "\n")
			}

			if repo.Description != "" && selected {
				desc := "  " + styles.HelpStyle.Render("  └─ " + repo.Description)
				sb.WriteString(desc + "\n")
			}
		}
		sb.WriteString("\n")
	}

	controlsWidget := controls.NewRepositoryListControls()
	sb.WriteString(controlsWidget.Render())

	return sb.String()
}