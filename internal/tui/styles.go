package tui

import "github.com/charmbracelet/lipgloss"

// Estilos con tema verde
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF7F")). // Verde brillante
		// Padding(1, 2).                         // Espaciado vertical y horizontal
		// Margin(1, 0, 1, 0).                         // Margen arriba y abajo
		BorderForeground(lipgloss.Color("#00A86B")) // Borde verde más oscuro

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#00A86B")).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E8F5E8"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90EE90"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF7F")).
			Bold(true)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2E8B57")).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#32CD32")).
			Bold(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#228B22")).
			Padding(1, 2)
)
