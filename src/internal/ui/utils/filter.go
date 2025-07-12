package utils

import (
	"fmt"
	"strings"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
)

// FilterState represents the current filter state
type FilterState struct {
	ActiveProvider     string // Empty string means no filter
	AvailableProviders []string
}

// FilterManager handles filtering operations
type FilterManager struct {
	state FilterState
}

// NewFilterManager creates a new filter manager
func NewFilterManager() *FilterManager {
	return &FilterManager{
		state: FilterState{
			ActiveProvider:     "",
			AvailableProviders: []string{},
		},
	}
}

// GetAvailableProviders extracts unique providers from results
func GetAvailableProviders(results []models.Torrent) []string {
	providers := make(map[string]bool)
	for _, torrent := range results {
		providers[torrent.Provider] = true
	}

	var uniqueProviders []string
	for provider := range providers {
		uniqueProviders = append(uniqueProviders, provider)
	}
	return uniqueProviders
}

// FilterResults filters results by provider if a filter is active
func FilterResults(results []models.Torrent, filterState FilterState) []models.Torrent {
	if filterState.ActiveProvider == "" {
		return results
	}

	var filtered []models.Torrent
	for _, torrent := range results {
		if strings.EqualFold(torrent.Provider, filterState.ActiveProvider) {
			filtered = append(filtered, torrent)
		}
	}
	return filtered
}

// GetFilterOptions returns numbered filter options
func GetFilterOptions(availableProviders []string) map[string]string {
	options := make(map[string]string)

	// Add "All" option
	options["0"] = "All"

	// Add numbered provider options
	for i, provider := range availableProviders {
		options[fmt.Sprintf("%d", i+1)] = provider
	}

	return options
}

// GetProviderBadge returns a styled badge for the provider
func GetProviderBadge(provider string) string {
	switch strings.ToLower(provider) {
	case "piratebay":
		return "[🏴‍☠️ PB]"
	case "1337x":
		return "[🔢 1337]"
	case "rarbg":
		return "[🎯 RB]"
	case "yts":
		return "[🎬 YTS]"
	case "nyaasi":
		return "[🐱 NY]"
	case "nyaa":
		return "[🐱 NY]"
	default:
		return "[📦 " + strings.ToUpper(provider[:2]) + "]"
	}
}

// GetFilterDisplayText returns formatted filter options with numbers
func GetFilterDisplayText(availableProviders []string, activeProvider string) string {
	var display string

	// Add "All" option
	if activeProvider == "" {
		display += "[0] All "
	} else {
		display += "0 All "
	}

	// Add numbered provider options
	for i, provider := range availableProviders {
		if provider == activeProvider {
			display += fmt.Sprintf("[%d] %s ", i+1, GetProviderBadge(provider))
		} else {
			display += fmt.Sprintf("%d %s ", i+1, GetProviderBadge(provider))
		}
	}

	return display
}

// SetFilterByNumber sets the active filter based on number selection
func SetFilterByNumber(availableProviders []string, number string) string {
	if number == "0" {
		return "" // Clear filter
	}

	// Convert number to index
	index := 0
	fmt.Sscanf(number, "%d", &index)

	// Adjust for 1-based indexing
	index--

	if index >= 0 && index < len(availableProviders) {
		return availableProviders[index]
	}

	return "" // Invalid number, clear filter
}
