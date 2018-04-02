package hugdb

import (
	"../serialization"
	"fmt"
	"github.com/boltdb/bolt"
)

var bucket_blocks = []byte("blocks")
var bucket_blocks_stats = []byte("blocks_stats")
var bucket_assets = []byte("asets")
var key_blocks_stats_last = []byte("l")
var key_blocks_stats_genesis = []byte("g")
var all_buckets = [][]byte{bucket_blocks, bucket_blocks_stats, bucket_assets}

type BoltDb struct {
	db *bolt.DB
}

func NewBoltDB(filePath string) (*BoltDb, error) {
	var db, err = bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, err
	}

	var result = BoltDb{db: db}
	return &result, nil
}

func (self *BoltDb) BucketBlockStats() []byte {
	return bucket_blocks_stats
}
func (self *BoltDb) BucketBlocks() []byte {
	return bucket_blocks
}
func (self *BoltDb) BucketAssets() []byte {
	return bucket_assets
}
func (self *BoltDb) KeyBlocksStatsLast() []byte {
	return key_blocks_stats_last
}
func (self *BoltDb) KeyBlocksStatsGenesis() []byte {
	return key_blocks_stats_genesis
}

func (self *BoltDb) Get(tx *bolt.Tx, bucket []byte, key []byte) []byte {
	var b = tx.Bucket(bucket)
	if b == nil {
		panic(fmt.Sprintf("unknown bucket %s", string(bucket)))
	}

	var result = b.Get(key)
	return result
}
func (self *BoltDb) GetDecode(tx *bolt.Tx, bucket []byte, key []byte, result interface{}) (bool, error) {
	var bin = self.Get(tx, bucket, key)
	if bin == nil || len(bin) == 0 {
		return false, nil
	}

	var err = serialization.ObjDecode(bin, result)

	return true, err
}

func (self *BoltDb) Set(tx *bolt.Tx, bucket []byte, key []byte, data []byte) error {
	b := tx.Bucket(bucket)
	err := b.Put(key, data)
	return err
}

func (self *BoltDb) SetEncoded(tx *bolt.Tx, bucket []byte, key []byte, data interface{}) error {
	bin, err := serialization.ObjEncode(data)
	if err != nil {
		return err
	}

	err = self.Set(tx, bucket, key, bin)

	return err
}

func (self *BoltDb) GetDb() *bolt.DB {
	return self.db
}

func (self *BoltDb) Bootstrap() error {
	var err = self.createBuckets(all_buckets...)

	return err
}

func (self *BoltDb) createBuckets(names ...[]byte) error {
	var err = self.db.Update(func(tx *bolt.Tx) error {
		for _, n := range names {
			var _, cErr = tx.CreateBucketIfNotExists(n)

			if cErr != nil {
				return cErr
			}
		}

		return nil
	})

	return err
}
