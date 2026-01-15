package mainview

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    textInput textinput.Model
    width     int
    height    int
}

func InitialModel() Model {
    ti := textinput.New()
    ti.Placeholder = "Ask a question or type a command..."
    ti.Focus()
    ti.CharLimit = 156
    ti.Width = 60
    ti.Prompt = " > "

    return Model{
        textInput: ti,
    }
}

func (m Model) Init() tea.Cmd {
    return textinput.Blink
}