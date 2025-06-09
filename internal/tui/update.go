// tui/update.go
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.CurrentView != FileView {
				m.CurrentView = FileView
				m.CommitMsg = ""
				m.RemoteName = ""
				m.RemoteURL = ""
				m.Message = ""
				return m, nil
			}
		case "enter":
			if m.CurrentView == CommitView && m.CommitMsg != "" {
				err := git.Commit(m.CommitMsg)
				if err != nil {
					m.Message = fmt.Sprintf("✗ Error al hacer commit: %s", err.Error())
					m.MessageType = "error"
				} else {
					m.Message = "✓ Commit realizado exitosamente"
					m.MessageType = "success"
					m.CurrentView = FileView
					m.CommitMsg = ""
					m.RefreshData()
				}
				return m, nil
			}
		case "backspace":
			if m.CurrentView == CommitView && len(m.CommitMsg) > 0 {
				m.CommitMsg = m.CommitMsg[:len(m.CommitMsg)-1]
				return m, nil
			}
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Files)-1 {
				m.Cursor++
			}
		default:
			// Si estamos en CommitView y la tecla es un caracter imprimible
			if m.CurrentView == CommitView {
				if len(msg.Runes) == 1 && isPrintableChar(msg.Runes[0]) {
					m.CommitMsg += string(msg.Runes)
					return m, nil
				}
			} else if m.CurrentView == FileView {
				// Comandos de una letra solo disponibles en FileView
				switch msg.String() {
				case " ":
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
				case "a":
					for i := range m.Files {
						if !m.Files[i].Staged {
							m.Files[i].Staged = true
							git.Add(m.Files[i].Name)
						}
					}
					m.Message = "✓ Todos los archivos añadidos al stage"
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
						m.Message = "⚠ No hay archivos en el stage para commitear"
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
					m.Message = "✓ Estado actualizado"
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
			m.Message = fmt.Sprintf("✗ Error al hacer push: %s", msg.err.Error())
			m.MessageType = "error"
		} else {
			m.Message = "✓ Push realizado exitosamente"
			m.MessageType = "success"
			m.RefreshData()
		}
	case fetchMsg:
		m.IsFetching = false
		if msg.err != nil {
			m.Message = "✗ Error al hacer fetch: " + msg.err.Error()
			m.MessageType = "error"
		} else {
			m.Message = "✓ Fetch realizado exitosamente"
			m.MessageType = "success"
			m.RefreshData()
		}
	case pullMsg:
		m.IsPulling = false
		if msg.err != nil {
			m.Message = "✗ Error al hacer pull: " + msg.err.Error()
			m.MessageType = "error"
		} else {
			m.Message = "✓ Pull realizado exitosamente"
			m.MessageType = "success"
			m.RefreshData()
		}
	}
	return m, nil
}

// isPrintableChar verifica si un caracter es imprimible
func isPrintableChar(r rune) bool {
	return r >= 32 && r <= 126 || r >= 128 && r <= 255
}
