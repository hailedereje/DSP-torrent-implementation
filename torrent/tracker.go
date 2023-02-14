package torrent

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"torrent/bencodeUtils"
	"torrent/peers"
)

// Port to listen on
const Port uint16 = 6881 // a standard port for bittorrent

func (t *TorrentFile) getTrackerReqURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.TrackerBaseURL)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

// tell the tracker that I am a new peer and get list of existing peers
func (t *TorrentFile) requestForPeers(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	url, err := t.getTrackerReqURL(peerID, port)
	if err != nil {
		return nil, err
	}
	log.Println("Tracker URL is", url)
	c := &http.Client{Timeout: 15 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close() // is this required?
	parsedRes, err := bencodeUtils.ParseTrackerResp(res.Body)
	if err != nil {
		return nil, err
	}
	return peers.Deserialize([]byte(parsedRes.Peers))
}
