package peers

import (
    "crypto/rand"
    "fmt"
    "net"
)

type PeerID [20]byte

var prefix = []byte("-P00001-")

// Node contains the details of a peer
type Node struct {
    IP   net.IP
    Port uint16
}

// Unmarshal returns the seeders as an array of nodes(ip:port)
func Unmarshal() ([]Node, error) {
    nodes := make([]Node, 1)
    node := make([]byte, 4)
    node[0] = 192
    node[1] = 168
    node[2] = 228
    node[3] = 142

    nodes[0].IP = net.IP(node)
    nodes[0].Port = 5858
    return nodes, nil
}

// String returns a string representation of the Node object
func (n Node) String() string {
    return net.JoinHostPort(n.IP.String(), fmt.Sprint(n.Port))
}

// GenerateNodeID generates a new node ID
func GenerateNodeID() (PeerID, error) {
    var id PeerID
    copy(id[:], prefix)
    _, err := rand.Read(id[len(prefix):])
    if err != nil {
        return PeerID{}, err
    }
    return id, nil
}
