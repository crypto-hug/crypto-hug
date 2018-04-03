package blocks

import (
	"../boltdb"
)

const (
	bucket_blocks = "blocks"
	key_last      = "last"
	key_genesis   = "genesis"
)

type BoltBlockStats struct {
	db *boltdb.BoltDB
}

func NewBoltBlockStats(dbFile string) (*BoltBlockStats, error) {
	db, err := boltdb.NewDB(dbFile)
	if err != nil {
		return nil, err
	}

	err = db.CreateBuckets(bucket_blocks)

	result := BoltBlockStats{db: db}
	return &result, nil
}

func (self *BoltBlockStats) PutTip(hash []byte) error {
	err := self.db.Put(bucket_blocks, key_last, hash)
	return err
}
func (self *BoltBlockStats) PutGenesis(hash []byte) error {
	err := self.db.Put(bucket_blocks, key_genesis, hash)
	return err
}
func (self *BoltBlockStats) GetTip() ([]byte, error) {
	result, err := self.db.Get(bucket_blocks, key_last)
	return result, err
}
func (self *BoltBlockStats) GetGenesis() ([]byte, error) {
	result, err := self.db.Get(bucket_blocks, key_genesis)
	return result, err
}
