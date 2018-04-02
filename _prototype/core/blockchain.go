package core

import (
	"errors"
	"strings"

	"../common/prompt"
)

type Blockchain struct {
	sink  BlockStore
	proof ProofAlgorithm
}

func (self Blockchain) Cursor() (*BlockCursor, error) {
	var sink = self.sink
	var result, err = sink.Cursor()
	return result, err
}

func (self Blockchain) ProofAlgorithm() ProofAlgorithm {
	return self.proof
}

func (self *Blockchain) AddNewBlock(data string) error {
	var sink = self.sink
	var last, err = sink.Tip()
	if err != nil {
		return err
	}
	if last == nil {
		return errors.New("add new block failed: no genesis block")
	}

	return err
}

func (self *Blockchain) HasGenesisBlock() (bool, error) {
	var sink = self.sink
	var block, err = sink.GenesisBlock()
	var result = block != nil
	return result, err
}

func (self *Blockchain) CreateGenesisBlockIfNotExists(rewardAddress string) error {
	var hasGenesisBlock, err = self.HasGenesisBlock()
	if err != nil {
		return err
	}
	if hasGenesisBlock {
		return err
	}

	rewardAddress = strings.TrimSpace(rewardAddress)
	if len(rewardAddress) <= 0 {
		return errors.New("no reward address for the genesis coinbase transaction specified")
	}

	prompt.Shared().Info("mining genesis block with reward to %v ...", rewardAddress)
	var coinbaseTx = NewCoinbaseTransaction(rewardAddress, "genesis coinbase tx")
	var genesisBlock = NewGenesisBlock(coinbaseTx)
	err = self.proof.Proof(genesisBlock)
	if err != nil {
		return err
	}

	prompt.Shared().Success("genesis block mined, reward of %vℍ was transferred to %v", coinbaseTx.getOutValue(), rewardAddress)

	prompt.Shared().Info("store blockchain ...")
	err = self.sink.Add(genesisBlock)
	if err != nil {
		return err
	}

	prompt.Shared().Success("blockchain genesis block created")
	return nil
}

func NewBlockchain(sink BlockStore) *Blockchain {
	var self = new(Blockchain)
	self.sink = sink
	self.proof = NewHashCash()

	return self
}
