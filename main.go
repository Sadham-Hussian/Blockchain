package main

import (
	"fmt"
)

type Blockchain struct {
	blocks []*Block
}

type Block struct {
	hash     []byte
	data     []byte
	prevHash []byte
}

func main() {
	fmt.Println("Hello, World.")
}
