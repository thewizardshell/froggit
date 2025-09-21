package tui

import (
	"fmt"
	"froggit/internal/config"
	"froggit/internal/tui/branding"
	"froggit/internal/tui/model"
	"froggit/internal/tui/styles"
	view "froggit/internal/tui/views"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func Render(m model.Model, cfg config.Config) string {
	var sb strings.Builder

	if cfg.Ui.Branding && m.CurrentView != model.QuickStartView {
		sb.WriteString(styles.TitleStyle.Render(branding.RenderTitle()) + "\n\n")
	}

	if m.CurrentView != model.QuickStartView {
		sb.WriteString(fmt.Sprintf(" current branch: %s\n\n",
			styles.HeaderStyle.Render(m.CurrentBranch),
		))
	}

	switch m.CurrentView {
	case model.QuickStartView:
		sb.WriteString(view.RenderQuickStartView(m))
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
	case model.MergeView:
		sb.WriteString(view.RenderMergeView(m))
	case model.RebaseView:
		sb.WriteString(view.RenderRebaseView(m))
	case model.StashView:
		sb.WriteString(view.RenderStashView(m))
	case model.StashMessageView:
		sb.WriteString(view.RenderStashMessageView(m))
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

	if m.IsFetching {
		sb.WriteString("\n" + styles.SpinnerStyle.Render(
			fmt.Sprintf(" Fetching... %s", m.SpinnerFrames[m.SpinnerIndex]),
		))
	}

	if m.IsPulling {
		sb.WriteString("\n" + styles.SpinnerStyle.Render(
			fmt.Sprintf(" Pulling... %s", m.SpinnerFrames[m.SpinnerIndex]),
		))
	}

	content := sb.String()

	switch strings.ToLower(cfg.Ui.Position) {
	case "center":
		return renderCentered(content)
	case "right":
		return renderRight(content)
	default:
		return content
	}
}

func renderCentered(content string) string {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 25
	}

	lines := strings.Split(content, "\n")
	numLines := len(lines)

	topPadding := (height - numLines) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	padding := strings.Repeat("\n", topPadding)
	content = padding + content

	centerStyle := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	return centerStyle.Render(content)
}

func renderRight(content string) string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80
	}

	rightStyle := lipgloss.NewStyle().Width(width).Align(lipgloss.Right)
	return rightStyle.Render(content)
}
