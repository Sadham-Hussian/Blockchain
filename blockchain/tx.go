package blockchain

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
