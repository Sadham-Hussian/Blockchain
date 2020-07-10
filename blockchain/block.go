package blockchain

import (
	"bytes"
	"crypto/sha256"
)

// Blockchain : struct to maintain the main chain of the block chain
type Blockchain struct {
	Blocks []*Block
}

// Block : struct to define a block
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

// ComputeHash : method to compute the hash of the block with data and prevHash
func (b *Block) ComputeHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hashGenerated := sha256.Sum256(info)
	b.Hash = hashGenerated[:]
}

// CreateBlock : method to create a new block
func CreateBlock(data string, PrevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), PrevHash}
	block.ComputeHash()
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
