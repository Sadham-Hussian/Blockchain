package blockchain

import (
	"bytes"
	"encoding/gob"
)

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

// Serialize : method to serialize the entire block to bytes to store the
// block in BadgerDB
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	return res.Bytes()
}

// Deserialize : method to deserialize the block data retrieved from the BadgerDB
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	return &block
}
