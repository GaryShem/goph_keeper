package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"goph_keeper/client/internal/views"
)

func main() {
	p := tea.NewProgram(views.NewRootModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
