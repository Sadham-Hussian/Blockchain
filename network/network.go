package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"

	"github.com/Sadham-Hussian/Blockchain/blockchain"
	"gopkg.in/vrecan/death.v3"
)

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

// Version represent the version of the blockchain version is incremented based on
// no of blocks in the chain. BestHeight holds the height of the blockchain
type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

// CmdToBytes function to convert cmd in string to bytes
func CmdToBytes(cmd string) []byte {
	var bytes [commandLength]byte

	for i, c := range cmd {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// BytesToCmd function converts bytes to string
func BytesToCmd(bytes []byte) string {
	var cmd []byte

	for _, b := range bytes {
		if b != 0x0 {
			cmd = append(cmd, b)
		}
	}

	return fmt.Sprintf("%s", cmd)
}

// RequestBlocks function request blocks from the knownNodes
func RequestBlocks() {
	for _, node := range knownNodes {
		SendGetBlocks(node)
	}
}

// SendGetBlocks function sends request to send blocks to another node
func SendGetBlocks(address string) {
	payload := GobEncode(GetBlocks{nodeAddress})
	request := append(CmdToBytes("getblocks"), payload...)

	SendData(address, request)
}

// SendGetData function sends request to send data to another node
func SendGetData(address, kind string, ID []byte) {
	payload := GobEncode(GetData{nodeAddress, kind, ID})
	request := append(CmdToBytes("getdata"), payload...)

	SendData(address, request)
}

// SendAddr sends the knownAddress of a node to another node
func SendAddr(address string) {
	nodes := Addr{knownNodes}
	nodes.AddrList = append(knownNodes, nodeAddress)
	payload := GobEncode(nodes)
	request := append(CmdToBytes("addr"), payload...)

	SendData(address, request)
}

// SendBlock sends the block from one node to another
func SendBlock(address string, b *blockchain.Block) {
	blockData := Block{address, b.Serialize()}
	payload := GobEncode(blockData)
	request := append(CmdToBytes("block"), payload...)

	SendData(address, request)
}

// SendInv sends the inventory data from one node to another
func SendInv(address, kind string, items [][]byte) {
	inventory := Inventory{address, kind, items}
	payload := GobEncode(inventory)
	request := append(CmdToBytes("inv"), payload...)

	SendData(address, request)
}

// SendTx sends a transaction from one node to another node
func SendTx(address string, tnx *blockchain.Transaction) {
	data := Tx{nodeAddress, tnx.Serialize()}
	payload := GobEncode(data)
	request := append(CmdToBytes("tx"), payload...)

	SendData(address, request)
}

// SendVersion function sends the version from a node to another node
func SendVersion(address string, chain *blockchain.Blockchain) {
	bestHeight := chain.BestHeight()
	version := Version{address, bestHeight, nodeAddress}
	payload := GobEncode(version)
	request := append(CmdToBytes("version"), payload...)

	SendData(address, request)
}

// SendData sends data to a node
func SendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)

	if err != nil {
		fmt.Printf("%s is not available\n", addr)

		var updateNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updateNodes = append(updateNodes, node)
			}
		}

		knownNodes = updateNodes

		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

// HandleAddr adds the nodeAddress sent by another node into its own node
func HandleAddr(request []byte) {
	var buff bytes.Buffer
	var payload Addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes\n", len(knownNodes))
	RequestBlocks()
}

// HandleBlock functions handles the blocks received from another node
func HandleBlock(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload Block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := blockchain.Deserialize(blockData)

	fmt.Println("Received a new block!")
	chain.AddBlock(block)

	fmt.Printf("Added block: %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		SendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{chain}
		UTXOSet.Reindex()
	}
}

// HandleGetBlocks functions handles the request from a node which is requesting
// a block
func HandleGetBlocks(request []byte, chain *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload GetBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := chain.GetBlockHashes()
	SendInv(payload.AddrFrom, "block", blocks)
}

// GobEncode encodes the commands to be passed through the network
func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// CloseDB used to close DB before terminating
func CloseDB(chain *blockchain.Blockchain) {
	d := death.NewDeath(syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d.WaitForDeathWithFunc(func() {
		defer os.Exit(1)
		defer runtime.Goexit()
		chain.Database.Close()
	})
}
