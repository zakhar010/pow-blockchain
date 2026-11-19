package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zakhar010/pow-blockchain/internal/blockchain"
)

func main() {
	chain, err := blockchain.NewChain()

	if err != nil {
		log.Fatalf("failed to create blockchain: %v", err)
	}

	timestamp := chain.Blocks()[0].Header.Timestamp

	for i := 1; i <= 101; i++ {
		if i <= blockchain.RetargetInterval {
			timestamp += 2 * blockchain.TargetBlockTime
		} else {
			timestamp += blockchain.TargetBlockTime
		}

		data := fmt.Sprintf("block %d", i)

		log.Printf("starting mining %s\n", data)
		startTime := time.Now()

		if err := chain.AddBlock(data, timestamp); err != nil {
			log.Fatal(err)
		}

		log.Printf("mined %s in %s\n", data, time.Since(startTime))
	}

	for i, block := range chain.Blocks() {
		fmt.Printf("block #%d:\n%s\n\n", i, block)
	}

	if err := chain.Validate(); err != nil {
		log.Fatalf("blockchain validation failed: %v", err)
	}

	fmt.Println("blockchain is valid")
}
