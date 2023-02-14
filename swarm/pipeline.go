package swarm

import (
	"time"

	"torrent/worker"
)

// MaxBlockSize is the largest number of bytes a request can ask for
const MaxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const MaxBacklog = 5

// download a piece as many blocks in a pipelined fashion for efficiency
func attemptDownload(w *worker.Worker, piece *pieceOfWork) ([]byte, error) {
	state := progressTracker{
		index:  piece.index,
		worker: w,
		buf:    make([]byte, piece.length),
	}

	// Setting a deadline helps get unresponsive peers unstuck.
	// 30 seconds is more than enough time to download a 262 KB piece
	w.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer w.Conn.SetDeadline(time.Time{}) // to disable the deadline

	// state.downloaded defaults to zero since it's not provided
	for state.downloaded < piece.length {
		// If unchoked, send requests until we have enough unfulfilled requests
		if !state.worker.Choked {
			for state.backlog < MaxBacklog && state.requested < piece.length {
				blockSize := MaxBlockSize
				// Last block might be shorter than the typical block
				if piece.length-state.requested < blockSize {
					blockSize = piece.length - state.requested
				}

				err := w.SendRequest(piece.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}

		// increments state.downloaded appropriately
		err := state.readMessage()
		if err != nil {
			return nil, err
		}
	}

	return state.buf, nil
}
