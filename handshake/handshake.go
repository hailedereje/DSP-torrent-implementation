package handshake

import (
	"fmt"
	"io"
)

type Handshake struct {
	ProtocolName string
	InfoHash     [20]byte
	PeerID       [20]byte
}

func NewHandshake(infoHash [20]byte) *Handshake {
	return &Handshake{
		ProtocolName: "BitTorrent protocol",
		InfoHash:     infoHash,
	}
}

func ReadHandshake(r io.Reader) (*Handshake, error) {
	bufferLength := make([]byte, 1)
	_, err := io.ReadFull(r, bufferLength)
	if err != nil {
		return nil, err
	}
	protocolNameLen := int(bufferLength[0])

	if protocolNameLen == 0 {
		err := fmt.Errorf("Protocol Names length cannot be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, 48+protocolNameLen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[protocolNameLen+8:protocolNameLen+8+20])
	copy(peerID[:], handshakeBuf[protocolNameLen+8+20:])

	handshake := Handshake{
		ProtocolName: string(handshakeBuf[0:protocolNameLen]),
		InfoHash:     infoHash,
		PeerID:       peerID,
	}

	return &handshake, nil
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.ProtocolName)+49)

	buf[0] = byte(len(h.ProtocolName))
	curr := 1
	curr += copy(buf[curr:], h.ProtocolName)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	return buf
}
