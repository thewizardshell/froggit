package main

import (
	"fmt"
	"log"
	"os"

	"giteasy/internal/git"
	"giteasy/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Verificar si estamos en un repositorio Git o si necesitamos inicializar
	if !git.IsGitRepository() {
		fmt.Println("🔧 No estás en un repositorio Git")
		fmt.Print("¿Deseas inicializar un repositorio Git aquí? (y/n): ")

		var response string
		fmt.Scanln(&response)

		if response == "y" || response == "Y" || response == "yes" {
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

	// Inicializar la aplicación TUI
	model := tui.InitialModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
