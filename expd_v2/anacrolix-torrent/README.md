# Anacrolix Torrent Library Examples

This project demonstrates various capabilities of the `github.com/anacrolix/torrent` library in Go.

## Features Demonstrated

### 1. Download Torrent (`01.download_torrent.go`)

- Basic torrent downloading from magnet links
- Progress monitoring with real-time updates
- Download completion detection
- Configurable download directory and ports

### 2. Seed Torrent (`02.seed_torrent.go`)

- Creating and seeding torrents
- Upload monitoring and statistics
- Peer connection tracking
- Torrent file creation (placeholder implementation)

### 3. Torrent Information (`03.torrent_info.go`)

- Detailed torrent metadata extraction
- File information and structure
- Peer statistics and connection details
- Real-time monitoring and statistics

### 4. Selective Download (`04.selective_download.go`)

- Download specific files from a torrent
- File priority management
- Pattern-based file selection
- Progress tracking for selected files

### 5. Comprehensive Client (`05.torrent_client.go`)

- Full-featured torrent client implementation
- Multiple torrent management
- Pause/resume functionality
- Advanced statistics and monitoring
- Graceful shutdown handling

## Getting Started

### Prerequisites

- Go 1.24.4 or later
- Internet connection for downloading dependencies

### Installation

1. Navigate to the project directory:

```bash
cd expd_v2/anacrolix-torrent
```

2. Download dependencies:

```bash
go mod tidy
```

### Usage

#### Running Individual Examples

Each example file can be run independently:

```bash
# Download torrent example
go run 01.download_torrent.go

# Seed torrent example
go run 02.seed_torrent.go

# Get torrent information
go run 03.torrent_info.go

# Selective download example
go run 04.selective_download.go

# Comprehensive client example
go run 05.torrent_client.go
```

#### Before Running Examples

**Important**: You need to replace the placeholder magnet links in each example file with real torrent magnet links:

1. Open any of the example files (e.g., `01.download_torrent.go`)
2. Find the line with the magnet link:
   ```go
   magnetLink := "magnet:?xt=urn:btih:YOUR_TORRENT_HASH_HERE&dn=example+torrent"
   ```
3. Replace it with a real magnet link, for example:
   ```go
   magnetLink := "magnet:?xt=urn:btih:08ada5a7a6183aae1e09d831df6748d566095a10&dn=Sintel"
   ```

## Example Magnet Links for Testing

You can use these legal torrent magnet links for testing:

- **Sintel (Open Movie)**: `magnet:?xt=urn:btih:08ada5a7a6183aae1e09d831df6748d566095a10&dn=Sintel`
- **Big Buck Bunny (Open Movie)**: `magnet:?xt=urn:btih:dd8255ecdc7ca55fb3739c5821ac8cb9e6a42074&dn=Big+Buck+Bunny`

## Configuration

### Download Directories

Each example uses different download directories:

- `01.download_torrent.go`: `./downloads`
- `02.seed_torrent.go`: `./uploads`
- `03.torrent_info.go`: `./info_cache`
- `04.selective_download.go`: `./selective_downloads`
- `05.torrent_client.go`: `./comprehensive_downloads`

### Port Configuration

Each example uses different ports to avoid conflicts:

- `01.download_torrent.go`: Port 6881
- `02.seed_torrent.go`: Port 6882
- `03.torrent_info.go`: Default port
- `04.selective_download.go`: Default port
- `05.torrent_client.go`: Port 6883

## Key Features of the Library

### Torrent Client Configuration

```go
config := torrent.NewDefaultClientConfig()
config.DataDir = "./downloads"
config.ListenPort = 6881
```

### Adding Torrents

```go
// From magnet link
t, err := client.AddMagnet(magnetLink)

// From .torrent file
t, err := client.AddTorrentFromFile("path/to/file.torrent")
```

### Download Control

```go
// Download all files
t.DownloadAll()

// Download specific file
file.Download()

// Pause download
t.StopDataDownload()
```

### Progress Monitoring

```go
stats := t.Stats()
progress := float64(t.BytesCompleted()) / float64(t.Length()) * 100
```

### File Management

```go
files := t.Files()
for _, file := range files {
    file.SetPriority(torrent.PiecePriorityNormal) // Download
    file.SetPriority(torrent.PiecePriorityNone)   // Skip
}
```

## Error Handling

The examples include proper error handling for:

- Client creation failures
- Invalid magnet links
- Network connectivity issues
- File system errors

## Security Considerations

- Always verify torrent sources before downloading
- Use legal torrents for testing
- Be aware of your local laws regarding torrent usage
- Consider using a VPN for privacy

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the `ListenPort` in the configuration
2. **Permission denied**: Ensure the download directory is writable
3. **No peers found**: Check your firewall settings and ensure port forwarding
4. **Slow download**: This is normal for torrents; speed depends on available peers

### Debug Information

Enable debug logging by setting the log level:

```go
import "log"
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## Contributing

Feel free to extend these examples with additional features:

- Web interface
- Database storage for torrent metadata
- Advanced peer management
- Bandwidth limiting
- Scheduling downloads

## License

This project is for educational purposes. Please respect copyright laws and only download legal content.
