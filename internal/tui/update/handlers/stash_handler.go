package handlers

import (
	"fmt"
	"froggit/internal/git"
	"froggit/internal/tui/model"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleStashView(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.Cursor > 0 {
			m.Cursor--
		}
		return m, nil

	case "down":
		if m.Cursor < len(m.Stashes)-1 {
			m.Cursor++
		}
		return m, nil

	case "s", "S":
		// Save stash - only if there are changes
		if len(m.Files) > 0 {
			m.CurrentView = model.StashMessageView
			m.StashMessage = ""
			m.Message = ""
			m.MessageType = ""
			return m, nil
		} else {
			m.Message = "⚠ No changes to stash"
			m.MessageType = "warning"
			return m, nil
		}

	case "enter":
		// Apply stash
		if len(m.Stashes) > 0 && m.Cursor < len(m.Stashes) {
			selectedStash := m.Stashes[m.Cursor]
			stashRef := git.GetStashRef(selectedStash)

			if err := git.StashApply(stashRef); err != nil {
				m.Message = fmt.Sprintf("✗ Error applying stash: %s", err)
				m.MessageType = "error"
			} else {
				m.Message = fmt.Sprintf("✓ Stash %s applied successfully", stashRef)
				m.MessageType = "success"
				m.RefreshData()
			}
		}
		return m, nil

	case "p", "P":
		// Pop stash (apply and remove)
		if len(m.Stashes) > 0 && m.Cursor < len(m.Stashes) {
			selectedStash := m.Stashes[m.Cursor]
			stashRef := git.GetStashRef(selectedStash)

			if err := git.StashPop(); err != nil {
				m.Message = fmt.Sprintf("✗ Error popping stash: %s", err)
				m.MessageType = "error"
			} else {
				m.Message = fmt.Sprintf("✓ Stash %s popped successfully", stashRef)
				m.MessageType = "success"
				m.RefreshData()
			}
		}
		return m, nil

	case "d", "D":
		if len(m.Stashes) > 0 && m.Cursor < len(m.Stashes) {
			selectedStash := m.Stashes[m.Cursor]
			stashRef := git.GetStashRef(selectedStash)

			m.DialogType = "drop_stash"
			m.DialogTarget = stashRef
			m.CurrentView = model.ConfirmDialog
		}
		return m, nil

	case "v", "V":
		if len(m.Stashes) > 0 && m.Cursor < len(m.Stashes) {
			selectedStash := m.Stashes[m.Cursor]
			stashRef := git.GetStashRef(selectedStash)

			diff, err := git.StashShow(stashRef)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error viewing stash: %s", err)
				m.MessageType = "error"
			} else {
				// Store diff in LogLines for display (split by lines)
				diffLines := strings.Split(diff, "\n")
				m.LogLines = append([]string{fmt.Sprintf("Stash %s:", stashRef), ""}, diffLines...)
				m.CurrentView = model.LogGraphView
				m.Message = fmt.Sprintf("Viewing stash %s (Press Esc to return)", stashRef)
				m.MessageType = "info"
			}
		}
		return m, nil

	case "esc":
		m.CurrentView = model.FileView
		m.Message = ""
		m.MessageType = ""
		return m, nil

	case "?":
		m.Message = "Stash Help: [S]ave [Enter]Apply [P]op [D]rop [V]iew [↑/↓]Navigate [Esc]Back"
		m.MessageType = "info"
		return m, nil
	}

	return m, nil
}

func HandleStashMessageView(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Save stash with message
		message := m.StashMessage
		if message == "" {
			message = "Work in progress"
		}

		if err := git.SaveStash(message); err != nil {
			m.Message = fmt.Sprintf("✗ Error saving stash: %s", err)
			m.MessageType = "error"
		} else {
			m.Message = fmt.Sprintf("✓ Stash saved: %s", message)
			m.MessageType = "success"
			m.CurrentView = model.StashView
			m.StashMessage = ""
			m.RefreshData()
		}
		return m, nil

	case "esc":
		m.CurrentView = model.StashView
		m.StashMessage = ""
		m.Message = ""
		m.MessageType = ""
		return m, nil

	case "backspace":
		if len(m.StashMessage) > 0 {
			m.StashMessage = m.StashMessage[:len(m.StashMessage)-1]
		}
		return m, nil

	default:
		if len(msg.Runes) == 1 && msg.Runes[0] >= 32 && msg.Runes[0] <= 126 {
			m.StashMessage += string(msg.Runes)
		}
		return m, nil
	}
}
