package message

import (
	"encoding/binary"
	"io"
)

type messageID uint8

// defines what the message represents
const (
	// MsgChoke chokes the receiver
	Choke messageID = 0
	// MsgUnchoke unchokes the receiver
	Unchoke messageID = 1
	// MsgInterested expresses interest in receiving data
	Interested messageID = 2
	// MsgNotInterested expresses disinterest in receiving data
	NotInterested messageID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	Have messageID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	Bitfield messageID = 5
	// MsgRequest requests a block of data from the receiver
	Request messageID = 6
	// MsgPiece delivers a block of data to fulfill a request
	Piece messageID = 7
	// MsgCancel cancels a request
	Cancel messageID = 8
)

// Message stores ID and payload of a message
type Message struct {
	ID      messageID
	Payload []byte
}

// Read parses a message from a stream. Returns `nil` on keep-alive message
func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	// the first 4 bytes in the standard message contains the length of ID + payload in big endian format
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep-alive message
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(messageBuf[0]), // length of ID is 1 byte
		Payload: messageBuf[1:],
	}

	return &m, nil
}

// Serialize serializes a message into a buffer of the form
// <length prefix><message ID><payload>
// Interprets `nil` as a keep-alive message
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}
