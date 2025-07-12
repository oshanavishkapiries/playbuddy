package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
