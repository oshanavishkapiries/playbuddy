package ui

import "github.com/oshanavishkapiries/playbuddy/src/internal/models"

// Message types for async operations
type SearchResultMsg struct {
	Torrents []models.Torrent
	Error    string
}

// DebouncedSearchMsg represents a debounced search message
type DebouncedSearchMsg struct {
	Query string
}
