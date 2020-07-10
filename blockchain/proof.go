package blockchain

import "math/big"

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
