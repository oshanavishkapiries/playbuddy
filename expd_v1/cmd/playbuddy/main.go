package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oshanavishkapiries/playbuddy/internal/models"
)

const (
	TMDB_API_KEY = "6b4357c41d9c606e4d7ebe2f4a8850ea"
	TMDB_BASE_URL = "https://api.themoviedb.org/3"
)

// Application states
type AppState int

const (
	StateWelcome AppState = iota
	StateMediaTypeSelection
	StateSearch
	StateResults
	StateSeasonSelection
	StateEpisodeSelection
	StateStreaming
)

// Media types
type MediaType int

const (
	MediaMovie MediaType = iota
	MediaTVSeries
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center).
		Padding(1, 2)

	asciiStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Align(lipgloss.Center)

	versionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center).
		Padding(0, 2)

	inputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(0, 1)

	menuStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 2)

	selectedMenuStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("86")).
		Padding(0, 1)

	resultStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 2)

	selectedResultStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("86")).
		Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(1, 2)

	// Minimal streaming styles
	providerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		Padding(0, 1)

	linkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 2)

	separatorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	// Professional CLI styles
	selectedLinkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("86")).
		Bold(true).
		Padding(0, 1)

	unselectedLinkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 1)

	urlDisplayStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true).
		Padding(0, 1)

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true)

	warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	// Professional header style
	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(0, 2).
		Margin(0, 2)

	// Professional subtitle style
	subtitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("244")).
		Italic(true).
		Align(lipgloss.Center).
		Padding(0, 1)
)

// ASCII Art for PlayBuddy
const asciiArt = `
██████╗ ██╗      █████╗ ██╗   ██╗██████╗ ██╗   ██╗██████╗ ██████╗ ██╗   ██╗
██╔══██╗██║     ██╔══██╗╚██╗ ██╔╝██╔══██╗██║   ██║██╔══██╗██╔══██╗╚██╗ ██╔╝
██████╔╝██║     ███████║ ╚████╔╝ ██████╔╝██║   ██║██║  ██║██║  ██║ ╚████╔╝ 
██╔═══╝ ██║     ██╔══██║  ╚██╔╝  ██╔══██╗██║   ██║██║  ██║██║  ██║  ╚██╔╝  
██║     ███████╗██║  ██║   ██║   ██████╔╝╚██████╔╝██████╔╝██████╔╝   ██║   
╚═╝     ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝    ╚═╝   
`

// Model for the application
type Model struct {
	state       AppState
	mediaType   MediaType
	input       string
	results     []models.Movie
	tvResults   []models.TVShow
	cursor      int
	loading     bool
	error       string
	selectedID  int
	seasons     []models.Season
	episodes    []models.Episode
	selectedSeason int
	selectedEpisode int
	streamingURL string
	lastSearch  string
	searchDebounce time.Time
}

