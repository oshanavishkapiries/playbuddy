package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
	"github.com/oshanavishkapiries/playbuddy/src/internal/services"
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
	searchDebounce  time.Time
	searchTimer     *time.Timer

	// Navigation history
	history      []HistoryEntry
	historyIndex int
	maxHistory   int
}

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
	}
}

// Add to navigation history
func (m *Model) addToHistory(state AppState, data interface{}) {
	// Remove any future history if we're not at the end
	if m.historyIndex < len(m.history)-1 {
		m.history = m.history[:m.historyIndex+1]
	}

	entry := HistoryEntry{
		State: state,
		Data:  data,
	}

	m.history = append(m.history, entry)
	m.historyIndex++

	// Keep history size manageable
	if len(m.history) > m.maxHistory {
		m.history = m.history[1:]
		m.historyIndex--
	}
}

// Navigate back
func (m *Model) navigateBack() bool {
	if m.historyIndex > 0 {
		m.historyIndex--
		entry := m.history[m.historyIndex]
		m.state = entry.State

		// Restore state-specific data
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

// Navigate forward
func (m *Model) navigateForward() bool {
	if m.historyIndex < len(m.history)-1 {
		m.historyIndex++
		entry := m.history[m.historyIndex]
		m.state = entry.State

		// Restore state-specific data
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

// Save current state before navigation
func (m *Model) saveCurrentState() {
	data := map[string]interface{}{
		"input":   m.input,
		"cursor":  m.cursor,
		"results": m.results,
	}
	m.addToHistory(m.state, data)
}

// Init method for Bubble Tea
func (m Model) Init() tea.Cmd {
	return nil
}

// Update method for Bubble Tea
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
				case 0: // Search torrent
					m.saveCurrentState()
					m.state = StateSearch
					m.cursor = 0
					m.input = ""
				case 1: // Download hub
					// TODO: Implement download hub
					m.error = "Download hub not implemented yet"
				case 2: // Settings
					// TODO: Implement settings
					m.error = "Settings not implemented yet"
				}
			case KeyQuit:
				return m, tea.Quit
			case KeyEscape:
				m.saveCurrentState()
				m.state = StateWelcome
				m.cursor = 0
			case "ctrl+c":
				return m, tea.Quit
			}

		case StateSearch:
			switch msg.String() {
			case KeyEnter:
				if strings.TrimSpace(m.input) != "" {
					m.saveCurrentState()
					m.lastSearch = strings.TrimSpace(m.input)
					m.loading = true
					return m, m.performSearch()
				}
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
				// Trigger debounced search
				if strings.TrimSpace(m.input) != "" {
					return m, m.debouncedSearch(strings.TrimSpace(m.input))
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
					// Trigger debounced search for real-time results
					if strings.TrimSpace(m.input) != "" {
						return m, m.debouncedSearch(strings.TrimSpace(m.input))
					}
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
				m.cursor = 0
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
		case KeyBack: // Left arrow
			if m.navigateBack() {
				return m, nil
			}
		case KeyForward: // Right arrow
			if m.navigateForward() {
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

	case DebouncedSearchMsg:
		if msg.Query == strings.TrimSpace(m.input) {
			m.lastSearch = msg.Query
			m.loading = true
			return m, m.performSearch()
		}
	}

	return m, nil
}

// View method for Bubble Tea
func (m Model) View() string {
	switch m.state {
	case StateWelcome:
		return m.viewWelcome()
	case StateMainMenu:
		return m.viewMainMenu()
	case StateSearch:
		return m.viewSearch()
	case StateResults:
		return m.viewResults()
	case StateTorrentDetails:
		return m.viewTorrentDetails()
	default:
		return "Unknown state"
	}
}

func (m Model) viewWelcome() string {
	s := AsciiStyle.Render(AsciiArt) + "\n"
	s += VersionStyle.Render(Version) + "\n\n"
	s += TitleStyle.Render("Welcome to PlayBuddy Torrent Search") + "\n\n"
	s += HelpStyle.Render("Press Enter to continue or Q to quit")
	return s
}

func (m Model) viewMainMenu() string {
	s := AsciiStyle.Render(AsciiArt) + "\n"
	s += VersionStyle.Render(Version) + "\n\n"
	s += TitleStyle.Render("Main Menu") + "\n\n"

	menuItems := []string{"Search Torrent", "Download Hub", "Settings"}
	for i, item := range menuItems {
		if m.cursor == i {
			s += SelectedMenuStyle.Render("> "+item) + "\n"
		} else {
			s += MenuStyle.Render("  "+item) + "\n"
		}
	}

	if m.error != "" {
		s += "\n" + ErrorStyle.Render("Error: "+m.error) + "\n"
	}

	s += "\n" + HelpStyle.Render("Use arrow keys to navigate, Enter to select, Esc to go back")
	s += "\n" + NavigationStyle.Render("← Back  → Forward")
	return s
}

func (m Model) viewSearch() string {
	s := AsciiStyle.Render(AsciiArt) + "\n"
	s += VersionStyle.Render(Version) + "\n\n"
	s += TitleStyle.Render("Search Torrents") + "\n\n"

	if m.loading {
		s += LoadingStyle.Render("Searching... Please wait") + "\n\n"
	} else {
		s += InputStyle.Render("Search: "+m.input+"█") + "\n\n"
	}

	// Show real-time results if available
	if len(m.results) > 0 && !m.loading {
		s += fmt.Sprintf("Found %d results:\n", len(m.results))
		for i, torrent := range m.results {
			if i < 10 { // Show first 10 results
				if m.cursor == i {
					s += SelectedResultStyle.Render(fmt.Sprintf("> %s (%s)", torrent.Name, torrent.Size)) + "\n"
				} else {
					s += ResultStyle.Render(fmt.Sprintf("  %s (%s)", torrent.Name, torrent.Size)) + "\n"
				}
			}
		}
		if len(m.results) > 10 {
			s += ResultStyle.Render(fmt.Sprintf("  ... and %d more results", len(m.results)-10)) + "\n"
		}
		s += "\n"
	}

	if m.error != "" {
		s += ErrorStyle.Render("Error: "+m.error) + "\n"
	}

	s += "\n" + HelpStyle.Render("Type to search, use ↑↓ to select, Enter to view details, Esc to go back")
	s += "\n" + NavigationStyle.Render("← Back  → Forward")
	return s
}

func (m Model) viewResults() string {
	s := AsciiStyle.Render(AsciiArt) + "\n"
	s += VersionStyle.Render(Version) + "\n\n"
	s += TitleStyle.Render(fmt.Sprintf("Search Results (%d found)", len(m.results))) + "\n\n"

	if len(m.results) == 0 {
		s += ResultStyle.Render("No results found") + "\n"
	} else {
		for i, torrent := range m.results {
			if m.cursor == i {
				s += SelectedResultStyle.Render(fmt.Sprintf("> %s (%s)", torrent.Name, torrent.Size)) + "\n"
			} else {
				s += ResultStyle.Render(fmt.Sprintf("  %s (%s)", torrent.Name, torrent.Size)) + "\n"
			}
		}
	}

	s += "\n" + HelpStyle.Render("Use arrow keys to navigate, Enter to view details, Esc to go back")
	s += "\n" + NavigationStyle.Render("← Back  → Forward")
	return s
}

func (m Model) viewTorrentDetails() string {
	if m.selectedTorrent == nil {
		return "No torrent selected"
	}

	t := m.selectedTorrent
	s := AsciiStyle.Render(AsciiArt) + "\n"
	s += VersionStyle.Render(Version) + "\n\n"
	s += TitleStyle.Render("Torrent Details") + "\n\n"

	s += DetailStyle.Render(DetailLabelStyle.Render("Name: ")+DetailValueStyle.Render(t.Name)) + "\n"
	s += DetailStyle.Render(DetailLabelStyle.Render("Size: ")+DetailValueStyle.Render(t.Size)) + "\n"
	s += DetailStyle.Render(DetailLabelStyle.Render("Category: ")+DetailValueStyle.Render(t.Category)) + "\n"
	s += DetailStyle.Render(DetailLabelStyle.Render("Uploaded: ")+DetailValueStyle.Render(t.DateUploaded)) + "\n"
	s += DetailStyle.Render(DetailLabelStyle.Render("Uploader: ")+DetailValueStyle.Render(t.UploadedBy)) + "\n"
	s += DetailStyle.Render(DetailLabelStyle.Render("Seeders: ")+DetailValueStyle.Render(t.Seeders)) + "\n"
	s += DetailStyle.Render(DetailLabelStyle.Render("Leechers: ")+DetailValueStyle.Render(t.Leechers)) + "\n"
	s += DetailStyle.Render(DetailLabelStyle.Render("URL: ")+DetailValueStyle.Render(t.Url)) + "\n"

	if t.Magnet != "" {
		s += DetailStyle.Render(DetailLabelStyle.Render("Magnet: ")+DetailValueStyle.Render(t.Magnet)) + "\n"
	}

	s += "\n" + HelpStyle.Render("Press Esc to go back")
	s += "\n" + NavigationStyle.Render("← Back  → Forward")
	return s
}

// performSearch performs the torrent search
func (m Model) performSearch() tea.Cmd {
	return func() tea.Msg {
		torrents := m.searchService.GetAllTorrents(m.lastSearch)
		return SearchResultMsg{Torrents: torrents}
	}
}

// debouncedSearch performs a debounced search
func (m Model) debouncedSearch(query string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond) // Debounce delay
		return DebouncedSearchMsg{Query: query}
	}
}
