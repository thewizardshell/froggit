package view

import (
	"strings"

	"froggit/internal/tui/styles"
)

// RenderHelpView shows the extended list of controls.
func RenderHelpView() string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render("Additional Controls:") + "\n\n")

	lines := []string{
		"[a] stage all",
		"[b] branches",
		"[m] remotes",
		"[p] push",
		"[f] fetch",
		"[l] pull (only when remote changes)",
		"[L] log graph",
		"[r] refresh",
		"[A] advanced (logs, merge, stash, rebase)",
		"[q] quit",
		"[esc] back",
	}

	for _, l := range lines {
		s.WriteString(styles.HelpStyle.Render("  "+l) + "\n")
	}

	return s.String()
}
