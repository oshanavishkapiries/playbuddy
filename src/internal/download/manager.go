package download

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/oshanavishkapiries/playbuddy/src/internal/database"
	"github.com/oshanavishkapiries/playbuddy/src/internal/models"
)

// Manager handles torrent downloads with database integration
type Manager struct {
	client    *torrent.Client
	db        *database.Database
	downloads map[string]*ActiveDownload
	mu        sync.RWMutex
	config    *Config
}

// Config holds download manager configuration
type Config struct {
	DataDir           string
	Port              int
	MaxConcurrent     int
	DownloadDirectory string
}

// ActiveDownload represents an active download with real-time stats
type ActiveDownload struct {
	ID              int
	TorrentHash     string
	Name            string
	Magnet          string
	Provider        string
	TotalSize       int64
	Status          string
	Progress        float64
	DownloadedBytes int64
	UploadSpeed     int64
	DownloadSpeed   int64
	PeersConnected  int
	Seeders         int
	Leechers        int
	Torrent         *torrent.Torrent
	SelectedFiles   []int
	CancelFunc      context.CancelFunc
	mu              sync.RWMutex
}

// NewManager creates a new download manager
func NewManager(config *Config, db *database.Database) (*Manager, error) {
	torrentConfig := torrent.NewDefaultClientConfig()
	torrentConfig.DataDir = config.DataDir
	torrentConfig.ListenPort = config.Port

	client, err := torrent.NewClient(torrentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create torrent client: %v", err)
	}

	// Ensure download directory exists
	if err := os.MkdirAll(config.DownloadDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create download directory: %v", err)
	}

	return &Manager{
		client:    client,
		db:        db,
		downloads: make(map[string]*ActiveDownload),
		config:    config,
	}, nil
}

// StartDownload starts downloading a torrent
func (m *Manager) StartDownload(torrentModel models.Torrent, selectedFiles []int) error {
	// Check if already downloading
	m.mu.Lock()
	if _, exists := m.downloads[torrentModel.Magnet]; exists {
		m.mu.Unlock()
		return fmt.Errorf("torrent already being downloaded")
	}
	m.mu.Unlock()

	// Add torrent to client
	t, err := m.client.AddMagnet(torrentModel.Magnet)
	if err != nil {
		return fmt.Errorf("failed to add magnet link: %v", err)
	}

	// Wait for metadata
	<-t.GotInfo()

	// Create download record
	download := &database.Download{
		TorrentHash:     t.InfoHash().String(),
		Name:            torrentModel.Name,
		Magnet:          torrentModel.Magnet,
		Provider:        torrentModel.Provider,
		TotalSize:       t.Length(),
		Status:          "pending",
		Progress:        0.0,
		DownloadedBytes: 0,
		DownloadPath:    filepath.Join(m.config.DownloadDirectory, torrentModel.Name),
		SelectedFiles:   selectedFiles,
	}

	// Add to database
	if err := m.db.AddDownload(download); err != nil {
		return fmt.Errorf("failed to add download to database: %v", err)
	}

	// Create active download
	activeDownload := &ActiveDownload{
		ID:              download.ID,
		TorrentHash:     download.TorrentHash,
		Name:            download.Name,
		Magnet:          download.Magnet,
		Provider:        download.Provider,
		TotalSize:       download.TotalSize,
		Status:          "downloading",
		Progress:        0.0,
		DownloadedBytes: 0,
		Torrent:         t,
		SelectedFiles:   selectedFiles,
	}

	// Add to active downloads
	m.mu.Lock()
	m.downloads[torrentModel.Magnet] = activeDownload
	m.mu.Unlock()

	// Start download
	if len(selectedFiles) > 0 {
		// Selective download
		files := t.Files()
		for i, file := range files {
			if contains(selectedFiles, i) {
				file.SetPriority(torrent.PiecePriorityNormal)
			} else {
				file.SetPriority(torrent.PiecePriorityNone)
			}
		}
	} else {
		// Download all files
		t.DownloadAll()
	}

	// Start monitoring
	ctx, cancel := context.WithCancel(context.Background())
	activeDownload.CancelFunc = cancel
	go m.monitorDownload(ctx, activeDownload)

	return nil
}

// monitorDownload monitors download progress and updates database
func (m *Manager) monitorDownload(ctx context.Context, download *ActiveDownload) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			stats := download.Torrent.Stats()

			download.mu.Lock()
			download.DownloadedBytes = stats.BytesRead.Int64()
			download.DownloadSpeed = stats.BytesReadData.Int64()
			download.UploadSpeed = stats.BytesWrittenData.Int64()
			download.PeersConnected = len(download.Torrent.PeerConns())

			if download.TotalSize > 0 {
				download.Progress = float64(download.DownloadedBytes) / float64(download.TotalSize) * 100
			}

			// Check if completed
			if download.DownloadedBytes >= download.TotalSize && download.TotalSize > 0 {
				download.Status = "completed"
				download.Progress = 100.0

				// Update database
				dbDownload := &database.Download{
					ID:              download.ID,
					Status:          download.Status,
					Progress:        download.Progress,
					DownloadedBytes: download.DownloadedBytes,
					UploadSpeed:     download.UploadSpeed,
					DownloadSpeed:   download.DownloadSpeed,
					PeersConnected:  download.PeersConnected,
				}

				if err := m.db.UpdateDownload(dbDownload); err != nil {
					log.Printf("Failed to update download: %v", err)
				}

				// Mark as completed in database
				if err := m.db.MarkDownloadCompleted(download.ID, len(download.Torrent.Files())); err != nil {
					log.Printf("Failed to mark download completed: %v", err)
				}

				// Remove from active downloads
				m.mu.Lock()
				delete(m.downloads, download.Magnet)
				m.mu.Unlock()

				return
			}
			download.mu.Unlock()

			// Update database
			dbDownload := &database.Download{
				ID:              download.ID,
				Status:          download.Status,
				Progress:        download.Progress,
				DownloadedBytes: download.DownloadedBytes,
				UploadSpeed:     download.UploadSpeed,
				DownloadSpeed:   download.DownloadSpeed,
				PeersConnected:  download.PeersConnected,
			}

			if err := m.db.UpdateDownload(dbDownload); err != nil {
				log.Printf("Failed to update download: %v", err)
			}
		}
	}
}

