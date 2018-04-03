package core

import (
	"../formatters"
	"../log"
	"github.com/pkg/errors"
)

type Blockchain struct {
	sink         *BlockStore
	validatorReg TxValidatorsRegistry
	log          *log.Logger
}

func NewBlockchain(sink *BlockStore, validators TxValidatorsRegistry) *Blockchain {
	var logger = log.NewLog("blockchain")

	result := Blockchain{sink: sink, validatorReg: validators, log: logger}
	return &result
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
		genesis, err = self.createGenesisBlock()
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

func (self *Blockchain) createGenesisBlock() (*Block, error) {
	myAddr := "VB6QzPAL7P83N48MhoFdLXuroxPmUiphp"
	self.log.Info("create genesis block", log.More{"reward": myAddr})
	genesisOwnerAddress, err := NewAddressFromString(myAddr)
	if err != nil {
		return nil, err
	}
	spawnHugs := Transactions{}

	for i := 0; i < 3; i++ {
		var spawnTx *Transaction = nil
		spawnTx, err = NewSpawnHugTransaction(genesisOwnerAddress)
		if err != nil {
			return nil, err
		}
		spawnHugs = append(spawnHugs, spawnTx)
	}

	result := NewGenesisBlock(spawnHugs)

	err = self.sink.Add(result)

	if err != nil {
		self.log.Warn("genesis block creation failed", nil)
		return nil, errors.Wrap(err, "createGenesisBlock")
	}

	self.log.Info("genesis block created", log.More{"hash": formatters.HexStringFromRaw(result.Hash)})

	return result, err
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
