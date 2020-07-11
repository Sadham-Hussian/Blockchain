package blockchain

// Blockchain : struct to maintain the main chain of the block chain
type Blockchain struct {
	Blocks []*Block
}

// Block : struct to define a block
type Block struct {
	Hash     []byte // Hash of the block
	Data     []byte // Data inside the block
	PrevHash []byte // Hash of the previous block
	Nonce    int    // a number to be found to find the Hash of the block
}

// CreateBlock : method to create a new block
func CreateBlock(data string, PrevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), PrevHash}
	pow := NewProof(block)
	nonce, hash = pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

// AddBlock : method to add new block to the chain
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}

// Genesis : method to create the genesis block.
// genesis block does not have PrevHash
func Genesis() *Block {
	return CreateBlock("genesis", []byte{})
}

// InitBlockchain : method to create a Blockchain
func InitBlockchain() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}
