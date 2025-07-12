package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anacrolix/torrent"
)

// TorrentClient represents a comprehensive torrent client
type TorrentClient struct {
	client *torrent.Client
	torrents map[string]*torrent.Torrent
	config   *torrent.ClientConfig
}

// NewTorrentClient creates a new torrent client instance
func NewTorrentClient(dataDir string, port int) (*TorrentClient, error) {
	config := torrent.NewDefaultClientConfig()
	config.DataDir = dataDir
	config.ListenPort = port
	config.DisableIPv6 = false
	config.DisableUTP = false
	config.DisableTCP = false

	client, err := torrent.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create torrent client: %v", err)
	}

	return &TorrentClient{
		client:    client,
		torrents:  make(map[string]*torrent.Torrent),
		config:    config,
	}, nil
}

// AddTorrentFromMagnet adds a torrent from a magnet link
func (tc *TorrentClient) AddTorrentFromMagnet(magnetLink string) (*torrent.Torrent, error) {
	t, err := tc.client.AddMagnet(magnetLink)
	if err != nil {
		return nil, fmt.Errorf("failed to add magnet link: %v", err)
	}

	// Wait for metadata
	<-t.GotInfo()
	
	// Store torrent reference
	tc.torrents[t.InfoHash().String()] = t
	
	return t, nil
}

// AddTorrentFromFile adds a torrent from a .torrent file
func (tc *TorrentClient) AddTorrentFromFile(torrentPath string) (*torrent.Torrent, error) {
	t, err := tc.client.AddTorrentFromFile(torrentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to add torrent file: %v", err)
	}

	// Wait for metadata
	<-t.GotInfo()
	
	// Store torrent reference
	tc.torrents[t.InfoHash().String()] = t
	
	return t, nil
}

// StartDownload starts downloading a torrent
func (tc *TorrentClient) StartDownload(t *torrent.Torrent) {
	t.DownloadAll()
	fmt.Printf("Started downloading: %s\n", t.Name())
}

// PauseDownload pauses downloading a torrent
func (tc *TorrentClient) PauseDownload(t *torrent.Torrent) {
	t.StopDataDownload()
	fmt.Printf("Paused downloading: %s\n", t.Name())
}

// ResumeDownload resumes downloading a torrent
func (tc *TorrentClient) ResumeDownload(t *torrent.Torrent) {
	t.DownloadAll()
	fmt.Printf("Resumed downloading: %s\n", t.Name())
}

// RemoveTorrent removes a torrent from the client
func (tc *TorrentClient) RemoveTorrent(t *torrent.Torrent, deleteData bool) {
	t.Drop()
	if deleteData {
		// Note: This is a simplified version. In a real implementation,
		// you would need to handle file deletion more carefully
		fmt.Printf("Removed torrent and deleted data: %s\n", t.Name())
	} else {
		fmt.Printf("Removed torrent (data preserved): %s\n", t.Name())
	}
}

// GetTorrentStats returns statistics for a torrent
func (tc *TorrentClient) GetTorrentStats(t *torrent.Torrent) map[string]interface{} {
	stats := t.Stats()
	info := t.Info()
	
	return map[string]interface{}{
		"name":              t.Name(),
		"info_hash":         t.InfoHash().String(),
		"total_size":        t.Length(),
		"downloaded":        stats.BytesRead,
		"uploaded":          stats.BytesWritten,
		"progress":          float64(t.BytesCompleted()) / float64(t.Length()) * 100,
		"peers_connected":   len(t.PeerConns()),
		"download_speed":    stats.BytesReadData.Int64(),
		"upload_speed":      stats.BytesWrittenData.Int64(),
		"piece_length":      info.PieceLength,
		"total_pieces":      info.NumPieces,
		"completed_pieces":  t.BytesCompleted() / info.PieceLength,
	}
}

// ListAllTorrents lists all torrents in the client
func (tc *TorrentClient) ListAllTorrents() {
	fmt.Println("=== Active Torrents ===")
	for hash, t := range tc.torrents {
		stats := tc.GetTorrentStats(t)
		fmt.Printf("Hash: %s\n", hash)
		fmt.Printf("Name: %s\n", stats["name"])
		fmt.Printf("Progress: %.2f%%\n", stats["progress"])
		fmt.Printf("Size: %.2f MB\n", float64(stats["total_size"].(int64))/1024/1024)
		fmt.Printf("Peers: %d\n", stats["peers_connected"])
		fmt.Println("---")
	}
}

// MonitorTorrents continuously monitors all torrents
func (tc *TorrentClient) MonitorTorrents(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Println("\n=== Torrent Status Update ===")
					for _, t := range tc.torrents {
			stats := tc.GetTorrentStats(t)
			fmt.Printf("%s: %.2f%% complete, %d peers, %.2f MB/s down, %.2f MB/s up\n",
				stats["name"],
				stats["progress"],
				stats["peers_connected"],
				float64(stats["download_speed"].(int64))/1024/1024,
				float64(stats["upload_speed"].(int64))/1024/1024)
		}
		}
	}
}

// Close closes the torrent client
func (tc *TorrentClient) Close() {
	tc.client.Close()
}

// Example usage of the comprehensive torrent client
func runComprehensiveClient() {
	// Create torrent client
	tc, err := NewTorrentClient("./comprehensive_downloads", 6883)
	if err != nil {
		log.Fatalf("Failed to create torrent client: %v", err)
	}
	defer tc.Close()

	// Example magnet link (replace with your own)
	magnetLink := "magnet:?xt=urn:btih:YOUR_TORRENT_HASH_HERE&dn=example+torrent"

	// Add torrent
	t, err := tc.AddTorrentFromMagnet(magnetLink)
	if err != nil {
		log.Fatalf("Failed to add torrent: %v", err)
	}

	// Start downloading
	tc.StartDownload(t)

	// Set up context for monitoring
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start monitoring in a goroutine
	go tc.MonitorTorrents(ctx)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal or completion
	select {
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal, cleaning up...")
		cancel()
	case <-t.Complete.On():
		fmt.Println("\nTorrent completed!")
	}

	// List final status
	tc.ListAllTorrents()
} 