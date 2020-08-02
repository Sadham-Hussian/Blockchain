package wallet

import (
	"fmt"
	"log"

	"github.com/mr-tron/base58"
)

// Base58Encode encodes the input data to return a address
func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

// Base58Decode decodes the input data
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	fmt.Println(decode)
	if err != nil {
		log.Panic(err)
	}

	return decode
}
