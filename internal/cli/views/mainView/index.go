package	mainview

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
)

func ShowMainView() {
    p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error running PlayBuddy: %v", err)
        os.Exit(1)
    }
}