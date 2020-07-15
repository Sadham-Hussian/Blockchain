package blockchain

// Transaction : struct to handle details of a transaction
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// TxInput : struct to handle transaction input of a transaction
type TxInput struct {
	ID  []byte
	Sig []byte
	Out int
}

// TxOutput : struct to handle transaction output of a transaction
type TxOutput struct {
	Value  int
	PubKey string
}
