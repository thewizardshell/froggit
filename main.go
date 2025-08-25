package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"froggit/internal/config"
	"froggit/internal/git"
	tui "froggit/internal/tui"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	VERSION            = "0.4.1 - beta"
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
	- log: Show commit logs
	- merge: Merge branches
	- rebase: Reapply commits on top of another base tip
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
    L: Show log graph
    esc: Go back/cancel
    enter: Confirm/execute
`
)

type App struct {
	M model.Model
	C config.Config
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := update.Update(a.M, msg)
	a.M = newModel
	return a, cmd
}

func (a App) View() string {
	return tui.Render(a.M, a.C)
}

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	helpFlag := flag.Bool("help", false, "Print help information")
	commandsFlag := flag.Bool("commands", false, "List supported Git commands")
	keyboardFlag := flag.Bool("keys", false, "List keyboard shortcuts")
	flag.Parse()

	cfg, err := config.LoadConfig("froggit.yml")
	if err != nil {
		log.Printf("Warning: Could not load config file: %v\n", err)
		cfg = config.Config{
			Ui:  config.UiConfig{Branding: true, Position: "center"},
			Git: config.GitConfig{DefaultBranch: "main", AutoFetch: true},
		}
	}

	if *versionFlag {
		fmt.Printf("Version: %s\nAuthor: %s\nRepository: %s\n", VERSION, AUTHOR, REPO)
		os.Exit(0)
	}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *commandsFlag {
		fmt.Print(SUPPORTED_COMMANDS)
		os.Exit(0)
	}

	if *keyboardFlag {
		fmt.Print(KEYBOARD_SHORTCUTS)
		os.Exit(0)
	}

	if !git.IsGitRepository() {
		tui.QuickStartFlow()
	}

	app := tui.App{
		M: model.InitialModel(),
		C: cfg,
	}

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
