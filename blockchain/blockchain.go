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

// BlockchainIterator : struct to get the blockchain and blocks from the
// BadgerDB
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// InitBlockchain : method to create a Blockchain
func InitBlockchain() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	// db, err := badger.Open(badger.DefaultOptions("/temp/badger"))
	db, err := badger.Open(opts)
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
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return err
		})
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})

	Handle(err)
}

// Iterator : It returns the current Blockchain
func (chain *Blockchain) Iterator() *BlockchainIterator {
	iter := &BlockchainIterator{chain.LastHash, chain.Database}

	return iter
}

// Next : function to iterate the blockchain in BadgerDB and return a block
func (iter *BlockchainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return err
		})
		return err
	})
	Handle(err)
	iter.CurrentHash = block.PrevHash

	return block
}
