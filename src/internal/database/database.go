package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database represents the SQLite database manager
type Database struct {
	db *sql.DB
}

// Download represents a download record
type Download struct {
	ID              int        `json:"id"`
	TorrentHash     string     `json:"torrent_hash"`
	Name            string     `json:"name"`
	Magnet          string     `json:"magnet"`
	Provider        string     `json:"provider"`
	TotalSize       int64      `json:"total_size"`
	Status          string     `json:"status"`
	Progress        float64    `json:"progress"`
	DownloadedBytes int64      `json:"downloaded_bytes"`
	UploadSpeed     int64      `json:"upload_speed"`
	DownloadSpeed   int64      `json:"download_speed"`
	PeersConnected  int        `json:"peers_connected"`
	Seeders         int        `json:"seeders"`
	Leechers        int        `json:"leechers"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	DownloadPath    string     `json:"download_path"`
	SelectedFiles   []int      `json:"selected_files"`
}

// DownloadFile represents a file in a torrent
type DownloadFile struct {
	ID              int    `json:"id"`
	DownloadID      int    `json:"download_id"`
	FileIndex       int    `json:"file_index"`
	FileName        string `json:"file_name"`
	FileSize        int64  `json:"file_size"`
	FilePath        string `json:"file_path"`
	Priority        string `json:"priority"`
	DownloadedBytes int64  `json:"downloaded_bytes"`
	IsSelected      bool   `json:"is_selected"`
}

// DownloadHistory represents a completed download
type DownloadHistory struct {
	ID           int       `json:"id"`
	TorrentHash  string    `json:"torrent_hash"`
	Name         string    `json:"name"`
	Provider     string    `json:"provider"`
	TotalSize    int64     `json:"total_size"`
	CompletedAt  time.Time `json:"completed_at"`
	DownloadPath string    `json:"download_path"`
	FileCount    int       `json:"file_count"`
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Initialize database schema
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return &Database{db: db}, nil
}

// initSchema initializes the database schema
func initSchema(db *sql.DB) error {
	schemaSQL, err := os.ReadFile("internal/database/schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	_, err = db.Exec(string(schemaSQL))
	return err
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// AddDownload adds a new download to the database
func (d *Database) AddDownload(download *Download) error {
	selectedFilesJSON, err := json.Marshal(download.SelectedFiles)
	if err != nil {
		return fmt.Errorf("failed to marshal selected files: %v", err)
	}

	query := `
		INSERT INTO downloads (
			torrent_hash, name, magnet, provider, total_size, status,
			progress, downloaded_bytes, download_path, selected_files
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := d.db.Exec(query,
		download.TorrentHash, download.Name, download.Magnet, download.Provider,
		download.TotalSize, download.Status, download.Progress, download.DownloadedBytes,
		download.DownloadPath, string(selectedFilesJSON))
	if err != nil {
		return fmt.Errorf("failed to add download: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}

	download.ID = int(id)
	return nil
}

// UpdateDownload updates a download record
func (d *Database) UpdateDownload(download *Download) error {
	query := `
		UPDATE downloads SET
			status = ?, progress = ?, downloaded_bytes = ?,
			upload_speed = ?, download_speed = ?, peers_connected = ?,
			seeders = ?, leechers = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := d.db.Exec(query,
		download.Status, download.Progress, download.DownloadedBytes,
		download.UploadSpeed, download.DownloadSpeed, download.PeersConnected,
		download.Seeders, download.Leechers, download.ID)

	return err
}

// GetDownloadByHash gets a download by torrent hash
func (d *Database) GetDownloadByHash(hash string) (*Download, error) {
	query := `
		SELECT id, torrent_hash, name, magnet, provider, total_size,
		       status, progress, downloaded_bytes, upload_speed, download_speed,
		       peers_connected, seeders, leechers, created_at, updated_at,
		       completed_at, download_path, selected_files
		FROM downloads WHERE torrent_hash = ?
	`

	download := &Download{}
	var selectedFilesJSON string
	var completedAt sql.NullTime

	err := d.db.QueryRow(query, hash).Scan(
		&download.ID, &download.TorrentHash, &download.Name, &download.Magnet,
		&download.Provider, &download.TotalSize, &download.Status, &download.Progress,
		&download.DownloadedBytes, &download.UploadSpeed, &download.DownloadSpeed,
		&download.PeersConnected, &download.Seeders, &download.Leechers,
		&download.CreatedAt, &download.UpdatedAt, &completedAt, &download.DownloadPath,
		&selectedFilesJSON)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		download.CompletedAt = &completedAt.Time
	}

	if selectedFilesJSON != "" {
		if err := json.Unmarshal([]byte(selectedFilesJSON), &download.SelectedFiles); err != nil {
			log.Printf("Warning: failed to unmarshal selected files: %v", err)
		}
	}

	return download, nil
}

// GetActiveDownloads gets all active downloads
func (d *Database) GetActiveDownloads() ([]*Download, error) {
	query := `
		SELECT id, torrent_hash, name, magnet, provider, total_size,
		       status, progress, downloaded_bytes, upload_speed, download_speed,
		       peers_connected, seeders, leechers, created_at, updated_at,
		       completed_at, download_path, selected_files
		FROM downloads WHERE status IN ('pending', 'downloading', 'paused')
		ORDER BY created_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var downloads []*Download
	for rows.Next() {
		download := &Download{}
		var selectedFilesJSON string
		var completedAt sql.NullTime

		err := rows.Scan(
			&download.ID, &download.TorrentHash, &download.Name, &download.Magnet,
			&download.Provider, &download.TotalSize, &download.Status, &download.Progress,
			&download.DownloadedBytes, &download.UploadSpeed, &download.DownloadSpeed,
			&download.PeersConnected, &download.Seeders, &download.Leechers,
			&download.CreatedAt, &download.UpdatedAt, &completedAt, &download.DownloadPath,
			&selectedFilesJSON)

		if err != nil {
			return nil, err
		}

		if completedAt.Valid {
			download.CompletedAt = &completedAt.Time
		}

		if selectedFilesJSON != "" {
			if err := json.Unmarshal([]byte(selectedFilesJSON), &download.SelectedFiles); err != nil {
				log.Printf("Warning: failed to unmarshal selected files: %v", err)
			}
		}

		downloads = append(downloads, download)
	}

	return downloads, nil
}

// GetDownloadHistory gets download history
func (d *Database) GetDownloadHistory(limit int) ([]*DownloadHistory, error) {
	query := `
		SELECT id, torrent_hash, name, provider, total_size, completed_at, download_path, file_count
		FROM download_history ORDER BY completed_at DESC LIMIT ?
	`

	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*DownloadHistory
	for rows.Next() {
		item := &DownloadHistory{}
		err := rows.Scan(
			&item.ID, &item.TorrentHash, &item.Name, &item.Provider,
			&item.TotalSize, &item.CompletedAt, &item.DownloadPath, &item.FileCount)
		if err != nil {
			return nil, err
		}
		history = append(history, item)
	}

	return history, nil
}

// MarkDownloadCompleted marks a download as completed
func (d *Database) MarkDownloadCompleted(downloadID int, fileCount int) error {
	// Update download status
	query := `
		UPDATE downloads SET
			status = 'completed', progress = 100.0, completed_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := d.db.Exec(query, downloadID)
	if err != nil {
		return err
	}

	// Get download info for history
	download, err := d.GetDownloadByID(downloadID)
	if err != nil {
		return err
	}

	// Add to history
	historyQuery := `
		INSERT INTO download_history (torrent_hash, name, provider, total_size, download_path, file_count)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = d.db.Exec(historyQuery,
		download.TorrentHash, download.Name, download.Provider,
		download.TotalSize, download.DownloadPath, fileCount)

	return err
}

// GetDownloadByID gets a download by ID
func (d *Database) GetDownloadByID(id int) (*Download, error) {
	query := `
		SELECT id, torrent_hash, name, magnet, provider, total_size,
		       status, progress, downloaded_bytes, upload_speed, download_speed,
		       peers_connected, seeders, leechers, created_at, updated_at,
		       completed_at, download_path, selected_files
		FROM downloads WHERE id = ?
	`

	download := &Download{}
	var selectedFilesJSON string
	var completedAt sql.NullTime

	err := d.db.QueryRow(query, id).Scan(
		&download.ID, &download.TorrentHash, &download.Name, &download.Magnet,
		&download.Provider, &download.TotalSize, &download.Status, &download.Progress,
		&download.DownloadedBytes, &download.UploadSpeed, &download.DownloadSpeed,
		&download.PeersConnected, &download.Seeders, &download.Leechers,
		&download.CreatedAt, &download.UpdatedAt, &completedAt, &download.DownloadPath,
		&selectedFilesJSON)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		download.CompletedAt = &completedAt.Time
	}

	if selectedFilesJSON != "" {
		if err := json.Unmarshal([]byte(selectedFilesJSON), &download.SelectedFiles); err != nil {
			log.Printf("Warning: failed to unmarshal selected files: %v", err)
		}
	}

	return download, nil
}

// GetSetting gets a setting value
func (d *Database) GetSetting(key string) (string, error) {
	query := `SELECT value FROM settings WHERE key = ?`
	var value string
	err := d.db.QueryRow(query, key).Scan(&value)
	return value, err
}

// SetSetting sets a setting value
func (d *Database) SetSetting(key, value string) error {
	query := `
		INSERT OR REPLACE INTO settings (key, value, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`
	_, err := d.db.Exec(query, key, value)
	return err
}
