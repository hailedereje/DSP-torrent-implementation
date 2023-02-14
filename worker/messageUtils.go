package worker

import "torrent/message"

// Read reads and consumes a message from the connection
func (w *Worker) Read() (*message.Message, error) {
	msg, err := message.Read(w.Conn)
	return msg, err
}

// SendRequest sends a Request message to the peer
func (w *Worker) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := w.Conn.Write(req.Serialize())
	return err
}

// SendInterested sends an Interested message to the peer
func (w *Worker) SendInterested() error {
	msg := message.Message{ID: message.Interested}
	_, err := w.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer
func (w *Worker) SendNotInterested() error {
	msg := message.Message{ID: message.NotInterested}
	_, err := w.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
func (w *Worker) SendUnchoke() error {
	msg := message.Message{ID: message.Unchoke}
	_, err := w.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have message to the peer
func (w *Worker) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := w.Conn.Write(msg.Serialize())
	return err
}
