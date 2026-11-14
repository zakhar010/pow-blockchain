package main

import (
	"fmt"
)

func main() {
	chain := NewChain()

	chain.AddBlock("block 1")
	chain.AddBlock("block 2")

	for i, block := range chain.Blocks {
		fmt.Printf("block №%d:\n%s\n", i, block)
	}

	if err := chain.Validate(); err != nil {
		fmt.Println("error validating chain:", err)
	} else {
		fmt.Println("chain is valid")
	}

	fmt.Println("tampering with block 1")

	chain.Blocks[1].Data = "tampered data"

	if err := chain.Validate(); err != nil {
		fmt.Println("error validating chain:", err)
	} else {
		fmt.Println("chain is valid")
	}

	// remine the modified block
	fmt.Println("remining block 1")

	chain.Blocks[1] = MineBlock(chain.Blocks[1].Data, chain.Blocks[0].Hash)

	if err := chain.Validate(); err != nil {
		fmt.Println("error validating chain:", err)
	} else {
		fmt.Println("chain is valid")
	}

	fmt.Println("remining all remaining blocks")

	for i := 2; i < len(chain.Blocks); i++ {
		chain.Blocks[i] = MineBlock(
			chain.Blocks[i].Data,
			chain.Blocks[i-1].Hash,
		)
	}

	fmt.Println("validating chain after remining all blocks")
	if err := chain.Validate(); err != nil {
		fmt.Println("error validating chain:", err)
	} else {
		fmt.Println("chain is valid")
	}
}
