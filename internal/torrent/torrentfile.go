package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

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

func (i *bencodeInfo) Hash() ([20]byte, error) {
	// conver the bencodeInfo struct to a map ...
	infoMap := map[string]interface{}{
		"pieces":       i.Pieces,
		"piece length": i.PieceLength,
		"length":       i.Length,
		"name":         i.Name,
	}

	var buffer bytes.Buffer
	err := bencode.Marshal(&buffer, infoMap)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buffer.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func OpenTorrentFile(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()
	var bto bencodeTorrent
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	infoHash, err := bto.Info.Hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	torrentfile := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}

	return torrentfile, nil
}

func (t *TorrentFile) PrintDetails() {
	fmt.Println("Torrent Name:", t.Name)
	fmt.Println("Announce URL:", t.Announce)
	fmt.Printf("InfoHash: %x\n", t.InfoHash)
	fmt.Println("Piece Length:", t.PieceLength)
	fmt.Println("Total Length:", t.Length)
	fmt.Println("Number of Pieces:", len(t.PieceHashes))
}