// GetActiveDownloads returns all active downloads
func (m *Manager) GetActiveDownloads() []*ActiveDownload {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var downloads []*ActiveDownload
	for _, download := range m.downloads {
		downloads = append(downloads, download)
	}
	return downloads
}

// GetDownloadByHash returns a download by torrent hash
func (m *Manager) GetDownloadByHash(hash string) *ActiveDownload {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, download := range m.downloads {
		if download.TorrentHash == hash {
			return download
		}
	}
	return nil
}

// PauseDownload pauses a download
func (m *Manager) PauseDownload(hash string) error {
	download := m.GetDownloadByHash(hash)
	if download == nil {
		return fmt.Errorf("download not found")
	}

	download.mu.Lock()
	download.Status = "paused"
	download.mu.Unlock()

	// Update database
	dbDownload := &database.Download{
		ID:     download.ID,
		Status: "paused",
	}

	return m.db.UpdateDownload(dbDownload)
}

// ResumeDownload resumes a download
func (m *Manager) ResumeDownload(hash string) error {
	download := m.GetDownloadByHash(hash)
	if download == nil {
		return fmt.Errorf("download not found")
	}

	download.mu.Lock()
	download.Status = "downloading"
	download.mu.Unlock()

	// Resume download
	download.Torrent.DownloadAll()

	// Update database
	dbDownload := &database.Download{
		ID:     download.ID,
		Status: "downloading",
	}

	return m.db.UpdateDownload(dbDownload)
}

// RemoveDownload removes a download
func (m *Manager) RemoveDownload(hash string, deleteData bool) error {
	download := m.GetDownloadByHash(hash)
	if download == nil {
		return fmt.Errorf("download not found")
	}

	// Cancel monitoring
	if download.CancelFunc != nil {
		download.CancelFunc()
	}

	// Remove from client
	download.Torrent.Drop()

	// Remove from active downloads
	m.mu.Lock()
	delete(m.downloads, download.Magnet)
	m.mu.Unlock()

	// Delete data if requested
	if deleteData {
		downloadPath := filepath.Join(m.config.DownloadDirectory, download.Name)
		if err := os.RemoveAll(downloadPath); err != nil {
			log.Printf("Failed to delete download data: %v", err)
		}
	}

	return nil
}

// GetDownloadHistory returns download history from database
func (m *Manager) GetDownloadHistory(limit int) ([]*database.DownloadHistory, error) {
	return m.db.GetDownloadHistory(limit)
}

// RecoverDownloads recovers downloads from database after crash
func (m *Manager) RecoverDownloads() error {
	downloads, err := m.db.GetActiveDownloads()
	if err != nil {
		return fmt.Errorf("failed to get active downloads: %v", err)
	}

	for _, download := range downloads {
		if download.Status == "downloading" || download.Status == "pending" {
			// Try to resume the download
			t, err := m.client.AddMagnet(download.Magnet)
			if err != nil {
				log.Printf("Failed to recover download %s: %v", download.Name, err)
				continue
			}

			// Wait for metadata
			select {
			case <-t.GotInfo():
			case <-time.After(30 * time.Second):
				log.Printf("Timeout waiting for metadata for %s", download.Name)
				continue
			}

			// Create active download
			activeDownload := &ActiveDownload{
				ID:              download.ID,
				TorrentHash:     download.TorrentHash,
				Name:            download.Name,
				Magnet:          download.Magnet,
				Provider:        download.Provider,
				TotalSize:       download.TotalSize,
				Status:          "downloading",
				Progress:        download.Progress,
				DownloadedBytes: download.DownloadedBytes,
				Torrent:         t,
				SelectedFiles:   download.SelectedFiles,
			}

			// Add to active downloads
			m.mu.Lock()
			m.downloads[download.Magnet] = activeDownload
			m.mu.Unlock()

			// Resume download
			t.DownloadAll()

			// Start monitoring
			ctx, cancel := context.WithCancel(context.Background())
			activeDownload.CancelFunc = cancel
			go m.monitorDownload(ctx, activeDownload)

			log.Printf("Recovered download: %s", download.Name)
		}
	}

	return nil
}

// Close closes the download manager
func (m *Manager) Close() error {
	// Cancel all downloads
	m.mu.Lock()
	for _, download := range m.downloads {
		if download.CancelFunc != nil {
			download.CancelFunc()
		}
	}
	m.mu.Unlock()

	// client.Close() returns []error, so we need to handle it properly
	errors := m.client.Close()
	if len(errors) > 0 {
		return fmt.Errorf("client close errors: %v", errors)
	}
	return nil
}

// Helper function to check if slice contains value
func contains(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// FormatBytes formats bytes to human readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatSpeed formats speed to human readable format
func FormatSpeed(bytesPerSecond int64) string {
	return FormatBytes(bytesPerSecond) + "/s"
}
