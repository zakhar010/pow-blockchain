package main

import (
	"fmt"
	"log"

	"github.com/zakhar010/pow-blockchain/internal/blockchain"
)

func main() {
	chain := blockchain.NewChain()

	for _, data := range []string{"block 1", "block 2"} {
		if err := chain.AddBlock(data); err != nil {
			log.Fatal(err)
		}
	}

	for i, block := range chain.Blocks() {
		fmt.Printf("block #%d:\n%s\n\n", i, block)
	}

	if err := chain.Validate(); err != nil {
		log.Fatalf("blockchain validation failed: %v", err)
	}

	fmt.Println("blockchain is valid")
}
