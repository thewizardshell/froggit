package update

import (
	"time"

	"froggit/internal/git"
	"froggit/internal/tui/update/messages"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	pushMsg        struct{ err error }
	fetchMsg       struct{ err error }
	pullMsg        struct{ err error }
	spinnerTickMsg struct{}
)

// spinner returns a Cmd that emits spinnerTickMsg every 100ms.
func spinner() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

// performPush runs git.Push asynchronously and returns a pushMsg.
func performPush() tea.Cmd {
	return func() tea.Msg {
		return pushMsg{err: git.Push()}
	}
}

// performFetch runs git.Fetch asynchronously and returns a fetchMsg.
func performFetch() tea.Cmd {
	return func() tea.Msg {
		return fetchMsg{err: git.Fetch()}
	}
}

// performPull runs git.Pull asynchronously and returns a pullMsg.
func performPull() tea.Cmd {
	return func() tea.Msg {
		return pullMsg{err: git.Pull()}
	}
}

// performSwitchAndMerge switches to target branch and then merges source branch.
func performSwitchAndMerge(targetBranch, sourceBranch string) tea.Cmd {
	return func() tea.Msg {
		err := git.Checkout(targetBranch)
		return messages.SwitchBranchMsg{
			Err:          err,
			TargetBranch: targetBranch,
			NextAction:   "merge",
			SourceBranch: sourceBranch,
		}
	}
}

// performSwitchAndRebase switches to source branch and then rebases onto target.
func performSwitchAndRebase(sourceBranch, targetBranch string) tea.Cmd {
	return func() tea.Msg {
		err := git.Checkout(sourceBranch)
		return messages.SwitchBranchMsg{
			Err:          err,
			TargetBranch: targetBranch,
			NextAction:   "rebase",
			SourceBranch: sourceBranch,
		}
	}
}
