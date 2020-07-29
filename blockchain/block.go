package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block : struct to define a block
type Block struct {
	Timestamp    int64          // Time block created
	Hash         []byte         // Hash of the block
	Transactions []*Transaction // Transactions in the block
	PrevHash     []byte         // Hash of the previous block
	Nonce        int            // a number to be found to find the Hash of the block
	Height       int            // Height of the blockchain
}

// HashTransactions : Function to hash all the transaction in the block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Serialize())
	}

	tree := NewMerkleTree(txHashes)

	return tree.RootNode.Data
}

// CreateBlock : method to create a new block
func CreateBlock(txs []*Transaction, PrevHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), []byte{}, txs, PrevHash, 0, height}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}

// Genesis : method to create the genesis block.
// genesis block does not have PrevHash
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{}, 0)
}

// Serialize : method to serialize the entire block to bytes to store the
// block in BadgerDB
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

// Deserialize : method to deserialize the block data retrieved from the BadgerDB
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

// Handle : method to handle errors
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
