package blockchain

var (
	utxoPrefix   = []byte("utxo-")
	prefixLength = len(utxoPrefix)
)

// UTXOSet struct to access blockchain database
type UTXOSet struct {
	BChain *Blockchain
}
