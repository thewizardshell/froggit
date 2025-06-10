package view

import (
	"strings"

	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

// RenderAddRemoteView dibuja la vista para agregar un nuevo remoto
func RenderAddRemoteView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render("➕ Add new remote:") + "\n\n")

	// Campo Name
	nameLabel := "Name:"
	if m.InputField == "name" {
		nameLabel = " " + nameLabel
	} else {
		nameLabel = "  " + nameLabel
	}

	nameStyle := styles.NormalStyle
	if m.InputField == "name" {
		nameStyle = styles.InputStyle
	}

	s.WriteString(styles.HelpStyle.Render(nameLabel) + "\n")
	s.WriteString(nameStyle.Render(m.RemoteName+"_") + "\n\n")

	// Campo URL
	urlLabel := "URL:"
	if m.InputField == "url" {
		urlLabel = " " + urlLabel
	} else {
		urlLabel = "  " + urlLabel
	}

	urlStyle := styles.NormalStyle
	if m.InputField == "url" {
		urlStyle = styles.InputStyle
	}

	s.WriteString(styles.HelpStyle.Render(urlLabel) + "\n")
	s.WriteString(urlStyle.Render(m.RemoteURL+"_") + "\n\n")

	s.WriteString(styles.BorderStyle.Render(
		styles.HelpStyle.Render("Controls:\n") +
			styles.HelpStyle.Render("  [Tab] switch field  [Enter] confirm/next  [Esc] cancel"),
	))

	return s.String()
}
