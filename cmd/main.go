package main

import (
	"fmt"
	"os"

	"mytorrent/internal/torrent"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-torrent-file>")
		return
	}

	torrentPath := os.Args[1]
	torrentFile, err := torrent.OpenTorrentFile(torrentPath)
	if err != nil {
		fmt.Println("Error opening torrent file:", err)
		return
	}

	torrentFile.PrintDetails()
}
