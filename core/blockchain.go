package core

import (
	"../common/prompt"
	"errors"
	"strings"
)

type Blockchain struct {
	sink 	BlockStore
	proof  	ProofAlgorithm
}



func (self Blockchain) Cursor() (*BlockCursor, error) {
	var sink = self.sink
	var result, err = sink.Cursor()
	return result, err
}

func (self Blockchain) ProofAlgorithm() ProofAlgorithm{
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


	//var newBlock = NewBlock(data, last.Hash)
	//err = self.proof.Proof(newBlock)
	//if err != nil{
	//	return err
	//}
	//
	//err = sink.Add(newBlock)

	return err
}


func NewBlockchain(sink BlockStore, address string) (*Blockchain, error) {

	var last, err = sink.Tip()
	if err != nil {
		return nil, err
	}

	var self = new(Blockchain)
	self.sink = sink
	self.proof = NewHashCash()

	if last == nil {
		address = strings.TrimSpace(address)
		if len(address) <= 0{
			return nil, errors.New("no reward address for the genesis coinbase transaction specified")
		}

		prompt.Shared().Info("mining genesis block with reward to %v ...", address)
		var coinbaseTx = NewCoinbaseTransaction(address, "genesis coinbase tx")
		last = NewGenesisBlock(coinbaseTx)
		err = self.proof.Proof(last)
		if err != nil {
			return nil, err
		}

		prompt.Shared().Success("genesis block mined, reward of %vℍ was transferred to %v", coinbaseTx.getOutValue(), address)

		prompt.Shared().Info("store blockchain ...")
		err = sink.Add(last)
		if err != nil {
			return nil, err
		}

		prompt.Shared().Success("blockchain stored")

	}



	return self, nil
}
