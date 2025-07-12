package providers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
)

// ExampleProvider is a template for adding new torrent providers
type ExampleProvider struct {
	baseURL string
	client  *http.Client
}

// NewExampleProvider creates a new example provider
func NewExampleProvider() *ExampleProvider {
	return &ExampleProvider{
		baseURL: "http://localhost:3001/api/example",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Search searches for torrents on the example provider
func (p *ExampleProvider) Search(query string) ([]models.Torrent, error) {
	// This is a template - implement actual API calls here
	// Example implementation:
	// 1. Make HTTP request to the provider's API
	// 2. Parse the response
	// 3. Convert to our Torrent model
	// 4. Return results

	return []models.Torrent{}, fmt.Errorf("example provider not implemented")
}

// GetName returns the provider name
func (p *ExampleProvider) GetName() string {
	return "ExampleProvider"
}

/*
To add a new provider:

1. Create a new file in internal/providers/ (e.g., rarbg.go)
2. Implement the Provider interface
3. Add the provider to the SearchService in internal/services/search.go

Example:
```go
// In internal/services/search.go, add to NewSearchService():
func NewSearchService() *SearchService {
	return &SearchService{
		providers: []models.Provider{
			providers.NewPirateBayProvider(),
			providers.NewRarbgProvider(), // Add your new provider here
		},
	}
}
```
*/
