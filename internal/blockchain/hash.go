package blockchain

import (
	"crypto/sha256"
)

type Hash [sha256.Size]byte

func hashBytes(data []byte) Hash {
	return sha256.Sum256(data)
}
