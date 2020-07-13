package main

import (
	"fmt"
	"strconv"

	"github.com/Sadham-Hussian/Blockchain/blockchain"
)

// CommandLine : Struct to interact with the user
type CommandLine struct {
	blockchain *blockchain.Blockchain
}

func main() {
	chain := blockchain.InitBlockchain()

	chain.AddBlock("Alice pay 20 BTC to Bob")
	chain.AddBlock("Bob pay 40 BTC to Charlie")

	for _, blocks := range chain.Blocks {
		fmt.Printf("\n\nBlock data : %s", blocks.Data)
		fmt.Printf("\nBlock Hash : %x", blocks.Hash)
		fmt.Printf("\nPrev Block Hash : %x", blocks.PrevHash)

		pow := blockchain.NewProof(blocks)
		fmt.Printf("\nPoW: %s\n", strconv.FormatBool(pow.Validate()))
	}
	fmt.Println()
}
