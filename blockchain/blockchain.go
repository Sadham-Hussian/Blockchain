package blockchain

// Blockchain : struct to maintain the main chain of the block chain
type Blockchain struct {
	Blocks []*Block
}

// InitBlockchain : method to create a Blockchain
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

// AddBlock : method to add new block to the chain
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}
