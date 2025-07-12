package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

func main() {
	// Create a new torrent client configuration
	config := torrent.NewDefaultClientConfig()
	config.DataDir = "./downloads" // Set download directory
	config.ListenPort = 6881       // Set listen port for incoming connections

	// Create the torrent client
	client, err := torrent.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating torrent client: %v", err)
	}
	defer client.Close()

	// Example magnet link (replace with your own)
	magnetLink := "magnet:?xt=urn:btih:YOUR_TORRENT_HASH_HERE&dn=example+torrent"

	// Add the torrent to the client
	t, err := client.AddMagnet(magnetLink)
	if err != nil {
		log.Fatalf("Error adding magnet link: %v", err)
	}

	// Wait for torrent metadata to be downloaded
	<-t.GotInfo()

	fmt.Printf("Torrent name: %s\n", t.Name())
	fmt.Printf("Total size: %d bytes\n", t.Length())
	fmt.Printf("Number of files: %d\n", len(t.Files()))

	// Download all files in the torrent
	t.DownloadAll()

	// Monitor download progress
	go func() {
		for {
			time.Sleep(1 * time.Second)
			stats := t.Stats()
			fmt.Printf("Downloaded: %d bytes, Progress: %.2f%%\n", 
				stats.BytesRead, float64(stats.BytesRead)/float64(t.Length())*100)
			
			if stats.BytesRead == t.Length() {
				fmt.Println("Download completed!")
				break
			}
		}
	}()

	// Wait for download to complete or context to be cancelled
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Println("Download timeout or cancelled")
	case <-t.Complete.On():
		fmt.Println("Download completed successfully!")
	}
}

// Helper function to create downloads directory
func ensureDownloadDir(dir string) error {
	return os.MkdirAll(dir, 0755)
} 