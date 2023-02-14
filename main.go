package main

import (
	"log"
	"os"

	"torrent/torrent"
)

func main() {
	// path of torrent file
	if len(os.Args) < 2 {
		log.Fatalln("Please pass the paths of the torrent file and location to store downloaded file as command line arguments")

	}
	inPath := os.Args[1]
	// path to save file
	outPath := os.Args[2]

	tor, err := torrent.Deserialize(inPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = tor.DownloadToFile(outPath)
	if err != nil {
		log.Println("Download Failed")
		log.Fatalln(err)
	}
}
