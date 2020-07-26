package network

import "github.com/Sadham-Hussian/Blockchain/blockchain"

const (
	protocol      = "tcp"
	version       = 1
	commandLength = 12
)

var (
	nodeAddress     string
	mineAddress     string
	knownNodes      = []string{"localhost:3000"}
	blocksInTransit = [][]byte{}
	memoryPool      = make(map[string]blockchain.Transaction)
)

// Addr handles the nodeAddress of nodes that are connected to the node in the network
type Addr struct {
	AddrList []string
}

// Block handles the value of the block and the node address from where the block
// is shared
type Block struct {
	AddrFrom string
	Block    []byte
}

// GetBlocks copies blockchain from another node
type GetBlocks struct {
	AddrFrom string
}

// GetData gets info(ID) about the data
type GetData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

// Inventory gets data from a node
type Inventory struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

// Tx represents a transaction
type Tx struct {
	AddrFrom    string
	Transaction []byte
}

// Version represent the version
type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}
