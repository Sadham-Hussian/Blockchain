package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
)

// Difficulty : measure of how difficult it is to find the hash within the target
const Difficulty = 16

// ProofOfWork : struct to perform proofofwork to achieve consensus
type ProofOfWork struct {
	Block  *Block   // Block to which proofofwork is computed
	Target *big.Int // Target within which the hash value must be found
}

// NewProof : computes new proofofwork for the block.
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

func (p *ProofOfWork) Run() (int, []byte) {

}

// ToHex : method to encode the nonce and difficulty
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
