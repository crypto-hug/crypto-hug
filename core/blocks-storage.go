package core

import (
	"github.com/crypto-hug/crypto-hug/errors"
)

type BlockStore struct {
	stat BlockStats
	sink BlockSink
}
type BlockCursor struct {
	store   BlockStore
	current *Block
}

func NewBlockStore(sink BlockSink, stat BlockStats) *BlockStore {
	var result = BlockStore{stat: stat, sink: sink}
	return &result
}

func (self *BlockStore) Add(block *Block) error {

	err := self.sink.Put(block)
	if err != nil {
		return errors.Wrap(err, "BlockStore:Add")
	}

	err = self.stat.PutTip(block.Hash)
	if err != nil {
		return errors.Wrap(err, "BlockStore:Add")
	}

	if block.IsGenesisBlock() {
		err = self.stat.PutGenesis(block.Hash)
	}

	return errors.Wrap(err, "BlockStore:Add")
}

func (self *BlockStore) GenesisBlock() (*Block, error) {

	gen, err := self.stat.GetGenesis()
	if err != nil {
		return nil, errors.Wrap(err, "BlockStore:GenesisBlock")
	}

	if gen == nil || len(gen) <= 0 {
		return nil, nil
	}

	genB, err := self.sink.Get(gen)
	if err != nil {
		return nil, errors.Wrap(err, "BlockStore:GenesisBlock")
	}

	return genB, nil
}

func (self *BlockStore) Tip() (*Block, error) {
	tip, err := self.stat.GetTip()
	if err != nil {
		return nil, errors.Wrap(err, "BlockStore:Tip")
	}

	if tip == nil || len(tip) <= 0 {
		return nil, nil
	}

	tipB, err := self.sink.Get(tip)
	if err != nil {
		return nil, errors.Wrap(err, "BlockStore:Tip")
	}

	return tipB, nil
}

func (self *BlockStore) Cursor() (*BlockCursor, error) {
	var tip, err = self.Tip()
	if err != nil {
		return nil, errors.Wrap(err, "BlockStore:Cursor")
	}

	var result = BlockCursor{store: *self, current: tip}

	return &result, nil
}

func (self *BlockCursor) Current() *Block {
	return self.current
}

func (self *BlockCursor) Reset() error {
	var tip, err = self.store.Tip()
	if err != nil {
		return errors.Wrap(err, "BoltBlockStore:Reset")
	}

	self.current = tip
	return nil
}

func (self *BlockCursor) Next() (bool, error) {
	var tip, err = self.prev(self.current)
	if err != nil {
		return false, errors.Wrap(err, "BoltBlockStore:Next")
	}
	if tip == nil {
		return false, nil
	}

	self.current = tip
	return true, nil
}

func (self *BlockCursor) prev(current *Block) (*Block, error) {
	if current == nil || current.PrevHash == nil || len(current.PrevHash) <= 0 {
		return nil, nil
	}

	result, err := self.store.sink.Get(current.PrevHash)
	if err != nil {
		return nil, errors.Wrap(err, "BlockCursor:prev")
	}

	return result, nil
}
