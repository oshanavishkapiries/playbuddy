package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
)

// PirateBayProvider implements the Provider interface for PirateBay
type PirateBayProvider struct {
	baseURL string
	client  *http.Client
}

// NewPirateBayProvider creates a new PirateBay provider
func NewPirateBayProvider() *PirateBayProvider {
	return &PirateBayProvider{
		baseURL: "http://localhost:3001/api/piratebay",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Search searches for torrents on PirateBay
func (p *PirateBayProvider) Search(query string) ([]models.Torrent, error) {
	// Encode the query for URL
	encodedQuery := url.QueryEscape(query)

	// Construct the URL
	searchURL := fmt.Sprintf("%s/%s", p.baseURL, encodedQuery)

	// Make the HTTP request
	resp, err := p.client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Decode the JSON response
	var torrents []models.Torrent
	if err := json.NewDecoder(resp.Body).Decode(&torrents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return torrents, nil
}

// GetName returns the provider name
func (p *PirateBayProvider) GetName() string {
	return "PirateBay"
}
