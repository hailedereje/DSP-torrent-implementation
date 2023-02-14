package bencodeUtils

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

type trackerResponse struct {
	Interval int    // how often to poll tracker and refresh peer list
	Peers    string // serialized list of peers to download from
}
