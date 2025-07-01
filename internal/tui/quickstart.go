package tui

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"froggit/internal/git"
	"froggit/internal/tui/model"
	"froggit/internal/tui/update"

	tea "github.com/charmbracelet/bubbletea"
)

// QuickStartFlow handles the initial onboarding and repo selection/creation
func QuickStartFlow() {
	hasGh := true
	if _, err := exec.LookPath("gh"); err != nil {
		hasGh = false
	}

	fmt.Println("\n\033[1;36m████ Froggit Quick Start ████\033[0m\n")
	fmt.Println("You are not in a Git repository. What would you like to do?")
	fmt.Println("  \033[1;32m1)\033[0m Initialize a new Git repository here")
	if hasGh {
		fmt.Println("  \033[1;34m2)\033[0m Clone a repository from GitHub (requires GitHub CLI)")
		fmt.Println("  \033[1;35m3)\033[0m Create a new repository on GitHub (requires GitHub CLI)")
		fmt.Print("\nEnter 1, 2 or 3: ")
	} else {
		fmt.Println("\n\033[1;31mGitHub CLI (gh) not found. Only option 1 is available.\033[0m")
		fmt.Print("\nEnter 1: ")
	}

	var resp string
	fmt.Scanln(&resp)
	if resp == "1" {
		if err := git.InitRepository(); err != nil {
			fmt.Printf("❌ Failed to initialize Git repository: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Git repository successfully initialized.")
	} else if hasGh && resp == "2" {
		m := model.InitialModel()
		m.CurrentView = model.RepositoryListView
		m = update.ShowRepositoryList(m, update.GetGhClient())
		app := App{M: m}
		p := tea.NewProgram(app, tea.WithAltScreen())
		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	} else if hasGh && resp == "3" {
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
			fmt.Printf("❌ Failed to create and clone repository: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n✅ Repository '%s' created and cloned!\n", repoName)
		if err := os.Chdir(repoName); err != nil {
			fmt.Printf("⚠️  Could not change directory to %s: %s\n", repoName, err)
		} else {
			fmt.Printf("\n📂 Changed directory to %s\n", repoName)
		}
	} else {
		fmt.Println("👋 Exiting...")
		os.Exit(0)
	}
}
