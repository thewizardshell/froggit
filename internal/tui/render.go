package tui

import (
	"fmt"
	"strings"

	"froggit/internal/tui/branding"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
	view "froggit/internal/tui/views"
)

// Render builds all the TUI output given a model.Model.
// You call it from your main: fmt.Println(tui.Render(m))
func Render(m model.Model) string {
	var sb strings.Builder

	sb.WriteString(styles.TitleStyle.Render(branding.RenderTitle()) + "\n\n")
	sb.WriteString(fmt.Sprintf(" current branch: %s\n\n",
		styles.HeaderStyle.Render(m.CurrentBranch),
	))

	// Selección de vista
	switch m.CurrentView {
	case model.FileView:
		sb.WriteString(view.RenderFileView(m))
	case model.CommitView:
		sb.WriteString(view.RenderCommitView(m))
	case model.BranchView:
		sb.WriteString(view.RenderBranchView(m))
	case model.RemoteView:
		sb.WriteString(view.RenderRemoteView(m))
	case model.AddRemoteView:
		sb.WriteString(view.RenderAddRemoteView(m))
	case model.NewBranchView:
		sb.WriteString(view.RenderNewBranchView(m))
	case model.ConfirmDialog:
		sb.WriteString(view.RenderConfirmDialog(m))
	case model.HelpView:
		sb.WriteString(view.RenderHelpView())
	case model.LogGraphView:
		sb.WriteString(view.RenderLogGraphView(m))
	case model.RepositoryListView:
		sb.WriteString(view.RenderRepositoryListView(m))
	case model.ConfirmCloneRepoView:
		sb.WriteString(view.RenderConfirmCloneRepoView(m))
	case model.GitHubControlsView:
		sb.WriteString(view.RenderGitHubControlsView())
	}

	if m.Message != "" {
		sb.WriteString("\n")
		switch m.MessageType {
		case "error":
			sb.WriteString(styles.ErrorStyle.Render(m.Message))
		case "success":
			sb.WriteString(styles.SuccessStyle.Render(m.Message))
		default:
			sb.WriteString(styles.NormalStyle.Render(m.Message))
		}
	}

	// Spinners
	if m.IsFetching {
		sb.WriteString("\n" + styles.SpinnerStyle.Render(
			fmt.Sprintf(" Fetching... %s", m.SpinnerFrames[m.SpinnerIndex]),
		))
	}
	if m.IsPulling {
		sb.WriteString("\n" + styles.SpinnerStyle.Render(
			fmt.Sprintf(" Pulling... %s", m.SpinnerFrames[m.SpinnerIndex]),
		))
	}

	return sb.String()
}
