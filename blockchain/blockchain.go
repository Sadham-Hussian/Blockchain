package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./temp/blocks"
)

// Blockchain : struct to maintain the main chain of the block chain
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

// InitBlockchain : method to create a Blockchain
func InitBlockchain() *Blockchain {
	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions("/temp/blocks"))
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis approved")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)

			err = item.Value(func(val []byte) error {
				lastHash = val
				return err
			})

			return err
		}
	})
	Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

// AddBlock : method to add new block to the chain
func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, newBlock)
}
