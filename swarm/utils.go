package swarm

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"

	"torrent/peers"
	"torrent/worker"
)

func (meta *DownloadMeta) calculateBoundsForPiece(index int) (int, int) {
	begin := index * meta.PieceSize
	end := begin + meta.PieceSize
	if end > meta.FileSize {
		end = meta.FileSize
	}
	return begin, end
}

// just to handle case of last piece, else always returns pieceSize
func (meta *DownloadMeta) calculatePieceSize(index int) int {
	begin, end := meta.calculateBoundsForPiece(index)
	return end - begin
}

func (meta *DownloadMeta) startDownloadWorker(peer peers.Peer, workQueue chan *pieceOfWork, results chan *pieceOfResult) {
	w, err := worker.New(peer, meta.PeerID, meta.InfoHash)
	if err != nil {
		log.Printf("peer %s failed to HandShake\n", peer.IP)
		return
	}
	log.Printf("peer %s connected with HandShake\n", peer.IP)

	defer w.Conn.Close()

	w.SendUnchoke()
	w.SendInterested()

	for piece := range workQueue {
		// if this peer doesn't have this piece, put this piece back in workQueue to retry with another
		if !w.Bitfield.HasPiece(piece.index) {
			workQueue <- piece
			continue
		}

		buf, err := attemptDownload(w, piece)
		// possible network failure, if so, close connection
		if err != nil {
			log.Printf("Downloading piece #%d from %s failed due to %s, exiting\n", piece.index, peer.IP, err)
		}

		err = checkIntegrity(piece, buf)
		if err != nil {
			log.Printf("Piece #%d from %s failed integrity check, will retry\n", piece.index, peer.IP)
			workQueue <- piece
			continue
		}

		w.SendHave(piece.index)
		results <- &pieceOfResult{piece.index, buf}
	}
}

func checkIntegrity(piece *pieceOfWork, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], piece.hash[:]) {
		return fmt.Errorf("Index %d failed integrity check", piece.index)
	}
	return nil
}
