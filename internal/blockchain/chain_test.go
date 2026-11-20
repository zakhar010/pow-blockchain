package blockchain

import (
	"reflect"
	"testing"
)

func TestNewChain(t *testing.T) {
	chain, err := NewChain()

	if err != nil {
		t.Fatalf("chain creation error: %v", err)
	}

	if chain.Len() != 1 {
		t.Fatalf(
			"chain length = %d, want 1",
			chain.Len(),
		)
	}

	if chain.Blocks()[0] != genesisBlock() {
		t.Fatalf("first block is not the canonical genesis block")
	}

	if err := chain.Validate(); err != nil {
		t.Fatalf("new chain is invalid: %v", err)
	}
}

func TestRejectInvalidBlock(t *testing.T) {
	chain, err := NewChain()
	if err != nil {
		t.Fatalf("failed to create chain: %v", err)
	}

	block, err := chain.MineNextBlock(
		"original data",
		genesisTimestamp+TargetBlockTime,
	)
	if err != nil {
		t.Fatalf("failed to mine block: %v", err)
	}

	block.Data = "modified data"

	before := chain.Blocks()

	if err := chain.AddBlock(block); err == nil {
		t.Fatal("expected invalid block to be rejected")
	}

	after := chain.Blocks()

	if len(before) != len(after) {
		t.Fatalf(
			"chain length changed after rejection: got %d, want %d",
			len(after),
			len(before),
		)
	}

	for i := range before {
		if before[i] != after[i] {
			t.Fatalf("chain changed at block %d after rejection", i)
		}
	}

	if err := chain.Validate(); err != nil {
		t.Errorf("chain became invalid after rejecting block: %v", err)
	}
}

func TestMineNextBlockDoesNotMutateChain(t *testing.T) {
	chain, err := NewChain()
	if err != nil {
		t.Fatalf("chain creation error: %v", err)
	}

	before := chain.Blocks()
	timestamp := before[len(before)-1].Header.Timestamp + 1

	if _, err := chain.MineNextBlock("test data", timestamp); err != nil {
		t.Fatalf("block mining error: %v", err)
	}

	if got := chain.Blocks(); !reflect.DeepEqual(got, before) {
		t.Fatal("MineNextBlock() mutated chain")
	}
}

func TestAddBlock(t *testing.T) {
	chain, err := NewChain()
	
	if err != nil {
		t.Fatalf("chain creation error: %v", err)
	}

	timestamp := chain.Blocks()[len(chain.Blocks())-1].Header.Timestamp + 1
	
	block, err  := chain.MineNextBlock("block 1", timestamp)

	if err != nil {
		t.Fatalf("block mining error - %v", err)
	}

	if err := chain.AddBlock(block); err != nil {
		t.Fatalf("block adding error - %v", err)
	}

	if chain.Len() != 2 {
		t.Fatalf(
			"chain length = %d, want 2",
			chain.Len(),
		)
	}

	if chain.Blocks()[1] != block {
		t.Fatalf("added block differs from supplied one")
	}

	if err := chain.Validate(); err != nil {
		t.Fatalf("chain invalid after adding block: %v", err)
	}
}

func TestAddBlockRejectsDuplicate(t *testing.T) {
	chain, err := NewChain()
	
	if err != nil {
		t.Fatalf("chain creation error: %v", err)
	}

	block, err := chain.MineNextBlock("duplicate", chain.Blocks()[0].Header.Timestamp+1)

	if err != nil {
		t.Fatalf("block mining error - %v", err)
	}

	if err := chain.AddBlock(block); err != nil {
		t.Fatalf("first add block error = %v", err)
	}
	if err := chain.AddBlock(block); err == nil {
		t.Fatal("second add block accepted the same block")
	}
	if chain.Len() != 2 {
		t.Fatalf(
			"chain length = %d after duplicate rejection, want 2",
			chain.Len(),
		)
	}
}

func TestEmptyChainOperationsFail(t *testing.T) {
	chain := &Chain{}

	if err := chain.Validate(); err == nil {
		t.Error("Validate() accepted empty chain")
	}

	if _, err := chain.MineNextBlock("data", 1); err == nil {
		t.Error("MineNextBlock() accepted empty chain")
	}
	
	if err := chain.AddBlock(Block{}); err == nil {
		t.Error("AddBlock() accepted empty chain")
	}
}