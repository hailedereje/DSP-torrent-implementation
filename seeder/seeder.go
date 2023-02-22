package main

import (
	"Bittorrent-client/bitfield"
	"Bittorrent-client/handshake"
	"Bittorrent-client/message"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const (
	ServerHost = "127.0.0.1"
	ServerPort = "8080"
	ServerType = "tcp"
)

type Torrent struct {
	Bitfield    []byte  `json:"bitfield"`
	Path        string  `json:"path"`
	PieceLength float64 `json:"piecelength"`
	Length      float64 `json:"length"`
}

func main() {
	listener, err := net.Listen(ServerType, ":"+ServerPort)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("listening on %s:%s", ServerHost, ServerPort))
	// close listener
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	torrent, err := completeHandShake(conn)
	if err != nil {
		return
	}

	sendBitfeild(torrent, conn)
	sendUnchoke(conn)

	pieceLength := torrent.PieceLength

	f, err := os.Open(torrent.Path)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	for {
		req, err := message.Read(conn)
		if err != nil {
			return
		}

		go serveRequest(req, pieceLength, f, conn)
	}
}

func serveRequest(req *message.Message, pieceLength float64, file *os.File, conn net.Conn) {
	if req.ID == message.MsgRequest {

		index := binary.BigEndian.Uint32(req.Payload[0:4])
		begin := binary.BigEndian.Uint32(req.Payload[4:8])
		length := binary.BigEndian.Uint32(req.Payload[8:])

		content := make([]byte, length)

		_, err := file.ReadAt(content, int64(int64(index)*int64(pieceLength)+int64(begin)))
		if err != nil {
			log.Fatal("Error Reading File.")
		}
		piece := getPiece(content, index, begin, length)
		conn.Write(piece.Serialize())
	}
}

func getPiece(content []byte, index uint32, begin uint32, length uint32) *message.Message {
	buf := make([]byte, 8+length)
	binary.BigEndian.PutUint32(buf[0:4], uint32(index))
	binary.BigEndian.PutUint32(buf[4:8], uint32(begin))
	copy(buf[8:], content)
	msg := &message.Message{
		ID:      message.MsgPiece,
		Payload: buf,
	}

	return msg
}

func completeHandShake(conn net.Conn) (*Torrent, error) {
	res, err := handshake.Read(conn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	infoHash := hex.EncodeToString(res.InfoHash[:])

	// Open our torrent jsonFile
	jsonFile, err := os.Open(fmt.Sprintf("files/%s.json", infoHash))
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	torrent := Torrent{}

	byteValue, err := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &torrent)

	conn.Write(res.Serialize())

	return &torrent, err
}