// Initialize the model
func initialModel() Model {
	return Model{
		state:   StateWelcome,
		cursor:  0,
		loading: false,
	}
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
			if msg.String() == "enter" || msg.String() == " " {
				m.state = StateMediaTypeSelection
				m.cursor = 0
			} else if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}

		case StateMediaTypeSelection:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < 1 {
					m.cursor++
				}
			case "enter":
				m.mediaType = MediaType(m.cursor)
				m.state = StateSearch
				m.cursor = 0
				m.input = ""
			case "q":
				return m, tea.Quit
			case "esc":
				m.state = StateWelcome
				m.cursor = 0
			case "ctrl+c":
				return m, tea.Quit
			}

		case StateSearch:
			switch msg.String() {
			case "enter":
				// Enter now just moves to results if we have any
				if len(m.results) > 0 || len(m.tvResults) > 0 {
					m.state = StateResults
					m.cursor = 0
				}
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
					// Trigger real-time search on backspace
					if strings.TrimSpace(m.input) != "" && strings.TrimSpace(m.input) != m.lastSearch {
						m.searchDebounce = time.Now()
						m.loading = true
						return m, m.performDebouncedSearch()
					} else if strings.TrimSpace(m.input) == "" {
						// Clear results when input is empty
						m.results = []models.Movie{}
						m.tvResults = []models.TVShow{}
						m.loading = false
						m.lastSearch = ""
					}
				}
			case "esc":
				m.state = StateMediaTypeSelection
				m.cursor = int(m.mediaType)
				m.input = ""
				m.results = []models.Movie{}
				m.tvResults = []models.TVShow{}
				m.lastSearch = ""
			case "q":
				return m, tea.Quit
			case "ctrl+c":
				return m, tea.Quit
			default:
				if len(msg.String()) == 1 {
					m.input += msg.String()
					// Trigger real-time search as user types with debouncing
					if strings.TrimSpace(m.input) != "" && strings.TrimSpace(m.input) != m.lastSearch {
						m.searchDebounce = time.Now()
						m.loading = true
						return m, m.performDebouncedSearch()
					}
				}
			}

		case StateResults:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				maxResults := len(m.results)
				if m.mediaType == MediaTVSeries {
					maxResults = len(m.tvResults)
				}
				if m.cursor < maxResults-1 {
					m.cursor++
				}
			case "enter":
				if m.mediaType == MediaMovie && len(m.results) > 0 {
					m.selectedID = m.results[m.cursor].ID
					m.streamingURL = m.generateMovieStreamURL(m.selectedID)
					m.state = StateStreaming
					m.cursor = 0  // Reset cursor for link selection
				} else if m.mediaType == MediaTVSeries && len(m.tvResults) > 0 {
					m.selectedID = m.tvResults[m.cursor].ID
					m.state = StateSeasonSelection
					m.cursor = 0
					return m, m.fetchSeasons()
				}
			case "q":
				return m, tea.Quit
			case "esc":
				m.state = StateSearch
				m.cursor = 0
			case "ctrl+c":
				return m, tea.Quit
			}

		case StateSeasonSelection:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.seasons)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.seasons) > 0 {
					m.selectedSeason = m.seasons[m.cursor].SeasonNumber
					m.state = StateEpisodeSelection
					m.cursor = 0
					return m, m.fetchEpisodes()
				}
			case "q":
				return m, tea.Quit
			case "esc":
				m.state = StateResults
				m.cursor = 0
			case "ctrl+c":
				return m, tea.Quit
			}

		case StateEpisodeSelection:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.episodes)-1 {
					m.cursor++
				}
			case "enter":
				if len(m.episodes) > 0 {
					m.selectedEpisode = m.episodes[m.cursor].EpisodeNumber
					m.streamingURL = m.generateTVStreamURL(m.selectedID, m.selectedSeason, m.selectedEpisode)
					m.state = StateStreaming
					m.cursor = 0  // Reset cursor for link selection
				}
			case "q":
				return m, tea.Quit
			case "esc":
				m.state = StateSeasonSelection
				m.cursor = 0
			case "ctrl+c":
				return m, tea.Quit
			}

		case StateStreaming:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				urls := m.generateAllStreamingURLs()
				if m.cursor < len(urls)-1 {
					m.cursor++
				}
			case "enter", " ":
				urls := m.generateAllStreamingURLs()
				if len(urls) > 0 && m.cursor < len(urls) {
					return m, m.openInBrowser(urls[m.cursor].url)
				}
			case "q":
				return m, tea.Quit
			case "esc":
				if m.mediaType == MediaMovie {
					m.state = StateResults
				} else {
					m.state = StateEpisodeSelection
				}
				m.cursor = 0
			case "ctrl+c":
				return m, tea.Quit
			}
		}

	case searchResultMsg:
		m.loading = false
		if m.mediaType == MediaMovie {
			m.results = msg.movies
		} else {
			m.tvResults = msg.tvShows
		}
		m.error = ""
		m.lastSearch = strings.TrimSpace(m.input)

	case searchErrorMsg:
		m.loading = false
		m.error = msg.error

	case seasonsResultMsg:
		m.seasons = msg.seasons
		m.loading = false

	case episodesResultMsg:
		m.episodes = msg.episodes
		m.loading = false

	case browserOpenMsg:
		if !msg.success {
			m.error = "Failed to open browser: " + msg.error
		} else {
			m.error = "✓ Opening in browser..."
		}
	}

	return m, nil
}

