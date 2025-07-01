package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderRepositoryListView(m model.Model) string {
	const viewport = 12 // number of repos to show at once
	var sb strings.Builder
	sb.WriteString(styles.TitleStyle.Render("  GitHub Repositories") + "\n\n")
	total := len(m.Repositories)
	if total == 0 {
		sb.WriteString(styles.HelpStyle.Render("No repositories found.\n"))
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
			icon := "" // folder icon
			name := repo.Name
			owner := repo.Owner.Login
			desc := repo.Description
			line := icon + " " + owner + "/" + name
			if desc != "" {
				line += "  " + styles.HelpStyle.Render(desc)
			}
			if selected {
				sb.WriteString(styles.SelectedStyle.Render("  "+line) + "\n")
			} else {
				sb.WriteString(styles.NormalStyle.Render("  "+line) + "\n")
			}
		}
	}
	controls := "[↑/↓] navigate   [c] clone   [esc] back"
	if total > 0 {
		controls += fmt.Sprintf("   %d/%d", m.SelectedRepoIndex+1, total)
	}
	sb.WriteString("\n" + styles.BorderStyle.Render(styles.HelpStyle.Render(controls)))
	return sb.String()
}
