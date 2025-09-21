package view

import (
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
)

func RenderAddRemoteView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.HeaderStyle.Render("➕ Add new remote:") + "\n\n")

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

	controlsWidget := controls.NewAddRemoteViewControls()
	s.WriteString(controlsWidget.Render())

	return s.String()
}
