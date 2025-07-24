package async

import (
	"time"

	"froggit/internal/git"
	"froggit/internal/tui/update/messages"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	PushMsg        struct{ Err error }
	FetchMsg       struct{ Err error }
	PullMsg        struct{ Err error }
	SpinnerTickMsg struct{}
)

// spinner returns a Cmd that emits spinnerTickMsg every 100ms.
func Spinner() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

// performPush runs git.Push asynchronously and returns a pushMsg.
func PerformPush() tea.Cmd {
	return func() tea.Msg {
		return PushMsg{Err: git.Push()}
	}
}

// performFetch runs git.Fetch asynchronously and returns a fetchMsg.
func PerformFetch() tea.Cmd {
	return func() tea.Msg {
		return FetchMsg{Err: git.Fetch()}
	}
}

// performPull runs git.Pull asynchronously and returns a pullMsg.
func PerformPull() tea.Cmd {
	return func() tea.Msg {
		return PullMsg{Err: git.Pull()}
	}
}

// performSwitchAndMerge switches to target branch and then merges source branch.
func PerformSwitchAndMerge(targetBranch, sourceBranch string) tea.Cmd {
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
func PerformSwitchAndRebase(sourceBranch, targetBranch string) tea.Cmd {
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
