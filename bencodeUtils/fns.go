package bencodeUtils

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"torrent/bencode-go"
)

func ParseTorrent(r io.Reader) (*bencodeTorrent, error) {
	tor := bencodeTorrent{}
	err := bencode.Unmarshal(r, &tor)
	if err != nil {
		return nil, err
	}
	return &tor, nil
}

func ParseTrackerResp(r io.Reader) (*trackerResponse, error) {
	res := trackerResponse{}
	err := bencode.Unmarshal(r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (i *bencodeInfo) Hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) SplitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}
