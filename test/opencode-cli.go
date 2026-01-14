package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// model represents the state of our application
type model struct {
	textInput textinput.Model
	width     int // Terminal width
	height    int // Terminal height
}

// initialModel sets up the default state
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Ask a question or type a command..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 60
	ti.Prompt = " > "

	return model{
		textInput: ti,
	}
}

// Init is called when the program starts
func (m model) Init() tea.Cmd {
	return textinput.Blink // Command to make the cursor blink
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit // Exit the application
		}
	case tea.WindowSizeMsg:
		// Capture terminal resizing
		m.width = msg.Width
		m.height = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the UI into a string
func (m model) View() string {
	if m.width == 0 {
		return "Initializing PlayBuddy..."
	}

	// Defining Colors (OpenCode Style)
	green := lipgloss.Color("#00FF41")
	gray := lipgloss.Color("#767676")
	dimGray := lipgloss.Color("#353535")

	// 1. Header Section
	title := lipgloss.NewStyle().Foreground(green).Bold(true).Render("PLAYBUDDY")
	version := lipgloss.NewStyle().Foreground(gray).Render("v0.1.0")
	header := lipgloss.JoinVertical(lipgloss.Center, title, version)

	// 2. Command Menu Section
	cmdStyle := lipgloss.NewStyle().Foreground(green).Width(12).Align(lipgloss.Left)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Width(25).Align(lipgloss.Left)
	keyStyle := lipgloss.NewStyle().Foreground(gray).Width(10).Align(lipgloss.Right)

	// Build the menu rows
	menu := lipgloss.JoinVertical(lipgloss.Left,
		renderRow(cmdStyle, "/search", descStyle, "search for torrents", keyStyle, "ctrl+s"),
		renderRow(cmdStyle, "/help", descStyle, "view command list", keyStyle, "ctrl+h"),
		renderRow(cmdStyle, "/status", descStyle, "check downloads", keyStyle, "ctrl+p"),
		renderRow(cmdStyle, "/exit", descStyle, "close application", keyStyle, "ctrl+q"),
	)

	// 3. Input Section
	inputBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(dimGray).
		Padding(0, 1).
		Render(m.textInput.View())

	inputFooter := lipgloss.NewStyle().Foreground(gray).Render("enter send")

	// Assemble everything vertically
	mainLayout := lipgloss.JoinVertical(lipgloss.Center,
		header,
		"\n",
		menu,
		"\n\n",
		inputBox,
		inputFooter,
	)

	// Position everything in the center of the terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, mainLayout)
}

// renderRow is a helper to join three columns into a single line
func renderRow(s1 lipgloss.Style, t1 string, s2 lipgloss.Style, t2 string, s3 lipgloss.Style, t3 string) string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom, s1.Render(t1), s2.Render(t2), s3.Render(t3))
}

func main() {
	// tea.WithAltScreen() opens the UI in a dedicated full-screen buffer
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running PlayBuddy: %v", err)
		os.Exit(1)
	}
}