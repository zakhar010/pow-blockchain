package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const (
	blockVersion = 1
	headerSize = 4 + 32 + 32 + 8 + 8 // version + previous hash + hash + nonce + timestamp
)

type BlockHeader struct {
	Version      uint32
	PreviousHash Hash
	DataHash     Hash
	Nonce        uint64
	Timestamp    int64
}

func (b BlockHeader) String() string {
	return fmt.Sprintf("version: %d\nprevious hash: %x\ndata hash: %x\nnonce: %d\ntimestamp: %d",
		b.Version,
		b.PreviousHash,
		b.DataHash,
		b.Nonce,
		b.Timestamp)
}

func (b BlockHeader) Hash() Hash {
	return sha256.Sum256(b.serialize())
}

func (b BlockHeader) serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, b.Version)
	buf.Write(b.PreviousHash[:])
	buf.Write(b.DataHash[:])
	binary.Write(buf, binary.BigEndian, b.Nonce)
	binary.Write(buf, binary.BigEndian, b.Timestamp)

	return buf.Bytes()
}
