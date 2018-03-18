package storage

import (
	"../../common/serialization"
	"../../core"
	"github.com/boltdb/bolt"
	"errors"
)

var bucket_blocks = []byte("blocks")

var key_last = []byte("l")

type BoltBlockStore struct {
	db *bolt.DB
}
type BoltBlockCursor struct {
	store   *BoltBlockStore
	current *core.Block
}

func NewBoltBlockStore(filePath string) (*BoltBlockStore, error) {
	var result = new(BoltBlockStore)
	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		var _, createErr = tx.CreateBucketIfNotExists(bucket_blocks)
		return createErr
	})

	if err != nil {
		return nil, err
	}

	result.db = db

	return result, nil
}

func (self *BoltBlockStore) Close(){
	self.db.Close()
	self.db = nil
}

func (self *BoltBlockStore) Add(block *core.Block) error {
	if self.db == nil{
		return errors.New("db connection was closed")
	}
	var encodedBlock, err = serialization.SimpleEncode(block)
	if err != nil {
		return err
	}
	err = self.db.Update(func(tx *bolt.Tx) error {
		var blocks = tx.Bucket(bucket_blocks)
		err = blocks.Put(block.Hash, encodedBlock.Bytes())
		if err != nil {
			return err
		}

		err = blocks.Put(key_last, block.Hash)

		return err
	})

	return err
}

func (self *BoltBlockStore) Tip() (*core.Block, error) {
	if self.db == nil{
		return nil, errors.New("db connection was closed")
	}

	var result *core.Block = nil
	var err = self.db.View(func(tx *bolt.Tx) error {
		var blocks = tx.Bucket(bucket_blocks)

		var hash = blocks.Get(key_last)
		if hash == nil {
			return nil
		}

		var lastBin = blocks.Get(hash)
		if lastBin == nil {
			return nil
		}

		var last = new(core.Block)
		var err = serialization.SimpleDecode(lastBin, &last)
		result = last

		return err
	})

	if err != nil {
		result = nil
	}

	return result, err
}

func (self *BoltBlockStore) Cursor() (*core.BlockCursor, error) {
	if self.db == nil{
		return nil, errors.New("db connection was closed")
	}

	var tip, err = self.Tip()
	if err != nil {
		return nil, err
	}
	if tip == nil{
		return nil, errors.New("genesis block not found")
	}

	var cursor = new(BoltBlockCursor)
	cursor.store = self
	cursor.current = tip

	var result core.BlockCursor = cursor
	return &result, nil
}

func (self *BoltBlockCursor) Current() *core.Block {
	if self.store.db == nil{
		return nil
	}

	return self.current
}

func (self *BoltBlockCursor) Reset() error {
	if self.store.db == nil{
		return errors.New("db connection was closed")
	}

	var tip, err = self.store.Tip()
	if err != nil {
		return err
	}

	self.current = tip
	return nil
}

func (self *BoltBlockCursor) Next() (bool, error) {
	if self.store.db == nil{
		return false, errors.New("db connection was closed")
	}

	var tip, err = self.store.prev(self.current)
	if err != nil {
		return false, err
	}
	if tip == nil {
		return false, nil
	}

	self.current = tip
	return true, nil
}

func (self *BoltBlockStore) prev(current *core.Block) (*core.Block, error) {

	if current.PrevHash == nil || len(current.PrevHash) <= 0 {
		return nil, nil
	}
	var prev = new(core.Block)
	var err = self.db.View(func(tx *bolt.Tx) error {
		var blocks = tx.Bucket(bucket_blocks)
		var bin = blocks.Get(current.PrevHash)
		if bin == nil {
			return nil
		}

		var err = serialization.SimpleDecode(bin, &prev)

		return err
	})

	if err != nil {
		prev = nil
	}

	return prev, err
}