// View method for Bubble Tea
func (m Model) View() string {
	var b strings.Builder

	switch m.state {
	case StateWelcome:
		b.WriteString(asciiStyle.Render(asciiArt))
		b.WriteString("\n")
		b.WriteString(versionStyle.Render("v1.0.0"))
		b.WriteString("\n\n")
		b.WriteString(titleStyle.Render("🎬 Welcome to PlayBuddy - Your Movie & TV Streaming Companion"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press ENTER or SPACE to start • Press Q to quit"))

	case StateMediaTypeSelection:
		b.WriteString(titleStyle.Render("🎯 Choose What You Want to Watch"))
		b.WriteString("\n\n")
		
		options := []string{"🎬 Movies", "📺 TV Series"}
		for i, option := range options {
			if i == m.cursor {
				b.WriteString(selectedMenuStyle.Render("> " + option))
			} else {
				b.WriteString(menuStyle.Render("  " + option))
			}
			b.WriteString("\n")
		}
		
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("↑/↓ navigate • ENTER select • ESC back • Q quit"))

	case StateSearch:
		mediaTypeStr := "Movies"
		if m.mediaType == MediaTVSeries {
			mediaTypeStr = "TV Series"
		}
		
		b.WriteString(titleStyle.Render(fmt.Sprintf("🔍 Search %s", mediaTypeStr)))
		b.WriteString("\n\n")
		
		prompt := fmt.Sprintf("Type %s name: ", strings.ToLower(mediaTypeStr))
		input := inputStyle.Render(prompt + m.input + "█")
		b.WriteString(input)
		b.WriteString("\n\n")
		
		// Show real-time search results
		if m.loading {
			b.WriteString("🔍 Searching...")
			b.WriteString("\n")
		} else if m.error != "" {
			b.WriteString(errorStyle.Render("❌ " + m.error))
			b.WriteString("\n")
		} else if strings.TrimSpace(m.input) != "" {
			// Show results in real-time
			if m.mediaType == MediaMovie && len(m.results) > 0 {
				b.WriteString(versionStyle.Render("Results:"))
				b.WriteString("\n")
				for i, movie := range m.results {
					if i >= 5 { break } // Limit to 5 results in search view
					line := fmt.Sprintf("• %s", movie.Title)
					if movie.ReleaseDate != "" {
						line += fmt.Sprintf(" (%s)", movie.GetReleaseYear())
					}
					b.WriteString(resultStyle.Render(line))
					b.WriteString("\n")
				}
			} else if m.mediaType == MediaTVSeries && len(m.tvResults) > 0 {
				b.WriteString(versionStyle.Render("Results:"))
				b.WriteString("\n")
				for i, show := range m.tvResults {
					if i >= 5 { break } // Limit to 5 results in search view
					line := fmt.Sprintf("• %s", show.Name)
					if show.FirstAirDate != "" {
						line += fmt.Sprintf(" (%s)", show.GetFirstAirYear())
					}
					b.WriteString(resultStyle.Render(line))
					b.WriteString("\n")
				}
			} else if strings.TrimSpace(m.input) != "" && !m.loading {
				b.WriteString(versionStyle.Render("No results found"))
				b.WriteString("\n")
			}
		}
		
		b.WriteString("\n")
		if len(m.results) > 0 || len(m.tvResults) > 0 {
			b.WriteString(helpStyle.Render("Type to search • ENTER to select • ESC back • Q quit"))
		} else {
			b.WriteString(helpStyle.Render("Type to search • ESC back • Q quit"))
		}

	case StateResults:
		if m.loading {
			b.WriteString("🔍 Searching...")
			b.WriteString("\n\n")
		} else if m.error != "" {
			b.WriteString(errorStyle.Render("❌ " + m.error))
			b.WriteString("\n\n")
		} else {
			mediaTypeStr := "Movies"
			if m.mediaType == MediaTVSeries {
				mediaTypeStr = "TV Series"
			}
			
			b.WriteString(titleStyle.Render(fmt.Sprintf("📊 %s Results", mediaTypeStr)))
			b.WriteString("\n\n")
			
			if m.mediaType == MediaMovie && len(m.results) > 0 {
				for i, movie := range m.results {
					if i >= 10 { break } // Limit to 10 results
					line := fmt.Sprintf("%d. %s", i+1, movie.Title)
					if movie.ReleaseDate != "" {
						line += fmt.Sprintf(" (%s)", movie.GetReleaseYear())
					}
					line += fmt.Sprintf(" - ⭐ %.1f", movie.VoteAverage)
					
					if i == m.cursor {
						line = selectedResultStyle.Render(line)
					} else {
						line = resultStyle.Render(line)
					}
					
					b.WriteString(line)
					b.WriteString("\n")
				}
			} else if m.mediaType == MediaTVSeries && len(m.tvResults) > 0 {
				for i, show := range m.tvResults {
					if i >= 10 { break } // Limit to 10 results
					line := fmt.Sprintf("%d. %s", i+1, show.Name)
					if show.FirstAirDate != "" {
						line += fmt.Sprintf(" (%s)", show.GetFirstAirYear())
					}
					line += fmt.Sprintf(" - ⭐ %.1f", show.VoteAverage)
					
					if i == m.cursor {
						line = selectedResultStyle.Render(line)
					} else {
						line = resultStyle.Render(line)
					}
					
					b.WriteString(line)
					b.WriteString("\n")
				}
			} else {
				b.WriteString("No results found.")
			}
			
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("↑/↓ navigate • ENTER select • ESC back • Q quit"))
		}

	case StateSeasonSelection:
		if m.loading {
			b.WriteString("🔍 Loading seasons...")
		} else {
			b.WriteString(titleStyle.Render("📺 Select Season"))
			b.WriteString("\n\n")
			
			for i, season := range m.seasons {
				line := fmt.Sprintf("Season %d (%d episodes)", season.SeasonNumber, season.EpisodeCount)
				
				if i == m.cursor {
					line = selectedResultStyle.Render(line)
				} else {
					line = resultStyle.Render(line)
				}
				
				b.WriteString(line)
				b.WriteString("\n")
			}
			
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("↑/↓ navigate • ENTER select • ESC back • Q quit"))
		}

	case StateEpisodeSelection:
		if m.loading {
			b.WriteString("🔍 Loading episodes...")
		} else {
			b.WriteString(titleStyle.Render(fmt.Sprintf("📺 Season %d Episodes", m.selectedSeason)))
			b.WriteString("\n\n")
			
			for i, episode := range m.episodes {
				line := fmt.Sprintf("E%d: %s", episode.EpisodeNumber, episode.Name)
				
				if i == m.cursor {
					line = selectedResultStyle.Render(line)
				} else {
					line = resultStyle.Render(line)
				}
				
				b.WriteString(line)
				b.WriteString("\n")
			}
			
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("↑/↓ navigate • ENTER select • ESC back • Q quit"))
		}

	case StateStreaming:
		// Get media title for context
		var mediaTitle string
		if m.mediaType == MediaMovie && len(m.results) > 0 {
			mediaTitle = m.results[m.cursor].Title
		} else if m.mediaType == MediaTVSeries && len(m.tvResults) > 0 {
			mediaTitle = m.tvResults[m.cursor].Name
			if m.selectedSeason > 0 && m.selectedEpisode > 0 {
				mediaTitle += fmt.Sprintf(" S%02dE%02d", m.selectedSeason, m.selectedEpisode)
			}
		}
		
		b.WriteString(titleStyle.Render("🎬 Streaming Options"))
		if mediaTitle != "" {
			b.WriteString("\n")
			b.WriteString(versionStyle.Render(mediaTitle))
		}
		b.WriteString("\n\n")
		
		// Show error if any
		if m.error != "" {
			if strings.Contains(m.error, "✓") {
				b.WriteString(successStyle.Render(m.error))
			} else {
				b.WriteString(errorStyle.Render("⚠️  " + m.error))
			}
			b.WriteString("\n\n")
		}
		
		// Generate multiple streaming URLs
		urls := m.generateAllStreamingURLs()
		
		// Professional streaming options display
		b.WriteString(versionStyle.Render("Select a streaming provider:"))
		b.WriteString("\n")
		
		// Show current selection status
		if len(urls) > 0 {
			selectedProvider := urls[m.cursor].name
			statusText := fmt.Sprintf("Current: %s (%d of %d)", selectedProvider, m.cursor+1, len(urls))
			b.WriteString(subtitleStyle.Render(statusText))
		}
		b.WriteString("\n\n")
		
		for i, urlInfo := range urls {
			var line string
			
			// Create professional-looking option with icons
			if i == m.cursor {
				// Selected item with arrow and highlighting
				line = fmt.Sprintf("▶  %s", urlInfo.name)
				line = selectedLinkStyle.Render(line)
			} else {
				// Unselected item
				line = fmt.Sprintf("   %s", urlInfo.name)
				line = unselectedLinkStyle.Render(line)
			}
			
			b.WriteString(line)
			b.WriteString("\n")
			
			// Show URL for selected item
			if i == m.cursor {
				urlDisplay := urlDisplayStyle.Render(fmt.Sprintf("   🔗 %s", urlInfo.url))
				b.WriteString(urlDisplay)
				b.WriteString("\n")
			}
			
			b.WriteString("\n")
		}
		
		b.WriteString(separatorStyle.Render("─────────────────────────────────────────────"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("↑/↓ navigate • ENTER/SPACE open in browser • ESC back • Q quit"))
	}

	return b.String()
}

