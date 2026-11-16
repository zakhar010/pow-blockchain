package blockchain

import (
	"errors"
	"fmt"
	"slices"
)

type Chain struct {
	blocks []Block
}

func NewChain() *Chain {
	genesis := mineBlock("genesis Block", Hash{})

	return &Chain{
		blocks: []Block{genesis},
	}
}

func (c *Chain) Blocks() []Block {
	return slices.Clone(c.blocks)
}

func (c *Chain) AddBlock(data string) error {
	if (c.blocks == nil) || len(c.blocks) == 0 {
		return errors.New("chain is empty")
	}

	previousHash := c.blocks[len(c.blocks)-1].Hash()
	newBlock := mineBlock(data, previousHash)

	c.blocks = append(c.blocks, newBlock)
	return nil
}

func (c *Chain) Len() int {
	return len(c.blocks)
}

func (c *Chain) Validate() error {
	if len(c.blocks) == 0 {
		return errors.New("chain is empty")
	}

	for i, block := range c.blocks {
		if block.Header.Version != blockVersion {
			return fmt.Errorf("block №%d: unsupported version %d", i, block.Header.Version)
		}

		if block.Header.DataHash != hashBytes([]byte(block.Data)) {
			return fmt.Errorf("block №%d: invalid data hash", i)
		}

		if !meetsTarget(block.Hash(), initialTarget) {
			return fmt.Errorf("block №%d: invalid proof of work", i)
		}

		if i == 0 && block.Header.PreviousHash != (Hash{}) {
			return fmt.Errorf("block №%d: invalid genesis previous hash", i)
		}

		if i > 0 && block.Header.PreviousHash != c.blocks[i-1].Hash() {
			return fmt.Errorf("block №%d: invalid previous hash", i)
		}
	}

	return nil
}
