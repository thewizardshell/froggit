package handlers

import (
	"fmt"
	"froggit/internal/git"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update/async"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleMergeView(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	if m.CurrentView != model.MergeView {
		return m, nil
	}

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
				m.Message = "⚠ Cannot merge a branch into itself"
				m.MessageType = "warning"
				return m, nil
			}
			if m.DialogTarget == selected {
				m.DialogTarget = ""
				m.Message = "Selection unmarked"
				m.MessageType = "info"
			} else {
				m.DialogTarget = selected
				m.Message = fmt.Sprintf("Will merge %s into %s (current: %s)", m.CurrentBranch, selected, m.CurrentBranch)
				m.MessageType = "info"
			}
		}
		return m, nil
	case "M", "m":
		if m.DialogTarget != "" {
			if m.DialogTarget == m.CurrentBranch {
				m.Message = "⚠ Cannot merge a branch into itself"
				m.MessageType = "warning"
				return m, nil
			}

			m.Message = fmt.Sprintf("Switching to %s and merging %s into it...", m.DialogTarget, m.CurrentBranch)
			m.MessageType = "info"

			return m, async.PerformSwitchAndMerge(m.DialogTarget, m.CurrentBranch)
		} else {
			m.Message = "Select a target branch first by pressing space"
			m.MessageType = "warning"
			return m, nil
		}
	case "P", "p":
		if m.AwaitingPush {
			m.Message = "Pushing..."
			m.MessageType = "info"
			m.AwaitingPush = false
			return m, tea.Batch(async.PerformPush(), async.Spinner())
		}
		if len(m.LogLines) > 0 {
			err := git.MergeContinue()
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error continuing merge: %s", err)
				m.MessageType = "error"
				return m, nil
			}
			conflicts, _ := git.GetConflictFiles()
			if len(conflicts) > 0 {
				m.LogLines = conflicts
				m.Message = "Conflicts still present. Please resolve all conflicts."
				m.MessageType = "warning"
				return m, nil
			}
			m.Message = "✓ Merge completed successfully. Press [P] to push to remote."
			m.MessageType = "success"
			m.AwaitingPush = true
			return m, nil
		}
	case "X", "x":
		if len(m.LogLines) > 0 {
			err := git.MergeAbort()
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error aborting merge: %s", err)
				m.MessageType = "error"
				return m, nil
			}
			m.Message = "Merge aborted."
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
	return m, nil
}
