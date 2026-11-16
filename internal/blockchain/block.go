package blockchain

import (
	"fmt"
)

type Block struct {
	Header BlockHeader
	Data   string
}

func (b Block) Hash() [32]byte {
	return b.Header.Hash()
}

func (b Block) String() string {
	return fmt.Sprintf("data: %s\nhash: %x\nblock header:\n%s",
		b.Data,
		b.Hash(),
		b.Header)
}
