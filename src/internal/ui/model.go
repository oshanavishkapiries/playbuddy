package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
	"github.com/oshanavishkapiries/playbuddy/src/internal/services"
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/views"
)

// Navigation history entry
type HistoryEntry struct {
	State AppState
	Data  interface{} // Can store state-specific data
}

// Model for the application
type Model struct {
	state           AppState
	input           string
	results         []models.Torrent
	cursor          int
	loading         bool
	error           string
	selectedTorrent *models.Torrent
	searchService   *services.SearchService
	lastSearch      string
	// Navigation history
	history      []HistoryEntry
	historyIndex int
	maxHistory   int
	// Loader animation
	loaderFrame int
}

// Loader frames for animation
var loaderFrames = []string{"|", "/", "-", "\\"}

// Initialize the model
func InitialModel() Model {
	return Model{
		state:         StateWelcome,
		cursor:        0,
		loading:       false,
		searchService: services.NewSearchService(),
		history:       make([]HistoryEntry, 0),
		historyIndex:  -1,
		maxHistory:    50,
		loaderFrame:   0,
	}
}

// Add to navigation history
func (m *Model) addToHistory(state AppState, data interface{}) {
	if m.historyIndex < len(m.history)-1 {
		m.history = m.history[:m.historyIndex+1]
	}
	entry := HistoryEntry{State: state, Data: data}
	m.history = append(m.history, entry)
	m.historyIndex++
	if len(m.history) > m.maxHistory {
		m.history = m.history[1:]
		m.historyIndex--
	}
}

func (m *Model) navigateBack() bool {
	if m.historyIndex > 0 {
		m.historyIndex--
		entry := m.history[m.historyIndex]
		m.state = entry.State
		if data, ok := entry.Data.(map[string]interface{}); ok {
			if input, exists := data["input"]; exists {
				m.input = input.(string)
			}
			if cursor, exists := data["cursor"]; exists {
				m.cursor = cursor.(int)
			}
			if results, exists := data["results"]; exists {
				m.results = results.([]models.Torrent)
			}
		}
		return true
	}
	return false
}

func (m *Model) saveCurrentState() {
	data := map[string]interface{}{
		"input":   m.input,
		"cursor":  m.cursor,
		"results": m.results,
	}
	m.addToHistory(m.state, data)
}

// Loader animation command
func loaderTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return loaderTickMsg{}
	})
}

type loaderTickMsg struct{}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case StateWelcome:
			if msg.String() == KeyEnter || msg.String() == " " {
				m.saveCurrentState()
				m.state = StateMainMenu
				m.cursor = 0
			} else if msg.String() == KeyQuit || msg.String() == KeyEscape || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		case StateMainMenu:
			switch msg.String() {
			case KeyUp, "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case KeyDown, "j":
				if m.cursor < 2 {
					m.cursor++
				}
			case KeyEnter:
				switch m.cursor {
				case 0:
					m.saveCurrentState()
					m.state = StateSearch
					m.cursor = 0
					m.input = ""
					m.results = nil
					m.error = ""
				case 1:
					m.error = "Download hub not implemented yet"
				case 2:
					m.error = "Settings not implemented yet"
				}
			case KeyQuit:
				return m, tea.Quit
			case KeyEscape:
				m.saveCurrentState()
				m.state = StateWelcome
				m.cursor = 0
			}
		case StateSearch:
			switch msg.String() {
			case KeyEnter:
				if strings.TrimSpace(m.input) != "" {
					m.saveCurrentState()
					m.lastSearch = strings.TrimSpace(m.input)
					m.loading = true
					m.results = nil
					return m, tea.Batch(m.performSearch(), loaderTick())
				}
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			case KeyEscape:
				m.saveCurrentState()
				m.state = StateMainMenu
				m.cursor = 0
				m.input = ""
				m.results = nil
				m.error = ""
			case "ctrl+c":
				return m, tea.Quit
			default:
				if len(msg.String()) == 1 {
					m.input += msg.String()
				}
			}
		case StateResults:
			switch msg.String() {
			case KeyUp, "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case KeyDown, "j":
				if m.cursor < len(m.results)-1 {
					m.cursor++
				}
			case KeyEnter:
				if len(m.results) > 0 && m.cursor < len(m.results) {
					m.saveCurrentState()
					m.selectedTorrent = &m.results[m.cursor]
					m.state = StateTorrentDetails
					m.cursor = 0
				}
			case KeyEscape:
				m.saveCurrentState()
				m.state = StateSearch
				m.cursor = 0 // focus input/results, do not clear
			case "ctrl+c":
				return m, tea.Quit
			}
		case StateTorrentDetails:
			switch msg.String() {
			case KeyEscape:
				m.saveCurrentState()
				m.state = StateResults
				m.selectedTorrent = nil
			case KeyQuit:
				return m, tea.Quit
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		// Global navigation keys
		switch msg.String() {
		case KeyBack:
			if m.navigateBack() {
				return m, nil
			}
		}
	case SearchResultMsg:
		m.loading = false
		if msg.Error != "" {
			m.error = msg.Error
		} else {
			m.results = msg.Torrents
			if len(m.results) > 0 {
				m.state = StateResults
				m.cursor = 0
			}
		}
	case loaderTickMsg:
		if m.loading {
			m.loaderFrame = (m.loaderFrame + 1) % len(loaderFrames)
			return m, loaderTick()
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case StateWelcome:
		return views.ViewWelcome()
	case StateMainMenu:
		return views.ViewMainMenu(m.cursor, m.error)
	case StateSearch:
		return views.ViewSearch(m.input, m.loading, m.loaderFrame, loaderFrames, m.error)
	case StateResults:
		return views.ViewResults(m.results, m.cursor, truncateString)
	case StateTorrentDetails:
		return views.ViewTorrentDetails(m.selectedTorrent)
	default:
		return "Unknown state"
	}
}

// performSearch performs the torrent search
func (m Model) performSearch() tea.Cmd {
	return func() tea.Msg {
		torrents := m.searchService.GetAllTorrents(m.lastSearch)
		return SearchResultMsg{Torrents: torrents}
	}
}

func truncateString(s string, max int) string {
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max-3]) + "..."
	}
	return s
}
