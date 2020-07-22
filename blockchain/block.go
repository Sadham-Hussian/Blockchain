package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// Block : struct to define a block
type Block struct {
	Hash         []byte         // Hash of the block
	Transactions []*Transaction // Transactions in the block
	PrevHash     []byte         // Hash of the previous block
	Nonce        int            // a number to be found to find the Hash of the block
}

// HashTransactions : Function to hash all the transaction in the block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Hash())
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// CreateBlock : method to create a new block
func CreateBlock(txs []*Transaction, PrevHash []byte) *Block {
	block := &Block{[]byte{}, txs, PrevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}

// Genesis : method to create the genesis block.
// genesis block does not have PrevHash
func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
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
