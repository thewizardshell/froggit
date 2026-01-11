package async

import (
	"time"

	"froggit/internal/copilot"
	"froggit/internal/git"
	"froggit/internal/tui/update/messages"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	PushMsg               struct{ Err error }
	FetchMsg              struct{ Err error }
	PullMsg               struct{ Err error }
	SpinnerTickMsg        struct{}
	RemoteChangesCheckMsg struct{ HasChanges bool; Err error }
	AICommitMsg           struct{ Message string; Err error }
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

func PerformPushWithConfig(defaultBranch string) tea.Cmd {
	return func() tea.Msg {
		return PushMsg{Err: git.PushWithBranch(defaultBranch)}
	}
}

// performFetch runs git.Fetch asynchronously and returns a fetchMsg.
func PerformFetch() tea.Cmd {
	return func() tea.Msg {
		return FetchMsg{Err: git.Fetch()}
	}
}

func PerformAutoFetch() tea.Cmd {
	return func() tea.Msg {
		return FetchMsg{Err: git.FetchWithConfig(true)}
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

// PerformRemoteChangesCheck checks for remote changes without blocking startup
func PerformRemoteChangesCheck(branch string) tea.Cmd {
	return func() tea.Msg {
		hasChanges, err := git.HasRemoteChangesWithFetch(branch, false)
		return RemoteChangesCheckMsg{HasChanges: hasChanges, Err: err}
	}
}

// PerformAICommitGeneration generates a commit message using Copilot
func PerformAICommitGeneration() tea.Cmd {
	return func() tea.Msg {
		diff, err := git.GetStagedDiff()
		if err != nil {
			return AICommitMsg{Err: err}
		}
		msg, err := copilot.GenerateCommitMessage(diff)
		return AICommitMsg{Message: msg, Err: err}
	}
}
