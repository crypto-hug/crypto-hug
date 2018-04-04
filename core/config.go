package core

import (
	"./errors"
)

type BlockchainConfig struct {
	CreateGenesisTransactions func() (Transactions, error)
}

func NewBlockchainConfig(createGenesisTransactions func() (Transactions, error)) *BlockchainConfig {
	result := BlockchainConfig{CreateGenesisTransactions: createGenesisTransactions}
	return &result
}

func (self *BlockchainConfig) Assert() {
	errors.AssertNotNil("BlockchainConfig.CreateGenesisTransactions", self.CreateGenesisTransactions)
}
