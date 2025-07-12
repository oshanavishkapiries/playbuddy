package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
)

// NyaaSiTorrent represents a torrent from Nyaa.si API
type NyaaSiTorrent struct {
	Name         string `json:"Name"`
	Category     string `json:"Category"`
	Url          string `json:"Url"`
	Size         string `json:"Size"`
	DateUploaded string `json:"DateUploaded"`
	Seeders      string `json:"Seeders"`
	Leechers     string `json:"Leechers"`
	Downloads    string `json:"Downloads"`
	Torrent      string `json:"Torrent"`
	Magnet       string `json:"Magnet"`
}

// NyaaSiProvider implements the Provider interface for Nyaa.si
type NyaaSiProvider struct {
	baseURL string
	client  *http.Client
}

// NewNyaaSiProvider creates a new Nyaa.si provider
func NewNyaaSiProvider() *NyaaSiProvider {
	return &NyaaSiProvider{
		baseURL: "http://68.183.184.162:5362/api/nyaasi",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Search searches for torrents on Nyaa.si
func (p *NyaaSiProvider) Search(query string) ([]models.Torrent, error) {
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
	var nyaaTorrents []NyaaSiTorrent
	if err := json.NewDecoder(resp.Body).Decode(&nyaaTorrents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert Nyaa.si torrents to Torrent models
	var torrents []models.Torrent
	for _, nyaaTorrent := range nyaaTorrents {
		torrent := models.Torrent{
			Name:         nyaaTorrent.Name,
			Size:         nyaaTorrent.Size,
			DateUploaded: nyaaTorrent.DateUploaded,
			Category:     nyaaTorrent.Category,
			Seeders:      nyaaTorrent.Seeders,
			Leechers:     nyaaTorrent.Leechers,
			UploadedBy:   "Nyaa.si",
			Url:          nyaaTorrent.Url,
			Magnet:       nyaaTorrent.Magnet,
			TorrentFile:  nyaaTorrent.Torrent,
			Provider:     "NyaaSi",
		}
		torrents = append(torrents, torrent)
	}

	return torrents, nil
}

// GetName returns the provider name
func (p *NyaaSiProvider) GetName() string {
	return "NyaaSi"
}
