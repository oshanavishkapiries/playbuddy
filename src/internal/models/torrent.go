package models

// Torrent represents a torrent search result
type Torrent struct {
	Name         string `json:"Name"`
	Size         string `json:"Size"`
	DateUploaded string `json:"DateUploaded"`
	Category     string `json:"Category"`
	Seeders      string `json:"Seeders"`
	Leechers     string `json:"Leechers"`
	UploadedBy   string `json:"UploadedBy"`
	Url          string `json:"Url"`
	Magnet       string `json:"Magnet"`
	TorrentFile  string `json:"TorrentFile"`	
}

// Provider represents a torrent provider interface
type Provider interface {
	Search(query string) ([]Torrent, error)
	GetName() string
}

// SearchResult represents the result of a torrent search
type SearchResult struct {
	Provider string    `json:"provider"`
	Torrents []Torrent `json:"torrents"`
	Error    string    `json:"error,omitempty"`
}
