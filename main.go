package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"froggit/internal/git"
	tui "froggit/internal/tui"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	VERSION            = "beta-0.1.0"
	AUTHOR             = "Vicente Roa | Github: @thewizardshell"
	REPO               = "https://github.com/thewizardshell/froggit"
	SUPPORTED_COMMANDS = `
Supported Git Commands:
    - init: Initialize a new Git repository
    - status: Show working tree status
    - add: Add files to staging area
    - commit: Record changes to the repository
    - branch: List, create, or delete branches
    - checkout: Switch branches
    - remote: Manage remote repositories
    - push: Push changes to remote repository
    - pull: Fetch from remote repository
`
	KEYBOARD_SHORTCUTS = `
Keyboard Shortcuts:
    q, ctrl+c: Quit
    h: Show help
    j/k: Navigate down/up
    space: Select/deselect file
    a: Stage all files
    c: Commit changes
    b: Create new branch
    r: Add remote
    p: Push changes
    l: Pull changes
    esc: Go back/cancel
    enter: Confirm/execute
`
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
	//flags
	versionFlag := flag.Bool("version", false, "Print version information")
	helpFlag := flag.Bool("help", false, "Print help information")
	commandsFlag := flag.Bool("commands", false, "List supported Git commands")
	keyboardFlag := flag.Bool("keys", false, "List keyboard shortcuts")

	// Parse flags
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("Version: %s\nAuthor: %s\nRepository: %s\n", VERSION, AUTHOR, REPO)
		os.Exit(0)
	}

	// Handle help flag
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Handle commands flag
	if *commandsFlag {
		fmt.Print(SUPPORTED_COMMANDS)
		os.Exit(0)
	}

	if *keyboardFlag {
		fmt.Print(KEYBOARD_SHORTCUTS)
		os.Exit(0)
	}
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
