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
	VERSION            = "beta-0.1.1"
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

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := update.Update(a.m, msg)
	a.m = newModel
	return a, cmd
}

func (a App) View() string {
	return tui.Render(a.m)
}

func main() {
	versionFlag := flag.Bool("version", false, "Print version information")
	helpFlag := flag.Bool("help", false, "Print help information")
	commandsFlag := flag.Bool("commands", false, "List supported Git commands")
	keyboardFlag := flag.Bool("keys", false, "List keyboard shortcuts")

	flag.Parse()

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
		fmt.Println("🔧 You are not in a Git repository.")
		fmt.Print("Do you want to initialize a Git repository here? (y/n): ")

		var resp string
		fmt.Scanln(&resp)
		if resp == "y" || resp == "Y" || resp == "yes" {
			if err := git.InitRepository(); err != nil {
				fmt.Printf("❌ Failed to initialize Git repository: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Git repository successfully initialized.")
		} else {
			fmt.Println("👋 Exiting...")
			os.Exit(0)
		}
	}

	app := App{m: model.InitialModel()}

	p := tea.NewProgram(app, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
