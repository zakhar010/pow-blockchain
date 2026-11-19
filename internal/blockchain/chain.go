package blockchain

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

const (
    genesisData      = "Genesis Block"
	genesisTimestamp int64 = 1767225600 // 2026-01-01 00:00:00 UTC

	MaxFutureSeconds = 120
)

type Chain struct {
	blocks []Block
}

func NewChain() (*Chain, error) {
	genesis, err := mineBlock(
        genesisData,
        Hash{},
        InitialTarget(),
        genesisTimestamp,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to mine genesis block: %w", err)
    }

    return &Chain{
		blocks: []Block{genesis},
	}, nil
}

func (c *Chain) Blocks() []Block {
	return slices.Clone(c.blocks)
}

func (c *Chain) AddBlock(data string, timestamp int64) error {
	if (c.blocks == nil) || len(c.blocks) == 0 {
		return errors.New("chain is empty")
	}

	previousBlock := c.blocks[len(c.blocks)-1]

    if err := validateTimestamp(previousBlock.Header.Timestamp, timestamp, time.Now().Unix()); err != nil {
        return err
    }

	target := c.RequiredTarget()

	previousHash := c.blocks[len(c.blocks)-1].Hash()
	newBlock, err := mineBlock(data, previousHash, target, timestamp)

	if err != nil {
		return fmt.Errorf("failed to mine block: %v", err)
	}
	
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
		if i > 0 {
			if err := validateTimestamp(c.blocks[i-1].Header.Timestamp, block.Header.Timestamp, time.Now().Unix()); err != nil {
				return fmt.Errorf("block #%d: %w", i, err)
			}
		}
		if block.Header.Version != blockVersion {
			return fmt.Errorf("block №%d: unsupported version %d", i, block.Header.Version)
		}

		if block.Header.DataHash != hashBytes([]byte(block.Data)) {
			return fmt.Errorf("block №%d: invalid data hash", i)
		}

		var requiredTarget Target

		if i == 0 {
			requiredTarget = InitialTarget()
		} else {
			previousChain := Chain{
				blocks: c.blocks[:i],
			}

			requiredTarget = previousChain.RequiredTarget()
		}

		if block.Header.Target != requiredTarget {
			return fmt.Errorf("block #%d: invalid target", i)
		}

		if !meetsTarget(block.Hash(), block.Header.Target) {
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

func (c *Chain) RequiredTarget() Target {
	nextHeight := len(c.blocks)

	if nextHeight <= RetargetInterval {
		return InitialTarget()
	}

	lastBlock := c.blocks[nextHeight - 1]

	if ((nextHeight - 1) % RetargetInterval) != 0 {
		return lastBlock.Header.Target
	}

	lastIndex := nextHeight - 1
	firstIndex := lastIndex - RetargetInterval

	actualTimestamp := c.blocks[lastIndex].Header.Timestamp - c.blocks[firstIndex].Header.Timestamp

	return calculateNextTarget(lastBlock.Header.Target, actualTimestamp)
}

func validateTimestamp(previousTimestamp, timestamp, now int64) error {
	if timestamp <= previousTimestamp {
		return errors.New("timestamp must be greater than previous block's timestamp")
	}
	if timestamp > now + MaxFutureSeconds {
		return errors.New("timestamp is too far in the future")
	}
	return nil
}