package core

import (
	"../formatters"
	"../log"
	"./errors"
)

type Blockchain struct {
	sink         *BlockStore
	cfg          *BlockchainConfig
	validatorReg TxValidatorsRegistry
	log          *log.Logger
}

func NewBlockchain(config *BlockchainConfig, sink *BlockStore, validators TxValidatorsRegistry) *Blockchain {
	var logger = log.NewLog("blockchain")
	config.Assert()

	result := Blockchain{sink: sink, validatorReg: validators, log: logger, cfg: config}
	return &result
}

func (self *Blockchain) createGenesisBlock(transactions Transactions) (*Block, error) {
	result := NewGenesisBlock(transactions)

	err := self.sink.Add(result)

	if err != nil {
		self.log.Warn("genesis block creation failed", nil)
		return nil, errors.Wrap(err, "CreateGenesisBlock")
	}

	self.log.Info("genesis block created", log.More{"hash": formatters.HexStringFromRaw(result.Hash)})

	return result, nil
}

func (self Blockchain) Cursor() (*BlockCursor, error) {
	var sink = self.sink
	var result, err = sink.Cursor()
	return result, err
}

func (self *Blockchain) AddTransaction(tx *Transaction) error {
	err := self.validateTx(tx)
	if err != nil {
		return err
	}

	err = processTx(tx)
	if err != nil {
		return err
	}

	err = self.addTransactionToChain(tx)

	return err
}

func (self *Blockchain) addTransactionToChain(tx *Transaction) error {
	var genesis, err = self.sink.GenesisBlock()
	if err != nil {
		return err
	}

	if genesis == nil {
		genesisTx, err := self.cfg.CreateGenesisTransactions()
		if err != nil {
			return err
		}

		genesis, err = self.createGenesisBlock(genesisTx)
	}

	if err != nil {
		return err
	}

	tip, err := self.sink.Tip()
	if err != nil {
		return err
	}

	var newBlock = NewBlock(Transactions{tx}, tip.Hash)
	err = self.sink.Add(newBlock)

	return err
}

func (self *Blockchain) validateTx(tx *Transaction) error {
	validators, err := self.validatorReg.Get(tx)
	if err != nil {
		return err
	}

	for _, validator := range validators {
		err = validator.Validate(tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func processTx(tx *Transaction) error {
	return nil
}
