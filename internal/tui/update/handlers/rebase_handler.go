package handlers

import (
	"fmt"
	"froggit/internal/git"
	"froggit/internal/tui/model"

	tea "github.com/charmbracelet/bubbletea"
)

// HandleRebaseView processes key messages in the rebase view.
func HandleRebaseView(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	if m.CurrentView == model.RebaseView {
		switch msg.String() {
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil
		case "down":
			if m.Cursor < len(m.Branches)-1 {
				m.Cursor++
			}
			return m, nil
		case " ", "space":
			if len(m.Branches) > 0 && m.Cursor < len(m.Branches) {
				selected := m.Branches[m.Cursor]
				if selected == m.CurrentBranch {
					m.Message = "⚠ Cannot rebase a branch onto itself"
					m.MessageType = "warning"
					return m, nil
				}
				if m.DialogTarget == selected {
					m.DialogTarget = ""
					m.Message = "Selection unmarked"
					m.MessageType = "info"
				} else {
					m.DialogTarget = selected
					m.Message = fmt.Sprintf("Will rebase %s onto %s", m.CurrentBranch, selected)
					m.MessageType = "info"
				}
			}
			return m, nil
		case "R", "r":
			if m.DialogTarget != "" {
				if m.DialogTarget == m.CurrentBranch {
					m.Message = "⚠ Cannot rebase a branch onto itself"
					m.MessageType = "warning"
					return m, nil
				}

				m.Message = fmt.Sprintf("Rebasing %s onto %s...", m.CurrentBranch, m.DialogTarget)
				m.MessageType = "info"

				err := git.Rebase(m.DialogTarget)
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error rebasing: %s", err)
					m.MessageType = "error"
					return m, nil
				}
				conflicts, _ := git.GetConflictFiles()
				if len(conflicts) > 0 {
					m.LogLines = conflicts
					m.Message = "Conflicts detected. Please resolve them and use [P] Proceed or [X] Cancel."
					m.MessageType = "warning"
					return m, nil
				} else {
					m.Message = fmt.Sprintf("✓ Rebase of %s onto %s successful", m.CurrentBranch, m.DialogTarget)
					m.MessageType = "success"
					m.CurrentView = model.FileView
					m.DialogTarget = ""
					m.LogLines = nil
					m.RefreshData()
					return m, nil
				}
			} else {
				m.Message = "Select a target branch first by pressing space"
				m.MessageType = "warning"
				return m, nil
			}
		case "P", "p":
			if len(m.LogLines) > 0 {
				err := git.RebaseContinue()
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error continuing rebase: %s", err)
					m.MessageType = "error"
					return m, nil
				}
				conflicts, _ := git.GetConflictFiles()
				if len(m.LogLines) > 0 {
					m.LogLines = conflicts
					m.Message = "Conflicts still present. Please resolve all conflicts."
					m.MessageType = "warning"
					return m, nil
				}
				m.Message = "✓ Rebase completed successfully"
				m.MessageType = "success"
				m.CurrentView = model.FileView
				m.DialogTarget = ""
				m.LogLines = nil
				m.RefreshData()
				return m, nil
			}
		case "X", "x":
			if len(m.LogLines) > 0 {
				err := git.RebaseAbort()
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error aborting rebase: %s", err)
					m.MessageType = "error"
					return m, nil
				}
				m.Message = "Rebase aborted."
				m.MessageType = "info"
				m.CurrentView = model.FileView
				m.DialogTarget = ""
				m.LogLines = nil
				m.RefreshData()
				return m, nil
			}
		case "esc":
			m.CurrentView = model.FileView
			m.DialogTarget = ""
			m.Message = ""
			m.Cursor = 0
			m.LogLines = nil
			return m, nil
		}
	}
	return m, nil
}
