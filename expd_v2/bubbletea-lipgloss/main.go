package main 

import (
	"fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    PaddingTop(2).
    PaddingLeft(4).
    Width(22)

type model struct {
    choices  []string           
    cursor   int                
    selected map[int]struct{} 
}

func (m model) Init() tea.Cmd {
	fmt.Println("Init")
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//fmt.Println("💢 Update" , msg)
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }
    return m, nil
}

func (m model) View() string {
    s := "What should we buy at the market?\n\n"
    for i, choice := range m.choices {
        cursor := " "

        if m.cursor == i {
            cursor = ">"
        }

        checked := " "

		_, ok := m.selected[i]

		if ok {
            checked = "x"
        }

        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }
    s += "\nPress q to quit.\n"
    return style.Render(s)
}


func initialModel() model {
	return model{
		choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		selected: make(map[int]struct{}),
	}
}


func main() {
    p := tea.NewProgram(initialModel())
	_, err := p.Run()
	if err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}