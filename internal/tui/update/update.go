// Package update handles the core message loop for user interactions
// in the Froggit TUI application, using the Bubble Tea framework.
package update

import (
	"fmt"
	"strings"
	"time"

	"froggit/internal/git"
	"froggit/internal/tui/model"

	tea "github.com/charmbracelet/bubbletea"
)

// Internal messages for async operations
type pushMsg struct{ err error }
type fetchMsg struct{ err error }
type pullMsg struct{ err error }
type spinnerTickMsg struct{}

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

// Update handles Bubble Tea messages and returns the updated model
// along with the next command to execute.
func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.CurrentView == model.ConfirmDialog {
			switch msg.String() {
			case "y":
				switch m.DialogType {
				case "delete_branch":
					if err := git.DeleteBranch(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error deleting branch: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Branch deleted successfully"
						m.MessageType = "success"
						m.RefreshData()
					}
				case "discard_changes":
					if err := git.DiscardChanges(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error discarding changes: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Changes discarded"
						m.MessageType = "success"
						m.RefreshData()
					}
				}
				m.CurrentView = model.FileView
				return m, nil

			case "n", "esc":
				m.CurrentView = model.FileView
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			// Allow exiting with "q" only when we are NOT in an input-focused view.
			// When the user is typing (e.g. commit message, new branch name, adding remotes)
			// the character should be treated as normal text instead of a quit signal.
			if m.CurrentView == model.FileView || m.CurrentView == model.BranchView || m.CurrentView == model.RemoteView || m.CurrentView == model.ConfirmDialog || m.CurrentView == model.HelpView {
				return m, tea.Quit
			}
			// Otherwise, fall through so the rune is processed as normal input.

		case "esc":
			if m.CurrentView != model.FileView {
				m.CurrentView = model.FileView
				m.CommitMsg = ""
				m.RemoteName = ""
				m.RemoteURL = ""
				m.NewBranchName = ""
				m.Message = ""
				m.MessageType = ""
			}
			return m, nil

		case "enter":
			switch m.CurrentView {
			case model.CommitView:
				if m.CommitMsg != "" {
					if err := git.Commit(m.CommitMsg); err != nil {
						m.Message = fmt.Sprintf("✗ Error committing: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = "✓ Changes committed successfully"
						m.MessageType = "success"
						m.CurrentView = model.FileView
						m.CommitMsg = ""
						m.RefreshData()
					}
				}
				return m, nil

			case model.NewBranchView:
				if m.NewBranchName != "" {
					if err := git.CreateBranch(m.NewBranchName); err != nil {
						m.Message = fmt.Sprintf("✗ Error creating branch: %s", err)
						m.MessageType = "error"
					} else {
						m.Message = fmt.Sprintf("✓ Branch %s created successfully", m.NewBranchName)
						m.MessageType = "success"
						m.CurrentView = model.BranchView
						m.NewBranchName = ""
						m.RefreshData()
					}
				}
				return m, nil

			case model.BranchView:
				if len(m.Branches) > 0 {
					if m.Cursor >= len(m.Branches) {
						m.Cursor = len(m.Branches) - 1
					}
					selected := m.Branches[m.Cursor]
					if selected != m.CurrentBranch {
						if err := git.Checkout(selected); err != nil {
							m.Message = fmt.Sprintf("✗ Error switching to branch %s: %s", selected, err)
							m.MessageType = "error"
						} else {
							m.Message = fmt.Sprintf("✓ Switched to branch %s", selected)
							m.MessageType = "success"
							m.CurrentBranch = selected
							m.RefreshData()
						}
					} else {
						m.Message = "⚠ You are already on this branch"
						m.MessageType = "info"
					}
				}
				return m, nil
			}

		case "backspace":
			switch m.CurrentView {
			case model.CommitView:
				if len(m.CommitMsg) > 0 {
					m.CommitMsg = m.CommitMsg[:len(m.CommitMsg)-1]
				}
			case model.NewBranchView:
				if len(m.NewBranchName) > 0 {
					m.NewBranchName = m.NewBranchName[:len(m.NewBranchName)-1]
				}
			}
			return m, nil

		case "up":
			if m.CurrentView != model.CommitView && m.CurrentView != model.NewBranchView && m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		case "down":
			if m.CurrentView != model.CommitView && m.CurrentView != model.NewBranchView {
				switch m.CurrentView {
				case model.FileView:
					if m.Cursor < len(m.Files)-1 {
						m.Cursor++
					}
				case model.BranchView:
					if m.Cursor < len(m.Branches)-1 {
						m.Cursor++
					}
				case model.RemoteView:
					if m.Cursor < len(m.Remotes)-1 {
						m.Cursor++
					}
				case model.LogGraphView:
					if m.Cursor < len(m.LogLines)-1 {
						m.Cursor++
					}
				}
			}
			return m, nil
		}

		if m.CurrentView == model.CommitView {
			if len(msg.Runes) == 1 && isPrintableChar(msg.Runes[0]) {
				m.CommitMsg += string(msg.Runes)
				return m, nil
			}
		}

		if m.CurrentView == model.NewBranchView {
			if len(msg.Runes) == 1 && isPrintableChar(msg.Runes[0]) {
				m.NewBranchName += string(msg.Runes)
				return m, nil
			}
		}

		if m.CurrentView == model.BranchView {
			switch msg.String() {
			case "n":
				m.CurrentView = model.NewBranchView
				m.NewBranchName = ""
				return m, nil
			case "d":
				if len(m.Branches) > 0 && m.Cursor < len(m.Branches) {
					toDel := m.Branches[m.Cursor]
					if toDel == m.CurrentBranch {
						m.Message = "✗ Cannot delete current branch"
						m.MessageType = "error"
					} else {
						m.DialogType = "delete_branch"
						m.DialogTarget = toDel
						m.CurrentView = model.ConfirmDialog
					}
				}
				return m, nil
			}
		}

		if m.CurrentView == model.FileView {
			switch msg.String() {
			case " ":
				if len(m.Files) > 0 {
					f := &m.Files[m.Cursor]
					f.Staged = !f.Staged
					if f.Staged {
						git.Add(f.Name)
						m.Message = fmt.Sprintf("✓ File %s added to stage", f.Name)
					} else {
						git.Reset(f.Name)
						m.Message = fmt.Sprintf("✓ File %s removed from stage", f.Name)
					}
					m.MessageType = "success"
				}
			case "a":
				for i := range m.Files {
					if !m.Files[i].Staged {
						m.Files[i].Staged = true
						git.Add(m.Files[i].Name)
					}
				}
				m.Message = "✓ All files added to stage"
				m.MessageType = "success"
			case "c":
				ok := false
				for _, f := range m.Files {
					if f.Staged {
						ok = true
						break
					}
				}
				if ok {
					m.CurrentView = model.CommitView
					m.Message = ""
				} else {
					m.Message = "⚠ No staged files to commit"
					m.MessageType = "error"
				}
			case "b":
				m.CurrentView = model.BranchView
				m.Cursor = 0
				m.Message = ""
			case "m":
				m.CurrentView = model.RemoteView
				m.Cursor = 0
				m.Message = ""
			case "r":
				m.RefreshData()
				m.Message = "✓ Status updated"
				m.MessageType = "success"
			case "p":
				if !m.IsPushing {
					m.IsPushing = true
					m.Message = "Pushing..."
					m.MessageType = "info"
					return m, tea.Batch(performPush(), spinner())
				}
			case "f":
				if !m.IsFetching {
					m.IsFetching = true
					m.Message = "Fetching..."
					m.MessageType = "info"
					return m, tea.Batch(performFetch(), spinner())
				}
			case "l":
				if !m.IsPulling {
					m.IsPulling = true
					m.Message = "Pulling..."
					m.MessageType = "info"
					return m, tea.Batch(performPull(), spinner())
				}
			case "L":
				// Open interactive Git log graph view
				graph, err := git.LogsGraph()
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error retrieving log graph: %s", err)
					m.MessageType = "error"
				} else {
					m.LogLines = strings.Split(strings.TrimSpace(graph), "\n")
					m.Cursor = 0
					m.CurrentView = model.LogGraphView
					m.Message = ""
				}
				return m, nil
			case "A":
				m.Message = "Advanced features (logs, merge, stash, rebase) are coming soon"
				m.MessageType = "info"
				return m, nil
			case "x":
				if len(m.Files) > 0 {
					m.DialogType = "discard_changes"
					m.DialogTarget = m.Files[m.Cursor].Name
					m.CurrentView = model.ConfirmDialog
				}
			case "?":
				if m.CurrentView == model.FileView {
					m.CurrentView = model.HelpView
					return m, nil
				} else if m.CurrentView == model.HelpView {
					m.CurrentView = model.FileView
					return m, nil
				}
			}
		}

		if m.CurrentView == model.LogGraphView && msg.String() == "esc" {
			m.CurrentView = model.FileView
			return m, nil
		}

	case spinnerTickMsg:
		if m.IsPushing || m.IsFetching || m.IsPulling {
			m.SpinnerIndex = (m.SpinnerIndex + 1) % len(m.SpinnerFrames)
			return m, spinner()
		}

	case pushMsg:
		m.IsPushing = false
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error pushing changes: %s", msg.err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pushed successfully"
			m.MessageType = "success"
			m.RefreshData()
		}

	case fetchMsg:
		m.IsFetching = false
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error fetching changes: %s", msg.err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes fetched successfully"
			m.MessageType = "success"
			m.RefreshData()
		}

	case pullMsg:
		m.IsPulling = false
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error pulling changes: %s", msg.err)
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pulled successfully"
			m.MessageType = "success"
			m.RefreshData()
		}
	}

	return m, nil
}

// isPrintableChar checks whether a rune is a printable ASCII or extended character.
func isPrintableChar(r rune) bool {
	return (r >= 32 && r <= 126) || (r >= 128 && r <= 255)
}
