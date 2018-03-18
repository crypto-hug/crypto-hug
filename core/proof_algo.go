package core

type ProofAlgorithm interface {
	Validate(block *Block) bool
	Proof(block *Block) error
}