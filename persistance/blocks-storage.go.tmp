package hugdb

import (
	"../core"
	"../core/errors"
	"./blocks"
	"github.com/boltdb/bolt"
)

type BoltBlockStore struct {
	db     *BoltDb
	blocks *blocks.FsBlockSink
}
type BoltBlockCursor struct {
	store   *BoltBlockStore
	current *core.Block
}

func NewBoltBlockStore(db *BoltDb) *BoltBlockStore {
	blocks := blocks.NewFsBlockSink("./blockhain_data/blocks")
	var result = BoltBlockStore{db: db, blocks: blocks}
	return &result
}

func (self *BoltBlockStore) Add(block *core.Block) error {
	err := self.db.GetDb().Update(func(tx *bolt.Tx) error {
		err := self.blocks.Put(block)

		if err != nil {
			return err
		}

		err = self.db.Set(tx, self.db.BucketBlockStats(), self.db.KeyBlocksStatsLast(), block.Hash)
		if err != nil {
			return err
		}
		if block.IsGenesisBlock() {
			err = self.db.Set(tx, self.db.BucketBlockStats(), self.db.KeyBlocksStatsGenesis(), block.Hash)
		}

		return err
	})

	return errors.Wrap(err, "BoltBlockStore:Add")
}

func (self *BoltBlockStore) GenesisBlock() (*core.Block, error) {
	var result *core.Block = nil
	var err = self.db.GetDb().View(func(tx *bolt.Tx) error {
		genHsh := self.db.Get(tx, self.db.BucketBlockStats(), self.db.KeyBlocksStatsGenesis())
		if len(genHsh) <= 0 {
			return nil
		}

		block, err := self.blocks.Get(genHsh)
		result = block

		return err
	})

	if err != nil {
		result = nil
	}

	return result, errors.Wrap(err, "BoltBlockStore:GenesisBlock")
}

func (self *BoltBlockStore) Tip() (*core.Block, error) {
	var result *core.Block = nil
	var err = self.db.GetDb().View(func(tx *bolt.Tx) error {
		tipHsh := self.db.Get(tx, self.db.BucketBlockStats(), self.db.KeyBlocksStatsLast())
		if len(tipHsh) <= 0 {
			return errors.NoGenesisBlock
		}

		block, err := self.blocks.Get(tipHsh)
		result = block

		return err

	})

	if err != nil {
		result = nil
	}

	return result, errors.Wrap(err, "BoltBlockStore:Tip")
}

func (self *BoltBlockStore) Cursor() (*core.BlockCursor, error) {
	var tip, err = self.Tip()
	if err != nil {
		return nil, errors.Wrap(err, "BoltBlockStore:Cursor")
	}

	var cursor = new(BoltBlockCursor)
	cursor.store = self
	cursor.current = tip

	var result core.BlockCursor = cursor
	return &result, nil
}

func (self *BoltBlockCursor) Current() *core.Block {
	if self.store == nil {
		return nil
	}

	return self.current
}

func (self *BoltBlockCursor) Reset() error {
	var tip, err = self.store.Tip()
	if err != nil {
		return errors.Wrap(err, "BoltBlockStore:Reset")
	}

	self.current = tip
	return nil
}

func (self *BoltBlockCursor) Next() (bool, error) {
	var tip, err = self.store.prev(self.current)
	if err != nil {
		return false, errors.Wrap(err, "BoltBlockStore:Next")
	}
	if tip == nil {
		return false, nil
	}

	self.current = tip
	return true, nil
}

func (self *BoltBlockStore) prev(current *core.Block) (*core.Block, error) {
	if current == nil || current.PrevHash == nil || len(current.PrevHash) <= 0 {
		return nil, nil
	}
	var result *core.Block = nil
	var err = self.db.GetDb().View(func(tx *bolt.Tx) error {
		block, err := self.blocks.Get(current.PrevHash)
		result = block

		return err
	})

	return result, errors.Wrap(err, "BoltBlockStore:prev")
}
