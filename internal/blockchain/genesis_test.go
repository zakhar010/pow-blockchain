package blockchain

import "testing"

func TestGenesisBlock(t *testing.T) {
	genesis := genesisBlock()

	if got := genesis.Hash(); got != genesisHash {
		t.Fatalf(
			"incorrect genesis hash: got - %x, want - %x",
			got,
			genesisHash,
		)
	}

	if !meetsTarget(genesis.Hash(), genesis.Header.Target) {
		t.Fatalf(
			"genesis hash %x doesn't satisfy target %x",
			genesis.Hash(),
			genesis.Header.Target,
		)
	}
} 

func TestNewChainContainsOnlyGenesis(t *testing.T) {
	chain, err := NewChain()

	if err != nil {
		t.Fatalf("chain creation error: %v", err)
	}

	if got := chain.Len(); got != 1 {
		t.Fatalf(
			"chain length: got - %d, want - 1",
			got,
		)
	}

	if got := chain.blocks[0]; got != genesisBlock() {
		t.Fatalf("first block is not the canonical genesis block")
	}
}