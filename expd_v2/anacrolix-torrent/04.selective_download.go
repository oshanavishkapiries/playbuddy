package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anacrolix/torrent"
)

func selectiveDownload() {
	// Create a new torrent client configuration
	config := torrent.NewDefaultClientConfig()
	config.DataDir = "./selective_downloads"

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

	fmt.Println("=== Available Files ===")
	files := t.Files()
	for i, file := range files {
		fmt.Printf("%d. %s (%.2f MB)\n", i+1, file.DisplayPath(), float64(file.Length())/1024/1024)
	}

	// Example: Download only specific files (indices 0 and 2)
	selectedFiles := []int{0, 2} // Change these indices as needed
	
	fmt.Printf("\n=== Downloading Selected Files ===")
	for _, fileIndex := range selectedFiles {
		if fileIndex >= 0 && fileIndex < len(files) {
			file := files[fileIndex]
			fmt.Printf("Downloading: %s\n", file.DisplayPath())
			
			// Set file priority to high (download)
			file.SetPriority(torrent.PiecePriorityNormal)
			
			// Download this specific file
			file.Download()
		}
	}

	// Set other files to not download (skip)
	for i, file := range files {
		isSelected := false
		for _, selectedIndex := range selectedFiles {
			if i == selectedIndex {
				isSelected = true
				break
			}
		}
		
		if !isSelected {
			file.SetPriority(torrent.PiecePriorityNone)
			fmt.Printf("Skipping: %s\n", file.DisplayPath())
		}
	}

	// Monitor download progress for selected files
	go func() {
		for {
			time.Sleep(2 * time.Second)
			
			totalSelectedSize := int64(0)
			totalDownloadedSize := int64(0)
			
			for _, fileIndex := range selectedFiles {
				if fileIndex >= 0 && fileIndex < len(files) {
					file := files[fileIndex]
					totalSelectedSize += file.Length()
					totalDownloadedSize += file.BytesCompleted()
				}
			}
			
			if totalSelectedSize > 0 {
				progress := float64(totalDownloadedSize) / float64(totalSelectedSize) * 100
				fmt.Printf("Selected files progress: %.2f%% (%.2f MB / %.2f MB)\n",
					progress,
					float64(totalDownloadedSize)/1024/1024,
					float64(totalSelectedSize)/1024/1024)
			}
		}
	}()

	// Wait for selected files to complete or timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Check if all selected files are completed
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Download timeout")
			return
		default:
			allCompleted := true
			for _, fileIndex := range selectedFiles {
				if fileIndex >= 0 && fileIndex < len(files) {
					file := files[fileIndex]
					if file.BytesCompleted() < file.Length() {
						allCompleted = false
						break
					}
				}
			}
			
			if allCompleted {
				fmt.Println("All selected files downloaded successfully!")
				return
			}
			
			time.Sleep(1 * time.Second)
		}
	}
}

// Example function to download files by name pattern
func downloadFilesByName(t *torrent.Torrent, patterns []string) {
	files := t.Files()
	
	for _, pattern := range patterns {
		for _, file := range files {
			if containsPattern(file.DisplayPath(), pattern) {
				fmt.Printf("Downloading file matching pattern '%s': %s\n", pattern, file.DisplayPath())
				file.SetPriority(torrent.PiecePriorityNormal)
				file.Download()
			}
		}
	}
}

// Helper function to check if string contains pattern
func containsPattern(s, pattern string) bool {
	// Simple pattern matching - you could use regex for more complex patterns
	return len(pattern) > 0 && len(s) >= len(pattern) && 
		   (s == pattern || 
		    (len(s) > len(pattern) && s[:len(pattern)] == pattern) ||
		    (len(s) > len(pattern) && s[len(s)-len(pattern):] == pattern))
} 