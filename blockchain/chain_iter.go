package blockchain

import "github.com/dgraph-io/badger"

// BlockchainIterator : struct to get the blockchain and blocks from the
// BadgerDB
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
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
