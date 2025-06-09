package tui

import (
	"fmt"

	"giteasy/internal/git"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.CurrentView {
		case FileView:
			return m.updateFileView(msg)
		case CommitView:
			return m.updateCommitView(msg)
		case BranchView:
			return m.updateBranchView(msg)
		case RemoteView:
			return m.updateRemoteView(msg)
		case AddRemoteView:
			return m.updateAddRemoteView(msg)
		}
	}
	return m, nil
}

func (m Model) updateFileView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Files)-1 {
			m.Cursor++
		}

	case " ": // Spacebar para stagear/unstagear
		if len(m.Files) > 0 {
			m.Files[m.Cursor].Staged = !m.Files[m.Cursor].Staged
			if m.Files[m.Cursor].Staged {
				git.Add(m.Files[m.Cursor].Name)
				m.Message = fmt.Sprintf("✓ Archivo %s añadido al stage", m.Files[m.Cursor].Name)
				m.MessageType = "success"
			} else {
				git.Reset(m.Files[m.Cursor].Name)
				m.Message = fmt.Sprintf("✓ Archivo %s removido del stage", m.Files[m.Cursor].Name)
				m.MessageType = "success"
			}
		}

	case "a": // Añadir todos los archivos
		for i := range m.Files {
			if !m.Files[i].Staged {
				m.Files[i].Staged = true
				git.Add(m.Files[i].Name)
			}
		}
		m.Message = "✓ Todos los archivos añadidos al stage"
		m.MessageType = "success"

	case "c": // Ir a vista de commit
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
			m.Message = "⚠ No hay archivos en el stage para commitear"
			m.MessageType = "error"
		}

	case "b": // Ir a vista de ramas
		m.CurrentView = BranchView
		m.Cursor = 0
		m.Message = ""

	case "m": // Ir a vista de remotes
		m.CurrentView = RemoteView
		m.Cursor = 0
		m.Message = ""

	case "r": // Refresh
		m.RefreshData()
		m.Message = "✓ Estado actualizado"
		m.MessageType = "success"

	case "p": // Push
		err := git.Push()
		if err != nil {
			m.Message = fmt.Sprintf("✗ Error al hacer push: %s", err.Error())
			m.MessageType = "error"
		} else {
			m.Message = "✓ Push realizado exitosamente"
			m.MessageType = "success"
		}
	}

	return m, nil
}

func (m Model) updateCommitView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.CurrentView = FileView
		m.CommitMsg = ""
		m.Message = ""

	case "enter":
		if m.CommitMsg != "" {
			err := git.Commit(m.CommitMsg)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error al hacer commit: %s", err.Error())
				m.MessageType = "error"
			} else {
				m.Message = "✓ Commit realizado exitosamente"
				m.MessageType = "success"
				m.CommitMsg = ""
				m.RefreshData()
				m.CurrentView = FileView
			}
		}

	case "backspace":
		if len(m.CommitMsg) > 0 {
			m.CommitMsg = m.CommitMsg[:len(m.CommitMsg)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.CommitMsg += msg.String()
		}
	}

	return m, nil
}

func (m Model) updateBranchView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.CurrentView = FileView
		m.Message = ""

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Branches)-1 {
			m.Cursor++
		}

	case "enter":
		if len(m.Branches) > 0 {
			selectedBranch := m.Branches[m.Cursor]
			if selectedBranch != m.CurrentBranch {
				err := git.Checkout(selectedBranch)
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error al cambiar rama: %s", err.Error())
					m.MessageType = "error"
				} else {
					m.CurrentBranch = selectedBranch
					m.Message = fmt.Sprintf("✓ Cambiado a rama '%s'", selectedBranch)
					m.MessageType = "success"
					m.RefreshData()
				}
			}
			m.CurrentView = FileView
		}
	}

	return m, nil
}

func (m Model) updateRemoteView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.CurrentView = FileView
		m.Message = ""

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(m.Remotes)-1 {
			m.Cursor++
		}

	case "n": // Nuevo remote
		m.CurrentView = AddRemoteView
		m.InputField = "name"
		m.RemoteName = ""
		m.RemoteURL = ""
		m.Message = ""

	case "d": // Eliminar remote
		if len(m.Remotes) > 0 && m.Cursor < len(m.Remotes) {
			remoteLine := m.Remotes[m.Cursor]
			remoteName := remoteLine[:len(remoteLine)-len(" -> "+remoteLine[len(remoteLine):])]
			err := git.RemoveRemote(remoteName)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error al eliminar remote: %s", err.Error())
				m.MessageType = "error"
			} else {
				m.Message = fmt.Sprintf("✓ Remote '%s' eliminado", remoteName)
				m.MessageType = "success"
				m.RefreshData()
			}
		}
	}

	return m, nil
}

func (m Model) updateAddRemoteView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.CurrentView = RemoteView
		m.InputField = ""
		m.RemoteName = ""
		m.RemoteURL = ""
		m.Message = ""

	case "tab":
		if m.InputField == "name" {
			m.InputField = "url"
		} else {
			m.InputField = "name"
		}

	case "enter":
		if m.InputField == "name" && m.RemoteName != "" {
			m.InputField = "url"
		} else if m.InputField == "url" && m.RemoteURL != "" && m.RemoteName != "" {
			err := git.AddRemote(m.RemoteName, m.RemoteURL)
			if err != nil {
				m.Message = fmt.Sprintf("✗ Error al añadir remote: %s", err.Error())
				m.MessageType = "error"
			} else {
				m.Message = fmt.Sprintf("✓ Remote '%s' añadido exitosamente", m.RemoteName)
				m.MessageType = "success"
				m.RefreshData()
				m.CurrentView = RemoteView
				m.InputField = ""
				m.RemoteName = ""
				m.RemoteURL = ""
			}
		}

	case "backspace":
		if m.InputField == "name" && len(m.RemoteName) > 0 {
			m.RemoteName = m.RemoteName[:len(m.RemoteName)-1]
		} else if m.InputField == "url" && len(m.RemoteURL) > 0 {
			m.RemoteURL = m.RemoteURL[:len(m.RemoteURL)-1]
		}

	default:
		if len(msg.String()) == 1 {
			if m.InputField == "name" {
				m.RemoteName += msg.String()
			} else if m.InputField == "url" {
				m.RemoteURL += msg.String()
			}
		}
	}

	return m, nil
}
