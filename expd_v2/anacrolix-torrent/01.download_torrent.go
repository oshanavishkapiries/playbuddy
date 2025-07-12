package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

func DownloadTorrent() {
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
	magnetLink := "magnet:?xt=urn:btih:223F7484D326AD8EFD3CF1E548DED524833CB77E&dn=Avengers.Endgame.2019.1080p.BRRip.x264-MP4&tr=http%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2F47.ip-51-68-199.eu%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2780%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2710%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2730%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2920%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Ftracker.pirateparty.gr%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce"

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
			bytesRead := stats.BytesRead.Int64()
			fmt.Printf("Downloaded: %d bytes, Progress: %.2f%%\n", 
				bytesRead, float64(bytesRead)/float64(t.Length())*100)
			
			if bytesRead == t.Length() {
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