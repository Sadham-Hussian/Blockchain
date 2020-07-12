package blockchain

// Block : struct to define a block
type Block struct {
	Hash     []byte // Hash of the block
	Data     []byte // Data inside the block
	PrevHash []byte // Hash of the previous block
	Nonce    int    // a number to be found to find the Hash of the block
}

// CreateBlock : method to create a new block
func CreateBlock(data string, PrevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), PrevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}

// Genesis : method to create the genesis block.
// genesis block does not have PrevHash
func Genesis() *Block {
	return CreateBlock("genesis", []byte{})
}
