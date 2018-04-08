package core

import (
	"github.com/crypto-hug/crypto-hug/errors"
	"github.com/crypto-hug/crypto-hug/formatters"
	"github.com/crypto-hug/crypto-hug/log"
)

type Blockchain struct {
	sink       *BlockStore
	cfg        *BlockchainConfig
	walletSink WalletSink
	assetSink  AssetSink
	log        *log.Logger
}

func NewBlockchain(config *BlockchainConfig, sink *BlockStore, walletSink WalletSink, assetSink AssetSink) *Blockchain {
	var logger = log.NewLog("blockchain")
	config.Assert()

	result := Blockchain{sink: sink, log: logger, cfg: config,
		walletSink: walletSink, assetSink: assetSink}
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
	processors := self.createProcessorsFor(tx)
	err := self.validateTx(tx, processors)
	if err != nil {
		return err
	}

	err = self.processTx(tx, processors)
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

func (self *Blockchain) validateTx(tx *Transaction, processors *TransactionProcessors) error {

	for _, processor := range *processors {
		processor.Setup(self.walletSink, self.assetSink)
		self.log.Debug("begin tx validation", log.More{"processor": processor.Name(), "tx": formatters.HexStringFromRaw(tx.Hash)})

		processor.Setup(self.walletSink, self.assetSink)

		if err := processor.Prepare(tx); err != nil {
			self.log.Error("tx failed validation", log.More{"processor": processor.Name(), "tx": formatters.HexStringFromRaw(tx.Hash), "err": err.Error()})
			return err
		}

		self.log.Info("tx validated", log.More{"processor": processor.Name(), "tx": formatters.HexStringFromRaw(tx.Hash)})
	}

	return nil
}

func (self *Blockchain) processTx(tx *Transaction, processors *TransactionProcessors) error {
	for _, processor := range *processors {
		self.log.Debug("begin tx processing", log.More{"processor": processor.Name(), "tx": formatters.HexStringFromRaw(tx.Hash)})
		err := processor.Process(tx)

		if err != nil {
			self.log.Error("tx failed processing", log.More{"processor": processor.Name(), "tx": formatters.HexStringFromRaw(tx.Hash), "err": err.Error()})
			return err
		}

		self.log.Info("tx processed", log.More{"processor": processor.Name(), "tx": formatters.HexStringFromRaw(tx.Hash)})
	}

	return nil
}

func (self *Blockchain) createProcessorsFor(tx *Transaction) *TransactionProcessors {
	result := TransactionProcessors{}
	for _, processor := range self.cfg.CreateTxProcessors() {
		if processor.ShouldProcess(tx) {
			result = append(result, processor)
		}
	}

	return &result
}
