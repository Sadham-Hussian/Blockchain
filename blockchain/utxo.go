package blockchain

import (
	"bytes"
	"encoding/hex"
	"log"

	"github.com/dgraph-io/badger"
)

var (
	utxoPrefix   = []byte("utxo-")
	prefixLength = len(utxoPrefix)
)

// UTXOSet struct to access blockchain database
type UTXOSet struct {
	BChain *Blockchain
}

// CountTransactions counts the no of keys with utxoPrefix and
// returns the no of unspent transaction output
func (u UTXOSet) CountTransactions() int {
	db := u.BChain.Database
	counter := 0

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			counter++
		}
		return nil
	})
	Handle(err)

	return counter
}

// FindUTXO retrieves the transaction in the database with utxoPrefix that is
// finds the Unspent Transaction Output
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput

	db := u.BChain.Database

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			var outs TxOutputs
			item := it.Item()

			err := item.Value(func(val []byte) error {
				outs = DeserializeOutputs(val)
				return nil
			})
			Handle(err)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	Handle(err)
	return UTXOs
}

// FindSpendableOutputs : function to find the outputs that a address can spend
// in another transaction
func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0
	db := u.BChain.Database

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(utxoPrefix); it.ValidForPrefix(utxoPrefix); it.Next() {
			var outs TxOutputs
			item := it.Item()
			k := item.Key()
			err := item.Value(func(val []byte) error {
				outs = DeserializeOutputs(val)
				return nil
			})
			Handle(err)

			k = bytes.TrimPrefix(k, utxoPrefix)
			txID := hex.EncodeToString(k)

			for outIndex, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOuts[txID] = append(unspentOuts[txID], outIndex)
				}
			}
		}
		return nil
	})
	Handle(err)
	return accumulated, unspentOuts
}

// Reindex removes utxoPrefix of the UTXO in the database and adds utxoPrefix
// to new UTXOs
func (u UTXOSet) Reindex() {
	db := u.BChain.Database

	u.DeleteByPrefix(utxoPrefix)

	UTXO := u.BChain.FindUTXO()

	err := db.Update(func(txn *badger.Txn) error {
		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			key = append(utxoPrefix, key...)

			err = txn.Set(key, outs.Serialize())
			Handle(err)
		}
		return nil
	})
	Handle(err)
}

// Update updates the UTXOSet in the database
func (u *UTXOSet) Update(b *Block) {
	db := u.BChain.Database

	err := db.Update(func(txn *badger.Txn) error {
		for _, tx := range b.Transactions {
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					var outs TxOutputs
					updatedOuts := TxOutputs{}
					inID := append(utxoPrefix, in.ID...)
					item, err := txn.Get(inID)
					Handle(err)
					err = item.Value(func(val []byte) error {
						outs = DeserializeOutputs(val)
						return err
					})
					Handle(err)

					for outIndex, out := range outs.Outputs {
						if outIndex != in.Out {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						if err := txn.Delete(inID); err != nil {
							log.Panic(err)
						}
					} else {
						if err := txn.Set(inID, updatedOuts.Serialize()); err != nil {
							log.Panic(err)
						}
					}
				}
			}

			newOutputs := TxOutputs{}
			for _, out := range tx.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			txID := append(utxoPrefix, tx.ID...)
			if err := txn.Set(txID, newOutputs.Serialize()); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	Handle(err)
}

// DeleteByPrefix deletes key in database with prefix utxoPrefix
func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := u.BChain.Database.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 10000
	u.BChain.Database.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++

			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					log.Panic(err)
				}
				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}
		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}
