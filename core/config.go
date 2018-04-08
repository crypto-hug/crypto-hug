package core

import (
	"github.com/crypto-hug/crypto-hug/errors"
)

type BlockchainConfig struct {
	CreateGenesisTransactions func() (Transactions, error)
	CreateTxProcessors        func() TransactionProcessors
}

func NewBlockchainConfig() *BlockchainConfig {
	result := BlockchainConfig{}
	return &result
}

func (self *BlockchainConfig) Assert() {
	errors.AssertNotNil("BlockchainConfig.CreateGenesisTransactions", self.CreateGenesisTransactions)
	errors.AssertNotNil("BlockchainConfig.CreateTxProcessors", self.CreateTxProcessors)
}