// Streaming URL info
type StreamingURLInfo struct {
	name string
	url  string
}

// Generate all streaming URLs
func (m Model) generateAllStreamingURLs() []StreamingURLInfo {
	var urls []StreamingURLInfo
	
	if m.mediaType == MediaMovie {
		urls = append(urls, StreamingURLInfo{
			name: "111Movies",
			url:  fmt.Sprintf("https://111movies.com/movie/%d", m.selectedID),
		})
		urls = append(urls, StreamingURLInfo{
			name: "Embed.su",
			url:  fmt.Sprintf("https://embed.su/embed/movie/%d", m.selectedID),
		})
		urls = append(urls, StreamingURLInfo{
			name: "Videasy",
			url:  fmt.Sprintf("https://player.videasy.net/movie/%d", m.selectedID),
		})
	} else {
		urls = append(urls, StreamingURLInfo{
			name: "111Movies",
			url:  fmt.Sprintf("https://111movies.com/tv/%d/%d/%d", m.selectedID, m.selectedSeason, m.selectedEpisode),
		})
		urls = append(urls, StreamingURLInfo{
			name: "Embed.su",
			url:  fmt.Sprintf("https://embed.su/embed/tv/%d/%d/%d", m.selectedID, m.selectedSeason, m.selectedEpisode),
		})
	}
	
	return urls
}

