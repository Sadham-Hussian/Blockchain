package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/Sadham-Hussian/Blockchain/wallet"
)

// TxInput : struct to handle transaction input of a transaction
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

// TxOutputs : struct to trace unspent transactions
type TxOutputs struct {
	Outputs []TxOutput
}

// TxOutput : struct to handle transaction output of a transaction
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// UsesKey checks whether the given pubkeyhash matches the pubkey in TxInput
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// Lock converts the address to PubKeyHash to lock the output transaction
func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks whether the Output Transaction is locked with correct pubKeyHash
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTxOutput creates new TxOutput by locking the transaction output
// with pubkeyHash of the given address
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// Serialize serializes TxOutputs to []byte
func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)

	return buffer.Bytes()
}

// DeserializeOutputs deserializes TxOutputs in byte to TxOutputs
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	Handle(err)

	return outputs
}
