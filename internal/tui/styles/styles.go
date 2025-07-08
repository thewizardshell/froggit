package styles

import "github.com/charmbracelet/lipgloss"

const (
	Green        = "76"
	DarkGreen    = "34"
	GrayDark     = "240"
	GrayMid      = "248"
	GrayLight    = "250"
	Red          = "196"
	GrayLightest = "252"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Bold(true).
			PaddingTop(1).
			PaddingBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DarkGreen)).
			Background(lipgloss.Color(GrayDark)).
			Bold(true).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(GrayLight))

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Background(lipgloss.Color(GrayDark)).
			Bold(true).
			Padding(0, 1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Red)).
			Background(lipgloss.Color(GrayDark)).
			Bold(true).
			Padding(0, 1)

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true).
			Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Background(lipgloss.Color(GrayDark)).
			Bold(true).
			Padding(0, 1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(GrayLightest)).
			Italic(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(GrayMid)).
			Padding(1, 2).
			Margin(1, 0)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DarkGreen)).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(GrayMid)).
			Padding(0, 1)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Green)).
			Bold(true)

	GraphSymbolStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

	CommitHashStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("81")).
			Bold(true)

	ControlTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(Green)).
				Bold(true).
				MarginBottom(1)

	ControlSectionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("254")).
				Bold(true).
				MarginTop(1)

	ControlKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DarkGreen)).
			Bold(true)

	ControlDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(GrayLightest))
	ControlsBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(GrayMid)).
				Background(lipgloss.Color("239")).
				Padding(1, 2).
				Margin(1, 0).
				Width(60)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(GrayMid)).
			Background(lipgloss.Color("236")).
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
			Background(lipgloss.Color(GrayDark))

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(GrayLight)).
			Background(lipgloss.Color(GrayDark)).
			Padding(0, 1).
			Width(100)

	MainContentStyle = lipgloss.NewStyle().
				Width(80).
				Height(25)
)
