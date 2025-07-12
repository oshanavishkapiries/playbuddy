package services

import (
	"sync"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
	"github.com/oshanavishkapiries/playbuddy/src/internal/providers"
)

// SearchService handles torrent searches across multiple providers
type SearchService struct {
	providers []models.Provider
}

// NewSearchService creates a new search service with default providers
func NewSearchService() *SearchService {
	return &SearchService{
		providers: []models.Provider{
			providers.NewPirateBayProvider(),
		},
	}
}

// AddProvider adds a new provider to the search service
func (s *SearchService) AddProvider(provider models.Provider) {
	s.providers = append(s.providers, provider)
}

// SearchAll searches all providers in parallel
func (s *SearchService) SearchAll(query string) []models.SearchResult {
	var wg sync.WaitGroup
	results := make([]models.SearchResult, len(s.providers))

	// Start a goroutine for each provider
	for i, provider := range s.providers {
		wg.Add(1)
		go func(index int, p models.Provider) {
			defer wg.Done()

			torrents, err := p.Search(query)
			result := models.SearchResult{
				Provider: p.GetName(),
			}

			if err != nil {
				result.Error = err.Error()
			} else {
				result.Torrents = torrents
			}

			results[index] = result
		}(i, provider)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return results
}

// GetAllTorrents returns all torrents from all providers as a flat list
func (s *SearchService) GetAllTorrents(query string) []models.Torrent {
	results := s.SearchAll(query)
	var allTorrents []models.Torrent

	for _, result := range results {
		if result.Error == "" {
			allTorrents = append(allTorrents, result.Torrents...)
		}
	}

	return allTorrents
}
