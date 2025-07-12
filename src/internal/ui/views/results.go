package views

import (
	"fmt"
	"strings"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/sharedstyles"
	"github.com/oshanavishkapiries/playbuddy/src/internal/ui/utils"
)

func ViewResults(results []models.Torrent, cursor int, truncateString func(string, int) string, filterState utils.FilterState) string {

	s := sharedstyles.AsciiStyle.Render(sharedstyles.AsciiArt) + "\n"
	s += sharedstyles.VersionStyle.Render(sharedstyles.Version) + "\n\n"

	// Update available providers if not set
	if len(filterState.AvailableProviders) == 0 {
		filterState.AvailableProviders = utils.GetAvailableProviders(results)
	}

	// Apply filter
	filteredResults := utils.FilterResults(results, filterState)

	if len(filteredResults) == 0 {
		if filterState.ActiveProvider != "" {
			s += sharedstyles.ResultStyle.Render(fmt.Sprintf("No results found for provider: %s", filterState.ActiveProvider)) + "\n"
		} else {
			s += sharedstyles.ResultStyle.Render("No results found") + "\n"
		}
	} else {

		// Calculate visible window (show 8 items at a time)
		windowSize := 6
		startIndex := 0
		endIndex := len(filteredResults)

		// If we have more results than window size, implement scrolling
		if len(filteredResults) > windowSize {
			// Calculate the center of the window around the cursor
			halfWindow := windowSize / 2
			startIndex = cursor - halfWindow
			endIndex = cursor + halfWindow + 1

			// Adjust if we're near the beginning
			if startIndex < 0 {
				startIndex = 0
				endIndex = windowSize
			}

			// Adjust if we're near the end
			if endIndex > len(filteredResults) {
				endIndex = len(filteredResults)
				startIndex = endIndex - windowSize
			}
		}

		// Find the maximum title length for the visible items
		maxTitleLength := 0
		for i := startIndex; i < endIndex; i++ {
			titleLength := len(truncateString(filteredResults[i].Name, 40)) // Reduced to make room for provider
			if titleLength > maxTitleLength {
				maxTitleLength = titleLength
			}
		}

		// Add header
		header := fmt.Sprintf("%-*s %-8s %s", maxTitleLength+2, "Title", "Provider", "Size")
		s += sharedstyles.DetailLabelStyle.Render(header) + "\n"
		s += sharedstyles.DetailLabelStyle.Render(strings.Repeat("-", len(header))) + "\n"

		// Display visible items
		for i := startIndex; i < endIndex; i++ {
			icon := "[🧲]"
			if cursor == i {
				icon = "[➡️]"
			}

			// Format title, provider badge, and size in columns
			title := truncateString(filteredResults[i].Name, 40)
			providerBadge := utils.GetProviderBadge(filteredResults[i].Provider)
			row := fmt.Sprintf("%s %-*s %-8s %s", icon, maxTitleLength, title, providerBadge, filteredResults[i].Size)

			if cursor == i {
				s += sharedstyles.SelectedResultStyle.Render(row) + "\n"
			} else {
				s += sharedstyles.ResultStyle.Render(row) + "\n"
			}
		}

		// Show current position indicator
		if len(filteredResults) > windowSize {
			position := fmt.Sprintf("Showing %d-%d of %d results (Position: %d)",
				startIndex+1, endIndex, len(filteredResults), cursor+1)
			s += sharedstyles.HelpStyle.Render(position) + "\n"
		}
	}

	// Show available filters with numbers
	if len(filterState.AvailableProviders) > 1 {
		s += "\n" + sharedstyles.HelpStyle.Render("Filter options:") + "\n"
		filterDisplay := utils.GetFilterDisplayText(filterState.AvailableProviders, filterState.ActiveProvider)
		s += sharedstyles.ResultStyle.Render(filterDisplay) + "\n"
	}

	s += "\n" + sharedstyles.NavigationStyle.Render("Use arrow keys to navigate, Enter to select, Esc to go back, 0-9 to filter")
	return s
}
