package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var s strings.Builder

	asciiTitle := `
█████████████╗ ██████╗ ██████╗ ██████╗██████████╗
██╔════██╔══████╔═══████╔════╝██╔════╝██╚══██╔══╝
█████╗ ██████╔██║   ████║  █████║  █████║  ██║   
██╔══╝ ██╔══████║   ████║   ████║   ████║  ██║   
██║    ██║  ██╚██████╔╚██████╔╚██████╔██║  ██║   
╚═╝    ╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚═════╝╚═╝  ╚═╝  𓆏  
`

	s.WriteString(TitleStyle.Render(asciiTitle) + "\n\n")
	s.WriteString(fmt.Sprintf(" current branch: %s\n\n", HeaderStyle.Render(m.CurrentBranch)))

	switch m.CurrentView {
	case FileView:
		s.WriteString(m.renderFileView())
	case CommitView:
		s.WriteString(m.renderCommitView())
	case BranchView:
		s.WriteString(m.renderBranchView())
	case RemoteView:
		s.WriteString(m.renderRemoteView())
	case AddRemoteView:
		s.WriteString(m.renderAddRemoteView())
	case NewBranchView:
		s.WriteString(m.renderNewBranchView())
	case ConfirmDialog:
		s.WriteString(m.renderConfirmDialog())
	}

	if m.Message != "" {
		s.WriteString("\n")
		switch m.MessageType {
		case "error":
			s.WriteString(ErrorStyle.Render(m.Message))
		case "success":
			s.WriteString(SuccessStyle.Render(m.Message))
		default:
			s.WriteString(NormalStyle.Render(m.Message))
		}
	}

	if m.IsFetching {
		s.WriteString("\n" + SpinnerStyle.Render(fmt.Sprintf(" Fetching... %s", m.SpinnerFrames[m.SpinnerIndex])))
	}
	if m.IsPulling {
		s.WriteString("\n" + SpinnerStyle.Render(fmt.Sprintf(" Pulling... %s", m.SpinnerFrames[m.SpinnerIndex])))
	}

	return s.String()
}

func (m Model) renderFileView() string {
	var s strings.Builder

	stagedCount := 0
	unstagedCount := 0
	for _, file := range m.Files {
		if file.Staged {
			stagedCount++
		} else {
			unstagedCount++
		}
	}

	s.WriteString(HeaderStyle.Render(" Git Status:") + "\n")
	s.WriteString(fmt.Sprintf(" Stage: %d files\n", stagedCount))
	s.WriteString(fmt.Sprintf(" Unstaged: %d files\n", unstagedCount))
	s.WriteString("\n")
	s.WriteString(HeaderStyle.Render(" Modified files:") + "\n\n")

	if len(m.Files) == 0 {
		s.WriteString(HelpStyle.Render("No modified files\n"))
	} else {
		for i, file := range m.Files {
			cursor := "  "
			if m.Cursor == i {
				cursor = ""
			}

			staged := " "
			if file.Staged {
				staged = "✓"
			}

			style := NormalStyle
			if m.Cursor == i {
				style = SelectedStyle
			}

			icon := GetIconForFile(file.Name)
			line := fmt.Sprintf("%s [%s] %s %s", cursor, staged, icon, file.Name)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	s.WriteString("\n" + BorderStyle.Render(
		HelpStyle.Render("Controls:\n")+
			HelpStyle.Render("  [↑/↓] navigate  [space] stage/unstage  [a] stage all  [x] discard changes")+
			HelpStyle.Render("  [c] commit  [b] branches  [m] remotes  [p] push  [f] fetch  [l] pull  [r] refresh  [q] quit"),
	))

	return s.String()
}

func (m Model) renderCommitView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render(" Commit message:") + "\n\n")
	s.WriteString(InputStyle.Render(m.CommitMsg+"_") + "\n\n")

	s.WriteString(BorderStyle.Render(
		HelpStyle.Render("Type your message and press [Enter] to confirm\n") +
			HelpStyle.Render("[Esc] to cancel"),
	))

	return s.String()
}

func (m Model) renderBranchView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("Branches:") + "\n\n")

	for i, branch := range m.Branches {
		cursor := "  "
		if m.Cursor == i {
			cursor = ""
		}

		current := " "
		if branch == m.CurrentBranch {
			current = "●"
		}

		style := NormalStyle
		if m.Cursor == i {
			style = SelectedStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, current, branch)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n" + BorderStyle.Render(
		HelpStyle.Render("Controls:\n")+
			HelpStyle.Render("  [↑/↓] navigate  [Enter] switch branch  [n] new branch  [d] delete branch  [Esc] back"),
	))

	return s.String()
}

func (m Model) renderRemoteView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render(" Remote repositories:") + "\n\n")

	if len(m.Remotes) == 0 {
		s.WriteString(HelpStyle.Render("No remote repositories configured\n"))
	} else {
		for i, remote := range m.Remotes {
			cursor := "  "
			if m.Cursor == i {
				cursor = " "
			}

			style := NormalStyle
			if m.Cursor == i {
				style = SelectedStyle
			}

			line := fmt.Sprintf("%s %s", cursor, remote)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	s.WriteString("\n" + BorderStyle.Render(
		HelpStyle.Render("Controls:\n")+
			HelpStyle.Render("  [↑/↓] navigate  [n] new remote  [d] delete  [Esc] back"),
	))

	return s.String()
}

func (m Model) renderAddRemoteView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("➕ Add new remote:") + "\n\n")

	nameLabel := "Name:"
	if m.InputField == "name" {
		nameLabel = " " + nameLabel
	} else {
		nameLabel = "  " + nameLabel
	}

	nameStyle := NormalStyle
	if m.InputField == "name" {
		nameStyle = InputStyle
	}

	s.WriteString(HelpStyle.Render(nameLabel) + "\n")
	s.WriteString(nameStyle.Render(m.RemoteName+"_") + "\n\n")

	urlLabel := "URL:"
	if m.InputField == "url" {
		urlLabel = " " + urlLabel
	} else {
		urlLabel = "  " + urlLabel
	}

	urlStyle := NormalStyle
	if m.InputField == "url" {
		urlStyle = InputStyle
	}

	s.WriteString(HelpStyle.Render(urlLabel) + "\n")
	s.WriteString(urlStyle.Render(m.RemoteURL+"_") + "\n\n")

	s.WriteString(BorderStyle.Render(
		HelpStyle.Render("Controls:\n") +
			HelpStyle.Render("  [Tab] switch field  [Enter] confirm/next  [Esc] cancel"),
	))

	return s.String()
}

func (m Model) renderNewBranchView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("🌿 New Branch:") + "\n\n")
	s.WriteString(InputStyle.Render(m.NewBranchName+"_") + "\n\n")

	s.WriteString(BorderStyle.Render(
		HelpStyle.Render("Type the branch name and press [Enter] to create\n") +
			HelpStyle.Render("[Esc] to cancel"),
	))

	return s.String()
}

func (m Model) renderConfirmDialog() string {
	var s strings.Builder
	var message string

	switch m.DialogType {
	case "delete_branch":
		message = fmt.Sprintf("Are you sure you want to delete branch '%s'?", m.DialogTarget)
	case "discard_changes":
		message = fmt.Sprintf("Are you sure you want to discard changes in '%s'?", m.DialogTarget)
	}

	s.WriteString("\n\n")
	s.WriteString(BorderStyle.Render(
		HeaderStyle.Render(" Confirm action") + "\n\n" +
			NormalStyle.Render(message) + "\n\n" +
			HelpStyle.Render("[y] Yes  [n] No"),
	))

	return s.String()
}
