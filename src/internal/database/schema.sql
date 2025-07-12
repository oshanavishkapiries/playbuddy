-- PlayBuddy Database Schema
-- SQLite database for tracking downloads and history

-- Downloads table - tracks active and completed downloads
CREATE TABLE IF NOT EXISTS downloads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    torrent_hash TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    magnet TEXT NOT NULL,
    provider TEXT NOT NULL,
    total_size INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, downloading, completed, paused, failed
    progress REAL DEFAULT 0.0,
    downloaded_bytes INTEGER DEFAULT 0,
    upload_speed INTEGER DEFAULT 0,
    download_speed INTEGER DEFAULT 0,
    peers_connected INTEGER DEFAULT 0,
    seeders INTEGER DEFAULT 0,
    leechers INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    download_path TEXT,
    selected_files TEXT -- JSON array of selected file indices
);

-- Download files table - tracks individual files in torrents
CREATE TABLE IF NOT EXISTS download_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    download_id INTEGER NOT NULL,
    file_index INTEGER NOT NULL,
    file_name TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    file_path TEXT,
    priority TEXT DEFAULT 'normal', -- normal, high, low, none
    downloaded_bytes INTEGER DEFAULT 0,
    is_selected BOOLEAN DEFAULT 1,
    FOREIGN KEY (download_id) REFERENCES downloads(id) ON DELETE CASCADE
);

-- Download history table - for completed downloads
CREATE TABLE IF NOT EXISTS download_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    torrent_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    provider TEXT NOT NULL,
    total_size INTEGER NOT NULL,
    completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    download_path TEXT,
    file_count INTEGER DEFAULT 0
);

-- Settings table - application settings
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_downloads_status ON downloads(status);
CREATE INDEX IF NOT EXISTS idx_downloads_hash ON downloads(torrent_hash);
CREATE INDEX IF NOT EXISTS idx_download_files_download_id ON download_files(download_id);
CREATE INDEX IF NOT EXISTS idx_download_history_hash ON download_history(torrent_hash);

-- Insert default settings
INSERT OR IGNORE INTO settings (key, value) VALUES 
    ('download_directory', './downloads'),
    ('max_concurrent_downloads', '3'),
    ('default_port', '6881'),
    ('auto_start_downloads', 'true'),
    ('max_upload_speed', '0'),
    ('max_download_speed', '0'); 