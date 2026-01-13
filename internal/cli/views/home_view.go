package views

import (
	"github.com/oshanavishkapiries/playbuddy/internal/cli/ui"
	"github.com/charmbracelet/lipgloss"
	"fmt"
)


func ShowHomeView() {
    
	renderedTitle := ui.TitleStyle.Render(ui.AsciiArt)

	description := ui.DescStyle.Render("Welcome to PlayBuddy - Your Ultimate Torrent CLI Companion")
	
	menu := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(ui.Green).Render("(Available Commands):"),
		"• search   - Search for torrents online",
		"• download - Start downloading from a magnet link",
		"• status   - Check current download progress",
		"• exit     - Close the application",
	)

	content := ui.BorderStyle.Render(menu)

	
	fmt.Println(renderedTitle)
	fmt.Println("  ", description) 
	fmt.Println("\n", content)

}