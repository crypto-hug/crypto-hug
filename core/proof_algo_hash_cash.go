package core

import (
	"math/big"
	"math"
	"../common/prompt"

	"fmt"

	"errors"
)

const PROOF_TARGET_BITS = 5



type HashCash struct {
	proofTarget *big.Int
}

func NewHashCash() *HashCash {
	var proofTarget = big.NewInt(1)
	proofTarget.Lsh(proofTarget, uint(256-PROOF_TARGET_BITS))

	var result = &HashCash{proofTarget}

	return result
}

func (self *HashCash) Validate(block *Block) bool{
	var hashAsInt big.Int
	var hash = block.GetHash(PROOF_TARGET_BITS)
	hashAsInt.SetBytes(hash)
	var valid = hashAsInt.Cmp(self.proofTarget) == -1
	return valid
}

func (self *HashCash) Proof(block *Block) error{
	var hashAsInt big.Int
	var hash [] byte
	var nonce = 0
	var max = math.MaxInt64
	for nonce < max{

		prompt.Shared().Debug("Mining block with nonce %v", nonce)
		block.Nonce = nonce
		hash = block.GetHash(PROOF_TARGET_BITS)
		hashAsInt.SetBytes(hash)
		if hashAsInt.Cmp(self.proofTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	if nonce == max{
		return errors.New(fmt.Sprintf("reached max nonce of %v", max))
	}

	block.Nonce = nonce
	block.Hash = hash

	return nil
}