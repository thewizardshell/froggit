package styles

import "github.com/charmbracelet/lipgloss"

// Paleta base
const (
	Green     = "76"  // Verde menta
	DarkGreen = "34"  // Verde oscuro
	Gray      = "240" // Gris
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
)
