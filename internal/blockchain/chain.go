package blockchain

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

const (
	MaxFutureSeconds = 120
)

type Chain struct {
	blocks []Block
}

func NewChain() (*Chain, error) {
	genesis := genesisBlock()

	if genesis.Hash() != genesisHash {
		return nil, errors.New("invalid genesis block hash")
	}

	if !meetsTarget(genesis.Hash(), genesis.Header.Target) {
		return nil, errors.New("genesis block does not meet target")
	}

	return &Chain{
		blocks: []Block{genesis},
	}, nil
}

func (c *Chain) Blocks() []Block {
	return slices.Clone(c.blocks)
}

func (c *Chain) MineNextBlock(data string, timestamp int64) (Block, error) {
	if (c.blocks == nil) || len(c.blocks) == 0 {
		return Block{}, errors.New("chain is empty")
	}

	previousBlock := c.blocks[len(c.blocks)-1]

	if err := validateTimestamp(
		previousBlock.Header.Timestamp,
		timestamp,
	); err != nil {
		return Block{}, err
	}

	if err := validateFutureTimestamp(
		timestamp,
		time.Now().Unix(),
	); err != nil {
		return Block{}, err
	}

	block, err := mineBlock(
		data,
		previousBlock.Hash(),
		c.RequiredTarget(),
		timestamp,
	)

	if err != nil {
		return Block{}, fmt.Errorf("failed to mine block: %v", err)
	}

	return block, nil
}

func (c *Chain) AddBlock(block Block) error {
	if len(c.blocks) == 0 {
		return errors.New("chain is empty")
	}

	if err := c.validateNextBlock(block); err != nil {
		return err
	}

	c.blocks = append(c.blocks, block)
	return nil
}

func (c *Chain) Len() int {
	return len(c.blocks)
}

func (c *Chain) Validate() error {
	if len(c.blocks) == 0 {
		return errors.New("chain is empty")
	}

	if c.blocks[0] != genesisBlock() {
		return errors.New("invalid genesis block")
	}

	if c.blocks[0].Hash() != genesisHash {
		return errors.New("invalid genesis block hash")
	}

	for i, block := range c.blocks {
		if i > 0 {
			if err := validateTimestamp(
				c.blocks[i-1].Header.Timestamp,
				block.Header.Timestamp,
			); err != nil {
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

	lastBlock := c.blocks[nextHeight-1]

	if ((nextHeight - 1) % RetargetInterval) != 0 {
		return lastBlock.Header.Target
	}

	lastIndex := nextHeight - 1
	firstIndex := lastIndex - RetargetInterval

	actualTimestamp := c.blocks[lastIndex].Header.Timestamp - c.blocks[firstIndex].Header.Timestamp

	return calculateNextTarget(lastBlock.Header.Target, actualTimestamp)
}

func (c *Chain) validateNextBlock(block Block) error {
	previous := c.blocks[len(c.blocks)-1]

	if block.Header.Version != blockVersion {
		return fmt.Errorf("block №%d: unsupported version %d", len(c.blocks), block.Header.Version)
	}

	if block.Header.PreviousHash != previous.Hash() {
		return fmt.Errorf("block №%d: invalid previous hash", len(c.blocks))
	}

	if block.Header.DataHash != hashBytes([]byte(block.Data)) {
		return fmt.Errorf("block №%d: invalid data hash", len(c.blocks))
	}

	if err := validateTimestamp(
		previous.Header.Timestamp,
		block.Header.Timestamp,
	); err != nil {
		return err
	}

	if err := validateFutureTimestamp(
		block.Header.Timestamp,
		time.Now().Unix(),
	); err != nil {
		return err
	}

	requiredTarget := c.RequiredTarget()

	if block.Header.Target != requiredTarget {
		return fmt.Errorf("block №%d: invalid target", len(c.blocks))
	}

	if !meetsTarget(block.Hash(), block.Header.Target) {
		return fmt.Errorf("block №%d: invalid proof of work", len(c.blocks))
	}

	return nil
}

func validateTimestamp(previous, timestamp int64) error {
	if timestamp <= previous {
		return errors.New(
			"timestamp must be greater than previous block's timestamp",
		)
	}

	return nil
}

func validateFutureTimestamp(timestamp, now int64) error {
	if timestamp > now+MaxFutureSeconds {
		return errors.New("timestamp is too far in the future")
	}

	return nil
}
