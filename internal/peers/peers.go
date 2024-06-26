package peers

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Peer struct {
	IP   net.IP
	port uint16
}

func Unmarshal(peerBinary []byte) ([]Peer, error) {
	const peerSize = 6
	numPeers := len(peerBinary) / peerSize
	if len(peerBinary)%peerSize != 0 {
		err := fmt.Errorf("invalid Peers received")
		return nil, err
	}
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peerBinary[offset : offset+4])
		peers[i].port = binary.BigEndian.Uint16(peerBinary[offset+4 : offset+6])
	}
	return peers, nil
}
