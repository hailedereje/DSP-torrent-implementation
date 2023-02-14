package swarm

import (
	"log"
	"runtime"

	"torrent/peers"
)

type DownloadMeta struct {
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceSize   int
	FileSize    int
	Name        string
}

type pieceOfWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceOfResult struct {
	index int
	buf   []byte
}

func (meta *DownloadMeta) Download() ([]byte, error) {
	log.Println("Downloading", meta.Name)

	// stores info on remaining pieces to download
	workQueue := make(chan *pieceOfWork, len(meta.PieceHashes))
	// stores downloaded data
	results := make(chan *pieceOfResult)

	for index, hash := range meta.PieceHashes {
		length := meta.calculatePieceSize(index)
		workQueue <- &pieceOfWork{index, hash, length}
	}

	log.Println("Starting download workers")

	// start goroutine workers, one for each peer
	// orchestration is simplified by using common channels to communicate
	for _, peer := range meta.Peers {
		go meta.startDownloadWorker(peer, workQueue, results)
	}

	// store result in memory -> maybe better to save peices to disk for big files
	resultBuf := make([]byte, meta.FileSize)
	donePieces := 0
	for donePieces < len(meta.PieceHashes) {
		piece := <-results
		begin, end := meta.calculateBoundsForPiece(piece.index)
		copy(resultBuf[begin:end], piece.buf)

		percentComplete := float64(donePieces) / float64(len(meta.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		log.Printf("(%0.2f%% done) Downloaded piece #%d, %d peers online\n", percentComplete, piece.index, numWorkers)
		donePieces++
	}
	close(workQueue)
	return resultBuf, nil
}
