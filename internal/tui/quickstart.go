package tui

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"froggit/internal/config"
	"froggit/internal/git"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update"

	tea "github.com/charmbracelet/bubbletea"
)

func QuickStartFlow() {
	hasGh := true
	if _, err := exec.LookPath("gh"); err != nil {
		hasGh = false
	}

	cfg, err := config.LoadConfig("froggit.yml")
	if err != nil {
		cfg = config.Config{
			Ui:  config.UiConfig{Branding: true, Position: "center"},
			Git: config.GitConfig{DefaultBranch: "main", AutoFetch: true},
		}
	}

	m := model.Model{
		CurrentView:   model.QuickStartView,
		Cursor:        0,
		HasGitHubCLI:  hasGh,
		SpinnerFrames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}

	app := App{M: m, C: cfg}
	p := tea.NewProgram(app, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	if finalApp, ok := finalModel.(App); ok {
		HandleQuickStartAction(finalApp.M.Cursor, hasGh)
	}
}

func HandleQuickStartAction(option int, hasGh bool) {
	switch option {
	case 0:
		if err := git.InitRepository(); err != nil {
			fmt.Printf("✗ Failed to initialize Git repository: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Git repository successfully initialized.")

	case 1:
		if !hasGh {
			fmt.Println("✗ GitHub CLI not available")
			os.Exit(1)
		}

		cfg, _ := config.LoadConfig("froggit.yml")
		if cfg.Ui.Branding == false {
			cfg = config.Config{
				Ui:  config.UiConfig{Branding: true, Position: "center"},
				Git: config.GitConfig{DefaultBranch: "main", AutoFetch: true},
			}
		}

		m := model.InitialModel()
		m.CurrentView = model.RepositoryListView
		m = update.ShowRepositoryList(m, update.GetGhClient())
		app := App{M: m, C: cfg}
		p := tea.NewProgram(app, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

	case 2:
		if !hasGh {
			fmt.Println("✗ GitHub CLI not available")
			os.Exit(1)
		}

		fmt.Print("Repository name: ")
		var repoName string
		fmt.Scanln(&repoName)
		fmt.Print("Make it private? (y/n): ")
		var priv string
		fmt.Scanln(&priv)
		visibility := "public"
		if priv == "y" || priv == "Y" {
			visibility = "private"
		}
		fmt.Println("\nCreating repository on GitHub...")
		cmd := exec.Command("gh", "repo", "create", repoName, "--"+visibility, "--clone")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("✗ Failed to create and clone repository: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n✓ Repository '%s' created and cloned!\n", repoName)
		if err := os.Chdir(repoName); err != nil {
			fmt.Printf("! Could not change directory to %s: %s\n", repoName, err)
		} else {
			fmt.Printf("\n> Changed directory to %s\n", repoName)
		}
	}
}
