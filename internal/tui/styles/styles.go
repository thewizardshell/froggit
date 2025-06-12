package styles

import "github.com/charmbracelet/lipgloss"

// Paleta base
const (
	Green     = "76"
	DarkGreen = "34"
	Gray      = "240"
	LightGray = "252"
	Red       = "196"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DarkGreen)).
			Background(lipgloss.Color(Gray)).
			Bold(true).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(LightGray))

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Background(lipgloss.Color(Gray)).
			Bold(true).
			Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Red)).
			Background(lipgloss.Color("235")).
			Bold(true).
			Padding(0, 1)

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
		// Background(lipgloss.Color("235")).
		Bold(true).
		Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Background(lipgloss.Color("236")).
			Bold(true).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(Gray)).
			Padding(1, 2).
			Margin(1, 0)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DarkGreen)).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(Gray)).
			Padding(0, 1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Bold(true)

	ControlTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(Green)).
				Bold(true).
				MarginBottom(1)

	ControlSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Bold(true).
				MarginTop(1)

	ControlKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DarkGreen)).
			Bold(true)

	ControlDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("249"))

	// Estilo para el contenedor de controles
	ControlsBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(Gray)).
				Padding(1, 2).
				Margin(1, 0).
				Width(60) // Ancho fijo para mejor alineación
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(Gray)).
			Padding(0, 1).
			Margin(0, 1)

	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(Green)).
				Padding(0, 1).
				Margin(0, 1)

	PanelTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color(Gray))

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(LightGray)).
			Background(lipgloss.Color("235")).
			Padding(0, 1).
			Width(100) // Ajustar según el ancho de tu terminal

	// Para el layout horizontal
	MainContentStyle = lipgloss.NewStyle().
				Width(80).
				Height(25) // Ajustar según la altura de tu terminal
)
