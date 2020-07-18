package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// Transaction : struct to handle details of a transaction
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// TxInput : struct to handle transaction input of a transaction
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

// TxOutput : struct to handle transaction output of a transaction
type TxOutput struct {
	Value  int
	PubKey string
}

// SetID : function to set the transaction ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var Hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	Handle(err)

	Hash = sha256.Sum256(encoded.Bytes())
	tx.ID = Hash[:]
}

// CoinbaseTx : firsttx in the block. Miner collects the block reward.
func CoinbaseTx(to string, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{100, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

// NewTransaction : function to create new transaction
func NewTransaction(from, to string, amount int, chain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulatedAmount, validOutputs := chain.FindSpendableOutputs(from, amount)

	if accumulatedAmount < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if accumulatedAmount > amount {
		outputs = append(outputs, TxOutput{accumulatedAmount - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

// IsCoinbase : function to check if the transaction is Coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// InputCanBeUnlock : function to check if the input of the transaction
// comes from the correct owner
func (in *TxInput) InputCanBeUnlock(data string) bool {
	return in.Sig == data
}

// OutputCanBeUnlocked : function to check whether the reciever can retrieve
// the transaction
func (out *TxOutput) OutputCanBeUnlocked(key string) bool {
	return out.PubKey == key
}
