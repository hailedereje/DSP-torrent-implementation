package client

import (
	"Bittorrent-client/bitfield"
	"Bittorrent-client/handshake"
	"Bittorrent-client/message"
	"Bittorrent-client/peers"
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	conn     net.Conn
	choked   bool
	bitfield bitfield.Bitfield
	peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func (c *Client) connect() error {
	conn, err := net.DialTimeout("tcp", c.peer.String(), 10*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) completeHandshake() error {
	c.conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.conn.SetDeadline(time.Time{})

	req := handshake.New(c.infoHash, c.peerID)
	_, err := c.conn.Write(req.Serialize())
	if err != nil {
		return err
	}

	res, err := handshake.Read(c.conn)
	if err != nil {
		return err
	}
	if !bytes.Equal(res.InfoHash[:], c.infoHash[:]) {
		return fmt.Errorf("File Mismatch: expected infohash %x but got %x", res.InfoHash, c.infoHash)
	}
	return nil
}

func (c *Client) recvBitfield() error {
	c.conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer c.conn.SetDeadline(time.Time{})

	msg, err := message.Read(c.conn)
	if err != nil {
		return err
	}
	if msg == nil {
		err := fmt.Errorf("Expected bitfield but got %s", msg)
		return err
	}
	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("Expected bitfield but got ID %d", msg.ID)
		return err
	}

	c.bitfield = msg.Payload
	return nil
}

func NewClient(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	c := &Client{
		peer:     peer,
		peerID:   peerID,
		infoHash: infoHash,
	}
	if err := c.connect(); err != nil {
		return nil, err
	}
	if err := c.completeHandshake(); err != nil {
		return nil, err
	}
	if err := c.recvBitfield(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.conn)
	return msg, err
}

func (c *Client) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := c.conn.Write(req.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	
