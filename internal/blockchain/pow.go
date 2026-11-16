package blockchain

import (
	"bytes"
	"time"
)

type Target [32]byte

var initialTarget = Target{
	0x00, 0x00,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff,
	0xff, 0xff,
}

func meetsTarget(hash Hash, target Target) bool {
	return bytes.Compare(hash[:], target[:]) <= 0
}

func mineBlock(data string, previousHash Hash) Block {
	block := Block{
		Data: data,
		Header: BlockHeader{
			Version:      blockVersion,
			PreviousHash: previousHash,
			DataHash:     hashBytes([]byte(data)),
			Timestamp:    time.Now().Unix(),
		},
	}

	for {
		hash := block.Header.Hash()

		if meetsTarget(hash, initialTarget) {
			return block
		}

		block.Header.Nonce++
	}
}
