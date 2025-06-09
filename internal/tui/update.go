package tui

import (
	"fmt"
	"time"

	"giteasy/internal/git"

	tea "github.com/charmbracelet/bubbletea"
)

type pushMsg struct {
	err error
}

type fetchMsg struct {
	err error
}

type pullMsg struct {
	err error
}

type spinnerTickMsg struct{}

func spinner() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

func performPush() tea.Cmd {
	return func() tea.Msg {
		err := git.Push()
		return pushMsg{err: err}
	}
}

func performFetch() tea.Cmd {
	return func() tea.Msg {
		err := git.Fetch()
		return fetchMsg{err: err}
	}
}

func performPull() tea.Cmd {
	return func() tea.Msg {
		err := git.Pull()
		return pullMsg{err: err}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.CurrentView == ConfirmDialog {
			switch msg.String() {
			case "y":
				switch m.DialogType {
				case "delete_branch":
					if err := git.DeleteBranch(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error deleting branch: %s", err.Error())
						m.MessageType = "error"
					} else {
						m.Message = "✓ Branch deleted successfully"
						m.MessageType = "success"
						m.RefreshData()
					}
				case "discard_changes":
					if err := git.DiscardChanges(m.DialogTarget); err != nil {
						m.Message = fmt.Sprintf("✗ Error discarding changes: %s", err.Error())
						m.MessageType = "error"
					} else {
						m.Message = "✓ Changes discarded"
						m.MessageType = "success"
						m.RefreshData()
					}
				}
				m.CurrentView = FileView
				return m, nil
			case "n", "esc":
				m.CurrentView = FileView
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.CurrentView != FileView {
				m.CurrentView = FileView
				m.CommitMsg = ""
				m.RemoteName = ""
				m.RemoteURL = ""
				m.NewBranchName = ""
				m.Message = ""
				return m, nil
			}
		case "enter":
			if m.CurrentView == CommitView && m.CommitMsg != "" {
				err := git.Commit(m.CommitMsg)
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error committing: %s", err.Error())
					m.MessageType = "error"
				} else {
					m.Message = "✓ Changes committed successfully"
					m.MessageType = "success"
					m.CurrentView = FileView
					m.CommitMsg = ""
					m.RefreshData()
				}
				return m, nil
			} else if m.CurrentView == NewBranchView && m.NewBranchName != "" {
				err := git.CreateBranch(m.NewBranchName)
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error creating branch: %s", err.Error())
					m.MessageType = "error"
				} else {
					m.Message = fmt.Sprintf("✓ Branch %s created successfully", m.NewBranchName)
					m.MessageType = "success"
					m.RefreshData()
					m.CurrentView = BranchView
					m.NewBranchName = ""
				}
				return m, nil
			} else if m.CurrentView == BranchView && len(m.Branches) > 0 {
				if m.Cursor >= len(m.Branches) {
					m.Cursor = len(m.Branches) - 1
					return m, nil
				}
				selectedBranch := m.Branches[m.Cursor]
				if selectedBranch != m.CurrentBranch {
					err := git.Checkout(selectedBranch)
					if err != nil {
						m.Message = fmt.Sprintf("✗ Error switching to branch %s: %s", selectedBranch, err.Error())
						m.MessageType = "error"
					} else {
						m.Message = fmt.Sprintf("✓ Switched to branch %s", selectedBranch)
						m.MessageType = "success"
						m.CurrentBranch = selectedBranch
						m.RefreshData()
					}
				} else {
					m.Message = "⚠ You are already on this branch"
					m.MessageType = "info"
				}
				return m, nil
			}
		case "backspace":
			if m.CurrentView == CommitView && len(m.CommitMsg) > 0 {
				m.CommitMsg = m.CommitMsg[:len(m.CommitMsg)-1]
				return m, nil
			} else if m.CurrentView == NewBranchView && len(m.NewBranchName) > 0 {
				m.NewBranchName = m.NewBranchName[:len(m.NewBranchName)-1]
				return m, nil
			}
		case "up":
			if m.CurrentView != CommitView && m.CurrentView != NewBranchView {
				if m.Cursor > 0 {
					m.Cursor--
				}
			}
		case "down":
			if m.CurrentView != CommitView && m.CurrentView != NewBranchView {
				if m.CurrentView == FileView && m.Cursor < len(m.Files)-1 {
					m.Cursor++
				} else if m.CurrentView == BranchView && m.Cursor < len(m.Branches)-1 {
					m.Cursor++
				} else if m.CurrentView == RemoteView && m.Cursor < len(m.Remotes)-1 {
					m.Cursor++
				}
			}
		default:
			if m.CurrentView == CommitView {
				if len(msg.Runes) == 1 && isPrintableChar(msg.Runes[0]) {
					m.CommitMsg += string(msg.Runes)
					return m, nil
				}
			} else if m.CurrentView == NewBranchView {
				if len(msg.Runes) == 1 && isPrintableChar(msg.Runes[0]) {
					m.NewBranchName += string(msg.Runes)
					return m, nil
				}
			} else if m.CurrentView == BranchView {
				switch msg.String() {
				case "n":
					m.CurrentView = NewBranchView
					m.NewBranchName = ""
					return m, nil
				case "d":
					if len(m.Branches) > 0 && m.Cursor < len(m.Branches) {
						branchToDelete := m.Branches[m.Cursor]
						if branchToDelete == m.CurrentBranch {
							m.Message = "✗ Cannot delete current branch"
							m.MessageType = "error"
						} else {
							m.DialogType = "delete_branch"
							m.DialogTarget = branchToDelete
							m.CurrentView = ConfirmDialog
						}
					}
				}
			} else if m.CurrentView == FileView {
				switch msg.String() {
				case " ":
					if len(m.Files) > 0 {
						m.Files[m.Cursor].Staged = !m.Files[m.Cursor].Staged
						if m.Files[m.Cursor].Staged {
							git.Add(m.Files[m.Cursor].Name)
							m.Message = fmt.Sprintf("✓ File %s added to stage", m.Files[m.Cursor].Name)
							m.MessageType = "success"
						} else {
							git.Reset(m.Files[m.Cursor].Name)
							m.Message = fmt.Sprintf("✓ File %s removed from stage", m.Files[m.Cursor].Name)
							m.MessageType = "success"
						}
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
					hasStaged := false
					for _, file := range m.Files {
						if file.Staged {
							hasStaged = true
							break
						}
					}
					if hasStaged {
						m.CurrentView = CommitView
						m.Message = ""
					} else {
						m.Message = "⚠ No staged files to commit"
						m.MessageType = "error"
					}
				case "b":
					m.CurrentView = BranchView
					m.Cursor = 0
					m.Message = ""
				case "m":
					m.CurrentView = RemoteView
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
					if m.CurrentView == FileView && !m.IsFetching {
						m.IsFetching = true
						m.Message = "Fetching..."
						m.MessageType = "info"
						return m, tea.Batch(performFetch(), spinner())
					}
				case "l":
					if m.CurrentView == FileView && !m.IsPulling {
						m.IsPulling = true
						m.Message = "Pulling..."
						m.MessageType = "info"
						return m, tea.Batch(performPull(), spinner())
					}
				case "x":
					if len(m.Files) > 0 {
						m.DialogType = "discard_changes"
						m.DialogTarget = m.Files[m.Cursor].Name
						m.CurrentView = ConfirmDialog
						return m, nil
					}
				}
			}
		}
	case spinnerTickMsg:
		if m.IsPushing || m.IsFetching || m.IsPulling {
			m.SpinnerIndex = (m.SpinnerIndex + 1) % len(m.SpinnerFrames)
			return m, spinner()
		}
	case pushMsg:
		m.IsPushing = false
		if msg.err != nil {
			m.Message = fmt.Sprintf("✗ Error pushing changes: %s", msg.err.Error())
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pushed successfully"
			m.MessageType = "success"
			m.RefreshData()
		}
	case fetchMsg:
		m.IsFetching = false
		if msg.err != nil {
			m.Message = "✗ Error fetching changes: " + msg.err.Error()
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes fetched successfully"
			m.MessageType = "success"
			m.RefreshData()
		}
	case pullMsg:
		m.IsPulling = false
		if msg.err != nil {
			m.Message = "✗ Error pulling changes: " + msg.err.Error()
			m.MessageType = "error"
		} else {
			m.Message = "✓ Changes pulled successfully"
			m.MessageType = "success"
			m.RefreshData()
		}
	}
	return m, nil
}

func isPrintableChar(r rune) bool {
	return r >= 32 && r <= 126 || r >= 128 && r <= 255
}
