package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anacrolix/torrent"
)

func getTorrentInfo() {
	// Create a new torrent client configuration
	config := torrent.NewDefaultClientConfig()
	config.DataDir = "./info_cache"

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

	// Get torrent information
	info := t.Info()
	if info == nil {
		log.Fatal("Failed to get torrent info")
	}

	fmt.Println("=== Torrent Information ===")
	fmt.Printf("Name: %s\n", t.Name())
	fmt.Printf("Total Size: %d bytes (%.2f MB)\n", t.Length(), float64(t.Length())/1024/1024)
	fmt.Printf("Piece Length: %d bytes\n", info.PieceLength)
	fmt.Printf("Number of Pieces: %d\n", info.NumPieces)
	fmt.Printf("Number of Files: %d\n", len(t.Files()))

	// Display file information
	fmt.Println("\n=== File Information ===")
	for i, file := range t.Files() {
		fmt.Printf("File %d: %s\n", i+1, file.DisplayPath())
		fmt.Printf("  Size: %d bytes (%.2f MB)\n", file.Length(), float64(file.Length())/1024/1024)
		fmt.Printf("  Priority: %d\n", file.Priority())
	}

	// Get peer information
	fmt.Println("\n=== Peer Information ===")
	peers := t.PeerConns()
	fmt.Printf("Connected Peers: %d\n", len(peers))
	
	for i, peer := range peers {
		fmt.Printf("Peer %d:\n", i+1)
		fmt.Printf("  Address: %s\n", peer.RemoteAddr())
		fmt.Printf("  Client: %s\n", peer.ClientName())
		fmt.Printf("  Upload Speed: %d bytes/s\n", peer.Stats().BytesWrittenData.Int64())
		fmt.Printf("  Download Speed: %d bytes/s\n", peer.Stats().BytesReadData.Int64())
	}

	// Get torrent statistics
	fmt.Println("\n=== Torrent Statistics ===")
	stats := t.Stats()
	fmt.Printf("Bytes Read: %d (%.2f MB)\n", stats.BytesRead, float64(stats.BytesRead)/1024/1024)
	fmt.Printf("Bytes Written: %d (%.2f MB)\n", stats.BytesWritten, float64(stats.BytesWritten)/1024/1024)
	fmt.Printf("Bytes ReadData: %d (%.2f MB)\n", stats.BytesReadData.Int64(), float64(stats.BytesReadData.Int64())/1024/1024)
	fmt.Printf("Bytes WrittenData: %d (%.2f MB)\n", stats.BytesWrittenData.Int64(), float64(stats.BytesWrittenData.Int64())/1024/1024)
	fmt.Printf("ChunksRead: %d\n", stats.ChunksRead)
	fmt.Printf("ChunksWritten: %d\n", stats.ChunksWritten)

	// Get piece information
	fmt.Println("\n=== Piece Information ===")
	pieces := t.Pieces()
	fmt.Printf("Total Pieces: %d\n", len(pieces))
	fmt.Printf("Completed Pieces: %d\n", t.BytesCompleted()/info.PieceLength)
	fmt.Printf("Completion: %.2f%%\n", float64(t.BytesCompleted())/float64(t.Length())*100)

	// Monitor real-time statistics
	fmt.Println("\n=== Real-time Monitoring ===")
	go func() {
		for {
			time.Sleep(2 * time.Second)
			currentStats := t.Stats()
			completion := float64(t.BytesCompleted()) / float64(t.Length()) * 100
			fmt.Printf("Progress: %.2f%%, Downloaded: %.2f MB, Uploaded: %.2f MB, Peers: %d\n",
				completion,
				float64(currentStats.BytesReadData.Int64())/1024/1024,
				float64(currentStats.BytesWrittenData.Int64())/1024/1024,
				len(t.PeerConns()))
		}
	}()

	// Keep monitoring for 30 seconds
	time.Sleep(30 * time.Second)
} 