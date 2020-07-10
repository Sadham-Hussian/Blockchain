package main

import (
	"fmt"

	"github.com/Sadham-Hussian/Blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockchain()

	chain.AddBlock("Alice pay 20 BTC to Bob")
	chain.AddBlock("Bob pay 40 BTC to Charlie")

	for _, blocks := range chain.Blocks {
		fmt.Printf("\n\nBlock data : %s", blocks.Data)
		fmt.Printf("\nBlock Hash : %x", blocks.Hash)
		fmt.Printf("\nPrev Block Hash : %x", blocks.PrevHash)
	}
}
