package mainview

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/oshanavishkapiries/playbuddy/internal/cli/ui"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing PlayBuddy..."
	}

	renderedTitle := ui.TitleStyle.Render(ui.AsciiArt())

	menu := lipgloss.JoinVertical(lipgloss.Left,
		"/search   - Search for torrents online",
		"/download - Start downloading from a magnet link",
		"/status   - Check current download progress",
		"/settings - Change settings",
		"/help     - Show this help message",
		"/exit     - Close the application",
	)

	// Input
	inputBox := ui.InputBoxStyle.Render(m.textInput.View())

	mainLayout := lipgloss.JoinVertical(lipgloss.Center,
		renderedTitle, "\n", menu, "\n\n", inputBox,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, mainLayout)
}
