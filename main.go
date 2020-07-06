package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// struct to maintain the main chain of the block chain
type Blockchain struct {
	blocks []*Block
}

// struct to define a block
type Block struct {
	hash     []byte
	data     []byte
	prevHash []byte
}

// method to compute the hash of the block with data and prevHash
func (b *Block) computeHash() {
	info := bytes.Join([][]byte{b.data, b.prevHash}, []byte{})
	hash_generated := sha256.Sum256(info)
	b.hash = hash_generated[:]
}

// method to create a new block
func createBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte{data}, prevHash}
	block.computeHash()
	return block
}

func main() {
	fmt.Println("Hello, World.")
}