// Generate movie streaming URL
func (m Model) generateMovieStreamURL(id int) string {
	return fmt.Sprintf("https://111movies.com/movie/%d", id)
}

// Generate TV streaming URL
func (m Model) generateTVStreamURL(id, season, episode int) string {
	return fmt.Sprintf("https://111movies.com/tv/%d/%d/%d", id, season, episode)
}

// Messages for async operations
type searchResultMsg struct {
	movies  []models.Movie
	tvShows []models.TVShow
}

type searchErrorMsg struct {
	error string
}

type seasonsResultMsg struct {
	seasons []models.Season
}

type episodesResultMsg struct {
	episodes []models.Episode
}

// Perform search
func (m Model) performSearch() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		if m.mediaType == MediaMovie {
			movies, err := searchMovies(m.input)
			if err != nil {
				return searchErrorMsg{error: err.Error()}
			}
			return searchResultMsg{movies: movies}
		} else {
			tvShows, err := searchTVShows(m.input)
			if err != nil {
				return searchErrorMsg{error: err.Error()}
			}
			return searchResultMsg{tvShows: tvShows}
		}
	})
}

// Perform debounced search
func (m Model) performDebouncedSearch() tea.Cmd {
	debounceTime := m.searchDebounce
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		// Only perform search if no new input has been received
		if debounceTime == m.searchDebounce {
			if m.mediaType == MediaMovie {
				movies, err := searchMovies(m.input)
				if err != nil {
					return searchErrorMsg{error: err.Error()}
				}
				return searchResultMsg{movies: movies}
			} else {
				tvShows, err := searchTVShows(m.input)
				if err != nil {
					return searchErrorMsg{error: err.Error()}
				}
				return searchResultMsg{tvShows: tvShows}
			}
		}
		return nil
	})
}

