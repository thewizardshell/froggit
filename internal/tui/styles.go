package tui

import "github.com/charmbracelet/lipgloss"

// Estilos con tema verde
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("76"))

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("34")).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("76")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("76")).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("34"))

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("76")).
			Bold(true)
)
