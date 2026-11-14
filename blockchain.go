package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Hash [32]byte
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

type Block struct {
	Data         string
	Hash         Hash
	PreviousHash Hash
	Nonce        uint64
}

func (b Block) CalculateHash() [32]byte {
	input := fmt.Sprintf(
		"%x:%s:%d",
		b.PreviousHash,
		b.Data,
		b.Nonce,
	)

	hash := sha256.Sum256([]byte(input))

	return hash
}

func (b Block) String() string {
	return fmt.Sprintf("data: %s\nhash: %x\nprevious hash: %x\nnonce: %d", b.Data, b.Hash, b.PreviousHash, b.Nonce)
}

type Chain struct {
	Blocks []Block
}

func NewChain() Chain {
	genesis := MineBlock("genesis Block", Hash{})

	return Chain{
		Blocks: []Block{genesis},
	}
}

func (c *Chain) AddBlock(data string) {
	previousHash := c.Blocks[len(c.Blocks)-1].Hash
	newBlock := MineBlock(data, previousHash)

	c.Blocks = append(c.Blocks, newBlock)
}

func (c *Chain) Validate() error {
	if len(c.Blocks) == 0 {
		return fmt.Errorf("chain is empty")
	}

	for i, block := range c.Blocks {
		expectedHash := block.CalculateHash()

		if expectedHash != block.Hash {
			return fmt.Errorf("block №%d: invalid hash", i)
		}

		if !MeetsTarget(block.Hash, initialTarget) {
			return fmt.Errorf("block №%d: invalid proof of work", i)
		}

		if i == 0 && block.PreviousHash != (Hash{}) {
			return fmt.Errorf("block №%d: invalid genesis previous hash", i)
		}

		if i > 0 && block.PreviousHash != c.Blocks[i-1].Hash {
			return fmt.Errorf("block №%d: invalid previous hash", i)
		}
	}

	return nil
}

func MeetsTarget(hash Hash, target Target) bool {
	return bytes.Compare(hash[:], target[:]) <= 0
}

func MineBlock(data string, previousHash [32]byte) Block {
	block := Block{
		Data:         data,
		PreviousHash: previousHash,
	}

	for {
		hash := block.CalculateHash()

		if MeetsTarget(hash, initialTarget) {
			block.Hash = hash
			return block
		}

		block.Nonce++
	}
}
