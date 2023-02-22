package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jackpal/bencode-go"
	"example.com/torrent/bitfield"
	"example.com/torrent/peer"
	"example.com/torrent/peer2peer"
)

const DefaultPort uint16 = 6881

// TorrentMetaInfo encodes the metadata from a .torrent file
type TorrentMetaInfo struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type bencodeInfo struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
}

func OpenTorrentFile(path string) (TorrentMetaInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentMetaInfo{}, err
	}
	defer file.Close()

	bct := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bct)
	if err != nil {
		return TorrentMetaInfo{}, err
	}

	return bct.toTorrentMetaInfo()
}

func (bct *bencodeTorrent) toTorrentMetaInfo() (TorrentMetaInfo, error) {
	infoHash, pieceHashes, err := bct.Info.hashInfo()
	if err != nil {
		return TorrentMetaInfo{}, err
	}
	return TorrentMetaInfo{
		Announce:    bct.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bct.Info.PieceLength,
		Length:      bct.Info.Length,
		Name:        bct.Info.Name,
	}, nil
}

func (bci *bencodeInfo) hashInfo() ([20]byte, [][20]byte, error) {
	pieces := []byte(bci.Pieces)
	hashLen := 20
	numHashes := len(pieces) / hashLen
	pieceHashes := make([][20]byte, numHashes)

	if len(pieces)%hashLen != 0 {
		err := fmt.Errorf("reading hash info failed: invalid hash length (length: %d - expected: %d", len(pieces), hashLen)
		return [20]byte{}, [][20]byte{}, err
	}
	for i := range pieceHashes {
		copy(pieceHashes[i][:], pieces[i*hashLen:(i+1)*hashLen])
	}

	var info bytes.Buffer
	err := bencode.Marshal(&info, *bci)
	if err != nil {
		return [20]byte{}, [][20]byte{}, err
	}
	infoHash := sha1.Sum(info.Bytes())
	return infoHash, pieceHashes, nil
}

func (tmi *TorrentMetaInfo) Download(path string) error {
	peerID, err := peer.GeneratePeerID()
	if err != nil {
		return err
	}

	seed, _ := peer.UnmarshalPeers()

	log.Println(seed)

	pieceLength := len(tmi.PieceHashes)
	byteSize := 8
	bitfield := make(bitfield.Bitfield, int(math.Ceil(float64(pieceLength)/float64(byteSize))))

	outFile, _ := os.OpenFile(tmi.Name, os.O_RDWR|os.O_CREATE, 0666)

	torrent

}