// Fetch seasons
func (m Model) fetchSeasons() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		seasons, err := getTVSeasons(m.selectedID)
		if err != nil {
			return searchErrorMsg{error: err.Error()}
		}
		return seasonsResultMsg{seasons: seasons}
	})
}

// Fetch episodes
func (m Model) fetchEpisodes() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		episodes, err := getTVEpisodes(m.selectedID, m.selectedSeason)
		if err != nil {
			return searchErrorMsg{error: err.Error()}
		}
		return episodesResultMsg{episodes: episodes}
	})
}

// API functions
func searchMovies(query string) ([]models.Movie, error) {
	u, err := url.Parse(TMDB_BASE_URL + "/search/movie")
	if err != nil {
		return nil, err
	}
	
	params := url.Values{}
	params.Add("api_key", TMDB_API_KEY)
	params.Add("query", query)
	params.Add("page", "1")
	u.RawQuery = params.Encode()
	
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var searchResponse models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	
	return searchResponse.Results, nil
}

func searchTVShows(query string) ([]models.TVShow, error) {
	u, err := url.Parse(TMDB_BASE_URL + "/search/tv")
	if err != nil {
		return nil, err
	}
	
	params := url.Values{}
	params.Add("api_key", TMDB_API_KEY)
	params.Add("query", query)
	params.Add("page", "1")
	u.RawQuery = params.Encode()
	
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var searchResponse models.TVSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	
	return searchResponse.Results, nil
}

func getTVSeasons(tvID int) ([]models.Season, error) {
	u, err := url.Parse(fmt.Sprintf("%s/tv/%d", TMDB_BASE_URL, tvID))
	if err != nil {
		return nil, err
	}
	
	params := url.Values{}
	params.Add("api_key", TMDB_API_KEY)
	u.RawQuery = params.Encode()
	
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var tvDetails models.TVDetails
	if err := json.NewDecoder(resp.Body).Decode(&tvDetails); err != nil {
		return nil, err
	}
	
	return tvDetails.Seasons, nil
}

func getTVEpisodes(tvID, seasonNumber int) ([]models.Episode, error) {
	u, err := url.Parse(fmt.Sprintf("%s/tv/%d/season/%d", TMDB_BASE_URL, tvID, seasonNumber))
	if err != nil {
		return nil, err
	}
	
	params := url.Values{}
	params.Add("api_key", TMDB_API_KEY)
	u.RawQuery = params.Encode()
	
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var seasonDetails models.SeasonDetails
	if err := json.NewDecoder(resp.Body).Decode(&seasonDetails); err != nil {
		return nil, err
	}
	
	return seasonDetails.Episodes, nil
}

// Browser opening message
type browserOpenMsg struct {
	success bool
	error   string
}

// Open URL in browser
func (m Model) openInBrowser(url string) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			return browserOpenMsg{success: false, error: "unsupported platform"}
		}
		
		if err != nil {
			return browserOpenMsg{success: false, error: err.Error()}
		}
		return browserOpenMsg{success: true}
	})
}

func main() {
	// Initialize the model
	model := initialModel()
	
	// Start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}