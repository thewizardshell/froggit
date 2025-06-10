package main

import (
	"fmt"
	"log"
	"os"

	"froggit/internal/git"
	tui "froggit/internal/tui"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	m model.Model
}

func (a App) Init() tea.Cmd {
	return nil
}

// Update delega en tui.Update y actualiza el estado interno.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := update.Update(a.m, msg)
	a.m = newModel
	return a, cmd
}

// View delega en tui.Render para pintar la TUI.
func (a App) View() string {
	return tui.Render(a.m)
}

func main() {
	// Verificar si estamos en un repositorio Git o inicializar
	if !git.IsGitRepository() {
		fmt.Println("🔧 No estás en un repositorio Git")
		fmt.Print("¿Deseas inicializar un repositorio Git aquí? (y/n): ")

		var resp string
		fmt.Scanln(&resp)
		if resp == "y" || resp == "Y" || resp == "yes" {
			if err := git.InitRepository(); err != nil {
				fmt.Printf("❌ Error al inicializar repositorio: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Repositorio Git inicializado exitosamente")
		} else {
			fmt.Println("👋 Saliendo...")
			os.Exit(0)
		}
	}

	// Construir el wrapper App con el modelo inicial
	app := App{m: model.InitialModel()}

	// Crear y arrancar el programa Bubble Tea en pantalla alterna
	p := tea.NewProgram(app, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
