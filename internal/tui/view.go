package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var s strings.Builder

	asciiTitle := `
                   ___           ___           ___           ___                                
     ___          /  /\         /  /\         /  /\         /  /\           ___         ___     
    /  /\        /  /::\       /  /::\       /  /::\       /  /::\         /__/\       /__/\    
   /  /::\      /  /:/\:\     /  /:/\:\     /  /:/\:\     /  /:/\:\        \__\:\      \  \:\   
  /  /:/\:\    /  /::\ \:\   /  /:/  \:\   /  /:/  \:\   /  /:/  \:\       /  /::\      \__\:\  
 /  /::\ \:\  /__/:/\:\_\:\ /__/:/ \__\:\ /__/:/_\_ \:\ /__/:/_\_ \:\   __/  /:/\/      /  /::\ 
/__/:/\:\ \:\ \__\/~|::\/:/ \  \:\ /  /:/ \  \:\__/\_\/ \  \:\__/\_\/  /__/\/:/~~      /  /:/\:\
\__\/  \:\_\/    |  |:|::/   \  \:\  /:/   \  \:\ \:\    \  \:\ \:\    \  \::/        /  /:/__\/
     \  \:\      |  |:|\/     \  \:\/:/     \  \:\/:/     \  \:\/:/     \  \:\       /__/:/     
      \__\/      |__|:|~       \  \::/       \  \::/       \  \::/       \__\/       \__\/      
                  \__\|         \__\/         \__\/         \__\/                               
			    @thewizardshell - jun 2025 - version beta 0.0.1 - in development
`
	// Título principal con ASCII art usando tu TitleStyle
	s.WriteString(TitleStyle.Render(asciiTitle) + "\n\n")

	// Información de la rama actual
	s.WriteString(fmt.Sprintf("current branch: %s\n\n", HeaderStyle.Render(m.CurrentBranch)))

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
	case NewBranchView:
		s.WriteString(m.renderNewBranchView())
	case ConfirmDialog:
		s.WriteString(m.renderConfirmDialog())
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

	// Agregar indicadores de estado para fetch y pull
	if m.IsFetching {
		s.WriteString("\n" + SpinnerStyle.Render(fmt.Sprintf("%s Fetching...", m.SpinnerFrames[m.SpinnerIndex])))
	}
	if m.IsPulling {
		s.WriteString("\n" + SpinnerStyle.Render(fmt.Sprintf("%s Pulling...", m.SpinnerFrames[m.SpinnerIndex])))
	}

	return s.String()
}

func (m Model) renderFileView() string {
	var s strings.Builder

	// Añadir resumen del estado
	stagedCount := 0
	unstagedCount := 0
	for _, file := range m.Files {
		if file.Staged {
			stagedCount++
		} else {
			unstagedCount++
		}
	}

	s.WriteString(HeaderStyle.Render("Git Status:") + "\n")
	s.WriteString(fmt.Sprintf("📦 Stage: %d archivos\n", stagedCount))
	s.WriteString(fmt.Sprintf("📝 Sin staging: %d archivos\n", unstagedCount))
	s.WriteString("\n")

	s.WriteString(HeaderStyle.Render("Modified files:") + "\n\n")

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

	// Modificar la sección de controles
	s.WriteString("\n" + BorderStyle.Render(
		HelpStyle.Render("Controles:\n")+
			HelpStyle.Render("  [↑/↓] navegar  [espacio] stage/unstage  [a] stage todos  [x] descartar cambios")+
			HelpStyle.Render("  [c] commit  [b] ramas  [m] remotes  [p] push  [f] fetch  [l] pull  [r] refresh  [q] salir"),
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

	s.WriteString(HeaderStyle.Render("🌿 Ramas:") + "\n\n")

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
			HelpStyle.Render("  [↑/↓] navegar  [Enter] cambiar rama  [n] nueva rama  [d] eliminar rama  [Esc] volver"),
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
				cursor = "▶ "
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

func (m Model) renderNewBranchView() string {
	var s strings.Builder

	s.WriteString(HeaderStyle.Render("🌿 Nueva Rama:") + "\n\n")
	s.WriteString(InputStyle.Render(m.NewBranchName+"_") + "\n\n")

	s.WriteString(BorderStyle.Render(
		HelpStyle.Render("Escribe el nombre de la rama y presiona [Enter] para crear\n") +
			HelpStyle.Render("[Esc] para cancelar"),
	))

	return s.String()
}

func (m Model) renderConfirmDialog() string {
	var s strings.Builder
	var message string

	switch m.DialogType {
	case "delete_branch":
		message = fmt.Sprintf("¿Estás seguro de que deseas eliminar la rama '%s'?", m.DialogTarget)
	case "discard_changes":
		message = fmt.Sprintf("¿Estás seguro de que deseas descartar los cambios en '%s'?", m.DialogTarget)
	}

	// Crear un "modal" con bordes
	s.WriteString("\n\n")
	s.WriteString(BorderStyle.Render(
		HeaderStyle.Render("⚠ Confirmar acción") + "\n\n" +
			NormalStyle.Render(message) + "\n\n" +
			HelpStyle.Render("[y] Sí  [n] No"),
	))

	return s.String()
}
