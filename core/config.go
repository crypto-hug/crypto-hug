package core

import (
	"github.com/crypto-hug/crypto-hug/errors"
)

type BlockchainConfig struct {
	CreateGenesisTransactions func() (Transactions, error)
	TransactionProcessors     TransactionProcessors
}

func NewBlockchainConfig() *BlockchainConfig {
	result := BlockchainConfig{}
	return &result
}

func (self *BlockchainConfig) Assert() {
	errors.AssertNotNil("BlockchainConfig.CreateGenesisTransactions", self.CreateGenesisTransactions)
	errors.AssertTrue("BlockchainConfig.TransactionProcessors", len(self.TransactionProcessors) > 0, "is empty")
}
