package swarm

import (
	"torrent/message"
	"torrent/worker"
)

type progressTracker struct {
	index      int
	worker     *worker.Worker
	buf        []byte
	downloaded int
	requested  int
	backlog    int // keeps track of number of blocks left to process
}

func (state *progressTracker) readMessage() error {
	msg, err := state.worker.Read() // blocking call
	if err != nil {
		return err
	}
	if msg == nil { // keep-alive
		return nil
	}

	switch msg.ID {
	case message.Unchoke:
		state.worker.Choked = false
	case message.Choke:
		state.worker.Choked = true
	case message.Have:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.worker.Bitfield.SetPiece(index)
	case message.Piece:
		n, err := message.ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}
