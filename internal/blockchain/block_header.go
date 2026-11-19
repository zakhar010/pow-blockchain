package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const (
	blockVersion = 1
	headerSize   = 4 + 32 + 32 + 8 + 8 + 32 // version + previous hash + hash + nonce + timestamp + target
)

type BlockHeader struct {
	Version      uint32
	PreviousHash Hash
	DataHash     Hash
	Timestamp    int64
	Target       Target
	Nonce        uint64
}

func (b BlockHeader) String() string {
	return fmt.Sprintf("version: %d\nprevious hash: %x\ndata hash: %x\nnonce: %d\ntimestamp: %d\ntarget: %x",
		b.Version,
		b.PreviousHash,
		b.DataHash,
		b.Nonce,
		b.Timestamp,
		b.Target)
}

func (b BlockHeader) Hash() Hash {
	return sha256.Sum256(b.serialize())
}

func (b BlockHeader) serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, b.Version)
	buf.Write(b.PreviousHash[:])
	buf.Write(b.DataHash[:])
	binary.Write(buf, binary.BigEndian, b.Timestamp)
	buf.Write(b.Target[:])
	binary.Write(buf, binary.BigEndian, b.Nonce)

	return buf.Bytes()
}
