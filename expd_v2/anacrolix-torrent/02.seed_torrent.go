package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

func seedTorrent() {
	// Create a new torrent client configuration
	config := torrent.NewDefaultClientConfig()
	config.DataDir = "./uploads" // Set upload directory
	config.ListenPort = 6882     // Set listen port for incoming connections

	// Create the torrent client
	client, err := torrent.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating torrent client: %v", err)
	}
	defer client.Close()

	// Path to the directory or file you want to seed
	pathToSeed := "./files_to_seed"

	// Check if the path exists
	if _, err := os.Stat(pathToSeed); os.IsNotExist(err) {
		log.Fatalf("Path %s does not exist", pathToSeed)
	}

	// Create a torrent file for the directory
	torrentFilePath := "./my_torrent.torrent"
	err = createTorrentFile(pathToSeed, torrentFilePath)
	if err != nil {
		log.Fatalf("Error creating torrent file: %v", err)
	}

	// Add the torrent to the client for seeding
	t, err := client.AddTorrentFromFile(torrentFilePath)
	if err != nil {
		log.Fatalf("Error adding torrent file: %v", err)
	}

	fmt.Printf("Seeding torrent: %s\n", t.Name())
	fmt.Printf("Total size: %d bytes\n", t.Length())
	fmt.Printf("Number of files: %d\n", len(t.Files()))

	// Start seeding (uploading)
	t.DownloadAll() // This will start seeding if files are already present

	// Monitor seeding progress
	go func() {
		for {
			time.Sleep(5 * time.Second)
			stats := t.Stats()
			fmt.Printf("Uploaded: %d bytes, Downloaded: %d bytes, Peers: %d\n", 
				stats.BytesWritten, stats.BytesRead, len(t.PeerConns()))
		}
	}()

	// Keep seeding indefinitely
	fmt.Println("Seeding started. Press Ctrl+C to stop...")
	select {}
}

// createTorrentFile creates a .torrent file for the given path
func createTorrentFile(path, torrentFilePath string) error {
	// This is a simplified example. In a real implementation,
	// you would use the torrent library's functionality to create
	// a proper .torrent file with correct piece size, trackers, etc.
	
	fmt.Printf("Creating torrent file for: %s\n", path)
	fmt.Printf("Torrent file will be saved as: %s\n", torrentFilePath)
	
	// Note: The actual torrent file creation would require more complex logic
	// including piece hashing, tracker information, etc.
	// This is just a placeholder for demonstration purposes.
	
	return nil
} 