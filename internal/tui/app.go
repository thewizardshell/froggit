package tui

import (
	"froggit/internal/config"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update"
	"froggit/internal/tui/update/async"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	M model.Model
	C config.Config
}

func (a App) Init() tea.Cmd {
	if a.M.IsFetching {
		return async.Spinner()
	}
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := update.Update(a.M, a.C, msg)
	a.M = newModel
	return a, cmd
}

func (a App) View() string {
	return Render(a.M, a.C)
}
