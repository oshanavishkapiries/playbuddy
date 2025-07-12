package services

import (
	"context"
	"sync"
	"time"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
	"github.com/oshanavishkapiries/playbuddy/src/internal/providers"
)

// SearchService handles torrent searches across multiple providers
type SearchService struct {
	providers []models.Provider
	timeout   time.Duration
}

// NewSearchService creates a new search service with default providers
func NewSearchService() *SearchService {
	return &SearchService{
		providers: []models.Provider{
			providers.NewPirateBayProvider(),
			providers.NewYTSProvider(),
			providers.NewNyaaSiProvider(),
		},
		timeout: 15 * time.Second, // 15 second timeout for all providers
	}
}

// AddProvider adds a new provider to the search service
func (s *SearchService) AddProvider(provider models.Provider) {
	s.providers = append(s.providers, provider)
}

// SearchAll searches all providers in parallel with timeout and error handling
func (s *SearchService) SearchAll(query string) []models.SearchResult {
	var wg sync.WaitGroup
	results := make([]models.SearchResult, len(s.providers))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	// Start a goroutine for each provider
	for i, provider := range s.providers {
		wg.Add(1)
		go func(index int, p models.Provider) {
			defer wg.Done()

			// Create a channel for the result
			resultChan := make(chan models.SearchResult, 1)

			// Start the search in a goroutine
			go func() {
				torrents, err := p.Search(query)
				result := models.SearchResult{
					Provider: p.GetName(),
				}

				if err != nil {
					result.Error = err.Error()
				} else {
					result.Torrents = torrents
				}

				resultChan <- result
			}()

			// Wait for either the result or timeout
			select {
			case result := <-resultChan:
				results[index] = result
			case <-ctx.Done():
				results[index] = models.SearchResult{
					Provider: p.GetName(),
					Error:    "Search timeout",
				}
			}
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
			// Add provider information to each torrent
			for i := range result.Torrents {
				result.Torrents[i].Provider = result.Provider
			}
			allTorrents = append(allTorrents, result.Torrents...)
		}
	}

	return allTorrents
}

// GetSearchStats returns statistics about the search results
func (s *SearchService) GetSearchStats(query string) map[string]interface{} {
	results := s.SearchAll(query)
	stats := map[string]interface{}{
		"total_providers":      len(s.providers),
		"successful_providers": 0,
		"failed_providers":     0,
		"total_torrents":       0,
		"provider_results":     make(map[string]interface{}),
	}

	for _, result := range results {
		if result.Error == "" {
			stats["successful_providers"] = stats["successful_providers"].(int) + 1
			stats["total_torrents"] = stats["total_torrents"].(int) + len(result.Torrents)
			stats["provider_results"].(map[string]interface{})[result.Provider] = map[string]interface{}{
				"torrents": len(result.Torrents),
				"status":   "success",
			}
		} else {
			stats["failed_providers"] = stats["failed_providers"].(int) + 1
			stats["provider_results"].(map[string]interface{})[result.Provider] = map[string]interface{}{
				"error":  result.Error,
				"status": "failed",
			}
		}
	}

	return stats
}
