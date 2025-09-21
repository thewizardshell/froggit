package view

import (
	"fmt"
	"strings"

	"froggit/internal/tui/controls"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

func RenderQuickStartView(m model.Model) string {
	var s strings.Builder

	s.WriteString(styles.TitleStyle.Render("▄▄▄ Froggit Quick Start ▄▄▄") + "\n\n")

	s.WriteString(styles.NormalStyle.Render("Welcome to Froggit! You're not in a Git repository yet.") + "\n")
	s.WriteString(styles.HelpStyle.Render("Choose one of the options below to get started:") + "\n\n")

	options := []struct {
		number      string
		title       string
		description string
		icon        string
		available   bool
	}{
		{
			number:      "1",
			title:       "Initialize New Repository",
			description: "Create a new Git repository in the current directory",
			icon:        "[+]",
			available:   true,
		},
		{
			number:      "2",
			title:       "Clone from GitHub",
			description: "Browse and clone repositories from your GitHub account",
			icon:        "[<]",
			available:   m.HasGitHubCLI,
		},
		{
			number:      "3",
			title:       "Create on GitHub",
			description: "Create a new repository on GitHub and clone it locally",
			icon:        "[^]",
			available:   m.HasGitHubCLI,
		},
	}

	for i, option := range options {
		cursor := "  "
		if m.Cursor == i {
			cursor = "> "
		}

		var optionStyle, titleStyle, descStyle lipgloss.Style

		if !option.available {
			optionStyle = styles.HelpStyle
			titleStyle = styles.HelpStyle
			descStyle = styles.HelpStyle
		} else if m.Cursor == i {
			optionStyle = styles.SelectedStyle
			titleStyle = styles.SuccessStyle
			descStyle = styles.NormalStyle
		} else {
			optionStyle = styles.NormalStyle
			titleStyle = styles.HeaderStyle
			descStyle = styles.HelpStyle
		}

		numberPart := titleStyle.Render(option.number + ")")
		iconPart := optionStyle.Render(option.icon)
		titlePart := titleStyle.Render(option.title)

		line := fmt.Sprintf("%s%s %s %s", cursor, numberPart, iconPart, titlePart)

		if m.Cursor == i {
			s.WriteString(styles.SelectedStyle.Render(line) + "\n")
		} else {
			s.WriteString(optionStyle.Render(line) + "\n")
		}

		if m.Cursor == i || !option.available {
			desc := "     " + descStyle.Render("└─ " + option.description)
			if !option.available {
				desc += descStyle.Render(" (GitHub CLI required)")
			}
			s.WriteString(desc + "\n")
		}
		s.WriteString("\n")
	}

	if !m.HasGitHubCLI {
		s.WriteString(styles.WarningStyle.Render("! GitHub CLI not found. Install 'gh' to enable GitHub features.") + "\n\n")
	}

	controlsWidget := controls.NewQuickStartControls()
	s.WriteString(controlsWidget.Render())

	return s.String()
}