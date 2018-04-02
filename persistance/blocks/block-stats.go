package blocks

import (
	"../../core/errors"
	"github.com/boltdb/bolt"
	"os"
	"path/filepath"
)

const (
	bucket_blocks = "blocks"
	key_last      = "last"
	key_genesis   = "genesis"
)

type BoltBlockStats struct {
	db *bolt.DB
}

func NewBoltBlockStats(dbFile string) (*BoltBlockStats, error) {
	dbFile, err := filepath.Abs(dbFile)
	if err != nil {
		return nil, errors.Wrap(err, "NewBoltBlockStats")
	}

	dir := filepath.Dir(dbFile)
	os.MkdirAll(dir, 0700)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewBoltBlockStats:OpenDb")
	}

	result := BoltBlockStats{db: db}
	return &result, nil
}

func (self *BoltBlockStats) PutTip(hash []byte) error {
	err := self.put(bucket_blocks, key_last, hash)
	return err
}
func (self *BoltBlockStats) PutGenesis(hash []byte) error {
	err := self.put(bucket_blocks, key_genesis, hash)
	return err
}
func (self *BoltBlockStats) GetTip() ([]byte, error) {
	result, err := self.get(bucket_blocks, key_last)
	return result, err
}
func (self *BoltBlockStats) GetGenesis() ([]byte, error) {
	result, err := self.get(bucket_blocks, key_genesis)
	return result, err
}

func (self *BoltBlockStats) get(bucket string, key string) ([]byte, error) {
	var result []byte = nil

	var err = self.db.Update(func(tx *bolt.Tx) error {
		hsh, err := self.getFromTx(tx, bucket, key)
		result = hsh
		return err
	})

	return result, err
}

func (self *BoltBlockStats) getFromTx(tx *bolt.Tx, bucket string, key string) ([]byte, error) {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, errors.NewErrorFromString("unknown bucket %s", bucket)
	}

	var result = b.Get([]byte(key))
	return result, nil
}

func (self *BoltBlockStats) put(bucket string, key string, data []byte) error {
	var err = self.db.Update(func(tx *bolt.Tx) error {
		err := self.putToTx(tx, bucket, key, data)
		return err
	})

	return err
}
func (self *BoltBlockStats) putToTx(tx *bolt.Tx, bucket string, key string, data []byte) error {
	b := tx.Bucket([]byte(bucket))
	err := b.Put([]byte(key), data)
	return err
}
