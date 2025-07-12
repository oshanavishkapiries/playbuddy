package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
)

// YTSMovie represents a movie from YTS API
type YTSMovie struct {
	Name         string    `json:"Name"`
	ReleasedDate string    `json:"ReleasedDate"`
	Genre        string    `json:"Genre"`
	Rating       string    `json:"Rating"`
	Likes        string    `json:"Likes"`
	Runtime      string    `json:"Runtime"`
	Language     string    `json:"Language"`
	Url          string    `json:"Url"`
	Poster       string    `json:"Poster"`
	Files        []YTSFile `json:"Files"`
}

// YTSFile represents a file/torrent from YTS
type YTSFile struct {
	Quality string `json:"Quality"`
	Type    string `json:"Type"`
	Size    string `json:"Size"`
	Torrent string `json:"Torrent"`
	Magnet  string `json:"Magnet"`
}

// YTSProvider implements the Provider interface for YTS
type YTSProvider struct {
	baseURL string
	client  *http.Client
}

// NewYTSProvider creates a new YTS provider
func NewYTSProvider() *YTSProvider {
	return &YTSProvider{
		baseURL: "http://68.183.184.162:5362/api/yts",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Search searches for torrents on YTS
func (p *YTSProvider) Search(query string) ([]models.Torrent, error) {
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
	var movies []YTSMovie
	if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert YTS movies to Torrent models
	var torrents []models.Torrent
	for _, movie := range movies {
		for _, file := range movie.Files {
			torrent := models.Torrent{
				Name:         fmt.Sprintf("%s (%s) [%s]", movie.Name, movie.ReleasedDate, file.Quality),
				Size:         file.Size,
				DateUploaded: movie.ReleasedDate,
				Category:     movie.Genre,
				Seeders:      "N/A", // YTS doesn't provide seeder info
				Leechers:     "N/A", // YTS doesn't provide leecher info
				UploadedBy:   "YTS",
				Url:          movie.Url,
				Magnet:       file.Magnet,
				TorrentFile:  file.Torrent,
				Provider:     "YTS",
			}
			torrents = append(torrents, torrent)
		}
	}

	return torrents, nil
}

// GetName returns the provider name
func (p *YTSProvider) GetName() string {
	return "YTS"
}
