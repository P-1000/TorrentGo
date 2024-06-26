package torrent

import (
	"mytorrent/internal/peers"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (torrent *TorrentFile) getTrackerUrl(peerId [20]byte, port uint16) (string, error) {
	baseUrl, err := url.Parse(torrent.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(torrent.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(torrent.Length)},
	}
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String(), nil
}

func (torrent *TorrentFile) getPeers(peerId [20]byte, port uint16) ([]peers.Peer, error) {
	url, err := torrent.getTrackerUrl(peerId, port)
	if err != nil {
		return nil, err
	}
	c := &http.Client{Timeout: 20 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var trackerResp TrackerResponse
	err = bencode.Unmarshal(resp.Body, &trackerResp)
	if err != nil {
		return nil, err
	}
	return peers.Unmarshal([]byte(trackerResp.Peers))
}
