package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
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

// FindDataToHash : Hash of the block in a blockchain includes Data of the block,
// previous block's hash, nonce and difficulty. This function returns all these data
// appended in the form of []byte
func (pow *ProofOfWork) FindDataToHash(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.Data,
			pow.Block.PrevHash,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

// Run : returns nonce and hash of the block after computing nonce to the required target
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.FindDataToHash(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// Validate : method to validate the hash of the block
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	var hash [32]byte

	data := pow.FindDataToHash(pow.Block.Nonce)
	hash = sha256.Sum256(data)

	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

// ToHex : method to encode the nonce and difficulty
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)

	Handle(err)

	return buff.Bytes()
}
