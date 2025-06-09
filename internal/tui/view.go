package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var s strings.Builder

	// Título principal
	s.WriteString(TitleStyle.Render("🌿 GitEasy") + "\n\n")

	// Información de la rama actual
	s.WriteString(fmt.Sprintf("🌱 Rama actual: %s\n\n", HeaderStyle.Render(m.CurrentBranch)))

	// Renderizar vista actual
	switch m.CurrentView {
	case FileView:
		s.WriteString(m.renderFileView())
	case CommitView:
		s.WriteString(m.renderCommitView())
	case BranchView:
		s.WriteString(m.renderBranchView())
	case RemoteView:
		s.WriteString(m.renderRemoteView())
	case AddRemoteView:
		s.WriteString(m.renderAddRemoteView())
	}

	// Mensaje de estado
	if m.Message != "" {
		s.WriteString("\n")
		switch m.MessageType {
		case "error":
			s.WriteString(ErrorStyle.Render(m.Message))
		case "success":
			s.WriteString(SuccessStyle.Render(m.Message))
		default:
			s.WriteString(NormalStyle.Render(m.Message))
		}
	}

	return s.String()
}

func (m Model) renderFileView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("📁 Archivos modificados:") + "\n\n")

	if len(m.Files) == 0 {
		s.WriteString(HelpStyle.Render("No hay archivos modificados\n"))
	} else {
		for i, file := range m.Files {
			cursor := " "
			if m.Cursor == i {
				cursor = "▶"
			}

			staged := " "
			if file.Staged {
				staged = "✓"
			}

			style := NormalStyle
			if m.Cursor == i {
				style = SelectedStyle
			}

			line := fmt.Sprintf("%s [%s] %s %s", cursor, staged, file.Status, file.Name)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	s.WriteString("\n" + BorderStyle.Render(
		HelpStyle.Render("Controles:\n")+
			HelpStyle.Render("  [↑/↓] navegar  [espacio] stage/unstage  [a] stage todos\n")+
			HelpStyle.Render("  [c] commit  [b] ramas  [m] remotes  [p] push  [r] refresh  [q] salir"),
	))

	return s.String()
}

func (m Model) renderCommitView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("💬 Mensaje de commit:") + "\n\n")
	s.WriteString(InputStyle.Render(m.CommitMsg+"_") + "\n\n")

	s.WriteString(BorderStyle.Render(
		HelpStyle.Render("Escribe tu mensaje y presiona [Enter] para confirmar\n") +
			HelpStyle.Render("[Esc] para cancelar"),
	))

	return s.String()
}

func (m Model) renderBranchView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("🌿 Seleccionar rama:") + "\n\n")

	for i, branch := range m.Branches {
		cursor := " "
		if m.Cursor == i {
			cursor = "▶"
		}

		current := " "
		if branch == m.CurrentBranch {
			current = "●"
		}

		style := NormalStyle
		if m.Cursor == i {
			style = SelectedStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, current, branch)
		s.WriteString(style.Render(line) + "\n")
	}

	s.WriteString("\n" + BorderStyle.Render(
		HelpStyle.Render("Controles:\n")+
			HelpStyle.Render("  [↑/↓] navegar  [Enter] cambiar rama  [Esc] volver"),
	))

	return s.String()
}

func (m Model) renderRemoteView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("🔗 Repositorios remotos:") + "\n\n")

	if len(m.Remotes) == 0 {
		s.WriteString(HelpStyle.Render("No hay repositorios remotos configurados\n"))
	} else {
		for i, remote := range m.Remotes {
			cursor := " "
			if m.Cursor == i {
				cursor = "▶"
			}

			style := NormalStyle
			if m.Cursor == i {
				style = SelectedStyle
			}

			line := fmt.Sprintf("%s %s", cursor, remote)
			s.WriteString(style.Render(line) + "\n")
		}
	}

	s.WriteString("\n" + BorderStyle.Render(
		HelpStyle.Render("Controles:\n")+
			HelpStyle.Render("  [↑/↓] navegar  [n] nuevo remote  [d] eliminar  [Esc] volver"),
	))

	return s.String()
}

func (m Model) renderAddRemoteView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("➕ Añadir nuevo remote:") + "\n\n")

	// Campo nombre
	nameLabel := "Nombre:"
	if m.InputField == "name" {
		nameLabel = "▶ " + nameLabel
	} else {
		nameLabel = "  " + nameLabel
	}

	nameStyle := NormalStyle
	if m.InputField == "name" {
		nameStyle = InputStyle
	}

	s.WriteString(HelpStyle.Render(nameLabel) + "\n")
	s.WriteString(nameStyle.Render(m.RemoteName+"_") + "\n\n")

	// Campo URL
	urlLabel := "URL:"
	if m.InputField == "url" {
		urlLabel = "▶ " + urlLabel
	} else {
		urlLabel = "  " + urlLabel
	}

	urlStyle := NormalStyle
	if m.InputField == "url" {
		urlStyle = InputStyle
	}

	s.WriteString(HelpStyle.Render(urlLabel) + "\n")
	s.WriteString(urlStyle.Render(m.RemoteURL+"_") + "\n\n")

	s.WriteString(BorderStyle.Render(
		HelpStyle.Render("Controles:\n") +
			HelpStyle.Render("  [Tab] cambiar campo  [Enter] confirmar/siguiente  [Esc] cancelar"),
	))

	return s.String()
}
